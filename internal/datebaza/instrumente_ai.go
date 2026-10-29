package datebaza

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"pontix/internal/ai"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

type ExplicatieSelectieAI struct {
	ID     string `json:"id"`
	Marca  string `json:"marca"`
	Nume   string `json:"nume"`
	Inclus bool   `json:"inclus"`
	Motiv  string `json:"motiv"`
}

// Instantaneele sunt exclusiv în baza locală; clientul nu poate alege ce se restaurează.
type InstantaneuCelulaAI struct {
	AngajatID string
	Data      string
	Ciorna    *acces.PontajCiorne
	Pozitii   []acces.Planificare
	Efective  []acces.ListaPontajePrecompletareRow
	Dupa      *InstantaneuCelulaAI `json:",omitempty"`
}

func citesteInstantaneuAI(ctx context.Context, tx *sql.Tx, id, data string) (InstantaneuCelulaAI, error) {
	r := InstantaneuCelulaAI{AngajatID: id, Data: data, Pozitii: []acces.Planificare{}, Efective: []acces.ListaPontajePrecompletareRow{}}
	var c acces.PontajCiorne
	err := tx.QueryRowContext(ctx, `SELECT angajat_id,data,valoare,modificat_la,revizie FROM pontaj_ciorne WHERE angajat_id=? AND data=?`, id, data).Scan(&c.AngajatID, &c.Data, &c.Valoare, &c.ModificatLa, &c.Revizie)
	if err != nil && err != sql.ErrNoRows {
		return r, err
	}
	if err == nil {
		r.Ciorna = &c
	}
	pozitii, err := acces.New(tx).ListaPlanificare(ctx, acces.ListaPlanificareParams{DataDeLa: data, DataPanaLa: data})
	if err != nil {
		return r, err
	}
	for _, p := range pozitii {
		if p.AngajatID == id {
			r.Pozitii = append(r.Pozitii, p)
		}
	}
	efective, err := acces.New(tx).ListaPontajePrecompletare(ctx, acces.ListaPontajePrecompletareParams{DataDeLa: data, DataPanaLa: data})
	if err != nil {
		return r, err
	}
	for _, p := range efective {
		if p.AngajatID == id {
			r.Efective = append(r.Efective, p)
		}
	}
	return r, nil
}

type RandAnulareAI struct {
	AngajatID     string   `json:"angajat_id"`
	Nume          string   `json:"nume"`
	Data          string   `json:"data"`
	Acum          string   `json:"acum"`
	Restaurat     string   `json:"restaurat"`
	PlanAcum      []string `json:"plan_acum"`
	PlanRestaurat []string `json:"plan_restaurat"`
	Motiv         string   `json:"motiv"`
}
type AnulareAI struct {
	OperatieID string          `json:"operatie_id"`
	CreatLa    string          `json:"creat_la"`
	Randuri    []RandAnulareAI `json:"randuri"`
	Token      string          `json:"token"`
	Conflicte  int             `json:"conflicte"`
	Anulata    bool            `json:"anulata"`
}

func pregatesteAnulareAI(ctx context.Context, tx *sql.Tx) (AnulareAI, []InstantaneuCelulaAI, error) {
	r := AnulareAI{Randuri: []RandAnulareAI{}}
	var imagini, continut string
	err := tx.QueryRowContext(ctx, `SELECT id,creat_la,plan,instantanee FROM operatii_ai WHERE anulat_la='' ORDER BY rowid DESC LIMIT 1`).Scan(&r.OperatieID, &r.CreatLa, &continut, &imagini)
	if err == sql.ErrNoRows {
		return r, nil, nil
	}
	if err != nil {
		return r, nil, err
	}
	var p PlanAI
	var instantanee []InstantaneuCelulaAI
	if err = json.Unmarshal([]byte(continut), &p); err != nil {
		return r, nil, err
	}
	if err = json.Unmarshal([]byte(imagini), &instantanee); err != nil {
		return r, nil, err
	}
	nume := map[string]string{}
	for _, a := range p.Persoane {
		nume[a.ID] = a.Nume + " (" + a.Marca + ")"
	}
	curente := []InstantaneuCelulaAI{}
	prima, ultima := p.Intentie.DeLa, p.Intentie.PanaLa
	surse, err := sursePentruPlanificare(ctx, acces.New(tx), prima, ultima)
	if err != nil {
		return r, nil, err
	}
	ignorate := map[string]bool{}
	for _, v := range instantanee {
		ignorate["c/"+v.AngajatID+"/"+v.Data] = true
		if v.Dupa != nil {
			for _, poz := range v.Dupa.Pozitii {
				ignorate["p/"+poz.ID] = true
			}
		}
	}
	intervale := append([]calcul.IntervalPlanificare(nil), surse.intervale...)
	for _, poz := range surse.pozitii {
		if poz.Activ == 1 {
			pr := surse.schimburi[poz.SchimbID]
			d, _ := time.Parse(time.DateOnly, poz.Data)
			start := d.Add(time.Duration(pr.OraStart) * time.Minute)
			intervale = append(intervale, calcul.IntervalPlanificare{AngajatID: poz.AngajatID, Cheie: "p/" + poz.ID, Inceput: start, Sfarsit: start.Add(time.Duration(pr.DurataMinute) * time.Minute)})
		}
	}
	for _, v := range instantanee {
		acum, err := citesteInstantaneuAI(ctx, tx, v.AngajatID, v.Data)
		if err != nil {
			return r, nil, err
		}
		curente = append(curente, acum)
		rand := RandAnulareAI{AngajatID: v.AngajatID, Nume: nume[v.AngajatID], Data: v.Data}
		rand.PlanAcum = []string{}
		rand.PlanRestaurat = []string{}
		for _, poz := range acum.Pozitii {
			if poz.Activ == 1 {
				rand.PlanAcum = append(rand.PlanAcum, surse.schimburi[poz.SchimbID].Cod)
			}
		}
		for _, poz := range v.Pozitii {
			if poz.Activ == 1 {
				rand.PlanRestaurat = append(rand.PlanRestaurat, surse.schimburi[poz.SchimbID].Cod)
			}
		}
		if acum.Ciorna != nil {
			rand.Acum = acum.Ciorna.Valoare
		}
		if v.Ciorna != nil {
			rand.Restaurat = v.Ciorna.Valoare
		}
		// Verificăm și restricțiile curente: o absență sau categorie nouă nu
		// trebuie ocolită prin restaurarea unui program anterior.
		valori := []string{}
		if strings.TrimSpace(rand.Restaurat) != "" {
			valori = append(valori, rand.Restaurat)
		}
		for _, poz := range v.Pozitii {
			if poz.Activ == 1 {
				valori = append(valori, surse.schimburi[poz.SchimbID].Cod)
			}
		}
		for _, valoare := range valori {
			ang, exista := surse.oameni[v.AngajatID]
			if !exista || ang.Activ != 1 {
				rand.Motiv = "restrictii"
				continue
			}
			if surse.absente[v.AngajatID+"/"+v.Data] {
				rand.Motiv = "restrictii"
				continue
			}
			pr, ok := surse.coduri[strings.ToUpper(strings.TrimSpace(valoare))]
			if !ok {
				if _, manual := calcul.OreManuale(valoare); !manual {
					rand.Motiv = "restrictii"
				}
				continue
			}
			if !pr.Activ {
				rand.Motiv = "restrictii"
				continue
			}
			if id := textSauGol(ang.CategorieProgramID); id != "" {
				cat, ok := surse.categorii[id]
				permis := false
				for _, id := range cat.SchimburiID {
					permis = permis || id == pr.ID
				}
				if !ok || !cat.Activ || !permis {
					rand.Motiv = "restrictii"
					continue
				}
			}
			propriu := "restaurat/" + v.AngajatID + "/" + v.Data
			ignorate[propriu] = true
			motiv := calcul.VerificaTuraPlanificata(v.AngajatID, v.Data, pr, surse.absente, intervale, ignorate)
			delete(ignorate, propriu)
			if motiv != "" {
				rand.Motiv = "restrictii"
				continue
			}
			d, _ := time.Parse(time.DateOnly, v.Data)
			start := d.Add(time.Duration(pr.OraStart) * time.Minute)
			intervale = append(intervale, calcul.IntervalPlanificare{AngajatID: v.AngajatID, Cheie: "restaurat/" + v.AngajatID + "/" + v.Data, Inceput: start, Sfarsit: start.Add(time.Duration(pr.DurataMinute) * time.Minute)})
		}
		if v.Dupa == nil || !reflect.DeepEqual(acum, *v.Dupa) {
			rand.Motiv = "modificata"
		}
		d, _ := time.Parse(time.DateOnly, v.Data)
		stare := "DESCHISA"
		err = tx.QueryRowContext(ctx, `SELECT stare FROM stare_luni WHERE an=? AND luna=?`, d.Year(), int(d.Month())).Scan(&stare)
		if err != nil && err != sql.ErrNoRows {
			return r, nil, err
		}
		if stare != "DESCHISA" {
			rand.Motiv = "luna_inchisa"
		}
		if rand.Motiv != "" {
			r.Conflicte++
		}
		r.Randuri = append(r.Randuri, rand)
	}
	date, err := json.Marshal(struct {
		Rezultat AnulareAI
		Curente  []InstantaneuCelulaAI
	}{r, curente})
	if err != nil {
		return r, nil, err
	}
	suma := sha256.Sum256(date)
	r.Token = hex.EncodeToString(suma[:])
	return r, instantanee, nil
}

func (s *Serviciu) PrevizualizeazaAnulareAI(ctx context.Context) (AnulareAI, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return AnulareAI{}, traduceEroare(err)
	}
	defer tx.Rollback()
	r, _, err := pregatesteAnulareAI(ctx, tx)
	if err == nil {
		err = tx.Commit()
	}
	return r, traduceEroare(err)
}

func (s *Serviciu) AnuleazaUltimaOperatieAI(ctx context.Context, token string) (AnulareAI, error) {
	var rezultat AnulareAI
	err := s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		r, imagini, err := pregatesteAnulareAI(ctx, tx)
		if err != nil {
			return err
		}
		if r.OperatieID == "" || token == "" || token != r.Token || r.Conflicte != 0 {
			return eroare("ai.anulare_conflict")
		}
		for _, v := range imagini {
			valoare := ""
			if v.Ciorna != nil {
				valoare = v.Ciorna.Valoare
			}
			if _, err := q.SalveazaCiorna(ctx, acces.SalveazaCiornaParams{AngajatID: v.AngajatID, Data: v.Data, Valoare: valoare, ModificatLa: moment()}); err != nil {
				return err
			}
			anterioare := map[string]acces.Planificare{}
			for _, p := range v.Pozitii {
				anterioare[p.ID] = p
			}
			// Dezactivăm întâi pozițiile noi, pentru a respecta indexul unic al celor active.
			for _, p := range v.Dupa.Pozitii {
				if _, ok := anterioare[p.ID]; !ok {
					if _, err := tx.ExecContext(ctx, `UPDATE planificare SET activ=0,revizie=revizie+1,modificat_la=? WHERE id=?`, moment(), p.ID); err != nil {
						return err
					}
				}
			}
			for _, p := range v.Pozitii {
				var dupa *acces.Planificare
				for n := range v.Dupa.Pozitii {
					if v.Dupa.Pozitii[n].ID == p.ID {
						dupa = &v.Dupa.Pozitii[n]
					}
				}
				if dupa != nil && !reflect.DeepEqual(p, *dupa) {
					if _, err := tx.ExecContext(ctx, `UPDATE planificare SET schimb_id=?,sectiune=?,activ=?,publicat_plan_revizie=?,publicat_pontaj_revizie=?,publicat_valoare=?,publicat_la=?,revizie=revizie+1,modificat_la=? WHERE id=?`, p.SchimbID, p.Sectiune, p.Activ, p.PublicatPlanRevizie, p.PublicatPontajRevizie, p.PublicatValoare, p.PublicatLa, moment(), p.ID); err != nil {
						return err
					}
				}
			}
			actual, err := citesteInstantaneuAI(ctx, tx, v.AngajatID, v.Data)
			if err != nil {
				return err
			}
			if err := jurnal(ctx, q, "pontaj_ciorne", v.AngajatID+"/"+v.Data, v.Dupa.Ciorna, actual.Ciorna); err != nil {
				return err
			}
			if !reflect.DeepEqual(v.Dupa.Pozitii, actual.Pozitii) {
				if err := jurnal(ctx, q, "planificare", v.AngajatID+"/"+v.Data, v.Dupa.Pozitii, actual.Pozitii); err != nil {
					return err
				}
			}
		}
		if _, err := tx.ExecContext(ctx, `UPDATE operatii_ai SET anulat_la=? WHERE id=?`, moment(), r.OperatieID); err != nil {
			return err
		}
		r.Anulata = true
		if err := jurnal(ctx, q, "operatii_ai", r.OperatieID, nil, r); err != nil {
			return err
		}
		rezultat = r
		return nil
	})
	return rezultat, err
}

type SablonAI struct {
	ID       string          `json:"id"`
	Nume     string          `json:"nume"`
	Intentie ai.IntentiePlan `json:"intentie"`
}

func (s *Serviciu) ListaSabloaneAI(ctx context.Context) ([]SablonAI, error) {
	r := []SablonAI{}
	rows, err := s.db.QueryContext(ctx, `SELECT id,nume,intentie FROM sabloane_ai ORDER BY nume,id`)
	if err != nil {
		return r, traduceEroare(err)
	}
	defer rows.Close()
	for rows.Next() {
		var v SablonAI
		var continut string
		if err := rows.Scan(&v.ID, &v.Nume, &continut); err != nil {
			return r, traduceEroare(err)
		}
		if err := json.Unmarshal([]byte(continut), &v.Intentie); err != nil {
			return r, traduceEroare(err)
		}
		r = append(r, v)
	}
	return r, traduceEroare(rows.Err())
}
func (s *Serviciu) SalveazaSablonAI(ctx context.Context, nume string, i ai.IntentiePlan) error {
	nume = strings.TrimSpace(nume)
	if len([]rune(nume)) < 1 || len([]rune(nume)) > 100 {
		return eroare("ai.sablon_nume_invalid")
	}
	// Șablonul nu păstrează cererea liberă și nu transportă un token de aplicare.
	i.MesajOriginal = ""
	i.Clarificari = []string{}
	date, err := json.Marshal(i)
	if err != nil {
		return traduceEroare(err)
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		var n int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM sabloane_ai`).Scan(&n); err != nil {
			return err
		}
		if n >= 100 {
			return eroare("ai.sablon_limita")
		}
		id := idNou()
		if _, err := tx.ExecContext(ctx, `INSERT INTO sabloane_ai(id,nume,intentie,modificat_la) VALUES(?,?,?,?)`, id, nume, string(date), moment()); err != nil {
			return err
		}
		return jurnal(ctx, q, "sabloane_ai", id, nil, SablonAI{id, nume, i})
	})
}
func (s *Serviciu) StergeSablonAI(ctx context.Context, id string) error {
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		var continut string
		err := tx.QueryRowContext(ctx, `SELECT intentie FROM sabloane_ai WHERE id=?`, id).Scan(&continut)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM sabloane_ai WHERE id=?`, id); err != nil {
			return err
		}
		return jurnal(ctx, q, "sabloane_ai", id, continut, nil)
	})
}

// PlanSalvatAI servește exportul din instantaneul verificat, nu din rânduri furnizate de UI.
func (s *Serviciu) PlanSalvatAI(ctx context.Context, id string) (PlanAI, error) {
	var r PlanAI
	var continut, anulat string
	err := s.db.QueryRowContext(ctx, `SELECT plan,anulat_la FROM operatii_ai WHERE id=?`, id).Scan(&continut, &anulat)
	if err == nil {
		err = json.Unmarshal([]byte(continut), &r)
	}
	r.Anulat = anulat != ""
	return r, traduceEroare(err)
}
