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
	ContractID          string `json:"contract_id"`
	GrupID              string `json:"grup_id"`
	Conducere           bool   `json:"conducere"`
	Garda               bool   `json:"garda"`
	NormaDefinita       bool   `json:"norma_definita"`
	ID                  string `json:"id"`
	Nume                string `json:"nume"`
	Prenume             string `json:"prenume"`
	Marca               string `json:"marca"`
	Functie             string `json:"functie"`
	FunctieID           string `json:"functie_id"`
	AngajatPrincipalID  string `json:"angajat_principal_id"`
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
	GrupuriProfesionale []ElementNomenclator `json:"grupuri_profesionale"`
	TipuriContract      []ElementNomenclator `json:"tipuri_contract"`
	CategoriiSchimburi  []ElementNomenclator `json:"categorii_schimburi"`
	Functii             []ElementNomenclator `json:"functii"`
	Angajati            []Angajat            `json:"angajati"`
	Sectii              []Sectie             `json:"sectii"`
	Departamente        []Departament        `json:"departamente"`
	CategoriiProgram    []CategorieProgram   `json:"categorii_program"`
	NormaIntreagaMinute int64                `json:"norma_intreaga_minute"`
	DemoDisponibil      bool                 `json:"demo_disponibil"`
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
	rez.CategoriiSchimburi, err = citesteNomenclator(ctx, tx, "categorii")
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	rez.Functii, err = citesteNomenclator(ctx, tx, "functii")
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	legaturi, err := tx.QueryContext(ctx, "SELECT angajat_id,functie_id FROM angajati_functii")
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	asocieri := map[string]string{}
	for legaturi.Next() {
		var a, f string
		if err := legaturi.Scan(&a, &f); err != nil {
			legaturi.Close()
			return Date{}, traduceEroare(err)
		}
		asocieri[a] = f
	}
	err = legaturi.Err()
	legaturi.Close()
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	for i := range rez.Angajati {
		rez.Angajati[i].FunctieID = asocieri[rez.Angajati[i].ID]
		for _, f := range rez.Functii {
			if f.ID == rez.Angajati[i].FunctieID {
				rez.Angajati[i].Functie = numeElement(f)
			}
		}
	}
	garzi, err := citesteLegaturiGarzi(ctx, tx)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	for i := range rez.Angajati {
		rez.Angajati[i].AngajatPrincipalID = garzi[rez.Angajati[i].ID]
	}
	rez.GrupuriProfesionale, err = citesteNomenclator(ctx, tx, "grupuri")
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	rez.TipuriContract, err = citesteNomenclator(ctx, tx, "contracte")
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	contracte, err := citesteContracteAngajati(ctx, tx)
	if err != nil {
		return Date{}, traduceEroare(err)
	}
	for i := range rez.Angajati {
		a := &rez.Angajati[i]
		a.NormaDefinita = true
		a.Garda = a.AngajatPrincipalID != ""
		a.ContractID = contracte[a.ID]
		if a.ContractID == "" {
			a.ContractID = a.TipContract
		}
		for _, f := range rez.Functii {
			if f.ID == a.FunctieID {
				a.GrupID = f.GrupID
				a.Conducere = f.Conducere
			}
		}
		for _, c := range rez.TipuriContract {
			if c.ID == a.ContractID {
				a.Garda = a.Garda || c.Garda
				if c.CategorieProgramID != "" {
					a.CategorieProgramID = c.CategorieProgramID
				}
				switch c.Normare {
				case "GLOBALA":
					a.NormaMinute = rez.NormaIntreagaMinute
				case "INDIVIDUALA":
					if a.NormaPartialaMinute == nil {
						a.NormaDefinita = false
						a.NormaMinute = 0
					} else {
						a.NormaMinute = *a.NormaPartialaMinute
					}
				case "FARA", "LUNARA":
					a.NormaDefinita = false
					a.NormaMinute = 0
				case "ZILNICA":
					a.NormaMinute = *c.NormaMinute
				}
			}
		}
	}
	indexAngajati := map[string]Angajat{}
	for _, a := range rez.Angajati {
		indexAngajati[a.ID] = a
	}
	for i := range rez.Angajati {
		a := &rez.Angajati[i]
		if a.AngajatPrincipalID != "" && a.FunctieID == "" {
			p := indexAngajati[a.AngajatPrincipalID]
			a.FunctieID = p.FunctieID
			a.Functie = p.Functie
			a.GrupID = p.GrupID
			a.Conducere = p.Conducere
		}
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
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		if (a.ContractID == "INTREAGA" || a.ContractID == "PARTIALA") && (a.TipContract == "INTREAGA" || a.TipContract == "PARTIALA") && a.TipContract != a.ContractID {
			a.ContractID = a.TipContract
		}
		if a.ContractID != "" {
			c, err := contractPentruSalvare(ctx, tx, a)
			if err != nil {
				return err
			}
			if c.Normare == "INDIVIDUALA" {
				a.TipContract = "PARTIALA"
			} else {
				a.TipContract = "INTREAGA"
			}
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

		if a.ContractID != "" {
			if err := valideazaContractActiv(ctx, tx, a); err != nil {
				return err
			}
		}
		if err := valideazaLegaturaGarda(ctx, tx, a); err != nil {
			return err
		}
		// Profilul returnat UI poate include functia mostenita; o salvare neutra nu o fixeaza separat.
		if a.AngajatPrincipalID != "" && a.FunctieID != "" {
			var proprie, principala sql.NullString
			if err := tx.QueryRowContext(ctx, "SELECT (SELECT functie_id FROM angajati_functii WHERE angajat_id=?),(SELECT functie_id FROM angajati_functii WHERE angajat_id=?)", a.ID, a.AngajatPrincipalID).Scan(&proprie, &principala); err != nil {
				return err
			}
			if !proprie.Valid && principala.Valid && principala.String == a.FunctieID {
				a.FunctieID = ""
				a.Functie = ""
			}
		}
		if a.FunctieID != "" {
			lista, err := citesteNomenclator(ctx, tx, "functii")
			if err != nil {
				return err
			}
			gasit := false
			for _, f := range lista {
				if f.ID == a.FunctieID {
					gasit = true
					var vechiID string
					_ = tx.QueryRowContext(ctx, "SELECT functie_id FROM angajati_functii WHERE angajat_id=?", a.ID).Scan(&vechiID)
					if !f.Activ && vechiID != f.ID {
						return eroare("nomenclator.inactiv")
					}
					a.Functie = numeElement(f)
				}
			}
			if !gasit {
				return eroare("nomenclator.invalid")
			}
		}
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
		if _, err := tx.ExecContext(ctx, "DELETE FROM angajati_functii WHERE angajat_id=?", a.ID); err != nil {
			return err
		}
		if a.FunctieID != "" {
			if _, err := tx.ExecContext(ctx, "INSERT INTO angajati_functii VALUES(?,?)", a.ID, a.FunctieID); err != nil {
				return err
			}
		}
		contractID := a.ContractID
		if contractID == "" {
			contractID = a.TipContract
		}
		var contractVechi string
		errContract := tx.QueryRowContext(ctx, "SELECT contract_id FROM angajati_contracte WHERE angajat_id=?", a.ID).Scan(&contractVechi)
		if errContract != nil && errContract != sql.ErrNoRows {
			return errContract
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO angajati_contracte VALUES(?,?) ON CONFLICT(angajat_id) DO UPDATE SET contract_id=excluded.contract_id", a.ID, contractID); err != nil {
			return err
		}
		if contractVechi != contractID && a.ContractID != "" {
			if err := jurnal(ctx, q, "angajati_contracte", a.ID, contractVechi, contractID); err != nil {
				return err
			}
		}
		var principalVechi string
		errLegatura := tx.QueryRowContext(ctx, "SELECT principal_id FROM angajati_garzi WHERE angajat_id=?", a.ID).Scan(&principalVechi)
		if errLegatura != nil && errLegatura != sql.ErrNoRows {
			return errLegatura
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM angajati_garzi WHERE angajat_id=?", a.ID); err != nil {
			return err
		}
		if a.AngajatPrincipalID != "" {
			if _, err := tx.ExecContext(ctx, "INSERT INTO angajati_garzi VALUES(?,?)", a.ID, a.AngajatPrincipalID); err != nil {
				return err
			}
		}
		if principalVechi != a.AngajatPrincipalID {
			if err := jurnal(ctx, q, "angajati_garzi", a.ID, principalVechi, a.AngajatPrincipalID); err != nil {
				return err
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
	return q.AdaugaJurnal(ctx, acces.AdaugaJurnalParams{Tabel: tabel, RandID: id, ValoareVeche: inainte, ValoareNoua: &v, Cand: moment(), Cine: Actor(ctx)})
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
