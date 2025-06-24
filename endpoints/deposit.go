package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MartianPay/go-binance/client"
	"github.com/MartianPay/go-binance/models"
)

type DepositService struct {
	client *client.Client
}

func NewDepositService(c *client.Client) *DepositService {
	return &DepositService{client: c}
}

func (s *DepositService) GetDepositAddress(req models.DepositAddressRequest) (*models.DepositAddress, error) {
	params := make(map[string]string)
	params["coin"] = req.Coin
	
	if req.Network != "" {
		params["network"] = req.Network
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Get("/sapi/v1/capital/deposit/address", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit address: %w", err)
	}
	
	var address models.DepositAddress
	if err := json.Unmarshal(resp, &address); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &address, nil
}

func (s *DepositService) GetDepositHistory(req models.DepositHistoryRequest) ([]models.DepositHistory, error) {
	params := make(map[string]string)
	
	if req.Coin != "" {
		params["coin"] = req.Coin
	}
	
	if req.Status != 0 {
		params["status"] = strconv.Itoa(req.Status)
	}
	
	if !req.StartTime.IsZero() {
		params["startTime"] = strconv.FormatInt(req.StartTime.UnixMilli(), 10)
	}
	
	if !req.EndTime.IsZero() {
		params["endTime"] = strconv.FormatInt(req.EndTime.UnixMilli(), 10)
	}
	
	if req.Offset > 0 {
		params["offset"] = strconv.Itoa(req.Offset)
	}
	
	if req.Limit > 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	if req.TxId != "" {
		params["txId"] = req.TxId
	}
	
	resp, err := s.client.Get("/sapi/v1/capital/deposit/hisrec", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit history: %w", err)
	}
	
	var history []models.DepositHistory
	if err := json.Unmarshal(resp, &history); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return history, nil
}