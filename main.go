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
    let suppressClick = false;
    let startX = 0;
    let startY = 0;
    document.addEventListener("mousedown", function(event) {
        if (event.button !== 0) return;
        pressed = true;
        dragging = false;
        suppressClick = false;
        startX = event.clientX;
        startY = event.clientY;
    }, true);
    document.addEventListener("mousemove", function(event) {
        if (!pressed || dragging) return;
        if (Math.hypot(event.clientX - startX, event.clientY - startY) < 5) return;
        dragging = true;
        suppressClick = true;
        if (window.chrome && window.chrome.webview) window.chrome.webview.postMessage("wails:drag");
    }, true);
    document.addEventListener("mouseup", function(event) {
        if (event.button === 0) {
            pressed = false;
            dragging = false;
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
