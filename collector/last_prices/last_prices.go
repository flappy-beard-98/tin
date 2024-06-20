package last_prices

import (
	"context"
	_ "embed"
	"errors"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"time"
	"tinkoff/adapter"
)

// LastPrice Информация о цене последней сделки.
type LastPrice struct {
	Figi          string  `json:"figi,omitempty"`           //Figi инструмента.
	Price         float64 `json:"price,omitempty"`          //Цена последней сделки за 1 инструмент. Для получения стоимости лота требуется умножить на лотность инструмента. Для перевод цен в валюту рекомендуем использовать [информацию со страницы](https://tinkoff.github.io/investAPI/faq_marketdata/)
	Time          string  `json:"time,omitempty"`           //Время получения последней цены в часовом поясе UTC по времени биржи.
	InstrumentUid string  `json:"instrument_uid,omitempty"` //Uid инструмента
}

func New(lastPrice *investapi.LastPrice) *LastPrice {
	return &LastPrice{
		Figi:          lastPrice.Figi,
		Price:         lastPrice.Price.ToFloat(),
		Time:          lastPrice.Time.AsTime().Format("2006-01-02"),
		InstrumentUid: lastPrice.InstrumentUid,
	}
}

func (o *LastPrice) GetTime() time.Time {
	t, _ := time.Parse("2006-01-02", o.Time)
	return t
}

type LastPrices struct {
	prices []LastPrice
	uids   []string
}

func NewLastPrices(uids ...string) *LastPrices {
	return &LastPrices{
		prices: make([]LastPrice, 0),
		uids:   uids,
	}
}

func (o *LastPrices) Get() []LastPrice {
	return o.prices
}

func (o *LastPrices) GetUids() []string {
	return o.uids
}

//go:embed schema.sql
var schema string

func (o *LastPrices) Scheme(ctx context.Context, db *adapter.Db) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}

//go:embed insert.sql
var insert string

func (o *LastPrices) Insert(ctx context.Context, db *adapter.Db) error {
	for _, s := range o.prices {
		_, err := db.DB.NamedExecContext(ctx, insert, s)
		if err != nil {
			return err
		}
	}
	return nil
}

//go:embed read.sql
var read string

func (o *LastPrices) Read(ctx context.Context, db *adapter.Db) error {
	o.prices = make([]LastPrice, 0)
	err := db.SelectContext(ctx, &o.prices, read)
	return err
}

func (o *LastPrices) Import(client *adapter.TApi) error {
	service := client.NewMarketDataServiceClient()

	response, err := service.GetLastPrices(o.uids)
	o.prices = make([]LastPrice, 0)

	if err != nil {
		return err
	}

	if response == nil {
		return errors.New("empty last prices response")
	}

	for _, v := range response.LastPrices {
		o.prices = append(o.prices, *New(v))
	}
	return nil
}
