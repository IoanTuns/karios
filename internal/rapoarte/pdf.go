package rapoarte

import (
	"bytes"
	"context"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Din F09-T23 toate rapoartele se deseneaza cu pdfme in interfata (D087);
// motorul vechi (maroto pentru rapoartele standard si desenarea gofpdf a celor
// flexibile) a fost sters. Aici raman doar functiile comune pentru PDF.

func NumaraPaginiPDF(ctx context.Context, pdf []byte) (int, error) {
	n, err := api.PageCount(ctx, bytes.NewReader(pdf), nil)
	return n, traduce(err)
}

// UnestePDF pune PDF-urile unul dupa altul, fiecare incepand pe pagina noua (F09-T02, D104:
// „o pagina pe sectie”). Intoarce PDF-ul unit si numarul de pagini.
func UnestePDF(ctx context.Context, parti [][]byte) ([]byte, int, error) {
	if len(parti) == 0 {
		return nil, 0, eroare("rapoarte.parametri_invalizi")
	}
	if len(parti) == 1 {
		n, err := NumaraPaginiPDF(ctx, parti[0])
		return parti[0], n, err
	}
	surse := make([]io.ReadSeeker, len(parti))
	for i, p := range parti {
		surse[i] = bytes.NewReader(p)
	}
	var rez bytes.Buffer
	if err := api.MergeRaw(ctx, surse, &rez, false, nil); err != nil {
		return nil, 0, traduce(err)
	}
	n, err := NumaraPaginiPDF(ctx, rez.Bytes())
	return rez.Bytes(), n, err
}
