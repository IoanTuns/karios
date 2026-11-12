package main

import (
	"testing"

	"pontix/internal/exportsalarizare"
)

// F10-T03 (D108): modelele generice exista, un model nou se salveaza doar de administrator,
// iar exportul ia un rand pe marca.
func TestExportSalarizareModeleSiRanduri(t *testing.T) {
	a := appCuConturi(t)
	r := a.ListaModeleExportSalarizare()
	if r.Eroare != nil || len(r.Modele) != 2 || len(r.Surse) < 20 {
		t.Fatalf("modele: %+v", r)
	}
	nou := exportsalarizare.Model{Denumire: "Program X", Format: exportsalarizare.FormatCSV, Antet: true, Coloane: []exportsalarizare.Coloana{{Sursa: "marca", Titlu: "Marca"}, {Sursa: "total:OL", Titlu: "Ore"}}, Activ: true}
	if s := a.SalveazaModelExportSalarizare(nou); s.Eroare == nil || s.Eroare.Cod != "cont.doar_administrator" {
		t.Fatalf("fara cont: %+v", s.Eroare)
	}
	if c := a.CreeazaPrimulAdministrator(cerereCont("ana")); c.Eroare != nil {
		t.Fatal(c.Eroare)
	}
	s := a.SalveazaModelExportSalarizare(nou)
	if s.Eroare != nil || s.Model == nil || s.Model.ID == "" || s.Model.Revizie != "1" || len(s.Modele) != 3 {
		t.Fatalf("salvare: %+v", s)
	}
	m := *s.Model
	m.Revizie = "7"
	if c := a.SalveazaModelExportSalarizare(m); c.Eroare == nil || c.Eroare.Cod != "export_salarizare.conflict" {
		t.Fatalf("conflict: %+v", c.Eroare)
	}
	if r := a.IncarcaDemo(); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	e := a.ExportaSalarizare("export-csv-generic", 2026, 10, "")
	if e.Eroare != nil || e.Randuri == 0 {
		t.Fatalf("export: %+v", e)
	}
	if x := a.ExportaSalarizare("inexistent", 2026, 10, ""); x.Eroare == nil || x.Eroare.Cod != "export_salarizare.model_inexistent" {
		t.Fatalf("model inexistent: %+v", x.Eroare)
	}
}
