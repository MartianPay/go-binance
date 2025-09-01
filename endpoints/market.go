package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MartianPay/go-binance/client"
	"github.com/MartianPay/go-binance/models"
)

type MarketDataService struct {
	client *client.Client
}

func NewMarketDataService(c *client.Client) *MarketDataService {
	return &MarketDataService{client: c}
}

// GetExchangeInfo retrieves exchange trading rules and symbol information
// API endpoint: GET /api/v3/exchangeInfo
func (s *MarketDataService) GetExchangeInfo(req models.ExchangeInfoRequest) (*models.ExchangeInfo, error) {
	params := make(map[string]string)
	
	if req.Symbol != "" {
		params["symbol"] = req.Symbol
	} else if len(req.Symbols) > 0 {
		// Format: ["BTCUSDT","BNBUSDT"]
		symbols := `["` + strings.Join(req.Symbols, `","`) + `"]`
		params["symbols"] = symbols
	}
	
	if len(req.Permissions) > 0 {
		// Format: ["SPOT","MARGIN"]
		permissions := `["` + strings.Join(req.Permissions, `","`) + `"]`
		params["permissions"] = permissions
	}
	
	resp, err := s.client.Get("/api/v3/exchangeInfo", params, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange info: %w", err)
	}
	
	var info models.ExchangeInfo
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exchange info: %w", err)
	}
	
	return &info, nil
}

// GetKlines retrieves kline/candlestick data for a symbol
// API endpoint: GET /api/v3/klines
func (s *MarketDataService) GetKlines(req models.KlineRequest) ([]models.Kline, error) {
	params := make(map[string]string)
	
	// Required parameters
	params["symbol"] = req.Symbol
	params["interval"] = string(req.Interval)
	
	// Optional parameters
	if req.StartTime > 0 {
		params["startTime"] = strconv.FormatInt(req.StartTime, 10)
	}
	
	if req.EndTime > 0 {
		params["endTime"] = strconv.FormatInt(req.EndTime, 10)
	}
	
	if req.TimeZone != "" {
		params["timeZone"] = req.TimeZone
	}
	
	if req.Limit > 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	
	// Klines endpoint doesn't require authentication
	resp, err := s.client.Get("/api/v3/klines", params, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}
	
	var klines []models.Kline
	if err := json.Unmarshal(resp, &klines); err != nil {
		return nil, fmt.Errorf("failed to unmarshal klines response: %w", err)
	}
	
	return klines, nil
}

// GetUIKlines retrieves kline/candlestick data optimized for presentation
// API endpoint: GET /api/v3/uiKlines
func (s *MarketDataService) GetUIKlines(req models.KlineRequest) ([]models.Kline, error) {
	params := make(map[string]string)
	
	// Required parameters
	params["symbol"] = req.Symbol
	params["interval"] = string(req.Interval)
	
	// Optional parameters
	if req.StartTime > 0 {
		params["startTime"] = strconv.FormatInt(req.StartTime, 10)
	}
	
	if req.EndTime > 0 {
		params["endTime"] = strconv.FormatInt(req.EndTime, 10)
	}
	
	if req.TimeZone != "" {
		params["timeZone"] = req.TimeZone
	}
	
	if req.Limit > 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	
	// UIKlines endpoint doesn't require authentication
	resp, err := s.client.Get("/api/v3/uiKlines", params, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get UI klines: %w", err)
	}
	
	var klines []models.Kline
	if err := json.Unmarshal(resp, &klines); err != nil {
		return nil, fmt.Errorf("failed to unmarshal UI klines response: %w", err)
	}
	
	return klines, nil
}