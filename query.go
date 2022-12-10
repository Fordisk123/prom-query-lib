package prom_query_tool

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

type Query struct {
	api    v1.API
	premQl string
	time   time.Time
	start  time.Time
	end    time.Time
	step   time.Duration
	opts   []v1.Option
}

func (d *Query) PromQl(template string, args ...any) *Query {
	d.premQl = fmt.Sprintf(template, args...)
	return d
}

func (d *Query) Time(time time.Time) *Query {
	d.time = time
	return d
}

func (d *Query) Options(opts ...v1.Option) *Query {
	d.opts = opts
	return d
}

func (d *Query) StartTime(time time.Time) *Query {
	d.start = time
	return d
}

func (d *Query) EndTime(time time.Time) *Query {
	d.end = time
	return d
}

func (d *Query) Range(start, end time.Time, step time.Duration) *Query {
	return d.StartTime(start).EndTime(end).Step(step)
}

func (d *Query) Step(step time.Duration) *Query {
	if step == 0 {
		d.step = 1 * time.Minute
	} else {
		d.step = step
	}
	return d
}

func (d *Query) DoQuery(ctx context.Context) (model.Value, v1.Warnings, error) {
	if len(d.premQl) == 0 {
		return nil, nil, errors.New("请传入PromQL语句")
	}
	return d.api.Query(ctx, d.premQl, d.time, d.opts...)
}

func (d *Query) DoQueryRange(ctx context.Context) (model.Value, v1.Warnings, error) {
	if len(d.premQl) == 0 {
		return nil, nil, errors.New("请传入PromQL语句")
	}
	return d.api.QueryRange(ctx, d.premQl, v1.Range{
		Start: d.start,
		End:   d.end,
		Step:  d.step,
	}, d.opts...)
}
