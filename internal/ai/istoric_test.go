package ai

import (
	"strings"
	"testing"
)

func TestActiunileAnterioareIntraInPromptulDeSistem(t *testing.T) {
	info := contextGrup()
	info.ActiuniAnterioare = []string{"[aplicată] Planificare schimb Z7 pentru 3 angajați"}
	p := construiestePromptSistem(info)
	if !strings.Contains(p, "Acțiuni din această conversație") || !strings.Contains(p, "- [aplicată] Planificare schimb Z7") {
		t.Fatalf("lista de acțiuni lipsește din prompt:\n%s", p)
	}
	if !strings.Contains(p, "NU scrie în text că ai planificat") {
		t.Fatal("regula împotriva confirmărilor în text lipsește")
	}
	if strings.Contains(construiestePromptSistem(contextGrup()), "Acțiuni din această conversație") {
		t.Fatal("secțiunea apare și fără acțiuni")
	}
}
