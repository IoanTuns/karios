package datebaza

import (
	"context"
	"testing"
)

func schitaConversieTest(t *testing.T, m *mediuTest, id string) []CelulaSchita {
	t.Helper()
	ctx := context.Background()
	if e := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-05", Valoare: "Z7"}, {AngajatID: id, Data: "2026-10-06", Valoare: "Z7"}}); e != nil {
		t.Fatal(e)
	}
	s, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil {
		t.Fatal(e)
	}
	return s.Celule
}

func TestConversiePontajAbsenteSiSchimbareTip(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	celule := schitaConversieTest(t, m, a.ID)
	z := ZiAbsenta{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true, StareAprobare: "APROBATA", AprobatDe: "Operator test"}
	if e := m.serviciu.InlocuiestePontajCuAbsente(ctx, 2026, 10, []ZiAbsenta{z}, celule[:1]); e != nil {
		t.Fatal(e)
	}
	s, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Absente) != 1 || s.Absente[0].Cod != "CO" || s.Absente[0].StareAprobare != "APROBATA" || s.Celule[0].Valoare != "" || s.Celule[1].Valoare != "Z7" {
		t.Fatal(s)
	}
	z = s.Absente[0]
	z.Cod = "CM"
	if e := m.serviciu.InlocuiestePontajCuAbsente(ctx, 2026, 10, []ZiAbsenta{z}, s.Celule[:1]); e != nil {
		t.Fatal(e)
	}
	if absenteTest(t, m, 2026).Zile[0].Cod != "CM" {
		t.Fatal("tip nemodificat")
	}
}

func TestConversieAbsenteRollbackProtectii(t *testing.T) {
	for _, caz := range []string{"tip_invalid", "conflict", "luna_inchisa", "selectie_invalida", "respinsa"} {
		t.Run(caz, func(t *testing.T) {
			m := pregateste(t)
			a := persoanaAbsente(t, m)
			ctx := context.Background()
			celule := schitaConversieTest(t, m, a.ID)
			zile := []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true}, {AngajatID: a.ID, Data: "2026-10-06", Cod: "CO", Activ: true}}
			cod := "absente.tip_invalid"
			switch caz {
			case "tip_invalid":
				zile[1].Cod = "NECUNOSCUT"
			case "conflict":
				celule[1].Revizie = "veche"
				cod = "pontaj.conflict"
			case "luna_inchisa":
				if _, e := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "BLOCATA", "test", "operator"); e != nil {
					t.Fatal(e)
				}
				cod = "pontaj.luna_inchisa"
			case "selectie_invalida":
				celule[1].Data = "2026-10-07"
				cod = "absente.lot_invalid"
			case "respinsa":
				zile[1].StareAprobare = "respinsa"
				cod = "absente.lot_invalid"
			}
			verificaCod(t, m.serviciu.InlocuiestePontajCuAbsente(ctx, 2026, 10, zile, celule), cod)
			s, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
			if e != nil {
				t.Fatal(e)
			}
			if len(s.Absente) != 0 || len(s.Celule) != 2 || s.Celule[0].Valoare != "Z7" || s.Celule[1].Valoare != "Z7" {
				t.Fatal("salvare partiala", s)
			}
		})
	}
}
