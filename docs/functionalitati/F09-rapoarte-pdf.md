# F09 · Rapoarte PDF

Etapă: V1 · Stare: în lucru · Depinde de: F04, F08

## Scop

Rapoartele se tipăresc dintr-un PDF generat, cu același aspect pe orice calculator, în limba aleasă pentru rapoarte.

## Taskuri

F09-T08 extinde modelele cu imagini locale în antet și subsol și verifică rapoarte cu una sau mai multe pagini; conținutul continuă automat pe paginile necesare.

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F09-T09 | Ore în pontaj, centralizator și raport detaliat, reguli configurabile | gata | Codex | 0.22.1 din 6502c12, migrare 0015, ramură F09-T09-pontaj-ore din integrarea F02-T11. Coduri păstrate, ore afișate, alegere program sau ore manuale; atribuire între zile configurabilă și previzualizare fără salvare. T total, ZNL weekend/sărbători/libere 00–24, Ore norm luni–vineri fără acele zile, CO/CM active; celelalte surse și factori sunt expliciți în Configurare pontaj. Centralizator v3, fișă v3, detaliat v2, absențe v2; referințele vechi păstrate. Go/vet, TypeScript/Vite, Edge, pdfcpu/pypdf strict și randarea Windows trecute. Executabil verificat, livrat în build/bin/pontix-0.22.1-rapoarte și lângă proba 0.20.0; baza utilizatorului nu a fost deschisă. Pas următor: regulile F09-T10/F01 și proba fizică de tipărire la utilizator. [D028](../decizii/D028-pontaj-ore-doua-rapoarte.md), [ghid](../ghid-configurare-pontaj.md) |
| F09-T10 | Confirmarea formulelor CEV/CFS/Form prof/L-FIN/Deleg/Gardă 200% | blocat | Codex | Ioan Tuns nu are încă regulile; configurarea este oferită de F09-T09, fără a deduce surse, unități sau coeficienți. Întrebarea și răspunsurile sunt în F01. |
| F09-T08 | Imagini în antet/subsol și rapoarte cu una sau mai multe pagini | gata | Codex | 0.21.0, ramură F09-T08-imagini-paginare; include F07-T12 și nota de predare bf7b4e7. PNG/JPEG în cele șase zone, înălțime și proporții, persistență/duplicare/eliminare, repetare pe pagini și număr real de pagini. Migrarea 0014 și actualizarea strictă a probei 0.20.0 păstrează modelele; backup testat înaintea completării 0012. Go/vet, TypeScript/Vite, Edge, pdfcpu, pypdf strict și randarea PDF au trecut. Executabil verificat din a390773, livrat separat și lângă proba 0.20.0 fără deschiderea datelor. Pas următor: probă de tipărire fizică de către utilizator și specificațiile contabilității pentru centralizatorul final. [D026](../decizii/D026-imagini-rapoarte.md), [ghid](../ghid-rapoarte.md) |
| F09-T01 | Prototip: o pagină A3 a centralizatorului cu diacritice, celule colorate și antet repetat | de făcut | | Confirmă `maroto` sau alege altă bibliotecă; fontul TTF trebuie inclus pentru ș, ț, ă |
| F09-T02 | Centralizatorul lunar (A4/A3 peisaj, o pagină pe secție, antet repetat) | de făcut | | Doar pentru V1 |
| F09-T03 | Fișa individuală pe angajat și lună | de făcut | | |
| F09-T04 | Graficul de schimburi | de făcut | | |
| F09-T05 | Textele rapoartelor din cataloage, în limba din `limba_rapoarte` | de făcut | | F03-T09 |
| F09-T06 | Versiunea fiecărui raport, subsolul cu versiuni și testele cu fișier de referință | de făcut | | D009 punctul 2 |
| F09-T07 | Definiții de rapoarte cu antet/conținut/subsol, meniu separat, previzualizare și tipărire PDF | gata | Codex | 0.20.0 include F07-T11 (bacd37b) în ramura izolată. Antet/conținut/subsol independente, trei modele, RO/EN, A4/A3, revizii și export local; verificate Go/vet, TypeScript/Vite, Edge și PDF cu patru pagini. Dialogurile native și tipărirea fizică nu au fost probate interactiv. Pas următor: probă la imprimanta utilizatorului, apoi specificațiile contabilității pentru centralizatorul final F09-T02; editorul și modelele de lucru sunt disponibile. [D025](../decizii/D025-definitii-rapoarte.md), [ghid](../ghid-rapoarte.md) |
