# D090 · Rapoartele standard ca tabele calculate în SQL

Stare: propusă · Data: 2026-11-06 · Taskuri: F09-T21

## Context

D087 cere ca toate rapoartele să fie generate din SQL, cu calculele rămase în Go. Rapoartele standard (pontaj lunar, pontaj detaliat, fișă individuală, evidența anuală a absențelor, grafic gărzi) au fiecare câte un calcul Go propriu, cu teste de referință: norma, totalurile, coloanele configurate, culorile de weekend și absență, rândurile de grup și legenda. Rescrierea acestor calcule în SQL ar muta reguli de pontaj în afara testelor (contrar D087 punctul 3 și AGENTS.md §2).

## Decizie propusă

1. **Calculul rămâne cel existent.** Partea din `Serviciu.Pregateste` care construiește documentul unei surse standard devine o funcție separată, apelată atât de motorul vechi, cât și de setul de date SQL. Nu se schimbă nicio regulă.
2. **Rezultatul calculului devine tabel în setul de date SQL** (D082), cu numele sursei: `pontaj_lunar`, `pontaj_detaliat`, `fisa_individuala`, `absente_anuale`, `grafic_garzi`. Un tabel se calculează numai dacă SQL-ul raportului îl folosește. Calculul folosește filtrele raportului (lună, secție, compartiment, angajat, filtre profesionale).
3. **Coloanele:** `rand` (poziția rândului în documentul calculat), `grup` (1 pentru rândurile de grup), apoi coloanele documentului, cu codurile existente. Zilele devin `z01`…`z31`, ca SQL-ul să nu depindă de lună; în lunile mai scurte coloanele în plus sunt goale. Valorile sunt textele afișate de motorul vechi, deci identice prin construcție.
4. **Prezentarea trece prin `rand`.** Dacă rezultatul SQL conține coloana `rand` și raportul folosește un singur tabel standard, aplicația preia din documentul calculat: titlurile coloanelor, culorile pe celulă (weekend, absență), culoarea pe coloană, rândurile de grup și aldine și legenda (explicațiile coloanelor). `rand` se ascunde în tabelul afișat.
5. **Modelele predefinite** (faza 2): câte un model SQL pentru fiecare sursă standard, cu șablon pdfme, adăugat prin migrare. Modelele vechi rămân până la conversia din F09-T23.
6. **Coloanele zilelor** în PDF (faza 2): Go trimite în `pontix_stiluri` indicii coloanelor `z29`…`z31` inexistente în luna raportului, iar desenarea tabelului le scoate și redistribuie lățimile. Culorile de pe coloane și celule se trimit tot prin `pontix_stiluri`.

## Variante respinse

- **Calculul rescris în SQL** (norma, totalurile pe reguli configurate, gărzile): ar duplica reguli de pontaj în afara `internal/calcul` și a testelor de referință.
- **Tabele brute și pivotare în SQL** (`MAX(CASE WHEN zi=…)` pentru 31 de zile): SQL lung, fragil, și tot fără culori sau totaluri pe reguli.
- **Câte un tabel pe lună, cu coloane după dată:** codurile s-ar schimba lunar și un SQL salvat n-ar mai funcționa luna următoare.

## Consecințe

- Un raport standard poate fi adaptat din SQL (coloane alese, ordine, filtre, coloane calculate) fără cod Go, iar valorile rămân cele testate.
- Datele noi în rapoartele standard cer în continuare cod Go și teste.
- Metadatele de prezentare se pierd dacă SQL-ul nu păstrează `rand` sau combină mai multe tabele standard; raportul rămâne corect, doar fără culori.
