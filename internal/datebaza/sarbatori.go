package datebaza

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ncruces/go-sqlite3"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type Sarbatoare struct {
	Data     string    `json:"data"`
	Denumire i18n.Text `json:"denumire"`
}

type CalendarAn struct {
	An               int          `json:"an"`
	Sarbatori        []Sarbatoare `json:"sarbatori"`
	ImportDisponibil bool         `json:"import_disponibil"`
}

//go:embed calendare/romania-2025.json
var calendar2025 []byte

//go:embed calendare/romania-2026.json
var calendar2026 []byte

//go:embed calendare/romania-2027.json
var calendar2027 []byte

func (s *Serviciu) CitesteCalendar(ctx context.Context, an int) (CalendarAn, error) {
	if an < 1 || an > 9999 {
		return CalendarAn{}, eroare("sarbatori.an_invalid")
	}
	anSelectat := int64(an)
	randuri, err := acces.New(s.db).ListaSarbatori(ctx, &anSelectat)
	if err != nil {
		return CalendarAn{}, traduceEroare(err)
	}
	rez := CalendarAn{An: an, Sarbatori: []Sarbatoare{}, ImportDisponibil: an == 2025 || an == 2026 || an == 2027}
	for _, r := range randuri {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return CalendarAn{}, traduceEroare(err)
		}
		rez.Sarbatori = append(rez.Sarbatori, Sarbatoare{Data: r.Data, Denumire: denumire})
	}
	return rez, nil
}

func dataValida(data string) bool {
	t, err := time.Parse(time.DateOnly, data)
	return err == nil && t.Year() >= 1 && t.Format(time.DateOnly) == data
}

func (s *Serviciu) SalveazaSarbatoare(ctx context.Context, initiala string, sarbatoare Sarbatoare) error {
	if !dataValida(sarbatoare.Data) {
		return eroare("sarbatori.data_invalida")
	}
	denumire := i18n.Text{}
	for l, text := range sarbatoare.Denumire {
		if text = strings.TrimSpace(text); text != "" {
			denumire[l] = text
		}
	}
	if len(denumire) == 0 {
		return eroare("sarbatori.denumire_obligatorie")
	}
	sarbatoare.Denumire = denumire
	return s.modifica(ctx, func(q *acces.Queries) error {
		text, err := json.Marshal(denumire)
		if err != nil {
			return err
		}
		var veche any
		var n int64
		if initiala == "" {
			n, err = q.AdaugaSarbatoare(ctx, acces.AdaugaSarbatoareParams{Data: sarbatoare.Data, Denumire: string(text), CreatLa: moment(), ModificatLa: moment()})
		} else {
			r, citireErr := q.CitesteSarbatoare(ctx, initiala)
			if errors.Is(citireErr, sql.ErrNoRows) {
				return eroare("sarbatori.inexistenta")
			}
			if citireErr != nil {
				return citireErr
			}
			veche = r
			n, err = q.ModificaSarbatoare(ctx, acces.ModificaSarbatoareParams{DataNoua: sarbatoare.Data, DataInitiala: initiala, Denumire: string(text), ModificatLa: moment()})
		}
		if errors.Is(err, sqlite3.CONSTRAINT_PRIMARYKEY) || errors.Is(err, sqlite3.CONSTRAINT_UNIQUE) || (err == nil && n == 0) {
			return eroare("sarbatori.data_duplicata")
		}
		if err != nil {
			return err
		}
		return jurnal(ctx, q, "sarbatori", sarbatoare.Data, veche, sarbatoare)
	})
}

func (s *Serviciu) CompleteazaSarbatoriAn(ctx context.Context, an int) (int64, error) {
	var dateCalendar []byte
	switch an {
	case 2025:
		dateCalendar = calendar2025
	case 2026:
		dateCalendar = calendar2026
	case 2027:
		dateCalendar = calendar2027
	default:
		return 0, eroare("sarbatori.an_fara_import")
	}
	var sarbatori []Sarbatoare
	if err := json.Unmarshal(dateCalendar, &sarbatori); err != nil {
		return 0, traduceEroare(err)
	}
	var adaugate int64
	err := s.modifica(ctx, func(q *acces.Queries) error {
		for _, sarbatoare := range sarbatori {
			text, err := json.Marshal(sarbatoare.Denumire)
			if err != nil {
				return err
			}
			n, err := q.AdaugaSarbatoare(ctx, acces.AdaugaSarbatoareParams{Data: sarbatoare.Data, Denumire: string(text), CreatLa: moment(), ModificatLa: moment()})
			if err != nil {
				return err
			}
			if n == 0 {
				continue
			}
			adaugate += n
			if err := jurnal(ctx, q, "sarbatori", sarbatoare.Data, nil, sarbatoare); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return adaugate, nil
}

func (s *Serviciu) CompleteazaSarbatori2026(ctx context.Context) (int64, error) {
	return s.CompleteazaSarbatoriAn(ctx, 2026)
}
