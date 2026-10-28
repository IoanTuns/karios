# D044 · Trasabilitate: jurnal înlănțuit, sigiliul datelor și amprenta lunii

Stare: acceptată · Data: 2026-10-28 · Decisă de: Ioan Tuns (alegerile din conversația cu Claude, 2026-10-28) · Taskuri: F06-T04, F06-T08, F06-T09, F12-T04

## Context

Aplicația e portabilă și rulează pe câte un calculator. Cine are acces la sistem și la aplicație are acces la date, deci confidențialitatea față de acesta nu se poate garanta (Ioan Tuns, 2026-10-28). Cerința reală: **datele nu se pot deschide sau modifica fără urmă**, adică orice modificare făcută în afara aplicației trebuie să fie detectată.

Până acum:

- `jurnal_modificari` exista din migrarea 0001 și era scris de majoritatea fluxurilor, dar nu era protejat: un rând se putea șterge sau modifica fără urmă.
- Închiderea lunii și regulile configurabile scriau în afara tranzacției și fără jurnal.
- Cine are cheia de date (azi: oricine are `.exe`, [D011](D011-cheia-aplicatiei-in-cod.md)) putea modifica baza cu alt program fără ca aplicația să observe.

## Decizie

### 1. Cheia de audit

`cheie_audit = HMAC-SHA256(cheia de date, "pontix/audit/v1")`. Nu se stochează nicăieri; se derivă la deschiderea bazei. Când cheia de date va fi deschisă doar cu parola (etapa 2, [I006](../idei.md)), cheia de audit devine automat inaccesibilă fără parolă, fără schimbarea formatului.

### 2. Sigiliul tabelelor

- Pentru fiecare tabel de date, `sigilii(tabel, valoare)` păstrează suma modulo 2^256 a valorilor `HMAC(cheie_audit, tabel ‖ rând)` pentru toate rândurile, cu o codificare canonică pe tipuri (NULL, întreg, real, text, BLOB). Suma, spre deosebire de XOR, nu se anulează la două rânduri identice.
- Aplicația creează pe fiecare conexiune **trigger-e temporare** (`TEMP`, invizibile în fișier) care actualizează sigiliul în aceeași tranzacție cu scrierea. Orice scriere a aplicației e acoperită, indiferent de cod.
- O modificare făcută cu alt program nu are trigger-ele și nici cheia, deci sigiliul nu mai corespunde.
- Tabelele acoperite: toate, mai puțin `sigilii` și tabelele interne `sqlite_*`, inclusiv `goose_db_version`. Un tabel nou dintr-o migrare intră automat.
- Se sigilează și schema (`sqlite_schema`, fără `rootpage`): un trigger, o vedere sau un index adăugat din afară e incident.
- Funcțiile SQL ale auditului sunt `DIRECTONLY`: merg din trigger-ele TEMP ale aplicației, dar nu dintr-un trigger salvat în fișier, care ar putea „repara” sigiliul la următoarea scriere a aplicației.

### 3. Lanțul jurnalului

Fiecare rând din `jurnal_modificari` primește `mac = HMAC(cheie_audit, mac-ul rândului anterior ‖ conținutul rândului)`, calculat de un trigger temporar la inserare. La verificare se găsește primul rând rupt. Ștergerea ultimelor rânduri nu rupe lanțul, dar schimbă sigiliul tabelului `jurnal_modificari`.

### 4. Verificarea la pornire

1. Înainte de migrări, baza se verifică (sigiliile și lanțul).
2. După migrări se resigilează doar dacă verificarea a trecut, iar o migrare nu poate ascunde o modificare externă.
3. Prima sigilare se face o singură dată: când migrarea care introduce sigiliile (0021) se aplică chiar la această pornire. Un tabel `sigilii` gol în orice alt moment este incident.

### 5. Reacția la incident (alegerea Ioan Tuns)

- **Avertizare și blocarea editării.** Banner permanent cu ce s-a detectat. Datele se pot citi și exporta, dar orice scriere e refuzată (`audit.editare_blocata`), atât în serviciu, cât și în trigger.
- **Confirmarea incidentului** cere un motiv, se scrie în jurnal cu detaliile detectate și resigilează starea curentă. Până la conturi (F15), oricine folosește aplicația poate confirma. Risc asumat, notat la F06-T08.

### 6. Amprenta lunii (alegerea Ioan Tuns)

- La închiderea lunii se calculează `HMAC(cheie_audit, "luna" ‖ an ‖ lună ‖ rândurile lunii)` din `pontaj_ciorne`, `absente_zile` și `planificare`, cu rândurile ordonate după cheie. Rezultatul se păstrează în `stare_luni.amprenta`.
- Amprenta (16 caractere hex, grupate câte 4) apare în subsolul tuturor rapoartelor lunii închise. Pe hârtia semnată devine ancora din afara calculatorului: o rescriere completă a bazei, chiar cu cheia, nu poate schimba ce s-a tipărit.
- Ecranul de pontaj arată, pentru luna închisă, amprenta și dacă datele curente încă îi corespund.

### 7. Neinclus acum

Detectarea înlocuirii bazei cu o copie mai veche, validă (contor autentificat în seif) este F06-T09.

## Variante respinse

- **Doar lanțul jurnalului** — nu vede o modificare a datelor care ocolește jurnalul.
- **Coloană de sigiliu pe fiecare rând** — schimbă toate tabelele și tot codul de scriere; sigiliul pe tabel prin trigger-e temporare nu atinge codul existent.
- **Sigiliu recalculat integral la fiecare salvare** — cost liniar cu toată baza la fiecare salvare; actualizarea incrementală a sumei are cost constant.
- **XOR în loc de sumă** — două rânduri identice s-ar anula.
- **Trigger-e permanente în schemă** — ar fi vizibile și ar rula și pentru alt program, care le-ar putea șterge; nu pot apela funcțiile cu cheie ale aplicației.
- **Blocare totală sau doar avertizare** la incident — alese de Ioan Tuns: avertizare și blocarea editării.
- **Amprenta cu capătul lanțului jurnalului** — ar face amprenta dependentă de jurnalul de după închidere; pentru hârtia semnată contează doar datele lunii.

## Consecințe

- Verificarea la pornire citește toată baza o dată. Costul crește liniar cu volumul și trebuie măsurat pe o bază mare (F06-T08).
- Orice scriere directă în bază, făcută intenționat (reparare, script), declanșează incident. Așa e intenționat; după reparare incidentul se confirmă cu motiv.
- Restaurarea manuală a unui backup produce incident numai dacă backup-ul a fost modificat; o copie veche validă nu e detectată până la F06-T09.
- Protecția reală împotriva cuiva care are `.exe` vine abia cu etapa 2 (I006/F15). Până atunci sigiliul oprește modificările naive și pe cele făcute de cine nu cunoaște codul.
- Rapoartele lunare își cresc versiunea: subsolul conține amprenta pentru lunile închise.
