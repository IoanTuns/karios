package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"pontix/internal/i18n"
)

//go:embed all:frontend/dist
var resurse embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     i18n.T(i18n.Implicita, "aplicatie.titlu", nil),
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: resurse,
		},
		OnStartup: app.pornire,
		Bind: []any{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
