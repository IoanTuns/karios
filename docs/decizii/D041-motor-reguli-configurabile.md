# D041 · Arhitectură pentru motorul de reguli configurabile (Rule Engine)

Stare: propusă · Data: 2026-10-26 · Propusă de: Antigravity · Task: F11-T08, I061

## Context

Utilizatorul a solicitat evaluarea creării unui proces care să ofere libertate de configurare pentru reguli care nu sunt încă definite sau hardcodate în sistem:
> *„Crezi că putem crea un proces care să ofere libertate de configurare și pentru regulile care nu sunt definite?”*

În reglementările muncii din România (Codul Muncii, contracte colective la nivel de ramură/unitate, regulamente interne de ordine interioară), fiecare companie sau spital are adesea politici proprii:
- Număr minim de ore de repaus între ture (ex: 24h, 48h după o gardă de 24h).
- Limită de nopți sau weekenduri consecutive lucrate.
- Reguli de rotație specifice pentru sărbători (Paști, Crăciun, Revelion, 1 Mai, Rusalii).
- Proporții și plafoane de ore suplimentare / gărzi pe lună sau trimestru.
- Niveluri diferite de severitate dorite de conducere: informativ, avertisment sau blocant.

Hardcodarea fiecărei reguli particulare în Go limitează flexibilitatea și încalcă principiul enunțat în `AGENTS.md` („Regulile de plată nu se inventează... vin doar de la contabilitate / utilizator”).

## Arhitectură propusă

1. **Blocuri logice de construcție (Building Blocks):**
   În loc de cod Go rigid, definim o structură de regulă compusă din:
   - **Domeniu:** `PLANIFICARE_GARZI`, `PLANIFICARE_TURE`, `PONTAJ_SALARIZARE`.
   - **Condiție (Trigger):**
     - `SUCCESIUNE`: Tura A urmată de Tura B într-un interval de N ore / zile.
     - `REPAUS_MINIM`: Minim H ore libere după tura specificată (ex: 24h după `G24`).
     - `PLAFON_LUNAR`: Maxim N ture / ore dintr-un set de coduri pe lună per angajat.
     - `ACOPERIRE_MINIMA`: În zilele de tip T (`WEEKEND`, `SARBATOARE`, `LUCRATOARE`), minim N angajați cu schimbul S.
     - `ROTATIE_EVENIMENTE`: Angajatul alocat la Evenimentul 1 nu poate / are prioritate la Evenimentul 2.
   - **Severitate:**
     - `INFO`: Notă informativă în sumarul angajatului.
     - `AVERTISMENT`: Atenționare vizuală în grilă (inspecție), permițând salvarea/publicarea.
     - `BLOCANT`: Interzice salvarea sau publicarea dacă regula este încălcată.

2. **Stocare flexibilă în baza de date:**
   - Tabel `reguli_configurabile` (goose migration):
     `id`, `cod`, `denumire`, `descriere`, `domeniu`, `severitate`, `conditii_json`, `activ`, `revizie`.

3. **Evaluator unificat:**
   - Un evaluator generic în `internal/calcul/reguli.go` care primește lista de reguli active și contextul (planificare, date calendaristice, istoric) și întoarce o listă de rezultate / încălcări.
   - În interfața React: feedback instantaneu în grilă și în sumarul angajatului, plus un ecran de configurare vizuală a regulilor din Setări.
