# D079: simulare de absente si coordonarea fluxului lunar

Stare: propusa. Data: 2026-10-31. Autor: Codex. Task: F16-T31.

Cerinta utilizatorului: implementarea tuturor celor sase recomandari UX din F16-T30.

Propunere tehnica: simularea absentelor refoloseste operatia existenta cu un colector limitat la sursele si campurile necesare comparatiei. Tranzactia se anuleaza inclusiv pentru documente, revizii si audit; salvarea finala ramane operatia normala, cu validarile existente. Nu se copiaza baza de date si nu se transmit documentele sau alte date personale suplimentare in rezultatul simularii. Fara biblioteca noua.

Alternativa: reconstructie separata in frontend; respinsa ca propunere deoarece ar dubla regulile de medical, ture si asocieri. Previziunea nu este o rezervare a datelor; confirmarea revalideaza starea curenta. Panoul lunar este consultativ si foloseste API-urile existente, fara reguli contabile noi. Dialogul comun de navigare acopera grilele si sumarul Planificarii, iar inchiderea nativa a aplicatiei isi pastreaza protectia existenta.

Implementarea este autorizata de cererea utilizatorului. Decizia ramane propusa pana la evaluarea explicita de catre un om.
