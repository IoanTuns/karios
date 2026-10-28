// Package datebaza gestionează datele angajaților și validările pentru interfață.
package datebaza

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ncruces/go-sqlite3"

	"pontix/internal/audit"
	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type Angajat struct {
	ID                  string `json:"id"`
	Nume                string `json:"nume"`
	Prenume             string `json:"prenume"`
	Marca               string `json:"marca"`
	Functie             string `json:"functie"`
	SectieID            string `json:"sectie_id"`
	DepartamentID       string `json:"departament_id"`
	CategorieProgramID  string `json:"categorie_program_id"`
	TipContract         string `json:"tip_contract"`
	NormaPartialaMinute *int64 `json:"norma_partiala_minute"`
	NormaMinute         int64  `json:"norma_minute"`
	ConcediuZileAn      *int64 `json:"concediu_zile_an"`
	Activ               bool   `json:"activ"`
}

type Sectie struct {
	ID       string    `json:"id"`
	Denumire i18n.Text `json:"denumire"`
}

type Departament struct {
	ID       string    `json:"id"`
	SectieID string    `json:"sectie_id"`
	Denumire i18n.Text `json:"denumire"`
}

type Date struct {
	Angajati            []Angajat          `json:"angajati"`
	Sectii              []Sectie           `json:"sectii"`
	Departamente        []Departament      `json:"departamente"`
	CategoriiProgram    []CategorieProgram `json:"categorii_program"`
	NormaIntreagaMinute int64              `json:"norma_intreaga_minute"`
	DemoDisponibil      bool               `json:"demo_disponibil"`
}

type Serviciu struct{ db *sql.DB }

func Nou(db *sql.DB) *Serviciu { return &Serviciu{db: db} }

func (s *Serviciu) Citeste(ctx context.Context) (Date, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	defer tx.Rollback()
	q := acces.New(tx)
	rez := Date{Angajati: []Angajat{}, Sectii: []Sectie{}, Departamente: []Departament{}}
	valoare, err := q.CitesteSetare(ctx, "norma_intreaga_minute")
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	rez.NormaIntreagaMinute, err = strconv.ParseInt(valoare, 10, 64)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	angajati, err := q.ListaAngajati(ctx)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	for _, a := range angajati {
		norma := rez.NormaIntreagaMinute
		if a.TipContract == "PARTIALA" && a.NormaPartialaMinute != nil {
			norma = *a.NormaPartialaMinute
		}
		rez.Angajati = append(rez.Angajati, Angajat{ID: a.ID, Nume: a.Nume, Prenume: a.Prenume,
			Marca: a.Marca, Functie: textSauGol(a.Functie), SectieID: textSauGol(a.SectieID),
			DepartamentID:      textSauGol(a.DepartamentID),
			CategorieProgramID: textSauGol(a.CategorieProgramID),
			TipContract:        a.TipContract, NormaPartialaMinute: a.NormaPartialaMinute,
			NormaMinute: norma, ConcediuZileAn: a.ConcediuZileAn, Activ: a.Activ == 1})
	}
	sectii, err := q.ListaSectii(ctx)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	for _, r := range sectii {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return Date{}, traduceEroare(err)
		}
		rez.Sectii = append(rez.Sectii, Sectie{ID: r.ID, Denumire: denumire})
	}
	departamente, err := q.ListaDepartamente(ctx)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	for _, r := range departamente {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return Date{}, traduceEroare(err)
		}
		rez.Departamente = append(rez.Departamente, Departament{ID: r.ID, SectieID: r.SectieID, Denumire: denumire})
	}
	rez.CategoriiProgram, err = citesteCategoriiProgram(ctx, q)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	n, err := q.NumaraAngajati(ctx)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	p, err := q.NumaraPontaje(ctx)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	rez.DemoDisponibil = n == 0 && p == 0
	if err := tx.Commit(); err != nil {
		return Date{}, traduceEroare(err)
	}
	return rez, nil
}

func (s *Serviciu) SalveazaAngajat(ctx context.Context, a Angajat) error {
	a.Nume, a.Prenume, a.Marca = strings.TrimSpace(a.Nume), strings.TrimSpace(a.Prenume), strings.TrimSpace(a.Marca)
	a.Functie = strings.TrimSpace(a.Functie)
	if a.Nume == "" || a.Prenume == "" || a.Marca == "" {
		return eroare("date.campuri_obligatorii")
	}
	if a.TipContract != "INTREAGA" && a.TipContract != "PARTIALA" {
		return eroare("date.contract_invalid")
	}
	if a.TipContract == "INTREAGA" {
		a.NormaPartialaMinute = nil
	} else if a.NormaPartialaMinute == nil || *a.NormaPartialaMinute < 1 || *a.NormaPartialaMinute > 1440 {
		return eroare("date.norma_invalida")
	}
	if a.ConcediuZileAn != nil && *a.ConcediuZileAn < 0 {
		return eroare("date.concediu_invalid")
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		if a.CategorieProgramID != "" {
			c, err := q.CitesteCategorieProgram(ctx, a.CategorieProgramID)
			if err != nil {
				return eroare("program.categorie_invalida")
			}
			if c.Activ != 1 {
				vechi, err := q.CitesteAngajat(ctx, a.ID)
				if err != nil || textSauGol(vechi.CategorieProgramID) != a.CategorieProgramID {
					return eroare("program.categorie_inactiva")
				}
			}
		}
		if a.DepartamentID != "" {
			d, err := q.CitesteDepartamentActiv(ctx, a.DepartamentID)
			if errors.Is(err, sql.ErrNoRows) || (err == nil && d.SectieID != a.SectieID) {
				return eroare("date.departament_invalid")
			}
			if err != nil {
				return err
			}
		}
		acum := moment()
		var vechi any
		if a.ID == "" {
			a.ID = idNou()
			if err := q.AdaugaAngajat(ctx, acces.AdaugaAngajatParams{ID: a.ID, Nume: a.Nume, Prenume: a.Prenume,
				Marca: a.Marca, Functie: textOptional(a.Functie), SectieID: textOptional(a.SectieID),
				DepartamentID:      textOptional(a.DepartamentID),
				CategorieProgramID: textOptional(a.CategorieProgramID),
				TipContract:        a.TipContract, NormaPartialaMinute: a.NormaPartialaMinute,
				ConcediuZileAn: a.ConcediuZileAn, Activ: numarBool(a.Activ), CreatLa: acum, ModificatLa: acum}); err != nil {
				return err
			}
		} else {
			r, err := q.CitesteAngajat(ctx, a.ID)
			if err != nil {
				return err
			}
			vechi = r
			n, err := q.ModificaAngajat(ctx, acces.ModificaAngajatParams{ID: a.ID, Nume: a.Nume, Prenume: a.Prenume,
				Marca: a.Marca, Functie: textOptional(a.Functie), SectieID: textOptional(a.SectieID),
				DepartamentID:      textOptional(a.DepartamentID),
				CategorieProgramID: textOptional(a.CategorieProgramID),
				TipContract:        a.TipContract, NormaPartialaMinute: a.NormaPartialaMinute,
				ConcediuZileAn: a.ConcediuZileAn, Activ: numarBool(a.Activ), ModificatLa: acum})
			if err != nil {
				return err
			}
			if n != 1 {
				return sql.ErrNoRows
			}
		}
		nou, err := q.CitesteAngajat(ctx, a.ID)
		if err != nil {
			return err
		}
		return jurnal(ctx, q, "angajati", a.ID, vechi, nou)
	})
}

func (s *Serviciu) SchimbaActiv(ctx context.Context, id string, activ bool) error {
	return s.modifica(ctx, func(q *acces.Queries) error {
		vechi, err := q.CitesteAngajat(ctx, id)
		if err != nil {
			return err
		}
		if _, err := q.SchimbaActivAngajat(ctx, acces.SchimbaActivAngajatParams{ID: id, Activ: numarBool(activ), ModificatLa: moment()}); err != nil {
			return err
		}
		nou, err := q.CitesteAngajat(ctx, id)
		if err != nil {
			return err
		}
		return jurnal(ctx, q, "angajati", id, vechi, nou)
	})
}

func (s *Serviciu) AtribuieCategorieProgramLot(ctx context.Context, angajatIDs []string, categorieID string) error {
	categorieID = strings.TrimSpace(categorieID)
	return s.modifica(ctx, func(q *acces.Queries) error {
		if categorieID != "" {
			if _, err := q.CitesteCategorieProgram(ctx, categorieID); err != nil {
				return eroare("program.categorie_invalida")
			}
		}
		for _, id := range angajatIDs {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			vechi, err := q.CitesteAngajat(ctx, id)
			if err != nil {
				return err
			}
			catVal := textOptional(categorieID)
			n, err := q.ModificaAngajat(ctx, acces.ModificaAngajatParams{
				ID:                  vechi.ID,
				Nume:                vechi.Nume,
				Prenume:             vechi.Prenume,
				Marca:               vechi.Marca,
				Functie:             vechi.Functie,
				SectieID:            vechi.SectieID,
				DepartamentID:       vechi.DepartamentID,
				CategorieProgramID:  catVal,
				TipContract:         vechi.TipContract,
				NormaPartialaMinute: vechi.NormaPartialaMinute,
				ConcediuZileAn:      vechi.ConcediuZileAn,
				Activ:               vechi.Activ,
				ModificatLa:         moment(),
			})
			if err != nil {
				return err
			}
			if n != 1 {
				return sql.ErrNoRows
			}
			nou, err := q.CitesteAngajat(ctx, id)
			if err != nil {
				return err
			}
			if err := jurnal(ctx, q, "angajati", id, vechi, nou); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Serviciu) SchimbaActivLot(ctx context.Context, ids []string, activ bool) error {
	return s.modifica(ctx, func(q *acces.Queries) error {
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			vechi, err := q.CitesteAngajat(ctx, id)
			if err != nil {
				return err
			}
			if _, err := q.SchimbaActivAngajat(ctx, acces.SchimbaActivAngajatParams{ID: id, Activ: numarBool(activ), ModificatLa: moment()}); err != nil {
				return err
			}
			nou, err := q.CitesteAngajat(ctx, id)
			if err != nil {
				return err
			}
			if err := jurnal(ctx, q, "angajati", id, vechi, nou); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Serviciu) SalveazaNorma(ctx context.Context, minute int64) error {
	if minute < 1 || minute > 1440 {
		return eroare("date.norma_invalida")
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		veche, err := q.CitesteSetare(ctx, "norma_intreaga_minute")
		if err != nil {
			return err
		}
		noua := strconv.FormatInt(minute, 10)
		if _, err := q.ActualizeazaSetare(ctx, acces.ActualizeazaSetareParams{Cheie: "norma_intreaga_minute", Valoare: noua, ModificatLa: moment()}); err != nil {
			return err
		}
		return jurnal(ctx, q, "setari", "norma_intreaga_minute", veche, noua)
	})
}

func (s *Serviciu) AdaugaSectie(ctx context.Context, denumire i18n.Text) error {
	curata := i18n.Text{}
	for limba, text := range denumire {
		if text = strings.TrimSpace(text); text != "" {
			curata[limba] = text
		}
	}
	if len(curata) == 0 {
		return eroare("date.sectie_obligatorie")
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		id := idNou()
		if err := adaugaSectie(ctx, q, Sectie{ID: id, Denumire: curata}); err != nil {
			return err
		}
		return jurnal(ctx, q, "sectii", id, nil, curata)
	})
}

func adaugaSectie(ctx context.Context, q *acces.Queries, s Sectie) error {
	text, err := json.Marshal(s.Denumire)
	if err != nil {
		return err
	}
	acum := moment()
	return q.AdaugaSectie(ctx, acces.AdaugaSectieParams{ID: s.ID, Denumire: string(text), CreatLa: acum, ModificatLa: acum})
}

func (s *Serviciu) AdaugaDepartament(ctx context.Context, sectieID string, denumire i18n.Text) error {
	curata := i18n.Text{}
	for l, text := range denumire {
		if text = strings.TrimSpace(text); text != "" {
			curata[l] = text
		}
	}
	if len(curata) == 0 {
		return eroare("date.departament_obligatoriu")
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		if _, err := q.CitesteSectieActiva(ctx, sectieID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return eroare("date.sectie_invalida")
			}
			return err
		}
		id := idNou()
		text, err := json.Marshal(curata)
		if err != nil {
			return err
		}
		if err := q.AdaugaDepartament(ctx, acces.AdaugaDepartamentParams{ID: id, SectieID: sectieID, Denumire: string(text), CreatLa: moment(), ModificatLa: moment()}); err != nil {
			return err
		}
		return jurnal(ctx, q, "departamente", id, nil, Departament{ID: id, SectieID: sectieID, Denumire: curata})
	})
}

func (s *Serviciu) modifica(ctx context.Context, actiune func(*acces.Queries) error) error {
	return s.modificaTx(ctx, func(_ *sql.Tx, q *acces.Queries) error { return actiune(q) })
}

// modificaTx dă acțiunii și tranzacția, pentru interogările scrise direct în SQL.
func (s *Serviciu) modificaTx(ctx context.Context, actiune func(*sql.Tx, *acces.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return traduceEroare(err)
	}
	defer tx.Rollback()
	if err := actiune(tx, acces.New(tx)); err != nil {
		return traduceEroare(err)
	}
	return traduceEroare(tx.Commit())
}

func jurnal(ctx context.Context, q *acces.Queries, tabel, id string, vechi, nou any) error {
	var inainte *string
	if vechi != nil {
		text, err := json.Marshal(vechi)
		if err != nil {
			return err
		}
		v := string(text)
		inainte = &v
	}
	text, err := json.Marshal(nou)
	if err != nil {
		return err
	}
	v := string(text)
	return q.AdaugaJurnal(ctx, acces.AdaugaJurnalParams{Tabel: tabel, RandID: id, ValoareVeche: inainte, ValoareNoua: &v, Cand: moment(), Cine: "local"})
}

func eroare(cod string) error { return &i18n.Eroare{Cod: cod} }

func traduceEroare(err error) error {
	if err == nil {
		return nil
	}
	var e *i18n.Eroare
	if errors.As(err, &e) {
		return e
	}
	switch {
	case audit.EsteEditareBlocata(err):
		return eroare(audit.CodEditareBlocata)
	case errors.Is(err, sql.ErrNoRows):
		return eroare("date.inexistent")
	case errors.Is(err, sqlite3.CONSTRAINT_UNIQUE):
		return eroare("date.marca_duplicata")
	case errors.Is(err, sqlite3.CONSTRAINT_FOREIGNKEY):
		return eroare("date.sectie_invalida")
	case errors.Is(err, sqlite3.BUSY), errors.Is(err, sqlite3.LOCKED):
		return eroare("date.ocupata")
	default:
		return eroare("date.eroare")
	}
}

func textOptional(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func textSauGol(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func numarBool(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func moment() string { return time.Now().UTC().Format(time.RFC3339Nano) }

func idNou() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6], b[8] = (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
