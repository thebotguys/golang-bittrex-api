package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// BaseURL represents the base URL for all requests
const (
	BaseURL    = "https://bittrex.com/api"
	Private    = "auth"
	Public     = "pub"
	APIVersion = "2.0"
)

type Auth struct {
	PublicKey  string // The public key to connect to bittrex API.
	PrivateKey string // The private key to connect to bittrex API.
}

// APIOptions lol
type APIOptions struct {
	hmacSignature interface{}
	Auth          Auth
}

func checkOptions() *APIOptions {
	return nil
}

// apiCall performs a generic API call.
func apiCall(Version, Visibility, Entity, Feature string, options *APIOptions) (*json.RawMessage, error) {
	client := http.DefaultClient
	URL := fmt.Sprintf("%s/v%s/%s/%s/%s", BaseURL, Version, Visibility, Entity, Feature)
	return nil, nil
}

// publicCall performs a call to the public bittrex API. It does not need API Keys.
func publicCall(Entity, Feature string, options *APIOptions) (*json.RawMessage, error) {
	options = checkOptions()
	return apiCall(APIVersion, Public, Entity, Feature, options)
}

func authCall(Entity, Feature string, options *APIOptions) (*json.RawMessage, error) {
	//options = checkOptions()
	if options.Auth.PublicKey == "" || options.Auth.PrivateKey == "" {
		return nil, errors.New("Cannot perform private api requst without authentication keys")
	}
	//createHMAC signature
	return apiCall(APIVersion, Private, Entity, Feature, options)
}
