package main

import (
	"fmt"
	"log"
	"os"

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
	
	fmt.Println("🏦 Binance Account Information")
	fmt.Println("==============================\n")
	
	// 1. 获取账户基本信息
	fmt.Println("📊 Account Status:")
	accountInfo, err := client.Account.GetAccountInfo()
	if err != nil {
		log.Printf("Error getting account info: %v", err)
	} else {
		fmt.Printf("VIP Level: %d\n", accountInfo.VipLevel)
		fmt.Printf("Margin Trading: %v\n", accountInfo.IsMarginEnable)
		fmt.Printf("Futures Trading: %v\n", accountInfo.IsFutureEnable)
	}
	
	// 2. 获取24小时提现限额
	fmt.Println("\n💸 24-Hour Withdrawal Quota:")
	quota, err := client.Withdrawal.GetWithdrawalQuota()
	if err != nil {
		log.Printf("Error getting withdrawal quota: %v", err)
	} else {
		fmt.Printf("Total Limit: $%s USD\n", quota.WdQuota)
		fmt.Printf("Used Amount: $%s USD\n", quota.UsedWdQuota)
		
		// 计算剩余额度
		var totalQuota, usedQuota float64
		fmt.Sscanf(quota.WdQuota, "%f", &totalQuota)
		fmt.Sscanf(quota.UsedWdQuota, "%f", &usedQuota)
		remaining := totalQuota - usedQuota
		fmt.Printf("Remaining:   $%.2f USD\n", remaining)
	}
	
	// 3. 获取所有资产余额
	fmt.Println("\n💰 Asset Balances:")
	assets, err := client.Account.GetUserAsset(models.UserAssetRequest{})
	if err != nil {
		log.Fatal("Error getting user assets:", err)
	}
	
	hasAssets := false
	for _, asset := range assets {
		// 只显示有余额的资产
		if asset.Free != "0" || asset.Locked != "0" || asset.Freeze != "0" || asset.Withdrawing != "0" {
			if !hasAssets {
				fmt.Println("\nCoin    | Free Balance    | Locked         | Freeze         | Withdrawing")
				fmt.Println("--------|-----------------|----------------|----------------|----------------")
				hasAssets = true
			}
			
			fmt.Printf("%-7s | %-15s | %-14s | %-14s | %-14s\n",
				asset.Asset,
				formatAmount(asset.Free),
				formatAmount(asset.Locked),
				formatAmount(asset.Freeze),
				formatAmount(asset.Withdrawing))
		}
	}
	
	if !hasAssets {
		fmt.Println("No assets found")
	}
	
	// 4. 获取保存的提现地址
	fmt.Println("\n📋 Saved Withdrawal Addresses:")
	addresses, err := client.Withdrawal.GetWithdrawalAddressList()
	if err != nil {
		log.Printf("Error getting withdrawal addresses: %v", err)
	} else {
		if len(addresses) == 0 {
			fmt.Println("No saved withdrawal addresses")
		} else {
			for _, addr := range addresses {
				fmt.Printf("\n• %s (%s - %s)\n", addr.Name, addr.Coin, addr.Network)
				fmt.Printf("  Address: %s\n", addr.Address)
				fmt.Printf("  Whitelisted: %v\n", addr.WhiteStatus)
			}
		}
	}
	
	fmt.Println("\n✅ Done!")
}

func formatAmount(amount string) string {
	if amount == "0" {
		return "-"
	}
	// 截断过长的小数
	if len(amount) > 15 {
		return amount[:15]
	}
	return amount
}