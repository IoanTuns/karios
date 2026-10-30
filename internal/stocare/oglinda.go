package stocare

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/seif"
)

// Copia seifului de lângă backup-uri (D055): fără seif, nici baza, nici backup-urile, nici
// cheia de recuperare nu mai deschid datele.
const (
	sufixDeteriorat = ".deteriorat"
	// marcajRefacere rămâne până când refacerea e notată în jurnalul de audit, chiar dacă datele
	// se deblochează abia la o pornire ulterioară.
	marcajRefacere = "seif.refacut"
)

func caleOglinda(radacina string) string { return filepath.Join(radacina, "backup", NumeSeif) }

// refaSeiful pune oglinda în locul seifului care lipsește sau e deteriorat. Seiful stricat se
// păstrează ca seif.deteriorat; fără oglindă validă nu se schimbă nimic.
func refaSeiful(radacina string) (bool, error) {
	data, err := os.ReadFile(caleOglinda(radacina))
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if _, err := seif.Incarca(caleOglinda(radacina)); err != nil {
		return false, nil // nici oglinda nu e validă: rămâne eroarea seifului
	}
	caleSeif := filepath.Join(radacina, NumeSeif)
	if existaFisier(caleSeif) {
		if err := os.Rename(caleSeif, caleSeif+sufixDeteriorat); err != nil {
			return false, err
		}
	}
	if err := os.WriteFile(filepath.Join(radacina, marcajRefacere), nil, 0o600); err != nil {
		return false, err
	}
	tmp := caleSeif + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return false, err
	}
	return true, os.Rename(tmp, caleSeif)
}

// noteazaRefacerea scrie în jurnalul de audit că seiful a fost refăcut din oglindă.
func noteazaRefacerea(ctx context.Context, radacina, caleBaza string, cheie []byte) (audit.Rezultat, bool, error) {
	marcaj := filepath.Join(radacina, marcajRefacere)
	if !existaFisier(marcaj) {
		return audit.Rezultat{}, false, nil
	}
	db, err := bd.Deschide(caleBaza, cheie)
	if err != nil {
		return audit.Rezultat{}, false, err
	}
	defer db.Close()
	if err := audit.NoteazaEveniment(ctx, db, "seif_refacut", nil, map[string]string{"din": "backup/" + NumeSeif}); err != nil {
		return audit.Rezultat{}, false, err
	}
	rez, err := audit.Verifica(ctx, db)
	if err != nil {
		return audit.Rezultat{}, false, err
	}
	return rez, true, os.Remove(marcaj)
}
