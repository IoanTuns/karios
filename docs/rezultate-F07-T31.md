# F07-T31 - Zile lucratoare fara sarbatori legale

Responsabil: Codex. Ramura F07-T31-zile-lucratoare. Aplicatie 0.92.7, manual 1.18. [D067](decizii/D067-calendar-lucrator.md) propusa pentru mecanism; regula L-V fara sarbatori este cerinta utilizatorului.

## Modificare

AI individual/grup foloseste calendarul configurat pentru filtrul lucratoare; nu mai verifica numai L-V. Weekend ramane S-D, Toate zilele si datele explicite permit planificarea separata. Calendarul intra in amprenta previzualizarii AI. Textele UI si instructiunile AI explica excluderea sarbatorilor.

Testul de integrare a descoperit ca ListaSarbatori(ctx,nil) returna zero randuri, deoarece SQL-ul filtrat pe an compara cu NULL. Apelurile fara an au fost inlocuite cu ListaToateSarbatori in AI, calculul Pontaj si import. Query-urile generate si schema BD sunt neschimbate. Precompletarea existenta foloseste deja calendarul filtrat corect pe an; norma zilnica si motorul orelor exclud corect weekend/sarbatoare atunci cand primesc calendarul.

Rapoarte pontaj/detaliat/fisa v9; referintele vechi sunt pastrate. Nu se sterge munca existenta in zile nelucratoare. Clasificarea efectiva pe intervale 00-24 si fara dubla numarare se pastreaza. Nu se inventeaza sporuri, compensari sau reguli pentru cazuri neconfirmate; integrarea completa a T/garzilor/soldurilor ramane task separat.

## Verificari si predare

Teste noi: sarbatoare in cursul saptamanii, suprapunere cu sambata, clasificare separata duminica, selectie normala/weekend/toate, trecere intre ani si invalidarea tokenului dupa editarea calendarului. Integrare Pontaj: 7 ore vineri sarbatoare plus 7 ore sambata sarbatoare = 14 ZNL, zero ore normale si 14 ore efective, fara dublare. TypeScript/Vite, Edge focus/zoom, 18 teste frontend, 3 teste manual, 7 teste livrare si control documentatie trecute. gofmt, go vet ./... si suita Go completa au trecut. Testul tehnic pentru octombrie-noiembrie a fost actualizat conform cerintei: 42 zile/persoana, cu excluderea explicita a datei 30 noiembrie din calendarul demo. Fara modificarea referintelor contabile. Build oficial complet din 9c78447: `build/bin/pontix-0.92.7-20261031-143417-8cd4a59988a9457492a396cc92ce97b8/pontix.exe`. Manual 1.18 MD/HTML/PDF alaturi; manifest, hashuri executabil/manuale/arhiva si integritate ZIP verificate.

Calendarul trebuie verificat de operator pentru anii vizati. Datele reale nu au fost deschise. Proba interactiva desktop si imprimarea fizica raman la operator. Pas urmator: review si proba operator. Fara push/merge in main.

Dependente verificate: Go 1.27.3 si npm actuale; 62 actualizari tranzitive Go raman in F14-T21. Fara dependente noi.
