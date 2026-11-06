package rapoarte

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"pontix/internal/audit"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"pontix/internal/identitate"
	raportAbsente "pontix/internal/rapoarte/absente"
	"pontix/internal/rapoarte/detaliat"
	raportGarzi "pontix/internal/rapoarte/garzi"
	"pontix/internal/rapoarte/individual"
	"pontix/internal/rapoarte/pontaj"
	"pontix/internal/versiune"
)

type Cerere struct {
	GrupID        string `json:"grup_id"`
	FunctieID     string `json:"functie_id"`
	ContractID    string `json:"contract_id"`
	Pozitie       string `json:"pozitie"`
	Garda         string `json:"garda"`
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
	Aldin  bool     `json:"aldin,omitempty"`
	Grup   bool     `json:"grup,omitempty"`
	Valori []string `json:"valori"`
	Culori []string `json:"culori"`
}
type Document struct {
	Stiluri       map[string]StilText `json:"stiluri,omitempty"`
	Personalizat  bool                `json:"personalizat,omitempty"`
	valoriRanduri []map[string]string
	// legendaPentru reface legenda pentru un set de coloane (rapoartele standard
	// folosite din SQL afiseaza doar explicatiile coloanelor alese, D090).
	legendaPentru func([]Coloana) []string
	// zileLuna: zilele lunii raportului, pentru ascunderea coloanelor z29-z31 (D090).
	zileLuna int
	// variabileStandard: variabilele textelor raportului standard (luna ca nume,
	// perioada localizata, persoana aleasa), preluate de rapoartele SQL (D090).
	variabileStandard map[string]string
	VariabileRaport   map[string]string  `json:"variabile_raport,omitempty"`
	Flexibil          *RaportFlexibil    `json:"flexibil,omitempty"`
	Schema            []TabelRaport      `json:"schema,omitempty"`
	Limba             i18n.Limba         `json:"-"`
	Legenda           []string           `json:"legenda,omitempty"`
	Titlu             string             `json:"titlu"`
	Sursa             string             `json:"sursa"`
	VersiuneRaport    int                `json:"versiune_raport"`
	RevizieModel      string             `json:"revizie_model"`
	Pagina            Pagina             `json:"pagina"`
	Antet             []string           `json:"antet"`
	AntetImagini      map[string]Imagine `json:"antet_imagini,omitempty"`
	Introducere       string             `json:"introducere"`
	Coloane           []Coloana          `json:"coloane"`
	Randuri           []Rand             `json:"randuri"`
	Subsol            []string           `json:"subsol"`
	SubsolImagini     map[string]Imagine `json:"subsol_imagini,omitempty"`
	Metadate          string             `json:"metadate"`
	Nota              string             `json:"nota"`
	FaraDate          string             `json:"fara_date"`
	Paginare          string             `json:"paginare"`
	GeneratLa         time.Time          `json:"-"`
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
	if esteFlexibil(d.Continut.Sursa) {
		return preparaFlexibil(ctx, tx, d, c, acum)
	}
	doc, err := pregatesteStandard(ctx, tx, q, d, c, acum)
	if err != nil {
		return Document{}, err
	}
	if err := tx.Commit(); err != nil {
		return Document{}, traduce(err)
	}
	return doc, nil
}

// ponderiColoane: latimile relative ale coloanelor rapoartelor standard (si in conversie).
var ponderiColoane = map[string]int{"angajat": 28, "marca": 10, "sectie": 18, "departament": 18, "functie": 22, "contract": 22, "norma_contract": 18, "data": 14, "zi": 10, "cod": 10, "ore": 10, "descriere": 40, "categorie": 22, "total_garzi": 10, "total_ore": 10}

// pregatesteStandard calculeaza documentul unei surse standard (pontaj lunar,
// detaliat, fisa individuala, absente anuale, grafic garzi) in tranzactia data,
// fara sa o inchida. O folosesc motorul vechi si setul de date SQL (D090).
func pregatesteStandard(ctx context.Context, tx *sql.Tx, q *acces.Queries, d Definitie, c Cerere, acum time.Time) (Document, error) {
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
	if (c.Pozitie != "" && c.Pozitie != "CONDUCERE" && c.Pozitie != "EXECUTIE") || (c.Garda != "" && c.Garda != "DA" && c.Garda != "NU") {
		return Document{}, eroare("rapoarte.parametri_invalizi")
	}
	profile, asocieri, err := profileRapoarte(ctx, tx, limba)
	if err != nil {
		return Document{}, traduce(err)
	}
	prefixe, err := prefixeRapoarte(ctx, tx, limba)
	if err != nil {
		return Document{}, traduce(err)
	}
	var persoana string
	for _, a := range oameni {
		if a.ID == c.AngajatID {
			persoana = identitate.CuPrefix(prefixe[a.ID], identitate.Nume(a.Nume, a.Prenume, profile[a.ID].Garda || asocieri[a.ID] != "", limba))
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
	denumiriCategorii, err := q.DenumiriCategoriiAbsente(ctx)
	if err != nil {
		return Document{}, traduce(err)
	}
	descrieri, categorii := map[string]string{}, map[string]string{}
	for _, v := range tipuri {
		descrieri[v.Cod] = textJSON(v.Denumire, limba)
		categorii[v.Cod] = i18n.T(limba, "absente."+v.Categorie, nil)
		if denumiriCategorii[v.Cod] != "" {
			categorii[v.Cod] = textJSON(denumiriCategorii[v.Cod], limba)
		}
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
		if coduri := calcul.CoduriPrograme(v.Valoare); len(coduri) > 1 {
			texte, cat := []string{}, []string{}
			for _, cod := range coduri {
				texte = append(texte, descrieri[cod])
				cat = append(cat, categorii[cod])
			}
			descrieri[v.Valoare] = strings.Join(texte, " + ")
			categorii[v.Valoare] = strings.Join(cat, " + ")
		}
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
	legaturiGarzi := asocieri
	var garziRaport []acces.Planificare
	if d.Continut.Sursa == "PONTAJ_DETALIAT" {
		rows, err := tx.QueryContext(ctx, "SELECT angajat_id,principal_id FROM angajati_garzi")
		if err != nil {
			return Document{}, traduce(err)
		}
		for rows.Next() {
			var id, principal string
			if err := rows.Scan(&id, &principal); err != nil {
				rows.Close()
				return Document{}, traduce(err)
			}
			legaturiGarzi[id] = principal
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return Document{}, traduce(err)
		}
		inceput, _ := time.Parse(time.DateOnly, deLa)
		plan, err := q.ListaPlanificare(ctx, acces.ListaPlanificareParams{DataDeLa: inceput.AddDate(0, 0, -2).Format(time.DateOnly), DataPanaLa: panaLa})
		if err != nil {
			return Document{}, traduce(err)
		}
		for _, g := range plan {
			if legaturiGarzi[g.AngajatID] != "" && g.Sectiune == "GARZI" && g.Activ == 1 {
				garziRaport = append(garziRaport, g)
			}
		}
		// Ordinea initiala a persoanelor principale ramane stabila; subtipurile urmeaza imediat dupa ele.
		pozitie := map[string]int{}
		for i, a := range oameni {
			pozitie[a.ID] = i
		}
		sort.SliceStable(oameni, func(i, j int) bool {
			a, b := oameni[i].ID, oameni[j].ID
			if a == b {
				return false
			}
			pa, pb := a, b
			if legaturiGarzi[a] != "" {
				pa = legaturiGarzi[a]
			}
			if legaturiGarzi[b] != "" {
				pb = legaturiGarzi[b]
			}
			if pa != pb {
				return pozitie[pa] < pozitie[pb]
			}
			if a == pa {
				return true
			}
			if b == pb {
				return false
			}
			return pozitie[a] < pozitie[b]
		})
	}
	if d.Continut.Sursa != "FISA_INDIVIDUALA" {
		legaturiGarzi = asocieri
		ordoneazaProfesional(oameni, profile, asocieri)
	}
	ore := map[string]string{}
	totaluri := map[string]map[string]calcul.ValoarePontaj{}
	var configuratie calcul.ConfiguratiePontaj
	if d.Continut.Sursa != "ABSENTE_ANUALE" && d.Continut.Sursa != "GRAFIC_GARZI" {
		rezultat, cfg, err := datebaza.CalculeazaPontajCuGarzi(ctx, q, c.An, c.Luna, garziRaport)
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
	planificari, err := q.ListaPlanificare(ctx, acces.ListaPlanificareParams{DataDeLa: deLa, DataPanaLa: panaLa})
	if err != nil {
		return Document{}, traduce(err)
	}
	codSchimb := map[string]string{}
	durataSchimb := map[string]int64{}
	for _, sc := range schimburi {
		codSchimb[sc.ID] = sc.Cod
		durataSchimb[sc.ID] = sc.DurataMinute
	}
	planificariAngajat := map[string]map[string][]acces.Planificare{}
	cuPlanificare := map[string]bool{}
	totaluriAngajatGarzi := map[string]int{}
	totaluriAngajatMinute := map[string]int64{}
	prezenteZi := map[string]int{}
	for _, p := range planificari {
		if d.Continut.Sursa == "GRAFIC_GARZI" && p.Sectiune != "GARZI" {
			continue
		}
		if p.Activ == 1 {
			if planificariAngajat[p.AngajatID] == nil {
				planificariAngajat[p.AngajatID] = map[string][]acces.Planificare{}
			}
			planificariAngajat[p.AngajatID][p.Data] = append(planificariAngajat[p.AngajatID][p.Data], p)
			cuPlanificare[p.AngajatID] = true
			totaluriAngajatGarzi[p.AngajatID]++
			totaluriAngajatMinute[p.AngajatID] += durataSchimb[p.SchimbID]
			prezenteZi[p.Data]++
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
	doc := Document{Limba: limba, Titlu: d.Denumire.In(limba), Sursa: d.Continut.Sursa, RevizieModel: d.Revizie, Pagina: d.Pagina, Antet: bloc(d.Antet), Introducere: inlocuieste(d.Continut.Text.In(limba)), Subsol: bloc(d.Subsol), Coloane: []Coloana{}, Randuri: []Rand{}, Nota: t("nota_schita"), FaraDate: t("fara_date"), Paginare: t("paginare"), GeneratLa: acum}
	doc.variabileStandard = valori
	if len(garziRaport) > 0 {
		doc.Nota += " " + t("nota_garzi")
	}
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
	case "GRAFIC_GARZI":
		doc.VersiuneRaport = raportGarzi.Versiune
	}
	ponderi := ponderiColoane
	for _, cod := range d.Continut.Coloane {
		if cod == "totaluri" {
			for _, col := range calcul.ColoanePontaj() {
				if col == "ore_lucrate" && d.Continut.Sursa == "PONTAJ_DETALIAT" {
					col = "ore_lucratoare"
				}
				if col == "ore_norma" {
					continue
				}
				titlu := i18n.T(limba, "pontaj.total_"+col, nil)
				unitate := calcul.UnitateColoana(col)
				if col == "ore_lucratoare" || col == "L_FIN" {
					unitate = "ORE"
				}
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
			if cod == "marca" && d.Continut.Sursa == "PONTAJ_DETALIAT" {
				doc.Coloane = append(doc.Coloane, Coloana{Cod: "total_L_INCEPUT", Titlu: i18n.T(limba, "pontaj.total_L_INCEPUT", nil), Pondere: 6})
			}
		}
	}
	familie := func(id string) string {
		if p := legaturiGarzi[id]; p != "" {
			return p
		}
		return id
	}
	inclusi := map[string]bool{}
	grupCurent := "__initial"
	areGrupuri := false
	for _, p := range profile {
		if p.Grup != "" {
			areGrupuri = true
		}
	}
	for _, a := range oameni {
		if !profile[a.ID].corespunde(c) {
			continue
		}
		selectat := c.AngajatID == "" || a.ID == c.AngajatID || (d.Continut.Sursa == "PONTAJ_DETALIAT" && familie(a.ID) == familie(c.AngajatID))
		if (c.SectieID != "" && optional(a.SectieID) != c.SectieID) || (c.DepartamentID != "" && optional(a.DepartamentID) != c.DepartamentID) || !selectat {
			continue
		}
		if (d.Continut.Sursa == "PONTAJ_LUNAR" || d.Continut.Sursa == "PONTAJ_DETALIAT") && a.Activ == 0 && !cuDate[a.ID] {
			continue
		}
		if d.Continut.Sursa == "GRAFIC_GARZI" && a.Activ == 0 && !cuPlanificare[a.ID] {
			continue
		}
		identitate := map[string]string{"angajat": identitate.CuPrefix(prefixe[a.ID], identitate.Nume(a.Nume, a.Prenume, profile[a.ID].Garda || asocieri[a.ID] != "", limba)), "marca": a.Marca, "sectie": numeSectii[optional(a.SectieID)], "departament": numeDepartamente[optional(a.DepartamentID)], "functie": optional(a.Functie), "contract": profile[a.ID].NumeContract, "norma_contract": profile[a.ID].NormaContract}
		if totaluri[a.ID] == nil && d.Continut.Sursa == "PONTAJ_DETALIAT" {
			gol, err := calcul.CalculeazaPontaj(deLa, panaLa, configuratie, nil, []calcul.IntrarePontaj{{AngajatID: a.ID, Data: deLa, Valoare: "@M0"}}, nil)
			if err != nil {
				return Document{}, traduce(err)
			}
			totaluri[a.ID] = gol.Totaluri[0].Valori
		}
		adauga := func(data, cod string) {
			inclusi[a.ID] = true
			if areGrupuri && d.Continut.Sursa != "FISA_INDIVIDUALA" {
				id := a.ID
				if asocieri[id] != "" {
					id = asocieri[id]
				}
				p := profile[id]
				if grupCurent != p.Grup {
					grupCurent = p.Grup
					nume := p.NumeGrup
					if nume == "" {
						nume = i18n.T(limba, "structura.fara_grup", nil)
					}
					v := make([]string, len(doc.Coloane))
					cul := make([]string, len(doc.Coloane))
					v[0] = nume
					doc.Randuri = append(doc.Randuri, Rand{Grup: true, Valori: v, Culori: cul})
				}
			}

			rand := Rand{Valori: []string{}, Culori: []string{}}
			for _, col := range doc.Coloane {
				v, culoare := identitate[col.Cod], ""
				switch col.Cod {
				case "total_garzi":
					if g := totaluriAngajatGarzi[a.ID]; g > 0 {
						v = strconv.Itoa(g)
					} else {
						v = "—"
					}
				case "total_ore":
					if m := totaluriAngajatMinute[a.ID]; m > 0 {
						if m%60 == 0 {
							v = strconv.FormatInt(m/60, 10)
						} else {
							v = fmt.Sprintf("%.2f", float64(m)/60.0)
						}
					} else {
						v = "—"
					}
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
					if d.Continut.Sursa == "GRAFIC_GARZI" {
						if poz, exista := planificariAngajat[a.ID][col.Cod]; exista {
							coduri := []string{}
							for _, p := range poz {
								coduri = append(coduri, codSchimb[p.SchimbID])
							}
							v = calcul.ValoarePrograme(coduri)
						} else {
							v = "—"
						}
					} else {
						v = ore[a.ID+"/"+col.Cod]
						if categorii[celule[a.ID+"/"+col.Cod]] != "" {
							culoare = "ABSENTA"
						}
					}
				}
				if strings.HasPrefix(col.Cod, "total_") && d.Continut.Sursa != "GRAFIC_GARZI" {
					v = calcul.TextValoarePontaj(totaluri[a.ID][strings.TrimPrefix(col.Cod, "total_")])
				}
				rand.Valori = append(rand.Valori, v)
				rand.Culori = append(rand.Culori, culoare)
			}
			doc.Randuri = append(doc.Randuri, rand)
		}
		switch d.Continut.Sursa {
		case "PONTAJ_LUNAR", "PONTAJ_DETALIAT", "GRAFIC_GARZI":
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
	if d.Continut.Sursa == "GRAFIC_GARZI" && len(doc.Randuri) > 0 {
		prezenteZi = map[string]int{}
		for id, zile := range planificariAngajat {
			if inclusi[id] {
				for data := range zile {
					prezenteZi[data]++
				}
			}
		}
		randTotal := Rand{Valori: []string{}, Culori: []string{}}
		totalToateGarzi := 0
		totalToateMinute := int64(0)
		for id, g := range totaluriAngajatGarzi {
			if !inclusi[id] {
				continue
			}
			totalToateGarzi += g
		}
		for id, m := range totaluriAngajatMinute {
			if !inclusi[id] {
				continue
			}
			totalToateMinute += m
		}
		for _, col := range doc.Coloane {
			v, culoare := "", col.Culoare
			switch col.Cod {
			case "angajat":
				v = t("total_prezente_zi")
			case "marca", "functie", "sectie", "departament":
				v = ""
			case "total_garzi":
				if totalToateGarzi > 0 {
					v = strconv.Itoa(totalToateGarzi)
				} else {
					v = "0"
				}
			case "total_ore":
				if totalToateMinute > 0 {
					if totalToateMinute%60 == 0 {
						v = strconv.FormatInt(totalToateMinute/60, 10)
					} else {
						v = fmt.Sprintf("%.2f", float64(totalToateMinute)/60.0)
					}
				} else {
					v = "0"
				}
			default:
				if len(col.Cod) == 10 && col.Cod[4] == '-' {
					nr := prezenteZi[col.Cod]
					v = strconv.Itoa(nr)
				}
			}
			randTotal.Valori = append(randTotal.Valori, v)
			randTotal.Culori = append(randTotal.Culori, culoare)
		}
		doc.Randuri = append(doc.Randuri, randTotal)
	}
	stareReguli := t("reguli_neaplicate")
	if configuratie.Revizie != "" {
		stareReguli = t("salarizare_neaplicata")
	}
	doc.Metadate = i18n.T(limba, "rapoarte.metadate", map[string]string{"nume": doc.Titlu, "raport": strconv.Itoa(doc.VersiuneRaport), "model": d.Revizie, "aplicatie": versiune.Aplicatie, "reguli": stareReguli, "generat": valori["generat_la"], "utilizator": t("actor_local")})
	if configuratie.Revizie != "" {
		doc.Metadate += " · " + i18n.T(limba, "pontaj.revizie_reguli", map[string]string{"revizie": configuratie.Revizie})
	}
	if d.Continut.Sursa != "ABSENTE_ANUALE" {
		amprenta, err := amprentaLunii(ctx, tx, c.An, c.Luna, limba)
		if err != nil {
			return Document{}, traduce(err)
		}
		if amprenta != "" {
			doc.Metadate += " · " + amprenta
		}
	}
	configLegenda := configuratie
	configLegenda.Conditionate = nil
	for _, r := range configuratie.Conditionate {
		if !r.Activ {
			continue
		}
		for id := range inclusi {
			if profile[id].corespunde(Cerere{GrupID: r.GrupID, FunctieID: r.FunctieID, ContractID: r.ContractID, Pozitie: r.Pozitie, Garda: r.Garda}) {
				configLegenda.Conditionate = append(configLegenda.Conditionate, r)
				break
			}
		}
	}
	etichete := map[string]string{}
	if len(configLegenda.Conditionate) > 0 {
		rows, err := tx.QueryContext(ctx, `SELECT 'grup/'||id,denumire FROM grupuri_profesionale UNION ALL SELECT 'functie/'||id,denumire FROM functii UNION ALL SELECT 'contract/'||id,denumire FROM tipuri_contract`)
		if err != nil {
			return Document{}, traduce(err)
		}
		for rows.Next() {
			var k, n string
			if err := rows.Scan(&k, &n); err != nil {
				rows.Close()
				return Document{}, traduce(err)
			}
			etichete[k] = textJSON(n, limba)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return Document{}, traduce(err)
		}
	}
	if !d.Continut.AscundeExplicatii {
		doc.Legenda = legendaColoane(doc, configLegenda, limba, etichete)
	}
	doc.legendaPentru = func(cols []Coloana) []string {
		sub := doc
		sub.Coloane = cols
		return legendaColoane(sub, configLegenda, limba, etichete)
	}
	aplicaPrezentare(&doc, d.Continut.Prezentare)
	return doc, nil
}

// amprentaLunii întoarce textul amprentei pentru subsol, dacă luna a fost închisă (D044
// punctul 6). Dacă datele s-au schimbat după închidere, textul o spune.
func amprentaLunii(ctx context.Context, tx *sql.Tx, an, luna int, limba i18n.Limba) (string, error) {
	var pastrata string
	err := tx.QueryRowContext(ctx, `SELECT amprenta FROM stare_luni WHERE an = ? AND luna = ?`, an, luna).Scan(&pastrata)
	if err == sql.ErrNoRows || (err == nil && pastrata == "") {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	curenta, err := audit.AmprentaLuna(ctx, tx, an, luna)
	if err != nil {
		return "", err
	}
	if curenta == pastrata {
		return i18n.T(limba, "rapoarte.amprenta_luna", map[string]string{"amprenta": pastrata}), nil
	}
	return i18n.T(limba, "rapoarte.amprenta_diferita", map[string]string{"amprenta": pastrata, "curenta": curenta}), nil
}
