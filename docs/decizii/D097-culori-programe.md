# D097 · Culoare pe fiecare program (tură, gardă)

Stare: acceptată (cerința) · propusă (detaliile tehnice) · Data: 2026-11-08 · Taskuri: F07-T42

Implementare: implementată din 0.92.79; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

Decisă de: Ioan Tuns, 2026-11-08 (cerere și răspunsuri explicite în conversație: culoare pe fiecare program; în grilele Planificare și Pontaj și în rapoartele PDF).

## Context

Programele (ture, gărzi, din Setări › Programe) se deosebeau în grile și rapoarte doar prin cod. Ioan Tuns a cerut culori pentru fiecare tip de tură și gardă.

## Decizie

- Fiecare program are o culoare opțională (`#RRGGBB`), aleasă în Setări › Programe. Fără alegere se folosește **culoarea implicită după tipul programului**: Normal `#DDE9F8` (albastru), Tură `#FCE6CC` (piersică), Noapte `#E2DCF5` (lavandă), Gardă `#F8D7DC` (roz). Nuanțele sunt deschise, lizibile cu text închis și diferite de weekend (gri), sărbătoare (galben), absență (verde) și zi liberă întreagă (hașură).
- **Planificare:** fundalul celulei planificate are culoarea programului (primul, dacă sunt mai multe); butoanele programelor speciale au mostra culorii; o legendă „Culori programe” arată programele lunii.
- **Pontaj:** celula cu cod de program are culoarea lui; absențele și zilele libere întregi își păstrează marcajul (D095).
- **Rapoarte PDF:** celulele zilelor cu program (Pontaj lunar, Pontaj detaliat, Fișa individuală, Grafic gărzi) au fundalul programului; absențele și zilele libere rămân cu culorile lor.
- Selecția, avertismentele și modificările nesalvate rămân deasupra culorii programului. În tema întunecată culoarea devine o nuanță peste fundal (28%).

### Detalii tehnice (propuse)

- Migrarea `0036_culori_programe.sql`: coloana `tipuri_schimb.culoare` (NULL = implicită), validată `#RRGGBB` cu majuscule.
- `datebaza.CuloareProgram(aleasa, tip)` dă culoarea afișată; `Schimb.culoare_afisata` o trimite interfeței, ca regula să fie într-un singur loc.
- În PDF, `Rand.Culori` acceptă și `#RRGGBB` (`fundalCelulaPdfme`).
- Rapoartele cu celule de program primesc versiune nouă (D009).

## Variante respinse

- **O culoare pe tip de program** — mai puțin flexibil; rămâne ca implicită.
- **Fără culoare implicită** — utilizatorul ar trebui să aleagă pentru fiecare program ca să vadă vreun efect.
- **Culoarea textului în loc de fundal** — în Planificare, culoarea textului arată deja stadiul (nepublicat, publicat).

## Consecințe

- Rapoartele SQL personalizate nu primesc culori automat (rândurile vin din SQL); poate fi o idee separată.
- Exportul Excel nu are culori.
