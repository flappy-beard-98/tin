package portfolio

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
)

type Save struct {
	db  *sqlx.DB
	api *investgo.Client
	log *zap.Logger
}

func NewSave(db *sqlx.DB, api *investgo.Client, log *zap.Logger) *Save {
	return &Save{db, api, log}
}

func (o *Save) Execute(ctx context.Context, accountIds []string) error {
	data, err := o.getPortfolio(accountIds)
	if err != nil {
		return err
	}
	err = o.savePortfolio(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (o *Save) getPortfolio(accountIds []string) ([]*investapi.PortfolioResponse, error) {
	o.log.Debug("get portfolios", zap.Strings("accountIds", accountIds))

	service := o.api.NewOperationsServiceClient()
	errs := make([]error, 0)
	result := make([]*investapi.PortfolioResponse, 0)

	for _, v := range accountIds {
		response, err := service.GetPortfolio(v, investapi.PortfolioRequest_RUB)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if response == nil {
			errs = append(errs, errors.New("empty shares response"))
			continue
		}

		o.log.Debug("got portfolio", zap.String("accountId", v))

		result = append(result, response.PortfolioResponse)
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return result, nil
}

//go:embed save_portfolio.sql
var savePortfolio string

//go:embed save_positions.sql
var savePositions string

//go:embed save_positions_virtual.sql
var savePositionsVirtual string

func (o *Save) savePortfolio(ctx context.Context, data []*investapi.PortfolioResponse) error {
	o.log.Debug("save portfolios")
	for _, v := range data {
		errs := make([]error, 0)

		tx, err := o.db.BeginTxx(ctx, nil)

		if err != nil {
			continue
		}

		o.log.Debug("save portfolio", zap.String("accountId", v.AccountId))
		_, err = tx.NamedExecContext(ctx, savePortfolio, portfolioToMap(v))
		if err != nil {
			errs = append(errs, err)
		}

		o.log.Debug("save virtual positions", zap.String("accountId", v.AccountId))
		for _, vp := range v.VirtualPositions {
			if _, err = tx.NamedExecContext(ctx, savePositionsVirtual, virtualPositionToMap(v.AccountId, vp)); err != nil {
				errs = append(errs, err)
				continue
			}
		}

		o.log.Debug("save positions", zap.String("accountId", v.AccountId))
		for _, p := range v.Positions {
			if _, err = tx.NamedExecContext(ctx, savePositions, positionToMap(v.AccountId, p)); err != nil {
				errs = append(errs, err)
				continue
			}
		}

		if len(errs) != 0 {
			err = errors.Join(errs...)
		}

		if err != nil {
			if e := tx.Rollback(); e != nil {
				return e
			}
			return err
		} else {
			if e := tx.Commit(); e != nil {
				return e
			}
		}
	}
	o.log.Debug("portfolios saved")
	return nil
}

func portfolioToMap(v *investapi.PortfolioResponse) map[string]interface{} {
	return map[string]interface{}{
		"totalamountshares":     v.TotalAmountShares.ToFloat(),
		"totalamountbonds":      v.TotalAmountBonds.ToFloat(),
		"totalamountetf":        v.TotalAmountEtf.ToFloat(),
		"totalamountcurrencies": v.TotalAmountCurrencies.ToFloat(),
		"totalamountfutures":    v.TotalAmountFutures.ToFloat(),
		"expectedyield":         v.ExpectedYield.ToFloat(),
		"accountid":             v.AccountId,
		"totalamountoptions":    v.TotalAmountOptions.ToFloat(),
		"totalamountsp":         v.TotalAmountSp.ToFloat(),
		"totalamountportfolio":  v.TotalAmountPortfolio.ToFloat(),
	}
}

func virtualPositionToMap(accountId string, v *investapi.VirtualPortfolioPosition) map[string]interface{} {
	return map[string]interface{}{
		"accountid":                accountId,
		"positionuid":              v.PositionUid,
		"instrumentuid":            v.InstrumentUid,
		"figi":                     v.Figi,
		"instrumenttype":           v.InstrumentType,
		"quantity":                 v.Quantity.ToFloat(),
		"averagepositionprice":     v.AveragePositionPrice.ToFloat(),
		"expectedyield":            v.ExpectedYield.ToFloat(),
		"expectedyieldfifo":        v.ExpectedYieldFifo.ToFloat(),
		"expiredate":               v.ExpireDate.AsTime().Format("2006-01-02"),
		"currentprice":             v.CurrentPrice.ToFloat(),
		"averagepositionpricefifo": v.AveragePositionPriceFifo.ToFloat(),
	}
}

func positionToMap(accountId string, v *investapi.PortfolioPosition) map[string]interface{} {
	return map[string]interface{}{
		"accountid":                accountId,
		"figi":                     v.Figi,
		"instrumenttype":           v.InstrumentType,
		"quantity":                 v.Quantity.ToFloat(),
		"averagepositionprice":     v.AveragePositionPrice.ToFloat(),
		"expectedyield":            v.ExpectedYield.ToFloat(),
		"currentnkd":               v.CurrentNkd.ToFloat(),
		"currentprice":             v.CurrentPrice.ToFloat(),
		"averagepositionpricefifo": v.AveragePositionPriceFifo.ToFloat(),
		"blocked":                  v.Blocked,
		"blockedlots":              v.BlockedLots.ToFloat(),
		"positionuid":              v.PositionUid,
		"instrumentuid":            v.InstrumentUid,
		"varmargin":                v.VarMargin.ToFloat(),
		"expectedyieldfifo":        v.ExpectedYieldFifo.ToFloat(),
	}
}
