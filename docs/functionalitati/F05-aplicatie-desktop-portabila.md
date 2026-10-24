# F05 · Aplicația desktop portabilă

Etapă: V1 · Stare: în lucru · Depinde de: F02

## Scop

Un singur `.exe` care pornește pe orice Windows, fără instalare și fără drepturi de administrator, și ține totul în folderul lui.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — Wails + React
- [D002](../decizii/D002-portabilitate.md) — regulile pentru portabilitate (propusă)

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F05-T01 | Stabilește ultima versiune stabilă de Wails și verifică ce implică pentru proiect | gata | Claude | 2026-10-24: **Wails v2.17.0** (v3 e încă beta, exclus de D007; vezi I007). Node.js **24.22.1 LTS**. Wails v2 are opțiunea `WebviewUserDataPath` cerută de D002 |
| F05-T02 | Scheletul Wails + React + TypeScript | gata | Claude | 2026-10-24, versiunea 0.2.0: Wails 2.17.0, React 19.3.1, TypeScript 7.0.2, Vite 8.3.5; `pontix.exe` 11,5 MB. Nepornit încă: întâi F05-T03 (altfel WebView2 scrie în `%APPDATA%`) |
| F05-T03 | Datele în `karios/` lângă `.exe`; cache-ul WebView2 în `karios/webview/` | gata | Claude | 2026-10-24, versiunea 0.3.0: `internal/portabil`, `WebviewUserDataPath`. Verificat: după 8 s de rulare, niciun folder nou în `%APPDATA%` sau `%LOCALAPPDATA%`; cache-ul în `karios/webview/` |
| F05-T04 | Setările aplicației | de făcut | | Include F03-T05. Dacă D008 e acceptată: în baza criptată, nu în `config.toml` |
| F05-T05 | Verificarea dreptului de scriere în `karios/`, cu mesaj tradus | gata | Claude | 2026-10-24, versiunea 0.4.0: fereastră Windows de eroare în limba Windows. Verificat pe `.exe`: folder fără drept de scriere (ACL) și fișier `karios` care blochează folderul; în ambele cazuri apare mesajul și nu se creează nimic. Modul WAL a trecut la F06-T02 |
| F05-T06 | Backup cu `VACUUM INTO`, automat înainte de migrări și la închiderea lunii | de făcut | | |
| F05-T07 | Probă pe un PC Windows curat și de pe un stick USB | de făcut | | |
| F05-T08 | Semnarea digitală a `.exe` | de făcut | | De ales furnizorul și costul; fără ea apar avertismente SmartScreen |
| F05-T09 | Prototip SQLite criptat cu `ncruces/go-sqlite3`: WAL, `VACUUM INTO`, `sqlc`, performanță | gata | Claude | 2026-10-24, versiunea 0.3.1: toate verificările trec pe `adiantum` și `xts`; propus `xts`. Limite: fără detectarea modificărilor (backup-uri semnate), criptare deterministă. Rezultate în D008 |
| F05-T10 | Seiful: cheia de date cu copiile A (cheia din `.exe`), P și R (pe cont) | de făcut | | D008 punctul 2; conturile sunt în F15 |
| F05-T11 | Fișierul de blocare: a doua instanță pornește doar în citire | de făcut | | D002 punctul 6 (aprobat) |
| F05-T12 | Gestionarea cheii aplicației: generare, păstrare în afara repo-ului, injectare la build | de făcut | | Aceeași cheie la toate versiunile; pierderea ei = date vechi ilizibile. Cine o păstrează: de decis |
| F05-T13 | Script de build: versiunea din `versiune.go` în proprietățile `.exe` și commit-ul prin `-ldflags` | gata | Claude | 2026-10-24, versiunea 0.3.3: `build.ps1`; versiunea în proprietățile `.exe` (verificată în Explorer: 0.3.3, Pontix, copyright) și commit-ul prin `-ldflags`. Scriptul verifică numele și versiunea prin Windows Shell, aceeași sursă ca Proprietăți din Explorer |
| F05-T14 | Semnătura backup-urilor (HMAC-SHA256), verificată înainte de restaurare | de făcut | | D008 punctul 4: criptarea nu detectează modificările |
