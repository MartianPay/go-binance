package models

import (
	"encoding/json"
	"strconv"
)

// KlineInterval represents the interval for kline data
type KlineInterval string

const (
	Interval1s  KlineInterval = "1s"
	Interval1m  KlineInterval = "1m"
	Interval3m  KlineInterval = "3m"
	Interval5m  KlineInterval = "5m"
	Interval15m KlineInterval = "15m"
	Interval30m KlineInterval = "30m"
	Interval1h  KlineInterval = "1h"
	Interval2h  KlineInterval = "2h"
	Interval4h  KlineInterval = "4h"
	Interval6h  KlineInterval = "6h"
	Interval8h  KlineInterval = "8h"
	Interval12h KlineInterval = "12h"
	Interval1d  KlineInterval = "1d"
	Interval3d  KlineInterval = "3d"
	Interval1w  KlineInterval = "1w"
	Interval1M  KlineInterval = "1M"
)

// KlineRequest represents the request parameters for kline data
type KlineRequest struct {
	Symbol    string        `json:"symbol"`
	Interval  KlineInterval `json:"interval"`
	StartTime int64         `json:"startTime,omitempty"`
	EndTime   int64         `json:"endTime,omitempty"`
	TimeZone  string        `json:"timeZone,omitempty"`
	Limit     int           `json:"limit,omitempty"` // Default 500; max 1500
}

// Kline represents a single kline/candlestick data point
type Kline struct {
	OpenTime                 int64  `json:"openTime"`
	Open                     string `json:"open"`
	High                     string `json:"high"`
	Low                      string `json:"low"`
	Close                    string `json:"close"`
	Volume                   string `json:"volume"`
	CloseTime                int64  `json:"closeTime"`
	QuoteAssetVolume         string `json:"quoteAssetVolume"`
	NumberOfTrades           int    `json:"numberOfTrades"`
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"`
	Ignore                   string `json:"ignore"`
}

// UnmarshalJSON custom unmarshaler for Kline to handle array response
func (k *Kline) UnmarshalJSON(data []byte) error {
	var raw []interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < 12 {
		return nil
	}

	// Parse open time
	if v, ok := raw[0].(float64); ok {
		k.OpenTime = int64(v)
	}

	// Parse open price
	if v, ok := raw[1].(string); ok {
		k.Open = v
	}

	// Parse high price
	if v, ok := raw[2].(string); ok {
		k.High = v
	}

	// Parse low price
	if v, ok := raw[3].(string); ok {
		k.Low = v
	}

	// Parse close price
	if v, ok := raw[4].(string); ok {
		k.Close = v
	}

	// Parse volume
	if v, ok := raw[5].(string); ok {
		k.Volume = v
	}

	// Parse close time
	if v, ok := raw[6].(float64); ok {
		k.CloseTime = int64(v)
	}

	// Parse quote asset volume
	if v, ok := raw[7].(string); ok {
		k.QuoteAssetVolume = v
	}

	// Parse number of trades
	if v, ok := raw[8].(float64); ok {
		k.NumberOfTrades = int(v)
	}

	// Parse taker buy base asset volume
	if v, ok := raw[9].(string); ok {
		k.TakerBuyBaseAssetVolume = v
	}

	// Parse taker buy quote asset volume
	if v, ok := raw[10].(string); ok {
		k.TakerBuyQuoteAssetVolume = v
	}

	// Parse ignore (unused field)
	if v, ok := raw[11].(string); ok {
		k.Ignore = v
	}

	return nil
}

// GetOpenFloat returns the open price as float64
func (k *Kline) GetOpenFloat() float64 {
	v, _ := strconv.ParseFloat(k.Open, 64)
	return v
}

// GetHighFloat returns the high price as float64
func (k *Kline) GetHighFloat() float64 {
	v, _ := strconv.ParseFloat(k.High, 64)
	return v
}

// GetLowFloat returns the low price as float64
func (k *Kline) GetLowFloat() float64 {
	v, _ := strconv.ParseFloat(k.Low, 64)
	return v
}

// GetCloseFloat returns the close price as float64
func (k *Kline) GetCloseFloat() float64 {
	v, _ := strconv.ParseFloat(k.Close, 64)
	return v
}

// GetVolumeFloat returns the volume as float64
func (k *Kline) GetVolumeFloat() float64 {
	v, _ := strconv.ParseFloat(k.Volume, 64)
	return v
}

// GetQuoteAssetVolumeFloat returns the quote asset volume as float64
func (k *Kline) GetQuoteAssetVolumeFloat() float64 {
	v, _ := strconv.ParseFloat(k.QuoteAssetVolume, 64)
	return v
}