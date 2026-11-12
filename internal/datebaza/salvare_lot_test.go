package datebaza

import (
	"context"
	"fmt"
	"testing"
)

// F08-T17 (I066): salvarea in loturi de cate 200 de celule pastreaza reviziile, jurnalul pe
// celula si verificarea conflictelor, si peste mai multe loturi (15 x 30 = 450 de celule).
func TestSalvarePontajInLoturi(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	for i := 0; i < 15; i++ {
		if _, err := m.db.Exec(`INSERT INTO angajati (id, nume, prenume, marca, creat_la, modificat_la) VALUES (?, ?, 'Test', ?, 'x', 'x')`, fmt.Sprintf("lot-%02d", i), fmt.Sprintf("Lot%02d", i), fmt.Sprintf("L-%02d", i)); err != nil {
			t.Fatal(err)
		}
	}
	celule := func(valoare string, revizii map[string]string) []CelulaSchita {
		rez := []CelulaSchita{}
		for i := 0; i < 15; i++ {
			for zi := 1; zi <= 30; zi++ {
				id, data := fmt.Sprintf("lot-%02d", i), fmt.Sprintf("2026-11-%02d", zi)
				rez = append(rez, CelulaSchita{AngajatID: id, Data: data, Valoare: valoare, Revizie: revizii[id+"/"+data]})
			}
		}
		return rez
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 11, celule("8", nil)); err != nil {
		t.Fatal(err)
	}
	revizii := func(asteptata string) map[string]string {
		t.Helper()
		s, err := m.serviciu.CitesteSchita(ctx, 2026, 11)
		if err != nil {
			t.Fatal(err)
		}
		rez := map[string]string{}
		for _, c := range s.Celule {
			if c.Revizie != asteptata {
				t.Fatalf("%s/%s: revizia %s, asteptata %s", c.AngajatID, c.Data, c.Revizie, asteptata)
			}
			rez[c.AngajatID+"/"+c.Data] = c.Revizie
		}
		if len(rez) != 450 {
			t.Fatalf("%d celule", len(rez))
		}
		return rez
	}
	r := revizii("1")
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 11, celule("7", r)); err != nil {
		t.Fatal(err)
	}
	revizii("2")
	// Revizia veche (1) dupa actualizare: conflict, nimic scris.
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 11, celule("6", r)), "pontaj.conflict")
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel='pontaj_ciorne'`).Scan(&n); err != nil || n != 900 {
		t.Fatal("jurnal", n, err)
	}
	var vechi, nou string
	if err := m.db.QueryRow(`SELECT coalesce(valoare_veche,''), valoare_noua FROM jurnal_modificari WHERE tabel='pontaj_ciorne' AND rand_id='lot-14/2026-11-30' ORDER BY id DESC LIMIT 1`).Scan(&vechi, &nou); err != nil {
		t.Fatal(err)
	}
	if vechi == "" || nou == "" {
		t.Fatal("jurnal fara valori", vechi, nou)
	}
}
