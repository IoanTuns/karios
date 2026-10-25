# D018 · Departamente în cadrul secțiilor

Stare: propusă · Data: 2026-10-25 · Taskuri: F07-T08

## Cerință confirmată

Ioan Tuns: fiecare secție poate avea mai multe departamente.

## Propunere tehnică

- Tabel `departamente` cu ID, secția părinte obligatorie, denumire `i18n.Text`, stare și coloanele de sincronizare existente în proiect.
- Coloană opțională `angajati.departament_id`; angajații existenți își păstrează secția și rămân fără departament până la atribuire explicită. Nu deducem departamente din funcții sau nume.
- Cheie străină și validare în serviciu/trigger: departamentul angajatului trebuie să aparțină exact secției lui. Schimbarea secției în formular golește selecția de departament.
- În Setări: creare departament sub o secție și listă ierarhică. În formularul angajatului: departamente filtrate după secție. În listă și în grila lunară: filtre pe secție și departament.
- Nu mutăm automat angajați și nu ștergem secții/departamente. Schimburile existente nu sunt modificate de migrarea organizării.
- Migrare nouă, cod SQL generat cu sqlc, operații tranzacționale și jurnal cu actor `local`.

Fără dependențe noi în aplicație. Specificația separată a setărilor inițiale ale schimburilor vine din PDF (F07-T09).

## Citirea documentului local

Pentru extragerea textului folosim auxiliar `pypdf==6.21.2`, ultima versiune stabilă verificată pe [PyPI](https://pypi.org/project/pypdf/), licență BSD-3-Clause. Se descarcă izolat în `build/bin/unelte-pdf`, fără instalare globală și fără includere în dependențele aplicației; PDF-ul local nu este trimis unui serviciu extern.
