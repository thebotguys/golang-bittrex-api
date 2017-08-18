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
