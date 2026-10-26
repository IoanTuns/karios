package calcul

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

const AtribuireInceput = "ZI_INCEPUT"
const AtribuireCalendar = "ZILE_CALENDARISTICE"

type RegulaColoana struct {
	Coloana       string   `json:"coloana"`
	Sursa         string   `json:"sursa"`
	Coduri        []string `json:"coduri"`
	FactorSutimi  int64    `json:"factor_sutimi"`
	NumaiAprobate bool     `json:"numai_aprobate"`
}
type ConfiguratiePontaj struct {
	Atribuire string          `json:"atribuire"`
	Reguli    []RegulaColoana `json:"reguli"`
	Revizie   string          `json:"revizie"`
}

func ColoanePontaj() []string {
	return []string{"ore_lucrate", "T", "ZNL", "CO", "CEV", "CFS", "ore_norma", "form_prof", "CM", "L_FIN", "deleg", "ore_garda"}
}
func UnitateColoana(c string) string {
	switch c {
	case "ore_lucrate", "T", "ZNL", "ore_norma", "ore_garda":
		return "ORE"
	case "CO", "CEV", "CFS", "CM", "deleg":
		return "ZILE"
	}
	return ""
}
func ConfiguratiePontajInitiala() ConfiguratiePontaj {
	c := ConfiguratiePontaj{Atribuire: AtribuireInceput, Revizie: "1", Reguli: []RegulaColoana{}}
	for _, cod := range ColoanePontaj() {
		r := RegulaColoana{Coloana: cod, Sursa: "NEDEFINITA", Coduri: []string{}, FactorSutimi: 100}
		switch cod {
		case "ore_lucrate", "T":
			r.Sursa = "TOTAL_ORE"
		case "ZNL":
			r.Sursa = "ORE_ZNL"
		case "ore_norma":
			r.Sursa = "ORE_NORMA"
		case "CO", "CM":
			r.Sursa = "ZILE_ABSENTA"
			r.Coduri = []string{cod}
		}
		c.Reguli = append(c.Reguli, r)
	}
	return c
}
func ValideazaConfiguratiePontaj(c ConfiguratiePontaj) error {
	gresit := func() error { return &Eroare{Cod: "pontaj.reguli_invalide"} }
	if c.Atribuire != AtribuireInceput && c.Atribuire != AtribuireCalendar {
		return gresit()
	}
	coloane := ColoanePontaj()
	if len(c.Reguli) != len(coloane) {
		return gresit()
	}
	vazute := map[string]bool{}
	for _, r := range c.Reguli {
		gasit := false
		for _, col := range coloane {
			gasit = gasit || col == r.Coloana
		}
		if !gasit || vazute[r.Coloana] {
			return gresit()
		}
		vazute[r.Coloana] = true
		unitate := "ORE"
		switch r.Sursa {
		case "NEDEFINITA":
			unitate = ""
		case "TOTAL_ORE", "ORE_ZNL", "ORE_NORMA", "ORE_SCHIMBURI":
		case "ZILE_ABSENTA":
			unitate = "ZILE"
		default:
			return gresit()
		}
		if u := UnitateColoana(r.Coloana); unitate != "" && u != "" && unitate != u {
			return gresit()
		}
		if r.FactorSutimi < 1 || r.FactorSutimi > 1000 || unitate != "ORE" && r.FactorSutimi != 100 {
			return gresit()
		}
		selectie := r.Sursa == "ZILE_ABSENTA" || r.Sursa == "ORE_SCHIMBURI"
		if selectie != (len(r.Coduri) > 0) || len(r.Coduri) > 200 || r.NumaiAprobate && r.Sursa != "ZILE_ABSENTA" {
			return gresit()
		}
		coduri := map[string]bool{}
		for _, cod := range r.Coduri {
			if cod == "" || len(cod) > 16 || strings.TrimSpace(cod) != cod || coduri[cod] {
				return gresit()
			}
			coduri[cod] = true
		}
	}
	return nil
}

type ProgramPontaj struct {
	Cod                                 string
	OraStart, DurataMinute, PauzaMinute int64
}
type IntrarePontaj struct{ AngajatID, Data, Valoare, Absenta, StareAprobare string }
type OreZi struct {
	AngajatID  string   `json:"angajat_id"`
	Data       string   `json:"data"`
	Minute     int64    `json:"minute"`
	Calculabil bool     `json:"calculabil"`
	Surse      []string `json:"surse"`
}

// ORE se păstrează în sutimi de minut, ZILE în sutimi de zi; evităm rotunjirea calculelor.
type ValoarePontaj struct {
	Sutimi   int64  `json:"sutimi"`
	Unitate  string `json:"unitate"`
	Definita bool   `json:"definita"`
}
type TotalPontaj struct {
	AngajatID string                   `json:"angajat_id"`
	Valori    map[string]ValoarePontaj `json:"valori"`
}
type RezultatPontaj struct {
	Celule   []OreZi       `json:"celule"`
	Totaluri []TotalPontaj `json:"totaluri"`
}

func OreManuale(text string) (int64, bool) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "@M") {
		n, e := strconv.ParseInt(strings.TrimPrefix(text, "@M"), 10, 64)
		return n, e == nil && n >= 0 && n <= 2880
	}
	if text == "" {
		return 0, false
	}
	for _, r := range text {
		if !(r >= '0' && r <= '9' || r == '.' || r == ',') {
			return 0, false
		}
	}
	ore, e := strconv.ParseFloat(strings.ReplaceAll(text, ",", "."), 64)
	minute := math.Round(ore * 60)
	return int64(minute), e == nil && ore >= 0 && ore <= 48 && math.Abs(ore*60-minute) < 0.00001
}
func TextValoarePontaj(v ValoarePontaj) string {
	if !v.Definita {
		return "—"
	}
	divizor := 100.0
	if v.Unitate == "ORE" {
		divizor = 6000
	}
	text := strconv.FormatFloat(float64(v.Sutimi)/divizor, 'f', 4, 64)
	return strings.TrimRight(strings.TrimRight(text, "0"), ".")
}

type acumularePontaj struct {
	total, norma, znl int64
	complet, calendar bool
	schimburi         map[string]int64
	absente           []IntrarePontaj
}

// CalculeazaPontaj folosește durata nominală a programelor, fără reguli salariale.
// Calendarul clasifică întotdeauna fragmentele 00–24; atribuirea decide perioada creditată.
func CalculeazaPontaj(deLa, panaLa string, c ConfiguratiePontaj, programe []ProgramPontaj, intrari []IntrarePontaj, libere map[string]bool) (RezultatPontaj, error) {
	rez := RezultatPontaj{Celule: []OreZi{}, Totaluri: []TotalPontaj{}}
	if err := ValideazaConfiguratiePontaj(c); err != nil {
		return rez, err
	}
	if _, e := time.Parse(time.DateOnly, deLa); e != nil {
		return rez, &Eroare{Cod: "pontaj.perioada_invalida"}
	}
	if _, e := time.Parse(time.DateOnly, panaLa); e != nil || panaLa < deLa {
		return rez, &Eroare{Cod: "pontaj.perioada_invalida"}
	}
	mapPrograme := map[string]ProgramPontaj{}
	for _, p := range programe {
		mapPrograme[strings.ToUpper(p.Cod)] = p
	}
	celule := map[string]*OreZi{}
	oameni := map[string]*acumularePontaj{}
	persoana := func(id string) *acumularePontaj {
		if oameni[id] == nil {
			oameni[id] = &acumularePontaj{complet: true, calendar: true, schimburi: map[string]int64{}}
		}
		return oameni[id]
	}
	celula := func(id, data string) *OreZi {
		k := id + "/" + data
		if celule[k] == nil {
			celule[k] = &OreZi{AngajatID: id, Data: data, Calculabil: true, Surse: []string{}}
		}
		return celule[k]
	}
	for _, in := range intrari {
		data, e := time.Parse(time.DateOnly, in.Data)
		if e != nil {
			return rez, &Eroare{Cod: "pontaj.celula_invalida"}
		}
		inPerioada := in.Data >= deLa && in.Data <= panaLa
		if in.Absenta != "" {
			if inPerioada {
				a := persoana(in.AngajatID)
				a.absente = append(a.absente, in)
				celula(in.AngajatID, in.Data)
			}
			continue
		}
		if strings.TrimSpace(in.Valoare) == "" {
			continue
		}
		if c.Atribuire == AtribuireInceput && !inPerioada {
			continue
		}
		p, program := mapPrograme[strings.ToUpper(strings.TrimSpace(in.Valoare))]
		minute, manual := OreManuale(in.Valoare)
		cunoscut := manual
		if program {
			minute = p.DurataMinute
			cunoscut = minute > 0 && minute <= 2880 && p.PauzaMinute == 0
		}
		faraStart := program && (p.OraStart < 0 || p.OraStart >= 1440)
		// Lipsa orei de început nu afectează totalul nominal, dar împiedică clasificarea calendaristică.
		if faraStart && c.Atribuire == AtribuireCalendar {
			cunoscut = false
		}
		tipFragment := func(zi string, n int64) {
			if zi < deLa || zi > panaLa {
				return
			}
			a := persoana(in.AngajatID)
			z := celula(in.AngajatID, zi)
			prezent := false
			for _, d := range z.Surse {
				prezent = prezent || d == in.Data
			}
			if !prezent {
				z.Surse = append(z.Surse, in.Data)
			}
			if !cunoscut {
				z.Calculabil = false
				a.complet = false
				a.calendar = false
				return
			}
			z.Minute += n
			a.total += n
			if program {
				a.schimburi[p.Cod] += n
			}
		}
		clasifica := func(zi time.Time, n int64) {
			a := persoana(in.AngajatID)
			if faraStart {
				a.calendar = false
				return
			}
			d := zi.Format(time.DateOnly)
			if zi.Weekday() == time.Saturday || zi.Weekday() == time.Sunday || libere[d] {
				a.znl += n
			} else {
				a.norma += n
			}
		}
		if !program || faraStart {
			if inPerioada {
				tipFragment(in.Data, minute)
				if cunoscut {
					clasifica(data, minute)
				}
			}
			continue
		}
		inceput := data.Add(time.Duration(p.OraStart) * time.Minute)
		sfarsit := inceput.Add(time.Duration(p.DurataMinute) * time.Minute)
		for curent := inceput; curent.Before(sfarsit); {
			limita := time.Date(curent.Year(), curent.Month(), curent.Day()+1, 0, 0, 0, 0, time.UTC)
			if limita.After(sfarsit) {
				limita = sfarsit
			}
			n := int64(limita.Sub(curent) / time.Minute)
			zi := curent.Format(time.DateOnly)
			credit := zi
			if c.Atribuire == AtribuireInceput {
				credit = in.Data
			}
			if credit >= deLa && credit <= panaLa {
				tipFragment(credit, n)
				if cunoscut {
					clasifica(curent, n)
				}
			}
			curent = limita
		}
	}
	for _, z := range celule {
		sort.Strings(z.Surse)
		rez.Celule = append(rez.Celule, *z)
	}
	sort.Slice(rez.Celule, func(i, j int) bool {
		a, b := rez.Celule[i], rez.Celule[j]
		return a.AngajatID < b.AngajatID || a.AngajatID == b.AngajatID && a.Data < b.Data
	})
	for id, a := range oameni {
		t := TotalPontaj{AngajatID: id, Valori: map[string]ValoarePontaj{}}
		for _, r := range c.Reguli {
			v := ValoarePontaj{Unitate: UnitateColoana(r.Coloana), Definita: r.Sursa != "NEDEFINITA"}
			switch r.Sursa {
			case "TOTAL_ORE":
				v.Unitate = "ORE"
				v.Sutimi = a.total * r.FactorSutimi
				v.Definita = a.complet
			case "ORE_ZNL":
				v.Unitate = "ORE"
				v.Sutimi = a.znl * r.FactorSutimi
				v.Definita = a.complet && a.calendar
			case "ORE_NORMA":
				v.Unitate = "ORE"
				v.Sutimi = a.norma * r.FactorSutimi
				v.Definita = a.complet && a.calendar
			case "ORE_SCHIMBURI":
				v.Unitate = "ORE"
				v.Definita = a.complet
				for _, cod := range r.Coduri {
					v.Sutimi += a.schimburi[cod] * r.FactorSutimi
				}
			case "ZILE_ABSENTA":
				v.Unitate = "ZILE"
				for _, zi := range a.absente {
					if r.NumaiAprobate && zi.StareAprobare != "APROBATA" {
						continue
					}
					for _, cod := range r.Coduri {
						if zi.Absenta == cod {
							v.Sutimi += 100
							break
						}
					}
				}
			}
			t.Valori[r.Coloana] = v
		}
		rez.Totaluri = append(rez.Totaluri, t)
	}
	sort.Slice(rez.Totaluri, func(i, j int) bool { return rez.Totaluri[i].AngajatID < rez.Totaluri[j].AngajatID })
	return rez, nil
}
