package rapoarte

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"unicode/utf8"
)

const MaxOctetiImagine = 2 * 1024 * 1024
const MaxPixeliImagine = 4_000_000

// Imagine păstrează conținutul în model, fără dependență de fișierul importat.
type Imagine struct {
	Nume       string `json:"nume"`
	Date       string `json:"date"`
	InaltimeMM int64  `json:"inaltime_mm"`
}

func dateImagine(i Imagine) ([]byte, image.Config, string, error) {
	if len(i.Date) > base64.StdEncoding.EncodedLen(MaxOctetiImagine) {
		return nil, image.Config{}, "", eroare("rapoarte.imagine_prea_mare")
	}
	date, err := base64.StdEncoding.Strict().DecodeString(i.Date)
	if err != nil || len(date) == 0 {
		return nil, image.Config{}, "", eroare("rapoarte.imagine_invalida")
	}
	if len(date) > MaxOctetiImagine {
		return nil, image.Config{}, "", eroare("rapoarte.imagine_prea_mare")
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(date))
	if err != nil || (format != "png" && format != "jpeg") || cfg.Width < 1 || cfg.Height < 1 {
		return nil, image.Config{}, "", eroare("rapoarte.imagine_invalida")
	}
	if cfg.Width > 4096 || cfg.Height > 4096 || int64(cfg.Width)*int64(cfg.Height) > MaxPixeliImagine {
		return nil, image.Config{}, "", eroare("rapoarte.imagine_prea_mare")
	}
	return date, cfg, format, nil
}

func valideazaImagine(i Imagine) error {
	if strings.TrimSpace(i.Nume) == "" || utf8.RuneCountInString(i.Nume) > 200 || strings.ContainsAny(i.Nume, "\r\n\x00/\\") {
		return eroare("rapoarte.imagine_invalida")
	}
	if i.InaltimeMM < 5 || i.InaltimeMM > 35 {
		return eroare("rapoarte.inaltime_imagine_invalida")
	}
	date, _, _, err := dateImagine(i)
	if err != nil {
		return err
	}
	// Decodarea completă respinge și fișierele trunchiate cu antet valid.
	if _, _, err := image.Decode(bytes.NewReader(date)); err != nil {
		return eroare("rapoarte.imagine_invalida")
	}
	return nil
}
