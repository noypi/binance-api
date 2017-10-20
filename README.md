# Golang Binance API
binance-api is a lightweight Golang implementation for [Binance API](https://www.binance.com/restapipub.html), providing complete API coverage, and supports both REST API and websockets API

This project is designed to help you interact with the Binance API, streaming candlestick charts data, market depth, or use other advanced features binance exposes via API. 

# Installation
```
go get github.com/eranyanay/binance-api
```

# Getting started
```golang
// Generate default client
client := binance.NewBinanceClient("API-KEY", "SECRET")

// Generate client with custom window size
client := binance.NewBinanceClientWindow("API-KEY", "SECRET", 5000)
```

# Examples
## REST API usage examples

### Get depth for symbol
```golang
depth, err := client.Depth(&binance.DepthOpts{Symbol: "ETHBTC"})
```

### Get candlesticks for symbol
```golang
klines, err := client.Klines(&binance.KlinesOpts{
		Symbol: "ETHBTC",
		Interval: binance.KlineInterval15m,
})
```

### Get 24 hour price change statistics for symbol
```golang
stats, err := client.Ticker(&binance.TickerOpts{Symbol: "ETHBTC"})
```

### Get latest prices for all symbols
```golang
prices, err := client.Prices()
```

### Create new order for ETHBTC, purchase 1 quantity at price 0.05BTC
```golang
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
```golang
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
```golang
query, err := client.QueryOrder(&binance.QueryOrderOpts{Symbol:"ETHBTC", OrderID:order.OrderID})
```

### Cancel an open order
```golang
cancel, err := client.CancelOrder(&binance.CancelOrderOpts{Symbol:"ETHBTC", OrderID:order.OrderID})
```

### Open orders for a symbol
```golang
orders, err := client.OpenOrders(&binance.OpenOrdersOpts{Symbol:"ETHBTC"})
```

### Get all account orders for symbol ETHBTC, with ID greater than 5
```golang
orders, err := client.AllOrders(&binance.AllOrdersOpts{Symbol:"ETHBTC", OrderID: 5})
```

### Get account information
```golang
info, err := client.Account()
```

### Get account trades for symbol
```golang
trades, err := client.Trades(&binance.TradesOpts{Symbol:"ETHBTC"})
```

### Create a new datastream key
```golang
key, err := client.DataStream()
```

### Set datastream key keep-alive to prevent from it to timeout
```golang
err := client.DataStreamKeepAlive(key)
```

### Close datastream
```golang
err := client.DataStreamClose(key)
```

