package portfolio

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

//go:embed schema_drop.sql
var schemaDrop string

func (o *Schema) Execute(ctx context.Context, drop bool) error {
	var err error
	if drop {
		_, err = o.db.ExecContext(ctx, schemaDrop)
	}
	if err != nil {
		return err
	}
	_, err = o.db.ExecContext(ctx, schema)
	return err
}
