# F16-T27 - Filtre pentru starea activa

Responsabil: Codex. Ramura F16-T27-filtre-stare. Aplicatie 0.92.14, manual 1.25. D074 propusa.

Componenta comuna Active/Inactive/Toate, implicit Active, folosita in nomenclatoare, Schimburi, Categorii program, zile si tipuri de absenta, conturi, modele de raport si reguli conditionate. Angajati pastreaza controlul existent, implicit activi. Pontaj filtreaza angajatii; Planificare filtreaza separat angajatii si pozitiile active/retrase in lista. Filtrele sunt locale per lista, fara stergerea istoricului. Selectiile pentru actiuni in lot se reseteaza la schimbarea starii; modificarile nesalvate se pastreaza.

Ecranele fara camp activ (jurnal, backupuri, sarbatori) nu primesc stari artificiale. Generarea rapoartelor si formulele istorice nu se schimba; modelele inactive se gestioneaza in Definire rapoarte. Referintele inactive deja asociate raman vizibile in formularele de editare. Adaugarea unei reguli initial inactive trece explicit la Toate pentru ca editorul nou sa fie vizibil.

22 teste frontend si TypeScript/Vite trecute. Edge RO/EN nomenclatoare: implicit Active, dezactivare, disparitie din lista activa, regasire la Inactive si Toate. Edge RO/EN absente: filtre zile/tipuri si reactivare. Scenariile anterioare au fost adaptate pentru afisarea explicita a istoricului. 3 teste manual, 7 livrare si control documentatie trecute; gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute; Edge focus/zoom trecut. Grila Pontaj verificata RO/EN pentru implicit Active, Inactive si revenire, inclusiv resetarea selectiei. Integrarea frontend final reverificata cu go test .; indiciul vechi despre inactivi actualizat.

Un script de aplicare s-a oprit la un marker diferit in Rapoarte, dupa primele module; continuarea a fost aplicata si compilata cu succes. O extindere a scenariului UI s-a oprit la numele capturii, corectat ulterior; scenariul final a trecut in ambele limbi. Fara modificari de calcul, migrari, dependente noi sau referinte contabile. Go stabil 1.27.3 si npm actuale; actualizarile tranzitive Go sunt urmarite separat prin F14-T21.

Pas urmator: build oficial, apoi proba operator WebView2. Datele reale nu au fost accesate; fara push/merge in main.
