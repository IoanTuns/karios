# D077 - Solduri lunare si initializare

Stare: acceptată. Data: 2026-10-31. Task F07-T27 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Reguli confirmate explicit de utilizator: norma lunara scade concediile CO/CEV/CFS/CM; L Fin = L Inceput + ore lucrate - ore lucratoare. Nu se modifica procente sau reguli de plata.

Propunere tehnica: initializare pe marca si luna intr-o tabela separata (migrare 0028), minute semnate, zile CO consumate anterior in acelasi an, revizie si audit. Soldurile se calculeaza in minute, cronologic, in adaptorul comun grilei si rapoartelor. Fara initializare, prima luna cu date porneste de la zero. Orele neinterpretabile fac soldul final nedefinit. Garda nu are norma presupusa.

CEV/CFS folosesc codurile explicite ale rubricilor de absente configurate; categoria Special singura nu identifica o reducere de norma. Se pastreaza contorul CO initial separat de absentele calendaristice, cu refuzul dublarii si al recalcularii lunilor inchise prin schimbarea initializarii. Soldul initial nu este o absenta inventata si nu creeaza date calendaristice.

Raport detaliat v10: Ore lucratoare in locul vechii rubrici afisate Ore lucrate; L Inceput dupa Marca si legenda explicita. Calculul intern al orelor efective ramane separat. Referintele v9 sunt pastrate, iar cele tehnice v10 sunt generate separat.

Limita: se foloseste norma contractuala curenta, deoarece istoricul normelor nu este definit in schema actuala. Nu se inventeaza perioade sau norme istorice.
