package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/adshao/go-binance/v2"
	restc "github.com/veska-io/streams-connectors/consumers/rest"
)

type Consumer struct {
	restc.Consumer
}

func New(ctx context.Context, logger *slog.Logger,
	markets []string,
	start, end time.Time,
	rps uint8,
	taskQuant time.Duration,
) *Consumer {
	c := Consumer{}
	c.UpdateConfig(ctx, logger, markets, start, end, rps, taskQuant)
	c.ResponseChan = make(chan *restc.ResponseMessage)
	c.ApiCall = ApiCall

	return &c
}

func ApiCall(task restc.Task) (*restc.ResponseMessage, error) {
	maxItems := 1000
	itemsToGet := maxItems

	// TODO: may be an edge case
	if task.End.Sub(task.Start) < time.Hour*time.Duration(maxItems) {
		itemsToGet = int(task.End.Sub(task.Start).Hours()) + 1
	}

	restClient := binance.NewFuturesClient("", "")
	klines, err := restClient.NewKlinesService().Symbol(task.Market).Limit(itemsToGet).
		Interval("1h").EndTime(task.End.UnixMilli()).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	start_datetime := task.Start
	end_datetime := task.End

	if len(klines) > 0 {
		dirty_start_datetime := time.UnixMilli(klines[0].OpenTime).UTC()
		start_datetime = time.Date(
			dirty_start_datetime.Year(),
			dirty_start_datetime.Month(),
			dirty_start_datetime.Day(),
			dirty_start_datetime.Hour(),
			0, 0, 0, time.UTC,
		)

		dirty_end_datetime := time.UnixMilli(klines[len(klines)-1].OpenTime).UTC()
		end_datetime = time.Date(
			dirty_end_datetime.Year(),
			dirty_end_datetime.Month(),
			dirty_end_datetime.Day(),
			dirty_end_datetime.Hour(),
			0, 0, 0, time.UTC,
		)
	}

	msg := restc.ResponseMessage{
		Task:  task,
		Start: start_datetime,
		End:   end_datetime,
		Data:  klines,
		Last:  len(klines) < maxItems,
	}

	return &msg, nil
}
