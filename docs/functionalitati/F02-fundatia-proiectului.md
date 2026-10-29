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
| F02-T12 | Curățarea artefactelor vechi și director propriu pentru fiecare build | gata | Codex | 0.22.2 din b563e46, ramură F02-T12-builduri-separate. Eliminate 36 de ținte generate (340,86 MiB), cu directoarele de date, checkouturile și uneltele păstrate. Build real și verificări PowerShell trecute: directoare unice, manifest/hash, fișier existent blocat, temporare și restaurare la eșec. Go/vet trecute. Checkoutul IDE a trecut curat de pe F07-T12 pe ramura taskului; noul build.ps1 este disponibil în rădăcină. Livrare build/bin/pontix-0.22.2-20261026-095522-fa27441b630046eab765d2c924acfc7c/pontix.exe. Pas următor: utilizarea noului build.ps1; dacă se dorește continuarea datelor, copierea întregului karios după închiderea aplicației. [D029](../decizii/D029-builduri-in-directoare-separate.md), [ghid](../ghid-build.md), [inventar](../curatare-artefacte-2026-10-26.md) |
| F02-T13 | Verificări automate pe GitHub (CI): `go vet`, `go test`, `gofmt`, `npm run build` la push și PR | gata | Claude | `.github/workflows/ci.yml`; prima rulare reală are loc la următorul push |
| F02-T14 | `gofmt` în hook-ul de pre-commit | de făcut | | Trei fișiere F16 neformatate au ajuns în repo |
| F02-T15 | Repararea `wails generate module`, care eșuează fără mesaj | de făcut | | Legăturile 0.60.0 au fost scrise manual |
| F02-T11 | Reunirea tuturor ramurilor în main prin PR și buildul versiunii finale | blocat | Codex | Integrarea locală include F09-T09 și F02-T12: ore în pontaj, două rapoarte, configurarea calculelor și builduri separate. Buildul 0.22.2 din b563e46 verificat și livrat în director propriu. Checkoutul IDE este pe F02-T12-builduri-separate; main rămâne nemodificat. Origin este public, contrar AGENTS.md §4; confirmarea trecerii la privat rămâne în așteptare. Nu s-au făcut push, PR sau merge main. Pas următor: răspunsul utilizatorului privind vizibilitatea, reverificarea vârfurilor, apoi push/PR/merge și build din main; datele operaționale sunt păstrate. [D027](../decizii/D027-reunire-main-pr.md), [descriere PR](../pr-F02-T11-reunire-main.md) |
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
