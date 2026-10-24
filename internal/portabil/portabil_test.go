package portabil

import (
	"os"
	"path/filepath"
	"testing"
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
	}
	if f != asteptat {
		t.Errorf("obținut %+v, așteptat %+v", f, asteptat)
	}
	for _, d := range []string{f.Radacina, f.WebView, f.Backup} {
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
	if _, err := PregatesteIn(dir); err == nil {
		t.Error("aștepta eroare când folderul de date nu poate fi creat")
	}
}
