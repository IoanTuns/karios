package datebaza

import (
	"context"
	"pontix/internal/ai"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"testing"
)

func TestPlanAINormalExcludeSarbatoriSiInvalideazaCalendarul(t *testing.T) {
	m, _, _ := mediuPlanificare(t)
	ctx := CuAsistentAI(context.Background())
	i := ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti"}, DeLa: "2026-10-29", PanaLa: "2026-11-02", CodSchimb: "Z7", FiltruZile: "lucratoare", ModExistente: "pastreaza"}
	vechi, e := m.serviciu.PregatestePlanAI(ctx, i)
	if e != nil {
		t.Fatal(e)
	}
	if len(vechi.Randuri) != 3 {
		t.Fatal(vechi)
	}
	for _, data := range []string{"2026-10-30", "2026-10-31"} {
		if _, e := m.db.Exec(`INSERT INTO sarbatori(data,denumire,creat_la,modificat_la) VALUES(?,'{"ro":"Sarbatoare fictiva"}','x','x')`, data); e != nil {
			t.Fatal(e)
		}
	}
	nou, e := m.serviciu.PregatestePlanAI(ctx, i)
	if e != nil {
		t.Fatal(e)
	}
	if len(nou.Randuri) != 2 || nou.Token == vechi.Token {
		t.Fatal(nou)
	}
	for _, r := range nou.Randuri {
		if r.Data != "2026-10-29" && r.Data != "2026-11-02" {
			t.Fatal(r)
		}
	}
	if _, e := m.serviciu.AplicaPlanAI(ctx, i, vechi.Token); e == nil {
		t.Fatal("calendar modificat, token vechi acceptat")
	}
	i.FiltruZile = "weekend"
	weekend, e := m.serviciu.PregatestePlanAI(ctx, i)
	if e != nil {
		t.Fatal(e)
	}
	if len(weekend.Randuri) != 2 || weekend.Randuri[0].Data != "2026-10-31" || weekend.Randuri[1].Data != "2026-11-01" {
		t.Fatal(weekend)
	}
	i.FiltruZile = "toate"
	toate, e := m.serviciu.PregatestePlanAI(ctx, i)
	if e != nil {
		t.Fatal(e)
	}
	if len(toate.Randuri) != 5 {
		t.Fatal(toate)
	}
}

func TestPontajCitesteSarbatorileFaraFiltruDeAn(t *testing.T) {
	m, id, _ := mediuPlanificare(t)
	ctx := context.Background()
	for _, data := range []string{"2026-10-30", "2026-10-31"} {
		if _, e := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,?,'Z7','x')`, id, data); e != nil {
			t.Fatal(e)
		}
		if _, e := m.db.Exec(`INSERT INTO sarbatori(data,denumire,creat_la,modificat_la) VALUES(?,'{"ro":"Sarbatoare fictiva"}','x','x')`, data); e != nil {
			t.Fatal(e)
		}
	}
	r, _, e := CalculeazaPontaj(ctx, acces.New(m.db), 2026, 10, nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, v := range r.Totaluri {
		if v.AngajatID == id {
			if calcul.TextValoarePontaj(v.Valori["ore_norma"]) != "0" || calcul.TextValoarePontaj(v.Valori["ZNL"]) != "14" || calcul.TextValoarePontaj(v.Valori["T"]) != "14" {
				t.Fatal(v)
			}
			return
		}
	}
	t.Fatal("total lipsa")
}
