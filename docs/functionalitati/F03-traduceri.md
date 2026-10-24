# F03 · Traduceri

Etapă: V1 · Stare: în lucru · Depinde de: F02

## Scop

Interfața, mesajele, rapoartele și datele introduse de utilizator sunt disponibile cel puțin în română și engleză. Traducerile se pot corecta și se pot adăuga limbi noi fără recompilare.

## Decizii

- [D004](../decizii/D004-traduceri.md) — totul se poate traduce și configura

## Taskuri

| ID | Task | Stare | Note |
| --- | --- | --- | --- |
| F03-T01 | Pachetul `i18n` cu cataloagele RO și EN incluse în `.exe` | gata | `internal/i18n/`, 2026-10-24 |
| F03-T02 | Erori traductibile în motor (`calcul.Eroare`) | gata | `internal/calcul/eroare.go`, 2026-10-24 |
| F03-T03 | Folder de traduceri lângă `.exe`: corecturi, limbi noi, lista traducerilor lipsă | gata | `AdaugaFolder`, `Limbi`, `Lipsa`, 2026-10-24 |
| F03-T04 | Tipul `Text` pentru denumirile introduse de utilizator | gata | `i18n.Text`, 2026-10-24 |
| F03-T05 | Setările `limba_interfata`, `limba_rapoarte` și folderul de traduceri în `config.toml` | de făcut | Odată cu F05-T04 |
| F03-T06 | Interfața React folosește aceleași cataloage (react-i18next) | de făcut | Odată cu F05-T02 |
| F03-T07 | Formatarea datelor și a numerelor după limbă | de făcut | 24.10.2026 / 10/24/2026; virgulă / punct |
| F03-T08 | Ecran de setări pentru limbă și traducerile lipsă | de făcut | Folosește `i18n.Lipsa` |
| F03-T09 | Rapoartele PDF și Excel în limba aleasă | de făcut | Se face în F09 și F10 |
