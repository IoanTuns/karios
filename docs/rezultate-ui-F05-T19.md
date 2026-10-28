# F05-T19 · Rezultate UI 0.28.0

26 octombrie 2026 · Codex · Ramură `F05-T19-ergonomie-pontaj` · [lista de lucru](lista-lucru-ui.md).

Implementate I045–I048, în ordinea acceptată prin D035. Nu sunt introduse biblioteci, migrări, reguli contabile sau modificări ale rapoartelor generate.

## Verificări

| Criteriu | Rezultat |
| --- | --- |
| Spațiul grilei la 1024×640, 50×31 | Grila începe la y=248, are 342 px și afișează 7 rânduri complete cu antetul; pagina are exact 1024×640. În audit începea la y=629, cu înălțime 130 px |
| T / Ore norm | 9 ore manuale sâmbătă, plus 12,5 ore în zile normale: T=21,5 fără avertisment contractual. Filtrul de avertismente nu include acest rând; lucrul după G24 rămâne evidențiat și filtrabil |
| Context la navigare | Pontaj → Rapoarte → Pontaj păstrează 11/2028, căutarea TEST-0, departamentul, 25 rânduri filtrate, selecția și derularea x=600/y=180 |
| Modificări nesalvate | Refuzul confirmării păstrează schița și ecranul. Acceptarea renunțării permite navigarea și reîncarcă ulterior valoarea salvată, nu cele 5 ore abandonate |
| Tastatură și panouri | F2 deschide Completare și focalizează orele manuale; Escape închide Ajutorul și întoarce focusul pe buton. Fluxurile existente de selecție, Undo/Redo și TSV trec |
| Dialog de raport lung | Dialog de 592 px; corp derulant de 439 px, cu 3062 px de conținut. Antetul și butonul de salvare sunt vizibile și accesibile la începutul și sfârșitul formularului |
| Salvarea regulilor | Butonul rămâne vizibil și accesibil după derularea paginii până la sfârșit |
| Contrast în grilă | Ore pe selecție 6,43:1; pe weekend 6,24:1; valori necalculate 5,64:1. Peste pragul 4,5:1 verificat pentru aceste perechi; nu este un audit WCAG complet |

`npm run build`, `gofmt -l -w .`, `go vet ./...` și `go test ./...` au trecut. Verificarea completă `node verificari/interfata.mjs` a trecut în Edge, inclusiv conflictele, RO/EN și fluxurile existente de angajați, absențe, configurare și rapoarte. Scriptul creează acum director separat pentru fiecare rulare.

Artefactele locale ale verificării finale sunt în `build/bin/verificare-ui-2026-10-26T12-57-19-658Z/`: [grilă compactă](../build/bin/verificare-ui-2026-10-26T12-57-19-658Z/pontaj-compact-1024.png), [editor raport](../build/bin/verificare-ui-2026-10-26T12-57-19-658Z/definire-rapoarte.png), [măsurători](../build/bin/verificare-ui-2026-10-26T12-57-19-658Z/masuratori-ui.json). Valorile de mai sus sunt păstrate în documentație chiar dacă artefactele generate vor fi curățate.

## Limite și predare

Proba UI folosește date fictive și punte Wails simulată; backendul real are teste Go separate pe SQLite criptat. Nu s-au deschis bazele operaționale, nu s-a probat tipărirea fizică și nu s-a efectuat o sesiune cu un operator. Contextul se păstrează numai până la închiderea aplicației. Structura datelor și referințele rapoartelor sunt păstrate.

Următorul task este F05-T20: motive distincte pentru „—” și validare localizată a configurării. F05-T21–T23 sunt în lista de lucru, nepreluate.

Buildul final din commitul `4016532` este [Pontix 0.28.0](../build/bin/pontix-0.28.0-20261026-145951-001bf64d73dc445b8ad81a2b0f2e9bb6/pontix.exe), construit cu `build.ps1` într-un director nou. Produsul și versiunea Windows sunt verificate; manifestul este complet, SHA256 `C1C4B87D08E321AAF43F8020EB48AAD3242ECF76B8D40B6086A00321A49B6CEA`. Marcajul `4016532-modificat` reflectă exclusiv documentul PR modificat în paralel și necuprins în task; sursele aplicației corespund commitului. Ghidurile și lista UI sunt copiate lângă executabil. Buildul nu transferă automat datele din alte versiuni și executabilul nou nu a fost pornit pe date operaționale.
