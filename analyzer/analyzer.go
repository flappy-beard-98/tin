package analyzer

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"tin/analyzer/dividend_hunter"
)

type Analyzer struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func New(db *sqlx.DB, logger *zap.Logger) *Analyzer {
	return &Analyzer{db: db, logger: logger}
}

func (o *Analyzer) Schema(ctx context.Context, drop bool) {
	o.logger.Info("schema", zap.Bool("drop", drop))

	if err := dividend_hunter.NewSchema(o.db, o.logger).Execute(ctx, drop); err != nil {
		o.logger.Error("dividend hunter, error", zap.Error(err))
	} else {
		o.logger.Info("dividend hunter, schema completed")
	}
}

func (o *Analyzer) HuntForDividends(ctx context.Context, balance float64, topResults int) {

	o.logger.Info("hunting for dividends")

	hunter := dividend_hunter.NewAnalyze(o.db, o.logger)

	if err := hunter.Execute(ctx, balance, topResults); err != nil {
		o.logger.Error("dividend hunter, analyze, error", zap.Error(err))
	} else {
		o.logger.Info("hunting for dividends, analyzed")
	}

	read := dividend_hunter.NewRead(o.db, o.logger)

	if r, err := read.Best(ctx); err != nil {
		o.logger.Error("dividend hunter, read, error", zap.Error(err))
	} else {
		o.logger.Info("dividend hunter, read complete" + r.GetReport(balance))
	}
}
