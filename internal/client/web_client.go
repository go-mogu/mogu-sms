package client

import (
	"context"
	baseClient "github.com/go-mogu/mgu-sms/pkg/client"
)

type webClient struct{}

var WebClient = &webClient{}

const (
	web            = "mogu-web"
	getSearchModel = "/search/getSearchModel"
)

func (c *webClient) GetSearchModel(ctx context.Context) (result []byte, err error) {
	result, err = baseClient.Get(ctx, web, getSearchModel, nil)
	return
}
