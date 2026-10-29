package ai

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SelectiePlan este un selector, nu lista persoanelor aleasă de model.
// Grupurile sunt rezolvate din baza curentă la previzualizare și confirmare.
type SelectiePlan struct {
	Tip        string   `json:"tip"`
	Valoare    string   `json:"valoare"`
	Marci      []string `json:"marci"`
	FaraPontaj bool     `json:"fara_pontaj"`
}

type IntentiePlan struct {
	Selectie      SelectiePlan `json:"selectie"`
	DeLa          string       `json:"de_la"`
	PanaLa        string       `json:"pana_la"`
	CodSchimb     string       `json:"cod_schimb"`
	FiltruZile    string       `json:"filtru_zile"`
	ModExistente  string       `json:"mod_existente"`
	MesajOriginal string       `json:"mesaj_original"`
	Clarificari   []string     `json:"clarificari"`
}

var reDataISO = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`)
var reAnPlan = regexp.MustCompile(`\b20\d{2}\b`)
var reZiPlan = regexp.MustCompile(`\b(?:pe|din|on) (\d{1,2}) (ianuarie|februarie|martie|aprilie|mai|iunie|iulie|august|septembrie|octombrie|noiembrie|decembrie)\b`)
var luniEngleza = []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"}

func contineExpresie(m, expresie string) bool {
	return strings.Contains(" "+m+" ", " "+normalizeaza(expresie)+" ")
}

// InterpreteazaPlan recunoaște intenția fără model. Câmpurile nesigure rămân
// necompletate și se cer în formular; nu se presupune L–V sau suprascrierea.
func InterpreteazaPlan(mesaj string, info ContextPontix) (IntentiePlan, bool) {
	m, _ := corecteazaMesaj(mesaj, info)
	m = strings.TrimSpace(regexp.MustCompile(`[,;!?]`).ReplaceAllString(m, " "))
	i := IntentiePlan{MesajOriginal: mesaj, Selectie: SelectiePlan{Marci: []string{}}}
	if regexp.MustCompile(`\b(?:sterge|sterg|stergeti|delete|remove)\b`).MatchString(m) {
		return i, false
	}
	if esteListareAngajati(m) || esteListareSchimburi(m) || regexp.MustCompile(`^(?:cine|cati|cate|care|ce|who|how many)\b`).MatchString(m) {
		return i, false
	}
	schimb, areSchimb := gasesteSchimb(mesaj, m, info.Schimburi)
	actiune := areSchimb || (esteToti(m) && reInterval.MatchString(m))
	for _, cuvant := range []string{"planifica", "pune", "lucreaza", "completeaza", "schedule", "assign", "work"} {
		actiune = actiune || strings.Contains(m, cuvant)
	}
	if !actiune {
		return i, false
	}
	if areSchimb {
		i.CodSchimb = schimb.Cod
	}
	// Mai multe programe sau grupuri diferite nu se reduc la primul găsit.
	coduri := 0
	for _, s := range info.Schimburi {
		if contineExpresie(m, s.Cod) {
			coduri++
		}
	}
	if coduri > 1 {
		i.CodSchimb = ""
	}
	// Un grup numit are prioritate față de „toți” din „toți din secția ATI”.
	for _, camp := range []struct {
		tip, prefix string
		nume        func(InfoAngajat) string
	}{
		{"sectie", "sectia", func(a InfoAngajat) string { return a.Sectie }},
		{"departament", "departamentul", func(a InfoAngajat) string { return a.Departament }},
		{"categorie", "categoria", func(a InfoAngajat) string { return a.Categorie }},
	} {
		if !strings.Contains(m, camp.prefix+" ") {
			continue
		}
		for _, a := range info.Angajati {
			nume := camp.nume(a)
			if nume != "" && contineExpresie(m, camp.prefix+" "+nume) && len(nume) > len(i.Selectie.Valoare) {
				i.Selectie.Tip, i.Selectie.Valoare = camp.tip, nume
			}
		}
		if i.Selectie.Tip == "" {
			i.Selectie.Tip = camp.tip
		}
	}
	if i.Selectie.Tip == "" && (esteToti(m) || strings.Contains(m, "all employees") || strings.Contains(m, "all staff")) {
		i.Selectie.Tip = "toti"
	}
	grupuriNumite := 0
	for _, prefix := range []string{"sectia ", "departamentul ", "categoria "} {
		if strings.Contains(m, prefix) {
			grupuriNumite++
		}
	}
	if i.Selectie.Tip == "" {
		for _, a := range info.Angajati {
			if contineExpresie(m, a.Marca) || contineExpresie(m, a.Nume+" "+a.Prenume) || contineExpresie(m, a.Prenume+" "+a.Nume) {
				i.Selectie.Marci = append(i.Selectie.Marci, a.Marca)
			}
		}
		if len(i.Selectie.Marci) > 0 {
			i.Selectie.Tip = "persoane"
		}
	}
	if grupuriNumite > 1 || regexp.MustCompile(`\b(?:exceptand|exceptia|except|fara angajat\w*|nu toti)\b`).MatchString(m) {
		i.Selectie = SelectiePlan{Marci: []string{}}
	}
	i.Selectie.FaraPontaj = esteFaraPontaj(m) || strings.Contains(m, "without timesheets")
	an := info.An
	if x := reAnPlan.FindString(m); x != "" {
		an, _ = strconv.Atoi(x)
	}
	date := reDataISO.FindAllString(m, -1)
	if len(date) > 0 {
		i.DeLa, i.PanaLa = date[0], date[len(date)-1]
	} else {
		type lunaGasita struct{ poz, luna int }
		var luni []lunaGasita
		for l, nume := range luniNorm {
			for _, n := range []string{nume, luniEngleza[l]} {
				if n == "mai" && !reLunaMaiNeambiguu.MatchString(m) && !strings.Contains(m, "si mai") {
					continue
				}
				if p := regexp.MustCompile(`\b` + n + `\b`).FindStringIndex(m); p != nil {
					luni = append(luni, lunaGasita{p[0], l + 1})
				}
			}
		}
		for x := range luni {
			for y := x + 1; y < len(luni); y++ {
				if luni[y].poz < luni[x].poz {
					luni[x], luni[y] = luni[y], luni[x]
				}
			}
		}
		prima, ultima := info.Luna, info.Luna
		if len(luni) > 0 {
			prima, ultima = luni[0].luna, luni[len(luni)-1].luna
		}
		anFinal := an
		if ultima < prima {
			anFinal++
		}
		if an >= 2000 && prima >= 1 && prima <= 12 {
			i.DeLa = time.Date(an, time.Month(prima), 1, 0, 0, 0, 0, time.UTC).Format(time.DateOnly)
			i.PanaLa = time.Date(anFinal, time.Month(ultima)+1, 0, 0, 0, 0, 0, time.UTC).Format(time.DateOnly)
		}
		if len(luni) > 1 && !strings.Contains(m, "pana") && !strings.Contains(m, "through") && !strings.Contains(m, " to ") {
			for x := 1; x < len(luni); x++ {
				if luni[x].luna != luni[x-1].luna%12+1 {
					i.DeLa, i.PanaLa = "", ""
				}
			}
		}
		if zi := reZiPlan.FindStringSubmatch(m); zi != nil {
			n, _ := strconv.Atoi(zi[1])
			d := time.Date(an, time.Month(prima), n, 0, 0, 0, 0, time.UTC)
			if d.Day() == n {
				i.DeLa, i.PanaLa = d.Format(time.DateOnly), d.Format(time.DateOnly)
			} else {
				i.DeLa, i.PanaLa = "", ""
			}
		}
		// Nu transformăm intervale zilnice sau date relative în luni întregi.
		if regexp.MustCompile(`\b(?:maine|azi|saptamana|tomorrow|today|next week)\b|\b(?:intre|din|de pe) \d`).MatchString(m) && reZiPlan.FindString(m) == "" {
			i.DeLa, i.PanaLa = "", ""
		}
	}
	switch {
	case strings.Contains(m, "fara weekend") || strings.Contains(m, "nu in weekend"):
		i.FiltruZile = "lucratoare"
	case strings.Contains(m, "toate zilele") || strings.Contains(m, "inclusiv weekend") || strings.Contains(m, "every day") || strings.Contains(m, "all days"):
		i.FiltruZile = "toate"
	case strings.Contains(m, "luni pana vineri") || strings.Contains(m, "luni-vineri") || strings.Contains(m, "luni–vineri") || strings.Contains(m, "l-v") || strings.Contains(m, "lucratoare") || strings.Contains(m, "weekdays"):
		i.FiltruZile = "lucratoare"
	case strings.Contains(m, "weekend"):
		i.FiltruZile = "weekend"
	}
	switch {
	case strings.Contains(m, "doar celulele goale") || strings.Contains(m, "numai celulele goale") || strings.Contains(m, "pastreaza") || strings.Contains(m, "empty cells") || strings.Contains(m, "keep existing"):
		i.ModExistente = "pastreaza"
	case strings.Contains(m, "inlocuieste") || strings.Contains(m, "suprascrie") || strings.Contains(m, "replace existing"):
		i.ModExistente = "inlocuieste"
	}
	i.CompleteazaClarificari()
	return i, true
}

func (i *IntentiePlan) CompleteazaClarificari() {
	i.Clarificari = []string{}
	if i.Selectie.Tip == "" || (i.Selectie.Tip != "toti" && i.Selectie.Tip != "persoane" && i.Selectie.Valoare == "") || (i.Selectie.Tip == "persoane" && len(i.Selectie.Marci) == 0) {
		i.Clarificari = append(i.Clarificari, "selectie")
	}
	if i.DeLa == "" || i.PanaLa == "" {
		i.Clarificari = append(i.Clarificari, "perioada")
	}
	if i.CodSchimb == "" {
		i.Clarificari = append(i.Clarificari, "program")
	}
	if i.FiltruZile == "" {
		i.Clarificari = append(i.Clarificari, "zile")
	}
	if i.ModExistente == "" {
		i.Clarificari = append(i.Clarificari, "existente")
	}
}
