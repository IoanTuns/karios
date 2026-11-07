// Package portabil găsește și pregătește folderele de lângă .exe (D002, D091):
// karios/ pentru datele operaționale și ai/ pentru motorul și modelele LLM.
// Aplicația nu scrie nimic în afara acestor două foldere.
package portabil

import (
	"os"
	"path/filepath"

	"pontix/internal/i18n"
)

// NumeFolderDate este folderul de lângă .exe în care stau toate datele operaționale.
const NumeFolderDate = "karios"

// NumeFolderAI este folderul de lângă .exe, la același nivel cu karios/, pentru motorul
// și modelele LLM locale (D091): baza de date și LLM-ul stau separat.
const NumeFolderAI = "ai"

// Foldere sunt căile folderului de date și ale subfolderelor lui.
type Foldere struct {
	Radacina string // karios/
	WebView  string // karios/webview/ — cache-ul WebView2
	Backup   string // karios/backup/
	AI       string // ai/, lângă karios/ — motorul și modelele LLM locale (F14, D042, D091)
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

// PregatesteIn creează folderul de date în dir, cu subfolderele lui, și verifică
// faptul că aplicația poate scrie în el. Dacă nu poate, întoarce *i18n.Eroare
// cu codul "portabil.fara_scriere" (D002 punctul 5).
func PregatesteIn(dir string) (Foldere, error) {
	radacina := filepath.Join(dir, NumeFolderDate)
	MutaFolderAIVechi(dir)
	f := Foldere{
		Radacina: radacina,
		WebView:  filepath.Join(radacina, "webview"),
		Backup:   filepath.Join(radacina, "backup"),
		AI:       CaleAI(dir),
	}
	faraScriere := &i18n.Eroare{Cod: "portabil.fara_scriere", Param: map[string]string{"folder": radacina}}

	for _, d := range []string{f.Radacina, f.WebView, f.Backup, f.AI} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return Foldere{}, faraScriere
		}
	}
	// Un folder existent poate fi totuși fără drept de scriere.
	proba, err := os.CreateTemp(f.Radacina, "verificare-*")
	if err != nil {
		return Foldere{}, faraScriere
	}
	proba.Close()
	os.Remove(proba.Name())
	return f, nil
}

// folderAIVechi este locul de dinainte de D091: karios/ai.
func folderAIVechi(dir string) string { return filepath.Join(dir, NumeFolderDate, "ai") }

func esteFolder(cale string) bool {
	info, err := os.Stat(cale)
	return err == nil && info.IsDir()
}

// MutaFolderAIVechi mută karios/ai în ai/, lângă karios/, dacă noul folder nu există încă
// (D091). Mutarea e o redenumire pe același disc, deci modelele de câțiva GB nu se copiază.
// Dacă redenumirea nu reușește (de exemplu motorul rulează din vechiul folder într-o altă
// instanță), nu se pierde nimic: CaleAI folosește vechiul folder, iar mutarea se reia la
// următoarea pornire. Când ambele există, nu se amestecă nimic; se folosește ai/.
func MutaFolderAIVechi(dir string) bool {
	vechi, nou := folderAIVechi(dir), filepath.Join(dir, NumeFolderAI)
	if !esteFolder(vechi) {
		return false
	}
	if _, err := os.Stat(nou); err == nil || !os.IsNotExist(err) {
		return false
	}
	return os.Rename(vechi, nou) == nil
}

// CaleAI întoarce folderul AI folosit: ai/ dacă există, altfel karios/ai dacă a rămas
// acolo (mutare amânată), altfel ai/ (care se creează).
func CaleAI(dir string) string {
	nou := filepath.Join(dir, NumeFolderAI)
	if !esteFolder(nou) && esteFolder(folderAIVechi(dir)) {
		return folderAIVechi(dir)
	}
	return nou
}
