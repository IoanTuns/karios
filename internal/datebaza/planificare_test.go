package datebaza

import (
	"context"
	"testing"
)

func mediuPlanificare(t *testing.T) (*mediuTest, string, map[string]string) {
	t.Helper()
	m := pregateste(t)
	ctx := context.Background()
	a := angajatTest("PLAN")
	a.CategorieProgramID = "program-mixt"
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	schimburi, err := m.serviciu.ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]string{}
	for _, p := range schimburi {
		ids[p.Cod] = p.ID
	}
	return m, a.ID, ids
}
func planTest(t *testing.T, m *mediuTest, an, luna int, pozitii ...PozitiePlanificare) PlanificareLuna {
	t.Helper()
	if err := m.serviciu.SalveazaPlanificare(context.Background(), an, luna, pozitii); err != nil {
		t.Fatal(err)
	}
	p, err := m.serviciu.CitestePlanificare(context.Background(), an, luna)
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func cerereTest(p PlanificareLuna, sectiune string) CererePublicarePlan {
	c := CererePublicarePlan{An: p.An, Luna: p.Luna, Sectiune: sectiune}
	for _, r := range p.Pozitii {
		if r.Activ && r.Sectiune == sectiune {
			c.PozitiiID = append(c.PozitiiID, r.ID)
		}
	}
	return c
}
func previzualizareTest(t *testing.T, m *mediuTest, c *CererePublicarePlan) PublicarePlan {
	t.Helper()
	p, err := m.serviciu.PrevizualizeazaPublicarePlan(context.Background(), *c)
	if err != nil {
		t.Fatal(err)
	}
	c.Token = p.Token
	return p
}
func TestPlanificareSeparataPublicareIdempotentaSiRetragere(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-01", SchimbID: schimburi["S7"], Sectiune: "DUPA_AMIAZA", Activ: true})
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&n); err != nil || n != 0 {
		t.Fatal("salvarea planului a modificat pontajul", n, err)
	}
	c := cerereTest(p, "DUPA_AMIAZA")
	r := previzualizareTest(t, m, &c)
	if r.Publicabile != 1 || r.Randuri[0].Stare != "NOUA" {
		t.Fatalf("previzualizare: %+v", r)
	}
	if _, err := m.serviciu.PublicaPlanificare(ctx, c); err != nil {
		t.Fatal(err)
	}
	p, err := m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if err != nil || p.Pozitii[0].Stare != "PUBLICATA" {
		t.Fatalf("stare: %+v %v", p, err)
	}
	var revizie int
	if err := m.db.QueryRow(`SELECT revizie FROM pontaj_ciorne`).Scan(&revizie); err != nil {
		t.Fatal(err)
	}
	r = previzualizareTest(t, m, &c)
	if r.Identice != 1 || r.Publicabile != 0 {
		t.Fatalf("republicare: %+v", r)
	}
	if _, err := m.serviciu.PublicaPlanificare(ctx, c); err != nil {
		t.Fatal(err)
	}
	var noua int
	if err := m.db.QueryRow(`SELECT revizie FROM pontaj_ciorne`).Scan(&noua); err != nil || noua != revizie {
		t.Fatalf("ore duplicate: %d %d %v", revizie, noua, err)
	}
	retras := p.Pozitii[0]
	retras.Activ = false
	planTest(t, m, 2026, 10, retras)
	var valoare string
	if err := m.db.QueryRow(`SELECT valoare FROM pontaj_ciorne`).Scan(&valoare); err != nil || valoare != "S7" {
		t.Fatal("retragerea planului a sters pontajul", valoare, err)
	}
}
func TestPublicareInlocuireExplicitaSiConflictPontaj(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "Z7"}}); err != nil {
		t.Fatal(err)
	}
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-01", SchimbID: schimburi["S7"], Sectiune: "DUPA_AMIAZA", Activ: true})
	c := cerereTest(p, "DUPA_AMIAZA")
	r := previzualizareTest(t, m, &c)
	if r.Blocate != 1 || r.Randuri[0].Motiv != "celula_ocupata" {
		t.Fatalf("protectie: %+v", r)
	}
	c.Inlocuieste = true
	r = previzualizareTest(t, m, &c)
	if r.Publicabile != 1 || r.Randuri[0].Stare != "INLOCUIRE" || r.Randuri[0].Inainte != "Z7" {
		t.Fatalf("inlocuire: %+v", r)
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "@M480", Revizie: "1"}}); err != nil {
		t.Fatal(err)
	}
	_, err := m.serviciu.PublicaPlanificare(ctx, c)
	verificaCod(t, err, "planificare.conflict_publicare")
	previzualizareTest(t, m, &c)
	if _, err := m.serviciu.PublicaPlanificare(ctx, c); err != nil {
		t.Fatal(err)
	}
	var valoare string
	if err := m.db.QueryRow(`SELECT valoare FROM pontaj_ciorne`).Scan(&valoare); err != nil || valoare != "S7" {
		t.Fatalf("valoare: %s %v", valoare, err)
	}
}
func TestPublicareReverificaAbsentaDinLunaUrmatoare(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-31", SchimbID: schimburi["G24"], Sectiune: "GARZI", Activ: true})
	c := cerereTest(p, "GARZI")
	previzualizareTest(t, m, &c)
	if _, err := m.db.Exec(`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES(?,'2026-11-01','CM',1,'x','x')`, id); err != nil {
		t.Fatal(err)
	}
	_, err := m.serviciu.PublicaPlanificare(ctx, c)
	verificaCod(t, err, "planificare.conflict_publicare")
	r := previzualizareTest(t, m, &c)
	if r.Blocate != 1 || r.Randuri[0].Motiv != "absenta" {
		t.Fatalf("absenta: %+v", r)
	}
}
func TestPlanificareSuprapunereSiCategorie(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	if _, err := m.db.Exec(`UPDATE tipuri_schimb SET ora_start=1320 WHERE id=?`, schimburi["G24"]); err != nil {
		t.Fatal(err)
	}
	planTest(t, m, 2026, 9, PozitiePlanificare{AngajatID: id, Data: "2026-09-30", SchimbID: schimburi["G24"], Sectiune: "GARZI", Activ: true})
	err := m.serviciu.SalveazaPlanificare(ctx, 2026, 10, []PozitiePlanificare{{AngajatID: id, Data: "2026-10-01", SchimbID: schimburi["S7"], Sectiune: "DUPA_AMIAZA", Activ: true}})
	verificaCod(t, err, "planificare.motiv_suprapunere")
	if _, err := m.db.Exec(`UPDATE angajati SET categorie_program_id='program-dimineata' WHERE id=?`, id); err != nil {
		t.Fatal(err)
	}
	err = m.serviciu.SalveazaPlanificare(ctx, 2026, 10, []PozitiePlanificare{{AngajatID: id, Data: "2026-10-05", SchimbID: schimburi["S7"], Sectiune: "DUPA_AMIAZA", Activ: true}})
	verificaCod(t, err, "planificare.motiv_program_nepermis")
}
func TestPublicareLotAtomicSiPontajEfectivProtejat(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-01", SchimbID: schimburi["S7"], Sectiune: "DUPA_AMIAZA", Activ: true}, PozitiePlanificare{AngajatID: id, Data: "2026-10-02", SchimbID: schimburi["S7"], Sectiune: "DUPA_AMIAZA", Activ: true})
	c := cerereTest(p, "DUPA_AMIAZA")
	previzualizareTest(t, m, &c)
	if _, err := m.db.Exec(`CREATE TRIGGER publicare_blocata BEFORE INSERT ON pontaj_ciorne WHEN NEW.data='2026-10-02' BEGIN SELECT RAISE(ABORT,'proba'); END`); err != nil {
		t.Fatal(err)
	}
	if _, err := m.serviciu.PublicaPlanificare(ctx, c); err == nil {
		t.Fatal("publicarea trebuia sa esueze")
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&n); err != nil || n != 0 {
		t.Fatal("lot partial", n, err)
	}
	if err := m.db.QueryRow(`SELECT count(*) FROM planificare WHERE publicat_la<>''`).Scan(&n); err != nil || n != 0 {
		t.Fatal("marcare partiala", n, err)
	}
	if _, err := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('ef',?,'2026-10-01','2026-10-01T12:00:00Z','2026-10-01T19:00:00Z','x','x')`, id); err != nil {
		t.Fatal(err)
	}
	c.Inlocuieste = true
	r := previzualizareTest(t, m, &c)
	if r.Blocate != 1 || r.Publicabile != 1 {
		t.Fatalf("efectiv: %+v", r)
	}
}
func TestPublicareUnificataToateSectiunile(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10,
		PozitiePlanificare{AngajatID: id, Data: "2026-10-01", SchimbID: schimburi["S7"], Sectiune: "DUPA_AMIAZA", Activ: true},
		PozitiePlanificare{AngajatID: id, Data: "2026-10-03", SchimbID: schimburi["G24"], Sectiune: "GARZI", Activ: true},
	)
	c := CererePublicarePlan{
		An:        2026,
		Luna:      10,
		Sectiune:  "TOATE",
		PozitiiID: []string{p.Pozitii[0].ID, p.Pozitii[1].ID},
	}
	r := previzualizareTest(t, m, &c)
	if r.Publicabile != 2 {
		t.Fatalf("previzualizare unificata: %+v", r)
	}
	if _, err := m.serviciu.PublicaPlanificare(ctx, c); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne WHERE angajat_id=? AND valoare IN ('S7','G24')`, id).Scan(&n); err != nil || n != 2 {
		t.Fatalf("publicare mixta: %d %v", n, err)
	}
}

func TestCitesteAcoperireSiEchitatePlanificare(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	planTest(t, m, 2026, 10,
		PozitiePlanificare{AngajatID: id, Data: "2026-10-03", SchimbID: schimburi["G24"], Sectiune: "GARZI", Activ: true},
		PozitiePlanificare{AngajatID: id, Data: "2026-10-04", SchimbID: schimburi["Z12"], Sectiune: "GARZI", Activ: true},
	)

	rap, err := m.serviciu.CitesteAcoperirePlanificare(ctx, 2026, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if rap.TotalZileNelucratoare == 0 || rap.TotalPosturiAcoperite == 0 {
		t.Fatalf("raport acoperire invalid: %+v", rap)
	}

	avertismente, err := m.serviciu.CitesteEchitatePlanificare(ctx, 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(avertismente) == 0 {
		t.Fatal("asteptam cel putin un avertisment consecutiv pentru 3 si 4 octombrie")
	}
}
