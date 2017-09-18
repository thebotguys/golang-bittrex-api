package bittrex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

// IsAPIAlive returns true if the bittrex api is reachable with the current network connection.
func IsAPIAlive() error {
	var pingResponse pingResponse
	timestamp := time.Now().UTC().Unix()
	URL := fmt.Sprintf("https://socket.bittrex.com/signalr/ping?_=%d", timestamp)
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Status Code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &pingResponse)
	if err != nil {
		return err
	}
	if pingResponse.Response == "pong" {
		return nil
	}
	return errors.New("API is not live")
}

// GetServerAPIVersion returns the version which is currently running on the server.
func GetServerAPIVersion() (string, error) {
	var versionResponse versionResponse
	timestamp := time.Now().UTC().Unix()
	URL := fmt.Sprintf("https://bittrex.com/Content/version.txt?_=%d", timestamp)
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Status Code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	err = json.Unmarshal(body, &versionResponse)
	if err != nil {
		return "", err
	}

	return versionResponse.Version.String(), nil
}

// GetBTCPrice returns the current BTC Price.
func GetBTCPrice() (*BTCPrice, error) {
	result, err := publicCall("currencies", "GetBTCPrice", nil, nil)
	if err != nil {
		return nil, err
	}

	var btcPrice BTCPrice

	err = json.Unmarshal(*result, &btcPrice)
	if err != nil {
		return nil, err
	}
	return &btcPrice, nil
}

// GetLatestTick returns the latest tick of the
// specified market's candlestick chart,
// following the specified tick interval.
func GetLatestTick(marketName, tickInterval string) (*CandleStick, error) {
	ticks, err := tickFunc(marketName, tickInterval, "GetLatestTick")
	if err != nil {
		return nil, err
	}
	return &ticks[0], nil
}

// GetTicks returns the ticks of the
// specified market's candlestick chart,
// following the specified tick interval.
func GetTicks(marketName, tickInterval string) (CandleSticks, error) {
	return tickFunc(marketName, tickInterval, "GetTicks")
}

// GetMarketSummaries gets the summary of all markets.
func GetMarketSummaries() (MarketSummaries, error) {
	now := time.Now().Unix()
	GetParameters := publicParams{
		Timestamp: &now,
	}
	result, err := publicCall("markets", "GetMarketSummaries", &GetParameters, nil)
	if err != nil {
		return nil, err
	}
	var resp marketSummariesResult
	err = json.Unmarshal(*result, &resp)
	if err != nil {
		return nil, err
	}
	ret := make(MarketSummaries, len(resp))
	for i, respItem := range resp {
		ret[i] = respItem.Summary
	}
	return ret, nil
}

// GetMarketSummary gets the summary of a single market.
func GetMarketSummary(marketName string) (*MarketSummary, error) {
	now := time.Now().UnixNano()
	GetParameters := publicParams{
		MarketName: &marketName,
		Timestamp:  &now,
	}
	result, err := publicCall("market", "GetMarketSummary", &GetParameters, nil)
	if err != nil {
		return nil, err
	}

	var ret MarketSummary
	err = json.Unmarshal(*result, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// tickFunc is a common pattern for GetTicks and GetLatestTick functions.
func tickFunc(marketName, tickInterval, tickFeature string) (CandleSticks, error) {
	now := time.Now().Unix()
	GetParameters := publicParams{
		MarketName:   &marketName,
		TickInterval: &tickInterval,
		Timestamp:    &now,
	}
	result, err := publicCall("market", tickFeature, &GetParameters, nil)
	if err != nil {
		return nil, err
	}
	var ret CandleSticks
	err = json.Unmarshal(*result, &ret)
	if err != nil {
		return nil, err
	}
	sort.Sort(csByTimestamp{ret})
	return ret, nil
}

// GetOrderBook gets the current order book, made up by currently open orders.
func GetOrderBook(marketName string) (OrderBook, error) {
	return nil, nil
}

//GetMarkets gets all markets data.
func GetMarkets() (Markets, error) {
	now := time.Now().Unix()
	GetParameters := publicParams{
		Timestamp: &now,
	}
	result, err := publicCall("markets", "GetMarkets", &GetParameters, nil)
	if err != nil {
		return nil, err
	}
	var resp Markets
	err = json.Unmarshal(*result, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
