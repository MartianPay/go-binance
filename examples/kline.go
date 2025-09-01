package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MartianPay/go-binance"
	"github.com/MartianPay/go-binance/models"
)

func main() {
	// 初始化客户端（K线数据不需要API密钥）
	client := binance.NewClient("", "")
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📊 Binance Kline/Candlestick Data Fetcher")
	fmt.Println("==========================================\n")

	for {
		// 1. 输入交易对
		fmt.Print("📌 Enter trading pair (e.g., BTCUSDT, ETHUSDT): ")
		symbol, _ := reader.ReadString('\n')
		symbol = strings.TrimSpace(strings.ToUpper(symbol))

		if symbol == "" {
			symbol = "BTCUSDT" // 默认值
			fmt.Printf("   Using default: %s\n", symbol)
		}

		// 2. 选择时间间隔
		fmt.Println("\n⏱️  Select interval:")
		intervals := []models.KlineInterval{
			models.Interval1m,
			models.Interval3m,
			models.Interval5m,
			models.Interval15m,
			models.Interval30m,
			models.Interval1h,
			models.Interval2h,
			models.Interval4h,
			models.Interval6h,
			models.Interval8h,
			models.Interval12h,
			models.Interval1d,
			models.Interval3d,
			models.Interval1w,
			models.Interval1M,
		}

		for i, interval := range intervals {
			fmt.Printf("%2d. %s", i+1, interval)
			if (i+1)%5 == 0 {
				fmt.Println()
			} else {
				fmt.Print("    ")
			}
		}
		fmt.Println()

		fmt.Print("\n✏️  Select interval number (default 6 for 1h): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		intervalIndex := 6 // 默认1小时
		if input != "" {
			if idx, err := strconv.Atoi(input); err == nil && idx >= 1 && idx <= len(intervals) {
				intervalIndex = idx
			}
		}
		selectedInterval := intervals[intervalIndex-1]
		fmt.Printf("   Selected: %s\n", selectedInterval)

		// 3. 选择获取方式
		fmt.Println("\n📅 Select data range:")
		fmt.Println("1. Get recent N candles")
		fmt.Println("2. Get data by date range")
		fmt.Println("3. Get UI optimized klines")
		fmt.Print("\n✏️  Select option (1-3, default 1): ")
		
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)
		if option == "" {
			option = "1"
		}

		var klines []models.Kline
		var err error

		switch option {
		case "1":
			// 获取最近N根K线
			fmt.Print("\n📊 How many candles? (1-1500, default 20): ")
			input, _ = reader.ReadString('\n')
			limit := 20
			if input = strings.TrimSpace(input); input != "" {
				if l, err := strconv.Atoi(input); err == nil && l > 0 && l <= 1500 {
					limit = l
				}
			}

			req := models.KlineRequest{
				Symbol:   symbol,
				Interval: selectedInterval,
				Limit:    limit,
			}

			fmt.Printf("\n⏳ Fetching %d %s klines for %s...\n", limit, selectedInterval, symbol)
			klines, err = client.Market.GetKlines(req)

		case "2":
			// 按日期范围获取
			fmt.Println("\n📅 Enter date range (leave empty for defaults):")
			
			fmt.Print("   Start date (YYYY-MM-DD, default 7 days ago): ")
			startInput, _ := reader.ReadString('\n')
			startInput = strings.TrimSpace(startInput)
			
			var startTime time.Time
			if startInput == "" {
				startTime = time.Now().AddDate(0, 0, -7)
			} else {
				startTime, _ = time.Parse("2006-01-02", startInput)
			}

			fmt.Print("   End date (YYYY-MM-DD, default today): ")
			endInput, _ := reader.ReadString('\n')
			endInput = strings.TrimSpace(endInput)
			
			var endTime time.Time
			if endInput == "" {
				endTime = time.Now()
			} else {
				endTime, _ = time.Parse("2006-01-02", endInput)
				endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second) // 包含整天
			}

			req := models.KlineRequest{
				Symbol:    symbol,
				Interval:  selectedInterval,
				StartTime: startTime.Unix() * 1000,
				EndTime:   endTime.Unix() * 1000,
			}

			fmt.Printf("\n⏳ Fetching %s klines from %s to %s...\n", 
				selectedInterval, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
			klines, err = client.Market.GetKlines(req)

		case "3":
			// UI优化的K线
			fmt.Print("\n📊 How many UI klines? (1-1500, default 50): ")
			input, _ = reader.ReadString('\n')
			limit := 50
			if input = strings.TrimSpace(input); input != "" {
				if l, err := strconv.Atoi(input); err == nil && l > 0 && l <= 1500 {
					limit = l
				}
			}

			req := models.KlineRequest{
				Symbol:   symbol,
				Interval: selectedInterval,
				Limit:    limit,
			}

			fmt.Printf("\n⏳ Fetching %d UI optimized %s klines for %s...\n", limit, selectedInterval, symbol)
			klines, err = client.Market.GetUIKlines(req)

		default:
			fmt.Println("❌ Invalid option")
			continue
		}

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			// 4. 显示结果
			displayKlines(klines, symbol, selectedInterval)
		}

		// 5. 询问是否继续
		fmt.Print("\n🔄 Continue? (yes/no, default yes): ")
		continueInput, _ := reader.ReadString('\n')
		continueInput = strings.TrimSpace(strings.ToLower(continueInput))
		if continueInput == "no" || continueInput == "n" {
			fmt.Println("\n👋 Goodbye!")
			break
		}
		fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	}
}

func displayKlines(klines []models.Kline, symbol string, interval models.KlineInterval) {
	if len(klines) == 0 {
		fmt.Println("❌ No data available")
		return
	}

	fmt.Printf("\n✅ Found %d klines for %s (%s)\n", len(klines), symbol, interval)
	fmt.Println(strings.Repeat("-", 100))

	// 选择显示格式
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n📈 Display format:")
	fmt.Println("1. Detailed view (show all data)")
	fmt.Println("2. Compact view (OHLCV only)")
	fmt.Println("3. Statistical summary")
	fmt.Print("\n✏️  Select format (1-3, default 2): ")
	
	formatInput, _ := reader.ReadString('\n')
	formatInput = strings.TrimSpace(formatInput)
	if formatInput == "" {
		formatInput = "2"
	}

	switch formatInput {
	case "1":
		// 详细视图
		displayCount := len(klines)
		if displayCount > 10 {
			displayCount = 10
			fmt.Printf("\n📊 Showing first %d klines (out of %d):\n\n", displayCount, len(klines))
		}

		for i := 0; i < displayCount; i++ {
			kline := klines[i]
			openTime := time.Unix(kline.OpenTime/1000, 0)
			closeTime := time.Unix(kline.CloseTime/1000, 0)
			
			fmt.Printf("Kline #%d:\n", i+1)
			fmt.Printf("  🕐 Open Time:  %s\n", openTime.Format("2006-01-02 15:04:05"))
			fmt.Printf("  🕑 Close Time: %s\n", closeTime.Format("2006-01-02 15:04:05"))
			fmt.Printf("  💹 OHLC: O=%s, H=%s, L=%s, C=%s\n", 
				kline.Open, kline.High, kline.Low, kline.Close)
			fmt.Printf("  📊 Volume: %s\n", kline.Volume)
			fmt.Printf("  💰 Quote Volume: %s\n", kline.QuoteAssetVolume)
			fmt.Printf("  🔢 Trades: %d\n", kline.NumberOfTrades)
			fmt.Printf("  📈 Taker Buy Base: %s\n", kline.TakerBuyBaseAssetVolume)
			fmt.Printf("  💵 Taker Buy Quote: %s\n", kline.TakerBuyQuoteAssetVolume)
			
			// 计算涨跌
			change := kline.GetCloseFloat() - kline.GetOpenFloat()
			changePercent := 0.0
			if kline.GetOpenFloat() > 0 {
				changePercent = (change / kline.GetOpenFloat()) * 100
			}
			
			changeSymbol := "➡️"
			if change > 0 {
				changeSymbol = "📈"
			} else if change < 0 {
				changeSymbol = "📉"
			}
			
			fmt.Printf("  %s Change: %.4f (%.2f%%)\n", changeSymbol, change, changePercent)
			fmt.Println(strings.Repeat("-", 50))
		}

	case "2":
		// 紧凑视图
		fmt.Printf("\n%-20s %-10s %-10s %-10s %-10s %-15s %-8s\n", 
			"Time", "Open", "High", "Low", "Close", "Volume", "Change%")
		fmt.Println(strings.Repeat("-", 100))

		displayCount := len(klines)
		if displayCount > 30 {
			displayCount = 30
			fmt.Printf("(Showing first %d of %d klines)\n", displayCount, len(klines))
		}

		for i := 0; i < displayCount; i++ {
			kline := klines[i]
			openTime := time.Unix(kline.OpenTime/1000, 0)
			
			change := kline.GetCloseFloat() - kline.GetOpenFloat()
			changePercent := 0.0
			if kline.GetOpenFloat() > 0 {
				changePercent = (change / kline.GetOpenFloat()) * 100
			}

			changeStr := fmt.Sprintf("%.2f%%", changePercent)
			if change > 0 {
				changeStr = "+" + changeStr
			}

			fmt.Printf("%-20s %-10s %-10s %-10s %-10s %-15.2f %-8s\n",
				openTime.Format("2006-01-02 15:04"),
				kline.Open[:min(10, len(kline.Open))],
				kline.High[:min(10, len(kline.High))],
				kline.Low[:min(10, len(kline.Low))],
				kline.Close[:min(10, len(kline.Close))],
				kline.GetVolumeFloat(),
				changeStr)
		}

	case "3":
		// 统计摘要
		var totalVolume, avgPrice, maxHigh, minLow float64
		var totalTrades int
		positiveCount := 0
		negativeCount := 0

		maxHigh = klines[0].GetHighFloat()
		minLow = klines[0].GetLowFloat()

		for _, kline := range klines {
			totalVolume += kline.GetVolumeFloat()
			avgPrice += (kline.GetOpenFloat() + kline.GetCloseFloat()) / 2
			totalTrades += kline.NumberOfTrades

			if kline.GetHighFloat() > maxHigh {
				maxHigh = kline.GetHighFloat()
			}
			if kline.GetLowFloat() < minLow {
				minLow = kline.GetLowFloat()
			}

			if kline.GetCloseFloat() > kline.GetOpenFloat() {
				positiveCount++
			} else if kline.GetCloseFloat() < kline.GetOpenFloat() {
				negativeCount++
			}
		}

		avgPrice /= float64(len(klines))

		firstKline := klines[0]
		lastKline := klines[len(klines)-1]
		startTime := time.Unix(firstKline.OpenTime/1000, 0)
		endTime := time.Unix(lastKline.CloseTime/1000, 0)

		overallChange := lastKline.GetCloseFloat() - firstKline.GetOpenFloat()
		overallChangePercent := 0.0
		if firstKline.GetOpenFloat() > 0 {
			overallChangePercent = (overallChange / firstKline.GetOpenFloat()) * 100
		}

		fmt.Println("\n📊 Statistical Summary")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("📅 Period: %s to %s\n", 
			startTime.Format("2006-01-02 15:04"), 
			endTime.Format("2006-01-02 15:04"))
		fmt.Printf("📈 Total Klines: %d\n", len(klines))
		fmt.Printf("💹 First Open: %.4f\n", firstKline.GetOpenFloat())
		fmt.Printf("💹 Last Close: %.4f\n", lastKline.GetCloseFloat())
		fmt.Printf("📊 Overall Change: %.4f (%.2f%%)\n", overallChange, overallChangePercent)
		fmt.Printf("⬆️  Highest: %.4f\n", maxHigh)
		fmt.Printf("⬇️  Lowest: %.4f\n", minLow)
		fmt.Printf("📊 Average Price: %.4f\n", avgPrice)
		fmt.Printf("📈 Positive Candles: %d (%.1f%%)\n", 
			positiveCount, float64(positiveCount)/float64(len(klines))*100)
		fmt.Printf("📉 Negative Candles: %d (%.1f%%)\n", 
			negativeCount, float64(negativeCount)/float64(len(klines))*100)
		fmt.Printf("💰 Total Volume: %.2f\n", totalVolume)
		fmt.Printf("🔢 Total Trades: %d\n", totalTrades)
		fmt.Printf("📊 Avg Volume per Candle: %.2f\n", totalVolume/float64(len(klines)))
		fmt.Printf("🔢 Avg Trades per Candle: %d\n", totalTrades/len(klines))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}