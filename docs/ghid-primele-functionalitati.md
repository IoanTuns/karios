# Primele funcționalități · Pontix 0.10.0

## Pornire

Executabilul se construiește cu `./build.ps1` și se găsește în `build/bin/pontix.exe`. Pornește-l dintr-un folder în care poți scrie; datele sunt în folderul `karios/` de lângă executabil.

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

Istoricul contractual, conturile, pontajul lunar, rapoartele și administrarea sărbătorilor rămân în taskurile lor. Modificările acestor prime fluxuri sunt jurnalizate cu actorul provizoriu `local`, până la F15.

## Verificări pentru dezvoltatori

- `gofmt -l -w .`, `go vet ./...`, `go test ./...`: validări, tranzacții, persistență pe SQLite criptat, păstrarea pontajului, import demo atomic și refuzul suprascrierii datelor.
- `./genereaza-date.ps1`: regenerează accesul SQL cu sqlc v1.31.1 fixat.
- `./build.ps1`: construiește și verifică versiunea din proprietățile executabilului.
- `node verificari/interfata.mjs`, după build-ul frontend: Edge headless instalat deja pe sistem, profil temporar, fără biblioteci suplimentare. Verifică demo, căutare, filtre, formulare, erori, norme, secții, modificări nesalvate, reîncărcare, lățime 1024px și limba EN. Puntea Wails este simulată numai în acest test; backendul real este testat separat prin Go. Captura este generată în `build/bin/interfata-demo.png`.

Rularea interactivă a executabilului Windows nu a fost verificată în această sesiune.
