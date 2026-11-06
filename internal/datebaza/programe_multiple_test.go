package datebaza

import (
	"context"
	"testing"
)

func TestProgrameMultiplePublicareAnulareSiSuprapuneri(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	// S7 15-22 si IN9 22-07 sunt adiacente, pe aceeasi marca.
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["S7"], Sectiune: "DUPA_AMIAZA", Activ: true}, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["IN9"], Sectiune: "GARZI", Activ: true})
	if len(p.Pozitii) != 2 {
		t.Fatal(p)
	}
	c := CererePublicarePlan{An: 2026, Luna: 10, PozitiiID: []string{p.Pozitii[0].ID}, Sectiune: "TOATE"}
	prev := previzualizareTest(t, m, &c)
	if len(prev.Randuri) != 2 {
		t.Fatal(prev)
	}
	if _, e := m.serviciu.PublicaPlanificare(ctx, c); e != nil {
		t.Fatal(e)
	}
	sch, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil || len(sch.Celule) != 1 || sch.Celule[0].Valoare != "IN9+S7" {
		t.Fatal(sch, e)
	}
	p, e = m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if e != nil {
		t.Fatal(e)
	}
	for _, v := range p.Pozitii {
		if v.Stare != "PUBLICATA" {
			t.Fatal(v)
		}
	}
	prev = previzualizareTest(t, m, &c)
	if prev.Identice != 2 {
		t.Fatal(prev)
	}
	if _, e = m.serviciu.PublicaPlanificare(ctx, c); e != nil {
		t.Fatal(e)
	}
	err := m.serviciu.SalveazaPlanificare(ctx, 2026, 10, []PozitiePlanificare{{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["G24"], Sectiune: "GARZI", Activ: true}})
	verificaCod(t, err, "planificare.motiv_suprapunere")
	// O alta zi care incepe la07 nu se suprapune cu tura pana la07.
	if _, e = m.db.Exec(`UPDATE tipuri_schimb SET ora_start=360,durata_minute=60 WHERE id=?`, coduri["G24"]); e != nil {
		t.Fatal(e)
	}
	err = m.serviciu.SalveazaPlanificare(ctx, 2026, 10, []PozitiePlanificare{{AngajatID: id, Data: "2026-10-31", SchimbID: coduri["G24"], Sectiune: "GARZI", Activ: true}})
	verificaCod(t, err, "planificare.motiv_suprapunere")
	p.Pozitii[0].Activ = false
	if e = m.serviciu.SalveazaPlanificare(ctx, 2026, 10, []PozitiePlanificare{p.Pozitii[0]}); e != nil {
		t.Fatal(e)
	}
	sch, e = m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil || len(sch.Celule) != 1 || sch.Celule[0].Valoare == "IN9+S7" || sch.Celule[0].Valoare == "" {
		t.Fatal(sch, e)
	}
}

func TestProgrameMultipleSincronizareInversa(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := context.Background()
	if e := m.serviciu.SalveazaSchitaCuPlanificare(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-30", Valoare: "S7+IN9"}}); e != nil {
		t.Fatal(e)
	}
	p, e := m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if e != nil || len(p.Pozitii) != 2 {
		t.Fatal(p, e)
	}
	sch, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil {
		t.Fatal(e)
	}
	c := sch.Celule[0]
	c.Valoare = "S7"
	if e = m.serviciu.SalveazaSchitaCuPlanificare(ctx, 2026, 10, []CelulaSchita{c}); e != nil {
		t.Fatal(e)
	}
	p, e = m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if e != nil {
		t.Fatal(e)
	}
	n := 0
	for _, v := range p.Pozitii {
		if v.Activ {
			n++
			if v.Stare != "PUBLICATA" {
				t.Fatal(v)
			}
		}
	}
	if n != 1 {
		t.Fatal(p)
	}
	c = sch.Celule[0]
	c.Valoare = "G24+S7"
	sch, _ = m.serviciu.CitesteSchita(ctx, 2026, 10)
	c.Revizie = sch.Celule[0].Revizie
	err := m.serviciu.SalveazaSchitaCuPlanificare(ctx, 2026, 10, []CelulaSchita{c})
	verificaCod(t, err, "planificare.motiv_suprapunere")
}

func TestProgrameMultipleMedicalPastreazaIntervalulNeintersectat(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-04", SchimbID: coduri["S7"], Sectiune: "DUPA_AMIAZA", Activ: true}, PozitiePlanificare{AngajatID: id, Data: "2026-10-04", SchimbID: coduri["IN9"], Sectiune: "GARZI", Activ: true})
	c := cerereTest(p, "GARZI")
	previzualizareTest(t, m, &c)
	if _, e := m.serviciu.PublicaPlanificare(ctx, c); e != nil {
		t.Fatal(e)
	}
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: id, Data: "2026-10-05", Cod: "CM", Activ: true, StareAprobare: "APROBATA"}}); e != nil {
		t.Fatal(e)
	}
	var val string
	if e := m.db.QueryRow(`SELECT valoare FROM pontaj_ciorne WHERE angajat_id=? AND data='2026-10-04'`, id).Scan(&val); e != nil || val != "S7" {
		t.Fatal(val, e)
	}
	var n int
	if e := m.db.QueryRow(`SELECT count(*) FROM planificare WHERE angajat_id=? AND activ=1`, id).Scan(&n); e != nil || n != 1 {
		t.Fatal(n, e)
	}
}

func TestProgrameMultipleAnulareIntregSetulPublicat(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["S7"], Sectiune: "DUPA_AMIAZA", Activ: true}, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["IN9"], Sectiune: "GARZI", Activ: true})
	c := cerereTest(p, "GARZI")
	previzualizareTest(t, m, &c)
	if _, e := m.serviciu.PublicaPlanificare(ctx, c); e != nil {
		t.Fatal(e)
	}
	for i := range p.Pozitii {
		p.Pozitii[i].Activ = false
	}
	if e := m.serviciu.SalveazaPlanificare(ctx, 2026, 10, p.Pozitii); e != nil {
		t.Fatal(e)
	}
	var val string
	if e := m.db.QueryRow(`SELECT valoare FROM pontaj_ciorne WHERE angajat_id=?`, id).Scan(&val); e != nil || val != "" {
		t.Fatal(val, e)
	}
}

func TestProgrameMultiplePublicareBlocheazaSetCuProgramInactiv(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["S7"], Sectiune: "DUPA_AMIAZA", Activ: true}, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["IN9"], Sectiune: "GARZI", Activ: true})
	if _, e := m.db.Exec(`UPDATE tipuri_schimb SET activ=0 WHERE id=?`, coduri["IN9"]); e != nil {
		t.Fatal(e)
	}
	c := cerereTest(p, "DUPA_AMIAZA")
	r := previzualizareTest(t, m, &c)
	if r.Blocate != 2 || r.Publicabile != 0 {
		t.Fatal(r)
	}
	if _, e := m.serviciu.PublicaPlanificare(context.Background(), c); e != nil {
		t.Fatal(e)
	}
	var n int
	if e := m.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&n); e != nil || n != 0 {
		t.Fatal(n, e)
	}
}

func TestProgrameMultiplePontajRefuzaSuprapunereVecinaInLot(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	err := m.serviciu.SalveazaSchita(context.Background(), 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-30", Valoare: "S7+IN9"}, {AngajatID: id, Data: "2026-10-31", Valoare: "G24"}})
	// G24 incepe la07, adiacenta este permisa; mutam inceputul la06 pentru conflict.
	if err != nil {
		t.Fatal(err)
	}
	if _, err = m.db.Exec(`UPDATE tipuri_schimb SET ora_start=360 WHERE cod='G24'`); err != nil {
		t.Fatal(err)
	}
	sch, err := m.serviciu.CitesteSchita(context.Background(), 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	err = m.serviciu.SalveazaSchita(context.Background(), 2026, 10, sch.Celule)
	verificaCod(t, err, "planificare.motiv_suprapunere")
}

func TestProgrameMultipleProgramNormalSiTura(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["Z7"], Sectiune: "DUPA_AMIAZA", Activ: true}, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["S7"], Sectiune: "DUPA_AMIAZA", Activ: true})
	c := cerereTest(p, "DUPA_AMIAZA")
	previzualizareTest(t, m, &c)
	if _, e := m.serviciu.PublicaPlanificare(ctx, c); e != nil {
		t.Fatal(e)
	}
	sch, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil || sch.Celule[0].Valoare != "S7+Z7" {
		t.Fatal(sch, e)
	}
	if e = m.serviciu.SalveazaSchitaCuPlanificare(ctx, 2026, 10, sch.Celule); e != nil {
		t.Fatal(e)
	}
	p, e = m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if e != nil {
		t.Fatal(e)
	}
	n := 0
	for _, v := range p.Pozitii {
		if v.Activ {
			n++
		}
	}
	if n != 2 {
		t.Fatal(p)
	}
}

func TestProgramNormalPlanificareExcludeWeekendSiSarbatori(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	for _, zi := range []string{"2026-10-31", "2026-12-01"} {
		luna := 10
		if zi[5:7] == "12" {
			luna = 12
		}
		err := m.serviciu.SalveazaPlanificare(context.Background(), 2026, luna, []PozitiePlanificare{{AngajatID: id, Data: zi, SchimbID: coduri["Z7"], Sectiune: "DUPA_AMIAZA", Activ: true}})
		verificaCod(t, err, "planificare.motiv_zi_nelucratoare")
	}
}
