// Package diagnoza gestionează captarea erorilor, jurnalele locale rotative și
// profilarea timpilor de execuție a interogărilor (D043, F16).
package diagnoza

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	// PragInterogareLentaMs este durata peste care o operație SQL este considerată lentă.
	PragInterogareLentaMs = 50.0

	// DimensiuneMaximaJurnalOcteti: 5 MB per fișier de log înainte de rotire.
	DimensiuneMaximaJurnalOcteti = 5 * 1024 * 1024

	// MaxFisierePastrate: câte fișiere de log arhivate se păstrează.
	MaxFisierePastrate = 3
)

// StatisticiDiagnoza conține metricile de sănătate, erori și latență ale aplicației.
type StatisticiDiagnoza struct {
	NumarErori         int     `json:"numar_erori"`
	UltimaEroare       string  `json:"ultima_eroare"`
	TotalInterogari    int64   `json:"total_interogari"`
	TotalLente         int64   `json:"total_lente"`
	LentaInterogare    string  `json:"lenta_interogare"`
	TimpMaximMs        float64 `json:"timp_maxim_ms"`
	TimpMediuMs        float64 `json:"timp_mediu_ms"`
	UltimaInterogare   string  `json:"ultima_interogare"`
	UltimulTimpMs      float64 `json:"ultimul_timp_ms"`
	CaleFolderJurnale  string  `json:"cale_folder_jurnale"`
	FisierJurnalCurent string  `json:"fisier_jurnal_curent"`
	DimensiuneJurnal   string  `json:"dimensiune_jurnal"`
}

// Manager gestionează stocarea pe disc a jurnalului și metadatele de profilare.
type Manager struct {
	mu sync.Mutex

	folderJurnale string
	acum          func() time.Time

	numarErori      int
	ultimaEroare    string
	totalInterogari int64
	totalLente      int64
	sumaTimpMs      float64
	timpMaximMs     float64
	lentaInterogare string
	ultimaOp        string
	ultimulTimpMs   float64
}

var (
	implicitMu sync.RWMutex
	implicit   *Manager
)

// Initializeaza configurează managerul de diagnoză cu folderul de date rădăcină.
func Initializeaza(radacina string) *Manager {
	m := &Manager{
		folderJurnale: filepath.Join(radacina, "jurnale"),
		acum:          time.Now,
	}
	_ = os.MkdirAll(m.folderJurnale, 0o700)

	implicitMu.Lock()
	implicit = m
	implicitMu.Unlock()

	return m
}

// ObtineManager întoarce instanța implicită sau una temporară dacă nu a fost inițializat.
func ObtineManager() *Manager {
	implicitMu.RLock()
	m := implicit
	implicitMu.RUnlock()
	if m == nil {
		m = Initializeaza(os.TempDir())
	}
	return m
}

// Cronometreaza măsoară durata unei interogări sau operații din backend.
// Utilizare:
//
//	stop := diagnoza.Cronometreaza("CitesteAngajati")
//	defer stop(&err)
func Cronometreaza(nume string) func(errRef *error) {
	start := time.Now()
	m := ObtineManager()

	return func(errRef *error) {
		durata := time.Since(start)
		durataMs := float64(durata.Microseconds()) / 1000.0

		var err error
		if errRef != nil {
			err = *errRef
		}
		m.inregistreazaMetrica(nume, durataMs, err)
	}
}

func (m *Manager) inregistreazaMetrica(nume string, durataMs float64, err error) {
	m.mu.Lock()
	m.totalInterogari++
	m.sumaTimpMs += durataMs
	m.ultimaOp = nume
	m.ultimulTimpMs = durataMs

	if durataMs > m.timpMaximMs {
		m.timpMaximMs = durataMs
		m.lentaInterogare = fmt.Sprintf("%s (%.2f ms)", nume, durataMs)
	}
	eLenta := durataMs >= PragInterogareLentaMs
	if eLenta {
		m.totalLente++
	}
	m.mu.Unlock()

	if err != nil {
		m.InregistreazaEroare("interogare/"+nume, err, fmt.Sprintf("durata: %.2f ms", durataMs))
	} else if eLenta {
		m.scrieLinieJurnal("AVERTISMENT", "performanta/"+nume, fmt.Sprintf("Interogare lentă detectată: %.2f ms (prag: %.0f ms)", durataMs, PragInterogareLentaMs))
	}
}

// InregistreazaEroare adaugă o eroare în fișierul local de diagnostic și crește contorul.
func (m *Manager) InregistreazaEroare(sursa string, err error, detalii ...string) {
	if err == nil {
		return
	}
	textDetalii := strings.Join(detalii, " | ")
	mesaj := err.Error()
	if textDetalii != "" {
		mesaj += " [" + textDetalii + "]"
	}

	m.mu.Lock()
	m.numarErori++
	m.ultimaEroare = fmt.Sprintf("[%s] %s", sursa, err.Error())
	m.mu.Unlock()

	m.scrieLinieJurnal("EROARE", sursa, mesaj)
}

// InregistreazaPanic capturează o panică Go cu stivă de execuție completă.
func (m *Manager) InregistreazaPanic(sursa string, p any, stiva []byte) {
	mesaj := fmt.Sprintf("PANIC: %v\nStivă:\n%s", p, string(stiva))

	m.mu.Lock()
	m.numarErori++
	m.ultimaEroare = fmt.Sprintf("[%s] PANIC: %v", sursa, p)
	m.mu.Unlock()

	m.scrieLinieJurnal("CRITIC", sursa, mesaj)
}

// InregistreazaEroareFrontend salvează o eroare survenită în browser / React.
func (m *Manager) InregistreazaEroareFrontend(sursa, mesaj, stack string) {
	text := mesaj
	if stack != "" {
		text += "\nStack JS:\n" + stack
	}

	m.mu.Lock()
	m.numarErori++
	m.ultimaEroare = fmt.Sprintf("[Frontend/%s] %s", sursa, mesaj)
	m.mu.Unlock()

	m.scrieLinieJurnal("FRONTEND", sursa, text)
}

func (m *Manager) scrieLinieJurnal(nivel, sursa, continut string) {
	if m.folderJurnale == "" {
		return
	}
	_ = os.MkdirAll(m.folderJurnale, 0o700)

	acum := m.acum()
	numeFisier := fmt.Sprintf("pontix-%04d-%02d.log", acum.Year(), acum.Month())
	caleFisier := filepath.Join(m.folderJurnale, numeFisier)

	// Rotire dacă depășește limita
	if fi, err := os.Stat(caleFisier); err == nil && fi.Size() >= DimensiuneMaximaJurnalOcteti {
		caleRotita := filepath.Join(m.folderJurnale, fmt.Sprintf("pontix-%04d-%02d-%d.log", acum.Year(), acum.Month(), acum.Unix()))
		_ = os.Rename(caleFisier, caleRotita)
	}

	linie := fmt.Sprintf("[%s] [%s] [%s] %s\n", acum.Format("2006-01-02 15:04:05.000"), nivel, sursa, continut)

	f, err := os.OpenFile(caleFisier, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString(linie)
}

// Statistici generează sumarul pentru interfață.
func (m *Manager) Statistici() StatisticiDiagnoza {
	m.mu.Lock()
	defer m.mu.Unlock()

	medie := 0.0
	if m.totalInterogari > 0 {
		medie = m.sumaTimpMs / float64(m.totalInterogari)
	}

	acum := m.acum()
	numeJurnal := fmt.Sprintf("pontix-%04d-%02d.log", acum.Year(), acum.Month())
	caleJurnal := filepath.Join(m.folderJurnale, numeJurnal)

	dimensiune := "0 B"
	if fi, err := os.Stat(caleJurnal); err == nil {
		dimensiune = formateazaDimensiune(fi.Size())
	}

	return StatisticiDiagnoza{
		NumarErori:         m.numarErori,
		UltimaEroare:       m.ultimaEroare,
		TotalInterogari:    m.totalInterogari,
		TotalLente:         m.totalLente,
		LentaInterogare:    m.lentaInterogare,
		TimpMaximMs:        m.timpMaximMs,
		TimpMediuMs:        medie,
		UltimaInterogare:   m.ultimaOp,
		UltimulTimpMs:      m.ultimulTimpMs,
		CaleFolderJurnale:  m.folderJurnale,
		FisierJurnalCurent: numeJurnal,
		DimensiuneJurnal:   dimensiune,
	}
}

// UltimeleLinii întoarce ultimele linii din jurnalul lunii curente. Un jurnal lipsă dă
// fs.ErrNotExist; textul pentru utilizator îl alege apelantul (F03-T10).
func (m *Manager) UltimeleLinii(numarLinii int) (string, error) {
	acum := m.acum()
	caleJurnal := filepath.Join(m.folderJurnale, fmt.Sprintf("pontix-%04d-%02d.log", acum.Year(), acum.Month()))

	continut, err := os.ReadFile(caleJurnal)
	if err != nil {
		return "", err
	}
	linii := strings.Split(strings.TrimSpace(string(continut)), "\n")
	if len(linii) > numarLinii {
		linii = linii[len(linii)-numarLinii:]
	}
	return strings.Join(linii, "\n"), nil
}

// RecupereazaPanic este o funcție utilitară pentru defer în funcțiile apelate prin Wails:
//
//	defer diagnoza.RecupereazaPanic("App.Metoda")
func RecupereazaPanic(sursa string) {
	if r := recover(); r != nil {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		ObtineManager().InregistreazaPanic(sursa, r, buf[:n])
	}
}

func formateazaDimensiune(octeti int64) string {
	const unitate = 1024
	if octeti < unitate {
		return fmt.Sprintf("%d B", octeti)
	}
	div, exp := int64(unitate), 0
	for n := octeti / unitate; n >= unitate; n /= unitate {
		div *= unitate
		exp++
	}
	val := float64(octeti) / float64(div)
	litere := []string{"KB", "MB", "GB"}
	if exp < len(litere) {
		return fmt.Sprintf("%.1f %s", val, litere[exp])
	}
	return fmt.Sprintf("%.1f MB", val)
}
