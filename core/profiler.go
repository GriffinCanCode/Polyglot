package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Profiler tracks performance across language runtimes
type Profiler struct {
	metrics map[string]*Metrics
	mu      sync.RWMutex
	enabled bool
}

// Metrics holds performance data for a runtime
type Metrics struct {
	RuntimeName     string
	CallCount       int64
	TotalDuration   time.Duration
	AverageDuration time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
	ErrorCount      int64
	LastCalled      time.Time
}

// NewProfiler creates a profiler instance
func NewProfiler() *Profiler {
	return &Profiler{
		metrics: make(map[string]*Metrics),
		enabled: true,
	}
}

// Enable activates profiling
func (p *Profiler) Enable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = true
}

// Disable deactivates profiling
func (p *Profiler) Disable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = false
}

// Track wraps a function call with profiling
func (p *Profiler) Track(runtime, function string, fn func() error) error {
	if !p.enabled {
		return fn()
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	key := fmt.Sprintf("%s.%s", runtime, function)
	p.record(key, runtime, duration, err)

	return err
}

// TrackCall wraps a runtime call with profiling
func (p *Profiler) TrackCall(ctx context.Context, runtime, function string, call func(context.Context) (interface{}, error)) (interface{}, error) {
	if !p.enabled {
		return call(ctx)
	}

	start := time.Now()
	result, err := call(ctx)
	duration := time.Since(start)

	key := fmt.Sprintf("%s.%s", runtime, function)
	p.record(key, runtime, duration, err)

	return result, err
}

// record updates metrics for a function call
func (p *Profiler) record(key, runtime string, duration time.Duration, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	m, exists := p.metrics[key]
	if !exists {
		m = &Metrics{
			RuntimeName: runtime,
			MinDuration: duration,
			MaxDuration: duration,
		}
		p.metrics[key] = m
	}

	m.CallCount++
	m.TotalDuration += duration
	m.AverageDuration = time.Duration(int64(m.TotalDuration) / m.CallCount)
	m.LastCalled = time.Now()

	if duration < m.MinDuration {
		m.MinDuration = duration
	}
	if duration > m.MaxDuration {
		m.MaxDuration = duration
	}

	if err != nil {
		m.ErrorCount++
	}
}

// GetMetrics returns metrics for a specific function
func (p *Profiler) GetMetrics(runtime, function string) (*Metrics, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	key := fmt.Sprintf("%s.%s", runtime, function)
	m, exists := p.metrics[key]
	if !exists {
		return nil, fmt.Errorf("no metrics for %s", key)
	}

	// Return a copy to prevent external modification
	return &Metrics{
		RuntimeName:     m.RuntimeName,
		CallCount:       m.CallCount,
		TotalDuration:   m.TotalDuration,
		AverageDuration: m.AverageDuration,
		MinDuration:     m.MinDuration,
		MaxDuration:     m.MaxDuration,
		ErrorCount:      m.ErrorCount,
		LastCalled:      m.LastCalled,
	}, nil
}

// GetAllMetrics returns all collected metrics
func (p *Profiler) GetAllMetrics() map[string]*Metrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]*Metrics, len(p.metrics))
	for k, v := range p.metrics {
		result[k] = &Metrics{
			RuntimeName:     v.RuntimeName,
			CallCount:       v.CallCount,
			TotalDuration:   v.TotalDuration,
			AverageDuration: v.AverageDuration,
			MinDuration:     v.MinDuration,
			MaxDuration:     v.MaxDuration,
			ErrorCount:      v.ErrorCount,
			LastCalled:      v.LastCalled,
		}
	}

	return result
}

// Reset clears all metrics
func (p *Profiler) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.metrics = make(map[string]*Metrics)
}

// Report generates a performance report
func (p *Profiler) Report() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.metrics) == 0 {
		return "No profiling data collected"
	}

	report := "Performance Report\n"
	report += "==================\n\n"

	for key, m := range p.metrics {
		report += fmt.Sprintf("%s:\n", key)
		report += fmt.Sprintf("  Calls: %d\n", m.CallCount)
		report += fmt.Sprintf("  Total: %v\n", m.TotalDuration)
		report += fmt.Sprintf("  Avg:   %v\n", m.AverageDuration)
		report += fmt.Sprintf("  Min:   %v\n", m.MinDuration)
		report += fmt.Sprintf("  Max:   %v\n", m.MaxDuration)
		report += fmt.Sprintf("  Errors: %d\n", m.ErrorCount)
		report += "\n"
	}

	return report
}
