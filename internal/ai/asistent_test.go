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
	if p.Descriere != "Planificare schimb G24 – Gardă 24h pentru Popescu Ion (marca DEMO-001) pe 2026-10-15" {
		t.Fatalf("descriere neașteptată: %s", p.Descriere)
	}
}

func TestGasesteAngajatInfoSiDescriere(t *testing.T) {
	angajati := []InfoAngajat{
		{ID: "1", Nume: "Dobre", Prenume: "Rares", Marca: "DEMO-008", Functie: "Operator"},
	}
	schimburi := []InfoSchimb{
		{Cod: "Z7", Denumire: "Zi 07–14", Ore: 7, Interval: "07:00–14:00"},
	}
	info := ContextPontix{
		Angajati:  angajati,
		Schimburi: schimburi,
	}

	// 1. Verificare căutare cu variantă fără cratimă ("DEM008")
	a, gasit := gasesteAngajatInfo("DEM008", angajati)
	if !gasit || a.Marca != "DEMO-008" {
		t.Fatalf("eșec identificare marcă DEM008: %+v", a)
	}

	// 2. Verificare căutare după nume parțial ("Dobre Rares")
	a, gasit = gasesteAngajatInfo("Dobre Rares", angajati)
	if !gasit || a.Marca != "DEMO-008" {
		t.Fatalf("eșec identificare nume: %+v", a)
	}

	// 3. Verificare descriere completă planifica_interval
	args := map[string]any{
		"marca":           "DEM008",
		"cod_schimb":      "Z7",
		"data_inceput":    "2026-10-01",
		"data_sfarsit":    "2026-10-31",
		"doar_lucratoare": true,
	}
	desc := construiesteDescrierePropunere("planifica_interval", args, info)
	asteptat := "Planificare schimb Z7 – Zi 07–14 (07:00–14:00, 7 ore) pentru Dobre Rares (marca DEMO-008) în intervalul 2026-10-01 – 2026-10-31 (doar zile lucrătoare L–V)"
	if desc != asteptat {
		t.Fatalf("descriere generată incorect:\nobținut:  %s\nașteptat: %s", desc, asteptat)
	}
}

func TestCorecteazaSchimbDacaEroare(t *testing.T) {
	schimburi := []InfoSchimb{
		{Cod: "Z7", Denumire: "Zi 07–14", Ore: 7, Interval: "07:00–14:00"},
		{Cod: "S7", Denumire: "Seară 15–22", Ore: 7, Interval: "15:00–22:00"},
	}

	// Cazul utilizatorului: utilizatorul cere "7 la 14", dar modelul returnase greșit "S7"
	istoric := []MesajConversatie{
		{Rol: "user", Continut: "Dobre rares e toata luna 7 la 14"},
	}
	argsGresite := map[string]any{
		"marca":        "DEMO-008",
		"cod_schimb":   "S7",
		"data_inceput": "2026-10-01",
		"data_sfarsit": "2026-10-31",
	}

	info := ContextPontix{
		An:        2026,
		Luna:      10,
		Schimburi: schimburi,
	}

	corectat := corecteazaSchimbDacaEroare(argsGresite, istoric, schimburi, info)
	if corectat["cod_schimb"] != "Z7" {
		t.Fatalf("așteptat cod Z7 pentru 7 la 14, obținut: %v", corectat["cod_schimb"])
	}
	if corectat["doar_lucratoare"] != true {
		t.Fatalf("așteptat doar_lucratoare true pentru toata luna, obținut: %v", corectat["doar_lucratoare"])
	}

	// Cazul weekend alternat: "Dumitru andrei pe octombrie lucreaza tot al doilea weekend pe zi, schimb de 12 ore, incepand cu data de 3"
	istoricWeekend := []MesajConversatie{
		{Rol: "user", Continut: "Dumitru andrei pe octombrie lucreaza tot al doilea weekend pe zi, schimb de 12 ore, incepand cu data de 3"},
	}
	argsWeekend := map[string]any{
		"marca":      "DEMO-002",
		"cod_schimb": "Z7",
	}
	corectatWeekend := corecteazaSchimbDacaEroare(argsWeekend, istoricWeekend, schimburi, info)
	if corectatWeekend["cod_schimb"] != "Z12" {
		t.Fatalf("așteptat cod Z12 pentru 12 ore pe zi, obținut: %v", corectatWeekend["cod_schimb"])
	}
	if corectatWeekend["filtru_zile"] != "tot_al_doilea_weekend" {
		t.Fatalf("așteptat filtru_zile tot_al_doilea_weekend, obținut: %v", corectatWeekend["filtru_zile"])
	}
	if corectatWeekend["doar_lucratoare"] != false {
		t.Fatalf("așteptat doar_lucratoare false pentru weekend, obținut: %v", corectatWeekend["doar_lucratoare"])
	}
	if corectatWeekend["data_inceput"] != "2026-10-03" {
		t.Fatalf("așteptat data_inceput 2026-10-03, obținut: %v", corectatWeekend["data_inceput"])
	}
}

func TestEsteCerereToataLuna(t *testing.T) {
	teste := []struct {
		mesaj     string
		toataLuna bool
		nrLuna    int
	}{
		{"Dobre in octombrie e pe zi de la 7 la 14", true, 10},
		{"Dobre rares e toata luna 7 la 14", true, 0},
		{"Dobre rares e tata luna 7 la 14", true, 0},
		{"Pune-l pe Popescu in garda pe 15 octombrie", false, 10},
		{"Pune-l pe Popescu in garda pe data de 20 noiembrie", false, 11},
		{"Cine e liber in acest weekend?", false, 0},
	}

	for _, tt := range teste {
		tl, luna := esteCerereToataLuna(tt.mesaj)
		if tl != tt.toataLuna || luna != tt.nrLuna {
			t.Errorf("pentru %q: asteptat (%v, %d), obtinut (%v, %d)", tt.mesaj, tt.toataLuna, tt.nrLuna, tl, luna)
		}
	}
}
