package datebaza

import (
	"context"
	"testing"
)

func TestCitesteJurnalFiltre(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	for _, marca := range []string{"M1", "M2"} {
		if err := m.serviciu.SalveazaAngajat(ctx, angajatTest(marca)); err != nil {
			t.Fatal(err)
		}
	}
	d := citeste(t, m.serviciu)
	a1 := d.Angajati[0].ID
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a1, Data: "2026-10-01", Valoare: "D8"}, {AngajatID: d.Angajati[1].ID, Data: "2026-10-02", Valoare: "D8"}}); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 11, []CelulaSchita{{AngajatID: a1, Data: "2026-11-03", Valoare: "N12"}}); err != nil {
		t.Fatal(err)
	}
	pagina := func(f FiltruJurnal) PaginaJurnal {
		t.Helper()
		p, err := m.serviciu.CitesteJurnal(ctx, f)
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	if p := pagina(FiltruJurnal{}); p.Total != 5 || p.Intrari[0].RandID != a1+"/2026-11-03" || len(p.Tabele) != 2 {
		t.Fatalf("toate: %d, primul %+v, tabele %v", p.Total, p.Intrari[0], p.Tabele)
	}
	if p := pagina(FiltruJurnal{AngajatID: a1}); p.Total != 3 {
		t.Fatalf("angajat: %d", p.Total)
	}
	if p := pagina(FiltruJurnal{Luna: "2026-10", Tabel: "pontaj_ciorne"}); p.Total != 2 {
		t.Fatalf("lună: %d", p.Total)
	}
	if p := pagina(FiltruJurnal{Text: "n12"}); p.Total != 1 {
		t.Fatalf("text: %d", p.Total)
	}
	if p := pagina(FiltruJurnal{Limita: 2, Deplasare: 4}); len(p.Intrari) != 1 || p.Total != 5 {
		t.Fatalf("paginare: %d / %d", len(p.Intrari), p.Total)
	}
	if tot, err := m.serviciu.TotJurnalul(ctx, FiltruJurnal{AngajatID: a1}); err != nil || len(tot) != 3 {
		t.Fatalf("export: %d, %v", len(tot), err)
	}
}
