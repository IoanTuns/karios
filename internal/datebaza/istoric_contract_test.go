package datebaza

import (
	"context"
	"testing"
)

func p64(v int64) *int64 { return &v }

// F07-T43 (I010, D101): norma pe zile, concediul proportional, lunile inchise protejate, jurnal.
func TestIstoricContractNormaConcediuSiLuniInchise(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	a.ConcediuZileAn = p64(21)
	if e := m.serviciu.SalveazaAngajat(ctx, a); e != nil {
		t.Fatal(e)
	}
	// Fara istoric: octombrie 2026 are 22 de zile lucratoare x 7 h (norma globala din test).
	if v := valoareSoldTest(t, m, a.ID, 2026, 10, "ore_lucratoare"); v.Sutimi != 22*420*100 {
		t.Fatal("fara istoric", v)
	}
	if e := m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-10-15", NormaMinute: p64(240)}); e != nil {
		t.Fatal(e)
	}
	// 1-14 octombrie: 10 zile x 7 h; 15-31 octombrie: 12 zile x 4 h.
	if v := valoareSoldTest(t, m, a.ID, 2026, 10, "ore_lucratoare"); v.Sutimi != (10*420+12*240)*100 {
		t.Fatal("norma pe zile", v)
	}
	if v := valoareSoldTest(t, m, a.ID, 2026, 9, "ore_lucratoare"); v.Sutimi != 22*420*100 {
		t.Fatal("luna dinaintea schimbarii", v)
	}
	if n := citeste(t, m.serviciu).Angajati[0].NormaMinute; n != 240 {
		t.Fatal("norma afisata azi", n)
	}
	if e := m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-07-01", ConcediuZileAn: p64(25)}); e != nil {
		t.Fatal(e)
	}
	abs, e := m.serviciu.CitesteAbsente(ctx, 2026)
	if e != nil {
		t.Fatal(e)
	}
	// 21 x 181/365 + 25 x 184/365 = 23,02 -> 23
	if abs.DreptConcediu[a.ID] != 23 {
		t.Fatal("drept proportional", abs.DreptConcediu)
	}
	if abs, _ = m.serviciu.CitesteAbsente(ctx, 2027); abs.DreptConcediu[a.ID] != 25 {
		t.Fatal("drept anul urmator", abs.DreptConcediu)
	}
	verificaCod(t, m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-10-15", NormaMinute: p64(300)}), "istoric.exista")
	verificaCod(t, m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-10-20"}), "istoric.valoare_lipsa")
	verificaCod(t, m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-02-30", NormaMinute: p64(300)}), "istoric.data_invalida")
	verificaCod(t, m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-10-20", NormaMinute: p64(1500)}), "istoric.norma_invalida")
	verificaCod(t, m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: "inexistent", DeLa: "2026-10-20", NormaMinute: p64(300)}), "istoric.angajat_inexistent")
	// Luna inchisa: o schimbare care ar atinge noiembrie e refuzata, una din decembrie nu.
	if _, e := m.serviciu.SeteazaStareLuna(ctx, 2026, 11, "BLOCATA", "test", "operator"); e != nil {
		t.Fatal(e)
	}
	verificaCod(t, m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-10-20", NormaMinute: p64(300)}), "istoric.luna_inchisa")
	verificaCod(t, m.serviciu.StergeSchimbareContract(ctx, a.ID, "2026-10-15"), "istoric.luna_inchisa")
	if e := m.serviciu.AdaugaSchimbareContract(ctx, SchimbareContract{AngajatID: a.ID, DeLa: "2026-12-01", NormaMinute: p64(480)}); e != nil {
		t.Fatal(e)
	}
	verificaCod(t, m.serviciu.StergeSchimbareContract(ctx, a.ID, "2026-12-02"), "istoric.inexistenta")
	if e := m.serviciu.StergeSchimbareContract(ctx, a.ID, "2026-12-01"); e != nil {
		t.Fatal(e)
	}
	lista, e := m.serviciu.ListaIstoricContract(ctx, a.ID)
	if e != nil || len(lista) != 2 || lista[0].DeLa != "2026-07-01" || lista[1].DeLa != "2026-10-15" {
		t.Fatal("lista", lista, e)
	}
	var jurnale int
	if e := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel='istoric_contract'`).Scan(&jurnale); e != nil || jurnale != 4 {
		t.Fatal("jurnal", jurnale, e)
	}
}
