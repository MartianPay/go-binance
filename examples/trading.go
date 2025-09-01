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
	fmt.Println("7. View My Trades")
	fmt.Println("8. Test Order (Simulation)")
	fmt.Println("9. Exit")
}

func viewAccountInfo(client *binance.BinanceClient) {
	fmt.Println("\n💰 Account Information")
	fmt.Println(strings.Repeat("-", 50))

	info, err := client.Trading.GetAccountInfo(0)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("Account Type: %s\n", info.AccountType)
	fmt.Printf("Can Trade: %v\n", info.CanTrade)
	fmt.Printf("Can Withdraw: %v\n", info.CanWithdraw)
	fmt.Printf("Can Deposit: %v\n", info.CanDeposit)
	
	fmt.Printf("\nCommissions:\n")
	fmt.Printf("  Maker: %d (0.%02d%%)\n", info.MakerCommission, info.MakerCommission)
	fmt.Printf("  Taker: %d (0.%02d%%)\n", info.TakerCommission, info.TakerCommission)

	fmt.Printf("\nPermissions: %v\n", info.Permissions)

	fmt.Printf("\n💼 Balances (non-zero only):\n")
	hasBalance := false
	for _, balance := range info.Balances {
		free, _ := strconv.ParseFloat(balance.Free, 64)
		locked, _ := strconv.ParseFloat(balance.Locked, 64)
		
		if free > 0 || locked > 0 {
			hasBalance = true
			fmt.Printf("  %s:\n", balance.Asset)
			if free > 0 {
				fmt.Printf("    Free: %s\n", balance.Free)
			}
			if locked > 0 {
				fmt.Printf("    Locked: %s\n", balance.Locked)
			}
		}
	}

	if !hasBalance {
		fmt.Println("  No balances found")
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

	fmt.Printf("\n📊 Found %d open orders:\n\n", len(orders))
	
	for i, order := range orders {
		displayOrder(i+1, order)
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

	// 4. Quantity
	fmt.Print("\nEnter quantity: ")
	quantity, _ := reader.ReadString('\n')
	quantity = strings.TrimSpace(quantity)

	req := models.NewOrderRequest{
		Symbol:           symbol,
		Side:             side,
		Type:             orderType,
		Quantity:         quantity,
		NewOrderRespType: models.OrderResponseTypeFULL,
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

	// 7. Confirm
	fmt.Println("\n📋 Order Summary:")
	fmt.Printf("  Symbol: %s\n", req.Symbol)
	fmt.Printf("  Side: %s\n", req.Side)
	fmt.Printf("  Type: %s\n", req.Type)
	fmt.Printf("  Quantity: %s\n", req.Quantity)
	if req.Price != "" {
		fmt.Printf("  Price: %s\n", req.Price)
	}
	if req.StopPrice != "" {
		fmt.Printf("  Stop Price: %s\n", req.StopPrice)
	}

	fmt.Print("\n⚠️  Confirm order? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(confirm)) != "yes" {
		fmt.Println("❌ Order cancelled")
		return
	}

	// 8. Execute
	fmt.Println("\n🚀 Placing order...")
	order, err := client.Trading.NewOrder(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Order placed successfully!\n")
	fmt.Printf("  Order ID: %d\n", order.OrderId)
	fmt.Printf("  Client Order ID: %s\n", order.ClientOrderId)
	fmt.Printf("  Status: %s\n", order.Status)
	fmt.Printf("  Price: %s\n", order.Price)
	fmt.Printf("  Original Qty: %s\n", order.OrigQty)
	fmt.Printf("  Executed Qty: %s\n", order.ExecutedQty)

	if len(order.Fills) > 0 {
		fmt.Println("\n📊 Fills:")
		for _, fill := range order.Fills {
			fmt.Printf("  Price: %s, Qty: %s, Commission: %s %s\n",
				fill.Price, fill.Qty, fill.Commission, fill.CommissionAsset)
		}
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

	fmt.Println("\n📊 Order Details:")
	displayOrderDetails(*order)

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

	fmt.Printf("\n✅ Order cancelled successfully!\n")
	fmt.Printf("  Order ID: %d\n", result.OrderId)
	fmt.Printf("  Symbol: %s\n", result.Symbol)
	fmt.Printf("  Status: %s\n", result.Status)
	fmt.Printf("  Original Qty: %s\n", result.OrigQty)
	fmt.Printf("  Executed Qty: %s\n", result.ExecutedQty)
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

	// Group by status
	statusGroups := make(map[models.OrderStatus][]models.Order)
	for _, order := range orders {
		statusGroups[order.Status] = append(statusGroups[order.Status], order)
	}

	fmt.Printf("\n📊 Found %d orders:\n", len(orders))
	
	// Display summary
	for status, orderList := range statusGroups {
		fmt.Printf("  %s: %d orders\n", status, len(orderList))
	}

	// Display orders
	fmt.Println("\n📋 Order Details:")
	for i, order := range orders {
		if i >= 20 { // Show only first 20 in detail
			fmt.Printf("\n... and %d more orders\n", len(orders)-20)
			break
		}
		displayOrder(i+1, order)
	}
}

func viewMyTrades(client *binance.BinanceClient, reader *bufio.Reader) {
	fmt.Println("\n💱 My Trades")
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

	req := models.MyTradesRequest{
		Symbol:    symbol,
		StartTime: time.Now().AddDate(0, 0, -days),
		EndTime:   time.Now(),
		Limit:     100,
	}

	fmt.Printf("\n⏳ Fetching trades from last %d days...\n", days)
	trades, err := client.Trading.GetMyTrades(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(trades) == 0 {
		fmt.Println("📭 No trades found in the specified period")
		return
	}

	fmt.Printf("\n📊 Found %d trades:\n\n", len(trades))
	
	var totalVolume, totalCommission float64
	commissionAssets := make(map[string]float64)

	for i, trade := range trades {
		tradeTime := time.Unix(trade.Time/1000, 0)
		
		side := "SELL"
		if trade.IsBuyer {
			side = "BUY"
		}
		
		maker := "TAKER"
		if trade.IsMaker {
			maker = "MAKER"
		}

		fmt.Printf("%d. %s | %s %s | Price: %s | Qty: %s | %s\n",
			i+1,
			tradeTime.Format("2006-01-02 15:04:05"),
			side,
			trade.Symbol,
			trade.Price,
			trade.Qty,
			maker)

		if trade.Commission != "" {
			fmt.Printf("   Commission: %s %s | Order ID: %d\n",
				trade.Commission, trade.CommissionAsset, trade.OrderId)
			
			commission, _ := strconv.ParseFloat(trade.Commission, 64)
			commissionAssets[trade.CommissionAsset] += commission
		}

		volume, _ := strconv.ParseFloat(trade.QuoteQty, 64)
		totalVolume += volume
	}

	// Summary
	fmt.Printf("\n📈 Summary:\n")
	fmt.Printf("  Total trades: %d\n", len(trades))
	fmt.Printf("  Total volume: %.2f\n", totalVolume)
	fmt.Printf("  Commissions paid:\n")
	for asset, amount := range commissionAssets {
		fmt.Printf("    %s: %.8f\n", asset, amount)
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

	fmt.Print("Quantity: ")
	quantity, _ := reader.ReadString('\n')
	quantity = strings.TrimSpace(quantity)

	req := models.NewOrderRequest{
		Symbol:   symbol,
		Side:     side,
		Type:     orderType,
		Quantity: quantity,
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