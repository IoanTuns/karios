package rapoarte

import (
	"context"
	"strings"
	"testing"

	"pontix/internal/audit"
)

func TestAmprentaInSubsolulLuniiInchise(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	amprenta, err := audit.AmprentaLuna(ctx, m.db, 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.db.Exec(`INSERT INTO stare_luni (an, luna, stare, modificat_la, amprenta) VALUES (2026, 10, 'INCHISA', 'x', ?)`, amprenta); err != nil {
		t.Fatal(err)
	}
	metadate := func(model string) string {
		t.Helper()
		c := Cerere{ModelID: model, An: 2026, Luna: 10, Limba: "ro", AngajatID: "a"}
		doc, err := m.s.Pregateste(ctx, c, momentFix)
		if err != nil {
			t.Fatal(err)
		}
		return doc.Metadate
	}
	for _, model := range []string{"raport-pontaj", "raport-fisa", "raport-pontaj-detaliat", "raport-garzi"} {
		if md := metadate(model); !strings.Contains(md, "amprentă lună închisă "+amprenta) {
			t.Errorf("%s: amprenta lipsește din subsol: %s", model, md)
		}
	}
	if md := metadate("raport-absente"); strings.Contains(md, "amprent") {
		t.Errorf("raportul anual de absențe nu are amprentă de lună: %s", md)
	}
	if _, err := m.db.Exec(`UPDATE pontaj_ciorne SET valoare = 'Z7' WHERE angajat_id = 'a' AND data = '2026-10-23'`); err != nil {
		t.Fatal(err)
	}
	if md := metadate("raport-pontaj"); !strings.Contains(md, "ATENȚIE: datele diferă de amprenta închiderii "+amprenta) {
		t.Errorf("modificarea după închidere nu e semnalată: %s", md)
	}
}
