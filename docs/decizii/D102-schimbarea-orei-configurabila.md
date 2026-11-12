# D102 · Turele din nopțile cu schimbarea orei: opțiune în setări

Stare: propusă · Data: 2026-11-12 · Taskuri: F04-T08 · Cerută de: Ioan Tuns, 2026-11-12 („Pentru F04-T08 poate fi implementat configurabil?” → „da”)

## Context

Orele din pontaj vin din programe cu oră de început și durată nominală. O tură 19–07 are 12 ore. În noaptea în care ora se dă înapoi (ultima duminică din octombrie, 04:00 devine 03:00), aceeași tură durează efectiv 13 ore. În noaptea în care ora se dă înainte (ultima duminică din martie, 03:00 devine 04:00), durează 11 ore. Până acum Pontix număra durata nominală, iar întrebarea „nominal sau efectiv” era deschisă la F01.

Cum se numără orele e o regulă de ore, nu de plată (D100). Ea nu se inventează: o alege utilizatorul, eventual cu contabilitatea.

## Decizie

O opțiune nouă în **Setări › Configurare pontaj**, lângă repartizarea turelor peste miezul nopții: **Turele din nopțile cu schimbarea orei**.

- **Durata din program** (`NOMINALA`, implicit). Tura are mereu durata din program. Este comportamentul de până acum, deci datele existente nu se schimbă. Configurațiile salvate înainte, fără câmpul nou, se citesc ca `NOMINALA`.
- **Orele efective** (`EFECTIVA`). Tura începe la ora de ceas din program și se termină la ora de ceas a sfârșitului, în fusul Europe/Bucharest. Orele sunt timpul scurs între cele două momente: 13 ore în octombrie, 11 ore în martie.
  - Împărțirea pe zile (repartizarea pe zile calendaristice și ZNL) se face la miezul nopții, ora României.
  - **Ora care se repetă toamna** (03:00–04:00 apare de două ori): un capăt care cade acolo se socotește la prima trecere.
  - **Ora care lipsește primăvara** (03:00–04:00 nu există): un capăt care cade acolo se mută la primul moment care există, adică 04:00, ora de vară.

**Ce nu se schimbă, în niciun mod:**

- orele scrise de mână (de exemplu `8`);
- programele fără oră de început;
- absențele, inclusiv creditul de ore egal cu programul;
- norma și orele lucrătoare;
- verificarea suprapunerilor la planificare;
- gărzile efective din rapoartele flexibile, care deja numără timpul real (D082).

**Lunile închise.** Opțiunea nu se poate schimba cât timp e închisă sau blocată o lună în care ora se schimbă (martie, octombrie). Luna se redeschide întâi, cu motiv, la fel ca la istoricul normei (D101). Celelalte setări ale pontajului se salvează ca până acum.

**Rapoartele.** Când opțiunea este „Orele efective”, legenda rapoartelor de pontaj spune că turele din aceste nopți au orele efective. Rapoartele care folosesc calculul pontajului primesc versiune nouă (D009): Pontaj 23, Pontaj detaliat 26, Fișa individuală 23 și rapoartele flexibile 21. Cu opțiunea implicită, conținutul lor rămâne același.

## Variante respinse

- **Doar durata nominală, fără opțiune.** Ar fi închis întrebarea fără răspunsul utilizatorului; unele unități numără orele efective.
- **Doar orele efective.** Ar fi schimbat orele deja raportate în lunile vechi, fără ca utilizatorul să aleagă.
- **Opțiune valabilă de la o anumită lună.** Ar fi permis schimbarea fără redeschiderea lunilor vechi, dar adaugă un istoric al setărilor pentru un caz care apare de două ori pe an. Refuzul în lunile închise e mai simplu și protejează la fel datele închise.
- **Câmpuri de oră exactă la introducere.** Pontix nu are ore de intrare și ieșire pe celulă, ci programe. Ora ambiguă de la introducere, titlul inițial al F04-T08, nu apare.

## Consecințe

- Cine nu schimbă opțiunea nu vede nicio diferență.
- Cu „Orele efective”, o tură 19–07 din noaptea de 24/25 octombrie 2026 are 13 ore (5 vineri/sâmbătă și 8 duminică, la împărțirea pe zile). Cea din 27/28 martie 2027 are 11 ore.
- Calculul e în `internal/calcul` (`MomentLocal`, `LunaCuSchimbareOra`), cu teste pe ambele nopți, pe orele lipsă și repetate și pe ZNL.
