# D020 · Concedii și absențe pe zile

Stare: propusă · Data: 2026-10-25 · Task: F07-T04

Ioan Tuns cere un modul pentru concediu de odihnă, zile libere speciale și concediu medical; tipurile speciale se definesc după nevoie.

## Reprezentare propusă

- Extindem nomenclatorul existent `coduri_absenta` cu categoriile ODIHNA, MEDICAL, SPECIAL; adăugăm CO/CM numai dacă aceste coduri nu sunt deja ocupate. Tipurile speciale au cod și denumiri RO/EN introduse de utilizator; fără biblioteci noi.
- O înregistrare reprezintă o zi explicit selectată pentru un angajat și un tip. Formularul permite intervale și alegerea zilelor din ele; nu deducem zilele consumate din programul contractual. Weekendurile și sărbătorile sunt vizibile. Selectarea luni–vineri fără sărbători este o acțiune explicită, nu o regulă de salarizare.
- Migrarea este `0008_concedii_absente`; codurile vechi sunt păstrate cu categoria tehnică SPECIAL până la verificarea utilizatorului. Pontajele istorice din tabelul `pontaj` și textele din schițe nu sunt reinterpretate automat ca înregistrări ale modulului.
- Salvarea unui lot este atomică, cu revizii pentru conflicte. O singură absență activă pe angajat/zi. Editarea și anularea păstrează jurnalul; anularea nu șterge fizic date. Tipurile se dezactivează fără pierderea istoricului; codul și categoria unui tip folosit nu se schimbă.
- Pontajul afișează absențele salvate și protejează celulele lor. Nu suprascriem schițe necompletate sau cu text existent; o schiță nevidă ori un pontaj efectiv pe ziua vizată blochează salvarea absenței. Salvarea schiței verifică și absențele din backend, inclusiv dacă s-au adăugat după încărcarea grilei.
- Sumarul anual numără zilele înregistrate pe categorii și afișează separat dreptul contractual curent introdus la angajat. Nu produce sold legal, reportări, plăți sau indemnizații medicale. Nu se colectează diagnostice ori documente medicale.
- Verificarea turelor efective care traversează ziua folosește limitele calendaristice Europe/Bucharest convertite în UTC, inclusiv schimbarea orei. `time/tzdata` din biblioteca standard Go este inclusă în executabil pentru portabilitate; nu adăugăm o bibliotecă externă.

## Limite și pași următori

Regulile de plată, zilele consumate în funcție de program și istoricul drepturilor contractuale așteaptă contabilitatea (F01 și I010). Aceste calcule rămân blocate, însă evidența zilelor și configurarea tipurilor pot funcționa independent. Codurile istorice C/FP/cev nu sunt reinterpretate automat.
