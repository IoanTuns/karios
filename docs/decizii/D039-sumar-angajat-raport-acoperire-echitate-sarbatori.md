# D039 · Sumar editabil pe angajat, raport de acoperire weekend și verificare echitate sărbători

Stare: acceptată · Data: 2026-10-26 · Propusă de: Antigravity · Decisă de: Ioan Tuns · Task: F11-T04, F11-T05, F11-T06

## Context

În gestionarea planificării gărzilor și a turelor de weekend/după-amiază, utilizatorul a solicitat trei capabilități critice pentru organizarea optimă a muncii:
1. **Sumar editabil per angajat:** Posibilitatea de a deschide, pentru oricare angajat, un sumar dedicat care centralizează toate turele planificate, totalul de ore și permite modificarea rapidă direct din sumar.
2. **Raport complet de acoperire weekend / zile nelucrătoare:** Zilele de weekend și sărbătorile legale nu pot rămâne neacoperite. Este necesar un raport și o vedere dedicată care verifică prezența personalului pe posturile obligatorii:
   - Gardă de 24 ore (`G24`)
   - Tura de zi weekend 12 ore (`Z12`)
   - Tura de noapte weekend 12 ore (`N12`)
   - Intrare weekend / zi liberă (`IN9`)
   și semnalează vizibil orice zi sau post neacoperit.
3. **Echitate și rotație la sărbători majore:**
   - Cine a fost planificat de Crăciun în anul anterior (`N-1`) este recomandat/prioritar pentru Revelion în anul curent (`N`).
   - Cine este planificat de Crăciun în anul curent (`N`) nu poate/nu ar trebui să fie planificat de Paști în anul viitor (`N+1`) pe exact același schimb.
   - Dacă o persoană are ture consecutive pe zile libere/gărzi (sau gărzi apropiate), se afișează un marcaj vizual distinctiv.

## Decizie

1. **Sumar editabil pe angajat:**
   - Accesibil atât ca dialog modal direct din grilă (prin clic pe numele angajatului sau buton dedicat), cât și ca vedere dedicată.
   - Afișează datele angajatului, statistici detaliate (număr de gărzi G24, ture Z12/N12, intrări IN9, după-amiezi S7, total ore planificate) și lista completă a zilelor cu schimburi atribuite.
   - Permite adăugarea, modificarea sau eliminarea rapidă a oricărui schimb pentru acel angajat, cu salvare tranzacțională.

2. **Raport și panou de acoperire a zilelor nelucrătoare:**
   - Vedere dedicată `Acoperire weekend și sărbători` integrată în ecranul de Planificare.
   - Analizează automat toate sâmbetele, duminicile și sărbătorile legale din luna selectată.
   - Verifică acoperirea fiecărui post cheie (`G24`, `Z12`, `N12`, `IN9`).
   - Evidențiază cu alertă roșie/insignă discretă orice post neacoperit și indică numele persoanei planificate pe posturile acoperite.
   - Include un rezumat statistic al gradului de acoperire al lunii (ex: `40/40 posturi acoperite` sau `38/40 posturi acoperite · 2 neacoperite`).

3. **Reguli de echitate și marcaje de ture consecutive:**
   - Verificările sunt **avertismente vizuale neblocante** (marcate cu indicator de avertisment / atenție și text explicativ în inspecție și sumar), permițând salvarea dacă planificatorul decide derogări.
   - Sistemul corelează automat istoricul anilor anteriori și viitori pentru sărbătorile de Crăciun (25–26 decembrie), Revelion (31 decembrie – 1–2 ianuarie) și Paști (Vinerea Mare, Duminica și Lunea Paștelui).
   - Pentru ture consecutive (gărzi de 24h consecutive sau gărzi urmate imediat de alte ture ori două zile de weekend consecutive fără repaus), celulele primesc un marcaj de alertă.

