# F04 · Motorul de calcul

Etapă: V1 · Stare: în lucru · Depinde de: F01, F02

## Scop

Fiecare tură introdusă se împarte automat pe categorii de tarif, corect și în cazurile grele: peste miezul nopții, gărzi de 24 h, sărbători, schimbarea orei.

## Decizii

- [D003](../decizii/D003-motor-calcul.md) — principiile motorului
- [D002](../decizii/D002-portabilitate.md) — fusul orar inclus în `.exe`
- [D105](../decizii/D105-stergere-tarife.md) — ștergerea tabelelor și a codului de tarif nefolosite
- [D102](../decizii/D102-schimbarea-orei-configurabila.md) — turele din nopțile cu schimbarea orei, opțiune în setări (acceptată, Ioan Tuns)

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F04-T09 | Corectarea citirii calendarului pe ani în adaptorul de calcul al pontajului | gata | Claude | I054, găsit la testul F07-T14: pontaj_calcul.go apelează ListaSarbatori(ctx,nil), dar interogarea filtrează an=?. Motorul pur primește corect calendarul dacă este furnizat; adaptorul trebuie să încarce anii de la granița lunii. Cazuri necesare: sărbătoare în cursul săptămânii și tură peste an, concordanță grilă/PDF. Task separat conform regulii de scop. Rezolvat in F07-T31 (9c78447, 2026-10-31): pontaj_calcul.go foloseste ListaToateSarbatori, fara filtru de an. Acoperit de TestPontajCitesteSarbatorileFaraFiltruDeAn (sarbatoare in timpul saptamanii) si TestDetaliatFolosesteAceleasiOreSiReguliCaGrila (grila/raport). Inchis 2026-11-06 de Claude dupa verificare; nu exista test dedicat pentru o tura care trece peste schimbarea anului. |
| F04-T10 | Ștergerea tabelelor și a codului de tarif nefolosite din D003 (I131) | gata | Claude | 2026-11-12, la cererea lui Ioan Tuns, 0.92.107, ramura F04-T10-stergere-tarife. Migrarea 0038 șterge `categorii_ore`, `seturi_reguli`, `reguli_interval`; șterse `segmenteaza.go`/`segmenteaza_test.go`, `Interval`/`SetReguli` din `reguli.go` și mesajele `reguli.*` ale validării. [D105](../decizii/D105-stergere-tarife.md). Teste: migrările (număr nou, tabelele lipsesc), i18n (alt mesaj de exemplu), toată suita. |
| F04-T01 | `Segmenteaza`, `Valideaza` și calendarul cu sărbători | gata | Claude | `internal/calcul/`, 2026-10-24 2026-11-12: `Segmenteaza` și `Valideaza` șterse ca nefolosite (F04-T10, [D105](../decizii/D105-stergere-tarife.md)); calendarul (`TipZi`, `Calendar`) rămâne. |
| F04-T02 | Teste: exemplele din documentul de design, gărzi, sărbători, limite de interval, schimbarea orei | gata | Claude | 20 de teste, toate trec |
| F04-T03 | Modelul de sporuri: categorie de bază + sporuri cumulate (de exemplu noapte peste weekend) | anulat | | Așteaptă F01-T02; probabil o decizie nouă Anulat 2026-11-09: Pontix nu calculeaza salarii ([D100](../decizii/D100-fara-calcul-salarii.md)). |
| F04-T04 | Orele suplimentare, calculate separat față de normă | blocat | | Așteaptă F01-T02 |
| F04-T05 | Scăderea pauzei din segmentele turei | blocat | | Așteaptă F01-T02 |
| F04-T06 | Reguli versionate în timp, inclusiv procentele | anulat | | Procentul stă acum în `categorii_ore`, care nu e versionat Anulat 2026-11-09: procentele nu se calculeaza ([D100](../decizii/D100-fara-calcul-salarii.md)); versionarea normelor ramane in I010. |
| F04-T07 | Testele de referință din luna reală primită de la contabilitate | blocat | | Așteaptă F01-T02 |
| F04-T08 | Schimbarea orei: durata turelor de noapte în nopțile de trecere la ora de vară/iarnă, configurabilă | gata | Claude | Analiza 2026-11-09: aplicația nu are câmpuri pentru ora exactă; orele vin din programe cu durată nominală, deci rămânea o regulă de ore. 2026-11-12, la cererea lui Ioan Tuns: opțiunea **Turele din nopțile cu schimbarea orei** în Setări › Configurare pontaj (implicit durata din program; „Orele efective” dă 13 h în octombrie și 11 h în martie pentru 19–07), refuzată cât timp e închisă o lună martie/octombrie. [D102](../decizii/D102-schimbarea-orei-configurabila.md) (acceptată de Ioan Tuns, 2026-11-12), 0.92.93, ramura F04-T08-schimbarea-orei. Teste: calcul (5: nopțile din oct./mar. în ambele repartizări, ZNL, ore lipsă/repetate, lunile cu schimbare, validare), datebaza (1), legenda rapoartelor (1), Edge (scenariul implicit). Rapoarte: Pontaj 23, Pontaj detaliat 26, Fișa individuală 23, flexibile 21. Neverificat pe date reale. |
