package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type mockTransport struct {
	handler http.Handler
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	m.handler.ServeHTTP(rec, req)
	return rec.Result(), nil
}

func TestDetecteazaMediu(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Folder gol -> neinstalat
	stare := DetecteazaMediu(tempDir, 28499)
	if stare.Instalat || stare.Activ {
		t.Fatalf("așteptat neinstalat și inactiv, obținut: %+v", stare)
	}

	// 2. Doar executabil -> neinstalat (lipsește modelul)
	binarPath := filepath.Join(tempDir, NumeBinarServer())
	if err := os.WriteFile(binarPath, []byte("fake-binar"), 0o755); err != nil {
		t.Fatal(err)
	}
	stare = DetecteazaMediu(tempDir, 28499)
	if stare.Instalat || stare.CaleBinar == "" || len(stare.ModeleDisponibile) != 0 {
		t.Fatalf("așteptat binar prezent dar model lipsă, obținut: %+v", stare)
	}

	// 3. Executabil + Model GGUF -> instalat
	modelPath := filepath.Join(tempDir, "qwen2.5-3b.gguf")
	if err := os.WriteFile(modelPath, []byte("fake-model"), 0o644); err != nil {
		t.Fatal(err)
	}
	stare = DetecteazaMediu(tempDir, 28499)
	if !stare.Instalat || stare.CaleBinar == "" || stare.CaleModel == "" || len(stare.ModeleDisponibile) != 1 {
		t.Fatalf("așteptat instalat cu model disponibil, obținut: %+v", stare)
	}
}

func TestInterogheazaAsistentConversatie(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/chat/completions" {
			w.Header().Set("Content-Type", "application/json")
			raspuns := map[string]any{
				"choices": []map[string]any{
					{
						"message": map[string]any{
							"role":    "assistant",
							"content": "Bună ziua! Cu ce vă pot ajuta la pontaj astăzi?",
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(raspuns)
			return
		}
		http.NotFound(w, r)
	})

	client := NouClientCuURL("http://in-memory-test").CuTransport(&mockTransport{handler: handler})
	ctx := context.Background()

	istoric := []MesajConversatie{
		{Rol: "user", Continut: "Salut!"},
	}
	info := ContextPontix{
		Limba: "ro",
		An:    2026,
		Luna:  10,
	}

	rez, err := client.Interogheaza(ctx, istoric, info)
	if err != nil {
		t.Fatalf("eroare interogare asistent: %v", err)
	}

	if rez.Mesaj != "Bună ziua! Cu ce vă pot ajuta la pontaj astăzi?" {
		t.Fatalf("răspuns neașteptat: %s", rez.Mesaj)
	}
	if len(rez.Propuneri) != 0 {
		t.Fatalf("nu erau așteptate propuneri, obținut: %d", len(rez.Propuneri))
	}
}

func TestInterogheazaAsistentFunctionCalling(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/chat/completions" {
			w.Header().Set("Content-Type", "application/json")
			raspuns := map[string]any{
				"choices": []map[string]any{
					{
						"message": map[string]any{
							"role":    "assistant",
							"content": "Am identificat solicitarea și am pregătit acțiunea de planificare.",
							"tool_calls": []map[string]any{
								{
									"id":   "call_12345",
									"type": "function",
									"function": map[string]any{
										"name":      "planifica_tura",
										"arguments": `{"marca":"DEMO-001","data":"2026-10-15","cod_schimb":"G24"}`,
									},
								},
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(raspuns)
			return
		}
		http.NotFound(w, r)
	})

	client := NouClientCuURL("http://in-memory-test").CuTransport(&mockTransport{handler: handler})
	ctx := context.Background()

	istoric := []MesajConversatie{
		{Rol: "user", Continut: "Pune-l pe Popescu în gardă pe 15 octombrie"},
	}
	info := ContextPontix{
		Limba: "ro",
		An:    2026,
		Luna:  10,
		Angajati: []InfoAngajat{
			{ID: "1", Nume: "Popescu", Prenume: "Ion", Marca: "DEMO-001", Functie: "Medic"},
		},
		Schimburi: []InfoSchimb{
			{Cod: "G24", Denumire: "Gardă 24h", Ore: 24},
		},
	}

	rez, err := client.Interogheaza(ctx, istoric, info)
	if err != nil {
		t.Fatalf("eroare interogare asistent: %v", err)
	}

	if len(rez.Propuneri) != 1 {
		t.Fatalf("așteptată o propunere, obținut %d", len(rez.Propuneri))
	}

	p := rez.Propuneri[0]
	if p.Tip != "planifica_tura" || !p.NecesitaConfirmare {
		t.Fatalf("propunere invalidă: %+v", p)
	}
	if p.Parametri["marca"] != "DEMO-001" || p.Parametri["data"] != "2026-10-15" || p.Parametri["cod_schimb"] != "G24" {
		t.Fatalf("parametri neașteptați: %+v", p.Parametri)
	}
}
