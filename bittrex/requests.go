package bittrex

import "net/url"

// publicParams represents the possible public parameters that can
// be passed for public Api calls.
type publicParams struct {
	MarketName   *string
	TickInterval *string
}

// AddToQueryString adds the non empty fields of the publicParams struct
// to the specified query string.
func (pp publicParams) AddToQueryString(queryString url.Values) {
	if queryString != nil {
		if pp.MarketName != nil {
			queryString.Add("marketName", *pp.MarketName)
		}
		if pp.TickInterval != nil {
			queryString.Add("tickInterval", *pp.TickInterval)
		}
	}
}

// privateParams represents the possible private parameters that can
// be passed for auth API calls.
type privateParams struct {
}

// AddToPostForm adds the non empty fields of the publicParams struct
// to the specified post form.
func (pp privateParams) AddToPostForm(postForm url.Values) {
	if postForm != nil {

	}
}
