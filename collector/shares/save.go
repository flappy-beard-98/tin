package shares

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

func (o *Save) Execute(ctx context.Context) error {
	data, err := o.getShares()
	if err != nil {
		return err
	}
	err = o.saveShares(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (o *Save) getShares() ([]*investapi.Share, error) {
	o.log.Debug("get shares")

	service := o.api.NewInstrumentsServiceClient()

	response, err := service.Shares(investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE)

	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("empty shares response")
	}

	return response.Instruments, nil
}

//go:embed save.sql
var save string

func (o *Save) saveShares(ctx context.Context, data []*investapi.Share) error {
	o.log.Debug("save shares")
	for _, v := range data {
		_, err := o.db.NamedExecContext(ctx, save,
			map[string]interface{}{
				"figi":                  v.Figi,
				"ticker":                v.Ticker,
				"classcode":             v.ClassCode,
				"isin":                  v.Isin,
				"lot":                   v.Lot,
				"currency":              v.Currency,
				"klong":                 v.Klong.ToFloat(),
				"kshort":                v.Kshort.ToFloat(),
				"dlong":                 v.Dlong.ToFloat(),
				"dshort":                v.Dshort.ToFloat(),
				"dlongmin":              v.DlongMin.ToFloat(),
				"dshortmin":             v.DshortMin.ToFloat(),
				"shortenabledflag":      v.ShortEnabledFlag,
				"name":                  v.Name,
				"exchange":              v.Exchange,
				"ipodate":               v.IpoDate.AsTime().Format("2006-01-02"),
				"issuesize":             v.IssueSize,
				"countryofrisk":         v.CountryOfRisk,
				"countryofriskname":     v.CountryOfRiskName,
				"sector":                v.Sector,
				"issuesizeplan":         v.IssueSizePlan,
				"nominal":               v.Nominal.ToFloat(),
				"tradingstatus":         v.TradingStatus.String(),
				"otcflag":               v.OtcFlag,
				"buyavailableflag":      v.BuyAvailableFlag,
				"sellavailableflag":     v.SellAvailableFlag,
				"divyieldflag":          v.DivYieldFlag,
				"sharetype":             v.ShareType.String(),
				"minpriceincrement":     v.MinPriceIncrement.ToFloat(),
				"apitradeavailableflag": v.ApiTradeAvailableFlag,
				"uid":                   v.Uid,
				"realexchange":          v.RealExchange.String(),
				"positionuid":           v.PositionUid,
				"foriisflag":            v.ForIisFlag,
				"forqualinvestorflag":   v.ForQualInvestorFlag,
				"weekendflag":           v.WeekendFlag,
				"blockedtcaflag":        v.BlockedTcaFlag,
				"liquidityflag":         v.LiquidityFlag,
				"first_1mincandledate":  v.First_1MinCandleDate.AsTime().Format("2006-01-02"),
				"first_1daycandledate":  v.First_1DayCandleDate.AsTime().Format("2006-01-02"),
			})

		if err != nil {
			return err
		}
	}
	o.log.Debug("shares saved")
	return nil
}
