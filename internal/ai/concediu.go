package ai

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TipDeschideConcediu e propunerea care deschide formularul de concediu precompletat (F14-T26,
// D052). Asistentul nu salvează concediul: numărul de înregistrare îl scrie utilizatorul, iar
// cererea se salvează prin formularul obișnuit, cu starea „în așteptare”.
const TipDeschideConcediu = "deschide_concediu"

var (
	cuvinteConcediu = []string{"concediu", "concedii", "odihna", "medical", "boala", "absenta", "absent", "liber"}
	verbeConcediu   = []string{"merge", "pleaca", "intra", "ia", "are", "va fi", "este in", "e in", "sta", "lipseste", "vrea", "cere", "pune", "trece"}
	reIntrebare     = regexp.MustCompile(`^(?:cine|cati|cate|care|ce|exista|este cineva|e cineva|cand|who|how many|which|is anyone)\b`)
	reNumarZi       = regexp.MustCompile(`\b(\d{1,2})\b(?:\s+(ianuarie|februarie|martie|aprilie|mai|iunie|iulie|august|septembrie|octombrie|noiembrie|decembrie))?`)
)

// esteIntrebare deosebește „cine e în concediu?” de „X merge în concediu”.
func esteIntrebare(original, m string) bool {
	return strings.HasSuffix(strings.TrimSpace(original), "?") || reIntrebare.MatchString(m)
}

func contineOricare(m string, lista []string) bool {
	for _, x := range lista {
		if contineExpresie(m, x) {
			return true
		}
	}
	return false
}

// EsteComandaConcediu: un mesaj care nu e întrebare, pomenește un concediu și are un verb de acțiune.
func EsteComandaConcediu(mesaj string, info ContextPontix) bool {
	m, _ := corecteazaMesaj(mesaj, info)
	return !esteIntrebare(mesaj, m) && contineOricare(m, cuvinteConcediu) && contineOricare(m, verbeConcediu)
}

// gasesteAngajatInMesaj caută angajatul al cărui nume complet apare în mesaj, cuvânt cu cuvânt,
// cu toleranța la greșeli din F14-T15 („Maria Daria” → Marin Daria), sau după marcă.
// Întoarce false dacă nu e niciunul sau sunt mai mulți la fel de apropiați.
func gasesteAngajatInMesaj(mesaj, m string, angajati []InfoAngajat) (InfoAngajat, []InfoAngajat) {
	cuvinte := strings.Fields(m)
	marci := strings.Fields(strings.ToUpper(mesaj))
	type candidat struct {
		a InfoAngajat
		d int
	}
	var buni []candidat
	for _, a := range angajati {
		for _, x := range marci {
			if a.Marca != "" && strings.Trim(x, ".,;:") == strings.ToUpper(a.Marca) {
				return a, nil
			}
		}
		total, gasit := 0, true
		for _, parte := range strings.Fields(normalizeaza(a.Nume + " " + a.Prenume)) {
			best := -1
			for _, c := range cuvinte {
				d := distantaDL(parte, c)
				if d <= pragGreseli(len([]rune(parte))) && (best < 0 || d < best) {
					best = d
				}
			}
			if best < 0 {
				gasit = false
				break
			}
			total += best
		}
		if gasit {
			buni = append(buni, candidat{a, total})
		}
	}
	if len(buni) == 0 {
		return InfoAngajat{}, nil
	}
	min := buni[0].d
	for _, c := range buni {
		if c.d < min {
			min = c.d
		}
	}
	var cei []InfoAngajat
	for _, c := range buni {
		if c.d == min {
			cei = append(cei, c.a)
		}
	}
	if len(cei) == 1 {
		return cei[0], nil
	}
	return InfoAngajat{}, cei
}

// intervalDin citește primele două zile din mesaj („26 pana pe 31”, „din 26 octombrie pana pe 3
// noiembrie”), cu luna scrisă lângă zi sau, altfel, luna numită în mesaj ori luna de lucru.
func intervalDin(m string, info ContextPontix) (time.Time, time.Time, bool) {
	gasite := reNumarZi.FindAllStringSubmatch(m, -1)
	var zile []int
	var luni []int
	for _, g := range gasite {
		z, _ := strconv.Atoi(g[1])
		if z < 1 || z > 31 {
			continue
		}
		luna := 0
		for i, l := range luniNorm {
			if g[2] == l {
				luna = i + 1
			}
		}
		zile, luni = append(zile, z), append(luni, luna)
		if len(zile) == 2 {
			break
		}
	}
	if len(zile) == 0 {
		return time.Time{}, time.Time{}, false
	}
	if len(zile) == 1 {
		zile, luni = append(zile, zile[0]), append(luni, luni[0])
	}
	lunaImplicita := info.Luna
	if l, ok := lunaDin(m); ok {
		lunaImplicita = l
	}
	if luni[1] == 0 {
		luni[1] = lunaImplicita
	}
	if luni[0] == 0 {
		luni[0] = luni[1]
	}
	if info.An <= 0 || luni[0] < 1 || luni[0] > 12 {
		return time.Time{}, time.Time{}, false
	}
	deLa := time.Date(info.An, time.Month(luni[0]), zile[0], 0, 0, 0, 0, time.UTC)
	panaLa := time.Date(info.An, time.Month(luni[1]), zile[1], 0, 0, 0, 0, time.UTC)
	if panaLa.Before(deLa) {
		panaLa = panaLa.AddDate(1, 0, 0) // „din 28 decembrie până pe 3 ianuarie”
	}
	if deLa.Day() != zile[0] || panaLa.Day() != zile[1] || panaLa.Sub(deLa) > 366*24*time.Hour {
		return time.Time{}, time.Time{}, false
	}
	return deLa, panaLa, true
}

// AnalizeazaCerereConcediu recunoaște o cerere de concediu („Maria Daria merge în concediu
// în octombrie, 26–31”) și propune deschiderea formularului precompletat (F14-T26). Întoarce
// true și pentru o comandă incompletă, cu o întrebare de clarificare, ca mesajul să nu ajungă
// la lista concediilor sau la model.
func AnalizeazaCerereConcediu(mesaj string, info ContextPontix) (RaspunsAsistent, bool) {
	if !EsteComandaConcediu(mesaj, info) {
		return RaspunsAsistent{}, false
	}
	m, corecturi := corecteazaMesaj(mesaj, info)
	// Numele se caută în textul necorectat: corectura vocabularului nu trebuie să atingă numele.
	a, candidati := gasesteAngajatInMesaj(mesaj, normalizeaza(mesaj), info.Angajati)
	if len(candidati) > 1 {
		nume := make([]string, len(candidati))
		for i, c := range candidati {
			nume[i] = strings.TrimSpace(c.Nume+" "+c.Prenume) + " (" + c.Marca + ")"
		}
		return RaspunsAsistent{Mesaj: "Am găsit mai mulți angajați potriviți: " + strings.Join(nume, ", ") + ". Scrie marca sau numele complet."}, true
	}
	if a.ID == "" {
		return RaspunsAsistent{Mesaj: "Pentru ce angajat este concediul? Scrie numele complet sau marca."}, true
	}
	deLa, panaLa, ok := intervalDin(m, info)
	if !ok {
		return RaspunsAsistent{Mesaj: fmt.Sprintf("Pentru ce zile este concediul lui %s %s? De exemplu: „din 26 până pe 31 octombrie”.", a.Nume, a.Prenume)}, true
	}
	cod := "CO"
	if contineOricare(m, []string{"medical", "boala", "cm"}) {
		cod = "CM"
	}
	nume := strings.TrimSpace(a.Nume + " " + a.Prenume)
	interpretare := fmt.Sprintf("%s (%s) · %s · %s – %s", nume, a.Marca, cod, deLa.Format(time.DateOnly), panaLa.Format(time.DateOnly))
	return RaspunsAsistent{
		Mesaj: "Deschid formularul de concediu precompletat. Completează numărul de înregistrare și salvează; cererea rămâne în așteptare până la aprobare.",
		Propuneri: []PropunereActiune{{
			ID:        fmt.Sprintf("concediu_%d", time.Now().UnixNano()),
			Tip:       TipDeschideConcediu,
			Descriere: "Cerere de concediu: " + interpretare,
			Parametri: map[string]any{
				"angajat_id":   a.ID,
				"marca":        a.Marca,
				"nume":         nume,
				"cod":          cod,
				"data_inceput": deLa.Format(time.DateOnly),
				"data_sfarsit": panaLa.Format(time.DateOnly),
				"interpretare": interpretare,
				"corecturi":    textCorecturi(corecturi),
			},
			NecesitaConfirmare: true,
		}},
	}, true
}
