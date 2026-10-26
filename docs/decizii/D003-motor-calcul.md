# D003 · Principiile motorului de calcul

Stare: acceptată · Data: 2026-10-24 · Taskuri: F04-T01, F04-T02, F04-T03

## Context

Orele lucrate trebuie împărțite pe categorii de tarif (zi, noapte, weekend, sărbătoare) corect și în cazurile grele: ture peste miezul nopții, gărzi de 24 h, sărbători, schimbarea orei.

## Decizie

- Orele lucrate se salvează ca interval (start, sfârșit). Totalurile se calculează, nu se introduc.
- `Segmenteaza` taie fiecare tură la miezul nopții și la limitele intervalelor din reguli.
- **Tipul zilei** se ia din ziua calendaristică a fiecărui segment. Sărbătoarea are prioritate față de sâmbătă și duminică.
- **Minutele sunt timp real** în fusul `Europe/Bucharest`. O tură de 12 h pe ceas are 13 h în noaptea în care ora se dă înapoi și 11 h în cea în care se dă înainte. *Dacă se plătesc orele reale sau cele din grafic e întrebare deschisă (F01-T02).*
- **Regulile** acoperă exact 00:00–24:00 pe fiecare tip de zi, fără goluri și fără suprapuneri. Setul se verifică la salvare (`Valideaza`).
- **Provizoriu:** fiecare minut intră într-o singură categorie, ca în documentul de design. Modelul se revizuiește după răspunsurile contabilității despre cumularea sporurilor (F04-T03).

## Variante respinse

- **Câmp `prioritate` pe intervale** — eliminat deja în documentul de design; intervalele exclusive sunt mai ușor de verificat.
- **Calcul în ore de ceas** — ar da 12 h în nopțile cu schimbarea orei, ceea ce nu e timpul lucrat efectiv.

## Consecințe

- Motorul nu știe nimic de bază de date sau interfață; se testează separat (`internal/calcul`).
- Sporurile cumulate și orele suplimentare nu se pot exprima încă. Vor cere probabil o decizie nouă care o înlocuiește parțial pe aceasta.
