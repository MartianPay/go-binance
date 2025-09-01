package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MartianPay/go-binance/client"
	"github.com/MartianPay/go-binance/models"
)

type TradingService struct {
	client *client.Client
}

func NewTradingService(c *client.Client) *TradingService {
	return &TradingService{client: c}
}

// TestNewOrder tests new order creation without actually sending it
// API endpoint: POST /api/v3/order/test
func (s *TradingService) TestNewOrder(req models.NewOrderRequest) error {
	params := s.buildOrderParams(req)
	
	_, err := s.client.Post("/api/v3/order/test", params, nil, true)
	if err != nil {
		return fmt.Errorf("failed to test new order: %w", err)
	}
	
	return nil
}

// NewOrder creates a new order
// API endpoint: POST /api/v3/order
func (s *TradingService) NewOrder(req models.NewOrderRequest) (*models.OrderResponse, error) {
	params := s.buildOrderParams(req)
	
	resp, err := s.client.Post("/api/v3/order", params, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create new order: %w", err)
	}
	
	var order models.OrderResponse
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}
	
	return &order, nil
}

// QueryOrder checks an order's status
// API endpoint: GET /api/v3/order
func (s *TradingService) QueryOrder(req models.QueryOrderRequest) (*models.Order, error) {
	params := make(map[string]string)
	params["symbol"] = req.Symbol
	
	if req.OrderId > 0 {
		params["orderId"] = strconv.FormatInt(req.OrderId, 10)
	}
	
	if req.OrigClientOrderId != "" {
		params["origClientOrderId"] = req.OrigClientOrderId
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Get("/api/v3/order", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to query order: %w", err)
	}
	
	var order models.Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}
	
	return &order, nil
}

// CancelOrder cancels an active order
// API endpoint: DELETE /api/v3/order
func (s *TradingService) CancelOrder(req models.CancelOrderRequest) (*models.CancelOrderResponse, error) {
	params := make(map[string]string)
	params["symbol"] = req.Symbol
	
	if req.OrderId > 0 {
		params["orderId"] = strconv.FormatInt(req.OrderId, 10)
	}
	
	if req.OrigClientOrderId != "" {
		params["origClientOrderId"] = req.OrigClientOrderId
	}
	
	if req.NewClientOrderId != "" {
		params["newClientOrderId"] = req.NewClientOrderId
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Delete("/api/v3/order", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}
	
	var result models.CancelOrderResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel response: %w", err)
	}
	
	return &result, nil
}

// CancelAllOpenOrders cancels all open orders on a symbol
// API endpoint: DELETE /api/v3/openOrders
func (s *TradingService) CancelAllOpenOrders(symbol string, recvWindow int64) ([]models.CancelOrderResponse, error) {
	params := make(map[string]string)
	params["symbol"] = symbol
	
	if recvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow, 10)
	}
	
	resp, err := s.client.Delete("/api/v3/openOrders", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel all open orders: %w", err)
	}
	
	var results []models.CancelOrderResponse
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel responses: %w", err)
	}
	
	return results, nil
}

// GetOpenOrders gets all open orders
// API endpoint: GET /api/v3/openOrders
func (s *TradingService) GetOpenOrders(req models.OpenOrdersRequest) ([]models.Order, error) {
	params := make(map[string]string)
	
	if req.Symbol != "" {
		params["symbol"] = req.Symbol
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Get("/api/v3/openOrders", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}
	
	var orders []models.Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders: %w", err)
	}
	
	return orders, nil
}

// GetAllOrders gets all orders (active, canceled, or filled)
// API endpoint: GET /api/v3/allOrders
func (s *TradingService) GetAllOrders(req models.AllOrdersRequest) ([]models.Order, error) {
	params := make(map[string]string)
	params["symbol"] = req.Symbol
	
	if req.OrderId > 0 {
		params["orderId"] = strconv.FormatInt(req.OrderId, 10)
	}
	
	if !req.StartTime.IsZero() {
		params["startTime"] = strconv.FormatInt(req.StartTime.Unix()*1000, 10)
	}
	
	if !req.EndTime.IsZero() {
		params["endTime"] = strconv.FormatInt(req.EndTime.Unix()*1000, 10)
	}
	
	if req.Limit > 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Get("/api/v3/allOrders", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	
	var orders []models.Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders: %w", err)
	}
	
	return orders, nil
}

// GetAccountInfo gets current account information
// API endpoint: GET /api/v3/account
func (s *TradingService) GetAccountInfo(recvWindow int64) (*models.TradingAccountInfo, error) {
	params := make(map[string]string)
	
	if recvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow, 10)
	}
	
	resp, err := s.client.Get("/api/v3/account", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}
	
	var info models.TradingAccountInfo
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account info: %w", err)
	}
	
	return &info, nil
}

// GetMyTrades gets trades for a specific account and symbol
// API endpoint: GET /api/v3/myTrades
func (s *TradingService) GetMyTrades(req models.MyTradesRequest) ([]models.Trade, error) {
	params := make(map[string]string)
	params["symbol"] = req.Symbol
	
	if req.OrderId > 0 {
		params["orderId"] = strconv.FormatInt(req.OrderId, 10)
	}
	
	if !req.StartTime.IsZero() {
		params["startTime"] = strconv.FormatInt(req.StartTime.Unix()*1000, 10)
	}
	
	if !req.EndTime.IsZero() {
		params["endTime"] = strconv.FormatInt(req.EndTime.Unix()*1000, 10)
	}
	
	if req.FromId > 0 {
		params["fromId"] = strconv.FormatInt(req.FromId, 10)
	}
	
	if req.Limit > 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Get("/api/v3/myTrades", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get my trades: %w", err)
	}
	
	var trades []models.Trade
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trades: %w", err)
	}
	
	return trades, nil
}

// buildOrderParams builds parameters for order requests
func (s *TradingService) buildOrderParams(req models.NewOrderRequest) map[string]string {
	params := make(map[string]string)
	params["symbol"] = req.Symbol
	params["side"] = string(req.Side)
	params["type"] = string(req.Type)
	
	if req.TimeInForce != "" {
		params["timeInForce"] = string(req.TimeInForce)
	}
	
	if req.Quantity != "" {
		params["quantity"] = req.Quantity
	}
	
	if req.QuoteOrderQty != "" {
		params["quoteOrderQty"] = req.QuoteOrderQty
	}
	
	if req.Price != "" {
		params["price"] = req.Price
	}
	
	if req.NewClientOrderId != "" {
		params["newClientOrderId"] = req.NewClientOrderId
	}
	
	if req.StopPrice != "" {
		params["stopPrice"] = req.StopPrice
	}
	
	if req.IcebergQty != "" {
		params["icebergQty"] = req.IcebergQty
	}
	
	if req.NewOrderRespType != "" {
		params["newOrderRespType"] = string(req.NewOrderRespType)
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	return params
}