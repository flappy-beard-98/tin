package dividend_hunter

import (
	"fmt"
	"time"
)

type BaseItem struct {
	FIGI        string  `db:"figi"`
	DividendNet float64 `db:"dividendnet"`
	LastBuyDate string  `db:"lastbuydate"`
	RecordDate  string  `db:"recorddate"`
	PaymentDate string  `db:"paymentdate"`
	Price       float64 `db:"price"`
	Lot         int     `db:"lot"`
	UnitPrice   float64 `db:"unitprice"`
	Ticker      string  `db:"ticker"`
	Name        string  `db:"name"`
}

func (o *BaseItem) GetLastBuyDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.LastBuyDate)
	return t
}

func (o *BaseItem) GetRecordDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.PaymentDate)
	return t
}

func (o *BaseItem) GetPaymentDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.PaymentDate)
	return t
}

func (o *BaseItem) GetID() string {
	return o.FIGI
}

func (o *BaseItem) GetFrom() time.Time {
	return o.GetLastBuyDate()
}

func (o *BaseItem) GetTo() time.Time {
	return o.GetRecordDate()
}

func (o *BaseItem) ToResultItem(tag string, balance float64) ResultItem {
	expectation := (balance / o.UnitPrice) * float64(o.Lot) * o.DividendNet
	return ResultItem{
		Tag:         tag,
		FIGI:        o.FIGI,
		Ticker:      o.Ticker,
		LastBuyDate: o.LastBuyDate,
		RecordDate:  o.RecordDate,
		PaymentDate: o.PaymentDate,
		Expectation: expectation,
	}
}

type ResultItem struct {
	Tag         string  `db:"tag"`
	FIGI        string  `db:"figi"`
	Ticker      string  `db:"ticker"`
	LastBuyDate string  `db:"lastbuydate"`
	RecordDate  string  `db:"recorddate"`
	PaymentDate string  `db:"paymentdate"`
	Expectation float64 `db:"expectation"`
}

func (o *ResultItem) GetLastBuyDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.LastBuyDate)
	return t
}

func (o *ResultItem) String() string {
	return fmt.Sprintf("%-7s %10s %10s %10s %10.2f ", o.Ticker, o.LastBuyDate, o.RecordDate, o.PaymentDate, o.Expectation)
}

type Results []ResultItem

func (o Results) GetReport(balance float64) string {

	s := fmt.Sprintf("\n\n%-7s %10s %10s %10s %10s\n", "Ticker", "LastBuy", "Record", "Payment", "Income")
	total := 0.0
	for _, r := range o {
		s = fmt.Sprintf("%s%s\n", s, r.String())
		total += r.Expectation
	}
	s = fmt.Sprintf("%sTOTAL: %10.2f for BALANCE %10.2f\n", s, total, balance)
	return s
}
