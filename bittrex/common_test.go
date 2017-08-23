package bittrex_test

import (
	"testing"

	"github.com/thebotguys/golang-bittrex-api/bittrex"
)

func testIsAPIAlive(t *testing.T) bool {
	err := bittrex.IsAPIAlive()
	if err != nil {
		t.Log("API is not reachable, test is invalid")
	}
	return err == nil
}
