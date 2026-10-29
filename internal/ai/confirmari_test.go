package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfirmareFalsaEsteInlocuita(t *testing.T) {
	raspuns := "Planificare interval Z7 – Zi 07–14 (07:00–14:00, 7 ore) pentru toți angajații (marca) în intervalul 2026-10-01 – 2026-10-31 (zile lucrătoare L–V) (deja aplicat)"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]any{"role": "assistant", "content": raspuns}}}})
	}))
	defer srv.Close()
	info := contextGrup()
	info.Angajati = nil // ca textul să nu poată fi transformat într-o propunere
	r, err := NouClientCuURL(srv.URL).Interogheaza(context.Background(), []MesajConversatie{{Rol: "user", Continut: "aplica in pontaj"}}, info)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Propuneri) != 0 || r.Mesaj != MesajFaraActiune {
		t.Fatalf("răspuns: %+v", r)
	}
}

func TestEsteConfirmareFalsa(t *testing.T) {
	for _, m := range []string{"Am planificat tura pentru Pop.", "Garda a fost planificată pe 15.", "Acțiuni propuse: Planificare schimb Z7"} {
		if !esteConfirmareFalsa(m) {
			t.Errorf("nerecunoscut: %q", m)
		}
	}
	for _, m := range []string{"", "Cine este liber în weekend? Iată lista:", "Nu știu ce schimb înseamnă „de la 9 la 17”."} {
		if esteConfirmareFalsa(m) {
			t.Errorf("recunoscut greșit: %q", m)
		}
	}
}
