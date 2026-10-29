package stocare

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/audit"
	"pontix/internal/seif"
)

func copiazaFisier(t *testing.T, sursa, tinta string) {
	t.Helper()
	in, err := os.Open(sursa)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(tinta)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		t.Fatal(err)
	}
}

// porneste deschide baza, scrie o setare (dacă e cazul) și o închide ca aplicația: cu contorul la zi.
func porneste(t *testing.T, dir, cheie string) *audit.Paza {
	t.Helper()
	ctx := context.Background()
	db, paza, err := Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if cheie != "" && !paza.Blocata() {
		if _, err := db.Exec(`INSERT INTO jurnal_modificari (tabel, rand_id, valoare_noua, cand, cine) VALUES ('setari', ?, '{}', 'x', 'local')`, cheie); err != nil {
			t.Fatal(err)
		}
	}
	if !paza.Blocata() {
		if err := ActualizeazaSecventa(ctx, db, dir); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()
	return paza
}

func TestCopiaVecheEsteDetectata(t *testing.T) {
	dir := t.TempDir()
	baza := filepath.Join(dir, NumeBaza)
	porneste(t, dir, "")
	porneste(t, dir, "a")
	veche := filepath.Join(t.TempDir(), "veche.db")
	copiazaFisier(t, baza, veche)
	porneste(t, dir, "b")
	porneste(t, dir, "c")
	if p := porneste(t, dir, ""); p.Blocata() {
		t.Fatalf("baza curentă dă incident: %+v", p.Rezultat())
	}
	// Baza înlocuită cu copia de dinainte de „b”: sigiliile ei sunt valide, dar e mai veche.
	copiazaFisier(t, veche, baza)
	p := porneste(t, dir, "")
	if r := p.Rezultat(); !p.Blocata() || !r.CopieVeche {
		t.Fatalf("copia veche nu e detectată: %+v", r)
	}
}

func TestContorFalsificat(t *testing.T) {
	dir := t.TempDir()
	porneste(t, dir, "")
	porneste(t, dir, "a")
	s, err := seif.Deschide(filepath.Join(dir, NumeSeif))
	if err != nil {
		t.Fatal(err)
	}
	c := s.Secventa()
	if c == nil || c.Numar == 0 {
		t.Fatalf("contorul nu e scris: %+v", c)
	}
	// Cineva coboară contorul ca să ascundă o copie veche: codul de autentificare nu mai corespunde.
	if err := s.SeteazaSecventa(seif.Secventa{Numar: 0, Mac: c.Mac}); err != nil {
		t.Fatal(err)
	}
	if p := porneste(t, dir, ""); !p.Blocata() || !p.Rezultat().CopieVeche {
		t.Fatalf("contorul falsificat nu e detectat: %+v", p.Rezultat())
	}
}

func TestRestaurareaNuEsteCopieVeche(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	cale, err := BackupManual(ctx, db, dir)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	porneste(t, dir, "a")
	porneste(t, dir, "b")
	if err := ProgrameazaRestaurare(dir, filepath.Base(cale)); err != nil {
		t.Fatal(err)
	}
	if p := porneste(t, dir, ""); p.Blocata() {
		t.Fatalf("restaurarea cerută a dat incident: %+v", p.Rezultat())
	}
	// După restaurare, contorul pornește de la baza restaurată.
	if p := porneste(t, dir, ""); p.Blocata() {
		t.Fatalf("pornirea de după restaurare: %+v", p.Rezultat())
	}
}
