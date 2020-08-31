package main

import (
	"emoji/backend"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {
	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")
	app := wails.CreateApp(&wails.AppConfig{
		Width:            1024,
		Height:           768,
		Resizable:        true,
		Title:            "Cantor",
		JS:               js,
		CSS:              css,
		Colour:           "#18181f",
		DisableInspector: true,
	})
	app.Bind(&backend.App{})
	app.Run()
}
