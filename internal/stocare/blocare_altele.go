//go:build !windows

package stocare

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

// incearcaBlocarea foloseste flock (consultativ, pe tot fisierul); continutul
// ramane lizibil pentru alta instanta.
func incearcaBlocarea(f *os.File) (bool, error) {
	err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, unix.EWOULDBLOCK) {
		return false, nil
	}
	return false, err
}
