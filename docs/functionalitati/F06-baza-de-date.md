# F06 · Baza de date și migrările

Etapă: V1 · Stare: în lucru · Depinde de: F05

## Scop

Datele stau într-un singur fișier SQLite, iar un fișier dintr-o versiune mai veche se actualizează automat la pornire.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — SQLite, `modernc.org/sqlite`, `sqlc`
- [D004](../decizii/D004-traduceri.md) — denumirile ca `Text` (JSON)

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F06-T01 | Alege sistemul de migrări: `goose` sau `PRAGMA user_version`, nu ambele | gata | Claude | 2026-10-24: propus goose v3.29.0, doar înainte, cu `istoric_migrari` și protecție la downgrade; verificat în prototip ([D010](../decizii/D010-migrari-goose.md)) |
| F06-T02 | Schema inițială: tabelele din documentul de design, denumiri `Text`, coloane de sincronizare; deschiderea bazei criptate în modul WAL | gata | Claude | 2026-10-25, versiunea 0.5.0: `internal/bd` (`Deschide`, `Migreaza`) și migrarea `0001_schema_initiala.sql` cu 11 tabele. Cheia prin `PRAGMA hexkey` la fiecare conexiune; WAL, cheile străine active. Încă nefolosit de `main.go`: are nevoie de seif (F05-T10) |
| F06-T03 | Codul de acces la date generat cu `sqlc` | gata | Codex | 2026-10-25, 0.9.1: `sqlc.yaml`, `genereaza-date.ps1`, interogări și cod generat în `internal/bd/acces`. Prima tranșă: angajați, secții, setări și jurnal; generatorul citește direct migrările 0001–0003. Pas următor: consumarea în F07-T01; interogările pentru celelalte fluxuri se adaugă odată cu ele; [D013](../decizii/D013-primele-functionalitati-demo.md) |
| F06-T04 | `jurnal_modificari` de la prima schemă | de făcut | | Tabelul există din migrarea 0001. F07-T01 scrie tranzacțional pentru angajați/secții/norma globală cu actor `local`; rămâne acoperirea celorlalte fluxuri și identificarea contului prin F15 |
| F06-T05 | Absențele: în `pontaj` sau într-un tabel separat; zi cu ore lucrate și absență parțială | de făcut | | Depinde de F01-T02 |
| F06-T06 | Istoricul migrărilor (dată, versiunea aplicației) și refuzul unei baze migrate de o versiune mai nouă | de făcut | | D009 punctul 3 |
| F06-T07 | Revizii unice independente de rezoluția ceasului pentru schițe și absențe | de făcut | | Descoperit la verificarea F07-T10: două salvări foarte apropiate pot primi aceeași modificat_la și rata conflictul. Se corectează separat, fără schimbarea testelor de referință |
