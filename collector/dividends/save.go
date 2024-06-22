package dividends

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"time"
)

type Save struct {
	db  *sqlx.DB
	api *investgo.Client
}

func NewSave(db *sqlx.DB, api *investgo.Client) *Save {
	return &Save{db, api}
}

func (o *Save) Execute(ctx context.Context, from time.Time, to time.Time, figis []string) error {
	data, err := o.getDividends(from, to, figis)
	if err != nil {
		return err
	}
	err = o.saveDividends(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (o *Save) getDividends(from time.Time, to time.Time, figis []string) (map[string]*investapi.Dividend, error) {
	service := o.api.NewInstrumentsServiceClient()
	errs := make([]error, 0)
	result := make(map[string]*investapi.Dividend)

	for _, v := range figis {
		response, err := service.GetDividents(v, from, to)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if response == nil {
			errs = append(errs, errors.New("empty shares response"))
			continue
		}
		for _, d := range response.Dividends {
			result[v] = d
		}
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return result, nil
}

//go:embed save.sql
var save string

func (o *Save) saveDividends(ctx context.Context, data map[string]*investapi.Dividend) error {
	for f, v := range data {
		_, err := o.db.NamedExecContext(ctx, save,
			map[string]interface{}{
				"figi":         f,
				"dividendnet":  v.DividendNet.ToFloat(),
				"paymentdate":  v.PaymentDate.AsTime().Format("2006-01-02"),
				"declareddate": v.DeclaredDate.AsTime().Format("2006-01-02"),
				"lastbuydate":  v.LastBuyDate.AsTime().Format("2006-01-02"),
				"dividendtype": v.DividendType,
				"recorddate":   v.RecordDate.AsTime().Format("2006-01-02"),
				"regularity":   v.Regularity,
				"closeprice":   v.ClosePrice.ToFloat(),
				"yieldvalue":   v.YieldValue.ToFloat(),
				"createdat":    v.CreatedAt.AsTime().Format("2006-01-02"),
			})
		if err != nil {
			return err
		}
	}
	return nil
}
