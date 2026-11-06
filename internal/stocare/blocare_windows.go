//go:build windows

package stocare

import (
	"errors"
	"os"

	"golang.org/x/sys/windows"
)

// incearcaBlocarea blocheaza exclusiv un octet aflat departe de continut (la
// 4 GiB), ca alta instanta sa poata citi cine tine fisierul.
func incearcaBlocarea(f *os.File) (bool, error) {
	ol := windows.Overlapped{Offset: 0, OffsetHigh: 1}
	err := windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &ol)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, windows.ERROR_LOCK_VIOLATION) || errors.Is(err, windows.ERROR_IO_PENDING) {
		return false, nil
	}
	return false, err
}
