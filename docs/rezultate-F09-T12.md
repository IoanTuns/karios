# F09-T12 - Legenda coloanelor raportului

Responsabil: Codex. Ramura: F09-T12-legenda-raport. Aplicatie 0.92.3, manual 1.14.

## Modificare

Ore norma este eliminata numai din raportul detaliat. Fiecare coloana vizibila are explicatie in subsolul PDF, pe rand separat la 7 pt; zilele lunii au o explicatie comuna. Legenda include semnificatia, sursa activa, codurile, factorul, aprobarile si variantele conditionate aplicabile persoanelor incluse. Este generata in RO/EN si inclusa sub tabel in Excel-ul documentului previzualizat. Semnaturile, imaginile si metadatele sunt pastrate.

Rapoartele pontaj/detaliat/individual sunt v6, absente/garzi v4. Referintele anterioare sunt pastrate; asteptarile financiare existente nu au fost modificate. Fara schimbari de calcul sau schema BD. [D064](decizii/D064-legenda-raport.md) este propusa.

## Verificari

Testele pachetului rapoarte au trecut: referinte, PDF, paginare, imagini, grupare, marci distincte, legenda RO/EN si Excel redeschis. PDF-ul detaliat fictiv a fost randat cu motorul Windows si inspectat: legenda nu se suprapune cu tabelul sau semnaturile. TypeScript/Vite si 18 teste frontend trecute; 3 teste manual, 7 teste livrare si control documentatie trecute. Verificarea diff nu raporteaza spatii finale. gofmt, go vet ./... si suita completa go test -p 1 -timeout 120s ./... au trecut. Build oficial complet din b50555a: `build/bin/pontix-0.92.3-20261031-140414-7c8084c7219d4b7fb049e54a88653dde/pontix.exe`. Manual 1.14 MD/HTML/PDF alaturi si ZIP separat; manifest, versiune/commit, hashurile executabilului si manualelor, integritatea ZIP, SHA256SUMS din arhiva si hashul arhivei verificate.

Dependente: Go 1.27.3 si npm actuale; 62 actualizari tranzitive Go raman urmarite in F14-T21. Fara biblioteci noi.

## Limite si predare

Legenda descrie configuratia actuala. Norma datorata, marcajul T, separarea noua a orelor de garda si soldurile L FIN/L Inceput raman de integrat in F07-T27/T28. Absentele dupa zile programate si introducerea explicita a orelor nemotivate raman in lista existenta. Nu sunt prezentate drept implementate. Proba fizica la imprimanta si desktopul real nu au fost efectuate. Datele reale nu au fost deschise. Fara push sau merge in main.
