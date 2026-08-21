package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/icon.ico
var trayIcon []byte

const liveURL = "https://orion.moreno.land/live.html"

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
		InitialPosition: application.WindowCentered,
		URL:             liveURL,
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
