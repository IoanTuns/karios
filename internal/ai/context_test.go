package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestReluareContextDepasit(t *testing.T) {
	for _, tip := range []string{"exceed_context_size_error", "alta_eroare"} {
		t.Run(tip, func(t *testing.T) {
			apeluri := 0
			client := NouClientCuURL("http://localhost")
			client.httpClient = &http.Client{Transport: &mockTransport{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				apeluri++
				var cerere struct {
					Mesaje []struct {
						Rol      string `json:"role"`
						Continut string `json:"content"`
					} `json:"messages"`
				}
				if err := json.NewDecoder(r.Body).Decode(&cerere); err != nil {
					t.Fatal(err)
				}
				if apeluri == 2 && (len(cerere.Mesaje) != 3 || cerere.Mesaje[0].Rol != "system" || !strings.HasPrefix(cerere.Mesaje[1].Continut, "DATE DE LUCRU") || cerere.Mesaje[2].Continut != "cererea curenta") {
					t.Fatalf("mesaje la reluare: %+v", cerere.Mesaje)
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"type": tip}})
			})}}
			_, err := client.Interogheaza(context.Background(), []MesajConversatie{{Rol: "user", Continut: "vechi"}, {Rol: "assistant", Continut: "raspuns vechi"}, {Rol: "user", Continut: "cererea curenta"}}, ContextPontix{An: 2026, Luna: 10})
			asteptate := 1
			if tip == "exceed_context_size_error" {
				asteptate = 2
			}
			if err == nil || apeluri != asteptate {
				t.Fatalf("apeluri=%d, eroare=%v", apeluri, err)
			}
		})
	}
}
