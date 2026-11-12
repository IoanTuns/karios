# D108 · Exportul pentru salarizare, configurabil

Stare: propusă · Data: 2026-11-12 · Taskuri: F10-T01, F10-T03 · Cerută de: Ioan Tuns, 2026-11-12 (chestionarul F01-T02: E1 „nu este nimic stabilit; poți face formate pentru câteva programe uzuale?”, E2 „Excel sau CSV; poate fi și alt format”; apoi „da” la exportul configurabil)

## Context

Pontix predă orele și zilele programului de salarizare; nu calculează bani (D100). Programul de salarizare nu e ales, iar formatele exacte de import ale programelor uzuale nu sunt în proiect. A le ghici ar da fișiere care par corecte, dar pe care programul le respinge sau, mai rău, le citește greșit.

## Decizie

**Un export configurabil, cu modele salvate.** Fiecare model descrie un fișier:

- **Formatul:** CSV sau Excel (`.xlsx`).
- **Pentru CSV:**
  - separatorul de câmpuri: `;` (implicit, potrivit pentru Excel în română), `,` sau tab;
  - semnul BOM la început, pentru diacritice în Excel; implicit da.
- **Numerele:**
  - separatorul zecimal (`,` sau `.`);
  - numărul de zecimale (0–4);
  - forma orelor: zecimal (`7,5`) sau ore:minute (`7:30`).
- **Antetul:** primul rând cu denumirile coloanelor; se poate omite.
- **Coloanele,** în ordinea dorită, fiecare cu titlul ei. Ce poate conține o coloană:
  - identitatea: marca, numele, prenumele, numele complet, secția, departamentul, funcția;
  - perioada: anul, luna, perioada `AAAA-LL`;
  - orice coloană de total din Pontaj (OL, OS, ZNL, CO, CM, CFP, CVE, DLG, FMP, T, Ore normă, Ore gardă 200%, TOCN, TZLN…), în unitatea ei: orele ca ore, zilele ca zile;
  - un text fix, de exemplu un cod cerut de program.

**Angajații.** Se identifică prin **marcă**; Pontix nu păstrează CNP. Un rând pe marcă, pentru angajații activi sau care au date în lună, ordonat după marcă. Exportul se poate restrânge la o secție.

**Valorile necalculate („—”)** rămân celule goale. După export, aplicația spune câte au fost, ca să nu treacă neobservate.

**Modelele** stau în baza criptată (tabelul `modele_export_salarizare`), cu revizie și jurnal.
- La instalare există două modele generice: „CSV generic” și „Excel generic”. Ele conțin marca, numele, OL, OS, ZNL, CO, CM.
- Modelele se editează în **Setări › General › Export salarizare**.
- Exportul se face din **Pontaj**: **Mai multe** › **Export pentru salarizare**.

**Modelele pentru programe anume** (de exemplu un program de salarizare ales de unitate) se adaugă după documentația lor de import, ca modele predefinite cu numele programului. Nu se inventează structuri.

**Versiunea formatului:** `VersiuneExportSalarizare` = 1.0.0, cu teste pe fișier de referință (D009) pentru CSV și Excel.

## Variante respinse

- **Formate „uzuale” scrise din memorie** pentru programe de salarizare. Fără documentația de import, riscul de date citite greșit e prea mare.
- **XML sau alte formate acum.** Nu e cerut de un program anume; se adaugă când apare documentația.
- **Identificare prin CNP.** Pontix nu păstrează CNP (date personale minime, D008).

## Consecințe

- Migrarea 0042 adaugă tabelul modelelor, cu cele două modele generice.
- F10-T01 rămâne deschis pentru alegerea programului de salarizare și pentru documentația lui de import.
