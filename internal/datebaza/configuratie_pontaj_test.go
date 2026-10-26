package datebaza

import (
	"context"
	"pontix/internal/calcul"
	"testing"
)

func TestReguliPontajPersistentaConflictSiJurnal(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	c, _, e := m.serviciu.ConfiguratiePontaj(ctx)
	if e != nil {
		t.Fatal(e)
	}
	c.Atribuire = calcul.AtribuireCalendar
	if e = m.serviciu.SalveazaConfiguratiePontaj(ctx, c); e != nil {
		t.Fatal(e)
	}
	verificaCod(t, m.serviciu.SalveazaConfiguratiePontaj(ctx, c), "pontaj.reguli_conflict")
	actual, _, e := m.serviciu.ConfiguratiePontaj(ctx)
	if e != nil || actual.Revizie != "2" || actual.Atribuire != calcul.AtribuireCalendar {
		t.Fatal(actual, e)
	}
	actual.Reguli[3].Coduri = []string{"INEXISTENT"}
	verificaCod(t, m.serviciu.SalveazaConfiguratiePontaj(ctx, actual), "pontaj.reguli_cod_lipsa")
	var n int
	if e = m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel='configuratie_pontaj'`).Scan(&n); e != nil || n != 1 {
		t.Fatal(n, e)
	}
	actual.Reguli[3].Coduri = []string{"CO"}
	if _, e = m.db.Exec(`CREATE TRIGGER refuza_reguli BEFORE INSERT ON jurnal_modificari WHEN NEW.tabel='configuratie_pontaj' BEGIN SELECT RAISE(ABORT,'proba'); END`); e != nil {
		t.Fatal(e)
	}
	actual.Atribuire = calcul.AtribuireInceput
	if e = m.serviciu.SalveazaConfiguratiePontaj(ctx, actual); e == nil {
		t.Fatal("salvare fara audit")
	}
	final, _, e := m.serviciu.ConfiguratiePontaj(ctx)
	if e != nil || final.Revizie != "2" || final.Atribuire != calcul.AtribuireCalendar {
		t.Fatal(final, e)
	}
}

func TestPrevizualizarePontajNuSalveazaSiIncludeLunaPrecedenta(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if e := m.serviciu.SalveazaAngajat(ctx, angajatTest("demo-calc")); e != nil {
		t.Fatal(e)
	}
	id := citeste(t, m.serviciu).Angajati[0].ID
	if e := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-31", Valoare: "N9"}}); e != nil {
		t.Fatal(e)
	}
	c, _, _ := m.serviciu.ConfiguratiePontaj(ctx)
	c.Atribuire = calcul.AtribuireCalendar
	if e := m.serviciu.SalveazaConfiguratiePontaj(ctx, c); e != nil {
		t.Fatal(e)
	}
	luna, e := m.serviciu.CitesteSchita(ctx, 2026, 11)
	if e != nil || len(luna.Calcul.Celule) != 1 || luna.Calcul.Celule[0].Minute != 420 {
		t.Fatal(luna, e)
	}
	prev, e := m.serviciu.CitesteSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-31", Valoare: "@M450"}})
	if e != nil || prev.Celule[0].Valoare != "N9" || prev.Calcul.Celule[0].Minute != 450 {
		t.Fatal(prev, e)
	}
	real, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil || real.Celule[0].Valoare != "N9" || real.Calcul.Celule[0].Minute != 120 {
		t.Fatal(real, e)
	}
}
