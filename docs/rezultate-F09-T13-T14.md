# F09-T13 / F09-T14 - Rapoarte configurabile

Data: 2026-11-02. Responsabil: Codex. Ramura: F09-T13-editor-flexibil. Aplicatie 0.92.25, raport flexibil v1, manual 1.35. Propunerea tehnica D082 ramane propusa; formulele de domeniu sunt confirmarile explicite ale utilizatorului din F01.

## Implementare

Setari > Definire rapoarte permite trei surse noi: SQL liber pe date de raportare, calcule lunare si lista de garzi. Schema este consultabila din editor. SQL parametrizat ruleaza pe un instantaneu izolat, filtrat inaintea interogarii, cu o singura instructiune de citire si limite de timp/rezultat. Nu citeste conturi, chei, documente sau jurnalul operational.

Coloane noi per model: formule numerice sau SQL scalar, unitate, precizie si optiune explicita pentru lipsa ca zero; dependente in ordinea coloanelor. Coloane concediu_<cod> generate din tipurile definite, inclusiv tipuri istorice inactive. Total concedii pe tipuri si categorii, fara dublare pe marca/data. Ore lucratoare = ore contractuale minus concedii care reduc obligatia; deplasarile creditate nu se scad.

Garzile efective active inlocuiesc planul aceleiasi marci/date. Lipsa efectivelor foloseste planul activ. Intersectii cu luna, weekend/sarbatoare 00-24 fara dublare; efectivele respecta fusul orar si DST. Tabelul garzi_zile permite reguli suplimentare prin SQL.

Aspect tabel automat, subset Markdown documentat sau elemente pozitionate numeric/prin tragere, inclusiv repetare pe fiecare rand si paginare. Campurile ascunse raman disponibile sablonului. Antet/subsol/imaginile existente sunt pastrate. Explicatiile SQL/formulelor si istoricul sunt pe anexa finala. Salvarea din editor cere previzualizarea definitiei curente; modificarile invalideaza previzualizarea.

Datele se exporta prin fluxurile existente PDF/Excel. Excel pastreaza tabelul rezultat, nu coordonatele vizuale PDF. Problema deschiderii Excel de pe sistemul utilizatorului ramane amanata conform cererii lui.

## Probe

- gofmt -l -w ., go vet ./..., go test -p 1 -timeout 120s ./...: trecute.
- TypeScript/Vite: trecut. Vitest: 24 teste, 9 fisiere, trecute.
- Edge: editor flexibil si regresie UX30, in romana si engleza, trecute; fara erori JavaScript.
- Go: formule, SQL izolat/interdictii/limite, previzualizare fara scrieri, persistenta modelului, deduplicare concedii si coloane dinamice, garzi efective/DST, export si validare PDF/paginare pe toate cele trei aspecte, campuri repetate/ascunse.
- Catalog tehnic nou modele-v6; modele-v5 si cazurile confirmate anterior sunt pastrate.
- Documentatie: 0 erori; manuale: 3 teste; livrare: 7 teste, trecute.
- Randare vizuala PDF Markdown verificata; exemple anonime in build/bin/verificare-F09-T13. Nu s-a accesat baza reala si nu s-a probat imprimanta fizica.
- Dependente: Go 1.27.3 si dependentele npm directe actuale. Cele 67 actualizari indirecte Go raman in F14-T21 pentru compatibilitate. gofpdf 1.4.3 existent este acum import direct; licenta MIT existenta reverificata, fara instalare noua.

## Predare

Ghid: [rapoarte flexibile](ghid-rapoarte-flexibile.md). Pasul urmator: proba operator pe copii de date si revizuire pentru integrare; fara push sau merge. Buildul oficial si sumele de control se consemneaza dupa compilare.
