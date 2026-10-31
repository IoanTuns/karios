# F07-T25 - Plan de implementare

Responsabil Codex. Data 2026-10-31. Cerere utilizator: implementarea listei discutate.

1. Documentare cerinte/decizii si intrebari F01.
2. Grupuri profesionale cu ordine configurabila.
3. Functii asociate grupurilor si pozitii conducere/executie.
4. Atribuire prin functia angajatului, in sectia/departamentul existent.
5. Tipuri de contract configurabile, normare zilnica/lunara/fara norma.
6. Parametri contract: minute implicite, categorie de programe permise.
7. Contract de garda si marca legata explicit de principal.
8. Filtre grup, functie, pozitie, contract, garda in liste/grile/rapoarte.
9. Rapoarte grupate, sefii primii, apoi alfabetic; marci normal/garda impreuna.
10. Reguli conditionate de grup/functie/contract/program/zi.
11. Medic garda weekend/sarbatoare: plata TOTALA 200%, nu spor suplimentar 200%; exprimare in ore echivalente, fara valoare lei pana la definirea bazei.
12. Asistenti: ZNL evidentiaza orele eligibile; conversia in zile libere si evidenta acordarii sunt blocate pana la formula confirmata.
13. Migrari fara clasificare automata dupa nume si fara pierderea datelor.
14. Teste, manual, build versionat.

Ordinea din rapoarte nu determina drepturi salariale. Configurarea este curenta; versionarea pe intervale si istoricul regulilor nu vor fi pretinse implementate daca lipsesc.

## Stare la predare

Punctele 1-11 si 13 sunt implementate pentru configurare, filtrare, grupare si calcul in ore/ore echivalente. La punctul 12 exista sursa ZNL si reguli conditionate, dar conversia si acordarea efectiva a timpului liber raman in F07-T26 (blocat, formule lipsa). Punctul 14: verificari automate si manual actualizat; build oficial dupa commit. Calculul in lei este tot in F07-T26. Instructiunile detaliate sunt in manualul 0.92.0 si in rezultate-F07-T25.md.
