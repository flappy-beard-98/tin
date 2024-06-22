package shares

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

//go:embed read_by_currency.sql
var readByCurrency string

func (o *Read) ByCurrency(ctx context.Context, currency string) (Shares, error) {
	result := make([]Share, 0)
	err := o.db.SelectContext(ctx, &result, readByCurrency, currency)
	return result, err
}

type Share struct {
	Figi string
	Uid  string
}

type Shares []Share

func (o Shares) GetFigis() []string {
	r := make([]string, 0)
	for _, v := range o {
		r = append(r, v.Figi)
	}
	return r
}

func (o Shares) GetUids() []string {
	r := make([]string, 0)
	for _, v := range o {
		r = append(r, v.Uid)
	}
	return r
}
