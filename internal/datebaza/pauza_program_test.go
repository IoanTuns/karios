package datebaza

import (
	"context"
	"testing"

	"pontix/internal/calcul"
	"pontix/internal/i18n"
)

// F04-T05 (D107): modul pauzei se salveaza si se valideaza; programele vechi cu pauza raman
// „nestabilit” (necalculate) pana la alegere.
func TestPauzaProgramSalvareSiValidare(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	ora := func(v int64) *int64 { return &v }
	baza := func() Schimb {
		return Schimb{Cod: "P8", Denumire: i18n.Text{"ro": "Program cu pauza"}, OraStart: 7 * 60, DurataMinute: 480, PauzaMinute: 30, Tip: "NORMAL", Activ: true}
	}
	s := baza()
	s.ModPauza = "ALTUL"
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, s), "schimburi.pauza_mod_invalid")
	s = baza()
	s.ModPauza, s.PauzaMinute = calcul.PauzaDinTotal, 300
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, s), "schimburi.pauza_prea_mare")
	s = baza()
	s.ModPauza = calcul.PauzaLaOra
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, s), "schimburi.pauza_ora_lipsa")
	s.PauzaOra = ora(16 * 60) // dupa sfarsitul turei (15:00)
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, s), "schimburi.pauza_in_afara_turei")
	s.PauzaOra = ora(12 * 60)
	if e := m.serviciu.SalveazaSchimb(ctx, s); e != nil {
		t.Fatal(e)
	}
	// Implicit (mod gol): nestabilit.
	n := baza()
	n.Cod = "P9"
	if e := m.serviciu.SalveazaSchimb(ctx, n); e != nil {
		t.Fatal(e)
	}
	lista, e := m.serviciu.ListaSchimburi(ctx)
	if e != nil {
		t.Fatal(e)
	}
	gasite := 0
	for _, x := range lista {
		switch x.Cod {
		case "P8":
			gasite++
			if x.ModPauza != calcul.PauzaLaOra || x.PauzaOra == nil || *x.PauzaOra != 720 {
				t.Errorf("P8: %+v", x)
			}
		case "P9":
			gasite++
			if x.ModPauza != calcul.PauzaNestabilita || x.PauzaOra != nil {
				t.Errorf("P9: %+v", x)
			}
		}
	}
	if gasite != 2 {
		t.Fatal("programe", gasite)
	}
}
