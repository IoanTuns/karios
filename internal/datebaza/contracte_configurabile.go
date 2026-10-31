package datebaza

import (
	"context"
	"database/sql"
)

func citesteContracteAngajati(ctx context.Context, tx *sql.Tx) (map[string]string, error) {
	rows, err := tx.QueryContext(ctx, "SELECT angajat_id,contract_id FROM angajati_contracte")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]string{}
	for rows.Next() {
		var a, c string
		if err := rows.Scan(&a, &c); err != nil {
			return nil, err
		}
		rez[a] = c
	}
	return rez, rows.Err()
}
func contractPentruSalvare(ctx context.Context, tx *sql.Tx, a Angajat) (ElementNomenclator, error) {
	lista, err := citesteNomenclator(ctx, tx, "contracte")
	if err != nil {
		return ElementNomenclator{}, traduceEroare(err)
	}
	for _, c := range lista {
		if c.ID == a.ContractID {
			return c, nil
		}
	}
	return ElementNomenclator{}, eroare("date.contract_invalid")
}
func valideazaContractActiv(ctx context.Context, tx *sql.Tx, a Angajat) error {
	var activ int
	err := tx.QueryRowContext(ctx, "SELECT activ FROM tipuri_contract WHERE id=?", a.ContractID).Scan(&activ)
	if err != nil {
		return eroare("date.contract_invalid")
	}
	if activ == 1 {
		return nil
	}
	var vechi string
	err = tx.QueryRowContext(ctx, "SELECT contract_id FROM angajati_contracte WHERE angajat_id=?", a.ID).Scan(&vechi)
	if err == nil && vechi == a.ContractID {
		return nil
	}
	return eroare("nomenclator.inactiv")
}
