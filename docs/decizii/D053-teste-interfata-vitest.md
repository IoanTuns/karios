# D053 · Teste automate pentru interfață: Vitest

Stare: acceptată · Data: 2026-10-30 · Decisă de: Ioan Tuns (recomandarea 7 din 2026-10-29, „începe implementarea propunerilor”) · Taskuri: F02-T16

## Context

Interfața nu avea niciun test automat. Logica din `asistentSesiune.ts` (repetări), validarea parolei sau diferențele din jurnal erau verificate doar de tipuri (`tsc`), deci o greșeală de logică ajungea neobservată până la proba manuală.

## Decizie

- **Vitest 5.0.5** (licență MIT, verificată pe npm), ca dependență de dezvoltare cu versiunea fixată exact. Folosește configurația Vite existentă; nu adaugă un al doilea sistem de build.
- Testele rulează în mediul Node, fără DOM simulat: se testează funcțiile pure, mutate în module mici (`parola.ts`, `repetari.ts`, `jurnalDiferente.ts`). Componentele React nu se testează acum.
- `npm test` rulează testele; CI le rulează după build ([D046](D046-verificari-automate-github.md)).

## Variante respinse

- **Jest** — cere configurare separată pentru TypeScript și ESM; Vitest refolosește Vite.
- **Testing Library + jsdom acum** — mai multe dependențe; se reia când apar teste de componente.

## Consecințe

- Logica nouă din interfață se scrie, pe cât se poate, ca funcții pure testabile.
- Vitest intră în verificarea dependențelor de la fiecare modificare (AGENTS.md §6).
