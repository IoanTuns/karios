package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/i18n"
)

// F03-T11 (D106): fisierul de lucru, verificarea, importul si stergerea traducerilor prin App,
// doar pentru administrator; textele importate se vad imediat in i18n.
func TestTraduceriDinFolder(t *testing.T) {
	defer i18n.SeteazaCorecturi(nil)
	a := appCuConturi(t)
	if r := a.ImportaFolderTraduceri(); r.Eroare == nil || r.Eroare.Cod != "cont.doar_administrator" {
		t.Fatalf("fara cont: %+v", r.Eroare)
	}
	if r := a.CreeazaPrimulAdministrator(cerereCont("ana")); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	if r := a.VerificaFolderTraduceri(); r.Eroare == nil || r.Eroare.Cod != "traduceri.folder_gol" {
		t.Fatalf("folder gol: %+v", r.Eroare)
	}
	r := a.ScrieFisierTraduceri("en")
	if r.Eroare != nil || filepath.Base(r.Cale) != "en.json" {
		t.Fatalf("fisier de lucru: %+v", r)
	}
	brut, err := os.ReadFile(r.Cale)
	if err != nil {
		t.Fatal(err)
	}
	var texte map[string]string
	if err := json.Unmarshal(brut, &texte); err != nil || texte["ui.salveaza"] != i18n.TextInclus(i18n.Engleza, "ui.salveaza") {
		t.Fatal("continut", err)
	}
	texte["ui.salveaza"] = "Store"
	brut, _ = json.Marshal(texte)
	if err := os.WriteFile(r.Cale, brut, 0o600); err != nil {
		t.Fatal(err)
	}
	if v := a.VerificaFolderTraduceri(); v.Eroare != nil || len(v.Rezumate) != 1 || v.Rezumate[0].Modificate != 1 {
		t.Fatalf("verificare: %+v", v)
	}
	if i := a.ImportaFolderTraduceri(); i.Eroare != nil || i.Corecturi["en"] != 1 {
		t.Fatalf("import: %+v", i)
	}
	if got := i18n.T(i18n.Engleza, "ui.salveaza", nil); got != "Store" {
		t.Fatalf("textul importat nu se vede: %q", got)
	}
	// Un fisier cu nume care nu e cod de limba opreste verificarea si numeste fisierul.
	if err := os.WriteFile(filepath.Join(a.folderTraduceri(), "romana.json"), []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if v := a.VerificaFolderTraduceri(); v.Eroare == nil || v.Eroare.Cod != "traduceri.fisier_invalid" || v.Eroare.Param["fisier"] != "romana.json" {
		t.Fatalf("fisier invalid: %+v", v.Eroare)
	}
	if s := a.StergeTraduceriLimba("en"); s.Eroare != nil || len(s.Corecturi) != 0 {
		t.Fatalf("stergere: %+v", s)
	}
	if got := i18n.T(i18n.Engleza, "ui.salveaza", nil); got != i18n.TextInclus(i18n.Engleza, "ui.salveaza") {
		t.Fatalf("dupa stergere: %q", got)
	}
}
