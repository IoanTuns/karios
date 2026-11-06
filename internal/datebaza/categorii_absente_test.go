package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
	"testing"
)

func categorieTest(t *testing.T, m *mediuTest, cod string) CategorieAbsenta {
	t.Helper()
	for _, c := range absenteTest(t, m, 2026).Categorii {
		if c.Cod == cod {
			return c
		}
	}
	t.Fatal("categorie lipsa", cod)
	return CategorieAbsenta{}
}
func TestCategoriiAbsenteMostenireProtectieSiStare(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	c := CategorieAbsenta{Cod: "DEPLASARI", Denumire: i18n.Text{"ro": "Deplasari"}, Natura: "SPECIAL", EfectOre: "NORMA", Activ: true}
	if e := m.serviciu.SalveazaCategorieAbsenta(ctx, c); e != nil {
		t.Fatal(e)
	}
	tip := TipAbsenta{Cod: "DPL", Denumire: i18n.Text{"ro": "Deplasare"}, CategorieID: c.Cod, Categorie: "MEDICAL", EfectOre: "FARA_ORE", Activ: true}
	if e := m.serviciu.SalveazaTipAbsenta(ctx, "", tip); e != nil {
		t.Fatal(e)
	}
	d := absenteTest(t, m, 2026)
	for _, v := range d.Tipuri {
		if v.Cod == tip.Cod && (v.Categorie != "SPECIAL" || v.EfectOre != "NORMA" || v.CategorieID != c.Cod) {
			t.Fatalf("mostenire %+v", v)
		}
	}
	c = categorieTest(t, m, c.Cod)
	veche := c
	c.Denumire = i18n.Text{"ro": "Deplasari redenumite"}
	if e := m.serviciu.SalveazaCategorieAbsenta(ctx, c); e != nil {
		t.Fatal(e)
	}
	verificaCod(t, m.serviciu.SalveazaCategorieAbsenta(ctx, veche), "absente.conflict")
	c = categorieTest(t, m, c.Cod)
	c.EfectOre = "FARA_ORE"
	if e := m.serviciu.SalveazaCategorieAbsenta(ctx, c); e != nil {
		t.Fatal(e)
	}
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: tip.Cod, Activ: true, StareAprobare: "APROBATA"}}); e != nil {
		t.Fatal(e)
	}
	r, _, e := CalculeazaPontaj(ctx, acces.New(m.db), 2026, 10, nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, v := range r.Totaluri {
		if v.AngajatID == a.ID && v.Valori["ore_lucratoare"].Sutimi != 21*420*100 {
			t.Fatalf("reducere norma %+v", v)
		}
	}
	c = categorieTest(t, m, c.Cod)
	c.EfectOre = "NORMA"
	verificaCod(t, m.serviciu.SalveazaCategorieAbsenta(ctx, c), "absente.categorie_folosita")
	c = categorieTest(t, m, c.Cod)
	c.Activ = false
	c.Denumire = i18n.Text{"ro": "Categorie inactiva"}
	if e := m.serviciu.SalveazaCategorieAbsenta(ctx, c); e != nil {
		t.Fatal(e)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-06", Cod: tip.Cod, Activ: true}}), "absente.tip_invalid")
	tip.Cod = "ALTA"
	verificaCod(t, m.serviciu.SalveazaTipAbsenta(ctx, "", tip), "absente.tip_invalid")
	c = categorieTest(t, m, c.Cod)
	c.Activ = true
	if e := m.serviciu.SalveazaCategorieAbsenta(ctx, c); e != nil {
		t.Fatal(e)
	}
}
func TestCategorieMedicalaPersonalizataPastreazaPrioritatea(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	c := CategorieAbsenta{Cod: "MED_PERSONALIZAT", Denumire: i18n.Text{"ro": "Medical personalizat"}, Natura: "MEDICAL", EfectOre: "FARA_ORE", Activ: true}
	if e := m.serviciu.SalveazaCategorieAbsenta(ctx, c); e != nil {
		t.Fatal(e)
	}
	tip := TipAbsenta{Cod: "CMP", Denumire: i18n.Text{"ro": "Medical"}, CategorieID: c.Cod, Activ: true}
	if e := m.serviciu.SalveazaTipAbsenta(ctx, "", tip); e != nil {
		t.Fatal(e)
	}
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true}}); e != nil {
		t.Fatal(e)
	}
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CMP", Activ: true}}); e != nil {
		t.Fatal(e)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true}}), "absente.medical_prioritar")
}
