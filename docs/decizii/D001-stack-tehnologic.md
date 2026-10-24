# D001 · Stack-ul tehnologic

Stare: acceptată · Data: 2026-10-23 · Taskuri: F02-T01, F05-T01, F05-T02

> Versiunile fixe de mai jos (Wails v2, maroto v2) sunt înlocuite de [D007](D007-versiuni-stabile.md): se folosește mereu ultima versiune stabilă.

## Context

Aplicația trebuie să fie un singur `.exe` pentru Windows, fără instalare, care lucrează offline și tipărește rapoarte cu același aspect pe orice calculator.

## Decizie

| Strat | Tehnologie |
| --- | --- |
| Aplicație desktop | Wails (Go), cu WebView2 din Windows |
| Logică | Go |
| Interfață | React + TypeScript |
| Bază de date | SQLite prin `modernc.org/sqlite` (Go pur, fără CGO), un fișier lângă `.exe` |
| Acces la date | `sqlc` |
| Rapoarte PDF | `maroto` v2 (de confirmat prin F09-T01) |
| Export Excel | `excelize` |
| Configurare | `config.toml` lângă `.exe` |

Versiunea Wails (v2 sau v3) se confirmă în F05-T01, după starea lor la data începerii.

## Variante respinse

- **Electron** — executabil prea mare.
- **Tauri / Rust** — mai complex, fără un beneficiu real aici.
- **PostgreSQL / MySQL** — au nevoie de server.
- **Fișiere JSON / CSV** — fără relații și fără siguranța datelor.

## Consecințe

- Toată logica de calcul stă în Go și se poate testa fără interfață.
- Biblioteca pentru grilă (F08-T01) și cea pentru PDF (F09-T01) se aleg prin prototip, nu din documentație.

Sursa: documentul de design din 2026-10-23, secțiunea 2.
