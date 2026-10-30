package main

import (
	"pontix/internal/datebaza"
	"testing"
)

func TestSalveazaSchitaDinInterfataActualizeazaPlanificarea(t *testing.T) {
	a := appCuDemo(t)
	var id string
	if err := a.db.QueryRow(`SELECT id FROM angajati WHERE sters_la IS NULL LIMIT 1`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	if _, err := a.db.Exec(`UPDATE angajati SET categorie_program_id='program-mixt' WHERE id=?`, id); err != nil {
		t.Fatal(err)
	}
	r := a.SalveazaSchita(2026, 10, []datebaza.CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "S7"}})
	if r.Eroare != nil || r.Schita == nil {
		t.Fatalf("salvare: %+v", r)
	}
	var cod, sectiune string
	if err := a.db.QueryRow(`SELECT s.cod,p.sectiune FROM planificare p JOIN tipuri_schimb s ON s.id=p.schimb_id WHERE p.angajat_id=? AND p.data='2026-10-01' AND p.activ=1`, id).Scan(&cod, &sectiune); err != nil {
		t.Fatal(err)
	}
	if cod != "S7" || sectiune != "DUPA_AMIAZA" {
		t.Fatal(cod, sectiune)
	}
}
