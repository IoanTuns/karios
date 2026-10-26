package bd

import (
	"context"
	"database/sql"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func bazaContracte(t *testing.T) *sql.DB {
	t.Helper()
	db, err := Deschide(filepath.Join(t.TempDir(), "contracte.db"), cheieNoua(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := Migreaza(context.Background(), db, Migrari); err != nil {
		t.Fatal(err)
	}
	return db
}

func executa(t *testing.T, db *sql.DB, comanda string, argumente ...any) {
	t.Helper()
	if _, err := db.Exec(comanda, argumente...); err != nil {
		t.Fatal(err)
	}
}

func verificaNorma(t *testing.T, db *sql.DB, id string, asteptata int) {
	t.Helper()
	var minute int
	if err := db.QueryRow(`SELECT norma_efectiva_minute FROM angajati_norme WHERE id = ?`, id).Scan(&minute); err != nil {
		t.Fatal(err)
	}
	if minute != asteptata {
		t.Errorf("norma angajatului %s: %d, așteptat %d minute", id, minute, asteptata)
	}
}

func TestContracteSiNormaGlobala(t *testing.T) {
	db := bazaContracte(t)
	executa(t, db, `INSERT INTO angajati (id, nume, prenume, marca, creat_la, modificat_la)
		VALUES ('intreg', 'Exemplu', 'Ana', '001', 'x', 'x')`)
	executa(t, db, `INSERT INTO angajati
		(id, nume, prenume, marca, tip_contract, norma_partiala_minute, concediu_zile_an, creat_la, modificat_la)
		VALUES ('partial', 'Exemplu', 'Ion', '002', 'PARTIALA', 225, 24, 'x', 'x')`)

	var nume, prenume, marca, tip string
	var concediu, normaPartiala, normaVeche sql.NullInt64
	if err := db.QueryRow(`SELECT nume, prenume, marca, tip_contract, concediu_zile_an, norma_partiala_minute, norma_minute_anterioara
		FROM angajati WHERE id = 'intreg'`).Scan(&nume, &prenume, &marca, &tip, &concediu, &normaPartiala, &normaVeche); err != nil {
		t.Fatal(err)
	}
	if nume != "Exemplu" || prenume != "Ana" || marca != "001" || tip != "INTREAGA" || concediu.Valid || normaPartiala.Valid || normaVeche.Valid {
		t.Fatalf("datele și valorile implicite: %q %q, marca %q, tip %q, concediu %v, normă parțială %v",
			nume, prenume, marca, tip, concediu, normaPartiala)
	}
	verificaNorma(t, db, "intreg", 420)
	verificaNorma(t, db, "partial", 225)

	// Modificarea globală se aplică și angajaților deja existenți, fără a
	// schimba norma individuală a contractelor parțiale sau concediul.
	executa(t, db, `UPDATE setari SET valoare = '450', modificat_la = 'y' WHERE cheie = 'norma_intreaga_minute'`)
	verificaNorma(t, db, "intreg", 450)
	verificaNorma(t, db, "partial", 225)
	if err := db.QueryRow(`SELECT concediu_zile_an FROM angajati WHERE id = 'partial'`).Scan(&concediu); err != nil || !concediu.Valid || concediu.Int64 != 24 {
		t.Fatalf("concediul contractual: %v, %v", concediu, err)
	}

	// Schimbarea tipului și a normei se face împreună, fără stare ambiguă.
	executa(t, db, `UPDATE angajati SET tip_contract = 'INTREAGA', norma_partiala_minute = NULL WHERE id = 'partial'`)
	verificaNorma(t, db, "partial", 450)
	executa(t, db, `UPDATE angajati SET tip_contract = 'PARTIALA', norma_partiala_minute = 180 WHERE id = 'intreg'`)
	verificaNorma(t, db, "intreg", 180)
}

func TestConstrangeriContractuale(t *testing.T) {
	db := bazaContracte(t)
	cazuri := []struct {
		nume string
		sql  string
	}{
		{"tip necunoscut", `INSERT INTO angajati (id, nume, marca, tip_contract, creat_la, modificat_la) VALUES ('a', 'Test', '1', 'ALTUL', 'x', 'x')`},
		{"parțial fără ore", `INSERT INTO angajati (id, nume, marca, tip_contract, creat_la, modificat_la) VALUES ('a', 'Test', '1', 'PARTIALA', 'x', 'x')`},
		{"normă zero", `INSERT INTO angajati (id, nume, marca, tip_contract, norma_partiala_minute, creat_la, modificat_la) VALUES ('a', 'Test', '1', 'PARTIALA', 0, 'x', 'x')`},
		{"normă negativă", `INSERT INTO angajati (id, nume, marca, tip_contract, norma_partiala_minute, creat_la, modificat_la) VALUES ('a', 'Test', '1', 'PARTIALA', -60, 'x', 'x')`},
		{"normă peste o zi", `INSERT INTO angajati (id, nume, marca, tip_contract, norma_partiala_minute, creat_la, modificat_la) VALUES ('a', 'Test', '1', 'PARTIALA', 1441, 'x', 'x')`},
		{"întreg cu normă individuală", `INSERT INTO angajati (id, nume, marca, norma_partiala_minute, creat_la, modificat_la) VALUES ('a', 'Test', '1', 240, 'x', 'x')`},
		{"concediu negativ", `INSERT INTO angajati (id, nume, marca, concediu_zile_an, creat_la, modificat_la) VALUES ('a', 'Test', '1', -1, 'x', 'x')`},
		{"concediu fracționar", `INSERT INTO angajati (id, nume, marca, concediu_zile_an, creat_la, modificat_la) VALUES ('a', 'Test', '1', 20.5, 'x', 'x')`},
		{"normă globală eliminată", `DELETE FROM setari WHERE cheie = 'norma_intreaga_minute'`},
		{"normă globală redenumită", `UPDATE setari SET cheie = 'altceva' WHERE cheie = 'norma_intreaga_minute'`},
		{"normă invalidă la inserare", `INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('norma_intreaga_minute', 'abc', 'x') ON CONFLICT (cheie) DO UPDATE SET valoare = excluded.valoare`},
	}
	for _, caz := range cazuri {
		t.Run(caz.nume, func(t *testing.T) {
			if _, err := db.Exec(caz.sql); err == nil {
				t.Fatal("date invalide acceptate")
			}
		})
	}
	for _, valoare := range []string{"0", "-1", "1441", "șapte", "420.5", "7h", "", "4.2e2"} {
		t.Run("norma globala "+valoare, func(t *testing.T) {
			if _, err := db.Exec(`UPDATE setari SET valoare = ? WHERE cheie = 'norma_intreaga_minute'`, valoare); err == nil {
				t.Fatal("normă globală invalidă acceptată")
			}
		})
	}
	executa(t, db, `INSERT INTO angajati (id, nume, marca, concediu_zile_an, creat_la, modificat_la) VALUES ('valid', 'Test', '1', 0, 'x', 'x')`)
	verificaNorma(t, db, "valid", 420)
	// Validarea normei nu restrânge celelalte setări ale aplicației.
	executa(t, db, `INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('limba', 'ro', 'x')`)
}

func TestMigrareContractePastreazaDatele(t *testing.T) {
	for _, normaGlobala := range []string{"", "450"} {
		t.Run("setare existenta "+normaGlobala, func(t *testing.T) {
			db, err := Deschide(filepath.Join(t.TempDir(), "veche.db"), cheieNoua(t))
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			initiala, err := fs.ReadFile(Migrari, "0001_schema_initiala.sql")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Migreaza(context.Background(), db, fstest.MapFS{"0001_schema_initiala.sql": {Data: initiala}}); err != nil {
				t.Fatal(err)
			}
			executa(t, db, `INSERT INTO angajati (id, nume, marca, norma_minute, creat_la, modificat_la)
				VALUES ('vechi', 'Nume complet vechi', '0007', 480, 'creat', 'modificat')`)
			executa(t, db, `INSERT INTO pontaj (id, angajat_id, data, start, sfarsit, creat_la, modificat_la)
				VALUES ('p1', 'vechi', '2026-10-25', '2026-10-25T05:00:00Z', '2026-10-25T12:00:00Z', 'x', 'x')`)
			executa(t, db, `INSERT INTO sarbatori (data, denumire, creat_la, modificat_la)
				VALUES ('2026-12-25', '{"ro":"Exemplu"}', 'creat', 'modificat')`)
			if normaGlobala != "" {
				executa(t, db, `INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('norma_intreaga_minute', ?, 'initial')`, normaGlobala)
			}
			if n, err := Migreaza(context.Background(), db, Migrari); err != nil || n != 17 {
				t.Fatalf("actualizarea bazei vechi: %d migrări, %v", n, err)
			}
			var nume, prenume, marca, creat, modificat string
			var normaVeche int
			if err := db.QueryRow(`SELECT nume, prenume, marca, norma_minute_anterioara, creat_la, modificat_la
				FROM angajati WHERE id = 'vechi'`).Scan(&nume, &prenume, &marca, &normaVeche, &creat, &modificat); err != nil {
				t.Fatal(err)
			}
			if nume != "Nume complet vechi" || prenume != "" || marca != "0007" || normaVeche != 480 || creat != "creat" || modificat != "modificat" {
				t.Fatal("datele vechi au fost modificate sau deduse automat")
			}
			asteptata := 420
			if normaGlobala == "450" {
				asteptata = 450
			}
			verificaNorma(t, db, "vechi", asteptata)
			var n int
			if err := db.QueryRow(`SELECT count(*) FROM pontaj WHERE angajat_id = 'vechi'`).Scan(&n); err != nil || n != 1 {
				t.Fatalf("pontajul existent: %d, %v", n, err)
			}
			if err := db.QueryRow(`SELECT count(*) FROM sarbatori WHERE an = 2026 AND denumire = '{"ro":"Exemplu"}'
				AND creat_la = 'creat' AND modificat_la = 'modificat'`).Scan(&n); err != nil || n != 1 {
				t.Fatalf("sărbătoarea existentă: %d, %v", n, err)
			}
			randuri, err := db.Query(`PRAGMA foreign_key_check`)
			if err != nil {
				t.Fatal(err)
			}
			defer randuri.Close()
			if randuri.Next() {
				t.Error("migrarea a deteriorat o legătură între tabele")
			}
			if err := randuri.Err(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSarbatoriPeAni(t *testing.T) {
	db := bazaContracte(t)
	for _, data := range []string{"2026-12-24", "2027-12-25", "2027-04-01"} {
		executa(t, db, `INSERT INTO sarbatori (data, denumire, creat_la, modificat_la) VALUES (?, '{"ro":"Exemplu"}', 'x', 'x')`, data)
	}
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM sarbatori WHERE an = 2027`).Scan(&n); err != nil || n != 2 {
		t.Fatalf("sărbători în 2027: %d, %v", n, err)
	}
	executa(t, db, `UPDATE sarbatori SET data = '2028-04-01' WHERE data = '2027-04-01'`)
	if err := db.QueryRow(`SELECT an FROM sarbatori WHERE data = '2028-04-01'`).Scan(&n); err != nil || n != 2028 {
		t.Fatalf("anul derivat după modificare: %d, %v", n, err)
	}
	if _, err := db.Exec(`UPDATE sarbatori SET an = 2029 WHERE data = '2028-04-01'`); err == nil {
		t.Error("a fost acceptat un an care contrazice data")
	}
	if _, err := db.Exec(`INSERT INTO sarbatori (data, denumire, creat_la, modificat_la) VALUES ('2026-12-24', '{"ro":"Duplicat"}', 'x', 'x')`); err == nil {
		t.Error("sărbătoare duplicată acceptată în același an")
	}
}
