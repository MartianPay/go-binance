package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MartianPay/go-binance/client"
	"github.com/MartianPay/go-binance/models"
)

type WithdrawalService struct {
	client *client.Client
}

func NewWithdrawalService(c *client.Client) *WithdrawalService {
	return &WithdrawalService{client: c}
}

func (s *WithdrawalService) Withdraw(req models.WithdrawalRequest) (*models.WithdrawalResponse, error) {
	params := make(map[string]string)
	params["coin"] = req.Coin
	params["address"] = req.Address
	params["amount"] = req.Amount
	
	if req.Network != "" {
		params["network"] = req.Network
	}
	
	if req.AddressTag != "" {
		params["addressTag"] = req.AddressTag
	}
	
	if req.WithdrawOrderId != "" {
		params["withdrawOrderId"] = req.WithdrawOrderId
	}
	
	if req.TransactionFeeFlag {
		params["transactionFeeFlag"] = "true"
	}
	
	if req.Name != "" {
		params["name"] = req.Name
	}
	
	if req.WalletType != 0 {
		params["walletType"] = strconv.Itoa(req.WalletType)
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Post("/sapi/v1/capital/withdraw/apply", params, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw: %w", err)
	}
	
	var result models.WithdrawalResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &result, nil
}

func (s *WithdrawalService) GetWithdrawalHistory(req models.WithdrawalHistoryRequest) ([]models.WithdrawalHistory, error) {
	params := make(map[string]string)
	
	if req.Coin != "" {
		params["coin"] = req.Coin
	}
	
	if req.WithdrawOrderId != "" {
		params["withdrawOrderId"] = req.WithdrawOrderId
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
	
	resp, err := s.client.Get("/sapi/v1/capital/withdraw/history", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal history: %w", err)
	}
	
	var history []models.WithdrawalHistory
	if err := json.Unmarshal(resp, &history); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return history, nil
}

// GetWithdrawalQuota 获取24小时提现限额信息
func (s *WithdrawalService) GetWithdrawalQuota() (*models.WithdrawalQuota, error) {
	params := make(map[string]string)
	
	resp, err := s.client.Get("/sapi/v1/capital/withdraw/quota", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal quota: %w", err)
	}
	
	var quota models.WithdrawalQuota
	if err := json.Unmarshal(resp, &quota); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &quota, nil
}

// GetWithdrawalAddressList 获取提现地址列表
func (s *WithdrawalService) GetWithdrawalAddressList() ([]models.WithdrawalAddress, error) {
	params := make(map[string]string)
	
	resp, err := s.client.Get("/sapi/v1/capital/withdraw/address/list", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal address list: %w", err)
	}
	
	var addresses []models.WithdrawalAddress
	if err := json.Unmarshal(resp, &addresses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return addresses, nil
}