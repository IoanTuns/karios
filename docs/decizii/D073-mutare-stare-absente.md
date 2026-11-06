# D073 - Mutarea si starea absentelor

Stare: propusa. Codex, 2026-10-31. F07-T36.

Cerinta utilizatorului: schimbarea datei si eliminarea/activarea/dezactivarea absentelor. Data era readOnly; anularea exista dar ascundea inregistrarea, fara reactivare accesibila.

Se reutilizeaza tranzactia de absente: mutarea anuleaza vechea zi si inregistreaza noua zi cu revizii, in acelasi lot. Destinatiile ocupate sunt protejate; orice eroare revine integral la starea initiala. Filtru pentru zile anulate si reactivare; starea este editabila in formular. Stergerea din evidenta activa este anulare jurnalizata, nu stergere fizica. Tipurile au actiuni directe activare/dezactivare. Reguli de calendar, luni inchise si conflicte raman valide. Fara restaurarea implicita a programului anterior.
