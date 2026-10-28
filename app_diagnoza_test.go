package main

import (
	"strings"
	"testing"

	"pontix/internal/diagnoza"
)

func TestAppDiagnoza(t *testing.T) {
	tempDir := t.TempDir()
	diagnoza.Initializeaza(tempDir)

	app := &App{radacina: tempDir}

	app.InregistreazaEroareFrontend("UnitTester", "Eroare simulat frontend", "Stack dummy")

	stat := app.StatisticiDiagnoza()
	if stat.NumarErori < 1 {
		t.Errorf("așteptat cel puțin 1 eroare, obținut: %d", stat.NumarErori)
	}

	linii := app.CitesteUltimeleLiniiJurnal(10)
	if !strings.Contains(linii, "UnitTester") {
		t.Errorf("jurnalul ar trebui să conțină UnitTester, obținut: %s", linii)
	}
}

