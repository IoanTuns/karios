package rapoarte

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

// pdfMinim construieste un PDF valid cu n pagini goale (pentru testul unirii).
func pdfMinim(n int) []byte {
	var b bytes.Buffer
	offset := []int{}
	obiect := func(text string) {
		offset = append(offset, b.Len())
		fmt.Fprintf(&b, "%d 0 obj\n%s\nendobj\n", len(offset), text)
	}
	b.WriteString("%PDF-1.4\n")
	obiect("<< /Type /Catalog /Pages 2 0 R >>")
	kids := ""
	for i := 0; i < n; i++ {
		kids += fmt.Sprintf("%d 0 R ", 3+i)
	}
	obiect(fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", kids, n))
	for i := 0; i < n; i++ {
		obiect("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] >>")
	}
	xref := b.Len()
	fmt.Fprintf(&b, "xref\n0 %d\n0000000000 65535 f \n", len(offset)+1)
	for _, o := range offset {
		fmt.Fprintf(&b, "%010d 00000 n \n", o)
	}
	fmt.Fprintf(&b, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(offset)+1, xref)
	return b.Bytes()
}

// F09-T02 (D104): partile pe sectii se unesc in ordine, fiecare pe pagina ei.
func TestUnestePDFSectii(t *testing.T) {
	ctx := context.Background()
	pdf, pagini, err := UnestePDF(ctx, [][]byte{pdfMinim(2), pdfMinim(1), pdfMinim(3)})
	if err != nil || pagini != 6 {
		t.Fatal(pagini, err)
	}
	if n, err := NumaraPaginiPDF(ctx, pdf); err != nil || n != 6 {
		t.Fatal(n, err)
	}
	if _, n, err := UnestePDF(ctx, [][]byte{pdfMinim(2)}); err != nil || n != 2 {
		t.Fatal("o singura parte", n, err)
	}
	if _, _, err := UnestePDF(ctx, nil); err == nil {
		t.Fatal("fara parti")
	}
}

// F09-T02 (D104): SectieID "-" alege doar angajatii fara sectie.
func TestRaportFaraSectie(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	if _, err := m.db.Exec(`INSERT INTO angajati(id,nume,prenume,marca,creat_la,modificat_la) VALUES ('e','Fictiv','FaraSectie','DEMO-R04','x','x')`); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	marci := func(sectie string) []string {
		t.Helper()
		doc, err := m.s.Pregateste(ctx, Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "ro", SectieID: sectie}, momentFix)
		if err != nil {
			t.Fatal(sectie, err)
		}
		col := -1
		for i, c := range doc.Coloane {
			if c.Cod == "marca" {
				col = i
			}
		}
		rez := []string{}
		for _, r := range doc.Randuri {
			if !r.Grup && col >= 0 {
				rez = append(rez, r.Valori[col])
			}
		}
		return rez
	}
	if r := marci("s"); fmt.Sprint(r) != "[DEMO-R01 DEMO-R02]" {
		t.Fatal("sectia s", r)
	}
	if r := marci(FaraSectie); fmt.Sprint(r) != "[DEMO-R04]" {
		t.Fatal("fara sectie", r)
	}
	if r := marci(""); len(r) != 3 {
		t.Fatal("toate", r)
	}
}
