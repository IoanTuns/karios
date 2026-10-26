// Package rapoarte gestionează definițiile și documentele, separat de calculul salarial.
package rapoarte

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type Bloc struct {
	Stanga  i18n.Text          `json:"stanga"`
	Centru  i18n.Text          `json:"centru"`
	Dreapta i18n.Text          `json:"dreapta"`
	Imagini map[string]Imagine `json:"imagini,omitempty"`
}
type Continut struct {
	Sursa   string    `json:"sursa"`
	Coloane []string  `json:"coloane"`
	Text    i18n.Text `json:"text"`
}
type Pagina struct {
	Format    string `json:"format"`
	Orientare string `json:"orientare"`
	MargineMM int64  `json:"margine_mm"`
	FontPT    int64  `json:"font_pt"`
}
type Definitie struct {
	ID       string    `json:"id"`
	Denumire i18n.Text `json:"denumire"`
	Antet    Bloc      `json:"antet"`
	Continut Continut  `json:"continut"`
	Subsol   Bloc      `json:"subsol"`
	Pagina   Pagina    `json:"pagina"`
	Activ    bool      `json:"activ"`
	Revizie  string    `json:"revizie"`
}
type Sursa struct {
	Cod     string   `json:"cod"`
	Coloane []string `json:"coloane"`
}
type Configuratie struct {
	Modele    []Definitie `json:"modele"`
	Surse     []Sursa     `json:"surse"`
	Variabile []string    `json:"variabile"`
}
type Serviciu struct{ db *sql.DB }

func Nou(db *sql.DB) *Serviciu { return &Serviciu{db: db} }

var surse = []Sursa{
	{Cod: "PONTAJ_LUNAR", Coloane: []string{"angajat", "marca", "sectie", "departament", "functie", "zile"}},
	{Cod: "PONTAJ_DETALIAT", Coloane: []string{"angajat", "marca", "sectie", "departament", "functie", "zile", "totaluri"}},
	{Cod: "FISA_INDIVIDUALA", Coloane: []string{"data", "zi", "ore", "descriere", "categorie"}},
	{Cod: "ABSENTE_ANUALE", Coloane: []string{"data", "angajat", "marca", "sectie", "departament", "cod", "descriere", "categorie"}},
}
var variabile = []string{"titlu", "an", "luna", "perioada", "sectie", "departament", "angajat", "generat_la"}
var variabila = regexp.MustCompile(`\{\{([^{}]+)\}\}`)

func eroare(cod string) error { return &i18n.Eroare{Cod: cod} }
func traduce(err error) error {
	if err == nil {
		return nil
	}
	var e *i18n.Eroare
	if errors.As(err, &e) {
		return err
	}
	return eroare("rapoarte.eroare")
}
func IDNou() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func SursaPentru(cod string) *Sursa {
	for _, s := range surse {
		if s.Cod == cod {
			return &s
		}
	}
	return nil
}
func contine(lista []string, valoare string) bool {
	for _, v := range lista {
		if v == valoare {
			return true
		}
	}
	return false
}

func Valideaza(d Definitie) error {
	for _, bloc := range []Bloc{d.Antet, d.Subsol} {
		for zona, imagine := range bloc.Imagini {
			if !contine([]string{"stanga", "centru", "dreapta"}, zona) {
				return eroare("rapoarte.imagine_invalida")
			}
			if err := valideazaImagine(imagine); err != nil {
				return err
			}
		}
	}
	nume := false
	for _, v := range d.Denumire {
		if strings.TrimSpace(v) != "" {
			nume = true
		}
		if utf8.RuneCountInString(v) > 100 || strings.ContainsAny(v, "\r\n\x00") {
			return eroare("rapoarte.nume_invalid")
		}
	}
	if !nume {
		return eroare("rapoarte.nume_invalid")
	}
	sursa := SursaPentru(d.Continut.Sursa)
	if sursa == nil || len(d.Continut.Coloane) == 0 {
		return eroare("rapoarte.continut_invalid")
	}
	vazute := map[string]bool{}
	for _, c := range d.Continut.Coloane {
		if !contine(sursa.Coloane, c) || vazute[c] {
			return eroare("rapoarte.continut_invalid")
		}
		vazute[c] = true
	}
	if d.Continut.Sursa == "PONTAJ_LUNAR" && !vazute["angajat"] {
		return eroare("rapoarte.continut_invalid")
	}
	if d.Continut.Sursa == "PONTAJ_DETALIAT" && (!vazute["angajat"] || !vazute["zile"] || d.Continut.Coloane[len(d.Continut.Coloane)-1] != "totaluri") {
		return eroare("rapoarte.continut_invalid")
	}
	if (d.Pagina.Format != "A4" && d.Pagina.Format != "A3") || (d.Pagina.Orientare != "PORTRET" && d.Pagina.Orientare != "PEISAJ") || d.Pagina.MargineMM < 5 || d.Pagina.MargineMM > 25 || d.Pagina.FontPT < 6 || d.Pagina.FontPT > 14 {
		return eroare("rapoarte.pagina_invalida")
	}
	for _, texte := range []i18n.Text{d.Antet.Stanga, d.Antet.Centru, d.Antet.Dreapta, d.Continut.Text, d.Subsol.Stanga, d.Subsol.Centru, d.Subsol.Dreapta} {
		for _, v := range texte {
			if utf8.RuneCountInString(v) > 1000 || strings.Count(v, "\n") > 5 || strings.ContainsRune(v, 0) {
				return eroare("rapoarte.text_invalid")
			}
			for _, m := range variabila.FindAllStringSubmatch(v, -1) {
				if !contine(variabile, m[1]) {
					return eroare("rapoarte.variabila_invalida")
				}
			}
			rest := variabila.ReplaceAllString(v, "")
			if strings.Contains(rest, "{{") || strings.Contains(rest, "}}") {
				return eroare("rapoarte.variabila_invalida")
			}
		}
	}
	return nil
}

func decodeaza(r acces.DefinitiiRapoarte) (Definitie, error) {
	d := Definitie{ID: r.ID, Activ: r.Activ == 1, Revizie: strconv.FormatInt(r.Revizie, 10)}
	for _, camp := range []struct {
		text string
		dest any
	}{{r.Denumire, &d.Denumire}, {r.Antet, &d.Antet}, {r.Continut, &d.Continut}, {r.Subsol, &d.Subsol}, {r.Pagina, &d.Pagina}} {
		if err := json.Unmarshal([]byte(camp.text), camp.dest); err != nil {
			return Definitie{}, err
		}
	}
	return d, nil
}
func (s *Serviciu) Lista(ctx context.Context) (Configuratie, error) {
	r, err := acces.New(s.db).ListaDefinitiiRapoarte(ctx)
	if err != nil {
		return Configuratie{}, traduce(err)
	}
	rez := Configuratie{Modele: []Definitie{}, Surse: surse, Variabile: variabile}
	for _, v := range r {
		d, err := decodeaza(v)
		if err != nil {
			return Configuratie{}, traduce(err)
		}
		rez.Modele = append(rez.Modele, d)
	}
	return rez, nil
}
func (s *Serviciu) Salveaza(ctx context.Context, d Definitie) (string, error) {
	if err := Valideaza(d); err != nil {
		return "", err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", traduce(err)
	}
	defer tx.Rollback()
	q := acces.New(tx)
	codifica := func(v any) string { b, _ := json.Marshal(v); return string(b) }
	acum := time.Now().UTC().Format(time.RFC3339Nano)
	activ := int64(0)
	if d.Activ {
		activ = 1
	}
	var inainte any
	if d.ID == "" {
		if d.Revizie != "" {
			return "", eroare("rapoarte.conflict")
		}
		d.ID = IDNou()
		d.Revizie = "1"
		err = q.AdaugaDefinitieRaport(ctx, acces.AdaugaDefinitieRaportParams{ID: d.ID, Denumire: codifica(d.Denumire), Antet: codifica(d.Antet), Continut: codifica(d.Continut), Subsol: codifica(d.Subsol), Pagina: codifica(d.Pagina), Activ: activ, CreatLa: acum, ModificatLa: acum})
	} else {
		veche, citireErr := q.CitesteDefinitieRaport(ctx, d.ID)
		if errors.Is(citireErr, sql.ErrNoRows) {
			return "", eroare("rapoarte.inexistent")
		}
		if citireErr != nil {
			return "", traduce(citireErr)
		}
		revizie, convErr := strconv.ParseInt(d.Revizie, 10, 64)
		if convErr != nil || revizie != veche.Revizie {
			return "", eroare("rapoarte.conflict")
		}
		inainte = veche
		var noua int64
		noua, err = q.ModificaDefinitieRaport(ctx, acces.ModificaDefinitieRaportParams{ID: d.ID, Denumire: codifica(d.Denumire), Antet: codifica(d.Antet), Continut: codifica(d.Continut), Subsol: codifica(d.Subsol), Pagina: codifica(d.Pagina), Activ: activ, ModificatLa: acum, Revizie: revizie})
		if errors.Is(err, sql.ErrNoRows) {
			return "", eroare("rapoarte.conflict")
		}
		d.Revizie = strconv.FormatInt(noua, 10)
	}
	if err != nil {
		return "", traduce(err)
	}
	var vechi *string
	if inainte != nil {
		text := codifica(inainte)
		vechi = &text
	}
	nou := codifica(d)
	if err := q.AdaugaJurnal(ctx, acces.AdaugaJurnalParams{Tabel: "definitii_rapoarte", RandID: d.ID, ValoareVeche: vechi, ValoareNoua: &nou, Cand: acum, Cine: "local"}); err != nil {
		return "", traduce(err)
	}
	return d.ID, traduce(tx.Commit())
}
