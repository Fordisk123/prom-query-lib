package prom_query_tool

import (
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

type String struct {
	result model.Value
}

func NewString(result model.Value) *Scalar {
	return &Scalar{result: result}
}

func (s *String) GetValue() (string, error) {
	if v, ok := s.result.(*model.String); !ok {
		return "", errors.New("查询结果不是一个String类型")
	} else {
		return v.Value, nil
	}
}

func (s *String) GetPoint() (Point[string], error) {
	if v, ok := s.result.(*model.String); !ok {
		return Point[string]{}, errors.New("查询结果不是一个String类型")
	} else {
		return Point[string]{
			Timestamp: v.Timestamp.Time(),
			Value:     v.Value,
		}, nil
	}
}
