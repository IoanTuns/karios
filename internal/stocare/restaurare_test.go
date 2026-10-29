package stocare

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/seif"
)

func valoare(t *testing.T, db *sql.DB, cheie string) string {
	t.Helper()
	var v string
	_ = db.QueryRow(`SELECT valoare FROM setari WHERE cheie = ?`, cheie).Scan(&v)
	return v
}

func TestRestaurareaBackupuluiLaPornire(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, paza, err := Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	paza.Sistem(func() error {
		_, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('inainte', 'da', 'x')`)
		return err
	})
	cale, err := BackupManual(ctx, db, dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('dupa', 'da', 'x')`); err != nil {
		t.Fatal(err)
	}
	db.Close()

	nume := filepath.Base(cale)
	if r, err := VerificaBackup(ctx, dir, nume); err != nil || r.Stare != audit.StareOK {
		t.Fatalf("verificarea backup-ului: %+v %v", r, err)
	}
	if err := ProgrameazaRestaurare(dir, "../"+nume); err == nil {
		t.Fatal("nume cu cale acceptat")
	}
	if err := ProgrameazaRestaurare(dir, nume); err != nil || RestaurareProgramata(dir) != nume {
		t.Fatalf("programare: %v", err)
	}

	db, paza, err = Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if r := paza.Rezultat(); r.Stare != audit.StareOK {
		t.Fatalf("după restaurare: %+v", r)
	}
	if valoare(t, db, "inainte") != "da" || valoare(t, db, "dupa") != "" {
		t.Fatal("baza nu a fost înlocuită cu backup-ul")
	}
	var jurnal string
	if err := db.QueryRow(`SELECT valoare_noua FROM jurnal_modificari WHERE tabel = 'audit' AND rand_id = 'restaurare'`).Scan(&jurnal); err != nil || !strings.Contains(jurnal, nume) {
		t.Fatalf("restaurarea în jurnal: %q %v", jurnal, err)
	}
	if RestaurareProgramata(dir) != "" {
		t.Fatal("programarea a rămas")
	}
	// Baza dinainte de restaurare e păstrată ca backup de siguranță, cu datele ei.
	lista, err := ListaBackupuri(dir)
	if err != nil {
		t.Fatal(err)
	}
	var siguranta string
	for _, b := range lista {
		if strings.HasPrefix(b.Nume, prefixInainteRestau) {
			siguranta = b.Nume
		}
	}
	if siguranta == "" {
		t.Fatalf("lipsește copia de siguranță: %+v", lista)
	}
	s, _ := seif.Deschide(filepath.Join(dir, NumeSeif))
	copie, err := bd.Deschide(filepath.Join(dir, "backup", siguranta), s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	defer copie.Close()
	if valoare(t, copie, "dupa") != "da" {
		t.Fatal("copia de siguranță nu conține datele dinainte de restaurare")
	}
}

func TestBackupModificatEsteRefuzat(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	cale, err := BackupManual(ctx, db, dir)
	db.Close()
	if err != nil {
		t.Fatal(err)
	}
	s, _ := seif.Deschide(filepath.Join(dir, NumeSeif))
	b, err := bd.Deschide(cale, s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('strecurat', 'da', 'x')`); err != nil {
		t.Fatal(err)
	}
	b.Close()
	if r, err := VerificaBackup(ctx, dir, filepath.Base(cale)); err != nil || r.Stare != audit.StareIncident {
		t.Fatalf("backup modificat: %+v %v", r, err)
	}
}
