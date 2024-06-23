package dividend_hunter

import (
	"context"
	_ "embed"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Read struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewRead(db *sqlx.DB, log *zap.Logger) *Read {
	return &Read{db, log}
}

//go:embed read_best.sql
var readBest string

func (o *Read) Best(ctx context.Context) (Results, error) {
	o.log.Debug("read best result of dividend hunting")
	result := make([]ResultItem, 0)
	err := o.db.SelectContext(ctx, &result, readBest)
	return result, err
}
