package main

import "golang.org/x/sys/windows"

// arataEroare afișează o fereastră Windows de eroare. Se folosește înainte de
// pornirea interfeței, când aplicația nu poate continua.
func arataEroare(titlu, text string) {
	t, _ := windows.UTF16PtrFromString(titlu)
	m, _ := windows.UTF16PtrFromString(text)
	windows.MessageBox(0, m, t, windows.MB_OK|windows.MB_ICONERROR)
}
