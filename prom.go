package prom_query_tool

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type PromQuery struct {
	api v1.API
}

func NewPromQuery(address string) (*PromQuery, error) {
	//TODO 根据option初始化
	if len(address) == 0 {
		return nil, errors.New("请输入prometheus的地址")
	}

	c, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "创建prometheus客户端出错")
	}

	return &PromQuery{
		api: v1.NewAPI(c),
	}, nil
}

func (p *PromQuery) Query() *Query {
	return &Query{
		api: p.api,
	}
}
