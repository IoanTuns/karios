# D093 · Gărzile: fără ZNL, ore gardă 200% și ore lucrătoare

Stare: acceptată (regulile) · propusă (detaliile tehnice) · Data: 2026-11-07 · Taskuri: F08-T14

Decisă de: Ioan Tuns, 2026-11-07 (cerere și răspunsuri explicite în conversație).

## Context

Rândurile mărcii de gardă (D060: contractul sau marca de gardă a unei persoane, afișată ca rând separat) nu au normă zilnică. Până acum, pe aceste rânduri ZNL număra orele din weekend și sărbători ca la orice rând, „Ore gardă 200%” era neconfigurată, iar „Ore lucrătoare” rămânea goală (norma lipsește).

## Decizie

Numai pe **rândurile mărcii de gardă**:

1. **ZNL nu se calculează**: valoarea e „—”, cu motivul afișat. Pe celelalte rânduri ZNL rămâne neschimbat.
2. **Ore gardă 200%** = orele de gardă care cad sâmbăta, duminica sau într-o sărbătoare legală (aceeași împărțire pe zile calendaristice ca la ZNL, D028).
3. **Ore lucrătoare** = total ore − ore gardă 200%.

### Detalii tehnice (propuse)

- Calculul (`internal/calcul`) primește mulțimea rândurilor de gardă (`ConfiguratiePontaj.Garzi`, din profilul pontajului) și aplică regulile peste valorile configurate în Configurare pontaj pentru coloanele ZNL și Ore gardă 200%.
- Soldurile lunare nu suprascriu „Ore lucrătoare” pe rândurile de gardă; L început și L/FIN rămân neschimbate (fără normă, nedefinite), pentru că nu s-a cerut altceva.
- Rapoartele care afișează aceste coloane primesc versiune nouă (D009).

## Consecințe

- Orele de gardă din weekend și sărbători apar o singură dată, în „Ore gardă 200%”, nu și în ZNL.
- Dacă o marcă de gardă are și programe care nu sunt gărzi, ele sunt tratate la fel (regula e pe rând, nu pe tipul programului).
