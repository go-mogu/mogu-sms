package client

import (
	"context"
	baseClient "github.com/go-mogu/mgu-sms/pkg/client"
	"net/url"
)

type searchClient struct{}

var SearchClient = &searchClient{}

const (
	search                      = "mogu-search"
	deleteElasticSearchByUidStr = "/search/deleteElasticSearchByUids"
	addElasticSearchByUid       = "/search/addElasticSearchIndexByUid"
)

func (c *searchClient) DeleteElasticSearchByUidStr(ctx context.Context, uidStr string) (result []byte, err error) {
	result, err = baseClient.Post(ctx, search, deleteElasticSearchByUidStr, url.Values{
		"uid": []string{uidStr},
	}, map[string]string{
		"uid": uidStr,
	}, nil)
	return
}

func (c *searchClient) AddElasticSearchByUid(ctx context.Context, uidStr string) (result []byte, err error) {
	result, err = baseClient.Post(ctx, search, addElasticSearchByUid, url.Values{
		"uid": []string{uidStr},
	}, map[string]string{
		"uid": uidStr,
	}, nil)
	return
}
