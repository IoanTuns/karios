# F16 · Sistem de Diagnoză, Captare Erori și Rafinări UI/UX

Stare: gata · Responsabil: Antigravity · Decizii: [D043](../decizii/D043-diagnoza-erori-profilare-local.md)

## Descriere

Sistem integrat de diagnoză locală, captare robustă a erorilor pe frontend și backend, măsurare a timpului de execuție a interogărilor (profilare performanță), urmat de pachetul de rafinări ergonomice UI (bară laterală compactă, zoom pe grilă, tooltips rapide, notificări toast, scurtături de tastatură, mod nocturn și ghid de pornire).

## Taskuri

| ID | Task | Responsabil | Stare | Legături | Note |
| --- | --- | --- | --- | --- | --- |
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
