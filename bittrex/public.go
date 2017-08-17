package bittrex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
		return "", errors.New("Status Code" + string(response.StatusCode))
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
	result, err := publicCall("currencies", "GetBTCPrice", nil)
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
