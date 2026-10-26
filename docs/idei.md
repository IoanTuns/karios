# Idei

Orice idee se scrie aici, chiar dacă nu se face acum. O idee nu se șterge: când se decide ce se întâmplă cu ea, i se schimbă starea și se leagă de taskul sau decizia rezultată.

Stări: `nouă`, `în analiză`, `transformată` (în task sau decizie), `amânată`, `respinsă`. O idee amânată sau respinsă spune de ce.

| ID | Data | Idee | Cine | Stare | Legături / motiv |
| --- | --- | --- | --- | --- | --- |
| I022 | 2026-10-26 | Reunirea tuturor ramurilor în main printr-un PR și construirea ultimei versiuni | Ioan Tuns | transformată | F02-T11, [D027](decizii/D027-reunire-main-pr.md) |
| I021 | 2026-10-25 | Rapoarte cu una sau mai multe pagini și imagini configurabile în antet/subsol | Ioan Tuns | transformată | F09-T08, [D026](decizii/D026-imagini-rapoarte.md) |
| I020 | 2026-10-25 | Separarea inițializării bazei de date de generarea legăturilor Wails și izolarea temporarelor de build per checkout | Codex | nouă | F09-T07: generatorul Wails rulează main dintr-un executabil temporar; main deschide baza înaintea wails.Run. Două ramuri pot ajunge la același folder temporar karios și la istorice de migrare diferite. Soluție temporară în această sesiune: TEMP/TMP proprii doar pentru procesul de build; corecția pornirii necesită task separat |
| I019 | 2026-10-25 | Fiecare raport are antet, conținut și subsol configurabile separat, într-un meniu dedicat, cu previzualizare și tipărire | Ioan Tuns | transformată | F09-T07, [D025](decizii/D025-definitii-rapoarte.md) |
| I018 | 2026-10-25 | Revizii de concurență separate de timestamp, ca două salvări în același tic al ceasului să rămână distincte | Codex | transformată | F06-T07; observat prin două teste de conflict la verificarea F07-T10 |
| I017 | 2026-10-25 | Programe reale 07–14/15–22/22–07 și intrare 22–07 distinctă pentru weekend/sărbătoare, urmată de opțiuni de 12/24 ore | Ioan Tuns | transformată | F07-T10, [D021](decizii/D021-programe-lucru.md) |
| I016 | 2026-10-25 | Evidență concedii de odihnă, medicale și libere speciale, cu tipuri speciale configurabile | Ioan Tuns | transformată | F07-T04, [D020](decizii/D020-concedii-absente.md) |
| I015 | 2026-10-25 | Fiecare secție are mai multe departamente; schimburile inițiale sunt cele din PDF | Ioan Tuns | transformată | F07-T08, F07-T09 |
| I014 | 2026-10-25 | Build-ul să păstreze întotdeauna datele portabile și să permită alt nume de executabil când aplicația rulează | Codex | transformată | F05-T15, D017; `-clean` poate șterge date din `build/bin/karios` |
| I013 | 2026-10-25 | Modelarea persoanei cu mai multe contracte/mărci și rânduri separate de gardă, plus grupări profesionale ca în raport | Codex | în analiză | F01: captura arată rânduri GARDĂ cu mărci proprii; relația contractuală nu este confirmată. Nu deducem legătura din nume |
| I012 | 2026-10-25 | Sărbători legale România 2026 cu interfață și începutul grilei după modelul de raport primit | Ioan Tuns | transformată | F07-T03, F08-T01; numai structura anonimizată a raportului este descrisă în documentație |
| I011 | 2026-10-25 | Primul flux complet pentru administrarea angajaților și un set de date demo fictive, încărcat la cerere | Ioan Tuns | transformată | F06-T03, F07-T01, F07-T07; [D013](decizii/D013-primele-functionalitati-demo.md) |
| I009 | 2026-10-25 | Date angajați separate, tip contract implicit întreg, normă globală 7 ore, normă parțială individuală, concediu anual contractual și sărbători pe ani | Ioan Tuns | transformată | F07-T06, [D012](decizii/D012-date-angajati-norme-sarbatori.md) |
| I010 | 2026-10-25 | Istoricul normelor și al dreptului contractual la concediu, cu perioade de valabilitate | Codex | nouă | Înaintea calculelor istorice și a închiderii lunii (F12); setarea curentă nu trebuie folosită pentru recalculări istorice |
| I008 | 2026-10-25 | Cheia aplicației în afara codului, injectată la build dintr-un fișier păstrat separat | Claude | amânată | Decizia D011: deocamdată în cod. Mutarea se poate face fără recriptarea bazei |
| I007 | 2026-10-24 | Trecerea la Wails v3 când apare prima versiune stabilă | Claude | amânată | Pe 2026-10-24, v3 era la 3.0.0-beta.34; D007 cere versiuni stabile. v3 schimbă API-ul, deci va fi un task separat |
| I006 | 2026-10-24 | Etapa 2 a criptării: cheia de date deschisă doar cu parola sau răspunsurile unui cont; fără cheia din `.exe` și fără auto-login | Claude | amânată | Copiile P și R există de la început, deci trecerea = ștergerea copiei A ([D008](decizii/D008-criptarea-datelor.md)) |
| I005 | 2026-10-24 | Folderul de date pe un disc din rețea, atașat de utilizator | Ioan Tuns | amânată | Se adaugă ulterior ca opțiune a utilizatorului ([D002](decizii/D002-portabilitate.md) punctul 4). Atenție: WAL nu merge pe rețea, risc de corupere cu mai mulți utilizatori |
| I004 | 2026-10-24 | Avertisment când după o gardă de 24 h nu urmează ziua liberă obligatorie | Claude | nouă | Depinde de răspunsul la întrebarea 3.4 din chestionar (F01-T02) |
| I003 | 2026-10-24 | GitHub Issues legate de ID-urile de task | Claude | amânată | Repo-ul nu are remote; vezi [D005](decizii/D005-documentare.md) |
| I002 | 2026-10-24 | Coduri scurte diferite pe limbă în grilă (de exemplu `CM` / `SL`) | Claude | respinsă | Ar complica importul și exportul; se poate redeschide la cerere ([D004](decizii/D004-traduceri.md)) |
| I001 | 2026-10-24 | Răspunsurile contabilității introduse direct ca teste automate (tabelul de cazuri din chestionar) | Claude | transformată | F01-T03, F04-T07 |
