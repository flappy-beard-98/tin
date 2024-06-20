package analyzer

import (
	"context"
	"tin/adapter"
	"tin/analyzer/dividend_hunter"
)

type Analyzer struct {
	db     *adapter.Db
	logger adapter.Logger
}

func New(dbFileName string) (*Analyzer, error) {
	db, err := adapter.NewSqliteDb(dbFileName)
	if err != nil {
		return nil, err
	}

	logger, err := adapter.NewLogger()
	if err != nil {
		return nil, err
	}
	return &Analyzer{db: db, logger: logger}, nil
}

func (o *Analyzer) HuntForDividends(ctx context.Context, balance float64, topResults int) error {

	o.logger.Infof("hunting for dividends")

	items := dividend_hunter.NewBase(balance, topResults)

	if err := items.Schema(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("hunting for dividends, schema recreated")

	if err := items.Prepare(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("hunting for dividends, prepared")

	if err := items.Read(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("hunting for dividends, loaded")

	if err := items.Analyze(); err != nil {
		return err
	}

	o.logger.Infof("hunting for dividends, analyzed")

	if err := items.Results(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("hunting for dividends, result saved")

	if err := items.Best(ctx, o.db); err != nil {
		return err
	}

	o.logger.Infof("hunting for dividends, best loaded")

	o.logger.Infof("result is %s", items.GetBestReport())

	return nil
}
