package main

import (
	"context"
	"os"
	"path/filepath"
	"pontix/internal/ai"
	"pontix/internal/versiune"
	"strings"
	"testing"
)

func TestAjutorFaraBazaSauMotor(t *testing.T) {
	a := NewApp(nil)
	a.SeteazaLimbaInterfata("en")
	r := a.interogheazaAI(context.Background(), CerereInterogareAI{Istoric: []ai.MesajConversatie{{Rol: "user", Continut: "How do I record leave?"}}})
	if r.Eroare != nil || r.Raspuns == nil || len(r.Raspuns.ArticoleAjutor) == 0 || a.aiManager != nil {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(r.Raspuns.Mesaj, "Record a leave request") {
		t.Fatal(r.Raspuns)
	}
}
func TestManualLocalNumaiVersiuneaCurentaSiIDValid(t *testing.T) {
	dir := t.TempDir()
	radacina := filepath.Join(dir, "karios")
	if _, e := caleManualAjutor(radacina, "../../alta"); e == nil || e.Cod != "ajutor.articol_lipsa" {
		t.Fatal(e)
	}
	if err := os.WriteFile(filepath.Join(dir, "manual-utilizare-ro-0.0.0.html"), []byte("fictiv"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, e := caleManualAjutor(radacina, "absente.inregistreaza"); e == nil || e.Cod != "ajutor.manual_lipsa" {
		t.Fatal(e)
	}
	if err := os.WriteFile(filepath.Join(dir, "manual-utilizare-ro-"+versiune.Aplicatie+".html"), []byte("fictiv"), 0600); err != nil {
		t.Fatal(err)
	}
	cale, e := caleManualAjutor(radacina, "absente.inregistreaza")
	if e != nil || !strings.HasPrefix(cale, "file:") || !strings.HasSuffix(cale, "#ajutor-absente.inregistreaza") {
		t.Fatal(cale, e)
	}
}
