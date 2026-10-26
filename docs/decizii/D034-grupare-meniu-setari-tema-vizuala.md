# D034 · Gruparea ecranelor de configurare sub Setări și temă vizuală distinctă

Data: 2026-10-26 · Stare: acceptată · Decisă de: Ioan Tuns
Legături: F05-T17, I044

## Context

Meniul lateral principal conținea 9 intrări de același nivel (`Angajați`, `Pontaj`, `Concedii și absențe`, `Sărbători`, `Schimburi`, `Rapoarte`, `Definire rapoarte`, `Configurare pontaj`, `Setări`), amestecând operațiunile zilnice/lunare cu nomenclatoarele și regulile de configurare. Utilizatorul a cerut:

1. Mutarea ecranelor de configurare (`Sărbători` redenumit în `Zile libere legale`, `Schimburi`, `Configurare pontaj` și `Definire rapoarte`) în interiorul meniului **Setări**, cu acces atât prin taburi în ecranul **Setări**, cât și prin sub-meniu vizibil sub **Setări** în bara laterală când zona **Setări** este activă.
2. Folosirea unei palete de culori distincte (albastru-ardezie instituțional `#2b4c7e` / `#eef3f9`) pentru meniul și zona de **Setări**, pentru a o diferenția clar de verdele operațional (`#17655a`).

## Variante

### Varianta A — 5 intrări principale în bara laterală + sub-meniu și taburi în Setări + temă albastru-ardezie

- Bara laterală principală păstrează 5 intrări principale: **Angajați**, **Pontaj**, **Concedii și absențe**, **Rapoarte** și **Setări**.
- Când utilizatorul se află în **Setări** (sau în oricare dintre sub-secțiunile sale), sub intrarea **Setări** din bara laterală apare sub-meniul cu cele 5 secțiuni (`General`, `Zile libere legale`, `Schimburi`, `Configurare pontaj`, `Definire rapoarte`), iar în partea de sus a ecranului apar taburile corespunzătoare.
- Denumirea `ui.sarbatori` în română devine **Zile libere legale** (`Public holidays` în engleză).
- Zona și meniul **Setări** folosesc accentul albastru-ardezie (`#2b4c7e`, fundal activ `#eef3f9`).

### Varianta B — Meniu plat cu 9 intrări în bara laterală

- Respiniă de utilizator: aglomerează meniul principal cu ecrane accesate rar.

## Decizie

Se aplică **Varianta A**, conform alegerii explicite a utilizatorului.

## Consecințe

- Ecranul principal rămâne curat, cu doar 5 destinații de nivel 1.
- Orice modificare nesalvată din `Configurare pontaj` rămâne protejată cu confirmare la schimbarea secțiunii sau a ecranului.
- Culoarea albastru-ardezie (`#2b4c7e`) semnalează imediat că utilizatorul se află în zona de administrare și configurare a aplicației.

