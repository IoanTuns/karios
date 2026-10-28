package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"pontix/internal/diagnoza"
)

// StatisticiDiagnoza întoarce metricile curente de sănătate, numărul de erori și timpii de interogare.
func (a *App) StatisticiDiagnoza() diagnoza.StatisticiDiagnoza {
	return diagnoza.ObtineManager().Statistici()
}

// InregistreazaEroareFrontend permite interfeței React/TypeScript să trimită erori neprinse în jurnalul local.
func (a *App) InregistreazaEroareFrontend(sursa, mesaj, stack string) {
	diagnoza.ObtineManager().InregistreazaEroareFrontend(sursa, mesaj, stack)
}

// CitesteUltimeleLiniiJurnal întoarce ultimele linii din fișierul de log pentru copiere rapidă.
func (a *App) CitesteUltimeleLiniiJurnal(numar int) string {
	if numar <= 0 {
		numar = 50
	}
	return diagnoza.ObtineManager().CitesteUltimeleLinii(numar)
}

// DeschideFolderJurnale deschide folderul local cu fișiere de log în Windows Explorer.
func (a *App) DeschideFolderJurnale() string {
	m := diagnoza.ObtineManager()
	cale := m.Statistici().CaleFolderJurnale
	if cale != "" {
		_ = os.MkdirAll(cale, 0o700)
		_ = exec.Command("explorer.exe", filepath.FromSlash(cale)).Start()
	}
	return cale
}

