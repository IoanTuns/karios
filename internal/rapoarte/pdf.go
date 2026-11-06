package rapoarte

import (
	"bytes"
	"context"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Din F09-T23 toate rapoartele se deseneaza cu pdfme in interfata (D087);
// motorul vechi (maroto pentru rapoartele standard si desenarea gofpdf a celor
// flexibile) a fost sters. Aici raman doar functiile comune pentru PDF.

func NumaraPaginiPDF(ctx context.Context, pdf []byte) (int, error) {
	n, err := api.PageCount(ctx, bytes.NewReader(pdf), nil)
	return n, traduce(err)
}
