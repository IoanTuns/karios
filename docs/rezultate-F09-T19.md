# F09-T19 - Prezentare, copii editabile si calendar extins

Responsabil Codex. Implementare2026-11-03/04. Ramura F09-T19-editare-prezentare. Aplicatie0.92.30, manual1.40. Rapoarte pontaj15, detaliat16, individual15, absente12, garzi12, flexibil4. Fara migrare sau biblioteca noua. D085 propusa; regula intervalului si editarea celulelor confirmate explicit de utilizator.

## Implementare

Modelele au stiluri pe camp (bold, italic, marime6-40pt, aliniere L/C/R, mostenire) si etichete RO/EN pentru coloane. PDF standard/flexibil si Excel folosesc aceleasi valori si optiuni; fonturile Go italic/bolditalic provin din pachetul/licenta deja folosite. In Markdown, stilurile titlului/paragrafelor se configureaza din editor; nu se introduce HTML executabil sau un editor Word pentru fragmente de text. Elementele LIBER au stiluri proprii; campurile repetate preiau stilul randului si modificarile celulelor corespunzatoare.

Raportul generat are editor de titlu, antet/subsol, introducere, etichete si celule, cu25 randuri pe pagina. Tastarea pastreaza cursorul. Modificarile se aplica atomic pe o copie profunda server, nu pe Document furnizat de client. Sunt acceptate doar campuri existente, maximum1000 modificari/2MiB,4000 caractere per camp. Validare stil/aspect/PDF inaintea jurnalizarii; token nou asociat jurnalului si istoricului PDF, originalul din cache neschimbat. Contextul cererii si valorile/stilurile inainte/dupa sunt in jurnal. Istoricul este nemodificabil din acest editor, cu autor si moment local; marcaj scurt de personalizare in nota, detalii la final. Exportul PDF/Excel foloseste copia aplicata. Cu modificari neaplicate, exportul si schimbarea filtrelor sunt dezactivate; se poate aplica sau renunta.

Editarea copiei nu schimba datele operationale si nu recalculeaza alte celule. Copia se pastreaza prin export; regenerarea reia sursa, iar stilurile/etichetele modelului raman persistente. Exemplul SQL principal include data_calendar ca data ISO, ascunsa din document. Boldul optional pentru weekend/sarbatori este calculat inaintea ascunderii coloanei, nu din eticheta numerica a zilei.

Optiunea de extindere pastreaza grupurile consecutive nelucratoare la ambele margini ale lunii, conform calendarului configurat. Noiembrie2026 acopera31octombrie-1decembrie; decembrie2026 acopera28noiembrie-31decembrie. Parametrii SQL de_la/pana_la reflecta intervalul; totalurile lunare/contractuale raman ale lunii alese. Nu se sare peste zile lucratoare. Limita tehnica366 zile suplimentare, cu eroare. SQL principal pentru zile cu garda si varianta calendar completa pentru randuri goale, ambele verificate in motorul efectiv si incluse in manual.

## Verificari

- Teste calcul: ambele margini, grupul noiembrie/decembrie2026, lipsa salturilor peste zile lucratoare, interval invalid si limita tehnica. Trecute.
- Teste raport: copie profunda/atomica, patch invalid/istoric protejat, persistenta stiluri/etichete RO/EN, inainte/dupa in jurnal, formule Excel tratate ca text, fonturi/stiluri PDF/Excel si paginare, bold pe data ascunsa, extindere cu totaluri lunare identice. Trecute.
- Test API: token invalid refuzat, token nou, PDF original neschimbat, legatura copie-jurnal, lipsa jurnalizarii/aplicarii unui PDF care nu incape. Trecut.
- TypeScript/Vite,24 teste frontend si scenarii Edge RO/EN trecute: model formatat, bold/calendar extins persistate, editare celule, cursor la fiecare tastare, blocarea exportului pana la aplicare, aplicare si renuntare.
- Exemplu PDF anonim randat Windows si inspectat vizual: titlu central mare, rand obisnuit normal si28-30noiembrie/1decembrie bold, ore24h. Probe in build/bin/verificare-F09-T19. Date reale/imprimanta fizica neaccesate.
- Documentatie0 erori,3 teste manuale,7 teste livrare trecute. Gofmt/vet si intreaga suita Go trecute dupa corectia i18n (go test -p 1 -timeout 180s ./...). Regresia UX30 RO/EN trecuta.
- Prima rulare completa: coliziune i18n coloana_data cu eticheta existenta Data/Date; reparata in cod cu cheie distincta camp_data_evidentiere. O regenerare a referintelor a fost respinsa automat; verificarea fara modificari a demonstrat randuri si structura identice cu versiunile salvate, numai versiunea din metadate diferita. Reincercarea limitata la doua referinte NOI ne urmarite/absente din main a fost aprobata automat. Referintele existente/cazurile contabile protejate nu au fost modificate.
- Timeout izolat in stocare, la FindFirstFile Windows; reluarea a trecut in15.407s. Testele/codul de stocare neschimbate.
- Dependente directe npm actuale; Go1.27.3/Wails2.17.1/gofpdf1.4.3 verificate. Excelize2.11.1 ramane amanat pentru dependenta efp prerelease; indirectele in F14-T21. Fara instalari noi.

## Predare

Pas urmator: proba operator pe modelul propriu si review pentru integrare. Fara push sau merge. Buildul oficial si sumele de verificare se consemneaza dupa compilare.
