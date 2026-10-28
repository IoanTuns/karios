package audit_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/i18n"
)

var ctx = context.Background()

// baza creează o bază migrată și sigilată și întoarce calea și cheia ei.
func baza(t *testing.T) (string, []byte) {
	t.Helper()
	cale := filepath.Join(t.TempDir(), "pontix.db")
	cheie := make([]byte, bd.LungimeCheie)
	for i := range cheie {
		cheie[i] = byte(i * 7)
	}
	db, err := bd.Deschide(cale, cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(ctx, db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	if err := audit.Initializeaza(ctx, db); err != nil {
		t.Fatal(err)
	}
	return cale, cheie
}

// lucru deschide baza ca aplicația: cu trigger-e și pază.
func lucru(t *testing.T, cale string, cheie []byte) (*sql.DB, *audit.Paza) {
	t.Helper()
	brut := extern(t, cale, cheie)
	rez, err := audit.Verifica(ctx, brut)
	if err != nil {
		t.Fatal(err)
	}
	brut.Close()
	paza := audit.NouaPaza(rez)
	db, err := bd.DeschideCuAudit(cale, cheie, paza)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db, paza
}

// extern deschide baza fără trigger-e, ca un alt program care are cheia de date.
func extern(t *testing.T, cale string, cheie []byte) *sql.DB {
	t.Helper()
	db, err := bd.Deschide(cale, cheie)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func exec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.ExecContext(ctx, q, args...); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
}

func verifica(t *testing.T, db *sql.DB) audit.Rezultat {
	t.Helper()
	rez, err := audit.Verifica(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	return rez
}

func scrieLucru(t *testing.T, db *sql.DB) {
	t.Helper()
	exec(t, db, `INSERT INTO sectii (id, denumire, creat_la, modificat_la) VALUES ('s1', '{"ro":"ATI"}', 'x', 'x')`)
	exec(t, db, `UPDATE sectii SET denumire = '{"ro":"ATI 2"}' WHERE id = 's1'`)
	exec(t, db, `INSERT INTO jurnal_modificari (tabel, rand_id, valoare_noua, cand, cine) VALUES ('sectii', 's1', '{}', 'x', 'local')`)
	exec(t, db, `INSERT INTO jurnal_modificari (tabel, rand_id, valoare_noua, cand, cine) VALUES ('sectii', 's1', '{"a":1}', 'y', 'local')`)
	exec(t, db, `INSERT INTO sectii (id, denumire, creat_la, modificat_la) VALUES ('s2', '{"ro":"Chirurgie"}', 'x', 'x')`)
	exec(t, db, `DELETE FROM sectii WHERE id = 's2'`)
}

func TestScrierileAplicatieiPastreazaSigiliile(t *testing.T) {
	cale, cheie := baza(t)
	if rez := verifica(t, extern(t, cale, cheie)); rez.Stare != audit.StareOK {
		t.Fatalf("după inițializare: %+v", rez)
	}
	db, paza := lucru(t, cale, cheie)
	if paza.Blocata() {
		t.Fatal("paza blocată fără incident")
	}
	scrieLucru(t, db)
	if rez := verifica(t, db); rez.Stare != audit.StareOK {
		t.Fatalf("după scrierile aplicației: %+v", rez)
	}
	var faraMac int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM jurnal_modificari WHERE mac IS NULL`).Scan(&faraMac); err != nil || faraMac != 0 {
		t.Fatalf("rânduri de jurnal fără mac: %d, %v", faraMac, err)
	}
}

func TestModificareaExternaEsteDetectata(t *testing.T) {
	cazuri := []struct {
		nume    string
		sql     string
		tabel   string
		randRup bool
	}{
		{"modificare date", `UPDATE sectii SET denumire = '{"ro":"Altceva"}' WHERE id = 's1'`, "sectii", false},
		{"inserare date", `INSERT INTO sectii (id, denumire, creat_la, modificat_la) VALUES ('s9', '{"ro":"X"}', 'x', 'x')`, "sectii", false},
		{"ștergere date", `DELETE FROM sectii WHERE id = 's1'`, "sectii", false},
		{"modificare jurnal", `UPDATE jurnal_modificari SET valoare_noua = '{"fals":1}' WHERE id = (SELECT min(id) FROM jurnal_modificari WHERE tabel = 'sectii')`, "jurnal_modificari", true},
		{"ștergerea ultimului rând din jurnal", `DELETE FROM jurnal_modificari WHERE id = (SELECT max(id) FROM jurnal_modificari)`, "jurnal_modificari", false},
		{"trigger plantat", `CREATE TRIGGER t AFTER INSERT ON sectii BEGIN SELECT 1; END`, "sqlite_schema", false},
		{"sigiliu șters", `DELETE FROM sigilii WHERE tabel = 'sectii'`, "sectii", false},
	}
	for _, c := range cazuri {
		t.Run(c.nume, func(t *testing.T) {
			cale, cheie := baza(t)
			db, _ := lucru(t, cale, cheie)
			scrieLucru(t, db)
			db.Close()
			exec(t, extern(t, cale, cheie), c.sql)
			rez := verifica(t, extern(t, cale, cheie))
			if rez.Stare != audit.StareIncident || len(rez.Tabele) != 1 || rez.Tabele[0] != c.tabel {
				t.Fatalf("rezultat: %+v, așteptat incident pe %s", rez, c.tabel)
			}
			if (rez.RandRupt > 0) != c.randRup {
				t.Fatalf("rând rupt: %d", rez.RandRupt)
			}
		})
	}
}

func TestSigiliiGoaleSuntIncident(t *testing.T) {
	cale, cheie := baza(t)
	exec(t, extern(t, cale, cheie), `DELETE FROM sigilii`)
	if rez := verifica(t, extern(t, cale, cheie)); rez.Stare != audit.StareIncident {
		t.Fatalf("rezultat: %+v", rez)
	}
}

func TestFunctiileNuMergDinTriggerPermanent(t *testing.T) {
	cale, cheie := baza(t)
	db := extern(t, cale, cheie)
	// Un trigger salvat în fișier care ar „repara” sigiliul la următoarea scriere a aplicației.
	exec(t, db, `CREATE TRIGGER repara AFTER INSERT ON setari BEGIN UPDATE sigilii SET valoare = pontix_h('x', 1) WHERE tabel = 'sectii'; END`)
	if _, err := db.ExecContext(ctx, `INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('k', 'v', 'x')`); err == nil {
		t.Fatal("funcția auditului a rulat dintr-un trigger permanent")
	}
}

func TestIncidentBlocheazaSiConfirmareaDeblocheaza(t *testing.T) {
	cale, cheie := baza(t)
	exec(t, extern(t, cale, cheie), `INSERT INTO sectii (id, denumire, creat_la, modificat_la) VALUES ('s9', '{"ro":"X"}', 'x', 'x')`)
	db, paza := lucru(t, cale, cheie)
	if !paza.Blocata() {
		t.Fatal("incidentul nu blochează")
	}
	_, err := db.ExecContext(ctx, `UPDATE sectii SET activ = 0 WHERE id = 's9'`)
	if !audit.EsteEditareBlocata(err) {
		t.Fatalf("scrierea nu e blocată: %v", err)
	}
	// Citirea merge.
	var n int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM sectii`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("citire: %d, %v", n, err)
	}
	var e *i18n.Eroare
	if _, err := audit.Confirma(ctx, db, paza, "  "); !errors.As(err, &e) || e.Cod != "audit.motiv_obligatoriu" {
		t.Fatalf("confirmare fără motiv: %v", err)
	}
	rez, err := audit.Confirma(ctx, db, paza, "Corectură făcută de IT, tichet 12")
	if err != nil || rez.Stare != audit.StareOK || paza.Blocata() {
		t.Fatalf("confirmare: %+v, %v", rez, err)
	}
	exec(t, db, `UPDATE sectii SET activ = 0 WHERE id = 's9'`)
	var incident string
	if err := db.QueryRowContext(ctx, `SELECT valoare_veche FROM jurnal_modificari WHERE tabel = 'audit' AND rand_id = 'incident'`).Scan(&incident); err != nil || incident == "" {
		t.Fatalf("incidentul nu e în jurnal: %q, %v", incident, err)
	}
	if rez := verifica(t, db); rez.Stare != audit.StareOK {
		t.Fatalf("după confirmare și scriere: %+v", rez)
	}
}

func TestConfirmareaReparaLantulRupt(t *testing.T) {
	cale, cheie := baza(t)
	db, _ := lucru(t, cale, cheie)
	scrieLucru(t, db)
	db.Close()
	exec(t, extern(t, cale, cheie), `UPDATE jurnal_modificari SET cine = 'altcineva' WHERE id = (SELECT min(id) FROM jurnal_modificari)`)
	db, paza := lucru(t, cale, cheie)
	if paza.Rezultat().RandRupt == 0 {
		t.Fatalf("lanțul rupt nu e detectat: %+v", paza.Rezultat())
	}
	if rez, err := audit.Confirma(ctx, db, paza, "verificat"); err != nil || rez.Stare != audit.StareOK {
		t.Fatalf("confirmare: %+v, %v", rez, err)
	}
}

func TestAmprentaLunii(t *testing.T) {
	cale, cheie := baza(t)
	db, _ := lucru(t, cale, cheie)
	exec(t, db, `INSERT INTO angajati (id, nume, marca, sectie_id, creat_la, modificat_la) VALUES ('a', 'Pop', 'M1', NULL, 'x', 'x')`)
	exec(t, db, `INSERT INTO pontaj_ciorne (angajat_id, data, valoare, modificat_la) VALUES ('a', '2026-10-05', 'D8', 'x'), ('a', '2026-11-02', 'D8', 'x')`)
	amprenta := func() string {
		t.Helper()
		a, err := audit.AmprentaLuna(ctx, db, 2026, 10)
		if err != nil {
			t.Fatal(err)
		}
		return a
	}
	a1 := amprenta()
	if len(a1) != 19 || a1[4] != '-' {
		t.Fatalf("format: %q", a1)
	}
	if a2 := amprenta(); a2 != a1 {
		t.Fatalf("amprenta nu e deterministă: %s / %s", a1, a2)
	}
	exec(t, db, `UPDATE pontaj_ciorne SET valoare = 'N12' WHERE data = '2026-11-02'`)
	if a := amprenta(); a != a1 {
		t.Fatal("o altă lună schimbă amprenta")
	}
	exec(t, db, `UPDATE pontaj_ciorne SET valoare = 'N12' WHERE data = '2026-10-05'`)
	if a := amprenta(); a == a1 {
		t.Fatal("modificarea lunii nu schimbă amprenta")
	}
	if rez := verifica(t, db); rez.Stare != audit.StareOK {
		t.Fatalf("sigilii: %+v", rez)
	}
}
