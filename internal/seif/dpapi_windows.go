package seif

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// entropieW leagă copia W de Pontix: alt program al aceluiași utilizator Windows nu o
// deschide doar apelând DPAPI fără această valoare.
var entropieW = []byte("pontix/seif/W/D054")

func blob(b []byte) *windows.DataBlob {
	if len(b) == 0 {
		return &windows.DataBlob{}
	}
	return &windows.DataBlob{Size: uint32(len(b)), Data: &b[0]}
}

func dinBlob(out windows.DataBlob) []byte {
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return append([]byte(nil), unsafe.Slice(out.Data, out.Size)...)
}

// protejeaza criptează cu DPAPI pentru utilizatorul Windows curent.
func protejeaza(date []byte) ([]byte, error) {
	var out windows.DataBlob
	if err := windows.CryptProtectData(blob(date), nil, blob(entropieW), 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, err
	}
	return dinBlob(out), nil
}

func deprotejeaza(date []byte) ([]byte, error) {
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(blob(date), nil, blob(entropieW), 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, err
	}
	return dinBlob(out), nil
}

// sidCurent întoarce SID-ul utilizatorului Windows al procesului.
func sidCurent() (string, error) {
	u, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return "", err
	}
	return u.User.Sid.String(), nil
}
