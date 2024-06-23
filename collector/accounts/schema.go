package accounts

import (
	"context"
	_ "embed"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Schema struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewSchema(db *sqlx.DB, log *zap.Logger) *Schema {
	return &Schema{db, log}
}

//go:embed schema.sql
var schema string

//go:embed schema_drop.sql
var schemaDrop string

func (o *Schema) Execute(ctx context.Context, drop bool) error {
	var err error
	if drop {
		o.log.Debug("drop schema")
		_, err = o.db.ExecContext(ctx, schemaDrop)
	}
	if err != nil {
		return err
	}
	o.log.Debug("create schema")
	_, err = o.db.ExecContext(ctx, schema)
	return err
}
