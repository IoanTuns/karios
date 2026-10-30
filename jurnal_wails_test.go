package main

import (
	"strings"
	"testing"
)

func TestJurnalWailsFaraArgumente(t *testing.T) {
	sursa, text := descrieMesajWails(`process message error: C{"name":"main.App.SalveazaAngajat","args":[{"nume":"Pop","prenume":"Ana","cnp":"2850101123456"}],"callbackID":"x"} -> runtime error: index out of range`)
	if sursa != "wails/main.App.SalveazaAngajat" || text != "runtime error: index out of range" {
		t.Fatalf("%q %q", sursa, text)
	}
	if strings.Contains(sursa+text, "Pop") || strings.Contains(sursa+text, "2850101123456") {
		t.Fatal("argumentele au ajuns în jurnal")
	}
	if s, x := descrieMesajWails(`process message error: Cnu e json -> panică`); s != "wails" || x != "panică" {
		t.Fatalf("mesaj nerecunoscut: %q %q", s, x)
	}
	if s, x := descrieMesajWails("alt mesaj"); s != "wails" || x != "alt mesaj" {
		t.Fatalf("alt mesaj: %q %q", s, x)
	}
}
