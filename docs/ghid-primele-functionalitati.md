# Funcționalitățile disponibile · Pontix 0.13.0

## Pornire

Executabilul se construiește cu `./build.ps1` și se găsește în `build/bin/pontix.exe`. Pornește-l dintr-un folder în care poți scrie; datele sunt în folderul `karios/` de lângă executabil.

Livrarea verificată în această sesiune este **`build/bin/pontix-0.13.0/pontix.exe`**, cu date demo deja regenerate într-un folder separat, la cererea utilizatorului. Baza veche a fost păstrată în `build/bin/karios/`; seiful ei a fost eliminat de vechiul build cu `-clean`, deci nu se reutilizează cu o cheie nouă. Incidentul și corecția sunt în [D017](decizii/D017-build-fara-stergerea-datelor.md).

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
- Filtrează după secție și activ/inactiv.
- **Adaugă angajat**: nume, prenume și marcă sunt obligatorii; marca este text și păstrează zerourile inițiale.
- Contractul este implicit cu normă întreagă și preia norma globală. Pentru norma parțială completează orele individuale pe zi, inclusiv fracții care reprezintă minute întregi (de exemplu 4,5 ore = 4 h 30 min).
- Concediul anual este cel din contract. Câmp gol înseamnă necompletat, nu zero. Nu se calculează încă soldul concediului.
- **Editează** modifică datele; la anularea unui formular modificat se cere confirmarea abandonării.
- **Dezactivează / Reactivează** schimbă starea fără a șterge persoana sau pontajul existent.

Datele migrate cu prenume necompletat trebuie verificate manual înainte de salvarea formularului. Nu se deduce automat prenumele din numele vechi.

## Setări

- Modifică norma întreagă în ore/zi (implicit 7). Modificarea se vede imediat la angajații cu normă întreagă; normele parțiale rămân individuale.
- Adaugă secții, apoi atribuie-le din formularul angajatului. Denumirea se păstrează ca text traductibil în limba interfeței.

## Sărbători

La actualizarea aplicației, migrarea 0004 adaugă lista națională pentru 2026, fără a suprascrie înregistrări existente: 17 sărbători nominale în 16 date distincte; la 1 iunie coincid Ziua Copilului și a doua zi de Rusalii. Setul folosește datele mobile ortodoxe, indicate în denumiri și în ecran. Sursele și limitele calendarului comun sunt în [D014](decizii/D014-sarbatori-2026.md).

Alege **Sărbători**, anul și **Arată anul**. Poți adăuga sau edita o dată și denumirile RO/EN. **Completează lista legală 2026** adaugă numai date lipsă și păstrează denumirile editate. Nu se introduc automat zile-punte sau reguli de plată.

## Schimburi

Din **Schimburi → Adaugă schimb**, introdu codul (de exemplu N12), denumirea, ora de început, durata și pauza în minute și tipul activității. Nu există intervale presupuse din cod. De exemplu, dacă alegi o tură de 12 ore, durata este 720 minute; ora de început rămâne cea stabilită în organizație. Codurile se normalizează la majuscule și nu pot fi duplicate, nici între schimburile inactive. Dezactivarea păstrează referințele din pontaj.

## Grila lunară · schiță

1. Deschide **Pontaj**, alege anul și luna; opțional filtrează secția.
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
- `node verificari/interfata.mjs`, după build-ul frontend: Edge headless deja instalat, profil temporar, fără biblioteci suplimentare. Verifică angajați/demo/setări, calendar 2026 și import repetat, nomenclator N12 și grilă 50×31, tastatură, TSV, luni/ani bisecți, salvare/conflict, lățime 1024px și EN. Puntea Wails este simulată numai în acest test; backendul real, inclusiv lotul 50×31 și conflictele, este testat separat prin Go pe SQLite criptat. Capturi în `build/bin/interfata-demo.png`, `calendar-2026.png` și `pontaj-schita.png`.

Noul executabil a trecut o probă de pornire ascunsă de 8 secunde, cu baza demo și seiful prezente. Fereastra desktop nu a fost testată interactiv în această sesiune; verificările de interacțiune sunt cele din Edge cu puntea simulată.
