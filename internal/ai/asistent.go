package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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
	Mesaj     string             `json:"mesaj"`
	Propuneri []PropunereActiune `json:"propuneri"`
	Eroare    string             `json:"eroare,omitempty"`
}

// ContextPontix oferă asistentului informații despre personal și schimburile definite.
type ContextPontix struct {
	Limba     string        `json:"limba"`
	An        int           `json:"an"`
	Luna      int           `json:"luna"`
	Angajati  []InfoAngajat `json:"angajati"`
	Schimburi []InfoSchimb  `json:"schimburi"`
}

type InfoAngajat struct {
	ID      string `json:"id"`
	Nume    string `json:"nume"`
	Prenume string `json:"prenume"`
	Marca   string `json:"marca"`
	Functie string `json:"functie"`
}

type InfoSchimb struct {
	Cod      string `json:"cod"`
	Denumire string `json:"denumire"`
	Ore      int    `json:"ore"`
	Interval string `json:"interval,omitempty"`
}

// Client gestionează comunicarea HTTP cu endpoint-ul local al modelului.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NouClient creează un client pentru serverul local (de ex: http://127.0.0.1:28491).
func NouClient(port int) *Client {
	if port <= 0 {
		port = PortImplicit
	}
	return &Client{
		baseURL: fmt.Sprintf("http://127.0.0.1:%d", port),
		httpClient: &http.Client{
			Timeout: 180 * time.Second,
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
				},
				"required": []string{"marca", "data_inceput", "data_sfarsit", "cod_schimb"},
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
	// Construiește promptul de sistem cu instrucțiuni și contextul disponibil
	promptSistem := construiestePromptSistem(info)

	mesaje := []map[string]any{
		{"role": "system", "content": promptSistem},
	}
	for _, m := range istoric {
		mesaje = append(mesaje, map[string]any{
			"role":    m.Rol,
			"content": m.Continut,
		})
	}

	corpCerere := map[string]any{
		"messages":    mesaje,
		"tools":       definitiiUnelte,
		"tool_choice": "auto",
		"temperature": 0.1, // Răspunsuri deterministe și precise
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RaspunsAsistent{}, fmt.Errorf("eroare comunicare cu motorul AI local: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		corpEroare, _ := io.ReadAll(resp.Body)
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

		descriere := construiesteDescrierePropunere(tc.Function.Name, args, info)
		rezultat.Propuneri = append(rezultat.Propuneri, PropunereActiune{
			ID:                 tc.ID,
			Tip:                tc.Function.Name,
			Descriere:          descriere,
			Parametri:          args,
			NecesitaConfirmare: true,
		})
	}

	return rezultat, nil
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

func construiestePromptSistem(info ContextPontix) string {
	var sb strings.Builder
	sb.WriteString("Ești Asistentul inteligent al aplicației Pontix pentru pontaj și planificare de gărzi/ture.\n")
	sb.WriteString("Răspunde întotdeauna concis, politicos și în limba română (sau engleză dacă ți se cere).\n")
	sb.WriteString("NU inventa date și NU calcula formule salariale pe cont propriu.\n")
	sb.WriteString("Dacă utilizatorul îți cere să planifici, să modifici sau să ștergi ture, apelează uneltele puse la dispoziție.\n")
	sb.WriteString("IMPORTANT PENTRU PERIOADE / TOATĂ LUNA:\n")
	sb.WriteString("- Când utilizatorul cere atribuirea unei ture pe o perioadă (ex: 'toată luna', 'toata luna', 'tata luna', 'toată săptămâna', 'de luni până vineri'), folosește neapărat unealta `planifica_interval` sau `sterge_interval`, specificând data_inceput (prima zi a lunii) și data_sfarsit (ultima zi a lunii). NU apela `planifica_tura` pentru o singură zi când se cere toată luna!\n")
	sb.WriteString("- Pentru 'pe dimineață' sau 'dimineața' sau 'de la 7 la 14', alege schimbul de dimineață disponibil (ex: Z7 sau primul schimb de zi).\n")
	sb.WriteString("- Folosește EXACT numărul de marcă așa cum este indicat în lista de personal de mai jos (ex: 'DEMO-008', NU prescurta și nu elimina cratima).\n")
	sb.WriteString(fmt.Sprintf("Anul curent de lucru: %d, Luna: %d.\n", info.An, info.Luna))

	if len(info.Schimburi) > 0 {
		sb.WriteString("Schimburi disponibile: ")
		var sList []string
		for _, s := range info.Schimburi {
			sList = append(sList, fmt.Sprintf("%s (%s, %dh)", s.Cod, s.Denumire, s.Ore))
		}
		sb.WriteString(strings.Join(sList, ", ") + ".\n")
	}

	if len(info.Angajati) > 0 {
		sb.WriteString("Personal identificabil (Folosește marca corectă când apelezi funcții):\n")
		for _, a := range info.Angajati {
			sb.WriteString(fmt.Sprintf("- %s %s (Marca: %s, Funcție: %s)\n", a.Nume, a.Prenume, a.Marca, a.Functie))
		}
	}

	return sb.String()
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
		schimbText := formateazaSchimb(schimb, info.Schimburi)
		angajatText := formateazaAngajat(marca, info.Angajati)
		supliment := ""
		if lucratoare {
			supliment = " (doar zile lucrătoare L–V)"
		}
		return fmt.Sprintf("Planificare schimb %s pentru %s în intervalul %s – %s%s", schimbText, angajatText, d1, d2, supliment)
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
