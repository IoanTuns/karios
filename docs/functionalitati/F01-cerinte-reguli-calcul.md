# F01 · Cerințe și reguli de calcul

Etapă: V1 · Stare: în lucru · Depinde de: —

## Scop

Regulile de plată și formatul rapoartelor sunt confirmate de contabilitate și transformate în teste automate, înainte de a fixa schema și motorul.

## Decizii

- [D003](../decizii/D003-motor-calcul.md) — principiile motorului; se revizuiește după răspunsuri

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F01-T01 | Chestionar pentru contabilitate | gata | Claude | [Pontaj — întrebări pentru contabilitate](https://claude.ai/artifact/E8SLB7zT2WAzKTyFSchmK6), 2026-10-24 |
| F01-T02 | Trimite chestionarul; primește răspunsurile și un centralizator real completat pe o lună | de făcut | | Deblochează F04-T03 … F04-T05 și F04-T07 |
| F01-T03 | Transformă răspunsurile în reguli și cazuri de test | blocat | | Așteaptă F01-T02 |
| F01-T04 | Validează cu contabilul formatul centralizatorului și al exportului de salarizare | de făcut | | Coloanele T, ZNL, L/FIN, celulele roșii, chenarele |
| F01-T05 | Verifică disponibilitatea numelui Pontix (OSIM, EUIPO, domeniu) | de făcut | | Alternativa: Cadran |

## Modelul de raport primit la 2026-10-25

Ioan Tuns a confirmat funcționarea versiunii 0.10.0 și a trimis în conversație capturi de centralizator pentru octombrie 2026. Preluăm numai structura, fără numele, mărcile sau valorile reale ale angajaților: zilele 1–31, weekenduri gri, legendă T = TURA / N = NOPȚI și coloane de total Ore lucrat, T, ZNL, CO, CEV, CFS, Ore norm, Form prof, CM, L/FIN, Deleg, Ore gardă 200%. Captura mai largă arată rânduri distincte „GARDĂ”, cu mărci proprii, și grupări pe categorii profesionale; ipoteza inițială de două rânduri pentru fiecare angajat nu este confirmată.

Întrebări trimise în conversație: semnificația celor două rânduri și a codurilor C, FP, cev; confirmarea coloanelor T, ZNL, CO, CEV, CFS, Form prof, L/FIN. Nu deducem reguli de calcul sau procente din antetul raportului. Captura cu date potențial reale nu este copiată în repo.

Răspuns Ioan Tuns: specificațiile acestor coduri/coloane nu sunt disponibile momentan. Calendarul solicitat este lista completă de sărbători legale din România. Codul de schimb `N12` este acceptabil, cu necesitatea unui nomenclator configurabil; orele exacte nu sunt încă precizate. Rămâne de stabilit relația dintre persoană, rândul de gardă cu altă marcă și contract; grila inițială păstrează rânduri distincte după înregistrarea angajatului și nu calculează salarizarea.

Precizare ulterioară, 2026-10-25: setările inițiale ale schimburilor vin din PDF. Pentru codurile fără început explicit în document, Ioan Tuns a stabilit 7 = 07:00–14:00 și 12 = 14:00–02:00 a doua zi. Sursa și setul complet sunt în [D019](../decizii/D019-schimburi-initiale-pdf.md). Această precizare nu stabilește sporuri sau regulile coloanelor de salarizare.
