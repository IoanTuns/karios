package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// I080: consumul de context vine doar din câmpul usage raportat de motor.
func TestUtilizareRaportataDeMotor(t *testing.T) {
	ultimaUtilizare.Store(nil)
	raspuns := func(usage any) *Client {
		client := NouClientCuURL("http://localhost")
		client.httpClient = &http.Client{Transport: &mockTransport{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			corp := map[string]any{"choices": []any{map[string]any{"message": map[string]any{"role": "assistant", "content": "ok"}}}}
			if usage != nil {
				corp["usage"] = usage
			}
			json.NewEncoder(w).Encode(corp)
		})}}
		return client
	}
	cere := func(c *Client) {
		t.Helper()
		if _, err := c.Interogheaza(context.Background(), []MesajConversatie{{Rol: "user", Continut: "salut"}}, ContextPontix{An: 2026, Luna: 10}); err != nil {
			t.Fatal(err)
		}
	}
	cere(raspuns(nil))
	if u := UltimaUtilizare(); u != nil {
		t.Fatalf("fara usage nu se inventeaza consum: %+v", u)
	}
	cere(raspuns(map[string]int{"prompt_tokens": 2100, "completion_tokens": 40, "total_tokens": 2140}))
	u := UltimaUtilizare()
	if u == nil || u.Prompt != 2100 || u.Raspuns != 40 || u.Total != 2140 || u.Limita != ContextConfigurat {
		t.Fatalf("consum: %+v", u)
	}
	// Un răspuns ulterior fără usage nu șterge și nu modifică ultima valoare raportată.
	cere(raspuns(nil))
	if v := UltimaUtilizare(); v == nil || *v != *u {
		t.Fatalf("consum modificat: %+v", v)
	}
	ultimaUtilizare.Store(nil)
}
