package main

import (
	"sync"
	"time"
)

type MetricsCollector struct {
	mu            sync.RWMutex
	requestCount  int64
	errorCount    int64
	responseTime  time.Duration
	lastUpdated   time.Time
}

var metrics = &MetricsCollector{}

func (m *MetricsCollector) RecordRequest(duration time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCount++
	if isError {
		m.errorCount++
	}
	m.responseTime += duration
	m.lastUpdated = time.Now()
}

func (m *MetricsCollector) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return map[string]interface{}{
		"requests":     m.requestCount,
		"errors":       m.errorCount,
		"avg_latency":  m.responseTime.Milliseconds() / m.requestCount,
		"last_updated": m.lastUpdated.Format(time.RFC3339),
	}
}



