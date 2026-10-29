package datebaza

import (
	"context"
	"database/sql"
	"strings"
)

// IntrareJurnal e un rând din jurnal_modificari, pentru ecranul jurnalului (F12-T02).
type IntrareJurnal struct {
	ID     int64  `json:"id"`
	Tabel  string `json:"tabel"`
	RandID string `json:"rand_id"`
	Camp   string `json:"camp"`
	Vechi  string `json:"vechi"`
	Nou    string `json:"nou"`
	Cand   string `json:"cand"`
	Cine   string `json:"cine"`
}

// FiltruJurnal restrânge intrările. AngajatID caută rândurile acelui angajat (inclusiv
// zilele lui din pontaj și absențe); Luna (AAAA-LL) caută rândurile care privesc luna.
type FiltruJurnal struct {
	Tabel     string `json:"tabel"`
	AngajatID string `json:"angajat_id"`
	Luna      string `json:"luna"`
	Text      string `json:"text"`
	Limita    int    `json:"limita"`
	Deplasare int    `json:"deplasare"`
}

// PaginaJurnal întoarce intrările, cele mai noi primele, numărul total și tabelele existente.
type PaginaJurnal struct {
	Intrari []IntrareJurnal `json:"intrari"`
	Total   int             `json:"total"`
	Tabele  []string        `json:"tabele"`
}

const limitaJurnal = 200

func (f FiltruJurnal) conditii() (string, []any) {
	var c []string
	var args []any
	if f.Tabel != "" {
		c = append(c, "tabel = ?")
		args = append(args, f.Tabel)
	}
	if f.AngajatID != "" {
		c = append(c, `(rand_id = ? OR rand_id LIKE ? ESCAPE '\')`)
		args = append(args, f.AngajatID, escapeLike(f.AngajatID)+"/%")
	}
	if f.Luna != "" {
		c = append(c, `(rand_id = ? OR rand_id LIKE ? ESCAPE '\' OR rand_id LIKE ? ESCAPE '\')`)
		l := escapeLike(f.Luna)
		args = append(args, f.Luna, "%/"+l+"-%", l+"-%")
	}
	if t := strings.TrimSpace(f.Text); t != "" {
		c = append(c, "(instr(lower(coalesce(valoare_veche, '')), lower(?)) > 0 OR instr(lower(coalesce(valoare_noua, '')), lower(?)) > 0 OR instr(lower(rand_id), lower(?)) > 0)")
		args = append(args, t, t, t)
	}
	if len(c) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(c, " AND "), args
}

func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\`, "%", `\%`, "_", `\_`).Replace(s)
}

// CitesteJurnal întoarce o pagină din jurnalul de modificări.
func (s *Serviciu) CitesteJurnal(ctx context.Context, f FiltruJurnal) (PaginaJurnal, error) {
	if f.Limita <= 0 || f.Limita > limitaJurnal {
		f.Limita = limitaJurnal
	}
	if f.Deplasare < 0 {
		f.Deplasare = 0
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return PaginaJurnal{}, traduceEroare(err)
	}
	defer tx.Rollback()
	rez := PaginaJurnal{Intrari: []IntrareJurnal{}, Tabele: []string{}}
	unde, args := f.conditii()
	if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM jurnal_modificari"+unde, args...).Scan(&rez.Total); err != nil {
		return rez, traduceEroare(err)
	}
	r, err := tx.QueryContext(ctx, "SELECT id, tabel, rand_id, coalesce(camp, ''), coalesce(valoare_veche, ''), coalesce(valoare_noua, ''), cand, cine FROM jurnal_modificari"+
		unde+" ORDER BY id DESC LIMIT ? OFFSET ?", append(args, f.Limita, f.Deplasare)...)
	if err != nil {
		return rez, traduceEroare(err)
	}
	for r.Next() {
		var i IntrareJurnal
		if err := r.Scan(&i.ID, &i.Tabel, &i.RandID, &i.Camp, &i.Vechi, &i.Nou, &i.Cand, &i.Cine); err != nil {
			r.Close()
			return rez, traduceEroare(err)
		}
		rez.Intrari = append(rez.Intrari, i)
	}
	r.Close()
	if err := r.Err(); err != nil {
		return rez, traduceEroare(err)
	}
	t, err := tx.QueryContext(ctx, "SELECT DISTINCT tabel FROM jurnal_modificari ORDER BY tabel")
	if err != nil {
		return rez, traduceEroare(err)
	}
	defer t.Close()
	for t.Next() {
		var x string
		if err := t.Scan(&x); err != nil {
			return rez, traduceEroare(err)
		}
		rez.Tabele = append(rez.Tabele, x)
	}
	return rez, traduceEroare(t.Err())
}

// TotJurnalul întoarce toate intrările care trec de filtru, fără paginare, pentru export.
func (s *Serviciu) TotJurnalul(ctx context.Context, f FiltruJurnal) ([]IntrareJurnal, error) {
	var tot []IntrareJurnal
	f.Limita, f.Deplasare = limitaJurnal, 0
	for {
		p, err := s.CitesteJurnal(ctx, f)
		if err != nil {
			return nil, err
		}
		tot = append(tot, p.Intrari...)
		if len(p.Intrari) < limitaJurnal {
			return tot, nil
		}
		f.Deplasare += limitaJurnal
	}
}

// UltimaModificare întoarce cea mai nouă intrare din jurnal pentru un rând (F08-T10).
func (s *Serviciu) UltimaModificare(ctx context.Context, tabel, randID string) (*IntrareJurnal, error) {
	var i IntrareJurnal
	err := s.db.QueryRowContext(ctx, `SELECT id, tabel, rand_id, coalesce(camp, ''), coalesce(valoare_veche, ''), coalesce(valoare_noua, ''), cand, cine
		FROM jurnal_modificari WHERE tabel = ? AND rand_id = ? ORDER BY id DESC LIMIT 1`, tabel, randID).
		Scan(&i.ID, &i.Tabel, &i.RandID, &i.Camp, &i.Vechi, &i.Nou, &i.Cand, &i.Cine)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, traduceEroare(err)
	}
	return &i, nil
}
