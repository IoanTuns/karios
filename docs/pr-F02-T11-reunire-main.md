Reunește în `main` aplicația Pontix 0.22.2 și istoricul tuturor ramurilor locale și remote. Integrarea include administrarea angajaților și normelor, calendarul legal 2026, secții/departamente, schimburi și categorii LAB/CUSTOM, pontaj lunar în ore, concedii cu cereri atașate și aprobare, precum și rapoarte configurabile cu imagini și paginare. Fiecare build produce acum un director unic și manifest cu hash, cu temporare izolate și datele existente păstrate; curățarea locală a eliminat 340,86 MiB de artefacte vechi.

Definițiile rapoartelor au antet/conținut/subsol independente, texte RO/EN, PDF A4/A3, previzualizare și deschidere pentru tipărire. Centralizatorul lunar și pontajul detaliat sunt modele separate; cel detaliat are cele 12 coloane finale. Configurare globală pentru atribuirea turelor peste miezul nopții și sursele totalurilor, cu revizii și audit. T este totalul, ZNL clasifică weekenduri și zile libere 00–24, iar Ore norm clasifică luni–vineri fără zile libere. Formulele neconfirmate rămân nedefinite până la configurare explicită. Codurile programelor se păstrează în date; grila, copierea și rapoartele lunare afișează ore.

Baza ajunge la migrarea 0015; actualizarea păstrează modelele și imaginile, inclusiv proba separată 0.20.0, cu backup înaintea completării istoricului cunoscut. Fișele individuale înlocuiesc coloana de cod cu ore. Centralizator v3, fișă individuală v3, pontaj detaliat v2 și absențe v2, cu referințele anterioare păstrate.

Validare pe ramura de integrare:

- Toate vârfurile ramurilor locale și `origin` sunt strămoși ai integrării.
- `gofmt -l -w .`, `go vet ./...`, `go test ./...` și TypeScript/Vite au trecut.
- Verificarea Edge cu punte Wails simulată a trecut fluxurile reunite, inclusiv 50×31, cereri/atașamente/aprobări, imagini, configurare/conflicte, coduri păstrate și ore afișate, precum și cele două rapoarte lunare. Backendul folosește baze criptate fictive în testele Go; calculele verifică trecerea lună/an, clasificarea calendaristică, absențele și sursele/factorii expliciți.
- PDF-urile cu imagini și paginare au validare pdfcpu, verificare independentă pypdf strict și randare Windows documentate în F09-T08.
- Go 1.27.3 și dependențele declarate Go/npm verificate la 2026-10-26, fără actualizări necesare.

Rapoartele sunt documente de lucru; calculele salariale finale așteaptă specificațiile contabilității. Dialogurile desktop native și tipărirea fizică nu au fost probate interactiv. Buildul final se face din `main` după merge, prin `build.ps1`, cu temporare izolate și păstrarea datelor portabile.
