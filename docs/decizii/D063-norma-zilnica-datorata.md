# D063 - Calculul de baza al normei zilnice datorate

Stare: propusa. Propusa de Codex, 2026-10-31. Task F07-T27.

Reguli confirmate de utilizator: reducere pentru CO/CEV/CFS/CM aprobate; formarea si delegatia nu sunt in lista reducerilor; liberele compensatorii nu au inca regula. CVE din ultimul raspuns este interpretat ca CEV din definitia anterioara.

Calcul pur pentru contracte cu norma zilnica si calendar luni-vineri: excludere date libere calendaristice si concedii aprobate, prin reuniune de date, nu prin adunarea numaratorilor. Tipul de concediu transmis motorului este clasificarea explicita de raport, nu o denumire libera. Norma nil ramane nedefinita, distincta de zero datorat pe un interval fara zile lucratoare. Nu se aplica unei norme lunare fixe prin impartire presupusa.

Functia nu este conectata inca la raport si nu stabileste calendarul turelor, valoarea T sau soldul initial. Pentru acestea raman necesare clarificarile deja transmise. Fara biblioteci sau migrari noi. Testele de referinta existente nu se modifica.
