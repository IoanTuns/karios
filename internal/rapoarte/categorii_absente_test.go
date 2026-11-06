package rapoarte

import (
	"context"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"testing"
)

func TestRaportAbsenteCategoriePersonalizata(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	s := datebaza.Nou(m.db)
	if e := s.SalveazaCategorieAbsenta(ctx, datebaza.CategorieAbsenta{Cod: "ALTA", Denumire: i18n.Text{"ro": "Categorie demo personalizata", "en": "Custom demo category"}, Natura: "SPECIAL", EfectOre: "FARA_ORE", Activ: true}); e != nil {
		t.Fatal(e)
	}
	if e := s.SalveazaTipAbsenta(ctx, "", datebaza.TipAbsenta{Cod: "ALTAZI", Denumire: i18n.Text{"ro": "Zi demo"}, CategorieID: "ALTA", Activ: true}); e != nil {
		t.Fatal(e)
	}
	if e := s.SalveazaAbsente(ctx, []datebaza.ZiAbsenta{{AngajatID: "a", Data: "2026-10-05", Cod: "ALTAZI", Activ: true}}); e != nil {
		t.Fatal(e)
	}
	for limba, nume := range map[string]string{"ro": "Categorie demo personalizata", "en": "Custom demo category"} {
		d, e := m.s.Pregateste(ctx, Cerere{ModelID: "raport-absente", An: 2026, Luna: 10, Limba: limba}, momentFix)
		if e != nil {
			t.Fatal(e)
		}
		gasit := false
		for _, r := range d.Randuri {
			for _, v := range r.Valori {
				gasit = gasit || v == nume
			}
		}
		if !gasit {
			t.Fatal("denumirea categoriei lipseste", limba)
		}
	}
}
