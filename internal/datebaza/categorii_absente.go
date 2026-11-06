package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
	"strconv"
	"strings"
)

type CategorieAbsenta struct {
	Cod      string    `json:"cod"`
	Denumire i18n.Text `json:"denumire"`
	Natura   string    `json:"natura"`
	EfectOre string    `json:"efect_ore"`
	Activ    bool      `json:"activ"`
	Revizie  string    `json:"revizie"`
}

func citesteCategoriiAbsente(ctx context.Context, tx *sql.Tx) ([]CategorieAbsenta, error) {
	rows, e := tx.QueryContext(ctx, "SELECT cod,denumire,natura,efect_ore,activ,revizie FROM categorii_absente ORDER BY cod")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	rez := []CategorieAbsenta{}
	for rows.Next() {
		var c CategorieAbsenta
		var text string
		if e := rows.Scan(&c.Cod, &text, &c.Natura, &c.EfectOre, &c.Activ, &c.Revizie); e != nil {
			return nil, e
		}
		if e := json.Unmarshal([]byte(text), &c.Denumire); e != nil {
			return nil, e
		}
		rez = append(rez, c)
	}
	return rez, rows.Err()
}
func (s *Serviciu) SalveazaCategorieAbsenta(ctx context.Context, c CategorieAbsenta) error {
	c.Cod = strings.ToUpper(strings.TrimSpace(c.Cod))
	if c.Cod == "" || len(c.Cod) > 48 || strings.IndexFunc(c.Cod, func(r rune) bool { return !(r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-') }) >= 0 {
		return eroare("absente.cod_categorie_invalid")
	}
	if c.Natura != "ODIHNA" && c.Natura != "MEDICAL" && c.Natura != "SPECIAL" {
		return eroare("absente.categorie_invalida")
	}
	if c.EfectOre != "ANTERIOR" && c.EfectOre != "FARA_ORE" && c.EfectOre != "NORMA" && c.EfectOre != "PROGRAM" {
		return eroare("absente.efect_invalid")
	}
	if c.Natura != "SPECIAL" && (c.EfectOre == "NORMA" || c.EfectOre == "PROGRAM") {
		return eroare("absente.efect_invalid")
	}
	den := i18n.Text{}
	for l, t := range c.Denumire {
		if t = strings.TrimSpace(t); t != "" {
			den[l] = t
		}
	}
	if len(den) == 0 {
		return eroare("sarbatori.denumire_obligatorie")
	}
	c.Denumire = den
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		lista, e := citesteCategoriiAbsente(ctx, tx)
		if e != nil {
			return e
		}
		var vechi *CategorieAbsenta
		for _, v := range lista {
			if v.Cod == c.Cod {
				copia := v
				vechi = &copia
			}
		}
		if vechi == nil && c.Revizie != "" || vechi != nil && c.Revizie != vechi.Revizie {
			return eroare("absente.conflict")
		}
		if vechi != nil && (c.Natura != vechi.Natura || c.EfectOre != vechi.EfectOre) {
			var folosit bool
			e := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM coduri_absenta t WHERE t.categorie_id=? AND (EXISTS(SELECT 1 FROM absente_zile z WHERE z.cod=t.cod) OR EXISTS(SELECT 1 FROM pontaj p WHERE p.cod_absenta=t.cod)))`, c.Cod).Scan(&folosit)
			if e != nil {
				return e
			}
			if folosit {
				return eroare("absente.categorie_folosita")
			}
		}
		text, e := json.Marshal(c.Denumire)
		if e != nil {
			return e
		}
		if vechi == nil {
			if c.EfectOre == "ANTERIOR" {
				return eroare("absente.efect_invalid")
			}
			_, e = tx.ExecContext(ctx, "INSERT INTO categorii_absente(cod,denumire,natura,efect_ore,activ) VALUES(?,?,?,?,?)", c.Cod, string(text), c.Natura, c.EfectOre, c.Activ)
			c.Revizie = "1"
		} else {
			_, e = tx.ExecContext(ctx, "UPDATE categorii_absente SET denumire=?,natura=?,efect_ore=?,activ=?,revizie=revizie+1 WHERE cod=?", string(text), c.Natura, c.EfectOre, c.Activ, c.Cod)
			n, _ := strconv.Atoi(vechi.Revizie)
			c.Revizie = strconv.Itoa(n + 1)
		}
		if e != nil {
			return e
		}
		// Propagarea se face numai pentru reguli fara utilizari istorice; denumirea nu schimba calculul.
		_, e = tx.ExecContext(ctx, "UPDATE coduri_absenta SET categorie=?,efect_ore=?,modificat_la=? WHERE categorie_id=?", c.Natura, c.EfectOre, moment(), c.Cod)
		if e != nil {
			return e
		}
		return jurnal(ctx, q, "categorii_absente", c.Cod, vechi, c)
	})
}
