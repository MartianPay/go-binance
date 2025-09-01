package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MartianPay/go-binance/client"
	"github.com/MartianPay/go-binance/models"
)

type MarketDataService struct {
	client *client.Client
}

func NewMarketDataService(c *client.Client) *MarketDataService {
	return &MarketDataService{client: c}
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