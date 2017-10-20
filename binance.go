package binance

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

type BinanceClient struct {
	client *client
	dialer *websocket.Dialer
}

func NewBinanceClient(apikey, secret string) *BinanceClient {
	return &BinanceClient{
		client: &client{
			window: 5000,
			apikey: apikey,
			secret: secret,
			client: http.DefaultClient,
		},
		dialer: websocket.DefaultDialer,
	}
}

func NewBinanceClientWindow(apikey, secret string, window int) (*BinanceClient, error) {
	if window <= 0 {
		return nil, fmt.Errorf("window value is invalid")
	}
	return &BinanceClient{
		client: &client{
			window: window,
			apikey: apikey,
			secret: secret,
			client: http.DefaultClient,
		},
		dialer: websocket.DefaultDialer,
	}, nil
}

// General endpoints

// Ping tests connectivity to the Rest API
func (b *BinanceClient) Ping() error {
	_, err := b.client.do(http.MethodGet, "api/v1/ping", nil, false, false)
	return err
}

// Time tests connectivity to the Rest API and get the current server time
func (b *BinanceClient) Time() (*ServerTime, error) {
	res, err := b.client.do(http.MethodGet, "api/v1/time", nil, false, false)
	if err != nil {
		return nil, err
	}
	serverTime := &ServerTime{}
	return serverTime, json.Unmarshal(res, serverTime)
}

// Market Data endpoints
// Depth retrieves the order book for the given symbol
func (b *BinanceClient) Depth(opts *DepthOpts) (*Depth, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	if opts.Limit == 0 || opts.Limit > 100 {
		opts.Limit = 100
	}
	res, err := b.client.do(http.MethodGet, "api/v1/depth", opts, false, false)
	if err != nil {
		return nil, err
	}
	depth := &Depth{}
	return depth, json.Unmarshal(res, &depth)
}

// AggregatedTrades gets compressed, aggregate trades.
// Trades that fill at the time, from the same order, with the same price will have the quantity aggregated
// Remark: If both startTime and endTime are sent, limit should not be sent AND the distance between startTime and endTime must be less than 24 hours.
// Remark: If frondId, startTime, and endTime are not sent, the most recent aggregate trades will be returned.
func (b *BinanceClient) AggregatedTrades(opts *AggregatedTradeOpts) ([]*AggregatedTrade, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	if opts.Limit == 0 || opts.Limit > 500 {
		opts.Limit = 500
	}
	res, err := b.client.do(http.MethodGet, "api/v1/aggTrades", opts, false, false)
	if err != nil {
		return nil, err
	}
	trades := []*AggregatedTrade{}
	return trades, json.Unmarshal(res, &trades)
}

// Klines returns kline/candlestick bars for a symbol. Klines are uniquely identified by their open time
func (b *BinanceClient) Klines(opts *KlinesOpts) ([]*Klines, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	if opts.Symbol == "" || opts.Interval == "" {
		return nil, fmt.Errorf("symbol or interval are missing")
	}
	if opts.Limit == 0 || opts.Limit > 500 {
		opts.Limit = 500
	}
	res, err := b.client.do(http.MethodGet, "api/v1/klines", opts, false, false)
	if err != nil {
		return nil, err
	}
	klines := []*Klines{}
	return klines, json.Unmarshal(res, &klines)

}

// Ticker returns 24 hour price change statistics
func (b *BinanceClient) Ticker(opts *TickerOpts) (*TickerStats, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	res, err := b.client.do(http.MethodGet, "api/v1/ticker/24hr", opts, false, false)
	if err != nil {
		return nil, err
	}
	tickerStats := &TickerStats{}
	return tickerStats, json.Unmarshal(res, tickerStats)
}

// Prices calculates the latest price for all symbols
func (b *BinanceClient) Prices() ([]*SymbolPrice, error) {
	res, err := b.client.do(http.MethodGet, "api/v1/ticker/allPrices", nil, false, false)
	if err != nil {
		return nil, err
	}
	prices := []*SymbolPrice{}
	return prices, json.Unmarshal(res, &prices)
}

// AllBookTickers returns best price/qty on the order book for all symbols
func (b *BinanceClient) AllBookTickers() ([]*BookTicker, error) {
	res, err := b.client.do(http.MethodGet, "api/v1/ticker/allBookTickers", nil, false, false)
	if err != nil {
		return nil, err
	}
	resp := []*BookTicker{}
	return resp, json.Unmarshal(res, &resp)
}

// Signed endpoints, associated with an account

// NewOrder sends in a new order
func (b *BinanceClient) NewOrder(opts *NewOrderOpts) (*NewOrder, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	res, err := b.client.do(http.MethodPost, "api/v3/order", opts, true, false)
	if err != nil {
		return nil, err
	}
	resp := &NewOrder{}
	return resp, json.Unmarshal(res, resp)
}

// NewOrderTest tests new order creation and signature/recvWindow long. Creates and validates a new order but does not send it into the matching engine
func (b *BinanceClient) NewOrderTest(opts *NewOrderOpts) error {
	if opts == nil {
		return fmt.Errorf("opts is nil")
	}
	_, err := b.client.do(http.MethodPost, "api/v3/order/test", opts, true, false)
	return err
}

// QueryOrder checks an order's status
func (b *BinanceClient) QueryOrder(opts *QueryOrderOpts) (*QueryOrder, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	if opts.OrderID < 0 && opts.OrigClientOrderId == "" {
		return nil, fmt.Errorf("order id must be set")
	}
	res, err := b.client.do(http.MethodGet, "api/v3/order", opts, true, false)
	if err != nil {
		return nil, err
	}
	resp := &QueryOrder{}
	return resp, json.Unmarshal(res, resp)
}

// CancelOrder cancel an active order
func (b *BinanceClient) CancelOrder(opts *CancelOrderOpts) (*CancelOrder, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	if opts.OrderID < 0 || (opts.OrigClientOrderId == "" && opts.NewClientOrderId == "") {
		return nil, fmt.Errorf("order id must be set")
	}
	res, err := b.client.do(http.MethodDelete, "api/v3/order", opts, true, false)
	if err != nil {
		return nil, err
	}
	resp := &CancelOrder{}
	return resp, json.Unmarshal(res, resp)
}

// OpenOrders get all open orders on a symbol
func (b *BinanceClient) OpenOrders(opts *OpenOrdersOpts) ([]*QueryOrder, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	res, err := b.client.do(http.MethodGet, "api/v3/openOrders", opts, true, false)
	if err != nil {
		return nil, err
	}
	resp := []*QueryOrder{}
	return resp, json.Unmarshal(res, &resp)
}

// AllOrders get all account orders; active, canceled, or filled
func (b *BinanceClient) AllOrders(opts *AllOrdersOpts) ([]*QueryOrder, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	if opts.Limit == 0 {
		opts.Limit = 500
	}
	res, err := b.client.do(http.MethodGet, "api/v3/allOrders", opts, true, false)
	if err != nil {
		return nil, err
	}
	resp := []*QueryOrder{}
	return resp, json.Unmarshal(res, &resp)
}

// Account get current account information
func (b *BinanceClient) Account() (*AccountInfo, error) {
	res, err := b.client.do(http.MethodGet, "api/v3/account", nil, true, false)
	if err != nil {
		return nil, err
	}
	resp := &AccountInfo{}
	return resp, json.Unmarshal(res, &resp)
}

// Trades get trades for a specific account and symbol
func (b *BinanceClient) Trades(opts *TradesOpts) (*Trades, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}
	if opts.Limit == 0 || opts.Limit > 500 {
		opts.Limit = 500
	}
	res, err := b.client.do(http.MethodGet, "api/v3/myTrades", opts, true, false)
	if err != nil {
		return nil, err
	}
	resp := &Trades{}
	return resp, json.Unmarshal(res, &resp)
}

// User stream endpoint

// Datastream starts a new user datastream
func (b *BinanceClient) DataStream() (string, error) {
	res, err := b.client.do(http.MethodPost, "api/v1/userDataStream", nil, false, true)
	if err != nil {
		return "", err
	}

	resp := &Datastream{}
	return resp.ListenKey, json.Unmarshal(res, &resp)
}

// DataStreamKeepAlive pings the datastream key to prevent timeout
func (b *BinanceClient) DataStreamKeepAlive(listenKey string) error {
	_, err := b.client.do(http.MethodPut, "api/v1/userDataStream", Datastream{ListenKey: listenKey}, false, true)
	return err
}

// DataStreamClose closes the datastream key
func (b *BinanceClient) DataStreamClose(listenKey string) error {
	_, err := b.client.do(http.MethodDelete, "api/v1/userDataStream", Datastream{ListenKey: listenKey}, false, true)
	return err
}

// DepthWS opens websocket with depth updates for the given symbol
func (b *BinanceClient) DepthWS(symbol string) (*DepthWS, error) {
	addr := strings.ToLower(symbol) + "@depth"
	conn, _, err := b.dialer.Dial(wsAddress+addr, nil)
	if err != nil {
		return nil, err
	}
	return &DepthWS{wsWrapper{conn: conn}}, nil
}

// KlinesWS opens websocket with klines updates for the given symbol with the given interval
func (b *BinanceClient) KlinesWS(symbol string, interval KlineInterval) (*KlinesWS, error) {
	addr := fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
	conn, _, err := b.dialer.Dial(wsAddress+addr, nil)
	if err != nil {
		return nil, err
	}
	return &KlinesWS{wsWrapper{conn: conn}}, nil
}

// TradesWS opens websocket with trades updates for the given symbol
func (b *BinanceClient) TradesWS(symbol string) (*TradesWS, error) {
	addr := strings.ToLower(symbol) + "@aggTrade"
	conn, _, err := b.dialer.Dial(wsAddress+addr, nil)
	if err != nil {
		return nil, err
	}
	return &TradesWS{wsWrapper{conn: conn}}, nil
}

// AccountInfoWS opens websocket with account info updates
func (b *BinanceClient) AccountInfoWS(listenKey string) (*AccountInfoWS, error) {
	conn, _, err := b.dialer.Dial(wsAddress+listenKey, nil)
	if err != nil {
		return nil, err
	}
	return &AccountInfoWS{wsWrapper{conn: conn}}, nil
}
