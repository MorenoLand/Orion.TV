package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/icon.ico
var trayIcon []byte

const liveURL = "https://orion.moreno.land/live.html"

const dragRegionScript = `(function() {
    if (window.__oriontvDragInstalled) return;
    window.__oriontvDragInstalled = true;
    let pressed = false;
    let dragging = false;
    let resizing = false;
    let suppressClick = false;
    let resizeEdge = "";
    let startX = 0;
    let startY = 0;
    function edgeAt(x, y) {
        const width = document.documentElement.clientWidth;
        const height = document.documentElement.clientHeight;
        const left = x < 8;
        const right = x >= width - 8 && x < width;
        const top = y < 8;
        const bottom = y >= height - 8 && y < height;
        if (top && left) return "nw-resize";
        if (top && right) return "ne-resize";
        if (bottom && left) return "sw-resize";
        if (bottom && right) return "se-resize";
        if (left) return "w-resize";
        if (right) return "e-resize";
        if (top) return "n-resize";
        if (bottom) return "s-resize";
        return "";
    }
    document.addEventListener("mousedown", function(event) {
        if (event.button !== 0) return;
        pressed = true;
        dragging = false;
        resizing = false;
        suppressClick = false;
        resizeEdge = edgeAt(event.clientX, event.clientY);
        startX = event.clientX;
        startY = event.clientY;
    }, true);
    document.addEventListener("mousemove", function(event) {
        if (!pressed) return;
        if (resizing || dragging) return;
        if (resizeEdge) {
            if (Math.hypot(event.clientX - startX, event.clientY - startY) < 1) return;
            resizing = true;
            suppressClick = true;
            if (window.chrome && window.chrome.webview) window.chrome.webview.postMessage("wails:resize:" + resizeEdge);
            return;
        }
        if (Math.hypot(event.clientX - startX, event.clientY - startY) < 5) return;
        dragging = true;
        suppressClick = true;
        if (window.chrome && window.chrome.webview) window.chrome.webview.postMessage("wails:drag");
    }, true);
    document.addEventListener("mouseup", function(event) {
        if (event.button === 0) {
            pressed = false;
            dragging = false;
            resizing = false;
            resizeEdge = "";
        }
    }, true);
    document.addEventListener("click", function(event) {
        if (!suppressClick) return;
        event.preventDefault();
        event.stopImmediatePropagation();
        suppressClick = false;
    }, true);
    document.addEventListener("dragstart", function(event) { event.preventDefault(); }, true);
    window.addEventListener("blur", function() {
        pressed = false;
        dragging = false;
        resizing = false;
        resizeEdge = "";
    });
})();`

func main() {
	app := application.New(application.Options{
		Name:        "OrionTV",
		Description: "OrionTV desktop client for the Moreno Land live channel.",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:            "main",
		Title:           "OrionTV",
		Width:           1280,
		Height:          800,
		MinWidth:        640,
		MinHeight:       480,
		DisableResize:   false,
		Frameless:       true,
		InitialPosition: application.WindowCentered,
		URL:             liveURL,
	})
	window.OnWindowEvent(events.Windows.WebViewNavigationCompleted, func(_ *application.WindowEvent) {
		window.HandleMessage("wails:runtime:ready")
		window.ExecJS(dragRegionScript)
	})

	tray := app.SystemTray.New()
	tray.SetIcon(trayIcon)
	tray.SetTooltip("OrionTV")

	menu := app.NewMenu()
	menu.Add("Show").OnClick(func(_ *application.Context) {
		window.Show()
		window.Focus()
	})
	menu.Add("Hide").OnClick(func(_ *application.Context) {
		window.Hide()
	})
	menu.Add("Quit").OnClick(func(_ *application.Context) {
		app.Quit()
	})
	tray.SetMenu(menu)

	tray.OnClick(func() {
		if window.IsVisible() {
			window.Hide()
			return
		}
		window.Show()
		window.Focus()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
