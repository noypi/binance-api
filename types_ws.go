package binance

type UpdateType string

const (
	UpdateTypeDepth  UpdateType = "depthUpdate"
	UpdateTypeKline  UpdateType = "kline"
	UpdateTypeTrades UpdateType = "aggTrade"

	UpdateTypeOutboundAccountInfo UpdateType = "outboundAccountInfo"
	UpdateTypeExecutionReport     UpdateType = "executionReport"
)

// DepthUpdate represents the incoming messages for depth websocket updates
type DepthUpdate struct {
	EventType UpdateType  `json:"e"` // EventType represents the update type
	Time      uint64      `json:"E"` // Time represents the event time
	Symbol    string      `json:"s"` // Symbol represents the symbol related to the update
	UpdateID  int         `json:"u"` // UpdateID to sync up with updateid in /api/v1/depth
	Bids      []DepthElem `json:"b"` // Bids is a list of bids for symbol
	Asks      []DepthElem `json:"a"` // Asks is a list of asks for symbol
}

// KlinesUpdate represents the incoming messages for klines websocket updates
type KlinesUpdate struct {
	EventType UpdateType `json:"e"` // EventType represents the update type
	Time      uint64     `json:"E"` // Time represents the event time
	Symbol    string     `json:"s"` // Symbol represents the symbol related to the update
	Kline     struct {
		StartTime    uint64        `json:"t"` // StartTime is the start time of this bar
		EndTime      uint64        `json:"T"` // EndTime is the end time of this bar
		Symbol       string        `json:"s"` // Symbol represents the symbol related to this kline
		Interval     KlineInterval `json:"i"` // Interval is the kline interval
		FirstTradeID int           `json:"f"` // FirstTradeID is the first trade ID
		LastTradeID  int           `json:"L"` // LastTradeID is the first trade ID

		OpenPrice            string `json:"o"` // OpenPrice represents the open price for this bar
		ClosePrice           string `json:"c"` // ClosePrice represents the close price for this bar
		High                 string `json:"h"` // High represents the highest price for this bar
		Low                  string `json:"l"` // Low represents the lowest price for this bar
		Volume               string `json:"v"` // Volume is the trades volume for this bar
		Trades               int    `json:"n"` // Trades is the number of conducted trades
		Final                bool   `json:"x"` // Final indicates whether this bar is final or yet may receive updates
		VolumeQuote          string `json:"q"` // VolumeQuote indicates the quote volume for the symbol
		VolumeActiveBuy      string `json:"V"` // VolumeActiveBuy represents the volume of active buy
		VolumeQuoteActiveBuy string `json:"Q"` // VolumeQuoteActiveBuy represents the quote volume of active buy
	} `json:"k"` // Kline is the kline update
}

// TradesUpdate represents the incoming messages for aggregated trades websocket updates
type TradesUpdate struct {
	EventType             UpdateType `json:"e"` // EventType represents the update type
	Time                  uint64     `json:"E"` // Time represents the event time
	Symbol                string     `json:"s"` // Symbol represents the symbol related to the update
	TradeID               int        `json:"a"` // TradeID is the aggregated trade ID
	Price                 string     `json:"p"` // Price is the trade price
	Quantity              string     `json:"q"` // Quantity is the trade quantity
	FirstBreakDownTradeID int        `json:"f"` // FirstBreakDownTradeID is the first breakdown trade ID
	LastBreakDownTradeID  int        `json:"l"` // LastBreakDownTradeID is the last breakdown trade ID
	TradeTime             uint64     `json:"T"` // Time is the trade time
	Maker                 bool       `json:"m"` // Maker indicates whether buyer is a maker
}

// AccountUpdate represents the incoming messages for account info websocket updates
type AccountUpdate struct {
	EventType        UpdateType `json:"e"` // EventType represents the update type
	Time             uint64     `json:"E"` // Time represents the event time
	MakerCommission  int        `json:"m"` // MakerCommission is the maker commission for the account
	TakerCommission  int        `json:"t"` // TakerCommission is the taker commission for the account
	BuyerCommission  int        `json:"b"` // BuyerCommission is the buyer commission for the account
	SellerCommission int        `json:"s"` // SellerCommission is the seller commission for the account
	CanTrade         bool       `json:"T"`
	CanWithdraw      bool       `json:"W"`
	CanDeposit       bool       `json:"D"`
	Balances         []*struct {
		Asset  string `json:"a"`
		Free   string `json:"f"`
		Locked string `json:"l"`
	} `json:"B"`
}

// OrderUpdate represents the incoming messages for account orders websocket updates
type OrderUpdate struct {
	EventType        UpdateType   `json:"e"` // EventType represents the update type
	Time             uint64       `json:"E"` // Time represents the event time
	Symbol           string       `json:"s"` // Symbol represents the symbol related to the update
	NewClientOrderID string       `json:"c"` // NewClientOrderID is the new client order ID
	Side             OrderSide    `json:"S"` // Side is the order side
	OrderType        OrderType    `json:"o"` // OrderType represents the order type
	TimeInForce      TimeInForce  `json:"f"` // TimeInForce represents the order TIF type
	OrigQty          string       `json:"q"` // OrigQty represents the order original quantity
	Price            string       `json:"p"` // Price is the order price
	ExecutionType    OrderStatus  `json:"x"` // ExecutionType represents the execution type for the order
	Status           OrderStatus  `json:"X"` // Status represents the order status for the order
	Error            OrderFailure `json:"r"` // Error represents an order rejection reason
	OrderID          int          `json:"i"` // OrderID represents the order ID
	OrderTime        uint64       `json:"T"` // OrderTime represents the order time
	FilledQty        string       `json:"l"` // FilledQty represents the quantity of the last filled trade
	FilledPrice      string       `json:"L"` // FilledPrice is the price of last filled trade
	TotalFilledQty   string       `json:"z"` // TotalFilledQty is the accumulated quantity of filled trades on this order
	Commission       string       `json:"n"` // Commission is the commission for the trade
	CommissionAsset  string       `json:"N"` // CommissionAsset is the asset on which commission is taken
	TradeTime        uint64       `json:"T"` // TradeTime is the trade time
	TradeID          int          `json:"t"` // TradeID represents the trade ID
	Maker            bool         `json:"m"` // Maker represents whether buyer is maker or not
}
