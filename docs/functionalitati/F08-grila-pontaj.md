# F08 · Grila lunară de pontaj

Etapă: V1 · Stare: în lucru · Depinde de: F04, F07

## Scop

O grilă angajați × zile, la fel de rapidă ca Excel, în care totalurile se recalculează pe loc.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F08-T01 | Prototip 50 × 31 cu introducere din tastatură, copiere și completare pe intervale | gata | Codex | 2026-10-25, 0.13.0: tabel HTML nativ, selecție Shift, TSV și completare de intervale, schițe persistente cu revizii și jurnal; migrarea 0005. Teste Go și Edge la 50×31, luni/ani bisecți și conflicte. Pas următor: confirmarea F01, apoi conversia schiței în pontaj efectiv; [D016](../decizii/D016-grila-pontaj-schita.md) rămâne propusă |
| F08-T02 | Introducere din tastatură cu coduri scurte | de făcut | | Prototipul F08-T01 acceptă text și sugerează schimburile active; rămâne conversia în intervale și absențe după F01 |
| F08-T03 | Copiere și completare pe intervale | de făcut | | Funcționează în schița F08-T01; validarea semantică a loturilor de pontaj efectiv vine după F01 |
| F08-T04 | Weekendurile și sărbătorile marcate; totaluri recalculate pe loc | de făcut | | Marcarea calendarului există în F08-T01; totalurile rămân necalculate până la confirmarea semnificațiilor și regulilor |
| F08-T05 | Un clic pe celulă arată segmentele și categoriile de tarif | de făcut | | |
| F08-T06 | Afișarea turelor de noapte: împărțite (5 + 7) sau în ziua de început (12) | de făcut | | Depinde de F01-T02 |
