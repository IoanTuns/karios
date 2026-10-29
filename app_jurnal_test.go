package main

import (
	"strings"
	"testing"

	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

func TestCSVJurnal(t *testing.T) {
	b, err := csvJurnal([]datebaza.IntrareJurnal{{ID: 7, Tabel: "pontaj_ciorne", RandID: "a/2026-10-01", Vechi: `{"valoare":"D8"}`, Nou: `{"valoare":"N12; noapte"}`, Cand: "2026-10-29T10:00:00Z", Cine: "local"}}, i18n.Romana)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.HasPrefix(s, bomUTF8) || strings.Count(s, "\n") != 2 || !strings.Contains(s, `"{""valoare"":""N12; noapte""}"`) {
		t.Fatalf("CSV:\n%s", s)
	}
}
