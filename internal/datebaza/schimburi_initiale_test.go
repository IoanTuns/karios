package datebaza

import (
	"context"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	"pontix/internal/bd"
)

func TestSchimburiInitialeConformSurselor(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	lista, err := m.serviciu.ListaSchimburi(ctx)
	if err != nil || len(lista) != 11 {
		t.Fatalf("set inițial: %+v, %v", lista, err)
	}
	var sursa []Schimb
	if err := json.Unmarshal(schimburiInitiale, &sursa); err != nil {
		t.Fatal(err)
	}
	for _, s := range sursa {
		gasit := false
		for _, curent := range lista {
			if curent.Cod == s.Cod {
				gasit = true
				if curent.OraStart != s.OraStart || curent.DurataMinute != s.DurataMinute || curent.PauzaMinute != 0 || curent.Tip != s.Tip || curent.Regim != s.Regim || curent.Categorie != s.Categorie || !curent.Activ {
					t.Fatalf("migrarea diferă de sursă: %+v", curent)
				}
			}
		}
		if !gasit {
			t.Fatalf("lipsește codul %s", s.Cod)
		}
	}
	for _, s := range lista {
		if s.Cod == "7" && (s.OraStart != 420 || s.DurataMinute != 420) {
			t.Fatal("intervalul 7 diferă de cerința utilizatorului")
		}
		if s.Cod == "12" && (s.OraStart != 840 || s.DurataMinute != 720 || (s.OraStart+s.DurataMinute)%1440 != 120) {
			t.Fatal("intervalul 12 nu este 14–02")
		}
	}
	if n, err := m.serviciu.CompleteazaSchimburiInitiale(ctx); err != nil || n != 0 {
		t.Fatalf("completare repetată: %d, %v", n, err)
	}
	for _, s := range lista {
		if s.Cod == "N12" {
			s.OraStart = 1200
			s.PauzaMinute = 45
			s.Activ = false
			if err := m.serviciu.SalveazaSchimb(ctx, s); err != nil {
				t.Fatal(err)
			}
		}
	}
	if n, err := m.serviciu.CompleteazaSchimburiInitiale(ctx); err != nil || n != 0 {
		t.Fatalf("completare peste configurat: %d, %v", n, err)
	}
	lista, err = m.serviciu.ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range lista {
		if s.Cod == "N12" && (s.OraStart != 1200 || s.PauzaMinute != 45 || s.Activ) {
			t.Fatal("configurație suprascrisă")
		}
	}
}

func TestCompletareDupaRedenumireaUnuiSchimbInitial(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	lista, err := m.serviciu.ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var idRedenumit string
	for _, s := range lista {
		if s.Cod == "Z7" {
			idRedenumit = s.ID
			s.Cod = "Z7_CUSTOM"
			s.OraStart = 300
			if err := m.serviciu.SalveazaSchimb(ctx, s); err != nil {
				t.Fatal(err)
			}
		}
	}
	if n, err := m.serviciu.CompleteazaSchimburiInitiale(ctx); err != nil || n != 1 {
		t.Fatalf("completare după redenumire: %d, %v", n, err)
	}
	var pastrat, completat int
	if err := m.db.QueryRow(`SELECT count(*) FROM tipuri_schimb WHERE id = ? AND cod = 'Z7_CUSTOM' AND ora_start = 300`, idRedenumit).Scan(&pastrat); err != nil || pastrat != 1 {
		t.Fatalf("șablonul redenumit a fost modificat: %d, %v", pastrat, err)
	}
	if err := m.db.QueryRow(`SELECT count(*) FROM tipuri_schimb WHERE id != ? AND cod = 'Z7' AND ora_start = 420`, idRedenumit).Scan(&completat); err != nil || completat != 1 {
		t.Fatalf("cod lipsă nerecreat separat: %d, %v", completat, err)
	}
	if n, err := m.serviciu.CompleteazaSchimburiInitiale(ctx); err != nil || n != 0 {
		t.Fatalf("completare repetată: %d, %v", n, err)
	}
}

func TestMigrareInitialePastreazaCoduriExistente(t *testing.T) {
	ctx := context.Background()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, intrare := range intrari {
		if intrare.Name() >= "0007_" {
			continue
		}
		continut, err := fs.ReadFile(bd.Migrari, intrare.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[intrare.Name()] = &fstest.MapFile{Data: continut}
	}
	cheie := make([]byte, bd.LungimeCheie)
	for i := range cheie {
		cheie[i] = byte(i + 1)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "configurata.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(ctx, db, surse); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
		VALUES ('personalizat', 'n12', '{"ro":"Personalizat"}', 1230, 600, 30, 'TURA', 0, 'initial', 'initial')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la, sters_la)
		VALUES ('sters', 'G24', '{"ro":"Șters"}', 420, 1440, 0, 'GARDA', 1, 'initial', 'initial', 'initial')`); err != nil {
		t.Fatal(err)
	}
	migrare, err := fs.ReadFile(bd.Migrari, "0007_schimburi_initiale.sql")
	if err != nil {
		t.Fatal(err)
	}
	surse["0007_schimburi_initiale.sql"] = &fstest.MapFile{Data: migrare}
	if n, err := bd.Migreaza(ctx, db, surse); err != nil || n != 1 {
		t.Fatalf("migrare: %d, %v", n, err)
	}
	var cod, modificat string
	var start, pauza, activ int
	if err := db.QueryRow(`SELECT cod, ora_start, pauza_minute, activ, modificat_la FROM tipuri_schimb WHERE id = 'personalizat'`).Scan(&cod, &start, &pauza, &activ, &modificat); err != nil || cod != "n12" || start != 1230 || pauza != 30 || activ != 0 || modificat != "initial" {
		t.Fatalf("configurație veche modificată: %s, %d, %d, %d, %s, %v", cod, start, pauza, activ, modificat, err)
	}
	if _, err := bd.Migreaza(ctx, db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	if n, err := Nou(db).CompleteazaSchimburiInitiale(ctx); err != nil || n != 0 {
		t.Fatalf("reactivare necerută: %d, %v", n, err)
	}
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM tipuri_schimb WHERE upper(cod) = 'G24' AND sters_la IS NOT NULL`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("cod șters recreat: %d, %v", n, err)
	}
}

func TestCompletareInitialeEsteAtomica(t *testing.T) {
	m := pregateste(t)
	// Proba pornește fără Z7/S7; eliminăm întâi referințele din categoriile inițiale.
	if _, err := m.db.Exec(`UPDATE categorii_program SET schimb_implicit_id=NULL;
	DELETE FROM categorii_program_schimburi WHERE schimb_id IN (SELECT id FROM tipuri_schimb WHERE cod IN ('Z7','S7'))`); err != nil {
		t.Fatal(err)
	}
	if _, err := m.db.Exec(`DELETE FROM tipuri_schimb WHERE cod IN ('Z7', 'S7')`); err != nil {
		t.Fatal(err)
	}
	if _, err := m.db.Exec(`CREATE TRIGGER initiale_blocate BEFORE INSERT ON tipuri_schimb WHEN NEW.cod = 'S7' BEGIN SELECT RAISE(ABORT, 'proba'); END`); err != nil {
		t.Fatal(err)
	}
	n, err := m.serviciu.CompleteazaSchimburiInitiale(context.Background())
	verificaCod(t, err, "date.eroare")
	if n != 0 {
		t.Fatal("număr parțial raportat la eșec")
	}
	var ramase int
	if err := m.db.QueryRow(`SELECT count(*) FROM tipuri_schimb WHERE cod = 'Z7'`).Scan(&ramase); err != nil || ramase != 0 {
		t.Fatalf("inserare parțială: %d, %v", ramase, err)
	}
}
