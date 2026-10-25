# Funcționalitățile disponibile · Pontix 0.18.0

## Pornire

Executabilul se construiește cu `./build.ps1` și se găsește în `build/bin/pontix.exe`. Pornește-l dintr-un folder în care poți scrie; datele sunt în folderul `karios/` de lângă executabil.

Livrarea actualizată este **`build/bin/pontix-0.13.0/pontix-0.17.1.exe`**, lângă aceleași date demo regenerate anterior la cererea utilizatorului. Închide versiunea veche înainte de pornire. Numele folderului a fost păstrat pentru a continua cu baza existentă; actualizarea copiază numai executabilul. La prima pornire, aplicația face backup înaintea migrărilor lipsă (0009 pentru programe și 0010 pentru revizii) și păstrează configurațiile existente. Baza veche separată a rămas în `build/bin/karios/`; seiful ei a fost eliminat de vechiul build cu `-clean`, deci nu se reutilizează cu o cheie nouă. Incidentul și corecția sunt în [D017](decizii/D017-build-fara-stergerea-datelor.md).

Scriptul corectat nu mai curăță folderul de ieșire. Dacă executabilul este deschis, construiește cu alt nume: `./build.ps1 -NumeExecutabil pontix-verificare.exe`.

## Încearcă datele demo

1. Deschide **Angajați** într-o bază fără angajați sau pontaje.
2. Apasă **Încarcă date demo** și citește explicația.
3. Confirmă cu **Încarcă exemplele**.
4. Apar 12 angajați fictivi și 3 secții. Datele rămân salvate după repornire.

Setul are 10 persoane active, 2 inactive și 3 contracte parțiale (3, 4 și 5 ore/zi). Concediile sunt exemple; pentru o persoană concediul este necompletat, pentru alta este zero. Numele, funcțiile și mărcile sunt fictive; nu reprezintă reguli de plată. Fișierul poate fi citit separat în [angajati.json](../internal/datebaza/demo/angajati.json).

Importul nu rulează automat, nu suprascrie date și nu schimbă norma globală. Nu există o comandă de golire a bazei. Dacă există deja angajați sau pontaje, importul este refuzat și din backend, nu doar ascuns în interfață.

## Angajați

- Caută după nume, prenume sau marcă; căutarea ignoră diferențele de diacritice și majuscule.
- Filtrează după secție, departament și activ/inactiv.
- **Adaugă angajat**: nume, prenume și marcă sunt obligatorii; marca este text și păstrează zerourile inițiale.
- Contractul este implicit cu normă întreagă și preia norma globală. Pentru norma parțială completează orele individuale pe zi, inclusiv fracții care reprezintă minute întregi (de exemplu 4,5 ore = 4 h 30 min).
- Concediul anual este cel din contract. Câmp gol înseamnă necompletat, nu zero. Nu se calculează încă soldul concediului.
- **Editează** modifică datele; la anularea unui formular modificat se cere confirmarea abandonării.
- **Dezactivează / Reactivează** schimbă starea fără a șterge persoana sau pontajul existent.

Datele migrate cu prenume necompletat trebuie verificate manual înainte de salvarea formularului. Nu se deduce automat prenumele din numele vechi.

## Setări

- Modifică norma întreagă în ore/zi (implicit 7). Modificarea se vede imediat la angajații cu normă întreagă; normele parțiale rămân individuale.
- Adaugă secții, apoi atribuie-le din formularul angajatului. Denumirea se păstrează ca text traductibil în limba interfeței.
- Adaugă mai multe departamente sub fiecare secție. În formularul angajatului alegi numai dintre departamentele secției sale; schimbarea secției golește departamentul. Atribuirea este opțională, iar migrarea nu atribuie automat departamente angajaților existenți.

## Sărbători

La actualizarea aplicației, migrarea 0004 adaugă lista națională pentru 2026, fără a suprascrie înregistrări existente: 17 sărbători nominale în 16 date distincte; la 1 iunie coincid Ziua Copilului și a doua zi de Rusalii. Setul folosește datele mobile ortodoxe, indicate în denumiri și în ecran. Sursele și limitele calendarului comun sunt în [D014](decizii/D014-sarbatori-2026.md).

Alege **Sărbători**, anul și **Arată anul**. Poți adăuga sau edita o dată și denumirile RO/EN. **Completează lista legală 2026** adaugă numai date lipsă și păstrează denumirile editate. Nu se introduc automat zile-punte sau reguli de plată.

## Schimburi

Migrarea 0009 adaugă programele confirmate ulterior de utilizator, numai pentru codurile lipsă, iar migrarea 0011 împarte tabelul de schimburi pe categorii (`Lab` pentru precompletare și `Custom` pentru restul):

| Cod | Categorie | Început | Sfârșit | Durată | Regim |
| --- | --- | --- | --- | --- | --- |
| Z7 | Lab | 07:00 | 14:00 | 7 ore | Program obișnuit |
| S7 | Lab | 15:00 | 22:00 | 7 ore | Program obișnuit |
| N9 | Lab | 22:00 | 07:00 a doua zi | 9 ore | Program obișnuit |
| IN9 | Lab | 22:00 | 07:00 a doua zi | 9 ore | Intrare specială în zi liberă |
| Z12 | Lab | 07:00 | 19:00 | 12 ore | Weekend / sărbătoare |
| N12 | Lab | 19:00 | 07:00 a doua zi | 12 ore | Weekend / sărbătoare |
| G24 | Lab | 07:00 | 07:00 a doua zi | 24 ore | Weekend / sărbătoare |

**IN9 începe în seara precedentă zilei libere.** Pentru intrarea în sâmbătă se scrie IN9 în celula de vineri. N9 și IN9 au aceeași durată nominală, dar regimuri distincte. Găsești IN9 prin filtrul **Regim de lucru → Intrare specială în zi liberă**; regimul și categoria pot fi alese explicit și în formular. Z12/N12 sau G24 sunt opțiuni de acoperire; nu se repartizează automat pe persoane.

**Completează schimburile inițiale** adaugă lipsurile din categoria de precompletare aleasă (`Lab`), fără a suprascrie coduri existente, inclusiv inactive. Codurile vechi 7/12/S1/S3 și schimburile adăugate manual intră în categoria **Custom**, cu intervalele și regimurile lor păstrate. Șabloanele folosesc pauza tehnică implicită 0, editabilă; golul 14–15 nu se scade automat din vreo tură. Sursele și limitele sunt în [D023](decizii/D023-categorii-schimburi.md) și [D021](decizii/D021-programe-lucru.md); [D019](decizii/D019-schimburi-initiale-pdf.md) descrie setul istoric.

Din **Schimburi → Adaugă schimb**, introdu codul, denumirea, ora de început, durata și pauza în minute și tipul activității. Intervalele rămân editabile. Coloana **Sfârșit nominal** arată și trecerea în ziua următoare. Codurile se normalizează la majuscule și nu pot fi duplicate, nici între schimburile inactive. Dezactivarea păstrează referințele din pontaj.

## Concedii și absențe

1. Deschide **Concedii și absențe** și alege anul. Poți filtra lista după angajat și categorie.
2. Pentru un tip special nou, apasă **Adaugă tip de absență**: introdu un cod unic (de exemplu `EV_FAMILIE`), denumirea RO/EN și categoria **Zile libere speciale**. Poți defini câte tipuri sunt necesare. CO și CM sunt inițiale pentru odihnă și medical, dacă aceste coduri nu erau deja ocupate.
3. Apasă **Înregistrează zile**, alege angajatul, tipul și intervalul de cel mult 366 de zile. Inițial sunt bifate toate datele, inclusiv weekendurile și sărbătorile. Debifează zilele care nu trebuie înregistrate sau folosește explicit **Selectează luni–vineri, fără sărbători**, pe baza calendarului din aplicație. Pentru angajați cu ture, alege datele potrivite programului lor.
4. Verifică numărul de zile selectate și apasă **Salvează**. Lotul se salvează integral sau deloc. O zi deja ocupată de o absență, o schiță nevidă ori un pontaj efectiv trebuie corectată înainte de înregistrare; turele efective care traversează ziua sunt și ele verificate.
5. **Editează** schimbă tipul unei zile existente. **Anulează ziua** sau bifarea mai multor zile și **Anulează zilele selectate** retrag înregistrările cu păstrarea istoricului. În cazul unei modificări concurente, reîncarcă evidența înainte de reaplicare.
6. Zilele salvate apar automat în **Pontaj**, cu fundal verde și codul tipului. Celulele sunt protejate; modificarea se face din modulul de concedii. Un cod scris manual într-o schiță de pontaj rămâne text și nu înregistrează automat o absență.

Sumarul anual afișează separat numărul de zile de odihnă, medical și speciale înregistrate, plus dreptul contractual **curent** din fișa angajatului. Nu calculează soldul legal, reportările sau indemnizațiile. Înregistrările care traversează anul se văd în anii corespunzători.

Un tip poate fi editat sau dezactivat din tabelul **Tipuri de absență**. Dezactivarea nu elimină zilele existente. Codul nu se redenumește; categoria unui tip folosit nu se schimbă, pentru a păstra sensul istoricului. Schema și limitele sunt în [D020](decizii/D020-concedii-absente.md).

## Grila lunară · schiță

1. Deschide **Pontaj**, alege anul și luna; opțional filtrează secția și departamentul.
2. Introdu o valoare în celulă. Nomenclatorul sugerează schimburile active; alte coduri se păstrează ca text, fără interpretare automată.
3. Folosește săgeți, Enter sau Tab pentru deplasare; Shift+săgeți ori Shift+clic pentru selectarea unui dreptunghi.
4. **Completează selecția** aplică valoarea din câmpul alăturat; câmpul gol golește selecția. Ctrl+C/Ctrl+V copiază și lipește dreptunghiuri TSV compatibile cu Excel. O lipire care depășește luna sau rândurile disponibile este refuzată integral.
5. Apasă **Salvează schița**. Modificările sunt păstrate în baza criptată și rămân după repornire. La navigare sau schimbarea perioadei cu celule nesalvate apare confirmarea abandonării.

Weekendurile sunt gri, sărbătorile sunt marcate distinct, iar antetul și identificarea angajaților rămân fixe la derulare. Angajații inactivi apar în lunile în care au deja schițe salvate. Dacă o altă instanță a modificat aceeași celulă, salvarea lotului este refuzată; copiază modificările locale și reîncarcă luna înainte de reaplicare.

**Totalurile din dreapta nu sunt încă calculate.** Schița nu reprezintă calcul de salarizare sau raportul final. Semnificațiile C/FP/cev, coloanele de total și legătura dintre rândurile GARDĂ și persoană/contract așteaptă specificațiile (F01). Datele reale din capturile primite nu au fost copiate în repo.

Istoricul contractual, conturile și rapoartele finale rămân în taskurile lor. Modificările fluxurilor disponibile sunt jurnalizate cu actorul provizoriu `local`, până la F15.

## Verificări pentru dezvoltatori

- `gofmt -l -w .`, `go vet ./...`, `go test ./...`: validări, tranzacții, persistență pe SQLite criptat, păstrarea pontajului, import demo atomic și refuzul suprascrierii datelor.
- `./genereaza-date.ps1`: regenerează accesul SQL cu sqlc v1.31.1 fixat.
- `./build.ps1`: construiește și verifică versiunea din proprietățile executabilului.
- `node verificari/interfata.mjs`, după build-ul frontend: Edge headless deja instalat, profil temporar, fără biblioteci suplimentare. Verifică angajați/demo/setări, ierarhia secție/departament, atribuire și filtre, calendar 2026 și import repetat, cele șase schimburi inițiale, concedii/medical/tipuri speciale/intervale/sumar/anulare/conflicte și protecția în pontaj, plus grila 50×31, tastatură, TSV, luni/ani bisecți, salvare/conflict, lățime 1024px și EN. Puntea Wails este simulată numai în acest test; backendul real, inclusiv lotul 50×31 și conflictele, este testat separat prin Go pe SQLite criptat. Capturi în `build/bin/interfata-demo.png`, `calendar-2026.png`, `concedii-absente.png` și `pontaj-schita.png`.

Versiunea anterioară 0.13.0 a trecut o probă de pornire ascunsă de 8 secunde, cu baza demo și seiful prezente. Pentru 0.17.1 verificările de interacțiune sunt cele din Edge cu puntea simulată; fereastra desktop nu a fost testată interactiv. Corecția [D022](decizii/D022-revizii-concurenta.md) păstrează detectarea conflictelor chiar când două salvări au același timestamp.
