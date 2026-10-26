# D028 · Ore în pontaj, două rapoarte și configurarea calculelor

Stare: propusă · Data: 2026-10-26 · Task: F09-T09

Ioan Tuns cere ore afișate în pontaj, codurile rămânând în date/calcul, un centralizator și un raport detaliat cu coloanele finale din captura trimisă. Confirmă T total, ZNL weekend/sărbători/libere 00–24 și Ore norm luni–vineri fără zile libere. Pentru coloanele rămase cere configurare, fără reguli disponibile încă. Captura nu este copiată în repo.

Propunere: codul selectează programul; grila și PDF-ul afișează durata nominală în ore, fără să rescrie codul. Editarea alege programul sau ore introduse explicit, nu deduce programul dintr-un număr ambiguu (de exemplu N9/IN9). Copierea exportă ore. Valorile istorice neinterpretate sunt marcate „—”; absențele au ore zero, cu detalii în explicația celulei. Calculele sunt pure în internal/calcul și au teste.

Configurare globală de pontaj, cu revizie și jurnal: atribuirea turelor integral datei începerii (implicit compatibil cu schițele existente) sau pe zile calendaristice. ZNL/Ore norm se clasifică pe intervale calendaristice 00–24; în primul mod totalul perioadei include turele începute în ea, în al doilea numai fragmentele din ea, inclusiv ture din luna precedentă. Calendarul include zilele libere definite de utilizator.

Regulile celor 12 coloane pot alege surse limitate: nedefinită, total ore, ore ZNL, ore norm, ore pentru schimburi selectate sau zile de absență pentru coduri selectate. Orele pot avea un coeficient introdus explicit de om; pentru gardă nu se precompletează sursa sau factorul 2 din etichetă. CO/CM numără zilele active înregistrate. Sursele trebuie să respecte unitatea coloanei; pentru absențe se poate alege numai aprobate. Fără formule executabile, SQL sau interpretări legale generate de agent. Nedefinit înseamnă „—”, nu zero. Normele/indemnizațiile/soldurile legale rămân în F01.

Migrarea 0015 adaugă configurarea și un model separat PONTAJ_DETALIAT; nu schimbă antetul, subsolul sau imaginile centralizatorului existent. Fișa individuală înlocuiește coloana de cod cu ore. Rapoarte: centralizator v3, fișă individuală v3, pontaj detaliat v1; evidența absențelor rămâne v2. Referințele vechi se păstrează.

Orele manuale se salvează explicit în minute, în formatul intern `@M450` pentru 7,5 ore. Codurile programelor au prioritate față de valorile numerice istorice; `@` nu este permis în codurile schimburilor. Orele fără program se atribuie zilei selectate. Copierea folosește până la zece zecimale, suficient pentru relipire cu precizie de minut; afișarea folosește cel mult patru.

Calculul folosește minute nominale, fără ajustări pentru ora de vară sau poziționarea pauzelor. Programele cu pauză nenulă rămân necalculabile; lipsa orei de început permite numai totalul nominal în modul zilei de început. O absență activă elimină programul propriei zile; continuarea unei ture începute anterior rămâne vizibilă, fără a inventa reguli de compatibilitate concediu/muncă. Clasificarea ia toate zilele libere din calendar, inclusiv din anul următor. Configurația are revizie și audit atomic, iar PDF-ul include revizia. Coloanele finale rămân împreună la sfârșitul raportului detaliat; toate zilele au aceeași lățime.

Lucrul pornește din reunirea F02-T11, iar publicarea/PR-ul rămân în așteptarea răspunsului privind vizibilitatea GitHub. După finalizare, integrarea va fi avansată la noua funcționalitate, fără publicare pe repo-ul public și fără modificarea datelor operaționale.
