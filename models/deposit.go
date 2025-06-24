package models

import "time"

type DepositAddress struct {
	Address string `json:"address"`
	Coin    string `json:"coin"`
	Tag     string `json:"tag,omitempty"`
	URL     string `json:"url,omitempty"`
}

type DepositHistory struct {
	Id              string    `json:"id"`
	Amount          string    `json:"amount"`
	Coin            string    `json:"coin"`
	Network         string    `json:"network"`
	Status          int       `json:"status"`
	Address         string    `json:"address"`
	AddressTag      string    `json:"addressTag,omitempty"`
	TxId            string    `json:"txId"`
	InsertTime      int64     `json:"insertTime"`
	TransferType    int       `json:"transferType"`
	UnlockConfirm   int       `json:"unlockConfirm,omitempty"`
	ConfirmTimes    string    `json:"confirmTimes"`
	WalletType      int       `json:"walletType"`
}

type DepositAddressRequest struct {
	Coin       string `json:"coin"`
	Network    string `json:"network,omitempty"`
	RecvWindow int64  `json:"recvWindow,omitempty"`
}

type DepositHistoryRequest struct {
	Coin       string    `json:"coin,omitempty"`
	Status     int       `json:"status,omitempty"`
	StartTime  time.Time `json:"startTime,omitempty"`
	EndTime    time.Time `json:"endTime,omitempty"`
	Offset     int       `json:"offset,omitempty"`
	Limit      int       `json:"limit,omitempty"`
	RecvWindow int64     `json:"recvWindow,omitempty"`
	TxId       string    `json:"txId,omitempty"`
}