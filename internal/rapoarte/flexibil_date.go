package rapoarte

import (
	"context"
	"database/sql"
	"fmt"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"pontix/internal/identitate"
	"sort"
	"strings"
	"time"
)

func incarcatBazaRaport(ctx context.Context, tx *sql.Tx, c Cerere, f *RaportFlexibil, acum time.Time) (*bazaRaport, error) {
	b, e := nouaBazaRaport()
	if e != nil {
		return nil, e
	}
	succes := false
	defer func() {
		if !succes {
			b.db.Close()
		}
	}()
	q := acces.New(tx)
	limba := i18n.Limba(c.Limba)
	start := time.Date(c.An, time.Month(c.Luna), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	deLa := start.Format(time.DateOnly)
	panaLa := end.AddDate(0, 0, -1).Format(time.DateOnly)
	extins := start.AddDate(0, 0, -2).Format(time.DateOnly)
	oameni, e := q.ListaAngajati(ctx)
	if e != nil {
		return nil, e
	}
	profile, asocieri, e := profileRapoarte(ctx, tx, limba)
	if e != nil {
		return nil, e
	}
	norme, e := q.NormeSoldPontaj(ctx)
	if e != nil {
		return nil, e
	}
	rez, _, e := datebaza.CalculeazaPontaj(ctx, q, c.An, c.Luna, nil)
	if e != nil {
		return nil, e
	}
	valori := map[string]map[string]calcul.ValoarePontaj{}
	for _, v := range rez.Totaluri {
		valori[v.AngajatID] = v.Valori
	}
	sectii, e := q.ListaSectii(ctx)
	if e != nil {
		return nil, e
	}
	dep, e := q.ListaDepartamente(ctx)
	if e != nil {
		return nil, e
	}
	ns, nd := map[string]string{}, map[string]string{}
	for _, v := range sectii {
		ns[v.ID] = textJSON(v.Denumire, limba)
	}
	for _, v := range dep {
		nd[v.ID] = textJSON(v.Denumire, limba)
	}
	prefixe, e := prefixeRapoarte(ctx, tx, limba)
	if e != nil {
		return nil, e
	}
	// F07-T43 (D101): norma la inceputul lunii raportului si dreptul la concediu pe anul raportului.
	istoric, e := acces.New(tx).ListaIstoricContract(ctx)
	if e != nil {
		return nil, e
	}
	inceputLuna := fmt.Sprintf("%04d-%02d-01", c.An, c.Luna)
	inclus := map[string]bool{}
	nume := map[string]string{}
	marci := map[string]string{}
	garda := map[string]bool{}
	ra := [][]any{}
	for _, a := range oameni {
		p := profile[a.ID]
		if !p.corespunde(c) || (c.AngajatID != "" && a.ID != c.AngajatID) || !sectieCeruta(c.SectieID, a.SectieID) || (c.DepartamentID != "" && optional(a.DepartamentID) != c.DepartamentID) {
			continue
		}
		inclus[a.ID] = true
		garda[a.ID] = p.Garda || asocieri[a.ID] != ""
		nume[a.ID] = identitate.CuPrefix(prefixe[a.ID], identitate.Nume(a.Nume, a.Prenume, garda[a.ID], limba))
		marci[a.ID] = a.Marca
		var norma any
		sch := make([]calcul.SchimbareContract, 0, len(istoric[a.ID]))
		for _, x := range istoric[a.ID] {
			sch = append(sch, calcul.SchimbareContract{DeLa: x.DeLa, NormaMinute: x.NormaMinute, ConcediuZileAn: x.ConcediuZileAn})
		}
		if n := calcul.NormaLaData(norme[a.ID], sch, inceputLuna); n != nil {
			norma = float64(*n) / 60
		}
		concediu := calcul.DreptConcediuAn(c.An, a.ConcediuZileAn, sch)
		ra = append(ra, []any{a.ID, nume[a.ID], a.Marca, optional(a.Functie), optional(a.SectieID), ns[optional(a.SectieID)], optional(a.DepartamentID), nd[optional(a.DepartamentID)], p.NumeContract, garda[a.ID], norma, a.Activ == 1, a.Nume, a.Prenume, concediu, asocieri[a.ID], prefixe[a.ID]})
	}
	if e = b.adauga(ctx, "angajati", []string{"angajat_id", "angajat", "marca", "functie", "sectie_id", "sectie", "compartiment_id", "compartiment", "contract", "garda", "norma_zilnica", "activ", "nume", "prenume", "concediu_zile_an", "angajat_principal_id", "prefix"}, ra); e != nil {
		return nil, e
	}
	ri := [][]any{}
	ids := make([]string, 0, len(istoric))
	for id := range istoric {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if !inclus[id] {
			continue
		}
		for _, x := range istoric[id] {
			var ore any
			if x.NormaMinute != nil {
				ore = float64(*x.NormaMinute) / 60
			}
			ri = append(ri, []any{id, x.DeLa, ore, x.ConcediuZileAn})
		}
	}
	if e = b.adauga(ctx, "istoric_contract", []string{"angajat_id", "de_la", "norma_zilnica", "concediu_zile_an"}, ri); e != nil {
		return nil, e
	}
	if c.SectieID == FaraSectie {
		ns[FaraSectie] = i18n.T(limba, "ui.fara_sectie", nil)
	}
	b.variabile = map[string]string{"sectie": ns[c.SectieID], "departament": nd[c.DepartamentID], "compartiment": nd[c.DepartamentID], "angajat": nume[c.AngajatID]}
	if c.AngajatID != "" && !inclus[c.AngajatID] {
		return nil, eroare("rapoarte.parametri_invalizi")
	}
	sarb, e := q.ListaToateSarbatori(ctx)
	if e != nil {
		return nil, e
	}
	libere := map[string]bool{}
	for _, s := range sarb {
		libere[s.Data] = true
	}
	b.sarbatori = libere
	intervalStart, intervalEnd := start, end
	if f.ExtindeZileLibere {
		var err error
		intervalStart, intervalEnd, err = calcul.IntervalRaportZileLibere(start, end, libere)
		if err != nil {
			return nil, eroare("rapoarte.interval_invalid")
		}
	}
	b.parametri = parametriFlexibili(c)
	b.parametri["de_la"] = intervalStart.Format(time.DateOnly)
	b.parametri["pana_la"] = intervalEnd.Format(time.DateOnly)
	b.variabile["de_la"] = intervalStart.Format(time.DateOnly)
	b.variabile["pana_la"] = intervalEnd.AddDate(0, 0, -1).Format(time.DateOnly)
	extins = intervalStart.AddDate(0, 0, -2).Format(time.DateOnly)
	garziPana := intervalEnd.AddDate(0, 0, -1).Format(time.DateOnly)
	rc := [][]any{}
	zileLucr := 0
	for zi := intervalStart; zi.Before(intervalEnd); zi = zi.AddDate(0, 0, 1) {
		d := zi.Format(time.DateOnly)
		week := zi.Weekday() == time.Saturday || zi.Weekday() == time.Sunday
		lucr := !week && !libere[d]
		if lucr && !zi.Before(start) && zi.Before(end) {
			zileLucr++
		}
		rc = append(rc, []any{d, week, libere[d], lucr})
	}
	if e = b.adauga(ctx, "calendar", []string{"data", "weekend", "sarbatoare", "lucratoare"}, rc); e != nil {
		return nil, e
	}
	tipuri, e := q.ListaTipuriAbsenta(ctx)
	if e != nil {
		return nil, e
	}
	asoc, e := q.CategoriiTipuriAbsenta(ctx)
	if e != nil {
		return nil, e
	}
	efect, e := q.EfecteAbsente(ctx)
	if e != nil {
		return nil, e
	}
	rt := [][]any{}
	cat := map[string]string{}
	tipSelectat := map[string]bool{}
	for _, x := range f.ConcediiTipuri {
		tipSelectat[x] = true
		// Un model salvat cu un cod redenumit (D092) alege si codul nou.
		if nou, ok := calcul.CoduriRedenumite[x]; ok {
			tipSelectat[nou] = true
		}
	}
	for _, x := range f.ConcediiCategorii {
		cat[x] = x
	}
	for _, t := range tipuri {
		rt = append(rt, []any{t.Cod, textJSON(t.Denumire, limba), asoc[t.Cod], t.Categorie, efect[t.Cod], t.Activ == 1})
		if cat[asoc[t.Cod]] != "" {
			tipSelectat[t.Cod] = true
		}
	}
	if e = b.adauga(ctx, "tipuri_absenta", []string{"cod", "denumire", "categorie_id", "natura", "efect_ore", "activ"}, rt); e != nil {
		return nil, e
	}
	abs, e := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: deLa, PanaLa: panaLa})
	if e != nil {
		return nil, e
	}
	rabs := [][]any{}
	numarTip := map[string]map[string]int{}
	totalAbs := map[string]int{}
	for _, a := range abs {
		if !inclus[a.AngajatID] {
			continue
		}
		rabs = append(rabs, []any{a.AngajatID, a.Data, a.Cod, asoc[a.Cod], a.Activ == 1, a.StareAprobare})
		if a.Activ == 1 && a.StareAprobare == "APROBATA" {
			zi, _ := time.Parse(time.DateOnly, a.Data)
			if zi.Weekday() == time.Saturday || zi.Weekday() == time.Sunday || libere[a.Data] {
				continue
			}
			if numarTip[a.AngajatID] == nil {
				numarTip[a.AngajatID] = map[string]int{}
			}
			numarTip[a.AngajatID][a.Cod]++
			if tipSelectat[a.Cod] {
				totalAbs[a.AngajatID]++
			}
		}
	}
	if e = b.adauga(ctx, "absente", []string{"angajat_id", "data", "cod", "categorie_id", "activ", "stare_aprobare"}, rabs); e != nil {
		return nil, e
	}
	cr, e := tx.QueryContext(ctx, "SELECT cod,denumire,natura,efect_ore,activ FROM categorii_absente ORDER BY cod")
	if e != nil {
		return nil, e
	}
	rct := [][]any{}
	for cr.Next() {
		var cod, den, natura, efect string
		var activ bool
		if e = cr.Scan(&cod, &den, &natura, &efect, &activ); e != nil {
			cr.Close()
			return nil, e
		}
		rct = append(rct, []any{cod, textJSON(den, limba), natura, efect, activ})
	}
	if e = cr.Err(); e != nil {
		cr.Close()
		return nil, e
	}
	cr.Close()
	if e = b.adauga(ctx, "categorii_absenta", []string{"cod", "denumire", "natura", "efect_ore", "activ"}, rct); e != nil {
		return nil, e
	}
	ciorne, e := q.ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: deLa, DataPanaLa: panaLa})
	if e != nil {
		return nil, e
	}
	rci := [][]any{}
	for _, v := range ciorne {
		if inclus[v.AngajatID] {
			rci = append(rci, []any{v.AngajatID, v.Data, v.Valoare})
		}
	}
	if e = b.adauga(ctx, "pontaj_schita", []string{"angajat_id", "data", "valoare"}, rci); e != nil {
		return nil, e
	}
	programe, e := q.ListaSchimburi(ctx)
	if e != nil {
		return nil, e
	}
	pr := map[string]acces.ListaSchimburiRow{}
	rp := [][]any{}
	for _, p := range programe {
		pr[p.ID] = p
		rp = append(rp, []any{p.ID, p.Cod, textJSON(p.Denumire, limba), p.OraStart, p.DurataMinute, p.Tip, p.Regim, p.Activ == 1, culoarePrograme(p)})
	}
	if e = b.adauga(ctx, "programe", []string{"id", "cod", "denumire", "start_minute", "durata_minute", "tip", "regim", "activ", "culoare"}, rp); e != nil {
		return nil, e
	}
	plan, e := q.ListaPlanificare(ctx, acces.ListaPlanificareParams{DataDeLa: extins, DataPanaLa: garziPana})
	if e != nil {
		return nil, e
	}
	rplan, rg, rz := [][]any{}, [][]any{}, [][]any{}
	totalG, eligG := map[string]float64{}, map[string]float64{}
	efective := map[string]bool{}
	// F09-T36 (D099): culoarea afisata a programului; garda din pontajul efectiv, fara program, ia culoarea implicita a garzilor.
	culoareCod := map[string]string{"": datebaza.CuloareProgram("", "GARDA")}
	for _, p := range programe {
		if _, exista := culoareCod[p.Cod]; !exista || p.Activ == 1 {
			culoareCod[p.Cod] = culoarePrograme(p)
		}
	}
	adaugaGarda := func(id, data, cod string, a, z time.Time) {
		limitaStart := time.Date(intervalStart.Year(), intervalStart.Month(), intervalStart.Day(), 0, 0, 0, 0, a.Location())
		limitaEnd := time.Date(intervalEnd.Year(), intervalEnd.Month(), intervalEnd.Day(), 0, 0, 0, 0, a.Location())
		s, l := a, z
		if s.Before(limitaStart) {
			s = limitaStart
		}
		if l.After(limitaEnd) {
			l = limitaEnd
		}
		if !s.Before(l) {
			return
		}
		minute := float64(l.Sub(s) / time.Minute)
		elig := float64(calcul.MinuteGardaEligibile(a, z, limitaStart, limitaEnd, f.Weekend, f.Sarbatori, libere))
		for cur := s; cur.Before(l); {
			lim := time.Date(cur.Year(), cur.Month(), cur.Day()+1, 0, 0, 0, 0, cur.Location())
			if lim.After(l) {
				lim = l
			}
			week := cur.Weekday() == time.Saturday || cur.Weekday() == time.Sunday
			data := cur.Format(time.DateOnly)
			eligZi := float64(calcul.MinuteGardaEligibile(cur, lim, limitaStart, limitaEnd, f.Weekend, f.Sarbatori, libere)) / 60
			rz = append(rz, []any{id, nume[id], marci[id], data, cod, cur.Format("2006-01-02 15:04"), lim.Format("2006-01-02 15:04"), float64(lim.Sub(cur)/time.Minute) / 60, eligZi, week, libere[data], culoareCod[cod]})
			cur = lim
		}
		lunaStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, a.Location())
		lunaEnd := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, a.Location())
		ls, le := a, z
		if ls.Before(lunaStart) {
			ls = lunaStart
		}
		if le.After(lunaEnd) {
			le = lunaEnd
		}
		if ls.Before(le) {
			totalG[id] += float64(le.Sub(ls)/time.Minute) / 60
		}
		eligG[id] += float64(calcul.MinuteGardaEligibile(a, z, lunaStart, lunaEnd, f.Weekend, f.Sarbatori, libere)) / 60
		rg = append(rg, []any{id, nume[id], marci[id], data, cod, a.Format("2006-01-02 15:04"), z.Format("2006-01-02 15:04"), minute / 60, elig / 60, culoareCod[cod]})
	}
	// Pontajele efective active au prioritate pentru aceeasi marca si zi; cele retrase nu intra.
	rows, e := tx.QueryContext(ctx, "SELECT angajat_id,data,start,sfarsit FROM pontaj WHERE sters_la IS NULL AND start IS NOT NULL AND sfarsit IS NOT NULL AND data BETWEEN ? AND ?", extins, garziPana)
	if e != nil {
		return nil, e
	}
	reff := [][]any{}
	local, _ := time.LoadLocation("Europe/Bucharest")
	for rows.Next() {
		var id, d, a, z string
		if e = rows.Scan(&id, &d, &a, &z); e != nil {
			rows.Close()
			return nil, e
		}
		if !inclus[id] {
			continue
		}
		reff = append(reff, []any{id, d, a, z})
		if garda[id] {
			x, ex := time.Parse(time.RFC3339, a)
			y, ey := time.Parse(time.RFC3339, z)
			if ex == nil && ey == nil {
				x = x.In(local)
				y = y.In(local)
				efective[id+"/"+d] = true
				adaugaGarda(id, d, "", x, y)
			}
		}
	}
	if e = rows.Err(); e != nil {
		rows.Close()
		return nil, e
	}
	rows.Close()
	if e = b.adauga(ctx, "pontaj_efectiv", []string{"angajat_id", "data", "inceput", "sfarsit"}, reff); e != nil {
		return nil, e
	}
	for _, p := range plan {
		if !inclus[p.AngajatID] {
			continue
		}
		s := pr[p.SchimbID]
		rplan = append(rplan, []any{p.AngajatID, p.Data, p.SchimbID, s.Cod, p.Sectiune, p.Activ == 1, culoarePrograme(s)})
		if p.Activ != 1 || !garda[p.AngajatID] || efective[p.AngajatID+"/"+p.Data] {
			continue
		}
		d, _ := time.Parse(time.DateOnly, p.Data)
		a := d.Add(time.Duration(s.OraStart) * time.Minute)
		z := a.Add(time.Duration(s.DurataMinute) * time.Minute)
		adaugaGarda(p.AngajatID, p.Data, s.Cod, a, z)
	}
	if e = b.adauga(ctx, "planificare", []string{"angajat_id", "data", "schimb_id", "cod", "sectiune", "activ", "culoare"}, rplan); e != nil {
		return nil, e
	}
	if e = b.adauga(ctx, "garzi", []string{"angajat_id", "angajat", "marca", "data", "program", "inceput", "sfarsit", "ore", "ore_eligibile", "culoare"}, rg); e != nil {
		return nil, e
	}
	if e = b.adauga(ctx, "garzi_zile", []string{"angajat_id", "angajat", "marca", "data", "program", "inceput", "sfarsit", "ore", "ore_eligibile", "weekend", "sarbatoare", "culoare"}, rz); e != nil {
		return nil, e
	}
	cols := []string{"angajat_id", "angajat", "marca", "total_zile_cf_contract", "total_ore_cf_contract", "total_zile_concediu", "total_ore_lucratoare", "total_ore_lucrate", "total_ore_garda", "total_ore_garda_weekend_sarbatori", "L_INCEPUT", "L_FIN"}
	for _, t := range tipuri {
		cod := "concediu_" + strings.ReplaceAll(strings.ReplaceAll(t.Cod, "_", "_5f"), "-", "_2d")
		if !codFlexibil.MatchString(cod) {
			return nil, eroare("rapoarte.flex_invalid")
		}
		cols = append(cols, cod)
	}
	rl := [][]any{}
	for _, a := range oameni {
		if !inclus[a.ID] {
			continue
		}
		v := valori[a.ID]
		ore := func(c string) any {
			x, ok := v[c]
			if !ok || !x.Definita {
				return nil
			}
			return float64(x.Sutimi) / 6000
		}
		var zile, contract any
		if n := norme[a.ID]; n != nil {
			zile = zileLucr
			contract = float64(zileLucr) * float64(*n) / 60
		}
		r := []any{a.ID, nume[a.ID], a.Marca, zile, contract, totalAbs[a.ID], ore("ore_lucratoare"), ore("ore_lucrate"), totalG[a.ID], eligG[a.ID], ore("L_INCEPUT"), ore("L_FIN")}
		for _, t := range tipuri {
			r = append(r, numarTip[a.ID][t.Cod])
		}
		rl = append(rl, r)
	}
	if e = b.adauga(ctx, "lunar", cols, rl); e != nil {
		return nil, e
	}
	// F08-T13 (D092): totalurile Pontajului, inclusiv TZLN si TOCN, intr-un tabel separat, ca
	// rapoartele existente (SELECT * FROM lunar) sa nu se schimbe; apar doar daca SQL-ul le cere.
	colsTot := []string{"angajat_id", "angajat", "marca", "ol", "co", "cm", "cfp", "cve", "dlg", "fmp", "tzln", "tocn"}
	rt2 := [][]any{}
	for _, a := range oameni {
		if !inclus[a.ID] {
			continue
		}
		v := valori[a.ID]
		val := func(c string, divizor float64) any {
			x, ok := v[c]
			if !ok || !x.Definita {
				return nil
			}
			return float64(x.Sutimi) / divizor
		}
		rt2 = append(rt2, []any{a.ID, nume[a.ID], a.Marca, val("OL", 6000), val("CO", 100), val("CM", 100), val("CFP", 100), val("CVE", 100), val("DLG", 100), val("FMP", 100), val("TZLN", 100), val("TOCN", 6000)})
	}
	if e = b.adauga(ctx, "pontaj_totaluri", colsTot, rt2); e != nil {
		return nil, e
	}
	if e = b.adaugaStandard(ctx, tx, c, textFolositSQL(f, f.SQL), acum); e != nil {
		return nil, e
	}
	succes = true
	return b, nil
}
func parametriFlexibili(c Cerere) map[string]any {
	d := time.Date(c.An, time.Month(c.Luna), 1, 0, 0, 0, 0, time.UTC)
	return map[string]any{"an": c.An, "luna": c.Luna, "angajat_id": c.AngajatID, "sectie_id": c.SectieID, "compartiment_id": c.DepartamentID, "de_la": d.Format(time.DateOnly), "pana_la": d.AddDate(0, 1, 0).Format(time.DateOnly)}
}
func sqlImplicit(s string) string {
	switch s {
	case "LISTA_GARZI":
		return "SELECT angajat_id, angajat, marca, data, program, inceput, sfarsit, ore, ore_eligibile FROM garzi ORDER BY angajat, data"
	case "CALCULATE_LUNAR":
		return "SELECT * FROM lunar ORDER BY angajat"
	}
	return "SELECT angajat_id, angajat, marca FROM angajati ORDER BY angajat"
}

// culoarePrograme este culoarea afisata a programului (D097): cea aleasa sau cea implicita dupa tip.
func culoarePrograme(p acces.ListaSchimburiRow) string {
	c := ""
	if p.Culoare != nil {
		c = *p.Culoare
	}
	return datebaza.CuloareProgram(c, p.Tip)
}
