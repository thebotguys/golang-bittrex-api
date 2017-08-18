package bittrex_test

import (
	"testing"

	"github.com/saniales/golang-bittrex-api/bittrex"
)

func testIsAPIAlive(t *testing.T) {
	err := bittrex.IsAPIAlive()
	if err != nil {
		t.SkipNow()
	}
}
