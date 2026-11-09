package datebaza

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

// SchimbareContract: norma zilnica si/sau dreptul anual la concediu, valabile de la o data
// (F07-T43, I010, D101). Valorile din fisa angajatului sunt cele de dinaintea primei schimbari.
type SchimbareContract = acces.SchimbareContract

// ListaIstoricContract intoarce schimbarile unui angajat, in ordinea datei.
func (s *Serviciu) ListaIstoricContract(ctx context.Context, angajatID string) ([]SchimbareContract, error) {
	toate, err := acces.New(s.db).ListaIstoricContract(ctx)
	if err != nil {
		return nil, traduceEroare(err)
	}
	rez := toate[angajatID]
	if rez == nil {
		rez = []SchimbareContract{}
	}
	return rez, nil
}

func valideazaSchimbareContract(x *SchimbareContract) error {
	x.AngajatID = strings.TrimSpace(x.AngajatID)
	x.DeLa = strings.TrimSpace(x.DeLa)
	if d, err := time.Parse(time.DateOnly, x.DeLa); err != nil || d.Year() < 1900 || d.Year() > 2200 {
		return eroare("istoric.data_invalida")
	}
	if x.NormaMinute == nil && x.ConcediuZileAn == nil {
		return eroare("istoric.valoare_lipsa")
	}
	if x.NormaMinute != nil && (*x.NormaMinute <= 0 || *x.NormaMinute > 1440) {
		return eroare("istoric.norma_invalida")
	}
	if x.ConcediuZileAn != nil && (*x.ConcediuZileAn < 0 || *x.ConcediuZileAn > 366) {
		return eroare("istoric.concediu_invalid")
	}
	return nil
}

// verificaIstoricModificabil refuza o schimbare care ar atinge o luna inchisa (alegerea Ioan Tuns,
// 2026-11-09): o schimbare de la data D se aplica tuturor zilelor de dupa D.
func verificaIstoricModificabil(ctx context.Context, tx *sql.Tx, q *acces.Queries, angajatID, deLa string) error {
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM angajati WHERE id=? AND sters_la IS NULL`, angajatID).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		return eroare("istoric.angajat_inexistent")
	}
	inchisa, err := q.SoldAfecteazaLunaInchisa(ctx, deLa[:7])
	if err != nil {
		return err
	}
	if inchisa {
		return eroare("istoric.luna_inchisa")
	}
	return nil
}

// AdaugaSchimbareContract adauga o schimbare; doua schimbari ale aceluiasi angajat nu pot avea aceeasi data.
func (s *Serviciu) AdaugaSchimbareContract(ctx context.Context, x SchimbareContract) error {
	if err := valideazaSchimbareContract(&x); err != nil {
		return err
	}
	x.CreatLa = moment()
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		if err := verificaIstoricModificabil(ctx, tx, q, x.AngajatID, x.DeLa); err != nil {
			return err
		}
		n, err := q.AdaugaSchimbareContract(ctx, x)
		if err != nil {
			return err
		}
		if n == 0 {
			return eroare("istoric.exista")
		}
		return jurnal(ctx, q, "istoric_contract", x.AngajatID+"/"+x.DeLa, nil, x)
	})
}

// StergeSchimbareContract sterge schimbarea de la data deLa.
func (s *Serviciu) StergeSchimbareContract(ctx context.Context, angajatID, deLa string) error {
	x := SchimbareContract{AngajatID: strings.TrimSpace(angajatID), DeLa: strings.TrimSpace(deLa)}
	if _, err := time.Parse(time.DateOnly, x.DeLa); err != nil {
		return eroare("istoric.data_invalida")
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		if err := verificaIstoricModificabil(ctx, tx, q, x.AngajatID, x.DeLa); err != nil {
			return err
		}
		toate, err := q.ListaIstoricContract(ctx)
		if err != nil {
			return err
		}
		var vechi *SchimbareContract
		for _, v := range toate[x.AngajatID] {
			if v.DeLa == x.DeLa {
				v := v
				vechi = &v
			}
		}
		n, err := q.StergeSchimbareContract(ctx, x.AngajatID, x.DeLa)
		if err != nil {
			return err
		}
		if n == 0 || vechi == nil {
			return eroare("istoric.inexistenta")
		}
		return jurnal(ctx, q, "istoric_contract", x.AngajatID+"/"+x.DeLa, vechi, nil)
	})
}

// schimbariCalcul transforma randurile din istoric pentru motorul de calcul.
func schimbariCalcul(r []SchimbareContract) []calcul.SchimbareContract {
	rez := make([]calcul.SchimbareContract, 0, len(r))
	for _, x := range r {
		rez = append(rez, calcul.SchimbareContract{DeLa: x.DeLa, NormaMinute: x.NormaMinute, ConcediuZileAn: x.ConcediuZileAn})
	}
	return rez
}

// dreptConcediuAn calculeaza dreptul la concediu pe anul an pentru toti angajatii (D101).
func dreptConcediuAn(ctx context.Context, tx *sql.Tx, an int) (map[string]int64, error) {
	istoric, err := acces.New(tx).ListaIstoricContract(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT id,concediu_zile_an FROM angajati WHERE sters_la IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]int64{}
	for rows.Next() {
		var id string
		var zile *int64
		if err := rows.Scan(&id, &zile); err != nil {
			return nil, err
		}
		if d := calcul.DreptConcediuAn(an, zile, schimbariCalcul(istoric[id])); d != nil {
			rez[id] = *d
		}
	}
	return rez, rows.Err()
}
