package main

import (
	"errors"

	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

type RaspunsAbsente struct {
	Absente *datebaza.AbsenteAn `json:"absente"`
	Eroare  *MesajEroare        `json:"eroare"`
}

type RaspunsDocumentAbsenta struct {
	Document *datebaza.DocumentAbsenta `json:"document"`
	Eroare   *MesajEroare              `json:"eroare"`
}

func (a *App) CitesteAbsente(an int) RaspunsAbsente { return a.raspunsAbsente(an, nil) }

func (a *App) SalveazaTipAbsenta(initial string, tip datebaza.TipAbsenta, an int) RaspunsAbsente {
	return a.raspunsAbsente(an, datebaza.Nou(a.db).SalveazaTipAbsenta(a.contextLucru(), initial, tip))
}

func (a *App) SalveazaAbsente(zile []datebaza.ZiAbsenta, an int) RaspunsAbsente {
	return a.raspunsAbsente(an, datebaza.Nou(a.db).SalveazaAbsente(a.contextLucru(), zile))
}

func (a *App) InlocuiestePontajCuAbsente(an, luna int, zile []datebaza.ZiAbsenta, celule []datebaza.CelulaSchita) RaspunsSchita {
	err := datebaza.Nou(a.db).InlocuiestePontajCuAbsente(a.contextLucru(), an, luna, zile, celule)
	if err != nil {
		return RaspunsSchita{Eroare: eroareInterfata(err)}
	}
	return a.CitesteSchita(an, luna)
}

func (a *App) CitesteDocumentAbsenta(id string) RaspunsDocumentAbsenta {
	doc, err := datebaza.Nou(a.db).CitesteDocumentAbsenta(a.contextLucru(), id)
	if err != nil {
		mesaj := &MesajEroare{Cod: "date.eroare"}
		var e *i18n.Eroare
		if errors.As(err, &e) {
			mesaj.Cod, mesaj.Param = e.Cod, e.Param
		}
		return RaspunsDocumentAbsenta{Eroare: mesaj}
	}
	return RaspunsDocumentAbsenta{Document: &doc}
}

func (a *App) raspunsAbsente(an int, err error) RaspunsAbsente {
	var absente datebaza.AbsenteAn
	if err == nil {
		absente, err = datebaza.Nou(a.db).CitesteAbsente(a.contextLucru(), an)
	}
	if err != nil {
		mesaj := &MesajEroare{Cod: "date.eroare"}
		var e *i18n.Eroare
		if errors.As(err, &e) {
			mesaj.Cod, mesaj.Param = e.Cod, e.Param
		}
		return RaspunsAbsente{Eroare: mesaj}
	}
	return RaspunsAbsente{Absente: &absente}
}

// PrevizualizeazaEfecteAbsente nu persista date sau jurnal.
func (a *App) PrevizualizeazaEfecteAbsente(zile []datebaza.ZiAbsenta, an, luna int, celule []datebaza.CelulaSchita) RaspunsEfecteAbsente {
	efecte, err := datebaza.Nou(a.db).PrevizualizeazaAbsente(a.contextLucru(), zile, an, luna, celule)
	if err != nil {
		return RaspunsEfecteAbsente{Eroare: eroareInterfata(err)}
	}
	return RaspunsEfecteAbsente{Efecte: efecte}
}

type RaspunsEfecteAbsente struct {
	Efecte []datebaza.EfectAbsenta `json:"efecte"`
	Eroare *MesajEroare            `json:"eroare"`
}

func (a *App) SalveazaCategorieAbsenta(c datebaza.CategorieAbsenta, an int) RaspunsAbsente {
	return a.raspunsAbsente(an, datebaza.Nou(a.db).SalveazaCategorieAbsenta(a.contextLucru(), c))
}
