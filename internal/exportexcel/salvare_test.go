package exportexcel

import (
	"bytes"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"testing"
)

func TestSalvareVerificataPastreazaDestinatiaLaEroare(t *testing.T) {
	p := filepath.Join(t.TempDir(), "proba cu spatii.xlsx")
	vechi := []byte("vechi")
	if e := os.WriteFile(p, vechi, 0600); e != nil {
		t.Fatal(e)
	}
	if e := SalveazaVerificat(p, []byte("nu este XLSX")); e == nil {
		t.Fatal("fisier invalid acceptat")
	}
	b, e := os.ReadFile(p)
	if e != nil || !bytes.Equal(b, vechi) {
		t.Fatal("destinatie pierduta", e)
	}
	f := excelize.NewFile()
	defer f.Close()
	if e := f.SetCellStr("Sheet1", "A1", "Test anonim"); e != nil {
		t.Fatal(e)
	}
	buf, e := f.WriteToBuffer()
	if e != nil {
		t.Fatal(e)
	}
	if e := SalveazaVerificat(p, buf.Bytes()); e != nil {
		t.Fatal(e)
	}
	salvat, e := excelize.OpenFile(p)
	if e != nil {
		t.Fatal(e)
	}
	defer salvat.Close()
	if v, e := salvat.GetCellValue("Sheet1", "A1"); e != nil || v != "Test anonim" {
		t.Fatal(v, e)
	}
}
