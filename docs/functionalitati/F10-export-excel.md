# F10 · Export Excel pentru salarizare

Etapă: V1 · Stare: de făcut · Depinde de: F01, F04

## Scop

Datele lunii ajung în programul de salarizare fără să fie retastate.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F10-T01 | Află formatul acceptat de programul de salarizare | de făcut | | Din F01-T02 |
| F10-T02 | Exportul centralizatorului și al foii colective în Excel (`excelize`), cu formule și aspect conform normelor | gata | Antigravity | Implementat în 0.50.0 prin `internal/exportexcel`. Include foaia de pontaj lunar și foaia de centralizator salarii cu formule SUM. |
| F10-T03 | Exportul pentru salarizare în formatul de la F10-T01 | de făcut | | Limba din `limba_rapoarte` |
| F10-T04 | Import din Excel (opțional) | de făcut | | |
| F10-T05 | Versiunea fiecărui export și testele cu fișier de referință | de făcut | | D009 punctul 2 |
