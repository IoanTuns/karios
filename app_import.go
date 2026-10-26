package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"pontix/internal/datebaza"
)

type RaspunsAnalizaImportAngajati struct {
	Rezultat *datebaza.RezultatAnalizaImportAngajati `json:"rezultat"`
	Eroare   *MesajEroare                            `json:"eroare"`
}

type RaspunsExecutaImportAngajati struct {
	Date      *datebaza.Date `json:"date"`
	Importati int            `json:"importati"`
	Eroare    *MesajEroare   `json:"eroare"`
}

type RaspunsAnalizaImportSarbatori struct {
	Rezultat *datebaza.RezultatAnalizaImportSarbatori `json:"rezultat"`
	Eroare   *MesajEroare                             `json:"eroare"`
}

type RaspunsExecutaImportSarbatori struct {
	Calendar  *datebaza.CalendarAn `json:"calendar"`
	Importate int                  `json:"importate"`
	Eroare    *MesajEroare         `json:"eroare"`
}

type RaspunsDescarcaSablon struct {
	ContinutBase64 string       `json:"continut_base64"`
	NumeFisier     string       `json:"nume_fisier"`
	Eroare         *MesajEroare `json:"eroare"`
}

func (a *App) AnalizeazaImportAngajati(continutBase64 string, numeFisier string) RaspunsAnalizaImportAngajati {
	dateBytes, err := base64.StdEncoding.DecodeString(continutBase64)
	if err != nil {
		return RaspunsAnalizaImportAngajati{Eroare: eroareInterfata(fmt.Errorf("fisier_invalid: %w", err))}
	}
	rez, err := datebaza.Nou(a.db).AnalizeazaImportAngajati(a.contextLucru(), dateBytes, numeFisier)
	if err != nil {
		return RaspunsAnalizaImportAngajati{Eroare: eroareInterfata(err)}
	}
	return RaspunsAnalizaImportAngajati{Rezultat: &rez}
}

func (a *App) ExecutaImportAngajati(cerere datebaza.CerereExecutaImportAngajati) RaspunsExecutaImportAngajati {
	n, err := datebaza.Nou(a.db).ExecutaImportAngajati(a.contextLucru(), cerere.Randuri, cerere.ActualizeazaExistente)
	if err != nil {
		return RaspunsExecutaImportAngajati{Eroare: eroareInterfata(err)}
	}
	date, err := datebaza.Nou(a.db).Citeste(a.contextLucru())
	if err != nil {
		return RaspunsExecutaImportAngajati{Eroare: eroareInterfata(err)}
	}
	return RaspunsExecutaImportAngajati{Date: &date, Importati: n}
}

func (a *App) AnalizeazaImportSarbatori(continutBase64 string, numeFisier string) RaspunsAnalizaImportSarbatori {
	dateBytes, err := base64.StdEncoding.DecodeString(continutBase64)
	if err != nil {
		return RaspunsAnalizaImportSarbatori{Eroare: eroareInterfata(fmt.Errorf("fisier_invalid: %w", err))}
	}
	rez, err := datebaza.Nou(a.db).AnalizeazaImportSarbatori(a.contextLucru(), dateBytes, numeFisier)
	if err != nil {
		return RaspunsAnalizaImportSarbatori{Eroare: eroareInterfata(err)}
	}
	return RaspunsAnalizaImportSarbatori{Rezultat: &rez}
}

func (a *App) ExecutaImportSarbatori(cerere datebaza.CerereExecutaImportSarbatori, an int) RaspunsExecutaImportSarbatori {
	n, err := datebaza.Nou(a.db).ExecutaImportSarbatori(a.contextLucru(), cerere.Randuri, cerere.InlocuiesteExistente)
	if err != nil {
		return RaspunsExecutaImportSarbatori{Eroare: eroareInterfata(err)}
	}
	cal, err := datebaza.Nou(a.db).CitesteCalendar(a.contextLucru(), an)
	if err != nil {
		return RaspunsExecutaImportSarbatori{Eroare: eroareInterfata(err)}
	}
	return RaspunsExecutaImportSarbatori{Calendar: &cal, Importate: n}
}

func (a *App) DescarcaSablon(tip string, format string) RaspunsDescarcaSablon {
	var (
		b          []byte
		numeFisier string
		err        error
	)
	switch strings.ToLower(tip) {
	case "angajati":
		b, numeFisier, err = datebaza.GenereazaSablonAngajati(format)
	case "sarbatori":
		b, numeFisier, err = datebaza.GenereazaSablonSarbatori(format)
	default:
		return RaspunsDescarcaSablon{Eroare: eroareInterfata(fmt.Errorf("tip_sablon_necunoscut"))}
	}
	if err != nil {
		return RaspunsDescarcaSablon{Eroare: eroareInterfata(err)}
	}
	return RaspunsDescarcaSablon{
		ContinutBase64: base64.StdEncoding.EncodeToString(b),
		NumeFisier:     numeFisier,
	}
}
