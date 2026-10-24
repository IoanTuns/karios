# D007 · Mereu ultima versiune stabilă

Stare: acceptată · Data: 2026-10-24 · Decisă de: Ioan Tuns · Taskuri: F02-T09, F05-T01

## Context

Proiectul depinde de Go, Wails, React, biblioteci și unelte care evoluează. Documentul de design numea versiuni fixe (Wails v2, maroto v2), care se învechesc.

## Decizie

- Se folosește mereu cea mai nouă versiune stabilă (lansare finală, fără alpha, beta, rc) a fiecărei unelte și dependențe.
- La fiecare modificare se verifică dependențele și se actualizează cele rămase în urmă, în commit separat, cu testele rulate și cu jurnal.
- Versiunile se fixează exact în fișierele proiectului, ca build-ul să fie reproductibil.
- Versiunile numite în [D001](D001-stack-tehnologic.md) și în documentul de design nu mai sunt obligatorii; contează doar biblioteca aleasă.

Regulile de aplicare sunt în [AGENTS.md](../../AGENTS.md), secțiunea 6.

## Variante respinse

- **Versiuni fixe, actualizate rar** — actualizările se adună și devin greu de făcut; remedierile de securitate întârzie.
- **Intervale de versiuni** (`^1.2`) — două calculatoare pot construi cu versiuni diferite.

## Consecințe

- Puțină muncă la fiecare task, în locul unor migrări mari, rare.
- O versiune majoră nouă poate cere modificări de cod; se tratează ca task.
- Prima aplicare: Go actualizat la 1.27.3 (F02-T09).
