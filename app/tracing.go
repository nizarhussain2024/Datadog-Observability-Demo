package main

import (
	"context"
	"time"
)

type TraceSpan struct {
	TraceID   string
	SpanID    string
	ParentID  string
	Operation string
	StartTime time.Time
	Duration  time.Duration
	Tags      map[string]string
}

func startSpan(ctx context.Context, operation string) *TraceSpan {
	return &TraceSpan{
		TraceID:   generateTraceID(),
		SpanID:    generateSpanID(),
		Operation: operation,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
	}
}

func (s *TraceSpan) Finish() {
	s.Duration = time.Since(s.StartTime)
	// In production, send to Datadog APM
}

func (s *TraceSpan) SetTag(key, value string) {
	s.Tags[key] = value
}

func generateTraceID() string {
	return "trace-" + time.Now().Format("20060102150405")
}

func generateSpanID() string {
	return "span-" + time.Now().Format("20060102150405")
}




