package core

import (
	"context"
	_ "embed"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
)

type Api struct {
	client *investgo.Client
	cancel context.CancelFunc
}

func NewApi(ctx context.Context, config string, token string, logger *zap.Logger) (*Api, error) {
	cfg, err := investgo.LoadConfig(config)
	if err != nil {
		return nil, err
	}

	cfg.Token = token

	client, err := investgo.NewClient(ctx, cfg, logger.Sugar())

	if err != nil {
		return nil, err
	}
	return &Api{
		client: client,
	}, nil
}

func (o *Api) Get() *investgo.Client {
	return o.client
}

func (o *Api) Close() {
	err := o.client.Stop()
	if err != nil {
		println(err.Error())
	}
}
