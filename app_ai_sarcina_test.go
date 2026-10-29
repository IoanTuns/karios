package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"
	"time"

	"pontix/internal/ai"
	"pontix/internal/bd"
	"pontix/internal/datebaza"
)

func appCuDemo(t *testing.T) *App {
	t.Helper()
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "ai.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	a := NewApp(db)
	if r := a.IncarcaDemo(); r.Eroare != nil {
		t.Fatalf("demo: %+v", r.Eroare)
	}
	return a
}

func TestSarcinaAIRuleazaInFundalSiSeElibereaza(t *testing.T) {
	a := appCuDemo(t)
	r := a.PornesteSarcinaAI(CerereInterogareAI{An: 2026, Luna: 10, Istoric: []ai.MesajConversatie{{Rol: "user", Continut: "ajutor"}}})
	if r.Eroare != nil || r.ID == "" {
		t.Fatalf("pornire: %+v", r)
	}
	limita := time.Now().Add(10 * time.Second)
	for a.SarcinaAIActiva() != "" {
		if time.Now().After(limita) {
			t.Fatal("sarcina nu s-a terminat")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if a.OpresteSarcinaAI() {
		t.Fatal("oprire fără sarcină activă")
	}
}

func TestSarcinaAIUnaSinguraSiOprire(t *testing.T) {
	a := appCuDemo(t)
	ctx, anuleaza := context.WithCancel(context.Background())
	a.aiSarcina = &sarcinaAI{id: "x", anuleaza: anuleaza}
	if r := a.PornesteSarcinaAI(CerereInterogareAI{}); r.Eroare == nil || r.Eroare.Cod != "ai.sarcina_activa" {
		t.Fatalf("a doua sarcină: %+v", r)
	}
	if a.SarcinaAIActiva() != "x" || !a.OpresteSarcinaAI() {
		t.Fatal("oprirea sarcinii active")
	}
	if ctx.Err() == nil {
		t.Fatal("contextul sarcinii nu a fost anulat")
	}
}

func TestSarcinaAIOpritaNuIntoarceRaspuns(t *testing.T) {
	a := appCuDemo(t)
	ctx, anuleaza := context.WithCancel(context.Background())
	anuleaza()
	rez := a.ruleazaSarcinaAI(ctx, "y", CerereInterogareAI{An: 2026, Luna: 10, Istoric: []ai.MesajConversatie{{Rol: "user", Continut: "ajutor"}}})
	if !rez.Oprita || rez.Raspuns != nil || rez.Eroare != nil {
		t.Fatalf("rezultat: %+v", rez)
	}
}

func TestPropunereGrupToataLunaPentruTotiAngajatii(t *testing.T) {
	a := appCuDemo(t)
	date, err := datebaza.Nou(a.db).Citeste(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	activi := 0
	for _, ang := range date.Angajati {
		if ang.Activ {
			activi++
		}
	}
	r := a.InterogheazaAI(CerereInterogareAI{An: 2026, Luna: 10, Istoric: []ai.MesajConversatie{{Rol: "user", Continut: "Toti angajatii pe luna octombrie sunt pe zi de la 7 la 14"}}})
	if r.Eroare != nil || r.Raspuns == nil || len(r.Raspuns.Propuneri) != 1 {
		t.Fatalf("interogare: %+v", r)
	}
	p := r.Raspuns.Propuneri[0]
	if p.Tip != ai.TipPlanificaGrup || len(p.Parametri["marci"].([]any)) != activi {
		t.Fatalf("propunere: %+v (activi %d)", p, activi)
	}
	// Un angajat are concediu pe 2026-10-05: ziua se sare.
	idCuConcediu := ""
	for _, ang := range date.Angajati {
		if ang.Activ {
			idCuConcediu = ang.ID
			break
		}
	}
	if _, err := a.db.Exec(`INSERT INTO absente_zile (angajat_id, data, cod, activ, creat_la, modificat_la) VALUES (?, '2026-10-05', 'CO', 1, 'x', 'x')`, idCuConcediu); err != nil {
		t.Fatal(err)
	}
	if rez := a.AplicaPropunereAI(CerereAplicaPropunereAI{An: 2026, Luna: 10, Propunere: p}); !rez.Succes {
		t.Fatalf("aplicare: %+v", rez.Eroare)
	}
	schita, err := datebaza.Nou(a.db).CitesteSchita(context.Background(), 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	pe := map[string]int{}
	for _, c := range schita.Celule {
		if c.Valoare != "" { // pontajul păstrează valoarea afișată a schimbului, nu codul
			pe[c.AngajatID]++
			if d, _ := time.Parse(time.DateOnly, c.Data); d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				t.Fatalf("weekend completat: %s", c.Data)
			}
			if c.AngajatID == idCuConcediu && c.Data == "2026-10-05" {
				t.Fatal("ziua de concediu a fost suprascrisă")
			}
		}
	}
	// Octombrie 2026 are 22 de zile lucrătoare; 10 angajați × 22 − 1 zi de concediu.
	if len(pe) != activi || pe[idCuConcediu] != 21 {
		t.Fatalf("angajați completați: %d din %d; zile pentru cel cu concediu: %d", len(pe), activi, pe[idCuConcediu])
	}
}
