package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"pontix/internal/raportare"
)

func TestRaportSalvatFaraDatePersonale(t *testing.T) {
	a := NewApp(nil)
	a.radacina = t.TempDir()
	vechi := acumRaport
	acumRaport = func() time.Time { return time.Date(2026, 10, 30, 10, 15, 0, 0, time.UTC) }
	t.Cleanup(func() { acumRaport = vechi })

	if r := a.PregatesteRaportProblema(raportare.Cerere{}); r.Eroare == nil || r.Eroare.Cod != "raportare.descriere_lipsa" {
		t.Fatalf("fără descriere: %+v", r.Eroare)
	}
	r := a.PregatesteRaportProblema(raportare.Cerere{Tip: "altceva", Descriere: `Nu se salvează pontajul lui "Pop Ana"`})
	if r.Eroare != nil || strings.Contains(r.Text, "Pop Ana") || !strings.HasPrefix(r.Titlu, "[") {
		t.Fatalf("raport: %+v", r)
	}
	s := a.SalveazaRaportProblema(r.Text)
	if s.Eroare != nil || filepath.Base(s.Cale) != "raport-20261030-101500.md" || filepath.Dir(s.Cale) != filepath.Join(a.radacina, FolderRapoarteProbleme) {
		t.Fatalf("salvare: %+v", s)
	}
	if b, err := os.ReadFile(s.Cale); err != nil || string(b) != r.Text {
		t.Fatalf("fișierul: %v", err)
	}
	if s := a.SalveazaRaportProblema(strings.Repeat("x", LungimeMaximaRaport+1)); s.Eroare == nil {
		t.Fatal("text prea lung acceptat")
	}
	if g := a.TrimiteRaportProblemaPeGitHub("t", "x"); g.Eroare == nil {
		t.Fatal("fără fereastră Wails, trimiterea trebuia refuzată")
	}
}
