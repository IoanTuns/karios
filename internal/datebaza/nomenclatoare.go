package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
	"strings"
)

type ElementNomenclator struct {
	ID       string    `json:"id"`
	Denumire i18n.Text `json:"denumire"`
	Activ    bool      `json:"activ"`
	Revizie  string    `json:"revizie"`
}

func tabelaNomenclator(tip string) string {
	switch tip {
	case "categorii":
		return "categorii_schimburi"
	case "functii":
		return "functii"
	}
	return ""
}
func numeElement(e ElementNomenclator) string {
	if e.Denumire["ro"] != "" {
		return e.Denumire["ro"]
	}
	return e.Denumire["en"]
}
func citesteNomenclator(ctx context.Context, tx *sql.Tx, tip string) ([]ElementNomenclator, error) {
	tabela := tabelaNomenclator(tip)
	if tabela == "" {
		return nil, eroare("nomenclator.invalid")
	}
	rows, err := tx.QueryContext(ctx, "SELECT id,denumire,activ,modificat_la FROM "+tabela+" ORDER BY denumire,id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lista := []ElementNomenclator{}
	for rows.Next() {
		var e ElementNomenclator
		var text string
		var activ int
		if err := rows.Scan(&e.ID, &text, &activ, &e.Revizie); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(text), &e.Denumire); err != nil {
			return nil, err
		}
		e.Activ = activ == 1
		lista = append(lista, e)
	}
	return lista, rows.Err()
}
func (s *Serviciu) SalveazaElementNomenclator(ctx context.Context, tip string, e ElementNomenclator) error {
	tabela := tabelaNomenclator(tip)
	if tabela == "" {
		return eroare("nomenclator.invalid")
	}
	e.Denumire = i18n.Text{"ro": strings.TrimSpace(e.Denumire["ro"]), "en": strings.TrimSpace(e.Denumire["en"])}
	if numeElement(e) == "" {
		return eroare("sarbatori.denumire_obligatorie")
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		lista, err := citesteNomenclator(ctx, tx, tip)
		if err != nil {
			return err
		}
		var vechi *ElementNomenclator
		for i := range lista {
			if lista[i].ID == e.ID {
				copie := lista[i]
				vechi = &copie
			}
		}
		for _, element := range lista {
			if element.ID != e.ID && strings.EqualFold(numeElement(element), numeElement(e)) && (vechi == nil || !strings.EqualFold(numeElement(*vechi), numeElement(e))) {
				return eroare("nomenclator.duplicat")
			}
		}
		text, err := json.Marshal(e.Denumire)
		if err != nil {
			return err
		}
		acum := moment()
		if e.ID == "" {
			e.ID = idNou()
			if tip == "categorii" {
				e.ID = strings.ToUpper(e.ID)
			}
			_, err = tx.ExecContext(ctx, "INSERT INTO "+tabela+"(id,denumire,activ,modificat_la) VALUES(?,?,?,?)", e.ID, string(text), numarBool(e.Activ), acum)
		} else {
			if vechi == nil {
				return eroare("nomenclator.invalid")
			}
			if vechi.Revizie != e.Revizie {
				return eroare("nomenclator.conflict")
			}
			var r sql.Result
			r, err = tx.ExecContext(ctx, "UPDATE "+tabela+" SET denumire=?,activ=?,modificat_la=? WHERE id=? AND modificat_la=?", string(text), numarBool(e.Activ), acum, e.ID, e.Revizie)
			if err == nil {
				n, _ := r.RowsAffected()
				if n != 1 {
					return eroare("nomenclator.conflict")
				}
			}
		}
		if err != nil {
			return err
		}
		if tip == "functii" && vechi != nil && numeElement(*vechi) != numeElement(e) {
			rows, err := tx.QueryContext(ctx, "SELECT angajat_id FROM angajati_functii WHERE functie_id=?", e.ID)
			if err != nil {
				return err
			}
			ids := []string{}
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err != nil {
					rows.Close()
					return err
				}
				ids = append(ids, id)
			}
			err = rows.Err()
			rows.Close()
			if err != nil {
				return err
			}
			for _, id := range ids {
				inainte, err := q.CitesteAngajat(ctx, id)
				if err != nil {
					return err
				}
				if _, err := tx.ExecContext(ctx, "UPDATE angajati SET functie=?,modificat_la=? WHERE id=?", numeElement(e), acum, id); err != nil {
					return err
				}
				if _, err := tx.ExecContext(ctx, "INSERT INTO angajati_functii VALUES(?,?) ON CONFLICT(angajat_id) DO UPDATE SET functie_id=excluded.functie_id", id, e.ID); err != nil {
					return err
				}
				dupa, err := q.CitesteAngajat(ctx, id)
				if err != nil {
					return err
				}
				if err := jurnal(ctx, q, "angajati", id, inainte, dupa); err != nil {
					return err
				}
			}
		}
		e.Revizie = acum
		return jurnal(ctx, q, tabela, e.ID, vechi, e)
	})
}
