// Package audit face vizibilă orice modificare a bazei făcută în afara aplicației (D044).
//
// Pentru fiecare tabel, sigilii(tabel, valoare) păstrează suma modulo 2^256 a valorilor
// HMAC(cheie_audit, tabel ‖ rând) pentru toate rândurile. Aplicația creează pe fiecare
// conexiune trigger-e TEMP care țin suma la zi în aceeași tranzacție cu scrierea; un alt
// program nu le are și nici nu cunoaște cheia, deci sigiliul nu mai corespunde.
// Rândurile din jurnal_modificari sunt în plus înlănțuite: mac = HMAC(mac anterior ‖ rând).
package audit

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ncruces/go-sqlite3"

	"pontix/internal/i18n"
)

// Stările verificării.
const (
	StareInactiv  = "inactiv" // baza nu are încă tabelul sigilii (înainte de migrarea 0021)
	StareOK       = "ok"
	StareIncident = "incident"
)

// CodEditareBlocata e codul i18n al erorii date oricărei scrieri cât timp există un incident.
const CodEditareBlocata = "audit.editare_blocata"

const (
	tabelSigilii = "sigilii"
	tabelJurnal  = "jurnal_modificari"
	// pseudoSchema sigilează sqlite_schema: un trigger, o vedere sau un index adăugat
	// din afară schimbă sigiliul.
	pseudoSchema = "sqlite_schema"
	// lungimeSigiliu: suma se ține pe 256 de biți, cât un HMAC-SHA256.
	lungimeSigiliu = sha256.Size
)

// Rezultat descrie ultima verificare a bazei.
type Rezultat struct {
	Stare       string   `json:"stare"`
	Tabele      []string `json:"tabele"`    // tabelele al căror sigiliu nu corespunde
	RandRupt    int64    `json:"rand_rupt"` // primul rând din jurnal cu lanțul rupt; 0 = niciunul
	VerificatLa string   `json:"verificat_la"`
	DurataMs    int64    `json:"durata_ms"`
}

// Paza ține rezultatul verificării pentru o bază deschisă și blochează scrierile
// cât timp există un incident neconfirmat.
type Paza struct {
	mu         sync.RWMutex
	rez        Rezultat
	confirmare bool
	// Conturile (F15, D045): până la autentificare, scrierile sunt refuzate, mai puțin
	// operațiile de sistem (autentificarea, crearea primului cont).
	cereCont     bool
	autentificat bool
	sistem       int
}

// NouaPaza pornește cu rezultatul verificării de la deschidere.
func NouaPaza(r Rezultat) *Paza { return &Paza{rez: r} }

// Rezultat întoarce ultima verificare.
func (p *Paza) Rezultat() Rezultat {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.rez
}

// Blocata spune dacă scrierile sunt refuzate.
func (p *Paza) Blocata() bool {
	if p == nil {
		return false
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.sistem > 0 {
		return false
	}
	if p.cereCont && !p.autentificat {
		return true
	}
	return p.rez.Stare == StareIncident && !p.confirmare
}

// CereAutentificare activează refuzul scrierilor până la autentificare.
func (p *Paza) CereAutentificare(da bool) {
	if p == nil {
		return
	}
	p.mu.Lock()
	p.cereCont = da
	p.mu.Unlock()
}

// SeteazaAutentificat marchează dacă un cont e autentificat.
func (p *Paza) SeteazaAutentificat(da bool) {
	if p == nil {
		return
	}
	p.mu.Lock()
	p.autentificat = da
	p.mu.Unlock()
}

// Sistem rulează f cu scrierile permise, pentru operațiile de autentificare, care trebuie
// să scrie (încercări, jurnal) și înainte de autentificare sau în timpul unui incident.
func (p *Paza) Sistem(f func() error) error {
	if p == nil {
		return f()
	}
	p.mu.Lock()
	p.sistem++
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		p.sistem--
		p.mu.Unlock()
	}()
	return f()
}

func (p *Paza) seteaza(r Rezultat) {
	p.mu.Lock()
	p.rez = r
	p.mu.Unlock()
}

func (p *Paza) permiteConfirmarea(da bool) {
	p.mu.Lock()
	p.confirmare = da
	p.mu.Unlock()
}

// Cheie derivă cheia de audit din cheia de date. Nu se păstrează nicăieri.
func Cheie(cheieDate []byte) []byte {
	m := hmac.New(sha256.New, cheieDate)
	m.Write([]byte("pontix/audit/v1"))
	return m.Sum(nil)
}

// PregatesteConexiune înregistrează funcțiile SQL ale auditului pe conexiune și, dacă
// paza nu e nil și baza are deja tabelul sigilii, creează trigger-ele TEMP.
func PregatesteConexiune(c *sqlite3.Conn, cheieDate []byte, paza *Paza) error {
	if err := inregistreazaFunctii(c, Cheie(cheieDate), paza); err != nil {
		return err
	}
	if paza == nil {
		return nil
	}
	tabele, err := tabeleConexiune(c)
	if err != nil {
		return err
	}
	if !contine(tabele, tabelSigilii) {
		return nil
	}
	// REPLACE șterge rânduri fără trigger-ele de ștergere, dacă recursive_triggers e oprit.
	if err := c.Exec(`PRAGMA recursive_triggers=ON`); err != nil {
		return err
	}
	for _, t := range tabele {
		if t == tabelSigilii {
			continue
		}
		coloane, err := coloaneConexiune(c, t)
		if err != nil {
			return err
		}
		if err := c.Exec(declansatori(t, coloane)); err != nil {
			return fmt.Errorf("audit: trigger-e pentru %s: %w", t, err)
		}
	}
	return nil
}

// Funcțiile sunt DIRECTONLY: pot fi apelate din trigger-ele TEMP ale aplicației, dar nu
// dintr-un trigger sau o vedere salvată în fișier de altcineva.
const flagFunctii = sqlite3.DETERMINISTIC | sqlite3.DIRECTONLY

func inregistreazaFunctii(c *sqlite3.Conn, cheie []byte, paza *Paza) error {
	hashRand := func(prefix string, arg []sqlite3.Value) []byte {
		m := hmac.New(sha256.New, cheie)
		m.Write([]byte(prefix))
		for _, v := range arg {
			scrieValoare(m, v)
		}
		return m.Sum(nil)
	}
	if err := c.CreateFunction("pontix_h", -1, flagFunctii, func(ctx sqlite3.Context, arg ...sqlite3.Value) {
		ctx.ResultBlob(hashRand("rand\x00", arg))
	}); err != nil {
		return err
	}
	if err := c.CreateFunction("pontix_lant", -1, flagFunctii, func(ctx sqlite3.Context, arg ...sqlite3.Value) {
		ctx.ResultBlob(hashRand("lant\x00", arg))
	}); err != nil {
		return err
	}
	if err := c.CreateFunction("pontix_adauga", 2, flagFunctii, func(ctx sqlite3.Context, arg ...sqlite3.Value) {
		ctx.ResultBlob(aduna(arg[0].RawBlob(), arg[1].RawBlob(), false))
	}); err != nil {
		return err
	}
	if err := c.CreateFunction("pontix_scade", 2, flagFunctii, func(ctx sqlite3.Context, arg ...sqlite3.Value) {
		ctx.ResultBlob(aduna(arg[0].RawBlob(), arg[1].RawBlob(), true))
	}); err != nil {
		return err
	}
	return c.CreateFunction("pontix_permite", 0, sqlite3.DIRECTONLY, func(ctx sqlite3.Context, arg ...sqlite3.Value) {
		if paza.Blocata() {
			ctx.ResultError(errors.New(CodEditareBlocata))
			return
		}
		ctx.ResultNull()
	})
}

// scrieValoare codifică o valoare SQL fără ambiguitate: tipul, lungimea, octeții.
func scrieValoare(h hash.Hash, v sqlite3.Value) {
	var b [9]byte
	switch v.Type() {
	case sqlite3.INTEGER:
		b[0] = 'i'
		binary.BigEndian.PutUint64(b[1:], uint64(v.Int64()))
		h.Write(b[:])
	case sqlite3.FLOAT:
		b[0] = 'f'
		binary.BigEndian.PutUint64(b[1:], math.Float64bits(v.Float()))
		h.Write(b[:])
	case sqlite3.TEXT:
		scrieOcteti(h, 't', v.RawText())
	case sqlite3.BLOB:
		scrieOcteti(h, 'b', v.RawBlob())
	default:
		h.Write([]byte{'n'})
	}
}

func scrieOcteti(h hash.Hash, tip byte, o []byte) {
	var b [9]byte
	b[0] = tip
	binary.BigEndian.PutUint64(b[1:], uint64(len(o)))
	h.Write(b[:])
	h.Write(o)
}

// aduna întoarce a ± b modulo 2^256. O valoare lipsă sau de altă lungime contează ca zero.
func aduna(a, b []byte, scade bool) []byte {
	rez := make([]byte, lungimeSigiliu)
	if len(a) == lungimeSigiliu {
		copy(rez, a)
	}
	if len(b) != lungimeSigiliu {
		return rez
	}
	transport := 0
	for i := lungimeSigiliu - 1; i >= 0; i-- {
		var x int
		if scade {
			x = int(rez[i]) - int(b[i]) - transport
			transport = 0
			if x < 0 {
				x += 256
				transport = 1
			}
		} else {
			x = int(rez[i]) + int(b[i]) + transport
			transport = x >> 8
			x &= 0xff
		}
		rez[i] = byte(x)
	}
	return rez
}

func cita(id string) string   { return `"` + strings.ReplaceAll(id, `"`, `""`) + `"` }
func literal(s string) string { return `'` + strings.ReplaceAll(s, `'`, `''`) + `'` }

func lista(prefix string, coloane []string) string {
	parti := make([]string, len(coloane))
	for i, c := range coloane {
		parti[i] = prefix + cita(c)
	}
	return strings.Join(parti, ", ")
}

// declansatori întoarce trigger-ele TEMP pentru un tabel: sigiliul la inserare, modificare
// și ștergere și, pentru jurnal, mac-ul înlănțuit al rândului nou.
func declansatori(tabel string, coloane []string) string {
	t, l := cita(tabel), literal(tabel)
	h := func(rand string) string { return "pontix_h(" + l + ", " + lista(rand+".", coloane) + ")" }
	nume := func(op string) string { return cita("pontix_audit_" + tabel + "_" + op) }
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TEMP TRIGGER IF NOT EXISTS %s AFTER INSERT ON main.%s BEGIN SELECT pontix_permite(); UPDATE main.sigilii SET valoare = pontix_adauga(valoare, %s) WHERE tabel = %s; END;\n", nume("i"), t, h("NEW"), l)
	fmt.Fprintf(&b, "CREATE TEMP TRIGGER IF NOT EXISTS %s AFTER UPDATE ON main.%s BEGIN SELECT pontix_permite(); UPDATE main.sigilii SET valoare = pontix_adauga(pontix_scade(valoare, %s), %s) WHERE tabel = %s; END;\n", nume("u"), t, h("OLD"), h("NEW"), l)
	fmt.Fprintf(&b, "CREATE TEMP TRIGGER IF NOT EXISTS %s AFTER DELETE ON main.%s BEGIN SELECT pontix_permite(); UPDATE main.sigilii SET valoare = pontix_scade(valoare, %s) WHERE tabel = %s; END;\n", nume("d"), t, h("OLD"), l)
	if tabel == tabelJurnal {
		fmt.Fprintf(&b, "CREATE TEMP TRIGGER IF NOT EXISTS %s AFTER INSERT ON main.%s BEGIN UPDATE main.%s SET mac = %s WHERE id = NEW.id; END;\n",
			nume("lant"), t, t, expresieLant("NEW.id", lista("NEW.", faraMac(coloane))))
	}
	return b.String()
}

// expresieLant: mac-ul rândului cu id-ul dat, din mac-ul rândului anterior și conținutul lui.
func expresieLant(id, valori string) string {
	return "pontix_lant((SELECT j.mac FROM main.jurnal_modificari j WHERE j.id < " + id + " ORDER BY j.id DESC LIMIT 1), " + valori + ")"
}

func faraMac(coloane []string) []string {
	var rez []string
	for _, c := range coloane {
		if c != "mac" {
			rez = append(rez, c)
		}
	}
	return rez
}

func contine(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

const interogareTabele = `SELECT name FROM sqlite_schema WHERE type = 'table' AND name NOT LIKE 'sqlite\_%' ESCAPE '\' ORDER BY name`

func tabeleConexiune(c *sqlite3.Conn) ([]string, error) {
	return coloanaText(c, interogareTabele)
}

func coloaneConexiune(c *sqlite3.Conn, tabel string) ([]string, error) {
	return coloanaText(c, `SELECT name FROM pragma_table_info(`+literal(tabel)+`) ORDER BY cid`)
}

func coloanaText(c *sqlite3.Conn, interogare string) ([]string, error) {
	st, _, err := c.Prepare(interogare)
	if err != nil {
		return nil, err
	}
	defer st.Close()
	var rez []string
	for st.Step() {
		rez = append(rez, st.ColumnText(0))
	}
	return rez, st.Err()
}

// Interogator e *sql.DB sau *sql.Tx.
type Interogator interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func textColoana(ctx context.Context, q Interogator, interogare string, args ...any) ([]string, error) {
	r, err := q.QueryContext(ctx, interogare, args...)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var rez []string
	for r.Next() {
		var s string
		if err := r.Scan(&s); err != nil {
			return nil, err
		}
		rez = append(rez, s)
	}
	return rez, r.Err()
}

func coloane(ctx context.Context, q Interogator, tabel string) ([]string, error) {
	return textColoana(ctx, q, `SELECT name FROM pragma_table_info(?) ORDER BY cid`, tabel)
}

// sigiliuCalculat recalculează sigiliul unui tabel din conținutul lui curent.
func sigiliuCalculat(ctx context.Context, q Interogator, tabel string) ([]byte, error) {
	var interogare string
	if tabel == pseudoSchema {
		// rootpage se schimbă la VACUUM, deci nu intră în sigiliu.
		interogare = `SELECT pontix_h('sqlite_schema', type, name, tbl_name, sql) FROM sqlite_schema WHERE name NOT LIKE 'sqlite\_%' ESCAPE '\'`
	} else {
		col, err := coloane(ctx, q, tabel)
		if err != nil {
			return nil, err
		}
		interogare = "SELECT pontix_h(" + literal(tabel) + ", " + lista("", col) + ") FROM main." + cita(tabel)
	}
	r, err := q.QueryContext(ctx, interogare)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	suma := make([]byte, lungimeSigiliu)
	for r.Next() {
		var h []byte
		if err := r.Scan(&h); err != nil {
			return nil, err
		}
		suma = aduna(suma, h, false)
	}
	return suma, r.Err()
}

func tabeleSigilate(ctx context.Context, q Interogator) ([]string, error) {
	t, err := textColoana(ctx, q, interogareTabele)
	if err != nil {
		return nil, err
	}
	rez := []string{pseudoSchema}
	for _, x := range t {
		if x != tabelSigilii {
			rez = append(rez, x)
		}
	}
	sort.Strings(rez)
	return rez, nil
}

func areSigilii(ctx context.Context, q Interogator) (bool, error) {
	var n int
	err := q.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name = 'sigilii'`).Scan(&n)
	return n > 0, err
}

// primulRandRupt întoarce primul rând din jurnal al cărui mac nu corespunde; 0 dacă lanțul e întreg.
func primulRandRupt(ctx context.Context, q Interogator) (int64, error) {
	col, err := coloane(ctx, q, tabelJurnal)
	if err != nil {
		return 0, err
	}
	r, err := q.QueryContext(ctx, `SELECT id, mac, pontix_lant(lag(mac) OVER (ORDER BY id), `+lista("", faraMac(col))+`) FROM main.jurnal_modificari ORDER BY id`)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	for r.Next() {
		var id int64
		var mac, asteptat []byte
		if err := r.Scan(&id, &mac, &asteptat); err != nil {
			return 0, err
		}
		if !hmac.Equal(mac, asteptat) {
			return id, nil
		}
	}
	return 0, r.Err()
}

// Verifica recalculează toate sigiliile și lanțul jurnalului și le compară cu cele păstrate.
func Verifica(ctx context.Context, q Interogator) (rez Rezultat, err error) {
	start := time.Now()
	rez = Rezultat{Stare: StareOK, Tabele: []string{}, VerificatLa: start.UTC().Format(time.RFC3339)}
	defer func() { rez.DurataMs = time.Since(start).Milliseconds() }()
	are, err := areSigilii(ctx, q)
	if err != nil {
		return rez, err
	}
	if !are {
		rez.Stare = StareInactiv
		return rez, nil
	}
	pastrate := map[string][]byte{}
	r, err := q.QueryContext(ctx, `SELECT tabel, valoare FROM main.sigilii`)
	if err != nil {
		return rez, err
	}
	for r.Next() {
		var t string
		var v []byte
		if err := r.Scan(&t, &v); err != nil {
			r.Close()
			return rez, err
		}
		pastrate[t] = v
	}
	r.Close()
	if err := r.Err(); err != nil {
		return rez, err
	}
	if len(pastrate) == 0 {
		rez.Stare, rez.Tabele = StareIncident, []string{tabelSigilii}
		return rez, nil
	}
	tabele, err := tabeleSigilate(ctx, q)
	if err != nil {
		return rez, err
	}
	for _, t := range tabele {
		calculat, err := sigiliuCalculat(ctx, q, t)
		if err != nil {
			return rez, err
		}
		if v, ok := pastrate[t]; !ok || !hmac.Equal(v, calculat) {
			rez.Tabele = append(rez.Tabele, t)
		}
		delete(pastrate, t)
	}
	for t := range pastrate {
		rez.Tabele = append(rez.Tabele, t) // tabel sigilat care nu mai există
	}
	sort.Strings(rez.Tabele)
	if contine(tabele, tabelJurnal) {
		if rez.RandRupt, err = primulRandRupt(ctx, q); err != nil {
			return rez, err
		}
	}
	if len(rez.Tabele) > 0 || rez.RandRupt > 0 {
		rez.Stare = StareIncident
	}
	return rez, nil
}

// resigileaza reface lanțul jurnalului de la primul rând rupt și recalculează toate sigiliile.
func resigileaza(ctx context.Context, tx *sql.Tx) error {
	rupt, err := primulRandRupt(ctx, tx)
	if err != nil {
		return err
	}
	if rupt > 0 {
		col, err := coloane(ctx, tx, tabelJurnal)
		if err != nil {
			return err
		}
		ids, err := randuriDeLa(ctx, tx, rupt)
		if err != nil {
			return err
		}
		actualizare := `UPDATE main.jurnal_modificari SET mac = ` + expresieLant("?1", lista("", faraMac(col))) + ` WHERE id = ?1`
		for _, id := range ids {
			if _, err := tx.ExecContext(ctx, actualizare, id); err != nil {
				return err
			}
		}
	}
	tabele, err := tabeleSigilate(ctx, tx)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM main.sigilii`); err != nil {
		return err
	}
	for _, t := range tabele {
		v, err := sigiliuCalculat(ctx, tx, t)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO main.sigilii (tabel, valoare) VALUES (?, ?)`, t, v); err != nil {
			return err
		}
	}
	return nil
}

func randuriDeLa(ctx context.Context, tx *sql.Tx, id int64) ([]int64, error) {
	r, err := tx.QueryContext(ctx, `SELECT id FROM main.jurnal_modificari WHERE id >= ? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var rez []int64
	for r.Next() {
		var x int64
		if err := r.Scan(&x); err != nil {
			return nil, err
		}
		rez = append(rez, x)
	}
	return rez, r.Err()
}

func scrieJurnal(ctx context.Context, tx *sql.Tx, rand string, vechi, nou any) error {
	text := func(v any) (*string, error) {
		if v == nil {
			return nil, nil
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		s := string(b)
		return &s, nil
	}
	v, err := text(vechi)
	if err != nil {
		return err
	}
	n, err := text(nou)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO main.jurnal_modificari (tabel, rand_id, valoare_veche, valoare_noua, cand, cine) VALUES ('audit', ?, ?, ?, ?, 'local')`,
		rand, v, n, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}

func inTranzactie(ctx context.Context, db *sql.DB, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := f(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// Initializeaza face prima sigilare, imediat după migrarea care introduce sigiliile.
func Initializeaza(ctx context.Context, db *sql.DB) error {
	return inTranzactie(ctx, db, func(tx *sql.Tx) error {
		if err := scrieJurnal(ctx, tx, "initializare", nil, map[string]string{"decizie": "D044"}); err != nil {
			return err
		}
		return resigileaza(ctx, tx)
	})
}

// DupaMigrare notează migrarea în jurnal și resigilează. Se apelează doar dacă
// verificarea dinaintea migrării a trecut.
func DupaMigrare(ctx context.Context, db *sql.DB, deLa, panaLa int64) error {
	return inTranzactie(ctx, db, func(tx *sql.Tx) error {
		if err := scrieJurnal(ctx, tx, "migrare", map[string]int64{"versiune": deLa}, map[string]int64{"versiune": panaLa}); err != nil {
			return err
		}
		return resigileaza(ctx, tx)
	})
}

// NoteazaEveniment scrie în jurnal un eveniment de sistem (de ex. restaurarea unui backup)
// și resigilează. Se apelează doar pe o bază verificată.
func NoteazaEveniment(ctx context.Context, db *sql.DB, rand string, vechi, nou any) error {
	return inTranzactie(ctx, db, func(tx *sql.Tx) error {
		if err := scrieJurnal(ctx, tx, rand, vechi, nou); err != nil {
			return err
		}
		return resigileaza(ctx, tx)
	})
}

// Confirma acceptă starea curentă a bazei după un incident: scrie în jurnal ce s-a
// detectat și motivul, resigilează și deblochează editarea.
func Confirma(ctx context.Context, db *sql.DB, paza *Paza, motiv string) (Rezultat, error) {
	motiv = strings.TrimSpace(motiv)
	if motiv == "" {
		return paza.Rezultat(), &i18n.Eroare{Cod: "audit.motiv_obligatoriu"}
	}
	incident := paza.Rezultat()
	if incident.Stare != StareIncident {
		return incident, nil
	}
	paza.permiteConfirmarea(true)
	defer paza.permiteConfirmarea(false)
	err := inTranzactie(ctx, db, func(tx *sql.Tx) error {
		if err := scrieJurnal(ctx, tx, "incident", incident, map[string]string{"motiv": motiv}); err != nil {
			return err
		}
		return resigileaza(ctx, tx)
	})
	if err != nil {
		return incident, err
	}
	rez, err := Verifica(ctx, db)
	if err != nil {
		return incident, err
	}
	paza.seteaza(rez)
	return rez, nil
}

// EsteEditareBlocata recunoaște eroarea dată de trigger-e cât timp există un incident.
func EsteEditareBlocata(err error) bool {
	return err != nil && strings.Contains(err.Error(), CodEditareBlocata)
}
