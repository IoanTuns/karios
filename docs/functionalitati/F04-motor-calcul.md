# F04 · Motorul de calcul

Etapă: V1 · Stare: în lucru · Depinde de: F01, F02

## Scop

Fiecare tură introdusă se împarte automat pe categorii de tarif, corect și în cazurile grele: peste miezul nopții, gărzi de 24 h, sărbători, schimbarea orei.

## Decizii

- [D003](../decizii/D003-motor-calcul.md) — principiile motorului
- [D002](../decizii/D002-portabilitate.md) — fusul orar inclus în `.exe`

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F04-T09 | Corectarea citirii calendarului pe ani în adaptorul de calcul al pontajului | gata | Claude | I054, găsit la testul F07-T14: pontaj_calcul.go apelează ListaSarbatori(ctx,nil), dar interogarea filtrează an=?. Motorul pur primește corect calendarul dacă este furnizat; adaptorul trebuie să încarce anii de la granița lunii. Cazuri necesare: sărbătoare în cursul săptămânii și tură peste an, concordanță grilă/PDF. Task separat conform regulii de scop. Rezolvat in F07-T31 (9c78447, 2026-10-31): pontaj_calcul.go foloseste ListaToateSarbatori, fara filtru de an. Acoperit de TestPontajCitesteSarbatorileFaraFiltruDeAn (sarbatoare in timpul saptamanii) si TestDetaliatFolosesteAceleasiOreSiReguliCaGrila (grila/raport). Inchis 2026-11-06 de Claude dupa verificare; nu exista test dedicat pentru o tura care trece peste schimbarea anului. |
| F04-T01 | `Segmenteaza`, `Valideaza` și calendarul cu sărbători | gata | Claude | `internal/calcul/`, 2026-10-24 |
| F04-T02 | Teste: exemplele din documentul de design, gărzi, sărbători, limite de interval, schimbarea orei | gata | Claude | 20 de teste, toate trec |
| F04-T03 | Modelul de sporuri: categorie de bază + sporuri cumulate (de exemplu noapte peste weekend) | anulat | | Așteaptă F01-T02; probabil o decizie nouă Anulat 2026-11-09: Pontix nu calculeaza salarii ([D100](../decizii/D100-fara-calcul-salarii.md)). |
| F04-T04 | Orele suplimentare, calculate separat față de normă | blocat | | Așteaptă F01-T02 |
| F04-T05 | Scăderea pauzei din segmentele turei | blocat | | Așteaptă F01-T02 |
| F04-T06 | Reguli versionate în timp, inclusiv procentele | anulat | | Procentul stă acum în `categorii_ore`, care nu e versionat Anulat 2026-11-09: procentele nu se calculeaza ([D100](../decizii/D100-fara-calcul-salarii.md)); versionarea normelor ramane in I010. |
| F04-T07 | Testele de referință din luna reală primită de la contabilitate | blocat | | Așteaptă F01-T02 |
| F04-T08 | Ora ambiguă (03:00–04:00 în noaptea în care ora se dă înapoi) la introducere | blocat | Se rezolva in interfata; motorul primeste momente exacte. Analiza 2026-11-09 (Claude): aplicatia nu are campuri pentru ora exacta; orele vin din programe cu durata nominala, fara ajustare la schimbarea orei (ghid-configurare-pontaj.md). Ramane o regula de ore, nu o problema de introducere: intrebarea e la F01 („Schimbarea orei”). Blocat pana la raspuns. |
