# F07-T25 - Structura profesionala, contracte si reguli conditionate

Responsabil Codex. Ramura F07-T25-structura-contracte. Aplicatie 0.92.0, manual 1.11, migrare 0026. Cerintele si lista initiala sunt in plan-F07-T25.md; arhitectura este propusa in D062.

## Implementat

- Grupuri profesionale ordonate numeric; functii asociate grupurilor, conducere/executie. Editor auditat cu revizie si dezactivare.
- Contracte configurabile: normare globala, individuala, zilnica fixa, lunara si fara norma. Parametrii se salveaza in minute. Garda fara norma predefinita; marcile deja legate explicit sunt migrate la contract GARDA, fara potrivire dupa nume.
- Categoria programelor permise din contract este aplicata de backend in Pontaj si Planificare. Precompletarea foloseste norma zilnica aplicabila, fara inventarea unei norme pentru Garda sau contracte lunare. Contractul si norma pot fi selectate drept coloane de raport.
- Filtre de grup, functie, pozitie, contract si garda in Angajati/Pontaj/Planificare/Rapoarte. Lista din Planificare si graficul PDF primesc aceleasi filtre. Filtrul Garda priveste incadrarea, nu existenta unei ture.
- Rapoarte generale cu antete de grup, conducerea prima si restul alfabetic romanesc; marca principala urmata de marcile de garda. Totalurile graficului respecta persoanele filtrate. Exportul Excel din Rapoarte foloseste exact documentul previzualizat, inclusiv gruparea si filtrele; valorile sunt celule text, nu formule preluate din date.
- Reguli conditionate configurabile pe profil, coloana, sursa, coduri si factor. Sursa ORE_SCHIMBURI_ZNL numara numai intersectia programelor selectate cu weekend/sarbatori, pe fragmente calendaristice. Factor 2 produce total 200% in ore echivalente, fara dublarea T. Regulile conditionate inlocuiesc regula generala pentru coloana respectiva; suprapunerile active sunt raportate, nu rezolvate arbitrar.
- Rapoarte pontaj/detaliat/fisa v5, absente/garzi v3, referinte noi. Referintele anterioare raman nemodificate.

## Verificari

Testele noi pentru migrarea unei baze vechi, persistenta grup/contract/reguli, contract dezactivat, norma de garda, profilul real din baza si total 200% au trecut. Motorul testeaza tura vineri 22-07: 9 ore lucrate, 7 eligibile, 14 echivalente; sarbatoare/weekend suprapuse fara dubla numarare; profil neeligibil; reguli suprapuse. Testele rapoartelor confirma ordinea grupurilor, sef inaintea personalului alfabetic, perechea normal/garda si filtrele. PDF-ul grupat este validat cu pdfcpu, randat Windows si inspectat vizual. Excel verificat prin redeschidere, inclusiv absenta formulelor injectate.

UI Edge cu punte simulata RO/EN: creare grup, functie de conducere, atribuire angajat, marca de garda, filtre si salvarea unei reguli active cu factor 200%. Backendul real este verificat separat in Go. Focus/cursor/zoom si ajutor UI trecute. TypeScript/Vite, 18 teste frontend, 3 teste manual si 7 teste livrare trecute. Control documentatie fara erori. gofmt/vet trecute; suita Go completa (`go test -p 1 -timeout 120s ./...`) trecuta dupa corectarea reviziilor monotone.

Dependente verificate: Go 1.27.3 si npm actuale. Actualizarile tranzitive Go sunt urmarite in F14-T21; fara biblioteci noi.

## Limite si predare

Nu exista calcul in lei fara baza orara. Conversia ZNL in zile compensatorii, acumularea si acordarea lor sunt blocate in F07-T26, cu intrebarile in F01. ZNL continua sa reprezinte ore eligibile, nu zile libere deja efectuate. Nicio formula neprecizata nu a fost inventata. Grupurile/functiile si regulile se configureaza explicit de operator; nu sunt deduse din denumiri.

Configurarea este curenta, fara istoric contractual pe intervale; modificarile pot afecta regenerarea lunilor vechi. Exportul clasic din grila Pontaj ramane registrul lunar complet; exportul filtrat/grupat se face din Rapoarte dupa previzualizare. Proba interactiva desktop/WebView2 si tiparirea fizica raman la operator. Datele reale nu au fost deschise. Buildul oficial urmeaza dupa commit; fara push sau merge.

Verificare suplimentara: ceas egal sau regresat produce o revizie strict noua; functia proprie fara grup nu mosteneste grupul principal in profilul de calcul. Graficul de garzi exclude programul DUPA_AMIAZA, cu totaluri corecte dupa filtrare.
