# D005 · Cum se documentează proiectul

Stare: acceptată · Data: 2026-10-24 · Taskuri: F02-T04

## Context

Toate acțiunile și deciziile trebuie documentate. Lucrul e organizat ca un roadmap de funcționalități, fiecare cu unul sau mai multe taskuri.

## Decizie

- Documentația stă în repo, în `docs/`, ca markdown, și se schimbă în același commit cu codul.
- Patru părți: `roadmap.md`, `functionalitati/Fxx-*.md` (cu taskurile), `decizii/Dxxx-*.md`, `jurnal.md`. Regulile sunt în [docs/README.md](../README.md).
- `CLAUDE.md` din rădăcină cere aceleași reguli și sesiunilor de lucru cu Claude.

## Variante respinse

- **Claude Docs / Google Docs** — nu sunt versionate cu codul; ușor de uitat la o schimbare. Rămân pentru documente trimise altor persoane (de exemplu chestionarul pentru contabilitate).
- **GitHub Issues** — repo-ul nu are încă un remote. Se pot adăuga mai târziu, legate de ID-urile de task.

## Consecințe

- Starea proiectului se vede din repo, fără alte unelte.
- Fiecare schimbare cere puțină muncă în plus: jurnal, starea taskului, eventual o decizie.
