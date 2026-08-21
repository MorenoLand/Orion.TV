# OrionTV

OrionTV is a Wails3 desktop client for the Moreno Land live channel. It loads `https://orion.moreno.land/live.html` in a native WebView2 window and keeps the original tray controls.

## Features

- 1280×800 centered window with 640×480 minimum dimensions and resizing enabled.
- OrionTV application and tray icon using the supplied TV mark.
- Tray menu actions for Show, Hide, and Quit.
- Left-click tray toggle with focus restoration when showing the window.
- Window hide/show keeps the native window position and size because it uses Wails window visibility operations without tray-window attachment.

## Build

Requirements:

- Go 1.25 or newer.
- Wails3 beta.11 or a compatible Wails3 release.
- Node.js and npm for the local fallback page build.
- WebView2 on Windows.

Install the frontend dependencies and build the Windows executable:

```powershell
npm --prefix frontend install
wails3 task windows:build
```

The executable is written to `bin/oriontv.exe`. The Wails build configuration embeds the product metadata and icon from `build/`.

For development, use `wails3 dev` or `wails3 task dev`. The application itself has not been launched as part of the source and package verification in this repository.

## Packaging

The standard Wails3 Windows package task can create an installer when NSIS is installed:

```powershell
wails3 package GOOS=windows
```

UPX compression is optional and should be applied only to a verified Windows executable:

```powershell
upx --best bin/oriontv.exe
```

Windows amd64 is the validated build target for this checkout. The Wails task layout also contains desktop task definitions for macOS and Linux; those targets require their platform toolchains and were not built on Windows.

## License

This repository uses the all-rights-reserved notice in `LICENSE`. No additional redistribution rights are granted by the source tree.
