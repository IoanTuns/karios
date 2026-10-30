// Package raportare construiește raportul „Raportează o problemă / Părerea ta” (F16-T15, I093):
// un text Markdown anonimizat, salvat ca fișier și trimis, la alegere, ca issue GitHub precompletat.
// Aplicația nu are niciun token GitHub: repo-ul e public, deci un token din executabil ar fi public.
// Utilizatorul trimite raportul cu propriul cont GitHub sau, fără cont, trimite fișierul pe email.
package raportare

import (
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"pontix/internal/diagnoza"
	"pontix/internal/i18n"
)

// AdresaRepo e repo-ul public al proiectului.
const AdresaRepo = "https://github.com/IoanTuns/karios"

// Tipurile de raport.
const (
	TipProblema = "problema"
	TipParere   = "parere"
)

// LungimeMaximaCorpIssue: GitHub refuză adrese mai lungi de câțiva KB; un raport mai lung se
// scurtează în link, iar textul complet rămâne în fișierul salvat.
const LungimeMaximaCorpIssue = 6000

// Cerere e ce scrie utilizatorul în formular.
type Cerere struct {
	Tip       string `json:"tip"`
	Descriere string `json:"descriere"`
	Pasi      string `json:"pasi"`
	CuJurnal  bool   `json:"cu_jurnal"`
	Ecran     string `json:"ecran"`
}

// Context e ce adaugă aplicația, fără date ale angajaților.
type Context struct {
	Versiune   string
	Commit     string
	Sistem     string
	Moment     time.Time
	NumarErori int
	Jurnal     string // ultimele linii din jurnalul de diagnoză, deja curățate la scriere
}

// Text construiește raportul. Descrierea și pașii trec prin diagnoza.Curata, ca jurnalul: textele
// dintre ghilimele, CNP-urile, telefoanele și emailurile nu ajung în raport.
func Text(l i18n.Limba, c Cerere, x Context) string {
	t := func(cheie string) string { return i18n.T(l, cheie, nil) }
	tip := t("raportare.tip_problema")
	if c.Tip == TipParere {
		tip = t("raportare.tip_parere")
	}
	var sb strings.Builder
	sb.WriteString("## " + tip + "\n\n")
	sb.WriteString(diagnoza.Curata(c.Descriere) + "\n\n")
	if p := diagnoza.Curata(c.Pasi); p != "" {
		sb.WriteString("### " + t("raportare.sectiune_pasi") + "\n\n" + p + "\n\n")
	}
	sb.WriteString("### " + t("raportare.sectiune_mediu") + "\n\n")
	sb.WriteString("- " + t("raportare.camp_versiune") + ": " + x.Versiune)
	if x.Commit != "" {
		sb.WriteString(" (" + x.Commit + ")")
	}
	sb.WriteString("\n- " + t("raportare.camp_sistem") + ": " + x.Sistem + "\n")
	if c.Ecran != "" {
		sb.WriteString("- " + t("raportare.camp_ecran") + ": " + diagnoza.Curata(c.Ecran) + "\n")
	}
	sb.WriteString("- " + t("raportare.camp_moment") + ": " + x.Moment.UTC().Format(time.RFC3339) + "\n")
	sb.WriteString("- " + i18n.T(l, "raportare.camp_erori", map[string]string{"numar": strconv.Itoa(x.NumarErori)}) + "\n")
	if c.CuJurnal && strings.TrimSpace(x.Jurnal) != "" {
		sb.WriteString("\n### " + t("raportare.sectiune_jurnal") + "\n\n```\n" + strings.TrimSpace(diagnoza.Curata(x.Jurnal)) + "\n```\n")
	}
	return sb.String()
}

// Titlu e primul rând al descrierii, scurtat, cu tipul în față.
func Titlu(l i18n.Limba, c Cerere) string {
	prefix := i18n.T(l, "raportare.tip_problema", nil)
	if c.Tip == TipParere {
		prefix = i18n.T(l, "raportare.tip_parere", nil)
	}
	rand, _, _ := strings.Cut(strings.TrimSpace(diagnoza.Curata(c.Descriere)), "\n")
	if utf8.RuneCountInString(rand) > 70 {
		rand = string([]rune(rand)[:70]) + "…"
	}
	return "[" + prefix + "] " + rand
}

// AdresaIssue întoarce pagina GitHub de issue nou, completată cu titlul și textul raportului.
func AdresaIssue(l i18n.Limba, titlu, text string) string {
	if utf8.RuneCountInString(text) > LungimeMaximaCorpIssue {
		text = string([]rune(text)[:LungimeMaximaCorpIssue]) + "\n\n" + i18n.T(l, "raportare.text_scurtat", nil)
	}
	v := url.Values{}
	v.Set("title", titlu)
	v.Set("body", text)
	return AdresaRepo + "/issues/new?" + v.Encode()
}
