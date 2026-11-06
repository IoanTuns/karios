package datebaza

import (
	"context"
	"pontix/internal/i18n"
	"testing"
)

func TestClasificariProgramIntegrareSiProtectii(t *testing.T) {
	m, id, schimburi := mediuPlanificare(t)
	ctx := context.Background()
	for _, c := range []struct{ tip, nume, comportament string }{{"tipuri_program", "Noapte exemplu", "NOAPTE"}, {"regimuri_program", "Normal exemplu", "LUCRATOARE"}} {
		if err := m.serviciu.SalveazaElementNomenclator(ctx, c.tip, ElementNomenclator{Denumire: i18n.Text{"ro": c.nume, "en": c.nume}, Activ: true, Comportament: c.comportament}); err != nil {
			t.Fatal(err)
		}
	}
	d := citeste(t, m.serviciu)
	var tip, regim ElementNomenclator
	for _, e := range d.TipuriProgram {
		if e.Denumire.In("ro") == "Noapte exemplu" {
			tip = e
		}
	}
	for _, e := range d.RegimuriProgram {
		if e.Denumire.In("ro") == "Normal exemplu" {
			regim = e
		}
	}
	lista, err := m.serviciu.ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var program Schimb
	for _, s := range lista {
		if s.ID == schimburi["N12"] {
			program = s
		}
		if s.Regim == "GENERAL" {
			t.Fatal("general nemigrat")
		}
	}
	program.TipID, program.RegimID = tip.ID, regim.ID
	program.Tip = "NORMAL" // Comportamentul vine din definitie, nu din campul trimis de client.
	if err := m.serviciu.SalveazaSchimb(ctx, program); err != nil {
		t.Fatal(err)
	}
	lista, err = m.serviciu.ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range lista {
		if s.ID == program.ID {
			if s.Tip != "NOAPTE" || s.Regim != "LUCRATOARE" || s.TipID != tip.ID || s.RegimID != regim.ID {
				t.Fatal("legaturi/comportament pierdute", s)
			}
			program = s
		}
	}
	// Aceeasi identitate N12 ramane in categoria programului si in planificare.
	planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-05", SchimbID: program.ID, Sectiune: "GARZI", Activ: true})
	d = citeste(t, m.serviciu)
	for _, e := range d.TipuriProgram {
		if e.ID == tip.ID {
			tip = e
		}
	}
	if tip.Utilizari != 1 {
		t.Fatal("utilizari incorecte", tip.Utilizari)
	}
	veche := tip
	tip.Denumire = i18n.Text{"ro": "Noapte redenumita", "en": "Renamed night"}
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "tipuri_program", tip); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaElementNomenclator(ctx, "tipuri_program", veche), "nomenclator.conflict")
	d = citeste(t, m.serviciu)
	for _, e := range d.TipuriProgram {
		if e.ID == tip.ID {
			tip = e
		}
	}
	tip.Comportament = "NORMAL"
	verificaCod(t, m.serviciu.SalveazaElementNomenclator(ctx, "tipuri_program", tip), "nomenclator.comportament_utilizat")
	tip.Comportament = "NOAPTE"
	tip.Activ = false
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "tipuri_program", tip); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaSchimb(ctx, program); err != nil {
		t.Fatal("program existent cu definitie inactiva", err)
	}
	program.ID = ""
	program.Cod = "NOU-TEST"
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, program), "nomenclator.inactiv")
}
