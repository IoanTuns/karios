package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	if c.SchimbareOra == "" {
		c.SchimbareOra = calcul.SchimbareOraNominala
	}
	optiuni, err := optiuniReguli(ctx, q)
	if err != nil {
		return c, err
	}
	adapteazaConfiguratia(&c, optiuni)
	return c, calcul.ValideazaConfiguratiePontaj(c)
}

// adapteazaConfiguratia aduce o configuratie salvata la coloanele si codurile curente (D092):
// codurile redenumite (CEV->CVE etc.) se traduc daca vechiul cod nu mai exista, iar
// coloanele noi lipsa primesc regula implicita, doar cu codurile de absenta existente.
// O regula pe zile de absenta ramasa fara coduri devine neconfigurata (motiv vizibil).
func adapteazaConfiguratia(c *calcul.ConfiguratiePontaj, optiuni []OptiuneRegula) {
	exista := map[string]bool{}
	for _, o := range optiuni {
		if o.Tip == "ABSENTA" {
			exista[o.Cod] = true
		}
	}
	traduce := func(r *calcul.RegulaColoana) {
		if tipRegula(r.Sursa) != "ABSENTA" {
			return
		}
		for i, cod := range r.Coduri {
			if nou, ok := calcul.CoduriRedenumite[cod]; ok && !exista[cod] && exista[nou] {
				r.Coduri[i] = nou
			}
		}
	}
	for i := range c.Reguli {
		traduce(&c.Reguli[i])
	}
	for i := range c.Conditionate {
		traduce(&c.Conditionate[i].Regula)
	}
	completeazaTipuriNoi(c, exista)
	prezente := map[string]bool{}
	for _, r := range c.Reguli {
		prezente[r.Coloana] = true
	}
	for _, r := range calcul.ConfiguratiePontajInitiala().Reguli {
		if prezente[r.Coloana] {
			continue
		}
		if tipRegula(r.Sursa) == "ABSENTA" {
			coduri := []string{}
			for _, cod := range r.Coduri {
				if exista[cod] {
					coduri = append(coduri, cod)
				}
			}
			r.Coduri = coduri
			if r.Sursa == "ZILE_ABSENTA" && len(coduri) == 0 {
				r.Sursa = "NEDEFINITA"
			}
		}
		c.Reguli = append(c.Reguli, r)
	}
}

// tipRegula spune ce fel de coduri foloseste sursa: tipuri de absenta sau programe.
func tipRegula(sursa string) string {
	switch sursa {
	case "ZILE_ABSENTA", "ZILE_ABSENTA_ALTELE", "ORE_LUCRATE_CODURI":
		return "ABSENTA"
	}
	return "SCHIMB"
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
	reguli := append([]calcul.RegulaColoana(nil), c.Reguli...)
	for _, r := range c.Conditionate {
		reguli = append(reguli, r.Regula)
	}
	for _, r := range reguli {
		tip := tipRegula(r.Sursa)
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
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		for _, r := range c.Conditionate {
			for _, ref := range []struct{ tabel, id string }{{"grupuri_profesionale", r.GrupID}, {"functii", r.FunctieID}, {"tipuri_contract", r.ContractID}} {
				if ref.id != "" {
					var n int
					if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM "+ref.tabel+" WHERE id=?", ref.id).Scan(&n); err != nil {
						return err
					}
					if n != 1 {
						return eroare("nomenclator.invalid")
					}
				}
			}
		}
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
		if c.SchimbareOra == "" {
			c.SchimbareOra = calcul.SchimbareOraNominala
		}
		if c.SchimbareOra != vechi.SchimbareOra {
			if err := verificaSchimbareOraModificabila(ctx, tx); err != nil {
				return err
			}
		}
		date, err := json.Marshal(struct {
			Atribuire    string                      `json:"atribuire"`
			SchimbareOra string                      `json:"schimbare_ora"`
			Reguli       []calcul.RegulaColoana      `json:"reguli"`
			Conditionate []calcul.RegulaConditionata `json:"conditionate,omitempty"`
		}{c.Atribuire, c.SchimbareOra, c.Reguli, c.Conditionate})
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

// verificaSchimbareOraModificabila (F04-T08, D102): optiunea de schimbare a orei rescrie orele
// din noptile de trecere la ora de vara/iarna. Se refuza daca o luna care contine o astfel de
// noapte e inchisa; luna se redeschide intai, cu motiv, ca la istoricul contractului (D101).
func verificaSchimbareOraModificabila(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(ctx, `SELECT an,luna FROM stare_luni WHERE stare<>'DESCHISA' ORDER BY an,luna`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var an, luna int
		if err := rows.Scan(&an, &luna); err != nil {
			return err
		}
		if calcul.LunaCuSchimbareOra(an, luna) {
			return &calcul.Eroare{Cod: "pontaj.schimbare_ora_luna_inchisa", Param: map[string]string{"luna": fmt.Sprintf("%02d.%04d", luna, an)}}
		}
	}
	return rows.Err()
}

// completeazaTipuriNoi (F08-T15): o configuratie salvata inainte sa existe un tip de absenta
// din D092 (de ex. FMP) ramanea fara el: CVE il numara ca „alta absenta”, coloana lui era
// neconfigurata, iar orele lui lipseau din OL. La citire, regulile D092 primesc tipurile
// aparute intre timp: coloana proprie neconfigurata si fara coduri primeste codul ei, CVE
// exclude tipurile cu coloana proprie, iar OL include creditele FMP.
func completeazaTipuriNoi(c *calcul.ConfiguratiePontaj, exista map[string]bool) {
	adauga := func(r *calcul.RegulaColoana, cod string) {
		if !exista[cod] {
			return
		}
		for _, x := range r.Coduri {
			if x == cod {
				return
			}
		}
		r.Coduri = append(r.Coduri, cod)
	}
	proprii := map[string]bool{}
	for _, cod := range calcul.CoduriExcluseCVE {
		proprii[cod] = true
	}
	for i := range c.Reguli {
		r := &c.Reguli[i]
		if proprii[r.Coloana] && r.Sursa == "NEDEFINITA" && len(r.Coduri) == 0 && exista[r.Coloana] {
			r.Sursa, r.Coduri = "ZILE_ABSENTA", []string{r.Coloana}
		}
	}
	for i := range c.Reguli {
		r := &c.Reguli[i]
		switch {
		case r.Coloana == "CVE" && r.Sursa == "ZILE_ABSENTA_ALTELE":
			for _, cod := range calcul.CoduriExcluseCVE {
				adauga(r, cod)
			}
			for _, alta := range c.Reguli {
				if alta.Sursa == "ZILE_ABSENTA" {
					for _, cod := range alta.Coduri {
						adauga(r, cod)
					}
				}
			}
		case r.Coloana == "OL" && r.Sursa == "ORE_LUCRATE_CODURI":
			adauga(r, "FMP")
		}
	}
}
