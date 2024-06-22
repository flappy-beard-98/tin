package accounts

import (
	"context"
	_ "embed"
	"github.com/jmoiron/sqlx"
)

type Schema struct {
	db *sqlx.DB
}

func NewSchema(db *sqlx.DB) *Schema {
	return &Schema{db}
}

//go:embed schema.sql
var schema string

func (o *Schema) Execute(ctx context.Context, drop bool) error {
	var err error
	if drop {
		_, err = o.db.ExecContext(ctx, `drop table if exists collector_accounts;`)
	}
	if err != nil {
		return err
	}
	_, err = o.db.ExecContext(ctx, schema)
	return err
}
