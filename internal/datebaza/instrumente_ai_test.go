package datebaza

import (
	"context"
	"pontix/internal/ai"
	"pontix/internal/bd"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"strings"
	"testing"
)

func aplicaInstrumenteTest(t *testing.T, m *mediuTest, cod string) PlanAI {
	t.Helper()
	ctx := CuAsistentAI(context.Background())
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti"}, DeLa: "2026-10-30", PanaLa: "2026-11-02", CodSchimb: cod, FiltruZile: "lucratoare", ModExistente: "inlocuieste"}
	p, err := m.serviciu.PregatestePlanAI(ctx, i)
	if err != nil {
		t.Fatal(err)
	}
	p, err = m.serviciu.AplicaPlanAI(ctx, i, p.Token)
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func ciornaInstrumenteTest(t *testing.T, m *mediuTest, id, data string) acces.CitesteCiornaRow {
	t.Helper()
	r, err := acces.New(m.db).CitesteCiorna(context.Background(), acces.CitesteCiornaParams{AngajatID: id, Data: data})
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func TestAnulareAIPersistentaDouaLuniSiAudit(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := CuAsistentAI(context.Background())
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-30", Valoare: "6"}}); err != nil {
		t.Fatal(err)
	}
	p := aplicaInstrumenteTest(t, m, "G24")
	if p.OperatieID == "" || p.Aplicate != 2 {
		t.Fatalf("plan %+v", p)
	}
	// Redeschiderea dovedește că anularea nu depinde de sesiunea din frontend.
	m.db.Close()
	var err error
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	m.serviciu = Nou(m.db)
	u, err := m.serviciu.PrevizualizeazaAnulareAI(ctx)
	if err != nil || u.Conflicte != 0 || len(u.Randuri) != 2 {
		t.Fatalf("anulare %+v %v", u, err)
	}
	verificaCod(t, func() error { _, e := m.serviciu.AnuleazaUltimaOperatieAI(ctx, "fals"); return e }(), "ai.anulare_conflict")
	u, err = m.serviciu.AnuleazaUltimaOperatieAI(ctx, u.Token)
	if err != nil || !u.Anulata {
		t.Fatal(u, err)
	}
	if c := ciornaInstrumenteTest(t, m, id, "2026-10-30"); c.Valoare != "6" || c.Revizie != 3 {
		t.Fatal(c)
	}
	if c := ciornaInstrumenteTest(t, m, id, "2026-11-02"); c.Valoare != "" || c.Revizie != 2 {
		t.Fatal(c)
	}
	var active, n int
	if err := m.db.QueryRow(`SELECT count(*) FROM planificare WHERE activ=1`).Scan(&active); err != nil || active != 0 {
		t.Fatal(active, err)
	}
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel='operatii_ai' AND cine LIKE '%asistent_ai%'`).Scan(&n); err != nil || n != 2 {
		t.Fatal(n, err)
	}
	p, err = m.serviciu.PlanSalvatAI(ctx, p.OperatieID)
	if err != nil || !p.Anulat {
		t.Fatal(p, err)
	}
	verificaCod(t, func() error { _, e := m.serviciu.AnuleazaUltimaOperatieAI(ctx, u.Token); return e }(), "ai.anulare_conflict")
}
func TestAnulareAIConflicteAtomice(t *testing.T) {
	for _, caz := range []string{"celula", "aceeasi_valoare", "plan", "publicare", "luna", "pontaj_efectiv", "absenta"} {
		t.Run(caz, func(t *testing.T) {
			m, id, _ := mediuPlanificare(t)
			ctx := context.Background()
			if caz == "absenta" {
				if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-30", Valoare: "6"}}); err != nil {
					t.Fatal(err)
				}
			}
			p := aplicaInstrumenteTest(t, m, "G24")
			u, err := m.serviciu.PrevizualizeazaAnulareAI(ctx)
			if err != nil {
				t.Fatal(err)
			}
			switch caz {
			case "celula", "aceeasi_valoare":
				v := "4"
				if caz == "aceeasi_valoare" {
					v = "G24"
				}
				_, err = m.db.Exec(`UPDATE pontaj_ciorne SET valoare=?,revizie=revizie+1 WHERE angajat_id=? AND data='2026-10-30'`, v, id)
			case "plan":
				_, err = m.db.Exec(`UPDATE planificare SET revizie=revizie+1 WHERE angajat_id=? AND data='2026-10-30'`, id)
			case "publicare":
				_, err = m.db.Exec(`UPDATE planificare SET publicat_la='nou' WHERE angajat_id=? AND data='2026-10-30'`, id)
			case "luna":
				_, err = m.db.Exec(`INSERT INTO stare_luni(an,luna,stare,modificat_la) VALUES(2026,11,'INCHISA','x')`)
			case "pontaj_efectiv":
				_, err = m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('efectiv',?,'2026-10-30','2026-10-30T07:00:00Z','2026-10-30T14:00:00Z','x','x')`, id)
			case "absenta":
				_, err = m.db.Exec(`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES(?,'2026-10-30','CO',1,'x','x')`, id)
			}
			if err != nil {
				t.Fatal(err)
			}
			verificaCod(t, func() error { _, e := m.serviciu.AnuleazaUltimaOperatieAI(ctx, u.Token); return e }(), "ai.anulare_conflict")
			actual, err := m.serviciu.PrevizualizeazaAnulareAI(ctx)
			if err != nil || actual.Conflicte == 0 {
				t.Fatal(actual, err)
			}
			verificaCod(t, func() error { _, e := m.serviciu.AnuleazaUltimaOperatieAI(ctx, actual.Token); return e }(), "ai.anulare_conflict")
			if c := ciornaInstrumenteTest(t, m, id, "2026-11-02"); c.Valoare != "G24" || c.Revizie != 1 {
				t.Fatal("anulare parțială", c)
			}
			salvat, err := m.serviciu.PlanSalvatAI(ctx, p.OperatieID)
			if err != nil || salvat.Anulat {
				t.Fatal(salvat, err)
			}
		})
	}
}
func TestAnulareAIRestaureazaPlanAnterior(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-30", SchimbID: coduri["G24"], Sectiune: "GARZI", Activ: true})
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-30", Valoare: "G24"}}); err != nil {
		t.Fatal(err)
	}
	aplicaInstrumenteTest(t, m, "N12")
	u, err := m.serviciu.PrevizualizeazaAnulareAI(ctx)
	if err != nil || u.Conflicte != 0 {
		t.Fatal(u, err)
	}
	if len(u.Randuri[0].PlanAcum) != 1 || u.Randuri[0].PlanAcum[0] != "N12" || len(u.Randuri[0].PlanRestaurat) != 1 || u.Randuri[0].PlanRestaurat[0] != "G24" {
		t.Fatal("planificarea restaurată nu este vizibilă", u)
	}
	if _, err = m.serviciu.AnuleazaUltimaOperatieAI(ctx, u.Token); err != nil {
		t.Fatal(err)
	}
	var schimb string
	var rev int
	if err = m.db.QueryRow(`SELECT schimb_id,revizie FROM planificare WHERE angajat_id=? AND data='2026-10-30' AND activ=1`, id).Scan(&schimb, &rev); err != nil || schimb != coduri["G24"] || rev != 3 {
		t.Fatal(schimb, rev, err)
	}
}
func TestAnulareAINuAcceptaTokenDeLaAltaOperatie(t *testing.T) {
	m, _, _ := mediuPlanificare(t)
	ctx := context.Background()
	aplicaInstrumenteTest(t, m, "G24")
	u, err := m.serviciu.PrevizualizeazaAnulareAI(ctx)
	if err != nil {
		t.Fatal(err)
	}
	aplicaInstrumenteTest(t, m, "N12")
	verificaCod(t, func() error { _, e := m.serviciu.AnuleazaUltimaOperatieAI(ctx, u.Token); return e }(), "ai.anulare_conflict")
}

func TestAnulareAIRollbackLaEroareInALunaDoua(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := context.Background()
	p := aplicaInstrumenteTest(t, m, "G24")
	u, err := m.serviciu.PrevizualizeazaAnulareAI(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = m.db.Exec(`CREATE TRIGGER proba_eroare_anulare BEFORE UPDATE ON pontaj_ciorne WHEN NEW.data='2026-11-02' BEGIN SELECT RAISE(ABORT,'eroare simulata'); END`); err != nil {
		t.Fatal(err)
	}
	if _, err = m.serviciu.AnuleazaUltimaOperatieAI(ctx, u.Token); err == nil {
		t.Fatal("lipsește eroarea simulată")
	}
	for _, data := range []string{"2026-10-30", "2026-11-02"} {
		if c := ciornaInstrumenteTest(t, m, id, data); c.Valoare != "G24" || c.Revizie != 1 {
			t.Fatal("restaurare parțială", c)
		}
	}
	var active int
	if err = m.db.QueryRow(`SELECT count(*) FROM planificare WHERE activ=1`).Scan(&active); err != nil || active != 2 {
		t.Fatal(active, err)
	}
	salvat, err := m.serviciu.PlanSalvatAI(ctx, p.OperatieID)
	if err != nil || salvat.Anulat {
		t.Fatal(salvat, err)
	}
}
func TestSabloaneAIValidareSiPersistenta(t *testing.T) {
	m, _, _ := mediuPlanificare(t)
	ctx := context.Background()
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti"}, MesajOriginal: "nu se reține"}
	for _, n := range []string{" ", strings.Repeat("ă", 101)} {
		verificaCod(t, m.serviciu.SalveazaSablonAI(ctx, n, i), "ai.sablon_nume_invalid")
	}
	if err := m.serviciu.SalveazaSablonAI(ctx, " Secția dimineața ", i); err != nil {
		t.Fatal(err)
	}
	m.db.Close()
	var err error
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	m.serviciu = Nou(m.db)
	lista, err := Nou(m.db).ListaSabloaneAI(ctx)
	if err != nil || len(lista) != 1 || lista[0].Nume != "Secția dimineața" || lista[0].Intentie.MesajOriginal != "" {
		t.Fatal(lista, err)
	}
	if err = m.serviciu.StergeSablonAI(ctx, lista[0].ID); err != nil {
		t.Fatal(err)
	}
	lista, err = m.serviciu.ListaSabloaneAI(ctx)
	if err != nil || len(lista) != 0 {
		t.Fatal(lista, err)
	}
}
func TestSelectieExplicataSiSumarAI(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := context.Background()
	a := angajatTest("EXCLUS")
	a.Activ = false
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-30", Valoare: "6"}}); err != nil {
		t.Fatal(err)
	}
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti"}, DeLa: "2026-10-30", PanaLa: "2026-11-02", CodSchimb: "G24", FiltruZile: "lucratoare", ModExistente: "pastreaza"}
	p, err := m.serviciu.PregatestePlanAI(ctx, i)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.SelectieExplicata) != 2 || p.SelectieExplicata[0].Inclus || p.SelectieExplicata[0].Motiv != "inactiv" {
		t.Fatal(p.SelectieExplicata)
	}
	if len(p.SumarLuni) != 2 || p.SumarLuni[0].Selectati != 1 || p.SumarLuni[0].Afectati != 0 || p.SumarLuni[0].Pastrate != 1 || p.SumarLuni[1].Afectati != 1 || p.SumarPersoane[0].Noi != 1 {
		t.Fatal(p.SumarLuni, p.SumarPersoane)
	}
	i.DeLa = "2026-10-31"
	i.PanaLa = i.DeLa
	p, err = m.serviciu.PregatestePlanAI(ctx, i)
	if err != nil || p.SumarLuni[0].Selectati != 1 || p.SumarLuni[0].Afectati != 0 || p.SumarLuni[0].Zile != 0 {
		t.Fatal("selecția persoanelor nu dispare când nu sunt zile eligibile", p, err)
	}
}
func TestVerificarePontajAIFaraScrieri(t *testing.T) {
	m, id, coduri := mediuPlanificare(t)
	ctx := context.Background()
	for _, v := range []struct{ data, valoare string }{{"2026-10-01", "NECUNOSCUT"}, {"2026-10-02", "Z7"}, {"2026-10-05", "6"}, {"2026-10-06", "N12"}, {"2026-10-07", "Z7"}} {
		if _, err := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,?,?,'x')`, id, v.data, v.valoare); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := m.db.Exec(`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES(?,'2026-10-02','CO',1,'x','x')`, id); err != nil {
		t.Fatal(err)
	}
	planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-07", SchimbID: coduri["S7"], Sectiune: "DUPA_AMIAZA", Activ: true})
	// N12 se încheie exact la începutul Z7, deci nu este suprapunere.
	// Adăugăm un interval efectiv care se intersectează realmente cu Z7.
	if _, err := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('suprapus',?,'2026-10-07','2026-10-07T08:00:00Z','2026-10-07T12:00:00Z','x','x')`, id); err != nil {
		t.Fatal(err)
	}
	var inainte, dupa int
	m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari`).Scan(&inainte)
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti"}, DeLa: "2026-10-01", PanaLa: "2026-10-09", FiltruZile: "lucratoare"}
	r, err := m.serviciu.VerificaPontajAI(ctx, i)
	if err != nil {
		t.Fatal(err)
	}
	m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari`).Scan(&dupa)
	if dupa != inainte || r.Celule != 7 {
		t.Fatal(r, inainte, dupa)
	}
	motive := map[string]bool{}
	for _, s := range r.Semnale {
		motive[s.Motiv] = true
		if s.Data == "2026-10-06" && s.Motiv == "suprapunere" {
			t.Fatal("turele adiacente nu se suprapun", s)
		}
	}
	for _, motiv := range []string{"cod_invalid", "absenta", "ore_manuale", "suprapunere", "plan_diferit", "gol"} {
		if !motive[motiv] {
			t.Errorf("lipsește %s: %+v", motiv, r)
		}
	}
}

func TestVerificareAIReguliPeLunaCompleta(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaRegulaConfigurabila(ctx, calcul.RegulaConfigurabila{Cod: "LIMITA_TEST", Denumire: "Limită fictivă", TipRegula: "PLAFON_LUNAR", Severitate: "AVERTISMENT", CoduriSchimb: "Z7", ValoareNumerica: 1, Activ: true}); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "Z7"}, {AngajatID: id, Data: "2026-10-02", Valoare: "Z7"}}); err != nil {
		t.Fatal(err)
	}
	r, err := m.serviciu.VerificaPontajAI(ctx, ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti"}, DeLa: "2026-10-30", PanaLa: "2026-10-30", FiltruZile: "lucratoare"})
	if err != nil {
		t.Fatal(err)
	}
	gasit := false
	for _, s := range r.Semnale {
		if s.RegulaCod == "LIMITA_TEST" {
			gasit = s.Motiv == "regula" && s.Data == "2026-10" && s.Severitate == "AVERTISMENT"
		}
	}
	if !gasit {
		t.Fatal("plafonul lunar a fost restrâns la ziua verificată", r)
	}
}
