package datebaza

import (
	"context"
	"database/sql"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

// F03-T11 (D004, D106): traducerile corectate sau noi, importate din folderul karios/traduceri
// și păstrate în baza criptată. Se păstrează doar textele diferite de cataloagele incluse.

// FisierTraduceri: conținutul unui fișier <limba>.json (cheie → text).
type FisierTraduceri struct {
	Limba string            `json:"limba"`
	Texte map[string]string `json:"texte"`
}

// ProblemaTraducere: o cheie care nu se importă și de ce.
type ProblemaTraducere struct {
	Cheie string `json:"cheie"`
	Motiv string `json:"motiv"` // „necunoscuta”, „parametri”, „prea_lung”
}

// RezumatTraduceri: ce ar face importul pentru o limbă.
type RezumatTraduceri struct {
	Limba       string              `json:"limba"`
	Noi         int                 `json:"noi"`
	Modificate  int                 `json:"modificate"`
	Revenite    int                 `json:"revenite"` // corecturi șterse: textul e iar cel inclus
	Neschimbate int                 `json:"neschimbate"`
	Probleme    []ProblemaTraducere `json:"probleme"`
}

const lungimeMaximaTraducere = 4000

var reParametru = regexp.MustCompile(`\{\{[^{}]+\}\}`)

// parametri întoarce parametrii {{nume}} ai unui text, sortați și fără dubluri.
func parametri(text string) string {
	gasiti := reParametru.FindAllString(text, -1)
	sort.Strings(gasiti)
	unici := gasiti[:0]
	for i, p := range gasiti {
		if i == 0 || p != gasiti[i-1] {
			unici = append(unici, p)
		}
	}
	return strings.Join(unici, ",")
}

// ListaTraduceri întoarce corecturile salvate, pe limbă.
func (s *Serviciu) ListaTraduceri(ctx context.Context) (map[i18n.Limba]map[string]string, error) {
	return listaTraduceri(ctx, s.db)
}

type interogare interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func listaTraduceri(ctx context.Context, db interogare) (map[i18n.Limba]map[string]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT limba, cheie, text FROM traduceri ORDER BY limba, cheie`)
	if err != nil {
		return nil, traduceEroare(err)
	}
	defer rows.Close()
	rez := map[i18n.Limba]map[string]string{}
	for rows.Next() {
		var l, k, t string
		if err := rows.Scan(&l, &k, &t); err != nil {
			return nil, err
		}
		if rez[i18n.Limba(l)] == nil {
			rez[i18n.Limba(l)] = map[string]string{}
		}
		rez[i18n.Limba(l)][k] = t
	}
	return rez, rows.Err()
}

// schimbareTraducere: ce se scrie (Text != "") sau se șterge (Text == "").
type schimbareTraducere struct{ Limba, Cheie, Text string }

// analizeazaTraduceri compară fișierele cu cataloagele incluse și cu corecturile salvate.
func analizeazaTraduceri(fisiere []FisierTraduceri, salvate map[i18n.Limba]map[string]string) ([]RezumatTraduceri, []schimbareTraducere, error) {
	rezumate := []RezumatTraduceri{}
	schimbari := []schimbareTraducere{}
	vazute := map[string]bool{}
	for _, f := range fisiere {
		if !i18n.CodLimbaValid(f.Limba) || vazute[f.Limba] {
			return nil, nil, eroare("traduceri.limba_invalida")
		}
		vazute[f.Limba] = true
		l := i18n.Limba(f.Limba)
		r := RezumatTraduceri{Limba: f.Limba, Probleme: []ProblemaTraducere{}}
		chei := make([]string, 0, len(f.Texte))
		for k := range f.Texte {
			chei = append(chei, k)
		}
		sort.Strings(chei)
		for _, k := range chei {
			text := strings.TrimSpace(f.Texte[k])
			switch {
			case !i18n.CheieCunoscuta(k):
				r.Probleme = append(r.Probleme, ProblemaTraducere{k, "necunoscuta"})
				continue
			case text == "":
				r.Neschimbate++
				continue
			case utf8.RuneCountInString(text) > lungimeMaximaTraducere:
				r.Probleme = append(r.Probleme, ProblemaTraducere{k, "prea_lung"})
				continue
			case parametri(text) != parametri(i18n.TextInclus(i18n.Implicita, k)):
				r.Probleme = append(r.Probleme, ProblemaTraducere{k, "parametri"})
				continue
			}
			vechi, areCorectura := salvate[l][k]
			// Spatiile de la capete nu conteaza (textul importat se scrie fara ele).
			inclus := strings.TrimSpace(i18n.TextInclus(l, k))
			switch {
			case text == inclus && areCorectura:
				r.Revenite++
				schimbari = append(schimbari, schimbareTraducere{f.Limba, k, ""})
			case text == inclus || (areCorectura && text == vechi):
				r.Neschimbate++
			case areCorectura || inclus != "":
				r.Modificate++
				schimbari = append(schimbari, schimbareTraducere{f.Limba, k, text})
			default:
				r.Noi++
				schimbari = append(schimbari, schimbareTraducere{f.Limba, k, text})
			}
		}
		rezumate = append(rezumate, r)
	}
	return rezumate, schimbari, nil
}

// PrevizualizeazaTraduceri arată ce ar face importul, fără să scrie nimic.
func (s *Serviciu) PrevizualizeazaTraduceri(ctx context.Context, fisiere []FisierTraduceri) ([]RezumatTraduceri, error) {
	salvate, err := s.ListaTraduceri(ctx)
	if err != nil {
		return nil, err
	}
	r, _, err := analizeazaTraduceri(fisiere, salvate)
	return r, err
}

// ImportaTraduceri scrie corecturile valide, într-o singură tranzacție, cu un rând de jurnal
// pe limbă (ce s-a importat, nu fiecare text). Cheile cu probleme nu se scriu.
func (s *Serviciu) ImportaTraduceri(ctx context.Context, fisiere []FisierTraduceri) ([]RezumatTraduceri, error) {
	var rezumate []RezumatTraduceri
	err := s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		salvate, err := listaTraduceri(ctx, tx)
		if err != nil {
			return err
		}
		var schimbari []schimbareTraducere
		rezumate, schimbari, err = analizeazaTraduceri(fisiere, salvate)
		if err != nil {
			return err
		}
		acum := moment()
		for _, c := range schimbari {
			if c.Text == "" {
				_, err = tx.ExecContext(ctx, `DELETE FROM traduceri WHERE limba=? AND cheie=?`, c.Limba, c.Cheie)
			} else {
				_, err = tx.ExecContext(ctx, `INSERT INTO traduceri (limba, cheie, text, modificat_la) VALUES (?, ?, ?, ?)
					ON CONFLICT (limba, cheie) DO UPDATE SET text = excluded.text, modificat_la = excluded.modificat_la`, c.Limba, c.Cheie, c.Text, acum)
			}
			if err != nil {
				return err
			}
		}
		for _, r := range rezumate {
			if r.Noi+r.Modificate+r.Revenite == 0 {
				continue
			}
			if err := jurnal(ctx, q, "traduceri", r.Limba, nil, map[string]int{"noi": r.Noi, "modificate": r.Modificate, "revenite": r.Revenite}); err != nil {
				return err
			}
		}
		return nil
	})
	return rezumate, err
}

// StergeTraduceri șterge toate corecturile unei limbi; textele revin la cele incluse.
func (s *Serviciu) StergeTraduceri(ctx context.Context, limba string) error {
	if !i18n.CodLimbaValid(limba) {
		return eroare("traduceri.limba_invalida")
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		r, err := tx.ExecContext(ctx, `DELETE FROM traduceri WHERE limba=?`, limba)
		if err != nil {
			return err
		}
		n, err := r.RowsAffected()
		if err != nil || n == 0 {
			return err
		}
		return jurnal(ctx, q, "traduceri", limba, map[string]int64{"sterse": n}, nil)
	})
}
