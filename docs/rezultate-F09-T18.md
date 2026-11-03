# F09-T18 - Aliniere antet/subsol flexibil

Data 2026-11-03. Responsabil Codex. Ramura F09-T18-aliniere-antet-flexibil. Aplicatie0.92.29, raport flexibil3, manual1.39. D084 ramane propusa.

## Defect si corectie

Captura utilizatorului arata toate textele antetului stivuite la stanga, desi editorul are zone Stanga/Centru/Dreapta. PDF-ul flexibil TABEL/Markdown desena textele succesiv pe intreaga latime. Proba noua pe coordonatele reale din fluxul PDF reproduce defectul: textul central era la31.18pt in loc de276.08pt in cazul A4 portret testat.

Corectia deseneaza trei zone simultane egale, aliniate L/C/R. Inaltimea antetului este maximul textelor impachetate, pentru a pozitiona tabelul sub ultima linie. Imaginile sunt aliniate in zona lor, cu proportiile pastrate. Subsolul respecta aceleasi alinieri. Antetul si imaginile se repeta si pe anexa; explicatiile/istoricul incep dupa antet. Fara modificari de calcul, SQL, baza sau referinte contabile.

Coordonatele continutului LIBER raman explicite; nu se introduc texte noi peste elementele existente. Pentru un titlu unic central, folositi Antet > Centru si {{tabel}} in Markdown, fara un al doilea titlu # in corp.

## Verificari

- Testul nou verifica coordonatele textului si imaginilor din PDF-ul efectiv, nu numai parametrii rendererului. Centrare, aliniere dreapta, simultaneitate, tabel sub ultima linie a antetului, repetare pe toate paginile inclusiv anexa.
- 12 cazuri: RO/EN, A4 portret/A4 peisaj/A3 peisaj, TABEL/Markdown,160 randuri si imagini in toate cele sase zone. Trecute.
- Randare Windows a exemplului Markdown A4 peisaj verificata vizual: titlu central, zone distincte, tabel fara suprapunere si subsol aliniat. Probe anonime in build/bin/verificare-F09-T18. Directorul initial de probe creat accidental relativ la pachet a fost verificat si eliminat; fisierele livrate sunt in directorul corect build/bin.
- Documentatie0 erori,3 teste manuale,7 teste livrare trecute. Gofmt/vet si intreaga suita Go finala trecute.
- Dependente: react-i18next17.1.12 actualizat separat in0.92.28; build frontend,24 teste, Edge RO/EN si suita Go/vet au trecut. Go1.27.3/Wails2.17.1/gofpdf1.4.3 actuale. Excelize2.11.1 ramane amanat pentru dependenta efp prerelease; actualizari indirecte in F14-T21. Fara instalari noi, probe PDF cu pdfcpu si Windows existente.

## Predare

Pas urmator: proba utilizator pe raportul propriu, apoi review pentru integrare. Nu s-au accesat datele reale si nu s-a probat imprimanta fizica. Fara push sau merge. Buildul oficial se consemneaza dupa compilare.
