package main

import (
	"context"
	_ "embed"
	"os/signal"
	"syscall"
	"time"
	"tin/analyzer"
	"tin/collector"
	"tin/core"
)

var (
	config     = "config.yaml"
	dbFileName = ".temp/invests.sqlite3"
	now        = time.Now()
	past       = now.AddDate(-5, 0, 0)
	currency   = "rub"
	money      = 300_000.0
)

//go:embed apikey.txt
var apiKey string

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	logger, err := core.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	db, err := core.NewSqliteDb(dbFileName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	api, err := core.NewApi(ctx, config, apiKey, logger.Get())
	if err != nil {
		panic(err)
	}
	defer api.Close()

	c := collector.New(db.Get(), api.Get(), logger.Get())

	c.Schema(ctx, true)
	c.ImportAccounts(ctx)
	c.ImportPortfolio(ctx)
	c.ImportShares(ctx)
	c.ImportLastPrices(ctx, currency)
	c.ImportHistoricCandles(ctx, currency, past, now)
	c.ImportDividends(ctx, currency, past, now.AddDate(2, 0, 0))

	a := analyzer.New(db.Get(), logger.Get())

	a.Schema(ctx, true)
	a.HuntForDividends(ctx, money, 10)
}
