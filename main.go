package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"pontix/internal/i18n"
	"pontix/internal/portabil"
)

//go:embed all:frontend/dist
var resurse embed.FS

func main() {
	foldere, err := portabil.Pregateste()
	if err != nil {
		// Mesajul pentru utilizator și verificarea dreptului de scriere vin cu F05-T05.
		log.Fatal(err)
	}

	app := NewApp()

	err = wails.Run(&options.App{
		Title:     i18n.T(i18n.Implicita, "aplicatie.titlu", nil),
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: resurse,
		},
		Windows: &windows.Options{
			// Fără această opțiune, WebView2 își ține datele în %APPDATA% (D002).
			WebviewUserDataPath: foldere.WebView,
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
