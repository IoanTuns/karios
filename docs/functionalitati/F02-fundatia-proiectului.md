# F02 · Fundația proiectului

Etapă: V1 · Stare: în lucru · Depinde de: —

## Scop

Proiectul are structura, uneltele și regulile de lucru pe care se construiește restul.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — stack-ul tehnologic
- [D005](../decizii/D005-documentare.md) — cum se documentează

## Taskuri

| ID | Task | Stare | Note |
| --- | --- | --- | --- |
| F02-T01 | Modulul Go `pontix`, cod în `internal/` | gata | `go.mod`, 2026-10-24 |
| F02-T02 | Instalează Go permanent pe calculatorul de dezvoltare | de făcut | `winget install GoLang.Go`. Instalarea MSI a eșuat pe 2026-10-24 (alt installer rula); s-a folosit temporar Go 1.27.0 portabil |
| F02-T03 | `CLAUDE.md` cu regulile proiectului | gata | 2026-10-24 |
| F02-T04 | Structura de documentare: roadmap, funcționalități, decizii, jurnal | gata | [docs/README.md](../README.md), 2026-10-24 |
| F02-T05 | Primul commit cu codul și documentația | gata | 2026-10-24, pe ramura `fundatie` |
| F02-T06 | Verificare automată la fiecare push: `gofmt`, `go vet`, `go test` | de făcut | Cere un remote (de exemplu GitHub) |
