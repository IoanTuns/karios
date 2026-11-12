# D098 · Tipurile D092 create după salvarea configurației și totalurile raportului detaliat

Stare: propusă · Data: 2026-11-09 · Taskuri: F08-T15

Implementare: implementată din 0.92.80; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

## Context

Ioan Tuns a raportat (2026-11-09) două probleme după ce a definit formarea profesională:

1. Zilele **FMP** apăreau în coloana **CVE** din Pontaj. Cauza: dacă Configurare pontaj fusese salvată înainte să existe tipul FMP, lista „tipuri excluse” a CVE se salvase fără FMP, iar coloana FMP rămăsese neconfigurată. Tot așa, orele FMP lipseau din OL.
2. În **Pontaj detaliat** nu apărea sumarul: raportul avea încă coloanele de dinainte de D092 (CEV, CFS, Form. prof., Deleg.), fără reguli, deci mereu „—”.

## Decizie

- La citirea configurației, regulile D092 se completează cu tipurile apărute după salvare:
  - coloana proprie (CO, CM, CFP, DLG, FMP) neconfigurată și fără coduri primește codul ei, dacă tipul există;
  - **CVE** exclude mereu tipurile cu coloană proprie (CO, CM, CFP, DLG, FMP și codurile oricărei alte coloane pe zile de absență);
  - **OL** include creditele FMP.
- În rapoartele standard, o coloană veche **neconfigurată** arată coloana D092 care a înlocuit-o, cu titlul și explicația ei: CEV→CVE, CFS→CFP, Form. prof.→FMP, Deleg.→DLG. Codul coloanei (`total_CEV` etc.) rămâne, ca SQL-urile salvate să funcționeze. O coloană veche configurată de utilizator rămâne cum e.

## Variante respinse

- **Corectarea manuală a configurației** de fiecare dată când apare un tip nou — ușor de uitat; exact asta a produs problema.
- **Coloane noi în raport** (`total_FMP` etc.) în locul celor vechi — ar strica modelele SQL salvate.

## Consecințe

- Un utilizator nu mai poate scoate FMP din OL sau din excluderile CVE prin configurație (regula D092 e definitorie).
- Raportul Pontaj detaliat și rapoartele flexibile primesc versiune nouă (D009).
