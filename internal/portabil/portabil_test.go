package portabil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/i18n"
)

func TestPregatesteIn(t *testing.T) {
	dir := t.TempDir()
	f, err := PregatesteIn(dir)
	if err != nil {
		t.Fatal(err)
	}
	asteptat := Foldere{
		Radacina: filepath.Join(dir, "karios"),
		WebView:  filepath.Join(dir, "karios", "webview"),
		Backup:   filepath.Join(dir, "karios", "backup"),
		AI:       filepath.Join(dir, "ai"),
	}
	if f != asteptat {
		t.Errorf("obținut %+v, așteptat %+v", f, asteptat)
	}
	for _, d := range []string{f.Radacina, f.WebView, f.Backup, f.AI} {
		if info, err := os.Stat(d); err != nil || !info.IsDir() {
			t.Errorf("%s nu a fost creat: %v", d, err)
		}
	}

	// A doua pornire găsește folderele existente și nu dă eroare.
	if _, err := PregatesteIn(dir); err != nil {
		t.Errorf("a doua pregătire: %v", err)
	}
}

func TestPregatesteInEroare(t *testing.T) {
	// Un fișier cu numele folderului de date împiedică crearea lui.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, NumeFolderDate), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := PregatesteIn(dir)
	var e *i18n.Eroare
	if !errors.As(err, &e) || e.Cod != "portabil.fara_scriere" {
		t.Fatalf("aștepta eroarea portabil.fara_scriere, a primit %v", err)
	}
	if e.Param["folder"] != filepath.Join(dir, NumeFolderDate) {
		t.Errorf("parametrul folder: %q", e.Param["folder"])
	}
	for _, l := range i18n.Limbi() {
		if e.Text(l) == e.Cod {
			t.Errorf("%s: mesajul nu are traducere", l)
		}
	}
}

func TestPregatesteInFaraFisiereRamase(t *testing.T) {
	dir := t.TempDir()
	f, err := PregatesteIn(dir)
	if err != nil {
		t.Fatal(err)
	}
	intrari, err := os.ReadDir(f.Radacina)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range intrari {
		if !i.IsDir() {
			t.Errorf("fișierul de verificare a rămas: %s", i.Name())
		}
	}
}

// D091: un karios/ai existent se mută lângă karios/, cu tot conținutul, o singură dată.
func TestFolderAIMutatLangaKarios(t *testing.T) {
	dir := t.TempDir()
	vechi := filepath.Join(dir, "karios", "ai")
	if err := os.MkdirAll(vechi, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"llama-server.exe", "model.gguf", "fara-autostart"} {
		if err := os.WriteFile(filepath.Join(vechi, f), []byte(f), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	f, err := PregatesteIn(dir)
	if err != nil {
		t.Fatal(err)
	}
	if f.AI != filepath.Join(dir, "ai") {
		t.Fatalf("AI: %s", f.AI)
	}
	for _, n := range []string{"llama-server.exe", "model.gguf", "fara-autostart"} {
		if b, err := os.ReadFile(filepath.Join(f.AI, n)); err != nil || string(b) != n {
			t.Errorf("%s nu a fost mutat: %v", n, err)
		}
	}
	if _, err := os.Stat(vechi); !os.IsNotExist(err) {
		t.Errorf("karios/ai a ramas: %v", err)
	}
	// Datele raman in karios/.
	if !esteFolder(f.Radacina) || !esteFolder(f.Backup) {
		t.Error("folderul de date lipseste")
	}
	if _, err := PregatesteIn(dir); err != nil {
		t.Fatalf("a doua pornire: %v", err)
	}
}

// Cand ambele exista, nu se amesteca nimic: se foloseste ai/, iar karios/ai ramane neatins.
func TestFolderAIAmbeleExista(t *testing.T) {
	dir := t.TempDir()
	vechi, nou := filepath.Join(dir, "karios", "ai"), filepath.Join(dir, "ai")
	for _, d := range []string{vechi, nou} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(vechi, "vechi.gguf"), []byte("v"), 0o600); err != nil {
		t.Fatal(err)
	}
	if MutaFolderAIVechi(dir) {
		t.Fatal("nu trebuia mutat peste un ai/ existent")
	}
	f, err := PregatesteIn(dir)
	if err != nil || f.AI != nou {
		t.Fatalf("AI %s, %v", f.AI, err)
	}
	if _, err := os.Stat(filepath.Join(vechi, "vechi.gguf")); err != nil {
		t.Errorf("continutul vechi a fost atins: %v", err)
	}
}

// Daca mutarea nu se poate face (aici: un fisier cu numele ai), nu se pierde nimic,
// iar aplicatia foloseste in continuare karios/ai.
func TestFolderAIMutareAmanata(t *testing.T) {
	dir := t.TempDir()
	vechi := filepath.Join(dir, "karios", "ai")
	if err := os.MkdirAll(vechi, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ai"), []byte("blocheaza"), 0o600); err != nil {
		t.Fatal(err)
	}
	if MutaFolderAIVechi(dir) {
		t.Fatal("mutarea trebuia amanata")
	}
	if c := CaleAI(dir); c != vechi {
		t.Fatalf("CaleAI %s, astept %s", c, vechi)
	}
}
