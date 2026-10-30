package datebaza

import (
	"context"
	"testing"
)

func salveazaCelulaPlanTest(t *testing.T, m *mediuTest, id, data, valoare string) {
	t.Helper()
	rev := ""
	s, err := m.serviciu.CitesteSchita(context.Background(), 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range s.Celule {
		if c.AngajatID == id && c.Data == data {
			rev = c.Revizie
		}
	}
	if err = m.serviciu.SalveazaSchitaCuPlanificare(context.Background(), 2026, 10, []CelulaSchita{{AngajatID: id, Data: data, Valoare: valoare, Revizie: rev}}); err != nil {
		t.Fatal(err)
	}
}
func TestPontajManualSincronizeazaProgramSectiuneSiRevenire(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "Z7")
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "s7")
	p, err := m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if err != nil || len(p.Pozitii) != 1 || p.Pozitii[0].SchimbID != coduri["S7"] || p.Pozitii[0].Sectiune != "DUPA_AMIAZA" || p.Pozitii[0].Stare != "PUBLICATA" {
		t.Fatal(p, err)
	}
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "G24")
	p, err = m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	active := 0
	for _, v := range p.Pozitii {
		if v.Activ {
			active++
			if v.SchimbID != coduri["G24"] || v.Sectiune != "GARZI" || v.Stare != "PUBLICATA" {
				t.Fatal(v)
			}
		}
	}
	if active != 1 || len(p.Pozitii) != 2 {
		t.Fatal(p)
	}
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "Z7")
	var n int
	if err = m.db.QueryRow(`SELECT count(*) FROM planificare WHERE activ=1`).Scan(&n); err != nil || n != 0 {
		t.Fatal(n, err)
	}
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "S7")
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "")
	if err = m.db.QueryRow(`SELECT count(*) FROM planificare WHERE activ=1`).Scan(&n); err != nil || n != 0 {
		t.Fatal(n, err)
	}
}
func TestPontajManualReparaCelulaVecheFaraRescriereSiFaraDubluri(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "S7"}}); err != nil {
		t.Fatal(err)
	}
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "S7")
	var n, dupa int
	m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari`).Scan(&n)
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "S7")
	m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari`).Scan(&dupa)
	if dupa != n {
		t.Fatal("repetarea schimbă auditul", n, dupa)
	}
	c := ciornaInstrumenteTest(t, m, id, "2026-10-01")
	if c.Revizie != 1 || c.Valoare != "S7" {
		t.Fatal(c)
	}
	p, err := m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if err != nil || len(p.Pozitii) != 1 || !p.Pozitii[0].Activ {
		t.Fatal(p, err)
	}
}
func TestPontajManualProtejeazaPlanIndependent(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-01", SchimbID: coduri["G24"], Sectiune: "GARZI", Activ: true})
	salveazaCelulaPlanTest(t, m, id, "2026-10-01", "Z7")
	verificaCod(t, m.serviciu.SalveazaSchitaCuPlanificare(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "S7", Revizie: "1"}}), "pontaj.plan_independent")
	c := ciornaInstrumenteTest(t, m, id, "2026-10-01")
	if c.Valoare != "Z7" || c.Revizie != 1 {
		t.Fatal(c)
	}
	p, err := m.serviciu.CitestePlanificare(ctx, 2026, 10)
	if err != nil || len(p.Pozitii) != 1 || !p.Pozitii[0].Activ || p.Pozitii[0].SchimbID != coduri["G24"] {
		t.Fatal(p, err)
	}
}
func TestPontajManualRefuzAtomicAbsentaSuprapunereCategorieSiLuna(t *testing.T) {
	for _, caz := range []struct{ nume, cod string }{{"absenta", "planificare.motiv_absenta"}, {"suprapunere", "planificare.motiv_suprapunere"}, {"categorie", "program.program_nepermis"}, {"luna", "pontaj.luna_inchisa"}, {"revizie", "pontaj.conflict"}} {
		t.Run(caz.nume, func(t *testing.T) {
			m, id, _ := mediuPlanificare(t)
			ctx := context.Background()
			celule := []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "G24"}}
			switch caz.nume {
			case "absenta":
				if _, err := m.db.Exec(`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES(?,'2026-10-02','CO',1,'x','x')`, id); err != nil {
					t.Fatal(err)
				}
			case "suprapunere":
				celule = []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "S7"}, {AngajatID: id, Data: "2026-10-02", Valoare: "N12"}}
				if _, err := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('efectiv',?,'2026-10-02','2026-10-02T19:00:00Z','2026-10-03T04:00:00Z','x','x')`, id); err != nil {
					t.Fatal(err)
				}
			case "categorie":
				if _, err := m.db.Exec(`UPDATE angajati SET categorie_program_id='program-dimineata' WHERE id=?`, id); err != nil {
					t.Fatal(err)
				}
				celule[0].Valoare = "S7"
			case "luna":
				if _, err := m.db.Exec(`INSERT INTO stare_luni(an,luna,stare,modificat_la) VALUES(2026,10,'INCHISA','x')`); err != nil {
					t.Fatal(err)
				}
			case "revizie":
				celule[0].Revizie = "99"
			}

			verificaCod(t, m.serviciu.SalveazaSchitaCuPlanificare(ctx, 2026, 10, celule), caz.cod)
			var n int
			if err := m.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&n); err != nil || n != 0 {
				t.Fatal("salvare parțială", n, err)
			}
		})
	}
}
func TestPontajManualRollbackDacaPlanificareaNuPoateFiScrisa(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	if _, err := m.db.Exec(`CREATE TRIGGER eroare_plan BEFORE INSERT ON planificare WHEN NEW.data='2026-10-02' BEGIN SELECT RAISE(ABORT,'eroare simulata'); END`); err != nil {
		t.Fatal(err)
	}
	err := m.serviciu.SalveazaSchitaCuPlanificare(context.Background(), 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "S7"}, {AngajatID: id, Data: "2026-10-02", Valoare: "S7"}})
	if err == nil {
		t.Fatal("lipsește eroarea simulată")
	}
	for _, tabel := range []string{"planificare", "pontaj_ciorne"} {
		var n int
		if err := m.db.QueryRow("SELECT count(*) FROM " + tabel).Scan(&n); err != nil || n != 0 {
			t.Fatal(tabel, n, err)
		}
	}
}
