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
		var closed *string = nil
		if v.ClosedDate != nil && v.ClosedDate.AsTime().Unix() > 0 {
			*closed = v.ClosedDate.AsTime().Format("2006-01-02")
		}
		_, err := o.db.NamedExecContext(ctx, save,
			map[string]interface{}{
				"id":          v.Id,
				"type":        v.Type.String(),
				"name":        v.Name,
				"status":      v.Status.String(),
				"openeddate":  v.OpenedDate.AsTime().Format("2006-01-02"),
				"closeddate":  closed,
				"accesslevel": v.Status.String(),
			})

		if err != nil {
			return err
		}
	}
	return nil
}
