package ai

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Manager se ocupă de pornirea, monitorizarea și oprirea procesului local llama-server.
type Manager struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	port   int
	folder string
	model  string
	activ  bool
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

// EsteActiv verifică dacă serverul rulează și răspunde.
func (m *Manager) EsteActiv() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return VerificaSanatate(m.port)
}

// Porneste lansează procesul local în fundal fără consolă vizibilă.
func (m *Manager) Porneste(ctx context.Context, numeModel string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if VerificaSanatate(m.port) {
		m.activ = true
		return nil
	}

	stare := DetecteazaMediu(m.folder, m.port)
	if !stare.Instalat {
		return fmt.Errorf("mediul AI nu este complet: %s", stare.Mesaj)
	}

	caleModel := stare.CaleModel
	if numeModel != "" {
		caleModel = filepath.Join(m.folder, numeModel)
	}

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
	// -c 2048: lungime context optimizată pentru inferență rapidă și economisire RAM (50% mai puțin KV cache)
	// -b 512, -ub 512: evaluare prompt vectorială accelerată
	// -np 1: un singur slot pentru a maximiza viteza de inferență și a evita fragmentarea memoriei
	// --threads: numărul optim de nuclee alocate
	args := []string{
		"-m", caleModel,
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", m.port),
		"-c", "2048",
		"-b", "512",
		"-ub", "512",
		"-np", "1",
		"--threads", fmt.Sprintf("%d", nrFire),
	}

	cmd := exec.Command(stare.CaleBinar, args...)
	ascundeConsola(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("eroare la lansarea serverului AI: %w", err)
	}

	m.cmd = cmd
	m.model = caleModel

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
		case <-ticker.C:
			if VerificaSanatate(m.port) {
				m.activ = true
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
