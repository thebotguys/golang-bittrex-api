package bittrex

import (
	"encoding/json"
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
