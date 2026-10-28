package datebaza

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

// ListaReguliConfigurabile citește toate regulile din baza de date.
func (s *Serviciu) ListaReguliConfigurabile(ctx context.Context) ([]calcul.RegulaConfigurabila, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, cod, denumire, descriere, tip_regula, severitate, coduri_schimb, valoare_numerica, activ
		FROM reguli_configurabile
		ORDER BY cod ASC
	`)
	if err != nil {
		return nil, traduceEroare(err)
	}
	defer rows.Close()

	var rez []calcul.RegulaConfigurabila
	for rows.Next() {
		var r calcul.RegulaConfigurabila
		var activ int
		if err := rows.Scan(&r.ID, &r.Cod, &r.Denumire, &r.Descriere, &r.TipRegula, &r.Severitate, &r.CoduriSchimb, &r.ValoareNumerica, &activ); err != nil {
			return nil, traduceEroare(err)
		}
		r.Activ = activ == 1
		rez = append(rez, r)
	}
	return rez, nil
}

// SalveazaRegulaConfigurabila actualizează sau inserează o regulă configurabilă.
func (s *Serviciu) SalveazaRegulaConfigurabila(ctx context.Context, r calcul.RegulaConfigurabila) error {
	r.Cod = strings.ToUpper(strings.TrimSpace(r.Cod))
	if r.Cod == "" || strings.TrimSpace(r.Denumire) == "" {
		return eroare("reguli.campuri_obligatorii")
	}
	if r.TipRegula != "REPAUS_MINIM" && r.TipRegula != "PLAFON_LUNAR" && r.TipRegula != "SUCCESIUNE_INTERZISA" {
		return eroare("reguli.tip_invalid")
	}
	if r.Severitate != "INFO" && r.Severitate != "AVERTISMENT" && r.Severitate != "BLOCANT" {
		r.Severitate = "AVERTISMENT"
	}

	acum := time.Now().UTC().Format(time.RFC3339)
	activNumar := 0
	if r.Activ {
		activNumar = 1
	}

	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		if r.ID == "" {
			r.ID = idNou()
			_, err := tx.ExecContext(ctx, `
				INSERT INTO reguli_configurabile (id, cod, denumire, descriere, tip_regula, severitate, coduri_schimb, valoare_numerica, activ, creat_la, modificat_la)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, r.ID, r.Cod, r.Denumire, r.Descriere, r.TipRegula, r.Severitate, r.CoduriSchimb, r.ValoareNumerica, activNumar, acum, acum)
			if err != nil {
				return err
			}
			return jurnal(ctx, q, "reguli_configurabile", r.ID, nil, r)
		}

		var vechi calcul.RegulaConfigurabila
		var activVechi int
		if err := tx.QueryRowContext(ctx, `
			SELECT id, cod, denumire, descriere, tip_regula, severitate, coduri_schimb, valoare_numerica, activ
			FROM reguli_configurabile WHERE id = ?
		`, r.ID).Scan(&vechi.ID, &vechi.Cod, &vechi.Denumire, &vechi.Descriere, &vechi.TipRegula, &vechi.Severitate, &vechi.CoduriSchimb, &vechi.ValoareNumerica, &activVechi); err != nil {
			return err
		}
		vechi.Activ = activVechi == 1
		_, err := tx.ExecContext(ctx, `
			UPDATE reguli_configurabile SET
				denumire = ?,
				descriere = ?,
				tip_regula = ?,
				severitate = ?,
				coduri_schimb = ?,
				valoare_numerica = ?,
				activ = ?,
				modificat_la = ?
			WHERE id = ?
		`, r.Denumire, r.Descriere, r.TipRegula, r.Severitate, r.CoduriSchimb, r.ValoareNumerica, activNumar, acum, r.ID)
		if err != nil {
			return err
		}
		r.Cod = vechi.Cod // codul nu se modifică după creare
		return jurnal(ctx, q, "reguli_configurabile", r.ID, vechi, r)
	})
}

// VerificaReguliPlanificare evaluează toate regulile configurabile active peste planificarea lunii selectate.
func (s *Serviciu) VerificaReguliPlanificare(ctx context.Context, an, luna int) ([]calcul.IncalcareRegula, error) {
	reguli, err := s.ListaReguliConfigurabile(ctx)
	if err != nil {
		return nil, err
	}

	// Citim schimburi pentru a mapa codul și orele
	schimburi, err := s.ListaSchimburi(ctx)
	if err != nil {
		return nil, err
	}
	hartaSchimburi := make(map[string]Schimb)
	for _, sc := range schimburi {
		hartaSchimburi[sc.ID] = sc
	}

	// Citim angajații pentru a mapa numele
	date, err := s.Citeste(ctx)
	if err != nil {
		return nil, err
	}
	hartaAngajati := make(map[string]Angajat)
	for _, a := range date.Angajati {
		hartaAngajati[a.ID] = a
	}

	// Citim pozițiile planificate în lună
	plan, err := s.CitestePlanificare(ctx, an, luna)
	if err != nil {
		return nil, err
	}

	var pozitii []calcul.PozitiePentruRegula
	for _, p := range plan.Pozitii {
		if !p.Activ {
			continue
		}
		sc, ok := hartaSchimburi[p.SchimbID]
		if !ok {
			continue
		}
		ang := hartaAngajati[p.AngajatID]
		numeComplet := strings.TrimSpace(ang.Nume + " " + ang.Prenume)
		if numeComplet == "" {
			numeComplet = p.AngajatID
		}

		pozitii = append(pozitii, calcul.PozitiePentruRegula{
			AngajatID: p.AngajatID,
			Nume:      numeComplet,
			Data:      p.Data,
			CodSchimb: sc.Cod,
			OraStart:  int(sc.OraStart),
			DurataMin: int(sc.DurataMinute),
			Activ:     p.Activ,
		})
	}

	return calcul.EvalueazaReguli(reguli, pozitii), nil
}
