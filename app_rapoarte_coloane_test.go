package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/rapoarte"
)

// F09-T27: editorul afla coloanele SQL fara PDF. Previzualizarea in modul vechi
// „tabel” e refuzata din F09-T23, deci editorul nu se mai poate baza pe ea.
func TestColoaneModelRaportFaraMotorulVechi(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "coloane.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	a := NewApp(db)
	modele := a.ListaModeleRapoarte()
	if modele.Eroare != nil {
		t.Fatal(modele.Eroare)
	}
	var d rapoarte.Definitie
	for _, m := range modele.Configuratie.Modele {
		if m.ID == "raport-pontaj-sql" {
			d = m
		}
	}
	if d.Continut.Flexibil == nil {
		t.Fatal("modelul SQL predefinit lipseste")
	}
	f := *d.Continut.Flexibil
	f.ModAspect, f.Sablon = "TABEL", ""
	d.Continut.Flexibil = &f
	c := rapoarte.Cerere{An: 2026, Luna: 10, Limba: "ro"}

	r := a.ColoaneModelRaport(d, c)
	if r.Eroare != nil || r.Document == nil || len(r.Document.Coloane) == 0 || r.PDF != "" || r.Token != "" {
		t.Fatalf("coloane: %+v %+v", r.Eroare, r.Document)
	}
	if len(r.Document.Randuri) != 0 {
		t.Fatal("raspunsul pentru coloane nu trebuie sa poarte datele raportului")
	}
	// Contractul care a produs defectul: previzualizarea in modul tabel e refuzata.
	if p := a.PrevizualizeazaModelRaport(d, c); p.Eroare == nil || p.Eroare.Cod != "rapoarte.model_neconvertit" {
		t.Fatalf("previzualizarea in modul tabel: %+v", p.Eroare)
	}
	// SQL invalid: eroarea ajunge la editor.
	f.SQL = "SELECT coloana_inexistenta FROM angajati"
	if r := a.ColoaneModelRaport(d, c); r.Eroare == nil {
		t.Fatal("SQL invalid acceptat")
	}
}
