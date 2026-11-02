package rapoarte

import (
	"math"
	"pontix/internal/i18n"
	"regexp"
	"strings"
)

type ColoanaCalculata struct {
	Cod       string    `json:"cod"`
	Denumire  i18n.Text `json:"denumire"`
	Unitate   string    `json:"unitate"`
	Mod       string    `json:"mod"`
	Expresie  string    `json:"expresie"`
	Zecimale  int       `json:"zecimale"`
	LipsaZero bool      `json:"lipsa_zero"`
	Activ     bool      `json:"activ"`
}
type ElementRaport struct {
	RepetaRand bool    `json:"repeta_rand"`
	ID         string  `json:"id"`
	Text       string  `json:"text"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Latime     float64 `json:"latime"`
	Inaltime   float64 `json:"inaltime"`
	Font       float64 `json:"font"`
}
type RaportFlexibil struct {
	InaltimeRand      float64            `json:"inaltime_rand"`
	Ascunse           []string           `json:"ascunse"`
	SQL               string             `json:"sql"`
	ModAspect         string             `json:"mod_aspect"`
	Markdown          string             `json:"markdown"`
	Elemente          []ElementRaport    `json:"elemente"`
	Calculate         []ColoanaCalculata `json:"calculate"`
	ConcediiTipuri    []string           `json:"concedii_tipuri"`
	ConcediiCategorii []string           `json:"concedii_categorii"`
	Weekend           bool               `json:"weekend"`
	Sarbatori         bool               `json:"sarbatori"`
}
type TabelRaport struct {
	Cod     string   `json:"cod"`
	Campuri []string `json:"campuri"`
}
type DateFlexibile struct {
	Coloane []string      `json:"coloane"`
	Randuri [][]any       `json:"randuri"`
	Schema  []TabelRaport `json:"schema"`
}

var codFlexibil = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,63}$`)

func esteFlexibil(s string) bool {
	return s == "SQL_LIBER" || s == "LISTA_GARZI" || s == "CALCULATE_LUNAR"
}
func valideazaFlexibil(f *RaportFlexibil) error {
	if f == nil || len(f.SQL) > 32768 || len(f.Markdown) > 32000 || len(f.Elemente) > 60 || len(f.Calculate) > 80 {
		return eroare("rapoarte.flex_invalid")
	}
	if f.ModAspect != "" && f.ModAspect != "TABEL" && f.ModAspect != "MARKDOWN" && f.ModAspect != "LIBER" {
		return eroare("rapoarte.flex_invalid")
	}
	vaz := map[string]bool{}
	for _, c := range f.Calculate {
		if !codFlexibil.MatchString(c.Cod) || vaz[strings.ToLower(c.Cod)] || len(c.Expresie) > 32768 || c.Zecimale < 0 || c.Zecimale > 6 || (c.Mod != "SQL" && c.Mod != "FORMULA") || (c.Unitate != "ORE" && c.Unitate != "ZILE" && c.Unitate != "NUMAR") {
			return eroare("rapoarte.flex_invalid")
		}
		vaz[strings.ToLower(c.Cod)] = true
	}
	tabele := 0
	repetate := 0
	for _, e := range f.Elemente {
		if f.ModAspect == "LIBER" && e.RepetaRand {
			repetate++
		}
	}
	if repetate > 0 && (f.InaltimeRand < 5 || f.InaltimeRand > 400 || math.IsNaN(f.InaltimeRand) || math.IsInf(f.InaltimeRand, 0)) {
		return eroare("rapoarte.flex_invalid")
	}
	if f.ModAspect == "LIBER" && len(f.Elemente) == 0 {
		return eroare("rapoarte.flex_invalid")
	}
	for _, e := range f.Elemente {
		if f.ModAspect == "LIBER" && strings.TrimSpace(e.Text) == "{{tabel}}" {
			tabele++
			if tabele > 1 || repetate > 0 {
				return eroare("rapoarte.flex_invalid")
			}
		}
		for _, v := range []float64{e.X, e.Y, e.Latime, e.Inaltime, e.Font} {
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return eroare("rapoarte.flex_invalid")
			}
		}
		if e.X < 0 || e.Y < 0 || e.Latime < 5 || e.Inaltime < 5 || e.Font < 6 || e.Font > 40 || len(e.Text) > 4000 {
			return eroare("rapoarte.flex_invalid")
		}
	}
	return nil
}
