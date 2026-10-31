package acces

import "context"

type ProfilPontajRow struct {
	ID, GrupID, FunctieID, ContractID string
	Conducere, Garda                  int64
}

func (q *Queries) ListaProfilePontaj(ctx context.Context) ([]ProfilPontajRow, error) {
	rows, err := q.db.QueryContext(ctx, `SELECT a.id,coalesce(CASE WHEN af.functie_id IS NULL THEN fp.grup_id ELSE f.grup_id END,''),coalesce(af.functie_id,ap.functie_id,''),coalesce(ac.contract_id,a.tip_contract),coalesce(f.conducere,fp.conducere,0),CASE WHEN ag.principal_id IS NOT NULL OR coalesce(tc.garda,0)=1 THEN 1 ELSE 0 END FROM angajati a LEFT JOIN angajati_functii af ON af.angajat_id=a.id LEFT JOIN functii f ON f.id=af.functie_id LEFT JOIN angajati_contracte ac ON ac.angajat_id=a.id LEFT JOIN tipuri_contract tc ON tc.id=ac.contract_id LEFT JOIN angajati_garzi ag ON ag.angajat_id=a.id LEFT JOIN angajati_functii ap ON ap.angajat_id=ag.principal_id LEFT JOIN functii fp ON fp.id=ap.functie_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := []ProfilPontajRow{}
	for rows.Next() {
		var r ProfilPontajRow
		if err := rows.Scan(&r.ID, &r.GrupID, &r.FunctieID, &r.ContractID, &r.Conducere, &r.Garda); err != nil {
			return nil, err
		}
		rez = append(rez, r)
	}
	return rez, rows.Err()
}

func (q *Queries) CategorieContractAngajat(ctx context.Context, id string) (string, error) {
	var categorie string
	err := q.db.QueryRowContext(ctx, `SELECT coalesce(tc.categorie_program_id,a.categorie_program_id,'') FROM angajati a LEFT JOIN angajati_contracte ac ON ac.angajat_id=a.id LEFT JOIN tipuri_contract tc ON tc.id=ac.contract_id WHERE a.id=?`, id).Scan(&categorie)
	return categorie, err
}
