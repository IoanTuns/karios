# D099 · Culorile programelor în rapoartele SQL

Stare: propusă · Data: 2026-11-09 · Taskuri: F09-T36 · Idee: I127

## Context

Culorile programelor (D097) apar în Planificare, în Pontaj și în rapoartele PDF standard. Rapoartele SQL personalizate, de exemplu raportul lunar de gărzi G5 al lui Ioan Tuns, nu primeau culori. Ioan Tuns a cerut I127 pe 2026-11-09.

## Decizie

- **Date.** Tabelele `programe`, `planificare`, `garzi` și `garzi_zile` din setul de date al rapoartelor SQL au o coloană nouă, `culoare`: culoarea afișată a programului (`#RRGGBB`), adică cea aleasă în Setări sau cea implicită după tip. O gardă din pontajul efectiv, fără program, are culoarea implicită a gărzilor.
- **Formular.** Raportul flexibil are două câmpuri noi:
  - **Fundal din coloana**: coloana SQL care dă culoarea rândului;
  - **Colorează coloanele**: coloanele colorate; gol înseamnă tot rândul.
- **Afișare.**
  - Coloana de culoare nu apare în raport.
  - Valorile care nu sunt `#RRGGBB` lasă rândul necolorat.
  - Culorile puse deja de raportul standard (weekend, sărbători, absențe) au prioritate.
  - Culorile apar în PDF-ul editorului vizual (pdfme).
- **Erori.** O coloană de fundal care nu există în rezultatul SQL dă o eroare care o numește (ca la F09-T34).

## Variante respinse

- **Coloana colorată după un nume convențional** (`culoare_<coloana>`), fără câmp în formular: comportament ascuns, greu de descoperit.
- **Culoarea doar în tabelul `programe`**, obținută cu JOIN: corect, dar complicat pentru un raport de gărzi obișnuit.
- **Culori automate** în orice raport SQL care are o coloană cu codul programului: SQL-ul poate redenumi sau combina coloane, deci ghicirea ar fi nesigură.

## Consecințe

- Un SQL salvat cu `SELECT * FROM garzi` (sau `planificare`, `garzi_zile`, `programe`) primește o coloană în plus, `culoare`. Ea se ascunde alegând-o la **Fundal din coloana** sau la **Coloane ascunse**. Modelele implicite folosesc coloane numite și nu se schimbă.
- Rapoartele flexibile au versiunea 19 (D009).
