# F09 · Rapoarte PDF

Etapă: V1 · Stare: în lucru · Depinde de: F04, F08

## Scop

Rapoartele se tipăresc dintr-un PDF generat, cu același aspect pe orice calculator, în limba aleasă pentru rapoarte.

## Taskuri

F09-T08 extinde modelele cu imagini locale în antet și subsol și verifică rapoarte cu una sau mai multe pagini; conținutul continuă automat pe paginile necesare.

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F09-T08 | Imagini în antet/subsol și rapoarte cu una sau mai multe pagini | în lucru | Codex | Cerut de Ioan Tuns la 2026-10-25. Ramură F09-T08-imagini-paginare în worktree-ul rapoartelor; integrează F07-T12 din a132a27 numai aici. PNG/JPEG per zonă, înălțime configurabilă și proporții păstrate, repetare pe pagini, teste PDF/UI. [D026](../decizii/D026-imagini-rapoarte.md) |
| F09-T01 | Prototip: o pagină A3 a centralizatorului cu diacritice, celule colorate și antet repetat | de făcut | | Confirmă `maroto` sau alege altă bibliotecă; fontul TTF trebuie inclus pentru ș, ț, ă |
| F09-T02 | Centralizatorul lunar (A4/A3 peisaj, o pagină pe secție, antet repetat) | de făcut | | Doar pentru V1 |
| F09-T03 | Fișa individuală pe angajat și lună | de făcut | | |
| F09-T04 | Graficul de schimburi | de făcut | | |
| F09-T05 | Textele rapoartelor din cataloage, în limba din `limba_rapoarte` | de făcut | | F03-T09 |
| F09-T06 | Versiunea fiecărui raport, subsolul cu versiuni și testele cu fișier de referință | de făcut | | D009 punctul 2 |
| F09-T07 | Definiții de rapoarte cu antet/conținut/subsol, meniu separat, previzualizare și tipărire PDF | gata | Codex | 0.20.0 include F07-T11 (bacd37b) în ramura izolată. Antet/conținut/subsol independente, trei modele, RO/EN, A4/A3, revizii și export local; verificate Go/vet, TypeScript/Vite, Edge și PDF cu patru pagini. Dialogurile native și tipărirea fizică nu au fost probate interactiv. Pas următor: probă la imprimanta utilizatorului, apoi specificațiile contabilității pentru centralizatorul final F09-T02; editorul și modelele de lucru sunt disponibile. [D025](../decizii/D025-definitii-rapoarte.md), [ghid](../ghid-rapoarte.md) |
