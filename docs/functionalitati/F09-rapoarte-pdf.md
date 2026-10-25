# F09 · Rapoarte PDF

Etapă: V1 · Stare: în lucru · Depinde de: F04, F08

## Scop

Rapoartele se tipăresc dintr-un PDF generat, cu același aspect pe orice calculator, în limba aleasă pentru rapoarte.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F09-T01 | Prototip: o pagină A3 a centralizatorului cu diacritice, celule colorate și antet repetat | de făcut | | Confirmă `maroto` sau alege altă bibliotecă; fontul TTF trebuie inclus pentru ș, ț, ă |
| F09-T02 | Centralizatorul lunar (A4/A3 peisaj, o pagină pe secție, antet repetat) | de făcut | | Doar pentru V1 |
| F09-T03 | Fișa individuală pe angajat și lună | de făcut | | |
| F09-T04 | Graficul de schimburi | de făcut | | |
| F09-T05 | Textele rapoartelor din cataloage, în limba din `limba_rapoarte` | de făcut | | F03-T09 |
| F09-T06 | Versiunea fiecărui raport, subsolul cu versiuni și testele cu fișier de referință | de făcut | | D009 punctul 2 |
| F09-T07 | Definiții de rapoarte cu antet/conținut/subsol, meniu separat, previzualizare și tipărire PDF | în lucru | Codex | Implementare 0.19.0 verificată în Go, Edge și PDF; urmează integrarea commiturilor F07-T11 în ramura izolată, verificarea comună și executabilul. Antet/conținut/subsol independente, trei modele, RO/EN, A4/A3, revizii și export local. Tipărirea fizică nu a fost probată. [D024](../decizii/D024-definitii-rapoarte.md), [ghid](../ghid-rapoarte.md) |
