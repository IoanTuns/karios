# D006 · Reguli de lucru pentru agenți

Stare: acceptată · Data: 2026-10-24 · Decisă de: Ioan Tuns · Taskuri: F02-T07, F02-T08

## Context

Mai mulți agenți AI (Claude, Codex, Cursor…) și oameni vor lucra în același repo. Fără reguli comune, agenții pot lucra la același task, pot lua decizii pe care ceilalți nu le văd, pot inventa reguli de plată sau pot publica date personale. Un agent nu ține minte nimic de la o sesiune la alta; tot ce știe despre proiect trebuie să fie în repo.

Regulile complete sunt în [AGENTS.md](../../AGENTS.md). Mai jos e motivul și mecanismul fiecăreia.

## Decizie

### 1. Un singur fișier de reguli: `AGENTS.md`

- **De ce:** `AGENTS.md` e numele citit de majoritatea uneltelor. Reguli separate pe unealtă (`CLAUDE.md`, `.cursorrules`…) ajung să difere.
- **Cum:** `CLAUDE.md` conține doar `@AGENTS.md`. O unealtă care cere alt nume primește la fel un fișier care trimite la `AGENTS.md`.

### 2. Orice idee, decizie și acțiune se documentează

- **De ce:** ce nu e scris nu există pentru următorul agent. Ideile respinse rămân scrise, ca să nu fie propuse din nou fără motiv.
- **Cum:** `idei.md`, `decizii/`, `functionalitati/`, `jurnal.md`, în același commit cu schimbarea.

### 3. Agenții propun, oamenii decid

- **De ce:** un agent poate greși cu încredere, iar deciziile despre plată au consecințe legale. Răspunderea rămâne la un om.
- **Cum:**
  - Un agent scrie decizii doar cu starea `propusă`.
  - Când un om spune explicit că acceptă sau respinge (în conversație sau în scris), agentul poate schimba starea în numele lui și completează `Decisă de: <nume>`, ca la D007.
  - Un agent nu schimbă niciodată o decizie `acceptată`; scrie una nouă care o înlocuiește.

### 4. Regulile de plată vin doar de la contabilitate

- **De ce:** procentele, sporurile și interpretarea Codului Muncii țin de contractul colectiv și de practica firmei. Un agent le poate „ghici” plauzibil, dar greșit, iar eroarea ajunge în salarii.
- **Cum:** dacă o regulă lipsește, taskul devine `blocat`, iar întrebarea se adaugă la F01. Valorile provizorii din cod sunt marcate ca provizorii.

### 5. Testele de referință sunt protejate

- **De ce:** sunt singura dovadă că aplicația calculează la fel ca acum. Un agent care „repară” un test ca să treacă ascunde exact eroarea pe care testul trebuia s-o prindă.
- **Cum:** cazurile confirmate de contabilitate stau într-un fișier separat (de exemplu `internal/calcul/referinta_test.go`), cu un comentariu în antet. O modificare a lor cere aprobarea unui om, notată în jurnal.

### 6. Un task, un responsabil

- **De ce:** doi agenți pe același task produc două variante incompatibile și muncă pierdută.
- **Cum:** coloana `Responsabil`. Agentul își trece numele și starea `în lucru` înainte de a începe.
- **Limită (nou):** pe ramuri separate, rezervarea făcută de un agent nu e vizibilă celorlalți până la merge. **Până există un remote comun, taskurile le alocă un om** (completează `Responsabil` pe `main`). Cu GitHub, alocarea se poate face prin issues (I003).

### 7. O ramură per task; merge-ul îl face un om

- **De ce:** `main` rămâne mereu într-o stare verificată de un om; munca unui agent se poate arunca fără să afecteze restul.
- **Cum:** ramura se numește după task (`F04-T03-sporuri`). Pentru agenți care rulează în paralel pe același calculator: câte un `git worktree` separat. Push-ul și merge-ul în `main` cer aprobare.

### 8. Commit-uri mici, identificabile

- **Cum:** titlul începe cu ID-ul taskului (`F04-T03: …`); mesajul e în română; semnătura `Co-Authored-By` arată ce agent a lucrat. Actualizările de versiuni au commit-ul lor (D007).

### 9. Predarea la sfârșitul sesiunii

- **De ce:** următorul agent pornește fără context.
- **Cum:** intrare în jurnal (ce s-a făcut, ce a rămas) și pasul următor în coloana `Note` a taskului. Un task neterminat rămâne `în lucru` cu responsabilul lui sau e eliberat explicit.

### 10. Date personale și secrete

- **De ce:** centralizatoarele reale conțin date personale ale angajaților (GDPR). Un repo se copiază ușor și ce a intrat în istoria git e greu de șters.
- **Cum:** datele reale stau în `date-reale/`, ignorat de git; în teste intră doar date anonimizate (nume și mărci înlocuite). Secretele (tokenuri, parole) nu intră niciodată în repo. Un agent nu trimite date reale către servicii externe.

### 11. Calitate și cinste

- **Cum:** înainte de commit rulează `gofmt`, `go vet`, `go test`. Agentul spune ce a verificat și ce nu. Problemele observate în afara taskului devin idei sau taskuri, nu intră în același commit.

## Variante respinse

- **Reguli separate pe unealtă** — ar ajunge să difere.
- **Idei notate doar în conversații** — se pierd; alt agent nu le vede.
- **Agenți care decid singuri** — lipsă de răspundere pentru decizii cu efect în salarii.
- **Lucru direct pe `main`** — o greșeală a unui agent ajunge imediat la toți.

## Consecințe

- Puțină muncă în plus la fiecare task, în schimbul unui istoric complet și al lucrului sigur în paralel.
- **Conflicte frecvente în `jurnal.md` (nou):** fiecare ramură adaugă rânduri sus în același tabel. La merge, conflictul se rezolvă păstrând toate rândurile, în ordinea datei. Dacă devine supărător, jurnalul se poate împărți pe luni sau pe fișiere separate per intrare.
- Agenții au nevoie ca omul să aloce taskuri și să facă merge-uri, ceea ce poate încetini lucrul până la configurarea unui remote.
