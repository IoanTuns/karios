# Jurnal

Toate acțiunile proiectului, cele mai noi primele. O intrare = ce s-a făcut, de ce dacă nu e evident, cine a făcut-o și legăturile.

| Data | Acțiune | Cine | Legături |
| --- | --- | --- | --- |
| 2026-10-24 | Merge `fundatie` → `main`, la cererea lui Ioan Tuns | Ioan Tuns (aprobare), Claude | F02-T05 … T09 |
| 2026-10-24 | Două commit-uri pe `fundatie`: actualizarea Go (separat, după D007) și regulile pentru agenți, deciziile D002 și D006–D008, ideile. `gofmt`, `go vet` și testele trec | Claude | F02-T07 … T09 |
| 2026-10-24 | Decizie: cheia de date se împachetează de la început și cu cheia derivată din răspunsurile la întrebări (două copii în seif); trecerea la etapa 2 nu va cere recriptare | Ioan Tuns (decizie), Claude | [D008](decizii/D008-criptarea-datelor.md), [I006](idei.md) |
| 2026-10-24 | Decizie: criptarea începe cu cheie ascunsă în `.exe`, plus validarea ultimului utilizator sau întrebări de recuperare. D008 detaliată: chei pe două niveluri, cheia aplicației injectată la build în afara repo-ului, ce protejează și ce nu | Ioan Tuns (decizie), Claude | [D008](decizii/D008-criptarea-datelor.md), [I006](idei.md), F05-T10, F05-T12 |
| 2026-10-24 | Cerințe noi: datele operaționale în `karios/` lângă `.exe`, lizibile doar prin aplicație. D002 actualizată (punctul 4 amânat, punctul 6 aprobat: a doua instanță doar în citire). Propusă D008: criptarea datelor, cu variantele pentru cheie | Ioan Tuns (decizii), Claude | [D002](decizii/D002-portabilitate.md), [D008](decizii/D008-criptarea-datelor.md), [I005](idei.md), F05-T09 … T11 |
| 2026-10-24 | D002 și D006 detaliate (motiv, mecanism, verificare). Puncte noi: backup automat înainte de migrări, o singură instanță pe bază, foldere sincronizate (D002); alocarea taskurilor de către un om până la remote, conflicte în jurnal (D006) | Claude | [D002](decizii/D002-portabilitate.md), [D006](decizii/D006-reguli-agenti.md) |
| 2026-10-24 | Regulă nouă: mereu ultima versiune stabilă, actualizată la fiecare modificare. Go actualizat 1.27.0 → 1.27.3 (`toolchain go1.27.3` în `go.mod`); `go vet` și toate testele trec | Ioan Tuns (decizie), Claude | [D007](decizii/D007-versiuni-stabile.md), F02-T09 |
| 2026-10-24 | Reguli de lucru pentru toți agenții în `AGENTS.md` (`CLAUDE.md` îl importă). Cerință nouă: orice idee se documentează → registrul [idei.md](idei.md). Coloanele `Responsabil` și `Cine`. `.gitattributes` și `.gitignore` | Claude | [D006](decizii/D006-reguli-agenti.md), F02-T07, F02-T08 |
| 2026-10-24 | Primul commit cu motorul, traducerile și documentația, pe ramura `fundatie` | Claude | F02-T05 |
| 2026-10-24 | Structura de documentare: roadmap cu 14 funcționalități, taskuri, decizii D001–D005, jurnal, `CLAUDE.md` | Claude | [D005](decizii/D005-documentare.md), F02-T03, F02-T04 |
| 2026-10-24 | Traducerile extinse: folder de traduceri lângă `.exe`, limbi noi, lista traducerilor lipsă, `Text` pentru datele utilizatorului | Claude | [D004](decizii/D004-traduceri.md), F03-T03, F03-T04 |
| 2026-10-24 | Cerință nouă: aplicația trebuie să fie traductibilă, minim RO și EN. Pachetul `i18n` și erorile motorului transformate în coduri traductibile | Claude | [D004](decizii/D004-traduceri.md), F03-T01, F03-T02 |
| 2026-10-24 | Instalarea Go cu winget a eșuat (cod 1618, alt installer rula). Testele rulate cu Go 1.27.0 portabil, descărcat temporar | Claude | F02-T02 |
| 2026-10-24 | Motorul de calcul: `Segmenteaza`, `Valideaza`, calendarul și 20 de teste, toate trec | Claude | [D003](decizii/D003-motor-calcul.md), F04-T01, F04-T02 |
| 2026-10-24 | Chestionarul pentru contabilitate creat | Claude | F01-T01, [document](https://claude.ai/artifact/E8SLB7zT2WAzKTyFSchmK6) |
| 2026-10-24 | Revizuirea documentului de design: probleme de portabilitate, modelul de sporuri, alegerea bibliotecilor | Claude | [D002](decizii/D002-portabilitate.md), F04-T03 … T06, F08-T01, F09-T01 |
| 2026-10-23 | Documentul de design „Aplicația de pontaj — tehnologii și decizii” | Ioan Tuns | [PDF](Aplicația%20de%20pontaj%20—%20tehnologii%20și%20decizii.pdf), [D001](decizii/D001-stack-tehnologic.md) |
