package conturi

import (
	"os"

	"golang.org/x/sys/windows"
)

// IdentitateCurenta întoarce SID-ul utilizatorului Windows și numele calculatorului.
func IdentitateCurenta() Identitate {
	id := Identitate{}
	id.Calculator, _ = os.Hostname()
	if u, err := windows.GetCurrentProcessToken().GetTokenUser(); err == nil {
		id.SID = u.User.Sid.String()
	}
	return id
}
