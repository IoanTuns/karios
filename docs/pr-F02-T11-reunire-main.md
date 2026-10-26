Reunește în `main` aplicația **Pontix 0.27.0** și istoricul tuturor ramurilor locale și remote (`F02-T11`, `F02-T12`, `F07-T11`–`F07-T13`, `F08-T07`–`F08-T09`, `F09-T07`–`F09-T09`, `F05-T16`, `F05-T17`).

## Rezumat al funcționalităților incluse

- **Arhitectură de meniu și Setări (`0.27.0`, `F05-T17`, [D034](decizii/D034-grupare-meniu-setari-tema-vizuala.md)):**
  - Meniu lateral cu 5 destinații principale (`Angajați`, `Pontaj`, `Concedii și absențe`, `Rapoarte`, `Setări`).
  - Gruparea ecranelor de configurare în **Setări** (`General`, `Zile libere legale` — fostul `Sărbători` —, `Schimburi`, `Configurare pontaj`, `Definire rapoarte`), accesibile atât prin taburi în ecranul Setări, cât și prin sub-meniul lateral, cu temă vizuală albastru-ardezie (`#2b4c7e` / `#eef3f9`).
- **Fluxuri UX și ergonomie UI în Pontaj și Absențe (`0.24.0`–`0.26.0`, `F08-T07`–`F08-T09`, `F05-T16`, [D031](decizii/D031-precompletare-avertismente-pontaj-backup-manual.md), [D032](decizii/D032-imbunatatiri-ergonomie-ui.md), [D033](decizii/D033-imbunatatiri-ux-fluxuri-pontaj-absente.md)):**
  - Bară de inspecție a celulei active (`.inspectie-celula`), paletă rapidă de schimburi dintr-un clic, stivă **Undo / Redo** (`Ctrl+Z` / `Ctrl+Y`), coloană de total `T` înghețată în dreapta (`sticky right`).
  - Navigare între luni (`◀` / `▶`), căutare după nume/marcă și filtru rapid de rânduri (`Toți angajații`, `Doar cu avertismente`, `Doar cu modificări`, `Doar înlocuitori`).
  - Selecție dreptunghiulară cu mouse-ul (`Click & Drag`) și autoderulare (`scrollIntoView`) la navigarea din săgeți.
  - Butoane **Precompletează L–V** și **Precompletează selecția L–V**, plus avertismente vizuale neblocante după `G24`, pentru `>24h/zi` și diferență între `T` și `Ore norm`.
  - Copie de siguranță criptată la cerere (`stocare.BackupManual`) din `Setări → General` și selector de limbă **RO / EN** în bara laterală.
- **Concedii, înlocuitori și cereri atașate (`0.19.0`, `0.23.0`, `F07-T12`, `F07-T13`, [D024](decizii/D024-cereri-concediu-aprobare.md), [D030](decizii/D030-suprapuneri-inlocuitor-previzualizare-cereri.md)):**
  - Cereri semnate atașate criptat în SQLite (PDF/PNG/JPEG), cu previzualizare modală și descărcare, număr de înregistrare, stare de aprobare și persoană înlocuitoare.
  - Validare tranzacțională împotriva suprapunerii între absența activă și rolul de înlocuitor în aceeași zi, evidențierea înlocuirii în pontaj (`celula-inlocuitor`), diferența aritmetică `Rămas` în sumarul anual, căutare liberă și buton **Selectează cererea** + salt contextual din Pontaj în Concedii și absențe.
- **Pontaj în ore, configurare pontaj și rapoarte PDF cu imagini (`0.18.0`–`0.22.2`, `F07-T11`, `F09-T07`–`F09-T09`, `F02-T12`, [D023](decizii/D023-categorii-schimburi.md), [D025](decizii/D025-definitii-rapoarte.md), [D026](decizii/D026-imagini-rapoarte.md), [D028](decizii/D028-pontaj-ore-doua-rapoarte.md), [D029](decizii/D029-builduri-in-directoare-separate.md)):**
  - Categorii de schimburi `LAB` / `CUSTOM`, grilă lunară în ore cu păstrarea codurilor de program, configurare globală a atribuirii turelor peste miezul nopții și a celor 12 coloane finale.
  - Definiții de rapoarte cu antet/conținut/subsol și imagini PNG/JPEG configurabile, paginare automată și două modele lunare distincte (centralizator v3 și pontaj detaliat v2).
  - Fiecare invocare `./build.ps1` construiește executabilul și manifestul `build.json` într-un director unic `build/bin/pontix-<versiune>-<data-ora>-<id>/`, fără suprascrierea datelor portabile.

## Verificări

- Toate vârfurile ramurilor locale sunt strămoși ai `F05-T17-grupare-meniu-setari`.
- `gofmt -l -w .`, `go vet ./...`, `go test ./...` și `npm --prefix frontend run build` (`tsc && vite build`) trecute.
- `node verificari/interfata.mjs` (Edge headless cu puntea Wails simulată) trecut integral.
- Build Windows verificat: `Pontix 0.27.0` (`069f850`, SHA256 `F532395C21B078A1073745D2F3E60D60566295DF71EBD6C36A65F92AD1AFE3DC`).

