# D068 - Conversia schitei de Pontaj in absente

Stare: acceptată. Codex, 2026-10-31. F07-T32 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Utilizatorul solicita concedii si schimbarea tipului zilei din Completare Pontaj. Se reutilizeaza formularul si validarile existente de absente, pe selectia exacta de persoane/date. Noua operatie este explicita si salveaza imediat, in aceeasi tranzactie, golirea celulelor de schita si adaugarea/modificarea absentelor. Fluxul obisnuit SalveazaAbsente continua sa refuze suprapunerea cu o schita de lucru.

Se transmit instantaneele originale ale celulelor, cu valori si revizii; conflictul oricarei celule sau absente anuleaza intregul lot, inclusiv jurnalul. Numai absente active si nerespinse sunt eligibile. Se verifica luna deschisa; pontajul efectiv ramane protejat. Aprobarea, documentele, inlocuitorii si celelalte reguli de absente raman cele existente. Planificarea nu este stearsa, fiind evidenta intentiei, distincta de realizare.

Frontendul pastreaza modificarile nesalvate din afara selectiei si reseteaza istoricul local dupa salvarea conversiei, pentru a nu restaura programe peste absente salvate. Celula afiseaza codul absentei. Nu se schimba formulele rapoartelor si nu se adauga tipuri sau compensari presupuse. Fara migrare sau biblioteca noua.

Actualizare F07-T33: regulile ulterioare confirmate de utilizator si implementarea propusa in [D069](D069-prioritate-absente.md) inlocuiesc pastrarea programului planificat si refuzul suprascrierii din fluxul obisnuit.
