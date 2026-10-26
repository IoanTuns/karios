# D021 · Programe de lucru și intrarea în zile libere

Stare: propusă · Data: 2026-10-25 · Task: F07-T10

Ioan Tuns a precizat programele reale: 07–14, 15–22, 22–07. Pentru weekenduri/sărbători se începe în seara anterioară cu 22–07, cu statut special și cod separat, apoi schimburi de 12/24 ore. La clarificare a răspuns „da” pentru începutul în seara precedentă și exemplele 07–19, 19–07, gardă 07–07, respectiv „ok” pentru codurile Z7/S7/N9/IN9 și evidențierea separată. Nu a precizat reguli suplimentare de plată.

| Cod | Interval nominal | Durată | Regim |
| --- | --- | --- | --- |
| Z7 | 07:00–14:00 | 7 ore | LUCRATOARE |
| S7 | 15:00–22:00 | 7 ore | LUCRATOARE |
| N9 | 22:00–07:00 (+1 zi) | 9 ore | LUCRATOARE |
| IN9 | 22:00–07:00 (+1 zi) | 9 ore | INTRARE_ZI_LIBERA |
| Z12 | 07:00–19:00 | 12 ore | NELUCRATOARE |
| N12 | 19:00–07:00 (+1 zi) | 12 ore | NELUCRATOARE |
| G24 | 07:00–07:00 (+1 zi) | 24 ore | NELUCRATOARE |

## Reprezentare propusă

- Migrare nouă 0009: câmp `regim` separat de `tip`. Tipul NOAPTE/GARDA descrie activitatea; regimul distinge programul obișnuit, weekend/sărbătoare și intrarea specială. Valoarea GENERAL păstrează schimburile anterioare neclasificate.
- Codurile noi se adaugă numai dacă lipsesc, fără coliziuni cu absențele. Codurile vechi 7/12/S1/S3 se păstrează și nu își schimbă sensul. Setul de completare devine cel de mai sus; documentele D019 și migrarea 0007 rămân istorice.
- Pentru șabloanele originale N12/G24 se completează numai regimul, dacă întregul șablon și marcajul de nemodificare corespund valorii livrate inițial. Schimburile personalizate, inactive, redenumite sau șterse logic nu sunt suprascrise.
- Regimul este editabil explicit în formular și vizibil/filtrabil în listă și în sugestiile din pontaj. IN9 poate fi recunoscut și după cod, și după regim, fără procente sau categorii de plată deduse.
- Pontajul păstrează ziua începerii turei: IN9 pentru sâmbătă începe vineri la 22:00. Intervalele de 12 și 24 ore sunt opțiuni de acoperire; nu alocăm automat ture unor angajați și nu adunăm o gardă de 24 ore peste două ture de 12.
- Pauza rămâne valoarea tehnică editabilă 0. Golul 14–15 este explicit și nu este interpretat ca o pauză scăzută din ture. Duratele sunt nominale; calculul efectiv la schimbarea orei rămâne în motor, cu regulile confirmate separat.

Regimul și codurile sunt evidență operațională. Plata specială pentru IN9, repartizarea pe coloanele raportului și rotația personalului așteaptă F01. Reprezentarea tehnică rămâne propusă; orele și distincția IN9 sunt cerințe explicite.
