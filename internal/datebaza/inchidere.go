package datebaza

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"pontix/internal/bd/acces"
)

// StareLuna reprezintă starea de deschidere / blocare a unei luni calendaristice (F12-T01).
type StareLuna struct {
	An          int    `json:"an"`
	Luna        int    `json:"luna"`
	Stare       string `json:"stare"` // "DESCHISA" | "BLOCATA" | "INCHISA"
	ModificatLa string `json:"modificat_la"`
	ModificatDe string `json:"modificat_de"`
	Motiv       string `json:"motiv"`
}

// CitesteStareLuna returnează starea curentă a lunii. Dacă nu este setată explicit, returnează DESCHISA.
func (s *Serviciu) CitesteStareLuna(ctx context.Context, an, luna int) (StareLuna, error) {
	rez := StareLuna{
		An:          an,
		Luna:        luna,
		Stare:       "DESCHISA",
		ModificatLa: "",
		ModificatDe: "",
		Motiv:       "",
	}
	if an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return rez, eroare("pontaj.perioada_invalida")
	}

	var stare, modLa, modDe, motiv string
	err := s.db.QueryRowContext(ctx,
		`SELECT stare, modificat_la, modificat_de, motiv FROM stare_luni WHERE an = ? AND luna = ?`,
		an, luna,
	).Scan(&stare, &modLa, &modDe, &motiv)
	if err == nil {
		rez.Stare = stare
		rez.ModificatLa = modLa
		rez.ModificatDe = modDe
		rez.Motiv = motiv
		return rez, nil
	}
	if err == sql.ErrNoRows {
		return rez, nil
	}
	return rez, traduceEroare(err)
}

// SeteazaStareLuna modifică starea lunii (DESCHISA, BLOCATA, INCHISA) și înregistrează acțiunea în audit.
func (s *Serviciu) SeteazaStareLuna(ctx context.Context, an, luna int, stare, motiv, utilizator string) (StareLuna, error) {
	stare = strings.ToUpper(strings.TrimSpace(stare))
	if stare != "DESCHISA" && stare != "BLOCATA" && stare != "INCHISA" {
		return StareLuna{}, eroare("pontaj.stare_invalida")
	}
	if an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return StareLuna{}, eroare("pontaj.perioada_invalida")
	}
	if utilizator == "" {
		utilizator = "Utilizator"
	}
	acum := time.Now().UTC().Format(time.RFC3339)

	var rez StareLuna
	err := s.modifica(ctx, func(q *acces.Queries) error {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO stare_luni (an, luna, stare, modificat_la, modificat_de, motiv)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(an, luna) DO UPDATE SET
				stare = excluded.stare,
				modificat_la = excluded.modificat_la,
				modificat_de = excluded.modificat_de,
				motiv = excluded.motiv
		`, an, luna, stare, acum, utilizator, motiv)
		if err != nil {
			return err
		}

		rez = StareLuna{
			An:          an,
			Luna:        luna,
			Stare:       stare,
			ModificatLa: acum,
			ModificatDe: utilizator,
			Motiv:       motiv,
		}
		return nil
	})
	if err != nil {
		return StareLuna{}, traduceEroare(err)
	}
	return rez, nil
}

// EsteLunaInchisa verifică dacă luna este blocată sau închisă la modificări.
func (s *Serviciu) EsteLunaInchisa(ctx context.Context, an, luna int) (bool, error) {
	st, err := s.CitesteStareLuna(ctx, an, luna)
	if err != nil {
		return false, err
	}
	return st.Stare == "BLOCATA" || st.Stare == "INCHISA", nil
}
