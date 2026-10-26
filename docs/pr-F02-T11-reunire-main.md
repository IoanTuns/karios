Reunește în `main` aplicația Pontix 0.21.0 și istoricul tuturor ramurilor locale și remote. Integrarea include administrarea angajaților și normelor, calendarul legal 2026, secții/departamente, schimburi și categorii LAB/CUSTOM, pontaj lunar, concedii cu cereri atașate și aprobare, precum și rapoarte configurabile cu imagini și paginare.

Definițiile rapoartelor au antet/conținut/subsol independente, texte RO/EN, PDF A4/A3, previzualizare și deschidere pentru tipărire. Baza ajunge la migrarea 0014; actualizarea păstrează și modelele create în proba separată 0.20.0, cu backup înaintea completării istoricului cunoscut. Cele trei rapoarte sunt v2, cu referințele v1 păstrate.

Validare pe ramura de integrare:

- Toate vârfurile ramurilor locale și `origin` sunt strămoși ai integrării.
- `gofmt -l -w .`, `go vet ./...`, `go test ./...` și TypeScript/Vite au trecut.
- Verificarea Edge cu punte Wails simulată a trecut fluxurile reunite, inclusiv 50×31, cereri/atașamente/aprobări și imagini în rapoarte. Backendul folosește baze criptate fictive în testele Go.
- PDF-urile cu imagini și paginare au validare pdfcpu, verificare independentă pypdf strict și randare Windows documentate în F09-T08.
- Go 1.27.3 și dependențele declarate Go/npm verificate la 2026-10-26, fără actualizări necesare.

Rapoartele sunt documente de lucru; calculele salariale finale așteaptă specificațiile contabilității. Dialogurile desktop native și tipărirea fizică nu au fost probate interactiv. Buildul final se face din `main` după merge, prin `build.ps1`, cu temporare izolate și păstrarea datelor portabile.
