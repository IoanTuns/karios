# D103 · Limba interfeței și a rapoartelor, reținute local

Stare: acceptată · Data: 2026-11-12 · Taskuri: F03-T05, F03-T08, F05-T04 · Cerută de: Ioan Tuns, 2026-11-12 („Poți să le implementezi pe toate și să construiești o versiune nouă?”) · Decisă de: Ioan Tuns, 2026-11-12 („accept toate”)

Implementare: implementată din 0.92.95; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

## Context

Limba interfeței se alegea din bara laterală, dar la fiecare pornire aplicația revenea la română. F03-T05 cerea setările `limba_interfata` și `limba_rapoarte` în `config.toml`. F05-T04 nota că, după D008, setările stau în baza criptată.

Limba trebuie însă cunoscută înainte de deschiderea bazei: ecranul de deblocare și dialogurile de dinainte sunt deja afișate în ea. În plus, nu e o dată personală.

## Decizie

- **Unde se rețin.** Ambele alegeri stau în stocarea locală a WebView-ului, la fel ca tema (luminoasă sau întunecată). Folderul acesteia este `karios/webview`, lângă `.exe`. Alegerea călătorește cu stickul și e aceeași pentru toți utilizatorii instalării.
- **Limba interfeței.** La pornire se folosește ultima limbă aleasă, dacă e încă disponibilă; altfel, limba implicită (română). Se schimbă din bara laterală sau din **Setări › General › Limbă**.
- **Limba implicită a rapoartelor.** Valoarea implicită este **La fel ca interfața**, adică comportamentul de până acum. Se poate alege și **Română** sau **English**. Ecranul Rapoarte pornește cu această limbă, iar fiecare raport poate fi generat în continuare în altă limbă. Rapoartele flexibile o folosesc la fel.
- **Textele netraduse.** Fila **Limbă** arată, pentru fiecare limbă în afara celei implicite, câte texte nu sunt traduse și care sunt. Lista vine din `i18n.Lipsa`, iar aceste texte se afișează în limba implicită.
- **Fără `config.toml`.** Nu se adaugă un fișier de configurare doar pentru limbă.

## Variante respinse

- **În baza criptată (`setari`).** Limba nu ar fi cunoscută înainte de deblocare, iar ecranele de dinainte ar rămâne în limba implicită.
- **Fișier `config.toml` lângă `.exe`.** Ar fi un al doilea loc pentru preferințe, pe lângă stocarea WebView folosită deja pentru temă, plus o bibliotecă TOML (decizie și licență) pentru două valori.
- **Limba pe utilizator.** Nu există o cerere în acest sens, iar limba trebuie aleasă înainte de autentificare.

## Consecințe

- Cine nu alege nimic nu vede nicio schimbare.
- Ștergerea folderului `karios/webview` readuce limba implicită, la fel ca tema.
- Importul traducerilor corectate din folderul de lângă `.exe` (D004) rămâne separat: F03-T11.
