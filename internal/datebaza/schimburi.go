package datebaza

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type Schimb struct {
	Revizie      string    `json:"revizie"`
	TipID        string    `json:"tip_id"`
	RegimID      string    `json:"regim_id"`
	ID           string    `json:"id"`
	Cod          string    `json:"cod"`
	Denumire     i18n.Text `json:"denumire"`
	OraStart     int64     `json:"ora_start"`
	DurataMinute int64     `json:"durata_minute"`
	PauzaMinute  int64     `json:"pauza_minute"`
	Tip          string    `json:"tip"`
	Regim        string    `json:"regim"`
	Categorie    string    `json:"categorie"`
	Activ        bool      `json:"activ"`
	// Culoare: aleasa de utilizator (#RRGGBB) sau goala = implicita dupa tip (F07-T42, D097).
	Culoare string `json:"culoare"`
	// CuloareAfisata: culoarea folosita in grile si rapoarte (aleasa sau implicita).
	CuloareAfisata string `json:"culoare_afisata"`
}

// culoriImplicite (D097): nuante deschise, distincte de weekend (gri), sarbatoare (galben),
// absenta (verde) si zi libera (hasura), lizibile cu text inchis.
var culoriImplicite = map[string]string{"NORMAL": "#DDE9F8", "TURA": "#FCE6CC", "NOAPTE": "#E2DCF5", "GARDA": "#F8D7DC"}

// CuloareProgram da culoarea afisata a unui program: cea aleasa sau cea implicita dupa tip.
func CuloareProgram(culoare, tip string) string {
	if c := normalizeazaCuloare(culoare); c != "" {
		return c
	}
	if c := culoriImplicite[tip]; c != "" {
		return c
	}
	return culoriImplicite["NORMAL"]
}

// normalizeazaCuloare intoarce #RRGGBB cu majuscule sau "" daca valoarea nu e o culoare valida.
func normalizeazaCuloare(c string) string {
	c = strings.ToUpper(strings.TrimSpace(c))
	if len(c) != 7 || c[0] != '#' {
		return ""
	}
	for _, r := range c[1:] {
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'F')) {
			return ""
		}
	}
	return c
}

//go:embed initiale/schimburi.json
var schimburiInitiale []byte

func (s *Serviciu) CompleteazaSchimburiInitiale(ctx context.Context, categorie ...string) (int64, error) {
	cat := "LAB"
	if len(categorie) > 0 && strings.TrimSpace(categorie[0]) != "" {
		cat = strings.ToUpper(strings.TrimSpace(categorie[0]))
	}
	if cat != "LAB" {
		return 0, eroare("schimburi.categorie_invalida")
	}
	var initiale []Schimb
	if err := json.Unmarshal(schimburiInitiale, &initiale); err != nil {
		return 0, traduceEroare(err)
	}
	var adaugate int64
	err := s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		var activ int
		if err := tx.QueryRowContext(ctx, "SELECT activ FROM categorii_schimburi WHERE id=?", cat).Scan(&activ); err != nil {
			return err
		}
		if activ != 1 {
			return eroare("nomenclator.inactiv")
		}
		for _, schimb := range initiale {
			if schimb.Categorie == "" {
				schimb.Categorie = "LAB"
			}
			if schimb.Categorie != cat {
				continue
			}
			// Codul poate lipsi după redenumirea unui șablon deja salvat.
			schimb.ID = idNou()
			text, err := json.Marshal(schimb.Denumire)
			if err != nil {
				return err
			}
			n, err := q.AdaugaSchimbInitial(ctx, acces.AdaugaSchimbInitialParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, Tip: schimb.Tip, Regim: schimb.Regim, Categorie: schimb.Categorie, CreatLa: moment(), ModificatLa: moment()})
			if err != nil {
				return err
			}
			if n == 0 {
				continue
			}
			adaugate += n
			if err := jurnal(ctx, q, "tipuri_schimb", schimb.ID, nil, schimb); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return adaugate, nil
}

func (s *Serviciu) ListaSchimburi(ctx context.Context) ([]Schimb, error) {
	randuri, err := acces.New(s.db).ListaSchimburi(ctx)
	if err != nil {
		return nil, traduceEroare(err)
	}
	clasificari, err := acces.New(s.db).ListaClasificariSchimburi(ctx)
	if err != nil {
		return nil, traduceEroare(err)
	}
	rez := []Schimb{}
	for _, r := range randuri {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return nil, traduceEroare(err)
		}
		c, exista := clasificari[r.ID]
		if !exista {
			c = acces.ClasificareSchimb{TipID: r.Tip, RegimID: r.Regim}
		}
		if c.RegimID == "GENERAL" {
			c.RegimID = "LUCRATOARE"
		}
		culoare := ""
		if r.Culoare != nil {
			culoare = *r.Culoare
		}
		rez = append(rez, Schimb{Revizie: c.Revizie, TipID: c.TipID, RegimID: c.RegimID, ID: r.ID, Cod: r.Cod, Denumire: denumire, OraStart: r.OraStart, DurataMinute: r.DurataMinute, PauzaMinute: r.PauzaMinute, Tip: r.Tip, Regim: r.Regim, Categorie: r.Categorie, Activ: r.Activ == 1, Culoare: culoare, CuloareAfisata: CuloareProgram(culoare, r.Tip)})
	}
	return rez, nil
}

func (s *Serviciu) SalveazaSchimb(ctx context.Context, schimb Schimb) error {
	schimb.Cod = strings.ToUpper(strings.TrimSpace(schimb.Cod))
	if schimb.Cod == "" || len(schimb.Cod) > 16 || strings.IndexFunc(schimb.Cod, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_')
	}) >= 0 {
		return eroare("schimburi.cod_invalid")
	}
	denumire := i18n.Text{}
	for l, text := range schimb.Denumire {
		if text = strings.TrimSpace(text); text != "" {
			denumire[l] = text
		}
	}
	if len(denumire) == 0 {
		return eroare("sarbatori.denumire_obligatorie")
	}
	schimb.Denumire = denumire
	if schimb.OraStart < 0 || schimb.OraStart > 1439 || schimb.DurataMinute < 1 || schimb.DurataMinute > 2880 || schimb.PauzaMinute < 0 || schimb.PauzaMinute > schimb.DurataMinute {
		return eroare("schimburi.interval_invalid")
	}
	switch schimb.Tip {
	case "NORMAL", "TURA", "NOAPTE", "GARDA":
	default:
		return eroare("schimburi.tip_invalid")
	}
	var culoare *string
	if c := strings.TrimSpace(schimb.Culoare); c != "" {
		n := normalizeazaCuloare(c)
		if n == "" {
			return eroare("schimburi.culoare_invalida")
		}
		culoare = &n
		schimb.Culoare = n
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		if _, err := q.CautaCodAbsenta(ctx, schimb.Cod); err == nil {
			return eroare("absente.cod_duplicat")
		} else if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if schimb.Regim == "" || schimb.Regim == "GENERAL" {
			schimb.Regim = "LUCRATOARE"
		}
		switch schimb.Regim {
		case "GENERAL", "LUCRATOARE", "NELUCRATOARE", "INTRARE_ZI_LIBERA":
		default:
			return eroare("schimburi.regim_invalid")
		}
		schimb.Categorie = strings.ToUpper(strings.TrimSpace(schimb.Categorie))
		if schimb.Categorie == "" {
			schimb.Categorie = "CUSTOM"
		}
		var activ int
		if err := tx.QueryRowContext(ctx, "SELECT activ FROM categorii_schimburi WHERE id=?", schimb.Categorie).Scan(&activ); err != nil {
			return eroare("schimburi.categorie_invalida")
		}
		if activ != 1 {
			vechi, err := q.CitesteSchimb(ctx, schimb.ID)
			if err != nil || vechi.Categorie != schimb.Categorie {
				return eroare("nomenclator.inactiv")
			}
		}

		idExistent, err := q.CitesteCodSchimb(ctx, schimb.Cod)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if err == nil && idExistent != schimb.ID {
			return eroare("schimburi.cod_duplicat")
		}
		text, err := json.Marshal(denumire)
		if err != nil {
			return err
		}
		clasificari, err := q.ListaClasificariSchimburi(ctx)
		if err != nil {
			return err
		}
		anterioara := clasificari[schimb.ID]
		if schimb.ID != "" && anterioara.TipID == "" {
			if v, e := q.CitesteSchimb(ctx, schimb.ID); e == nil {
				anterioara = acces.ClasificareSchimb{TipID: v.Tip, RegimID: v.Regim}
			}
		}
		if anterioara.RegimID == "GENERAL" {
			anterioara.RegimID = "LUCRATOARE"
		}
		if schimb.TipID == "" {
			schimb.TipID = schimb.Tip
		}
		if schimb.RegimID == "" {
			schimb.RegimID = schimb.Regim
		}
		for _, c := range []struct {
			tabela, id, vechi string
			dest              *string
		}{
			{"tipuri_program", schimb.TipID, anterioara.TipID, &schimb.Tip},
			{"regimuri_program", schimb.RegimID, anterioara.RegimID, &schimb.Regim},
		} {
			var activ int
			if err := tx.QueryRowContext(ctx, "SELECT comportament,activ FROM "+c.tabela+" WHERE id=?", c.id).Scan(c.dest, &activ); err != nil {
				return eroare("nomenclator.invalid")
			}
			if activ != 1 && c.id != c.vechi {
				return eroare("nomenclator.inactiv")
			}
		}
		acum := moment()
		var vechi any
		if schimb.ID == "" {
			schimb.ID = idNou()
			err = q.AdaugaSchimb(ctx, acces.AdaugaSchimbParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, PauzaMinute: schimb.PauzaMinute, Tip: schimb.Tip, Regim: schimb.Regim, Categorie: schimb.Categorie, Activ: numarBool(schimb.Activ), Culoare: culoare, CreatLa: acum, ModificatLa: acum})
		} else {
			r, citireErr := q.CitesteSchimb(ctx, schimb.ID)
			if errors.Is(citireErr, sql.ErrNoRows) {
				return eroare("schimburi.inexistent")
			}
			if citireErr != nil {
				return citireErr
			}
			vechi = r
			_, err = q.ModificaSchimb(ctx, acces.ModificaSchimbParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, PauzaMinute: schimb.PauzaMinute, Tip: schimb.Tip, Regim: schimb.Regim, Categorie: schimb.Categorie, Activ: numarBool(schimb.Activ), Culoare: culoare, ModificatLa: acum})
		}
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO clasificari_schimburi(schimb_id,tip_id,regim_id) VALUES(?,?,?) ON CONFLICT(schimb_id) DO UPDATE SET tip_id=excluded.tip_id,regim_id=excluded.regim_id", schimb.ID, schimb.TipID, schimb.RegimID); err != nil {
			return err
		}
		if err := jurnal(ctx, q, "clasificari_schimburi", schimb.ID, anterioara, acces.ClasificareSchimb{TipID: schimb.TipID, RegimID: schimb.RegimID}); err != nil {
			return err
		}
		return jurnal(ctx, q, "tipuri_schimb", schimb.ID, vechi, schimb)
	})
}

func (s *Serviciu) SchimbaActivSchimb(ctx context.Context, id string, activ bool) error {
	return s.modifica(ctx, func(q *acces.Queries) error {
		vechi, err := q.CitesteSchimb(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return eroare("schimburi.inexistent")
		}
		if err != nil {
			return err
		}
		if _, err := q.SchimbaActivSchimb(ctx, acces.SchimbaActivSchimbParams{ID: id, Activ: numarBool(activ), ModificatLa: moment()}); err != nil {
			return err
		}
		nou, err := q.CitesteSchimb(ctx, id)
		if err != nil {
			return err
		}
		return jurnal(ctx, q, "tipuri_schimb", id, vechi, nou)
	})
}
