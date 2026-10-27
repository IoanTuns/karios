package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppAutoPornireAI(t *testing.T) {
	tempDir := t.TempDir()
	a := &App{radacina: tempDir}

	folder := a.folderAI()
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}

	// Implicit trebuie să fie activă (fișierul de dezactivare nu există)
	if !a.EsteAutoPornireAI() {
		t.Fatalf("așteptat EsteAutoPornireAI() true inițial")
	}

	// Dezactivare
	if a.SeteazaAutoPornireAI(false) != false {
		t.Fatalf("așteptat false după SeteazaAutoPornireAI(false)")
	}
	if a.EsteAutoPornireAI() {
		t.Fatalf("așteptat EsteAutoPornireAI() false după dezactivare")
	}
	if _, err := os.Stat(filepath.Join(folder, fisierFaraAutoPornire)); err != nil {
		t.Fatalf("fișierul de dezactivare nu a fost creat: %v", err)
	}

	// Reactivare
	if a.SeteazaAutoPornireAI(true) != true {
		t.Fatalf("așteptat true după SeteazaAutoPornireAI(true)")
	}
	if !a.EsteAutoPornireAI() {
		t.Fatalf("așteptat EsteAutoPornireAI() true după reactivare")
	}

	// Verificare StareAI pe folder gol
	stare := a.StareAI()
	if stare.Stare.Instalat {
		t.Fatalf("folderul gol nu trebuia să fie marcat ca instalat")
	}
}
