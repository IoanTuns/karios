package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

// Prefixul angajatului (D089): titulatura optionala din nomenclatorul
// prefixe, legata prin angajati_prefixe. Numele salvat nu se schimba.

func citestePrefixeAngajati(ctx context.Context, tx *sql.Tx) (map[string]string, error) {
	rows, err := tx.QueryContext(ctx, "SELECT angajat_id,prefix_id FROM angajati_prefixe")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]string{}
	for rows.Next() {
		var a, p string
		if err := rows.Scan(&a, &p); err != nil {
			return nil, err
		}
		rez[a] = p
	}
	return rez, rows.Err()
}

// completeazaPrefixe pune prefixul (ID si denumire) pe fiecare angajat.
func completeazaPrefixe(angajati []Angajat, prefixe []ElementNomenclator, legaturi map[string]string) {
	for i := range angajati {
		angajati[i].PrefixID = legaturi[angajati[i].ID]
		angajati[i].Prefix = nil
		for _, p := range prefixe {
			if p.ID == angajati[i].PrefixID {
				angajati[i].Prefix = p.Denumire
			}
		}
	}
}

// salveazaPrefixAngajat valideaza si inlocuieste legatura; un prefix inactiv
// ramane la angajatii care il au deja, dar nu poate fi ales nou.
func salveazaPrefixAngajat(ctx context.Context, tx *sql.Tx, q *acces.Queries, angajatID, prefixID string) error {
	var vechi string
	if err := tx.QueryRowContext(ctx, "SELECT prefix_id FROM angajati_prefixe WHERE angajat_id=?", angajatID).Scan(&vechi); err != nil && err != sql.ErrNoRows {
		return err
	}
	if vechi == prefixID {
		return nil
	}
	if prefixID != "" {
		var activ int
		if err := tx.QueryRowContext(ctx, "SELECT activ FROM prefixe WHERE id=?", prefixID).Scan(&activ); err != nil {
			return eroare("nomenclator.invalid")
		}
		if activ != 1 {
			return eroare("nomenclator.inactiv")
		}
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM angajati_prefixe WHERE angajat_id=?", angajatID); err != nil {
		return err
	}
	if prefixID != "" {
		if _, err := tx.ExecContext(ctx, "INSERT INTO angajati_prefixe VALUES(?,?)", angajatID, prefixID); err != nil {
			return err
		}
	}
	return jurnal(ctx, q, "angajati_prefixe", angajatID, vechi, prefixID)
}

// prefixeImport potriveste textul din coloana Prefix a importului cu
// nomenclatorul (F07-T41, D089), dupa denumirea RO sau EN, fara diferente de
// majuscule, spatii sau punct ("Dr." = "dr" = "DR.").
type prefixeImport struct {
	dupaText map[string]prefixImport
	// existente: marca normalizata -> prefixul pe care angajatul il are deja.
	existente map[string]string
}

type prefixImport struct {
	id    string
	activ bool
}

func cheiePrefix(s string) string {
	return strings.ReplaceAll(curataAntet(s), ".", "")
}

type interogator interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func citestePrefixeImport(ctx context.Context, db interogator) (prefixeImport, error) {
	rez := prefixeImport{dupaText: map[string]prefixImport{}, existente: map[string]string{}}
	rows, err := db.QueryContext(ctx, "SELECT id,denumire,activ FROM prefixe")
	if err != nil {
		return rez, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, den string
		var activ int
		if err := rows.Scan(&id, &den, &activ); err != nil {
			return rez, err
		}
		var text i18n.Text
		_ = json.Unmarshal([]byte(den), &text)
		for _, v := range text {
			if k := cheiePrefix(v); k != "" {
				// Un prefix activ are intaietate fata de unul inactiv cu acelasi text.
				if vechi, ok := rez.dupaText[k]; !ok || !vechi.activ {
					rez.dupaText[k] = prefixImport{id: id, activ: activ == 1}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return rez, err
	}
	rows.Close()
	rows2, err := db.QueryContext(ctx, "SELECT a.marca,ap.prefix_id FROM angajati_prefixe ap JOIN angajati a ON a.id=ap.angajat_id")
	if err != nil {
		return rez, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var marca, p string
		if err := rows2.Scan(&marca, &p); err != nil {
			return rez, err
		}
		rez.existente[strings.ToUpper(strings.TrimSpace(marca))] = p
	}
	return rez, rows2.Err()
}

// rezolva intoarce ID-ul prefixului sau codul erorii. Textul gol nu cere nimic.
// Un prefix inactiv e acceptat doar daca angajatul cu aceasta marca il are deja.
func (p prefixeImport) rezolva(text, marca string) (string, string) {
	if strings.TrimSpace(text) == "" {
		return "", ""
	}
	g, ok := p.dupaText[cheiePrefix(text)]
	if !ok {
		return "", "import.prefix_necunoscut"
	}
	if !g.activ && p.existente[marca] != g.id {
		return "", "import.prefix_inactiv"
	}
	return g.id, ""
}

// salveaza leaga prefixul de angajat; textul gol lasa legatura existenta.
func (p prefixeImport) salveaza(ctx context.Context, tx *sql.Tx, q *acces.Queries, angajatID, text, marca string) error {
	id, cod := p.rezolva(text, marca)
	if cod != "" {
		return &i18n.Eroare{Cod: cod, Param: map[string]string{"prefix": text}}
	}
	if id == "" {
		return nil
	}
	return salveazaPrefixAngajat(ctx, tx, q, angajatID, id)
}
