package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"pontix/internal/audit"
	"pontix/internal/diagnoza"
	"pontix/internal/i18n"
	"pontix/internal/portabil"
	"pontix/internal/stocare"
)

//go:embed all:frontend/dist
var resurse embed.FS

func main() {
	foldere, err := portabil.Pregateste()
	if err != nil {
		oprireCuEroare(err)
	}
	diagnoza.Initializeaza(foldere.Radacina)

	db, paza, err := stocare.Deschide(context.Background(), foldere.Radacina)
	if err != nil {
		diagnoza.ObtineManager().InregistreazaEroare("pornire/stocare", err)
		oprireCuEroare(err)
	}
	if r := paza.Rezultat(); r.Stare == audit.StareIncident {
		// Doar numele tabelelor și id-ul rândului din jurnal, fără date (D043).
		diagnoza.ObtineManager().InregistreazaEroare("pornire/audit", fmt.Errorf("incident: tabele %v, rând rupt %d", r.Tabele, r.RandRupt))
	}

	app := NewApp(db)
	app.radacina = foldere.Radacina
	app.paza = paza
	// Până la autentificare (sau crearea primului cont), scrierile sunt refuzate (D045).
	paza.CereAutentificare(true)

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
		OnStartup:     app.pornire,
		OnShutdown:    app.oprire,
		OnBeforeClose: app.inainteDeInchidere,
		Bind: []any{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// oprireCuEroare arată eroarea într-o fereastră Windows și închide aplicația.
// Setările nu se pot citi încă, deci mesajul e în limba Windows (D008).
func oprireCuEroare(err error) {
	limba := i18n.LimbaSistem()
	text := i18n.T(limba, "aplicatie.eroare_pornire", map[string]string{"detalii": err.Error()})
	var e *i18n.Eroare
	if errors.As(err, &e) {
		text = e.Text(limba)
	}
	arataEroare(i18n.T(limba, "aplicatie.titlu", nil), text)
	os.Exit(1)
}
