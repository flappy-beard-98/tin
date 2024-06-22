package account

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type Save struct {
	db  *sqlx.DB
	api *investgo.Client
}

func NewSave(db *sqlx.DB, api *investgo.Client) *Save {
	return &Save{db, api}
}

func (o *Save) Execute(ctx context.Context, uids []string) error {
	data, err := o.getLastPrices(uids)
	if err != nil {
		return err
	}
	err = o.saveLastPrices(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (o *Save) getLastPrices(uids []string) ([]*investapi.LastPrice, error) {
	service := o.api.NewMarketDataServiceClient()

	response, err := service.GetLastPrices(uids)

	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("empty last prices response")
	}

	return response.LastPrices, nil
}

//go:embed save.sql
var save string

func (o *Save) saveLastPrices(ctx context.Context, data []*investapi.LastPrice) error {
	for _, v := range data {
		_, err := o.db.NamedExecContext(ctx, save,
			map[string]interface{}{
				"figi":          v.Figi,
				"price":         v.Price.ToFloat(),
				"time":          v.Time.AsTime().Format("2006-01-02"),
				"instrumentuid": v.InstrumentUid,
			})

		if err != nil {
			return err
		}
	}
	return nil
}
