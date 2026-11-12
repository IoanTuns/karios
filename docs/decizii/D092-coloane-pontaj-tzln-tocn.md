# D092 · Coloanele Pontajului, TZLN și TOCN

Stare: acceptată (regulile) · propusă (detaliile tehnice) · Data: 2026-11-07 · Taskuri: F08-T11, F08-T12, F08-T13

Implementare: implementată din 0.92.57; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

Decisă de: Ioan Tuns, 2026-11-07, pentru: definițiile de mai jos, coloanele păstrate în Pontaj și redenumirea codurilor (răspunsuri explicite la întrebări în conversație).

## Context

Grila Pontaj are 12 coloane de total (Ore lucrate, T, ZNL, CO, CEV, CFS, Ore normă, Form. prof., CM, L/FIN, Deleg., Ore gardă 200%), configurate în Setări › Configurare pontaj. Ioan Tuns a cerut un set nou, cu definiții precise, și o coloană calculată ascunsă pentru orele datorate pe lună. Aceste definiții numără zile și ore; nu stabilesc procente, sporuri sau plata (acelea rămân la contabilitate, AGENTS.md §2).

## Decizie

### Definiții

| Cod | Denumire | Regulă |
| --- | --- | --- |
| TZLN | Total zile lucrătoare nete | Zilele lunii, minus weekendurile, minus sărbătorile legale care cad luni–vineri. Ascunsă. |
| TOCN | Total ore calculat | TZLN × norma zilnică a angajatului (7 h la normă întreagă implicit, orele lui la normă parțială). Fixă pe lună: concediile nu o scad. Ascunsă în Pontaj și în rapoarte. Fără normă zilnică (gărzi, contracte fără normă zilnică) rămâne necalculată, cu motiv. |
| OL | Ore lucrate | Orele salvate în Pontaj (programe și ore manuale) plus orele creditate ale zilelor FMP. DLG nu intră. Ce e nesalvat nu se numără. |
| CO | Zile concediu | Zilele cu codul CO. |
| CM | Concediu medical | Zilele cu codul CM. |
| CFP | Concediu fără plată | Zilele cu codul CFP. |
| CVE | Concediu evenimente speciale | Zilele oricărui alt tip de absență, în afară de CO, CM, CFP, DLG și FMP. |
| DLG | Delegații | Zilele cu codul DLG. |
| FMP | Formare profesională | Zilele cu codul FMP; orele lor intră și în OL. |

### Afișare

Grila Pontaj păstrează doar OL, CO, CM, CFP, CVE, DLG și FMP. TZLN și TOCN se calculează, dar nu se afișează în grilă; sunt disponibile pentru rapoarte, în SQL. Coloanele vechi (T, ZNL, Ore normă, L/FIN, Ore gardă etc.) rămân în calcul pentru rapoartele și exportul care le folosesc.

### Coduri

Codurile de absență se redenumesc: CEV → CVE, CFS → CFP, DEPL → DLG, FP → FMP. Redenumirea se face o singură dată, printr-o migrare, numai când codul vechi există și cel nou nu; zilele înregistrate trec pe codul nou. Regulile salvate care pomenesc codurile vechi sunt traduse la citire.

### Detalii tehnice (propuse)

- Coloanele noi se adaugă în calcul (`internal/calcul`), cu surse noi: ore lucrate cu codurile creditate alese (OL), zilele altor absențe decât cele excluse (CVE), zile lucrătoare nete (TZLN), TZLN × normă (TOCN). Toate rămân configurabile în Configurare pontaj.
- O configurație salvată fără coloanele noi le primește automat la citire, doar cu codurile care există în bază; o coloană fără niciun cod existent rămâne „neconfigurată”, cu motivul afișat (F05-T20).
- Reducerea normei lunare (D063, D077) folosește și codurile noi (CVE, CFP), pe lângă cele vechi.

## Variante respinse

- **Ștergerea coloanelor vechi din calcul:** ar strica rapoartele standard, exportul Excel și soldurile, care le folosesc; s-au ascuns doar din grilă.
- **Înlocuirea textuală a codurilor în JSON-ul regulilor și al rapoartelor:** codurile vechi (CEV, CFS) sunt și identificatori de coloană; o înlocuire textuală le-ar redenumi și pe acelea.

## Consecințe

- Interogările SQL scrise de mână în rapoartele flexibile care conțin codurile vechi ('CEV', 'CFS', 'DEPL', 'FP') trebuie actualizate de utilizator; listele de tipuri alese în modele se traduc automat.
- Dacă baza conține deja și codul vechi și pe cel nou, migrarea nu le amestecă; ambele rămân, iar utilizatorul le unifică manual.
