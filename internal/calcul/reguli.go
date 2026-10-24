package calcul

import (
	"fmt"
	"sort"
	"time"
)

type TipZi string

const (
	Lucratoare TipZi = "LUCRATOARE"
	Sambata    TipZi = "SAMBATA"
	Duminica   TipZi = "DUMINICA"
	Sarbatoare TipZi = "SARBATOARE"
)

var TipuriZi = []TipZi{Lucratoare, Sambata, Duminica, Sarbatoare}

const minuteZi = 24 * 60

// Interval atribuie o categorie de ore unui interval dintr-un tip de zi.
// Start și Sfarsit sunt minute de la miezul nopții (0–1440), după ora de pe ceas.
type Interval struct {
	TipZi     TipZi
	Start     int
	Sfarsit   int
	Categorie string
}

type SetReguli struct {
	Denumire  string
	Intervale []Interval
}

// Valideaza verifică faptul că, pentru fiecare tip de zi, intervalele acoperă
// exact 00:00–24:00, fără goluri și fără suprapuneri.
func (s SetReguli) Valideaza() error {
	for _, iv := range s.Intervale {
		if !tipCunoscut(iv.TipZi) {
			return &Eroare{"reguli.tip_zi_necunoscut", map[string]string{"tip": string(iv.TipZi)}}
		}
		if iv.Start < 0 || iv.Sfarsit > minuteZi || iv.Sfarsit <= iv.Start {
			return eroareInterval("reguli.interval_invalid", iv.TipZi, iv.Start, iv.Sfarsit)
		}
	}
	for _, tip := range TipuriZi {
		cursor := 0
		for _, iv := range s.intervale(tip) {
			if iv.Start > cursor {
				return eroareInterval("reguli.gol", tip, cursor, iv.Start)
			}
			if iv.Start < cursor {
				return eroareInterval("reguli.suprapunere", tip, iv.Start, cursor)
			}
			cursor = iv.Sfarsit
		}
		if cursor != minuteZi {
			return eroareInterval("reguli.gol", tip, cursor, minuteZi)
		}
	}
	return nil
}

func eroareInterval(cod string, tip TipZi, deLa, panaLa int) *Eroare {
	return &Eroare{cod, map[string]string{"tip": string(tip), "de_la": hm(deLa), "pana_la": hm(panaLa)}}
}

// intervale întoarce intervalele unui tip de zi, ordonate după ora de start.
func (s SetReguli) intervale(tip TipZi) []Interval {
	var rez []Interval
	for _, iv := range s.Intervale {
		if iv.TipZi == tip {
			rez = append(rez, iv)
		}
	}
	sort.Slice(rez, func(i, j int) bool { return rez[i].Start < rez[j].Start })
	return rez
}

func tipCunoscut(tip TipZi) bool {
	for _, t := range TipuriZi {
		if t == tip {
			return true
		}
	}
	return false
}

func hm(minute int) string {
	return fmt.Sprintf("%02d:%02d", minute/60, minute%60)
}

// Calendar stabilește tipul fiecărei zile calendaristice.
type Calendar struct {
	Sarbatori map[string]string // "2006-01-02" → denumire
}

// TipZi întoarce tipul zilei. Sărbătoarea are prioritate față de weekend.
func (c Calendar) TipZi(zi time.Time) TipZi {
	if _, ok := c.Sarbatori[zi.Format(time.DateOnly)]; ok {
		return Sarbatoare
	}
	switch zi.Weekday() {
	case time.Saturday:
		return Sambata
	case time.Sunday:
		return Duminica
	}
	return Lucratoare
}
