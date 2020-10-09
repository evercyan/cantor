package main

import (
	"github.com/evercyan/cantor/backend"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {
	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")
	app := wails.CreateApp(&wails.AppConfig{
		Width:            800,
		Height:           600,
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
