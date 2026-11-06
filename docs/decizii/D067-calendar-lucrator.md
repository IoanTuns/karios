# D067 - Calendarul programului normal

Stare: propusa pentru mecanismul de implementare. Codex, 2026-10-31. F07-T31.

Regula functionala este precizata explicit de utilizator: program normal L-V = zile lucratoare fara sarbatori legale. Nu se presupun reguli noi de plata. Weekendul ramane filtru S-D, distinct de sarbatorile care pot cadea in timpul saptamanii; o suprapunere nu dubleaza orele eligibile.

Se reutilizeaza calcul.Calendar.TipZi, unde sarbatoarea are prioritate fata de weekend. Filtrele AI consultau anterior numai ziua saptamanii; se conecteaza la calendarul configurat. Calendarul intra in amprenta planului AI, pentru a impune regenerarea dupa modificari.

Auditul a identificat si folosirea ListaSarbatori(ctx,nil): SQL-ul generat compara an=NULL si nu intoarce nimic. Apelurile fara an sunt inlocuite cu o interogare explicita ListaToateSarbatori, fara modificarea fisierelor generate, migrarilor sau interogarii filtrate pe an. Se corecteaza astfel citirea pentru AI, calcul Pontaj si verificarea existentei sarbatorilor la import. Fara biblioteca noua.
