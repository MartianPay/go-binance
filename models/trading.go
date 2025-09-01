package models

import "time"

// OrderSide represents BUY or SELL
type OrderSide string

const (
	SideBuy  OrderSide = "BUY"
	SideSell OrderSide = "SELL"
)

// OrderType represents different order types
type OrderType string

const (
	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeStopLoss        OrderType = "STOP_LOSS"
	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
)

// TimeInForce represents time in force options
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Cancel
	TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill or Kill
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusPendingCancel   OrderStatus = "PENDING_CANCEL"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
	OrderStatusExpiredInMatch  OrderStatus = "EXPIRED_IN_MATCH"
)

// OrderResponseType represents response type for orders
type OrderResponseType string

const (
	OrderResponseTypeACK    OrderResponseType = "ACK"
	OrderResponseTypeRESULT OrderResponseType = "RESULT"
	OrderResponseTypeFULL   OrderResponseType = "FULL"
)

// NewOrderRequest represents a new order request
type NewOrderRequest struct {
	Symbol           string            `json:"symbol"`
	Side             OrderSide         `json:"side"`
	Type             OrderType         `json:"type"`
	TimeInForce      TimeInForce       `json:"timeInForce,omitempty"`
	Quantity         string            `json:"quantity,omitempty"`
	QuoteOrderQty    string            `json:"quoteOrderQty,omitempty"`
	Price            string            `json:"price,omitempty"`
	NewClientOrderId string            `json:"newClientOrderId,omitempty"`
	StopPrice        string            `json:"stopPrice,omitempty"`
	IcebergQty       string            `json:"icebergQty,omitempty"`
	NewOrderRespType OrderResponseType `json:"newOrderRespType,omitempty"`
	RecvWindow       int64             `json:"recvWindow,omitempty"`
}

// OrderResponse represents the response from placing an order
type OrderResponse struct {
	Symbol              string      `json:"symbol"`
	OrderId             int64       `json:"orderId"`
	OrderListId         int64       `json:"orderListId"`
	ClientOrderId       string      `json:"clientOrderId"`
	TransactTime        int64       `json:"transactTime"`
	Price               string      `json:"price"`
	OrigQty             string      `json:"origQty"`
	ExecutedQty         string      `json:"executedQty"`
	CummulativeQuoteQty string      `json:"cummulativeQuoteQty"`
	Status              OrderStatus `json:"status"`
	TimeInForce         TimeInForce `json:"timeInForce"`
	Type                OrderType   `json:"type"`
	Side                OrderSide   `json:"side"`
	WorkingTime         int64       `json:"workingTime"`
	Fills               []OrderFill `json:"fills,omitempty"`
}

// OrderFill represents a fill in an order
type OrderFill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	TradeId         int64  `json:"tradeId"`
}

// QueryOrderRequest represents a query order request
type QueryOrderRequest struct {
	Symbol            string `json:"symbol"`
	OrderId           int64  `json:"orderId,omitempty"`
	OrigClientOrderId string `json:"origClientOrderId,omitempty"`
	RecvWindow        int64  `json:"recvWindow,omitempty"`
}

// Order represents an order with full details
type Order struct {
	Symbol              string      `json:"symbol"`
	OrderId             int64       `json:"orderId"`
	OrderListId         int64       `json:"orderListId"`
	ClientOrderId       string      `json:"clientOrderId"`
	Price               string      `json:"price"`
	OrigQty             string      `json:"origQty"`
	ExecutedQty         string      `json:"executedQty"`
	CummulativeQuoteQty string      `json:"cummulativeQuoteQty"`
	Status              OrderStatus `json:"status"`
	TimeInForce         TimeInForce `json:"timeInForce"`
	Type                OrderType   `json:"type"`
	Side                OrderSide   `json:"side"`
	StopPrice           string      `json:"stopPrice"`
	IcebergQty          string      `json:"icebergQty"`
	Time                int64       `json:"time"`
	UpdateTime          int64       `json:"updateTime"`
	IsWorking           bool        `json:"isWorking"`
	WorkingTime         int64       `json:"workingTime"`
	OrigQuoteOrderQty   string      `json:"origQuoteOrderQty"`
}

// CancelOrderRequest represents a cancel order request
type CancelOrderRequest struct {
	Symbol            string `json:"symbol"`
	OrderId           int64  `json:"orderId,omitempty"`
	OrigClientOrderId string `json:"origClientOrderId,omitempty"`
	NewClientOrderId  string `json:"newClientOrderId,omitempty"`
	RecvWindow        int64  `json:"recvWindow,omitempty"`
}

// CancelOrderResponse represents cancel order response
type CancelOrderResponse struct {
	Symbol              string      `json:"symbol"`
	OrigClientOrderId   string      `json:"origClientOrderId"`
	OrderId             int64       `json:"orderId"`
	OrderListId         int64       `json:"orderListId"`
	ClientOrderId       string      `json:"clientOrderId"`
	TransactTime        int64       `json:"transactTime"`
	Price               string      `json:"price"`
	OrigQty             string      `json:"origQty"`
	ExecutedQty         string      `json:"executedQty"`
	CummulativeQuoteQty string      `json:"cummulativeQuoteQty"`
	Status              OrderStatus `json:"status"`
	TimeInForce         TimeInForce `json:"timeInForce"`
	Type                OrderType   `json:"type"`
	Side                OrderSide   `json:"side"`
}

// OpenOrdersRequest represents open orders request
type OpenOrdersRequest struct {
	Symbol     string `json:"symbol,omitempty"`
	RecvWindow int64  `json:"recvWindow,omitempty"`
}

// AllOrdersRequest represents all orders request
type AllOrdersRequest struct {
	Symbol     string    `json:"symbol"`
	OrderId    int64     `json:"orderId,omitempty"`
	StartTime  time.Time `json:"startTime,omitempty"`
	EndTime    time.Time `json:"endTime,omitempty"`
	Limit      int       `json:"limit,omitempty"` // Default 500; max 1000
	RecvWindow int64     `json:"recvWindow,omitempty"`
}

// TradingAccountInfoRequest represents account information request
type TradingAccountInfoRequest struct {
	RecvWindow int64 `json:"recvWindow,omitempty"`
}

// TradingAccountInfo represents trading account information
type TradingAccountInfo struct {
	MakerCommission  int64     `json:"makerCommission"`
	TakerCommission  int64     `json:"takerCommission"`
	BuyerCommission  int64     `json:"buyerCommission"`
	SellerCommission int64     `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       int64     `json:"updateTime"`
	AccountType      string    `json:"accountType"`
	Balances         []Balance `json:"balances"`
	Permissions      []string  `json:"permissions"`
}

// Balance represents account balance
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// MyTradesRequest represents account trade list request
type MyTradesRequest struct {
	Symbol     string    `json:"symbol"`
	OrderId    int64     `json:"orderId,omitempty"`
	StartTime  time.Time `json:"startTime,omitempty"`
	EndTime    time.Time `json:"endTime,omitempty"`
	FromId     int64     `json:"fromId,omitempty"`
	Limit      int       `json:"limit,omitempty"` // Default 500; max 1000
	RecvWindow int64     `json:"recvWindow,omitempty"`
}

// Trade represents a trade
type Trade struct {
	Symbol          string `json:"symbol"`
	Id              int64  `json:"id"`
	OrderId         int64  `json:"orderId"`
	OrderListId     int64  `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}

// ExchangeInfoRequest represents exchange info request
type ExchangeInfoRequest struct {
	Symbol     string   `json:"symbol,omitempty"`     // Single symbol
	Symbols    []string `json:"symbols,omitempty"`    // Multiple symbols
	Permissions []string `json:"permissions,omitempty"` // Filter by permissions
}

// ExchangeInfo represents exchange information
type ExchangeInfo struct {
	Timezone        string           `json:"timezone"`
	ServerTime      int64            `json:"serverTime"`
	RateLimits      []RateLimit      `json:"rateLimits"`
	ExchangeFilters []interface{}   `json:"exchangeFilters"`
	Symbols         []SymbolInfo     `json:"symbols"`
}

// RateLimit represents rate limit info
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// SymbolInfo represents symbol trading rules and information
type SymbolInfo struct {
	Symbol                     string                   `json:"symbol"`
	Status                     string                   `json:"status"`
	BaseAsset                  string                   `json:"baseAsset"`
	BaseAssetPrecision         int                      `json:"baseAssetPrecision"`
	QuoteAsset                 string                   `json:"quoteAsset"`
	QuotePrecision             int                      `json:"quotePrecision"`
	QuoteAssetPrecision        int                      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int                      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int                      `json:"quoteCommissionPrecision"`
	OrderTypes                 []string                 `json:"orderTypes"`
	IcebergAllowed             bool                     `json:"icebergAllowed"`
	OcoAllowed                 bool                     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool                     `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop          bool                     `json:"allowTrailingStop"`
	CancelReplaceAllowed       bool                     `json:"cancelReplaceAllowed"`
	IsSpotTradingAllowed       bool                     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool                     `json:"isMarginTradingAllowed"`
	Filters                    []map[string]interface{} `json:"filters"`
	Permissions                []string                 `json:"permissions"`
	DefaultSelfTradePreventionMode string              `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string          `json:"allowedSelfTradePreventionModes"`
}

// Common filter types for symbol
type FilterType string

const (
	FilterTypePriceFilter       FilterType = "PRICE_FILTER"
	FilterTypePercentPrice      FilterType = "PERCENT_PRICE"
	FilterTypeLotSize           FilterType = "LOT_SIZE"
	FilterTypeMinNotional       FilterType = "MIN_NOTIONAL"
	FilterTypeNotional          FilterType = "NOTIONAL"
	FilterTypeIcebergParts      FilterType = "ICEBERG_PARTS"
	FilterTypeMarketLotSize     FilterType = "MARKET_LOT_SIZE"
	FilterTypeMaxNumOrders      FilterType = "MAX_NUM_ORDERS"
	FilterTypeMaxNumAlgoOrders  FilterType = "MAX_NUM_ALGO_ORDERS"
	FilterTypeMaxNumIcebergOrders FilterType = "MAX_NUM_ICEBERG_ORDERS"
	FilterTypeMaxPosition       FilterType = "MAX_POSITION"
	FilterTypeTrailingDelta     FilterType = "TRAILING_DELTA"
)

// PriceFilter represents PRICE_FILTER
type PriceFilter struct {
	FilterType string `json:"filterType"`
	MinPrice   string `json:"minPrice"`
	MaxPrice   string `json:"maxPrice"`
	TickSize   string `json:"tickSize"`
}

// LotSizeFilter represents LOT_SIZE filter
type LotSizeFilter struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty"`
	MaxQty     string `json:"maxQty"`
	StepSize   string `json:"stepSize"`
}

// MinNotionalFilter represents MIN_NOTIONAL filter
type MinNotionalFilter struct {
	FilterType        string `json:"filterType"`
	MinNotional       string `json:"minNotional"`
	ApplyToMarket     bool   `json:"applyToMarket"`
	AvgPriceMins      int    `json:"avgPriceMins"`
}

// NotionalFilter represents NOTIONAL filter
type NotionalFilter struct {
	FilterType        string `json:"filterType"`
	MinNotional       string `json:"minNotional"`
	ApplyMinToMarket  bool   `json:"applyMinToMarket"`
	MaxNotional       string `json:"maxNotional"`
	ApplyMaxToMarket  bool   `json:"applyMaxToMarket"`
	AvgPriceMins      int    `json:"avgPriceMins"`
}