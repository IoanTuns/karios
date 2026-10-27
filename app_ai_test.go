package main

import (
	"context"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/ai"
	"pontix/internal/bd"
	"pontix/internal/datebaza"
)

func TestAppAutoPornireAI(t *testing.T) {
	tempDir := t.TempDir()
	a := &App{radacina: tempDir}

	folder := a.folderAI()
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}

	// Implicit trebuie să fie activă (fișierul de dezactivare nu există)
	if !a.EsteAutoPornireAI() {
		t.Fatalf("așteptat EsteAutoPornireAI() true inițial")
	}

	// Dezactivare
	if a.SeteazaAutoPornireAI(false) != false {
		t.Fatalf("așteptat false după SeteazaAutoPornireAI(false)")
	}
	if a.EsteAutoPornireAI() {
		t.Fatalf("așteptat EsteAutoPornireAI() false după dezactivare")
	}
	if _, err := os.Stat(filepath.Join(folder, fisierFaraAutoPornire)); err != nil {
		t.Fatalf("fișierul de dezactivare nu a fost creat: %v", err)
	}

	// Reactivare
	if a.SeteazaAutoPornireAI(true) != true {
		t.Fatalf("așteptat true după SeteazaAutoPornireAI(true)")
	}
	if !a.EsteAutoPornireAI() {
		t.Fatalf("așteptat EsteAutoPornireAI() true după reactivare")
	}

	// Verificare StareAI pe folder gol
	stare := a.StareAI()
	if stare.Stare.Instalat {
		t.Fatalf("folderul gol nu trebuia să fie marcat ca instalat")
	}

	// Verificare că nu dublează folderul karios (karios/karios/ai)
	a2 := &App{radacina: filepath.Join(tempDir, "karios")}
	asteptat := filepath.Join(tempDir, "karios", "ai")
	if a2.folderAI() != asteptat {
		t.Fatalf("folderAI incorect: primit %s, așteptat %s", a2.folderAI(), asteptat)
	}
}

func TestAplicaPropunereAITotAlDoileaWeekend(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "ai_test.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}

	a := NewApp(db)
	rDemo := a.IncarcaDemo()
	if rDemo.Eroare != nil {
		t.Fatalf("eroare încărcare demo: %+v", rDemo.Eroare)
	}

	// Test aplicare: Dumitru Andrei (DEMO-002), octombrie 2026, tot al doilea weekend, Z12, începând cu 3 oct
	prop := ai.PropunereActiune{
		ID:  "p1",
		Tip: "planifica_interval",
		Parametri: map[string]any{
			"marca":           "DEMO-002",
			"cod_schimb":      "Z12",
			"data_inceput":    "2026-10-03",
			"data_sfarsit":    "2026-10-31",
			"filtru_zile":     "tot_al_doilea_weekend",
			"doar_lucratoare": false,
		},
	}

	resp := a.AplicaPropunereAI(CerereAplicaPropunereAI{
		An:        2026,
		Luna:      10,
		Propunere: prop,
	})
	if resp.Eroare != nil {
		t.Fatalf("eroare la aplicare propunere: %+v", resp.Eroare)
	}
	if !resp.Succes {
		t.Fatalf("așteptat succes true")
	}

	// Verificare în Pontaj (schita)
	schita, err := datebaza.Nou(db).CitesteSchita(context.Background(), 2026, 10)
	if err != nil {
		t.Fatalf("eroare citire schiță: %v", err)
	}
	zileAsteptate := map[string]bool{
		"2026-10-03": true,
		"2026-10-04": true,
		"2026-10-17": true,
		"2026-10-18": true,
		"2026-10-31": true,
	}

	zileGasite := make(map[string]string)
	for _, c := range schita.Celule {
		if c.Valoare != "" {
			zileGasite[c.Data] = c.Valoare
		}
	}

	if len(zileGasite) != len(zileAsteptate) {
		t.Fatalf("așteptat exact %d zile planificate, obținut %d: %+v", len(zileAsteptate), len(zileGasite), zileGasite)
	}
	for z := range zileAsteptate {
		if val, ok := zileGasite[z]; !ok || val != "Z12" {
			t.Errorf("ziua %s așteptată cu Z12, obținut: %s", z, val)
		}
	}

	// Verificare în Planificare
	plan, err := datebaza.Nou(db).CitestePlanificare(context.Background(), 2026, 10)
	if err != nil {
		t.Fatalf("eroare citire planificare: %v", err)
	}
	if len(plan.Pozitii) != 5 {
		t.Fatalf("așteptat 5 poziții în planificare pentru Z12 (GARZI), obținut: %d", len(plan.Pozitii))
	}
}
