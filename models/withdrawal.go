package models

import "time"

type WithdrawalRequest struct {
	Coin              string  `json:"coin"`
	Network           string  `json:"network,omitempty"`
	Address           string  `json:"address"`
	AddressTag        string  `json:"addressTag,omitempty"`
	Amount            string  `json:"amount"`
	WithdrawOrderId   string  `json:"withdrawOrderId,omitempty"`
	TransactionFeeFlag bool    `json:"transactionFeeFlag,omitempty"`
	Name              string  `json:"name,omitempty"`
	WalletType        int     `json:"walletType,omitempty"`
	RecvWindow        int64   `json:"recvWindow,omitempty"`
}

type WithdrawalResponse struct {
	Id string `json:"id"`
}

type WithdrawalHistory struct {
	Id              string    `json:"id"`
	Amount          string    `json:"amount"`
	TransactionFee  string    `json:"transactionFee"`
	Coin            string    `json:"coin"`
	Status          int       `json:"status"`
	Address         string    `json:"address"`
	TxId            string    `json:"txId"`
	ApplyTime       string    `json:"applyTime"`
	Network         string    `json:"network"`
	TransferType    int       `json:"transferType"`
	WithdrawOrderId string    `json:"withdrawOrderId,omitempty"`
	Info            string    `json:"info,omitempty"`
	ConfirmNo       int       `json:"confirmNo"`
	WalletType      int       `json:"walletType"`
	TxKey           string    `json:"txKey,omitempty"`
	CompleteTime    string    `json:"completeTime,omitempty"`
}

type WithdrawalHistoryRequest struct {
	Coin            string    `json:"coin,omitempty"`
	WithdrawOrderId string    `json:"withdrawOrderId,omitempty"`
	Status          int       `json:"status,omitempty"`
	StartTime       time.Time `json:"startTime,omitempty"`
	EndTime         time.Time `json:"endTime,omitempty"`
	Offset          int       `json:"offset,omitempty"`
	Limit           int       `json:"limit,omitempty"`
	IdList          string    `json:"idList,omitempty"` // 提现ID列表，以逗号分隔，最大支持45个
	RecvWindow      int64     `json:"recvWindow,omitempty"`
}

// WithdrawalQuota - 24小时提现限额信息
type WithdrawalQuota struct {
	WdQuota     string `json:"wdQuota"`     // 24小时总提现限额 (USD)
	UsedWdQuota string `json:"usedWdQuota"` // 24小时已使用限额 (USD)
}

// WithdrawalAddress - 提现地址信息
type WithdrawalAddress struct {
	Address      string `json:"address"`
	Coin         string `json:"coin"`
	Name         string `json:"name"`
	Network      string `json:"network"`
	OriginType   string `json:"originType"`
	WhiteStatus  bool   `json:"whiteStatus"`
}