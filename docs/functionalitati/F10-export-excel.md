# F10 · Export Excel pentru salarizare

Etapă: V1 · Stare: de făcut · Depinde de: F01, F04

## Scop

Datele lunii ajung în programul de salarizare fără să fie retastate.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F10-T01 | Află formatul acceptat de programul de salarizare | de făcut | | Din F01-T02 D100: exportul contine doar ore si zile. 2026-11-12: răspuns E1 — programul nu e ales; formatele exacte se adaugă ca modele predefinite după documentația de import a programului ales ([D108](../decizii/D108-export-salarizare-configurabil.md)). |
| F10-T02 | Exportul centralizatorului și al foii colective în Excel (`excelize`), cu formule și aspect conform normelor | gata | Antigravity | Implementat în 0.50.0 prin `internal/exportexcel`. Include foaia de pontaj lunar și foaia de centralizator salarii cu formule SUM. |
| F10-T03 | Exportul pentru salarizare în formatul de la F10-T01 | gata | Claude 2026-11-12, la cererea lui Ioan Tuns (E1/E2 din F01-T02 și „da”), 0.92.112, ramura F10-T03-export-salarizare, [D108](../decizii/D108-export-salarizare-configurabil.md) (propusă): export configurabil CSV/Excel după modele salvate în bază (migrarea 0042, două modele generice). Coloane: identitate (marcă, nume…), perioadă, orice total din Pontaj, text fix; separatori, zecimale (rotunjire 0,5 în sus), ore zecimal sau ore:minute, antet, BOM; în Excel totalurile sunt numere. Un rând pe marcă; valorile necalculate rămân goale și sunt numărate. UI: Setări › General › Export salarizare (editor, doar administratorul salvează) și Pontaj › Mai multe › Export pentru salarizare (doar cu datele salvate). Pachet `internal/exportsalarizare`, versiunea formatului 1.0.0, referință CSV/Excel. Modelele pentru programe anume așteaptă F10-T01. |
| F10-T04 | Import din Excel (opțional) | gata | | Reconciliere 2026-11-09 (Claude, F02-T26): importul angajatilor din Excel/CSV exista (internal/datebaza/import.go). |
| F10-T05 | Versiunea fiecărui export și testele cu fișier de referință | gata | Claude | 2026-11-12, la cererea lui Ioan Tuns, 0.92.96, ramura F10-T05-referinta-excel. Exportul Excel al pontajului (VersiuneExport 1.2.0) are referința `internal/exportexcel/testdata/export-excel-v1.2.0.json`: foile, valorile celulelor și celulele unite, cu date fictive și momentul generării fixat (`OptiuniExport.GeneratLa`). O schimbare de conținut pică testul până la o versiune nouă și o referință nouă (`go test ./internal/exportexcel -actualizeaza-export`). Exportul PDF al rapoartelor avea deja referințe (D009). Textele exportului sunt scrise direct în cod: I130. |
| F10-T06 | Textele exportului Excel prin cataloagele de traduceri, în limba rapoartelor (I130, D004) | gata | Claude | 2026-11-12, la cererea lui Ioan Tuns, 0.92.98, ramura F10-T06-export-excel-tradus. Titlurile, numele foilor, lunile și antetele vin din cheile `export.*` și `rapoarte.luna_*`; numele angajaților și sărbătorile urmează aceeași limbă. Pontaj trimite limba implicită a rapoartelor (D103). Versiunea exportului 1.3.0: referința română e identică cu 1.2.0 (aceleași texte), referința engleză e nouă (`export-excel-en-v1.3.0.json`). Teste: exportexcel (ro, en), app_export. |
| F10-T07 | Antetul coloanei de gardă din exportul Excel (I132) | gata | Claude | 2026-11-12, la cererea lui Ioan Tuns, 0.92.106, ramura F10-T07-antet-ore-garda. Coloana `ore_garda` (Ore gardă 200%, D093) avea în foaia de pontaj antetul „Ore Noapte / Gardă”, iar în centralizator „Ore Gărzi”; ambele au acum numele coloanei din aplicație, „Ore gardă 200%” („On-call hours 200%”). Exportul 1.4.0, referințe noi ro/en (singurele diferențe față de 1.3.0: cele două antete). |
