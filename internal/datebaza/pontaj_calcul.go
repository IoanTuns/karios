package datebaza

import (
	"context"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

// Adaptor comun grilei și rapoartelor; toate citirile pot folosi aceeași tranzacție.
func CalculeazaPontaj(ctx context.Context, q *acces.Queries, an, luna int, modificari []CelulaSchita) (calcul.RezultatPontaj, calcul.ConfiguratiePontaj, error) {
	var gol calcul.RezultatPontaj
	c, err := citesteConfiguratiePontaj(ctx, q)
	if err != nil {
		return gol, c, err
	}
	if err := ValideazaCoduriReguli(ctx, q, c); err != nil {
		return gol, c, err
	}
	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return gol, c, err
	}
	inceput, _ := time.Parse(time.DateOnly, deLa)
	extins := inceput.AddDate(0, 0, -2).Format(time.DateOnly)
	randuri, err := q.ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: extins, DataPanaLa: panaLa})
	if err != nil {
		return gol, c, err
	}
	intrari := map[string]calcul.IntrarePontaj{}
	for _, r := range randuri {
		intrari[r.AngajatID+"/"+r.Data] = calcul.IntrarePontaj{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Valoare}
	}
	verificate := map[string]bool{}
	vazute := map[string]bool{}
	if len(modificari) > 10000 {
		return gol, c, eroare("pontaj.lot_invalid")
	}
	for _, m := range modificari {
		if !dataValida(m.Data) || m.Data < deLa || m.Data > panaLa || m.AngajatID == "" || utf8.RuneCountInString(m.Valoare) > 16 || strings.ContainsAny(m.Valoare, "\x00\t\r\n") {
			return gol, c, eroare("pontaj.celula_invalida")
		}
		k := m.AngajatID + "/" + m.Data
		if vazute[k] {
			return gol, c, eroare("pontaj.lot_invalid")
		}
		vazute[k] = true
		if !verificate[m.AngajatID] {
			if _, err := q.CitesteAngajat(ctx, m.AngajatID); err != nil {
				return gol, c, err
			}
			verificate[m.AngajatID] = true
		}
		intrari[k] = calcul.IntrarePontaj{AngajatID: m.AngajatID, Data: m.Data, Valoare: m.Valoare}
	}
	absente, err := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: extins, PanaLa: panaLa})
	if err != nil {
		return gol, c, err
	}
	for _, a := range absente {
		if a.Activ == 1 {
			k := a.AngajatID + "/" + a.Data
			r := intrari[k]
			r.AngajatID = a.AngajatID
			r.Data = a.Data
			r.Absenta = a.Cod
			r.StareAprobare = a.StareAprobare
			intrari[k] = r
		}
	}
	schimburi, err := q.ListaSchimburi(ctx)
	if err != nil {
		return gol, c, err
	}
	programe := []calcul.ProgramPontaj{}
	for _, s := range schimburi {
		programe = append(programe, calcul.ProgramPontaj{Cod: s.Cod, OraStart: s.OraStart, DurataMinute: s.DurataMinute, PauzaMinute: s.PauzaMinute})
	}
	sarbatori, err := q.ListaSarbatori(ctx, nil)
	if err != nil {
		return gol, c, err
	}
	libere := map[string]bool{}
	for _, s := range sarbatori {
		libere[s.Data] = true
	}
	lista := []calcul.IntrarePontaj{}
	for _, r := range intrari {
		lista = append(lista, r)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].AngajatID+lista[i].Data < lista[j].AngajatID+lista[j].Data })
	rez, err := calcul.CalculeazaPontaj(deLa, panaLa, c, programe, lista, libere)
	return rez, c, err
}
