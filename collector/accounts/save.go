package accounts

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type Save struct {
	db  *sqlx.DB
	api *investgo.Client
}

func NewSave(db *sqlx.DB, api *investgo.Client) *Save {
	return &Save{db, api}
}

func (o *Save) Execute(ctx context.Context) error {
	data, err := o.getAccounts()
	if err != nil {
		return err
	}
	err = o.saveAccounts(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (o *Save) getAccounts() ([]*investapi.Account, error) {
	service := o.api.NewUsersServiceClient()

	response, err := service.GetAccounts()

	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("empty accounts response")
	}

	return response.Accounts, nil
}

//go:embed save.sql
var save string

func (o *Save) saveAccounts(ctx context.Context, data []*investapi.Account) error {
	for _, v := range data {
		_, err := o.db.NamedExecContext(ctx, save,
			map[string]interface{}{
				"id":          v.Id,
				"type":        v.Type,
				"name":        v.Name,
				"status":      v.Status,
				"openeddate":  v.OpenedDate.AsTime().Format("2006-01-02"),
				"closeddate":  v.ClosedDate.AsTime().Format("2006-01-02"),
				"accesslevel": v.Status,
			})

		if err != nil {
			return err
		}
	}
	return nil
}
