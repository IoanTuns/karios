# Audit UI/UX Pontix 0.92.18

Data: 2026-10-31. Responsabil: Codex. Task: F16-T28. Analiza si recomandari; fara implementare in acest task.

## Metoda si limite

Inspectie euristica a capturilor trimise de utilizator si a codului curent. Ecrane analizate: Angajati/formular, Schimburi/Categorii, Planificare/sumar, Pontaj, Concedii, actiuni de administrator si ajutor. Nu s-au accesat date reale. Nu s-au copiat capturi cu nume de angajati in repo. Nu s-au efectuat sesiuni observate cu operatori, masuratori noi de contrast sau teste cu cititor de ecran. Capturile pot proveni dintr-o versiune anterioara; constatarile confirmate in cod sunt diferentiate de recomandarile euristice. Nu este o certificare de accesibilitate.

Concluzie: aplicatia are controale utile pentru operare repetitiva, dar cere utilizatorului sa cunoasca prea multe reguli interne. Prioritatea este explicarea consecventa a identitatii angajatului si a efectelor actiunilor, apoi reducerea densitatii vizuale.

## Avantaje

- Grile lunare si sumar pe persoana, potrivite lucrului repetitiv cu multe zile; introducere rapida si selectie in lot.
- Filtre profesionale si cautare dupa identitate/marca; stare activa implicita, acces explicit la istoricul inactiv.
- Separarea marcilor de program normal si garda permite evidenta distincta.
- Inspector de celula, legenda si jurnal pentru investigarea modificarilor; previzualizare la publicare.
- Zoom global plus densitate S/M/L pentru grile; dialoguri redimensionabile, fluxuri prin tastatura si indicii de modificari nesalvate existente in cod.
- Ajutor accesibil din antet si manual versionat; regulile raportului sunt explicate in legenda/subsol.
- Protectii backend pentru luni inchise, revizii, referinte si stergere rezervata administratorilor. Acestea reduc riscul, chiar daca explicarea lor in UI poate fi imbunatatita.

## Probleme si recomandari prioritizate

P0 = interpretare operationala gresita posibila; P1 = eroare ori efort repetat; P2 = claritate si confort. Severitatea este o apreciere euristica, nu o incidenta masurata.

| ID | Prioritate | Observatie si dovada | Efect | Recomandare si criteriu de acceptare |
| --- | --- | --- | --- | --- |
| UI28-01 | P0 | Confirmat: numeAngajat.ts foloseste doar angajat_principal_id; FiltreProfesionale.tsx foloseste si garda derivata din contract. Rapoartele folosesc asocierea pentru sufix. | Persoana apare filtrata ca Garda, dar numele nu are acelasi marcaj. | O singura regula de identitate de afisare, bazata pe statutul efectiv de garda. Teste: contract Garda fara legatura, marca asociata, contract normal si contract personalizat cu atribut garda; aceeasi identitate in liste, selectoare si rapoarte. Nu modifica numele salvat. |
| UI28-02 | P0 | Captura arata persoane cu contract Garda fara indicatia marcii principale; in cod acestea sunt doua proprietati distincte. | Utilizatorul poate presupune ca exista asocierea si ca medicalul se aplica automat tuturor marcilor. | Afiseaza separat Contract: Garda si Marca principala: asociata/neatribuita. Cere asociere numai cand exista o marca principala pentru aceeasi persoana; permite contractul independent. La medical arata marcile efectiv afectate. Nu crea legaturi dupa similitudinea numelor. |
| UI28-03 | P0 | Confirmat: Planificarea are salvare si publicare separate; publicarea retrage acum si pontaje efective. | Salvarea poate fi confundata cu sincronizarea, iar efectul publicarii poate fi subestimat. | Stari explicite Nesalvat / Salvat in plan / Publicat in Pontaj. Actiuni Salveaza planul si Publica in Pontaj. Previzualizare cu valori vechi/noi, numar de pontaje efective retrase, absente prioritare si randuri blocate. O modificare concurenta cere previzualizare noua. |
| UI28-04 | P1 | Confirmat: programePermise filtreaza simultan activ, programPlanificat si categoria; denumirea programului nu stabileste tipul sau. | Un program activ si bifat poate lipsi fara explicatie; exemplu program de la 22:00 care nu este clasificat Noapte. | Arata de ce un program nu este disponibil: inactiv, in afara categoriei, clasificare incompatibila. Link contextual catre configurarea programului. Criteriu: niciun program bifat disparut fara motiv accesibil. Validarea sa fie comuna cu backendul. |
| UI28-05 | P1 | Confirmat: formularul angajatului reuneste identitate, asociere, functie, sectie, compartiment, contract, norma, categorie si initializare. | Formular lung; soldurile de configurare initiala concureaza cu editarea curenta. | Grupeaza in Identitate, Incadrare, Program, Solduri initiale. Ultimul grup retractabil, cu rezumat al valorilor salvate si luna de referinta. Pastreaza butoanele Salvare/Anulare vizibile si avertizarea la iesire. |
| UI28-06 | P1 | Capturi: multe selectoare si randuri de unelte inaintea tabelului. Codul confirma filtre profesionale si controale distincte pe ecrane. | Grila utila ocupa mai putin ecran; filtrele ramase active pot trece neobservate. | Cautare, perioada si stare vizibile permanent; filtre avansate intr-un panou retractabil. Afiseaza filtrele active, numarul rezultatelor si Resetare. Pastreaza explicite filtrele restrictive, inclusiv Inactive. |
| UI28-07 | P1 | Confirmat: completarea Pontajului ofera programe, ore, absente si actiuni rapide in acelasi context. Prioritatile sunt mai complexe decat simpla editare de celula. | Selectia, tipul actiunii si amploarea suprascrierii pot fi greu de anticipat. | Alege explicit Program / Concediu / Medical / Ore manuale; afiseaza numele/marcile, intervalul si numarul zilelor afectate. Pentru medical arata anularea turelor intersectate si marcile asociate; explica excluderea weekendurilor/sarbatorilor. Nu cere confirmari pentru fiecare celula banala. |
| UI28-08 | P1 | Confirmat: initializarea cere luna, ore semnate si CO anterior; raportul are L Inceput/L Fin, iar configurarea normelor ramane curenta. | Diferenta dintre drept anual, CO consumat si sold ore poate fi interpretata gresit. | Rezumat de formula langa solduri, cu exemple +2:30/-2:30 si echivalent zecimal cand este folosit. CO initial sa arate anul si perioada acoperita. Indica norma folosita la recalcularea istoricului. Teste pentru semn, minute, trecere de an si CO dublat. |
| UI28-09 | P2 | Capturile arata folosirea culorii/insignelor, dar nu toate diferentele semantice sunt explicite textual. | Verde, fundalul sau pictograma pot fi confundate cu salvarea/validarea. | Text de stare langa culoare; termeni distincti pentru Activ, Salvat, Aprobat, Publicat. Culoarea este un ajutor suplimentar. Masurare de contrast pe ambele teme, inclusiv selectii/erori si butoane dezactivate. |
| UI28-10 | P2 | Confirmat: ButonStergereAdmin este prezent si dezactivat pe configurari active; foloseste window.confirm. | Actiuni destructive repetate incarca listele, iar dialogul nativ are context limitat. | Grupeaza actiunile rare in meniul randului. Stergere definitiva pentru admin la inregistrari inactive; explica utilizarea blocanta cu navigare spre sursa cand exista drept de acces. Pastreaza autorizarea backend, verificarea reviziei si confirmarea numelui. |
| UI28-11 | P2 | Confirmat in CSS: unele indicii folosesc 10px; formularele si tabelele sunt dense. Nu s-a masurat lizibilitatea pe dispozitivul utilizatorului. | Textul secundar poate fi greu de citit la scalare ridicata sau in fotografii/ecrane mici. | Verifica un minim practic de 12px pentru indicii si 14px pentru formulare, pastrand densitate reglabila in grile. Probe la 1024x768 si 1366x768, Windows 125/150%, zoom 100/150/200%. Fara actiuni inaccesibile sau text taiat fara alternativa. |
| UI28-12 | P2 | Confirmat: exista ajutor contextual si configurari distribuite in mai multe sectiuni. | Manualul nu poate explica singur o lista care ascunde optiuni sau o configurare incompleta. | Ajutor contextual pentru campul problematic, linkuri intre contract/categorie/program, verificare de configurare cu probleme explicate si actiuni de reparare. Nu transforma fiecare operatiune intr-un asistent cu pasi inutili. |

## Fluxul recomandat de configurare

1. Defineste functiile si contractele, cu norma si atributul de garda.
2. Defineste programele: ora, durata, tip si regim; verifica interpretarile tip/regim inainte de alocare.
3. Grupeaza programele permise in categorii.
4. Incadreaza angajatii si asociaza explicit marcile aceleiasi persoane, unde exista.
5. Introdu soldurile initiale pentru luna de pornire, separat de dreptul anual CO.
6. Pregateste Planificarea, verifica previzualizarea si publica in Pontaj.
7. Inregistreaza concediile/medicalul, verificand efectele prioritatii.
8. Verifica raportul si soldurile, apoi inchide luna.

Acesta este un flux propus pentru claritate, nu un blocaj care sa impiedice corectarea ulterioara a configurarilor.

## Ordinea propusa de implementare

- Lot A: UI28-01, UI28-02, UI28-04. Identitate consecventa, asociere explicita si explicarea programelor lipsa. Beneficiu imediat pentru problemele raportate de utilizator.
- Lot B: UI28-03 si UI28-07. Clarifica efectele salvarii/publicarii si ale concediilor, cu previzualizare proportionata actiunii.
- Lot C: UI28-05, UI28-06, UI28-08. Reduce densitatea si separa configurarea initiala de operarea lunara.
- Lot D: UI28-09 ... UI28-12. Finisaje vizuale, accesibilitate, actiuni rare si ajutor pentru configurare.

Prioritatile sunt recomandari pentru review. Nu s-au creat automat taskuri de implementare si nu s-au schimbat formule, roluri sau date.

## Cum se valideaza rezultatul

Scenarii cu date anonime: creare angajat normal; contract Garda independent; doua marci asociate; gasirea unui program activ exclus din lista; programare, salvare si publicare; conversie in concediu si apoi medical cu tura de noapte; initializare sold negativ/CO; consultare Inactive/reactivare; inchidere de luna si incercare de editare.

Se masoara finalizarea fara asistenta, erorile de interpretare, pasii inutili si revenirea din eroare. Nu se inventeaza un scor de utilizabilitate sau estimari de imbunatatire procentuala. Probe cu tastatura, ambele teme si dimensiunile mentionate; cititor de ecran si contrast se verifica separat.

## Verificari pentru acest audit

Doar documentatie. Controlul documentatiei a trecut cu 0 erori, iar git diff --check a trecut; nu se afirma rerularea testelor aplicatiei sau o proba noua WebView2. Pasul urmator: utilizatorul alege lotul de implementare sau confirma prioritatile. Fara push/merge.
