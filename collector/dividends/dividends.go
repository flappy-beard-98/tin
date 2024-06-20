package dividends

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"time"
	"tinkoff/adapter"
	"tinkoff/core"
)

// Dividend Информация о выплате.
type Dividend struct {
	Figi         string  `json:"figi,omitempty"`          //Figi-идентификатор инструмента.
	DividendNet  float64 `json:"dividend_net,omitempty"`  //Величина дивиденда на 1 ценную бумагу (включая валюту).
	PaymentDate  string  `json:"payment_date,omitempty"`  //Дата фактических выплат в часовом поясе UTC.
	DeclaredDate string  `json:"declared_date,omitempty"` //Дата объявления дивидендов в часовом поясе UTC.
	LastBuyDate  string  `json:"last_buy_date,omitempty"` //Последний день (включительно) покупки для получения выплаты в часовом поясе UTC.
	DividendType string  `json:"dividend_type,omitempty"` //Тип выплаты. Возможные значения: Regular Cash – регулярные выплаты, Cancelled – выплата отменена, Daily Accrual – ежедневное начисление, Return of Capital – возврат капитала, прочие типы выплат.
	RecordDate   string  `json:"record_date,omitempty"`   //Дата фиксации реестра в часовом поясе UTC.
	Regularity   string  `json:"regularity,omitempty"`    //Регулярность выплаты. Возможные значения: Annual – ежегодная, Semi-Anl – каждые полгода, прочие типы выплат.
	ClosePrice   float64 `json:"close_price,omitempty"`   //Цена закрытия инструмента на момент ex_dividend_date.
	YieldValue   float64 `json:"yield_value,omitempty"`   //Величина доходности.
	CreatedAt    string  `json:"created_at,omitempty"`    //Дата и время создания записи в часовом поясе UTC.
}

func New(dividend *investapi.Dividend, figi string) *Dividend {
	return &Dividend{
		Figi:         figi,
		DividendNet:  dividend.DividendNet.ToFloat(),
		PaymentDate:  dividend.PaymentDate.AsTime().Format("2006-01-02"),
		DeclaredDate: dividend.DeclaredDate.AsTime().Format("2006-01-02"),
		LastBuyDate:  dividend.LastBuyDate.AsTime().Format("2006-01-02"),
		DividendType: dividend.DividendType,
		RecordDate:   dividend.RecordDate.AsTime().Format("2006-01-02"),
		Regularity:   dividend.Regularity,
		ClosePrice:   dividend.ClosePrice.ToFloat(),
		YieldValue:   dividend.YieldValue.ToFloat(),
		CreatedAt:    dividend.CreatedAt.AsTime().Format("2006-01-02"),
	}
}

func (o *Dividend) GetPaymentDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.PaymentDate)
	return t
}

func (o *Dividend) GetDeclaredDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.DeclaredDate)
	return t
}

func (o *Dividend) GetLastBuyDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.LastBuyDate)
	return t
}

func (o *Dividend) GetRecordDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.RecordDate)
	return t
}

func (o *Dividend) GetCreatedAt() time.Time {
	t, _ := time.Parse("2006-01-02", o.CreatedAt)
	return t
}

type Dividends struct {
	dividends []Dividend
	figis     []string
	from      time.Time
	to        time.Time
}

func NewDividends(from time.Time, to time.Time, figis ...string) *Dividends {
	return &Dividends{
		dividends: make([]Dividend, 0),
		figis:     figis,
		from:      from,
		to:        to,
	}
}

func (o *Dividends) Get() []Dividend {
	return o.dividends
}

func (o *Dividends) GetFrom() time.Time {
	return o.from
}

func (o *Dividends) GetTo() time.Time {
	return o.to
}

func (o *Dividends) GetFigis() []string {
	return o.figis
}

//go:embed schema.sql
var schema string

func (o *Dividends) Scheme(ctx context.Context, db *adapter.Db) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}

//go:embed insert.sql
var insert string

func (o *Dividends) Insert(ctx context.Context, db *adapter.Db) error {
	for _, s := range o.dividends {
		_, err := db.DB.NamedExecContext(ctx, insert, s)
		if err != nil {
			return err
		}
	}
	return nil
}

//go:embed read.sql
var read string

func (o *Dividends) ReadByLastBuyDateInterval(ctx context.Context, db *adapter.Db, from time.Time, to time.Time) error {
	byLastBuyDateInterval := fmt.Sprintf("%s\nWHERE LastBuyDate BETWEEN ? AND ?", read)
	o.dividends = make([]Dividend, 0)
	o.from = from
	o.to = to
	err := db.SelectContext(ctx, &o.dividends, byLastBuyDateInterval, from, to)

	figis := core.Set[string]{}
	for _, v := range o.dividends {
		figis.Add(v.Figi)
	}

	o.figis = figis.ToArray()
	return err
}

func (o *Dividends) Import(client *adapter.TApi) error {
	service := client.NewInstrumentsServiceClient()
	errs := make([]error, 0)

	for _, v := range o.figis {
		d, err := getDividends(v, o.from, o.to, service)
		if err != nil {
			errs = append(errs, err)
		}
		o.dividends = append(o.dividends, d...)
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}

func getDividends(figi string, from time.Time, to time.Time, client *investgo.InstrumentsServiceClient) ([]Dividend, error) {

	response, err := client.GetDividents(figi, from, to)
	dividends := make([]Dividend, 0)

	if err != nil {
		return dividends, err
	}

	if response == nil {
		return dividends, nil
	}

	for _, v := range response.Dividends {
		dividends = append(dividends, *New(v, figi))
	}
	return dividends, nil
}
