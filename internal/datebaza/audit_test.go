package datebaza

import (
	"context"
	"testing"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/calcul"
)

func TestJurnalReguliConfigurabile(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	r := calcul.RegulaConfigurabila{Cod: "proba", Denumire: "Probă", TipRegula: "PLAFON_LUNAR", Severitate: "INFO", CoduriSchimb: "G24", ValoareNumerica: 3, Activ: true}
	if err := m.serviciu.SalveazaRegulaConfigurabila(ctx, r); err != nil {
		t.Fatal(err)
	}
	var id string
	if err := m.db.QueryRow(`SELECT id FROM reguli_configurabile WHERE cod = 'PROBA'`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	r.ID, r.ValoareNumerica = id, 4
	if err := m.serviciu.SalveazaRegulaConfigurabila(ctx, r); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'reguli_configurabile' AND rand_id = ?`, id).Scan(&n); err != nil || n != 2 {
		t.Fatalf("intrări de jurnal: %d, %v", n, err)
	}
	r.ID = "inexistent"
	verificaCod(t, m.serviciu.SalveazaRegulaConfigurabila(ctx, r), "date.inexistent")
}

func TestIncidentulBlocheazaServiciul(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := audit.Initializeaza(ctx, m.db); err != nil {
		t.Fatal(err)
	}
	if _, err := m.db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('strecurat', 'da', 'x')`); err != nil {
		t.Fatal(err)
	}
	rez, err := audit.Verifica(ctx, m.db)
	if err != nil || rez.Stare != audit.StareIncident {
		t.Fatalf("verificare: %+v, %v", rez, err)
	}
	db, err := bd.DeschideCuAudit(m.cale, m.cheie, audit.NouaPaza(rez))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := Nou(db)
	verificaCod(t, s.SalveazaAngajat(ctx, angajatTest("M1")), audit.CodEditareBlocata)
	if _, err := s.Citeste(ctx); err != nil {
		t.Fatalf("citirea trebuie să meargă: %v", err)
	}
}
