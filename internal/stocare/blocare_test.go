package stocare

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"pontix/internal/audit"
	"pontix/internal/i18n"
)

// Prima instanta tine blocarea si scrie cine este; a doua afla detinatorul.
func TestBlocareFolderDateSiDetinator(t *testing.T) {
	cale := filepath.Join(t.TempDir(), NumeBlocare)
	t0 := time.Date(2026, 11, 6, 9, 30, 0, 0, time.UTC)
	f1, det, err := blocheaza(cale, t0)
	if err != nil || f1 == nil || det != nil {
		t.Fatalf("prima blocare: %v %v %v", f1, det, err)
	}
	f2, det, err := blocheaza(cale, t0.Add(time.Hour))
	if err != nil || f2 != nil || det == nil {
		t.Fatalf("a doua blocare trebuia refuzata: %v %v %v", f2, det, err)
	}
	asteptat := infoCurenta(t0)
	if *det != asteptat {
		t.Fatalf("detinator %+v, astept %+v", *det, asteptat)
	}
	f1.Close()
	f3, det, err := blocheaza(cale, t0)
	if err != nil || f3 == nil || det != nil {
		t.Fatalf("dupa eliberare: %v %v %v", f3, det, err)
	}
	f3.Close()
}

// A doua instanta deschide datele fara sa scrie nimic la pornire si refuza
// scrierile obisnuite; autentificarea (operatie de sistem) ramane posibila.
func TestDeschidereDoarCitire(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := deschideTest(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	seifInainte, _ := os.ReadFile(filepath.Join(dir, NumeSeif))
	backupuriInainte, _ := os.ReadDir(filepath.Join(dir, "backup"))
	det := &InfoBlocare{Calculator: "PC-1", Utilizator: "ana", DeLa: "2026-11-06T09:30:00Z"}
	d, err := deschideDoarCitire(ctx, dir, det)
	if err != nil {
		t.Fatal(err)
	}
	defer d.DB.Close()
	if !d.DoarCitire || d.Detinator != det || !d.Paza.EsteDoarCitire() {
		t.Fatalf("deschidere: %+v", d)
	}
	_, err = d.DB.Exec(`INSERT INTO setari(cheie,valoare,modificat_la) VALUES('proba','1','x')`)
	if err == nil || !audit.EsteEditareBlocata(err) || audit.CodBlocare(err) != audit.CodDoarCitire {
		t.Fatalf("scrierea trebuia refuzata cu doar_citire: %v", err)
	}
	if err := d.Paza.Sistem(func() error {
		_, e := d.DB.Exec(`INSERT INTO setari(cheie,valoare,modificat_la) VALUES('sistem','1','x')`)
		return e
	}); err != nil {
		t.Fatalf("operatia de sistem trebuia permisa: %v", err)
	}
	seifDupa, _ := os.ReadFile(filepath.Join(dir, NumeSeif))
	backupuriDupa, _ := os.ReadDir(filepath.Join(dir, "backup"))
	if !bytes.Equal(seifInainte, seifDupa) || len(backupuriDupa) != len(backupuriInainte) {
		t.Fatal("pornirea doar in citire a scris in seif sau a facut copii")
	}
}

// Daca baza trebuie actualizata, instanta doar in citire nu o atinge.
func TestDoarCitireRefuzaActualizarea(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := deschideTest(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`DELETE FROM goose_db_version WHERE version_id=(SELECT max(version_id) FROM goose_db_version)`); err != nil {
		t.Fatal(err)
	}
	db.Close()
	_, err = deschideDoarCitire(ctx, dir, &InfoBlocare{})
	var e *i18n.Eroare
	if !errors.As(err, &e) || e.Cod != "stocare.alta_instanta_actualizare" {
		t.Fatalf("astept stocare.alta_instanta_actualizare, primit %v", err)
	}
	// Folder fara date: a doua instanta asteapta ca prima sa le creeze.
	_, err = deschideDoarCitire(ctx, t.TempDir(), &InfoBlocare{})
	if !errors.As(err, &e) || !strings.HasPrefix(e.Cod, "stocare.alta_instanta") {
		t.Fatalf("folder gol: %v", err)
	}
}

// TestBlocareIntreProcese porneste un proces care tine blocarea (D002 punctul
// 6 este despre doua copii ale aplicatiei, nu doua fire).
func TestBlocareIntreProcese(t *testing.T) {
	if cale := os.Getenv("PONTIX_TINE_BLOCARE"); cale != "" {
		f, det, err := blocheaza(cale, time.Now())
		if err != nil || f == nil || det != nil {
			os.Exit(3)
		}
		os.Stdout.WriteString("gata\n")
		_, _ = io.Copy(io.Discard, os.Stdin) // tine blocarea pana se inchide intrarea
		return
	}
	cale := filepath.Join(t.TempDir(), NumeBlocare)
	cmd := exec.Command(os.Args[0], "-test.run=^TestBlocareIntreProcese$")
	cmd.Env = append(os.Environ(), "PONTIX_TINE_BLOCARE="+cale)
	intrare, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	iesire, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	linie, err := bufio.NewReader(iesire).ReadString('\n')
	if err != nil || strings.TrimSpace(linie) != "gata" {
		t.Fatalf("procesul copil: %q %v", linie, err)
	}
	f, det, err := blocheaza(cale, time.Now())
	if err != nil || f != nil || det == nil || det.Calculator == "" {
		t.Fatalf("blocarea altui proces nu a fost vazuta: %v %+v %v", f, det, err)
	}
	intrare.Close()
	if err := cmd.Wait(); err != nil {
		t.Fatalf("procesul copil: %v", err)
	}
	f, det, err = blocheaza(cale, time.Now())
	if err != nil || f == nil || det != nil {
		t.Fatalf("dupa iesirea celuilalt proces: %v %+v %v", f, det, err)
	}
	f.Close()
}
