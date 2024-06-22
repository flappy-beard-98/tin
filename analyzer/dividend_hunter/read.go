package dividend_hunter

import (
	"context"
	_ "embed"
	"github.com/jmoiron/sqlx"
)

type Read struct {
	db *sqlx.DB
}

func NewRead(db *sqlx.DB) *Read {
	return &Read{db}
}

//go:embed read_best.sql
var readBest string

func (o *Read) Best(ctx context.Context) (Results, error) {
	result := make([]ResultItem, 0)
	err := o.db.SelectContext(ctx, &result, readBest)
	return result, err
}
