# D002 · Reguli pentru portabilitate

Stare: propusă · Data: 2026-10-24 · Taskuri: F05-T03 … F05-T07

## Context

Documentul de design cere ca toate datele să stea lângă `.exe`, fără drepturi de administrator. Revizuirea din 2026-10-24 a găsit cinci locuri unde, implicit, regula nu se respectă sau datele se pot pierde.

## Decizie

1. **WebView2:** `WebviewUserDataPath` se setează la un folder lângă `.exe`. Altfel Wails scrie în profilul utilizatorului.
2. **Fusul orar:** baza de fusuri orare e inclusă în `.exe` (`import _ "time/tzdata"`), pentru că Windows nu are una pe care Go s-o poată folosi. *Aplicat deja în `internal/calcul/fus.go`.*
3. **Backup:** doar cu `VACUUM INTO`. Copierea fișierului `.db` cu aplicația deschisă nu e sigură în modul WAL.
4. **Disc din rețea:** aplicația detectează o cale de rețea și trece pe `journal_mode=DELETE` sau refuză pornirea. WAL nu funcționează pe discuri din rețea, nici pentru un singur utilizator.
5. **Folder fără drept de scriere:** dacă folderul `.exe` nu permite scrierea (de exemplu `Program Files`), aplicația spune clar asta la pornire.

## Variante respinse

- **Date în `%APPDATA%`** — contrazice cerința de portabilitate.

## Consecințe

- Mutarea folderului mută tot: date, setări, traduceri, cache-ul WebView2.
- Autentificarea viitoare (F13) păstrează tokenurile în Windows Credential Manager, deci nu se mută cu folderul. E o excepție acceptată.
