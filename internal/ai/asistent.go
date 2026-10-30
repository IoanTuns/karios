package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"pontix/internal/i18n"
)

// PropunereActiune reprezintă o acțiune structurată propusă de asistent,
// care necesită validare și confirmare explicită din partea utilizatorului.
type PropunereActiune struct {
	ID                 string         `json:"id"`
	Tip                string         `json:"tip"` // "planifica_tura", "sterge_tura", etc.
	Descriere          string         `json:"descriere"`
	Parametri          map[string]any `json:"parametri"`
	NecesitaConfirmare bool           `json:"necesita_confirmare"`
}

// MesajConversatie descrie un mesaj în format OpenAI-compatible.
type MesajConversatie struct {
	Rol      string `json:"role"`
	Continut string `json:"content"`
}

// RaspunsAsistent este rezultatul interogării asistentului AI.
type RaspunsAsistent struct {
	Intentie  *IntentiePlan      `json:"intentie,omitempty"`
	Mesaj     string             `json:"mesaj"`
	Propuneri []PropunereActiune `json:"propuneri"`
	Eroare    string             `json:"eroare,omitempty"`
}

// ContextPontix oferă asistentului informații despre personal și schimburile definite.
type ContextPontix struct {
	Domeniu          string                   `json:"domeniu,omitempty"`
	ContextIncomplet bool                     `json:"context_incomplet,omitempty"`
	Limba            string                   `json:"limba"`
	An               int                      `json:"an"`
	Luna             int                      `json:"luna"`
	Angajati         []InfoAngajat            `json:"angajati"`
	Schimburi        []InfoSchimb             `json:"schimburi"`
	Absente          []InfoAbsenta            `json:"absente,omitempty"`
	TurePlanificate  []InfoPozitiePlanificare `json:"ture_planificate,omitempty"`
	// CuPontaj întoarce angajații (după ID) care au ceva completat în pontajul lunii date;
	// folosit de filtrul „fără pontaj” (F14-T16). Nu intră în prompt.
	CuPontaj func(an, luna int) map[string]bool `json:"-"`
	// ActiuniAnterioare sunt propunerile din conversație, cu starea lor („aplicată” sau
	// „neconfirmată”). Intră în promptul de sistem, nu în mesajele asistentului, ca modelul
	// să nu le copieze ca răspuns (F14-T17).
	ActiuniAnterioare []string `json:"-"`
}

type InfoAngajat struct {
	ID      string `json:"id"`
	Nume    string `json:"nume"`
	Prenume string `json:"prenume"`
	Marca   string `json:"marca"`
	Functie string `json:"functie"`
	// Sectie și Categorie servesc comenzilor de grup (F14-T08); nu intră în prompt.
	Sectie      string `json:"sectie,omitempty"`
	Categorie   string `json:"categorie,omitempty"`
	Departament string `json:"departament,omitempty"`
}

type InfoSchimb struct {
	Cod      string `json:"cod"`
	Denumire string `json:"denumire"`
	Ore      int    `json:"ore"`
	Interval string `json:"interval,omitempty"`
}

type InfoAbsenta struct {
	AngajatID string `json:"angajat_id"`
	Marca     string `json:"marca"`
	Nume      string `json:"nume"`
	Data      string `json:"data"`
	Cod       string `json:"cod"`
}

type InfoPozitiePlanificare struct {
	AngajatID string `json:"angajat_id"`
	Marca     string `json:"marca"`
	Nume      string `json:"nume"`
	Data      string `json:"data"`
	CodSchimb string `json:"cod_schimb"`
	Ore       int    `json:"ore,omitempty"`
}

// Client gestionează comunicarea HTTP cu endpoint-ul local al modelului.
type Client struct {
	baseURL    string
	cheieAPI   string
	httpClient *http.Client
}

// CuCheie setează cheia API trimisă serverului local (F14-T12).
func (c *Client) CuCheie(cheie string) *Client {
	c.cheieAPI = cheie
	return c
}

// NouClient creează un client pentru serverul local (de ex: http://127.0.0.1:28491).
func NouClient(port int) *Client {
	if port <= 0 {
		port = PortImplicit
	}
	return &Client{
		baseURL: fmt.Sprintf("http://127.0.0.1:%d", port),
		httpClient: &http.Client{
			Timeout: 300 * time.Second, // prima cerere poate aștepta încălzirea cache-ului (F14-T29)
		},
	}
}

// NouClientCuURL creează un client folosind o adresă de bază arbitrară (util pentru teste cu httptest.Server).
func NouClientCuURL(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CuTransport setează un transport HTTP personalizat (util pentru teste in-memory).
func (c *Client) CuTransport(rt http.RoundTripper) *Client {
	c.httpClient.Transport = rt
	return c
}

// definitiiUnelte conține schema uneltelor (Function Calling) expuse modelului local.
var definitiiUnelte = []map[string]any{
	{
		"type": "function",
		"function": map[string]any{
			"name":        "planifica_tura",
			"description": "Planifică sau atribuie o tură/gardă unui angajat într-o anumită zi.",
			"parameters": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"marca": map[string]any{
						"type":        "string",
						"description": "Numărul de marcă al angajatului sau numele identificat.",
					},
					"data": map[string]any{
						"type":        "string",
						"description": "Data planificării în format AAAA-LL-ZZ (ex: 2026-10-15).",
					},
					"cod_schimb": map[string]any{
						"type":        "string",
						"description": "Codul scurt al schimbului sau gărzii (ex: G24, IN9, S7, Z12, N12).",
					},
					"motiv": map[string]any{
						"type":        "string",
						"description": "Motivul sau contextul solicitării.",
					},
				},
				"required": []string{"marca", "data", "cod_schimb"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]any{
			"name":        "planifica_interval",
			"description": "Planifică o tură/schimb pentru un angajat pe un interval de date, o săptămână sau toată luna dintr-o singură instrucțiune.",
			"parameters": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"marca": map[string]any{
						"type":        "string",
						"description": "Numărul de marcă al angajatului sau numele identificat.",
					},
					"data_inceput": map[string]any{
						"type":        "string",
						"description": "Data de început în format AAAA-LL-ZZ (ex: 2026-10-01).",
					},
					"data_sfarsit": map[string]any{
						"type":        "string",
						"description": "Data de sfârșit în format AAAA-LL-ZZ (ex: 2026-10-31).",
					},
					"cod_schimb": map[string]any{
						"type":        "string",
						"description": "Codul schimbului (ex: Z7 pentru dimineață/zi, S7 pentru după-amiază, G24 etc.).",
					},
					"doar_lucratoare": map[string]any{
						"type":        "boolean",
						"description": "Dacă este true, tura se atribuie doar în zilele lucrătoare (luni-vineri), excluzând weekendul.",
					},
					"filtru_zile": map[string]any{
						"type":        "string",
						"enum":        []string{"toate", "lucratoare", "weekend", "tot_al_doilea_weekend"},
						"description": "Filtru aplicat zilelor: 'lucratoare' (L–V), 'weekend' (doar S–D), 'tot_al_doilea_weekend' (la fiecare două săptămâni începând cu primul weekend din interval), 'toate' (zilnic).",
					},
				},
				"required": []string{"marca", "data_inceput", "data_sfarsit", "cod_schimb"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]any{
			"name":        "planifica_zile",
			"description": "Planifică o tură pentru un angajat într-o listă explicită de date calendaristice.",
			"parameters": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"marca": map[string]any{
						"type":        "string",
						"description": "Numărul de marcă al angajatului sau numele identificat.",
					},
					"cod_schimb": map[string]any{
						"type":        "string",
						"description": "Codul schimbului (ex: Z12, G24, Z7, S7).",
					},
					"zile": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Lista datelor calendaristice în format AAAA-LL-ZZ (ex: ['2026-10-03', '2026-10-04']).",
					},
				},
				"required": []string{"marca", "cod_schimb", "zile"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]any{
			"name":        "sterge_tura",
			"description": "Elimină o tură planificată a unui angajat dintr-o anumită zi.",
			"parameters": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"marca": map[string]any{
						"type":        "string",
						"description": "Numărul de marcă al angajatului sau numele identificat.",
					},
					"data": map[string]any{
						"type":        "string",
						"description": "Data în format AAAA-LL-ZZ din care se șterge tura.",
					},
				},
				"required": []string{"marca", "data"},
			},
		},
	},
	{
		"type": "function",
		"function": map[string]any{
			"name":        "sterge_interval",
			"description": "Elimină turele planificate ale unui angajat pe o perioadă sau toată luna.",
			"parameters": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"marca": map[string]any{
						"type":        "string",
						"description": "Numărul de marcă al angajatului sau numele identificat.",
					},
					"data_inceput": map[string]any{
						"type":        "string",
						"description": "Data de început în format AAAA-LL-ZZ.",
					},
					"data_sfarsit": map[string]any{
						"type":        "string",
						"description": "Data de sfârșit în format AAAA-LL-ZZ.",
					},
				},
				"required": []string{"marca", "data_inceput", "data_sfarsit"},
			},
		},
	},
}

// Interogheaza trimite mesajul utilizatorului către modelul local și parsează rezultatul.
func (c *Client) Interogheaza(ctx context.Context, istoric []MesajConversatie, info ContextPontix) (RaspunsAsistent, error) {
	if len(istoric) > 0 && len(istoric[len(istoric)-1].Continut) > 12000 {
		return RaspunsAsistent{}, &i18n.Eroare{Cod: "ai.cerere_prea_lunga"}
	}
	// Păstrăm cererea curentă întreagă și doar schimburi recente întregi.
	if len(istoric) > 1 {
		buget := 4000
		start := len(istoric) - 1
		for x := len(istoric) - 2; x >= 0; x-- {
			buget -= len(istoric[x].Continut)
			if buget < 0 {
				break
			}
			if istoric[x].Rol == "user" {
				start = x
			}
			if len(istoric)-x >= 5 {
				break
			}
		}
		istoric = istoric[start:]
	}
	// Rutele rezolvate din cod (concediu, planificare, întrebări) sunt alese înainte, de Ruteaza
	// (F14-T27); clientul trimite la model doar ce a rămas.

	corpCerere := map[string]any{
		"max_tokens":   512, // un apel de unealtă are ~40 de tokeni; pe CPU se generează ~4 tokeni/s
		"messages":     mesajeCerere(istoric, info),
		"cache_prompt": true, // llama.cpp: reia din cache începutul identic al cererii (F14-T29)
		"tools":        definitiiUnelte,
		"tool_choice":  "auto",
		"temperature":  0.1, // Răspunsuri deterministe și precise
	}

	dateJSON, err := json.Marshal(corpCerere)
	if err != nil {
		return RaspunsAsistent{}, err
	}

	endpoint := fmt.Sprintf("%s/v1/chat/completions", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(dateJSON))
	if err != nil {
		return RaspunsAsistent{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.cheieAPI != "" {
		req.Header.Set("Authorization", "Bearer "+c.cheieAPI)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RaspunsAsistent{}, fmt.Errorf("eroare comunicare cu motorul AI local: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		corpEroare, _ := io.ReadAll(resp.Body)
		var eroareContext struct {
			Eroare struct {
				Tip string `json:"type"`
			} `json:"error"`
		}
		// Reluăm numai eroarea explicită de context, eliminând un schimb vechi
		// întreg. Instrucțiunile, datele și ultima cerere nu sunt trunchiate.
		if resp.StatusCode == http.StatusBadRequest && json.Unmarshal(corpEroare, &eroareContext) == nil && eroareContext.Eroare.Tip == "exceed_context_size_error" {
			for i := 1; i < len(istoric); i++ {
				if istoric[i].Rol == "user" {
					resp.Body.Close()
					return c.Interogheaza(ctx, istoric[i:], info)
				}
			}
		}
		return RaspunsAsistent{}, fmt.Errorf("serverul AI a returnat codul %d: %s", resp.StatusCode, string(corpEroare))
	}

	var raspunsOpenAI struct {
		Choices []struct {
			Message struct {
				Role      string `json:"role"`
				Content   string `json:"content"`
				ToolCalls []struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raspunsOpenAI); err != nil {
		return RaspunsAsistent{}, fmt.Errorf("eroare decodare răspuns AI: %w", err)
	}

	if len(raspunsOpenAI.Choices) == 0 {
		return RaspunsAsistent{Mesaj: "Nu am primit niciun răspuns de la modelul local."}, nil
	}

	alegere := raspunsOpenAI.Choices[0].Message
	rezultat := RaspunsAsistent{
		Mesaj: strings.TrimSpace(alegere.Content),
	}

	for _, tc := range alegere.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)

		// Canonicizăm marca cu marca exactă din baza de date pentru a preveni diferențele de formatare
		if marcaVal, ok := args["marca"].(string); ok {
			if a, gasit := gasesteAngajatInfo(marcaVal, info.Angajati); gasit {
				args["marca"] = a.Marca
			}
		}

		// Corecție de siguranță pentru intervale și cod de schimb cerut de utilizator
		args = corecteazaSchimbDacaEroare(args, istoric, info.Schimburi, info)

		// Dacă utilizatorul a cerut o perioadă pe lună sau weekenduri, dar modelul a apelat planifica_tura
		tipFunctie := tc.Function.Name
		if tipFunctie == "planifica_tura" && len(istoric) > 0 {
			ultimul := istoric[len(istoric)-1].Continut
			ultimulNorm := strings.ToLower(ultimul)
			cereWeekend := strings.Contains(ultimulNorm, "weekend") || strings.Contains(ultimulNorm, "sambata") || strings.Contains(ultimulNorm, "duminica")
			toataLuna, nrLuna := esteCerereToataLuna(ultimul)

			if toataLuna || cereWeekend {
				an := info.An
				luna := info.Luna
				if nrLuna > 0 {
					luna = nrLuna
				}
				tipFunctie = "planifica_interval"
				ziStart := extrageZiInceput(ultimulNorm)
				if ziStart <= 0 {
					ziStart = 1
				}
				args["data_inceput"] = fmt.Sprintf("%04d-%02d-%02d", an, luna, ziStart)
				args["data_sfarsit"] = ultimaZiLuna(an, luna)
				if cereWeekend {
					args["doar_lucratoare"] = false
					if strings.Contains(ultimulNorm, "tot al doilea") || strings.Contains(ultimulNorm, "al doilea weekend") {
						args["filtru_zile"] = "tot_al_doilea_weekend"
					} else {
						args["filtru_zile"] = "weekend"
					}
				} else {
					args["doar_lucratoare"] = true
				}
				delete(args, "data")
			}
		}

		descriere := construiesteDescrierePropunere(tipFunctie, args, info)
		rezultat.Propuneri = append(rezultat.Propuneri, PropunereActiune{
			ID:                 tc.ID,
			Tip:                tipFunctie,
			Descriere:          descriere,
			Parametri:          args,
			NecesitaConfirmare: true,
		})
	}

	// Fallback defensiv dacă modelul a formulat text în loc să apeleze direct unelte
	if len(rezultat.Propuneri) == 0 && len(istoric) > 0 {
		ultimul := istoric[len(istoric)-1].Continut
		if prop, ok := incearcaConstruirePropunereDinText(ultimul, alegere.Content, info); ok {
			rezultat.Propuneri = append(rezultat.Propuneri, prop)
			rezultat.Mesaj = ""
		}
	}
	// Un răspuns fără nicio acțiune care pretinde totuși că a planificat ceva e o confirmare
	// falsă (F14-T18): utilizatorul ar crede că s-a făcut ceva. Se înlocuiește cu adevărul.
	if len(rezultat.Propuneri) == 0 && esteConfirmareFalsa(rezultat.Mesaj) {
		rezultat.Mesaj = MesajFaraActiune
	}

	return rezultat, nil
}

// MesajFaraActiune înlocuiește un răspuns care pretinde o planificare fără nicio acțiune.
const MesajFaraActiune = "Nu am putut construi o propunere pentru această cerere, deci nu s-a modificat nimic. " +
	"Reformulează cu marca sau numele complet al angajatului și cu schimbul (de ex. Z7, G24), " +
	"ori folosește „toți angajații …” / „secția …” pentru un grup."

var semneConfirmareFalsa = []string{
	"am planificat", "am programat", "am aplicat", "am salvat", "am adaugat", "am sters", "am completat",
	"a fost planificat", "au fost planificat", "a fost aplicat", "au fost aplicat", "a fost programat", "au fost programat",
	"s-a aplicat", "s-au aplicat", "deja aplicat", "actiuni propuse", "actiunea a fost",
}

// esteConfirmareFalsa recunoaște un text care afirmă că o modificare a fost făcută sau propusă.
func esteConfirmareFalsa(mesaj string) bool {
	m := normalizeaza(mesaj)
	if m == "" {
		return false
	}
	if strings.HasPrefix(m, "planificare ") || strings.HasPrefix(m, "stergere ") {
		return true
	}
	for _, s := range semneConfirmareFalsa {
		if strings.Contains(m, s) {
			return true
		}
	}
	return false
}

func esteCerereToataLuna(mesaj string) (bool, int) {
	m := strings.ToLower(mesaj)
	luni := map[string]int{
		"ianuarie": 1, "februarie": 2, "martie": 3, "aprilie": 4,
		"mai": 5, "iunie": 6, "iulie": 7, "august": 8,
		"septembrie": 9, "octombrie": 10, "noiembrie": 11, "decembrie": 12,
	}

	for _, term := range []string{"toata luna", "toată luna", "tata luna", "pe toata luna", "pe toată luna"} {
		if strings.Contains(m, term) {
			return true, 0
		}
	}

	for numeLuna, nrLuna := range luni {
		if strings.Contains(m, numeLuna) {
			if strings.Contains(m, "weekend") || strings.Contains(m, "sambata") || strings.Contains(m, "duminica") {
				return false, nrLuna
			}
			if contineZiSpecifica(m, numeLuna) {
				return false, nrLuna
			}
			return true, nrLuna
		}
	}
	return false, 0
}

func extrageZiInceput(mesaj string) int {
	m := strings.ToLower(mesaj)
	patterns := []string{
		"incepand cu data de ",
		"începând cu data de ",
		"cu data de ",
		"data de ",
		"din data de ",
		"incepand cu ",
		"începând cu ",
		"din ",
		"pe data de ",
	}
	for _, p := range patterns {
		if idx := strings.Index(m, p); idx != -1 {
			rest := strings.TrimSpace(m[idx+len(p):])
			cuvinte := strings.Fields(rest)
			if len(cuvinte) > 0 {
				curat := strings.Trim(cuvinte[0], ",.;:!?")
				if nr, err := strconv.Atoi(curat); err == nil && nr >= 1 && nr <= 31 {
					return nr
				}
			}
		}
	}
	return 0
}

func incearcaConstruirePropunereDinText(mesajUtilizator, raspunsModel string, info ContextPontix) (PropunereActiune, bool) {
	mUtilizator := strings.ToLower(mesajUtilizator)
	mRaspuns := strings.ToLower(raspunsModel)
	textCombinat := mUtilizator + "\n" + mRaspuns

	// 1. Identificăm angajatul
	var angajatGasit *InfoAngajat
	for _, a := range info.Angajati {
		nume := strings.ToLower(a.Nume)
		prenume := strings.ToLower(a.Prenume)
		marca := strings.ToLower(a.Marca)
		marcaFaraCratima := strings.ReplaceAll(marca, "-", "")
		if strings.Contains(textCombinat, marca) || strings.Contains(textCombinat, marcaFaraCratima) ||
			(len(nume) >= 3 && strings.Contains(textCombinat, nume)) ||
			(len(prenume) >= 3 && strings.Contains(textCombinat, prenume)) {
			ang := a
			angajatGasit = &ang
			break
		}
	}
	if angajatGasit == nil {
		return PropunereActiune{}, false
	}

	// 2. Identificăm schimbul
	codSchimb := "Z7"
	if strings.Contains(textCombinat, "z12") || (strings.Contains(textCombinat, "12 ore") && (strings.Contains(textCombinat, "zi") || strings.Contains(textCombinat, "weekend"))) {
		codSchimb = "Z12"
	} else if strings.Contains(textCombinat, "n12") || (strings.Contains(textCombinat, "12 ore") && strings.Contains(textCombinat, "noapte")) {
		codSchimb = "N12"
	} else if strings.Contains(textCombinat, "g24") || strings.Contains(textCombinat, "garda") || strings.Contains(textCombinat, "gardă") || strings.Contains(textCombinat, "24 ore") {
		codSchimb = "G24"
	} else if strings.Contains(textCombinat, "s7") || strings.Contains(textCombinat, "seara") || strings.Contains(textCombinat, "seară") || strings.Contains(textCombinat, "15 la 22") {
		codSchimb = "S7"
	} else if strings.Contains(textCombinat, "z7") || strings.Contains(textCombinat, "7 la 14") || strings.Contains(textCombinat, "dimineata") {
		codSchimb = "Z7"
	}

	// 3. Identificăm luna și anul
	an := info.An
	luna := info.Luna
	luni := map[string]int{
		"ianuarie": 1, "februarie": 2, "martie": 3, "aprilie": 4,
		"mai": 5, "iunie": 6, "iulie": 7, "august": 8,
		"septembrie": 9, "octombrie": 10, "noiembrie": 11, "decembrie": 12,
	}
	for numeLuna, nrLuna := range luni {
		if strings.Contains(textCombinat, numeLuna) {
			luna = nrLuna
			break
		}
	}

	// 4. Identificăm dacă este interval / weekend / tură unică
	cereWeekend := strings.Contains(textCombinat, "weekend") || strings.Contains(textCombinat, "sambata") || strings.Contains(textCombinat, "duminica")
	cereAlternat := strings.Contains(textCombinat, "tot al doilea") || strings.Contains(textCombinat, "al doilea weekend")

	args := map[string]any{
		"marca":      angajatGasit.Marca,
		"cod_schimb": codSchimb,
	}

	var tipFunctie string
	if cereWeekend || strings.Contains(textCombinat, "octombrie") || strings.Contains(textCombinat, "luna") || strings.Contains(textCombinat, "interval") {
		tipFunctie = "planifica_interval"
		ziStart := extrageZiInceput(textCombinat)
		if ziStart <= 0 {
			ziStart = 1
		}
		args["data_inceput"] = fmt.Sprintf("%04d-%02d-%02d", an, luna, ziStart)
		args["data_sfarsit"] = ultimaZiLuna(an, luna)
		if cereWeekend {
			args["doar_lucratoare"] = false
			if cereAlternat {
				args["filtru_zile"] = "tot_al_doilea_weekend"
			} else {
				args["filtru_zile"] = "weekend"
			}
		} else {
			args["doar_lucratoare"] = true
		}
	} else if zi := extrageZiInceput(textCombinat); zi > 0 {
		tipFunctie = "planifica_tura"
		args["data"] = fmt.Sprintf("%04d-%02d-%02d", an, luna, zi)
	} else {
		return PropunereActiune{}, false
	}

	desc := construiesteDescrierePropunere(tipFunctie, args, info)
	return PropunereActiune{
		ID:                 fmt.Sprintf("prop-%d", time.Now().UnixNano()),
		Tip:                tipFunctie,
		Descriere:          desc,
		Parametri:          args,
		NecesitaConfirmare: true,
	}, true
}

func contineZiSpecifica(mesaj, numeLuna string) bool {
	cuvinte := strings.Fields(mesaj)
	for i, c := range cuvinte {
		cCurat := strings.Trim(c, ",.;:!?")
		if cCurat == "pe" && i+1 < len(cuvinte) {
			urm := strings.Trim(cuvinte[i+1], ",.;:!?")
			if nr, err := strconv.Atoi(urm); err == nil && nr >= 1 && nr <= 31 {
				return true
			}
		}
		if cCurat == "data" && i+2 < len(cuvinte) && strings.Trim(cuvinte[i+1], ",.;:!?") == "de" {
			urm := strings.Trim(cuvinte[i+2], ",.;:!?")
			if nr, err := strconv.Atoi(urm); err == nil && nr >= 1 && nr <= 31 {
				return true
			}
		}
		if cCurat == numeLuna && i > 0 {
			prec := strings.Trim(cuvinte[i-1], ",.;:!?")
			if nr, err := strconv.Atoi(prec); err == nil && nr >= 1 && nr <= 31 {
				return true
			}
		}
	}
	return false
}

func ultimaZiLuna(an, luna int) string {
	t := time.Date(an, time.Month(luna)+1, 0, 0, 0, 0, 0, time.UTC)
	return t.Format(time.DateOnly)
}

func curataAlfanumeric(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func extrageCifre(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func extrageLitere(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func gasesteAngajatInfo(marcaSauNume string, angajati []InfoAngajat) (InfoAngajat, bool) {
	raw := strings.TrimSpace(marcaSauNume)
	if raw == "" || len(angajati) == 0 {
		return InfoAngajat{}, false
	}
	norm := strings.ToLower(raw)
	alfa := curataAlfanumeric(raw)

	// 1. Potrivire directă pe marcă
	for _, a := range angajati {
		if strings.ToLower(a.Marca) == norm {
			return a, true
		}
	}

	// 2. Potrivire alfanumerică pe marcă (ex: "DEM008" vs "DEMO-008")
	for _, a := range angajati {
		aAlfa := curataAlfanumeric(a.Marca)
		if aAlfa == alfa {
			return a, true
		}
		cifreRaw := extrageCifre(alfa)
		cifreA := extrageCifre(aAlfa)
		if cifreRaw != "" && cifreRaw == cifreA {
			litRaw := extrageLitere(alfa)
			litA := extrageLitere(aAlfa)
			if strings.HasPrefix(litA, litRaw) || strings.HasPrefix(litRaw, litA) {
				return a, true
			}
		}
	}

	// 3. Potrivire pe nume/prenume
	for _, a := range angajati {
		numeComplet := strings.ToLower(strings.TrimSpace(a.Nume + " " + a.Prenume))
		prenumeComplet := strings.ToLower(strings.TrimSpace(a.Prenume + " " + a.Nume))
		if numeComplet == norm || prenumeComplet == norm ||
			strings.Contains(numeComplet, norm) || strings.Contains(prenumeComplet, norm) ||
			(len(a.Nume) >= 3 && strings.Contains(norm, strings.ToLower(a.Nume))) {
			return a, true
		}
	}

	return InfoAngajat{}, false
}

func formateazaAngajat(marca string, angajati []InfoAngajat) string {
	if a, gasit := gasesteAngajatInfo(marca, angajati); gasit {
		numeComplet := strings.TrimSpace(a.Nume + " " + a.Prenume)
		if a.Marca != "" {
			return fmt.Sprintf("%s (marca %s)", numeComplet, a.Marca)
		}
		return numeComplet
	}
	if marca == "" {
		return "angajat nespecificat"
	}
	return marca
}

func formateazaSchimb(cod string, schimburi []InfoSchimb) string {
	codNorm := strings.TrimSpace(strings.ToUpper(cod))
	for _, s := range schimburi {
		if strings.ToUpper(s.Cod) == codNorm {
			parti := []string{s.Cod}
			if s.Denumire != "" && !strings.EqualFold(s.Denumire, s.Cod) {
				parti = append(parti, s.Denumire)
			}
			var detalii []string
			if s.Interval != "" && !strings.Contains(s.Denumire, s.Interval) {
				detalii = append(detalii, s.Interval)
			}
			if s.Ore > 0 && !strings.Contains(s.Denumire, fmt.Sprintf("%dh", s.Ore)) && !strings.Contains(s.Denumire, fmt.Sprintf("%d ore", s.Ore)) {
				detalii = append(detalii, fmt.Sprintf("%d ore", s.Ore))
			}
			if len(detalii) > 0 {
				return fmt.Sprintf("%s (%s)", strings.Join(parti, " – "), strings.Join(detalii, ", "))
			}
			return strings.Join(parti, " – ")
		}
	}
	if cod == "" {
		return "schimb nespecificat"
	}
	return cod
}

func corecteazaSchimbDacaEroare(args map[string]any, istoric []MesajConversatie, schimburi []InfoSchimb, info ContextPontix) map[string]any {
	if len(istoric) == 0 {
		return args
	}
	ultimulMesaj := strings.ToLower(istoric[len(istoric)-1].Continut)

	// Dacă utilizatorul a cerut "7 la 14", "dimineata", "de la 7", dar modelul a ales din greșeală "S7" (Seară)
	if strings.Contains(ultimulMesaj, "7 la 14") || strings.Contains(ultimulMesaj, "7-14") ||
		strings.Contains(ultimulMesaj, "07:00") || strings.Contains(ultimulMesaj, "dimineata") ||
		strings.Contains(ultimulMesaj, "dimineață") {
		if cod, ok := args["cod_schimb"].(string); ok && strings.EqualFold(cod, "S7") {
			args["cod_schimb"] = "Z7"
		}
	}

	// Dacă utilizatorul a cerut 12 ore pe zi / weekend
	if strings.Contains(ultimulMesaj, "12 ore") || strings.Contains(ultimulMesaj, "schimb de 12") {
		if strings.Contains(ultimulMesaj, "zi") || strings.Contains(ultimulMesaj, "ziua") || strings.Contains(ultimulMesaj, "pe zi") {
			args["cod_schimb"] = "Z12"
		} else if strings.Contains(ultimulMesaj, "noapte") {
			args["cod_schimb"] = "N12"
		}
	}

	// Detectare weekend și tot al doilea weekend
	esteWeekend := strings.Contains(ultimulMesaj, "weekend") || strings.Contains(ultimulMesaj, "sambata") || strings.Contains(ultimulMesaj, "sâmbătă") || strings.Contains(ultimulMesaj, "duminica") || strings.Contains(ultimulMesaj, "duminică")
	esteAlternat := strings.Contains(ultimulMesaj, "tot al doilea") || strings.Contains(ultimulMesaj, "al doilea weekend") || strings.Contains(ultimulMesaj, "la doua saptamani") || strings.Contains(ultimulMesaj, "la 2 saptamani")

	if esteWeekend {
		args["doar_lucratoare"] = false
		if esteAlternat {
			args["filtru_zile"] = "tot_al_doilea_weekend"
		} else {
			args["filtru_zile"] = "weekend"
		}
	} else if strings.Contains(ultimulMesaj, "toata luna") || strings.Contains(ultimulMesaj, "toată luna") || strings.Contains(ultimulMesaj, "tata luna") {
		if _, ok := args["doar_lucratoare"]; !ok {
			args["doar_lucratoare"] = true
		}
	}

	// Extragere dată de început dacă utilizatorul menționează "începând cu data de X"
	if ziStart := extrageZiInceput(ultimulMesaj); ziStart > 0 {
		an := info.An
		luna := info.Luna
		if _, nrLuna := esteCerereToataLuna(ultimulMesaj); nrLuna > 0 {
			luna = nrLuna
		}
		args["data_inceput"] = fmt.Sprintf("%04d-%02d-%02d", an, luna, ziStart)
	}

	return args
}

// Incalzeste trimite o dată instrucțiunile și uneltele, fără date, ca serverul să le aibă în cache
// înainte de prima întrebare a utilizatorului (F14-T29). Răspunsul nu se folosește.
func (c *Client) Incalzeste(ctx context.Context) error {
	corp, err := json.Marshal(map[string]any{
		"max_tokens":   1,
		"messages":     []map[string]any{{"role": "system", "content": promptInstructiuni()}, {"role": "user", "content": "."}},
		"tools":        definitiiUnelte,
		"tool_choice":  "auto",
		"cache_prompt": true,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(corp))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.cheieAPI != "" {
		req.Header.Set("Authorization", "Bearer "+c.cheieAPI)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.Body.Close()
}

// promptInstructiuni e partea fixă a promptului de sistem: nu depinde de cerere, deci serverul
// o citește o singură dată și o păstrează în cache între cereri (F14-T29). Pe CPU, citirea celor
// ~2.000 de tokeni de instrucțiuni și unelte dura peste 100 s la fiecare mesaj.
func promptInstructiuni() string {
	var sb strings.Builder
	sb.WriteString("Ești Asistentul inteligent al aplicației Pontix pentru pontaj și planificare de gărzi/ture.\n")
	sb.WriteString("Răspunde întotdeauna concis, politicos și în limba română (sau engleză dacă ți se cere).\n")
	sb.WriteString("NU inventa date și NU calcula formule salariale pe cont propriu.\n")
	sb.WriteString("Dacă utilizatorul îți cere să planifici, să modifici sau să ștergi ture, apelează uneltele puse la dispoziție.\n")
	sb.WriteString("IMPORTANT PENTRU PERIOADE / LUNI / INTERVALE:\n")
	sb.WriteString("- Când utilizatorul menționează o lună ('în octombrie', 'in octombrie', 'pe octombrie', 'luna octombrie') sau o perioadă ('toată luna', 'toata luna', 'tata luna', 'toată săptămâna') FĂRĂ o zi anume (fără 'pe 15' etc.), înseamnă ÎNTREAGA LUNĂ! Folosește neapărat unealta `planifica_interval` cu data_inceput (ex: 2026-10-01) și data_sfarsit (ex: 2026-10-31), și doar_lucratoare: true. NU apela `planifica_tura` pentru o singură zi când se cere toată luna!\n")
	sb.WriteString("- Pentru ture normale de zi ('toată luna', '7 la 14', etc.), setează ÎNTOTDEAUNA `doar_lucratoare: true` pentru a exclude weekendurile!\n")
	sb.WriteString("- Pentru comenzi pe weekend ('în weekend', 'tot al doilea weekend'), setează `filtru_zile: 'tot_al_doilea_weekend'` sau `'weekend'` și setează `doar_lucratoare: false`!\n")
	sb.WriteString("- ALEGEREA SCHIMBURILOR DUPĂ INTERVAL SAU ORE:\n")
	sb.WriteString("  * '7 la 14', '07:00–14:00', 'dimineața', 'pe zi' -> codul este Z7 (Zi 07–14). NICIODATĂ nu alege S7 pentru 7 la 14!\n")
	sb.WriteString("  * '12 ore pe zi', 'schimb de 12 ore', '12 ore ziua' -> codul este Z12 (Zi liberă 07–19, 12 ore).\n")
	sb.WriteString("  * '12 ore noaptea', '12 ore noapte' -> codul este N12 (Noapte 19–07, 12 ore).\n")
	sb.WriteString("  * '15 la 22', '15:00–22:00', 'seară', 'după-amiază' -> codul este S7 (Seară 15–22).\n")
	sb.WriteString("  * '22 la 7', 'noapte' -> codul este N9 sau N12.\n")
	sb.WriteString("  * 'gardă', '24 ore' -> codul este G24.\n")
	sb.WriteString("- REGULĂ STRICTĂ: Orice planificare cerută de utilizator se execută EXCLUSIV prin apel de unelte (tool calls). NU genera niciodată în text confirmări simulate sau 'Acțiuni propuse' fără tool call!\n")
	sb.WriteString("- Răspunde doar la ultima cerere a utilizatorului. NU scrie în text că ai planificat sau aplicat ceva: orice modificare se face doar prin apel de unealtă, iar utilizatorul o confirmă.\n")
	sb.WriteString("- Folosește EXACT numărul de marcă așa cum este indicat în lista de personal din datele de lucru (ex: 'DEMO-008', NU prescurta și nu elimina cratima).\n")
	sb.WriteString("Datele de lucru (anul, luna, schimburile, personalul, absențele și turele) vin într-un mesaj de sistem separat, chiar înaintea ultimei cereri a utilizatorului. Acționează doar dacă acea cerere o cere explicit.\n")
	return sb.String()
}

// promptDate e partea variabilă: datele relevante pentru cererea curentă (ContextPentruCerere).
func promptDate(info ContextPontix) string {
	var sb strings.Builder
	if info.Domeniu != "" {
		sb.WriteString("Context specific cererii: " + info.Domeniu + ". Selecțiile de grup și modificările sunt verificate ulterior de aplicație.\n")
	}
	if info.ContextIncomplet {
		sb.WriteString("CONTEXT PARȚIAL: lipsa unei persoane sau înregistrări din acest context NU înseamnă că ea nu există. Nu enumera și nu număra echipa din acest extras. Nu afirma rezultate exacte din date incomplete; cere precizări.\n")
	}
	sb.WriteString(fmt.Sprintf("Anul curent de lucru: %d, Luna: %d.\n", info.An, info.Luna))

	if len(info.Schimburi) > 0 {
		sb.WriteString("Schimburi disponibile:\n")
		for _, s := range info.Schimburi {
			interv := s.Interval
			if interv == "" && s.Ore > 0 {
				interv = fmt.Sprintf("%d ore", s.Ore)
			}
			sb.WriteString(fmt.Sprintf("- Cod: %s | Denumire: %s | Interval: %s | Durată: %d ore\n", s.Cod, s.Denumire, interv, s.Ore))
		}
	}

	if len(info.Angajati) > 0 {
		sb.WriteString("Personal identificabil (Folosește marca corectă când apelezi funcții):\n")
		for _, a := range info.Angajati {
			sb.WriteString(fmt.Sprintf("- %s %s (Marca: %s, Funcție: %s)\n", a.Nume, a.Prenume, a.Marca, a.Functie))
		}
	}

	if len(info.Absente) > 0 {
		sb.WriteString("Concedii și absențe programate în această lună:\n")
		for _, a := range info.Absente {
			sb.WriteString(fmt.Sprintf("- %s (%s): %s pe %s\n", a.Nume, a.Marca, a.Cod, a.Data))
		}
	}

	if len(info.ActiuniAnterioare) > 0 {
		sb.WriteString("Acțiuni din această conversație (informativ; nu le repeta și nu le reformula în răspuns):\n")
		for _, a := range info.ActiuniAnterioare {
			sb.WriteString("- " + a + "\n")
		}
	}

	if len(info.TurePlanificate) > 0 {
		sb.WriteString("Ture și gărzi deja programate în această lună:\n")
		for _, t := range info.TurePlanificate {
			sb.WriteString(fmt.Sprintf("- %s (%s): Schimb %s pe %s\n", t.Nume, t.Marca, t.CodSchimb, t.Data))
		}
	}

	return sb.String()
}

// mesajeCerere pune instrucțiunile fixe primele, apoi istoricul, iar datele de lucru într-un mesaj de
// sistem chiar înaintea ultimei cereri: așa, începutul cererii (instrucțiuni, unelte, istoric) e
// identic de la un mesaj la altul și serverul reia din cache tot ce nu s-a schimbat, iar cererea
// utilizatorului rămâne curată (lipită de date, modelul le trata ca parte din cerere).
func mesajeCerere(istoric []MesajConversatie, info ContextPontix) []map[string]any {
	mesaje := []map[string]any{{"role": "system", "content": promptInstructiuni()}}
	for x, m := range istoric {
		if x == len(istoric)-1 && m.Rol == "user" {
			mesaje = append(mesaje, map[string]any{"role": "system", "content": "DATE DE LUCRU:\n" + promptDate(info)})
		}
		mesaje = append(mesaje, map[string]any{"role": m.Rol, "content": m.Continut})
	}
	return mesaje
}

func construiesteDescrierePropunere(numeFunctie string, args map[string]any, info ContextPontix) string {
	switch numeFunctie {
	case "planifica_tura":
		marca, _ := args["marca"].(string)
		data, _ := args["data"].(string)
		schimb, _ := args["cod_schimb"].(string)
		schimbText := formateazaSchimb(schimb, info.Schimburi)
		angajatText := formateazaAngajat(marca, info.Angajati)
		return fmt.Sprintf("Planificare schimb %s pentru %s pe %s", schimbText, angajatText, data)
	case "planifica_interval":
		marca, _ := args["marca"].(string)
		d1, _ := args["data_inceput"].(string)
		d2, _ := args["data_sfarsit"].(string)
		schimb, _ := args["cod_schimb"].(string)
		lucratoare, _ := args["doar_lucratoare"].(bool)
		filtruZile, _ := args["filtru_zile"].(string)
		schimbText := formateazaSchimb(schimb, info.Schimburi)
		angajatText := formateazaAngajat(marca, info.Angajati)
		supliment := ""
		if filtruZile == "tot_al_doilea_weekend" {
			supliment = " (tot al doilea weekend)"
		} else if filtruZile == "weekend" {
			supliment = " (doar în weekend S–D)"
		} else if lucratoare || filtruZile == "lucratoare" {
			supliment = " (doar zile lucrătoare L–V)"
		}
		return fmt.Sprintf("Planificare schimb %s pentru %s în intervalul %s – %s%s", schimbText, angajatText, d1, d2, supliment)
	case "planifica_zile":
		marca, _ := args["marca"].(string)
		schimb, _ := args["cod_schimb"].(string)
		var zile []string
		if rz, ok := args["zile"].([]any); ok {
			for _, z := range rz {
				if s, ok := z.(string); ok {
					zile = append(zile, s)
				}
			}
		} else if sz, ok := args["zile"].([]string); ok {
			zile = sz
		}
		schimbText := formateazaSchimb(schimb, info.Schimburi)
		angajatText := formateazaAngajat(marca, info.Angajati)
		return fmt.Sprintf("Planificare schimb %s pentru %s în zilele: %s", schimbText, angajatText, strings.Join(zile, ", "))
	case "sterge_tura":
		marca, _ := args["marca"].(string)
		data, _ := args["data"].(string)
		angajatText := formateazaAngajat(marca, info.Angajati)
		return fmt.Sprintf("Ștergere tură pentru %s pe %s", angajatText, data)
	case "sterge_interval":
		marca, _ := args["marca"].(string)
		d1, _ := args["data_inceput"].(string)
		d2, _ := args["data_sfarsit"].(string)
		angajatText := formateazaAngajat(marca, info.Angajati)
		return fmt.Sprintf("Ștergere ture pentru %s în intervalul %s – %s", angajatText, d1, d2)
	default:
		return fmt.Sprintf("Acțiune %s cu parametrii %v", numeFunctie, args)
	}
}

// AnalizeazaInterogareSinteza verifică dacă mesajul utilizatorului este o interogare analitică
// despre concedii, disponibilitate în weekend, acoperirea posturilor sau statistici gărzi,
// și generează un răspuns factual direct pe baza contextului furnizat.
func AnalizeazaInterogareSinteza(mesaj string, info ContextPontix) (string, bool) {
	m := strings.ToLower(strings.TrimSpace(mesaj))
	if m == "" {
		return "", false
	}

	numeLuni := []string{"ianuarie", "februarie", "martie", "aprilie", "mai", "iunie", "iulie", "august", "septembrie", "octombrie", "noiembrie", "decembrie"}
	numeLuna := "luna curentă"
	if info.Luna >= 1 && info.Luna <= 12 {
		numeLuna = numeLuni[info.Luna-1]
	}

	// 0. Întrebări despre o zi anume (F14-T28), înaintea listei concediilor pe toată luna.
	if r, ok := sintetizeazaZi(mesaj, info); ok {
		return r, true
	}

	// 1. Interogare concedii / absențe ("cine este în concediu", "cine e în concediu", "concedii luna", "absențe")
	// O comandă („X merge în concediu…”) nu e o întrebare despre concedii (F14-T26).
	if (strings.Contains(m, "in concediu") || strings.Contains(m, "în concediu") ||
		strings.Contains(m, "concedii") || strings.Contains(m, "absente") || strings.Contains(m, "absențe")) && !EsteComandaConcediu(mesaj, info) {
		return sintetizeazaConcedii(info, numeLuna), true
	}

	// 1b. Listări de angajați și de schimburi (F14-T09); se compară fără diacritice
	if mn, _ := corecteazaMesaj(m, info); esteListareAngajati(mn) {
		return sintetizeazaListaAngajati(mn, info), true
	} else if esteListareSchimburi(mn) {
		return sintetizeazaListaSchimburi(info), true
	}

	// 2. Interogare disponibilitate weekend ("cine este liber", "cine e liber în acest weekend", "disponibilitate")
	if strings.Contains(m, "liber") || strings.Contains(m, "disponibil") {
		return sintetizeazaDisponibilitateWeekend(m, info, numeLuna), true
	}

	// 3. Verificare acoperire posturi obligatorii / gărzi weekend ("lipsesc posturi", "verifică acoperire", "neacoperit")
	if strings.Contains(m, "lipsesc posturi") || strings.Contains(m, "posturi obligatorii") ||
		strings.Contains(m, "acoperire") || strings.Contains(m, "neacoperit") {
		return sintetizeazaAcoperire(info, numeLuna), true
	}

	// 4. Statistici gărzi / ore ("statistici gărzi", "câte gărzi", "total ore planificate", "sumar gărzi")
	if strings.Contains(m, "statistici") || strings.Contains(m, "cate garzi") || strings.Contains(m, "câte gărzi") ||
		strings.Contains(m, "total ore") || strings.Contains(m, "sumar garzi") {
		return sintetizeazaStatisticiGarzi(info, numeLuna), true
	}

	// 5. Ghid de comenzi rapide și ajutor ("ajutor", "help", "ce comenzi", "ce poti face")
	if m == "ajutor" || m == "help" || strings.Contains(m, "ce poti face") || strings.Contains(m, "ce poți face") ||
		strings.Contains(m, "ce comenzi") || strings.Contains(m, "ce stii sa faci") || strings.Contains(m, "ce știi să faci") {
		return sintetizeazaGhidComenzi(numeLuna, info.An), true
	}

	return "", false
}

func sintetizeazaGhidComenzi(numeLuna string, an int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("💡 **Ghid Asistent Pontix (100%% Local & Offline)**\n\n"))
	sb.WriteString("Iată câteva exemple de comenzi pe care le poți folosi:\n\n")
	sb.WriteString("📋 **Planificare și atribuire:**\n")
	sb.WriteString("- *„Pune-l pe Popescu în gardă pe 15 a lunii”*\n")
	sb.WriteString(fmt.Sprintf("- *„Dobre în %s e pe zi de la 7 la 14”*\n", numeLuna))
	sb.WriteString("- *„Dumitru Andrei lucrează tot al doilea weekend pe zi, schimb de 12 ore, începând cu data de 3”*\n")
	sb.WriteString("- *„Șterge garda lui Popescu din 15”*\n\n")
	sb.WriteString("🔍 **Interogări analitice instantanee (fără timp de așteptare):**\n")
	sb.WriteString("- *„Cine este în concediu luna asta?”*\n")
	sb.WriteString("- *„Cine este liber în acest weekend?”*\n")
	sb.WriteString("- *„Verifică dacă lipsesc posturi obligatorii / gărzi weekend”*\n")
	sb.WriteString("- *„Statistici gărzi și ore planificate”*\n")
	sb.WriteString("- *„Listează angajații din secția ATI”*, *„Ce schimburi există?”*\n")
	return sb.String()
}

func sintetizeazaConcedii(info ContextPontix, numeLuna string) string {
	if len(info.Absente) == 0 {
		return fmt.Sprintf("🏖️ **Concedii și absențe (%s %d)**\n\nNu este înregistrat niciun concediu sau absență pentru luna %s %d.", numeLuna, info.An, numeLuna, info.An)
	}

	tipuri := map[string]string{
		"CO": "Concediu de odihnă",
		"CM": "Concediu medical",
		"CF": "Concediu fără plată",
		"CS": "Concediu de studii",
		"M":  "Misiune / Delegație",
		"FP": "Formare profesională",
		"D":  "Donare sânge",
	}

	grupuri := make(map[string][]InfoAbsenta)
	for _, a := range info.Absente {
		grupuri[a.AngajatID] = append(grupuri[a.AngajatID], a)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🏖️ **Concedii și absențe programate în %s %d**\n\n", numeLuna, info.An))

	totalZile := 0
	for _, a := range info.Angajati {
		lista, ok := grupuri[a.ID]
		if !ok || len(lista) == 0 {
			continue
		}
		sort.Slice(lista, func(i, j int) bool { return lista[i].Data < lista[j].Data })

		zileStr := make([]string, 0, len(lista))
		coduriUnice := make(map[string]bool)
		for _, abs := range lista {
			parti := strings.Split(abs.Data, "-")
			if len(parti) == 3 {
				zileStr = append(zileStr, parti[2])
			} else {
				zileStr = append(zileStr, abs.Data)
			}
			coduriUnice[abs.Cod] = true
			totalZile++
		}

		descTipuri := make([]string, 0, len(coduriUnice))
		for cod := range coduriUnice {
			if den, ex := tipuri[strings.ToUpper(cod)]; ex {
				descTipuri = append(descTipuri, fmt.Sprintf("%s (%s)", den, cod))
			} else {
				descTipuri = append(descTipuri, cod)
			}
		}

		intervalText := fmt.Sprintf("%s → %s", lista[0].Data, lista[len(lista)-1].Data)
		if len(lista) == 1 {
			intervalText = lista[0].Data
		}

		sb.WriteString(fmt.Sprintf("- 👤 **%s %s** (marca `%s`):\n  * Tip: %s\n  * Perioadă: %s (%d zile: %s)\n",
			a.Nume, a.Prenume, a.Marca, strings.Join(descTipuri, ", "), intervalText, len(lista), strings.Join(zileStr, ", ")))
	}

	sb.WriteString(fmt.Sprintf("\n📊 **Total secție:** %d zile de absență înregistrate în această lună.", totalZile))
	return sb.String()
}

func sintetizeazaDisponibilitateWeekend(mesaj string, info ContextPontix, numeLuna string) string {
	an := info.An
	luna := info.Luna
	if an <= 0 {
		an = time.Now().Year()
	}
	if luna <= 0 {
		luna = int(time.Now().Month())
	}

	ultZi := time.Date(an, time.Month(luna)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	type Weekend struct {
		Sambata  int
		Duminica int
	}
	var weekenduri []Weekend
	for zi := 1; zi <= ultZi; zi++ {
		t := time.Date(an, time.Month(luna), zi, 0, 0, 0, 0, time.UTC)
		if t.Weekday() == time.Saturday {
			dum := zi + 1
			if dum > ultZi {
				dum = 0
			}
			weekenduri = append(weekenduri, Weekend{Sambata: zi, Duminica: dum})
		}
	}

	if len(weekenduri) == 0 {
		return fmt.Sprintf("Nu există zile de weekend identificate pentru %s %d.", numeLuna, an)
	}

	weekendSelectat := weekenduri[0]
	ziCeruta := extrageZiInceput(strings.ToLower(mesaj))
	if ziCeruta > 0 {
		for _, w := range weekenduri {
			if w.Sambata == ziCeruta || w.Duminica == ziCeruta || (ziCeruta >= w.Sambata && (w.Duminica == 0 || ziCeruta <= w.Duminica)) {
				weekendSelectat = w
				break
			}
		}
	} else {
		acum := time.Now()
		if acum.Year() == an && int(acum.Month()) == luna {
			ziAzi := acum.Day()
			for _, w := range weekenduri {
				if w.Duminica >= ziAzi || w.Sambata >= ziAzi {
					weekendSelectat = w
					break
				}
			}
		}
	}

	dataSam := fmt.Sprintf("%04d-%02d-%02d", an, luna, weekendSelectat.Sambata)
	dataDum := ""
	if weekendSelectat.Duminica > 0 {
		dataDum = fmt.Sprintf("%04d-%02d-%02d", an, luna, weekendSelectat.Duminica)
	}

	ocupatSam := make(map[string]string)
	ocupatDum := make(map[string]string)

	for _, abs := range info.Absente {
		if abs.Data == dataSam {
			ocupatSam[abs.AngajatID] = fmt.Sprintf("Concediu (%s)", abs.Cod)
		}
		if dataDum != "" && abs.Data == dataDum {
			ocupatDum[abs.AngajatID] = fmt.Sprintf("Concediu (%s)", abs.Cod)
		}
	}

	for _, t := range info.TurePlanificate {
		if t.Data == dataSam {
			ocupatSam[t.AngajatID] = fmt.Sprintf("Schimb %s", t.CodSchimb)
		}
		if dataDum != "" && t.Data == dataDum {
			ocupatDum[t.AngajatID] = fmt.Sprintf("Schimb %s", t.CodSchimb)
		}
	}

	var liberiTotal []InfoAngajat
	var ocupati []string

	for _, a := range info.Angajati {
		mSam, areSam := ocupatSam[a.ID]
		mDum, areDum := ocupatDum[a.ID]

		if !areSam && !areDum {
			liberiTotal = append(liberiTotal, a)
		} else {
			var detalii []string
			if areSam {
				detalii = append(detalii, fmt.Sprintf("sâm. %s", mSam))
			}
			if areDum {
				detalii = append(detalii, fmt.Sprintf("dum. %s", mDum))
			}
			ocupati = append(ocupati, fmt.Sprintf("- 🔴 **%s %s** (`%s`): %s", a.Nume, a.Prenume, a.Marca, strings.Join(detalii, "; ")))
		}
	}

	var sb strings.Builder
	intervalWeekend := dataSam
	if dataDum != "" {
		intervalWeekend = fmt.Sprintf("%s – %s", dataSam, dataDum)
	}
	sb.WriteString(fmt.Sprintf("🗓️ **Disponibilitate weekend: %s**\n\n", intervalWeekend))

	sb.WriteString(fmt.Sprintf("🟢 **Angajați complet liberi (%d):**\n", len(liberiTotal)))
	if len(liberiTotal) == 0 {
		sb.WriteString("*(Niciun angajat complet liber în ambele zile)*\n")
	} else {
		for _, a := range liberiTotal {
			funcText := ""
			if a.Functie != "" {
				funcText = fmt.Sprintf(" — %s", a.Functie)
			}
			sb.WriteString(fmt.Sprintf("- 👤 **%s %s** (marca `%s`%s)\n", a.Nume, a.Prenume, a.Marca, funcText))
		}
	}

	sb.WriteString(fmt.Sprintf("\n🔴 **Angajați la serviciu sau absenți (%d):**\n", len(ocupati)))
	if len(ocupati) == 0 {
		sb.WriteString("*(Nimeni nu este planificat sau în concediu în acest weekend)*\n")
	} else {
		for _, o := range ocupati {
			sb.WriteString(o + "\n")
		}
	}

	return sb.String()
}

func sintetizeazaAcoperire(info ContextPontix, numeLuna string) string {
	an := info.An
	luna := info.Luna
	if an <= 0 {
		an = time.Now().Year()
	}
	if luna <= 0 {
		luna = int(time.Now().Month())
	}

	ultZi := time.Date(an, time.Month(luna)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	type ZiAcoperire struct {
		Data   string
		Zi     int
		NumeZi string
	}

	var zileWeekend []ZiAcoperire
	for zi := 1; zi <= ultZi; zi++ {
		t := time.Date(an, time.Month(luna), zi, 0, 0, 0, 0, time.UTC)
		if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
			numeZi := "Sâmbătă"
			if t.Weekday() == time.Sunday {
				numeZi = "Duminică"
			}
			data := fmt.Sprintf("%04d-%02d-%02d", an, luna, zi)
			zileWeekend = append(zileWeekend, ZiAcoperire{
				Data:   data,
				Zi:     zi,
				NumeZi: numeZi,
			})
		}
	}

	hartaTure := make(map[string][]string)
	for _, t := range info.TurePlanificate {
		if t.CodSchimb != "" {
			hartaTure[t.Data] = append(hartaTure[t.Data], fmt.Sprintf("%s (%s)", t.Nume, t.CodSchimb))
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🛡️ **Raport verificare acoperire weekend și gărzi (%s %d)**\n\n", numeLuna, an))

	zileAcoperite := 0
	zileNeacoperite := 0

	for _, zw := range zileWeekend {
		echipa := hartaTure[zw.Data]
		if len(echipa) == 0 {
			zileNeacoperite++
			sb.WriteString(fmt.Sprintf("- ⚠️ **%s (%s)**: ❌ **NEACOPERIT!** Nicio gardă sau tură planificată.\n", zw.Data, zw.NumeZi))
		} else {
			zileAcoperite++
			sb.WriteString(fmt.Sprintf("- ✅ **%s (%s)**: Asigurat de %s\n", zw.Data, zw.NumeZi, strings.Join(echipa, ", ")))
		}
	}

	sb.WriteString("\n📊 **Sinteză acoperire secție:**\n")
	sb.WriteString(fmt.Sprintf("- Zile weekend asigurate: **%d / %d**\n", zileAcoperite, len(zileWeekend)))
	if zileNeacoperite > 0 {
		sb.WriteString(fmt.Sprintf("- ⚠️ Zile cu deficit de personal: **%d** (necesită planificare)\n", zileNeacoperite))
	} else {
		sb.WriteString("- 🎉 Toate weekendurile lunii au asigurată prezența personalului!\n")
	}

	return sb.String()
}

func sintetizeazaStatisticiGarzi(info ContextPontix, numeLuna string) string {
	type StatAngajat struct {
		Angajat   InfoAngajat
		NrGarzi   int
		NrTure    int
		TotalOre  int
		Schimburi map[string]int
	}

	harta := make(map[string]*StatAngajat)
	for _, a := range info.Angajati {
		harta[a.ID] = &StatAngajat{
			Angajat:   a,
			Schimburi: make(map[string]int),
		}
	}

	oreSchimb := make(map[string]int)
	for _, s := range info.Schimburi {
		oreSchimb[strings.ToUpper(s.Cod)] = s.Ore
	}

	totalOreSectie := 0
	totalGarziSectie := 0

	for _, t := range info.TurePlanificate {
		codNorm := strings.ToUpper(strings.TrimSpace(t.CodSchimb))
		if codNorm == "" {
			continue
		}
		stat, ok := harta[t.AngajatID]
		if !ok {
			continue
		}
		stat.NrTure++
		stat.Schimburi[codNorm]++

		ore := t.Ore
		if ore <= 0 {
			ore = oreSchimb[codNorm]
		}
		if ore <= 0 {
			switch codNorm {
			case "G24":
				ore = 24
			case "Z12", "N12":
				ore = 12
			case "Z7", "S7":
				ore = 7
			case "IN9", "N9":
				ore = 9
			case "8":
				ore = 8
			}
		}
		stat.TotalOre += ore
		totalOreSectie += ore

		if strings.HasPrefix(codNorm, "G") {
			stat.NrGarzi++
			totalGarziSectie++
		}
	}

	var lista []*StatAngajat
	for _, s := range harta {
		lista = append(lista, s)
	}
	sort.Slice(lista, func(i, j int) bool {
		if lista[i].TotalOre != lista[j].TotalOre {
			return lista[i].TotalOre > lista[j].TotalOre
		}
		return lista[i].Angajat.Nume < lista[j].Angajat.Nume
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 **Statistici gărzi și ore planificate (%s %d)**\n\n", numeLuna, info.An))

	if totalOreSectie == 0 {
		sb.WriteString("Nu există nicio tură sau gardă planificată în această lună.\n")
		return sb.String()
	}

	sb.WriteString("| Angajat | Funcție | Gărzi | Total Ture | Ore Planificate | Detalii Schimburi |\n")
	sb.WriteString("| :--- | :--- | :---: | :---: | :---: | :--- |\n")

	for _, s := range lista {
		if s.NrTure == 0 {
			continue
		}
		var det []string
		for cod, cnt := range s.Schimburi {
			det = append(det, fmt.Sprintf("%dx %s", cnt, cod))
		}
		sort.Strings(det)
		sb.WriteString(fmt.Sprintf("| **%s %s** | %s | %d | %d | **%dh** | %s |\n",
			s.Angajat.Nume, s.Angajat.Prenume, s.Angajat.Functie, s.NrGarzi, s.NrTure, s.TotalOre, strings.Join(det, ", ")))
	}

	sb.WriteString(fmt.Sprintf("\n📌 **Total secție:** %d gărzi, %d ore planificate în total.\n", totalGarziSectie, totalOreSectie))
	return sb.String()
}
