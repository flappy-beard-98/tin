package accounts

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

//go:embed read.sql
var read string

func (o *Read) AccountIds(ctx context.Context) ([]string, error) {
	o.log.Debug("get account ids")
	result := make([]string, 0)
	err := o.db.SelectContext(ctx, &result, read)
	return result, err
}
