package webview

import "github.com/webview/webview"

// Hint type alias
type Hint = webview.Hint

// WebviewBackend defines the interface for webview implementations
type WebviewBackend interface {
	// SetTitle sets the window title
	SetTitle(title string)

	// SetSize sets the window dimensions
	SetSize(width, height int, hint Hint)

	// Navigate loads a URL
	Navigate(url string)

	// Run starts the event loop (blocking)
	Run()

	// Eval executes JavaScript
	Eval(script string)

	// Bind adds a Go function callable from JavaScript
	Bind(name string, fn interface{}) error

	// Init runs initialization JavaScript
	Init(script string)

	// Terminate stops the webview
	Terminate()

	// Destroy cleans up resources
	Destroy()
}

// Hint constants for window sizing
const (
	HintNone  = webview.HintNone
	HintMin   = webview.HintMin
	HintMax   = webview.HintMax
	HintFixed = webview.HintFixed
)

// NewBackend creates a native webview instance
var NewBackend func(debug bool) WebviewBackend = func(debug bool) WebviewBackend {
	return webview.New(debug)
}

// stubBackend is a no-op implementation for testing
type stubBackend struct{}

func (s *stubBackend) SetTitle(title string)                  {}
func (s *stubBackend) SetSize(w, h int, hint Hint)            {}
func (s *stubBackend) Navigate(url string)                    {}
func (s *stubBackend) Run()                                   {}
func (s *stubBackend) Eval(script string)                     {}
func (s *stubBackend) Bind(name string, fn interface{}) error { return nil }
func (s *stubBackend) Init(script string)                     {}
func (s *stubBackend) Terminate()                             {}
func (s *stubBackend) Destroy()                               {}

// ConfigureBackend allows setting a custom webview backend
func ConfigureBackend(factory func(debug bool) WebviewBackend) {
	NewBackend = factory
}
