package datebaza

import (
	"context"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	"pontix/internal/bd"
	"pontix/internal/i18n"
)

func persoanaAbsente(t *testing.T, m *mediuTest) Angajat {
	t.Helper()
	a := angajatTest("ABS-001")
	drept := int64(23)
	a.ConcediuZileAn = &drept
	if err := m.serviciu.SalveazaAngajat(context.Background(), a); err != nil {
		t.Fatal(err)
	}
	return citeste(t, m.serviciu).Angajati[0]
}

func absenteTest(t *testing.T, m *mediuTest, an int) AbsenteAn {
	t.Helper()
	d, err := m.serviciu.CitesteAbsente(context.Background(), an)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestAbsentePersistentaEditareAnulareSiAni(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	initial := absenteTest(t, m, 2026)
	if len(initial.Tipuri) != 2 || len(initial.Zile) != 0 {
		t.Fatalf("inițiale: %+v", initial)
	}
	tip := TipAbsenta{Cod: " ev_familie ", Denumire: i18n.Text{"ro": "Eveniment familial", "en": "Family event"}, Categorie: "SPECIAL", Activ: true}
	if err := m.serviciu.SalveazaTipAbsenta(ctx, "", tip); err != nil {
		t.Fatal(err)
	}
	// Date explicite, inclusiv weekend/sărbătoare: fără regulă implicită de consum.
	lot := []ZiAbsenta{
		{AngajatID: a.ID, Data: "2026-12-25", Cod: "CO", Activ: true},
		{AngajatID: a.ID, Data: "2026-12-26", Cod: "CM", Activ: true},
		{AngajatID: a.ID, Data: "2026-12-31", Cod: "EV_FAMILIE", Activ: true},
		{AngajatID: a.ID, Data: "2027-01-01", Cod: "CO", Activ: true},
	}
	if err := m.serviciu.SalveazaAbsente(ctx, lot); err != nil {
		t.Fatal(err)
	}
	d := absenteTest(t, m, 2026)
	if len(d.Zile) != 3 || len(absenteTest(t, m, 2027).Zile) != 1 {
		t.Fatal("an sau număr de zile incorect")
	}
	if got := citeste(t, m.serviciu).Angajati[0].ConcediuZileAn; got == nil || *got != 23 {
		t.Fatal("dreptul contractual a fost rescris")
	}
	veche := d.Zile[0]
	editata := veche
	editata.Cod = "CM"
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{editata}); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{veche}), "absente.conflict")
	actuala := absenteTest(t, m, 2026).Zile[0]
	actuala.Activ = false
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{actuala}); err != nil {
		t.Fatal(err)
	}
	anulata := absenteTest(t, m, 2026).Zile[0]
	if anulata.Activ {
		t.Fatal("anulare nesalvată")
	}
	anulata.Activ, anulata.Cod = true, "CO"
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{anulata}); err != nil {
		t.Fatal(err)
	}
	if err := m.db.Close(); err != nil {
		t.Fatal(err)
	}
	var err error
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	m.serviciu = Nou(m.db)
	d = absenteTest(t, m, 2026)
	if len(d.Zile) != 3 || !d.Zile[0].Activ || d.Zile[0].Cod != "CO" {
		t.Fatalf("redeschidere: %+v", d.Zile)
	}
	var jurnal int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'absente_zile'`).Scan(&jurnal); err != nil || jurnal != 7 {
		t.Fatalf("istoric: %d %v", jurnal, err)
	}
}

func TestAbsenteTipuriProtectieIstoricSiCoduri(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	tip := TipAbsenta{Cod: "EV", Denumire: i18n.Text{"ro": "Eveniment"}, Categorie: "SPECIAL", Activ: true}
	if err := m.serviciu.SalveazaTipAbsenta(ctx, "", tip); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaTipAbsenta(ctx, "", tip), "absente.cod_duplicat")
	tip.Cod = "n12"
	verificaCod(t, m.serviciu.SalveazaTipAbsenta(ctx, "", tip), "absente.cod_duplicat")
	tip.Cod = "EV"
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "EV", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	tip.Categorie = "ODIHNA"
	verificaCod(t, m.serviciu.SalveazaTipAbsenta(ctx, "EV", tip), "absente.categorie_folosita")
	tip.Categorie, tip.Activ = "SPECIAL", false
	if err := m.serviciu.SalveazaTipAbsenta(ctx, "EV", tip); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-06", Cod: "EV", Activ: true}}), "absente.tip_invalid")
	zi := absenteTest(t, m, 2026).Zile[0]
	zi.Activ = false
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{zi}); err != nil {
		t.Fatal(err)
	}
	tip.Categorie = "MEDICAL"
	verificaCod(t, m.serviciu.SalveazaTipAbsenta(ctx, "EV", tip), "absente.categorie_folosita")
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, Schimb{Cod: "EV", Denumire: i18n.Text{"ro": "Test"}, DurataMinute: 420, Tip: "NORMAL", Activ: true}), "absente.cod_duplicat")
}

func TestAbsenteConflicteAtomiceSiPontaj(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	lot := []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true}, {AngajatID: a.ID, Data: "2026-10-06", Cod: "CO", Activ: true}}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-06", Valoare: "7"}}); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, lot), "absente.pontaj_existent")
	if len(absenteTest(t, m, 2026).Zile) != 0 {
		t.Fatal("lot salvat parțial peste schiță")
	}
	schita, err := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	c := schita.Celule[0]
	c.Valoare = ""
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{c}); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAbsente(ctx, lot); err != nil {
		t.Fatal(err)
	}
	// Protecție inclusiv când absența este salvată după încărcarea grilei.
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-07", Valoare: "7"}, {AngajatID: a.ID, Data: "2026-10-05", Valoare: "12"}}), "absente.celula_protejata")
	schita, err = m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil || len(schita.Celule) != 1 || len(schita.Absente) != 2 {
		t.Fatalf("schiță/absențe: %+v %v", schita, err)
	}
	d := absenteTest(t, m, 2026)
	conflict := []ZiAbsenta{d.Zile[0], d.Zile[1]}
	conflict[0].Activ, conflict[1].Activ = false, false
	conflict[1].Revizie = "veche"
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, conflict), "absente.conflict")
	if !absenteTest(t, m, 2026).Zile[0].Activ {
		t.Fatal("anulare parțială la conflict")
	}
	// Tura începe în ziua precedentă; ora României se schimbă în 25 octombrie.
	if _, err := m.db.Exec(`INSERT INTO pontaj (id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES ('noapte',?,'2026-10-24','2026-10-24T19:00:00Z','2026-10-25T05:00:00Z','x','x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-25", Cod: "CM", Activ: true}}), "absente.pontaj_existent")
	// Terminarea exact la miezul nopții locale nu ocupă ziua următoare.
	if _, err := m.db.Exec(`UPDATE pontaj SET sfarsit='2026-10-24T21:00:00Z' WHERE id='noapte'`); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-25", Cod: "CM", Activ: true}}); err != nil {
		t.Fatal(err)
	}
}

func TestAbsenteValidariSiRollbackJurnal(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	buna := ZiAbsenta{AngajatID: a.ID, Data: "2026-02-28", Cod: "CO", Activ: true}
	for _, lot := range [][]ZiAbsenta{nil, {buna, buna}, {{AngajatID: a.ID, Data: "2026-02-29", Cod: "CO", Activ: true}}} {
		verificaCod(t, m.serviciu.SalveazaAbsente(ctx, lot), "absente.lot_invalid")
	}
	gresita := buna
	gresita.Cod = "LIPSA"
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{gresita}), "absente.tip_invalid")
	if _, err := m.db.Exec(`CREATE TRIGGER refuza_jurnal_absente BEFORE INSERT ON jurnal_modificari WHEN NEW.tabel = 'absente_zile' BEGIN SELECT RAISE(ABORT,'proba'); END`); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{buna}), "date.eroare")
	if len(absenteTest(t, m, 2026).Zile) != 0 {
		t.Fatal("salvare fără jurnal")
	}
}

func TestMigrareAbsentePastreazaDateleAnterioare(t *testing.T) {
	m := pregateste(t)
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "veche.db"), m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range intrari {
		if r.Name() >= "0008_" {
			continue
		}
		continut, err := fs.ReadFile(bd.Migrari, r.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[r.Name()] = &fstest.MapFile{Data: continut}
	}
	ctx := context.Background()
	if _, err := bd.Migreaza(ctx, db, surse); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO coduri_absenta(cod,denumire,activ,creat_la,modificat_la) VALUES ('co','{"ro":"Personalizat"}',0,'initial','initial')`); err != nil {
		t.Fatal(err)
	}
	if _, err := bd.Migreaza(ctx, db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	d, err := Nou(db).CitesteAbsente(ctx, 2026)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Tipuri) != 2 {
		t.Fatalf("cod duplicat la migrare: %+v", d.Tipuri)
	}
	for _, tip := range d.Tipuri {
		if tip.Cod == "co" && (tip.Activ || tip.Categorie != "SPECIAL" || tip.Denumire["ro"] != "Personalizat") {
			t.Fatal("configurație suprascrisă")
		}
	}
	if err := Nou(db).SalveazaTipAbsenta(ctx, "co", TipAbsenta{Cod: "co", Denumire: i18n.Text{"ro": "Personalizat actualizat"}, Categorie: "ODIHNA", Activ: true}); err != nil {
		t.Fatal(err)
	}
}
