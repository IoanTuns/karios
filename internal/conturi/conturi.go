// Package conturi gestionează conturile aplicației (F15, D008 punctul 3, D045).
//
// Parola unui cont deschide copia „P:<id>” din seif, iar răspunsurile la întrebările de
// recuperare copia „R:<id>”; baza nu păstrează niciun hash. Tabelul conturi ține numele,
// rolul, întrebările alese și blocarea, iar sesiune_locala ultimul cont pentru auto-login.
package conturi

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
	"unicode"

	"pontix/internal/audit"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"pontix/internal/seif"
)

// Roluri (D045 punctul 3).
const (
	RolAdministrator = "ADMINISTRATOR"
	RolOperator      = "OPERATOR"
)

// Intrebari sunt codurile întrebărilor de recuperare; textele sunt în cataloage („cont.<cod>”).
var Intrebari = []string{"intrebare_1", "intrebare_2", "intrebare_3", "intrebare_4", "intrebare_5", "intrebare_6", "intrebare_7", "intrebare_8", "intrebare_9", "intrebare_10"}

// NumarIntrebari e câte întrebări diferite alege fiecare cont.
const NumarIntrebari = 3

// pragBlocare: de la a câta încercare greșită consecutivă se blochează contul.
const pragBlocare = 5

var duratiBlocare = []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 120 * time.Minute}

// Cont e un cont, fără secrete.
type Cont struct {
	ID           string   `json:"id"`
	Nume         string   `json:"nume"`
	NumeAfisat   string   `json:"nume_afisat"`
	Rol          string   `json:"rol"`
	Activ        bool     `json:"activ"`
	Intrebari    []string `json:"intrebari"`
	BlocatPanaLa string   `json:"blocat_pana_la"`
}

// Administrator spune dacă contul are rolul de administrator.
func (c Cont) Administrator() bool { return c.Rol == RolAdministrator }

// CerereCont descrie un cont nou.
type CerereCont struct {
	Nume       string   `json:"nume"`
	NumeAfisat string   `json:"nume_afisat"`
	Rol        string   `json:"rol"`
	Parola     string   `json:"parola"`
	Intrebari  []string `json:"intrebari"`
	Raspunsuri []string `json:"raspunsuri"`
}

// Identitate e utilizatorul Windows și calculatorul, pentru auto-login.
type Identitate struct {
	SID        string
	Calculator string
}

// Serviciu lucrează cu baza deschisă și cu seiful din același folder de date.
type Serviciu struct {
	db   *sql.DB
	seif *seif.Seif
	paza *audit.Paza
	acum func() time.Time
}

// Nou primește seiful deschis la pornire (D054): parolele se verifică pe copiile lui, fără să
// redeschidă fișierul și fără nicio cheie din cod.
func Nou(db *sql.DB, sf *seif.Seif, paza *audit.Paza) *Serviciu {
	return &Serviciu{db: db, seif: sf, paza: paza, acum: time.Now}
}

func eroare(cod string, param ...string) error {
	e := &i18n.Eroare{Cod: cod}
	if len(param) == 2 {
		e.Param = map[string]string{param[0]: param[1]}
	}
	return e
}

// ValideazaParola aplică D045 punctul 2: cel puțin 8 caractere, o literă și o cifră.
func ValideazaParola(p string) error {
	var litera, cifra bool
	n := 0
	for _, r := range p {
		n++
		litera = litera || unicode.IsLetter(r)
		cifra = cifra || unicode.IsDigit(r)
	}
	if n < 8 || !litera || !cifra {
		return eroare("cont.parola_slaba")
	}
	return nil
}

var faraDiacritice = strings.NewReplacer("ă", "a", "â", "a", "î", "i", "ș", "s", "ş", "s", "ț", "t", "ţ", "t")

// normalizeazaRaspuns: litere mici, fără diacritice, spații strânse (D008 punctul 3).
func normalizeazaRaspuns(s string) string {
	return strings.Join(strings.Fields(faraDiacritice.Replace(strings.ToLower(s))), " ")
}

// secretRaspunsuri leagă cele trei răspunsuri normalizate, în ordinea întrebărilor.
func secretRaspunsuri(raspunsuri []string) string {
	n := make([]string, len(raspunsuri))
	for i, r := range raspunsuri {
		n[i] = normalizeazaRaspuns(r)
	}
	return strings.Join(n, "\x1f")
}

func valideazaIntrebari(intrebari, raspunsuri []string) error {
	if len(intrebari) != NumarIntrebari || len(raspunsuri) != NumarIntrebari {
		return eroare("cont.intrebari_invalide")
	}
	vazute := map[string]bool{}
	for i, q := range intrebari {
		valid := false
		for _, x := range Intrebari {
			valid = valid || x == q
		}
		if !valid || vazute[q] || normalizeazaRaspuns(raspunsuri[i]) == "" {
			return eroare("cont.intrebari_invalide")
		}
		vazute[q] = true
	}
	return nil
}

func valideazaNume(nume string) (string, error) {
	nume = strings.TrimSpace(nume)
	if n := len([]rune(nume)); n < 2 || n > 40 {
		return "", eroare("cont.nume_invalid")
	}
	return nume, nil
}

func idNou() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}

func (s *Serviciu) moment() string { return s.acum().UTC().Format(time.RFC3339Nano) }

// sistem rulează f ca operație de sistem: scrierile sunt permise și înainte de autentificare.
func (s *Serviciu) sistem(f func() error) error { return s.paza.Sistem(f) }

func (s *Serviciu) tranzactie(ctx context.Context, f func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := f(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Serviciu) jurnal(ctx context.Context, tx *sql.Tx, contID, camp string, vechi, nou any) error {
	text := func(v any) (*string, error) {
		if v == nil {
			return nil, nil
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		x := string(b)
		return &x, nil
	}
	v, err := text(vechi)
	if err != nil {
		return err
	}
	n, err := text(nou)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO jurnal_modificari (tabel, rand_id, camp, valoare_veche, valoare_noua, cand, cine) VALUES ('conturi', ?, ?, ?, ?, ?, ?)`,
		contID, camp, v, n, s.moment(), datebaza.Actor(ctx))
	return err
}

const coloaneCont = `id, nume, nume_afisat, rol, activ, intrebari, blocat_pana_la`

type scanner interface{ Scan(...any) error }

func scanCont(r scanner) (Cont, error) {
	var c Cont
	var activ int
	var intrebari string
	if err := r.Scan(&c.ID, &c.Nume, &c.NumeAfisat, &c.Rol, &activ, &intrebari, &c.BlocatPanaLa); err != nil {
		return c, err
	}
	c.Activ = activ == 1
	if err := json.Unmarshal([]byte(intrebari), &c.Intrebari); err != nil {
		return c, err
	}
	return c, nil
}

// Numar întoarce câte conturi există; 0 înseamnă că trebuie creat primul administrator.
func (s *Serviciu) Numar(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM conturi`).Scan(&n)
	return n, err
}

// Lista întoarce toate conturile, după nume.
func (s *Serviciu) Lista(ctx context.Context) ([]Cont, error) {
	r, err := s.db.QueryContext(ctx, `SELECT `+coloaneCont+` FROM conturi ORDER BY nume COLLATE NOCASE`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	rez := []Cont{}
	for r.Next() {
		c, err := scanCont(r)
		if err != nil {
			return nil, err
		}
		rez = append(rez, c)
	}
	return rez, r.Err()
}

func (s *Serviciu) dupaNume(ctx context.Context, q interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}, nume string) (Cont, error) {
	c, err := scanCont(q.QueryRowContext(ctx, `SELECT `+coloaneCont+` FROM conturi WHERE nume = ? COLLATE NOCASE`, strings.TrimSpace(nume)))
	if errors.Is(err, sql.ErrNoRows) {
		return c, eroare("cont.date_gresite")
	}
	return c, err
}

// DupaID întoarce contul cu id-ul dat.
func (s *Serviciu) DupaID(ctx context.Context, id string) (Cont, error) {
	c, err := scanCont(s.db.QueryRowContext(ctx, `SELECT `+coloaneCont+` FROM conturi WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return c, eroare("cont.inexistent")
	}
	return c, err
}

// Creeaza adaugă un cont. Primul cont e întotdeauna administrator (D045 punctul 4).
func (s *Serviciu) Creeaza(ctx context.Context, c CerereCont) (Cont, error) {
	nume, err := valideazaNume(c.Nume)
	if err != nil {
		return Cont{}, err
	}
	if err := ValideazaParola(c.Parola); err != nil {
		return Cont{}, err
	}
	if err := valideazaIntrebari(c.Intrebari, c.Raspunsuri); err != nil {
		return Cont{}, err
	}
	if c.Rol != RolAdministrator && c.Rol != RolOperator {
		c.Rol = RolOperator
	}
	cont := Cont{ID: idNou(), Nume: nume, NumeAfisat: strings.TrimSpace(c.NumeAfisat), Rol: c.Rol, Activ: true, Intrebari: c.Intrebari}
	err = s.sistem(func() error {
		return s.tranzactie(ctx, func(tx *sql.Tx) error {
			var n int
			if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM conturi`).Scan(&n); err != nil {
				return err
			}
			if n == 0 {
				cont.Rol = RolAdministrator
			}
			var exista int
			if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM conturi WHERE nume = ? COLLATE NOCASE`, nume).Scan(&exista); err != nil {
				return err
			}
			if exista > 0 {
				return eroare("cont.nume_existent")
			}
			intrebari, _ := json.Marshal(cont.Intrebari)
			acum := s.moment()
			if _, err := tx.ExecContext(ctx, `INSERT INTO conturi (id, nume, nume_afisat, rol, activ, intrebari, creat_la, modificat_la) VALUES (?, ?, ?, ?, 1, ?, ?, ?)`,
				cont.ID, cont.Nume, cont.NumeAfisat, cont.Rol, string(intrebari), acum, acum); err != nil {
				return err
			}
			if err := s.jurnal(ctx, tx, cont.ID, "creare", nil, cont); err != nil {
				return err
			}
			return s.scrieSecrete(cont.ID, c.Parola, c.Raspunsuri)
		})
	})
	return cont, err
}

// scrieSecrete scrie în seif copiile P și R ale contului. Se apelează în tranzacție,
// ca un seif nescris să anuleze și contul.
func (s *Serviciu) scrieSecrete(id, parola string, raspunsuri []string) error {
	if parola != "" {
		if err := s.seif.SeteazaSecret(seif.PrefixParola+id, parola); err != nil {
			return err
		}
		_ = s.seif.StergeCopie(seif.PrefixParolaVeche + id)
	}
	if raspunsuri != nil {
		if err := s.seif.SeteazaSecret(seif.PrefixRaspunsuri+id, secretRaspunsuri(raspunsuri)); err != nil {
			return err
		}
		_ = s.seif.StergeCopie(seif.PrefixRaspunsuriVechi + id)
	}
	return nil
}

// verificaSecret spune dacă secretul deschide copia P sau R a contului cu cheia de date curentă.
// Un cont care nu s-a mai autentificat de la trecerea D054 are doar copia veche (P0/R0): secretul
// se verifică pe ea și, dacă e corect, copia se reface pe cheia nouă, iar cea veche se șterge.
func (s *Serviciu) verificaSecret(id, secret string) bool {
	prefix, vechi := seif.PrefixParola, seif.PrefixParolaVeche
	if strings.HasPrefix(id, seif.PrefixRaspunsuri) {
		prefix, vechi = seif.PrefixRaspunsuri, seif.PrefixRaspunsuriVechi
	}
	cont := strings.TrimPrefix(id, prefix)
	if s.seif.AreCopia(prefix + cont) {
		return s.seif.VerificaSecret(prefix+cont, secret)
	}
	if !s.seif.AreCopia(vechi+cont) || !s.seif.VerificaSecret(vechi+cont, secret) {
		return false
	}
	if err := s.seif.SeteazaSecret(prefix+cont, secret); err == nil {
		_ = s.seif.StergeCopie(vechi + cont)
	}
	return true
}

func (s *Serviciu) verificaBlocare(c Cont) error {
	if c.BlocatPanaLa == "" {
		return nil
	}
	pana, err := time.Parse(time.RFC3339Nano, c.BlocatPanaLa)
	if err == nil && s.acum().Before(pana) {
		return eroare("cont.blocat", "pana_la", pana.Format(time.RFC3339))
	}
	return nil
}

// incercareEsuata crește contorul, blochează contul de la prag și notează în jurnal.
func (s *Serviciu) incercareEsuata(ctx context.Context, c Cont, motiv string) error {
	return s.sistem(func() error {
		return s.tranzactie(ctx, func(tx *sql.Tx) error {
			var n int
			if err := tx.QueryRowContext(ctx, `UPDATE conturi SET incercari_esuate = incercari_esuate + 1 WHERE id = ? RETURNING incercari_esuate`, c.ID).Scan(&n); err != nil {
				return err
			}
			blocat := ""
			if n >= pragBlocare {
				d := duratiBlocare[min(n-pragBlocare, len(duratiBlocare)-1)]
				blocat = s.acum().Add(d).UTC().Format(time.RFC3339Nano)
				if _, err := tx.ExecContext(ctx, `UPDATE conturi SET blocat_pana_la = ? WHERE id = ?`, blocat, c.ID); err != nil {
					return err
				}
			}
			return s.jurnal(datebaza.CuActor(ctx, c.Nume), tx, c.ID, motiv, nil, map[string]any{"incercari": n, "blocat_pana_la": blocat})
		})
	})
}

// autentificat resetează contorul, reține sesiunea locală și notează în jurnal.
func (s *Serviciu) autentificat(ctx context.Context, c Cont, id Identitate, camp string) error {
	return s.sistem(func() error {
		return s.tranzactie(ctx, func(tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, `UPDATE conturi SET incercari_esuate = 0, blocat_pana_la = '' WHERE id = ? AND (incercari_esuate <> 0 OR blocat_pana_la <> '')`, c.ID); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO sesiune_locala (id, cont_id, sid, calculator, autologin, modificat_la) VALUES (1, ?, ?, ?, 1, ?)
				ON CONFLICT (id) DO UPDATE SET cont_id = excluded.cont_id, sid = excluded.sid, calculator = excluded.calculator, autologin = 1, modificat_la = excluded.modificat_la`,
				c.ID, id.SID, id.Calculator, s.moment()); err != nil {
				return err
			}
			return s.jurnal(datebaza.CuActor(ctx, c.Nume), tx, c.ID, camp, nil, map[string]string{"calculator": id.Calculator})
		})
	})
}

// Autentifica verifică numele și parola. Un nume inexistent și o parolă greșită dau aceeași
// eroare, ca să nu se poată afla ce conturi există.
func (s *Serviciu) Autentifica(ctx context.Context, nume, parola string, id Identitate) (Cont, error) {
	c, err := s.dupaNume(ctx, s.db, nume)
	if err != nil {
		return Cont{}, err
	}
	if !c.Activ {
		return Cont{}, eroare("cont.date_gresite")
	}
	if err := s.verificaBlocare(c); err != nil {
		return Cont{}, err
	}
	if !s.verificaSecret("P:"+c.ID, parola) {
		if err := s.incercareEsuata(ctx, c, "autentificare_esuata"); err != nil {
			return Cont{}, err
		}
		return Cont{}, eroare("cont.date_gresite")
	}
	if err := s.autentificat(ctx, c, id, "autentificare"); err != nil {
		return Cont{}, err
	}
	c.BlocatPanaLa = ""
	return c, nil
}

// AutoLogin intră în ultimul cont dacă utilizatorul Windows și calculatorul sunt aceleași
// și dacă nu s-a ieșit din cont (D008 punctul 3).
func (s *Serviciu) AutoLogin(ctx context.Context, id Identitate) (Cont, bool, error) {
	var contID, sid, calculator string
	var autologin int
	err := s.db.QueryRowContext(ctx, `SELECT coalesce(cont_id, ''), sid, calculator, autologin FROM sesiune_locala WHERE id = 1`).Scan(&contID, &sid, &calculator, &autologin)
	if errors.Is(err, sql.ErrNoRows) {
		return Cont{}, false, nil
	}
	if err != nil {
		return Cont{}, false, err
	}
	if autologin != 1 || contID == "" || id.SID == "" || sid != id.SID || !strings.EqualFold(calculator, id.Calculator) {
		return Cont{}, false, nil
	}
	c, err := s.DupaID(ctx, contID)
	if err != nil || !c.Activ || s.verificaBlocare(c) != nil {
		return Cont{}, false, nil
	}
	if err := s.autentificat(ctx, c, id, "autentificare_automata"); err != nil {
		return Cont{}, false, err
	}
	return c, true, nil
}

// UltimulNume întoarce numele ultimului cont autentificat aici, pentru precompletarea login-ului.
func (s *Serviciu) UltimulNume(ctx context.Context) string {
	var nume string
	_ = s.db.QueryRowContext(ctx, `SELECT c.nume FROM sesiune_locala l JOIN conturi c ON c.id = l.cont_id WHERE l.id = 1`).Scan(&nume)
	return nume
}

// Iesire oprește auto-login-ul până la următoarea autentificare cu parolă.
func (s *Serviciu) Iesire(ctx context.Context, c Cont) error {
	return s.sistem(func() error {
		return s.tranzactie(ctx, func(tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, `UPDATE sesiune_locala SET autologin = 0, modificat_la = ? WHERE id = 1`, s.moment()); err != nil {
				return err
			}
			return s.jurnal(datebaza.CuActor(ctx, c.Nume), tx, c.ID, "iesire", nil, map[string]bool{"autologin": false})
		})
	})
}

// IntrebariPentru întoarce întrebările contului, pentru recuperarea parolei.
func (s *Serviciu) IntrebariPentru(ctx context.Context, nume string) ([]string, error) {
	c, err := s.dupaNume(ctx, s.db, nume)
	if err != nil {
		return nil, err
	}
	if !c.Activ || len(c.Intrebari) != NumarIntrebari {
		return nil, eroare("cont.fara_recuperare")
	}
	return c.Intrebari, nil
}

// Recupereaza schimbă parola pe baza răspunsurilor și autentifică contul. Răspunsurile
// greșite contează ca încercări greșite.
func (s *Serviciu) Recupereaza(ctx context.Context, nume string, raspunsuri []string, parolaNoua string, id Identitate) (Cont, error) {
	if err := ValideazaParola(parolaNoua); err != nil {
		return Cont{}, err
	}
	c, err := s.dupaNume(ctx, s.db, nume)
	if err != nil {
		return Cont{}, err
	}
	if !c.Activ {
		return Cont{}, eroare("cont.date_gresite")
	}
	if err := s.verificaBlocare(c); err != nil {
		return Cont{}, err
	}
	if len(raspunsuri) != NumarIntrebari || !s.verificaSecret("R:"+c.ID, secretRaspunsuri(raspunsuri)) {
		if err := s.incercareEsuata(ctx, c, "recuperare_esuata"); err != nil {
			return Cont{}, err
		}
		return Cont{}, eroare("cont.raspunsuri_gresite")
	}
	err = s.sistem(func() error {
		return s.tranzactie(ctx, func(tx *sql.Tx) error {
			if err := s.jurnal(datebaza.CuActor(ctx, c.Nume), tx, c.ID, "recuperare_parola", nil, map[string]bool{"parola_schimbata": true}); err != nil {
				return err
			}
			return s.scrieSecrete(c.ID, parolaNoua, nil)
		})
	})
	if err != nil {
		return Cont{}, err
	}
	if err := s.autentificat(ctx, c, id, "autentificare"); err != nil {
		return Cont{}, err
	}
	return c, nil
}

// ReseteazaParola: administratorul dă o parolă nouă unui cont și îl deblochează.
func (s *Serviciu) ReseteazaParola(ctx context.Context, contID, parolaNoua string) error {
	if err := ValideazaParola(parolaNoua); err != nil {
		return err
	}
	return s.tranzactie(ctx, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `UPDATE conturi SET incercari_esuate = 0, blocat_pana_la = '', modificat_la = ? WHERE id = ?`, s.moment(), contID)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return eroare("cont.inexistent")
		}
		if err := s.jurnal(ctx, tx, contID, "resetare_parola", nil, map[string]bool{"parola_schimbata": true}); err != nil {
			return err
		}
		return s.scrieSecrete(contID, parolaNoua, nil)
	})
}

// SchimbaParola: contul își schimbă parola, cu parola veche.
func (s *Serviciu) SchimbaParola(ctx context.Context, c Cont, veche, noua string) error {
	if !s.verificaSecret("P:"+c.ID, veche) {
		return eroare("cont.parola_veche_gresita")
	}
	if err := ValideazaParola(noua); err != nil {
		return err
	}
	return s.tranzactie(ctx, func(tx *sql.Tx) error {
		if err := s.jurnal(ctx, tx, c.ID, "schimbare_parola", nil, map[string]bool{"parola_schimbata": true}); err != nil {
			return err
		}
		return s.scrieSecrete(c.ID, noua, nil)
	})
}

// SchimbaIntrebari: contul își alege alte întrebări și răspunsuri, confirmând cu parola.
func (s *Serviciu) SchimbaIntrebari(ctx context.Context, c Cont, parola string, intrebari, raspunsuri []string) error {
	if !s.verificaSecret("P:"+c.ID, parola) {
		return eroare("cont.parola_veche_gresita")
	}
	if err := valideazaIntrebari(intrebari, raspunsuri); err != nil {
		return err
	}
	return s.tranzactie(ctx, func(tx *sql.Tx) error {
		j, _ := json.Marshal(intrebari)
		if _, err := tx.ExecContext(ctx, `UPDATE conturi SET intrebari = ?, modificat_la = ? WHERE id = ?`, string(j), s.moment(), c.ID); err != nil {
			return err
		}
		if err := s.jurnal(ctx, tx, c.ID, "schimbare_intrebari", c.Intrebari, intrebari); err != nil {
			return err
		}
		return s.scrieSecrete(c.ID, "", raspunsuri)
	})
}

// Modifica schimbă numele afișat, rolul și starea unui cont. Ultimul administrator activ
// nu poate fi retrogradat sau dezactivat.
func (s *Serviciu) Modifica(ctx context.Context, c Cont) (Cont, error) {
	if c.Rol != RolAdministrator && c.Rol != RolOperator {
		return Cont{}, eroare("cont.rol_invalid")
	}
	var rez Cont
	err := s.tranzactie(ctx, func(tx *sql.Tx) error {
		vechi, err := scanCont(tx.QueryRowContext(ctx, `SELECT `+coloaneCont+` FROM conturi WHERE id = ?`, c.ID))
		if errors.Is(err, sql.ErrNoRows) {
			return eroare("cont.inexistent")
		}
		if err != nil {
			return err
		}
		if vechi.Administrator() && vechi.Activ && (c.Rol != RolAdministrator || !c.Activ) {
			var admini int
			if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM conturi WHERE rol = 'ADMINISTRATOR' AND activ = 1`).Scan(&admini); err != nil {
				return err
			}
			if admini <= 1 {
				return eroare("cont.ultimul_administrator")
			}
		}
		activ := 0
		if c.Activ {
			activ = 1
		}
		if _, err := tx.ExecContext(ctx, `UPDATE conturi SET nume_afisat = ?, rol = ?, activ = ?, modificat_la = ? WHERE id = ?`,
			strings.TrimSpace(c.NumeAfisat), c.Rol, activ, s.moment(), c.ID); err != nil {
			return err
		}
		rez = vechi
		rez.NumeAfisat, rez.Rol, rez.Activ = strings.TrimSpace(c.NumeAfisat), c.Rol, c.Activ
		return s.jurnal(ctx, tx, c.ID, "modificare", vechi, rez)
	})
	return rez, err
}
