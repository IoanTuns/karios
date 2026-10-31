package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"testing"
)

func valoareSoldTest(t *testing.T, m *mediuTest, id string, an, luna int, cod string) calcul.ValoarePontaj {
	t.Helper()
	r, _, e := CalculeazaPontaj(context.Background(), acces.New(m.db), an, luna, nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, v := range r.Totaluri {
		if v.AngajatID == id {
			return v.Valori[cod]
		}
	}
	t.Fatal("sold lipsa")
	return calcul.ValoarePontaj{}
}
func TestSoldLunarConcediiReportSiInitializare(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	a.Initializare = &acces.InitializarePontaj{Luna: "2026-10", Minute: 120, COZile: 5}
	if e := m.serviciu.SalveazaAngajat(ctx, a); e != nil {
		t.Fatal(e)
	}
	if e := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-06", Valoare: "Z7"}}); e != nil {
		t.Fatal(e)
	}
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true, StareAprobare: "APROBATA"}}); e != nil {
		t.Fatal(e)
	}
	if v := valoareSoldTest(t, m, a.ID, 2026, 10, "ore_lucratoare"); !v.Definita || v.Sutimi != 21*420*100 {
		t.Fatal(v)
	}
	fin := valoareSoldTest(t, m, a.ID, 2026, 10, "L_FIN")
	if fin.Sutimi != (120+420-21*420)*100 {
		t.Fatal(fin)
	}
	if v := valoareSoldTest(t, m, a.ID, 2026, 11, "L_INCEPUT"); v.Sutimi != fin.Sutimi {
		t.Fatal(v, fin)
	}
	dec := valoareSoldTest(t, m, a.ID, 2026, 12, "L_FIN")
	if v := valoareSoldTest(t, m, a.ID, 2027, 1, "L_INCEPUT"); v.Sutimi != dec.Sutimi {
		t.Fatal(v, dec)
	}
	got := citeste(t, m.serviciu).Angajati[0].Initializare
	if got == nil || got.COZile != 5 || got.Revizie != 1 {
		t.Fatal(got)
	}
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, a), "pontaj.conflict")
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-09-01", Cod: "CO", Activ: true}}), "sold.co_dublat")
}
func TestSoldInitialFaraDublareCO(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-09-01", Cod: "CO", Activ: true}}); e != nil {
		t.Fatal(e)
	}
	a.Initializare = &acces.InitializarePontaj{Luna: "2026-10", COZile: 1}
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, a), "sold.co_dublat")
}

func TestSoldInitialProtejeazaLunileInchise(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	a.Initializare = &acces.InitializarePontaj{Luna: "2026-10", Minute: -15}
	if e := m.serviciu.SalveazaAngajat(ctx, a); e != nil {
		t.Fatal(e)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	if _, e := m.serviciu.SeteazaStareLuna(ctx, 2026, 11, "BLOCATA", "test", "operator"); e != nil {
		t.Fatal(e)
	}
	a.Initializare.Minute = 30
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, a), "pontaj.luna_inchisa")
	if v := citeste(t, m.serviciu).Angajati[0].Initializare.Minute; v != -15 {
		t.Fatal("sold partial", v)
	}
}
func TestSoldNormaPartialaSiGardaFaraNorma(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	n := int64(240)
	a.TipContract = "PARTIALA"
	a.ContractID = "PARTIALA"
	a.NormaPartialaMinute = &n
	if e := m.serviciu.SalveazaAngajat(ctx, a); e != nil {
		t.Fatal(e)
	}
	if v := valoareSoldTest(t, m, a.ID, 2026, 10, "ore_lucratoare"); v.Sutimi != 22*240*100 {
		t.Fatal(v)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	a.ContractID = "GARDA"
	if e := m.serviciu.SalveazaAngajat(ctx, a); e != nil {
		t.Fatal(e)
	}
	if v := valoareSoldTest(t, m, a.ID, 2026, 10, "ore_lucratoare"); v.Definita {
		t.Fatal(v)
	}
}
func TestSoldConcediuNeaprobatNuReduceNorma(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true, StareAprobare: "IN_ASTEPTARE"}}); e != nil {
		t.Fatal(e)
	}
	if v := valoareSoldTest(t, m, a.ID, 2026, 10, "ore_lucratoare"); v.Sutimi != 22*420*100 {
		t.Fatal(v)
	}
}
