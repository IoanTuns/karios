# F07-T37 - Stergere configurari de administrator

Responsabil: Codex. Aplicatie 0.92.15, manual 1.26. Ramura F07-T37-stergere-configurari-admin. Domeniul confirmat explicit de utilizator; D075 ramane propusa pentru reprezentarea tehnica.

Buton Sterge definitiv pentru administratori in nomenclatoare, Schimburi, Categorii program, tipuri de absenta si Definire rapoarte. Se selecteaza Inactive; o definitie activa trebuie dezactivata intai. Confirmarea afiseaza numele. Endpointul verifica sesiunea si rolul/starea actuale din baza. Nu exista exceptie de acces pe baze fara conturi.

Reviziile, referintele FK, codurile din schite si referintele JSON din setari/sabloane/istoric AI blocheaza stergerea. Definitii structurale sunt protejate. Metadatele proprii se elimina numai cu definitia; programele unei categorii raman. Snapshotul complet se pastreaza in jurnal. Tranzactia asigura rollback inclusiv la esec dupa eliminarea listei proprii.

## Verificari

- gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute; teste suplimentare TestStergere trecute pentru absente, modele si referinte AI.
- Endpoint: fara sesiune/operator refuzati; administrator autorizat; administrator dezactivat sau retrogradat dupa autentificare refuzat.
- Serviciu: activa, revizie veche, referinte FK/JSON/schite, definitii protejate, tip invalid, audit si rollback.
- 22 teste frontend, TypeScript/Vite, 3 teste manual, 7 teste livrare, control documentatie trecute.
- Edge RO/EN: confirmare anulata fara apel backend, refuz de utilizare tradus, succes si disparitie din lista, operator fara buton si refuz direct. Regresii nomenclatoare si absente din Pontaj RO/EN trecute.
- Prima proba UI a necesitat completarea simularii cu ListaBackupuri; nu a fost defect de productie. Prima implementare propusa fusese oprita de verificarea automata; domeniul a fost apoi confirmat explicit, inaintea implementarii.
- Dependente verificate: Go stabil 1.27.3, npm fara actualizari; actualizarile tranzitive Go identificate raman in taskul separat F14-T21, pentru compatibilitate si verificari dedicate. Nu s-au instalat biblioteci.

Nicio baza reala nu a fost accesata sau modificata. Nu se sterg angajati, conturi, pontaje, planificari, zile de concediu ori jurnal. Urmeaza build oficial cu manual versionat si verificarea integritatii pachetului. Fara push sau merge.
