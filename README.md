# Golang Binance API
binance-api is a lightweight Golang implementation for [Binance API](https://www.binance.com/restapipub.html), providing complete API coverage, and supports both REST API and websockets API

This project is designed to help you interact with the Binance API, streaming candlestick charts data, market depth, or use other advanced features binance exposes via API. 

# Installation
```
go get github.com/eranyanay/binance-api
```

# Getting started
```
// Generate default client
client := binance.NewBinanceClient("API-KEY", "SECRET")

// Generate client with custom window size
client := binance.NewBinanceClientWindow("API-KEY", "SECRET", 5000)
```

# Examples
## REST API usage examples

### Get depth for symbol
```
depth, err := client.Depth(&binance.DepthOpts{Symbol: "ETHBTC"})
```

### Get candlesticks for symbol
```
klines, err := client.Klines(&binance.KlinesOpts{
		Symbol: "ETHBTC",
		Interval: binance.KlineInterval15m,
})
```

### Get 24 hour price change statistics for symbol
```
stats, err := client.Ticker(&binance.TickerOpts{Symbol: "ETHBTC"})
```

### Get latest prices for all symbols
```
prices, err := client.Prices()
```

### Create new order for ETHBTC, purchase 1 quantity at price 0.05BTC
```
order, err := client.NewOrder(&binance.NewOrderOpts{
		Symbol: "ETHBTC",
		Type: binance.OrderTypeLimit,
		Price: "0.05",
		Quantity: "1",
		Side: binance.OrderSideBuy,
		TimeInForce: binance.TimeInForceGTC,
	})
```
### Test new order options
```
order, err := client.NewOrderTest(&binance.NewOrderOpts{
		Symbol: "ETHBTC",
		Type: binance.OrderTypeLimit,
		Price: "0.05",
		Quantity: "1",
		Side: binance.OrderSideBuy,
		TimeInForce: binance.TimeInForceGTC,
	})
```

### Query order status
Orders are assigned with order ID when issued and can later be queried using it
```
query, err := client.QueryOrder(&binance.QueryOrderOpts{Symbol:"ETHBTC", OrderID:order.OrderID})
```

### Cancel an open order
```
cancel, err := client.CancelOrder(&binance.CancelOrderOpts{Symbol:"ETHBTC", OrderID:order.OrderID})
```

### Open orders for a symbol
```
orders, err := client.OpenOrders(&binance.OpenOrdersOpts{Symbol:"ETHBTC"})
```
