package rapoarte

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

// Rapoartele standard ca tabele calculate in setul de date SQL (D090, F09-T21).
// Calculul este cel al motorului vechi (pregatesteStandard); SQL-ul doar alege,
// ordoneaza si filtreaza coloanele rezultate.

type tabelStandard struct {
	Tabel string
	Sursa string
}

var tabeleStandard = []tabelStandard{
	{"pontaj_lunar", "PONTAJ_LUNAR"},
	{"pontaj_detaliat", "PONTAJ_DETALIAT"},
	{"fisa_individuala", "FISA_INDIVIDUALA"},
	{"absente_anuale", "ABSENTE_ANUALE"},
	{"grafic_garzi", "GRAFIC_GARZI"},
}

// documentStandard pastreaza documentul calculat si codurile SQL ale coloanelor,
// pentru preluarea prezentarii prin coloana rand.
type documentStandard struct {
	tabel  string
	doc    Document
	coduri []string
}

var dataColoana = regexp.MustCompile(`^\d{4}-\d{2}-(\d{2})$`)

// codSQLStandard transforma codurile documentului in coloane SQL stabile:
// zilele (AAAA-LL-ZZ) devin zZZ, ca un SQL salvat sa mearga in orice luna.
func codSQLStandard(cod string) string {
	if m := dataColoana.FindStringSubmatch(cod); m != nil {
		return "z" + m[1]
	}
	return cod
}

// folosesteTabel spune daca textul SQL (interogarea si calculele) numeste tabelul.
func folosesteTabel(text, tabel string) bool {
	return regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(tabel) + `\b`).MatchString(text)
}

func definitieStandard(sursa string) Definitie {
	s := SursaPentru(sursa)
	coloane := []string{}
	if s != nil {
		coloane = append(coloane, s.Coloane...)
		// F09-T37: tabelul de date are si orele suplimentare (total_OS), dupa celelalte totaluri.
		if len(coloane) > 0 && coloane[len(coloane)-1] == "totaluri" {
			coloane = append(coloane, "totaluri_os")
		}
	}
	return Definitie{
		ID: "sql-" + strings.ToLower(sursa), Denumire: i18n.Text{"ro": sursa, "en": sursa},
		Continut: Continut{Sursa: sursa, Coloane: coloane, Text: i18n.Text{}},
		Antet:    Bloc{Stanga: i18n.Text{}, Centru: i18n.Text{}, Dreapta: i18n.Text{}},
		Subsol:   Bloc{Stanga: i18n.Text{}, Centru: i18n.Text{}, Dreapta: i18n.Text{}},
		Pagina:   Pagina{Format: "A4", Orientare: "PEISAJ", MargineMM: 10, FontPT: 8},
		Activ:    true, Revizie: "sql",
	}
}

// adaugaStandard calculeaza numai tabelele standard folosite de SQL-ul raportului.
func (b *bazaRaport) adaugaStandard(ctx context.Context, tx *sql.Tx, c Cerere, text string, acum time.Time) error {
	for _, t := range tabeleStandard {
		if !folosesteTabel(text, t.Tabel) {
			continue
		}
		doc, err := pregatesteStandard(ctx, tx, acces.New(tx), definitieStandard(t.Sursa), c, acum)
		if err != nil {
			return err
		}
		coduri := make([]string, len(doc.Coloane))
		cols := []string{"rand", "grup"}
		vazute := map[string]bool{"rand": true, "grup": true}
		for i, col := range doc.Coloane {
			coduri[i] = codSQLStandard(col.Cod)
			if !codFlexibil.MatchString(coduri[i]) || vazute[strings.ToLower(coduri[i])] {
				return eroare("rapoarte.flex_invalid")
			}
			vazute[strings.ToLower(coduri[i])] = true
			cols = append(cols, coduri[i])
		}
		randuri := make([][]any, 0, len(doc.Randuri))
		for i, r := range doc.Randuri {
			grup := 0
			if r.Grup {
				grup = 1
			}
			rand := []any{i, grup}
			for j := range doc.Coloane {
				var v any
				if j < len(r.Valori) {
					v = r.Valori[j]
				}
				rand = append(rand, v)
			}
			randuri = append(randuri, rand)
		}
		cols, randuri = completeazaZile(cols, randuri)
		if err := b.adauga(ctx, t.Tabel, cols, randuri); err != nil {
			return err
		}
		b.standard = append(b.standard, documentStandard{tabel: t.Tabel, doc: doc, coduri: coduri})
	}
	return nil
}

// aplicaMetadateStandard preia din documentul calculat titlurile, culorile,
// grupurile si legenda, daca rezultatul pastreaza coloana rand (D090 punctul 4).
// Se aplica inaintea ascunderii coloanelor, pe coloanele rezultatului SQL.
func aplicaMetadateStandard(doc *Document, b *bazaRaport) {
	if len(b.standard) != 1 {
		return
	}
	s := b.standard[0]
	if an, e1 := strconv.Atoi(doc.VariabileRaport["an"]); e1 == nil {
		if luna, e2 := strconv.Atoi(doc.VariabileRaport["luna"]); e2 == nil && luna >= 1 && luna <= 12 {
			doc.zileLuna = time.Date(an, time.Month(luna)+1, 0, 0, 0, 0, 0, time.UTC).Day()
		}
	}
	idxRand := -1
	for j, c := range doc.Coloane {
		if c.Cod == "rand" {
			idxRand = j
		}
	}
	orig := map[string]int{}
	for i, cod := range s.coduri {
		orig[cod] = i
	}
	for j, c := range doc.Coloane {
		if i, ok := orig[c.Cod]; ok {
			doc.Coloane[j].Titlu = s.doc.Coloane[i].Titlu
			doc.Coloane[j].Culoare = s.doc.Coloane[i].Culoare
		}
	}
	// Legenda explica doar coloanele alese, in ordinea rezultatului.
	alese := []Coloana{}
	for _, c := range doc.Coloane {
		if i, ok := orig[c.Cod]; ok {
			alese = append(alese, s.doc.Coloane[i])
		}
	}
	if s.doc.legendaPentru != nil && len(alese) > 0 {
		doc.Legenda = append(s.doc.legendaPentru(alese), doc.Legenda...)
	}
	if idxRand < 0 {
		return
	}
	for ri := range doc.Randuri {
		r := &doc.Randuri[ri]
		n, err := strconv.Atoi(strings.TrimSpace(r.Valori[idxRand]))
		if err != nil || n < 0 || n >= len(s.doc.Randuri) {
			continue
		}
		o := s.doc.Randuri[n]
		r.Grup, r.Aldin = o.Grup, o.Aldin
		culori := make([]string, len(doc.Coloane))
		are := false
		for j, c := range doc.Coloane {
			if i, ok := orig[c.Cod]; ok && i < len(o.Culori) && o.Culori[i] != "" {
				culori[j], are = o.Culori[i], true
			}
		}
		if are {
			r.Culori = culori
		}
	}
}

// textFolositSQL aduna textele care pot numi tabele: interogarea si calculele SQL.
func textFolositSQL(f *RaportFlexibil, sqlText string) string {
	parti := []string{sqlText}
	for _, c := range f.Calculate {
		if c.Activ && c.Mod == "SQL" {
			parti = append(parti, c.Expresie)
		}
	}
	return strings.Join(parti, "\n")
}

var codZi = regexp.MustCompile(`^z(\d{2})$`)

// completeazaZile adauga coloanele z29-z31 lipsa in lunile mai scurte, goale,
// imediat dupa ultima zi, ca un SQL salvat sa mearga in orice luna (D090).
func completeazaZile(cols []string, randuri [][]any) ([]string, [][]any) {
	ultima, maxZi := -1, 0
	for j, c := range cols {
		if m := codZi.FindStringSubmatch(c); m != nil {
			n, _ := strconv.Atoi(m[1])
			ultima = j
			if n > maxZi {
				maxZi = n
			}
		}
	}
	if ultima < 0 || maxZi >= 31 {
		return cols, randuri
	}
	lipsa := []string{}
	for n := maxZi + 1; n <= 31; n++ {
		lipsa = append(lipsa, fmt.Sprintf("z%02d", n))
	}
	noi := append(append(append([]string{}, cols[:ultima+1]...), lipsa...), cols[ultima+1:]...)
	for i, r := range randuri {
		goale := make([]any, len(lipsa))
		randuri[i] = append(append(append([]any{}, r[:ultima+1]...), goale...), r[ultima+1:]...)
	}
	return noi, randuri
}
