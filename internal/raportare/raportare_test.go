package raportare

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"pontix/internal/i18n"
)

func context() Context {
	return Context{Versiune: "0.88.12", Commit: "abc1234", Sistem: "windows/amd64", Moment: time.Date(2026, 10, 30, 9, 0, 0, 0, time.UTC), NumarErori: 2,
		Jurnal: "[2026-10-30] [EROARE] [interogare/Pontaj] baza blocată"}
}

func TestRaportulNuContineDatePersonale(t *testing.T) {
	c := Cerere{Tip: TipProblema, Descriere: `Grila nu salvează pentru "Pop Ana", CNP 2850101123456`, Pasi: "1. Deschid pontajul\n2. Scriu Z7 la ana.pop@spital.ro", CuJurnal: true, Ecran: "pontaj"}
	text := Text(i18n.Romana, c, context())
	for _, interzis := range []string{"Pop Ana", "2850101123456", "ana.pop", "spital.ro"} {
		if strings.Contains(text, interzis) {
			t.Fatalf("a rămas %q în raport:\n%s", interzis, text)
		}
	}
	for _, trebuie := range []string{"0.88.12", "abc1234", "windows/amd64", "Grila nu salvează", "Deschid pontajul", "baza blocată", "pontaj"} {
		if !strings.Contains(text, trebuie) {
			t.Fatalf("lipsește %q din raport:\n%s", trebuie, text)
		}
	}
	if strings.Contains(Text(i18n.Romana, Cerere{Tip: TipParere, Descriere: "x"}, context()), "baza blocată") {
		t.Fatal("jurnalul a intrat fără acord")
	}
}

func TestAdresaIssue(t *testing.T) {
	c := Cerere{Tip: TipParere, Descriere: "Aș vrea export pe săptămână\nși încă ceva"}
	titlu := Titlu(i18n.Romana, c)
	if !strings.HasSuffix(titlu, "Aș vrea export pe săptămână") || strings.Contains(titlu, "încă") {
		t.Fatalf("titlu: %q", titlu)
	}
	a := AdresaIssue(i18n.Romana, titlu, "corp")
	u, err := url.Parse(a)
	if err != nil || !strings.HasPrefix(a, AdresaRepo+"/issues/new?") || u.Query().Get("title") != titlu || u.Query().Get("body") != "corp" {
		t.Fatalf("adresa: %q %v", a, err)
	}
	lung := AdresaIssue(i18n.Romana, titlu, strings.Repeat("ă", LungimeMaximaCorpIssue+500))
	u, _ = url.Parse(lung)
	if n := len([]rune(u.Query().Get("body"))); n > LungimeMaximaCorpIssue+200 {
		t.Fatalf("corpul nu a fost scurtat: %d", n)
	}
}
