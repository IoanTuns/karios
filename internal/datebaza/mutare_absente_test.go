package datebaza

import (
	"context"
	"testing"
)

func TestMutareAbsentaAtomicaSiReactivare(t *testing.T) {
	for _, caz := range []string{"succes", "nelucratoare", "conflict_destinatie", "conflict_sursa", "luna_inchisa"} {
		t.Run(caz, func(t *testing.T) {
			m := pregateste(t)
			ctx := context.Background()
			a := persoanaAbsente(t, m)
			if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true}}); err != nil {
				t.Fatal(err)
			}
			veche := absenteTest(t, m, 2026).Zile[0]
			anulata := veche
			anulata.Activ = false
			noua := veche
			noua.Data = "2026-10-06"
			noua.Revizie = ""
			cod := ""
			switch caz {
			case "nelucratoare":
				noua.Data = "2026-10-10"
				cod = "absente.zi_nelucratoare"
			case "conflict_destinatie":
				if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{noua}); err != nil {
					t.Fatal(err)
				}
				cod = "absente.conflict"
			case "conflict_sursa":
				anulata.Revizie = "veche"
				cod = "absente.conflict"
			case "luna_inchisa":
				if _, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "BLOCATA", "test", "operator"); err != nil {
					t.Fatal(err)
				}
				cod = "pontaj.luna_inchisa"
			}
			err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{anulata, noua})
			if cod != "" {
				verificaCod(t, err, cod)
				for _, z := range absenteTest(t, m, 2026).Zile {
					if z.Data == veche.Data && !z.Activ {
						t.Fatal("vechea zi pierduta la eroare")
					}
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for _, z := range absenteTest(t, m, 2026).Zile {
				if z.Data == veche.Data && z.Activ {
					t.Fatal("ziua veche ramasa activa")
				}
				if z.Data == noua.Data {
					if !z.Activ {
						t.Fatal("destinatie inactiva")
					}
					z.Activ = false
					if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); err != nil {
						t.Fatal(err)
					}
				}
			}
			for _, z := range absenteTest(t, m, 2026).Zile {
				if z.Data == noua.Data {
					z.Activ = true
					if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); err != nil {
						t.Fatal("reactivare", err)
					}
				}
			}
		})
	}
}
