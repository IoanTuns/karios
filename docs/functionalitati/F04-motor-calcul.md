# F04 · Motorul de calcul

Etapă: V1 · Stare: în lucru · Depinde de: F01, F02

## Scop

Fiecare tură introdusă se împarte automat pe categorii de tarif, corect și în cazurile grele: peste miezul nopții, gărzi de 24 h, sărbători, schimbarea orei.

## Decizii

- [D003](../decizii/D003-motor-calcul.md) — principiile motorului
- [D002](../decizii/D002-portabilitate.md) — fusul orar inclus în `.exe`

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F04-T01 | `Segmenteaza`, `Valideaza` și calendarul cu sărbători | gata | Claude | `internal/calcul/`, 2026-10-24 |
| F04-T02 | Teste: exemplele din documentul de design, gărzi, sărbători, limite de interval, schimbarea orei | gata | Claude | 20 de teste, toate trec |
| F04-T03 | Modelul de sporuri: categorie de bază + sporuri cumulate (de exemplu noapte peste weekend) | blocat | | Așteaptă F01-T02; probabil o decizie nouă |
| F04-T04 | Orele suplimentare, calculate separat față de normă | blocat | | Așteaptă F01-T02 |
| F04-T05 | Scăderea pauzei din segmentele turei | blocat | | Așteaptă F01-T02 |
| F04-T06 | Reguli versionate în timp, inclusiv procentele | de făcut | | Procentul stă acum în `categorii_ore`, care nu e versionat |
| F04-T07 | Testele de referință din luna reală primită de la contabilitate | blocat | | Așteaptă F01-T02 |
| F04-T08 | Ora ambiguă (03:00–04:00 în noaptea în care ora se dă înapoi) la introducere | de făcut | | Se rezolvă în interfață; motorul primește momente exacte |
