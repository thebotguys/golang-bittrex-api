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

// BaseURL represents the base URL for all requests
const (
	BaseURL    = "https://bittrex.com/api"
	Private    = "auth"
	Public     = "pub"
	APIVersion = "2.0"
)

type Bittrex struct {
	PublicKey  string // The public key to connect to bittrex API.
	PrivateKey string // The private key to connect to bittrex API.
}

// APIOptions lol
type APIOptions struct{}

func checkOptions() *APIOptions {
	return nil
}

func (b Bittrex) apiCall(Version, Visibility, Entity, Feature string, options APIOptions) error {
	//client := http.DefaultClient
	//URL := fmt.Sprintf("%s/v%s/%s/%s/%s", BaseURL, Version, Visibility, Entity, Feature)
	//options = checkOptions()
	return nil
}

// IsAPIAlive returns true if the bittrex api is reachable with the current network connection.
func IsAPIAlive() error {
	var pingResponse struct {
		Response string `json:"response,required"`
	}
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
	var versionResponse struct {
		Version json.Number `json:"version,required,Number"`
	}
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
