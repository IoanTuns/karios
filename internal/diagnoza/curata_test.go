package diagnoza

import (
	"strings"
	"testing"
)

func TestCurataScoateDatelePersonale(t *testing.T) {
	text := Curata(`angajatul "Pop Ana" (marca 'M-17') nu există; CNP 1850101123456, tel 0722123456, ana.pop@spital.ro, „Rus Ema”`)
	for _, interzis := range []string{"Pop Ana", "M-17", "1850101123456", "0722123456", "ana.pop", "spital.ro", "Rus Ema"} {
		if strings.Contains(text, interzis) {
			t.Fatalf("a rămas %q în %q", interzis, text)
		}
	}
	if !strings.Contains(text, "nu există") || !strings.Contains(text, `"…"`) {
		t.Fatalf("contextul erorii s-a pierdut: %q", text)
	}
	// Datele, orele și locul din cod rămân.
	if t2 := Curata("eroare la 2026-10-15 12:30 în Grid.tsx:42:7"); t2 != "eroare la 2026-10-15 12:30 în Grid.tsx:42:7" {
		t.Fatalf("text util modificat: %q", t2)
	}
	if n := len([]rune(Curata(strings.Repeat("x", 10000)))); n > LungimeMaximaIntrare+10 {
		t.Fatalf("intrare nelimitată: %d", n)
	}
}

func TestJurnalulNuPastreazaNumele(t *testing.T) {
	m := Initializeaza(t.TempDir())
	m.InregistreazaEroareFrontend("Grila", `TypeError: "Popescu Ion" lipsește`, `at f ("Popescu Ion")`)
	linii, err := m.UltimeleLinii(10)
	if err != nil || strings.Contains(linii, "Popescu") || strings.Contains(m.Statistici().UltimaEroare, "Popescu") {
		t.Fatalf("numele a ajuns în jurnal: %q %v", linii, err)
	}
}
