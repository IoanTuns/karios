package rapoarte

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"strings"
	"testing"
)

func TestRaportProgrameMultipleIncludeAmbeleIntervale(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	ctx := context.Background()
	if _, e := m.db.Exec(`DELETE FROM planificare;UPDATE tipuri_schimb SET ora_start=840,durata_minute=300 WHERE cod='G24';UPDATE tipuri_schimb SET ora_start=1140,durata_minute=720 WHERE cod='N12';INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,creat_la,modificat_la) SELECT cod,'g','2026-10-30',id,'GARZI','x','x' FROM tipuri_schimb WHERE cod IN ('G24','N12');INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES('g','2026-10-30','G24+N12','x')`); e != nil {
		t.Fatal(e)
	}
	d := modelFlexTest()
	d.Continut.Sursa = "LISTA_GARZI"
	d.Continut.Flexibil.SQL = "SELECT data,program,inceput,sfarsit,ore FROM garzi ORDER BY inceput"
	doc, e := m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 10, Limba: "ro", AngajatID: "g"}, momentFix)
	if e != nil || len(doc.Randuri) != 2 {
		t.Fatal(doc.Randuri, e)
	}
	if doc.Randuri[0].Valori[1] != "G24" || doc.Randuri[1].Valori[1] != "N12" {
		t.Fatal(doc.Randuri)
	}
	for _, r := range doc.Randuri {
		if !strings.Contains(r.Valori[2], "2026-10-30") {
			t.Fatal(r)
		}
	}
	if _, e = GenereazaPDF(doc); e != nil {
		t.Fatal(e)
	}
	if _, e = GenereazaExcel(doc); e != nil {
		t.Fatal(e)
	}
	rez, _, e := datebaza.CalculeazaPontaj(ctx, acces.New(m.db), 2026, 10, nil)
	if e != nil {
		t.Fatal(e)
	}
	gasit := false
	for _, v := range rez.Totaluri {
		if v.AngajatID == "g" {
			gasit = true
			if calcul.TextValoarePontaj(v.Valori["T"]) != "17" {
				t.Fatal(v)
			}
		}
	}
	if !gasit {
		t.Fatal("total lipsa")
	}
}
