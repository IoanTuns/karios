// Package portabil găsește și pregătește folderul de date de lângă .exe (D002).
// Aplicația nu scrie nimic în afara acestui folder.
package portabil

import (
	"os"
	"path/filepath"
)

// NumeFolderDate este folderul de lângă .exe în care stau toate datele operaționale.
const NumeFolderDate = "karios"

// Foldere sunt căile folderului de date și ale subfolderelor lui.
type Foldere struct {
	Radacina string // karios/
	WebView  string // karios/webview/ — cache-ul WebView2
	Backup   string // karios/backup/
}

// Pregateste creează folderul de date lângă executabilul care rulează.
func Pregateste() (Foldere, error) {
	exe, err := os.Executable()
	if err != nil {
		return Foldere{}, err
	}
	if exe, err = filepath.EvalSymlinks(exe); err != nil {
		return Foldere{}, err
	}
	return PregatesteIn(filepath.Dir(exe))
}

// PregatesteIn creează folderul de date în dir și subfolderele lui.
func PregatesteIn(dir string) (Foldere, error) {
	radacina := filepath.Join(dir, NumeFolderDate)
	f := Foldere{
		Radacina: radacina,
		WebView:  filepath.Join(radacina, "webview"),
		Backup:   filepath.Join(radacina, "backup"),
	}
	for _, d := range []string{f.Radacina, f.WebView, f.Backup} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return Foldere{}, err
		}
	}
	return f, nil
}
