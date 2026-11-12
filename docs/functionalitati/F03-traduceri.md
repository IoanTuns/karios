# F03 · Traduceri

Etapă: V1 · Stare: în lucru · Depinde de: F02

## Scop

Interfața, mesajele, rapoartele și datele introduse de utilizator sunt disponibile cel puțin în română și engleză. Traducerile se pot corecta și se pot adăuga limbi noi fără recompilare.

## Decizii

- [D004](../decizii/D004-traduceri.md) — totul se poate traduce și configura

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F03-T01 | Pachetul `i18n` cu cataloagele RO și EN incluse în `.exe` | gata | Claude | `internal/i18n/`, 2026-10-24 |
| F03-T02 | Erori traductibile în motor (`calcul.Eroare`) | gata | Claude | `internal/calcul/eroare.go`, 2026-10-24 |
| F03-T03 | Folder de traduceri lângă `.exe`: corecturi, limbi noi, lista traducerilor lipsă | gata | Claude | `AdaugaFolder`, `Limbi`, `Lipsa`, 2026-10-24 |
| F03-T04 | Tipul `Text` pentru denumirile introduse de utilizator | gata | Claude | `i18n.Text`, 2026-10-24 |
| F03-T05 | Setările `limba_interfata` și `limba_rapoarte` | gata | Claude | 2026-11-12, la cererea lui Ioan Tuns: reținute în stocarea locală a WebView-ului din karios/webview, ca tema, nu în `config.toml` ([D103](../decizii/D103-limba-retinuta-local.md), acceptată de Ioan Tuns pe 2026-11-12), 0.92.95, ramura F03-T08-limba-setari. Limba interfeței revine la pornire; limba implicită a rapoartelor („La fel ca interfața”, română, engleză) e folosită de Rapoarte și de rapoartele flexibile. Folderul de traduceri a trecut la F03-T11. Teste: vitest limba.test.ts, Edge (scenariul implicit). |
| F03-T06 | Interfața React folosește aceleași cataloage (react-i18next) | gata | Claude | Schelet gata în 0.2.0: cataloagele vin din Go (`Catalog`, `Limbi`). Lipsește schimbarea limbii din setări (F03-T05) 2026-10-30: schimbarea limbii e în interfață (`App.tsx`, `SeteazaLimbaInterfata`); setările din `config.toml` rămân la F03-T05. |
| F03-T07 | Formatarea datelor și a numerelor după limbă | gata | | 24.10.2026 / 10/24/2026; virgulă / punct Reconciliere 2026-11-09 (Claude, F02-T26): formatarea datelor si numerelor dupa limba exista (formatData, Intl.NumberFormat in interfata). |
| F03-T08 | Ecran de setări pentru limbă și traducerile lipsă | gata | Claude | 2026-11-12: fila **Limbă** din Setări › General (limba interfeței, limba rapoartelor, textele netraduse pe limbă din `i18n.Lipsa`, prin `App.TraduceriLipsa`). [D103](../decizii/D103-limba-retinuta-local.md), 0.92.95. Edge: fila deschisă, alegerea limbii rapoartelor reținută. |
| F03-T09 | Rapoartele PDF și Excel în limba aleasă | gata | | Se face în F09 și F10 Reconciliere 2026-11-09 (Claude, F02-T26): rapoartele PDF si Excel au limba aleasa in cerere (Cerere.limba), cu cataloagele RO/EN. |
| F03-T10 | Mutarea în cataloage a textelor scrise direct în cod (Pontaj, Toast, export Excel, mesajele motorului AI, detaliile de dezvoltare) (UI-4) | gata | Claude | 0.75.0: frontend + dialogurile native + erorile AI; `SeteazaLimbaInterfata`. Rămân în română răspunsurile asistentului AI (pachetul `internal/ai`) |
| F03-T11 | Importul traducerilor corectate din folderul de lângă `.exe` în baza criptată (D004) | gata | Claude | Desprins din F03-T05 la 2026-11-12 (Claude). 2026-11-12, la cererea lui Ioan Tuns, 0.92.108, ramura F03-T11-import-traduceri, [D106](../decizii/D106-import-traduceri.md) (propusă): folderul `karios/traduceri`, fișiere `<limba>.json`; în fila Limbă: fișierul de lucru, verificarea, importul (doar administratorul) și ștergerea pe limbă. Migrarea 0039 (`traduceri`), stratul de corecturi din `i18n` (`SeteazaCorecturi`), încărcat la deschiderea bazei; interfața își reîncarcă textele după import și după deblocare. Teste: i18n (corecturi, limbă nouă), datebaza (import, revenire, probleme, jurnal), App (folder, admin, fișier invalid), Edge (scenariul implicit). |
