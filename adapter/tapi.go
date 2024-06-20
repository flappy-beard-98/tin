package adapter

import (
	"context"
	_ "embed"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"os/signal"
	"syscall"
)

type TApi struct {
	*investgo.Client
	cancel context.CancelFunc
}

func NewTapi(cfg investgo.Config, logger Logger) (*TApi, error) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	client, err := investgo.NewClient(ctx, cfg, logger)

	if err != nil {
		return nil, err
	}
	return &TApi{
		Client: client,
		cancel: cancel,
	}, nil
}

func (o *TApi) Close() {
	_ = o.Client.Stop()
	o.cancel()
}
