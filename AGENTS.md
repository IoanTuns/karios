# Reguli de lucru

Se aplică oricărui agent AI (Claude, Codex, Cursor, Copilot…) și oricărui om care lucrează în acest repo. Proiectul: Pontix, aplicație desktop portabilă pentru pontaj (Go + Wails + React + SQLite). Starea proiectului e în [docs/](docs/README.md).

## 1. Totul se documentează

Orice idee, decizie și acțiune se scrie în `docs/`, în același commit cu schimbarea. Ce nu e scris nu există pentru ceilalți agenți.

| Ce | Unde | Regulă |
| --- | --- | --- |
| Idee | [docs/idei.md](docs/idei.md) | Orice idee, chiar dacă nu se face acum. Cele respinse sau amânate rămân, cu motivul. |
| Decizie | [docs/decizii/](docs/decizii/) | Orice alegere între variante: tehnologie, bibliotecă, regulă de calcul, format, proces. |
| Muncă | [docs/functionalitati/](docs/functionalitati/) | Un task există **înainte** de a începe lucrul. |
| Acțiune | [docs/jurnal.md](docs/jurnal.md) | Orice acțiune: cod, document, instalare, discuție, încercare eșuată. Cu coloana `Cine`. |

Faptele despre proiect stau în repo, nu în memoria privată a unui agent, pe care ceilalți nu o văd. Regulile complete ale documentației sunt în [docs/README.md](docs/README.md).

## 2. Ce poate decide un agent

- **Deciziile:** un agent doar propune (`propusă`). Un om le acceptă sau le respinge; când omul o spune explicit, agentul poate schimba starea în numele lui, cu `Decisă de: <nume>`.
- **Regulile de plată nu se inventează.** Procentele, sporurile și interpretările legale vin doar de la contabilitate. Dacă lipsesc, taskul devine `blocat`, iar întrebarea se adaugă la F01.
- **Testele de referință sunt protejate.** Cazurile confirmate de contabilitate nu se schimbă fără aprobare. Dacă un test pică, se corectează codul, nu testul.
- **Cer aprobare:** push, merge în `main`, trimiterea sau publicarea de documente, instalarea de programe pe sistem, ștergerea de date sau de fișiere pe care agentul nu le-a creat.

## 3. Lucrul în paralel

- **Un task are un singur responsabil.** Înainte de a începe, agentul își trece numele în coloana `Responsabil` și starea `în lucru`. Nu preia un task care are deja responsabil. Până există un remote comun, taskurile le alocă un om, pe `main`.
- **O ramură per task**, numită după ID: `F04-T03-sporuri`. Ideal într-un worktree separat. Nimeni nu scrie direct pe `main`; merge-ul îl face un om.
- **Commit-uri mici**, cu ID-ul taskului în titlu (`F04-T03: …`) și cu semnătura agentului (`Co-Authored-By`).
- **Predarea:** la sfârșitul sesiunii, agentul scrie în jurnal ce a făcut și, în notele taskului, pasul următor, ca altcineva să poată continua fără context.

## 4. Date și securitate

- **Datele reale ale angajaților nu intră în repo** (date personale, GDPR). Stau în `date-reale/`, ignorat de git. În teste intră doar date anonimizate.
- **Niciun secret în repo:** tokenuri, parole, chei. Fișierele cu secrete sunt ignorate de git.
- **O bibliotecă nouă** cere o decizie scrisă, cu licența verificată.

## 5. Calitatea lucrului

- **Raportare cinstită:** ce teste s-au rulat și rezultatul; ce nu s-a verificat se spune.
- **Doar ce cere taskul.** O refactorizare sau o problemă observată pe parcurs devine o idee sau un task nou, nu intră în același commit.
- **Calculul are teste.** Orice schimbare în `internal/calcul` vine cu teste.
- Înainte de commit: `gofmt -l -w .`, `go vet ./...`, `go test ./...`.

## 6. Versiuni: mereu ultima versiune stabilă

Se folosește mereu cea mai nouă versiune stabilă pentru tot: Go, Wails, Node, bibliotecile Go și npm, uneltele (`sqlc`, migrări), fonturile, modelele LLM ([D007](docs/decizii/D007-versiuni-stabile.md)).

- **Stabilă** = lansare finală. Fără alpha, beta, rc, preview, nightly.
- **La fiecare modificare**, agentul verifică dependențele proiectului și le actualizează dacă a apărut o versiune stabilă nouă: Go din <https://go.dev/VERSION?m=text>, modulele cu `go list -m -u all`, pachetele npm cu `npm outdated`.
- **Versiunile se fixează exact** (`go.mod` cu `toolchain`, `go.sum`, `package-lock.json`), fără intervale. „Ultima” înseamnă ultima la momentul modificării, ca build-ul să dea același rezultat pe orice calculator.
- **Actualizarea are commit-ul ei**, în același task, cu testele rulate și o intrare în jurnal: ce s-a actualizat, de la ce versiune la ce versiune.
- **Versiune majoră nouă** (schimbări incompatibile): agentul citește notele de lansare și verifică licența. Dacă se schimbă licența → decizie nouă și aprobare.
- **Dacă ultima versiune strică ceva** ce nu se poate repara în task, se rămâne temporar la cea anterioară, cu motivul în jurnal și un task pentru actualizare.

## 7. Traduceri

Niciun text afișat utilizatorului nu se scrie direct în cod ([D004](docs/decizii/D004-traduceri.md)).

- Text nou → cheie nouă în `internal/i18n/locales/ro.json` și `en.json`.
- Erori pentru utilizator → `calcul.Eroare{Cod, Param}` sau un tip asemănător, nu `fmt.Errorf` cu text.
- Denumiri introduse de utilizator → `i18n.Text`.

## 8. Cod și mediu

- Identificatorii și comentariile sunt în română, fără diacritice în identificatori.
- Logica de calcul stă în `internal/calcul` și nu depinde de bază de date sau de interfață.
- Versiunea Go a proiectului e cea din `toolchain` în `go.mod` (secțiunea 6).
- Terminațiile de linie sunt LF (`.gitattributes`).

## 9. Limba

- Documentația, comentariile și mesajele de commit: română.
- Conversația cu utilizatorul: în limba în care scrie el.
- Textele din aplicație: prin cataloagele `i18n`.
