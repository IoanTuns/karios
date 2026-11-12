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

// Schimbarea orei (F04-T08, D102): cum se numara o tura din program in noptile in care
// ora se schimba. NOMINALA: durata din program (o tura 19-07 are mereu 12 h). EFECTIVA:
// de la ora de ceas a inceputului pana la ora de ceas a sfarsitului, in fusul Europe/Bucharest
// (13 h in octombrie, 11 h in martie). Gol inseamna NOMINALA (configuratiile salvate inainte).
const SchimbareOraNominala = "NOMINALA"
const SchimbareOraEfectiva = "EFECTIVA"

type RegulaColoana struct {
	Coloana       string   `json:"coloana"`
	Sursa         string   `json:"sursa"`
	Coduri        []string `json:"coduri"`
	FactorSutimi  int64    `json:"factor_sutimi"`
	NumaiAprobate bool     `json:"numai_aprobate"`
}
type ProfilPontaj struct {
	GrupID, FunctieID, ContractID string
	Conducere, Garda              bool
}
type RegulaConditionata struct {
	Regula     RegulaColoana `json:"regula"`
	GrupID     string        `json:"grup_id"`
	FunctieID  string        `json:"functie_id"`
	ContractID string        `json:"contract_id"`
	Pozitie    string        `json:"pozitie"`
	Garda      string        `json:"garda"`
	Activ      bool          `json:"activ"`
}

func (r RegulaConditionata) corespunde(p ProfilPontaj) bool {
	return (r.GrupID == "" || r.GrupID == p.GrupID) && (r.FunctieID == "" || r.FunctieID == p.FunctieID) && (r.ContractID == "" || r.ContractID == p.ContractID) && (r.Pozitie == "" || p.Conducere == (r.Pozitie == "CONDUCERE")) && (r.Garda == "" || p.Garda == (r.Garda == "DA"))
}

type ConfiguratiePontaj struct {
	Conditionate []RegulaConditionata    `json:"conditionate,omitempty"`
	Profile      map[string]ProfilPontaj `json:"-"`
	// Norme: norma zilnica in minute pe angajat, pentru TOCN (D092); nil = fara norma zilnica.
	Norme map[string]*int64 `json:"-"`
	// Garzi: randurile marcii de garda (D060), cu regulile din D093 pentru ZNL, Ore garda 200% si Ore lucratoare.
	Garzi        map[string]bool `json:"-"`
	Atribuire    string          `json:"atribuire"`
	SchimbareOra string          `json:"schimbare_ora"`
	Reguli       []RegulaColoana `json:"reguli"`
	Revizie      string          `json:"revizie"`
}

func ColoanePontaj() []string {
	return []string{"ore_lucrate", "T", "ZNL", "CO", "CEV", "CFS", "ore_norma", "form_prof", "CM", "L_FIN", "deleg", "ore_garda",
		"OL", "CFP", "CVE", "DLG", "FMP", "TZLN", "TOCN"}
}

// ColoaneRapoarte sunt coloanele de total pe care le arata rapoartele standard si
// conversia modelelor: lista de dinainte de D092, ca rapoartele de referinta sa nu se
// schimbe. Coloanele noi raman in calcul si in SQL (F08-T13).
func ColoaneRapoarte() []string { return ColoanePontaj()[:12] }

// ColoaneGrilaPontaj sunt coloanele de total afisate in grila Pontaj (D092); celelalte
// raman in calcul pentru rapoarte si export. TZLN si TOCN sunt ascunse.
func ColoaneGrilaPontaj() []string {
	return []string{"OL", "CO", "CM", "CFP", "CVE", "DLG", "FMP"}
}

// CoduriRedenumite: codurile de absenta redenumite prin D092 (vechi -> nou).
var CoduriRedenumite = map[string]string{"CEV": "CVE", "CFS": "CFP", "DEPL": "DLG", "FP": "FMP"}

// EchivalenteRapoarteD092 (F08-T15): coloanele de total vechi din rapoartele standard si
// coloanele D092 care le-au inlocuit. O coloana veche neconfigurata arata valoarea celei noi.
var EchivalenteRapoarteD092 = map[string]string{"CEV": "CVE", "CFS": "CFP", "form_prof": "FMP", "deleg": "DLG"}

// CoduriExcluseCVE sunt tipurile cu coloana proprie; CVE numara restul absentelor (D092).
var CoduriExcluseCVE = []string{"CO", "CM", "CFP", "DLG", "FMP"}

func UnitateColoana(c string) string {
	switch c {
	case "ore_lucrate", "T", "ZNL", "ore_norma", "ore_garda", "OL", "TOCN":
		return "ORE"
	case "CO", "CEV", "CFS", "CM", "deleg", "CFP", "CVE", "DLG", "FMP", "TZLN":
		return "ZILE"
	}
	return ""
}
func ConfiguratiePontajInitiala() ConfiguratiePontaj {
	c := ConfiguratiePontaj{Atribuire: AtribuireInceput, SchimbareOra: SchimbareOraNominala, Revizie: "1", Reguli: []RegulaColoana{}}
	for _, cod := range ColoanePontaj() {
		r := RegulaColoana{Coloana: cod, Sursa: "NEDEFINITA", Coduri: []string{}, FactorSutimi: 100}
		switch cod {
		case "ore_lucrate", "T":
			r.Sursa = "TOTAL_ORE"
		case "ZNL":
			r.Sursa = "ORE_ZNL"
		case "ore_norma":
			r.Sursa = "ORE_NORMA"
		case "CO", "CM", "CFP", "DLG", "FMP":
			r.Sursa = "ZILE_ABSENTA"
			r.Coduri = []string{cod}
		case "OL":
			// Orele salvate plus orele creditate ale formarii profesionale (D092).
			r.Sursa = "ORE_LUCRATE_CODURI"
			r.Coduri = []string{"FMP"}
		case "CVE":
			r.Sursa = "ZILE_ABSENTA_ALTELE"
			r.Coduri = append([]string(nil), CoduriExcluseCVE...)
		case "TZLN":
			r.Sursa = "ZILE_LUCRATOARE_NETE"
		case "TOCN":
			r.Sursa = "ORE_NORMA_LUNA"
		}
		c.Reguli = append(c.Reguli, r)
	}
	return c
}
func ValideazaConfiguratiePontaj(c ConfiguratiePontaj) error {
	gresit := func() error { return &Eroare{Cod: "pontaj.reguli_invalide"} }
	// gresitCamp indica rândul (coloana) și câmpul greșit (F05-T20, I049).
	gresitCamp := func(cod, coloana, camp string) error {
		return &Eroare{Cod: cod, Param: map[string]string{"coloana": coloana, "camp": camp}}
	}
	if c.Atribuire != AtribuireInceput && c.Atribuire != AtribuireCalendar {
		return gresit()
	}
	if c.SchimbareOra != "" && c.SchimbareOra != SchimbareOraNominala && c.SchimbareOra != SchimbareOraEfectiva {
		return gresit()
	}
	if len(c.Conditionate) > 100 {
		return gresit()
	}
	for _, r := range c.Conditionate {
		if (r.Pozitie != "" && r.Pozitie != "CONDUCERE" && r.Pozitie != "EXECUTIE") || (r.Garda != "" && r.Garda != "DA" && r.Garda != "NU") {
			return gresit()
		}
		baza := c
		baza.Conditionate = nil
		gasit := false
		for i, v := range baza.Reguli {
			if v.Coloana == r.Regula.Coloana {
				baza.Reguli = append([]RegulaColoana(nil), c.Reguli...)
				baza.Reguli[i] = r.Regula
				gasit = true
				break
			}
		}
		if !gasit {
			return gresit()
		}
		if err := ValideazaConfiguratiePontaj(baza); err != nil {
			return err
		}
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
		case "TOTAL_ORE", "ORE_ZNL", "ORE_NORMA", "ORE_SCHIMBURI", "ORE_SCHIMBURI_ZNL", "ORE_LUCRATE_CODURI", "ORE_NORMA_LUNA":
		case "ZILE_ABSENTA", "ZILE_ABSENTA_ALTELE", "ZILE_LUCRATOARE_NETE":
			unitate = "ZILE"
		default:
			return gresitCamp("pontaj.regula_sursa", r.Coloana, "sursa")
		}
		if u := UnitateColoana(r.Coloana); unitate != "" && u != "" && unitate != u {
			return gresitCamp("pontaj.regula_sursa", r.Coloana, "sursa")
		}
		if r.FactorSutimi < 1 || r.FactorSutimi > 1000 || unitate != "ORE" && r.FactorSutimi != 100 {
			return gresitCamp("pontaj.regula_factor", r.Coloana, "factor")
		}
		selectie := r.Sursa == "ZILE_ABSENTA" || (r.Sursa == "ORE_SCHIMBURI" || r.Sursa == "ORE_SCHIMBURI_ZNL")
		if selectie && len(r.Coduri) == 0 {
			return gresitCamp("pontaj.regula_fara_coduri", r.Coloana, "coduri")
		}
		// OL si CVE primesc o lista optionala: codurile creditate incluse, respectiv cele excluse.
		optionala := r.Sursa == "ORE_LUCRATE_CODURI" || r.Sursa == "ZILE_ABSENTA_ALTELE"
		if !selectie && !optionala && len(r.Coduri) > 0 || len(r.Coduri) > 200 {
			return gresitCamp("pontaj.regula_coduri", r.Coloana, "coduri")
		}
		if r.NumaiAprobate && r.Sursa != "ZILE_ABSENTA" && r.Sursa != "ZILE_ABSENTA_ALTELE" {
			return gresitCamp("pontaj.regula_coduri", r.Coloana, "aprobare")
		}
		coduri := map[string]bool{}
		for _, cod := range r.Coduri {
			if cod == "" || len(cod) > 16 || strings.TrimSpace(cod) != cod || coduri[cod] {
				return gresitCamp("pontaj.regula_coduri", r.Coloana, "coduri")
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
type IntrarePontaj struct {
	AngajatID, Data, Valoare, Absenta, StareAprobare string
	Crediteaza                                       bool
	MinuteCreditate                                  *int64
}
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
	// Motiv spune de ce valoarea nu e definită (F05-T20, I049), iar MotivCod
	// codul de program sau de absență în cauză; ambele lipsesc la valorile definite.
	Motiv    string `json:"motiv,omitempty"`
	MotivCod string `json:"motiv_cod,omitempty"`
}

// Motivele unei valori nedefinite, în ordinea în care se afișează când sunt mai multe.
// Motivul nu inventează reguli: arată doar ce lipsește din configurare sau din date.
const (
	MotivRegulaNedefinita   = "regula_nedefinita"
	MotivValoareNecunoscuta = "valoare_necunoscuta"
	MotivDurataInvalida     = "durata_invalida"
	MotivPauza              = "pauza"
	MotivAbsentaFaraOre     = "absenta_fara_ore"
	MotivFaraOraInceput     = "fara_ora_inceput"
	// MotivGardaFaraZNL: pe randul marcii de garda ZNL nu se calculeaza (D093).
	MotivGardaFaraZNL = "garda_fara_znl"
	// MotivFaraNorma: TOCN cere norma zilnica, pe care garzile si unele contracte nu o au (D092).
	MotivFaraNorma = "fara_norma"
)

var ordineMotive = []string{MotivValoareNecunoscuta, MotivDurataInvalida, MotivPauza, MotivAbsentaFaraOre, MotivFaraOraInceput}

// noteazaMotiv retine, pentru fiecare motiv, cel mai mic cod (rezultat stabil).
func noteazaMotiv(motive map[string]string, motiv, cod string) {
	if vechi, ok := motive[motiv]; !ok || cod < vechi {
		motive[motiv] = cod
	}
}

func primulMotiv(motive map[string]string) (string, string) {
	for _, m := range ordineMotive {
		if cod, ok := motive[m]; ok {
			return m, cod
		}
	}
	return "", ""
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
	schimburiZNL      map[string]int64
	absente           []IntrarePontaj
	// motiveComplet: de ce totalul nu e complet; motiveCalendar: de ce orele nu
	// se pot imparti pe tipuri de zi (F05-T20).
	motiveComplet, motiveCalendar map[string]string
	// credite: minutele creditate de absente, pe cod (pentru OL, D092).
	credite map[string]int64
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
			oameni[id] = &acumularePontaj{complet: true, calendar: true, schimburi: map[string]int64{}, schimburiZNL: map[string]int64{}, motiveComplet: map[string]string{}, motiveCalendar: map[string]string{}, credite: map[string]int64{}}
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
	// Fiecare interval contribuie separat; absenta se aplica o singura data pe zi.
	expandate := make([]IntrarePontaj, 0, len(intrari))
	for _, in := range intrari {
		if in.Absenta != "" {
			expandate = append(expandate, in)
			continue
		}
		coduri := CoduriPrograme(in.Valoare)
		if len(coduri) > 1 {
			intervale := []IntervalPlanificare{}
			d, e := time.Parse(time.DateOnly, in.Data)
			if e != nil {
				return rez, &Eroare{Cod: "pontaj.celula_invalida"}
			}
			for _, cod := range coduri {
				p, ok := mapPrograme[cod]
				if !ok {
					return rez, &Eroare{Cod: "pontaj.celula_invalida"}
				}
				program := ProgramPrecompletare{OraStart: p.OraStart, DurataMinute: p.DurataMinute, PauzaMinute: p.PauzaMinute}
				if motiv := VerificaTuraPlanificata(in.AngajatID, in.Data, program, nil, intervale, nil); motiv != "" {
					return rez, &Eroare{Cod: "planificare.motiv_" + motiv}
				}
				a := d.Add(time.Duration(p.OraStart) * time.Minute)
				intervale = append(intervale, IntervalPlanificare{AngajatID: in.AngajatID, Inceput: a, Sfarsit: a.Add(time.Duration(p.DurataMinute) * time.Minute)})
			}
		}
		for _, cod := range coduri {
			nou := in
			nou.Valoare = cod
			expandate = append(expandate, nou)
		}
	}
	for _, in := range expandate {
		data, e := time.Parse(time.DateOnly, in.Data)
		if e != nil {
			return rez, &Eroare{Cod: "pontaj.celula_invalida"}
		}
		inPerioada := in.Data >= deLa && in.Data <= panaLa
		if in.Absenta != "" {
			if inPerioada {
				a := persoana(in.AngajatID)
				a.absente = append(a.absente, in)
				z := celula(in.AngajatID, in.Data)
				if in.Crediteaza && in.StareAprobare == "APROBATA" {
					if in.MinuteCreditate == nil {
						z.Calculabil = false
						a.complet = false
						a.calendar = false
						noteazaMotiv(a.motiveComplet, MotivAbsentaFaraOre, in.Absenta)
					} else {
						z.Minute += *in.MinuteCreditate
						a.total += *in.MinuteCreditate
						a.norma += *in.MinuteCreditate
						a.credite[in.Absenta] += *in.MinuteCreditate
					}
				}
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
		motivNecunoscut, codValoare := "", strings.ToUpper(strings.TrimSpace(in.Valoare))
		switch {
		case !program && !manual:
			motivNecunoscut = MotivValoareNecunoscuta
		case program && (p.DurataMinute <= 0 || p.DurataMinute > 2880):
			motivNecunoscut = MotivDurataInvalida
		case program && p.PauzaMinute != 0:
			motivNecunoscut = MotivPauza
		case faraStart:
			motivNecunoscut = MotivFaraOraInceput
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
				noteazaMotiv(a.motiveComplet, motivNecunoscut, codValoare)
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
				noteazaMotiv(a.motiveCalendar, MotivFaraOraInceput, codValoare)
				return
			}
			d := zi.Format(time.DateOnly)
			if zi.Weekday() == time.Saturday || zi.Weekday() == time.Sunday || libere[d] {
				a.znl += n
				if program {
					a.schimburiZNL[p.Cod] += n
				}
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
		if c.SchimbareOra == SchimbareOraEfectiva {
			inceput = MomentLocal(data, p.OraStart)
			sfarsit = MomentLocal(data, p.OraStart+p.DurataMinute)
		}
		for curent := inceput; curent.Before(sfarsit); {
			// Miezul noptii in acelasi fus ca momentele turei (UTC nominal sau Europe/Bucharest).
			limita := time.Date(curent.Year(), curent.Month(), curent.Day()+1, 0, 0, 0, 0, curent.Location())
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
	// TZLN: zilele lunii fara weekend si fara sarbatorile care cad luni-vineri (D092).
	tzln := zileLucratoareNete(deLa, panaLa, libere)
	// Fiecare angajat cu norma cunoscuta primeste totaluri, ca TOCN sa existe si fara pontaj.
	for id := range c.Norme {
		persoana(id)
	}
	for id, a := range oameni {
		t := TotalPontaj{AngajatID: id, Valori: map[string]ValoarePontaj{}}
		reguli := append([]RegulaColoana(nil), c.Reguli...)
		aplicate := map[string]bool{}
		for _, rc := range c.Conditionate {
			if !rc.Activ || !rc.corespunde(c.Profile[id]) {
				continue
			}
			if aplicate[rc.Regula.Coloana] {
				return rez, &Eroare{Cod: "pontaj.reguli_suprapuse"}
			}
			aplicate[rc.Regula.Coloana] = true
			for i, r := range reguli {
				if r.Coloana == rc.Regula.Coloana {
					reguli[i] = rc.Regula
				}
			}
		}
		for _, r := range reguli {
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
			case "ORE_SCHIMBURI_ZNL":
				v.Unitate = "ORE"
				v.Definita = a.complet && a.calendar
				for _, cod := range r.Coduri {
					v.Sutimi += a.schimburiZNL[cod] * r.FactorSutimi
				}
			case "ORE_LUCRATE_CODURI":
				// Orele salvate (fara creditele absentelor) plus creditele codurilor alese.
				v.Unitate = "ORE"
				v.Definita = a.complet
				minute := a.total
				for _, n := range a.credite {
					minute -= n
				}
				for _, cod := range r.Coduri {
					minute += a.credite[cod]
				}
				v.Sutimi = minute * r.FactorSutimi
			case "ZILE_ABSENTA_ALTELE":
				v.Unitate = "ZILE"
				for _, zi := range a.absente {
					if r.NumaiAprobate && zi.StareAprobare != "APROBATA" {
						continue
					}
					exclus := false
					for _, cod := range r.Coduri {
						exclus = exclus || zi.Absenta == cod
					}
					if !exclus {
						v.Sutimi += 100
					}
				}
			case "ZILE_LUCRATOARE_NETE":
				v.Unitate = "ZILE"
				v.Sutimi = tzln * 100
			case "ORE_NORMA_LUNA":
				v.Unitate = "ORE"
				if n := c.Norme[id]; n != nil {
					v.Sutimi = tzln * *n * r.FactorSutimi
				} else {
					v.Definita = false
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
			if !v.Definita {
				switch {
				case r.Sursa == "NEDEFINITA":
					v.Motiv = MotivRegulaNedefinita
				case r.Sursa == "ORE_NORMA_LUNA":
					v.Motiv = MotivFaraNorma
				case !a.complet:
					v.Motiv, v.MotivCod = primulMotiv(a.motiveComplet)
				default:
					v.Motiv, v.MotivCod = primulMotiv(a.motiveCalendar)
				}
			}
			t.Valori[r.Coloana] = v
		}
		if c.Garzi[id] {
			aplicaReguliGarda(t.Valori, a)
		}
		rez.Totaluri = append(rez.Totaluri, t)
	}
	sort.Slice(rez.Totaluri, func(i, j int) bool { return rez.Totaluri[i].AngajatID < rez.Totaluri[j].AngajatID })
	return rez, nil
}

// zileLucratoareNete numara zilele din [deLa, panaLa] care nu sunt sambata, duminica
// sau sarbatoare legala (D092, TZLN). O sarbatoare in weekend nu se scade de doua ori.
func zileLucratoareNete(deLa, panaLa string, libere map[string]bool) int64 {
	inceput, e1 := time.Parse(time.DateOnly, deLa)
	sfarsit, e2 := time.Parse(time.DateOnly, panaLa)
	if e1 != nil || e2 != nil {
		return 0
	}
	var n int64
	for zi := inceput; !zi.After(sfarsit); zi = zi.AddDate(0, 0, 1) {
		if zi.Weekday() != time.Saturday && zi.Weekday() != time.Sunday && !libere[zi.Format(time.DateOnly)] {
			n++
		}
	}
	return n
}

// aplicaReguliGarda (D093): pe randul marcii de garda ZNL nu se calculeaza, Ore garda 200%
// sunt orele din weekend si sarbatori (aceeasi impartire pe zile ca ZNL), iar Ore lucratoare
// = total ore - ore garda 200%. Regulile au intaietate fata de Configurare pontaj pe aceste randuri.
func aplicaReguliGarda(v map[string]ValoarePontaj, a *acumularePontaj) {
	definit := a.complet && a.calendar
	v["ZNL"] = ValoarePontaj{Unitate: "ORE", Motiv: MotivGardaFaraZNL}
	g := ValoarePontaj{Sutimi: a.znl * 100, Unitate: "ORE", Definita: definit}
	l := ValoarePontaj{Sutimi: (a.total - a.znl) * 100, Unitate: "ORE", Definita: definit}
	if !definit {
		g.Motiv, g.MotivCod = primulMotiv(a.motiveComplet)
		if g.Motiv == "" {
			g.Motiv, g.MotivCod = primulMotiv(a.motiveCalendar)
		}
		l.Motiv, l.MotivCod = g.Motiv, g.MotivCod
	}
	v["ore_garda"] = g
	v["ore_lucratoare"] = l
}
