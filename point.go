package prom_query_tool

import "time"

type Point[V comparable] struct {
	Timestamp time.Time
	Value     V
}
