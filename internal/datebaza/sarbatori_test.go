package datebaza

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"pontix/internal/i18n"
)

func TestCalendar2026GeneratSiIdempotent(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	c, err := m.serviciu.CitesteCalendar(ctx, 2026)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Sarbatori) != 16 {
		t.Fatalf("date distincte: %d", len(c.Sarbatori))
	}
	var sursa []Sarbatoare
	if err := json.Unmarshal(calendar2026, &sursa); err != nil {
		t.Fatal(err)
	}
	for i, s := range c.Sarbatori {
		if s.Data != sursa[i].Data || s.Denumire.In(i18n.Romana) != sursa[i].Denumire.In(i18n.Romana) || s.Denumire.In(i18n.Engleza) != sursa[i].Denumire.In(i18n.Engleza) {
			t.Fatalf("migrarea și sursa diferă la %s", s.Data)
		}
	}
	var denumire string
	for _, s := range c.Sarbatori {
		if s.Data == "2026-06-01" {
			denumire = s.Denumire.In(i18n.Romana)
		}
	}
	if !strings.Contains(denumire, "Copilului") || !strings.Contains(denumire, "Rusalii") {
		t.Fatal("suprapunerea de la 1 iunie nu este reprezentată")
	}
	if n, err := m.serviciu.CompleteazaSarbatori2026(ctx); err != nil || n != 0 {
		t.Fatalf("import repetat: %d, %v", n, err)
	}
	if c, err := m.serviciu.CitesteCalendar(ctx, 2028); err != nil || len(c.Sarbatori) != 0 || c.ImportDisponibil {
		t.Fatalf("calendar 2028: %+v, %v", c, err)
	}
	if n, err := m.serviciu.CompleteazaSarbatoriAn(ctx, 2025); err != nil || n == 0 {
		t.Fatalf("import 2025: %d, %v", n, err)
	}
	if n, err := m.serviciu.CompleteazaSarbatoriAn(ctx, 2027); err != nil || n == 0 {
		t.Fatalf("import 2027: %d, %v", n, err)
	}
}

func TestEditareSarbatoriPastreazaDateleSiImportulNuSuprascrie(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaSarbatoare(ctx, "2026-01-01", Sarbatoare{Data: "2026-01-01", Denumire: i18n.Text{i18n.Romana: " Denumire personalizată "}}); err != nil {
		t.Fatal(err)
	}
	if n, err := m.serviciu.CompleteazaSarbatori2026(ctx); err != nil || n != 0 {
		t.Fatalf("completare: %d, %v", n, err)
	}
	c, err := m.serviciu.CitesteCalendar(ctx, 2026)
	if err != nil || c.Sarbatori[0].Denumire.In(i18n.Romana) != "Denumire personalizată" {
		t.Fatalf("suprascriere: %+v, %v", c, err)
	}
	if _, err := m.db.Exec(`UPDATE sarbatori SET creat_la = 'initial' WHERE data = '2026-01-01'`); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaSarbatoare(ctx, "2026-01-01", Sarbatoare{Data: "2027-01-01", Denumire: i18n.Text{i18n.Romana: "Mutată"}}); err != nil {
		t.Fatal(err)
	}
	var creat string
	if err := m.db.QueryRow(`SELECT creat_la FROM sarbatori WHERE data = '2027-01-01' AND an = 2027`).Scan(&creat); err != nil || creat != "initial" {
		t.Fatalf("data creării: %s, %v", creat, err)
	}
	if n, err := m.serviciu.CompleteazaSarbatori2026(ctx); err != nil || n != 1 {
		t.Fatalf("completarea datei lipsă: %d, %v", n, err)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'sarbatori'`).Scan(&n); err != nil || n != 3 {
		t.Fatalf("jurnal: %d, %v", n, err)
	}
}

func TestValidareSarbatori(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	for _, data := range []string{"2026-02-29", "2026-02-30", "2026-13-01", "0000-01-01", "2026-1-01", ""} {
		verificaCod(t, m.serviciu.SalveazaSarbatoare(ctx, "", Sarbatoare{Data: data, Denumire: i18n.Text{i18n.Romana: "Test"}}), "sarbatori.data_invalida")
	}
	verificaCod(t, m.serviciu.SalveazaSarbatoare(ctx, "", Sarbatoare{Data: "2026-01-01", Denumire: i18n.Text{i18n.Romana: "Duplicat"}}), "sarbatori.data_duplicata")
	verificaCod(t, m.serviciu.SalveazaSarbatoare(ctx, "2026-01-02", Sarbatoare{Data: "2026-01-01", Denumire: i18n.Text{i18n.Romana: "Duplicat"}}), "sarbatori.data_duplicata")
	verificaCod(t, m.serviciu.SalveazaSarbatoare(ctx, "", Sarbatoare{Data: "2026-02-01", Denumire: i18n.Text{i18n.Romana: " "}}), "sarbatori.denumire_obligatorie")
	verificaCod(t, m.serviciu.SalveazaSarbatoare(ctx, "2026-02-01", Sarbatoare{Data: "2026-02-01", Denumire: i18n.Text{i18n.Romana: "Test"}}), "sarbatori.inexistenta")
	if _, err := m.serviciu.CitesteCalendar(ctx, 0); err == nil {
		t.Fatal("an invalid acceptat")
	}
	if err := m.serviciu.SalveazaSarbatoare(ctx, "", Sarbatoare{Data: "2028-02-29", Denumire: i18n.Text{i18n.Romana: "An bisect"}}); err != nil {
		t.Fatal(err)
	}
}
