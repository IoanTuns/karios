# D019 · Schimburi inițiale din PDF și precizările utilizatorului

Stare: propusă · Data: 2026-10-25 · Taskuri: F07-T09

Înlocuiește propunerea D015 numai în privința absenței unui set inițial: Ioan Tuns a cerut acum folosirea valorilor din document și a precizat orele lipsă pentru 7 și 12. Validările și editarea ulterioară rămân aceleași.

## Surse și cerințe confirmate

[PDF-ul local de design](../Aplicația%20de%20pontaj%20—%20tehnologii%20și%20decizii.pdf), 2026-10-23:

- Pagina 2: exemplele S1 06–14 și S3 22–06.
- Pagina 4: exemplul de noapte 19–07 și garda de 24 ore, 07–07; nu preluăm procentele de plată ca reguli confirmate.
- Pagina 5: codurile 7, 12, N12, G24 pentru introducere rapidă.

Ioan Tuns, în conversația din 2026-10-25: schimbul de 7 ore începe la 07:00 și durează 7 ore; schimbul de 12 ore începe după terminarea celui de 7 ore. Rezultă 07–14 și 14–02 în ziua următoare.

| Cod | Început | Sfârșit nominal | Durată | Tip |
| --- | --- | --- | --- | --- |
| 7 | 07:00 | 14:00 | 420 min | NORMAL |
| 12 | 14:00 | 02:00 (+1 zi) | 720 min | TURA |
| N12 | 19:00 | 07:00 (+1 zi) | 720 min | NOAPTE |
| G24 | 07:00 | 07:00 (+1 zi) | 1440 min | GARDA |
| S1 | 06:00 | 14:00 | 480 min | TURA |
| S3 | 22:00 | 06:00 (+1 zi) | 480 min | NOAPTE |

Pauza nu este precizată în PDF. Păstrăm valoarea tehnică implicită deja existentă de 0 minute, editabilă; aceasta nu confirmă o regulă de plată sau de scădere a pauzei. Tipurile și asocierea codurilor N12/G24 la exemple sunt reprezentarea tehnică propusă.

## Implementare propusă

- Set JSON inclus în aplicație și migrarea de date `0007_schimburi_initiale.sql`.
- Adăugăm numai codurile care nu există, ignorând majusculele; păstrăm orice configurație existentă, inclusiv inactivă ori ștearsă logic. Nu revenim automat peste editările utilizatorului.
- Completarea explicită generează ID-uri noi: dacă un șablon a fost redenumit, recrearea codului lipsă păstrează separat schimbul personalizat.
- Interfața permite și completarea explicită a codurilor lipsă, tranzacțional și cu jurnal. Arată sursa setului și sfârșitul nominal calculat din început și durată, inclusiv +1/+2 zile.
- Intervalele sunt șabloane nominale de schimb. Transformarea în momente calendaristice exacte și calculul la schimbarea orei rămân în motorul de calcul, separat de acest task.

Starea propusă privește reprezentarea tehnică. Orele pentru 7 și 12 sunt cerințe explicite ale lui Ioan Tuns și nu au fost inventate de agent.
