package connector_test

import (
	"testing"

	connector "github.com/veska-io/streams-connectors/bybit/futures/funding-rate/exchange-connector/rest-clickhouse/src"
)

func TestRun(t *testing.T) {
	connector.MustRun()
}
