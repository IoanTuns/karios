# Decizii în așteptare · 2026-11-06

Pregătit de Claude (F02-T25), pentru Ioan Tuns. Un agent doar propune; starea o schimbă un om (AGENTS.md §2).

## Rezultat

Ioan Tuns a răspuns „da la toate” (2026-11-06): recomandările de mai jos s-au aplicat întocmai. 60 de decizii sunt acceptate integral, D084 e înlocuită de D087, D049 e acceptată doar pentru primul lot, D041 rămâne propusă și amânată (I061 amânată). Taskuri noi din D089: F07-T41 (prefix la import) și F14-T33 (prefix în asistentul AI). D090 e acceptată cu condiția probei pe date reale.

## Pe scurt

63 de decizii au starea `propusă`. Aproape toate descriu lucruri cerute de tine și deja implementate în `main` (0.92.45); le lipsește doar confirmarea formală. Le-am împărțit în trei grupe:

| Grupa | Câte | Ce îți cer |
| --- | --- | --- |
| A. Implementate, fără întrebări deschise | 51 | Un „da” pentru toată grupa, dacă nu vezi ceva greșit în tabel |
| B. Înlocuite parțial de editorul pdfme | 5 | Confirmarea că partea de aspect e înlocuită de D087 |
| C. Cer o alegere sau o verificare reală | 7 | Răspuns pe fiecare |

Poți răspunde scurt, de exemplu: „A da, B da, C1 da, C2 amână, …”. Eu actualizez apoi stările cu `Decisă de: Ioan Tuns`.

## C. Cer o alegere sau o verificare reală

| # | Decizie | De ce nu e doar formalitate | Recomandarea mea |
| --- | --- | --- | --- |
| C1 | [D063](decizii/D063-norma-zilnica-datorata.md) Norma zilnică datorată | Un agent a interpretat „CVE” din răspunsul tău ca **CEV**. Calculul normei depinde de asta. | Confirmă că era CEV. Dacă nu, devine întrebare pentru contabilitate (F01). |
| C2 | [D041](decizii/D041-motor-reguli-configurabile.md) Motor de reguli configurabile | Neimplementat. Arhitectură propusă pentru reguli definite liber (repaus minim, plafoane, nopți consecutive). Ideea I061 e „în analiză”. | **Amână** până vin regulile de la contabilitate și spital: fără ele nu știm ce trebuie să exprime motorul. |
| C3 | [D049](decizii/D049-prioritati-ui-asistent.md) Priorități UI pentru asistent | Ai aprobat doar primul lot (I072–I074, implementat în F14-T24). Loturile următoare (I075 și mai departe) nu sunt aprobate. | Acceptă primul lot; restul rămâne propus. |
| C4 | [D088](decizii/D088-fonturi-windows-rapoarte.md) Fonturi Windows în rapoarte | Fontul se încorporează în PDF. Unele fonturi comerciale interzic încorporarea; aplicația nu verifică licența fiecărui font. | Acceptă, cu regula: pentru rapoarte oficiale, fonturile Windows standard (Arial, Calibri, Times New Roman) sau fontul inclus. |
| C5 | [D087](decizii/D087-rapoarte-sql-pdfme.md) detaliile tehnice | Direcția e decisă de tine. Rămân de acceptat: pluginul propriu de tabel, statusul de „referință” al rapoartelor standard, limitele SQL. | Acceptă: sunt implementate și verificate pe 0.92.36. |
| C6 | [D089](decizii/D089-prefix-angajat.md) detaliile tehnice | Cerințele sunt decise. Rămân: tabelul de prefixe, ordinea, afișarea. Și o întrebare nouă: prefixul în asistentul AI și la import. | Acceptă detaliile. Pentru AI și import: asistentul recunoaște numele cu sau fără prefix; importul citește prefixul dintr-o coloană separată, opțională. |
| C7 | [D090](decizii/D090-rapoarte-standard-sql.md) Rapoartele standard în SQL | Calculul rămâne în Go, cu teste de referință; SQL doar alege și ordonează. Schimbă felul în care se întrețin rapoartele. | Acceptă, după ce proba ta pe date reale arată aceleași cifre ca înainte. |

## B. Înlocuite parțial de editorul pdfme (D087)

Partea de **aspect** (Markdown, poziționare liberă, desenarea cu gofpdf) a fost ștearsă în F09-T23. Partea de **date** (SQL pe instantaneu izolat, formule, coloane calculate) a rămas și e folosită de pdfme.

| Decizie | Ce rămâne valabil |
| --- | --- |
| [D025](decizii/D025-definitii-rapoarte.md) Definiții de rapoarte | Meniul de definire, antet/conținut/subsol |
| [D026](decizii/D026-imagini-rapoarte.md) Imagini în rapoarte | Imaginile, acum ca elemente pdfme |
| [D082](decizii/D082-rapoarte-flexibile.md) Rapoarte SQL și formule | SQL, formulele, coloanele calculate |
| [D084](decizii/D084-aliniere-antet-flexibil.md) Zonele antetului | Nimic separat: aranjarea se face acum vizual |
| [D085](decizii/D085-personalizare-rapoarte.md) Prezentare și copii editabile | Bold la weekend și sărbători, intervalul gărzilor peste lună |

Recomandare: le marchez „acceptată, aspectul înlocuit de D087”. D084 devine „înlocuită de D087”.

## A. Implementate, fără întrebări deschise

Cerute de tine și livrate; regulile de calcul din ele sunt confirmate explicit de tine în conversații.

| Domeniu | Decizii |
| --- | --- |
| Angajați, structură, contracte | [D012](decizii/D012-date-angajati-norme-sarbatori.md), [D013](decizii/D013-primele-functionalitati-demo.md), [D018](decizii/D018-sectii-departamente.md), [D060](decizii/D060-marci-garda.md), [D062](decizii/D062-structura-profesionala-contracte.md), [D066](decizii/D066-nume-ierarhic.md) |
| Calendar și sărbători | [D014](decizii/D014-sarbatori-2026.md), [D067](decizii/D067-calendar-lucrator.md) |
| Schimburi și programe | [D015](decizii/D015-nomenclator-schimburi.md), [D019](decizii/D019-schimburi-initiale-pdf.md), [D021](decizii/D021-programe-lucru.md), [D023](decizii/D023-categorii-schimburi.md), [D059](decizii/D059-nomenclatoare-schimburi-functii.md), [D072](decizii/D072-tipuri-regimuri-configurabile.md), [D086](decizii/D086-programe-multiple.md) |
| Pontaj | [D016](decizii/D016-grila-pontaj-schita.md), [D028](decizii/D028-pontaj-ore-doua-rapoarte.md), [D031](decizii/D031-precompletare-avertismente-pontaj-backup-manual.md), [D032](decizii/D032-imbunatatiri-ergonomie-ui.md), [D033](decizii/D033-imbunatatiri-ux-fluxuri-pontaj-absente.md), [D051](decizii/D051-pontaj-planificare-manuala.md) |
| Absențe și concedii | [D020](decizii/D020-concedii-absente.md), [D024](decizii/D024-cereri-concediu-aprobare.md), [D030](decizii/D030-suprapuneri-inlocuitor-previzualizare-cereri.md), [D068](decizii/D068-absente-din-pontaj.md), [D069](decizii/D069-prioritate-absente.md), [D070](decizii/D070-corectare-explicita-absente.md), [D073](decizii/D073-mutare-stare-absente.md), [D080](decizii/D080-efecte-absente.md), [D081](decizii/D081-categorii-absente.md) |
| Planificare și solduri | [D037](decizii/D037-planificare-garzi-dupa-amiaza-publicare.md), [D076](decizii/D076-prioritate-publicare-pontaj.md), [D077](decizii/D077-solduri-lunare-initializare.md) |
| Rapoarte | [D064](decizii/D064-legenda-raport.md), [D083](decizii/D083-explicatii-optionale.md) |
| Interfață | [D061](decizii/D061-fereastra-adaptiva.md), [D065](decizii/D065-zoom-interfata.md), [D071](decizii/D071-format-ore-24h.md), [D074](decizii/D074-filtre-stare.md), [D075](decizii/D075-stergere-configurari-admin.md), [D078](decizii/D078-remedii-ui-si-export-verificat.md), [D079](decizii/D079-simulare-absente-si-flux-lunar.md) |
| Asistent AI | [D047](decizii/D047-context-asistent-ai.md), [D048](decizii/D048-interactiune-ai-ghidata.md), [D050](decizii/D050-instrumente-asistent.md) |
| Tehnic și proces | [D017](decizii/D017-build-fara-stergerea-datelor.md), [D022](decizii/D022-revizii-concurenta.md), [D027](decizii/D027-reunire-main-pr.md), [D029](decizii/D029-builduri-in-directoare-separate.md), [D057](decizii/D057-manual-ajutor-contextual.md), [D058](decizii/D058-manuale-la-build.md) |

Note:

- D019 înlocuiește D015 doar în privința setului inițial de schimburi; ambele rămân valabile în rest.
- D027 descrie o integrare deja făcută (proces istoric).
- Nu am recitit integral fiecare decizie. Am verificat starea, rezumatul și dacă au taskul `gata`. Dacă vrei să revezi una anume înainte de „da”, spune care.
