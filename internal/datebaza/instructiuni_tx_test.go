package datebaza

import (
	"context"
	"strings"
	"testing"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

// I066: in aceeasi tranzactie o instructiune se pregateste o singura data; erorile
// raman aceleasi, iar Rollback anuleaza scrierile facute prin instructiunile pregatite.
func TestTxPregatitRefolosesteInstructiunile(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	p := nouTxPregatit(tx)
	const ins = `INSERT INTO setari(cheie,valoare,modificat_la) VALUES(?,?,'x')`
	for _, k := range []string{"proba1", "proba2", "proba3"} {
		if _, err := p.ExecContext(ctx, ins, k, "v"); err != nil {
			t.Fatal(err)
		}
	}
	var n int
	for i := 0; i < 2; i++ {
		if err := p.QueryRowContext(ctx, `SELECT count(*) FROM setari WHERE cheie LIKE 'proba%'`).Scan(&n); err != nil {
			t.Fatal(err)
		}
	}
	if n != 3 || len(p.cache) != 2 {
		t.Fatalf("randuri %d, instructiuni pregatite %d (astept 3 si 2)", n, len(p.cache))
	}
	// QueryContext nu foloseste cache-ul: randurile deschise nu blocheaza alte apeluri.
	rows, err := p.QueryContext(ctx, `SELECT cheie FROM setari WHERE cheie LIKE 'proba%'`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			t.Fatal(err)
		}
		if err := p.QueryRowContext(ctx, `SELECT count(*) FROM setari WHERE cheie=?`, k).Scan(&n); err != nil || n != 1 {
			t.Fatalf("interogare imbricata pentru %s: %d %v", k, n, err)
		}
	}
	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}
	if err := p.QueryRowContext(ctx, `SELECT nu_exista FROM tabel_inexistent`).Scan(&n); err == nil {
		t.Fatal("eroarea de pregatire trebuia sa ajunga la Scan")
	}
	if _, err := p.ExecContext(ctx, `INSERT INTO tabel_inexistent VALUES(1)`); err == nil {
		t.Fatal("eroarea de pregatire trebuia intoarsa de Exec")
	}
	// Sqlc functioneaza peste adaptor.
	if _, err := acces.New(p).ListaSchimburi(ctx); err != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	if err := m.db.QueryRowContext(ctx, `SELECT count(*) FROM setari WHERE cheie LIKE 'proba%'`).Scan(&n); err != nil || n != 0 {
		t.Fatalf("Rollback nu a anulat scrierile: %d %v", n, err)
	}
}

// Multimea angajatilor cu programe multiple, calculata o data pe lot, da acelasi
// raspuns ca vechea parcurgere a tuturor ciornelor pentru fiecare celula.
func TestMultipliCalculatiODataPeLot(t *testing.T) {
	s := sursePlanificare{
		ciorne: map[string]acces.PontajCiorne{
			"a/2026-10-01": {AngajatID: "a", Data: "2026-10-01", Valoare: "Z7"},
			"b/2026-10-01": {AngajatID: "b", Data: "2026-10-01", Valoare: "Z7+S7"},
		},
		coduri: map[string]calcul.ProgramPrecompletare{},
	}
	intervaleLotPontaj(&s, []CelulaSchita{{AngajatID: "a", Data: "2026-10-02", Valoare: "N12+Z7"}, {AngajatID: "c", Data: "2026-10-02", Valoare: "Z7"}})
	for _, id := range []string{"a", "b", "c"} {
		vechi := false
		for _, v := range s.ciorne {
			vechi = vechi || (v.AngajatID == id && strings.Contains(v.Valoare, "+"))
		}
		if s.multipli[id] != vechi {
			t.Errorf("%s: multipli %v, parcurgerea veche %v", id, s.multipli[id], vechi)
		}
	}
	if !s.multipli["a"] || !s.multipli["b"] || s.multipli["c"] {
		t.Fatalf("multipli: %v", s.multipli)
	}
}
