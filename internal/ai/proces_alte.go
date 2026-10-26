//go:build !windows

package ai

import "os/exec"

func ascundeConsola(cmd *exec.Cmd) {
	// Pe Linux/macOS nu este necesară ascunderea ferestrei de consolă
}

