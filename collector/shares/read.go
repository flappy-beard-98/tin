package shares

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

//go:embed read_by_currency.sql
var readByCurrency string

func (o *Read) SharesByCurrency(ctx context.Context, currency string) (Shares, error) {
	o.log.Debug("get shares by currency")

	result := make([]Share, 0)
	err := o.db.SelectContext(ctx, &result, readByCurrency, currency)
	return result, err
}
