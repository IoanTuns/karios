package ai

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"sync"
	"time"
)

const ContextConfigurat = 32768

// ModelActiv expune numai numele procesului propriu, fără căi sau cheia API.
func (m *Manager) ModelActiv() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cmd == nil {
		return ""
	}
	select {
	case <-m.terminat:
		return ""
	default:
		return filepath.Base(m.model)
	}
}

// Manager se ocupă de pornirea, monitorizarea și oprirea procesului local llama-server.
//
// Serverul pornește pe un port liber ales la fiecare pornire, cu o cheie API aleatoare
// (F14-T12, SEC-2): alte procese locale nu-l pot folosi, iar aplicația trimite date doar
// procesului pe care l-a lansat ea, nu oricui răspunde pe un port cunoscut.
type Manager struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	terminat chan struct{} // se închide când procesul copil se oprește
	port     int
	cheieAPI string
	folder   string
	model    string
	activ    bool
}

// NouManager inițializează un manager pentru procesul local.
func NouManager(folderAI string, port int) *Manager {
	if port <= 0 {
		port = PortImplicit
	}
	return &Manager{
		folder: folderAI,
		port:   port,
	}
}

// Port întoarce portul pe care rulează serverul.
func (m *Manager) Port() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.port
}

// EsteActiv verifică dacă procesul lansat de aplicație rulează și răspunde.
func (m *Manager) EsteActiv() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.alPropriuActiv()
}

func (m *Manager) alPropriuActiv() bool {
	if m.cmd == nil {
		return false
	}
	select {
	case <-m.terminat:
		return false
	default:
	}
	return VerificaSanatate(m.port)
}

// Client întoarce un client pentru serverul propriu, cu cheia API.
func (m *Manager) Client() *Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return NouClient(m.port).CuCheie(m.cheieAPI)
}

// portLiber cere sistemului un port liber pe 127.0.0.1.
func portLiber() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func cheieNoua() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}

// Porneste lansează procesul local în fundal fără consolă vizibilă.
func (m *Manager) Porneste(ctx context.Context, numeModel string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.alPropriuActiv() {
		m.activ = true
		return nil
	}
	_ = m.opresteIntern() // un proces propriu care nu mai răspunde

	stare := DetecteazaMediu(m.folder, 0)
	if !stare.Instalat {
		return fmt.Errorf("mediul AI nu este complet: %s", stare.Mesaj)
	}

	caleModel := stare.CaleModel
	if numeModel != "" {
		// Doar un model din folderul AI, fără căi (SEC-3).
		if numeModel != filepath.Base(numeModel) || !slices.Contains(stare.ModeleDisponibile, numeModel) {
			return fmt.Errorf("modelul %q nu există în folderul AI", numeModel)
		}
		caleModel = filepath.Join(m.folder, numeModel)
	}
	port, err := portLiber()
	if err != nil {
		return fmt.Errorf("nu s-a găsit un port local liber: %w", err)
	}
	m.port, m.cheieAPI = port, cheieNoua()

	// Calculăm numărul optim de fire CPU pentru inferență
	nrFire := runtime.NumCPU()
	if nrFire > 2 {
		nrFire -= 2
	}
	if nrFire < 2 {
		nrFire = 2
	}

	// Argumente optime pentru llama-server pe mașină desktop (D042):
	// -m: calea către fișierul GGUF
	// --port: portul de ascultare
	// --host: strict 127.0.0.1 (fără ascultare în rețea)
	// -c 32768: context pentru instrucțiuni, personal, pontaj și conversație (F14-T20).
	// -b 512, -ub 512: evaluare prompt vectorială accelerată
	// -np 1: un singur slot pentru a maximiza viteza de inferență și a evita fragmentarea memoriei
	// --threads: numărul optim de nuclee alocate
	args := []string{
		"-m", caleModel,
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", m.port),
		"-c", strconv.Itoa(ContextConfigurat),
		"-b", "512",
		"-ub", "512",
		"-np", "1",
		"--threads", fmt.Sprintf("%d", nrFire),
		"--api-key", m.cheieAPI,
	}

	cmd := exec.Command(stare.CaleBinar, args...)
	ascundeConsola(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("eroare la lansarea serverului AI: %w", err)
	}

	m.cmd = cmd
	m.model = caleModel
	terminat := make(chan struct{})
	m.terminat = terminat
	go func() {
		_ = cmd.Wait()
		close(terminat)
	}()

	// Așteaptă până când serverul răspunde la /health (până la 45 secunde)
	timeout := time.After(45 * time.Second)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			_ = m.opresteIntern()
			return ctx.Err()
		case <-timeout:
			_ = m.opresteIntern()
			return errors.New("timeout: serverul AI local nu a pornit în timpul alocat")
		case <-terminat:
			m.cmd = nil
			return errors.New("serverul AI local s-a oprit la pornire")
		case <-ticker.C:
			if VerificaSanatate(m.port) {
				m.activ = true
				// Instrucțiunile intră în cache în fundal; o întrebare pusă între timp așteaptă
				// în același slot și pornește apoi de la cache (F14-T29).
				go func(c *Client) { _ = c.Incalzeste(context.Background()) }(NouClient(m.port).CuCheie(m.cheieAPI))
				return nil
			}
		}
	}
}

// Opreste oprește procesul copil al serverului local.
func (m *Manager) Opreste() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.opresteIntern()
}

func (m *Manager) opresteIntern() error {
	if m.cmd == nil || m.cmd.Process == nil {
		m.activ = false
		return nil
	}
	err := m.cmd.Process.Kill()
	m.cmd = nil
	m.activ = false
	return err
}
