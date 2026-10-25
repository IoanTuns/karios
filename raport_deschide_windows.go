package main

import "golang.org/x/sys/windows"

func deschidePDFAplicatie(cale string) error {
	verb, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	fisier, err := windows.UTF16PtrFromString(cale)
	if err != nil {
		return err
	}
	return windows.ShellExecute(0, verb, fisier, nil, nil, windows.SW_SHOWNORMAL)
}
