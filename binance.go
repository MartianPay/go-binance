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
	Market     *endpoints.MarketDataService
	Trading    *endpoints.TradingService
}

func NewClient(apiKey, secretKey string) *BinanceClient {
	c := client.NewClient(apiKey, secretKey)
	
	return &BinanceClient{
		client:     c,
		Deposit:    endpoints.NewDepositService(c),
		Withdrawal: endpoints.NewWithdrawalService(c),
		Account:    endpoints.NewAccountService(c),
		Market:     endpoints.NewMarketDataService(c),
		Trading:    endpoints.NewTradingService(c),
	}
}

func (b *BinanceClient) SetBaseURL(url string) {
	b.client.SetBaseURL(url)
}