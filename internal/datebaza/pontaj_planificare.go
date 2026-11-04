package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"strings"
	"time"
)

type schimbarePontajPlan struct {
	Celula  CelulaSchita
	Inainte string
}

// Sincronizarea modifica setul zilei in aceeasi tranzactie cu Pontajul.
func sincronizeazaPontajPlan(ctx context.Context, q *acces.Queries, deLa, panaLa string, modificari []schimbarePontajPlan) error {
	surse, err := sursePentruPlanificare(ctx, q, deLa, panaLa)
	if err != nil {
		return err
	}
	verifica := map[string]bool{}
	for _, m := range modificari {
		c := m.Celula
		existente := []acces.Planificare{}
		for _, p := range surse.pozitii {
			if p.Activ == 1 && p.AngajatID == c.AngajatID && p.Data == c.Data {
				existente = append(existente, p)
			}
		}
		dorite := map[string]calcul.ProgramPrecompletare{}
		for _, cod := range calcul.CoduriPrograme(c.Valoare) {
			p, ok := surse.coduri[cod]
			pastrat := false
			for _, v := range existente {
				pastrat = pastrat || v.SchimbID == p.ID
			}
			if ok && (calcul.ProgramDePlanificat(p) || pastrat || strings.Contains(c.Valoare, "+")) {
				dorite[p.ID] = p
			}
		}
		corespunde := calcul.AceleasiPrograme(m.Inainte, valoareZiPlan(surse, c.AngajatID, c.Data))
		// Un plan independent de ciorna nu se pierde la editarea manuala.
		for _, v := range existente {
			if _, pastrat := dorite[v.SchimbID]; pastrat {
				continue
			}
			if !corespunde {
				if len(dorite) > 0 {
					return eroare("pontaj.plan_independent")
				}
				continue
			}
			n, e := q.ModificaPozitiePlanificare(ctx, acces.ModificaPozitiePlanificareParams{ID: v.ID, Revizie: v.Revizie, SchimbID: v.SchimbID, Activ: 0, ModificatLa: moment()})
			if e != nil {
				return e
			}
			if n != 1 {
				return eroare("planificare.conflict")
			}
			nou := v
			nou.Activ = 0
			nou.Revizie++
			nou.ModificatLa = moment()
			if e = jurnal(ctx, q, "planificare", v.ID, v, nou); e != nil {
				return e
			}
		}
		for id, program := range dorite {
			p := acces.Planificare{}
			for _, v := range existente {
				if v.SchimbID == id {
					p = v
					break
				}
			}
			if p.ID == "" {
				sectiune := "GARZI"
				if calcul.ProgramInSectiune(program, "DUPA_AMIAZA") {
					sectiune = "DUPA_AMIAZA"
				}
				p = acces.Planificare{ID: idNou(), AngajatID: c.AngajatID, Data: c.Data, SchimbID: id, Sectiune: sectiune, Activ: 1, Revizie: 1, CreatLa: moment(), ModificatLa: moment()}
				if e := q.AdaugaPozitiePlanificare(ctx, acces.AdaugaPozitiePlanificareParams{ID: p.ID, AngajatID: p.AngajatID, Data: p.Data, SchimbID: id, Sectiune: sectiune, Activ: 1, CreatLa: p.CreatLa, ModificatLa: p.ModificatLa}); e != nil {
					return e
				}
				if e := jurnal(ctx, q, "planificare", p.ID, nil, p); e != nil {
					return e
				}
			}
			verifica[p.ID] = true
			ciorna := surse.ciorne[c.AngajatID+"/"+c.Data]
			if p.PublicatPlanRevizie == p.Revizie && p.PublicatPontajRevizie == ciorna.Revizie && p.PublicatValoare == c.Valoare {
				continue
			}
			nou := p
			nou.PublicatPlanRevizie = p.Revizie
			nou.PublicatPontajRevizie = ciorna.Revizie
			nou.PublicatValoare = c.Valoare
			nou.PublicatLa = moment()
			n, e := q.MarcheazaPlanificarePublicata(ctx, acces.MarcheazaPlanificarePublicataParams{ID: p.ID, Revizie: p.Revizie, PublicatPlanRevizie: p.Revizie, PublicatPontajRevizie: ciorna.Revizie, PublicatValoare: c.Valoare, PublicatLa: nou.PublicatLa})
			if e != nil {
				return e
			}
			if n != 1 {
				return eroare("planificare.conflict")
			}
			if e = jurnal(ctx, q, "planificare", p.ID, p, nou); e != nil {
				return e
			}
		}
	}
	finale, err := sursePentruPlanificare(ctx, q, deLa, panaLa)
	if err != nil {
		return err
	}
	for _, p := range finale.pozitii {
		if verifica[p.ID] {
			if motiv := finale.verifica(p, true, true); motiv != "" {
				return eroare("planificare.motiv_" + motiv)
			}
		}
	}
	return nil
}

// Valideaza un set explicit; codurile simple necunoscute isi pastreaza comportamentul vechi.
func valideazaSetPrograme(s sursePlanificare, c CelulaSchita) error {
	coduri := calcul.CoduriPrograme(c.Valoare)
	if len(coduri) < 2 {
		vecinMultiplu := false
		for _, v := range s.ciorne {
			if v.AngajatID == c.AngajatID && strings.Contains(v.Valoare, "+") {
				vecinMultiplu = true
			}
		}
		if !vecinMultiplu || c.Valoare == "" {
			return nil
		}
		if _, ok := s.coduri[coduri[0]]; !ok {
			return nil
		}
	}
	vazute := map[string]bool{}
	intervale := []calcul.IntervalPlanificare{}
	for _, i := range s.intervale {
		if i.Cheie != "c/"+c.AngajatID+"/"+c.Data && s.cheiEfective[i.Cheie] != c.AngajatID+"/"+c.Data {
			intervale = append(intervale, i)
		}
	}
	for _, cod := range coduri {
		p, ok := s.coduri[cod]
		if !ok || !p.Activ || vazute[cod] {
			return eroare("pontaj.celula_invalida")
		}
		vazute[cod] = true
		if motiv := calcul.VerificaTuraPlanificata(c.AngajatID, c.Data, p, s.absente, intervale, nil); motiv != "" {
			return eroare("planificare.motiv_" + motiv)
		}
		d, _ := time.Parse(time.DateOnly, c.Data)
		a := d.Add(time.Duration(p.OraStart) * time.Minute)
		intervale = append(intervale, calcul.IntervalPlanificare{AngajatID: c.AngajatID, Inceput: a, Sfarsit: a.Add(time.Duration(p.DurataMinute) * time.Minute)})
	}
	return nil
}

// Intervalele lotului inlocuiesc proiectiile vechi inaintea validarii, fara scrieri.
func intervaleLotPontaj(s *sursePlanificare, celule []CelulaSchita) {
	chei := map[string]bool{}
	for _, c := range celule {
		chei["c/"+c.AngajatID+"/"+c.Data] = true
		s.ciorne[c.AngajatID+"/"+c.Data] = acces.PontajCiorne{AngajatID: c.AngajatID, Data: c.Data, Valoare: c.Valoare}
	}
	pastrate := []calcul.IntervalPlanificare{}
	for _, i := range s.intervale {
		if !chei[i.Cheie] {
			pastrate = append(pastrate, i)
		}
	}
	for _, c := range celule {
		d, _ := time.Parse(time.DateOnly, c.Data)
		for _, cod := range calcul.CoduriPrograme(c.Valoare) {
			if cod == "" {
				continue
			}
			p, ok := s.coduri[cod]
			a, z := d, d.AddDate(0, 0, 1)
			if ok && p.OraStart >= 0 && p.DurataMinute > 0 {
				a = d.Add(time.Duration(p.OraStart) * time.Minute)
				z = a.Add(time.Duration(p.DurataMinute) * time.Minute)
			}
			pastrate = append(pastrate, calcul.IntervalPlanificare{AngajatID: c.AngajatID, Cheie: "c/" + c.AngajatID + "/" + c.Data, Inceput: a, Sfarsit: z})
		}
	}
	s.intervale = pastrate
}
