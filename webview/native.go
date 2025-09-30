//go:build !stub
// +build !stub

package webview

import webview "github.com/webview/webview_go"

// NativeBackend implements WebviewBackend using the webview/webview library
type NativeBackend struct {
	wv webview.WebView
}

// NewNativeBackend creates a native webview instance
func NewNativeBackend(debug bool) WebviewBackend {
	return &NativeBackend{
		wv: webview.New(debug),
	}
}

func (n *NativeBackend) SetTitle(title string) {
	n.wv.SetTitle(title)
}

func (n *NativeBackend) SetSize(w, h int, hint Hint) {
	n.wv.SetSize(w, h, webview.Hint(hint))
}

func (n *NativeBackend) Navigate(url string) {
	n.wv.Navigate(url)
}

func (n *NativeBackend) Run() {
	n.wv.Run()
}

func (n *NativeBackend) Eval(script string) {
	n.wv.Eval(script)
}

func (n *NativeBackend) Bind(name string, fn interface{}) error {
	return n.wv.Bind(name, fn)
}

func (n *NativeBackend) Init(script string) {
	n.wv.Init(script)
}

func (n *NativeBackend) Terminate() {
	n.wv.Terminate()
}

func (n *NativeBackend) Destroy() {
	n.wv.Destroy()
}

func init() {
	NewBackend = NewNativeBackend
}
