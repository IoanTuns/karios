package ai

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TipPlanificaGrup e propunerea care aplică același schimb mai multor angajați (F14-T08).
const TipPlanificaGrup = "planifica_grup"

var faraDiacritice = strings.NewReplacer("ă", "a", "â", "a", "î", "i", "ș", "s", "ş", "s", "ț", "t", "ţ", "t")

// normalizeaza scoate diacriticele românești, trece la litere mici și strânge spațiile.
func normalizeaza(s string) string {
	return strings.Join(strings.Fields(faraDiacritice.Replace(strings.ToLower(s))), " ")
}

var expresiiToti = []string{"toti angajatii", "toti salariatii", "tot personalul", "toata lumea", "toate persoanele", "fiecare angajat", "toti colegii", "toti oamenii"}

var (
	reInterval = regexp.MustCompile(`\b(\d{1,2})(?:[:.](\d{2}))?\s*(?:-|–|—|la|pana la)\s*(\d{1,2})(?:[:.](\d{2}))?\b`)
	reZiAnume  = regexp.MustCompile(`\bpe (\d{1,2})\b(?:\s*(?:ianuarie|februarie|martie|aprilie|mai|iunie|iulie|august|septembrie|octombrie|noiembrie|decembrie|a lunii))`)
)

var luniNorm = []string{"ianuarie", "februarie", "martie", "aprilie", "mai", "iunie", "iulie", "august", "septembrie", "octombrie", "noiembrie", "decembrie"}

// grupComanda descrie cui se adresează o comandă de grup.
type grupComanda struct {
	eticheta string
	angajati []InfoAngajat
}

func gasesteGrup(m string, info ContextPontix) (grupComanda, bool) {
	for _, e := range expresiiToti {
		if strings.Contains(m, e) {
			return grupComanda{eticheta: "toți angajații activi", angajati: info.Angajati}, true
		}
	}
	potrivire := func(cuvant string, camp func(InfoAngajat) string) (grupComanda, bool) {
		if !strings.Contains(m, cuvant+" ") {
			return grupComanda{}, false
		}
		var gasit string
		for _, a := range info.Angajati {
			n := normalizeaza(camp(a))
			if n != "" && strings.Contains(m, cuvant+" "+n) && len(n) > len(gasit) {
				gasit = n
			}
		}
		if gasit == "" {
			return grupComanda{}, false
		}
		g := grupComanda{}
		for _, a := range info.Angajati {
			if normalizeaza(camp(a)) == gasit {
				g.angajati = append(g.angajati, a)
				if g.eticheta == "" {
					g.eticheta = cuvant + " " + camp(a)
				}
			}
		}
		return g, true
	}
	if g, ok := potrivire("sectia", func(a InfoAngajat) string { return a.Sectie }); ok {
		g.eticheta = "secția " + g.angajati[0].Sectie
		return g, true
	}
	if g, ok := potrivire("categoria", func(a InfoAngajat) string { return a.Categorie }); ok {
		g.eticheta = "categoria " + g.angajati[0].Categorie
		return g, true
	}
	return grupComanda{}, false
}

// gasesteSchimb recunoaște schimbul după cod („Z7”) sau după interval („de la 7 la 14”).
// Întoarce false dacă nu e niciunul sau dacă intervalul se potrivește cu mai multe schimburi.
func gasesteSchimb(mesaj, m string, schimburi []InfoSchimb) (InfoSchimb, bool) {
	cuvinte := strings.FieldsFunc(strings.ToUpper(mesaj), func(r rune) bool {
		return !(r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})
	for _, s := range schimburi {
		for _, c := range cuvinte {
			if s.Cod != "" && c == strings.ToUpper(s.Cod) {
				return s, true
			}
		}
	}
	g := reInterval.FindStringSubmatch(m)
	if g == nil {
		return InfoSchimb{}, false
	}
	ora := func(h, mi string) string {
		x, _ := strconv.Atoi(h)
		if mi == "" {
			mi = "00"
		}
		return fmt.Sprintf("%02d:%s", x, mi)
	}
	cautat := ora(g[1], g[2]) + "–" + ora(g[3], g[4])
	var gasite []InfoSchimb
	for _, s := range schimburi {
		if s.Interval == cautat {
			gasite = append(gasite, s)
		}
	}
	if len(gasite) != 1 {
		return InfoSchimb{}, false
	}
	return gasite[0], true
}

// AnalizeazaComandaGrup recunoaște fără model o comandă de planificare pentru un grup
// („toți angajații”, „secția X”, „categoria Y”) pe o lună întreagă și întoarce o singură
// propunere cu lista persoanelor. Comenzile pentru o zi anume rămân la model.
func AnalizeazaComandaGrup(mesaj string, info ContextPontix) (RaspunsAsistent, bool) {
	m := normalizeaza(mesaj)
	if m == "" || reZiAnume.MatchString(m) {
		return RaspunsAsistent{}, false
	}
	grup, ok := gasesteGrup(m, info)
	if !ok {
		return RaspunsAsistent{}, false
	}
	schimb, ok := gasesteSchimb(mesaj, m, info.Schimburi)
	if !ok {
		return RaspunsAsistent{}, false
	}
	if len(grup.angajati) == 0 {
		return RaspunsAsistent{Mesaj: fmt.Sprintf("Nu am găsit angajați activi pentru %s.", grup.eticheta)}, true
	}

	an, luna := info.An, info.Luna
	for i, l := range luniNorm {
		if regexp.MustCompile(`\b` + l + `\b`).MatchString(m) {
			luna = i + 1
		}
	}
	if an <= 0 || luna < 1 || luna > 12 {
		return RaspunsAsistent{}, false
	}
	filtru, textFiltru := "lucratoare", "doar zile lucrătoare (L–V)"
	switch {
	case strings.Contains(m, "inclusiv weekend") || strings.Contains(m, "toate zilele") || strings.Contains(m, "in fiecare zi") || strings.Contains(m, "si in weekend"):
		filtru, textFiltru = "toate", "toate zilele"
	case strings.Contains(m, "weekend"):
		filtru, textFiltru = "weekend", "doar weekendurile"
	}
	inceput := time.Date(an, time.Month(luna), 1, 0, 0, 0, 0, time.UTC)
	sfarsit := inceput.AddDate(0, 1, -1)

	sort.Slice(grup.angajati, func(i, j int) bool {
		return normalizeaza(grup.angajati[i].Nume+" "+grup.angajati[i].Prenume) < normalizeaza(grup.angajati[j].Nume+" "+grup.angajati[j].Prenume)
	})
	marci := make([]any, len(grup.angajati))
	nume := make([]any, len(grup.angajati))
	for i, a := range grup.angajati {
		marci[i] = a.Marca
		nume[i] = strings.TrimSpace(a.Nume+" "+a.Prenume) + " (" + a.Marca + ")"
	}
	denumire := schimb.Cod
	if schimb.Denumire != "" {
		denumire += " – " + schimb.Denumire
	}
	p := PropunereActiune{
		ID:  fmt.Sprintf("grup_%d", time.Now().UnixNano()),
		Tip: TipPlanificaGrup,
		Descriere: fmt.Sprintf("Planificare schimb %s pentru %d angajați (%s), %s – %s, %s; zilele cu concediu sau absență se sar",
			denumire, len(grup.angajati), grup.eticheta, inceput.Format(time.DateOnly), sfarsit.Format(time.DateOnly), textFiltru),
		Parametri: map[string]any{
			"marci":        marci,
			"nume":         nume,
			"grup":         grup.eticheta,
			"cod_schimb":   schimb.Cod,
			"data_inceput": inceput.Format(time.DateOnly),
			"data_sfarsit": sfarsit.Format(time.DateOnly),
			"filtru_zile":  filtru,
		},
		NecesitaConfirmare: true,
	}
	return RaspunsAsistent{
		Mesaj:     fmt.Sprintf("Am pregătit o propunere pentru %d angajați. Verifică lista și confirmă.", len(grup.angajati)),
		Propuneri: []PropunereActiune{p},
	}, true
}
