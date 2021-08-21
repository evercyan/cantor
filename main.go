package main

import (
	_ "embed"

	"github.com/wailsapp/wails"

	"github.com/evercyan/cantor/backend"
)

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

func main() {
	app := wails.CreateApp(&wails.AppConfig{
		Width:            888,
		Height:           666,
		Resizable:        true,
		Title:            "Cantor",
		JS:               js,
		CSS:              css,
		Colour:           "#18181f",
		DisableInspector: true,
	})
	app.Bind(new(backend.App))
	app.Run()
}
