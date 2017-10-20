package binance

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBinancePing(t *testing.T) {
	ctx := newBinanceCtx()
	require.NoError(t, ctx.api.Ping())
}

func TestBinanceClient_Time(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.Time()
	require.NoError(t, e)
}

func TestBinanceClient_Ticker(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.Ticker(&TickerOpts{"LTCBTC"})
	require.NoError(t, e)
}

func TestBinanceClient_Depth(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.Depth(&DepthOpts{Symbol: "NEOBTC"})
	require.NoError(t, e)
}

func TestBinanceClient_AggregatedTrades(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.AggregatedTrades(&AggregatedTradeOpts{Symbol: "NEOBTC"})
	require.NoError(t, e)
}

func TestBinanceClient_Klines(t *testing.T) {
	ctx := newBinanceCtx()
	s, e := ctx.api.Klines(&KlinesOpts{Symbol: "NEOBTC", Interval: KlineInterval1h, Limit: 5})
	require.NoError(t, e)
	require.Len(t, s, 5)
}

func TestBinanceClient_AllBookTickers(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.AllBookTickers()
	require.NoError(t, e)
}
func TestBinanceClient_Prices(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.Prices()
	require.NoError(t, e)
}

func TestBinanceClient_Order(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.NewOrder(&NewOrderOpts{
		Symbol:      "NEOBTC",
		Side:        OrderSideSell,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "1",
		Price:       "0.1",
	})
	require.NoError(t, e)
}

func TestBinanceClient_QueryCancelOrder(t *testing.T) {
	ctx := newBinanceCtx()
	s, e := ctx.api.NewOrder(&NewOrderOpts{
		Symbol:      "NEOBTC",
		Side:        OrderSideSell,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "1",
		Price:       "0.1",
	})
	require.NoError(t, e)
	q, e := ctx.api.QueryOrder(&QueryOrderOpts{
		Symbol:  "NEOBTC",
		OrderID: s.OrderID,
	})
	require.NoError(t, e)
	require.Equal(t, "NEOBTC", q.Symbol)
	c, e := ctx.api.CancelOrder(&CancelOrderOpts{
		Symbol:  "NEOBTC",
		OrderID: s.OrderID,
	})
	require.NoError(t, e)
	require.Equal(t, "NEOBTC", c.Symbol)
}

func TestBinanceClient_DataStream(t *testing.T) {
	ctx := newBinanceCtx()
	key, err := ctx.api.DataStream()
	require.NoError(t, err)
	require.NotEmpty(t, key)
	require.NoError(t, ctx.api.DataStreamKeepAlive(key))
	require.NoError(t, ctx.api.DataStreamClose(key))
}

func TestBinanceClient_AllOrders(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.AllOrders(&AllOrdersOpts{Symbol: "SNMBTC"})
	require.NoError(t, e)
}

func TestBinanceClient_OpenOrders(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.OpenOrders(&OpenOrdersOpts{Symbol: "SNMBTC"})
	require.NoError(t, e)
}

func TestBinanceClient_Account(t *testing.T) {
	ctx := newBinanceCtx()
	_, e := ctx.api.Account()
	require.NoError(t, e)
}

func TestBinanceClient_DepthWS(t *testing.T) {
	const symbol = "ETHBTC"
	ctx := newBinanceCtx()
	ws, err := ctx.api.DepthWS(symbol)
	defer ws.Close()
	require.NoError(t, err)
	u, err := ws.Read()
	require.NoError(t, err)
	require.Equal(t, symbol, u.Symbol)
}

func TestBinanceClient_KlinesWS(t *testing.T) {
	const symbol = "ETHBTC"
	ctx := newBinanceCtx()
	ws, err := ctx.api.KlinesWS(symbol, KlineInterval1m)
	defer ws.Close()
	require.NoError(t, err)
	u, err := ws.Read()
	require.NoError(t, err)
	require.Equal(t, symbol, u.Symbol)
}

func TestBinanceClient_TradesWS(t *testing.T) {
	const symbol = "ETHBTC"
	ctx := newBinanceCtx()
	ws, err := ctx.api.TradesWS(symbol)
	defer ws.Close()
	require.NoError(t, err)
	u, err := ws.Read()
	require.NoError(t, err)
	require.Equal(t, symbol, u.Symbol)
}

func TestBinanceClient_AccountInfoWS(t *testing.T) {
	ctx := newBinanceCtx()
	key, err := ctx.api.DataStream()
	require.NoError(t, err)
	defer ctx.api.DataStreamClose(key)
	ws, err := ctx.api.AccountInfoWS(key)
	require.NoError(t, err)
	defer ws.Close()
	u1, u2, err := ws.Read()
	require.NoError(t, err)
	if u1 == nil && u2 == nil {
		require.FailNow(t, "Expected to receive an update")
	}
}

type binanceCtx struct {
	api *BinanceClient
}

func newBinanceCtx() *binanceCtx {
	return &binanceCtx{
		api: NewBinanceClient("", ""),
	}
}
