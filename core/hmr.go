package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// HMR manages Hot Module Replacement across runtimes
type HMR struct {
	orchestrator *Orchestrator
	watcher      *fsnotify.Watcher
	watchers     map[string]*RuntimeWatcher
	mu           sync.RWMutex
	enabled      bool
}

// RuntimeWatcher tracks files for a specific runtime
type RuntimeWatcher struct {
	runtime string
	files   map[string]time.Time
	handler ReloadHandler
}

// ReloadHandler is called when a file changes
type ReloadHandler func(ctx context.Context, path string) error

// NewHMR creates an HMR manager
func NewHMR(orchestrator *Orchestrator) (*HMR, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &HMR{
		orchestrator: orchestrator,
		watcher:      watcher,
		watchers:     make(map[string]*RuntimeWatcher),
		enabled:      false,
	}, nil
}

// Enable activates HMR
func (h *HMR) Enable() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.enabled = true
}

// Disable deactivates HMR
func (h *HMR) Disable() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.enabled = false
}

// Watch starts watching files for a runtime
func (h *HMR) Watch(runtime, pattern string, handler ReloadHandler) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	watcher := &RuntimeWatcher{
		runtime: runtime,
		files:   make(map[string]time.Time),
		handler: handler,
	}

	// Find matching files
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	// Add files to watcher
	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			// Watch directory recursively
			if err := h.watchDir(path); err != nil {
				return err
			}
		} else {
			// Watch single file
			if err := h.watcher.Add(path); err != nil {
				return fmt.Errorf("failed to watch %s: %w", path, err)
			}
			watcher.files[path] = info.ModTime()
		}
	}

	h.watchers[runtime] = watcher
	return nil
}

// watchDir adds all files in a directory to the watcher
func (h *HMR) watchDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return h.watcher.Add(path)
		}

		return nil
	})
}

// Start begins monitoring for file changes
func (h *HMR) Start(ctx context.Context) error {
	if !h.enabled {
		return fmt.Errorf("HMR not enabled")
	}

	go h.monitor(ctx)
	return nil
}

// monitor watches for file system events
func (h *HMR) monitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-h.watcher.Events:
			if !ok {
				return
			}

			if h.shouldReload(event) {
				h.handleReload(ctx, event.Name)
			}

		case err, ok := <-h.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("HMR watcher error: %v\n", err)
		}
	}
}

// shouldReload determines if an event should trigger a reload
func (h *HMR) shouldReload(event fsnotify.Event) bool {
	// Only reload on write or create events
	return event.Op&fsnotify.Write == fsnotify.Write ||
		event.Op&fsnotify.Create == fsnotify.Create
}

// handleReload processes a file change
func (h *HMR) handleReload(ctx context.Context, path string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.enabled {
		return
	}

	// Find which runtime owns this file
	for _, watcher := range h.watchers {
		if _, tracked := watcher.files[path]; tracked {
			if err := watcher.handler(ctx, path); err != nil {
				fmt.Printf("HMR reload failed for %s: %v\n", path, err)
			} else {
				fmt.Printf("HMR reloaded: %s\n", path)
			}
			return
		}
	}
}

// Stop shuts down the HMR system
func (h *HMR) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.enabled = false

	if h.watcher != nil {
		return h.watcher.Close()
	}

	return nil
}

// ReloadPython creates a reload handler for Python modules
func (h *HMR) ReloadPython(moduleName string) ReloadHandler {
	return func(ctx context.Context, path string) error {
		code := fmt.Sprintf(`
import importlib
import sys
if '%s' in sys.modules:
    importlib.reload(sys.modules['%s'])
`, moduleName, moduleName)

		_, err := h.orchestrator.Execute(ctx, "python", code)
		return err
	}
}

// ReloadJavaScript creates a reload handler for JavaScript modules
func (h *HMR) ReloadJavaScript(modulePath string) ReloadHandler {
	return func(ctx context.Context, path string) error {
		// Read the updated file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Re-execute the module
		_, err = h.orchestrator.Execute(ctx, "javascript", string(content))
		return err
	}
}

// ReloadNative creates a reload handler for native libraries
func (h *HMR) ReloadNative(runtime, libraryPath string) ReloadHandler {
	return func(ctx context.Context, path string) error {
		// Native libraries require a restart - signal for app reload
		fmt.Printf("Native library changed: %s (restart required)\n", path)
		return nil
	}
}
