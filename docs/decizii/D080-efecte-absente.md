# D080 - Efectele tipurilor de absenta

Stare: acceptată. Data: 2026-11-02. Task: F07-T38 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Cerinta si reguli confirmate de utilizator: absente fara ore lucrate cu reducerea obligatiei lunare; absente cu ore lucrate (deplasari) din norma zilnica sau programul planificat. Alegerea sursei se face pe tip.

Propunere tehnica: efect separat de categoria de prioritate; valori FARA_ORE, NORMA, PROGRAM. ANTERIOR pastreaza clasificarea si regulile existente fara recalculare automata a istoricului. Un tip deja folosit isi pastreaza efectul; pentru schimbare se creeaza alt tip. Numai aprobatele participa la totaluri.

Durata creditata se calculeaza in backend si se pastreaza pe zi inaintea anularii programului, in aceeasi tranzactie si audit. Editarea aprobarii si reactivarea pastreaza durata. Norma sau programul lipsa produc eroare, nu zero presupus. Pentru corectarea duratei deja capturate se foloseste un tip distinct, cu sursa corecta; anularea/reactivarea aceluiasi cod nu rescrie instantaneul. Migratia 0029 adauga campuri fara a rescrie date existente.

Alternative: calcul din planul curent respins deoarece absenta anuleaza planul; deducerea efectului din nume respinsa deoarece denumirile sunt libere; schimbarea retroactiva a tipurilor folosite respinsa deoarece afecteaza lunile istorice.

Formularul existent, previzualizarea cu rollback, aprobarile, reviziile si prioritatea medicalului sunt reutilizate din Planificare. Nu sunt reguli monetare noi. Plata/remunerarea nu este echivalenta cu orele creditate.
