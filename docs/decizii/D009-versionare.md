# D009 · Versionarea aplicației, a rapoartelor și a bazei de date

Stare: acceptată · Data: 2026-10-24 · Decisă de: Ioan Tuns · Taskuri: F02-T10, F06-T01, F09-T06, F10-T05

## Context

Cerință (Ioan Tuns, 2026-10-24):

1. Fiecare modificare de cod generează o versiune nouă.
2. Fiecare raport are versiunea lui.
3. Fiecare modificare a bazei de date are versiunea ei.

Scopul: pentru orice `.exe`, orice raport tipărit și orice fișier de bază de date, să se poată spune exact din ce versiune vine.

## Decizie

### 1. Versiunea aplicației

- **Format:** `MAJOR.MINOR.PATCH`, cu o singură sursă: `const Aplicatie` din `internal/versiune/versiune.go`.
  - `PATCH` — corecturi, actualizări de dependențe (D007), schimbări interne.
  - `MINOR` — funcționalitate nouă.
  - `MAJOR` — schimbare incompatibilă (de exemplu un fișier `karios/` pe care versiunile vechi nu-l mai pot deschide).
  - Până la prima versiune utilizabilă (V1) se rămâne la `0.x`; V1 = `1.0.0`.
- **Orice commit care schimbă codul** (Go, SQL, TypeScript, cataloagele de traduceri, `go.mod`) crește versiunea **în același commit** și adaugă o secțiune în [docs/versiuni.md](../versiuni.md). Commit-urile doar cu documentație nu schimbă versiunea.
- **Ramuri în paralel:** dacă două ramuri ajung la același număr, cea adusă a doua în `main` primește la merge numărul următor.
- **Pe `main`**, fiecare versiune primește un tag git `vX.Y.Z` la merge (făcut de un om, D006).
- **În `.exe`:** versiunea și commit-ul (`-ldflags`) se văd în ecranul „Despre”, în rapoarte și în jurnalul aplicației.
- **Verificare automată:** hook-ul `.githooks/pre-commit` refuză un commit cu cod schimbat fără versiune nouă; testul `TestIstoric` verifică intrarea din `docs/versiuni.md`.

### 2. Versiunea fiecărui raport

- Fiecare raport (centralizator, fișă individuală, grafic de schimburi, export pentru salarizare) are **propria versiune**, un număr întreg (`const Versiune = 3` în pachetul raportului), independentă de versiunea aplicației.
- **Versiunea raportului crește** când se schimbă ce conține, cum arată sau cum se calculează ce afișează.
- **Fiecare raport generat poartă în subsol:** numele și versiunea raportului, versiunea aplicației, setul de reguli de calcul folosit, data generării și contul care l-a generat.
- **Verificare automată:** fiecare raport are un test cu fișier de referință (golden file). Dacă rezultatul se schimbă, testul pică până când se actualizează referința *și* crește versiunea raportului.
- Schimbările de raport se trec în `docs/versiuni.md`, în secțiunea versiunii aplicației care le aduce (de exemplu „Centralizator v3: coloana ZNL”).

### 3. Versiunea fiecărei schimbări a bazei de date

- **Orice schimbare a bazei e o migrare nouă**, un fișier numerotat în `internal/bd/migrari/` (`0007_adauga_sectii.sql`). O migrare = o singură schimbare.
- **O migrare salvată nu se mai modifică.** O greșeală se corectează cu o migrare nouă. Hook-ul refuză modificarea sau ștergerea unei migrări existente.
- **Baza ține istoricul:** ce migrări s-au aplicat, când și din ce versiune a aplicației.
- **Înainte de orice migrare:** backup automat cu `VACUUM INTO` (D002).
- **Protecție la downgrade:** o versiune veche a aplicației refuză să deschidă o bază migrată de o versiune mai nouă, cu un mesaj tradus, în loc s-o strice.
- **Unealta:** se alege în F06-T01. Cerința de istoric cu dată favorizează `goose`, care ține un tabel cu fiecare migrare aplicată, față de `PRAGMA user_version`, care ține doar un număr.

## Variante respinse

- **Versiune doar la lansări** — nu respectă cerința; între lansări n-ar fi clar din ce cod vine un `.exe`.
- **Versiune calculată din git** (numărul de commit-uri) — nu poate fi citită din cod fără git și nu spune nimic despre tipul schimbării.
- **Rapoartele versionate odată cu aplicația** — n-ar arăta care raport s-a schimbat efectiv.
- **Modificarea migrărilor existente** — bazele deja migrate n-ar mai corespunde schemei din cod.

## Consecințe

- Fiecare commit de cod are puțină muncă în plus: versiunea și o intrare în istoric.
- Hook-ul trebuie activat o dată pe fiecare clonă: `git config core.hooksPath .githooks`.
- Un raport tipărit sau un fișier de bază de date se poate lega oricând de codul exact care l-a produs.
