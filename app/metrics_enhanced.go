package main

import (
	"sync"
	"time"
)

type EnhancedMetrics struct {
	mu              sync.RWMutex
	requestCount    int64
	errorCount      int64
	successCount    int64
	totalLatency    time.Duration
	lastRequestTime time.Time
	endpoints       map[string]*EndpointMetrics
}

type EndpointMetrics struct {
	Count   int64
	Errors  int64
	Latency time.Duration
}

var enhancedMetrics = &EnhancedMetrics{
	endpoints: make(map[string]*EndpointMetrics),
}

func (em *EnhancedMetrics) RecordRequest(endpoint string, duration time.Duration, isError bool) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.requestCount++
	em.totalLatency += duration
	em.lastRequestTime = time.Now()

	if isError {
		em.errorCount++
	} else {
		em.successCount++
	}

	if em.endpoints[endpoint] == nil {
		em.endpoints[endpoint] = &EndpointMetrics{}
	}
	em.endpoints[endpoint].Count++
	em.endpoints[endpoint].Latency += duration
	if isError {
		em.endpoints[endpoint].Errors++
	}
}

func (em *EnhancedMetrics) GetEnhancedStats() map[string]interface{} {
	em.mu.RLock()
	defer em.mu.RUnlock()

	avgLatency := int64(0)
	if em.requestCount > 0 {
		avgLatency = em.totalLatency.Milliseconds() / em.requestCount
	}

	return map[string]interface{}{
		"total_requests":    em.requestCount,
		"success_count":     em.successCount,
		"error_count":       em.errorCount,
		"error_rate":        float64(em.errorCount) / float64(em.requestCount) * 100,
		"avg_latency_ms":    avgLatency,
		"last_request_time": em.lastRequestTime.Format(time.RFC3339),
		"endpoints":         em.endpoints,
	}
}


