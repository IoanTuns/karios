# D010 · Migrările bazei de date cu goose

Stare: acceptată · Data: 2026-10-24 · Decisă de: Ioan Tuns · Taskuri: F06-T01, F06-T06

## Context

D009 cere ca fiecare schimbare a bazei de date să aibă versiunea ei, cu istoric: ce migrare, când, din ce versiune a aplicației. Documentul de design numea atât `goose`, cât și `PRAGMA user_version`; trebuie ales unul.

## Decizie

- **Unealta:** `github.com/pressly/goose/v3` (v3.29.0, licență MIT), prin interfața `Provider`. Migrările sunt fișiere `.sql` incluse în `.exe` (`embed.FS`) din `internal/bd/migrari/`.
- **Numerotare:** `0001_descriere.sql`, `0002_…`, secvențial. O migrare = o singură schimbare.
- **Doar înainte:** migrările nu au secțiune `Down`. O schimbare greșită se corectează cu o migrare nouă; revenirea la o stare veche se face din backup-ul automat dinaintea migrării (D002).
- **Istoricul:** goose ține în `goose_db_version` versiunea migrării și data. Aplicația adaugă tabelul `istoric_migrari` cu versiunea aplicației care a aplicat fiecare migrare.
- **Protecția la downgrade:** la pornire, dacă versiunea bazei e mai mare decât ultima migrare cunoscută de aplicație, aplicația refuză deschiderea, cu un mesaj tradus.
- **„Salvată”** (D009: „o migrare salvată nu se mai modifică”) înseamnă **existentă în `main`**. Pe o ramură neadusă încă în `main`, o migrare se poate redenumi, de exemplu când altă ramură a folosit deja același număr. Hook-ul verifică exact asta.

**Verificat în prototip** (`prototipuri/sqlite-criptat/migrari_test.go`, pe baza criptată cu `xts`): aplicarea în ordine, istoricul cu dată, versiunea aplicației lângă migrare, nicio migrare la a doua pornire, detectarea downgrade-ului.

## Variante respinse

- **`PRAGMA user_version`** — ține un singur număr, fără dată și fără istoric; nu îndeplinește D009.
- **Migrări cu secțiune `Down`** — cod rar folosit și rar testat; backup-ul automat e o cale de întoarcere mai sigură.
- **Versiuni de tip dată-oră** (`20261024153000_…`) — evită coliziunile între ramuri, dar pot fi aplicate în altă ordine decât au fost scrise; redenumirea înainte de merge e mai simplă.

## Consecințe

- Dependențe noi în aplicație: goose și trei biblioteci mici (`interpolate`, `go-retry`, `multierr`), toate MIT, Apache 2.0 sau BSD.
- La merge, cine aduce a doua o ramură cu o migrare cu același număr o redenumește.
