//go:build windows

package ai

import (
	"os/exec"
	"syscall"
)

func ascundeConsola(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}

