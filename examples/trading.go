package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		fmt.Println("⚠️  WARNING: API keys not found in environment variables")
		fmt.Println("   Trading operations require BINANCE_API_KEY and BINANCE_SECRET_KEY")
		fmt.Println("   You can still use read-only features\n")
	}

	client := binance.NewClient(apiKey, secretKey)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📈 Binance Trading Interface")
	fmt.Println("============================\n")

	for {
		showMainMenu()
		
		fmt.Print("\n✏️  Select option: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			viewAccountInfo(client)
		case "2":
			viewOpenOrders(client, reader)
		case "3":
			createOrder(client, reader)
		case "4":
			queryOrderStatus(client, reader)
		case "5":
			cancelOrder(client, reader)
		case "6":
			viewOrderHistory(client, reader)
		case "7":
			viewMyTrades(client, reader)
		case "8":
			testOrder(client, reader)
		case "9":
			viewTradingRules(client, reader)
		case "10":
			fmt.Println("\n👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid option")
		}

		fmt.Print("\n🔄 Press Enter to continue...")
		reader.ReadString('\n')
		fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	}
}

func showMainMenu() {
	fmt.Println("📋 Main Menu:")
	fmt.Println("1. View Account Info & Balances")
	fmt.Println("2. View Open Orders")
	fmt.Println("3. Create New Order")
	fmt.Println("4. Query Order Status")
	fmt.Println("5. Cancel Order")
	fmt.Println("6. View Order History")
	fmt.Println("7. View My Trades (Account Trade List)")
	fmt.Println("8. Test Order (Simulation)")
	fmt.Println("9. View Trading Rules (LOT_SIZE, NOTIONAL, etc.)")
	fmt.Println("10. Exit")
}

func viewAccountInfo(client *binance.BinanceClient) {
	fmt.Println("\n💰 Account Information")
	fmt.Println(strings.Repeat("-", 50))

	info, err := client.Trading.GetAccountInfo(0)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("\n✅ Account info retrieved!")
	
	// Print raw JSON response
	jsonBytes, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
	} else {
		fmt.Println("\n📊 Full Response (JSON):")
		fmt.Println(string(jsonBytes))
	}
}

func viewOpenOrders(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n📂 View Open Orders")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Print("Enter symbol (e.g., BTCUSDT, or press Enter for all): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	req := models.OpenOrdersRequest{
		Symbol: symbol,
	}

	orders, err := client.Trading.GetOpenOrders(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("📭 No open orders found")
		return
	}

	fmt.Printf("\n✅ Found %d open orders!\n", len(orders))
	
	// Print raw JSON response
	jsonBytes, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
	} else {
		fmt.Println("\n📊 Full Response (JSON):")
		fmt.Println(string(jsonBytes))
	}
}

func createOrder(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n➕ Create New Order")
	fmt.Println(strings.Repeat("-", 50))

	// 1. Symbol
	fmt.Print("Enter symbol (e.g., BTCUSDT): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	// 2. Side
	fmt.Println("\nSelect side:")
	fmt.Println("1. BUY")
	fmt.Println("2. SELL")
	fmt.Print("✏️  Select (1-2): ")
	sideInput, _ := reader.ReadString('\n')
	sideInput = strings.TrimSpace(sideInput)
	
	var side models.OrderSide
	if sideInput == "1" {
		side = models.SideBuy
	} else {
		side = models.SideSell
	}

	// 3. Order Type
	fmt.Println("\nSelect order type:")
	fmt.Println("1. LIMIT (specify price)")
	fmt.Println("2. MARKET (market price)")
	fmt.Println("3. STOP_LOSS_LIMIT")
	fmt.Println("4. TAKE_PROFIT_LIMIT")
	fmt.Print("✏️  Select (1-4, default 1): ")
	typeInput, _ := reader.ReadString('\n')
	typeInput = strings.TrimSpace(typeInput)
	
	var orderType models.OrderType
	var needPrice bool
	var needStopPrice bool
	
	switch typeInput {
	case "2":
		orderType = models.OrderTypeMarket
	case "3":
		orderType = models.OrderTypeStopLossLimit
		needPrice = true
		needStopPrice = true
	case "4":
		orderType = models.OrderTypeTakeProfitLimit
		needPrice = true
		needStopPrice = true
	default:
		orderType = models.OrderTypeLimit
		needPrice = true
	}

	// 4. Quantity or Quote Order Quantity
	req := models.NewOrderRequest{
		Symbol:           symbol,
		Side:             side,
		Type:             orderType,
		NewOrderRespType: models.OrderResponseTypeFULL,
	}

	if orderType == models.OrderTypeMarket {
		// For MARKET orders, allow choice between quantity and quoteOrderQty
		fmt.Println("\nFor MARKET orders, choose quantity type:")
		fmt.Println("1. Specify base asset quantity (e.g., amount of BTC)")
		fmt.Println("2. Specify quote asset amount (e.g., amount of USDT to spend/receive)")
		fmt.Print("✏️  Select (1-2, default 1): ")
		qtyTypeInput, _ := reader.ReadString('\n')
		qtyTypeInput = strings.TrimSpace(qtyTypeInput)
		
		if qtyTypeInput == "2" {
			// Use quoteOrderQty
			// Parse symbol to get quote asset (last part of pair, e.g., USDT in BTCUSDT)
			quoteAsset := "quote asset"
			if len(symbol) >= 4 {
				// Common quote assets
				if strings.HasSuffix(symbol, "USDT") {
					quoteAsset = "USDT"
				} else if strings.HasSuffix(symbol, "BUSD") {
					quoteAsset = "BUSD"
				} else if strings.HasSuffix(symbol, "BTC") {
					quoteAsset = "BTC"
				} else if strings.HasSuffix(symbol, "ETH") {
					quoteAsset = "ETH"
				} else if strings.HasSuffix(symbol, "BNB") {
					quoteAsset = "BNB"
				}
			}
			
			if side == models.SideBuy {
				fmt.Printf("\nEnter %s amount to spend (quoteOrderQty): ", quoteAsset)
			} else {
				fmt.Printf("\nEnter %s amount to receive (quoteOrderQty): ", quoteAsset)
			}
			quoteOrderQty, _ := reader.ReadString('\n')
			req.QuoteOrderQty = strings.TrimSpace(quoteOrderQty)
		} else {
			// Use regular quantity
			fmt.Print("\nEnter base asset quantity: ")
			quantity, _ := reader.ReadString('\n')
			req.Quantity = strings.TrimSpace(quantity)
		}
	} else {
		// For LIMIT and other orders, use regular quantity
		fmt.Print("\nEnter quantity: ")
		quantity, _ := reader.ReadString('\n')
		req.Quantity = strings.TrimSpace(quantity)
	}

	// 5. Price (for limit orders)
	if needPrice {
		fmt.Print("Enter price: ")
		price, _ := reader.ReadString('\n')
		req.Price = strings.TrimSpace(price)
		req.TimeInForce = models.TimeInForceGTC
	}

	// 6. Stop Price (for stop orders)
	if needStopPrice {
		fmt.Print("Enter stop price: ")
		stopPrice, _ := reader.ReadString('\n')
		req.StopPrice = strings.TrimSpace(stopPrice)
	}

	// 7. Client Order ID (optional)
	fmt.Print("\nEnter custom Client Order ID (optional, press Enter to skip): ")
	clientOrderId, _ := reader.ReadString('\n')
	clientOrderId = strings.TrimSpace(clientOrderId)
	if clientOrderId != "" {
		req.NewClientOrderId = clientOrderId
	}

	// 8. Confirm
	fmt.Println("\n📋 Order Summary:")
	fmt.Printf("  Symbol: %s\n", req.Symbol)
	fmt.Printf("  Side: %s\n", req.Side)
	fmt.Printf("  Type: %s\n", req.Type)
	if req.Quantity != "" {
		fmt.Printf("  Quantity: %s\n", req.Quantity)
	}
	if req.QuoteOrderQty != "" {
		fmt.Printf("  Quote Order Qty: %s\n", req.QuoteOrderQty)
	}
	if req.Price != "" {
		fmt.Printf("  Price: %s\n", req.Price)
	}
	if req.StopPrice != "" {
		fmt.Printf("  Stop Price: %s\n", req.StopPrice)
	}
	if req.NewClientOrderId != "" {
		fmt.Printf("  Client Order ID: %s\n", req.NewClientOrderId)
	}

	fmt.Print("\n⚠️  Confirm order? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(confirm)) != "yes" {
		fmt.Println("❌ Order cancelled")
		return
	}

	// 9. Execute
	fmt.Println("\n🚀 Placing order...")
	order, err := client.Trading.NewOrder(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("\n✅ Order placed successfully!")
	
	// Print raw JSON response
	jsonBytes, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
	} else {
		fmt.Println("\n📊 Full Response (JSON):")
		fmt.Println(string(jsonBytes))
	}
}

func queryOrderStatus(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n🔍 Query Order Status")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Print("Enter symbol (e.g., BTCUSDT): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	fmt.Print("Enter Order ID (or press Enter to use Client Order ID): ")
	orderIdInput, _ := reader.ReadString('\n')
	orderIdInput = strings.TrimSpace(orderIdInput)

	req := models.QueryOrderRequest{
		Symbol: symbol,
	}

	if orderIdInput != "" {
		if orderId, err := strconv.ParseInt(orderIdInput, 10, 64); err == nil {
			req.OrderId = orderId
		} else {
			fmt.Print("Enter Client Order ID: ")
			clientOrderId, _ := reader.ReadString('\n')
			req.OrigClientOrderId = strings.TrimSpace(clientOrderId)
		}
	} else {
		fmt.Print("Enter Client Order ID: ")
		clientOrderId, _ := reader.ReadString('\n')
		req.OrigClientOrderId = strings.TrimSpace(clientOrderId)
	}

	fmt.Println("\n⏳ Querying order...")
	order, err := client.Trading.QueryOrder(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("\n✅ Order found!")

	// Display execution details for filled orders
	if order.Status == models.OrderStatusFilled {
		fmt.Println("\n💰 Execution Summary:")
		fmt.Println(strings.Repeat("-", 40))

		// Parse symbol to identify base and quote assets
		baseAsset := "base asset"
		quoteAsset := "quote asset"
		if len(symbol) >= 4 {
			// Common quote assets
			if strings.HasSuffix(symbol, "USDT") {
				quoteAsset = "USDT"
				baseAsset = strings.TrimSuffix(symbol, "USDT")
			} else if strings.HasSuffix(symbol, "BUSD") {
				quoteAsset = "BUSD"
				baseAsset = strings.TrimSuffix(symbol, "BUSD")
			} else if strings.HasSuffix(symbol, "BTC") {
				quoteAsset = "BTC"
				baseAsset = strings.TrimSuffix(symbol, "BTC")
			} else if strings.HasSuffix(symbol, "ETH") {
				quoteAsset = "ETH"
				baseAsset = strings.TrimSuffix(symbol, "ETH")
			} else if strings.HasSuffix(symbol, "BNB") {
				quoteAsset = "BNB"
				baseAsset = strings.TrimSuffix(symbol, "BNB")
			}
		}

		if order.Side == models.SideBuy {
			fmt.Printf("  💵 Spent: %s %s\n", order.CummulativeQuoteQty, quoteAsset)
			fmt.Printf("  📦 Received: %s %s\n", order.ExecutedQty, baseAsset)

			// Calculate average price if both values are available
			if order.ExecutedQty != "" && order.ExecutedQty != "0" &&
			   order.CummulativeQuoteQty != "" && order.CummulativeQuoteQty != "0" {
				executedQty, _ := strconv.ParseFloat(order.ExecutedQty, 64)
				cummulativeQuoteQty, _ := strconv.ParseFloat(order.CummulativeQuoteQty, 64)
				if executedQty > 0 {
					avgPrice := cummulativeQuoteQty / executedQty
					fmt.Printf("  📊 Average Price: %.8f %s/%s\n", avgPrice, quoteAsset, baseAsset)
				}
			}
		} else {
			fmt.Printf("  📦 Sold: %s %s\n", order.ExecutedQty, baseAsset)
			fmt.Printf("  💵 Received: %s %s\n", order.CummulativeQuoteQty, quoteAsset)

			// Calculate average price if both values are available
			if order.ExecutedQty != "" && order.ExecutedQty != "0" &&
			   order.CummulativeQuoteQty != "" && order.CummulativeQuoteQty != "0" {
				executedQty, _ := strconv.ParseFloat(order.ExecutedQty, 64)
				cummulativeQuoteQty, _ := strconv.ParseFloat(order.CummulativeQuoteQty, 64)
				if executedQty > 0 {
					avgPrice := cummulativeQuoteQty / executedQty
					fmt.Printf("  📊 Average Price: %.8f %s/%s\n", avgPrice, quoteAsset, baseAsset)
				}
			}
		}
		fmt.Println(strings.Repeat("-", 40))
	}

	// Print raw JSON response
	jsonBytes, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
	} else {
		fmt.Println("\n📊 Full Response (JSON):")
		fmt.Println(string(jsonBytes))
	}

	// Monitor status if still active
	if order.Status == models.OrderStatusNew || order.Status == models.OrderStatusPartiallyFilled {
		fmt.Print("\n⏳ Order is still active. Monitor status? (yes/no): ")
		monitor, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(monitor)) == "yes" {
			monitorOrderStatus(client, req)
		}
	}
}

func monitorOrderStatus(client *binance.BinanceClient, req models.QueryOrderRequest) {
	fmt.Println("\n⏳ Monitoring order status (press Ctrl+C to stop)...")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 60; i++ { // Monitor for up to 5 minutes
		<-ticker.C
		
		order, err := client.Trading.QueryOrder(req)
		if err != nil {
			fmt.Printf("\n❌ Error checking status: %v\n", err)
			break
		}

		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("[%s] Status: %s, Executed: %s/%s\n", 
			timestamp, order.Status, order.ExecutedQty, order.OrigQty)

		if order.Status == models.OrderStatusFilled || 
		   order.Status == models.OrderStatusCanceled || 
		   order.Status == models.OrderStatusRejected {
			fmt.Printf("\n✅ Order completed with status: %s\n", order.Status)
			break
		}
	}
}

func cancelOrder(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n❌ Cancel Order")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Print("Enter symbol (e.g., BTCUSDT): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	fmt.Print("Enter Order ID to cancel: ")
	orderIdInput, _ := reader.ReadString('\n')
	orderIdInput = strings.TrimSpace(orderIdInput)

	orderId, err := strconv.ParseInt(orderIdInput, 10, 64)
	if err != nil {
		fmt.Printf("❌ Invalid order ID: %v\n", err)
		return
	}

	req := models.CancelOrderRequest{
		Symbol:  symbol,
		OrderId: orderId,
	}

	fmt.Print("\n⚠️  Confirm cancel order? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(confirm)) != "yes" {
		fmt.Println("❌ Cancellation aborted")
		return
	}

	fmt.Println("\n⏳ Cancelling order...")
	result, err := client.Trading.CancelOrder(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("\n✅ Order cancelled successfully!")
	
	// Print raw JSON response
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
	} else {
		fmt.Println("\n📊 Full Response (JSON):")
		fmt.Println(string(jsonBytes))
	}
}

func viewOrderHistory(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n📜 Order History")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Print("Enter symbol (e.g., BTCUSDT): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	fmt.Print("Number of days to look back (default 7): ")
	daysInput, _ := reader.ReadString('\n')
	daysInput = strings.TrimSpace(daysInput)
	
	days := 7
	if d, err := strconv.Atoi(daysInput); err == nil && d > 0 {
		days = d
	}

	fmt.Print("Limit (max 1000, default 50): ")
	limitInput, _ := reader.ReadString('\n')
	limitInput = strings.TrimSpace(limitInput)
	
	limit := 50
	if l, err := strconv.Atoi(limitInput); err == nil && l > 0 && l <= 1000 {
		limit = l
	}

	req := models.AllOrdersRequest{
		Symbol:    symbol,
		StartTime: time.Now().AddDate(0, 0, -days),
		EndTime:   time.Now(),
		Limit:     limit,
	}

	fmt.Printf("\n⏳ Fetching orders from last %d days...\n", days)
	orders, err := client.Trading.GetAllOrders(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("📭 No orders found in the specified period")
		return
	}

	fmt.Printf("\n✅ Found %d orders!\n", len(orders))
	
	// Print raw JSON response
	jsonBytes, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
	} else {
		fmt.Println("\n📊 Full Response (JSON):")
		fmt.Println(string(jsonBytes))
	}
}

func viewMyTrades(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n💱 My Trades (Account Trade List)")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Print("Enter symbol (e.g., BTCUSDT): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	// Ask if user wants to query by Order ID
	fmt.Print("Query by specific Order ID? (yes/no, default no): ")
	queryByOrder, _ := reader.ReadString('\n')
	queryByOrder = strings.TrimSpace(strings.ToLower(queryByOrder))

	req := models.MyTradesRequest{
		Symbol: symbol,
	}

	if queryByOrder == "yes" || queryByOrder == "y" {
		// Query by Order ID
		fmt.Print("Enter Order ID: ")
		orderIdInput, _ := reader.ReadString('\n')
		orderIdInput = strings.TrimSpace(orderIdInput)

		if orderId, err := strconv.ParseInt(orderIdInput, 10, 64); err == nil {
			req.OrderId = orderId
			fmt.Printf("\n⏳ Fetching trades for Order ID %d...\n", orderId)
		} else {
			fmt.Printf("❌ Invalid Order ID: %v\n", err)
			return
		}
	} else {
		// Query by time range
		fmt.Print("Number of days to look back (default 7): ")
		daysInput, _ := reader.ReadString('\n')
		daysInput = strings.TrimSpace(daysInput)

		days := 7
		if d, err := strconv.Atoi(daysInput); err == nil && d > 0 {
			days = d
		}

		fmt.Print("From Trade ID (optional, press Enter to skip): ")
		fromIdInput, _ := reader.ReadString('\n')
		fromIdInput = strings.TrimSpace(fromIdInput)

		if fromIdInput != "" {
			if fromId, err := strconv.ParseInt(fromIdInput, 10, 64); err == nil {
				req.FromId = fromId
			}
		}

		fmt.Print("Limit (max 1000, default 500): ")
		limitInput, _ := reader.ReadString('\n')
		limitInput = strings.TrimSpace(limitInput)

		limit := 500
		if l, err := strconv.Atoi(limitInput); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
		req.Limit = limit

		req.StartTime = time.Now().AddDate(0, 0, -days)
		req.EndTime = time.Now()

		fmt.Printf("\n⏳ Fetching trades from last %d days...\n", days)
	}

	trades, err := client.Trading.GetMyTrades(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(trades) == 0 {
		if req.OrderId > 0 {
			fmt.Printf("📭 No trades found for Order ID %d\n", req.OrderId)
		} else {
			fmt.Println("📭 No trades found in the specified period")
		}
		return
	}

	fmt.Printf("\n✅ Found %d trades!\n", len(trades))

	// Calculate statistics
	var totalBuyQty, totalSellQty float64
	var totalBuyValue, totalSellValue float64
	var buyCount, sellCount int
	var totalCommission = make(map[string]float64)

	for _, trade := range trades {
		qty, _ := strconv.ParseFloat(trade.Qty, 64)
		quoteQty, _ := strconv.ParseFloat(trade.QuoteQty, 64)
		commission, _ := strconv.ParseFloat(trade.Commission, 64)

		if trade.IsBuyer {
			buyCount++
			totalBuyQty += qty
			totalBuyValue += quoteQty
		} else {
			sellCount++
			totalSellQty += qty
			totalSellValue += quoteQty
		}

		// Sum commissions by asset
		totalCommission[trade.CommissionAsset] += commission
	}

	// Display statistics
	fmt.Println("\n📊 Trade Statistics:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("  Total Trades: %d\n", len(trades))
	if buyCount > 0 {
		fmt.Printf("  Buy Trades: %d (Total Qty: %.8f)\n", buyCount, totalBuyQty)
		fmt.Printf("  Total Buy Value: %.2f\n", totalBuyValue)
	}
	if sellCount > 0 {
		fmt.Printf("  Sell Trades: %d (Total Qty: %.8f)\n", sellCount, totalSellQty)
		fmt.Printf("  Total Sell Value: %.2f\n", totalSellValue)
	}

	if len(totalCommission) > 0 {
		fmt.Println("\n  Commissions Paid:")
		for asset, amount := range totalCommission {
			if amount > 0 {
				fmt.Printf("    %s: %.8f\n", asset, amount)
			}
		}
	}
	fmt.Println(strings.Repeat("-", 40))

	// Print raw JSON response
	jsonBytes, err := json.MarshalIndent(trades, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
	} else {
		fmt.Println("\n📊 Full Response (JSON):")
		fmt.Println(string(jsonBytes))
	}
}

func testOrder(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n🧪 Test Order (Simulation)")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("This will validate your order without actually placing it")

	fmt.Print("\nEnter symbol (e.g., BTCUSDT): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	fmt.Print("BUY or SELL? ")
	sideInput, _ := reader.ReadString('\n')
	sideInput = strings.TrimSpace(strings.ToUpper(sideInput))
	
	var side models.OrderSide
	if sideInput == "BUY" {
		side = models.SideBuy
	} else {
		side = models.SideSell
	}

	fmt.Print("Order type (LIMIT/MARKET, default LIMIT): ")
	typeInput, _ := reader.ReadString('\n')
	typeInput = strings.TrimSpace(strings.ToUpper(typeInput))
	
	orderType := models.OrderTypeLimit
	needPrice := true
	if typeInput == "MARKET" {
		orderType = models.OrderTypeMarket
		needPrice = false
	}

	req := models.NewOrderRequest{
		Symbol:   symbol,
		Side:     side,
		Type:     orderType,
	}

	// Handle quantity input based on order type
	if orderType == models.OrderTypeMarket {
		fmt.Println("\nFor MARKET orders, choose quantity type:")
		fmt.Println("1. Specify base asset quantity")
		fmt.Println("2. Specify quote asset amount (quoteOrderQty)")
		fmt.Print("✏️  Select (1-2, default 1): ")
		qtyTypeInput, _ := reader.ReadString('\n')
		qtyTypeInput = strings.TrimSpace(qtyTypeInput)
		
		if qtyTypeInput == "2" {
			fmt.Print("Enter quote order quantity: ")
			quoteOrderQty, _ := reader.ReadString('\n')
			req.QuoteOrderQty = strings.TrimSpace(quoteOrderQty)
		} else {
			fmt.Print("Enter quantity: ")
			quantity, _ := reader.ReadString('\n')
			req.Quantity = strings.TrimSpace(quantity)
		}
	} else {
		fmt.Print("Quantity: ")
		quantity, _ := reader.ReadString('\n')
		req.Quantity = strings.TrimSpace(quantity)
	}

	if needPrice {
		fmt.Print("Price: ")
		price, _ := reader.ReadString('\n')
		req.Price = strings.TrimSpace(price)
		req.TimeInForce = models.TimeInForceGTC
	}

	fmt.Println("\n⏳ Testing order...")
	err := client.Trading.TestNewOrder(req)
	if err != nil {
		fmt.Printf("❌ Order validation failed: %v\n", err)
		fmt.Println("   Please check your parameters and try again")
	} else {
		fmt.Println("✅ Order is valid! This order would be accepted if placed.")
		fmt.Println("   Note: This was just a test, no actual order was placed.")
	}
}

func displayOrder(num int, order models.Order) {
	orderTime := time.Unix(order.Time/1000, 0)
	
	fmt.Printf("%d. Order #%d | %s | %s\n", 
		num, order.OrderId, order.Symbol, orderTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   %s %s @ %s | Qty: %s/%s\n",
		order.Side, order.Type, order.Price, order.ExecutedQty, order.OrigQty)
	fmt.Printf("   Status: %s", order.Status)
	
	if order.Status == models.OrderStatusPartiallyFilled {
		executed, _ := strconv.ParseFloat(order.ExecutedQty, 64)
		orig, _ := strconv.ParseFloat(order.OrigQty, 64)
		percent := (executed / orig) * 100
		fmt.Printf(" (%.1f%% filled)", percent)
	}
	fmt.Println()
	fmt.Println()
}

func displayOrderDetails(order models.Order) {
	fmt.Printf("  Order ID: %d\n", order.OrderId)
	fmt.Printf("  Client Order ID: %s\n", order.ClientOrderId)
	fmt.Printf("  Symbol: %s\n", order.Symbol)
	fmt.Printf("  Type: %s\n", order.Type)
	fmt.Printf("  Side: %s\n", order.Side)
	fmt.Printf("  Status: %s\n", order.Status)
	fmt.Printf("  Price: %s\n", order.Price)
	
	if order.StopPrice != "" {
		fmt.Printf("  Stop Price: %s\n", order.StopPrice)
	}
	
	fmt.Printf("  Original Qty: %s\n", order.OrigQty)
	fmt.Printf("  Executed Qty: %s\n", order.ExecutedQty)
	
	if order.ExecutedQty != "0" && order.ExecutedQty != "" {
		executed, _ := strconv.ParseFloat(order.ExecutedQty, 64)
		orig, _ := strconv.ParseFloat(order.OrigQty, 64)
		percent := (executed / orig) * 100
		fmt.Printf("  Filled: %.1f%%\n", percent)
	}
	
	fmt.Printf("  Cumulative Quote Qty: %s\n", order.CummulativeQuoteQty)
	
	if order.TimeInForce != "" {
		fmt.Printf("  Time in Force: %s\n", order.TimeInForce)
	}
	
	orderTime := time.Unix(order.Time/1000, 0)
	updateTime := time.Unix(order.UpdateTime/1000, 0)
	fmt.Printf("  Created: %s\n", orderTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Updated: %s\n", updateTime.Format("2006-01-02 15:04:05"))
	
	fmt.Printf("  Is Working: %v\n", order.IsWorking)
}


func viewTradingRules(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n📐 Trading Rules & Filters")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("View trading rules for a symbol (LOT_SIZE, NOTIONAL, etc.)")

	fmt.Print("\nEnter symbol (e.g., BTCUSDT): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	if symbol == "" {
		symbol = "BTCUSDT"
		fmt.Printf("Using default: %s\n", symbol)
	}

	fmt.Printf("\n⏳ Fetching trading rules for %s...\n", symbol)

	req := models.ExchangeInfoRequest{
		Symbol: symbol,
	}

	info, err := client.Market.GetExchangeInfo(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(info.Symbols) == 0 {
		fmt.Printf("❌ Symbol %s not found\n", symbol)
		return
	}

	symbolInfo := info.Symbols[0]

	fmt.Printf("\n✅ Trading Rules for %s:\n", symbolInfo.Symbol)
	fmt.Println(strings.Repeat("=", 60))

	// Basic Info
	fmt.Printf("\n📊 Basic Information:\n")
	fmt.Printf("  Status: %s\n", symbolInfo.Status)
	fmt.Printf("  Base Asset: %s (Precision: %d)\n", symbolInfo.BaseAsset, symbolInfo.BaseAssetPrecision)
	fmt.Printf("  Quote Asset: %s (Precision: %d)\n", symbolInfo.QuoteAsset, symbolInfo.QuoteAssetPrecision)
	fmt.Printf("  Allowed Order Types: %v\n", symbolInfo.OrderTypes)
	fmt.Printf("  Permissions: %v\n", symbolInfo.Permissions)

	// Parse and display filters
	fmt.Printf("\n🔧 Trading Filters:\n")
	for _, filter := range symbolInfo.Filters {
		filterType, ok := filter["filterType"].(string)
		if !ok {
			continue
		}

		switch filterType {
		case "PRICE_FILTER":
			fmt.Printf("\n  💰 PRICE_FILTER:\n")
			if minPrice, ok := filter["minPrice"].(string); ok {
				fmt.Printf("    Min Price: %s\n", minPrice)
			}
			if maxPrice, ok := filter["maxPrice"].(string); ok {
				fmt.Printf("    Max Price: %s\n", maxPrice)
			}
			if tickSize, ok := filter["tickSize"].(string); ok {
				fmt.Printf("    Tick Size (Price Increment): %s\n", tickSize)
				fmt.Printf("    ➡️  Price must be a multiple of %s\n", tickSize)
			}

		case "LOT_SIZE":
			fmt.Printf("\n  📦 LOT_SIZE (Quantity Rules):\n")
			if minQty, ok := filter["minQty"].(string); ok {
				fmt.Printf("    Min Quantity: %s\n", minQty)
			}
			if maxQty, ok := filter["maxQty"].(string); ok {
				fmt.Printf("    Max Quantity: %s\n", maxQty)
			}
			if stepSize, ok := filter["stepSize"].(string); ok {
				fmt.Printf("    Step Size (Quantity Increment): %s\n", stepSize)
				fmt.Printf("    ➡️  Quantity must be a multiple of %s\n", stepSize)
			}

		case "MIN_NOTIONAL":
			fmt.Printf("\n  💵 MIN_NOTIONAL (Minimum Order Value):\n")
			if minNotional, ok := filter["minNotional"].(string); ok {
				fmt.Printf("    Min Notional: %s %s\n", minNotional, symbolInfo.QuoteAsset)
				fmt.Printf("    ➡️  Order value (price × quantity) must be ≥ %s %s\n", minNotional, symbolInfo.QuoteAsset)
			}
			if applyToMarket, ok := filter["applyToMarket"].(bool); ok {
				fmt.Printf("    Apply to Market Orders: %v\n", applyToMarket)
			}
			if avgPriceMins, ok := filter["avgPriceMins"].(float64); ok {
				fmt.Printf("    Average Price Minutes: %.0f\n", avgPriceMins)
			}

		case "NOTIONAL":
			fmt.Printf("\n  💴 NOTIONAL (Order Value Range):\n")
			if minNotional, ok := filter["minNotional"].(string); ok {
				fmt.Printf("    Min Notional: %s %s\n", minNotional, symbolInfo.QuoteAsset)
			}
			if maxNotional, ok := filter["maxNotional"].(string); ok {
				fmt.Printf("    Max Notional: %s %s\n", maxNotional, symbolInfo.QuoteAsset)
			}
			if applyMinToMarket, ok := filter["applyMinToMarket"].(bool); ok {
				fmt.Printf("    Apply Min to Market: %v\n", applyMinToMarket)
			}
			if applyMaxToMarket, ok := filter["applyMaxToMarket"].(bool); ok {
				fmt.Printf("    Apply Max to Market: %v\n", applyMaxToMarket)
			}
			if avgPriceMins, ok := filter["avgPriceMins"].(float64); ok {
				fmt.Printf("    Average Price Minutes: %.0f\n", avgPriceMins)
			}

		case "ICEBERG_PARTS":
			fmt.Printf("\n  🧊 ICEBERG_PARTS:\n")
			if limit, ok := filter["limit"].(float64); ok {
				fmt.Printf("    Max Parts: %.0f\n", limit)
			}

		case "MARKET_LOT_SIZE":
			fmt.Printf("\n  📈 MARKET_LOT_SIZE (Market Order Quantity):\n")
			if minQty, ok := filter["minQty"].(string); ok {
				fmt.Printf("    Min Market Qty: %s\n", minQty)
			}
			if maxQty, ok := filter["maxQty"].(string); ok {
				fmt.Printf("    Max Market Qty: %s\n", maxQty)
			}
			if stepSize, ok := filter["stepSize"].(string); ok {
				fmt.Printf("    Market Step Size: %s\n", stepSize)
			}

		case "MAX_NUM_ORDERS":
			fmt.Printf("\n  📝 MAX_NUM_ORDERS:\n")
			if limit, ok := filter["limit"].(float64); ok {
				fmt.Printf("    Max Open Orders: %.0f\n", limit)
			}

		case "MAX_NUM_ALGO_ORDERS":
			fmt.Printf("\n  🤖 MAX_NUM_ALGO_ORDERS:\n")
			if limit, ok := filter["limit"].(float64); ok {
				fmt.Printf("    Max Algo Orders: %.0f\n", limit)
			}

		case "PERCENT_PRICE":
			fmt.Printf("\n  📊 PERCENT_PRICE:\n")
			if multiplierUp, ok := filter["multiplierUp"].(string); ok {
				fmt.Printf("    Multiplier Up: %s\n", multiplierUp)
			}
			if multiplierDown, ok := filter["multiplierDown"].(string); ok {
				fmt.Printf("    Multiplier Down: %s\n", multiplierDown)
			}
			if avgPriceMins, ok := filter["avgPriceMins"].(float64); ok {
				fmt.Printf("    Average Price Minutes: %.0f\n", avgPriceMins)
			}

		case "TRAILING_DELTA":
			fmt.Printf("\n  📉 TRAILING_DELTA:\n")
			if minTrailingAboveDelta, ok := filter["minTrailingAboveDelta"].(float64); ok {
				fmt.Printf("    Min Trailing Above Delta: %.0f\n", minTrailingAboveDelta)
			}
			if maxTrailingAboveDelta, ok := filter["maxTrailingAboveDelta"].(float64); ok {
				fmt.Printf("    Max Trailing Above Delta: %.0f\n", maxTrailingAboveDelta)
			}
			if minTrailingBelowDelta, ok := filter["minTrailingBelowDelta"].(float64); ok {
				fmt.Printf("    Min Trailing Below Delta: %.0f\n", minTrailingBelowDelta)
			}
			if maxTrailingBelowDelta, ok := filter["maxTrailingBelowDelta"].(float64); ok {
				fmt.Printf("    Max Trailing Below Delta: %.0f\n", maxTrailingBelowDelta)
			}
		}
	}

	// Trading tips based on rules
	fmt.Printf("\n💡 Quick Reference for %s:\n", symbol)
	fmt.Println(strings.Repeat("-", 50))
	
	// Find key filters for tips
	for _, filter := range symbolInfo.Filters {
		filterType, _ := filter["filterType"].(string)
		
		if filterType == "LOT_SIZE" {
			if minQty, ok := filter["minQty"].(string); ok {
				if stepSize, ok := filter["stepSize"].(string); ok {
					fmt.Printf("• Quantity: Min=%s, must be multiple of %s\n", minQty, stepSize)
				}
			}
		}
		
		if filterType == "PRICE_FILTER" {
			if tickSize, ok := filter["tickSize"].(string); ok {
				fmt.Printf("• Price: Must be multiple of %s\n", tickSize)
			}
		}
		
		if filterType == "MIN_NOTIONAL" || filterType == "NOTIONAL" {
			if minNotional, ok := filter["minNotional"].(string); ok {
				fmt.Printf("• Min Order Value: %s %s (price × quantity)\n", minNotional, symbolInfo.QuoteAsset)
			}
		}
	}

	// Option to check another symbol
	fmt.Print("\n🔍 Check another symbol? (yes/no): ")
	another, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(another)) == "yes" {
		viewTradingRules(client, reader)
	}
}