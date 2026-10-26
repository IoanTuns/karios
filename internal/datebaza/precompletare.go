package datebaza

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

func (s *Serviciu) PrevizualizeazaPrecompletare(ctx context.Context, an, luna int, ids []string, modificari []CelulaSchita) (calcul.PlanPrecompletare, error) {
	var gol calcul.PlanPrecompletare
	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return gol, err
	}
	if len(ids) == 0 || len(ids) > 1000 || len(modificari) > 10000 {
		return gol, eroare("pontaj.lot_invalid")
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return gol, traduceEroare(err)
	}
	defer tx.Rollback()
	q := acces.New(tx)
	oameni, err := q.ListaAngajati(ctx)
	if err != nil {
		return gol, traduceEroare(err)
	}
	norma, err := q.CitesteSetare(ctx, "norma_intreaga_minute")
	if err != nil {
		return gol, traduceEroare(err)
	}
	// Valoarea este validată la salvare și păstrată în minute întregi.
	normaIntreaga, err := strconv.ParseInt(norma, 10, 64)
	if err != nil {
		return gol, traduceEroare(err)
	}
	selectati := map[string]bool{}
	for _, id := range ids {
		if id == "" || selectati[id] {
			return gol, eroare("pontaj.lot_invalid")
		}
		selectati[id] = true
	}
	angajati := []calcul.AngajatPrecompletare{}
	for _, a := range oameni {
		if !selectati[a.ID] {
			continue
		}
		minute := normaIntreaga
		if a.TipContract == "PARTIALA" && a.NormaPartialaMinute != nil {
			minute = *a.NormaPartialaMinute
		}
		angajati = append(angajati, calcul.AngajatPrecompletare{ID: a.ID, CategorieID: textSauGol(a.CategorieProgramID), NormaMinute: minute, Activ: a.Activ == 1})
	}
	if len(angajati) != len(ids) {
		return gol, eroare("pontaj.celula_invalida")
	}
	categorii, err := citesteCategoriiProgram(ctx, q)
	if err != nil {
		return gol, traduceEroare(err)
	}
	profile := []calcul.CategoriePrecompletare{}
	for _, c := range categorii {
		profile = append(profile, calcul.CategoriePrecompletare{ID: c.ID, ImplicitID: c.SchimbImplicitID, Activ: c.Activ})
	}
	schimburi, err := q.ListaSchimburi(ctx)
	if err != nil {
		return gol, traduceEroare(err)
	}
	programe := []calcul.ProgramPrecompletare{}
	for _, p := range schimburi {
		programe = append(programe, calcul.ProgramPrecompletare{ID: p.ID, Cod: p.Cod, Tip: p.Tip, Regim: p.Regim, OraStart: p.OraStart, DurataMinute: p.DurataMinute, PauzaMinute: p.PauzaMinute, Activ: p.Activ == 1})
	}
	prima, _ := time.Parse(time.DateOnly, deLa)
	ultima, _ := time.Parse(time.DateOnly, panaLa)
	extinsDeLa, extinsPanaLa := prima.AddDate(0, 0, -2).Format(time.DateOnly), ultima.AddDate(0, 0, 2).Format(time.DateOnly)
	ciorne, err := q.ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: extinsDeLa, DataPanaLa: extinsPanaLa})
	if err != nil {
		return gol, traduceEroare(err)
	}
	celule := map[string]calcul.IntrarePontaj{}
	for _, c := range ciorne {
		celule[c.AngajatID+"/"+c.Data] = calcul.IntrarePontaj{AngajatID: c.AngajatID, Data: c.Data, Valoare: c.Valoare}
	}
	vazute := map[string]bool{}
	for _, c := range modificari {
		k := c.AngajatID + "/" + c.Data
		if !dataValida(c.Data) || c.Data < deLa || c.Data > panaLa || c.AngajatID == "" || utf8.RuneCountInString(c.Valoare) > 16 || strings.ContainsAny(c.Valoare, "\x00\t\r\n") || vazute[k] {
			return gol, eroare("pontaj.celula_invalida")
		}
		vazute[k] = true
		celule[k] = calcul.IntrarePontaj{AngajatID: c.AngajatID, Data: c.Data, Valoare: c.Valoare}
	}
	intrari := []calcul.IntrarePontaj{}
	for _, c := range celule {
		intrari = append(intrari, c)
	}
	absente, err := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: extinsDeLa, PanaLa: extinsPanaLa})
	if err != nil {
		return gol, traduceEroare(err)
	}
	zileAbsente := map[string]bool{}
	for _, a := range absente {
		if a.Activ == 1 {
			zileAbsente[a.AngajatID+"/"+a.Data] = true
		}
	}
	libere := map[string]bool{}
	for _, anCalendar := range []int{an - 1, an, an + 1} {
		anSelectat := int64(anCalendar)
		sarbatori, err := q.ListaSarbatori(ctx, &anSelectat)
		if err != nil {
			return gol, traduceEroare(err)
		}
		for _, z := range sarbatori {
			libere[z.Data] = true
		}
	}
	efective, err := q.ListaPontajePrecompletare(ctx, acces.ListaPontajePrecompletareParams{DataDeLa: extinsDeLa, DataPanaLa: extinsPanaLa})
	if err != nil {
		return gol, traduceEroare(err)
	}
	ocupate := map[string]bool{}
	for _, p := range efective {
		ocupate[p.AngajatID+"/"+p.Data] = true
		if p.Start == nil || p.Sfarsit == nil {
			continue
		}
		start, e1 := time.Parse(time.RFC3339, *p.Start)
		stop, e2 := time.Parse(time.RFC3339, *p.Sfarsit)
		if e1 != nil || e2 != nil || !stop.After(start) {
			continue
		}
		a := start.In(calcul.Locatie).Format(time.DateOnly)
		b := stop.Add(-time.Nanosecond).In(calcul.Locatie).Format(time.DateOnly)
		if a < extinsDeLa {
			a = extinsDeLa
		}
		if b > extinsPanaLa {
			b = extinsPanaLa
		}
		for z, _ := time.Parse(time.DateOnly, a); z.Format(time.DateOnly) <= b; z = z.AddDate(0, 0, 1) {
			ocupate[p.AngajatID+"/"+z.Format(time.DateOnly)] = true
		}
	}
	rez, err := calcul.PlanificaPrecompletare(deLa, panaLa, angajati, profile, programe, intrari, libere, zileAbsente, ocupate)
	if err != nil {
		return gol, err
	}
	if err = tx.Commit(); err != nil {
		return gol, traduceEroare(err)
	}
	return rez, nil
}
