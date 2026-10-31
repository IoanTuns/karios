# F07-T34 - Corectarea explicita a concediilor

Responsabil: Codex. Ramura F07-T34-corectare-concedii. Aplicatie 0.92.10, manual 1.21. D070 propusa.

Cerinta: operatiunile, inclusiv concediile, raman corectabile si anulabile. Verificarea modulului Absente a confirmat editarea individuala si in lot, anularea individuala si in lot, documentele editabile si jurnalul existent. Restrictia din F07-T33 bloca insa schimbarea explicita a medicalului in alt tip; a fost corectata.

O editare cu revizia curenta poate schimba medicalul in CO sau alt tip valid. Fara revizie, prioritatea medicalului se pastreaza. O revizie veche nu suprascrie inregistrarea curenta. Anularea foloseste inactivare, cu istoric. Lunile inchise si calendarul raman protejate; programele anulate de medical nu se restaureaza automat. Interfata explica distinctia.

Teste backend: CM-CO-CM, anulare medical, conflict de revizie; existentele teste confirma prioritatea fara revizie. Edge cu API simulat RO/EN: editare medical in CO, revenire la CM, anulare cu confirmare, plus conversie program, eroare/anulare formular si pastrarea modificarilor din afara selectiei. Backendul real este verificat separat prin testele Go.

TypeScript/Vite, 18 teste frontend, 3 teste manual, 7 teste livrare si control documentatie trecute. gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute. Go 1.27.3 si npm actuale; cele 62 actualizari tranzitive Go sunt urmarite separat in F14-T21. Fara migrari, dependente noi sau schimbarea formulelor/cazurilor contabile de referinta.

Pas urmator: build oficial; proba operator pentru aplicatia desktop. Verificarea curenta acopera concediile si controalele existente de corectare, nu certifica toate fluxurile aplicatiei. Fara acces la date reale, fara push/merge.
