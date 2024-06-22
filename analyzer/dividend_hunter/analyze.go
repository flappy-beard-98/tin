package dividend_hunter

import (
	"cmp"
	"context"
	_ "embed"
	"github.com/jmoiron/sqlx"
	"slices"
	"strconv"
	"tin/core"
)

type Analyze struct {
	db *sqlx.DB
}

func NewAnalyze(db *sqlx.DB) *Analyze {
	return &Analyze{db}
}

//go:embed analyze_prepare.sql
var prepare string

//go:embed analyze_save.sql
var save string

func (o *Analyze) Execute(ctx context.Context, balance float64, bestResultsCount int) error {

	items, err := o.prepare(ctx)
	if err != nil {
		return err
	}

	results := o.analyze(balance, bestResultsCount, items)

	err = o.save(ctx, results)

	return err
}

func (o *Analyze) analyze(balance float64, bestResultsCount int, items []BaseItem) []ResultItem {

	intervals := make([]core.Interval, 0)

	for _, item := range items {
		i := item
		intervals = append(intervals, &i)
	}

	var combos [][]core.Interval
	core.FindNonOverlappingIntervals(intervals, 0, []core.Interval{}, &combos)

	resultCombos := make([][]ResultItem, 0)

	for ci, combo := range combos {
		ri := make([]ResultItem, 0)
		for _, item := range combo {
			r := item.(*BaseItem).ToResultItem(strconv.Itoa(ci), balance)
			ri = append(ri, r)
		}
		resultCombos = append(resultCombos, ri)
	}

	slices.SortFunc(resultCombos, func(a, b []ResultItem) int {
		return int(sum(b) - sum(a))
	})

	results := make([]ResultItem, 0)
	if bestResultsCount < len(resultCombos) {
		resultCombos = resultCombos[:bestResultsCount]
	}

	for i, row := range resultCombos {

		if i == 0 {
			for j := 0; j < len(row); j++ {
				row[j].Tag = "best"
			}
		}

		slices.SortFunc(row, func(a, b ResultItem) int {
			return cmp.Compare(a.GetLastBuyDate().Unix(), b.GetLastBuyDate().Unix())
		})

		results = append(results, row...)
	}

	return results
}

func sum(results []ResultItem) float64 {
	total := 0.0
	for _, r := range results {
		total += r.Expectation
	}
	return total
}

func (o *Analyze) prepare(ctx context.Context) ([]BaseItem, error) {
	result := make([]BaseItem, 0)
	err := o.db.SelectContext(ctx, &result, prepare)
	return result, err
}

func (o *Analyze) save(ctx context.Context, data []ResultItem) error {
	for _, v := range data {
		_, err := o.db.NamedExecContext(ctx, save, v)
		if err != nil {
			return err
		}
	}
	return nil
}
