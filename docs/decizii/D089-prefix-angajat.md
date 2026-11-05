# D089 · Prefixul angajatului (Dr., Biol., As. …)

Stare: decisă (cerințele) · propusă (detaliile tehnice) · Data: 2026-11-05 · Taskuri: F07-T40

Decisă de: Ioan Tuns, 2026-11-05, pentru: listă predefinită din care se alege prefixul fiecărui angajat; prefix opțional; afișare în rapoarte (PDF/Excel) și în liste și selectoare, nu în grilele Pontaj și Planificare.

## Context

Ioan Tuns cere ca fiecare angajat să poată avea o titulatură înaintea numelui, de exemplu „Dr.” sau „Biol.” (I115). Angajatul nu are acum un astfel de câmp. Funcțiile au deja un nomenclator (D069/F07-T24), legat de angajat printr-o tabelă separată.

## Decizie

1. **Nomenclator nou `prefixe`**, gestionat în Setări ca celelalte nomenclatoare: denumire RO/EN (`i18n.Text`, D004), activ/inactiv, revizie și jurnal. Migrarea îl populează cu o listă de pornire, editabilă de utilizator: Dr., Biol., Chim., Farm., As., Ing., Ec., Psih. Forma engleză pornește identică cu cea română, iar utilizatorul o corectează unde e cazul. Aplicația nu inventează traduceri ale titulaturilor.
2. **Legătura cu angajatul** stă în tabela `angajati_prefixe(angajat_id, prefix_id)`, ca la funcții și contracte. Tabela `angajati` nu se schimbă. Prefixul este opțional; un prefix inactiv nu se poate alege nou, dar rămâne la angajații care îl au deja.
3. **Afișare:** „Dr. Popescu Ion” în rapoarte (PDF și Excel), în lista de angajați și în selectoare. Grilele Pontaj și Planificare rămân doar cu numele, cum a cerut utilizatorul. Marcajul de gardă se păstrează după nume.
4. **Rapoarte SQL:** în instantaneul `angajati` apar coloana nouă `prefix` (în limba raportului) și numele afișat cu prefix. `nume` și `prenume` rămân nemodificate. Filtrarea după prefix devine posibilă din SQL.
5. Numele salvat (nume, prenume) nu se modifică; prefixul este doar un atribut de afișare și de raportare.

## Variante respinse

- **Text liber în fișa angajatului:** aceeași titulatură ar apărea scrisă în mai multe feluri.
- **Coloană nouă în `angajati`:** schimbă interogările sqlc și citirea angajatului pentru o legătură opțională; tabela separată urmează modelul funcțiilor.
- **Prefixul în grilele lunare:** utilizatorul nu l-a cerut acolo, iar coloana numelui e deja îngustă.

## Consecințe

- Rapoartele care afișează numele angajatului își schimbă conținutul, deci primesc versiuni noi (D009).
- Asistentul AI și importurile nu primesc prefixul în acest task.
