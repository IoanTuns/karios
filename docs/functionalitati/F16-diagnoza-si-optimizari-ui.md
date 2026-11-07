# F16 · Sistem de Diagnoză, Captare Erori și Rafinări UI/UX

Stare: gata · Responsabil: Antigravity · Decizii: [D043](../decizii/D043-diagnoza-erori-profilare-local.md)

## Descriere

Sistem integrat de diagnoză locală, captare robustă a erorilor pe frontend și backend, măsurare a timpului de execuție a interogărilor (profilare performanță), urmat de pachetul de rafinări ergonomice UI (bară laterală compactă, zoom pe grilă, tooltips rapide, notificări toast, scurtături de tastatură, mod nocturn și ghid de pornire).

## Taskuri

| ID | Task | Responsabil | Stare | Legături | Note |
| --- | --- | --- | --- | --- | --- |
| F16-T39 | Planificare compactă pe un rând și optimizare înălțime ecran (I122) | Antigravity | gata | 0.92.78, ramura F16-T39-planificare-compacta | Cerere utilizator. Aliniere acțiuni pe un rând, meniu derulant „Mai multe ▾” pentru acțiuni secundare (raport gărzi, reguli, deschide în pontaj, reîncarcă, densitate S/M/L), insigne inline compacte pentru etape, scădere spațiere antet și alerte; primele rânduri ale grilei vizibile fără derulare la 1024×640. |
| F16-T38 | Ergonomie pontaj: navigare Numpad (Enter vertical), total OL în bara de inspecție și securizare Închide luna | Antigravity | gata | 0.92.70, ramura F16-T38-ergonomie-pontaj | Cerere utilizator. Navigare fluidă Numpad (Enter coboară vertical pe aceeași zi, Tab orizontal), buffer tastare numerică directă fără salt forțat, afișare OL în bara de inspecție celulă, stil dedicat buton-sigilare și sinteză cu avertismente la închiderea lunii. |
| F16-T33 | Analiza UI a rapoartelor cu editor vizual pdfme (0.92.34) | Claude | gata | [Analiza](../analiza-ui-rapoarte-pdfme-0.92.34.md), D087, D088 | Cerere utilizator. Scenariu Edge nou `--rapoarte-pdfme`, 12 constatari UIR-01..12, fara implementare. Pas urmator: utilizatorul alege ce se implementeaza; recomandat UIR-01+02. |
| F16-T34 | Navigarea din asistentul AI si confirmarea modificarilor nesalvate (I117) | Claude | gata | [I117](../idei.md), ramura F16-T34-confirmare-navigare-ai | Scenariul Edge --ajutor picat la integrare (0 confirmari in loc de 2). Cauza: scenariul numara window.confirm, dar din F16-T31 confirmaParasire deschide DialogParasire cand ecranul are salvare inregistrata (Pontajul are). Aplicatia trece corect prin poateAbandona la ambele butoane (Deschide in Pontaj, formular concediu). Schimbat doar verificari/ajutor-scenarii.mjs: asteapta dialogul, apasa Ramai, verifica perioada, ecranul si valoarea 7, si 0 apeluri window.confirm. Fara schimbare de cod in aplicatie, deci fara versiune noua. --ajutor trece. |
| F16-T35 | Remedii mici dupa F05-T20/F14-T35: coloana fixata din Configurare pontaj (I119) si proba invechita a sabloanelor AI (I120) | gata | Claude | 0.92.51, ramura F14-T35-plan-ai (la cererea utilizatorului, pasii 1-5 pe aceeasi ramura). I119: clasa tabel-reguli-pontaj si exceptie CSS pentru ultima coloana (corp si antet); verificat in captura Edge --nomenclatoare la 1024 px. I120: nu era defect; selectorul probei numara si comutatorul din F14-T32, corectat; ruleaza-instrumente.ps1 24/24. |
| F16-T36 | Analiza UX a versiunii 0.92.59: pro, contra, recomandari | gata | Claude | [Analiza](../analiza-ux-0.92.59.md), cerere Ioan Tuns 2026-11-07. Scenariu nou de capturi `--capturi-ux` (verificari/capturi-ux-scenarii.mjs), 7 ecrane la 1024×640 si 1440×900. 12 constatari (UX59-01..12), 9 recomandari (R1..R9); nimic implementat. Pas urmator: alegerea recomandarilor de catre utilizator. |
| F16-T37 | Remediile analizei UX 0.92.59 (I121), pe rand: R4 diacritice, R5 termeni si date, R7 motive la butoane dezactivate, R1 antet compact, R2 bara Pontaj, R3 totaluri vizibile, R6 Setari pe sectiuni, R8 asistent pe toata inaltimea, R9 finisari | gata | Claude | Cerere Ioan Tuns 2026-11-07 („le faci pe toate unul dupa altul”). Ramura F16-T37-remedii-ux (din F08-T11-coloane-pontaj), cate un commit si o versiune pe recomandare. R4 gata (0.92.61): 204 texte, test de diacritice, rapoartele de referinta cu versiune noua dupa verificarea ca diferenta e doar de diacritice. R5 gata (0.92.62): SelectieRezumat cu nume si date formatate, dataAfisata comuna (Pontaj, Planificare, Absente, asistent), „Densitate” si „Mărire interfață”; test vitest. R7 gata (0.92.63): motive (chei motiv.*) pe 5 butoane principale; verificare Edge pe Salvează schița. R1 gata (0.92.64): antet titlu-antet + unelte-antet, masurat in --capturi-ux (antetJos ≤145 px la 1024×640); toate scenariile Edge trec. R2 gata (0.92.65): bara Pontaj pe un rand, meniul Mai multe (details, inchidere la actiune/Escape), Stare si filtre profesionale in panoul Filtre, selectia in bara de inspectie; 5–6 randuri vizibile la 1024×640, verificat in --pontaj-absente. R3 gata (0.92.66): comutator Zile/Totaluri (clasa col-zi, tabel doar-totaluri, localStorage), verificat in --pontaj-absente (7 coloane in cadrul grilei). Varianta cu toate totalurile fixate la dreapta respinsa: la 1024 px ar fi lasat 2–3 zile vizibile. R6 gata (0.92.67): file in Setari › General (sectiuni hidden, display: contents), verificarea configurarii deschisa cand are probleme; 2416 → 1332 px. R8 gata (0.92.68): luna de lucru in randul de navigare, conversatia 334 px la 1024×640 (masurat in --capturi-ux), pasii de instalare deschisi si buton spre ghid; probele asistentului 24/24 si proba planului trec. R9 gata (0.92.69): randul de acoperire ascuns cand e gol, eticheta comuna a filtrului de stare, filtrele Planificarii pe un rand. Toate cele 9 recomandari (0.92.61–0.92.69) facute; masuratorile inainte/dupa in analiza. Ramas, ca idei noi: I122 (Planificare cu aceeasi bara compacta ca Pontajul), I123 (diacriticele din legendele rapoartelor, care cer versiuni noi de raport). |
| F16-T27 | Filtre Active/Inactive/Toate, implicit Active, in listele cu stare | Codex | gata | [D074](../decizii/D074-filtre-stare.md) | [Rezultate](../rezultate-F16-T27.md); proba operator. |
| F16-T26 | Afisare unitara a orelor in format 24h | Codex | gata | [D071](../decizii/D071-format-ore-24h.md) | [Rezultate](../rezultate-F16-T26.md); proba operator. |
| F16-T25 | Zoom global cu butoane minus/plus, resetare, tastatura si memorare locala | Codex | gata | [D065](../decizii/D065-zoom-interfata.md) | [Rezultate](../rezultate-F16-T25.md); proba operator pe monitorul propriu. |
| F16-T01 | Backend diagnoză, rotire jurnale și profilare interogări | Antigravity | gata | D043 | `internal/diagnoza` |
| F16-T02 | Frontend ErrorBoundary și captare globală a excepțiilor | Antigravity | gata | D043 | `ErrorBoundary.tsx` |
| F16-T03 | Afișare metrici diagnoză și export jurnale în Setări | Antigravity | gata | D043 | `Formulare.tsx` |
| F16-T04 | Bară laterală pliabilă (compact 60px) cu comutator | Antigravity | gata | UI | `App.tsx`, `style.css` |
| F16-T05 | Control densitate & Zoom pentru grilele de pontaj/planificare | Antigravity | gata | UI | `Pontaj.tsx`, `Planificare.tsx` |
| F16-T06 | Tooltips rapide la hover pe celulele din grilă | Antigravity | gata | UI | `Pontaj.tsx`, `Planificare.tsx` |
| F16-T07 | Sistem unificat de notificări Toast | Antigravity | gata | UI | `Toast.tsx` |
| F16-T08 | Scurtături universale tastatură (Escape, Ctrl+Enter) | Antigravity | gata | UI | `Formulare.tsx` |
| F16-T09 | Suport Dark Mode / Temă Nocturnă | Antigravity | gata | UI | `style.css` |
| F16-T10 | Ghid vizual de pornire pe ecrane goale (Empty state onboarding) | Antigravity | gata | UI | `Pontaj.tsx`, `Planificare.tsx` |
| F16-T11 | Jurnalul de diagnoză fără date personale și panicile metodelor Wails înregistrate (SEC-7, SEC-8) | gata | Claude | 0.88.4: SEC-7, SEC-8 |
| F16-T12 | Contrastul grilelor pe ambele teme, cu tokenuri de culoare pentru grile (UI-3, UI-2 pentru grile) | gata | Claude | 0.88.5: toate stările grilei ≥ 4,5:1 pe ambele teme |
| F16-T13 | Tema inițială după setarea din Windows (UI-1) | gata | Claude | 0.88.6 |
| F16-T14 | Formatterele de dată refolosite în loc de create la fiecare randare (PERF-2) | gata | Claude | 0.88.7: pas de selecție la 300 de angajați 27–53 ms |
| F16-T15 | Butonul „Raportează o problemă”: formular, raport anonimizat salvat ca fișier și link către issue GitHub precompletat (I093) | gata | Claude | 0.89.0: formular, raport anonimizat, fișier și link GitHub, fără token |
| F16-T16 | Tema întunecată pe tokenuri și pentru celelalte ecrane, cu verificare de contrast (restul UI-2) | gata | Claude | 0.88.12: 132 → 0 probleme de contrast pe toate ecranele, ambele teme |
| F16-T17 | Meniu contextual desktop la clic dreapta și copiere/lipire blocuri celule (Clipboard) în grilă | Antigravity | gata | UI | `Pontaj.tsx`, `style.css` |
| F16-T18 | Mini-hartă vizuală de acoperire pe zile (Heatmap acoperire / deficit în antetul grilei) | Antigravity | gata | UI | `Pontaj.tsx`, `style.css` |
| F16-T19 | Panou lateral retractabil (Drawer) pentru Asistentul AI integrat direct în Pontaj | Antigravity | gata | UI | `Pontaj.tsx`, `AsistentAI.tsx`, `style.css` |

| F16-T20 | Analiza UX pentru release pilot 0.88.14 | Codex | gata | docs/analiza-ux-release-0.88.14.md | Audit euristic, fluxuri si prioritati; fara modificari de interfata. |

| F16-T21 | Remedieri UX01-UX08 din auditul release pilot | Codex | gata | analiza-ux-release-0.88.14.md | Protectie navigare AI, indicator neutru, ajutor fara model, drawer, cautare si tastatura; validare recuperare pe date fictive. |
| F16-T22 | Internaționalizare panou eroare (ErrorBoundary, I095) și formatare defensivă a datelor în UI (I094) | Antigravity | gata | D004, I094, I095 | `ErrorBoundary.tsx`, `formatare.ts`, `Jurnal.tsx`, `Restaurare.tsx` |

| F16-T23 | Corectarea focusului/cursorului din formulare si verificarea controalelor zoom | Codex | gata | Formulare.tsx; verificari UI | Reproducere tastare reala si modificare in mijlocul textului, controale vizibile la rezolutii mici si mari. Cerere suplimentara: importul zilelor anuale de concediu; antetul oficial nu era recunoscut, nici Norma_Ore_Zi. Probe CSV/XLSX si previzualizare cu concediu. |

Predare F16-T23: cod 6f046a2, build complet 0.89.4 in build/bin/pontix-0.89.4-20261031-103113-d2f322809c2845d185e2bd69d4b8abb0, manual 1.7 si ZIP verificate. Pas urmator: proba operator si review PR daca se solicita. Fara push.

| F16-T24 | Dimensiunea initiala a ferestrei adaptata ecranului si DPI | gata | Codex | Limite dinamice dupa zona de lucru Windows, fara acoperirea barei de activitati. |

F16-T24: [rezultate si limite](../rezultate-F16-T24.md). Pas urmator: build oficial si proba utilizator la rezolutia/scalarea sa.

| F16-T28 | Audit UI/UX 0.92.18: avantaje, probleme si recomandari | gata | Codex | Capturi utilizator si inspectie cod; completat cu UI28-13 capete tabele la derulare si UI28-14 exporturi Excel care nu se deschid, raportate de utilizator si inca nereproduse. Fara modificari de interfata. [Audit](../analiza-ui-0.92.18.md). |

| F16-T29 | Implementarea tuturor remediilor UI28-01 ... UI28-14 | gata | Codex | 0.92.19: UI28-01...13 implementate/verificate; investigatia blocajului Excel UI28-14 amanata explicit de utilizator, protectiile de salvare deja implementate pastrate. Urmeaza build/proba operator, fara push/merge. [Plan si verificari](../plan-F16-T29.md). |

| F16-T30 | Analiza UX dupa remediile 0.92.19 | Codex | gata | [Audit si prioritati](../analiza-ux-0.92.19.md). 12 recomandari, limite explicite; pas urmator: lot de implementare, fara push/merge. |

| F16-T31 | Implementarea celor sase recomandari UX30 | Codex | gata | 0.92.20/manual 1.31. [Implementare si verificari](../plan-F16-T31.md). Urmeaza build/proba operator; fara push/merge, Excel amanat. |

| F16-T32 | Istoric si mod de calcul la sfarsitul rapoartelor | Codex | gata | 0.92.21: anexa finala PDF verificata RO/EN, rapoarte versionate, manual 1.32. [Rezultate](../plan-F16-T32.md). Urmeaza build/proba operator; fara push/merge. |
