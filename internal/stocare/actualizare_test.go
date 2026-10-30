package stocare

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// Copierea integrala foloseste numai date fictive si nu modifica sursa.
func TestActualizarePrinCopieCompleta(t *testing.T) {
	ctx := context.Background()
	sursa := t.TempDir()
	db, paza, err := deschideTest(ctx, sursa)
	if err != nil {
		t.Fatal(err)
	}
	if err := paza.Sistem(func() error {
		_, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('proba_actualizare', 'fictiv', 'x')`)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	backup, err := BackupManual(ctx, db, sursa, seifTest(t, sursa).CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	destinatie := filepath.Join(t.TempDir(), "karios")
	if err := os.CopyFS(destinatie, os.DirFS(sursa)); err != nil {
		t.Fatal(err)
	}
	copie, err := Deschide(ctx, destinatie)
	if err != nil {
		t.Fatal(err)
	}
	if copie.DB == nil {
		t.Fatal("copia completa nu poate fi deschisa")
	}
	defer copie.DB.Close()
	if valoare(t, copie.DB, "proba_actualizare") != "fictiv" {
		t.Fatal("datele nu au fost pastrate")
	}
	if _, err := os.Stat(filepath.Join(destinatie, "backup", filepath.Base(backup))); err != nil {
		t.Fatal(err)
	}
	original, _, err := deschideTest(ctx, sursa)
	if err != nil {
		t.Fatal(err)
	}
	defer original.Close()
	if valoare(t, original, "proba_actualizare") != "fictiv" {
		t.Fatal("sursa a fost modificata")
	}
}
