package versiune

import (
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
