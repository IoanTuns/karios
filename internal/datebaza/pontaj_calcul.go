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
	return calculeazaPontaj(ctx, q, an, luna, modificari, nil)
}

// CalculeazaPontajCuGarzi completeaza doar intrarile lipsa, fara sa scrie in baza.
func CalculeazaPontajCuGarzi(ctx context.Context, q *acces.Queries, an, luna int, garzi []acces.Planificare) (calcul.RezultatPontaj, calcul.ConfiguratiePontaj, error) {
	return calculeazaPontaj(ctx, q, an, luna, nil, garzi)
}
func calculeazaPontaj(ctx context.Context, q *acces.Queries, an, luna int, modificari []CelulaSchita, garzi []acces.Planificare) (calcul.RezultatPontaj, calcul.ConfiguratiePontaj, error) {
	rez, c, err := calculeazaPontajBaza(ctx, q, an, luna, modificari, garzi)
	if err != nil {
		return rez, c, err
	}
	err = completeazaSolduri(ctx, q, an, luna, &rez)
	return rez, c, err
}

func calculeazaPontajBaza(ctx context.Context, q *acces.Queries, an, luna int, modificari []CelulaSchita, garzi []acces.Planificare) (calcul.RezultatPontaj, calcul.ConfiguratiePontaj, error) {
	var gol calcul.RezultatPontaj
	c, err := citesteConfiguratiePontaj(ctx, q)
	if err != nil {
		return gol, c, err
	}
	if err := ValideazaCoduriReguli(ctx, q, c); err != nil {
		return gol, c, err
	}
	// Norma zilnica pe angajat, pentru TOCN (D092).
	if c.Norme, err = q.NormeSoldPontaj(ctx); err != nil {
		return gol, c, err
	}
	// Istoricul normei (D101), pentru norma fiecarei zile la orele suplimentare (D107).
	istoric, err := q.ListaIstoricContract(ctx)
	if err != nil {
		return gol, c, err
	}
	c.IstoricNorme = map[string][]calcul.SchimbareContract{}
	for id, r := range istoric {
		c.IstoricNorme[id] = schimbariCalcul(r)
	}
	// Randurile marcii de garda, pentru regulile din D093.
	if c.Garzi, err = randuriGarda(ctx, q); err != nil {
		return gol, c, err
	}
	if len(c.Conditionate) > 0 {
		oameni, err := q.ListaAngajati(ctx)
		if err != nil {
			return gol, c, err
		}
		// Citirea profilului este furnizata prin aceeasi interogare sqlc/tranzactie.
		profile, err := q.ListaProfilePontaj(ctx)
		if err != nil {
			return gol, c, err
		}
		c.Profile = map[string]calcul.ProfilPontaj{}
		for _, a := range oameni {
			c.Profile[a.ID] = calcul.ProfilPontaj{ContractID: a.TipContract}
		}
		for _, p := range profile {
			c.Profile[p.ID] = calcul.ProfilPontaj{GrupID: p.GrupID, FunctieID: p.FunctieID, ContractID: p.ContractID, Conducere: p.Conducere != 0, Garda: p.Garda != 0}
		}
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
	if len(garzi) > 0 {
		schimburi, err := q.ListaSchimburi(ctx)
		if err != nil {
			return gol, c, err
		}
		coduri := map[string]string{}
		for _, s := range schimburi {
			coduri[s.ID] = s.Cod
		}
		for _, g := range garzi {
			if g.Activ == 1 && g.Sectiune == "GARZI" && g.Data >= extins && g.Data <= panaLa {
				k := g.AngajatID + "/" + g.Data
				cod := coduri[g.SchimbID]
				if vechi := intrari[k].Valoare; vechi != "" {
					cod = calcul.ValoarePrograme(append(calcul.CoduriPrograme(vechi), cod))
				}
				intrari[k] = calcul.IntrarePontaj{AngajatID: g.AngajatID, Data: g.Data, Valoare: cod}
			}
		}
	}
	for _, r := range randuri {
		intrari[r.AngajatID+"/"+r.Data] = calcul.IntrarePontaj{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Valoare}
	}
	verificate := map[string]bool{}
	vazute := map[string]bool{}
	if len(modificari) > 10000 {
		return gol, c, eroare("pontaj.lot_invalid")
	}
	for _, m := range modificari {
		if !dataValida(m.Data) || m.Data < deLa || m.Data > panaLa || m.AngajatID == "" || (utf8.RuneCountInString(m.Valoare) > 1024 || (!strings.Contains(m.Valoare, "+") && utf8.RuneCountInString(m.Valoare) > 16)) || strings.ContainsAny(m.Valoare, "\x00\t\r\n") {
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
	efecte, err := q.EfecteAbsente(ctx)
	if err != nil {
		return gol, c, err
	}
	minute, err := q.MinuteAbsente(ctx, extins, panaLa)
	if err != nil {
		return gol, c, err
	}
	for _, a := range absente {
		if a.Activ == 1 {
			k := a.AngajatID + "/" + a.Data
			r := intrari[k]
			r.AngajatID = a.AngajatID
			r.Data = a.Data
			r.Crediteaza = efecte[a.Cod] == "NORMA" || efecte[a.Cod] == "PROGRAM"
			r.MinuteCreditate = minute[k]
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
	sarbatori, err := q.ListaToateSarbatori(ctx)
	if err != nil {
		return gol, c, err
	}
	libere := map[string]bool{}
	// F04-T11 (D107): zilele libere adaugate de utilizator, pentru optiunea ZNL.
	tipuri, err := q.TipuriSarbatori(ctx)
	if err != nil {
		return gol, c, err
	}
	c.ZileAdaugate = map[string]bool{}
	for d, tip := range tipuri {
		if tip == TipZiLibera {
			c.ZileAdaugate[d] = true
		}
	}
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

// randuriGarda: angajatii al caror rand este marca de garda (contract de garda sau marca
// asociata unei marci principale, D060), din profilul pontajului.
func randuriGarda(ctx context.Context, q *acces.Queries) (map[string]bool, error) {
	profile, err := q.ListaProfilePontaj(ctx)
	if err != nil {
		return nil, err
	}
	rez := map[string]bool{}
	for _, p := range profile {
		if p.Garda != 0 {
			rez[p.ID] = true
		}
	}
	return rez, nil
}
