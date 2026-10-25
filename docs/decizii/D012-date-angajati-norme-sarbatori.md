# D012 · Datele angajaților, normele și sărbătorile anuale

Stare: propusă · Data: 2026-10-25 · Taskuri: F07-T06, F07-T01, F07-T03

## Cerințe confirmate

Ioan Tuns a cerut explicit în conversația din 2026-10-25:

- Nume, prenume și număr de marcă pentru fiecare angajat.
- Tip contract: normă întreagă sau parțială; implicit întreagă.
- Norma întreagă configurabilă global, implicit 7 ore pe zi.
- Pentru norma parțială, ore pe zi individuale (confirmare explicită în aceeași conversație).
- Numărul anual de zile de concediu conform contractului, fără o valoare implicită inventată (confirmare explicită).
- Sărbători definite pentru fiecare an.

Aceste cerințe sunt confirmate; starea `propusă` privește reprezentarea tehnică de mai jos.

## Propunere de implementare

- Migrarea `0002_angajati_contracte.sql` adaugă `prenume`, `tip_contract` (`INTREAGA` / `PARTIALA`), `norma_partiala_minute` și `concediu_zile_an`.
- Duratele se păstrează în minute întregi, ca în schema inițială. Norma globală este cheia `norma_intreaga_minute` din `setari`, inițial `420`; o vedere SQL expune norma efectivă fiecărui angajat.
- Contractul parțial cere o normă individuală pozitivă; contractul întreg folosește valoarea globală și nu are normă parțială. Nu impunem o regulă juridică despre raportul dintre cele două norme.
- Concediul anual este un număr întreg nenegativ; `NULL` înseamnă necompletat. Nu presupunem un minim legal și nu calculăm solduri sau reportări.
- Numele vechi nu se desparte automat: rămâne în `nume`, iar `prenume` este inițial gol. Vechea coloană de normă este păstrată sub numele `norma_minute_anterioara`, exclusiv pentru verificarea și completarea datelor migrate; nu determină noul tip de contract și nu participă la norma efectivă. Contractele migrate au valoarea implicită `INTREAGA` și trebuie verificate manual.
- Copiem norma veche înainte de eliminarea coloanei originale, în aceeași tranzacție de migrare; angajații noi au `norma_minute_anterioara = NULL`, fără vechiul default de 480 de minute. Nu eliminăm nicio valoare existentă și nu refacem manual legăturile cu pontajul.
- Migrarea `0003_sarbatori_an.sql` adaugă anul derivat din dată și indexat. Data completă rămâne sursa unică: aceeași sărbătoare poate avea date distincte în ani diferiți. Nu precompletăm un calendar legal.
- Nu schimbăm migrarea 0001 și nu adăugăm biblioteci. Ecranele și accesul la date generat cu sqlc rămân în taskurile existente.

## Variante evaluate

- Norma întreagă copiată în fiecare angajat: ar permite divergențe față de setarea globală.
- Deducerea numelui/prenumelui sau a tipului de contract din datele vechi: ar introduce presupuneri neconfirmate.
- An editabil separat de data sărbătorii: ar permite nepotriviri între an și dată.

## Limite și continuare

Schimbarea setării globale afectează vederea normei curente. Istoricul normelor și înghețarea rezultatelor la închiderea lunii trebuie tratate înaintea calculelor istorice (F12), într-un task separat; această schimbare nu implementează reguli de salarizare.
