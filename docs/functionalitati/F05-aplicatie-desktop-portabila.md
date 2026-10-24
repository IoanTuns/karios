# F05 · Aplicația desktop portabilă

Etapă: V1 · Stare: de făcut · Depinde de: F02

## Scop

Un singur `.exe` care pornește pe orice Windows, fără instalare și fără drepturi de administrator, și ține totul în folderul lui.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — Wails + React
- [D002](../decizii/D002-portabilitate.md) — regulile pentru portabilitate (propusă)

## Taskuri

| ID | Task | Stare | Note |
| --- | --- | --- | --- |
| F05-T01 | Alege Wails v2 sau v3, după starea lor actuală | de făcut | Rezultatul se notează în D001 sau într-o decizie nouă |
| F05-T02 | Scheletul Wails + React + TypeScript | de făcut | |
| F05-T03 | Datele și cache-ul WebView2 lângă `.exe` | de făcut | `WebviewUserDataPath` |
| F05-T04 | `config.toml` lângă `.exe` | de făcut | Include setările din F03-T05 |
| F05-T05 | SQLite în modul WAL; detectarea discului din rețea și a folderului fără drept de scriere | de făcut | |
| F05-T06 | Backup cu `VACUUM INTO` | de făcut | |
| F05-T07 | Probă pe un PC Windows curat și de pe un stick USB | de făcut | |
| F05-T08 | Semnarea digitală a `.exe` | de făcut | De ales furnizorul și costul; fără ea apar avertismente SmartScreen |
