package connector_test

import (
	"context"
	"testing"

	connector "github.com/veska-io/streams-connectors/bybit/futures/kline/exchange-connector/rest-clickhouse/src"
)

func TestRun(t *testing.T) {
	connector.MustRun(context.Background())
}
