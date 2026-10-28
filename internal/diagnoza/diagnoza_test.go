package diagnoza

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDiagnozaFluxComplet(t *testing.T) {
	tempDir := t.TempDir()
	m := Initializeaza(tempDir)

	m.acum = func() time.Time {
		return time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)
	}

	// 1. Cronometrare interogare reușită
	var errNormal error
	func() {
		stop := Cronometreaza("SelectAngajati")
		defer stop(&errNormal)
	}()

	// 2. Cronometrare interogare cu eroare
	errEsuat := errors.New("tabela blocată")
	func() {
		stop := Cronometreaza("UpdatePontaj")
		defer stop(&errEsuat)
	}()

	// 3. Înregistrare eroare frontend
	m.InregistreazaEroareFrontend("GridCell", "TypeError: undefined is not an object", "at GridCell.render (Grid.tsx:42)")

	// 4. Panic recovery
	func() {
		defer RecupereazaPanic("TestFuncPanic")
		panic("eroare forțată de test")
	}()

	stat := m.Statistici()

	if stat.TotalInterogari != 2 {
		t.Errorf("așteptat 2 interogări, obținut: %d", stat.TotalInterogari)
	}
	if stat.NumarErori != 3 { // errEsuat + frontend error + panic
		t.Errorf("așteptat 3 erori, obținut: %d", stat.NumarErori)
	}
	if stat.UltimaInterogare != "UpdatePontaj" {
		t.Errorf("ultima operație ar trebui să fie UpdatePontaj, obținut: %s", stat.UltimaInterogare)
	}
	if !strings.Contains(stat.UltimaEroare, "PANIC") {
		t.Errorf("așteptat ca ultima eroare să conțină PANIC, obținut: %s", stat.UltimaEroare)
	}

	// Verificare fișier pe disc
	caleFisier := filepath.Join(tempDir, "jurnale", "pontix-2026-10.log")
	if _, err := os.Stat(caleFisier); err != nil {
		t.Fatalf("fișierul de log nu a fost creat la: %s", caleFisier)
	}

	linii := m.CitesteUltimeleLinii(20)
	if !strings.Contains(linii, "UpdatePontaj") {
		t.Errorf("jurnalul ar trebui să conțină UpdatePontaj")
	}
	if !strings.Contains(linii, "TypeError") {
		t.Errorf("jurnalul ar trebui să conțină TypeError")
	}
	if !strings.Contains(linii, "PANIC") {
		t.Errorf("jurnalul ar trebui să conțină PANIC")
	}
}
