package ai

import (
	"strings"
	"testing"
)

func TestIntrebariUtilizareNuDevinOperatii(t *testing.T) {
	for _, mesaj := range []string{"Cum inregistrez un concediu?", "Cum șterg o tură?", "Cum planific S7 pentru EX-001?", "How do I record leave?", "Unde gasesc backup?"} {
		d := Ruteaza(mesaj, ContextPontix{Limba: "ro", An: 2026, Luna: 10})
		if d.Ruta != RutaAjutor || d.Raspuns.Intentie != nil || len(d.Raspuns.Propuneri) != 0 || d.Raspuns.Mesaj == "" {
			t.Fatal(mesaj, d)
		}
	}
	for _, m := range []string{"Înregistrează concediu pentru EX-001", "Sterge tura pentru EX-001", "Toți angajații lucrează S7", "ajutor", "Multumesc"} {
		if _, ok := AjutorUtilizare(m, "ro"); ok {
			t.Fatal("comanda capturata ca intrebare", m)
		}
	}
}
func TestAjutorCiteazaVersiuneaSiNuInventeaza(t *testing.T) {
	r, ok := AjutorUtilizare("How do I record leave?", "en")
	if !ok || len(r.ArticoleAjutor) == 0 || r.VersiuneAjutor == "" || !strings.Contains(r.Mesaj, "Record a leave request") || !strings.Contains(r.Mesaj, "Source:") {
		t.Fatal(r)
	}
	r, ok = AjutorUtilizare("Cum telematicxyz?", "ro")
	if !ok || len(r.ArticoleAjutor) != 0 || r.Intentie != nil || len(r.Propuneri) != 0 || !strings.Contains(r.Mesaj, "articol verificat") {
		t.Fatal(r)
	}
}

func TestManualInteractivOferaPuncteDePornire(t *testing.T) {
	r, ok := AjutorUtilizare("manual interactiv", "ro")
	if !ok || len(r.ArticoleAjutor) != 3 || r.Intentie != nil || len(r.Propuneri) != 0 {
		t.Fatal(r)
	}
}
