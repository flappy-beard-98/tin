package historic_candles

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"time"
)

type Save struct {
	db  *sqlx.DB
	api *investgo.Client
	log *zap.Logger
}

func NewSave(db *sqlx.DB, api *investgo.Client, log *zap.Logger) *Save {
	return &Save{db, api, log}
}

func (o *Save) Execute(ctx context.Context, from time.Time, to time.Time, uids []string) error {
	data, err := o.getCandles(from, to, uids)
	if err != nil {
		return err
	}
	err = o.saveCandles(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

type candleKey struct {
	uid  string
	time string
}

func (o *Save) getCandles(from time.Time, to time.Time, uids []string) (map[candleKey]*investapi.HistoricCandle, error) {
	service := o.api.NewMarketDataServiceClient()
	errs := make([]error, 0)
	result := make(map[candleKey]*investapi.HistoricCandle)
	ifrom := to.AddDate(-1, 0, 0)
	ito := to

	o.log.Debug("get historic candles", zap.Time("from", from), zap.Time("to", to), zap.Strings("uids", uids))

	for ifrom.After(from) {
		o.log.Debug("get historic candles for sub interval", zap.Time("from", ifrom), zap.Time("to", ito))
		for _, v := range uids {
			response, err := service.GetCandles(v, investapi.CandleInterval_CANDLE_INTERVAL_DAY, ifrom, ito, investapi.GetCandlesRequest_CANDLE_SOURCE_UNSPECIFIED)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if response == nil {
				errs = append(errs, errors.New("empty shares response"))
				continue
			}
			o.log.Debug("got historic candles", zap.String("uid", v))
			for _, c := range response.Candles {
				k := candleKey{
					uid:  v,
					time: c.Time.AsTime().Format("2006-01-02"),
				}
				result[k] = c
			}
		}
		ito = ifrom
		ifrom = ifrom.AddDate(-1, 0, 0)
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return result, nil
}

//go:embed save.sql
var save string

func (o *Save) saveCandles(ctx context.Context, data map[candleKey]*investapi.HistoricCandle) error {
	i := 0
	o.log.Debug("save historic candles")
	for k, v := range data {
		if i%1000 == 0 {
			o.log.Debug("saved historic candles", zap.Int("count", i))
		}
		_, err := o.db.NamedExecContext(ctx, save,
			map[string]interface{}{
				"uid":          k.uid,
				"open":         v.Open.ToFloat(),
				"high":         v.High.ToFloat(),
				"low":          v.Low.ToFloat(),
				"close":        v.Close.ToFloat(),
				"volume":       v.Volume,
				"time":         k.time,
				"iscomplete":   v.IsComplete,
				"candlesource": v.CandleSource.String(),
			})

		if err != nil {
			return err
		}
		i++
	}
	o.log.Debug("historic candles saved")
	return nil
}
