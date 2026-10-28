// Package versiune ține versiunea aplicației (D009). Orice commit care schimbă
// codul crește Aplicatie și adaugă o intrare în docs/versiuni.md.
package versiune

// Aplicatie este versiunea aplicației, în format MAJOR.MINOR.PATCH.
const Aplicatie = "0.60.0"

// Commit este commit-ul din care s-a construit .exe; se completează la build:
// -ldflags "-X pontix/internal/versiune.Commit=<hash>".
var Commit = "necunoscut"
