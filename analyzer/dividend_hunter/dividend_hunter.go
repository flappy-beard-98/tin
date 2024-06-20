package dividend_hunter

import (
	"cmp"
	"context"
	_ "embed"
	"fmt"
	"slices"
	"time"
	"tin/adapter"
	"tin/core"
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

type ResultItem struct {
	Tag         string  `db:"tag"`
	FIGI        string  `db:"figi"`
	Ticker      string  `db:"ticker"`
	LastBuyDate string  `db:"lastbuydate"`
	RecordDate  string  `db:"recorddate"`
	PaymentDate string  `db:"paymentdate"`
	Expectation float64 `db:"expectation"`
}

func NewResultItem(v *BaseItem, tag string, balance float64) ResultItem {

	expectation := (balance / v.UnitPrice) * float64(v.Lot) * v.DividendNet
	return ResultItem{
		Tag:         tag,
		FIGI:        v.FIGI,
		Ticker:      v.Ticker,
		LastBuyDate: v.LastBuyDate,
		RecordDate:  v.RecordDate,
		PaymentDate: v.PaymentDate,
		Expectation: expectation,
	}
}

func (o *ResultItem) GetLastBuyDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.LastBuyDate)
	return t
}

func (o *ResultItem) GetRecordDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.PaymentDate)
	return t
}

func (o *ResultItem) GetPaymentDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.PaymentDate)
	return t
}

func (o *ResultItem) String() string {
	return fmt.Sprintf("%-7s %10s %10s %10s %10.2f ", o.Ticker, o.LastBuyDate, o.RecordDate, o.PaymentDate, o.Expectation)
}

type Base struct {
	items      []BaseItem
	results    []ResultItem
	best       []ResultItem
	balance    float64
	topResults int
}

func NewBase(balance float64, topResults int) *Base {
	return &Base{
		items:      make([]BaseItem, 0),
		results:    make([]ResultItem, 0),
		best:       make([]ResultItem, 0),
		balance:    balance,
		topResults: topResults,
	}
}

func (o *Base) Get() []BaseItem {
	return o.items
}

func (o *Base) GetResults() []ResultItem {
	return o.results
}

func (o *Base) GetBest() []ResultItem {
	return o.best
}

func (o *Base) GetBestReport() string {

	s := fmt.Sprintf("\n\n%-7s %10s %10s %10s %10s\n", "Ticker", "LastBuy", "Record", "Payment", "Income")
	total := 0.0
	for _, r := range o.best {
		s = fmt.Sprintf("%s%s\n", s, r.String())
		total += r.Expectation
	}
	s = fmt.Sprintf("%sTOTAL: %10.2f for BALANCE %10.2f\n", s, total, o.balance)
	return s
}

//go:embed schema.sql
var schema string

func (o *Base) Schema(ctx context.Context, db *adapter.Db) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}

//go:embed prepare.sql
var prepare string

func (o *Base) Prepare(ctx context.Context, db *adapter.Db) error {
	o.items = make([]BaseItem, 0)
	o.results = make([]ResultItem, 0)
	o.best = make([]ResultItem, 0)
	_, err := db.ExecContext(ctx, prepare)
	return err
}

//go:embed read.sql
var read string

func (o *Base) Read(ctx context.Context, db *adapter.Db) error {
	o.items = make([]BaseItem, 0)
	o.results = make([]ResultItem, 0)
	o.best = make([]ResultItem, 0)
	err := db.SelectContext(ctx, &o.items, read)
	return err
}

//go:embed best.sql
var best string

func (o *Base) Best(ctx context.Context, db *adapter.Db) error {
	o.best = make([]ResultItem, 0)
	err := db.SelectContext(ctx, &o.best, best)
	return err
}

//go:embed result.sql
var result string

func (o *Base) Results(ctx context.Context, db *adapter.Db) error {
	for _, s := range o.results {
		_, err := db.DB.NamedExecContext(ctx, result, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Base) Analyze() error {

	intervals := make([]core.Interval, 0)

	for _, item := range o.items {
		i := item
		intervals = append(intervals, &i)
	}

	var combos [][]core.Interval
	core.FindNonOverlappingIntervals(intervals, 0, []core.Interval{}, &combos)

	results := make([][]ResultItem, 0)

	for ci, combo := range combos {
		ri := make([]ResultItem, 0)
		for _, item := range combo {
			i := item.(*BaseItem)
			t := fmt.Sprintf("%d", ci)
			r := NewResultItem(i, t, o.balance)
			ri = append(ri, r)
		}
		results = append(results, ri)
	}

	slices.SortFunc(results, func(a, b []ResultItem) int {
		return int(sum(b) - sum(a))
	})

	o.results = make([]ResultItem, 0)

	for _, row := range results[:o.topResults] {

		slices.SortFunc(row, func(a, b ResultItem) int {
			return cmp.Compare(a.GetLastBuyDate().Unix(), b.GetLastBuyDate().Unix())
		})

		o.results = append(o.results, row...)
	}

	o.best = results[0]

	return nil
}

func sum(results []ResultItem) float64 {
	total := 0.0
	for _, r := range results {
		total += r.Expectation
	}
	return total
}
