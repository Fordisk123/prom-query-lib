package prom_query_tool

import (
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

func NewMatrix[L comparable, V comparable](result model.Value) *Matrix[L, V] {
	return &Matrix[L, V]{Result: result}
}

func NewMatrixWithDefaultValue[L comparable](result model.Value) *Matrix[L, float64] {
	v := &Matrix[L, float64]{Result: result}
	v.SetValueConvertFunc(func(v model.SampleValue) float64 {
		return float64(v)
	})
	return v
}

type Matrix[L comparable, V comparable] struct {
	Result           model.Value
	labelCovertFunc  func(metric model.Metric) L
	valueConvertFunc func(value model.SampleValue) V
}

func (v *Matrix[L, V]) SetLabelCovertFunc(f func(metric model.Metric) L) *Matrix[L, V] {
	v.labelCovertFunc = f
	return v
}

func (v *Matrix[L, V]) SetValueConvertFunc(f func(value model.SampleValue) V) *Matrix[L, V] {
	v.valueConvertFunc = f
	return v
}

func (v *Matrix[L, V]) ToMap() (map[L][]V, error) {
	if matrix, ok := v.Result.(model.Matrix); !ok {
		return nil, errors.New("查询结果不是一个Matrix类型")
	} else {
		resultMap := make(map[L][]V)
		for i := range matrix {
			s := matrix[i]
			label := v.labelCovertFunc(s.Metric)
			resultMap[label] = make([]V, 0)
			for j := range s.Values {
				sv := s.Values[j]
				resultMap[label] = append(resultMap[label], v.valueConvertFunc(sv.Value))
			}
		}
		return resultMap, nil
	}
}

func (v *Matrix[L, V]) ToPointMap() (map[L][]Point[V], error) {
	if matrix, ok := v.Result.(model.Matrix); !ok {
		return nil, errors.New("查询结果不是一个Matrix类型")
	} else {
		resultMap := make(map[L][]Point[V])
		for i := range matrix {
			s := matrix[i]
			label := v.labelCovertFunc(s.Metric)
			resultMap[label] = make([]Point[V], 0)
			for j := range s.Values {
				sv := s.Values[j]
				resultMap[label] = append(resultMap[label], Point[V]{
					Timestamp: sv.Timestamp.Time(),
					Value:     v.valueConvertFunc(sv.Value),
				})
			}
		}
		return resultMap, nil
	}
}
