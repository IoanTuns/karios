package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"
	"time"

	"pontix/internal/ai"
	"pontix/internal/bd"
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
