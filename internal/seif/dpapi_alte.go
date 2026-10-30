//go:build !windows

package seif

// Pe alte sisteme nu există DPAPI: datele se deblochează cu o parolă sau cu cheia de recuperare.
func protejeaza([]byte) ([]byte, error)   { return nil, errWindowsIndisponibil }
func deprotejeaza([]byte) ([]byte, error) { return nil, errWindowsIndisponibil }
func sidCurent() (string, error)          { return "", errWindowsIndisponibil }
