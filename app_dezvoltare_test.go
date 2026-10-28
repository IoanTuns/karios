package main

import (
	"testing"
)

func TestAppDetaliiDezvoltare(t *testing.T) {
	app := &App{radacina: t.TempDir()}
	detalii := app.DetaliiDezvoltare()

	if detalii.NumeRepo != "karios" {
		t.Errorf("așteptat repo karios, obținut: %s", detalii.NumeRepo)
	}
	if detalii.NumeAplicatie != "Pontix" {
		t.Errorf("așteptat Pontix, obținut: %s", detalii.NumeAplicatie)
	}
	if detalii.Versiune == "" {
		t.Errorf("versiunea este goală")
	}
	if detalii.CerinteMinime.SistemOperare == "" {
		t.Errorf("cerințele minime sunt goale")
	}
	if detalii.DetaliiLLM.FereastraContext != 8192 {
		t.Errorf("așteptat context 8192, obținut: %d", detalii.DetaliiLLM.FereastraContext)
	}
	if detalii.Licenta != "Apache License 2.0 (Apache-2.0)" {
		t.Errorf("așteptat Apache License 2.0 (Apache-2.0), obținut: %s", detalii.Licenta)
	}
}
