package ai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestServerStrainNuEsteMotorulPropriu(t *testing.T) {
	strain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer strain.Close()
	port, _ := strconv.Atoi(strain.URL[strings.LastIndex(strain.URL, ":")+1:])
	if !VerificaSanatate(port) {
		t.Fatal("serverul de probă nu răspunde")
	}
	m := NouManager(t.TempDir(), port)
	if m.EsteActiv() {
		t.Fatal("un proces străin care răspunde pe port a fost luat drept motorul aplicației")
	}
}

func TestClientulTrimiteCheiaAPI(t *testing.T) {
	var antet string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		antet = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"ok"}}]}`))
	}))
	defer srv.Close()
	c := NouClientCuURL(srv.URL).CuCheie("secret123")
	if _, err := c.Interogheaza(context.Background(), []MesajConversatie{{Rol: "user", Continut: "salut"}}, ContextPontix{An: 2026, Luna: 10}); err != nil {
		t.Fatal(err)
	}
	if antet != "Bearer secret123" {
		t.Fatalf("antet Authorization: %q", antet)
	}
}

func TestPornireRefuzaModelCuCale(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{NumeBinarServer(), "model.gguf"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	m := NouManager(dir, 0)
	for _, nume := range []string{"../model.gguf", `..\secret.gguf`, "altul.gguf"} {
		if err := m.Porneste(context.Background(), nume); err == nil || !strings.Contains(err.Error(), "nu există") {
			t.Fatalf("%q: %v", nume, err)
		}
	}
}

func TestPortLiber(t *testing.T) {
	p, err := portLiber()
	if err != nil || p <= 0 {
		t.Fatalf("port: %d %v", p, err)
	}
}
