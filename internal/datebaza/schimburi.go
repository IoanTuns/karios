package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type Schimb struct {
	ID           string    `json:"id"`
	Cod          string    `json:"cod"`
	Denumire     i18n.Text `json:"denumire"`
	OraStart     int64     `json:"ora_start"`
	DurataMinute int64     `json:"durata_minute"`
	PauzaMinute  int64     `json:"pauza_minute"`
	Tip          string    `json:"tip"`
	Activ        bool      `json:"activ"`
}

func (s *Serviciu) ListaSchimburi(ctx context.Context) ([]Schimb, error) {
	randuri, err := acces.New(s.db).ListaSchimburi(ctx)
	if err != nil {
		return nil, traduceEroare(err)
	}
	rez := []Schimb{}
	for _, r := range randuri {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return nil, traduceEroare(err)
		}
		rez = append(rez, Schimb{ID: r.ID, Cod: r.Cod, Denumire: denumire, OraStart: r.OraStart, DurataMinute: r.DurataMinute, PauzaMinute: r.PauzaMinute, Tip: r.Tip, Activ: r.Activ == 1})
	}
	return rez, nil
}

func (s *Serviciu) SalveazaSchimb(ctx context.Context, schimb Schimb) error {
	schimb.Cod = strings.ToUpper(strings.TrimSpace(schimb.Cod))
	if schimb.Cod == "" || len(schimb.Cod) > 16 || strings.IndexFunc(schimb.Cod, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_')
	}) >= 0 {
		return eroare("schimburi.cod_invalid")
	}
	denumire := i18n.Text{}
	for l, text := range schimb.Denumire {
		if text = strings.TrimSpace(text); text != "" {
			denumire[l] = text
		}
	}
	if len(denumire) == 0 {
		return eroare("sarbatori.denumire_obligatorie")
	}
	schimb.Denumire = denumire
	if schimb.OraStart < 0 || schimb.OraStart > 1439 || schimb.DurataMinute < 1 || schimb.DurataMinute > 2880 || schimb.PauzaMinute < 0 || schimb.PauzaMinute > schimb.DurataMinute {
		return eroare("schimburi.interval_invalid")
	}
	switch schimb.Tip {
	case "NORMAL", "TURA", "NOAPTE", "GARDA":
	default:
		return eroare("schimburi.tip_invalid")
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		idExistent, err := q.CitesteCodSchimb(ctx, schimb.Cod)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if err == nil && idExistent != schimb.ID {
			return eroare("schimburi.cod_duplicat")
		}
		text, err := json.Marshal(denumire)
		if err != nil {
			return err
		}
		acum := moment()
		var vechi any
		if schimb.ID == "" {
			schimb.ID = idNou()
			err = q.AdaugaSchimb(ctx, acces.AdaugaSchimbParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, PauzaMinute: schimb.PauzaMinute, Tip: schimb.Tip, Activ: numarBool(schimb.Activ), CreatLa: acum, ModificatLa: acum})
		} else {
			r, citireErr := q.CitesteSchimb(ctx, schimb.ID)
			if errors.Is(citireErr, sql.ErrNoRows) {
				return eroare("schimburi.inexistent")
			}
			if citireErr != nil {
				return citireErr
			}
			vechi = r
			_, err = q.ModificaSchimb(ctx, acces.ModificaSchimbParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, PauzaMinute: schimb.PauzaMinute, Tip: schimb.Tip, Activ: numarBool(schimb.Activ), ModificatLa: acum})
		}
		if err != nil {
			return err
		}
		return jurnal(ctx, q, "tipuri_schimb", schimb.ID, vechi, schimb)
	})
}

func (s *Serviciu) SchimbaActivSchimb(ctx context.Context, id string, activ bool) error {
	return s.modifica(ctx, func(q *acces.Queries) error {
		vechi, err := q.CitesteSchimb(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return eroare("schimburi.inexistent")
		}
		if err != nil {
			return err
		}
		if _, err := q.SchimbaActivSchimb(ctx, acces.SchimbaActivSchimbParams{ID: id, Activ: numarBool(activ), ModificatLa: moment()}); err != nil {
			return err
		}
		nou, err := q.CitesteSchimb(ctx, id)
		if err != nil {
			return err
		}
		return jurnal(ctx, q, "tipuri_schimb", id, vechi, nou)
	})
}
