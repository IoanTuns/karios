package datebaza

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"pontix/internal/calcul"
)

// F12-T06: „Verifică luna” adună într-un loc ce trebuie văzut înainte de închiderea lunii
// (procedura de control lunar, F02-T23). Nu schimbă nimic.

const (
	GravitateEroare      = "EROARE"
	GravitateAvertisment = "AVERTISMENT"
	GravitateInfo        = "INFO"
	maxExemple           = 10
)

// ProblemaLuna: o constatare, cu numărul de cazuri și câteva exemple (mărci sau coduri).
type ProblemaLuna struct {
	Cod       string            `json:"cod"`
	Gravitate string            `json:"gravitate"`
	Numar     int               `json:"numar"`
	Exemple   []string          `json:"exemple"`
	Param     map[string]string `json:"param,omitempty"`
}

type VerificareLuna struct {
	An       int            `json:"an"`
	Luna     int            `json:"luna"`
	Probleme []ProblemaLuna `json:"probleme"`
}

func exemple(lista []string) []string {
	sort.Strings(lista)
	if len(lista) > maxExemple {
		return lista[:maxExemple]
	}
	return lista
}

// VerificaLuna întoarce constatările pentru luna dată, cele mai grave primele.
func (s *Serviciu) VerificaLuna(ctx context.Context, an, luna int) (VerificareLuna, error) {
	rez := VerificareLuna{An: an, Luna: luna, Probleme: []ProblemaLuna{}}
	date, err := s.Citeste(ctx)
	if err != nil {
		return rez, err
	}
	schita, err := s.CitesteSchita(ctx, an, luna)
	if err != nil {
		return rez, err
	}
	marci := map[string]string{}
	for _, a := range date.Angajati {
		marci[a.ID] = a.Marca
	}
	adauga := func(p ProblemaLuna) {
		if p.Numar > 0 {
			rez.Probleme = append(rez.Probleme, p)
		}
	}

	// 1. Valori necalculate, pe motiv (cele care nu se aplică rândului nu sunt probleme).
	type cheie struct{ motiv, cod string }
	peMotiv := map[cheie]map[string]bool{}
	for _, t := range schita.Calcul.Totaluri {
		for _, v := range t.Valori {
			if v.Definita || v.NuSeAplica || v.Motiv == "" {
				continue
			}
			k := cheie{v.Motiv, v.MotivCod}
			if peMotiv[k] == nil {
				peMotiv[k] = map[string]bool{}
			}
			peMotiv[k][marci[t.AngajatID]] = true
		}
	}
	chei := make([]cheie, 0, len(peMotiv))
	for k := range peMotiv {
		chei = append(chei, k)
	}
	sort.Slice(chei, func(i, j int) bool { return chei[i].motiv+chei[i].cod < chei[j].motiv+chei[j].cod })
	for _, k := range chei {
		lista := []string{}
		for m := range peMotiv[k] {
			lista = append(lista, m)
		}
		adauga(ProblemaLuna{Cod: "necalculat", Gravitate: GravitateEroare, Numar: len(lista), Exemple: exemple(lista), Param: map[string]string{"motiv": k.motiv, "cod": k.cod}})
	}

	// 2. Programe active cu pauză, fără modul pauzei ales.
	programe, err := s.ListaSchimburi(ctx)
	if err != nil {
		return rez, err
	}
	pauze := []string{}
	for _, p := range programe {
		if p.Activ && p.PauzaMinute > 0 && (p.ModPauza == "" || p.ModPauza == calcul.PauzaNestabilita) {
			pauze = append(pauze, p.Cod)
		}
	}
	adauga(ProblemaLuna{Cod: "pauza_nestabilita", Gravitate: GravitateAvertisment, Numar: len(pauze), Exemple: exemple(pauze)})

	// 3. Absențe în așteptarea aprobării.
	asteptare, zileAsteptare := map[string]bool{}, 0
	for _, z := range schita.Absente {
		if z.Activ && z.StareAprobare == "IN_ASTEPTARE" {
			zileAsteptare++
			asteptare[marci[z.AngajatID]] = true
		}
	}
	lista := []string{}
	for m := range asteptare {
		lista = append(lista, m)
	}
	adauga(ProblemaLuna{Cod: "absente_in_asteptare", Gravitate: GravitateAvertisment, Numar: zileAsteptare, Exemple: exemple(lista)})

	// 4. Anul lunii fără nicio sărbătoare în calendar.
	cal, err := s.CitesteCalendar(ctx, an)
	if err != nil {
		return rez, err
	}
	if len(cal.Sarbatori) == 0 {
		adauga(ProblemaLuna{Cod: "fara_sarbatori", Gravitate: GravitateAvertisment, Numar: 1, Exemple: []string{}, Param: map[string]string{"an": strconv.Itoa(an)}})
	}

	// 5. Angajați activi fără nicio înregistrare în lună (nici ore, nici absențe).
	cuDate := map[string]bool{}
	for _, c := range schita.Celule {
		if strings.TrimSpace(c.Valoare) != "" {
			cuDate[c.AngajatID] = true
		}
	}
	for _, z := range schita.Absente {
		if z.Activ {
			cuDate[z.AngajatID] = true
		}
	}
	fara := []string{}
	for _, a := range date.Angajati {
		if a.Activ && !cuDate[a.ID] {
			fara = append(fara, a.Marca)
		}
	}
	adauga(ProblemaLuna{Cod: "fara_inregistrari", Gravitate: GravitateInfo, Numar: len(fara), Exemple: exemple(fara)})

	// 6. Luna deja închisă sau blocată.
	if schita.Stare != "" && schita.Stare != "DESCHISA" {
		adauga(ProblemaLuna{Cod: "luna_inchisa", Gravitate: GravitateInfo, Numar: 1, Exemple: []string{}, Param: map[string]string{"stare": schita.Stare}})
	}

	ordine := map[string]int{GravitateEroare: 0, GravitateAvertisment: 1, GravitateInfo: 2}
	sort.SliceStable(rez.Probleme, func(i, j int) bool { return ordine[rez.Probleme[i].Gravitate] < ordine[rez.Probleme[j].Gravitate] })
	return rez, nil
}
