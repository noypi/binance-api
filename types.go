package binance

import (
	"fmt"
	"strconv"
	"strings"
)

// OrderType represents the order type
type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

type OrderStatus string

const (
	OrderStatusNew      OrderStatus = "NEW"
	OrderStatusPartial  OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled   OrderStatus = "FILLED"
	OrderStatusCanceled OrderStatus = "CANCELED"
	OrderStatusPending  OrderStatus = "PENDING_CANCEL"
	OrderStatusRejected OrderStatus = "REJECTED"
	OrderStatusExpired  OrderStatus = "EXPIRED"
	OrderStatusReplaced OrderStatus = "REPLACED"
	OrderStatusTrade    OrderStatus = "TRADE"
)

type OrderFailure string

const (
	OrderFailureNone              OrderFailure = "NONE"
	OrderFailureUnknownInstrument OrderFailure = "UNKNOWN_INSTRUMENT"
	OrderFailureMarketClosed      OrderFailure = "MARKET_CLOSED"
	OrderFailurePriceExceed       OrderFailure = "PRICE_QTY_EXCEED_HARD_LIMITS"
	OrderFailureUnknownOrder      OrderFailure = "UNKNOWN_ORDER"
	OrderFailureDuplicate         OrderFailure = "DUPLICATE_ORDER"
	OrderFailureUnknownAccount    OrderFailure = "UNKNOWN_ACCOUNT"
	OrderFailureInsufficientFunds OrderFailure = "INSUFFICIENT_BALANCE"
	OrderFailureAccountInaactive  OrderFailure = "ACCOUNT_INACTIVE"
	OrderFailureAccountSettle     OrderFailure = "ACCOUNT_CANNOT_SETTLE"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Cancel
	TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
)

type NewOrderOpts struct {
	Symbol           string      `url:"symbol"`
	Side             OrderSide   `url:"side"`
	Type             OrderType   `url:"type"`
	TimeInForce      TimeInForce `url:"timeInForce,omitempty"`
	Quantity         string      `url:"quantity"`
	Price            string      `url:"price"`
	NewClientOrderId string      `url:"newClientOrderId,omitempty"`
	StopPrice        string      `url:"stopPrice,omitempty"`
	IcebergQty       string      `url:"icebergQty,omitempty"`
}

type NewOrder struct {
	Symbol            string `json:"symbol"`
	OrderID           int    `json:"orderId"`
	OrigClientOrderID string `json:"origClientOrderId"`
	TransactTime      uint64 `json:"transactTime"`
}

type ServerTime struct {
	ServerTime uint64 `json:"serverTime"`
}

type KlineInterval string

const (
	KlineInterval1m  KlineInterval = "1m"
	KlineInterval3m  KlineInterval = "3m"
	KlineInterval5m  KlineInterval = "5m"
	KlineInterval15m KlineInterval = "15m"
	KlineInterval30m KlineInterval = "30m"
	KlineInterval1h  KlineInterval = "1h"
	KlineInterval2h  KlineInterval = "2h"
	KlineInterval4h  KlineInterval = "4h"
	KlineInterval6h  KlineInterval = "6h"
	KlineInterval8h  KlineInterval = "8h"
	KlineInterval12h KlineInterval = "12h"
	KlineInterval1d  KlineInterval = "1d"
	KlineInterval3d  KlineInterval = "3d"
	KlineInterval1w  KlineInterval = "1w"
	KlineInterval1M  KlineInterval = "1M"
)

// DepthOpts are used to specify symbol to retrieve order book for
type DepthOpts struct {
	Symbol string `url:"symbol"` // Symbol is the symbol to fetch data for
	Limit  int    `url:"limit"`  // Limit is the number of order book items to retrieve. Max 100
}

// DepthElem represents a specific order in the order book
type DepthElem struct {
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

// UnmarshalJSON unmarshal the given depth raw data and converts to depth struct
func (b *DepthElem) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("UnmarshalJSON on nil pointer")
	}

	if len(data) == 0 {
		return nil
	}
	s := strings.Replace(string(data), `"`, "", -1)
	s = strings.Trim(s, "[]")
	tokens := strings.Split(s, ",")
	if len(tokens) < 2 {
		return fmt.Errorf("at least two fields are expected but got: %v", tokens)
	}
	b.Price = tokens[0]
	b.Quantity = tokens[1]
	return nil
}

type Depth struct {
	LastUpdateID int         `json:"lastUpdateId"`
	Bids         []DepthElem `json:"bids"`
	Asks         []DepthElem `json:"asks"`
}

type KlinesOpts struct {
	Symbol    string        `url:"symbol"`   // Symbol is the symbol to fetch data for
	Interval  KlineInterval `url:"interval"` // Interval is the interval for each kline/candlestick
	Limit     int           `url:"limit"`    // Limit is the maximal number of elements to receive. Max 500
	StartTime uint64        `url:"startTime,omitempty"`
	EndTime   uint64        `url:"endTime,omitempty"`
}

type Klines struct {
	OpenTime                 uint64
	OpenPrice                string
	High                     string
	Low                      string
	ClosePrice               string
	Volume                   string
	CloseTime                uint64
	QuoteAssetVolume         string
	Trades                   int
	TakerBuyBaseAssetVolume  string
	TakerBuyQuoteAssetVolume string
}

// UnmarshalJSON unmarshal the given depth raw data and converts to depth struct
func (b *Klines) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("UnmarshalJSON on nil pointer")
	}

	if len(data) == 0 {
		return nil
	}
	s := strings.Replace(string(data), `"`, "", -1)
	s = strings.Trim(s, "[]")
	tokens := strings.Split(s, ",")
	if len(tokens) < 11 {
		return fmt.Errorf("at least 11 fields are expected but got: %v", tokens)
	}
	var err error
	u, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse open time: %v", tokens[0])
	}
	b.OpenTime = uint64(u)
	b.OpenPrice = tokens[1]
	b.High = tokens[2]
	b.Low = tokens[3]
	b.ClosePrice = tokens[4]
	b.Volume = tokens[5]
	u, err = strconv.ParseInt(tokens[6], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse close time: %v", tokens[6])
	}
	b.CloseTime = uint64(u)
	b.QuoteAssetVolume = tokens[7]
	u, err = strconv.ParseInt(tokens[8], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse trades: %v", tokens[8])
	}
	b.Trades = int(u)
	b.TakerBuyBaseAssetVolume = tokens[9]
	b.TakerBuyQuoteAssetVolume = tokens[10]
	return nil
}

type BookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}

// TickerOpts represents the opts for a specified ticker
type TickerOpts struct {
	Symbol string `url:"symbol"`
}

// TickerStats is the stats for a specific symbol
type TickerStats struct {
	PriceChange           string `json:"priceChange"`
	PriceChangePercentage string `json:"priceChangePercent"`
	WeightedAvgPrice      string `json:"weightedAvgPrice"`
	PrevClosePrice        string `json:"prevClosePrice"`
	LastPrice             string `json:"lastPrice"`
	BidPrice              string `json:"bidPrice"`
	AskPrice              string `json:"askPrice"`
	OpenPrice             string `json:"openPrice"`
	HighPrice             string `json:"highPrice"` // HighPrice is 24hr high price
	LowPrice              string `json:"lowPrice"`  // LowPrice is 24hr low price
	Volume                string `json:"volume"`
	OpenTime              uint64 `json:"openTime"`
	CloseTime             uint64 `json:"closeTime"`
	FirstID               int    `json:"fristId"`
	LastID                int    `json:"lastId"`
	Count                 int    `json:"count"`
}

type SymbolPrice struct {
	Symbol string
	Price  string
}

type AllPrices struct {
	Prices []*SymbolPrice
}

// QueryOrderOpts represents the opts for querying an order
// Remark: Either OrderID or OrigOrderiD must be set
type QueryOrderOpts struct {
	Symbol            string `url:"symbol"`
	OrderID           int    `url:"orderId,omitempty"`
	OrigClientOrderId string `url:"origClientOrderId,omitempty"`
}

type QueryOrder struct {
	Symbol        string      `json:"symbol"`
	OrderID       int         `json:"orderId"`
	ClientOrderID string      `json:"clientOrderId"`
	Price         string      `json:"price"`
	OrigQty       string      `json:"origQty"`
	ExecutedQty   string      `json:"executedQty"`
	Status        OrderStatus `json:"status"`
	TimeInForce   TimeInForce `json:"timeInForce"`
	Type          OrderType   `json:"type"`
	Side          OrderSide   `json:"side"`
	StopPrice     string      `json:"stopPrice"`
	IcebergQty    string      `json:"IcebergQty"`
	Time          uint64      `json:"time"`
}

// Remark: Either OrderID or OrigOrderiD must be set
type CancelOrderOpts struct {
	Symbol            string `url:"symbol"`
	OrderID           int    `url:"orderId"`
	OrigClientOrderId string `url:"origClientOrderId,omitempty"`
	NewClientOrderId  string `url:"newClientOrderId,omitempty"`
}

type CancelOrder struct {
	Symbol            string `json:"symbol"`
	OrderID           int    `json:"orderId"`
	OrigClientOrderId string `json:"origClientOrderId"`
	ClientOrderId     string `json:"clientOrderId"`
}

type OpenOrdersOpts struct {
	Symbol string `url:"symbol"`
}

// AllOrdersOpts represents the opts used for querying orders of the given symbol
// Remark: If orderId is set, it will get orders >= that orderId. Otherwise most recent orders are returned
type AllOrdersOpts struct {
	Symbol  string `url:"symbol"`  // Symbol is the symbol to fetch orders for
	OrderID int    `url:"orderId"` // OrderID, if set, will filter all recent orders newer from the given ID
	Limit   int    `url:"limit"`   // Limit is the maximal number of elements to receive. Max 500
}

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type AccountInfo struct {
	MakerCommission  int        `json:"makerCommission"`
	TakerCommission  int        `json:"takerCommission"`
	BuyerCommission  int        `json:"buyerCommission"`
	SellerCommission int        `json:"sellerCommission"`
	CanTrade         bool       `json:"canTrade"`
	CanWithdraw      bool       `json:"canWithdraw"`
	CanDeposit       bool       `json:"canDeposit"`
	Balances         []*Balance `json:"balances"`
}

type TradesOpts struct {
	Symbol string `url:"symbol"`
	Limit  int    `url:"limit"`  // Limit is the maximal number of elements to receive. Max 500
	FromID int    `url:"fromId"` // FromID is trade ID to fetch from. Default gets most recent trades
}

type Trades struct {
	ID              int    `json:"id"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            uint64 `json:"time"`
	Buyer           bool   `json:"isBuyer"`
	Maker           bool   `json:"isMaker"`
	BestMatch       bool   `json:"isBestMatch"`
}

type Datastream struct {
	ListenKey string `json:"listenKey" url:"listenKey"`
}

type AggregatedTradeOpts struct {
	Symbol    string `url:"symbol"` // Symbol is the symbol to fetch data for
	FromID    int    `url:"fromId"` // Interval is the interval for each kline/candlestick
	Limit     int    `url:"limit"`  // Limit is the maximal number of elements to receive. Max 500
	StartTime uint64 `url:"startTime,omitempty"`
	EndTime   uint64 `url:"endTime,omitempty"`
}
type AggregatedTrade struct {
	TradeID      int    `json:"a"` // TradeID is the aggregate trade ID
	Price        string `json:"p"` // Price is the trade price
	Quantity     string `json:"q"` // Quantity is the trade quantity
	FirstTradeID int    `json:"f"`
	LastTradeID  int    `json:"l"`
	Time         uint64 `json:"T"`
	Maker        bool   `json:"m"` // Maker indicates if the buyer is the maker
	BestMatch    bool   `json:"M"` // BestMatch indicates if the trade was at the best price match
}

type ExchangeInfo struct {
	Symbols []SymbolInfo
}

type SymbolInfo struct {
	Symbol              string             `json:"symbol"`
	Status              SymbolStatus       `json:"status"`
	BaseAsset           string             `json:"baseAsset"`
	BaseAssetPrecision  int                `json:"baseAssetPrecision"`
	QuoteAsset          string             `json:"quoteAsset"`
	QuoteAssetPrecision int                `json:"quoteAssetPrecision"`
	OrderTypes          []OrderType        `json:"orderTypes"`
	Iceberg             bool               `json:"icebergAllowed"`
	Filters             []SymbolInfoFilter `json:"filters"`
}

type SymbolStatus string

const (
	SymbolStatusTrading SymbolStatus = "TRADING"
)

type FilterType string

const (
	FilterTypePrice       FilterType = "PRICE_FILTER"
	FilterTypeLotSize     FilterType = "LOT_SIZE"
	FilterTypeMinNotional FilterType = "MIN_NOTIONAL"
)

type SymbolInfoFilter struct {
	Type FilterType `json:"filterType"`

	// PRICE_FILTER paramters
	MinPrice string `json:"minPrice"`
	MaxPrice string `json:"maxPrice"`
	TickSize string `json:"tickSize"`

	// LOT_SIZE parameters
	MinQty   string `json:"minQty"`
	MaxQty   string `json:"maxQty"`
	StepSize string `json:"stepSize"`

	// MIN_NOTIONAL paramters
	MinNotional string `json:"minNotional"`
}
