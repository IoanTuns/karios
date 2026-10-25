package datebaza

import (
	"context"
	_ "embed"
	"encoding/json"

	"pontix/internal/bd/acces"
)

//go:embed demo/angajati.json
var dateDemo []byte

func (s *Serviciu) IncarcaDemo(ctx context.Context) error {
	var date struct {
		Sectii   []Sectie  `json:"sectii"`
		Angajati []Angajat `json:"angajati"`
	}
	if err := json.Unmarshal(dateDemo, &date); err != nil {
		return traduceEroare(err)
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		n, err := q.NumaraAngajati(ctx)
		if err != nil {
			return err
		}
		p, err := q.NumaraPontaje(ctx)
		if err != nil {
			return err
		}
		if n != 0 || p != 0 {
			return eroare("date.demo_nepermis")
		}
		for _, sectie := range date.Sectii {
			if err := adaugaSectie(ctx, q, sectie); err != nil {
				return err
			}
			if err := jurnal(ctx, q, "sectii", sectie.ID, nil, sectie.Denumire); err != nil {
				return err
			}
		}
		for _, a := range date.Angajati {
			acum := moment()
			if err := q.AdaugaAngajat(ctx, acces.AdaugaAngajatParams{ID: a.ID, Nume: a.Nume, Prenume: a.Prenume,
				Marca: a.Marca, Functie: textOptional(a.Functie), SectieID: textOptional(a.SectieID),
				DepartamentID: textOptional(a.DepartamentID),
				TipContract:   a.TipContract, NormaPartialaMinute: a.NormaPartialaMinute,
				ConcediuZileAn: a.ConcediuZileAn, Activ: numarBool(a.Activ), CreatLa: acum, ModificatLa: acum}); err != nil {
				return err
			}
			if err := jurnal(ctx, q, "angajati", a.ID, nil, a); err != nil {
				return err
			}
		}
		return nil
	})
}
