# D070 - Corectarea explicita a absentelor

Stare: acceptată. Codex, 2026-10-31. F07-T34 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Utilizatorul confirma ca operatiunile, inclusiv concediile, trebuie sa poata fi editate/modificate sau sterse. Prioritatea medicalului din D069 nu trebuie interpretata ca imposibilitate de corectie explicita.

Propunere de implementare: o revizie curenta permite editarea medicalului in alt tip; o cerere fara revizie nu poate suprascrie medicalul, iar o revizie veche produce conflict. Anularea deja existenta pastreaza inregistrarea inactiva si jurnalul. Se mentin lunile inchise, validarile de calendar si protectia pontajului efectiv. Nu se restaureaza implicit programele anulate. Nu se introduce stergere fizica a istoricului.
