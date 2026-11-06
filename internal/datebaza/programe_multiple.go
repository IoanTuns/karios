package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"strconv"
)

func valoareZiPlan(s sursePlanificare, angajat, data string) string {
	coduri := []string{}
	for _, p := range s.pozitii {
		if p.Activ == 1 && p.AngajatID == angajat && p.Data == data {
			coduri = append(coduri, s.schimburi[p.SchimbID].Cod)
		}
	}
	return calcul.ValoarePrograme(coduri)
}

// Anularea unei pozitii publicate retrage numai contributia sa din proiectie.
func retrageProgramPublicat(ctx context.Context, q *acces.Queries, s sursePlanificare, v acces.Planificare) error {
	c, ok := s.ciorne[v.AngajatID+"/"+v.Data]
	if !ok || v.PublicatLa == "" {
		return nil
	}
	cod := s.schimburi[v.SchimbID].Cod
	ramase := []string{}
	scos := false
	for _, x := range calcul.CoduriPrograme(c.Valoare) {
		if x == cod && !scos {
			scos = true
		} else {
			ramase = append(ramase, x)
		}
	}
	if !scos {
		return nil
	}
	val := calcul.ValoarePrograme(ramase)
	rev, err := q.SalveazaCiorna(ctx, acces.SalveazaCiornaParams{AngajatID: v.AngajatID, Data: v.Data, Valoare: val, ModificatLa: moment()})
	if err != nil {
		return err
	}
	if err = jurnal(ctx, q, "pontaj_ciorne", v.AngajatID+"/"+v.Data, c, CelulaSchita{AngajatID: v.AngajatID, Data: v.Data, Valoare: val, Revizie: strconv.FormatInt(rev, 10)}); err != nil {
		return err
	}
	// Setul sursei este partajat cu lotul, astfel incat doua anulari nu se restaureaza reciproc.
	c.Valoare = val
	c.Revizie = rev
	s.ciorne[v.AngajatID+"/"+v.Data] = c
	return nil
}
