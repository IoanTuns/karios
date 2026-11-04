package rapoarte

import (
	"context"
	"os"
	"path/filepath"
	"pontix/internal/i18n"
	"testing"
)

func TestExempluPublicGarziCuPrezentare(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	ctx := context.Background()
	for _, zi := range []string{"2026-11-27", "2026-11-28", "2026-11-29", "2026-11-30", "2026-12-01"} {
		if _, e := m.db.Exec(`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,creat_la,modificat_la) SELECT ?, 'g', ?, id,'GARZI','x','x' FROM tipuri_schimb WHERE cod='G24'`, zi, zi); e != nil {
			t.Fatal(e)
		}
	}
	q, e := os.ReadFile("../../docs/exemple/garzi-lunar.sql")
	if e != nil {
		t.Fatal(e)
	}
	d := modelFlexTest()
	d.Continut.Sursa = "LISTA_GARZI"
	f := d.Continut.Flexibil
	f.SQL = string(q)
	f.ModAspect = "MARKDOWN"
	f.Markdown = "{{tabel}}"
	f.Ascunse = []string{"data_calendar"}
	f.ColoanaData = "data_calendar"
	f.BoldWeekend = true
	f.BoldSarbatori = true
	f.ExtindeZileLibere = true
	d.Antet.Centru = i18n.Text{"ro": "TABEL CU ECHIPELE DE GARDA\n{{perioada}}"}
	d.Antet.Dreapta = i18n.Text{"ro": "SE APROBA:\nDIRECTOR MEDICAL"}
	d.Subsol.Stanga = i18n.Text{"ro": "Intocmit: ______________"}
	d.Subsol.Dreapta = i18n.Text{"ro": "Verificat: ______________"}
	d.Continut.AscundeExplicatii = true
	d.Continut.Prezentare = &PrezentareRaport{Stiluri: map[string]StilText{"antet_centru": {FontPT: 14, Aliniere: "C", Aldin: logic(true)}, "antet_dreapta": {FontPT: 10, Aliniere: "R"}}, Etichete: map[string]i18n.Text{"data": {"ro": "Data"}, "ziua": {"ro": "Ziua"}, "nume_si_prenume": {"ro": "Nume si prenume"}, "observatii": {"ro": "Observatii"}}}
	doc, e := m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 11, Limba: "ro", AngajatID: "g"}, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	if len(doc.Coloane) != 4 || len(doc.Randuri) != 5 || doc.Randuri[0].Aldin || !doc.Randuri[4].Aldin || doc.Randuri[4].Valori[0] != "1" {
		t.Fatal(doc.Randuri)
	}
	pdf, e := GenereazaPDF(doc)
	if e != nil {
		t.Fatal(e)
	}
	if *directorDemo != "" {
		if e := os.MkdirAll(*directorDemo, 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(filepath.Join(*directorDemo, "garzi-stilizate-noiembrie.pdf"), pdf, 0600); e != nil {
			t.Fatal(e)
		}
	}
	qCalendar, err := os.ReadFile("../../docs/exemple/garzi-lunar-calendar.sql")
	if err != nil {
		t.Fatal(err)
	}
	f.SQL = string(qCalendar)
	complet, err := m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 11, Limba: "ro", AngajatID: "g"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(complet.Randuri) != 32 || complet.Randuri[0].Valori[0] != "31" || complet.Randuri[0].Valori[2] != "" {
		t.Fatal("calendar complet incorect", complet.Randuri)
	}
	f.SQL = string(q)
	// Campurile repetate primesc stilul automat si pot fi editate individual.
	f.ModAspect = "LIBER"
	f.InaltimeRand = 15
	f.Elemente = []ElementRaport{{ID: "camp", Text: "{{nume_si_prenume}}", X: 5, Y: 20, Latime: 150, Inaltime: 12, Font: 10, RepetaRand: true, Stil: &StilText{Cursiv: logic(true)}}}
	doc, e = m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 11, Limba: "ro", AngajatID: "g"}, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	doc, e = EditeazaDocument(doc, []EditareTextRaport{{Cheie: "celula:0:nume_si_prenume", Text: "Nume editat", Stil: &StilText{FontPT: 14, Aldin: logic(true), Aliniere: "C"}}}, "Demo", momentFix)
	if e != nil {
		t.Fatal(e)
	}
	if _, e := GenereazaPDF(doc); e != nil {
		t.Fatal(e)
	}
}
