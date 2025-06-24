package binance

import (
	"github.com/MartianPay/go-binance/client"
	"github.com/MartianPay/go-binance/endpoints"
)

type BinanceClient struct {
	client     *client.Client
	Deposit    *endpoints.DepositService
	Withdrawal *endpoints.WithdrawalService
	Account    *endpoints.AccountService
}

func NewClient(apiKey, secretKey string) *BinanceClient {
	c := client.NewClient(apiKey, secretKey)
	
	return &BinanceClient{
		client:     c,
		Deposit:    endpoints.NewDepositService(c),
		Withdrawal: endpoints.NewWithdrawalService(c),
		Account:    endpoints.NewAccountService(c),
	}
}

func (b *BinanceClient) SetBaseURL(url string) {
	b.client.SetBaseURL(url)
}