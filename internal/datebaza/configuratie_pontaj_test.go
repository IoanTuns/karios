package datebaza

import (
	"context"
	"pontix/internal/calcul"
	"pontix/internal/i18n"
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

// F04-T08 (D102): optiunea pentru noptile cu schimbarea orei se salveaza, iar schimbarea ei
// e refuzata cat timp e inchisa o luna care contine o astfel de noapte.
func TestSchimbareOraSalvareSiLuniInchise(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	c, _, e := m.serviciu.ConfiguratiePontaj(ctx)
	if e != nil || c.SchimbareOra != calcul.SchimbareOraNominala {
		t.Fatal("implicit", c.SchimbareOra, e)
	}
	// Noiembrie nu are schimbarea orei: inchiderea ei nu blocheaza optiunea.
	if _, e = m.serviciu.SeteazaStareLuna(ctx, 2026, 11, "INCHISA", "test", "operator"); e != nil {
		t.Fatal(e)
	}
	c.SchimbareOra = calcul.SchimbareOraEfectiva
	if e = m.serviciu.SalveazaConfiguratiePontaj(ctx, c); e != nil {
		t.Fatal(e)
	}
	c, _, _ = m.serviciu.ConfiguratiePontaj(ctx)
	if c.SchimbareOra != calcul.SchimbareOraEfectiva {
		t.Fatal("salvata", c.SchimbareOra)
	}
	// Octombrie inchisa: optiunea nu se mai poate schimba, dar celelalte setari da.
	if _, e = m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "BLOCATA", "test", "operator"); e != nil {
		t.Fatal(e)
	}
	c.SchimbareOra = calcul.SchimbareOraNominala
	e = m.serviciu.SalveazaConfiguratiePontaj(ctx, c)
	verificaCod(t, e, "pontaj.schimbare_ora_luna_inchisa")
	if er, ok := e.(*calcul.Eroare); !ok || er.Param["luna"] != "10.2026" {
		t.Fatal("parametru", e)
	}
	c.SchimbareOra = calcul.SchimbareOraEfectiva
	c.Atribuire = calcul.AtribuireCalendar
	if e = m.serviciu.SalveazaConfiguratiePontaj(ctx, c); e != nil {
		t.Fatal("alte setari", e)
	}
	// Dupa redeschidere, optiunea se poate schimba.
	if _, e = m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "DESCHISA", "corectie", "operator"); e != nil {
		t.Fatal(e)
	}
	c, _, _ = m.serviciu.ConfiguratiePontaj(ctx)
	c.SchimbareOra = calcul.SchimbareOraNominala
	if e = m.serviciu.SalveazaConfiguratiePontaj(ctx, c); e != nil {
		t.Fatal("dupa redeschidere", e)
	}
}

// F04-T11 (D107): tipul zilei se salveaza, iar optiunea ZNL e refuzata cat timp o luna cu
// zile libere adaugate e inchisa.
func TestZileLibereZNLSalvareSiLuniInchise(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if e := m.serviciu.SalveazaSarbatoare(ctx, "", Sarbatoare{Data: "2026-05-04", Denumire: i18n.Text{"ro": "Punte"}, Tip: TipZiLibera}); e != nil {
		t.Fatal(e)
	}
	verificaCod(t, m.serviciu.SalveazaSarbatoare(ctx, "", Sarbatoare{Data: "2026-05-05", Denumire: i18n.Text{"ro": "X"}, Tip: "ALTUL"}), "sarbatori.tip_invalid")
	cal, e := m.serviciu.CitesteCalendar(ctx, 2026)
	if e != nil {
		t.Fatal(e)
	}
	gasit := false
	for _, s := range cal.Sarbatori {
		gasit = gasit || (s.Data == "2026-05-04" && s.Tip == TipZiLibera)
	}
	if !gasit {
		t.Fatal("tipul zilei", cal.Sarbatori)
	}
	c, _, e := m.serviciu.ConfiguratiePontaj(ctx)
	if e != nil || c.ZileLibereZNL != calcul.ZileLibereZNLToate {
		t.Fatal("implicit", c.ZileLibereZNL, e)
	}
	if _, e = m.serviciu.SeteazaStareLuna(ctx, 2026, 5, "INCHISA", "test", "operator"); e != nil {
		t.Fatal(e)
	}
	c.ZileLibereZNL = calcul.ZileLibereZNLLegale
	verificaCod(t, m.serviciu.SalveazaConfiguratiePontaj(ctx, c), "pontaj.zile_libere_luna_inchisa")
	if _, e = m.serviciu.SeteazaStareLuna(ctx, 2026, 5, "DESCHISA", "corectie", "operator"); e != nil {
		t.Fatal(e)
	}
	if e = m.serviciu.SalveazaConfiguratiePontaj(ctx, c); e != nil {
		t.Fatal(e)
	}
	if c, _, _ = m.serviciu.ConfiguratiePontaj(ctx); c.ZileLibereZNL != calcul.ZileLibereZNLLegale {
		t.Fatal("salvata", c.ZileLibereZNL)
	}
}
