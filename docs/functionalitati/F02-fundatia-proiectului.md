# F02 · Fundația proiectului

Etapă: V1 · Stare: în lucru · Depinde de: —

## Scop

Proiectul are structura, uneltele și regulile de lucru pe care se construiește restul.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — stack-ul tehnologic
- [D005](../decizii/D005-documentare.md) — cum se documentează

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F02-T12 | Curățarea artefactelor vechi și director propriu pentru fiecare build | în lucru | Codex | Cerut de Ioan Tuns la 2026-10-26. Din integrarea 9a4ef86, ramură F02-T12-builduri-separate. Inventar numai de metadate; istoricul Git nu include build/bin, dist sau executabile. Curățare locală a artefactelor generate, protejând checkouturile, uneltele și directoarele cu date. Fiecare build nou într-un director unic, cu temporare Wails izolate și manifest; fără rescriere de istoric sau publicare. [D029](../decizii/D029-builduri-in-directoare-separate.md) |
| F02-T11 | Reunirea tuturor ramurilor în main prin PR și buildul versiunii finale | blocat | Codex | Integrarea locală include F09-T09 până la 4da1617: ore în pontaj, două rapoarte și configurarea calculelor. Fetch origin repetat; toate cele 36 de referințe locale/origin verificate sunt strămoși ai integrării. Go/vet, TypeScript/Vite și Edge trecute; buildul 0.22.1 din 6502c12 verificat și livrat în build/bin/pontix-0.22.1-rapoarte. Origin este public, contrar AGENTS.md §4; confirmarea trecerii la privat rămâne în așteptare. Nu s-au făcut push, PR sau merge main. Pas următor: răspunsul utilizatorului privind vizibilitatea, reverificarea vârfurilor, apoi push/PR/merge și build din main; checkoutul principal și datele operaționale rămân neatinse. [D027](../decizii/D027-reunire-main-pr.md), [descriere PR](../pr-F02-T11-reunire-main.md) |
| F02-T01 | Modulul Go `pontix`, cod în `internal/` | gata | Claude | `go.mod`, 2026-10-24 |
| F02-T02 | Instalează Go permanent pe calculatorul de dezvoltare | gata | Claude | 2026-10-24: Go 1.27.3 în `D:\Apps\Go` (era deja instalat; winget a confirmat că nu există versiune mai nouă) |
| F02-T03 | `CLAUDE.md` cu regulile proiectului | gata | Claude | 2026-10-24 |
| F02-T04 | Structura de documentare: roadmap, funcționalități, decizii, jurnal | gata | Claude | [docs/README.md](../README.md), 2026-10-24 |
| F02-T05 | Primul commit cu codul și documentația | gata | Claude | 2026-10-24, pe ramura `fundatie` |
| F02-T06 | Verificare automată la fiecare push: `gofmt`, `go vet`, `go test` | de făcut | | Cere un remote (de exemplu GitHub) |
| F02-T07 | Reguli de lucru pentru toți agenții (`AGENTS.md`), registrul de idei, coloanele `Responsabil` și `Cine` | gata | Claude | [D006](../decizii/D006-reguli-agenti.md), 2026-10-24 |
| F02-T08 | `.gitattributes` (LF) și `.gitignore` (date reale, baze de date, secrete) | gata | Claude | 2026-10-24 |
| F02-T09 | Regula „mereu ultima versiune stabilă” și prima aplicare: Go 1.27.0 → 1.27.3 | gata | Claude | [D007](../decizii/D007-versiuni-stabile.md), `go.mod` cu `toolchain go1.27.3`, 2026-10-24 |
| F02-T10 | Versionarea aplicației: `internal/versiune`, `docs/versiuni.md`, hook-ul `.githooks/pre-commit`; prima versiune 0.1.0 | gata | Claude | [D009](../decizii/D009-versionare.md), 2026-10-24 |
