package datebaza

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd/acces"
)

// StareLuna reprezintă starea de deschidere / blocare a unei luni calendaristice (F12-T01).
// Amprenta se calculează la închidere (D044 punctul 6); AmprentaCorespunde spune dacă
// datele curente ale lunii mai dau aceeași amprentă.
type StareLuna struct {
	An                 int    `json:"an"`
	Luna               int    `json:"luna"`
	Stare              string `json:"stare"` // "DESCHISA" | "BLOCATA" | "INCHISA"
	ModificatLa        string `json:"modificat_la"`
	ModificatDe        string `json:"modificat_de"`
	Motiv              string `json:"motiv"`
	Amprenta           string `json:"amprenta"`
	AmprentaCorespunde bool   `json:"amprenta_corespunde"`
}

// CitesteStareLuna returnează starea curentă a lunii. Dacă nu este setată explicit, returnează DESCHISA.
func (s *Serviciu) CitesteStareLuna(ctx context.Context, an, luna int) (StareLuna, error) {
	rez := StareLuna{An: an, Luna: luna, Stare: "DESCHISA"}
	if an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return rez, eroare("pontaj.perioada_invalida")
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return rez, traduceEroare(err)
	}
	defer tx.Rollback()
	err = tx.QueryRowContext(ctx,
		`SELECT stare, modificat_la, modificat_de, motiv, amprenta FROM stare_luni WHERE an = ? AND luna = ?`,
		an, luna,
	).Scan(&rez.Stare, &rez.ModificatLa, &rez.ModificatDe, &rez.Motiv, &rez.Amprenta)
	if err == sql.ErrNoRows {
		return rez, nil
	}
	if err != nil {
		return rez, traduceEroare(err)
	}
	if rez.Amprenta != "" {
		curenta, err := audit.AmprentaLuna(ctx, tx, an, luna)
		if err != nil {
			return rez, traduceEroare(err)
		}
		rez.AmprentaCorespunde = curenta == rez.Amprenta
	}
	return rez, nil
}

// SeteazaStareLuna modifică starea lunii (DESCHISA, BLOCATA, INCHISA) și înregistrează acțiunea în jurnal.
// La închidere calculează amprenta lunii; la redeschidere o șterge (vechea rămâne în jurnal).
func (s *Serviciu) SeteazaStareLuna(ctx context.Context, an, luna int, stare, motiv, utilizator string) (StareLuna, error) {
	stare = strings.ToUpper(strings.TrimSpace(stare))
	if stare != "DESCHISA" && stare != "BLOCATA" && stare != "INCHISA" {
		return StareLuna{}, eroare("pontaj.stare_invalida")
	}
	if an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return StareLuna{}, eroare("pontaj.perioada_invalida")
	}
	// Motivul îl scrie utilizatorul și intră în jurnal (F12-T05).
	if motiv = strings.TrimSpace(motiv); motiv == "" {
		return StareLuna{}, eroare("pontaj.motiv_obligatoriu")
	}
	if utilizator == "" {
		utilizator = "Utilizator"
	}
	rez := StareLuna{An: an, Luna: luna, Stare: stare, ModificatLa: time.Now().UTC().Format(time.RFC3339), ModificatDe: utilizator, Motiv: motiv}

	err := s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		var vechi *StareLuna
		var v StareLuna
		err := tx.QueryRowContext(ctx, `SELECT stare, modificat_la, modificat_de, motiv, amprenta FROM stare_luni WHERE an = ? AND luna = ?`, an, luna).
			Scan(&v.Stare, &v.ModificatLa, &v.ModificatDe, &v.Motiv, &v.Amprenta)
		switch {
		case err == nil:
			v.An, v.Luna = an, luna
			vechi = &v
		case err != sql.ErrNoRows:
			return err
		}
		if stare == "INCHISA" {
			if rez.Amprenta, err = audit.AmprentaLuna(ctx, tx, an, luna); err != nil {
				return err
			}
			rez.AmprentaCorespunde = true
		} else if stare == "BLOCATA" && vechi != nil && vechi.Stare == "INCHISA" {
			// Blocarea unei luni închise nu schimbă datele, deci păstrează amprenta.
			rez.Amprenta, rez.AmprentaCorespunde = vechi.Amprenta, true
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO stare_luni (an, luna, stare, modificat_la, modificat_de, motiv, amprenta)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(an, luna) DO UPDATE SET
				stare = excluded.stare,
				modificat_la = excluded.modificat_la,
				modificat_de = excluded.modificat_de,
				motiv = excluded.motiv,
				amprenta = excluded.amprenta
		`, an, luna, stare, rez.ModificatLa, utilizator, motiv, rez.Amprenta); err != nil {
			return err
		}
		var j any
		if vechi != nil {
			j = vechi
		}
		return jurnal(ctx, q, "stare_luni", perioada(an, luna), j, rez)
	})
	if err != nil {
		return StareLuna{}, traduceEroare(err)
	}
	return rez, nil
}

func perioada(an, luna int) string {
	return time.Date(an, time.Month(luna), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
}

// EsteLunaInchisa verifică dacă luna este blocată sau închisă la modificări.
func (s *Serviciu) EsteLunaInchisa(ctx context.Context, an, luna int) (bool, error) {
	st, err := s.CitesteStareLuna(ctx, an, luna)
	if err != nil {
		return false, err
	}
	return st.Stare == "BLOCATA" || st.Stare == "INCHISA", nil
}
