package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/i18n"
)

type OptiuneRegula struct {
	Cod      string    `json:"cod"`
	Denumire i18n.Text `json:"denumire"`
	Tip      string    `json:"tip"`
}

func citesteConfiguratiePontaj(ctx context.Context, q *acces.Queries) (calcul.ConfiguratiePontaj, error) {
	r, err := q.CitesteConfiguratiePontaj(ctx)
	if err != nil {
		return calcul.ConfiguratiePontaj{}, err
	}
	var c calcul.ConfiguratiePontaj
	if err := json.Unmarshal([]byte(r.Configuratie), &c); err != nil {
		return c, err
	}
	c.Revizie = strconv.FormatInt(r.Revizie, 10)
	return c, calcul.ValideazaConfiguratiePontaj(c)
}
func CitesteConfiguratiePontaj(ctx context.Context, q *acces.Queries) (calcul.ConfiguratiePontaj, error) {
	return citesteConfiguratiePontaj(ctx, q)
}
func (s *Serviciu) ConfiguratiePontaj(ctx context.Context) (calcul.ConfiguratiePontaj, []OptiuneRegula, error) {
	c, err := citesteConfiguratiePontaj(ctx, acces.New(s.db))
	if err != nil {
		return c, nil, traduceEroare(err)
	}
	optiuni, err := optiuniReguli(ctx, acces.New(s.db))
	return c, optiuni, traduceEroare(err)
}
func optiuniReguli(ctx context.Context, q *acces.Queries) ([]OptiuneRegula, error) {
	rez := []OptiuneRegula{}
	schimburi, err := q.ListaSchimburi(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range schimburi {
		var d i18n.Text
		if err := json.Unmarshal([]byte(s.Denumire), &d); err != nil {
			return nil, err
		}
		rez = append(rez, OptiuneRegula{Cod: s.Cod, Denumire: d, Tip: "SCHIMB"})
	}
	absente, err := q.ListaTipuriAbsenta(ctx)
	if err != nil {
		return nil, err
	}
	for _, a := range absente {
		var d i18n.Text
		if err := json.Unmarshal([]byte(a.Denumire), &d); err != nil {
			return nil, err
		}
		rez = append(rez, OptiuneRegula{Cod: a.Cod, Denumire: d, Tip: "ABSENTA"})
	}
	return rez, nil
}
func ValideazaCoduriReguli(ctx context.Context, q *acces.Queries, c calcul.ConfiguratiePontaj) error {
	optiuni, err := optiuniReguli(ctx, q)
	if err != nil {
		return err
	}
	for _, r := range c.Reguli {
		tip := "SCHIMB"
		if r.Sursa == "ZILE_ABSENTA" {
			tip = "ABSENTA"
		}
		for _, cod := range r.Coduri {
			gasit := false
			for _, o := range optiuni {
				gasit = gasit || o.Cod == cod && o.Tip == tip
			}
			if !gasit {
				return eroare("pontaj.reguli_cod_lipsa")
			}
		}
	}
	return nil
}
func (s *Serviciu) SalveazaConfiguratiePontaj(ctx context.Context, c calcul.ConfiguratiePontaj) error {
	if err := calcul.ValideazaConfiguratiePontaj(c); err != nil {
		return err
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		vechi, err := citesteConfiguratiePontaj(ctx, q)
		if err != nil {
			return err
		}
		r, err := strconv.ParseInt(c.Revizie, 10, 64)
		if err != nil || c.Revizie != vechi.Revizie {
			return eroare("pontaj.reguli_conflict")
		}
		if err := ValideazaCoduriReguli(ctx, q, c); err != nil {
			return err
		}
		date, err := json.Marshal(struct {
			Atribuire string                 `json:"atribuire"`
			Reguli    []calcul.RegulaColoana `json:"reguli"`
		}{c.Atribuire, c.Reguli})
		if err != nil {
			return err
		}
		noua, err := q.SalveazaConfiguratiePontaj(ctx, acces.SalveazaConfiguratiePontajParams{Configuratie: string(date), ModificatLa: moment(), Revizie: r})
		if errors.Is(err, sql.ErrNoRows) {
			return eroare("pontaj.reguli_conflict")
		}
		if err != nil {
			return err
		}
		c.Revizie = strconv.FormatInt(noua, 10)
		return jurnal(ctx, q, "configuratie_pontaj", "1", vechi, c)
	})
}
