//go:build !webview_enabled
// +build !webview_enabled

package webview

import "fmt"

// StubBackend is a no-op implementation for testing or when webview is disabled
type StubBackend struct {
	title  string
	url    string
	width  int
	height int
}

// NewStubBackend creates a stub webview instance
func NewStubBackend(debug bool) WebviewBackend {
	fmt.Println("Using stub webview backend (webview not enabled in build)")
	return &StubBackend{
		width:  800,
		height: 600,
	}
}

func (s *StubBackend) SetTitle(title string) {
	s.title = title
	fmt.Printf("Stub: SetTitle(%s)\n", title)
}

func (s *StubBackend) SetSize(w, h int, hint Hint) {
	s.width = w
	s.height = h
	fmt.Printf("Stub: SetSize(%d, %d, %d)\n", w, h, hint)
}

func (s *StubBackend) Navigate(url string) {
	s.url = url
	fmt.Printf("Stub: Navigate(%s)\n", url)
}

func (s *StubBackend) Run() {
	fmt.Println("Stub: Run() - webview would start here")
	fmt.Printf("Stub: Window: %s - %dx%d - %s\n", s.title, s.width, s.height, s.url)
}

func (s *StubBackend) Eval(script string) {
	fmt.Printf("Stub: Eval(%s)\n", script)
}

func (s *StubBackend) Bind(name string, fn interface{}) error {
	fmt.Printf("Stub: Bind(%s, <func>)\n", name)
	return nil
}

func (s *StubBackend) Init(script string) {
	fmt.Printf("Stub: Init(%s)\n", script)
}

func (s *StubBackend) Terminate() {
	fmt.Println("Stub: Terminate()")
}

func (s *StubBackend) Destroy() {
	fmt.Println("Stub: Destroy()")
}

func init() {
	NewBackend = NewStubBackend
}
