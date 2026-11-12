package datebaza

import (
	"context"
	"testing"

	"pontix/internal/i18n"
)

// F03-T11 (D106): importul pastreaza doar textele diferite de cele incluse, refuza cheile
// necunoscute si parametrii schimbati, iar revenirea la textul inclus sterge corectura.
func TestImportTraduceri(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	inclusEn := i18n.TextInclus(i18n.Engleza, "ui.salveaza")
	fisiere := []FisierTraduceri{
		{Limba: "en", Texte: map[string]string{
			"ui.salveaza":                       "Store",
			"ui.rapoarte":                       i18n.TextInclus(i18n.Engleza, "ui.rapoarte"), // neschimbat
			"cheie.inexistenta":                 "x",
			"pontaj.schimbare_ora_luna_inchisa": "Month {{an}} is closed", // parametru gresit
		}},
		{Limba: "hu", Texte: map[string]string{"ui.salveaza": "Mentés", "ui.rapoarte": ""}},
	}
	prev, err := m.serviciu.PrevizualizeazaTraduceri(ctx, fisiere)
	if err != nil {
		t.Fatal(err)
	}
	if len(prev) != 2 || prev[0].Modificate != 1 || prev[0].Neschimbate != 1 || len(prev[0].Probleme) != 2 || prev[1].Noi != 1 || prev[1].Neschimbate != 1 {
		t.Fatalf("previzualizare %+v", prev)
	}
	if c, _ := m.serviciu.ListaTraduceri(ctx); len(c) != 0 {
		t.Fatal("previzualizarea a scris", c)
	}
	if _, err := m.serviciu.ImportaTraduceri(ctx, fisiere); err != nil {
		t.Fatal(err)
	}
	c, err := m.serviciu.ListaTraduceri(ctx)
	if err != nil || c["en"]["ui.salveaza"] != "Store" || c["hu"]["ui.salveaza"] != "Mentés" || len(c["en"]) != 1 || len(c["hu"]) != 1 {
		t.Fatal("dupa import", c, err)
	}
	// Al doilea import identic nu schimba nimic; revenirea la textul inclus sterge corectura.
	r, err := m.serviciu.ImportaTraduceri(ctx, []FisierTraduceri{{Limba: "en", Texte: map[string]string{"ui.salveaza": inclusEn}}})
	if err != nil || r[0].Revenite != 1 {
		t.Fatal("revenire", r, err)
	}
	if c, _ = m.serviciu.ListaTraduceri(ctx); len(c["en"]) != 0 {
		t.Fatal("corectura en ramasa", c)
	}
	verificaCod(t, func() error { _, e := m.serviciu.ImportaTraduceri(ctx, []FisierTraduceri{{Limba: "Romana"}}); return e }(), "traduceri.limba_invalida")
	if err := m.serviciu.StergeTraduceri(ctx, "hu"); err != nil {
		t.Fatal(err)
	}
	if c, _ = m.serviciu.ListaTraduceri(ctx); len(c) != 0 {
		t.Fatal("dupa stergere", c)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel='traduceri'`).Scan(&n); err != nil || n != 4 {
		t.Fatal("jurnal", n, err)
	}
}
