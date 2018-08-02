package bittrex

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/thebotguys/signalr"

	"github.com/shopspring/decimal"

	"compress/gzip"
)

// OrderDeltaType is the type of an order delta.
type OrderDeltaType uint8

// ExchangeDeltaType is the type of a market delta.
type ExchangeDeltaType uint8

const (
	// WebsocketHost represents the API Websocket endpoint.
	WebsocketHost = "socket.bittrex.com/signalr"
	// AllowedHub is the only allowed hub of bittrex signalr.
	AllowedHub = "c2"
	// SummariesEndpoint is the endpoint to subscribe for market summary updates.
	SummariesEndpoint = "summaries"
	// TickEndpoint is the endpoint to subscribe for tick updates.
	TickEndpoint = "tick"

	// OpenType is the type of an opened order.
	OpenType OrderDeltaType = 0
	// PartialType is the type of a partially filled order.
	PartialType OrderDeltaType = 1
	// FillType is the type of a filled order.
	FillType OrderDeltaType = 2
	// CancelType is the type of a canceled order.
	CancelType OrderDeltaType = 3

	// AddType represents an added market.
	AddType ExchangeDeltaType = 0
	// RemoveType represents a removed market.
	RemoveType ExchangeDeltaType = 1
	// UpdateType represents an updated market.
	UpdateType ExchangeDeltaType = 2
)

var minifiedJSONKeys = map[string]string{
	"A":  "Ask",
	"a":  "Available",
	"B":  "Bid",
	"b":  "Balance",
	"C":  "Closed",
	"c":  "Currency",
	"CI": "CancelInitiated",
	"D":  "Deltas",
	"d":  "Delta",
	"DT": "OrderDeltaType",
	"E":  "Exchange",
	"e":  "ExchangeDeltaType",
	"F":  "FillType",
	"FI": "FillId",
	"f":  "Fills",
	"G":  "OpenBuyOrders",
	"g":  "OpenSellOrders",
	"H":  "High",
	"h":  "AutoSell",
	"I":  "Id",
	"i":  "IsOpen",
	"J":  "Condition",
	"j":  "ConditionTarget",
	"K":  "ImmediateOrCancel",
	"k":  "IsConditional",
	"L":  "Low",
	"l":  "Last",
	"M":  "MarketName",
	"m":  "BaseVolume",
	"N":  "Nonce",
	"n":  "CommissionPaid",
	"O":  "Orders",
	"o":  "Order",
	"OT": "OrderType",
	"OU": "OrderUuid",
	"P":  "Price",
	"p":  "CryptoAddress",
	"PD": "PrevDay",
	"PU": "PricePerUnit",
	"Q":  "Quantity",
	"q":  "QuantityRemaining",
	"R":  "Rate",
	"r":  "Requested",
	"S":  "Sells",
	"s":  "Summaries",
	"T":  "TimeStamp",
	"t":  "Total",
	"TY": "Type",
	"U":  "Uuid",
	"u":  "Updated",
	"V":  "Volume",
	"W":  "AccountId",
	"w":  "AccountUuid",
	"X":  "Limit",
	"x":  "Created",
	"Y":  "Opened",
	"y":  "State",
	"Z":  "Buys",
	"z":  "Pending",
}

// ExchangeDelta is object returned by a SubscribeToExchangeDeltas call.
type ExchangeDelta struct {
	MarketName string               `json:"MarketName,required,string"`
	Nonce      int                  `json:"Nonce,required,number"`
	Buys       []exchangeDeltaOrder `json:"Buys,required"`
	Sells      []exchangeDeltaOrder `json:"Sells,required"`
	Fills      []exchangeDeltaFill  `json:"Fills,required"`
}

// LiteSummaryDelta is object returned by a SubscribeToLiteSummaryDeltas call.
type LiteSummaryDelta struct {
	Deltas []LiteSummary `json:"Deltas,required"`
}

// SummaryDelta is object returned by a SubscribeToSummaryDeltas call.
type SummaryDelta struct {
	Deltas []Summary `json:"Deltas,required"`
}

// LiteSummary is the summary object as returned from LiteSummaryDelta.
type LiteSummary struct {
	MarketName string          `json:"MarketName,required,string"`
	Last       decimal.Decimal `json:"Last,required,number"`
	BaseVolume decimal.Decimal `json:"BaseVolume,required,number"`
}

// Summary is the summary object as returned from SummaryDelta.
type Summary struct {
	MarketName     string          `json:"MarketName,required,string"`
	Last           decimal.Decimal `json:"Last,required,number"`
	BaseVolume     decimal.Decimal `json:"BaseVolume,required,number"`
	High           decimal.Decimal `json:"High,required,number"`
	Low            decimal.Decimal `json:"Low,required,number"`
	Volume         decimal.Decimal `json:"Volume,required,number"`
	TimeStamp      time.Time       `json:"TimeStamp,required,string"`
	Bid            decimal.Decimal `json:"Bid,required,number"`
	Ask            decimal.Decimal `json:"Ask,required,number"`
	OpenBuyOrders  decimal.Decimal `json:"OpenBuyOrders,required,number"`
	OpenSellOrders decimal.Decimal `json:"OpenSellOrders,required,number"`
	PrevDay        decimal.Decimal `json:"PrevDay,required,number"`
	Created        time.Time       `json:"Created,required,string"`
}

type exchangeDeltaOrder struct {
	Type     OrderDeltaType  `json:"Type,required,number"`
	Rate     decimal.Decimal `json:"Rate,required,number"`
	Quantity decimal.Decimal `json:"Quantity,required,number"`
}

type exchangeDeltaFill struct {
	FillID    int             `json:"FillId,required,number"`
	OrderType string          `json:"Type,required,string"`
	Rate      decimal.Decimal `json:"Rate,required,number"`
	Quantity  decimal.Decimal `json:"Quantity,required,number"`
	Timestamp time.Time       `json:"TimeStamp,required,string"`
}

// WebsocketService provides websocket features from bittrex signalr endpoint.
type WebsocketService struct{}

// Websocket is the global WebsocketService.
var Websocket WebsocketService
var ws *signalr.Client

// WSError contains errors coming from a websocket subscription handling.
type WSError struct {
	Channel string // The name of the channel.
	Error   error  // The Error to be shown.
}

// Channels

// SubscriptionChans contains all possible channels from the bittrex Websocket.
type SubscriptionChans struct {
	ExchangeDeltas    map[string]chan ExchangeDelta // Data from ExchangeDelta subscriptions, divided by market.
	SummaryDeltas     map[string]chan Summary       // Data from SummaryDelta subscriptions, divided by market.
	LiteSummaryDeltas map[string]chan LiteSummary   // Data from LiteSummaryDelta subscriptions, divided by market.
	Errors            chan WSError                  // Errors thrown by the websocket client method.
}

// CloseAll closes all channels.
func (sc *SubscriptionChans) CloseAll() {
	for market := range sc.ExchangeDeltas {
		close(sc.ExchangeDeltas[market])
	}
	for market := range sc.SummaryDeltas {
		close(sc.SummaryDeltas[market])
	}
	for market := range sc.LiteSummaryDeltas {
		close(sc.LiteSummaryDeltas[market])
	}
	close(sc.Errors)
}

var defaultChannelsSubscription = SubscriptionChans{
	ExchangeDeltas:    make(map[string]chan ExchangeDelta),
	SummaryDeltas:     make(map[string]chan Summary),
	LiteSummaryDeltas: make(map[string]chan LiteSummary),
	Errors:            make(chan WSError),
}

// Channels contains all possible channels for all bittrex websocket events.
var Channels *SubscriptionChans

// Connect connects the client to the bittrex websocket.
func (wss WebsocketService) Connect() error {
	Channels = new(SubscriptionChans)
	*Channels = defaultChannelsSubscription

	ws = signalr.NewWebsocketClient()
	ws.OnClientMethod = func(hub string, method string, arguments []json.RawMessage) {
		if hub != AllowedHub {
			return
		}

		err := parseSignalrResponse(method, arguments)
		if err != nil {
			Channels.Errors <- WSError{
				Channel: method,
				Error:   err,
			}
		}
	}
	return ws.Connect("https", WebsocketHost, []string{AllowedHub})
}

func parseSignalrResponse(method string, arguments []json.RawMessage) error {
	for _, msg := range arguments {
		gzipMsg, err := base64.StdEncoding.DecodeString(string(msg))
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(gzipMsg)
		gzipReader, err := gzip.NewReader(buf)
		if err != nil {
			return err
		}

		content, err := ioutil.ReadAll(gzipReader)
		if err != nil {
			return err
		}

		// TODO: decode correctly and populate the correct channel (every channel has its own variable)
		//       handle errors
		//       handle channel receival using a select case statement.
		//
		switch method {
		case "SubscribeToExchangeDeltas":
			var decoded ExchangeDelta
			err := json.Unmarshal(content, &decoded)
			if err != nil {
				return err
			}

			Channels.ExchangeDeltas[decoded.MarketName] <- decoded
			break
		case "SubscribeToSummaryDeltas":
			var decoded SummaryDelta
			err := json.Unmarshal(content, &decoded)
			if err != nil {
				return err
			}

			for _, delta := range decoded.Deltas {
				if channel, subscribed := Channels.SummaryDeltas[delta.MarketName]; subscribed {
					channel <- delta
				}
			}
			break
		case "SubscribeToLiteSummaryDeltas":
			var decoded LiteSummaryDelta
			err := json.Unmarshal(content, &decoded)
			if err != nil {
				return err
			}

			for _, delta := range decoded.Deltas {
				if channel, subscribed := Channels.LiteSummaryDeltas[delta.MarketName]; subscribed {
					channel <- delta
				}
			}

			break
		}
	}

	return nil
}

// Disconnect disconnects the websocket service.
func (wss WebsocketService) Disconnect() {
	if ws != nil {
		Channels.CloseAll()
		ws.Close()
		ws = nil
	}
}

// SubscribeToExchangeDeltas Subscribes to the ExchangeDeltas service.
func (wss WebsocketService) SubscribeToExchangeDeltas(market string) (chan ExchangeDelta, error) {
	if ws == nil {
		return nil, errors.New("Websocket not connected")
	}

	if _, exists := Channels.ExchangeDeltas[market]; !exists {
		_, err := ws.CallHub(AllowedHub, "SubscribeToExchangeDeltas", market)
		if err != nil {
			return nil, err
		}

		Channels.ExchangeDeltas[market] = make(chan (ExchangeDelta))
	}
	return Channels.ExchangeDeltas[market], nil
}

// SubscribeToSummaryDeltas Subscribes to the SummaryDeltas service.
func (wss WebsocketService) SubscribeToSummaryDeltas(market string) (chan Summary, error) {
	if ws == nil {
		return nil, errors.New("Websocket not connected")
	}

	if _, exists := Channels.ExchangeDeltas[market]; !exists {
		_, err := ws.CallHub(AllowedHub, "SubscribeToSummaryDeltas")
		if err != nil {
			return nil, err
		}

		Channels.SummaryDeltas[market] = make(chan (Summary))
	}
	return Channels.SummaryDeltas[market], nil
}

// SubscribeToLiteSummaryDeltas Subscribes to the LiteSummaryDeltas service.
func (wss WebsocketService) SubscribeToLiteSummaryDeltas(market string) (chan LiteSummary, error) {
	if ws == nil {
		return nil, errors.New("Websocket not connected")
	}

	if _, exists := Channels.ExchangeDeltas[market]; !exists {
		_, err := ws.CallHub(AllowedHub, "SubscribeToLiteSummaryDeltas")
		if err != nil {
			return nil, err
		}

		Channels.LiteSummaryDeltas[market] = make(chan (LiteSummary))
	}
	return Channels.LiteSummaryDeltas[market], nil
}
