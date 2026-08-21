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
    const id = "oriontv-drag-region";
    let style = document.getElementById(id);
    if (!style) {
        style = document.createElement("style");
        style.id = id;
        (document.head || document.documentElement).appendChild(style);
    }
    style.textContent = "html, body, body * { --wails-draggable: drag !important; }";
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
	window.OnWindowEvent(events.Common.WindowRuntimeReady, func(_ *application.WindowEvent) {
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
