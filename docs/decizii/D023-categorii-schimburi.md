# D023 · Categorii pentru nomenclatorul de schimburi (`LAB` și `CUSTOM`)

Stare: propusă · Data: 2026-10-25 · Task: F07-T11

Ioan Tuns a cerut împărțirea tabelului de schimburi pe categorii: codurile `G24`, `IN9`, `N12` (confirmat explicit ca forma corectă pentru `n2`), `N9`, `Z7`, `Z12`, `S7` intră în categoria de precompletare `LAB`, iar restul schimburilor intră în categoria `CUSTOM` pentru moment. La clarificare a confirmat că `categorie` este un câmp separat de `tip` și `regim` (cu valorile predefinite `LAB` și `CUSTOM`), cu filtru în tabel și alegere la precompletare, adăugare și editare.

## Reprezentare propusă

- Migrarea `0011_categorii_schimburi.sql` adaugă coloana `categorie TEXT NOT NULL DEFAULT 'CUSTOM' CHECK (categorie IN ('LAB', 'CUSTOM'))` în `tipuri_schimb`.
- Înregistrările neșterse cu codurile din setul `LAB` (`Z7`, `S7`, `N9`, `IN9`, `Z12`, `N12`, `G24`) primesc `categorie = 'LAB'`, fără a li se modifica intervalele, denumirile, starea `activ`, `tip` sau `regim`. Toate celelalte schimburi existente (`7`, `12`, `S1`, `S3` și orice schimb definit de utilizator) primesc `categorie = 'CUSTOM'`.
- Fișierul `internal/datebaza/initiale/schimburi.json` include `"categorie": "LAB"` pentru cele 7 șabloane.
- Precompletarea explicită (`CompleteazaSchimburiInitiale`) acceptă categoria dorită (implicit `LAB`), adăugând numai codurile lipsă din acea categorie, fără suprascrierea schimburilor existente sau coliziuni cu absențele.
- În interfață (`Schimburi.tsx`), tabelul afișează coloana `Categorie`, permite filtrarea după categorie (`LAB` / `CUSTOM`), alegerea categoriei de precompletare (`LAB`) și selectarea categoriei în formularul de adăugare/editare (implicit `CUSTOM` la adăugarea manuală a unui schimb nou). Sugestiile din pontaj includ și categoria.

Categoria organizează nomenclatorul și precompletarea șabloanelor; nu introduce reguli sau procente de plată noi. Reprezentarea tehnică rămâne propusă; încadrarea celor 7 coduri în `LAB` și a restului în `CUSTOM` este cerință explicită a lui Ioan Tuns.
