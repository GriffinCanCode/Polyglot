package webview

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

// Webview manages the native webview window
type Webview struct {
	config   core.WebviewConfig
	bridge   core.Bridge
	instance WebviewBackend
	mu       sync.Mutex
	running  bool
}

// New creates a new webview instance
func New(config core.WebviewConfig, bridge core.Bridge) *Webview {
	return &Webview{
		config: config,
		bridge: bridge,
	}
}

// Initialize creates the webview window
func (w *Webview) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return fmt.Errorf("webview already running")
	}

	// Create webview instance using configured backend
	w.instance = NewBackend(w.config.Debug)
	if w.instance == nil {
		return fmt.Errorf("failed to create webview")
	}

	// Configure window
	w.instance.SetTitle(w.config.Title)
	w.instance.SetSize(w.config.Width, w.config.Height, HintNone)

	// Bind bridge functions
	w.bindBridge()

	return nil
}

// Run starts the webview event loop
func (w *Webview) Run() error {
	w.mu.Lock()
	if w.instance == nil {
		w.mu.Unlock()
		return fmt.Errorf("webview not initialized")
	}
	w.running = true
	w.mu.Unlock()

	// Navigate to URL
	w.instance.Navigate(w.config.URL)

	// Run event loop (blocks)
	w.instance.Run()

	w.mu.Lock()
	w.running = false
	w.mu.Unlock()

	return nil
}

// Eval executes JavaScript in the webview
func (w *Webview) Eval(script string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.instance == nil {
		return fmt.Errorf("webview not initialized")
	}

	w.instance.Eval(script)
	return nil
}

// Bind adds a Go function callable from JavaScript
func (w *Webview) Bind(name string, fn interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.instance == nil {
		return fmt.Errorf("webview not initialized")
	}

	return w.instance.Bind(name, fn)
}

// Terminate closes the webview
func (w *Webview) Terminate() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.instance == nil {
		return nil
	}

	w.instance.Terminate()
	w.instance.Destroy()
	w.instance = nil

	return nil
}

// bindBridge sets up the JavaScript bridge
func (w *Webview) bindBridge() {
	if w.bridge == nil {
		return
	}

	// Create a unified bridge function
	w.instance.Bind("__polyglot_call__", func(name string, argsJSON string) (string, error) {
		// Parse arguments
		var args []interface{}
		if argsJSON != "" {
			if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}
		}

		// Call bridge function
		result, err := w.bridge.Call(nil, name, args...)
		if err != nil {
			return "", err
		}

		// Serialize result
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("failed to serialize result: %w", err)
		}

		return string(resultJSON), nil
	})

	// Inject bridge initialization script
	initScript := `
		window.polyglot = {
			call: async function(name, ...args) {
				const argsJSON = JSON.stringify(args);
				const resultJSON = await __polyglot_call__(name, argsJSON);
				return JSON.parse(resultJSON);
			}
		};
	`
	w.instance.Init(initScript)
}
