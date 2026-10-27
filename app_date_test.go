package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/datebaza"
)

func TestMetodeInterfataPastreazaCodurileSiValorileNule(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "ui.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	a := NewApp(db)
	refuzat := a.SalveazaAngajat(datebaza.Angajat{})
	if refuzat.Eroare == nil || refuzat.Eroare.Cod != "date.campuri_obligatorii" || refuzat.Date != nil {
		t.Fatalf("răspuns de eroare: %+v", refuzat)
	}
	r := a.IncarcaDemo()
	if r.Eroare != nil || r.Date == nil || len(r.Date.Angajati) != 12 {
		t.Fatalf("demo prin Wails: %+v", r)
	}
	text, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var decodat RaspunsDate
	if err := json.Unmarshal(text, &decodat); err != nil {
		t.Fatal(err)
	}
	gasit := false
	for _, angajat := range decodat.Date.Angajati {
		if angajat.Marca == "DEMO-006" {
			gasit = true
			if angajat.ConcediuZileAn != nil {
				t.Fatal("concediu NULL convertit la zero")
			}
		}
	}
	if !gasit {
		t.Fatal("angajatul de referință lipsește")
	}
}

func TestAppAtribuieCategorieProgramLot(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "lot.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	a := NewApp(db)
	rDemo := a.IncarcaDemo()
	if rDemo.Eroare != nil || rDemo.Date == nil || len(rDemo.Date.Angajati) < 3 {
		t.Fatalf("eroare la incarcare demo: %+v", rDemo.Eroare)
	}

	id1 := rDemo.Date.Angajati[0].ID
	id2 := rDemo.Date.Angajati[1].ID

	// Preluăm o categorie existentă
	if len(rDemo.Date.CategoriiProgram) == 0 {
		t.Fatalf("nu există categorii în demo")
	}
	catID := rDemo.Date.CategoriiProgram[0].ID

	resp := a.AtribuieCategorieProgramLot([]string{id1, id2}, catID)
	if resp.Eroare != nil || resp.Date == nil {
		t.Fatalf("eroare la atribuire lot prin App: %+v", resp.Eroare)
	}
	for _, ang := range resp.Date.Angajati {
		if ang.ID == id1 || ang.ID == id2 {
			if ang.CategorieProgramID != catID {
				t.Fatalf("angajatul %s nu a primit categoria %s, are %s", ang.Marca, catID, ang.CategorieProgramID)
			}
		}
	}
}
