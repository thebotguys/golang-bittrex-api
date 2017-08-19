package bittrex

import (
	"encoding/json"
	"fmt"
	"time"
)

type response struct {
	Success bool             `json:"success,required"`
	Message string           `json:"message,omitempty"`
	Result  *json.RawMessage `json:"result,required"`
}

type pingResponse struct {
	Response string `json:"response,required"`
}

type versionResponse struct {
	Version json.Number `json:"version,required,Number"`
}

// MarketSummary is the summary data of a market (usually 24h summary).
type MarketSummary struct {
	MarketName     string  `json:"MarketName,required"`     //The name of the market (e.g. BTC-ETH).
	High           float64 `json:"High,required"`           // The 24h high for the market.
	Low            float64 `json:"Low,required"`            // The 24h low for the market.
	Last           float64 `json:"Last,required"`           // The value of the last trade for the market (in base currency).
	Bid            float64 `json:"Bid,required"`            // The current highest bid value for the market.
	Ask            float64 `json:"Ask,required"`            // The current lowest ask value for the market.
	Volume         float64 `json:"Volume,required"`         // The 24h volume of the market, in market currency.
	BaseVolume     float64 `json:"BaseVolume,required"`     // The 24h volume for the market, in base currency.
	Timestamp      string  `json:"Timestamp,required"`      // The timestamp of the request.
	OpenBuyOrders  uint64  `json:"OpenBuyOrders,required"`  // The number of currently open buy orders.
	OpenSellOrders uint64  `json:"OpenSellOrders,required"` // The number of currently open sell orders.
	PrevDay        float64 `json:"PrevDay,required"`        // The closing price 24h before.
	Created        string  `json:"Created,required"`        // The timestamp of the creation of the market.
}

//MarketSummaries is a set of MarketSummary objects.
type MarketSummaries []MarketSummary

type marketSummariesResult []struct {
	IsVerified bool          `json:"IsVerified"`
	Market     Market        `json:"Market,required"`
	Summary    MarketSummary `json:"Summary,required"`
}

// OpenOrder represents a currently open order.
type OpenOrder struct {
}

// OrderBook represents a set of public open Orders, which compose the OrderBook.
type OrderBook []OpenOrder

// Market represents a market metadata (name, base currency, trade currency)
// and so forth.
type Market struct {
	BaseCurrency       string  `json:"BaseCurrency,required"`
	BaseCurrencyLong   string  `json:"BaseCurrencyLong,required"`
	MarketCurrency     string  `json:"MarketCurrency,required"`
	MarketCurrencyLong string  `json:"MarketCurrencyLong,required"`
	MarketName         string  `json:"MarketName,required"`
	MinTradeSize       float64 `json:"MinTradeSize,required"`
	IsActive           bool    `json:"IsActive,required"`
	Created            string  `json:"Created,required"`
	Notice             string  `json:"Notice,required"`
	IsSponsored        bool    `json:"IsSponsored,required"`
	LogoURL            string  `json:"LogoUrl,required"`
}

// Markets is a set of markets.
type Markets []Market

type btcPriceResult struct {
	Bpi struct {
		USD struct {
			Code        string      `json:"code,required"`
			Description string      `json:"description,required"`
			Rate        string      `json:"rate,required"`
			RateFloat   json.Number `json:"rate_float,required"`
		} `json:"USD,required"`
		Disclaimer string `json:"disclaimer,required"`
	} `json:"bpi,required"`
	Time struct {
		Updated    string `json:"updated,required"`
		UpdatedISO string `json:"updatedISO,omitempty"`
		UpdatedUK  string `json:"updateduk,omitempty"`
	} `json:"time,required"`
}

// CandleStick represents a single candlestick in a chart.
type CandleStick struct {
	High       float64    `json:"H,required"`
	Open       float64    `json:"O,required"`
	Close      float64    `json:"C,required"`
	Low        float64    `json:"L,required"`
	Volume     float64    `json:"V,required"`
	BaseVolume float64    `json:"BV,required"`
	Timestamp  candleTime `json:"T,required"`
}

// CandleSticks is an array of CandleStick objects. It is a result from GetTicks
// and GetLatestTick calls too.
type CandleSticks []CandleStick

// CandleIntervals represent all valid intervals supported
// by the GetTicks and GetLatestTick calls.
var CandleIntervals = map[string]bool{
	"oneMin":    true,
	"fiveMin":   true,
	"thirtyMin": true,
	"hour":      true,
	"day":       true,
}

type candleTime time.Time

func (t *candleTime) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("could not parse time %s", string(b))
	}
	// trim enclosing ""
	result, err := time.Parse("2006-01-02T15:04:05", string(b[1:len(b)-1]))
	if err != nil {
		return fmt.Errorf("could not parse time: %v", err)
	}
	*t = candleTime(result)
	return nil
}

func (result btcPriceResult) Compress() BTCPrice {
	value, _ := result.Bpi.USD.RateFloat.Float64()
	ts, _ := time.Parse(time.RFC3339, result.Time.UpdatedISO)
	return BTCPrice{
		USDValue:  value,
		Timestamp: ts,
	}
}

// BTCPrice represents the BTC price at the specified timestamp.
type BTCPrice struct {
	USDValue  float64
	Timestamp time.Time
}
