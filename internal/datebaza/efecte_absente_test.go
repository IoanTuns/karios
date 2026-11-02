package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
	"testing"
)

func TestEfecteAbsenteNormaProgramSiAnulare(t *testing.T) {
	for _, efect := range []string{"NORMA", "PROGRAM", "FARA_ORE"} {
		t.Run(efect, func(t *testing.T) {
			m := pregateste(t)
			ctx := context.Background()
			a := persoanaAbsente(t, m)
			tip := TipAbsenta{Cod: "DEPL", Denumire: i18n.Text{"ro": "Deplasare"}, Categorie: "SPECIAL", Activ: true, EfectOre: efect}
			if e := m.serviciu.SalveazaTipAbsenta(ctx, "", tip); e != nil {
				t.Fatal(e)
			}
			if _, e := m.db.Exec(`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,creat_la,modificat_la) SELECT 'efect',?,'2026-10-05',id,'GARZI','x','x' FROM tipuri_schimb WHERE cod='N12'`, a.ID); e != nil {
				t.Fatal(e)
			}
			z := ZiAbsenta{AngajatID: a.ID, Data: "2026-10-05", Cod: "DEPL", Activ: true, StareAprobare: "APROBATA"}
			before := instantaneuSimulare(t, m)
			if _, e := m.serviciu.PrevizualizeazaAbsente(ctx, []ZiAbsenta{z}, 2026, 10, nil); e != nil {
				t.Fatal(e)
			}
			if before != instantaneuSimulare(t, m) {
				t.Fatal("previzualizarea a scris")
			}
			if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); e != nil {
				t.Fatal(e)
			}
			q := acces.New(m.db)
			r, _, e := CalculeazaPontaj(ctx, q, 2026, 10, nil)
			if e != nil {
				t.Fatal(e)
			}
			n := int64(0)
			if efect == "NORMA" {
				n = 420
			}
			if efect == "PROGRAM" {
				n = 720
			}
			for _, total := range r.Totaluri {
				if total.AngajatID == a.ID {
					asteptat := int64(22 * 420)
					if efect == "FARA_ORE" {
						asteptat -= 420
					}
					if v := total.Valori["ore_lucratoare"]; !v.Definita || v.Sutimi != asteptat*100 {
						t.Fatalf("obligatie: %+v", v)
					}
				}
			}
			var gasit bool
			for _, c := range r.Celule {
				if c.AngajatID == a.ID && c.Data == z.Data {
					gasit = true
					if c.Minute != n || !c.Calculabil {
						t.Fatalf("celula %+v", c)
					}
				}
			}
			if !gasit {
				t.Fatal("celula lipsa")
			}
			// Programul este anulat; editarea aprobarii pastreaza instantaneul.
			z = absenteTest(t, m, 2026).Zile[0]
			z.StareAprobare = "IN_ASTEPTARE"
			if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); e != nil {
				t.Fatal(e)
			}
			z = absenteTest(t, m, 2026).Zile[0]
			z.StareAprobare = "APROBATA"
			if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); e != nil {
				t.Fatal(e)
			}
			tip.EfectOre = "ANTERIOR"
			verificaCod(t, m.serviciu.SalveazaTipAbsenta(ctx, "DEPL", tip), "absente.categorie_folosita")
			z = absenteTest(t, m, 2026).Zile[0]
			z.Activ = false
			if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); e != nil {
				t.Fatal(e)
			}
		})
	}
}
func TestAbsentaProgramLipsaNuScrie(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	tip := TipAbsenta{Cod: "DEPL", Denumire: i18n.Text{"ro": "Deplasare"}, Categorie: "SPECIAL", Activ: true, EfectOre: "PROGRAM"}
	if e := m.serviciu.SalveazaTipAbsenta(ctx, "", tip); e != nil {
		t.Fatal(e)
	}
	before := instantaneuSimulare(t, m)
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "DEPL", Activ: true}}), "absente.ore_lipsa")
	if before != instantaneuSimulare(t, m) {
		t.Fatal("scriere partiala")
	}
}
