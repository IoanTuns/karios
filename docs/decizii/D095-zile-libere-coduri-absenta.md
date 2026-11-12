# D095 · Zilele libere întregi și codul absenței în locul lui 0

Stare: acceptată · Data: 2026-11-07 · Taskuri: F09-T30 · Decisă de: Ioan Tuns, 2026-11-12 („accept toate”)

Implementare: implementată din 0.92.73; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

Decisă de: Ioan Tuns, 2026-11-07 (cerere și răspunsuri explicite în conversație).

## Context

În rapoarte, o zi de absență arăta orele zilei („0”), colorată ca absență, iar codul (CO, CM…) nu apărea. Weekendurile și sărbătorile dintr-un interval de absență arătau la fel, iar în grila Pontaj apărea codul absenței și în zilele libere.

## Decizie

Pentru fiecare persoană și zi, numai unde persoana **nu are ore** în acea zi:

1. **Zi liberă întreagă** (sâmbătă, duminică sau sărbătoare legală, 00–24): celula e **goală** și marcată distinct, chiar dacă ziua intră într-un interval de absență. Exemplu: CO de vineri până luni → Vi `CO`, S gol, D gol, Lu `CO`.
2. **Zi lucrătoare cu absență**: apare **codul absenței** în locul lui 0.
3. Unde persoana are ore (gardă, program în weekend, FMP cu prezență), rămân orele. Valorile necalculate („—”) rămân neschimbate.

Se aplică în grila Pontaj și în toate rapoartele cu valori pe zile (Pontaj lunar, Pontaj detaliat, Fișa individuală, coloana ore din Evidența absențelor și rapoartele flexibile care le folosesc).

### Detalii tehnice (acceptate)

- Rapoarte: `celulaZiFaraOre` în `internal/rapoarte/document.go`; culoarea nouă `ZI_LIBERA` (`#D2D8D3` în PDF), distinctă de `WEEKEND` (#E7EAE7), care rămâne pe zilele libere cu ore.
- Grila Pontaj: clasa `zi-libera-intreaga`, hașură peste culoarea weekendului sau a sărbătorii, cu intrare în legendă și explicație în tooltip (cu codul absenței, dacă există).
- Totalurile nu se schimbă; absența rămâne înregistrată în zilele libere.
- Rapoartele primesc versiune nouă (D009).

## Consecințe

- În PDF, aproape toate celulele de weekend fără ore au acum culoarea `ZI_LIBERA`.
- Exportul Excel primește aceleași texte (goale sau codul), fără culori.
