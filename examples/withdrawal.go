package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MartianPay/go-binance"
	"github.com/MartianPay/go-binance/models"
)

func main() {
	// 从环境变量读取密钥
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		log.Fatal("Please set BINANCE_API_KEY and BINANCE_SECRET_KEY environment variables")
	}

	client := binance.NewClient(apiKey, secretKey)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("💸 Binance Small Amount Withdrawal")
	fmt.Println("==================================\n")

	// 1. 启用快速提现
	fmt.Println("📌 Enabling fast withdraw...")
	err := client.Account.EnableFastWithdrawSwitch(0)
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		fmt.Println("✅ Fast withdraw enabled")
	}

	// 2. 显示提现限额
	quota, err := client.Withdrawal.GetWithdrawalQuota()
	if err == nil {
		fmt.Printf("\n💰 24h Quota: $%s (Used: $%s)\n", quota.WdQuota, quota.UsedWdQuota)
	}

	// 3. 获取并显示有余额的资产
	fmt.Println("\n📊 Your assets:")
	assets, err := client.Account.GetUserAsset(models.UserAssetRequest{})
	if err != nil {
		log.Fatal("Error getting assets:", err)
	}

	availableCoins := []string{}
	coinBalances := make(map[string]string)

	for _, asset := range assets {
		if asset.Free != "0" {
			availableCoins = append(availableCoins, asset.Asset)
			coinBalances[asset.Asset] = asset.Free
			fmt.Printf("%d. %s: %s\n", len(availableCoins), asset.Asset, asset.Free)
		}
	}

	if len(availableCoins) == 0 {
		log.Fatal("No assets available")
	}

	// 4. 选择币种
	fmt.Print("\n✏️  Select coin number: ")
	input, _ := reader.ReadString('\n')
	coinIndex, _ := strconv.Atoi(strings.TrimSpace(input))

	if coinIndex < 1 || coinIndex > len(availableCoins) {
		log.Fatal("Invalid selection")
	}

	selectedCoin := availableCoins[coinIndex-1]
	balance := coinBalances[selectedCoin]
	fmt.Printf("✅ Selected: %s (Balance: %s)\n", selectedCoin, balance)

	// 5. 获取该币种的网络
	fmt.Printf("\n🔍 Loading %s networks...\n", selectedCoin)
	coins, err := client.Account.GetAllCoins()
	if err != nil {
		log.Fatal("Error getting coin info:", err)
	}

	var networks []models.NetworkInfo
	for _, coin := range coins {
		if coin.Coin == selectedCoin {
			for _, net := range coin.NetworkList {
				if net.WithdrawEnable {
					networks = append(networks, net)
				}
			}
			break
		}
	}

	if len(networks) == 0 {
		log.Fatal("No networks available for withdrawal")
	}

	// 显示网络选项
	fmt.Println("\n📡 Available networks:")
	for i, net := range networks {
		fmt.Printf("%d. %s", i+1, net.Network)
		if net.Name != "" {
			fmt.Printf(" (%s)", net.Name)
		}
		if net.IsDefault {
			fmt.Print(" [DEFAULT]")
		}
		fmt.Printf("\n   Min: %s, Fee: %s\n", net.WithdrawMin, net.WithdrawFee)
	}

	// 6. 选择网络
	fmt.Print("\n✏️  Select network number: ")
	input, _ = reader.ReadString('\n')
	netIndex, _ := strconv.Atoi(strings.TrimSpace(input))

	if netIndex < 1 || netIndex > len(networks) {
		log.Fatal("Invalid selection")
	}

	selectedNetwork := networks[netIndex-1]
	fmt.Printf("✅ Selected: %s\n", selectedNetwork.Network)

	// 7. 输入地址
	fmt.Print("\n📬 Enter withdrawal address: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	if address == "" {
		log.Fatal("Address cannot be empty")
	}

	// 8. 输入金额（建议小额）
	fmt.Printf("\n💵 Enter amount (Min: %s, Balance: %s): ", selectedNetwork.WithdrawMin, balance)
	input, _ = reader.ReadString('\n')
	amount := strings.TrimSpace(input)

	// 验证金额
	amountFloat, _ := strconv.ParseFloat(amount, 64)
	minFloat, _ := strconv.ParseFloat(selectedNetwork.WithdrawMin, 64)
	if amountFloat < minFloat {
		log.Fatalf("Amount too small. Minimum: %s", selectedNetwork.WithdrawMin)
	}

	// 9. 生成订单ID
	withdrawOrderId := fmt.Sprintf("WD_%d", time.Now().UnixNano())

	// 10. 确认信息
	fmt.Println("\n📋 Withdrawal Summary:")
	fmt.Println("=====================")
	fmt.Printf("Coin:     %s\n", selectedCoin)
	fmt.Printf("Network:  %s\n", selectedNetwork.Network)
	fmt.Printf("Address:  %s\n", address)
	fmt.Printf("Amount:   %s\n", amount)
	fmt.Printf("Fee:      %s %s\n", selectedNetwork.WithdrawFee, selectedCoin)
	fmt.Printf("Order ID: %s\n", withdrawOrderId)

	fmt.Print("\n⚠️  Confirm? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(confirm)) != "yes" {
		fmt.Println("❌ Cancelled")
		return
	}

	// 11. 执行提现
	fmt.Println("\n🚀 Processing withdrawal...")
	withdrawal, err := client.Withdrawal.Withdraw(models.WithdrawalRequest{
		Coin:               selectedCoin,
		Network:            selectedNetwork.Network,
		Address:            address,
		Amount:             amount,
		WithdrawOrderId:    withdrawOrderId,
		TransactionFeeFlag: true,
	})

	if err != nil {
		log.Fatalf("❌ Failed: %v", err)
	}

	fmt.Printf("✅ Success! Withdrawal ID: %s\n", withdrawal.Id)

	// 12. 监控状态（最多30秒）
	fmt.Println("\n⏳ Checking status...")

	for i := 0; i < 6; i++ {
		time.Sleep(5 * time.Second)

		history, err := client.Withdrawal.GetWithdrawalHistory(models.WithdrawalHistoryRequest{
			WithdrawOrderId: withdrawOrderId,
			StartTime:       time.Now().Add(-10 * time.Minute),
			EndTime:         time.Now(),
		})

		if err == nil && len(history) > 0 {
			w := history[0]
			status := getStatusText(w.Status)
			fmt.Printf("[%s] %s", time.Now().Format("15:04:05"), status)

			if w.TxId != "" {
				fmt.Printf(" - TxID: %s", w.TxId)
			} else if w.Status == 6 {
				fmt.Print(" - Internal transfer")
			}
			fmt.Println()

			// 检查是否完成或失败
			if w.Status == 6 {
				fmt.Println("\n✅ Withdrawal completed!")
				break
			} else if w.Status == 1 || w.Status == 3 || w.Status == 5 {
				fmt.Printf("\n❌ Withdrawal failed: %s\n", status)
				break
			}
		}
	}

	fmt.Println("\n✅ Done!")
}

func getStatusText(status int) string {
	switch status {
	case 0:
		return "Email Sent"
	case 1:
		return "Cancelled"
	case 2:
		return "Awaiting Approval"
	case 3:
		return "Rejected"
	case 4:
		return "Processing"
	case 5:
		return "Failed"
	case 6:
		return "Completed"
	default:
		return fmt.Sprintf("Unknown(%d)", status)
	}
}
