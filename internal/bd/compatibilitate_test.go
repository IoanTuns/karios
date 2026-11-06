package bd

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"pontix/internal/i18n"
	"pontix/internal/versiune"
)

func bazaMigrata(t *testing.T) (*sql.DB, int) {
	t.Helper()
	db, err := Deschide(filepath.Join(t.TempDir(), "istoric.db"), make([]byte, LungimeCheie))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	n, err := Migreaza(context.Background(), db, Migrari)
	if err != nil {
		t.Fatal(err)
	}
	return db, n
}

// Fiecare migrare aplicata are in istoric versiunea aplicatiei (D009 punctul 3).
func TestIstoricMigrariCuVersiuneaAplicatiei(t *testing.T) {
	db, n := bazaMigrata(t)
	var randuri, aplicatie int
	if err := db.QueryRow(`SELECT count(*), sum(aplicatie = ?) FROM istoric_migrari`, versiune.Aplicatie).Scan(&randuri, &aplicatie); err != nil {
		t.Fatal(err)
	}
	if randuri != n || aplicatie != n {
		t.Fatalf("istoric: %d randuri, %d cu versiunea aplicatiei, %d migrari", randuri, aplicatie, n)
	}
	// O pornire fara migrari noi nu adauga nimic.
	if m, err := Migreaza(context.Background(), db, Migrari); err != nil || m != 0 {
		t.Fatal(m, err)
	}
	var dupa int
	_ = db.QueryRow(`SELECT count(*) FROM istoric_migrari`).Scan(&dupa)
	if dupa != randuri {
		t.Fatalf("istoric modificat fara migrari: %d", dupa)
	}
}

// O baza migrata de o versiune mai noua este refuzata, inainte de orice schimbare.
func TestBazaMaiNouaRefuzata(t *testing.T) {
	db, _ := bazaMigrata(t)
	if _, err := db.Exec(`INSERT INTO goose_db_version(version_id,is_applied) VALUES(9999,1)`); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	verifica := func(nume string, err error) {
		t.Helper()
		var e *i18n.Eroare
		if !errors.As(err, &e) || e.Cod != "bd.baza_mai_noua" || e.Param["baza"] != "9999" {
			t.Fatalf("%s: astept bd.baza_mai_noua, primit %v", nume, err)
		}
	}
	_, _, err := StareMigrari(ctx, db, Migrari)
	verifica("StareMigrari", err)
	_, err = Migreaza(ctx, db, Migrari)
	verifica("Migreaza", err)
}

func TestVersiuneMaximaMigrari(t *testing.T) {
	maxim, err := versiuneMaximaMigrari(Migrari)
	if err != nil || maxim < 34 {
		t.Fatalf("versiune maxima: %d %v", maxim, err)
	}
	if err := verificaCompatibilitate(maxim, Migrari); err != nil {
		t.Fatal(err)
	}
}
