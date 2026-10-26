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
		AI:       filepath.Join(dir, "karios", "ai"),
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
