package exportexcel

import (
	"archive/zip"
	"bytes"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"path/filepath"
	"pontix/internal/i18n"
)

// ValideazaRegistru verifica arhiva si redeschide toate foile, inclusiv sirurile partajate.
func ValideazaRegistru(b []byte) error {
	z, e := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if e != nil {
		return &i18n.Eroare{Cod: "excel.fisier_invalid"}
	}
	for _, f := range z.File {
		r, e := f.Open()
		if e != nil {
			return e
		}
		_, e = io.Copy(io.Discard, r)
		r.Close()
		if e != nil {
			return e
		}
	}
	f, e := excelize.OpenReader(bytes.NewReader(b))
	if e != nil {
		return e
	}
	defer f.Close()
	if len(f.GetSheetList()) == 0 {
		return &i18n.Eroare{Cod: "excel.fisier_invalid"}
	}
	for _, s := range f.GetSheetList() {
		if _, e := f.GetRows(s); e != nil {
			return e
		}
	}
	return nil
}

// SalveazaVerificat valideaza fisierul temporar final inainte de inlocuirea destinatiei.
func SalveazaVerificat(cale string, b []byte) error {
	if e := ValideazaRegistru(b); e != nil {
		return e
	}
	f, e := os.CreateTemp(filepath.Dir(cale), ".pontix-export-*.xlsx")
	if e != nil {
		return e
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, e = f.Write(b); e != nil {
		f.Close()
		return e
	}
	if e = f.Sync(); e != nil {
		f.Close()
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	citit, e := os.ReadFile(tmp)
	if e != nil {
		return e
	}
	if !bytes.Equal(citit, b) {
		return &i18n.Eroare{Cod: "excel.fisier_invalid"}
	}
	if e = ValideazaRegistru(citit); e != nil {
		return e
	}
	return os.Rename(tmp, cale)
}
