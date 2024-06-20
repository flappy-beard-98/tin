package collector

import (
	"context"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"time"
	"tinkoff/adapter"
	"tinkoff/collector/dividends"
	"tinkoff/collector/last_prices"
	"tinkoff/collector/shares"
)

type Collector struct {
	db     *adapter.Db
	logger adapter.Logger
	tapi   *adapter.TApi
}

func New(dbFileName string, apiConfig investgo.Config) (*Collector, error) {
	db, err := adapter.NewSqliteDb(dbFileName)
	if err != nil {
		return nil, err
	}

	logger, err := adapter.NewLogger()
	if err != nil {
		return nil, err
	}

	tapi, err := adapter.NewTapi(apiConfig, logger)
	if err != nil {
		return nil, err
	}

	return &Collector{
		db:     db,
		logger: logger,
		tapi:   tapi,
	}, nil
}

func (o *Collector) Close() {
	if o.db != nil {
		o.db.Close()
	}
	if o.logger != nil {
		o.logger.Close()
	}
	if o.tapi != nil {
		o.tapi.Close()
	}
}

func (o *Collector) ImportShares(ctx context.Context) error {

	o.logger.Infof("import shares")

	items := share.NewShares()
	if err := items.Scheme(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("import shares, scheme recreated")

	if err := items.Import(o.tapi); err != nil {
		return err
	}

	o.logger.Infof("import shares, data received")

	if err := items.Insert(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("import shares, data saved")

	return nil
}

func (o *Collector) ImportLastPrices(ctx context.Context, currency string) error {

	o.logger.Infof("import last prices")

	shares := share.NewShares()

	if err := shares.ReadByCurrency(ctx, o.db, currency); err != nil {
		return err
	}
	o.logger.Infof("import last prices, got UIDs for %s", currency)

	items := last_prices.NewLastPrices(shares.GetUids()...)
	if err := items.Scheme(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("import last prices, scheme recreated")

	if err := items.Import(o.tapi); err != nil {
		return err
	}

	o.logger.Infof("import last prices, data received")

	if err := items.Insert(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("import last prices, data saved")

	return nil
}

func (o *Collector) ImportDividends(ctx context.Context, currency string, from time.Time, to time.Time) error {
	o.logger.Infof("import dividends")

	shares := share.NewShares()

	if err := shares.ReadByCurrency(ctx, o.db, currency); err != nil {
		return err
	}

	o.logger.Infof("import dividends, got FIGIs for %s", currency)

	items := dividends.NewDividends(from, to, shares.GetFigis()...)
	if err := items.Scheme(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("import dividends, scheme recreated")

	if err := items.Import(o.tapi); err != nil {
		return err
	}

	o.logger.Infof("import dividends, data received")

	if err := items.Insert(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("import dividends, data saved")

	return nil
}
