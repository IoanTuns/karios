# F07-T36 - Data si starea absentelor

Responsabil: Codex. Ramura F07-T36-editare-data-stare-absente. Aplicatie 0.92.13, manual 1.24. D073 propusa.

Confirmat: campurile de data erau readOnly, zilele anulate erau ascunse fara reactivare, iar actiunile puteau fi in afara zonei vizibile a tabelului lat. Formularul permite acum schimbarea datei unei zile, cu sincronizarea limitelor. Vechea inregistrare devine inactiva si noua este salvata in aceeasi tranzactie folosind endpointul existent. Destinatiile ocupate sunt refuzate in formular; erorile serviciului anuleaza integral mutarea.

Sterge ziua anuleaza cu pastrarea istoricului. Filtrul Arata si zilele anulate permite Reactiveaza. Starea Activ se poate modifica in editorul unei zile/lot existente. Tipurile de absenta au butoane directe Dezactiveaza/Reactiveaza. Actiunile din tabel raman vizibile la derulare. Totalurile existente continua sa foloseasca doar zile active. Nu se introduc migrari, calcule sau reguli noi de plata; clasificarea separata CFS ramane distincta de aceasta corectie a editarii.

Teste backend real: mutare CO, anulare/reactivare, rollback la destinatie nelucratoare, conflict sursa/destinatie sau luna inchisa. Edge cu API simulat RO/EN: schimbare data, sincronizare limite, salvare, stergere, afisare anulate/reactivare si dezactivare/reactivare tip; regresiile conversiei programului si modificarii CM/CO raman acoperite. TypeScript/Vite, 21 teste frontend, 3 manual, 7 livrare, control documentatie si Edge focus/zoom trecute. gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute integral.

Dependente: Go stabil 1.27.3 si npm actuale; 62 actualizari tranzitive Go urmarite separat in F14-T21. Date reale neaccesate, fara modificarea referintelor contabile. Pas urmator: build oficial, apoi proba operator WebView2. Fara push/merge.
