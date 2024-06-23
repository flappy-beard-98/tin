package collector

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"time"
	"tin/collector/accounts"
	"tin/collector/dividends"
	"tin/collector/historic_candles"
	"tin/collector/last_prices"
	"tin/collector/portfolio"
	"tin/collector/shares"
)

type Collector struct {
	db     *sqlx.DB
	logger *zap.Logger
	client *investgo.Client
}

func New(db *sqlx.DB, client *investgo.Client, logger *zap.Logger) *Collector {
	return &Collector{
		db:     db,
		logger: logger,
		client: client,
	}
}

func (o *Collector) Schema(ctx context.Context, drop bool) {
	o.logger.Info("schema", zap.Bool("drop", drop))

	if err := shares.NewSchema(o.db, o.logger).Execute(ctx, drop); err != nil {
		o.logger.Error("shares, error", zap.Error(err))
	} else {
		o.logger.Info("shares, schema completed")
	}

	if err := last_prices.NewSchema(o.db, o.logger).Execute(ctx, drop); err != nil {
		o.logger.Error("last prices, error", zap.Error(err))
	} else {
		o.logger.Info("last prices, schema completed")
	}

	if err := dividends.NewSchema(o.db, o.logger).Execute(ctx, drop); err != nil {
		o.logger.Error("dividends, error", zap.Error(err))
	} else {
		o.logger.Info("dividends, schema completed")
	}

	if err := accounts.NewSchema(o.db, o.logger).Execute(ctx, drop); err != nil {
		o.logger.Error("accounts, error", zap.Error(err))
	} else {
		o.logger.Info("accounts, schema completed")
	}

	if err := portfolio.NewSchema(o.db, o.logger).Execute(ctx, drop); err != nil {
		o.logger.Error("portfolio, error", zap.Error(err))
	} else {
		o.logger.Info("portfolio, schema completed")
	}

	if err := historic_candles.NewSchema(o.db, o.logger).Execute(ctx, drop); err != nil {
		o.logger.Error("historic candles, error", zap.Error(err))
	} else {
		o.logger.Info("historic candles, schema completed")
	}

}

func (o *Collector) ImportPortfolio(ctx context.Context) {

	o.logger.Info("import portfolio")

	a, err := accounts.NewRead(o.db, o.logger).AccountIds(ctx)

	if err != nil {
		o.logger.Error("import portfolio, error", zap.Error(err))
		return
	} else {
		o.logger.Info("import portfolio, got account ids")
	}

	if err = portfolio.NewSave(o.db, o.client, o.logger).Execute(ctx, a); err != nil {
		o.logger.Error("import portfolio, error", zap.Error(err))
	} else {
		o.logger.Info("import portfolio, data received and saved")
	}
}

func (o *Collector) ImportAccounts(ctx context.Context) {

	o.logger.Info("import accounts")

	if err := accounts.NewSave(o.db, o.client, o.logger).Execute(ctx); err != nil {
		o.logger.Error("import accounts, error", zap.Error(err))
	} else {
		o.logger.Info("import accounts, data received and saved")
	}
}

func (o *Collector) ImportShares(ctx context.Context) {

	o.logger.Info("import shares")

	if err := shares.NewSave(o.db, o.client, o.logger).Execute(ctx); err != nil {
		o.logger.Error("import shares, error", zap.Error(err))
	} else {
		o.logger.Info("import shares, data received and saved")
	}
}

func (o *Collector) ImportLastPrices(ctx context.Context, currency string) {

	o.logger.Info("import last prices")

	s, err := shares.NewRead(o.db, o.logger).SharesByCurrency(ctx, currency)

	if err != nil {
		o.logger.Error("import last prices, error", zap.Error(err))
		return
	} else {
		o.logger.Info("import last prices, got shares by currency",
			zap.String("currency", currency),
			zap.Int("count", len(s)))
	}

	if err = last_prices.NewSave(o.db, o.client, o.logger).Execute(ctx, s.GetUids()); err != nil {
		o.logger.Error("import last prices, error", zap.Error(err))
	} else {
		o.logger.Info("import last prices, data received and saved")
	}
}

func (o *Collector) ImportDividends(ctx context.Context, currency string, from time.Time, to time.Time) {

	o.logger.Info("import dividends")

	s, err := shares.NewRead(o.db, o.logger).SharesByCurrency(ctx, currency)

	if err != nil {
		o.logger.Error("import dividends, error", zap.Error(err))
		return
	} else {
		o.logger.Info("import dividends, got shares by currency",
			zap.String("currency", currency),
			zap.Int("count", len(s)))
	}

	if err = dividends.NewSave(o.db, o.client, o.logger).Execute(ctx, from, to, s.GetFigis()); err != nil {
		o.logger.Error("import dividends, error", zap.Error(err))
	} else {
		o.logger.Info("import dividends, data received and saved")
	}
}

func (o *Collector) ImportHistoricCandles(ctx context.Context, currency string, from time.Time, to time.Time) {

	o.logger.Info("import historic candles")

	s, err := shares.NewRead(o.db, o.logger).SharesByCurrency(ctx, currency)

	if err != nil {
		o.logger.Error("import historic candles, error", zap.Error(err))
		return
	} else {
		o.logger.Info("import historic candles, got shares by currency",
			zap.String("currency", currency),
			zap.Int("count", len(s)))
	}

	if err = historic_candles.NewSave(o.db, o.client, o.logger).Execute(ctx, from, to, s.GetUids()); err != nil {
		o.logger.Error("import historic candles, error", zap.Error(err))
	} else {
		o.logger.Info("import historic candles, data received and saved")
	}
}
