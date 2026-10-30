package ai

import (
	"strings"
	"testing"
)

func TestActiunileAnterioareIntraInDateleDeLucru(t *testing.T) {
	info := contextGrup()
	info.ActiuniAnterioare = []string{"[aplicată] Planificare schimb Z7 pentru 3 angajați"}
	p := promptDate(info)
	if !strings.Contains(p, "Acțiuni din această conversație") || !strings.Contains(p, "- [aplicată] Planificare schimb Z7") {
		t.Fatalf("lista de acțiuni lipsește din prompt:\n%s", p)
	}
	if !strings.Contains(promptInstructiuni(), "NU scrie în text că ai planificat") {
		t.Fatal("regula împotriva confirmărilor în text lipsește")
	}
	if strings.Contains(promptDate(contextGrup()), "Acțiuni din această conversație") {
		t.Fatal("secțiunea apare și fără acțiuni")
	}
}

// Începutul cererii (instrucțiuni și istoric) nu depinde de date, ca serverul să-l reia din cache;
// datele de lucru stau într-un mesaj de sistem chiar înaintea ultimei cereri (F14-T29).
func TestDateleDeLucruStauInainteaUltimeiCereri(t *testing.T) {
	istoric := []MesajConversatie{{Rol: "user", Continut: "prima"}, {Rol: "assistant", Continut: "răspuns"}, {Rol: "user", Continut: "a doua"}}
	a, b := contextGrup(), contextGrup()
	b.An++
	ma, mb := mesajeCerere(istoric, a), mesajeCerere(istoric, b)
	if len(ma) != 5 || ma[0]["role"] != "system" || ma[3]["role"] != "system" {
		t.Fatalf("mesaje: %v", ma)
	}
	for x := 0; x < 3; x++ {
		if ma[x]["content"] != mb[x]["content"] {
			t.Fatalf("mesajul %d depinde de date: %q", x, ma[x]["content"])
		}
	}
	if strings.Contains(ma[0]["content"].(string), "Anul curent") || ma[1]["content"] != "prima" {
		t.Fatal("datele au intrat în instrucțiuni sau în istoric")
	}
	if !strings.Contains(ma[3]["content"].(string), "Anul curent de lucru") || ma[4]["content"] != "a doua" {
		t.Fatalf("datele și cererea: %v", ma[3:])
	}
}
