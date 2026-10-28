package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/calcul"
)

func TestAppReguliConfigurabile(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "test_reguli.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}

	app := &App{db: db}

	// 1. Citire reguli predefinite din migrarea 0020
	rLista := app.ListaReguliConfigurabile()
	if rLista.Eroare != nil {
		t.Fatalf("eroare citire reguli: %+v", rLista.Eroare)
	}
	if len(rLista.Reguli) < 3 {
		t.Fatalf("așteptat minim 3 reguli predefinite, obținut: %d", len(rLista.Reguli))
	}

	// 2. Modificare valoare regulă
	regulaG24 := rLista.Reguli[0]
	regulaG24.ValoareNumerica = 36 // creștere repaus la 36h
	rSalvare := app.SalveazaRegulaConfigurabila(regulaG24)
	if rSalvare.Eroare != nil {
		t.Fatalf("eroare salvare regulă: %+v", rSalvare.Eroare)
	}

	// 3. Adăugare regulă nouă
	rNoua := calcul.RegulaConfigurabila{
		Cod:             "REG_CUSTOM_NOAPTE",
		Denumire:        "Plafon ture noapte",
		Descriere:       "Maxim 4 ture de noapte pe lună",
		TipRegula:       "PLAFON_LUNAR",
		Severitate:      "INFO",
		CoduriSchimb:    "N12,N9",
		ValoareNumerica: 4,
		Activ:           true,
	}
	rSalvareNoua := app.SalveazaRegulaConfigurabila(rNoua)
	if rSalvareNoua.Eroare != nil {
		t.Fatalf("eroare adăugare regulă nouă: %+v", rSalvareNoua.Eroare)
	}

	// Verificare listă reîncărcată
	rLista2 := app.ListaReguliConfigurabile()
	if len(rLista2.Reguli) != len(rLista.Reguli)+1 {
		t.Errorf("așteptat %d reguli, obținut: %d", len(rLista.Reguli)+1, len(rLista2.Reguli))
	}

	// 4. Verificare reguli planificare pe lună goală
	rVerif := app.VerificaReguliPlanificare(2026, 10)
	if rVerif.Eroare != nil {
		t.Fatalf("eroare verificare planificare: %+v", rVerif.Eroare)
	}
}
