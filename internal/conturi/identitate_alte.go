//go:build !windows

package conturi

import (
	"os"
	"os/user"
)

// IdentitateCurenta întoarce utilizatorul sistemului și numele calculatorului.
func IdentitateCurenta() Identitate {
	id := Identitate{}
	id.Calculator, _ = os.Hostname()
	if u, err := user.Current(); err == nil {
		id.SID = u.Uid
	}
	return id
}
