# F11-T21 - Mai multe programe pe aceeasi marca si zi

## Cerinta confirmata de utilizator, 2026-11-04

Aceeasi persoana, aceeasi marca, aceeasi zi: se pot adauga doua sau mai multe programe. Al doilea se adauga, nu il inlocuieste pe primul. Suprapunerile sunt interzise. Exemplu:07:00-14:00 si19:00-07:00 a doua zi. Intervalele care doar se ating la capete nu se suprapun. Nu se deduc reguli de spor sau de repaus noi.

## Audit al implementarii existente

- `internal/bd/migrari/0017_planificare.sql`: index unic activ pe angajat/data. Nu se modifica migrarea existenta; implementarea necesita migrare noua.
- `internal/datebaza/planificare.go`: validare explicita a duplicatelor pe zi; publicarea scrie cate un cod in aceeasi ciorna. Publicarea succesiva a doua pozitii ar pierde prima valoare fara o schimbare de model.
- `internal/datebaza/pontaj_calcul.go`: maparea intrarilor si a garzilor pe angajat/data pastreaza o singura valoare.
- `internal/datebaza/pontaj_planificare.go`: sincronizarea inversa presupune o singura pozitie activa pe zi.
- `frontend/src/Planificare.tsx`: grila si sumarul mapeaza o singura pozitie pe zi. Un dropdown suplimentar singur nu rezolva stocarea/publicarea/calculul.

## Activitati pentru implementare

1. Pastrarea intervalelor ca pozitii distincte, cu identificator, revizie, audit si editare/anulare individuala. Migrare noua compatibila cu datele existente.
2. Validare atomica a lotului final: suprapuneri pe aceeasi marca, inclusiv cu ziua precedenta/urmatoare si peste limita lunii. Editarea ignora numai pozitia proprie. Intervalele adiacente sunt permise.
3. Planificare: buton Adauga program pentru zi, lista intervalelor, editare/stergere individuala; grila, sumarul si totalurile afiseaza toate pozitiile. Prevenirea inlocuirii accidentale prin actiunile existente.
4. Publicare si sincronizare: intregul set al zilei, fara pierderea celorlalte intervale; idempotenta, conflicte de revizie si retragerea pontajelor efective conform prioritatii deja confirmate.
5. Pontaj si calcul: insumarea intervalelor distincte fara dublare, pauze proprii, impartire pe zile calendaristice si preluarea soldurilor. Reprezentarea pentru afisare nu trebuie folosita drept substitut al intervalelor la calcul.
6. Absentele si medicalul: aplicarea prioritatii existente asupra tuturor intervalelor afectate. Medicalul anuleaza intreaga tura de noapte intersectata, pe toate marcile asociate. Editarea/anularea absentei pastreaza auditul si regulile existente.
7. Rapoarte standard, flexibile/SQL si exporturi: toate intervalele disponibile, lista garzilor cate un rand per interval, totaluri corecte fara duplicate. Comparativul planificat/realizat si acoperirea verifica intregul set.
8. Teste anonime de migrare, adaugare/editare/stergere, suprapuneri/adiacenta, ture nocturne, limite de luna, publicare repetata/partiala, absente/medical, calcul si export. Probe UI RO/EN si cursor pastrat.
9. Versiune aplicatie/rapoarte, manual revizuit si build oficial cu manuale regenerate.

## Stare si predare

Cerintele si auditul initial sunt documentate. Codul aplicatiei si schema nu sunt modificate in acest commit; functionalitatea nu este inca disponibila in0.92.30. Pasul urmator este proiectarea reprezentarii intervalelor din Pontaj si o migrare noua, inaintea modificarii interfetei. Nu exista o decizie tehnica acceptata sau reguli de plata noi.
