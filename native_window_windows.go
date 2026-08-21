//go:build windows

package main

import (
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

var nativeResizeHitTests = map[string]uintptr{
	"n-resize":  w32.HTTOP,
	"ne-resize": w32.HTTOPRIGHT,
	"e-resize":  w32.HTRIGHT,
	"se-resize": w32.HTBOTTOMRIGHT,
	"s-resize":  w32.HTBOTTOM,
	"sw-resize": w32.HTBOTTOMLEFT,
	"w-resize":  w32.HTLEFT,
	"nw-resize": w32.HTTOPLEFT,
}

func handleNativeWindowMessage(window application.Window, message string) bool {
	var hitTest uintptr
	switch {
	case strings.HasPrefix(message, "oriontv:resize:"):
		var ok bool
		hitTest, ok = nativeResizeHitTests[strings.TrimPrefix(message, "oriontv:resize:")]
		if !ok {
			return false
		}
	default:
		return false
	}
	hwnd := window.NativeWindow()
	if hwnd == nil {
		return false
	}
	w32.ReleaseCapture()
	return w32.PostMessage(w32.HWND(uintptr(hwnd)), w32.WM_NCLBUTTONDOWN, hitTest, 0)
}
