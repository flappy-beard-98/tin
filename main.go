package main

import (
	"context"
	"fmt"
	"time"
	"tinkoff/analyzer"
	"tinkoff/collector"
)

var (
	config   = "config.yaml"
	db       = ".temp/invests.sqlite3"
	from     = time.Now().AddDate(-10, 0, 0)
	to       = time.Now().AddDate(2, 0, 0)
	currency = "rub"
	money    = 300_000.0
)

func main() {

	ctx := context.Background()

	c, err := collector.New(db, config)
	if err != nil || c == nil {
		panic(fmt.Sprintf("collector not created or error occured : %v", err))
	}

	defer c.Close()

	if err = c.ImportShares(ctx); err != nil {
		panic(err)
	}

	if err = c.ImportLastPrices(ctx, currency); err != nil {
		panic(err)
	}

	if err = c.ImportDividends(ctx, currency, from, to); err != nil {
		panic(err)
	}

	a, err := analyzer.New(db)
	if err != nil || a == nil {
		panic(fmt.Sprintf("analyzer not created or error occured : %v", err))
	}

	if err = a.HuntForDividends(ctx, money, 10); err != nil {
		panic(err)
	}
}
