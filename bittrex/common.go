package bittrex

import (
	"encoding/json"

	"github.com/juju/errors"

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

// Auth represents the auth credentials to authenticate to the Bittrex API:
//
// It consists of a set of a private and a public key.
type Auth struct {
	PublicKey  string // The public key to connect to bittrex API.
	PrivateKey string // The private key to connect to bittrex API.
}

var (
	// defaultClient represents the default configuration for HTTP requests to the API.
	defaultClient = http.Client{
		Timeout: time.Second * 30,
	}
	// client represents the actual configuration for HTTP requests to the API.
	client = defaultClient
)

// SetCustomHTTPClient sets a custom client for requests.
func SetCustomHTTPClient(value http.Client) {
	client = value
}

// apiCall performs a generic API call.
func apiCall(Version, Visibility, Entity, Feature string, GetParameters *publicParams, PostParameters *privateParams) (*json.RawMessage, error) {
	URL := fmt.Sprintf("%s/v%s/%s/%s/%s", BaseURL, Version, Visibility, Entity, Feature)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, errors.Annotatef(err, "%s - URL: %s", Feature, URL)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Add("Cache-Control", "no-store")
	req.Header.Add("Cache-Control", "must-revalidate")

	if Visibility == Public && GetParameters != nil { // Add them to query string
		queryString := req.URL.Query()
		GetParameters.AddToQueryString(&queryString)
		req.URL.RawQuery = queryString.Encode()
	} else if Visibility == Private {
		addSecurityHeaders(req.Header)
		if PostParameters != nil {
			PostParameters.AddToPostForm(&req.PostForm)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status Code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret response
	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Success == false {
		return nil, fmt.Errorf("Error Response: %s", ret.Message)
	}

	return ret.Result, nil
}

// publicCall performs a call to the public bittrex API.
//
// It does not need API Keys.
func publicCall(Entity, Feature string, GetParameters *publicParams) (*json.RawMessage, error) {
	return apiCall(APIVersion, Public, Entity, Feature, GetParameters, nil)
}

// authCall performs a call to the private bittrex API.
//
// It needs an Auth struct to be passed with valid Keys.
func authCall(Entity, Feature string, PostParams *privateParams, auth Auth) (*json.RawMessage, error) {
	if auth.PublicKey == "" || auth.PrivateKey == "" {
		return nil, errors.New("Cannot perform private api request without authentication keys")
	}
	//createHMAC signature
	return apiCall(APIVersion, Private, Entity, Feature, nil, PostParams)
}

// addSecurityHeaders adds security headers, required for bittrex private API calls.
//
// Example of this headers which need to be added are HMAC signature.
func addSecurityHeaders(header http.Header) {

}
