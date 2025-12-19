package main

import (
	"sync"
	"time"
)

type MetricAggregator struct {
	mu           sync.RWMutex
	aggregations map[string]*Aggregation
	windowSize   time.Duration
}

type Aggregation struct {
	Count   int64
	Sum     float64
	Min     float64
	Max     float64
	Avg     float64
	LastUpdate time.Time
}

var aggregator = &MetricAggregator{
	aggregations: make(map[string]*Aggregation),
	windowSize:   5 * time.Minute,
}

func (ma *MetricAggregator) Record(metricName string, value float64) {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	agg, exists := ma.aggregations[metricName]
	if !exists {
		agg = &Aggregation{
			Min: value,
			Max: value,
		}
		ma.aggregations[metricName] = agg
	}

	agg.Count++
	agg.Sum += value
	if value < agg.Min {
		agg.Min = value
	}
	if value > agg.Max {
		agg.Max = value
	}
	agg.Avg = agg.Sum / float64(agg.Count)
	agg.LastUpdate = time.Now()
}

func (ma *MetricAggregator) GetAggregation(metricName string) *Aggregation {
	ma.mu.RLock()
	defer ma.mu.RUnlock()
	return ma.aggregations[metricName]
}

func (ma *MetricAggregator) GetAllAggregations() map[string]*Aggregation {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	result := make(map[string]*Aggregation)
	for k, v := range ma.aggregations {
		result[k] = v
	}
	return result
}

