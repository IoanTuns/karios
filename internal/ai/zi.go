package ai

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// reZiIntrebare: „pe 15 octombrie”, „pe 15”, „în 15 octombrie”, „la data de 15”.
var reZiIntrebare = regexp.MustCompile(`\b(?:pe|in|din|la data de|data de) (\d{1,2})\b(?: (ianuarie|februarie|martie|aprilie|mai|iunie|iulie|august|septembrie|octombrie|noiembrie|decembrie))?`)

// sintetizeazaZi răspunde din cod la întrebările despre o zi anume („Cine lipsește pe 15
// octombrie?”, „Cine lucrează pe 20?”) cu absenții și programările zilei (F14-T28). Pentru o zi
// din afara lunii de lucru spune asta, în loc să ghicească.
func sintetizeazaZi(mesaj string, info ContextPontix) (string, bool) {
	m, _ := corecteazaMesaj(mesaj, info)
	if !esteIntrebare(mesaj, m) {
		return "", false
	}
	g := reZiIntrebare.FindStringSubmatch(m)
	if g == nil {
		return "", false
	}
	zi, _ := strconv.Atoi(g[1])
	luna := info.Luna
	for i, l := range luniNorm {
		if g[2] == l {
			luna = i + 1
		}
	}
	if zi < 1 || zi > 31 || info.Luna < 1 || info.Luna > 12 {
		return "", false
	}
	if luna != info.Luna {
		return fmt.Sprintf("Pot răspunde pentru o zi anume doar din luna de lucru (%s %d). Schimbă luna din antetul asistentului și întreabă din nou.",
			luniAfisate[info.Luna-1], info.An), true
	}
	data := fmt.Sprintf("%04d-%02d-%02d", info.An, luna, zi)
	vreaAbsenti := contineOricare(m, []string{"lipseste", "lipsesc", "absent", "absenti", "concediu", "liber", "libere"})
	vreaProgram := contineOricare(m, []string{"lucreaza", "tura", "ture", "garda", "garzi", "program", "programat", "programati"})
	if !vreaAbsenti && !vreaProgram {
		vreaAbsenti, vreaProgram = true, true
	}
	var absenti, programati []string
	for _, a := range info.Absente {
		if a.Data == data {
			absenti = append(absenti, fmt.Sprintf("- %s (%s) · %s", a.Nume, a.Marca, a.Cod))
		}
	}
	for _, p := range info.TurePlanificate {
		if p.Data == data {
			programati = append(programati, fmt.Sprintf("- %s (%s) · %s", p.Nume, p.Marca, p.CodSchimb))
		}
	}
	sort.Strings(absenti)
	sort.Strings(programati)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📅 **%d %s %d**\n", zi, luniAfisate[luna-1], info.An))
	sectiune := func(titlu string, lista []string, gol string) {
		sb.WriteString(fmt.Sprintf("\n**%s (%d)**\n", titlu, len(lista)))
		if len(lista) == 0 {
			sb.WriteString(gol + "\n")
		}
		for _, l := range lista {
			sb.WriteString(l + "\n")
		}
	}
	if vreaAbsenti {
		sectiune("Absenți", absenti, "Nicio absență înregistrată.")
	}
	if vreaProgram {
		sectiune("Programați", programati, "Nicio tură completată sau planificată.")
	}
	return sb.String(), true
}
