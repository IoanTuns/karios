package stocare

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/seif"
)

// folderCuDate pornește un folder nou, scrie o setare și întoarce seiful deschis la pornire.
func folderCuDate(t *testing.T) (string, *seif.Seif) {
	t.Helper()
	dir := t.TempDir()
	d, err := deschideTot(context.Background(), dir, bd.Migrari, ceas(0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.DB.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('proba', 'da', 'x')`); err != nil {
		t.Fatal(err)
	}
	d.DB.Close()
	return dir, d.Seif
}

func TestOglindaSeifuluiUrmeazaSalvarile(t *testing.T) {
	dir, s := folderCuDate(t)
	if err := s.SeteazaSecret(seif.PrefixParola+"cont-ana", "parola123"); err != nil {
		t.Fatal(err)
	}
	seifData, _ := os.ReadFile(filepath.Join(dir, NumeSeif))
	oglinda, err := os.ReadFile(caleOglinda(dir))
	if err != nil || !bytes.Equal(seifData, oglinda) {
		t.Fatalf("oglinda nu e la zi: %v", err)
	}
	lista, _ := ListaBackupuri(dir)
	for _, b := range lista {
		if b.Nume == NumeSeif {
			t.Fatal("oglinda apare printre backup-uri")
		}
	}
}

func TestSeifulPierdutSeRefaceDinOglinda(t *testing.T) {
	for _, caz := range []struct {
		nume   string
		strica func(cale string) error
	}{
		{"lipsă", os.Remove},
		{"deteriorat", func(cale string) error { return os.WriteFile(cale, []byte("nu e json"), 0o600) }},
	} {
		t.Run(caz.nume, func(t *testing.T) {
			ctx := context.Background()
			dir, s := folderCuDate(t)
			cheie := s.CheieDate
			caleSeif := filepath.Join(dir, NumeSeif)
			if err := caz.strica(caleSeif); err != nil {
				t.Fatal(err)
			}
			d, err := deschideTot(ctx, dir, bd.Migrari, ceas(1))
			if err != nil {
				t.Fatal(err)
			}
			defer d.DB.Close()
			if d.Blocat || !bytes.Equal(d.Seif.CheieDate, cheie) {
				t.Fatalf("seiful refăcut: blocat=%v", d.Blocat)
			}
			var v string
			if err := d.DB.QueryRow(`SELECT valoare FROM setari WHERE cheie = 'proba'`).Scan(&v); err != nil || v != "da" {
				t.Fatalf("datele după refacere: %q %v", v, err)
			}
			var n int
			if err := d.DB.QueryRow(`SELECT COUNT(*) FROM jurnal_modificari WHERE tabel = 'audit' AND rand_id = 'seif_refacut'`).Scan(&n); err != nil || n != 1 {
				t.Fatalf("refacerea în jurnal: %d %v", n, err)
			}
			if r := d.Paza.Rezultat(); r.Stare != "ok" {
				t.Fatalf("auditul după refacere: %+v", r)
			}
			if existaFisier(filepath.Join(dir, marcajRefacere)) {
				t.Fatal("marcajul refacerii a rămas")
			}
			if caz.nume == "deteriorat" && !existaFisier(caleSeif+sufixDeteriorat) {
				t.Fatal("seiful stricat nu a fost păstrat")
			}
		})
	}
}

func TestFaraOglindaSeifulPierdutRamaneEroare(t *testing.T) {
	dir, _ := folderCuDate(t)
	if err := os.Remove(caleOglinda(dir)); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, NumeSeif)); err != nil {
		t.Fatal(err)
	}
	if _, err := deschideTot(context.Background(), dir, bd.Migrari, ceas(1)); !errors.Is(err, seif.ErrDeteriorat) {
		t.Fatalf("fără seif și fără oglindă: %v", err)
	}
}
