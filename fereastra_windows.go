package main

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

func dimensiuniPornire() (int, int, int, int) {
	user32 := windows.NewLazySystemDLL("user32.dll")
	var zona struct{ Stanga, Sus, Dreapta, Jos int32 }
	ok, _, _ := user32.NewProc("SystemParametersInfoW").Call(0x0030, 0, uintptr(unsafe.Pointer(&zona)), 0)
	if ok == 0 {
		return dimensiuniFereastra(0, 0, 96)
	}
	dpi := uintptr(96)
	procDPI := user32.NewProc("GetDpiForSystem")
	if procDPI.Find() == nil {
		if valoare, _, _ := procDPI.Call(); valoare > 0 {
			dpi = valoare
		}
	}
	return dimensiuniFereastra(int(zona.Dreapta-zona.Stanga), int(zona.Jos-zona.Sus), int(dpi))
}
