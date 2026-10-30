# Builduri separate

Din checkoutul versiunii dorite rulează:

```powershell
$buildNou = .\build.ps1
$buildNou.Executabil
```

Fiecare invocare creează un director nou:

```text
build/bin/pontix-X.Y.Z-AAAALLZZ-HHMMSS-<identificator>/
  pontix.exe
  build.json
  manual-utilizare-ro-X.Y.Z.md
  manual-utilizare-ro-X.Y.Z.html
  manual-utilizare-ro-X.Y.Z.pdf
```

Versiunea, data/ora și identificatorul sunt generate automat. Două builduri ale aceleiași versiuni au directoare diferite. Pentru alt nume al executabilului, folosește `-NumeExecutabil pontix-proba.exe`; directorul rămâne unic.

`build.json` înregistrează versiunea, commitul, momentul, numele executabilului și SHA256, plus ediția, identitatea și hashurile celor trei manuale. Din 0.88.11, fiecare build regenerează manualele și le pune lângă executabil, cu aceeași versiune și același commit. Un build complet are `stare: complet`; un eșec are `stare: esuat`. Sufixul `-modificat` al commitului indică surse cu modificări necomise. Proprietățile Windows ale executabilului și manualele sunt verificate înainte de marcarea buildului ca reușit.

Înainte de fiecare versiune, revizuiește `docs/manual-utilizare-ro.md` și actualizează `versiune_revizuita` din `docs/manual/manual.json`; altfel buildul se oprește înainte de compilare. Sunt necesare Python și Microsoft Edge existente pe calculatorul de construire. Nu se instalează nimic automat. Detalii și previzualizare: [manual/README.md](manual/README.md). Aplicația construită nu necesită Python.

Temporarele Wails folosesc `.temporar/` din directorul nou și sunt eliminate la încheiere. Fișierele blocate sau legăturile se păstrează cu avertizare. TEMP/TMP și directorul apelantului sunt restaurate inclusiv la eșec. Dependențele sunt cele fixate în proiect; buildul nu rulează actualizări sau tidy și nu utilizează `-clean`.

## Legăturile dintre interfață și Go

După orice schimbare a metodelor legate (`App`) sau a tipurilor lor, se regenerează `frontend/wailsjs`:

```powershell
wails generate module
```

Ieșirea generatorului e varianta de referință; nu se editează manual (F02-T15). Generatorul include doar tipurile care apar în semnăturile metodelor: un tip trimis numai prin eveniment (de ex. `RezultatSarcinaAI`, evenimentul `ai:sarcina_terminata`) se definește în frontend, lângă codul care îl folosește. `wails build` regenerează și el legăturile.

## Datele aplicației

Un director nou reprezintă o copie portabilă nouă. Buildul nu transferă date între versiuni. Pentru continuarea unei baze existente, închide aplicația și copiază întregul folder `karios/` lângă executabilul nou; nu doar fișierul bazei. Păstrează copia originală până la verificarea noii versiuni.

Curățarea F02-T12 păstrează directoarele vechi care conțin date, checkouturile și uneltele locale. Fișierele generate necesare construirii din surse (legături Wails, resurse Windows, referințe ale rapoartelor) rămân în Git.

## Verificarea scriptului

```powershell
$buildNou = .\build.ps1
.\verificari\build.ps1 -ExecutabilReferinta $buildNou.Executabil
```

Verificarea folosește executabilul real deja construit și simulează compilatorul pentru două invocări și un eșec. Verifică directoare diferite, manifest/hash, temporare, restaurarea mediului, refuzul căilor invalide și păstrarea unui fișier existent blocat. Elimină numai probele proprii. Testele Go obligatorii rămân `gofmt -l -w .`, `go vet ./...` și `go test ./...`.
