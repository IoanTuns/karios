package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
	"sort"
	"strings"
	"time"
)

type ElementNomenclator struct {
	Comportament       string    `json:"comportament"`
	Utilizari          int       `json:"utilizari"`
	GrupID             string    `json:"grup_id"`
	Ordine             int       `json:"ordine"`
	Conducere          bool      `json:"conducere"`
	Normare            string    `json:"normare"`
	NormaMinute        *int64    `json:"norma_minute"`
	CategorieProgramID string    `json:"categorie_program_id"`
	Garda              bool      `json:"garda"`
	ID                 string    `json:"id"`
	Denumire           i18n.Text `json:"denumire"`
	Activ              bool      `json:"activ"`
	Revizie            string    `json:"revizie"`
}

func tabelaNomenclator(tip string) string {
	switch tip {
	case "tipuri_program":
		return "tipuri_program"
	case "regimuri_program":
		return "regimuri_program"
	case "categorii":
		return "categorii_schimburi"
	case "grupuri":
		return "grupuri_profesionale"
	case "contracte":
		return "tipuri_contract"
	case "functii":
		return "functii"
	}
	return ""
}
func numeElement(e ElementNomenclator) string {
	if e.Denumire["ro"] != "" {
		return e.Denumire["ro"]
	}
	return e.Denumire["en"]
}
func citesteNomenclator(ctx context.Context, tx *sql.Tx, tip string) ([]ElementNomenclator, error) {
	tabela := tabelaNomenclator(tip)
	if tabela == "" {
		return nil, eroare("nomenclator.invalid")
	}
	extra := ",NULL,0,0,'',NULL,NULL,0"
	switch tip {
	case "grupuri":
		extra = ",NULL,ordine,0,'',NULL,NULL,0"
	case "functii":
		extra = ",grup_id,ordine,conducere,'',NULL,NULL,0"
	case "contracte":
		extra = ",NULL,0,0,normare,norma_minute,categorie_program_id,garda"
	}
	clasificare := ", '', 0"
	if tip == "tipuri_program" || tip == "regimuri_program" {
		camp, baza := "tip_id", "tip"
		if tip == "regimuri_program" {
			camp, baza = "regim_id", "regim"
		}
		clasificare = ",comportament,(SELECT count(*) FROM tipuri_schimb s LEFT JOIN clasificari_schimburi c ON c.schimb_id=s.id WHERE COALESCE(c." + camp + ",s." + baza + ")=" + tabela + ".id)"
	}
	rows, err := tx.QueryContext(ctx, "SELECT id,denumire,activ,modificat_la"+extra+clasificare+" FROM "+tabela+" ORDER BY denumire,id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lista := []ElementNomenclator{}
	for rows.Next() {
		var e ElementNomenclator
		var text string
		var activ, conducere, garda int
		var grup, categorie sql.NullString
		if err := rows.Scan(&e.ID, &text, &activ, &e.Revizie, &grup, &e.Ordine, &conducere, &e.Normare, &e.NormaMinute, &categorie, &garda, &e.Comportament, &e.Utilizari); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(text), &e.Denumire); err != nil {
			return nil, err
		}
		e.Activ = activ == 1
		e.GrupID = grup.String
		e.CategorieProgramID = categorie.String
		e.Conducere = conducere == 1
		e.Garda = garda == 1
		lista = append(lista, e)
	}
	sort.SliceStable(lista, func(i, j int) bool {
		if lista[i].Ordine != lista[j].Ordine {
			return lista[i].Ordine < lista[j].Ordine
		}
		return numeElement(lista[i]) < numeElement(lista[j])
	})
	return lista, rows.Err()
}
func (s *Serviciu) SalveazaElementNomenclator(ctx context.Context, tip string, e ElementNomenclator) error {
	tabela := tabelaNomenclator(tip)
	if tabela == "" {
		return eroare("nomenclator.invalid")
	}
	e.Denumire = i18n.Text{"ro": strings.TrimSpace(e.Denumire["ro"]), "en": strings.TrimSpace(e.Denumire["en"])}
	if e.Ordine < 0 || e.Ordine > 100000 {
		return eroare("nomenclator.invalid")
	}
	if numeElement(e) == "" {
		return eroare("sarbatori.denumire_obligatorie")
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		lista, err := citesteNomenclator(ctx, tx, tip)
		if err != nil {
			return err
		}
		var vechi *ElementNomenclator
		for i := range lista {
			if lista[i].ID == e.ID {
				copie := lista[i]
				vechi = &copie
			}
		}
		if tip == "tipuri_program" || tip == "regimuri_program" {
			permis := e.Comportament == "NORMAL" || e.Comportament == "TURA" || e.Comportament == "NOAPTE" || e.Comportament == "GARDA"
			if tip == "regimuri_program" {
				permis = e.Comportament == "LUCRATOARE" || e.Comportament == "NELUCRATOARE" || e.Comportament == "INTRARE_ZI_LIBERA"
			}
			if !permis {
				return eroare("nomenclator.invalid")
			}
			if vechi != nil && vechi.Utilizari > 0 && vechi.Comportament != e.Comportament {
				return eroare("nomenclator.comportament_utilizat")
			}
		}
		for _, element := range lista {
			if element.ID != e.ID && strings.EqualFold(numeElement(element), numeElement(e)) && (vechi == nil || !strings.EqualFold(numeElement(*vechi), numeElement(e))) {
				return eroare("nomenclator.duplicat")
			}
		}
		if tip == "functii" && e.GrupID != "" {
			var activ int
			if err := tx.QueryRowContext(ctx, "SELECT activ FROM grupuri_profesionale WHERE id=?", e.GrupID).Scan(&activ); err != nil {
				return eroare("nomenclator.invalid")
			}
			if activ != 1 && (vechi == nil || vechi.GrupID != e.GrupID) {
				return eroare("nomenclator.inactiv")
			}
		}
		if tip == "contracte" {
			if e.Normare == "INDIVIDUALA" && vechi != nil && vechi.Normare != "INDIVIDUALA" {
				var n int
				if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM angajati_contracte ac JOIN angajati a ON a.id=ac.angajat_id WHERE ac.contract_id=? AND a.norma_partiala_minute IS NULL", e.ID).Scan(&n); err != nil {
					return err
				}
				if n > 0 {
					return eroare("date.norma_invalida")
				}
			}
			switch e.Normare {
			case "GLOBALA", "INDIVIDUALA", "FARA":
				if e.NormaMinute != nil {
					return eroare("nomenclator.invalid")
				}
			case "ZILNICA", "LUNARA":
				if e.NormaMinute == nil || *e.NormaMinute <= 0 || *e.NormaMinute > 44640 || (e.Normare == "ZILNICA" && *e.NormaMinute > 1440) {
					return eroare("date.norma_invalida")
				}
			default:
				return eroare("nomenclator.invalid")
			}
			if e.CategorieProgramID != "" {
				c, err := q.CitesteCategorieProgram(ctx, e.CategorieProgramID)
				if err != nil {
					return eroare("program.categorie_invalida")
				}
				if c.Activ != 1 && (vechi == nil || vechi.CategorieProgramID != e.CategorieProgramID) {
					return eroare("program.categorie_inactiva")
				}
			}
		}
		text, err := json.Marshal(e.Denumire)
		if err != nil {
			return err
		}
		acum := moment()
		if vechi != nil {
			acum = revizieNomenclatorDupa(vechi.Revizie, acum)
		}
		if e.ID == "" {
			e.ID = idNou()
			if tip == "categorii" {
				e.ID = strings.ToUpper(e.ID)
			}
			_, err = tx.ExecContext(ctx, "INSERT INTO "+tabela+"(id,denumire,activ,modificat_la) VALUES(?,?,?,?)", e.ID, string(text), numarBool(e.Activ), acum)
		} else {
			if vechi == nil {
				return eroare("nomenclator.invalid")
			}
			if vechi.Revizie != e.Revizie {
				return eroare("nomenclator.conflict")
			}
			var r sql.Result
			r, err = tx.ExecContext(ctx, "UPDATE "+tabela+" SET denumire=?,activ=?,modificat_la=? WHERE id=? AND modificat_la=?", string(text), numarBool(e.Activ), acum, e.ID, e.Revizie)
			if err == nil {
				n, _ := r.RowsAffected()
				if n != 1 {
					return eroare("nomenclator.conflict")
				}
			}
		}
		if err != nil {
			return err
		}
		switch tip {
		case "tipuri_program", "regimuri_program":
			_, err = tx.ExecContext(ctx, "UPDATE "+tabela+" SET comportament=? WHERE id=?", e.Comportament, e.ID)
		case "grupuri":
			_, err = tx.ExecContext(ctx, "UPDATE grupuri_profesionale SET ordine=? WHERE id=?", e.Ordine, e.ID)
		case "functii":
			_, err = tx.ExecContext(ctx, "UPDATE functii SET grup_id=?,ordine=?,conducere=? WHERE id=?", textOptional(e.GrupID), e.Ordine, numarBool(e.Conducere), e.ID)
		case "contracte":
			_, err = tx.ExecContext(ctx, "UPDATE tipuri_contract SET normare=?,norma_minute=?,categorie_program_id=?,garda=? WHERE id=?", e.Normare, e.NormaMinute, textOptional(e.CategorieProgramID), numarBool(e.Garda), e.ID)
		}
		if err != nil {
			return err
		}
		if tip == "functii" && vechi != nil && numeElement(*vechi) != numeElement(e) {
			rows, err := tx.QueryContext(ctx, "SELECT angajat_id FROM angajati_functii WHERE functie_id=?", e.ID)
			if err != nil {
				return err
			}
			ids := []string{}
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err != nil {
					rows.Close()
					return err
				}
				ids = append(ids, id)
			}
			err = rows.Err()
			rows.Close()
			if err != nil {
				return err
			}
			for _, id := range ids {
				inainte, err := q.CitesteAngajat(ctx, id)
				if err != nil {
					return err
				}
				if _, err := tx.ExecContext(ctx, "UPDATE angajati SET functie=?,modificat_la=? WHERE id=?", numeElement(e), acum, id); err != nil {
					return err
				}
				if _, err := tx.ExecContext(ctx, "INSERT INTO angajati_functii VALUES(?,?) ON CONFLICT(angajat_id) DO UPDATE SET functie_id=excluded.functie_id", id, e.ID); err != nil {
					return err
				}
				dupa, err := q.CitesteAngajat(ctx, id)
				if err != nil {
					return err
				}
				if err := jurnal(ctx, q, "angajati", id, inainte, dupa); err != nil {
					return err
				}
			}
		}
		e.Revizie = acum
		return jurnal(ctx, q, tabela, e.ID, vechi, e)
	})
}

// Ceasul Windows poate intoarce aceeasi valoare pentru doua scrieri rapide.
func revizieNomenclatorDupa(veche, curenta string) string {
	v, e1 := time.Parse(time.RFC3339Nano, veche)
	c, e2 := time.Parse(time.RFC3339Nano, curenta)
	if e1 == nil && e2 == nil && !c.After(v) {
		return v.Add(time.Nanosecond).Format(time.RFC3339Nano)
	}
	return curenta
}
