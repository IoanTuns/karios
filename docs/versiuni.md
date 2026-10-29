# Versiuni

Istoricul versiunilor aplicației, cea mai nouă prima ([D009](decizii/D009-versionare.md)). Fiecare commit care schimbă codul crește versiunea din `internal/versiune/versiune.go` și adaugă aici o secțiune `## X.Y.Z · AAAA-LL-ZZ`.

Pe lângă versiunea aplicației, fiecare raport și fiecare schimbare a bazei de date au versiunea lor; ele se trec în secțiunea versiunii aplicației care le aduce.

## 0.85.2 · 2026-10-29

Întrebări despre o zi anume, cu răspuns din cod (F14-T28).

- „Cine lipsește pe 15 octombrie?”, „Cine lucrează pe 15?”, „Ce avem pe 15 octombrie?” primesc lista exactă a absenților și/sau a celor programați în acea zi, fără model.
- Pentru o zi din afara lunii de lucru, asistentul spune că poate răspunde doar pentru luna de lucru și cum o schimbi, în loc să ghicească.
- Setul de comenzi reale: întrebarea „Cine lipsește pe 15 octombrie?” trece de la model la răspunsul din cod.

Baza de date și rapoarte: nicio schimbare.

## 0.85.1 · 2026-10-29

O singură poartă pentru mesajele asistentului și setul de comenzi reale (F14-T27, F14-T01).

- `ai.Ruteaza` alege calea fiecărui mesaj într-o ordine fixă: cerere de concediu → planificare ghidată → întrebare la care răspunde codul → model. Înainte, aceeași sinteză rula și în aplicație, și în clientul modelului, iar căile se puteau contrazice.
- Setul de comenzi reale (`internal/ai/testdata/comenzi.json`, 30 de comenzi, inclusiv toate eșecurile din capturile Ioan Tuns) verifică ruta aleasă la fiecare `go test`. Fiecare eșec nou raportat se adaugă acolo.

Baza de date și rapoarte: nicio schimbare.

## 0.85.0 · 2026-10-29

Cererile de concediu din asistentul AI (F14-T26, [D052](decizii/D052-concedii-din-asistent.md)).

- „X merge în concediu …” e recunoscut ca cerere (nu ca întrebare) înaintea celorlalte căi. Angajatul se identifică după marcă sau nume complet, cu toleranță la greșeli („Maria Daria” → Marin Daria); intervalul, din zilele și luna din mesaj; tipul CO, sau CM la „medical”.
- Cardul are butonul „Deschide formularul de concediu”: ecranul Concedii și absențe se deschide cu formularul nou precompletat (angajat, tip, interval). Numărul de înregistrare e obligatoriu pentru aceste cereri; starea rămâne „în așteptare”. Salvarea trece prin validările obișnuite.
- Comenzile incomplete primesc o întrebare (ce angajat, ce zile), nu lista concediilor.
- Lista concediilor răspunde doar la întrebări și nu mai afirmă că „tot personalul este activ și la dispoziție”.

Baza de date și rapoarte: nicio schimbare.

## 0.84.0 · 2026-10-29

F14-T25: instrumentele asistentului cerute de utilizator.

- Anularea ultimei operații AI prin instantanee persistente, verificare și confirmare tranzacțională, cu refuz integral la conflicte sau luni închise; audit păstrat.
- Explicații pentru selecție, rezumate pe luni și persoane, căutare și filtre în detalii și selecția nominală.
- Confirmarea anului preluat din context, validarea perioadei și șabloane locale încărcate în editor.
- Verificare consultativă a pontajului existent, modelul efectiv pornit/context configurat și export CSV local verificat de backend.

Migrarea 0023 adaugă operațiile și șabloanele AI. Exportul CSV al planului are formatul propriu 1.0.0, inclus în fișier. Nu se schimbă regulile salariale sau versiunile rapoartelor existente. D050 propusă; probele și limitele sunt în documentația F14-T25.

## 0.83.0 · 2026-10-29

F14-T24, I072–I074: primul lot după review-ul UI al asistentului.

- Planul activ are zonă proprie pe lățimea utilă, cu etape Interpretare / Verificare / Rezultat, context vizibil și acțiuni în afara conținutului derulabil. Istoricul conversației rămâne disponibil separat.
- Editarea intenției păstrează poziția și focusul. Derularea automată a istoricului se face doar pentru mesaje noi când utilizatorul urmărește capătul, fără deplasarea paginii.
- În rezultatul salvat se afișează valoarea finală și stările Salvat / Păstrat / Neaplicat / Deja identic. Celulele excluse sau păstrate nu sunt prezentate ca valori propuse aplicate.

Verificarea tokenului, clarificările și confirmarea explicită rămân active. Fără schimbări ale bazei de date, calculelor sau rapoartelor. Îmbunătățirile P2 din review rămân în backlog.

## 0.82.1 · 2026-10-29

F14-T22: actualizare separată a dependențelor folosite de aplicație: `go-sqlite3` 0.35.6 → 0.36.0, `go-sqlite3-wasm/v6` 6.3.35401 → 6.4.35401, `mscfb` 1.0.7 → 1.0.9, `efp` 0.0.1 → 0.0.2 și `x/net` 0.61.0 → 0.62.0. Go 1.27.3 și pachetele npm sunt deja la versiunile stabile raportate la verificare. Fără schimbarea schemei bazei de date, rapoartelor sau comportamentului interacțiunii AI din 0.82.0.

## 0.82.0 · 2026-10-29

F14-T22: interacțiune ghidată a asistentului AI, D048 propusă.

- Interpretarea apare într-un formular editabil: persoane, grup, program și interval, inclusiv între luni și ani.
- Zilele și comportamentul pentru celulele ocupate se aleg explicit când lipsesc. Selecțiile disponibile: toți activii, secție, departament, categorie, persoane nominalizate și filtrul fără pontaj pe lună.
- Previzualizarea arată persoanele, lunile, valorile vechi/noi și excluderile. Confirmarea verifică amprenta datelor și salvează tranzacțional întregul plan; scrierea este recitită înainte de commit.
- Contextele modelului sunt limitate și separate pe planificare, absențe, pontaj sau ajutor; analizele deterministe folosesc datele complete. Cererea curentă nu se trunchiază.
- 40 de comenzi RO/EN de regresie și teste pentru date schimbate, selecții, absențe, celule ocupate și rollback între luni. Cardurile vechi trec prin revizuire înainte de aplicare.

Baza de date și rapoartele nu își schimbă schema sau versiunea. Nu au fost adăugate biblioteci. Calitatea modelului real și consumul pe calculator cu 8 GB rămân de măsurat separat.

## 0.81.2 · 2026-10-29

F14-T20: contextul motorului AI crește de la 8192 la 32768 tokenuri. La depășire, clientul reîncearcă fără schimburile cele mai vechi, păstrând ultima cerere și datele aplicației. Fără modificări ale bazei de date sau rapoartelor.

## 0.81.1 · 2026-10-29

Corecturi pe ecranele contului (F15-T09).

- Crearea primului administrator citește parolele direct din formular la trimitere: completarea automată a WebView2 poate scrie în câmp fără să actualizeze starea, iar formularul raporta „Parolele nu coincid” cu parole identice. Câmpurile de parolă și cele de autentificare ascultă și de `onInput`.
- Buton „Arată / Ascunde” lângă fiecare parolă.
- Chenarul întrebărilor de recuperare rămâne în lățimea cardului (`fieldset` avea lățimea minimă a conținutului), iar textul de ajutor e mai mic.

Baza de date și rapoarte: nicio schimbare.

## 0.81.0 · 2026-10-29

Propunerile repetate tratate după starea lor (F14-T19).

- Se omit doar propunerile identice cu unele încă neconfirmate din conversație (se confirmă cardul de mai sus); mesajul spune asta.
- O propunere identică cu una deja aplicată se arată, cu nota „O propunere identică a fost deja aplicată mai sus. Confirmă doar dacă vrei să o aplici din nou.” Utilizatorul nu mai rămâne fără răspuns când cere explicit reaplicarea.

Baza de date și rapoarte: nicio schimbare.

## 0.80.0 · 2026-10-29

Detectarea răspunsurilor care pretind o planificare fără nicio acțiune (F14-T18).

- Dacă modelul răspunde fără nicio propunere, dar textul afirmă că a planificat, aplicat sau salvat ceva (de ex. „Planificare interval Z7 … (deja aplicat)”), răspunsul se înlocuiește cu un mesaj cinstit: nu s-a modificat nimic, cu sugestii de reformulare.
- Verificarea rulează după încercarea existentă de a construi propunerea din text.

Baza de date și rapoarte: nicio schimbare.

## 0.79.0 · 2026-10-29

Istoricul trimis modelului fără text de imitat (F14-T17).

- Mesajele asistentului cu propuneri ajung la model doar cu textul lor (sau „(propuneri afișate utilizatorului)”), fără descrierea propunerilor.
- Propunerile din conversație, cu starea „[aplicată]” sau „[neconfirmată]”, se trimit separat (`CerereInterogareAI.actiuni`) și apar în promptul de sistem ca listă informativă.
- Promptul cere explicit să nu se scrie în text că s-a planificat sau aplicat ceva; modificările se fac doar prin unelte.

Baza de date și rapoarte: nicio schimbare.

## 0.78.0 · 2026-10-29

Filtrul „fără pontaj” și completarea doar a celulelor goale (F14-T16).

- „toți angajații fără pontaj …” (și „secția X fără pontaj”) restrânge grupul la angajații care nu au nimic completat în pontajul lunii vizate. Dacă toți au deja pontaj, asistentul o spune în loc să propună.
- Propunerile de grup completează implicit doar zilele goale: valorile deja existente în pontaj nu se suprascriu (`doar_goale`), iar concediile și absențele se sar ca înainte.

Baza de date și rapoarte: nicio schimbare.

## 0.77.0 · 2026-10-29

Comenzile asistentului recunoscute și cu greșeli de tastare (F14-T15).

- Înainte de recunoaștere, fiecare cuvânt fără cifre se compară cu un vocabular (luni, zile ale săptămânii, „toți”, „angajat”, „secția”, „categoria”, „weekend”, „fără”, „pontaj”… și numele secțiilor și categoriilor) prin distanța Damerau-Levenshtein, scrisă în cod. Prag: până la 4 litere exact, 5–7 o greșeală, de la 8 două; la egalitate între cuvinte diferite nu se corectează nimic, iar între formele aceluiași cuvânt se ia cea mai scurtă.
- Grupul „toți” se recunoaște după cuantificator + substantiv în orice formă („toti angajati”, „tot personalul”). „mai” e lună doar în context clar („în mai”, „luna mai”, „mai 2026”).
- Codurile de schimb și orele rămân exacte. Mesajul trimis modelului rămâne cel original.
- Cardul propunerii de grup arată „Am înțeles: …” și corecturile („octombie” → octombrie).

Baza de date și rapoarte: nicio schimbare.

## 0.76.2 · 2026-10-29

Măsurare pe o bază mare (F06-T10).

- Test de măsurare la cerere (`PONTIX_MASURARE=1 go test -run MasurareBazaMare -v ./internal/datebaza/`) pe 300 de angajați și un an de pontaj (78.300 de celule).
- Rezultate pe calculatorul de dezvoltare: prima sigilare 0,44 s; verificarea la pornire 0,45 s; amprenta unei luni 0,06 s; citirea pontajului unei luni 0,17 s; salvarea unei luni întregi (6.300 de celule) 2,4 s fără audit și 4,6 s cu trigger-ele auditului; verificarea după salvare 0,74 s.
- Concluzie: verificarea la pornire nu e o problemă; salvarea în bloc a unei luni întregi e lentă, iar trigger-ele aproape îi dublează durata (idee I066).

Baza de date și rapoarte: nicio schimbare.

## 0.76.1 · 2026-10-29

Legăturile Wails regenerate cu generatorul (F02-T15).

- `frontend/wailsjs` vine din nou din `wails generate module`, nu din editări manuale. Repară importul `conturi` care lipsea din `App.d.ts`; restul diferențelor sunt ordinea claselor și spațiile.
- Tipul evenimentului `ai:sarcina_terminata` (`RezultatSarcinaAI`) e definit în `asistentSesiune.ts`, fiindcă generatorul include doar tipurile din semnăturile metodelor.
- Ghidul de build descrie regenerarea.

Baza de date și rapoarte: nicio schimbare.

## 0.76.0 · 2026-10-29

Ultima modificare a celulei în bara de inspecție a pontajului (F08-T10).

- Pentru celula activă, bara de inspecție arată cine și când a modificat-o ultima dată, din jurnalul de modificări (cont, „utilizator local” sau „asistent AI, confirmat de …”). Se recitește după salvare.
- `datebaza.UltimaModificare` și `App.UltimaModificareCelula`; traducerea autorului e comună cu ecranul jurnalului (`frontend/src/autorJurnal.ts`).

Baza de date și rapoarte: nicio schimbare.

## 0.75.0 · 2026-10-29

Textele din interfață mutate în cataloage (F03-T10, UI-4 din review).

- Aproximativ 190 de chei noi în `ro.json` și `en.json`: bara laterală și tema, notificările, zoom-ul grilelor, tooltip-urile celulelor, ecranele goale din Pontaj și Planificare, cardul „Despre” și dialogul cu detaliile de dezvoltare (inclusiv raportul copiat), dialogul regulilor din planificare, comparația plan/realizat, harta concediilor, etichetele din cardurile asistentului AI.
- Valorile descriptive din detaliile de dezvoltare vin din Go în limba interfeței (`DetaliiDezvoltare(limba)`); dialogurile native (salvarea exportului Excel, închiderea cu modificări nesalvate) folosesc limba interfeței, trimisă cu `SeteazaLimbaInterfata`.
- Erorile la aplicarea propunerilor AI au coduri proprii (`ai.interval_invalid`, `ai.data_invalida`, `ai.tip_necunoscut`, `ai.fara_zile`) în locul textelor românești din `detalii`.
- Jurnalul de diagnoză: `diagnoza.UltimeleLinii` întoarce o eroare, iar mesajul pentru utilizator îl alege aplicația.

Rămâne în română: conținutul generat de asistentul AI (răspunsurile calculate din cod și cele ale modelului), care urmează limba asistentului.

Baza de date și rapoarte: nicio schimbare.

## 0.74.0 · 2026-10-29

Managerul motorului AI creat o singură dată (F14-T14, SEC-6 din review).

- `asiguratManagerAI` folosește `sync.Once`: pornirea automată din fundal și apelurile din interfață nu mai pot crea două managere și două procese `llama-server`, dintre care unul rămânea pornit după închidere.
- La închiderea aplicației se oprește procesul managerului unic.

Baza de date și rapoarte: nicio schimbare.

## 0.73.0 · 2026-10-29

Identificarea exactă a angajatului la aplicarea propunerilor AI (F14-T13, SEC-4 din review).

- La aplicare, angajatul se găsește doar după marca exactă (cu sau fără cratime) sau după numele complet („Nume Prenume” ori „Prenume Nume”), fără diacritice și majuscule. Potrivirile parțiale (de ex. „Ana” în „Ioana”) nu mai sunt acceptate.
- Mai mulți candidați dau eroarea `ai.angajat_ambiguu`, cu lista lor; niciunul dă `ai.angajat_negasit`. Angajații inactivi nu se iau în calcul.

Baza de date și rapoarte: nicio schimbare.

## 0.72.0 · 2026-10-29

Serverul AI local protejat (F14-T12, SEC-2 și SEC-3 din review).

- `llama-server` pornește pe un port liber ales la fiecare pornire și cu o cheie API aleatoare (`--api-key`), trimisă de client în antetul `Authorization`. Alte procese locale nu mai pot folosi motorul.
- Aplicația consideră activ doar procesul lansat de ea (urmărit până la oprire); un alt program care răspunde la `/health` pe portul 28491 nu mai primește datele angajaților. Starea motorului din ecranul asistentului și din detaliile de dezvoltare vine din procesul propriu.
- Numele modelului cerut la pornire se acceptă doar dacă e un fișier `.gguf` din folderul AI, fără căi.

Baza de date și rapoarte: nicio schimbare.

## 0.71.0 · 2026-10-29

Backup automat zilnic (F05-T26).

- La prima pornire din zi (UTC) se face o copie criptată `pontix-AAAA-LL-ZZ-hhmmss-zilnic-v….db`, dacă în acea zi nu există deja o copie automată (de ex. cea dinaintea unei migrări).
- Copiile zilnice intră în aceeași rotație cu backup-urile dinaintea migrărilor: se păstrează ultimele 7 automate; cele manuale nu se ating.
- Nu se face copie pentru o bază creată chiar la această pornire și nici peste o bază cu incident de audit, ca rotația să nu împingă afară copiile bune.

Baza de date și rapoarte: nicio schimbare.

## 0.70.0 · 2026-10-29

Restaurarea unui backup din aplicație (F05-T25).

- Setări → Backup arată copiile din `karios/backup/` (automate și manuale). Administratorul alege o copie; aplicația îi verifică sigiliile pe o copie temporară (D044) și refuză o copie modificată în afara aplicației.
- Restaurarea se programează și se aplică la următoarea pornire, cu baza închisă: baza curentă se salvează ca backup `pontix-manual-inainte-de-restaurare-…` (nu intră în rotația celor 7), copia îi ia locul, iar restaurarea se notează în jurnalul bazei restaurate și se resigilează.
- Programarea se poate anula; butonul „Închide aplicația acum” grăbește repornirea.
- Numele backup-ului se acceptă doar ca nume de fișier din `backup/`, fără căi.

Baza de date și rapoarte: nicio schimbare.

## 0.69.0 · 2026-10-29

Avertizare la închiderea ferestrei cu modificări nesalvate (F05-T24, UX-1 din review).

- `beforeunload` nu oprea închiderea ferestrei Wails, deci X-ul pierdea modificările. Acum `OnBeforeClose` întreabă, cu un dialog nativ în limba interfeței, dacă fereastra se închide fără salvare; implicit rămâi în aplicație.
- Pontajul, editorul de rapoarte și configurarea regulilor își raportează modificările nesalvate prin `frontend/src/modificariNesalvate.ts`.

Baza de date și rapoarte: nicio schimbare.

## 0.68.0 · 2026-10-29

Conturi și acces (F15, [D008](decizii/D008-criptarea-datelor.md) punctul 3, [D045](decizii/D045-conturi-detalii.md)).

- **Primul administrator:** la prima pornire fără conturi (inclusiv o bază existentă), aplicația cere crearea contului de administrator; datele rămân.
- **Autentificare:** nume și parolă; parola deschide copia `P:<id>` din seif (Argon2id), fără hash în bază. Parola: minim 8 caractere, cu literă și cifră.
- **Auto-login** pentru ultimul cont pe același utilizator Windows (SID) și calculator; „Ieșire din cont” îl oprește până la următoarea autentificare cu parolă.
- **Recuperare:** 3 întrebări alese dintr-o listă de 10; răspunsurile normalizate deschid copia `R:<id>`. Administratorul poate reseta parola oricărui cont.
- **Blocare progresivă:** de la a 5-a încercare greșită consecutivă: 1, 5, 30, apoi 120 de minute; fiecare încercare, reușită sau nu, intră în jurnal.
- **Roluri:** administratorul gestionează conturile (Setări → Conturi) și confirmă incidentele de audit; operatorul face restul. Ultimul administrator activ nu poate fi retrogradat sau dezactivat.
- **Contul meu:** schimbarea parolei și a întrebărilor; numele contului și ieșirea apar în bara laterală.
- **Jurnal:** autorul modificărilor devine numele contului (`ana`, `ana/asistent_ai`); închiderea lunii notează numele afișat.
- Până la autentificare, orice scriere în bază e refuzată de trigger-ele auditului, mai puțin operațiile de autentificare.

Baza de date: migrarea `0022_conturi` (tabelele `conturi` și `sesiune_locala`). Rapoarte: nicio schimbare.

## 0.67.0 · 2026-10-29

Autorul modificărilor aplicate de asistentul AI (F14-T11).

- Autorul din jurnal vine acum din context (`datebaza.CuActor`, `Actor`); implicit `local`, iar ce aplică asistentul după confirmare se notează `local/asistent_ai`. Codul stocat nu depinde de limbă; ecranul jurnalului îl afișează „asistent AI, confirmat de utilizator local”.
- Definițiile rapoartelor folosesc același autor din context.
- Când vin conturile (F15), `local` se înlocuiește cu numele contului, fără schimbarea formatului.

Baza de date și rapoarte: nicio schimbare.

## 0.66.0 · 2026-10-29

Motivul închiderii și al redeschiderii lunii (F12-T05).

- `window.confirm` și textele fixe din cod („Închidere lună salarizare”, „Redeschidere utilizator”, mesajele de succes) sunt înlocuite cu un dialog în care utilizatorul scrie motivul; la închidere e precompletat „Transmis la salarizare”, la redeschidere e gol.
- La redeschiderea unei luni închise, dialogul spune că amprenta tipărită nu va mai corespunde după modificări.
- `SeteazaStareLuna` refuză un motiv gol (`pontaj.motiv_obligatoriu`); motivul intră în jurnal.

Baza de date și rapoarte: nicio schimbare.

## 0.65.0 · 2026-10-29

Ecranul jurnalului de modificări (F12-T02), în Setări.

- Filtre: zona (tabelul), angajatul (inclusiv zilele lui din pontaj și absențe), luna vizată, căutare în valori; paginare câte 100, cele mai noi primele.
- Fiecare intrare arată când, cine, zona, rândul (angajatul și ziua, unde e cazul) și doar câmpurile schimbate, „înainte → după”; câmpurile tehnice (revizie, momente) se ascund.
- Starea sigiliului (verificat la pornire sau incident) apare deasupra listei.
- Export CSV (separator „;”, UTF-8 cu BOM, pentru Excel) al intrărilor filtrate, în folderul ales.

Baza de date și rapoarte: nicio schimbare.

## 0.64.0 · 2026-10-29

Luna din context și confirmare doar prin buton în asistentul AI (F14-T10, UX-2 și UX-3 din review).

- Luna de lucru a asistentului se preia din pontaj (ultima lună deschisă) sau, în lipsă, din planificare ori din data curentă; se vede și se poate schimba în antetul asistentului și rămâne în sesiune.
- La aplicare, luna se ia din datele propunerii, nu din ecran; un interval peste sfârșitul lunii se salvează pe fiecare lună.
- „da”, „ok”, „sigur” etc. nu mai aplică propunerile: asistentul răspunde local că se confirmă cu butonul de pe card.

Baza de date și rapoarte: nicio schimbare.

## 0.63.0 · 2026-10-29

Listări din cod și fără propuneri repetate în asistentul AI (F14-T09).

- „Listează / câți / ce angajați”, cu „secția X” sau „categoria Y” opțional, și „ce schimburi există” primesc răspuns calculat direct, fără model: angajații activi grupați pe secții (marcă, funcție, categorie), respectiv schimburile active.
- Propunerile identice (același tip și aceiași parametri) cu unele deja afișate sau aplicate în conversație nu se mai arată; mesajul spune câte au fost omise.
- În istoricul trimis modelului, propunerile apar ca „deja aplicată” sau „propusă, neconfirmată”, iar promptul cere să nu repropună ce e deja aplicat.

Baza de date și rapoarte: nicio schimbare.

## 0.62.0 · 2026-10-29

Propuneri de grup în asistentul AI (F14-T08).

- Comenzile pentru „toți angajații”, „secția X” sau „categoria Y” cu un schimb (cod sau interval, de ex. „de la 7 la 14”) pe o lună întreagă sunt recunoscute direct din cod (`ai.AnalizeazaComandaGrup`), fără model: o singură propunere `planifica_grup` cu lista persoanelor.
- Implicit doar zile lucrătoare (L–V); „toate zilele / inclusiv weekend” și „weekend” schimbă filtrul. Comenzile pentru o zi anume și intervalele ambigue rămân la model.
- La aplicare, angajații se identifică doar după marca exactă, zilele cu concediu sau absență activă se sar, iar pontajul se salvează într-o singură operație pentru tot grupul.
- Cardul propunerii arată numărul de persoane și lista lor (pliabilă).
- `AplicaPropunereAI`: calculul zilelor și salvarea extrase în `zileDinInterval` și `salveazaPropunere`, folosite de ambele tipuri de propuneri.

Baza de date și rapoarte: nicio schimbare.

## 0.61.0 · 2026-10-29

Asistentul AI lucrează în fundal (F14-T07).

- `PornesteSarcinaAI` întoarce imediat; răspunsul vine prin evenimentul Wails `ai:sarcina_terminata`. O singură sarcină o dată (`ai.sarcina_activa`); `OpresteSarcinaAI` întrerupe cererea către motor; `SarcinaAIActiva` permite reluarea după reîncărcarea interfeței.
- Conversația și sarcina în curs stau în `frontend/src/asistentSesiune.ts`, în afara ecranului: schimbarea ecranului nu mai pierde răspunsul. Pe alt ecran apare un toast, iar iconița asistentului din bara laterală are un semnal (în lucru / răspuns necitit).
- Buton „Oprește” lângă indicatorul de lucru. Propunerile nu se aplică singure.

Baza de date și rapoarte: nicio schimbare.

## 0.60.1 · 2026-10-28

Actualizarea dependențelor (AGENTS.md §6, F06-T08): `github.com/pdfcpu/pdfcpu` v0.16.3 → v0.16.4 (patch, aceeași licență Apache-2.0). `go mod tidy` a trecut `github.com/xuri/excelize/v2` v2.11.0 la dependență directă, cum e folosită de `internal/exportexcel`; versiunea nu se schimbă. Go 1.27.3 și pachetele npm erau deja la ultima versiune stabilă.

## 0.60.0 · 2026-10-28

Trasabilitate: orice modificare a datelor făcută în afara aplicației este detectată (F06-T04, F06-T08, F12-T04, [D044](decizii/D044-jurnal-audit-sigilat.md)).

- **Sigiliul tabelelor (`internal/audit`):** pentru fiecare tabel și pentru schemă, suma modulo 2^256 a HMAC-urilor rândurilor, cu cheia de audit derivată din cheia de date. Trigger-ele TEMP create de aplicație pe fiecare conexiune țin sigiliul la zi în aceeași tranzacție cu scrierea; funcțiile SQL sunt `DIRECTONLY`.
- **Lanțul jurnalului:** fiecare rând din `jurnal_modificari` are `mac` = HMAC(mac anterior ‖ rând).
- **Pornirea:** baza se verifică înainte de migrări; după migrări se resigilează doar o bază intactă, cu migrarea notată în jurnal. Prima sigilare se face o singură dată, la aplicarea migrării 0021.
- **Incident:** banner permanent cu tabelele afectate și primul rând rupt din jurnal; scrierile sunt refuzate (`audit.editare_blocata`) până la confirmarea cu motiv, notată în jurnal. Incidentul se scrie și în jurnalul de diagnoză, doar cu numele tabelelor.
- **Închiderea lunii și regulile configurabile** scriu acum în tranzacție și în jurnal (înainte scriau în afara tranzacției).
- **Amprenta lunii închise:** calculată la închidere din pontaj, absențe și planificare, păstrată în `stare_luni.amprenta`, afișată în bannerul lunii și în subsolul rapoartelor lunare, cu avertisment dacă datele s-au schimbat după închidere.

Baza de date: migrarea `0021_audit_sigilii` (tabelul `sigilii`, coloanele `jurnal_modificari.mac` și `stare_luni.amprenta`).

Rapoarte: pontaj lunar v3 → v4, pontaj detaliat v2 → v3, fișă individuală v3 → v4, grafic gărzi v1 → v2 (amprenta în subsol pentru lunile închise). Absențele anuale rămân la v2.

## 0.59.0 · 2026-10-28

Pachet complet de rafinare UX/UI (F16-T04..T10, prioritățile 1, 2 și 3):

- **Prioritatea 1 · Eficiență imediată & vizibilitate grilă:**
  * **Bară laterală pliabilă (Compact / Icon-only Sidebar, F16-T04):** Adăugat comutator toggle `[◀] / [▶]` în antetul meniului lateral care restrânge bara la 64px (doar pictogramele cu tooltips native). Eliberează peste 140px lățime direct pentru grilele de pontaj și planificare pe laptopuri. Starea este persistată în `localStorage`.
  * **Control de densitate & Zoom pentru grilă (F16-T05):** Comutator rapid `[S (44px)] [M (52px)] [L (62px)]` în bara de unelte a ecranelor de Pontaj și Planificare, permițând densitate maximă pe ecrane mici sau lizibilitate mărită pe monitoare de rezoluție înaltă.
  * **Tooltips rapide la hover (F16-T06):** La trecerea cursorului peste orice celulă, se afișează instant un rezumat clar: nume angajat, marcă, data, zi de weekend/sărbătoare, cod schimb cu denumire completă și ore, plus eventuale avertismente sau concedii.
- **Prioritatea 2 · Notificări unificate & scurtături universale:**
  * **Sistem unificat Toast Notifications (`Toast.tsx`, F16-T07):** Notificări neintruzive animate în colțul din dreapta-jos pentru confirmări de salvare (`succes`), copiere în clipboard (`info`), avertismente de gardă (`avertisment`) și erori (`eroare`).
  * **Scurtături de tastatură universale (F16-T08):** Închiderea oricărui dialog modal prin tasta `Escape`; confirmare și trimitere rapidă a formularelor prin `Ctrl + Enter`.
- **Prioritatea 3 · Confort vizual & ghid de start:**
  * **Suport Dark Mode / Temă Nocturnă (F16-T09):** Temă întunecată dedicată gărzilor de noapte cu paletă medicală odihnitoare (`[data-theme="dark"]`), comutator `☀️ Zi / 🌙 Noapte` în subsolul barei laterale și persistență a preferinței.
  * **Ghid vizual pe ecrane goale (Onboarding Empty State, F16-T10):** Înlocuirea textelor simple de lipsă date din Pontaj și Planificare cu carduri interactive ghidate, oferind butoane de reîmprospătare sau resetare rapidă a filtrelor.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.58.0 · 2026-10-28

Sistem de diagnoză locală, captare robustă a erorilor și profilare a timpilor de interogare (D043, F16).

- **Pachet backend dedicat (`internal/diagnoza`):**
  * Salvare persistentă a erorilor în fișiere de log rotative locale (`karios/jurnale/pontix-AAAA-LL.log`, max 5 MB).
  * Protecție strictă GDPR: jurnalele conțin exclusiv date tehnice (timestamp, nivel, modul, eroare, stack trace), fără PII sau date medicale.
  * Captare panici Go (`RecupereazaPanic`) cu stivă de execuție completă pentru a preveni crash-ul executabilului.
  * Măsurare a timpilor de execuție a interogărilor (`Cronometreaza`), semnalizare slow queries (>50ms) și calcul statistici (minim, mediu, maxim).
- **Captare globală a erorilor frontend (`ErrorBoundary.tsx`, `main.tsx`):**
  * Componentă React `ErrorBoundary` înfășurând fiecare ecran principal pentru recuperare graceful fără blocarea întregii aplicații.
  * Ascultători globali `window.addEventListener('error')` și `window.addEventListener('unhandledrejection')` trimițând automat excepțiile JS către backend prin `InregistreazaEroareFrontend`.
- **Panou de diagnoză și export în interfață (`Formulare.tsx`):**
  * Integrarea secțiunii de diagnoză în dialogul *Detalii dezvoltare, integritate și cerințe sistem*.
  * Butoane rapide: `📁 Deschide folderul de jurnale` (Explorer) și `📋 Copiază ultimele 100 linii din log`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.57.0 · 2026-10-28

Aliniere licență aplicație și UI cu licența repository-ului (Apache-2.0).

- **Aliniere termen licență cu `LICENSE` din repo (`app_dezvoltare.go`, `app_dezvoltare_test.go`):**
  * Corectat câmpul `Licenta` din raportul backend `DetaliiDezvoltare()` la `Apache License 2.0 (Apache-2.0)` conform fișierului `LICENSE` aflat la rădăcina repoului karios (anterior marcat eronat ca proprietar).
  * Adăugat test specific de validare a licenței în `app_dezvoltare_test.go`.
- **Actualizare UI și internaționalizare (`Formulare.tsx`, `ro.json`, `en.json`):**
  * Afișare insignă de licență `Licență: Apache-2.0` în cardul *Despre Pontix & Autor*.
  * Preluare dinamică a versiunii curente prin `Informatii()` pentru a preveni desincronizările.
  * Actualizare text drepturi autor: `© 2026 Ioan Tuns · Distribuit conform termenilor Apache License 2.0.`
  * Adăugat rând dedicat pentru `Licență Repository: Apache License 2.0 (Apache-2.0)` în fereastra modală tehnică `DialogDetaliiDezvoltare`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.56.0 · 2026-10-28

Fereastră dedicată și raport complet de dezvoltare, integritate și cerințe sistem în Setări (D009, D042).

- **Raport tehnic complet de dezvoltare (`app_dezvoltare.go`, `app_dezvoltare_test.go`):**
  * Noul endpoint backend `DetaliiDezvoltare()` colectează și structurează metadatele tehnice ale aplicației:
    - **Identificare:** Nume repository (`karios`), nume aplicație (`Pontix`), versiune compilată (`0.56.0`), Git commit hash (`versiune.Commit`), arhitectură (`windows/amd64`), versiune Go și Wails (`v2.17.0`).
    - **Validare integritate:** Calea binarului executabil și calculul dinamic al amprentei criptografice **SHA-256 Checksum** a executabilului pentru audit și securitate.
    - **Cerințe minime și recomandate:** OS (Windows 10/11 64-bit), CPU (x64 AVX2), RAM (4 GB fără AI / 8–16 GB cu AI), disc (100 MB de bază / 3.5 GB cu AI), dependințe (WebView2) și conectivitate (100% Offline).
    - **Modul Local AI & LLM:** Nume model (`model.gguf`), dimensiune fișier pe disc, format GGUF 4-bit, motor C++ `llama-server` cu AVX2, fereastră context (8192 tokeni), port local (28491) și izolare completă de rețea (127.0.0.1, 0 KB/s trafic exterior, GDPR).
- **Interfață utilizator & Fereastră Pop-up în Setări (`Formulare.tsx`, `style.css`):**
  * Secțiunea *Despre Pontix & Dezvoltare* din ecranul de Setări afișează insigne rapide (Repo, Versiune, Model AI, Offline) și butonul dedicat `ℹ️ Detalii dezvoltare, integritate și cerințe sistem`.
  * Dialogul modal `DialogDetaliiDezvoltare` prezintă raportul organizat pe categorii și include butonul `📋 Copiază raportul tehnic` pentru asistență sau audit.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.55.0 · 2026-10-27

Extindere fereastră context AI la 8192 tokeni (F14-T04 / D042).

- **Rezolvare eroare depășire context (`exceed_context_size_error`):**
  * Corectat parametrul `-c` al procesului `llama-server.exe` din `internal/ai/proces.go` la **8192 tokeni** (anterior limitat la 2048).
  * Permite încărcarea simultană fără erori 400 a tuturor datelor lunare agregate (personal medical complet, lista tuturor schimburilor, absențele și concediile lunare, turele deja planificate și istoricul conversațional extins, care cumulează peste 2800 tokeni).

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.54.0 · 2026-10-27

Pachet complet de optimizări de performanță, viteză de lansare, productivitate și consum de memorie (Optimizarile 1–5).

- **Optimizarea 1 · Code-Splitting asincron pe Frontend (`React.lazy` + `Suspense` în `App.tsx`):**
  * Divizarea modulelor mari (`Planificare`, `Absente`, `Pontaj`, `Rapoarte`, `AsistentAI`, `Schimburi`, `Sarbatori`) în chunk-uri asincrone încărcate la cerere.
  * Scăderea dimensiunii pachetului JavaScript inițial cu **51%** (de la 515 KB la 252 KB), eliminând complet avertismentul Vite de bundle și accelerând lansarea interfeței în WebView2.
- **Optimizarea 2 & 3 · Navigare rapidă și productivitate în grilă (Excel-style în `Pontaj.tsx` și `Planificare.tsx`):**
  * Adăugat suport pentru tastele `Home` (salt la ziua 1), `End` (salt la ultima zi din lună), `PageUp` / `PageDown` (salt 5 rânduri) și `Space` (deschidere introducere ore).
  * **Introducere directă a orelor din cifre:** tastarea tastelor `0..9` direct pe celulă introduce orele (`@M480` pentru 8h, curățare pentru 0) și avansează automat focusul la ziua următoare, multiplicând viteza de introducere a pontajului.
- **Optimizarea 4 · Eficiență inferență AI Locală & economisire RAM (`internal/ai/proces.go`, `asistent.go`):**
  * Redus fereastra de context de la 4096 la 2048 tokeni, diminuând alocarea KV cache în memorie cu 50% (~500 MB RAM economisiți pe desktop-urile cu resurse limitate).
  * Adăugat procesare de prompt accelerată vectorial (`-b 512`, `-ub 512`).
  * Extins `AnalizeazaInterogareSinteza` cu `sintetizeazaGhidComenzi` pentru comenzi de ajutor/ghid cu răspuns determinist instant (0 ms, 0% CPU).

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.53.0 · 2026-10-27

Motor de reguli configurabile (Rule Engine) pentru restricții medicale și de gardă (Lot 7 / D041 / I061).

- **Motor de evaluare a regulilor medicale (`internal/calcul/motor_reguli.go`):**
  * Evaluare automată a regulilor specifice unităților medicale:
    - `REPAUS_G24`: repaus minim obligatoriu de 24 de ore după o gardă completă de 24h.
    - `PLAFON_GARZI`: limită maximă configurabilă de gărzi/ture per angajat într-o lună (implicit 6).
    - `SUCCESIUNE_G24_Z7`: interdicție succesiune imediată (gardă 24h urmată direct de tură de dimineață a doua zi).
  * Raportare structurată a încălcărilor (`IncalcareRegula`) cu severitate (`AVERTISMENT` / `BLOCANT`), dată, angajat, mesaj explicativ și referință legală/medicală.
- **Persistență și configurare în baza de date (Migrarea 0020 `reguli_configurabile`):**
  * Tabelul `reguli_configurabile` cu cod, denumire, descriere, tip regulă, severitate, coduri schimb vizate, valoare numerică limită și stare de activare.
  * Funcții de citire și salvare (`ListaReguliConfigurabile`, `SalveazaRegulaConfigurabila`, `VerificaReguliPlanificare` în `internal/datebaza/reguli.go` și `app_reguli.go`).
- **Interfață de configurare și alertare în timp real (`Planificare.tsx`, `style.css`):**
  * Alertă vizuală în antetul planificării cu numărul de atenționări medicale și detalii per angajat.
  * Dialog dedicat `DialogConfigurareReguli` pentru ajustarea plafoanelor, schimbarea severității și activarea/dezactivarea regulilor.
  * Teste unitare și de integrare complete în `internal/calcul/reguli_test.go` și `app_reguli_test.go`.

Baza de date: migrarea 0020 adaugă tabelul `reguli_configurabile`. Rapoarte: nicio schimbare.

## 0.52.0 · 2026-10-27

Grafic Gantt anual de concedii și vizualizare panoramică a absențelor (F08-T05 / Lot 6).

- **Vizualizare panoramică și zoom lunar (`PanouGanttAbsente` în `GanttAbsente.tsx`, `Absente.tsx`):**
  * Modul Gantt integrat direct în ecranul de Concedii și absențe, accesibil prin fila dedicată „Cronologie (Gantt)”.
  * Oferă două moduri de afișare: **Panoramă anuală (12 luni)** pentru o perspectivă macro asupra distribuției concediilor și **Lună detaliată (zilele 1..31)** pentru inspecția detaliată a fiecărei zile lucrătoare și de weekend.
  * Grupează automat zilele consecutive în intervale continue per angajat, colorate conform categoriei legale (Odihnă = verde, Medical = portocaliu/roșu, Special = albastru).
- **Detector automat de suprapuneri critice pe departament (Overlap Engine):**
  * Algoritmul identifică în timp real situațiile în care doi sau mai mulți colegi din același departament sunt în concediu în aceleași zile.
  * Afișează caseta de avertisment cu posturile afectate, intervalele critice și numele colegilor implicați.
  * Evidențiază vizual zilele cu deficit de personal prin etichete galbene/portocalii în antetul coloanei și în celule.
- **Carduri sintetice și filtrare pe departament:** Total intervale înregistrate, număr de suprapuneri identificate și filtru rapid `[ ] Doar angajați cu suprapuneri`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.51.0 · 2026-10-27

Închiderea și blocarea lunii la editare (F12-T01).

- **Persistență stare luni calendaristice (Migrarea 0019 `stare_luni`):** Adăugat tabelul `stare_luni` cu câmpurile `an`, `luna`, `stare` (`DESCHISA`, `BLOCATA`, `INCHISA`), `modificat_la`, `modificat_de`, `motiv`.
- **Protecție strictă împotriva modificărilor (`SalveazaSchita`, `SalveazaPlanificare`):** Backend-ul verifică starea lunii înainte de salvare (`EsteLunaInchisa`). Dacă luna este `INCHISA` sau `BLOCATA`, orice încercare de editare este respinsă tranzacțional cu eroarea explicită `pontaj.luna_inchisa`.
- **Semnalizare vizuală și comenzi administrative (`Pontaj.tsx`, `style.css`):**
  * Banner de avertizare stilizat `banner-luna-blocata` cu pictograma 🔒, afișând mesajul „Lună blocată la editare (Transmis la salarii)” și buton rapid de redeschidere.
  * Buton de comutare a stării în antet („🔒 Închide luna (salarii)” / „🔓 Redeschide luna”) protejat prin dialog nativ de confirmare.
  * Dezactivare automată a aplicării celulelor și a butonului de salvare când luna este închisă.
- **Teste unitare și de integrare:** Adăugate `TestStareLunaSiBlocareEditare` în `internal/datebaza/inchidere_test.go` și `TestAppStareLuna` în `app_pontaj_test.go`.

Baza de date: migrarea 0019 adaugă tabelul `stare_luni`. Rapoarte: nicio schimbare.

## 0.50.0 · 2026-10-27

Export Excel nativ (.xlsx via excelize) al foii colective de pontaj și centralizatorului pentru salarizare (F10-T02).

- **Generare registru Excel profesional (`internal/exportexcel/export.go`):** Generator nativ bazat pe biblioteca `excelize/v2`, structurat pe două foi de calcul:
  * **Foaia 1 („Pontaj Lunar”):** Foaia colectivă de prezență completă (Nr. crt., Nume, Marcă, Funcție, zilele 1..31 cu distincție vizuală weekend/lucrătoare, total ore, ore gardă/noapte, ore ZNL, normă, zile CO, zile CM).
  * **Foaia 2 („Centralizator Salarii”):** Centralizatorul sintetizat pentru departamentul de contabilitate și salarizare (Ore normă, Ore efectiv lucrate, Ore gărzi, Ore ZNL, Zile CO, Zile CM) incluzând rândul de **TOTAL GENERAL** calculat cu formule native Excel `=SUM(...)`.
- **Salvare nativă și descărcare web (`app_export.go`, `frontend/wailsjs`):** Metoda `ExportaPontajExcel` integrează dialogul de salvare nativ pe disc (`SaveFileDialog`) sau descarcă automat fișierul generat prin blob Base64.
- **Butoane de export dedicate:** Adăugat butonul „📊 Export Excel (.xlsx)” în bara de acțiuni a ecranului de Pontaj (`Pontaj.tsx`) și în generatorul de rapoarte (`Rapoarte.tsx`).
- **Teste unitare și de integrare:** Adăugate `TestGenereazaRegistruExcel` în `internal/exportexcel/export_test.go` și `TestAppExportaPontajExcel` în `app_export_test.go`.

Baza de date și rapoarte: nicio schimbare de schemă. Format export Excel v1.0.0.

## 0.49.0 · 2026-10-27

Raport comparativ „Planificat vs. Realizat” în ecranul de Planificare și Rapoarte (F11-T02).

- **Vizualizare comparativă interactivă (`PanouComparativPlanificatVsRealizat` în `Planificare.tsx`):** Adăugat modulul comparativ accesibil direct din selectorul de vizualizare al ecranului de Planificare („Planificat vs. Realizat”). Afișează grila lunară completă punând în paralel tura planificată cu pontajul efectiv/ciornă.
- **Identificare automată a discrepanțelor și variațiilor:** Evidențiază stările celulelor prin coduri de culori și insigne dedicate:
  * 🟢 *Identic:* planul a fost respectat exact (`Plan == Realizat`).
  * 🟠 *Modificat:* schimb efectuat diferit de cel planificat inițial (ex: `Z7 ➔ G24`).
  * 🔵 *Neplanificat:* prezență la serviciu fără alocare prealabilă în plan (acoperire ad-hoc).
  * 🔴 *Neefectuat:* tură planificată, dar neînregistrată în pontaj.
  * 🟣 *Concediu / Absență:* zi acoperită de concediu (CO, CM, etc.).
- **Carduri sintetice și filtrare rapidă:** Carduri cu total ore planificate, total ore realizate, variație globală (+/- ore) și număr de angajați cu discrepanțe; filtru binar `[ ] Doar angajați cu discrepanțe` și buton de salt direct în Pontaj.
- **Modul pur de calcul backend (`internal/calcul/comparativ.go`, `internal/datebaza/planificare.go`, `app_planificare.go`):** `CalculeazaPlanificatVsRealizat` și `CitestePlanificatVsRealizat` rulează fără dependențe externe și generează structura de date completă pentru UI și export.
- **Teste unitare și de integrare:** Adăugate `TestCalculeazaPlanificatVsRealizat` în `internal/calcul/comparativ_test.go` și `TestAppCitestePlanificatVsRealizat` în `app_planificare_test.go`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.48.0 · 2026-10-27

Atribuire în masă (bulk) a categoriilor de program pe selecție multiplă de angajați (F07-T15 / I055).

- **Selecție multiplă în tabelul de angajați (`App.tsx`):** Adăugată coloană dedicată cu casete de bifare (checkboxes) per angajat și selector global „Selectează toți angajații vizibili” în antetul tabelului.
- **Bară de comenzi rapide pentru loturi (`bara-actiuni-lot` în `App.tsx`, `style.css`):** La selectarea unuia sau mai multor angajați, apare automat bara contextuală care afișează numărul de angajați selectați, selectorul categoriilor de program definite (cu opțiune de eliminare a categoriei), butonul de aplicare în masă și butonul de deselectare.
- **Operații backend optimizate (`internal/datebaza/date.go`, `app_date.go`):** Implementate metodele `AtribuieCategorieProgramLot` și `SchimbaActivLot` care actualizează atomic în baza de date categoria de program sau starea activ/inactiv pentru o listă de ID-uri, cu înregistrarea completă a modificărilor în jurnalul de audit (`jurnal_modificari`).
- **Traduceri complete:** Adăugate cheile `ui.actiuni_lot`, `ui.selectati_numar`, `ui.atribuie_categorie_lot`, `ui.fara_categorie_program`, `ui.aplica_categorie_lot`, `ui.deselecteaza_tot`, `ui.selecteaza_toti`, `ui.categorie_lot_aplicata` în română și engleză.
- **Teste unitare și de integrare:** Adăugat `TestAtribuieCategorieProgramLot` în `internal/datebaza/date_test.go` și `TestAppAtribuieCategorieProgramLot` în `app_date_test.go`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.47.0 · 2026-10-27

Comenzi analitice și interogări de sinteză în Asistentul AI (F14-T07, D042).

- **Sinteză analitică factuală instantanee (`internal/ai/asistent.go`):** Implementat modulul `AnalizeazaInterogareSinteza` care răspunde determinist, precis și fără latență la interogările frecvente despre personal:
  * *Disponibilitate weekend:* „Cine este liber în acest weekend?” — analizează toate schimburile și concediile din sâmbătă și duminică, raportând lista completă a angajaților liberi și a celor la serviciu sau absenți.
  * *Concedii și absențe:* „Cine este în concediu luna asta?” — sintetizează intervalele de absențe pe tipuri (CO, CM, CF etc.) și totalul de zile pe secție.
  * *Verificare acoperire posturi:* „Verifică dacă lipsesc posturi obligatorii” — analizează fiecare zi de weekend și semnalează zilele neacoperite sau cu deficit de personal.
  * *Statistici gărzi și ore:* „Statistici gărzi și ore” — calculează tabelul centralizat pe angajați cu numărul de gărzi, total ture și ore planificate.
- **Context extins cu absențe și planificare (`app_ai.go`):** `ContextPontix` colectează acum automat absențele active și turele planificate/completate în ciornă pentru luna curentă, oferindu-le modelului AI în promptul de sistem.
- **Teste unitare de integrare:** Adăugate `TestAnalizeazaInterogareSinteza` în `internal/ai/asistent_test.go` și `TestAppInterogheazaAIComenziAnalitice` în `app_ai_test.go`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.46.0 · 2026-10-27

Execuție automată la confirmări în chat („da” / „confirm”), suport pentru weekenduri alternate și unelte de planificare recurente (F14-T04, D042).

- **Execuție directă la răspuns afirmativ în chat (`AsistentAI.tsx`):** Când utilizatorul confirmă o propunere tastând „da”, „confirm”, „aplică”, „ok” etc. în caseta de mesaj, aplicația detectează automat intenția și execută direct salvarea (`AplicaPropunereAI`), eliminând confuzia în care modelul confirma doar verbal fără a opera modificarea în baza de date.
- **Suport complet pentru weekenduri alternate și recurență (`app_ai.go`, `internal/ai/asistent.go`):** Adăugat filtrul `tot_al_doilea_weekend` și `weekend` în `planifica_interval` și unealta dedicată `planifica_zile`. Comanda utilizatorului (ex: „Dumitru Andrei pe octombrie lucrează tot al doilea weekend pe zi, schimb de 12 ore, începând cu data de 3”) este mapată exact pe schimbul `Z12` și pe datele corecte ale weekendurilor (3–4 oct, 17–18 oct, 31 oct), fără a mai forța eronat zile lucrătoare L–V.
- **Rezolvare revizii existente în Planificare (`app_ai.go`):** La sincronizarea în planificare pentru turele de gărzi/weekend (Z12), sunt prelucrate reviziile și identificatorii pozițiilor existente pentru a preveni erorile de concurență.
- **Parser defensiv de propuneri (`internal/ai/asistent.go`):** În cazurile în care modelul generează accidental descrierea acțiunii sub formă de text în loc de apel nativ de unelte, backend-ul sintetizează automat propunerea și afișează cardul interactiv cu butoane de confirmare.
- **Teste unitare de regresie:** Adăugate `TestAplicaPropunereAITotAlDoileaWeekend` în `app_ai_test.go` și actualizate testele în `asistent_test.go`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.45.0 · 2026-10-27

Persistența conversației pe durata sesiunii, context îmbogățit pentru comenzi consecutive și buton de salt direct în pontaj (F14-T04, D042).

- **Persistență în sesiune (`AsistentAI.tsx`):** Conversația nu se mai pierde la navigarea între ecrane (de ex. la comutarea între Asistent AI, Pontaj și Planificare). Mesajele sunt salvate în `sessionStorage` și restaurate automat la revenire.
- **Buton „Conversație nouă” (`ai.conversatie_noua`):** Adăugat buton în bara de stare a asistentului pentru resetarea rapidă a fluxului conversațional când utilizatorul dorește să înceapă un dialog curat.
- **Memorie contextuală multi-turn (`istoricAPI`):** Răspunsurile asistentului trimise în istoricul către LLM includ acum explicit detaliile acțiunilor propuse anterior și starea acestora (`[confirmată și aplicată în pontaj]` vs `[în așteptare]`). Comenzi ulterioare precum „și Popescu la fel”, „de fapt pune-l pe noapte” sau „pe 15 e liber” păstrează contextul exact al persoanei, lunii și programului discutat.
- **Buton direct „Deschide în Pontaj” (`ai.deschide_in_pontaj`):** Pe cardurile de acțiune aplicate cu succes apare opțiunea de salt direct în ecranul de Pontaj, prefiltrat exact pe anul și luna vizată de comandă.
- **Traduceri:** Actualizate cataloagele de limbă `ro.json` și `en.json` pentru noile acțiuni și etichete.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.44.0 · 2026-10-27

Aplicare directă în pontaj și interpretare automată a cererilor pentru întreaga lună (F14-T04, D042).

- **Interpretare comenzi lună întreagă (`internal/ai/asistent.go`):** Când utilizatorul menționează o lună (ex: „în octombrie”, „in octombrie”) sau o perioadă („toată luna”, „tata luna”) fără a specifica o anumită zi, asistentul generează automat `planifica_interval` pentru prima și ultima zi a lunii cu `doar_lucratoare: true`. Adăugată funcția `esteCerereToataLuna` pentru corecție defensivă pe server și test de regresie `TestEsteCerereToataLuna`.
- **Metodă unificată de aplicare acțiuni AI (`AplicaPropunereAI` în `app_ai.go`):** Eliminat eroarea `Selecția pentru planificare nu este validă` cauzată de apelarea `SalveazaPlanificare` pentru schimburi normale de zi (cum ar fi `Z7`). Noul endpoint salvează direct în pontaj (`SalveazaSchita` cu revizii la zi) și sincronizează automat tabelul de planificare exclusiv pentru turele din categoriile gărzi (`GARZI`) sau după-amiază (`DUPA_AMIAZA`).
- **Actualizare client frontend (`AsistentAI.tsx`):** Simplificat handlerul `aplicaPropunere` pentru a delega execuția completă către backend-ul Go prin `AplicaPropunereAI`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.43.0 · 2026-10-27

Reguli stricte de mapare intervale orare și corecție automată schimb cerut (F14-T04, D042).

- **Mapare strictă intervale orare în promptul de sistem (`internal/ai/asistent.go`):** Instruit modelul cu reguli fără echivoc: „7 la 14” / „07:00–14:00” corespunde exclusiv schimbului `Z7` (Zi 07–14), iar `S7` este rezervat pentru „15 la 22” (Seară). Inclus intervalul orar explicit în lista fiecărui schimb din prompt.
- **Protecție defensivă de corecție automată (`corecteazaSchimbDacaEroare`):** Dacă utilizatorul a cerut „7 la 14” sau „dimineață”, dar modelul propusese eronat `S7`, backend-ul corectează automat codul de schimb la `Z7`. De asemenea, pentru cererile „toată luna” pe programe de zi, setează implicit `doar_lucratoare: true` pentru a proteja weekendurile.
- **Test unitar de regresie:** Adăugat `TestCorecteazaSchimbDacaEroare` în `internal/ai/asistent_test.go`.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.42.0 · 2026-10-27

Afișare completă nume angajat și detalii program în propuneri AI și rezoluție robustă marcă (F14-T04, D042).

- **Descriere extinsă a propunerilor (`internal/ai/asistent.go`):** Lângă codul schimbului este afișată denumirea completă, intervalul orar și durata (ex: `Z7 – Zi 07–14 (07:00–14:00, 7 ore)`), iar lângă marcă este inclus numele complet al angajatului (ex: `Dobre Rareș (marca DEMO-008)`).
- **Canonicizare și rezoluție tolerantă a mărcii:** Adăugată normalizare alfanumerică împotriva diferențelor de formatare (ex: `DEM008` vs `DEMO-008`), diacriticelor sau căutărilor după nume; backend-ul canonicizează automat parametrul mărcii la valoarea oficială din baza de date, iar frontend-ul (`gasesteAngajat`) găsește angajatul garantat la aplicare.
- **Insigne vizuale structurate în UI (`AsistentAI.tsx`):** Cardul de propunere afișează etichete dedicate pentru Program, Angajat și Dată/Interval.
- **Formatare interval orar schimburi (`app_ai.go`):** Calculat intervalul orar (start–sfârșit) din minute pentru contextul asistentului.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.41.0 · 2026-10-27

Optimizare parametri server AI local (`-np 1`, fire CPU dinamice) și extindere timeout HTTP (F14-T04, D042).

- **Optimizare slot unic (`-np 1`) în `internal/ai/proces.go`:** Eliminat supraîncărcarea memoriei și a nucleelor prin forțarea unui singur slot activ (`-np 1`), specific scenariului de asistent desktop monopost. Timpul de evaluare inițială a promptului pe CPU scade dramatic (de la >24s la ~7.5s).
- **Alocare dinamică de nuclee CPU:** Serverul folosește acum `runtime.NumCPU() - 2` (minimum 2 fire) în loc de o valoare fixă de 4 nuclee, lăsând resurse pentru sistem și interfață și utilizând la maximum arhitectura procesorului gazdă.
- **Creștere timeout HTTP client la 180 secunde (`internal/ai/asistent.go`):** Previne eroarea `context deadline exceeded` în timpul generărilor lungi sau pe procesoare mai lente la prima interogare.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.40.0 · 2026-10-27

Planificare pe perioade/luni întregi pentru Asistentul AI și extindere unelte structurate (F14-T04, D042).

- **Unelte noi `planifica_interval` și `sterge_interval` (`internal/ai/asistent.go`):** Suport pentru atribuirea sau ștergerea turelor pe un interval de date, o săptămână sau toată luna dintr-o singură instrucțiune rapidă (1–2 secunde), eliminând generarea a 30 de apeluri repetitive și riscul de expirare timeout.
- **Suport opțional zile lucrătoare:** Parametrul `doar_lucratoare` (Luni–Vineri) pentru programările pe intervale fără weekend.
- **Creștere timeout comunicare HTTP:** Mărit timeout-ul clientului la 120 secunde pentru a permite generări ample pe procesor.
- **Potrivire flexibilă angajat:** Recunoaștere în interfață după marcă, `Nume Prenume`, `Prenume Nume` sau nume de familie.
- **Mesaje tehnice explicite:** Înlocuit mesajul generic de eroare cu detalii clare de diagnostic (`ai.eroare_generare`).

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.39.0 · 2026-10-27

Corectare cale rezoluție folder local AI portabil și extindere timeout inițializare (F14-T02, F14-T05, D042).

- **Corectare rezoluție cale `karios/ai/` (`app_ai.go`):** Eliminat dublarea numelui de folder `karios` (`karios/karios/ai`) care apărea când `app.radacina` conținea deja calea folderului de date portabil; asistentul scanează acum corect folderul `karios/ai/` de lângă executabil.
- **Extindere timeout inițializare motor AI:** Crescut timeout-ul de pornire la 45 secunde în `internal/ai/proces.go` pentru a permite încărcarea fluentă în memoria RAM a modelelor LLM mari (GGUF ~2 GB+).
- **Test de regresie:** Adăugată verificare automată în `app_ai_test.go` împotriva dublării segmentului de cale portabilă.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.38.0 · 2026-10-27

Compilare nativă llama.cpp pe Windows cu toolchain portabil și integrare în procesul de Deploy (F14-T06, D042).

- **Toolchain portabil C/C++ și CMake (`scripts/pregateste-toolchain.ps1`):** Configurare complet autonomă a unui mediu de compilare în `tools/toolchain/` (CMake 3.31.5 și MinGW-w64 GCC 16.2.0 via `w64devkit`) fără modificări în sistemul de operare Windows, fără drepturi de administrator și ignorat complet de Git.
- **Compilare nativă a motorului AI (`scripts/compileaza-llama.ps1`):** Clonare și compilare a executabilului C++ `llama-server.exe` optimizat cu AVX2 direct pe Windows folosind toate nucleele procesorului (cu suport incremental de ~7 secunde la recompilare), generând binarul verificat independent și plasat în `karios/ai/` și `.cache-ai/`.
- **Integrare în procesul de Deploy (`build.ps1 -CuCompilareAI`):** Extinderea scriptului oficial de build cu parametrul `-CuCompilareAI` pentru a asigura compilarea proaspătă a binarului `llama-server.exe` și includerea automată a acestuia în structura livrată la fiecare lansare de versiune.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.37.0 · 2026-10-26

Pachet portabil Pontix AI Ready și pornire automată la lansare (F14-T05, D042).

- **Script automatizat de împachetare AI Ready (`scripts/creaza-pachet-ai.ps1`):** Script dedicat care produce o distribuție portabilă completă (sau arhivă `.zip`) conținând executabilul `pontix.exe`, directorul `karios/ai/` pre-populat cu binarul motorului C++ `llama-server.exe` și modelul LLM cuantizat (`qwen2.5-3b-instruct-q4_k_m.gguf`), script de descărcare la cerere, ghid detaliat `README-AI-READY.txt` și fișierul de integritate `build.json` cu sume de control SHA256 calculate pentru fiecare componentă.
- **Pornire automată supravegheată a motorului AI:** Funcționalitate nativă de pornire automată în fundal la lansarea aplicației dacă fișierele sunt prezente în `karios/ai/`, pentru a elimina timpul de încărcare atunci când utilizatorul deschide prima dată ecranul de asistență.
- **Comutator în interfață:** Adăugată bifa „Pornire automată la lansare” în bara de stare a ecranului Asistent AI, cu persistență portabilă locală (fișier `fara-autostart` în `karios/ai/`).
- **Suport extins teste unitare:** Adăugat `app_ai_test.go` pentru verificarea activării/dezactivării pornirii automate și detectării corecte a stării pe folder gol.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.36.0 · 2026-10-26

Integrare arhitectură și asistent AI local, offline și portabil (F14-T02, F14-T04, D042).

- **Arhitectură motor AI local portabil (`internal/ai`):** Modul nativ Go pentru supravegherea procesului `llama-server` în fundal fără consolă vizibilă (`HideWindow: true` pe Windows), cu verificare automată a integrității fișierelor și a stării de sănătate (`/health`) pe port local dinamic (`127.0.0.1:28491`).
- **Director portabil dedicat:** Adăugat folderul `karios/ai/` în structura portabilă a aplicației pentru binarul `llama-server.exe` și modelele cuantizate `.gguf`.
- **Client de inferență cu apelare structurată de unelte (Function Calling):** Client HTTP OpenAI-compatible local cu scheme JSON stricte pentru acțiunile de pontaj și planificare (`planifica_tura`, `sterge_tura`). Modelul nu rulează calcule matematice sau SQL direct, ci emite propuneri de acțiuni structurate.
- **Ecran dedicat Asistent AI (`AsistentAI.tsx`):** Punct de control în meniul principal cu indicator de stare a motorului local (activ / gata / fișiere lipsă cu ghidare cale folder), flux interactiv de conversație, cipuri de comenzi rapide și carduri de confirmare a acțiunilor propuse cu buton `Confirmă și aplică` și verificare atomică a salvarii în planificare.
- **Garantare GDPR și siguranță:** Oprire automată a procesului copil la părăsirea aplicației Pontix; nicio transmisie de date prin rețea externă sau cloud.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.35.0 · 2026-10-26

Import de date pentru angajați și sărbători legale din Excel (.xlsx) și CSV (.csv), șabloane demonstrative descărcabile, semnătura dezvoltatorului în interfață și definirea deciziei pentru asistent AI exclusiv local (F02-T13, F03-T08, F14, D042).

- **Import angajați din Excel (.xlsx) și CSV (.csv):** Analiză și validare automată la încărcarea fișierului, recunoaștere inteligentă antete indiferent de diacritice sau majuscule, detecție automată delimitator (virgulă / punct și virgulă), creare automată pentru secțiile și departamentele noi, opțiune de actualizare angajați existenți (după marcă) și execuție tranzacțională atomică în SQLite.
- **Import sărbători legale din Excel (.xlsx) și CSV (.csv):** Suport pentru formate multiple de dată (`AAAA-LL-ZZ`, `ZZ.LL.AAAA`, `ZZ/LL/AAAA`), identificare denumiri RO/EN, detectare duplicate și opțiune de actualizare/înlocuire sărbători existente.
- **Șabloane demonstrative descărcabile:** Generare directă din interfață a șabloanelor structurate pentru angajați și sărbători în formate Excel (.xlsx via excelize) și CSV (cu UTF-8 BOM pentru afișarea corectă a diacriticelor românești în Excel).
- **Semnătura dezvoltatorului:** Integrarea semnăturii oficiale a autorului („Dezvoltat de Ioan Tuns”) în subsolul bazei de navigație a aplicației și în ecranul Setări într-o secțiune dedicată Despre Pontix & Autor (cu copyright © 2026 Ioan Tuns).
- **Asistent AI exclusiv local:** Stabilirea deciziei arhitecturale D042 care impune operarea exclusivă pe modele LLM locale (fără trimitere de date în cloud sau dependențe de API-uri externe), protejând datele salariale și confidențialitatea personalului.

Baza de date și rapoarte: nicio schimbare de schemă.

## 0.34.0 · 2026-10-26

Raport dedicat pentru graficul lunar de gărzi și schimburi, calendare sărbători legale 2025 și 2027, eliminare avertisment zile consecutive în mod standard și propunere arhitectură motor reguli configurabile (F11-T07, D040, D041, I061, I062).

- **Raport oficial Grafic lunar de gărzi și schimburi (`raport-garzi`):** Sursă nouă de raportare `GRAFIC_GARZI`, generată direct din tabela de planificare, afișând matricea lunară cu personalul pe rânduri, zilele 1–31 cu evidențierea weekendurilor și a sărbătorilor, totalurile lunare de gărzi și ore, rândul de sinteză `TOTAL / ZI` pentru verificarea acoperirii zilnice și blocul oficial de semnături (Întocmit, Avizat, Aprobat).
- **Acces direct din Planificare:** Buton dedicat `Grafic gărzi (PDF)` în bara de acțiuni a ecranului Planificare, cu previzualizare instantanee în dialog modal maximizabil, opțiuni de descărcare și tipărire directă.
- **Sărbători legale 2025 și 2027:** Adăugarea calendarelor oficiale de sărbători legale pentru România pentru anii 2025 (Paști 20 aprilie, Rusalii 8–9 iunie) și 2027 (Paști 2 mai, Rusalii 20–21 iunie), cu import complet automat în interfață.
- **Ajustare avertismente grilă:** Eliminarea avertismentului neadecvat pentru ture în zile consecutive normale (`CONSECUTIV`), păstrând exclusiv rotația de sărbători (Crăciun/Revelion/Paști).

Baza de date: migrarea **0018** (`definitii_rapoarte` pentru `raport-garzi`).
Rapoarte: raport nou `raport-garzi` (sursa `GRAFIC_GARZI`, versiune **1**).

## 0.33.0 · 2026-10-26

Îmbunătățire ergonomie modal sumar planificare angajat: opțiune de mărire/maximizare dialog și bară de derulare completă pe verticală și orizontală (F11-T04).

- **Opțiune de mărire fereastră (maximizare):** Buton dedicat `🗖` / `🗗` în antetul dialogului pentru comutarea între dimensiunea standard mărită (`min(880px, calc(100vw - 36px))`) și modul pe tot ecranul (`calc(100vw - 32px)` / `calc(100vh - 32px)`).
- **Bară de derulare și lățime optimizată:** Zona tabelară a sumarului lunar include acum bară de derulare verticală dedicată cu `scrollbar-gutter: stable` și stil Pontix, permițând parcurgerea facilă a tuturor celor 31 de zile, iar lățimea generoasă elimină tăierea coloanei de acțiuni.

Baza de date și rapoarte: nicio schimbare.

## 0.32.0 · 2026-10-26

Sumar editabil pe angajat, raport complet de acoperire pentru weekenduri și sărbători legale, verificare echitate rotație și marcaje pentru ture consecutive (F11-T04, F11-T05, F11-T06, I058–I060, D039).

- **Sumar editabil pe angajat:** Vizualizare dedicată și fereastră modală accesibilă din grilă printr-un clic pe rândul angajatului. Afișează totalul de ture și ore planificate în lună, lista zilelor cu evidențierea weekendurilor și a sărbătorilor, selector rapid pentru schimburi permise sau ștergere tură și avertismentele specifice angajatului.
- **Raport complet de acoperire weekend și sărbători legale:** Vizualizare tabelară dedicată care verifică prezența personalului pe posturile obligatorii (`G24`, `Z12`, `N12`, `IN9`) pentru fiecare sâmbătă, duminică și sărbătoare legală. Afișează badge-uri cu angajații alocați sau semnalizare vizuală `Lipsă / Neacoperit`, carduri sintetice cu procentul de acoperire și filtru pentru afișarea exclusivă a zilelor cu deficit.
- **Echitate sărbători și marcaj ture consecutive:** Verificarea automată a istoricului de sărbători majore (Crăciun `N-1` -> Revelion `N`, Crăciun `N` -> Paști `N+1` pe același schimb) și a turelor/gărzilor planificate în zile consecutive (fără repaus după `G24` sau ture succesive). Semnalizare vizuală neblocantă în grilă (`celula-atentie`), tooltip-uri explicative și detalii în panoul de inspecție poziție.

Baza de date și rapoarte: nicio schimbare de schemă sau formule.

## 0.31.0 · 2026-10-26

Grilă lunară tabelară de planificare cu atribuire directă pe celule (pereche zi – tip program) pentru gărzi, intrare weekend și după-amiază, cu publicare unificată în schița de pontaj (F11-T03, I057, D038).

- **Grilă lunară interactivă:** Angajații pe rânduri și zilele 1–31 pe coloane, cu evidențierea weekendurilor, a sărbătorilor legale și a absențelor active (concedii de odihnă, medicale, evenimente).
- **Atribuire rapidă pereche zi – program:** Paletă rapidă cu cipurile schimburilor speciale (`S7`, `IN9`, `G24`, `Z12`, `N12`) și acțiune de ștergere, cu aplicare dintr-un clic pe celula activă sau pe o selecție multiplă (drag / Shift+Click).
- **Verificare categorii și stivă Undo / Redo:** Validare împotriva atribuirii schimburilor nepermise de categoria activă a angajatului; butoane și scurtături `Ctrl+Z` / `Ctrl+Y` pentru anularea/refacerea pașilor locali.
- **Publicare unificată:** Previzualizarea și publicarea în pontaj permit selectarea tuturor secțiunilor simultan (`TOATE`), cu verificarea tranzacțională a absențelor, schițelor existente și pontajelor efective.
- **Comutator Grilă / Listă:** Păstrarea tabelului detaliat de audit pentru consultarea istoricului și a datelor de publicare.

Baza de date și rapoarte: nicio schimbare de schemă sau formule.

## 0.30.0 · 2026-10-26

Secțiune operațională **Planificare**, cu **Gărzi** și **După-amiază**, salvare independentă și publicare explicită în schița de pontaj (F11-T01, I056, D037).

- Planificarea pe angajat, program și date se salvează fără a modifica pontajul. Listele au filtre, selecție și stări de publicare; retragerea din plan păstrează pontajul deja publicat.
- Publicarea se previzualizează, completează implicit numai celulele goale și permite înlocuirea explicită a schițelor cu afișarea înainte/după. Sunt reverificate tranzacțional absențele, programele permise, intervalele, pontajele efective și reviziile. Republicarea valorilor identice nu dublează orele.
- Precompletarea automată din Pontaj direcționează schimburile de după-amiază și gărzile către Planificare; corecțiile manuale în pontaj rămân disponibile.

Baza de date: migrarea **0017**, tabel de planificare și metadate de publicare. Formulele și versiunile rapoartelor sunt păstrate.

## 0.29.0 · 2026-10-26

Categorii de program pe angajat și precompletare explicită, conform confirmării utilizatorului (F07-T14, I053, D036).

- **Precompletează L–V** revine în bara vizibilă din Pontaj. Butonul **După categorii** propune programul implicit individual, cu previzualizare și aplicare în schiță, Undo/Redo și salvare ulterioară.
- **Setări → Categorii de program** permite creare/editare, denumiri RO/EN, schimburi permise, implicit opțional și dezactivare. Șase categorii inițiale sunt editabile; niciun angajat existent nu este atribuit automat. Combinațiile nu generează rotații; gărzile se aleg explicit.
- Fișa angajatului are categorie de program, vizibilă și în lista angajaților. Pontajul se filtrează după categorie, inclusiv angajații fără categorie.
- Precompletarea pe categorii respectă calendarul, toate absențele active, datele existente, norma individuală, suprapunerile și trecerile dintre luni. Programele incomplete, inactive sau cele speciale fără alegere explicită sunt omise cu explicații. Programele cunoscute din afara categoriei active sunt refuzate la aplicare și salvare; orele manuale rămân introduse explicit fără deducerea unui schimb.

Baza de date: migrarea **0016**, categorii, legături stabile la schimburi și atribuire opțională pe angajat. Rapoarte și reguli de plată: nicio schimbare. Corecția încărcării calendarului în adaptorul vechi de totaluri este urmărită separat în F04-T09.

## 0.28.0 · 2026-10-26

Primul lot UI după audit (F05-T19, I045–I048, D035): comenzile detaliate din Pontaj sunt grupate în Filtre, Completare și Ajutor, lăsând grilei mai mult spațiu. Căutarea și salvarea rămân la vedere; legenda textuală și contrastul orelor/valorilor necalculate clarifică stările.

- T diferit de Ore norm nu mai produce un avertisment de normă contractuală. Ore norm păstrează sensul de ore lucrate în zile normale; explicația este disponibilă în Ajutor și la total. Avertismentele existente după G24 și peste 24 h se păstrează.
- Navigarea din și în Pontaj păstrează în sesiune perioada, filtrele, selecția și derularea. Modificările nesalvate cer în continuare confirmare; la revenire datele se recitesc din bază.
- Dialogurile lungi au antet și acțiuni permanente, cu derulare în corp. Bara de salvare din Configurare pontaj rămâne accesibilă la derularea paginii.

Rapoarte, formule, biblioteci și baza de date: nicio schimbare. Lista următoarelor loturi este în [lista UI](lista-lucru-ui.md).

## 0.27.0 · 2026-10-26

Gruparea ecranelor de configurare în meniul Setări, redenumirea „Sărbători” în „Zile libere legale” și temă vizuală distinctă (albastru-ardezie) pentru zona Setări (F05-T17, I044, D034).

- Meniul principal din bara laterală păstrează cele 5 destinații principale: **Angajați**, **Pontaj**, **Concedii și absențe**, **Rapoarte** și **Setări**.
- Ecranele de configurare (**General**, **Zile libere legale** — fostul **Sărbători** —, **Schimburi**, **Configurare pontaj** și **Definire rapoarte**) sunt grupate în interiorul **Setări**, accesibile atât prin taburile din partea de sus a ecranului Setări, cât și prin sub-meniul afișat sub **Setări** în bara laterală.
- Meniul și zona **Setări** folosesc o paletă de culori proprie (albastru-ardezie instituțional `#2b4c7e` și fundal activ `#eef3f9`), distinctă de verdele operațional al aplicației.

Rapoarte: nicio schimbare. Baza de date: nicio schimbare.

## 0.26.0 · 2026-10-26

Fluxuri UX în grila de pontaj și în modulul de concedii și absențe: căutare după nume/marcă și navigare rapidă între luni, selecție dreptunghiulară cu mouse-ul (Click & Drag) și autoderulare la navigarea din tastatură, precompletare L–V doar pentru rândurile selectate, filtrare după angajați cu avertismente/modificări/înlocuitori, căutare și selectarea dintr-un clic a întregii cereri în Concedii și absențe, plus salt contextual din Pontaj către Concedii și absențe (F08-T09, I038–I043, D033).

- În **Pontaj**, s-au adăugat butoanele `◀` / `▶` pentru trecerea la luna anterioară/următoare (inclusiv trecerea automată între decembrie și ianuarie), câmpul de căutare după nume, prenume sau marcă și filtrul rapid de rânduri (`Toți angajații`, `Doar cu avertismente`, `Doar cu modificări`, `Doar înlocuitori`).
- Selecția dreptunghiulară în grilă poate fi extinsă și prin apăsarea și tragerea cu mouse-ul peste celule (`Click & Drag`), iar navigarea din săgeți derulează automat tabelul (`scrollIntoView`) pentru a menține celula activă vizibilă.
- Butonul **Precompletează selecția L–V** completează cu programul ales doar celulele goale din zilele lucrătoare (luni–vineri, fără sărbători legale și fără absențe active) pentru angajații care au cel puțin o celulă în selecția curentă.
- Când celula activă din **Pontaj** este o zi de absență, bara de inspecție oferă butonul **Deschide în Concedii și absențe**, care comută ecranul direct filtrat pe anul și angajatul respectiv.
- În **Concedii și absențe**, s-au adăugat câmpul de căutare liberă (după nume, marcă, număr de înregistrare, înlocuitor, cod, aprobator sau document) și butonul **Selectează cererea**, care bifează dintr-un clic toate zilele aparținând aceleiași cereri multizi (după documentul partajat sau numărul de înregistrare).

Rapoarte: nicio schimbare. Baza de date: nicio schimbare.

## 0.25.0 · 2026-10-26

Îmbunătățiri de ergonomie în interfață: bară de inspecție a celulei active și paletă rapidă de schimburi în grila de pontaj, stivă Undo/Redo, coloană `T` fixă la derulare, comutator de secțiuni în Concedii și absențe, cipuri de selecție rapidă în Configurare pontaj, selector de limbă RO/EN și pictograme distincte în meniu (F08-T08, I032–I037, D032).

- În **Pontaj**, panoul de inspecție de deasupra grilei afișează instantaneu (inclusiv la navigarea din tastatură) detaliile celulei active: angajatul, data și sărbătoarea legală, codul programului / orele manuale / absența, orele calculate, persoanele înlocuite și avertismentele; paleta rapidă permite aplicarea oricărui schimb activ sau golirea selecției dintr-un clic.
- Modificările din schița grilei de pontaj dispun de stivă **Undo / Redo** (`Ctrl+Z` / `Ctrl+Y` și butoanele **Anulează pasul** / **Refă pasul**).
- Coloana de total **`T`** este înghețată în dreapta tabelului (`sticky right`), rămânând vizibilă la derularea orizontală.
- În **Concedii și absențe**, butoanele de secțiune permit izolarea rapidă a sumarului anual, a zilelor înregistrate sau a tipurilor de absență; în **Configurare pontaj**, cipurile de sub selector permit bifarea/debifarea codurilor dintr-un clic; în bara laterală există selectorul de limbă **RO / EN** și pictograme SVG distincte pentru toate cele 9 ecrane.

Rapoarte: nicio schimbare. Baza de date: nicio schimbare.

## 0.24.0 · 2026-10-26

Precompletarea lunii pentru zilele lucrătoare luni–vineri în grila de pontaj, avertismente vizuale neblocante în grilă și copie de siguranță criptată la cerere din Setări (F08-T07, F05-T16, I029–I031, D031).

- Butonul **Precompletează L–V** din bara de unelte a grilei de pontaj completează cu programul ales toate celulele goale din zilele de luni până vineri care nu sunt sărbători legale și nu au absențe active, pentru rândurile vizibile.
- Celulele din grilă sunt evidențiate vizual (`celula-avertisment`) și explică la survol când există program/ore lucrate în ziua imediat următoare unei gărzi de 24 h (`G24`) sau când totalul zilei depășește 24 h; coloana de total `T` este evidențiată (`total-diferenta`) când diferă de `ore_norma`.
- În **Setări**, secțiunea **Copie de siguranță** permite salvarea imediată a unei copii criptate a bazei de date (`VACUUM INTO`) în `karios/backup/pontix-manual-<data-ora>-v<versiune>.db`, fără a afecta rotația celor 7 copii automate făcute înaintea migrărilor.

Rapoarte: nicio schimbare. Baza de date: nicio schimbare.

## 0.23.0 · 2026-10-26

Verificarea suprapunerilor pentru persoana înlocuitoare la concedii, evidențierea înlocuirii în pontaj, diferența aritmetică în sumarul anual și previzualizarea modală a cererii semnate (F07-T13, I025–I028, D030).

- Serviciul Go și formularul de concedii refuză desemnarea unui înlocuitor care are el însuși o absență activă nerespinsă în aceeași zi (`absente.inlocuitor_absent`), precum și înregistrarea unei absențe active nerespinse pentru un angajat care este deja desemnat înlocuitor în aceeași zi (`absente.titular_este_inlocuitor`); dubla înlocuire (când aceeași persoană înlocuiește mai mulți colegi în aceeași zi) este semnalată prin avertisment în formular și etichetă în tabel.
- Grila lunară de pontaj evidențiază vizual (`celula-inlocuitor`) și în explicația celulei (`title`) zilele în care un angajat înlocuiește colegi aflați în concediu.
- Tabelul de sumar din Concedii și absențe afișează, sub dreptul anual contractual, diferența aritmetică rămasă (`Rămas: X zile`, fără reportări istorice) și numărul de zile de odihnă aflate încă în așteptare.
- Cererile de concediu atașate (PDF, PNG, JPEG) pot fi previzualizate direct într-o fereastră modală, pe lângă descărcarea explicită.

Rapoarte: nicio schimbare. Baza de date: nicio schimbare.

## 0.22.2 · 2026-10-26

F02-T12: fiecare invocare build.ps1 produce un director unic `build/bin/pontix-<versiune>-<data-ora>-<identificator>/`, cu executabil și manifest build.json (versiune, commit, SHA256). Temporarele Wails sunt izolate per invocare și curățate numai din directorul nou. Buildurile existente și datele portabile nu se suprascriu. Dependențele fixate nu sunt modificate de build. Fără schimbări în rapoarte, calcule sau baza de date.

## 0.22.1 · 2026-10-26

F09-T09: corectează unitățile coloanelor finale din PDF-ul detaliat, folosind paranteze în locul separatorului de rând pe care motorul PDF îl afișa ca pătrat. Pontaj detaliat **v2**, cu referința v1 păstrată. Centralizator v3, fișă individuală v3, absențe v2 și migrarea 0015 rămân la aceleași versiuni.

## 0.22.0 · 2026-10-26

F09-T09: pontajul afișează ore, păstrând codurile programelor pentru calcule. Alegere explicită a programului sau ore manuale; previzualizarea modificărilor folosește același motor ca PDF-urile. Configurare separată pentru atribuirea turelor peste miezul nopții și cele 12 coloane finale. T este totalul, ZNL include weekenduri și zile libere 00–24, Ore norm include luni–vineri fără zile libere. CO/CM numără absențele active. Celelalte reguli sunt inițial nedefinite; nu se inventează sporuri.

Două rapoarte lunare distincte: centralizatorul existent și pontajul detaliat cu totaluri la final. Antetele, subsolurile și imaginile existente se păstrează. Rapoarte: centralizator **v3**, fișă individuală **v3**, detaliat **v1**, absențe **v2**. Baza: **0015_pontaj_ore_reguli**, configurare cu revizii și jurnal, conversia coloanei cod în ore numai la fișele individuale.

## 0.21.0 · 2026-10-25

F09-T08: imagini PNG/JPEG în antetul și subsolul fiecărui raport, în cele trei zone, alături de texte. Import local, înlocuire/eliminare, înălțime configurabilă cu proporții păstrate și validare a conținutului, dimensiunilor și volumului. Imaginile sunt păstrate în baza criptată și se repetă pe paginile documentului.

Paginarea automată permite rapoarte de una sau mai multe pagini; previzualizarea afișează numărul real de pagini al PDF-ului. Include cererile de concediu și aprobările F07-T12, integrate din a132a27.

Actualizarea recunoaște și istoricul exact al probei 0.20.0, completând migrarea cererilor lipsă înaintea imaginilor; modelele existente se păstrează. Celelalte istorice incomplete rămân refuzate.

Rapoarte: **Pontaj lunar v2, Fișă individuală v2, Evidența absențelor v2**, cu referințe noi, păstrând referințele v1. Baza: **0014_imagini_rapoarte**, fără schimbarea textelor, reviziilor sau momentelor editărilor modelelor existente.

## 0.20.2 · 2026-10-25

F09-T08, actualizare separată a dependențelor declarate: pdfcpu 0.16.2 → 0.16.3 (Apache-2.0), go-runewidth 0.0.30 → 0.0.31 (MIT). Licențele noilor versiuni verificate și păstrate. Go 1.27.3 și npm sunt la versiunile verificate; fără schimbări ale rapoartelor sau bazei.

## 0.20.1 · 2026-10-25

F09-T08: integrarea cererilor de concediu F07-T12 din a132a27 în ramura rapoartelor; migrările 0012 și 0013 coexistă, iar documentele și stările de aprobare rămân în aplicație. Pregătire pentru imagini și paginare în 0.21.0; rapoartele sunt încă v1.

## 0.20.0 · 2026-10-25

F09-T07: migrarea rapoartelor este renumerotată 0013, iar decizia tehnică devine D025, pentru a lăsa 0012/D024 taskului paralel F07-T12 (cereri de concediu). Versiunea 0.19.0 este folosită și de acel task, încă necomis; livrarea rapoartelor folosește 0.20.0. Nu include modificările nesalvate F07-T12. Integrarea lor este necesară înaintea folosirii ambelor funcționalități pe aceeași bază. Rapoartele rămân v1.

## 0.19.1 · 2026-10-25

Integrare F09-T07 cu categoriile LAB/CUSTOM din F07-T11 (bacd37b), pe ramura rapoartelor. Include migrările 0011 și 0012, păstrează filtrele și formularele categoriilor și regenerează accesul la date și puntea Wails. Versiunile rapoartelor rămân v1.

## 0.19.0 · 2026-10-25

Definire rapoarte și PDF pentru tipărire (F09-T07, D025).

- Meniu separat pentru modele cu antet, conținut și subsol proprii, texte RO/EN cu câmpuri de context, surse/coloane configurabile, format A4/A3, orientare, margini și font.
- Definiții persistente cu revizii și jurnal, creare/editare/duplicare/dezactivare. Modele inițiale pentru pontaj lunar, fișă individuală și evidența anuală a absențelor.
- Generare PDF locală cu fonturi incluse pentru diacritice, antet/subsol repetate, paginare și metadatele D009. Previzualizare, salvare explicită și deschidere în aplicația PDF pentru alegerea imprimantei.
- Datele sunt o fotografie consistentă a schițelor și absențelor; nu sunt deduse totaluri salariale sau reguli neconfirmate.

Rapoarte: **Pontaj lunar v1, Fișă individuală v1, Evidența absențelor v1.** Baza de date: **0013_definitii_rapoarte**, după migrarea 0011 rezervată taskului paralel pentru categorii.
## 0.19.0 · 2026-10-25

Cereri de concediu semnate atașate, număr de înregistrare, persoană înlocuitoare și stare de aprobare pentru concedii și absențe (F07-T12, D024).

- Zilele de concediu și absență includ opțional numărul de înregistrare al cererii, persoana care înlocuiește angajatul absent (validată ca fiind diferită de angajatul în concediu), starea aprobării (`IN_ASTEPTARE`, `APROBATA`, `RESPINSA`), numele celui care a aprobat și un document semnat atașat (PDF/PNG/JPG până la 10 MB), stocat în tabela `absente_documente` din baza SQLite criptată și partajat la salvarea unui interval.
- Ecranul Concedii și absențe afișează și filtrează după starea aprobării, permite atașarea și descărcarea cererii semnate și oferă completarea ulterioară a cererii, înlocuitorului sau aprobării pentru o zi ori pentru un lot de zile selectate.
- Grila de pontaj afișează în explicația celulei și prin stil vizual starea aprobării, numărul de înregistrare și persoana înlocuitoare pentru zilele de absență active.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0012_cereri_concediu`.**

## 0.18.0 · 2026-10-25

Categorii în nomenclatorul de schimburi: `LAB` pentru precompletare și `CUSTOM` pentru restul (F07-T11, D023).

- Câmpul `categorie` (`LAB` / `CUSTOM`) este separat de `tip` și `regim` în `tipuri_schimb`.
- Codurile `G24`, `IN9`, `N12`, `N9`, `Z7`, `Z12`, `S7` intră în categoria de precompletare `LAB`; celelalte coduri existente (`7`, `12`, `S1`, `S3`) și schimburile create manual intră în `CUSTOM`, cu păstrarea intervalelor și personalizărilor.
- Ecranul Schimburi oferă coloană și filtru pentru categorie, alegerea categoriei la precompletare și în formularul de adăugare/editare; sugestiile din pontaj includ categoria.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0011_categorii_schimburi`.**

## 0.17.3 · 2026-10-25

Pregătire F09-T07: maroto v2.4.4 și x/image v0.47.0 pentru PDF local și fonturi cu diacritice, versiuni exacte și licențe verificate în D025. Include dependențele tranzitive ale bibliotecii PDF și actualizarea i18next 26.5.10 → 26.5.11; fără instalare de programe pe sistem. Funcționalitatea și testele PDF urmează în 0.19.0, pe worktree separat de F07-T11.
## 0.17.2 · 2026-10-25

Actualizare de dependențe (F07-T11, D007).

- `react-i18next` actualizat 17.1.8 → 17.1.9 (MIT).

Rapoarte: nicio schimbare. Baza de date: nicio schimbare.

## 0.17.1 · 2026-10-25

Corecție a detectării conflictelor la salvări rapide (F06-T07, D022).

- Revizie numerică per rând, incrementată în SQL, separată de timestampul de audit în schițe și absențe.
- Două modificări în același tic al ceasului au revizii distincte; o salvare pe baza reviziei vechi refuză întregul lot. Conținutul și momentele înregistrărilor existente se păstrează la migrare.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0010_revizii_concurenta`.** Include programele și IN9 din 0.17.0.

## 0.17.0 · 2026-10-25

Programele reale și intrarea specială în weekend/sărbătoare (F07-T10, D021).

- Z7 07–14, S7 15–22, N9 22–07; IN9 22–07 distinct, început în seara precedentă zilei libere. Pentru weekenduri/sărbători: Z12 07–19, N12 19–07 și G24 07–07.
- Regim separat de tipul activității, vizibil și editabil în nomenclator, filtrabil în listă și afișat în sugestiile din pontaj. GENERAL păstrează configurațiile vechi neclasificate.
- Completarea adaugă numai lipsurile. Codurile istorice și personalizările se păstrează; numai metadatele șabloanelor originale nemodificate N12/G24 sunt încadrate automat în regimul pentru zile libere.
- Fără alocare automată de ture, reguli noi de plată sau interpretarea golului 14–15 ca pauză scăzută.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0009_programe_lucru`.**

## 0.16.0 · 2026-10-25

Modul Concedii și absențe (F07-T04, D020).

- Tipuri de odihnă, medical și zile speciale; CO/CM inițiale și tipuri noi cu cod/denumiri RO/EN, editare și dezactivare. Codurile nu se confundă cu schimburile, iar categoria unui tip folosit se păstrează.
- Înregistrare explicită pe zile sau intervale cu selecție individuală, weekenduri/sărbători vizibile și selectare opțională luni–vineri fără sărbători. Editare, anulare individuală/în lot cu istoric și revizii pentru conflicte.
- Sumar anual pe angajat și categorie, cu dreptul contractual curent afișat separat. Pontajul arată absențele și le protejează împotriva suprascrierii; înregistrarea refuză zilele cu pontaj existent, inclusiv ture care traversează ziua.
- Evidența nu produce solduri legale, reportări, indemnizații ori calcule de plată. Baza standard Go de fusuri orare este inclusă pentru executabilul portabil.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0008_concedii_absente`.**

## 0.15.0 · 2026-10-25

Schimburi inițiale din PDF și intervalele precizate de Ioan Tuns (F07-T09, D019).

- 7: 07–14; 12: 14–02 a doua zi; N12: 19–07; G24: 07–07; S1: 06–14; S3: 22–06.
- Migrarea adaugă numai codurile lipsă și păstrează configurațiile existente, inclusiv cele inactive ori șterse logic. Pauza păstrează implicitul tehnic existent de 0 minute, editabil; nu se introduc reguli de plată.
- Ecranul Schimburi arată sfârșitul nominal și trecerea în ziua următoare, plus completare explicită a lipsurilor din setul inițial.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0007_schimburi_initiale`.**

## 0.14.0 · 2026-10-25

Departamente în cadrul secțiilor (F07-T08, D018).

- Setări: creare de departamente sub o secție și afișare ierarhică.
- Angajați: departament opțional, ales numai dintre cele ale secției; schimbarea secției golește selecția. Filtre pe secție și departament în listă și grila lunară.
- Migrarea păstrează secțiile și angajații existenți, fără atribuiri automate; validarea relației se face în serviciu și prin constrângeri/trigger în bază. Schimburile existente se păstrează.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0006_departamente`.**

## 0.13.0 · 2026-10-25

Prima grilă lunară după structura modelului primit (F08-T01, D016).

- An/lună, filtru pe secție, zilele reale ale lunii, weekenduri și sărbători, antet și identificarea angajatului fixate la derulare.
- Introducere din tastatură, selecție cu Shift, copiere/lipire TSV și completare pe intervale; sugestii din schimburile active.
- Salvarea explicită și tranzacțională a schiței text, cu revizii per celulă, refuzul suprascrierii la conflict și jurnal. Avertizare pentru modificări nesalvate.
- Coloanele de total din model sunt prezente, necalculate. Codurile necunoscute nu se transformă în ore sau plăți.

Rapoarte: niciun PDF generat încă. **Baza de date: migrarea `0005_pontaj_ciorne`.**

## 0.12.1 · 2026-10-25

Corecția build-ului pentru păstrarea datelor portabile (F05-T15, D017).

- Eliminată curățarea `build/bin`, care poate conține folderul de date `karios/`.
- Parametru `-NumeExecutabil` pentru o ieșire alternativă; căile invalide și fișierele blocate sunt refuzate înainte de compilare.
- Proprietățile Windows sunt verificate pentru fișierul de ieșire ales. Test PowerShell de regresie pentru refuzuri.

Rapoarte: nicio schimbare. Baza de date: nicio migrare nouă.

## 0.12.0 · 2026-10-25

Nomenclatorul configurabil de schimburi (F07-T02, D015).

- Ecran Schimburi: cod, denumiri RO/EN, oră de început, durată și pauză în minute, tip și activ/inactiv.
- Creare/editare și dezactivare/reactivare, validări și jurnal tranzacțional. Coduri normalizate la majuscule, fără duplicate, inclusiv între schimburile inactive.
- Nu se introduc automat ore sau reguli de plată pentru N12 ori alte coduri.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.11.0 · 2026-10-25

Sărbători legale România 2026 și administrarea calendarului pe ani (F07-T03, D014).

- Migrarea de date include lista națională: 17 sărbători nominale în 16 date distincte; 1 iunie combină Ziua Copilului și a doua zi de Rusalii. Datele mobile sunt explicit ortodoxe. Înregistrările preexistente nu sunt suprascrise.
- Ecran Sărbători: selectare an, listare, creare/editare cu denumiri RO/EN, validarea datei reale și completarea repetabilă a datelor lipsă din 2026.
- Operațiile din interfață sunt tranzacționale și jurnalizate; nu se introduc reguli de plată.

Rapoarte: niciunul încă. **Baza de date: migrarea `0004_sarbatori_romania_2026`.**

## 0.10.0 · 2026-10-25

Primele fluxuri complete în interfața desktop (F07-T01, F07-T07).

- Ecran Angajați: căutare și filtre, adăugare, editare, dezactivare/reactivare, normă efectivă și concediu anual contractual. Formular cu validare și protecție la abandonarea modificărilor.
- Setări: norma întreagă globală și adăugarea secțiilor. Norma parțială rămâne individuală.
- Date demo fictive: 12 angajați și 3 secții; încărcare explicită, tranzacțională, numai în lipsa angajaților și pontajelor. Fără modificarea normei globale.
- Serviciu Go cu operații tranzacționale și jurnal pentru aceste fluxuri; erori structurate și texte RO/EN din cataloagele comune.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.9.1 · 2026-10-25

Acces la date generat cu sqlc v1.31.1 (F06-T03), pentru angajați, secții, setări și jurnal. Configurația folosește direct migrările existente; `genereaza-date.ps1` fixează versiunea generatorului fără instalare globală. Nu adaugă o dependență în executabil.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.9.0 · 2026-10-25

Schema datelor de bază cerute pentru angajați și sărbători (F07-T06, D012).

- Angajați: nume și prenume separate, marcă, contract `INTREAGA` implicit sau `PARTIALA`, norma parțială individuală în minute și zile anuale de concediu conform contractului (`NULL` până la completare).
- Norma întreagă globală: `setari.norma_intreaga_minute`, implicit `420` (7 ore/zi). Vederea `angajati_norme` expune `norma_efectiva_minute`: globală pentru norma întreagă, individuală pentru norma parțială.
- Datele existente se păstrează: numele vechi nu se desparte automat, norma veche rămâne în `norma_minute_anterioara`; noile informații contractuale necesită verificare/completare manuală. Norma veche nu intră în noul calcul al normei curente.
- Sărbători definite pe dată completă, cu `an` derivat și index pentru filtrare anuală; nu se introduc automat sărbători legale.
- Constrângeri pentru tipul contractului, norma parțială, concediu nenegativ și norma globală pozitivă; cheia normei globale nu poate fi eliminată sau redenumită.

Această versiune pregătește baza de date; formularele, accesul la date prin sqlc și istoricul contractual rămân în taskuri separate.

Rapoarte: niciunul încă. **Baza de date: migrările `0002_angajati_contracte` și `0003_sarbatori_an`.**

## 0.8.1 · 2026-10-25

Actualizarea dependențelor pentru F07-T06 (D007): SQLite WASM `v6.3.35304` → `v6.3.35401` în aplicație și prototip; `react-i18next` `17.1.7` → `17.1.8`. Versiuni fixate exact în manifestele și fișierele de blocare existente.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.8.0 · 2026-10-25

Backup automat înaintea actualizărilor bazei de date (F05-T06, D002 punctul 3).

- Înainte de orice migrare pe o bază care are deja date, aplicația face o copie criptată în `karios/backup/pontix-AAAA-LL-ZZ-HHMMSS-vX.Y.Z.db` (`bd.Backup`, cu `VACUUM INTO`).
- Dacă backup-ul eșuează, migrarea nu pornește și aplicația arată eroarea `stocare.backup_esuat`.
- Se păstrează cele mai noi **7** backup-uri automate; alte fișiere din `backup/` nu se ating.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.7.0 · 2026-10-25

Aplicația își deschide datele criptate la pornire (F05-T12, D011).

- Cheia aplicației stă în cod (decizia D011): datele din `karios/` nu se pot citi fără executabil.
- Pachetul `internal/stocare`: la prima pornire creează `karios/seif` și `karios/pontix.db`; la următoarele le deschide și aplică migrările lipsă.
- Dacă baza există, dar seiful lipsește, aplicația nu creează un seif nou (ar face baza ilizibilă), ci afișează eroarea `seif.deteriorat`.
- Baza se închide la ieșirea din aplicație.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă (prima pornire aplică `0001_schema_initiala`).

## 0.6.0 · 2026-10-25

Seiful cheii de date (F05-T10, D008 punctul 2).

- Pachetul `internal/seif`: generează cheia de date și o păstrează în `karios/seif`, în copii criptate cu AES-256-GCM: copia „A” cu cheia aplicației și copii deschise cu o parolă sau cu răspunsuri (Argon2id).
- Cheia aplicației se injectează la build (`seif.CheieAplicatie`); fără ea, seiful refuză să funcționeze.
- Textele noi: `seif.fara_cheie_aplicatie`, `seif.deteriorat`, `seif.secret_gresit`.
- Nu e încă folosit de aplicație; integrarea vine cu F05-T12.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.5.0 · 2026-10-25

Baza de date criptată și schema inițială (F06-T02).

- Pachetul `internal/bd`: `Deschide` (SQLite criptat AES-256-XTS, modul WAL, cheile străine active, cheia verificată la deschidere) și `Migreaza` (goose).
- Dependențe noi: `ncruces/go-sqlite3` v0.35.6 și `pressly/goose/v3` v3.29.0 (D008, D010).
- Textele noi: `bd.nu_se_deschide`, `bd.cheie_invalida`.
- Pachetul nu e încă folosit de aplicație; legătura vine cu seiful cheii (F05-T10).

Rapoarte: niciunul încă. **Baza de date: migrarea `0001_schema_initiala`** — tabelele `setari`, `sectii`, `angajati`, `tipuri_schimb`, `coduri_absenta`, `sarbatori`, `categorii_ore`, `seturi_reguli`, `reguli_interval`, `pontaj`, `jurnal_modificari`.

## 0.4.0 · 2026-10-24

Mesaj clar când folderul de date nu permite scrierea (F05-T05, D002 punctul 5).

- La pornire, aplicația verifică faptul că poate scrie în `karios/`. Dacă nu poate, afișează o fereastră Windows cu un mesaj tradus, în limba Windows, care spune unde să fie mutată aplicația.
- `i18n.Eroare` e acum tipul comun de eroare traductibilă; `calcul.Eroare` e un alias.
- `i18n.LimbaSistem`: limba interfeței Windows, dacă are catalog.
- Textele noi: `aplicatie.eroare_pornire`, `portabil.fara_scriere`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.3 · 2026-10-24

Scriptul de build (F05-T13).

- `build.ps1` construiește `pontix.exe` cu versiunea în proprietățile fișierului și commit-ul curent (cu sufixul `-modificat` dacă există fișiere necomise).
- `wails.json` are secțiunea `info` (nume, versiune, copyright); testul `TestWailsJSON` verifică faptul că versiunea e aceeași ca în `versiune.go`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.2 · 2026-10-24

Prototipul migrărilor (F06-T01). Aplicația nu se schimbă.

- `prototipuri/sqlite-criptat/migrari_test.go`: goose v3.29.0 pe baza criptată; rezultatele sunt în D010.
- Hook-ul `pre-commit` protejează doar migrările existente în `main`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.1 · 2026-10-24

Prototipul bazei de date criptate (F05-T09). Aplicația nu se schimbă.

- `prototipuri/sqlite-criptat/`, modul Go separat: `ncruces/go-sqlite3` v0.35.6 cu VFS-urile `adiantum` și `xts`. Toate verificările din D008 trec; rezultatele sunt în D008.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.0 · 2026-10-24

Datele lângă executabil (F05-T03, D002).

- Pachetul `internal/portabil` creează la pornire folderul `karios/` lângă `.exe`, cu `webview/` și `backup/`.
- Cache-ul WebView2 stă în `karios/webview/` (`WebviewUserDataPath`). Verificat: aplicația nu mai creează nimic în `%APPDATA%` sau `%LOCALAPPDATA%`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.2.0 · 2026-10-24

Scheletul aplicației desktop (F05-T02).

- Aplicația Wails v2.17.0 cu interfața React 19.3.1 + TypeScript 7.0.2, construită cu Vite 8.3.5. `pontix.exe` are 11,5 MB.
- Interfața ia textele din cataloagele Go prin i18next 26.5.10 / react-i18next 17.1.7 (F03-T06, parțial); funcția nouă `i18n.Catalog`.
- Fereastra arată numele, descrierea și versiunea aplicației; commit-ul se injectează la build.
- Textele noi: `aplicatie.titlu`, `aplicatie.subtitlu`, `aplicatie.versiune`.
- Dependențele Go actualizate la ultimele versiuni stabile (D007); pachetele npm fixate exact (`save-exact`).

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.1.0 · 2026-10-24

Prima versiune numerotată. Cuprinde tot codul de până acum.

- Motorul de calcul: `Segmenteaza`, `Valideaza`, calendarul cu sărbători (F04-T01, F04-T02).
- Traducerile: cataloagele RO și EN, folderul de traduceri, `i18n.Text`, erori traductibile (F03-T01 … T04).
- Go 1.27.3 (F02-T09).
- Versionarea: pachetul `internal/versiune`, verificarea la commit (F02-T10).

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.
