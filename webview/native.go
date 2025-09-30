//go:build webview_enabled
// +build webview_enabled

package webview

// NOTE: To use native webview, install the webview library:
// go get github.com/webview/webview
//
// Then uncomment the import and implementation below.
// For now, this file serves as a template for future native webview integration.

/*
import "github.com/webview/webview"

// Native webview backend using external library
type nativeBackend struct {
	wv webview.WebView
}

// NewNativeBackend creates a native webview instance
func NewNativeBackend(debug bool) WebviewBackend {
	return &nativeBackend{
		wv: webview.New(debug),
	}
}

func (n *nativeBackend) SetTitle(title string) {
	n.wv.SetTitle(title)
}

func (n *nativeBackend) SetSize(w, h int, hint Hint) {
	n.wv.SetSize(w, h, webview.Hint(hint))
}

func (n *nativeBackend) Navigate(url string) {
	n.wv.Navigate(url)
}

func (n *nativeBackend) Run() {
	n.wv.Run()
}

func (n *nativeBackend) Eval(script string) {
	n.wv.Eval(script)
}

func (n *nativeBackend) Bind(name string, fn interface{}) error {
	return n.wv.Bind(name, fn)
}

func (n *nativeBackend) Init(script string) {
	n.wv.Init(script)
}

func (n *nativeBackend) Terminate() {
	n.wv.Terminate()
}

func (n *nativeBackend) Destroy() {
	n.wv.Destroy()
}

func init() {
	NewBackend = NewNativeBackend
}
*/

// Placeholder implementation
func init() {
	NewBackend = NewStubBackend
}
