package accounts

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

//go:embed read.sql
var read string

func (o *Read) AccountIds(ctx context.Context) ([]string, error) {
	result := make([]string, 0)
	err := o.db.SelectContext(ctx, &result, read)
	return result, err
}
