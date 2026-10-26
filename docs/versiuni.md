# Versiuni

Istoricul versiunilor aplicației, cea mai nouă prima ([D009](decizii/D009-versionare.md)). Fiecare commit care schimbă codul crește versiunea din `internal/versiune/versiune.go` și adaugă aici o secțiune `## X.Y.Z · AAAA-LL-ZZ`.

Pe lângă versiunea aplicației, fiecare raport și fiecare schimbare a bazei de date au versiunea lor; ele se trec în secțiunea versiunii aplicației care le aduce.

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
