package versiune

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	if !regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`).MatchString(Aplicatie) {
		t.Errorf("versiunea %q nu are formatul MAJOR.MINOR.PATCH", Aplicatie)
	}
}

// Versiunea din proprietățile .exe (wails.json) trebuie să fie aceeași.
func TestWailsJSON(t *testing.T) {
	data, err := os.ReadFile("../../wails.json")
	if err != nil {
		t.Fatal(err)
	}
	var cfg struct {
		Info struct {
			ProductVersion string `json:"productVersion"`
		} `json:"info"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Info.ProductVersion != Aplicatie {
		t.Errorf("wails.json are productVersion %q, versiune.go are %q", cfg.Info.ProductVersion, Aplicatie)
	}
}

// Fiecare versiune trebuie să aibă intrarea ei în istoricul versiunilor.
func TestIstoric(t *testing.T) {
	data, err := os.ReadFile("../../docs/versiuni.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "## "+Aplicatie+" ") {
		t.Errorf("docs/versiuni.md nu are o secțiune „## %s · …” pentru versiunea curentă", Aplicatie)
	}
}
