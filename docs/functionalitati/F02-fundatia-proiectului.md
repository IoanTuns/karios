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
| F02-T21 | Ajutor contextual pe toate ecranele si manual interactiv din AI | gata | Codex | Integrare main f15e623 in 5438776, aplicatie 0.88.14, manual 1.3. Patru conflicte rezolvate cumulativ; citari functionale si din drawer Pontaj. gofmt/vet/suita Go, 13 teste frontend, 3 probe manual, UI RO/EN cu drawer trecute. Build oficial complet in build/bin/pontix-0.88.14-20261030-191707-c61cb74bc51b48739cddf17ed7422ce8; PDF valid, 19 pagini. Pas urmator: review PR si proba operator WebView2. D057 ramane propusa. |
| F02-T17 | Build din main și manual de utilizare în română, cu propunere pentru ajutor contextual/AI | gata | Codex | Main 4600f30, aplicație 0.88.7 compilată și verificată, fără modificări de cod. Manual ediția 1.0, 19 capitole, PDF valid de 13 pagini și HTML cu căutare verificată. Go/vet și cele 13 teste Vitest trecute. [Manual](../manual-utilizare-ro.md), [reproducere](../manual/README.md), D057 propusă. Pas următor: evaluare manual și proiectarea integrării contextuale; fără merge/push. 2026-10-30: manualul a fost reluat în F02-T18 (generat la fiecare build, D058); în main au intrat D057 și propunerea, fără HTML/PDF. |
| F02-T20 | Worktree permanent pentru Codex | gata | Codex | Cerere explicita a utilizatorului: de acum Codex lucreaza in D:/karios/karios-codex. Creat din main df8ecbc pe F02-T20-worktree-codex; pentru fiecare task nou se foloseste o ramura proprie in acest worktree. Celelalte worktree-uri si modificarile existente raman neatinse. Pas urmator: continuarea lucrului exclusiv aici. |
| F02-T18 | Manuale regenerate, versionate și livrate cu fiecare executabil | gata | Codex | 0.88.11 din commit curat b886a38; manual ediția 1.1, MD/HTML/PDF lângă executabil, versiune/commit/data/hashuri în manifest, control editorial înainte de build. Generator 3/3, probe build succes/erori, Edge HTML, PDF valid 14 pagini, Go/vet, frontend 13/13 și build Windows trecute. Livrare build/bin/Pontix-0.88.11-RO-b886a38.zip. Integrare în main conform cererii anterioare, fără push. Pas următor: la fiecare versiune se revizuiește sursa și se actualizează manual.json; ajutorul contextual/AI rămâne separat. D058 propusă. |
| F02-T19 | Scripturile din `verificari/` opresc toate procesele Edge pe care le pornesc | gata | Claude | `opresteBrowserul` în `contrast-grila.mjs`, `interfata.mjs`, `masurare-grila.mjs`: după rulare, 0 procese rămase (înainte 14 pe rulare de browser) |
| F02-T12 | Curățarea artefactelor vechi și director propriu pentru fiecare build | gata | Codex | 0.22.2 din b563e46, ramură F02-T12-builduri-separate. Eliminate 36 de ținte generate (340,86 MiB), cu directoarele de date, checkouturile și uneltele păstrate. Build real și verificări PowerShell trecute: directoare unice, manifest/hash, fișier existent blocat, temporare și restaurare la eșec. Go/vet trecute. Checkoutul IDE a trecut curat de pe F07-T12 pe ramura taskului; noul build.ps1 este disponibil în rădăcină. Livrare build/bin/pontix-0.22.2-20261026-095522-fa27441b630046eab765d2c924acfc7c/pontix.exe. Pas următor: utilizarea noului build.ps1; dacă se dorește continuarea datelor, copierea întregului karios după închiderea aplicației. [D029](../decizii/D029-builduri-in-directoare-separate.md), [ghid](../ghid-build.md), [inventar](../curatare-artefacte-2026-10-26.md) |
| F02-T13 | Verificări automate pe GitHub (CI): `go vet`, `go test`, `gofmt`, `npm run build` la push și PR | gata | Claude | `.github/workflows/ci.yml`; prima rulare reală are loc la următorul push |
| F02-T14 | `gofmt` în hook-ul de pre-commit | gata | Claude | Verifică versiunea din index (`git show :fișier | gofmt -l`), nu directorul de lucru |
| F02-T15 | Repararea `wails generate module`, care eșuează fără mesaj | gata | Claude | Comanda merge; eșecul inițial nereproductibil. Regula: legăturile se regenerează, nu se editează manual (ghid-build) |
| F02-T16 | Teste automate pentru interfață (Vitest) | gata | Claude | 0.86.1: [D053](../decizii/D053-teste-interfata-vitest.md); funcții pure în `parola.ts`, `repetari.ts`, `jurnalDiferente.ts`; `npm test` în CI |
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

| F02-T22 | Pachet ZIP, acceptanta operator, note de versiune si verificari traduceri/legaturi | gata | Codex | 0.88.16, cod 4bd9d7f; ZIP explicit 10 fisiere/11.9 MiB, manual 1.5/PDF 20 pagini, documente operator si SHA256. Go/vet, frontend 13 teste, 3 probe manual, 7 probe livrare si regresia buildului trecute. CI extins. Pas urmator: review si integrare impreuna cu munca celuilalt agent; fara push. |
