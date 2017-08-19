package bittrex_test

import (
	"testing"

	"github.com/saniales/golang-bittrex-api/bittrex"
)

var latestTestedVersion = "2.11"

func TestGetServerAPIVersion(t *testing.T) {
	if testIsAPIAlive(t) {
		version, err := bittrex.GetServerAPIVersion()
		if err != nil {
			t.Fatal(err)
		}
		if version != latestTestedVersion {
			t.Errorf(`Please check version you are testing, on server it is %s, 
while on this client it is %s.`,
				latestTestedVersion, version)
		}
	}
}

func TestGetBTCPrice(t *testing.T) {
	if testIsAPIAlive(t) {
		_, err := bittrex.GetBTCPrice()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetTicks(t *testing.T) {
	if testIsAPIAlive(t) {
		_, err := bittrex.GetTicks("BTC-ETH", "thirtyMin")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetLatestTick(t *testing.T) {
	if testIsAPIAlive(t) {
		_, err := bittrex.GetLatestTick("BTC-ETH", "thirtyMin")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetMarketSummaries(t *testing.T) {
	if testIsAPIAlive(t) {
		_, err := bittrex.GetMarketSummaries()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetMarketSummary(t *testing.T) {
	if testIsAPIAlive(t) {
		var err error
		_, err = bittrex.GetMarketSummary("BTC-ETH") // no error expected
		if err != nil {
			t.Fatal(err)
		}
		_, err = bittrex.GetMarketSummary("INVALID-MARKET") // error expected
		if err == nil {
			t.Fatal("Error expected, but function did not fail")
		}
	}
}
