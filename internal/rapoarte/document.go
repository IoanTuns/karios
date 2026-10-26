package rapoarte

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	raportAbsente "pontix/internal/rapoarte/absente"
	"pontix/internal/rapoarte/detaliat"
	"pontix/internal/rapoarte/individual"
	"pontix/internal/rapoarte/pontaj"
	"pontix/internal/versiune"
)

type Cerere struct {
	ModelID       string `json:"model_id"`
	Revizie       string `json:"revizie"`
	An            int    `json:"an"`
	Luna          int    `json:"luna"`
	SectieID      string `json:"sectie_id"`
	DepartamentID string `json:"departament_id"`
	AngajatID     string `json:"angajat_id"`
	Limba         string `json:"limba"`
}
type Coloana struct {
	Cod     string `json:"cod"`
	Titlu   string `json:"titlu"`
	Pondere int    `json:"pondere"`
	Culoare string `json:"culoare"`
}
type Rand struct {
	Valori []string `json:"valori"`
	Culori []string `json:"culori"`
}
type Document struct {
	Titlu          string             `json:"titlu"`
	Sursa          string             `json:"sursa"`
	VersiuneRaport int                `json:"versiune_raport"`
	RevizieModel   string             `json:"revizie_model"`
	Pagina         Pagina             `json:"pagina"`
	Antet          []string           `json:"antet"`
	AntetImagini   map[string]Imagine `json:"antet_imagini,omitempty"`
	Introducere    string             `json:"introducere"`
	Coloane        []Coloana          `json:"coloane"`
	Randuri        []Rand             `json:"randuri"`
	Subsol         []string           `json:"subsol"`
	SubsolImagini  map[string]Imagine `json:"subsol_imagini,omitempty"`
	Metadate       string             `json:"metadate"`
	Nota           string             `json:"nota"`
	FaraDate       string             `json:"fara_date"`
	Paginare       string             `json:"paginare"`
	GeneratLa      time.Time          `json:"-"`
}

func textJSON(v string, limba i18n.Limba) string {
	var t i18n.Text
	_ = json.Unmarshal([]byte(v), &t)
	return t.In(limba)
}
func optional(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func dataAfisata(data string, limba i18n.Limba) string {
	d, err := time.Parse(time.DateOnly, data)
	if err != nil {
		return data
	}
	if limba == i18n.Engleza {
		return d.Format("01/02/2006")
	}
	return d.Format("02.01.2006")
}

// Documentul citește definiția și datele din aceeași tranzacție de citire.
func (s *Serviciu) Pregateste(ctx context.Context, c Cerere, acum time.Time) (Document, error) {
	if c.An < 1 || c.An > 9999 || c.Luna < 1 || c.Luna > 12 || (c.Limba != "ro" && c.Limba != "en") {
		return Document{}, eroare("rapoarte.parametri_invalizi")
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return Document{}, traduce(err)
	}
	defer tx.Rollback()
	q := acces.New(tx)
	r, err := q.CitesteDefinitieRaport(ctx, c.ModelID)
	if err != nil {
		return Document{}, eroare("rapoarte.inexistent")
	}
	d, err := decodeaza(r)
	if err != nil {
		return Document{}, traduce(err)
	}
	if !d.Activ {
		return Document{}, eroare("rapoarte.inactiv")
	}
	if c.Revizie != "" && c.Revizie != d.Revizie {
		return Document{}, eroare("rapoarte.conflict")
	}
	if err := Valideaza(d); err != nil {
		return Document{}, err
	}
	limba := i18n.Limba(c.Limba)
	t := func(k string) string { return i18n.T(limba, "rapoarte."+k, nil) }
	oameni, err := q.ListaAngajati(ctx)
	if err != nil {
		return Document{}, traduce(err)
	}
	sectii, err := q.ListaSectii(ctx)
	if err != nil {
		return Document{}, traduce(err)
	}
	departamente, err := q.ListaDepartamente(ctx)
	if err != nil {
		return Document{}, traduce(err)
	}
	numeSectii, numeDepartamente := map[string]string{}, map[string]string{}
	for _, v := range sectii {
		numeSectii[v.ID] = textJSON(v.Denumire, limba)
	}
	for _, v := range departamente {
		numeDepartamente[v.ID] = textJSON(v.Denumire, limba)
		if v.ID == c.DepartamentID && c.SectieID != "" && c.SectieID != v.SectieID {
			return Document{}, eroare("rapoarte.parametri_invalizi")
		}
	}
	if (c.SectieID != "" && numeSectii[c.SectieID] == "") || (c.DepartamentID != "" && numeDepartamente[c.DepartamentID] == "") {
		return Document{}, eroare("rapoarte.parametri_invalizi")
	}
	var persoana string
	for _, a := range oameni {
		if a.ID == c.AngajatID {
			persoana = a.Nume + " " + a.Prenume
			if (c.SectieID != "" && optional(a.SectieID) != c.SectieID) || (c.DepartamentID != "" && optional(a.DepartamentID) != c.DepartamentID) {
				return Document{}, eroare("rapoarte.parametri_invalizi")
			}
		}
	}
	if (c.AngajatID != "" && persoana == "") || (d.Continut.Sursa == "FISA_INDIVIDUALA" && c.AngajatID == "") {
		return Document{}, eroare("rapoarte.alege_angajat")
	}
	deLa := fmt.Sprintf("%04d-%02d-01", c.An, c.Luna)
	panaLa := time.Date(c.An, time.Month(c.Luna)+1, 0, 12, 0, 0, 0, time.UTC).Format(time.DateOnly)
	perioada := t(fmt.Sprintf("luna_%02d", c.Luna)) + " " + strconv.Itoa(c.An)
	if d.Continut.Sursa == "ABSENTE_ANUALE" {
		deLa = fmt.Sprintf("%04d-01-01", c.An)
		panaLa = fmt.Sprintf("%04d-12-31", c.An)
		perioada = strconv.Itoa(c.An)
	}
	ciorne, err := q.ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: deLa, DataPanaLa: panaLa})
	if err != nil {
		return Document{}, traduce(err)
	}
	absente, err := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: deLa, PanaLa: panaLa})
	if err != nil {
		return Document{}, traduce(err)
	}
	tipuri, err := q.ListaTipuriAbsenta(ctx)
	if err != nil {
		return Document{}, traduce(err)
	}
	schimburi, err := q.ListaSchimburi(ctx)
	if err != nil {
		return Document{}, traduce(err)
	}
	anCalendar := int64(c.An)
	sarbatori, err := q.ListaSarbatori(ctx, &anCalendar)
	if err != nil {
		return Document{}, traduce(err)
	}
	descrieri, categorii := map[string]string{}, map[string]string{}
	for _, v := range tipuri {
		descrieri[v.Cod] = textJSON(v.Denumire, limba)
		categorii[v.Cod] = i18n.T(limba, "absente."+v.Categorie, nil)
	}
	for _, v := range schimburi {
		descrieri[v.Cod] = textJSON(v.Denumire, limba)
	}
	sarbatoare := map[string]bool{}
	for _, v := range sarbatori {
		sarbatoare[v.Data] = true
	}
	culoareZi := func(data string) string {
		if sarbatoare[data] {
			return "SARBATOARE"
		}
		d, _ := time.Parse(time.DateOnly, data)
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			return "WEEKEND"
		}
		return ""
	}
	celule := map[string]string{}
	cuDate := map[string]bool{}
	for _, v := range ciorne {
		celule[v.AngajatID+"/"+v.Data] = v.Valoare
		if v.Valoare != "" {
			cuDate[v.AngajatID] = true
		}
	}
	for _, v := range absente {
		if v.Activ == 1 {
			celule[v.AngajatID+"/"+v.Data] = v.Cod
			cuDate[v.AngajatID] = true
		}
	}
	ore := map[string]string{}
	totaluri := map[string]map[string]calcul.ValoarePontaj{}
	var configuratie calcul.ConfiguratiePontaj
	if d.Continut.Sursa != "ABSENTE_ANUALE" {
		rezultat, cfg, err := datebaza.CalculeazaPontaj(ctx, q, c.An, c.Luna, nil)
		if err != nil {
			return Document{}, traduce(err)
		}
		configuratie = cfg
		for _, z := range rezultat.Celule {
			ore[z.AngajatID+"/"+z.Data] = calcul.TextValoarePontaj(calcul.ValoarePontaj{Sutimi: z.Minute * 100, Unitate: "ORE", Definita: z.Calculabil})
			cuDate[z.AngajatID] = true
		}
		for _, r := range rezultat.Totaluri {
			totaluri[r.AngajatID] = r.Valori
		}
	}
	local, _ := time.LoadLocation("Europe/Bucharest")
	valori := map[string]string{"titlu": d.Denumire.In(limba), "an": strconv.Itoa(c.An), "luna": t(fmt.Sprintf("luna_%02d", c.Luna)), "perioada": perioada, "sectie": numeSectii[c.SectieID], "departament": numeDepartamente[c.DepartamentID], "angajat": persoana, "generat_la": acum.In(local).Format("02.01.2006 15:04")}
	inlocuieste := func(v string) string {
		return variabila.ReplaceAllStringFunc(v, func(k string) string { return valori[strings.TrimSuffix(strings.TrimPrefix(k, "{{"), "}}")] })
	}
	bloc := func(b Bloc) []string {
		return []string{inlocuieste(b.Stanga.In(limba)), inlocuieste(b.Centru.In(limba)), inlocuieste(b.Dreapta.In(limba))}
	}
	doc := Document{Titlu: d.Denumire.In(limba), Sursa: d.Continut.Sursa, RevizieModel: d.Revizie, Pagina: d.Pagina, Antet: bloc(d.Antet), Introducere: inlocuieste(d.Continut.Text.In(limba)), Subsol: bloc(d.Subsol), Coloane: []Coloana{}, Randuri: []Rand{}, Nota: t("nota_schita"), FaraDate: t("fara_date"), Paginare: t("paginare"), GeneratLa: acum}
	doc.AntetImagini, doc.SubsolImagini = d.Antet.Imagini, d.Subsol.Imagini
	switch d.Continut.Sursa {
	case "PONTAJ_LUNAR":
		doc.VersiuneRaport = pontaj.Versiune
	case "PONTAJ_DETALIAT":
		doc.VersiuneRaport = detaliat.Versiune
	case "FISA_INDIVIDUALA":
		doc.VersiuneRaport = individual.Versiune
	case "ABSENTE_ANUALE":
		doc.VersiuneRaport = raportAbsente.Versiune
	}
	ponderi := map[string]int{"angajat": 28, "marca": 10, "sectie": 18, "departament": 18, "functie": 22, "data": 14, "zi": 10, "cod": 10, "ore": 10, "descriere": 40, "categorie": 22}
	for _, cod := range d.Continut.Coloane {
		if cod == "totaluri" {
			for _, col := range calcul.ColoanePontaj() {
				titlu := i18n.T(limba, "pontaj.total_"+col, nil)
				unitate := calcul.UnitateColoana(col)
				for _, r := range configuratie.Reguli {
					if r.Coloana == col && r.Sursa != "NEDEFINITA" {
						if r.Sursa == "ZILE_ABSENTA" {
							unitate = "ZILE"
						} else {
							unitate = "ORE"
						}
					}
				}
				if unitate != "" {
					cheie := "unitate_ore"
					if unitate == "ZILE" {
						cheie = "unitate_zile"
					}
					titlu += " (" + i18n.T(limba, "pontaj."+cheie, nil) + ")"
				}
				doc.Coloane = append(doc.Coloane, Coloana{Cod: "total_" + col, Titlu: titlu, Pondere: 6})
			}
		} else if cod == "zile" {
			for zi := 1; zi <= 31; zi++ {
				data := fmt.Sprintf("%04d-%02d-%02d", c.An, c.Luna, zi)
				if data > panaLa {
					break
				}
				doc.Coloane = append(doc.Coloane, Coloana{Cod: data, Titlu: strconv.Itoa(zi), Pondere: 3, Culoare: culoareZi(data)})
			}
		} else {
			doc.Coloane = append(doc.Coloane, Coloana{Cod: cod, Titlu: t("coloana_" + cod), Pondere: ponderi[cod]})
		}
	}
	for _, a := range oameni {
		if (c.SectieID != "" && optional(a.SectieID) != c.SectieID) || (c.DepartamentID != "" && optional(a.DepartamentID) != c.DepartamentID) || (c.AngajatID != "" && a.ID != c.AngajatID) {
			continue
		}
		if (d.Continut.Sursa == "PONTAJ_LUNAR" || d.Continut.Sursa == "PONTAJ_DETALIAT") && a.Activ == 0 && !cuDate[a.ID] {
			continue
		}
		identitate := map[string]string{"angajat": a.Nume + " " + a.Prenume, "marca": a.Marca, "sectie": numeSectii[optional(a.SectieID)], "departament": numeDepartamente[optional(a.DepartamentID)], "functie": optional(a.Functie)}
		if totaluri[a.ID] == nil && d.Continut.Sursa == "PONTAJ_DETALIAT" {
			gol, err := calcul.CalculeazaPontaj(deLa, panaLa, configuratie, nil, []calcul.IntrarePontaj{{AngajatID: a.ID, Data: deLa, Valoare: "@M0"}}, nil)
			if err != nil {
				return Document{}, traduce(err)
			}
			totaluri[a.ID] = gol.Totaluri[0].Valori
		}
		adauga := func(data, cod string) {
			rand := Rand{Valori: []string{}, Culori: []string{}}
			for _, col := range doc.Coloane {
				v, culoare := identitate[col.Cod], ""
				switch col.Cod {
				case "data":
					v = dataAfisata(data, limba)
				case "zi":
					zi, _ := time.Parse(time.DateOnly, data)
					v = t(fmt.Sprintf("zi_%d", zi.Weekday()))
				case "cod":
					v = cod
				case "ore":
					v = ore[a.ID+"/"+data]
				case "descriere":
					v = descrieri[cod]
					if _, manual := calcul.OreManuale(cod); manual && v == "" {
						v = t("ore_manuale")
					}
					if cod != "" && v == "" {
						v = t("cod_neinterpretat")
					}
				case "categorie":
					v = categorii[cod]
				}
				if col.Culoare != "" {
					culoare = col.Culoare
				}
				if len(col.Cod) == 10 && col.Cod[4] == '-' {
					v = ore[a.ID+"/"+col.Cod]
					if categorii[celule[a.ID+"/"+col.Cod]] != "" {
						culoare = "ABSENTA"
					}
				}
				if strings.HasPrefix(col.Cod, "total_") {
					v = calcul.TextValoarePontaj(totaluri[a.ID][strings.TrimPrefix(col.Cod, "total_")])
				}
				rand.Valori = append(rand.Valori, v)
				rand.Culori = append(rand.Culori, culoare)
			}
			doc.Randuri = append(doc.Randuri, rand)
		}
		switch d.Continut.Sursa {
		case "PONTAJ_LUNAR", "PONTAJ_DETALIAT":
			adauga("", "")
		case "FISA_INDIVIDUALA":
			for zi := 1; zi <= 31; zi++ {
				data := fmt.Sprintf("%04d-%02d-%02d", c.An, c.Luna, zi)
				if data > panaLa {
					break
				}
				adauga(data, celule[a.ID+"/"+data])
			}
		case "ABSENTE_ANUALE":
			for _, v := range absente {
				if v.Activ == 1 && v.AngajatID == a.ID {
					adauga(v.Data, v.Cod)
				}
			}
		}
		if len(doc.Randuri) > 10000 {
			return Document{}, eroare("rapoarte.prea_mare")
		}
	}
	stareReguli := t("reguli_neaplicate")
	if configuratie.Revizie != "" {
		stareReguli = t("salarizare_neaplicata")
	}
	doc.Metadate = i18n.T(limba, "rapoarte.metadate", map[string]string{"nume": doc.Titlu, "raport": strconv.Itoa(doc.VersiuneRaport), "model": d.Revizie, "aplicatie": versiune.Aplicatie, "reguli": stareReguli, "generat": valori["generat_la"], "utilizator": t("actor_local")})
	if configuratie.Revizie != "" {
		doc.Metadate += " · " + i18n.T(limba, "pontaj.revizie_reguli", map[string]string{"revizie": configuratie.Revizie})
	}
	if err := tx.Commit(); err != nil {
		return Document{}, traduce(err)
	}
	return doc, nil
}
