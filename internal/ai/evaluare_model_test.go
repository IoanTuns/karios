package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// cazEvaluare e o comandă care ajunge la model (nu e rezolvată de Ruteaza) și ce trebuie să propună.
type cazEvaluare struct {
	mesaj  string
	tip    string
	marca  string
	data   string // pentru planifica_tura / sterge_tura
	deLa   string // pentru intervale
	panaLa string
	cod    string // codul schimbului, dacă se planifică
}

var cazuriEvaluare = []cazEvaluare{
	{mesaj: "Sterge garda lui Pop Ana din 15 octombrie", tip: "sterge_tura", marca: "M1", data: "2026-10-15"},
	{mesaj: "Anuleaza tura lui Ion Dan de pe 20 octombrie", tip: "sterge_tura", marca: "M2", data: "2026-10-20"},
	{mesaj: "Sterge programul lui Rus Ema intre 10 si 12 octombrie", tip: "sterge_interval", marca: "M3", deLa: "2026-10-10", panaLa: "2026-10-12"},
	{mesaj: "Pop Ana face garda pe 22 octombrie", tip: "planifica_tura", marca: "M1", data: "2026-10-22", cod: "G24"},
	{mesaj: "Ion Dan e de noapte pe 18 octombrie", tip: "planifica_tura", marca: "M2", data: "2026-10-18", cod: "N12"},
	{mesaj: "Rus Ema vine seara pe 14 octombrie", tip: "planifica_tura", marca: "M3", data: "2026-10-14", cod: "S7"},
	{mesaj: "Marin Daria e de garda pe 3 octombrie", tip: "planifica_tura", marca: "DEMO-001", data: "2026-10-03", cod: "G24"},
	{mesaj: "Multumesc!", tip: ""},
}

func contextEvaluare() ContextPontix {
	info := contextConcediu()
	info.Schimburi = append(info.Schimburi, InfoSchimb{Cod: "N12", Denumire: "Noapte 19–07", Ore: 12, Interval: "19:00–07:00"})
	return info
}

// rezultatCaz verifică propunerea față de așteptări și întoarce motivul eșecului ("" = corect).
func rezultatCaz(c cazEvaluare, r RaspunsAsistent) string {
	if c.tip == "" {
		if len(r.Propuneri) > 0 {
			return "a propus o acțiune la un mesaj fără cerere"
		}
		return ""
	}
	if len(r.Propuneri) == 0 {
		return "nicio propunere; mesaj: " + strings.ReplaceAll(r.Mesaj, "\n", " ")
	}
	p := r.Propuneri[0]
	text := func(k string) string { s, _ := p.Parametri[k].(string); return s }
	switch {
	case p.Tip != c.tip:
		return "unealta " + p.Tip
	case !strings.EqualFold(text("marca"), c.marca):
		return "marca " + text("marca")
	case c.data != "" && text("data") != c.data:
		return "data " + text("data")
	case c.deLa != "" && (text("data_inceput") != c.deLa || text("data_sfarsit") != c.panaLa):
		return "intervalul " + text("data_inceput") + "…" + text("data_sfarsit")
	case c.cod != "" && !strings.EqualFold(text("cod_schimb"), c.cod):
		return "schimbul " + text("cod_schimb")
	}
	return ""
}

// TestEvaluareModele compară modelele locale pe comenzile care ajung la model (F14-T03).
// Rulează doar la cerere, cu serverul llama.cpp și modelele date:
//
//	PONTIX_EVAL_SERVER=…/llama-server.exe PONTIX_EVAL_MODELE="a.gguf;b.gguf" go test -run EvaluareModele -v -timeout 60m ./internal/ai/
//
// Fișierele se leagă (hard link) într-un folder temporar de lângă model, ca managerul aplicației
// să le pornească exact ca în aplicație; de aceea serverul și modelele trebuie să fie pe același disc.
func TestEvaluareModele(t *testing.T) {
	server, modele := os.Getenv("PONTIX_EVAL_SERVER"), os.Getenv("PONTIX_EVAL_MODELE")
	if server == "" || modele == "" {
		t.Skip("evaluare la cerere: PONTIX_EVAL_SERVER și PONTIX_EVAL_MODELE")
	}
	info := contextEvaluare()
	for _, c := range cazuriEvaluare {
		if d := Ruteaza(c.mesaj, info); d.Ruta != RutaModel {
			t.Fatalf("%q nu ajunge la model (ruta %s); cazul nu evaluează modelul", c.mesaj, d.Ruta)
		}
	}
	for _, model := range strings.Split(modele, ";") {
		model = strings.TrimSpace(model)
		nume := strings.TrimSuffix(filepath.Base(model), filepath.Ext(model))
		dir := filepath.Join(filepath.Dir(model), "evaluare-"+nume)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		for sursa, tinta := range map[string]string{server: NumeBinarServer(), model: "model.gguf"} {
			_ = os.Remove(filepath.Join(dir, tinta))
			if err := os.Link(sursa, filepath.Join(dir, tinta)); err != nil {
				t.Fatalf("legătura %s: %v", tinta, err)
			}
		}
		m := NouManager(dir, 0)
		start := time.Now()
		if err := m.Porneste(context.Background(), ""); err != nil {
			t.Fatalf("%s: pornire: %v", nume, err)
		}
		pornire := time.Since(start)
		corecte, total := 0, time.Duration(0)
		var linii []string
		for _, c := range cazuriEvaluare {
			ctx, anuleaza := context.WithTimeout(context.Background(), 3*time.Minute)
			start := time.Now()
			r, err := m.Client().Interogheaza(ctx, []MesajConversatie{{Rol: "user", Continut: c.mesaj}}, ContextPentruCerere(c.mesaj, info))
			d := time.Since(start)
			anuleaza()
			total += d
			motiv := ""
			if err != nil {
				motiv = "eroare: " + err.Error()
			} else {
				motiv = rezultatCaz(c, r)
			}
			stare := "corect"
			if motiv == "" {
				corecte++
			} else {
				stare = "greșit: " + motiv
			}
			linii = append(linii, fmt.Sprintf("  %5.1fs  %-55q %s", d.Seconds(), c.mesaj, stare))
		}
		_ = m.Opreste()
		_ = os.RemoveAll(dir)
		t.Logf("%s: %d/%d corecte · pornire %.1fs · răspuns mediu %.1fs\n%s", nume, corecte, len(cazuriEvaluare),
			pornire.Seconds(), total.Seconds()/float64(len(cazuriEvaluare)), strings.Join(linii, "\n"))
	}
}
