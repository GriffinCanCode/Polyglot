package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/polyglot-framework/polyglot/core"
)

// TestProfilerBasic tests basic profiler functionality
func TestProfilerBasic(t *testing.T) {
	profiler := core.NewProfiler()

	// Test enable/disable
	profiler.Enable()
	profiler.Disable()
	profiler.Enable()

	// Test tracking
	err := profiler.Track("python", "test_function", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	if err != nil {
		t.Errorf("track failed: %v", err)
	}

	// Get metrics
	metrics, err := profiler.GetMetrics("python", "test_function")
	if err != nil {
		t.Errorf("failed to get metrics: %v", err)
	}

	if metrics.CallCount != 1 {
		t.Errorf("expected 1 call, got %d", metrics.CallCount)
	}

	if metrics.TotalDuration < 10*time.Millisecond {
		t.Errorf("duration too short: %v", metrics.TotalDuration)
	}
}

// TestProfilerMultipleCalls tests tracking multiple calls
func TestProfilerMultipleCalls(t *testing.T) {
	profiler := core.NewProfiler()

	// Make multiple calls
	for i := 0; i < 5; i++ {
		profiler.Track("javascript", "compute", func() error {
			time.Sleep(5 * time.Millisecond)
			return nil
		})
	}

	metrics, err := profiler.GetMetrics("javascript", "compute")
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}

	if metrics.CallCount != 5 {
		t.Errorf("expected 5 calls, got %d", metrics.CallCount)
	}

	if metrics.AverageDuration == 0 {
		t.Error("average duration should not be zero")
	}

	if metrics.MinDuration > metrics.MaxDuration {
		t.Error("min duration should be <= max duration")
	}
}

// TestProfilerErrors tests error tracking
func TestProfilerErrors(t *testing.T) {
	profiler := core.NewProfiler()

	testErr := errors.New("test error")

	// Track successful and failed calls
	profiler.Track("rust", "fallible_fn", func() error {
		return nil
	})

	profiler.Track("rust", "fallible_fn", func() error {
		return testErr
	})

	profiler.Track("rust", "fallible_fn", func() error {
		return testErr
	})

	metrics, err := profiler.GetMetrics("rust", "fallible_fn")
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}

	if metrics.CallCount != 3 {
		t.Errorf("expected 3 calls, got %d", metrics.CallCount)
	}

	if metrics.ErrorCount != 2 {
		t.Errorf("expected 2 errors, got %d", metrics.ErrorCount)
	}
}

// TestProfilerTrackCall tests context-aware tracking
func TestProfilerTrackCall(t *testing.T) {
	profiler := core.NewProfiler()
	ctx := context.Background()

	result, err := profiler.TrackCall(ctx, "python", "add", func(ctx context.Context) (interface{}, error) {
		return 42, nil
	})

	if err != nil {
		t.Errorf("track call failed: %v", err)
	}

	if result.(int) != 42 {
		t.Errorf("expected result 42, got %v", result)
	}

	metrics, err := profiler.GetMetrics("python", "add")
	if err != nil {
		t.Errorf("failed to get metrics: %v", err)
	}

	if metrics.CallCount != 1 {
		t.Errorf("expected 1 call, got %d", metrics.CallCount)
	}
}

// TestProfilerGetAllMetrics tests retrieving all metrics
func TestProfilerGetAllMetrics(t *testing.T) {
	profiler := core.NewProfiler()

	// Track calls to different runtimes
	profiler.Track("python", "fn1", func() error { return nil })
	profiler.Track("javascript", "fn2", func() error { return nil })
	profiler.Track("rust", "fn3", func() error { return nil })

	allMetrics := profiler.GetAllMetrics()

	if len(allMetrics) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(allMetrics))
	}
}

// TestProfilerReset tests resetting metrics
func TestProfilerReset(t *testing.T) {
	profiler := core.NewProfiler()

	// Track some calls
	profiler.Track("python", "test", func() error { return nil })
	profiler.Track("python", "test", func() error { return nil })

	metrics, _ := profiler.GetMetrics("python", "test")
	if metrics.CallCount != 2 {
		t.Errorf("expected 2 calls before reset, got %d", metrics.CallCount)
	}

	// Reset
	profiler.Reset()

	// Metrics should be gone
	_, err := profiler.GetMetrics("python", "test")
	if err == nil {
		t.Error("expected error after reset, got nil")
	}
}

// TestProfilerReport tests report generation
func TestProfilerReport(t *testing.T) {
	profiler := core.NewProfiler()

	// Track some calls
	profiler.Track("python", "process", func() error {
		time.Sleep(5 * time.Millisecond)
		return nil
	})

	report := profiler.Report()
	if report == "" {
		t.Error("report should not be empty")
	}

	if len(report) < 50 {
		t.Errorf("report seems too short: %d characters", len(report))
	}
}

// TestProfilerDisabled tests that disabled profiler doesn't track
func TestProfilerDisabled(t *testing.T) {
	profiler := core.NewProfiler()
	profiler.Disable()

	// Track a call while disabled
	profiler.Track("python", "test", func() error {
		return nil
	})

	// Should have no metrics
	_, err := profiler.GetMetrics("python", "test")
	if err == nil {
		t.Error("expected error for non-existent metrics")
	}
}
