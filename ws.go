package binance

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type wsWrapper struct {
	conn *websocket.Conn
}

func (w *wsWrapper) Close() error {
	return w.conn.Close()
}

// DepthWS is a wrapper for depth websocket
type DepthWS struct {
	wsWrapper
}

// Read reads a depth update message from the depth websocket
func (d *DepthWS) Read() (*DepthUpdate, error) {
	_, data, err := d.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	update := &DepthUpdate{}
	return update, json.Unmarshal(data, update)
}

// KlinesWS is a wrapper for klines websocket
type KlinesWS struct {
	wsWrapper
}

// Read reads a klines update message from the klines websocket
func (d *KlinesWS) Read() (*KlinesUpdate, error) {
	_, data, err := d.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	update := &KlinesUpdate{}
	return update, json.Unmarshal(data, update)
}

// TradesWS is a wrapper for trades websocket
type TradesWS struct {
	wsWrapper
}

// Read reads a trades update message from the trades websocket
func (d *TradesWS) Read() (*TradesUpdate, error) {
	_, data, err := d.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	update := &TradesUpdate{}
	return update, json.Unmarshal(data, update)
}

// AccountInfoWS is a wrapper for account info websocket
type AccountInfoWS struct {
	wsWrapper
}

// Read reads a account info update message from the account info websocket
// Remark: The websocket is used to update two different structs, which both are flat, hence every call to this function
// will return either one of the types initialized and the other one will be set to nil
func (d *AccountInfoWS) Read() (*AccountUpdate, *OrderUpdate, error) {
	_, data, err := d.conn.ReadMessage()
	if err != nil {
		return nil, nil, err
	}
	msgType := &struct {
		EventType UpdateType `json:"e"` // EventType represents the update type
		Time      uint64     `json:"E"` // Time represents the event time
	}{}
	if err := json.Unmarshal(data, msgType); err != nil {
		return nil, nil, err
	}
	if msgType.EventType == UpdateTypeOutboundAccountInfo {
		update := &AccountUpdate{}
		return update, nil, json.Unmarshal(data, update)
	}
	update := &OrderUpdate{}
	return nil, update, json.Unmarshal(data, update)
}
