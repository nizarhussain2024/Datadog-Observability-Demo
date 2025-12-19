package main

import (
	"sync"
	"time"
)

type ServiceDependency struct {
	FromService string
	ToService   string
	CallCount   int64
	LastCall    time.Time
	AvgLatency  time.Duration
}

type DependencyTracker struct {
	mu          sync.RWMutex
	dependencies map[string]*ServiceDependency
}

var dependencyTracker = &DependencyTracker{
	dependencies: make(map[string]*ServiceDependency),
}

func (dt *DependencyTracker) RecordCall(from, to string, latency time.Duration) {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	key := from + "->" + to
	dep, exists := dt.dependencies[key]
	if !exists {
		dep = &ServiceDependency{
			FromService: from,
			ToService:   to,
		}
		dt.dependencies[key] = dep
	}

	dep.CallCount++
	dep.LastCall = time.Now()
	
	// Simple moving average
	if dep.AvgLatency == 0 {
		dep.AvgLatency = latency
	} else {
		dep.AvgLatency = (dep.AvgLatency + latency) / 2
	}
}

func (dt *DependencyTracker) GetDependencies() map[string]*ServiceDependency {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	result := make(map[string]*ServiceDependency)
	for k, v := range dt.dependencies {
		result[k] = v
	}
	return result
}


