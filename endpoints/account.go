package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MartianPay/go-binance/client"
	"github.com/MartianPay/go-binance/models"
)

type AccountService struct {
	client *client.Client
}

func NewAccountService(c *client.Client) *AccountService {
	return &AccountService{client: c}
}

func (s *AccountService) GetAllCoins() ([]models.CoinInfo, error) {
	params := make(map[string]string)
	
	resp, err := s.client.Get("/sapi/v1/capital/config/getall", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get all coins: %w", err)
	}
	
	var coins []models.CoinInfo
	if err := json.Unmarshal(resp, &coins); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return coins, nil
}

func (s *AccountService) GetAccountInfo() (*models.AccountInfo, error) {
	params := make(map[string]string)
	
	resp, err := s.client.Get("/sapi/v1/account/info", params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}
	
	var info models.AccountInfo
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &info, nil
}

func (s *AccountService) UniversalTransfer(req models.AssetTransferRequest) (*models.AssetTransferResponse, error) {
	params := make(map[string]string)
	params["type"] = req.Type
	params["asset"] = req.Asset
	params["amount"] = req.Amount
	
	if req.FromSymbol != "" {
		params["fromSymbol"] = req.FromSymbol
	}
	
	if req.ToSymbol != "" {
		params["toSymbol"] = req.ToSymbol
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Post("/sapi/v1/asset/transfer", params, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer asset: %w", err)
	}
	
	var result models.AssetTransferResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &result, nil
}

func (s *AccountService) GetUserAsset(req models.UserAssetRequest) ([]models.UserAsset, error) {
	params := make(map[string]string)
	
	if req.Asset != "" {
		params["asset"] = req.Asset
	}
	
	if req.NeedBtcValuation {
		params["needBtcValuation"] = "true"
	}
	
	if req.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(req.RecvWindow, 10)
	}
	
	resp, err := s.client.Post("/sapi/v3/asset/getUserAsset", params, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get user asset: %w", err)
	}
	
	var assets []models.UserAsset
	if err := json.Unmarshal(resp, &assets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return assets, nil
}

func (s *AccountService) EnableFastWithdrawSwitch(recvWindow int64) error {
	params := make(map[string]string)
	
	if recvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow, 10)
	}
	
	_, err := s.client.Post("/sapi/v1/account/enableFastWithdrawSwitch", params, nil, true)
	if err != nil {
		return fmt.Errorf("failed to enable fast withdraw switch: %w", err)
	}
	
	return nil
}