# F01 · Cerințe și reguli de calcul

Etapă: V1 · Stare: în lucru · Depinde de: —

## Scop

Regulile de plată și formatul rapoartelor sunt confirmate de contabilitate și transformate în teste automate, înainte de a fixa schema și motorul.

## Decizii

- [D003](../decizii/D003-motor-calcul.md) — principiile motorului; se revizuiește după răspunsuri

## Taskuri

| ID | Task | Stare | Note |
| --- | --- | --- | --- |
| F01-T01 | Chestionar pentru contabilitate | gata | [Pontaj — întrebări pentru contabilitate](https://claude.ai/artifact/E8SLB7zT2WAzKTyFSchmK6), 2026-10-24 |
| F01-T02 | Trimite chestionarul; primește răspunsurile și un centralizator real completat pe o lună | de făcut | Deblochează F04-T03 … F04-T05 și F04-T07 |
| F01-T03 | Transformă răspunsurile în reguli și cazuri de test | blocat | Așteaptă F01-T02 |
| F01-T04 | Validează cu contabilul formatul centralizatorului și al exportului de salarizare | de făcut | Coloanele T, ZNL, L/FIN, celulele roșii, chenarele |
| F01-T05 | Verifică disponibilitatea numelui Pontix (OSIM, EUIPO, domeniu) | de făcut | Alternativa: Cadran |
