# F16 · Sistem de Diagnoză, Captare Erori și Rafinări UI/UX

Stare: gata · Responsabil: Antigravity · Decizii: [D043](../decizii/D043-diagnoza-erori-profilare-local.md)

## Descriere

Sistem integrat de diagnoză locală, captare robustă a erorilor pe frontend și backend, măsurare a timpului de execuție a interogărilor (profilare performanță), urmat de pachetul de rafinări ergonomice UI (bară laterală compactă, zoom pe grilă, tooltips rapide, notificări toast, scurtături de tastatură, mod nocturn și ghid de pornire).

## Taskuri

| ID | Task | Responsabil | Stare | Legături | Note |
| --- | --- | --- | --- | --- | --- |
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
