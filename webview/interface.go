package webview

// Hint represents window sizing hints
type Hint int

// Hint constants for window sizing
const (
	HintNone  Hint = 0
	HintMin   Hint = 1
	HintMax   Hint = 2
	HintFixed Hint = 3
)

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

// NewBackend creates a webview instance (implementation set by build tags)
var NewBackend func(debug bool) WebviewBackend

// ConfigureBackend allows setting a custom webview backend
func ConfigureBackend(factory func(debug bool) WebviewBackend) {
	NewBackend = factory
}
