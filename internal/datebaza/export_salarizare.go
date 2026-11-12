package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"pontix/internal/bd/acces"
	"pontix/internal/exportsalarizare"
)

// F10-T03 (D108): modelele exportului pentru salarizare, in baza criptata, cu revizie si jurnal.

// ListaModeleExport intoarce modelele, ordonate dupa denumire.
func (s *Serviciu) ListaModeleExport(ctx context.Context) ([]exportsalarizare.Model, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, denumire, configuratie, activ, revizie FROM modele_export_salarizare ORDER BY denumire, id`)
	if err != nil {
		return nil, traduceEroare(err)
	}
	defer rows.Close()
	rez := []exportsalarizare.Model{}
	for rows.Next() {
		var m exportsalarizare.Model
		var config string
		var activ, revizie int64
		if err := rows.Scan(&m.ID, &m.Denumire, &config, &activ, &revizie); err != nil {
			return nil, err
		}
		denumire, id := m.Denumire, m.ID
		if err := json.Unmarshal([]byte(config), &m); err != nil {
			return nil, err
		}
		m.ID, m.Denumire, m.Activ, m.Revizie = id, denumire, activ == 1, strconv.FormatInt(revizie, 10)
		rez = append(rez, m)
	}
	return rez, rows.Err()
}

// ModelExport intoarce un model dupa id.
func (s *Serviciu) ModelExport(ctx context.Context, id string) (exportsalarizare.Model, error) {
	lista, err := s.ListaModeleExport(ctx)
	if err != nil {
		return exportsalarizare.Model{}, err
	}
	for _, m := range lista {
		if m.ID == id {
			return m, nil
		}
	}
	return exportsalarizare.Model{}, eroare("export_salarizare.model_inexistent")
}

// SalveazaModelExport adauga (ID gol) sau modifica un model; Revizie trebuie sa fie cea citita.
func (s *Serviciu) SalveazaModelExport(ctx context.Context, m exportsalarizare.Model) (exportsalarizare.Model, error) {
	if err := exportsalarizare.Valideaza(&m); err != nil {
		return m, err
	}
	config := m
	config.ID, config.Denumire, config.Activ, config.Revizie = "", "", false, ""
	text, err := json.Marshal(config)
	if err != nil {
		return m, err
	}
	err = s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		acum := moment()
		if m.ID == "" {
			m.ID = idNou()
			if _, err := tx.ExecContext(ctx, `INSERT INTO modele_export_salarizare (id, denumire, configuratie, activ, creat_la, modificat_la) VALUES (?, ?, ?, ?, ?, ?)`, m.ID, m.Denumire, string(text), numarBool(m.Activ), acum, acum); err != nil {
				return err
			}
			m.Revizie = "1"
			return jurnal(ctx, q, "modele_export_salarizare", m.ID, nil, m)
		}
		var vechi string
		var revizie int64
		err := tx.QueryRowContext(ctx, `SELECT configuratie, revizie FROM modele_export_salarizare WHERE id=?`, m.ID).Scan(&vechi, &revizie)
		if errors.Is(err, sql.ErrNoRows) {
			return eroare("export_salarizare.model_inexistent")
		}
		if err != nil {
			return err
		}
		if strconv.FormatInt(revizie, 10) != m.Revizie {
			return eroare("export_salarizare.conflict")
		}
		if _, err := tx.ExecContext(ctx, `UPDATE modele_export_salarizare SET denumire=?, configuratie=?, activ=?, revizie=revizie+1, modificat_la=? WHERE id=?`, m.Denumire, string(text), numarBool(m.Activ), acum, m.ID); err != nil {
			return err
		}
		m.Revizie = strconv.FormatInt(revizie+1, 10)
		return jurnal(ctx, q, "modele_export_salarizare", m.ID, json.RawMessage(vechi), m)
	})
	return m, err
}
