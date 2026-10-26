# Lista de lucru UI după audit

Data: 2026-10-26 · Solicitant: Ioan Tuns · Bază: [audit 0.27.0](analiza-ui-2026-10-26.md), [D035](decizii/D035-prioritati-ui-dupa-audit.md).

Utilizatorul a cerut explicit lista și începerea îmbunătățirilor. Lucrăm în loturi verificabile, păstrând regulile de calcul confirmate.

| Ordine | Task | Conținut | Criteriu de acceptare | Stare |
| --- | --- | --- | --- | --- |
| 1 | F05-T19 | I045–I048: indicatori corecți, grilă compactă și contrast, context la navigare, salvare accesibilă în formulare lungi | T diferit de Ore norm nu produce avertisment contractual; la 1024×640 se văd antetul grilei și minimum 5 rânduri; luna/filtrele/selecția se păstrează în sesiune; antetul și salvarea rămân vizibile în dialogurile lungi; protecția modificărilor nesalvate rămâne activă | gata · Codex · 0.28.0 |
| 2 | F05-T20 | I049: motive distincte pentru valori necalculate și validare lângă câmp | Operatorul vede cauza și acțiunea posibilă; configurările invalide indică rândul/câmpul; fără inventarea regulilor contabile | de făcut |
| 3 | F05-T21 | I050: lista zilnică de absențe și gruparea cererilor | Cererile și zilele aferente se găsesc fără parcurgerea obligatorie a soldurilor anuale; selecția cererii rămâne disponibilă | de făcut |
| 4 | F05-T22 | I051: inserarea variabilelor în editorul de rapoarte | Variabila se inserează la cursor, cu validarea câmpului și diferențiere clară între structură și PDF | de făcut |
| 5 | F05-T23 | I052: simplificarea controalelor din Setări | Reducerea duplicării vizuale; păstrarea ambelor căi de navigare cerute prin D034; selector de coduri utilizabil cu tastatura | de făcut |

La finalul fiecărui lot se actualizează taskul, jurnalul și această listă. Verificările automate folosesc date fictive; o probă cu un operator și date anonimizate rămâne necesară pentru validarea fluxului zilnic. Persistența limbii este deja urmărită separat în F03-T05.

Primul lot a trecut criteriile automate: 7 rânduri complete la 1024×640, context restaurat și schiță abandonată neînregistrată, salvare accesibilă și contraste măsurate. [Rezultate și limite](rezultate-ui-F05-T19.md). Următorul lot este F05-T20.
