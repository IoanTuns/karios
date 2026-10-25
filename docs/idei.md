# Idei

Orice idee se scrie aici, chiar dacă nu se face acum. O idee nu se șterge: când se decide ce se întâmplă cu ea, i se schimbă starea și se leagă de taskul sau decizia rezultată.

Stări: `nouă`, `în analiză`, `transformată` (în task sau decizie), `amânată`, `respinsă`. O idee amânată sau respinsă spune de ce.

| ID | Data | Idee | Cine | Stare | Legături / motiv |
| --- | --- | --- | --- | --- | --- |
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
