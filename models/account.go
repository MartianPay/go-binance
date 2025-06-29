package models

type CoinInfo struct {
	Coin              string       `json:"coin"`
	DepositAllEnable  bool         `json:"depositAllEnable"`
	Free              string       `json:"free"`
	Freeze            string       `json:"freeze"`
	Ipoable           string       `json:"ipoable"`
	Ipoing            string       `json:"ipoing"`
	IsLegalMoney      bool         `json:"isLegalMoney"`
	Locked            string       `json:"locked"`
	Name              string       `json:"name"`
	Storage           string       `json:"storage"`
	Trading           bool         `json:"trading"`
	WithdrawAllEnable bool         `json:"withdrawAllEnable"`
	Withdrawing       string       `json:"withdrawing"`
	NetworkList       []NetworkInfo `json:"networkList"`
}

type NetworkInfo struct {
	AddressRegex            string `json:"addressRegex"`
	Coin                    string `json:"coin"`
	DepositDesc             string `json:"depositDesc,omitempty"`
	DepositEnable           bool   `json:"depositEnable"`
	IsDefault               bool   `json:"isDefault"`
	MemoRegex               string `json:"memoRegex"`
	MinConfirm              int    `json:"minConfirm"`
	Name                    string `json:"name"`
	Network                 string `json:"network"`
	ResetAddressStatus      bool   `json:"resetAddressStatus"`
	SpecialTips             string `json:"specialTips,omitempty"`
	UnLockConfirm           int    `json:"unLockConfirm"`
	WithdrawDesc            string `json:"withdrawDesc,omitempty"`
	WithdrawEnable          bool   `json:"withdrawEnable"`
	WithdrawFee             string `json:"withdrawFee"`
	WithdrawIntegerMultiple string `json:"withdrawIntegerMultiple"`
	WithdrawMax             string `json:"withdrawMax"`
	WithdrawMin             string `json:"withdrawMin"`
	SameAddress             bool   `json:"sameAddress"`
	EstimatedArrivalTime    int    `json:"estimatedArrivalTime"`
	Busy                    bool   `json:"busy"`
}

type AccountInfo struct {
	VipLevel                  int  `json:"vipLevel"`
	IsMarginEnable            bool `json:"isMarginEnable"`
	IsFutureEnable            bool `json:"isFutureEnable"`
}

type AssetTransferRequest struct {
	Type       string `json:"type"`
	Asset      string `json:"asset"`
	Amount     string `json:"amount"`
	FromSymbol string `json:"fromSymbol,omitempty"`
	ToSymbol   string `json:"toSymbol,omitempty"`
	RecvWindow int64  `json:"recvWindow,omitempty"`
}

type AssetTransferResponse struct {
	TranId int64 `json:"tranId"`
}

type UserAsset struct {
	Asset            string `json:"asset"`
	Free             string `json:"free"`
	Locked           string `json:"locked"`
	Freeze           string `json:"freeze"`
	Withdrawing      string `json:"withdrawing"`
	BtcValuation     string `json:"btcValuation"`
}

type UserAssetRequest struct {
	Asset      string `json:"asset,omitempty"`
	NeedBtcValuation bool   `json:"needBtcValuation,omitempty"`
	RecvWindow int64  `json:"recvWindow,omitempty"`
}