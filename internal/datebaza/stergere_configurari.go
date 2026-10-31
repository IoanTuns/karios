package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
	"strings"
)

// StergeConfigurare elimina doar definitii inactive si fara utilizari externe.
// Sesiunea administratorului este verificata de endpointul public App.
func (s *Serviciu) StergeConfigurare(ctx context.Context, tip, id, revizie string) error {
	tabela, cheie, campRev := tabelaNomenclator(tip), "id", "modificat_la"
	switch tip {
	case "schimburi":
		tabela = "tipuri_schimb"
	case "tipuri_absenta":
		tabela = "coduri_absenta"
		cheie = "cod"
	case "categorii_program":
		tabela = "categorii_program"
		campRev = "revizie"
	case "modele_raport":
		tabela = "definitii_rapoarte"
		campRev = "revizie"
	}
	if tabela == "" || id == "" || revizie == "" {
		return eroare("stergere.selectie_invalida")
	}
	protejate := map[string][]string{
		"categorii": {"LAB", "CUSTOM"}, "tipuri_program": {"NORMAL", "TURA", "NOAPTE", "GARDA"},
		"regimuri_program": {"LUCRATOARE", "NELUCRATOARE", "INTRARE_ZI_LIBERA"},
		"contracte":        {"INTREAGA", "PARTIALA", "GARDA"},
		"modele_raport":    {"raport-pontaj", "raport-pontaj-detaliat", "raport-absente", "raport-garzi", "raport-fisa"},
	}
	for _, p := range protejate[tip] {
		if id == p {
			return eroare("stergere.protejat")
		}
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		var activ int
		var actuala, denumire, cod string
		campCod := "''"
		if tip == "schimburi" || tip == "tipuri_absenta" {
			campCod = "cod"
		}
		if err := tx.QueryRowContext(ctx, "SELECT activ,CAST("+campRev+" AS TEXT),denumire,"+campCod+" FROM "+tabela+" WHERE "+cheie+"=?", id).Scan(&activ, &actuala, &denumire, &cod); err != nil {
			return err
		}
		if actuala != revizie {
			return eroare("stergere.conflict")
		}
		if activ != 0 {
			return eroare("stergere.activa")
		}
		folosit := func(sursa string) error {
			return &i18n.Eroare{Cod: "stergere.folosit", Param: map[string]string{"sursa": sursa}}
		}
		rows, err := tx.QueryContext(ctx, `SELECT m.name,f."from" FROM sqlite_master m JOIN pragma_foreign_key_list(m.name) f WHERE m.type='table' AND f."table"=?`, tabela)
		if err != nil {
			return err
		}
		type referinta struct{ tabela, camp string }
		var referinte []referinta
		for rows.Next() {
			var r referinta
			if err := rows.Scan(&r.tabela, &r.camp); err != nil {
				rows.Close()
				return err
			}
			referinte = append(referinte, r)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return err
		}
		citat := func(s string) string { return `"` + strings.ReplaceAll(s, `"`, `""`) + `"` }
		for _, r := range referinte {
			// Numai metadatele proprii pot disparea odata cu definitia, nu alte configurari.
			if tip == "schimburi" && r.tabela == "clasificari_schimburi" && r.camp == "schimb_id" {
				continue
			}
			if tip == "categorii_program" && r.tabela == "categorii_program_schimburi" && r.camp == "categorie_id" {
				continue
			}
			var n int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM "+citat(r.tabela)+" WHERE "+citat(r.camp)+"=?", id).Scan(&n); err != nil {
				return err
			}
			if n > 0 {
				return folosit(r.tabela)
			}
		}
		if tip == "categorii" {
			var n int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM tipuri_schimb WHERE categorie=?", id).Scan(&n); err != nil {
				return err
			}
			if n > 0 {
				return folosit("tipuri_schimb")
			}
		}
		if tip == "schimburi" {
			var n int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM pontaj_ciorne WHERE upper(trim(valoare))=upper(?)", cod).Scan(&n); err != nil {
				return err
			}
			if n > 0 {
				return folosit("pontaj_ciorne")
			}
		}
		for _, r := range []referinta{{"setari", "valoare"}, {"sabloane_ai", "intentie"}, {"operatii_ai", "plan"}, {"operatii_ai", "instantanee"}} {
			var n int
			sqlText := "SELECT count(*) FROM " + citat(r.tabela) + " s,json_tree(CASE WHEN json_valid(s." + citat(r.camp) + ") THEN s." + citat(r.camp) + " ELSE 'null' END) j WHERE j.type='text' AND (j.value=? OR (?<>'' AND j.value=?))"
			if err := tx.QueryRowContext(ctx, sqlText, id, cod, cod).Scan(&n); err != nil {
				return err
			}
			if n > 0 {
				return folosit(r.tabela)
			}
		}
		snapshot := map[string]any{}
		randuri, err := tx.QueryContext(ctx, "SELECT * FROM "+tabela+" WHERE "+cheie+"=?", id)
		if err != nil {
			return err
		}
		coloane, err := randuri.Columns()
		if err != nil {
			randuri.Close()
			return err
		}
		if randuri.Next() {
			valori := make([]any, len(coloane))
			dest := make([]any, len(coloane))
			for i := range valori {
				dest[i] = &valori[i]
			}
			if err := randuri.Scan(dest...); err != nil {
				randuri.Close()
				return err
			}
			for i, c := range coloane {
				snapshot[c] = valori[i]
			}
		}
		err = randuri.Err()
		randuri.Close()
		if err != nil {
			return err
		}
		if tip == "schimburi" {
			clasificari, err := q.ListaClasificariSchimburi(ctx)
			if err != nil {
				return err
			}
			snapshot["clasificare"] = clasificari[id]
		}
		if tip == "categorii_program" {
			ids, err := q.ListaSchimburiCategorieProgram(ctx, id)
			if err != nil {
				return err
			}
			snapshot["programe_permise"] = ids
		}
		if tip == "categorii_program" {
			if _, err := tx.ExecContext(ctx, "DELETE FROM categorii_program_schimburi WHERE categorie_id=?", id); err != nil {
				return err
			}
		}
		rez, err := tx.ExecContext(ctx, "DELETE FROM "+tabela+" WHERE "+cheie+"=? AND activ=0 AND CAST("+campRev+" AS TEXT)=?", id, revizie)
		if err != nil {
			return err
		}
		n, err := rez.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return eroare("stergere.conflict")
		}
		var nume i18n.Text
		if err := json.Unmarshal([]byte(denumire), &nume); err != nil {
			return err
		}
		return jurnal(ctx, q, tabela, id, snapshot, nil)
	})
}
