package datebaza

import (
	"context"
	"testing"

	"pontix/internal/ai"
)

func TestPlanAIVerificaSiPlanificareaSeparata(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := CuAsistentAI(context.Background())
	date := citeste(t, m.serviciu)
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "persoane", Marci: []string{date.Angajati[0].Marca}}, DeLa: "2026-10-30", PanaLa: "2026-11-02", CodSchimb: "G24", FiltruZile: "lucratoare", ModExistente: "pastreaza"}
	p, err := m.serviciu.PregatestePlanAI(ctx, i)
	if err != nil {
		t.Fatal(err)
	}
	if p.Noi != 2 {
		t.Fatalf("plan: %+v", p)
	}
	rez, err := m.serviciu.AplicaPlanAI(ctx, i, p.Token)
	if err != nil {
		t.Fatal(err)
	}
	if rez.Aplicate != 2 {
		t.Fatalf("aplicate: %d", rez.Aplicate)
	}
	for _, luna := range []int{10, 11} {
		p, err := m.serviciu.CitestePlanificare(ctx, 2026, luna)
		if err != nil {
			t.Fatal(err)
		}
		if len(p.Pozitii) != 1 || p.Pozitii[0].AngajatID != id || !p.Pozitii[0].Activ {
			t.Fatalf("planificare luna %d: %+v", luna, p)
		}
	}
}

func TestPlanAIFaraPontajPeFiecareLuna(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "G24"}}); err != nil {
		t.Fatal(err)
	}
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti", FaraPontaj: true}, DeLa: "2026-10-30", PanaLa: "2026-11-02", CodSchimb: "G24", FiltruZile: "lucratoare", ModExistente: "pastreaza"}
	p, err := m.serviciu.PregatestePlanAI(ctx, i)
	if err != nil {
		t.Fatal(err)
	}
	if p.Noi != 1 || p.Blocate != 1 || p.Randuri[0].Motiv != "fara_pontaj" || p.Randuri[1].Data != "2026-11-02" {
		t.Fatalf("plan: %+v", p)
	}
}
