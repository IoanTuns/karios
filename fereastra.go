package main

// dimensiuniFereastra converteste zona utila in unitati logice Wails.
func dimensiuniFereastra(latime, inaltime, dpi int) (int, int, int, int) {
	if dpi <= 0 {
		dpi = 96
	}
	if latime <= 0 || inaltime <= 0 {
		return 800, 600, 640, 400
	}
	l, h := max(1, latime*96/dpi-64), max(1, inaltime*96/dpi-64)
	return min(1280, l), min(800, h), min(1024, l), min(640, h)
}
