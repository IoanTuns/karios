# Analiza UI/UX · Pontix 0.27.0 · 26 octombrie 2026

F05-T18 · Codex · Bază analizată: 2f55173, implementare UI 069f850. Concluziile și prioritățile sunt propuneri; auditul nu modifică interfața sau regulile de calcul.

Interfața are o identitate vizuală coerentă și instrumente utile pentru operare repetitivă. Principalele probleme sunt spațiul ocupat înaintea datelor, pierderea contextului la navigare și câteva mesaje care nu mai descriu corect funcțiile disponibile. Acestea merită rezolvate înaintea unor schimbări decorative.

## Metodă și limite

- Revizuire App, Pontaj, Absente, ReguliPontaj, Rapoarte, Formulare, CSS, cataloage și D030–D034.
- Edge cu interfața compilată curentă și punte Wails simulată, numai date fictive. Au trecut fluxurile existente, inclusiv grila 50×31, selecție, Undo/Redo, configurare, conflicte și rapoarte. Probe vizuale la 1440×960, 1440×900 și 1024×640.
- Probă suplimentară de revenire între Pontaj și Rapoarte, cu an/lună/căutare selectate explicit.
- Capturile, măsurătorile și scriptul temporar sunt în `build/bin/audit-ui-20261026/`, separat de builduri. Nu s-au deschis bazele utilizatorului.
- PDF-ul din simularea UI este un substitut minimal: proba verifică acțiunile interfeței, nu fidelitatea documentului sau imprimarea fizică. Nu s-au făcut cercetare cu utilizatori, audit WCAG complet ori măsurători de performanță pe baze mari.

Probe vizuale locale: [pontaj 1024×640](../build/bin/audit-ui-20261026/pontaj-1024.png), [pontaj 1440×900](../build/bin/audit-ui-20261026/pontaj-1440.png), [editor raport](../build/bin/audit-ui-20261026/editor-rapoarte-1024.png), [configurare reguli](../build/bin/audit-ui-20261026/reguli-1024.png).

## Avantaje

| Zonă | Ce funcționează bine | Beneficiu |
| --- | --- | --- |
| Navigare | Cinci intrări principale, configurarea grupată în Setări, iconuri distincte și locație vizibilă în pagină | Separă lucrul curent de configurare |
| Aspect | Culori temperate, componente recognoscibile, Setări în albastru-ardezie, tipografie și spațiere consecvente | Interfața este ordonată și ușor de identificat |
| Pontaj | Ore afișate, identitatea angajatului și T fixe, navigare cu tastatura, selecție prin tragere, TSV, paletă rapidă, Undo/Redo | Susține introducerea repetitivă fără deschiderea unui formular pentru fiecare zi |
| Înțelegerea celulei | Bara de inspecție arată programul, orele, avertismentele și înlocuirile; salt contextual la concedii | Utilizatorul poate explica o valoare și continua verificarea |
| Protecție | Modificări nesalvate semnalate, absențe protejate, conflicte și mesaje de salvare | Reduce pierderea accidentală și suprascrierile concurente |
| Concedii | Căutare după înregistrare, selecția aceleiași cereri, stări de aprobare, previzualizarea atașamentului și sumar separat | Oferă verificare și administrare din același modul |
| Rapoarte | Centralizator și detaliat separate; antet/conținut/subsol, imagini, RO/EN, previzualizare și export | Modele adaptabile documentelor instituției |
| Acces de bază | Focus vizibil, etichete accesibile în multe controale, dialog nativ și selector de limbă vizibil | Există o bază bună pentru operare cu tastatura și îmbunătățirea accesibilității |

Funcțiile deja implementate prin D032–D034 nu sunt repropuneri: Undo/Redo, căutarea în pontaj, selecția prin tragere, totalul T fix și gruparea Setărilor sunt prezente.

## Probleme și îmbunătățiri

P1 = primul lot recomandat, impact mare în utilizarea zilnică. P2 = al doilea lot. Eforturile S/M sunt estimări relative, nu termene.

| ID / prioritate | Constatare și dovadă | Îmbunătățire propusă | Efort |
| --- | --- | --- | --- |
| I045 · P1 | Bannerul pontajului spune că totalurile rămân necalculate, deși T/ZNL/Ore norm și regulile definite se calculează. În plus, `T != ore_norma` este etichetat „Diferență față de normă”, deși Ore norm reprezintă orele lucrate în zile normale, nu obligația contractuală lunară | Actualizarea bannerului și delimitarea clară a totalului, orelor din zile normale și unei eventuale norme contractuale. Diferența T–Ore norm nu trebuie prezentată drept deficit sau ore suplimentare. Nu se inventează un calcul contractual | S pentru texte; validarea indicatorului M |
| I046 · P1 | La 1024×640, grila începe la y=629 și are doar 130 px înălțime totală; nu se vede niciun rând complet de la începutul paginii. La 1440×900 are 399 px. Filtrele, mesajele, comenzile și instrucțiunile ocupă mare parte din ecran | Bară compactă pentru perioadă/salvare, grupuri clare Filtre / Completare / Verificare, opțiuni avansate pliabile și ajutor la cerere. Păstrarea inspecției cu înălțime controlată. Contrast mai bun și legendă textuală pentru stări | M |
| I047 · P1 | Reproducere: Pontaj 11/2028, căutare TEST-001 → Rapoarte → Pontaj. Revenirea afișează 10/2026 și căutare goală. Starea locală se reinițializează la remontare | Păstrarea anului, lunii, filtrelor și poziției de lucru în sesiune; revenire explicită la context după inspectarea concediilor. Protecția modificărilor nesalvate rămâne activă | M |
| I048 · P1 | Editorul de raport are dialog de 592 px și conținut de 3213 px în proba cu imagini. Bara Salvare/Anulare este la y=3163 și are poziționare statică. În configurarea regulilor, tabelul are aproximativ 1184 px, iar salvarea este numai deasupra lui | Bară de acțiuni fixă în dialog, derulare în conținut, acces persistent la salvarea regulilor; editor cu secțiuni Antet / Conținut / Subsol / Pagină | M |
| I049 · P1 | „—” și mesajul comun „Regulă nedefinită sau înregistrări fără durată calculabilă” acoperă cauze diferite. Erorile din configurări sunt afișate global, fără evidențierea câmpului relevant | Motiv distinct și acțiune concretă: definește regula, selectează codurile, verifică programul/pauza sau reîncarcă revizia. Erori lângă câmp și rezumat accesibil; motive structurate de backend unde lipsesc | M |
| I050 · P2 | Concediile pornesc în „Toate secțiunile”. La 1024×640 se vede începutul sumarului, iar evidența zilelor rămâne mult mai jos. Datele unei cereri se repetă pe fiecare zi, chiar dacă selectarea întregii cereri există deja | Deschidere pe Zile înregistrate sau pe ultima secțiune; vedere grupată pe cereri cu detalii expandabile. Sumarul rămâne explicit aritmetic, fără interpretări legale suplimentare | S/M |
| I051 · P2 | Câmpurile `{{...}}` din editor sunt coduri statice, iar eroarea pentru un câmp necunoscut este globală. Previzualizarea din editor este structurală, separată de PDF-ul generat | Inserarea variabilei prin clic la cursor, validare pe zona greșită și diferențiere explicită între structură și previzualizarea PDF. Folosirea datelor fictive pentru verificarea unui model | M |
| I052 · P2 | Setările au simultan submeniu lateral și taburi; la 1024 px submeniul activ necesită derulare. Selecțiile regulilor sunt oferite atât ca listă multiplă, cât și ca o serie completă de cipuri | Submeniu pliabil sau compact, cu taburile păstrate la vedere; un singur control principal de selecție cu căutare și opțiunile selectate clar afișate. Păstrarea unei alternative accesibile cu tastatura | S/M |

## Lizibilitate: verificări punctuale

Textul din grilă are 12 px, identificarea și antetele 9–11 px, iar multe acțiuni secundare 11 px. Dimensiunile nu constituie singure o neconformitate; în combinație cu fundaluri colorate și contrast redus devin greu de urmărit în lucru prelungit. Este util un mod compact/confortabil, cu text mai mare pentru formular și instrucțiuni.

Perechi de culori preluate din stilurile curente, verificate și în browser:

| Text / fundal | Raport calculat | Observație |
| --- | --- | --- |
| Ore `#167d9d` / alb | 4,71:1 | Peste reperul pentru text mic |
| Ore `#167d9d` / weekend `#e5e8e5` | 3,82:1 | Sub reper; datele din weekend trebuie să rămână lizibile |
| Ore `#167d9d` / selecție `#dcefe5` | 3,93:1 | Sub reper; selecția reduce contrastul textului |
| Text nedefinit `#859087` / `#f8f9f6` | 3,14:1 | Valoarea informativă „—” trebuie să poată fi citită |
| Placeholder `#87948c` / alb | 3,16:1 | Informația utilă nu trebuie să depindă de un placeholder pal |

Reperul pentru text obișnuit este 4,5:1, conform [W3C, WCAG 2.2, contrast minim](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html). Acestea sunt măsurători punctuale, nu un verdict de conformitate asupra aplicației. Nu au fost tratate butoanele dezactivate ca erori de contrast.

## Ordinea recomandată

1. Claritatea mesajelor și a indicatorului T/Ore norm, apoi contrastul stărilor informative.
2. Spațiul de lucru în pontaj și păstrarea contextului la navigare.
3. Salvare accesibilă permanent și motive clare pentru „—”/erori.
4. Gruparea cererilor, simplificarea editorului și reducerea controalelor duplicate.

Persistența limbii selectate rămâne o limitare cunoscută F03-T05/I034; selectorul vizibil este implementat, iar inițializarea folosește în continuare limba întoarsă de backend.

## Criterii de acceptare pentru un viitor lot UI

- La 1024×640, fără derularea paginii, să fie vizibile antetul grilei și minimum cinci rânduri, păstrând selecția și inspecția utilizabile. Pragul este o țintă propusă de ergonomie.
- Pontaj → Rapoarte/Concedii → Pontaj păstrează perioada și filtrele. Modificările nesalvate nu sunt abandonate implicit.
- Un caz cu T=9, Ore norm=2 și ZNL=7 nu este etichetat ca deficit contractual; textele descriu calculul existent.
- Salvarea/Anularea editorului sunt vizibile indiferent de poziția derulării, inclusiv cu imagini în șase zone.
- Valorile, stările și câmpurile active îndeplinesc reperul de contrast; culoarea este însoțită de explicație textuală.
- O valoare necalculabilă indică motivul specific și pasul de corectare; configurările valide păstrează comportamentul și rezultatele actuale.

Următorul pas este alegerea unui lot de implementare din aceste propuneri, urmat de teste pe date fictive și o probă cu un operator. Auditul nu solicită modificarea regulilor de plată sau a cazurilor confirmate de contabilitate.
