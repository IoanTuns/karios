# F09-T11 - Marci distincte pentru program normal si garda

Responsabil: Codex. Data: 2026-10-31. Ramura: F09-T11-marci-garda. Aplicatie 0.91.0, raport detaliat v4, migrare 0025, manual 1.9. Include F07-T24 (77c2312) si corecturile F16-T23.

## Comportament

Angajati: inregistrarea cu marca de garda se leaga explicit de inregistrarea principala prin selectorul Angajat principal. Lista arata subtipul Garda si marca principala. Datele existente raman neasociate; omonimii nu se unesc. Se refuza legaturile catre sine, lanturile, ciclurile si referintele inexistente, inclusiv la nivelul bazei. Legatura este auditata si salvata tranzactional.

Pontaj detaliat: randuri distincte, grupate principal apoi garzi, fiecare cu marca sa. Filtrul de angajat include familia; filtrele de structura raman aplicate fiecarei inregistrari. Pentru subtipurile Garda se foloseste planificarea GARZI activa ca rezerva in zilele fara inregistrare Pontaj. Pontajul existent, inclusiv zero sau valoare goala, si absentele au prioritate. Publicarea nu dubleaza orele. Sunt folosite motorul si regulile deja configurate; nu sunt adaugate procente de salarizare. Nota raportului explica sursa. Turele din ultimele doua zile ale lunii precedente sunt incluse conform atribuirii configurate.

## Verificari

- gofmt -l -w ., go vet ./..., go test -p 1 -timeout 120s ./...: trecute.
- Teste noi: migrare 24->25 fara potrivire dupa nume; persistenta si audit; legaturi invalide protejate; doua marci distincte; 7h normal/24h garda; fara dublare dupa publicare; zero manual si absente prioritare; filtre, omonim si inactiv cu date; tura impartita 17h/7h si 7h venite din luna precedenta.
- Referintele vechi pastrate, referinte noi v4 pentru raportul de baza si proba cu marci de garda. PDF fictiv validat cu pdfcpu si randat cu Windows; inspectie vizuala trecuta.
- TypeScript/Vite si 18 teste frontend: trecute. UI Edge cu punte simulata: nomenclatoare si selectie marca principala RO/EN, salvare si afisare subtip; focus/cursor si zoom trecute. Backendul real verificat separat in Go.
- 3 teste manuale, 7 teste livrare si verificarea documentatiei: trecute. Ajutor si traduceri reverificate dupa actualizare.
- Go 1.27.3, dependente npm actuale. Actualizarile tranzitive Go observate raman in taskul existent F14-T21; fara biblioteci noi in aceasta implementare.

## Predare

Build oficial urmeaza dupa commit, pentru identificator curat in manifest. Nu s-au deschis datele reale. Ramane proba operator in aplicatia desktop/WebView2 si tiparirea fizica; legaturile dintre marcile existente trebuie selectate de operator. Nicio publicare sau integrare in main efectuata.
