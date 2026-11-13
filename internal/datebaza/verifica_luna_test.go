package datebaza

import (
	"context"
	"testing"

	"pontix/internal/i18n"
)

// F12-T06: verificarea lunii gaseste programele cu pauza nestabilita, celulele necalculate si
// angajatii fara inregistrari, cele mai grave primele.
func TestVerificaLuna(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	if e := m.serviciu.SalveazaSchimb(ctx, Schimb{Cod: "P8", Denumire: i18n.Text{"ro": "Cu pauza"}, OraStart: 420, DurataMinute: 480, PauzaMinute: 30, Tip: "NORMAL", Activ: true}); e != nil {
		t.Fatal(e)
	}
	v, err := m.serviciu.VerificaLuna(ctx, 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	coduri := map[string]ProblemaLuna{}
	for _, p := range v.Probleme {
		coduri[p.Cod] = p
	}
	if p := coduri["pauza_nestabilita"]; p.Numar != 1 || p.Exemple[0] != "P8" {
		t.Fatalf("pauza: %+v", v.Probleme)
	}
	if p := coduri["fara_inregistrari"]; p.Numar < 1 {
		t.Fatalf("fara inregistrari: %+v", v.Probleme)
	}
	// O celula cu programul cu pauza nestabilita devine valoare necalculata.
	if e := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-05", Valoare: "P8"}}); e != nil {
		t.Fatal(e)
	}
	v, _ = m.serviciu.VerificaLuna(ctx, 2026, 10)
	if len(v.Probleme) == 0 || v.Probleme[0].Cod != "necalculat" || v.Probleme[0].Gravitate != GravitateEroare || v.Probleme[0].Param["motiv"] != "pauza" {
		t.Fatalf("necalculat primul: %+v", v.Probleme)
	}
}
