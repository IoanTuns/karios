package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/diagnoza"
	"pontix/internal/raportare"
	"pontix/internal/versiune"
)

// FolderRapoarteProbleme e folderul din karios/ în care se salvează rapoartele (F16-T15).
const FolderRapoarteProbleme = "rapoarte-probleme"

// LungimeMaximaRaport limitează textul primit din interfață.
const LungimeMaximaRaport = 50000

// RaspunsRaportProblema: textul raportului pregătit, fișierul salvat sau eroarea.
type RaspunsRaportProblema struct {
	Text   string       `json:"text"`
	Titlu  string       `json:"titlu"`
	Cale   string       `json:"cale"`
	Eroare *MesajEroare `json:"eroare,omitempty"`
}

// acumRaport e ceasul rapoartelor; testele îl înlocuiesc.
var acumRaport = time.Now

// PregatesteRaportProblema construiește raportul anonimizat, pe care utilizatorul îl vede și îl poate
// modifica înainte de salvare sau de trimitere.
func (a *App) PregatesteRaportProblema(c raportare.Cerere) RaspunsRaportProblema {
	if strings.TrimSpace(c.Descriere) == "" {
		return RaspunsRaportProblema{Eroare: &MesajEroare{Cod: "raportare.descriere_lipsa"}}
	}
	if c.Tip != raportare.TipParere {
		c.Tip = raportare.TipProblema
	}
	m := diagnoza.ObtineManager()
	x := raportare.Context{Versiune: versiune.Aplicatie, Commit: versiune.Commit, Sistem: runtime.GOOS + "/" + runtime.GOARCH,
		Moment: acumRaport(), NumarErori: m.Statistici().NumarErori}
	if c.CuJurnal {
		x.Jurnal, _ = m.UltimeleLinii(40)
	}
	l := a.limba()
	return RaspunsRaportProblema{Text: raportare.Text(l, c, x), Titlu: raportare.Titlu(l, c)}
}

func (a *App) folderRapoarte() (string, *MesajEroare) {
	if a.radacina == "" {
		return "", &MesajEroare{Cod: "date.eroare"}
	}
	return filepath.Join(a.radacina, FolderRapoarteProbleme), nil
}

// SalveazaRaportProblema scrie raportul în karios/rapoarte-probleme/, pentru trimitere pe email.
func (a *App) SalveazaRaportProblema(text string) RaspunsRaportProblema {
	if strings.TrimSpace(text) == "" || len(text) > LungimeMaximaRaport {
		return RaspunsRaportProblema{Eroare: &MesajEroare{Cod: "raportare.text_invalid"}}
	}
	dir, e := a.folderRapoarte()
	if e != nil {
		return RaspunsRaportProblema{Eroare: e}
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return RaspunsRaportProblema{Eroare: eroareInterfata(err)}
	}
	cale := filepath.Join(dir, "raport-"+acumRaport().Format("20060102-150405")+".md")
	if err := os.WriteFile(cale, []byte(text), 0o600); err != nil {
		return RaspunsRaportProblema{Eroare: eroareInterfata(err)}
	}
	return RaspunsRaportProblema{Text: text, Cale: cale}
}

// TrimiteRaportProblemaPeGitHub deschide în browser pagina de issue nou, completată cu raportul; trimiterea
// o face utilizatorul, cu contul lui GitHub. Aplicația nu are niciun token.
func (a *App) TrimiteRaportProblemaPeGitHub(titlu, text string) RaspunsRaportProblema {
	if strings.TrimSpace(text) == "" || len(text) > LungimeMaximaRaport {
		return RaspunsRaportProblema{Eroare: &MesajEroare{Cod: "raportare.text_invalid"}}
	}
	if a.ctx == nil {
		return RaspunsRaportProblema{Eroare: &MesajEroare{Cod: "date.eroare"}}
	}
	wruntime.BrowserOpenURL(a.ctx, raportare.AdresaIssue(a.limba(), titlu, text))
	return RaspunsRaportProblema{Text: text, Titlu: titlu}
}

// DeschideFolderRapoarteProbleme deschide în Explorer folderul cu rapoartele salvate.
func (a *App) DeschideFolderRapoarteProbleme() string {
	dir, e := a.folderRapoarte()
	if e != nil {
		return ""
	}
	_ = os.MkdirAll(dir, 0o700)
	_ = exec.Command("explorer.exe", filepath.FromSlash(dir)).Start()
	return dir
}
