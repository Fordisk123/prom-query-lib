package prom_query_tool

import (
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

func NewVector[L comparable, V any](result model.Value) *Vector[L, V] {
	return &Vector[L, V]{Result: result}
}

func NewVectorWithDefaultValue[L comparable](result model.Value) *Vector[L, float64] {
	v := &Vector[L, float64]{Result: result}
	v.SetValueConvertFunc(func(v model.SampleValue) float64 {
		return float64(v)
	})
	return v
}

type Vector[L comparable, V any] struct {
	Result           model.Value
	labelCovertFunc  func(metric model.Metric) L
	valueConvertFunc func(value model.SampleValue) V
}

func (v *Vector[L, V]) SetLabelCovertFunc(f func(metric model.Metric) L) *Vector[L, V] {
	v.labelCovertFunc = f
	return v
}

func (v *Vector[L, V]) SetValueConvertFunc(f func(value model.SampleValue) V) *Vector[L, V] {
	v.valueConvertFunc = f
	return v
}

func (v *Vector[L, V]) ToMap() (map[L]V, error) {
	if vector, ok := v.Result.(model.Vector); !ok {
		return nil, errors.New("查询结果不是一个Vector类型")
	} else {
		resultMap := make(map[L]V)
		for i := 0; i < len(vector); i++ {
			sample := vector[i]
			if v.labelCovertFunc == nil {
				//TODO 这里是否可以使用反射解决问题？？？
				return nil, errors.New("必须设置一个label转换函数")
			}
			if v.valueConvertFunc == nil {
				return nil, errors.New("必须设置一个值转换函数")
			}
			resultMap[v.labelCovertFunc(sample.Metric)] = v.valueConvertFunc(sample.Value)
		}
		return resultMap, nil
	}
}

func (v *Vector[L, V]) ToList() (*TupleList[L, V], error) {
	if vector, ok := v.Result.(model.Vector); !ok {
		return nil, errors.New("查询结果不是一个Vector类型")
	} else {
		resultList := &TupleList[L, V]{
			List: make([]Tuple[L, V], 0),
		}
		for i := 0; i < len(vector); i++ {
			sample := vector[i]
			if v.labelCovertFunc == nil {
				//TODO 这里是否可以使用反射解决问题？？？
				return nil, errors.New("必须设置一个label转换函数")
			}
			if v.valueConvertFunc == nil {
				return nil, errors.New("必须设置一个值转换函数")
			}
			resultList.List = append(resultList.List, Tuple[L, V]{
				Label: v.labelCovertFunc(sample.Metric),
				Value: v.valueConvertFunc(sample.Value),
			})
		}
		return resultList, nil
	}
}

func (v *Vector[L, V]) ToPointMap() (map[L]Point[V], error) {
	if vector, ok := v.Result.(model.Vector); !ok {
		return nil, errors.New("查询结果不是一个Vector类型")
	} else {
		resultMap := make(map[L]Point[V])
		for i := 0; i < len(vector); i++ {
			sample := vector[i]
			if v.labelCovertFunc == nil {
				//TODO 这里是否可以使用反射解决问题？？？
				return nil, errors.New("必须设置一个label转换函数")
			}
			if v.valueConvertFunc == nil {
				return nil, errors.New("必须设置一个值转换函数")
			}
			resultMap[v.labelCovertFunc(sample.Metric)] = Point[V]{
				Timestamp: sample.Timestamp.Time(),
				Value:     v.valueConvertFunc(sample.Value),
			}
		}
		return resultMap, nil
	}
}

func (v *Vector[L, V]) ToPointList() (*TupleList[L, Point[V]], error) {
	if vector, ok := v.Result.(model.Vector); !ok {
		return nil, errors.New("查询结果不是一个Vector类型")
	} else {
		resultList := &TupleList[L, Point[V]]{
			List: make([]Tuple[L, Point[V]], 0),
		}
		for i := 0; i < len(vector); i++ {
			sample := vector[i]
			if v.labelCovertFunc == nil {
				//TODO 这里是否可以使用反射解决问题？？？
				return nil, errors.New("必须设置一个label转换函数")
			}
			if v.valueConvertFunc == nil {
				return nil, errors.New("必须设置一个值转换函数")
			}
			resultList.List = append(resultList.List, Tuple[L, Point[V]]{
				Label: v.labelCovertFunc(sample.Metric),
				Value: Point[V]{
					Timestamp: sample.Timestamp.Time(),
					Value:     v.valueConvertFunc(sample.Value),
				},
			})
		}
		return resultList, nil
	}
}
