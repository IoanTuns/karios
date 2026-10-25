package datebaza

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"pontix/internal/bd/acces"
)

type CelulaSchita struct {
	AngajatID string `json:"angajat_id"`
	Data      string `json:"data"`
	Valoare   string `json:"valoare"`
	Revizie   string `json:"revizie"`
}

type SchitaLuna struct {
	An      int            `json:"an"`
	Luna    int            `json:"luna"`
	Celule  []CelulaSchita `json:"celule"`
	Absente []ZiAbsenta    `json:"absente"`
}

func limiteLuna(an, luna int) (string, string, error) {
	if an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return "", "", eroare("pontaj.perioada_invalida")
	}
	return fmt.Sprintf("%04d-%02d-01", an, luna), time.Date(an, time.Month(luna)+1, 0, 0, 0, 0, 0, time.UTC).Format(time.DateOnly), nil
}

func (s *Serviciu) CitesteSchita(ctx context.Context, an, luna int) (SchitaLuna, error) {
	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return SchitaLuna{}, err
	}
	randuri, err := acces.New(s.db).ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: deLa, DataPanaLa: panaLa})
	if err != nil {
		return SchitaLuna{}, traduceEroare(err)
	}
	rez := SchitaLuna{An: an, Luna: luna, Celule: []CelulaSchita{}}
	for _, r := range randuri {
		rez.Celule = append(rez.Celule, CelulaSchita{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Valoare, Revizie: r.ModificatLa})
	}
	rez.Absente, err = citesteAbsenteInterval(ctx, acces.New(s.db), deLa, panaLa)
	if err != nil {
		return SchitaLuna{}, traduceEroare(err)
	}
	return rez, nil
}

func (s *Serviciu) SalveazaSchita(ctx context.Context, an, luna int, celule []CelulaSchita) error {
	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return err
	}
	if len(celule) > 10000 {
		return eroare("pontaj.lot_invalid")
	}
	vazute := map[string]bool{}
	for _, c := range celule {
		if !dataValida(c.Data) || c.Data < deLa || c.Data > panaLa || c.AngajatID == "" || utf8.RuneCountInString(c.Valoare) > 16 || strings.ContainsAny(c.Valoare, "\x00\t\r\n") {
			return eroare("pontaj.celula_invalida")
		}
		cheie := c.AngajatID + "/" + c.Data
		if vazute[cheie] {
			return eroare("pontaj.lot_invalid")
		}
		vazute[cheie] = true
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		angajati := map[string]bool{}
		for _, c := range celule {
			absenta, err := q.CitesteZiAbsenta(ctx, acces.CitesteZiAbsentaParams{AngajatID: c.AngajatID, Data: c.Data})
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			if err == nil && absenta.Activ == 1 {
				return eroare("absente.celula_protejata")
			}
			if !angajati[c.AngajatID] {
				if _, err := q.CitesteAngajat(ctx, c.AngajatID); err != nil {
					return err
				}
				angajati[c.AngajatID] = true
			}
			veche, err := q.CitesteCiorna(ctx, acces.CitesteCiornaParams{AngajatID: c.AngajatID, Data: c.Data})
			exista := err == nil
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			if (exista && c.Revizie != veche.ModificatLa) || (!exista && c.Revizie != "") {
				return eroare("pontaj.conflict")
			}
			c.Valoare = strings.TrimSpace(c.Valoare)
			if (exista && c.Valoare == veche.Valoare) || (!exista && c.Valoare == "") {
				continue
			}
			c.Revizie = moment()
			if err := q.SalveazaCiorna(ctx, acces.SalveazaCiornaParams{AngajatID: c.AngajatID, Data: c.Data, Valoare: c.Valoare, ModificatLa: c.Revizie}); err != nil {
				return err
			}
			var inainte any
			if exista {
				inainte = veche
			}
			if err := jurnal(ctx, q, "pontaj_ciorne", c.AngajatID+"/"+c.Data, inainte, c); err != nil {
				return err
			}
		}
		return nil
	})
}
