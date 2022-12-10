package prom_query_tool

import (
	"errors"
	"github.com/prometheus/common/model"
)

type Scalar struct {
	result model.Value
}

func NewScalar(result model.Value) *Scalar {
	return &Scalar{result: result}
}

func (s *Scalar) GetValue() (float64, error) {
	if v, ok := s.result.(*model.Scalar); !ok {
		return 0, errors.New("查询结果不是一个Scalar类型")
	} else {
		return float64(v.Value), nil
	}
}

func (s *Scalar) GetPoint() (Point[float64], error) {
	if v, ok := s.result.(*model.Scalar); !ok {
		return Point[float64]{}, errors.New("查询结果不是一个Scalar类型")
	} else {
		return Point[float64]{
			Timestamp: v.Timestamp.Time(),
			Value:     float64(v.Value),
		}, nil
	}
}
