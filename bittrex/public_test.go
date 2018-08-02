package bittrex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thebotguys/golang-bittrex-api/bittrex"
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
		time.Sleep(time.Second * 2)
		_, err := bittrex.GetTicks("BTC-ETH", "thirtyMin")
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 2)
		val, err := bittrex.GetTicks("INVALID-MARKET-HAHA", "thirtyMin")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetTicks(\"INVALID-MARKET-HAHA\", " +
				"\"thirtyMin\"))\n Value Returned : " +
				fmt.Sprint(val))
		}
		time.Sleep(time.Second * 2)
		val, err = bittrex.GetTicks("", "thirtyMin")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetTicks(\"\", \"thirtyMin\"))\n" +
				"Value Returned : " + fmt.Sprint(val))
		}
		time.Sleep(time.Second * 2)
		val, err = bittrex.GetTicks("BTC-LTC", "INVALID-INTERVAL-HAHA")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetTicks(\"BTC-LTC\", \"INVALID-INTERVAL-HAHA\"))\n" +
				"Value Returned : " + fmt.Sprint(val))
		}
		time.Sleep(time.Second * 2)
		val, err = bittrex.GetTicks("BTC-LTC", "")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetTicks(\"BTC-LTC\", \"\"))\n" +
				"Value Returned : " + fmt.Sprint(val))
		}
	}
}

func TestGetLatestTick(t *testing.T) {
	if testIsAPIAlive(t) {
		_, err := bittrex.GetLatestTick("BTC-ETH", "thirtyMin")
		if err != nil {
			t.Fatal(err)
		}
		val, err := bittrex.GetLatestTick("INVALID-MARKET-HAHA", "thirtyMin")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetLatestTick(\"INVALID-MARKET-HAHA\", " +
				"\"thirtyMin\"))\n Value Returned : " +
				fmt.Sprint(val))
		}
		val, err = bittrex.GetLatestTick("", "thirtyMin")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetLatestTick(\"\", \"thirtyMin\"))\n" +
				"Value Returned : " + fmt.Sprint(val))
		}
		val, err = bittrex.GetLatestTick("BTC-LTC", "INVALID-INTERVAL-HAHA")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetLatestTick(\"BTC-LTC\", \"INVALID-INTERVAL-HAHA\"))\n" +
				"Value Returned : " + fmt.Sprint(val))
		}
		val, err = bittrex.GetLatestTick("BTC-LTC", "")
		if err == nil {
			t.Fatal("Error expected, but function did " +
				"not fail (GetLatestTick(\"BTC-LTC\", \"\"))\n" +
				"Value Returned : " + fmt.Sprint(val))
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
		time.Sleep(time.Second)
		val, err := bittrex.GetMarketSummary("INVALID-MARKET-HAHA") // error expected
		if err == nil {
			t.Fatal("Error expected, but function did not fail (GetMarketSummary(\"INVALID-MARKET-HAHA\"))\n Value Returned : " + fmt.Sprint(val))
		}
		time.Sleep(time.Second)
		val, err = bittrex.GetMarketSummary("") // error expected
		if err == nil {
			t.Fatal("Error expected, but function did not fail (GetMarketSummary(\"\"))\n Value Returned : " + fmt.Sprint(val))
		}
	}
}

func TestGetMarkets(t *testing.T) {
	if testIsAPIAlive(t) {
		_, err := bittrex.GetMarkets() // no error expected
		if err != nil {
			t.Fatal(err)
		}
	}
}
