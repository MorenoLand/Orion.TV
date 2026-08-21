//go:build !windows

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func handleNativeWindowMessage(application.Window, string) bool {
	return false
}
