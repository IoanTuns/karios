package datebaza

import (
	"context"
	"pontix/internal/ai"
	"strings"
	"testing"
)

func TestMarcajAIUrmeazaLegaturaSauContractul(t *testing.T) {
	m, principal, _ := mediuPlanificare(t)
	ctx := CuAsistentAI(context.Background())
	a := angajatTest("GARDA-TEST")
	a.AngajatPrincipalID = principal
	a.ContractID = "GARDA"
	a.CategorieProgramID = "program-mixt"
	if e := m.serviciu.SalveazaAngajat(ctx, a); e != nil {
		t.Fatal(e)
	}
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "persoane", Marci: []string{a.Marca}}, DeLa: "2026-10-05", PanaLa: "2026-10-05", CodSchimb: "G24", FiltruZile: "toate", ModExistente: "pastreaza"}
	p, e := m.serviciu.PregatestePlanAI(ctx, i)
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Persoane) != 1 || !strings.Contains(p.Persoane[0].Nume, " (") || !strings.Contains(p.Randuri[0].Nume, " (") {
		t.Fatal(p)
	}
	for _, v := range citeste(t, m.serviciu).Angajati {
		if v.Marca == a.Marca {
			a = v
		}
	}
	a.AngajatPrincipalID = ""
	if e := m.serviciu.SalveazaAngajat(ctx, a); e != nil {
		t.Fatal(e)
	}
	nou, e := m.serviciu.PregatestePlanAI(ctx, i)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(nou.Persoane[0].Nume, " (") || strings.Count(nou.Randuri[0].Nume, " (") != 1 || p.Token == nou.Token {
		t.Fatal(nou)
	}
	if _, e := m.serviciu.AplicaPlanAI(ctx, i, p.Token); e == nil {
		t.Fatal("plan vechi acceptat dupa schimbarea ierarhiei")
	}
	if a.Nume != strings.TrimSpace(angajatTest("GARDA-TEST").Nume) {
		t.Fatal("nume stocat modificat")
	}
}
