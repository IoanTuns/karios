package calcul

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"pontix/internal/i18n"
)

// reguliDoc este tabelul din documentul de design (23.10.2026). Procentele nu
// contează pentru segmentare, doar categoriile.
var reguliDoc = SetReguli{
	Denumire: "document 2026-10-23",
	Intervale: []Interval{
		{Lucratoare, 0, 6 * 60, "NOAPTE"},
		{Lucratoare, 6 * 60, 22 * 60, "ZI"},
		{Lucratoare, 22 * 60, 24 * 60, "NOAPTE"},
		{Sambata, 0, 24 * 60, "WEEKEND"},
		{Duminica, 0, 24 * 60, "WEEKEND"},
		{Sarbatoare, 0, 24 * 60, "SARBATOARE"},
	},
}

var calendarTest = Calendar{Sarbatori: map[string]string{
	"2026-12-01": "Ziua Națională",
	"2026-12-25": "Crăciunul",
	"2026-12-26": "Crăciunul",
	"2027-01-02": "Anul Nou", // cade sâmbăta
}}

func ora(t *testing.T, s string) time.Time {
	t.Helper()
	rez, err := time.ParseInLocation("2006-01-02 15:04", s, Locatie)
	if err != nil {
		t.Fatal(err)
	}
	return rez
}

// rand este forma compactă a unui segment folosită în comparații.
type rand struct {
	Data      string
	Categorie string
	Minute    int
}

func TestSegmenteaza(t *testing.T) {
	cazuri := []struct {
		nume           string
		start, sfarsit string
		asteptat       []rand
	}{
		{
			nume:  "tură de zi 12 h",
			start: "2026-10-19 07:00", sfarsit: "2026-10-19 19:00",
			asteptat: []rand{{"2026-10-19", "ZI", 720}},
		},
		{
			nume:  "vineri spre sâmbătă (exemplul din document)",
			start: "2026-10-23 19:00", sfarsit: "2026-10-24 07:00",
			asteptat: []rand{
				{"2026-10-23", "ZI", 180},
				{"2026-10-23", "NOAPTE", 120},
				{"2026-10-24", "WEEKEND", 420},
			},
		},
		{
			nume:  "duminică spre luni (exemplul din document)",
			start: "2026-10-18 19:00", sfarsit: "2026-10-19 07:00",
			asteptat: []rand{
				{"2026-10-18", "WEEKEND", 300},
				{"2026-10-19", "NOAPTE", 360},
				{"2026-10-19", "ZI", 60},
			},
		},
		{
			nume:  "gardă 24 h sâmbătă spre duminică, toată la weekend",
			start: "2026-10-17 07:00", sfarsit: "2026-10-18 07:00",
			asteptat: []rand{
				{"2026-10-17", "WEEKEND", 1020},
				{"2026-10-18", "WEEKEND", 420},
			},
		},
		{
			nume:  "gardă 24 h în zi lucrătoare",
			start: "2026-10-21 07:00", sfarsit: "2026-10-22 07:00",
			asteptat: []rand{
				{"2026-10-21", "ZI", 900},
				{"2026-10-21", "NOAPTE", 120},
				{"2026-10-22", "NOAPTE", 360},
				{"2026-10-22", "ZI", 60},
			},
		},
		{
			nume:  "noapte spre o sărbătoare (1 Decembrie)",
			start: "2026-11-30 19:00", sfarsit: "2026-12-01 07:00",
			asteptat: []rand{
				{"2026-11-30", "ZI", 180},
				{"2026-11-30", "NOAPTE", 120},
				{"2026-12-01", "SARBATOARE", 420},
			},
		},
		{
			nume:  "sărbătoarea are prioritate față de sâmbătă",
			start: "2027-01-02 07:00", sfarsit: "2027-01-02 19:00",
			asteptat: []rand{{"2027-01-02", "SARBATOARE", 720}},
		},
		{
			nume:  "se termină exact la o limită de interval",
			start: "2026-10-20 14:00", sfarsit: "2026-10-20 22:00",
			asteptat: []rand{{"2026-10-20", "ZI", 480}},
		},
		{
			nume:  "începe exact la o limită de interval (S3 22–06)",
			start: "2026-10-20 22:00", sfarsit: "2026-10-21 06:00",
			asteptat: []rand{
				{"2026-10-20", "NOAPTE", 120},
				{"2026-10-21", "NOAPTE", 360},
			},
		},
		{
			nume:  "ora se dă înapoi: 12 h pe ceas = 13 h reale",
			start: "2026-10-24 19:00", sfarsit: "2026-10-25 07:00",
			asteptat: []rand{
				{"2026-10-24", "WEEKEND", 300},
				{"2026-10-25", "WEEKEND", 480},
			},
		},
		{
			nume:  "ora se dă înainte: 12 h pe ceas = 11 h reale",
			start: "2027-03-27 19:00", sfarsit: "2027-03-28 07:00",
			asteptat: []rand{
				{"2027-03-27", "WEEKEND", 300},
				{"2027-03-28", "WEEKEND", 360},
			},
		},
	}

	for _, c := range cazuri {
		t.Run(c.nume, func(t *testing.T) {
			segmente, err := Segmenteaza(ora(t, c.start), ora(t, c.sfarsit), reguliDoc, calendarTest)
			if err != nil {
				t.Fatal(err)
			}
			var obtinut []rand
			for _, s := range segmente {
				obtinut = append(obtinut, rand{s.Data.Format(time.DateOnly), s.Categorie, s.Minute})
			}
			if !reflect.DeepEqual(obtinut, c.asteptat) {
				t.Errorf("\nobținut:  %v\nașteptat: %v", obtinut, c.asteptat)
			}
		})
	}
}

func TestSegmenteazaSegmenteContigue(t *testing.T) {
	start, sfarsit := ora(t, "2026-10-24 19:00"), ora(t, "2026-10-25 07:00")
	segmente, err := Segmenteaza(start, sfarsit, reguliDoc, calendarTest)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for i, s := range segmente {
		if i > 0 && !s.Start.Equal(segmente[i-1].Sfarsit) {
			t.Errorf("segmentul %d nu continuă segmentul anterior", i)
		}
		total += s.Minute
	}
	if want := int(sfarsit.Sub(start) / time.Minute); total != want {
		t.Errorf("total %d minute, așteptat %d", total, want)
	}
}

func TestSegmenteazaTuraInversata(t *testing.T) {
	_, err := Segmenteaza(ora(t, "2026-10-20 19:00"), ora(t, "2026-10-20 07:00"), reguliDoc, calendarTest)
	var e *Eroare
	if !errors.As(err, &e) || e.Cod != "calcul.tura_inversata" {
		t.Fatalf("aștepta eroarea calcul.tura_inversata, a primit %v", err)
	}
	if got, want := e.Text(i18n.Engleza), "The shift ends (2026-10-20 07:00) before it starts (2026-10-20 19:00)."; got != want {
		t.Errorf("en: %q, așteptat %q", got, want)
	}
}

func TestValideaza(t *testing.T) {
	if err := reguliDoc.Valideaza(); err != nil {
		t.Fatalf("setul din document ar trebui să fie valid: %v", err)
	}

	// cuGresit înlocuiește intervalele de LUCRATOARE din setul din document.
	cuGresit := func(lucratoare ...Interval) SetReguli {
		s := SetReguli{Intervale: lucratoare}
		for _, iv := range reguliDoc.Intervale {
			if iv.TipZi != Lucratoare {
				s.Intervale = append(s.Intervale, iv)
			}
		}
		return s
	}

	lucr := func(cod, deLa, panaLa string) *Eroare {
		return &Eroare{cod, map[string]string{"tip": "LUCRATOARE", "de_la": deLa, "pana_la": panaLa}}
	}

	cazuri := []struct {
		nume   string
		set    SetReguli
		eroare *Eroare
	}{
		{"gol la mijloc", cuGresit(
			Interval{Lucratoare, 0, 360, "NOAPTE"},
			Interval{Lucratoare, 420, 1440, "ZI"},
		), lucr("reguli.gol", "06:00", "07:00")},
		{"suprapunere", cuGresit(
			Interval{Lucratoare, 0, 420, "NOAPTE"},
			Interval{Lucratoare, 360, 1440, "ZI"},
		), lucr("reguli.suprapunere", "06:00", "07:00")},
		{"nu ajunge la 24:00", cuGresit(
			Interval{Lucratoare, 0, 1320, "ZI"},
		), lucr("reguli.gol", "22:00", "24:00")},
		{"tip de zi fără intervale", cuGresit(), lucr("reguli.gol", "00:00", "24:00")},
		{"interval inversat", cuGresit(
			Interval{Lucratoare, 1320, 360, "NOAPTE"},
		), lucr("reguli.interval_invalid", "22:00", "06:00")},
		{"tip de zi necunoscut",
			SetReguli{Intervale: append([]Interval{{"LUNI", 0, 1440, "ZI"}}, reguliDoc.Intervale...)},
			&Eroare{"reguli.tip_zi_necunoscut", map[string]string{"tip": "LUNI"}}},
	}
	for _, c := range cazuri {
		t.Run(c.nume, func(t *testing.T) {
			err := c.set.Valideaza()
			var e *Eroare
			if !errors.As(err, &e) {
				t.Fatalf("aștepta *Eroare, a primit %v", err)
			}
			if !reflect.DeepEqual(e, c.eroare) {
				t.Errorf("\nobținut:  %+v\nașteptat: %+v", e, c.eroare)
			}
			for _, l := range i18n.Limbi() {
				if text := e.Text(l); text == e.Cod {
					t.Errorf("%s: codul %q nu are traducere", l, e.Cod)
				}
			}
		})
	}
}
