package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/datebaza"
)

func TestAppCitestePlanificatVsRealizat(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "comp.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}

	a := NewApp(db)
	rDemo := a.IncarcaDemo()
	if rDemo.Eroare != nil || rDemo.Date == nil || len(rDemo.Date.Angajati) == 0 {
		t.Fatalf("eroare la încărcare demo: %+v", rDemo.Eroare)
	}

	ang1 := rDemo.Date.Angajati[0].ID

	// Planificăm o gardă G24 pe 2026-10-03
	resPlan := a.SalveazaPlanificare(2026, 10, []datebaza.PozitiePlanificare{
		{
			ID:        "",
			AngajatID: ang1,
			Data:      "2026-10-03",
			SchimbID:  "10000000-0000-4000-8000-000000000004", // G24
			Sectiune:  "GARZI",
			Activ:     true,
		},
	})
	if resPlan.Eroare != nil {
		t.Fatalf("eroare la salvare planificare: %+v", resPlan.Eroare)
	}

	// Citim raportul comparativ
	respComp := a.CitestePlanificatVsRealizat(2026, 10)
	if respComp.Eroare != nil || respComp.Raport == nil {
		t.Fatalf("eroare la CitestePlanificatVsRealizat: %+v", respComp.Eroare)
	}

	if respComp.Raport.An != 2026 || respComp.Raport.Luna != 10 {
		t.Fatalf("perioadă incorectă: %d/%d", respComp.Raport.An, respComp.Raport.Luna)
	}
	if len(respComp.Raport.Randuri) == 0 {
		t.Fatalf("nu s-au generat rânduri de comparație")
	}

	// Angajatul 1 trebuie să aibă 24h planificate
	var gasit bool
	for _, r := range respComp.Raport.Randuri {
		if r.AngajatID == ang1 {
			gasit = true
			if r.TotalOrePlanificate != 24 {
				t.Errorf("așteptat 24h planificate pentru angajatul 1, obținut: %d", r.TotalOrePlanificate)
			}
			// Înainte de pontare, pe 3 oct starea trebuie să fie NEEFECTUAT
			for _, z := range r.Zile {
				if z.Data == "2026-10-03" {
					if z.Stare != "NEEFECTUAT" {
						t.Errorf("stare așteptată pe 2026-10-03: NEEFECTUAT, obținut: %s", z.Stare)
					}
				}
			}
		}
	}
	if !gasit {
		t.Fatalf("angajatul planificat nu a fost găsit în raportul comparativ")
	}
}
