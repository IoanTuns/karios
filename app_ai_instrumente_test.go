package main

import (
	"pontix/internal/ai"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"strings"
	"testing"
)

func TestInstrumenteAIApiSiExportVerificat(t *testing.T) {
	a := appCuDemo(t)
	if r := a.ListaSabloaneAI(); r.Eroare != nil || len(r.Sabloane) != 0 {
		t.Fatal(r)
	}
	i := intentieDemoAI()
	if r := a.SalveazaSablonAI("Dimineața", i); r.Eroare != nil || len(r.Sabloane) != 1 {
		t.Fatal(r)
	}
	p := a.PrevizualizeazaPlanAI(i)
	if p.Eroare != nil {
		t.Fatal(p.Eroare)
	}
	if r := a.ExportaPlanAICSV(i, "fals", "", "ro"); r.Eroare == nil || r.Eroare.Cod != "ai.plan_expirat" {
		t.Fatal(r)
	}
	if r := a.ExportaPlanAICSV(i, p.Plan.Token, "", "ro"); r.Eroare != nil || !r.Anulat {
		t.Fatal(r)
	}
	p = a.ConfirmaPlanAI(i, p.Plan.Token)
	if p.Eroare != nil {
		t.Fatal(p.Eroare)
	}
	if r := a.ExportaPlanAICSV(ai.IntentiePlan{}, "", p.Plan.OperatieID, "en"); r.Eroare != nil || !r.Anulat {
		t.Fatal(r)
	}
	u := a.PrevizualizeazaAnulareAI()
	if u.Eroare != nil || u.Anulare.Conflicte != 0 {
		t.Fatal(u)
	}
	if r := a.AnuleazaUltimaOperatieAI(u.Anulare.Token); r.Eroare != nil || !r.Anulare.Anulata {
		t.Fatal(r)
	}
}
func TestCSVPlanAILocalizareSiFormule(t *testing.T) {
	p := datebaza.PlanAI{Intentie: intentieDemoAI(), OperatieID: "istoric", Randuri: []datebaza.RandPlanAI{{Nume: " =HYPERLINK(test)", Marca: "+1", Data: "2026-10-01", Inainte: "6", Dupa: "Z7", Stare: "ocupata"}}}
	for _, limba := range []i18n.Limba{i18n.Romana, i18n.Engleza} {
		b, err := csvPlanAI(p, limba)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, "' =HYPERLINK(test)") || !strings.Contains(s, "'+1") || !strings.Contains(s, ";6;Z7;6;") || strings.Contains(s, "ai.") {
			t.Fatal(s)
		}
	}
}
