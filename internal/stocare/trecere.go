package stocare

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/seif"
)

// Trecerea la D054: seiful cu copia A (cheia aplicației, publică) devine un seif fără cheie în
// cod, cu o cheie de date nouă. Baza și seiful noi se scriu alături (".nou") și înlocuiesc
// fișierele vechi (păstrate ca ".vechi") abia la final; o trecere întreruptă se termină sau se
// reia la pornirea următoare (finalizeazaTrecerea).
const (
	sufixNou   = ".nou"
	sufixVechi = ".vechi"
)

// rezultatTrecere se notează în jurnalul bazei, după trecere.
type rezultatTrecere struct {
	Amprente []string `json:"amprente,omitempty"` // „AAAA-LL: veche → nouă”, pentru hârtiile deja tipărite
	Conturi  int      `json:"conturi_de_reconfirmat"`
}

func existaFisier(cale string) bool {
	_, err := os.Stat(cale)
	return err == nil
}

// mutaBaza redenumește baza și fișierele ei WAL/SHM.
func mutaBaza(de, la string) error {
	for _, sufix := range []string{"", "-wal", "-shm"} {
		if err := os.Rename(de+sufix, la+sufix); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}

func stergeBaza(cale string) {
	for _, sufix := range []string{"", "-wal", "-shm"} {
		_ = os.Remove(cale + sufix)
	}
}

type amprentaLuna struct {
	an, luna       int
	pastrata, noua string
	corespunde     bool
}

func amprenteInchise(ctx context.Context, db *sql.DB) []amprentaLuna {
	r, err := db.QueryContext(ctx, `SELECT an, luna, amprenta FROM stare_luni WHERE amprenta <> ''`)
	if err != nil {
		return nil // schema dinaintea amprentelor
	}
	defer r.Close()
	var rez []amprentaLuna
	for r.Next() {
		var a amprentaLuna
		if r.Scan(&a.an, &a.luna, &a.pastrata) == nil {
			rez = append(rez, a)
		}
	}
	return rez
}

// treceLaD054 face trecerea pentru folderul de date dat și întoarce seiful nou, deschis, și cheia
// de recuperare de afișat.
func treceLaD054(ctx context.Context, radacina string) (*seif.Seif, string, error) {
	caleSeif, caleBaza := filepath.Join(radacina, NumeSeif), filepath.Join(radacina, NumeBaza)
	vechi, err := seif.Deschide(caleSeif)
	if err != nil {
		return nil, "", err
	}
	_ = os.Remove(caleSeif + sufixNou)
	stergeBaza(caleBaza + sufixNou)
	nou, err := seif.CreeazaNou(caleSeif+sufixNou, bd.LungimeCheie)
	if err != nil {
		return nil, "", err
	}

	// 1. Cu cheia veche: starea sigiliilor, amprentele lunilor închise și copia bazei cu cheia nouă.
	areBaza := existaFisier(caleBaza)
	resigileaza := false
	var amprente []amprentaLuna
	if areBaza {
		db0, err := bd.Deschide(caleBaza, vechi.CheieDate)
		if err != nil {
			return nil, "", err
		}
		rez, err := audit.Verifica(ctx, db0)
		if err != nil {
			db0.Close()
			return nil, "", err
		}
		// O bază cu incident rămâne cu incident: sub cheia nouă sigiliile nu vor mai corespunde.
		resigileaza = rez.Stare == audit.StareOK
		amprente = amprenteInchise(ctx, db0)
		for i := range amprente {
			if a, err := audit.AmprentaLuna(ctx, db0, amprente[i].an, amprente[i].luna); err == nil {
				amprente[i].corespunde = a == amprente[i].pastrata
			}
		}
		err = bd.Backup(ctx, db0, caleBaza+sufixNou, nou.CheieDate)
		db0.Close()
		if err != nil {
			return nil, "", fmt.Errorf("trecerea D054, copia bazei: %w", err)
		}
	}

	// 2. Seiful nou, alături: conturile ca P0/R0, cheia de recuperare, copia W.
	nou.PreiaCopiiVechi(vechi)
	cod, err := nou.SeteazaRecuperare()
	if err != nil {
		return nil, "", err
	}
	_ = nou.AdaugaWindows() // lipsește doar în afara Windows; datele se deschid atunci cu parola

	// 3. Înlocuirea: de aici înainte, datele țin de cheia nouă.
	if err := os.Rename(caleSeif, caleSeif+sufixVechi); err != nil {
		return nil, "", err
	}
	if areBaza {
		if err := mutaBaza(caleBaza, caleBaza+sufixVechi); err != nil {
			return nil, "", err
		}
		if err := mutaBaza(caleBaza+sufixNou, caleBaza); err != nil {
			return nil, "", err
		}
	}
	if err := os.Rename(caleSeif+sufixNou, caleSeif); err != nil {
		return nil, "", err
	}
	nou.MutaLa(caleSeif)

	// 4. Backup-urile, apoi amprentele noi, jurnalul și sigiliile sub cheia nouă.
	if err := recripteazaBackupuri(ctx, radacina, vechi.CheieDate, nou.CheieDate); err != nil {
		return nil, "", err
	}
	if areBaza {
		db1, err := bd.Deschide(caleBaza, nou.CheieDate)
		if err != nil {
			return nil, "", err
		}
		rt := rezultatTrecere{Conturi: len(nou.CopiiCuPrefix(seif.PrefixParolaVeche))}
		for _, a := range amprente {
			if !a.corespunde {
				continue // datele erau deja schimbate față de amprentă: rămâne nepotrivirea
			}
			noua, err := audit.AmprentaLuna(ctx, db1, a.an, a.luna)
			if err != nil {
				db1.Close()
				return nil, "", err
			}
			if _, err := db1.ExecContext(ctx, `UPDATE stare_luni SET amprenta = ? WHERE an = ? AND luna = ?`, noua, a.an, a.luna); err != nil {
				db1.Close()
				return nil, "", err
			}
			rt.Amprente = append(rt.Amprente, fmt.Sprintf("%04d-%02d: %s → %s", a.an, a.luna, a.pastrata, noua))
		}
		if resigileaza {
			err = audit.NoteazaEveniment(ctx, db1, "trecere_d054", nil, rt)
		}
		db1.Close()
		if err != nil {
			return nil, "", err
		}
	}

	// 5. Curățenie: cheia veche nu mai deschide nimic din datele actuale.
	_ = os.Remove(caleSeif + sufixVechi)
	stergeBaza(caleBaza + sufixVechi)
	return nou, cod, nil
}

// recripteazaBackupuri trece backup-urile pe cheia nouă; cele deja trecute se sar.
func recripteazaBackupuri(ctx context.Context, radacina string, vechi, noua []byte) error {
	dir := filepath.Join(radacina, "backup")
	intrari, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, i := range intrari {
		if i.IsDir() || !strings.HasSuffix(i.Name(), sufixBackup) {
			continue
		}
		cale := filepath.Join(dir, i.Name())
		if db, err := bd.Deschide(cale, noua); err == nil {
			db.Close()
			continue
		}
		db, err := bd.Deschide(cale, vechi)
		if err != nil {
			continue // nici cheia veche nu-l deschide: nu e al acestui folder, se lasă neatins
		}
		err = bd.Backup(ctx, db, cale+sufixNou, noua)
		db.Close()
		if err != nil {
			return fmt.Errorf("trecerea D054, backup %s: %w", i.Name(), err)
		}
		stergeBaza(cale)
		if err := os.Rename(cale+sufixNou, cale); err != nil {
			return err
		}
	}
	return nil
}

// finalizeazaTrecerea rezolvă o trecere D054 întreruptă:
//   - înainte de înlocuire (seiful vechi încă la locul lui): fișierele ".nou" se aruncă, iar
//     trecerea se reia de la capăt;
//   - în timpul sau după înlocuire (există seif.vechi): înlocuirea se termină, backup-urile
//     rămase se recriptează, apoi fișierele vechi se șterg. Amprentele și sigiliile nu se mai
//     refac; pornirea va semnala incidentul, iar administratorul îl confirmă.
func finalizeazaTrecerea(ctx context.Context, radacina string) error {
	caleSeif, caleBaza := filepath.Join(radacina, NumeSeif), filepath.Join(radacina, NumeBaza)
	if !existaFisier(caleSeif + sufixVechi) {
		if existaFisier(caleSeif + sufixNou) {
			_ = os.Remove(caleSeif + sufixNou)
			stergeBaza(caleBaza + sufixNou)
		}
		return nil
	}
	if existaFisier(caleBaza+sufixNou) && !existaFisier(caleBaza) {
		if err := mutaBaza(caleBaza+sufixNou, caleBaza); err != nil {
			return err
		}
	}
	if existaFisier(caleSeif+sufixNou) && !existaFisier(caleSeif) {
		if err := os.Rename(caleSeif+sufixNou, caleSeif); err != nil {
			return err
		}
	}
	vechi, err := seif.Deschide(caleSeif + sufixVechi)
	if err != nil {
		return err
	}
	nou, err := seif.Incarca(caleSeif)
	if err != nil {
		return err
	}
	// Fără copia W (alt utilizator Windows), fișierele vechi rămân până la o pornire care o are,
	// ca backup-urile să poată fi încă recriptate.
	if err := nou.DeschideWindows(); err != nil {
		return nil
	}
	if err := recripteazaBackupuri(ctx, radacina, vechi.CheieDate, nou.CheieDate); err != nil {
		return err
	}
	_ = os.Remove(caleSeif + sufixVechi)
	stergeBaza(caleBaza + sufixVechi)
	return nil
}
