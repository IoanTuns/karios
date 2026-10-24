# F05 · Aplicația desktop portabilă

Etapă: V1 · Stare: de făcut · Depinde de: F02

## Scop

Un singur `.exe` care pornește pe orice Windows, fără instalare și fără drepturi de administrator, și ține totul în folderul lui.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — Wails + React
- [D002](../decizii/D002-portabilitate.md) — regulile pentru portabilitate (propusă)

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F05-T01 | Stabilește ultima versiune stabilă de Wails și verifică ce implică pentru proiect | de făcut | | După D007; o versiune majoră nouă față de documentul de design se notează în jurnal |
| F05-T02 | Scheletul Wails + React + TypeScript | de făcut | | |
| F05-T03 | Datele în `karios/` lângă `.exe`; cache-ul WebView2 în `karios/webview/` | de făcut | | `WebviewUserDataPath`; D002 |
| F05-T04 | Setările aplicației | de făcut | | Include F03-T05. Dacă D008 e acceptată: în baza criptată, nu în `config.toml` |
| F05-T05 | SQLite în modul WAL; verificarea dreptului de scriere în `karios/` | de făcut | | Discurile din rețea sunt amânate (D002 punctul 4, I005) |
| F05-T06 | Backup cu `VACUUM INTO`, automat înainte de migrări și la închiderea lunii | de făcut | | |
| F05-T07 | Probă pe un PC Windows curat și de pe un stick USB | de făcut | | |
| F05-T08 | Semnarea digitală a `.exe` | de făcut | | De ales furnizorul și costul; fără ea apar avertismente SmartScreen |
| F05-T09 | Prototip SQLite criptat cu `ncruces/go-sqlite3`: WAL, `VACUUM INTO`, `sqlc`, performanță | de făcut | | D008 |
| F05-T10 | Seiful: cheia de date cu copiile A (cheia din `.exe`), P și R (pe cont) | de făcut | | D008 punctul 2; conturile sunt în F15 |
| F05-T11 | Fișierul de blocare: a doua instanță pornește doar în citire | de făcut | | D002 punctul 6 (aprobat) |
| F05-T12 | Gestionarea cheii aplicației: generare, păstrare în afara repo-ului, injectare la build | de făcut | | Aceeași cheie la toate versiunile; pierderea ei = date vechi ilizibile. Cine o păstrează: de decis |
