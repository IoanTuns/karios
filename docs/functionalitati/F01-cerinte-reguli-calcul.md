# F01 · Cerințe și reguli de calcul

Etapă: V1 · Stare: în lucru · Depinde de: —

## Scop

Regulile de plată și formatul rapoartelor sunt confirmate de contabilitate și transformate în teste automate, înainte de a fixa schema și motorul.

## Decizii

- [D003](../decizii/D003-motor-calcul.md) — principiile motorului; se revizuiește după răspunsuri

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F01-T01 | Chestionar pentru contabilitate | gata | Claude | [Pontaj — întrebări pentru contabilitate](https://claude.ai/artifact/E8SLB7zT2WAzKTyFSchmK6), 2026-10-24 |
| F01-T02 | Trimite chestionarul; primește răspunsurile și un centralizator real completat pe o lună | de făcut | | Deblochează F04-T03 … F04-T05 și F04-T07 |
| F01-T03 | Transformă răspunsurile în reguli și cazuri de test | blocat | | Așteaptă F01-T02 |
| F01-T04 | Validează cu contabilul formatul centralizatorului și al exportului de salarizare | de făcut | | Coloanele T, ZNL, L/FIN, celulele roșii, chenarele |
| F01-T05 | Verifică disponibilitatea numelui Pontix (OSIM, EUIPO, domeniu) | de făcut | | Alternativa: Cadran |

## Modelul de raport primit la 2026-10-25

### Precizări 2026-10-26

Ioan Tuns cere ore în celulele de pontaj, păstrând codurile numai în date/calcul. Pontajul are două rapoarte: centralizatorul existent și unul detaliat cu Ore lucrate, T, ZNL, CO, CEV, CFS, Ore norm, Form prof, CM, L/FIN, Deleg, Ore gardă 200% la final.

Confirmări: **T = total**, **ZNL = weekend sau sărbători/zile libere, interval 00–24**, **Ore norm = toate orele de luni până vineri 00–24 când nu sunt sărbători/zile libere**. T este implementat ca total ore lucrate. CO/CM se pot număra din zilele active înregistrate. Pentru restul utilizatorul nu are încă regulile și cere o zonă de configurare folosită ulterior în calcule; nu aplicăm implicit 200% sau asocieri după abrevieri.

Pentru turele peste miezul nopții cere o configurare suplimentară. Se oferă atribuirea integrală zilei începerii sau distribuirea pe zile calendaristice, cu aceeași opțiune în grilă și rapoarte. Încadrarea ZNL/Ore norm respectă zilele calendaristice chiar când afișarea creditează tura zilei de început. Durata folosită este cea nominală configurată, în lipsa orelor de intrare/ieșire. Distribuția pauzelor, orele efective la schimbarea ceasului și formulele rămase necesită specificații; nu se inventează reguli de plată.

Ioan Tuns a confirmat funcționarea versiunii 0.10.0 și a trimis în conversație capturi de centralizator pentru octombrie 2026. Preluăm numai structura, fără numele, mărcile sau valorile reale ale angajaților: zilele 1–31, weekenduri gri, legendă T = TURA / N = NOPȚI și coloane de total Ore lucrat, T, ZNL, CO, CEV, CFS, Ore norm, Form prof, CM, L/FIN, Deleg, Ore gardă 200%. Captura mai largă arată rânduri distincte „GARDĂ”, cu mărci proprii, și grupări pe categorii profesionale; ipoteza inițială de două rânduri pentru fiecare angajat nu este confirmată.

Întrebări trimise în conversație: semnificația celor două rânduri și a codurilor C, FP, cev; confirmarea coloanelor T, ZNL, CO, CEV, CFS, Form prof, L/FIN. Nu deducem reguli de calcul sau procente din antetul raportului. Captura cu date potențial reale nu este copiată în repo.

Răspuns Ioan Tuns: specificațiile acestor coduri/coloane nu sunt disponibile momentan. Calendarul solicitat este lista completă de sărbători legale din România. Codul de schimb `N12` este acceptabil, cu necesitatea unui nomenclator configurabil; orele exacte nu sunt încă precizate. Rămâne de stabilit relația dintre persoană, rândul de gardă cu altă marcă și contract; grila inițială păstrează rânduri distincte după înregistrarea angajatului și nu calculează salarizarea.

Precizare ulterioară, 2026-10-25: setările inițiale ale schimburilor vin din PDF. Pentru codurile fără început explicit în document, Ioan Tuns a stabilit 7 = 07:00–14:00 și 12 = 14:00–02:00 a doua zi. Sursa și setul complet sunt în [D019](../decizii/D019-schimburi-initiale-pdf.md). Această precizare nu stabilește sporuri sau regulile coloanelor de salarizare.

Modulul de concedii solicitat ulterior (F07-T04): evidența numără numai zilele selectate explicit. Rămân de confirmat consumul concediului după programul individual, reportările și soldurile istorice, limitele/drepturile pentru evenimente speciale, plata și indemnizațiile medicale. Calculul automat al acestora este blocat până la răspunsul contabilității; nomenclatorul și evidența nu depind de ele.

Programe precizate ulterior (F07-T10, D021): Z7 07–14, S7 15–22, N9 22–07; IN9 distinct începe în seara dinaintea weekendului/sărbătorii, apoi opțiuni Z12 07–19, N12 19–07 și G24 07–07. Aceste cerințe înlocuiesc exemplele PDF pentru noile programe, fără reinterpretarea codurilor istorice. Rămân neprecizate plata specială pentru IN9, coloanele de total și rotația personalului; calculele aferente sunt blocate până la confirmarea contabilității.

## Clarificari F07-T25, 2026-10-31

Utilizatorul confirma plata TOTALA 200% (100% baza + 100% supliment), pentru medic/garda/weekend sau sarbatoare. Baza orara in lei, conversia ZNL in zile libere si urmarirea acordarii lor raman neprecizate; componentele monetare/compensatorii sunt blocate. Nu se presupune ca ZNL inseamna zile libere acordate. Regula nu este considerata confirmare contabila a unei interpretari legale.

## F07-T27 - Ore lucratoare / norma lunara datorata, 2026-10-31

Utilizatorul precizeaza: coloana Ore lucratoare este totalul de ore pe care fiecare angajat trebuie sa le efectueze conform contractului si numarului de zile libere. Nu reprezinta orele efectiv lucrate.

In aplicatia curenta exista Ore lucrate (total efectiv), T (total efectiv) si Ore norma (clasificarea orelor efectiv lucrate luni-vineri). Aceasta precizare necesita confirmarea coloanei vizate, pentru a nu schimba tacit alte semnificatii/teste confirmate.

Intrebari transmise: prima coloana Ore lucrate devine Ore lucratoare, T ramane efectiv, sau se adauga o coloana separata? Zilele libere sunt numai weekend/sarbatori din calendar, sau se scad si absente personale (care anume)?

Contractele zilnice/individuale furnizeaza durata zilnica; cele lunare furnizeaza tinta lunara explicita; Garda fara norma nu trebuie sa primeasca o norma dedusa. Tratamentul absentelor si prezentarea finala asteapta raspuns. Nicio formula modificata inca.

### Precizare suplimentara F07-T27: T, L FIN si L Inceput

Utilizatorul precizeaza ca T inseamna TURE: persoanele care lucreaza in ture sunt marcate in aceasta coloana a pontajului detaliat. T nu trebuie considerat implicit totalul orelor efective. Forma valorii (marcaj, ore sau numar de ture) asteapta confirmare.

L FIN este diferenta dintre totalul orelor efectiv lucrate si Ore lucratoare (norma datorata). In luna urmatoare, valoarea L FIN din luna precedenta devine L Inceput; coloana L Inceput este afisata imediat dupa Marca. Nu se presupune ca L Inceput intra automat in formula L FIN.

Intrebari transmise: ce se inscrie in T; confirmarea formulei efectiv-minus-norma fara adaugarea soldului initial si cu pastrarea valorilor negative; ce zile libere/absente reduc norma. Lipsa unei luni precedente nu va fi asimilata automat unui sold zero; soldul initial necesita o regula explicita.

Precizarea inlocuieste ipoteza anterioara ca T ar ramane totalul efectiv. Formulele existente si testele de referinta nu au fost modificate pana la raspuns.

### Definitiile coloanelor de absente, precizate de utilizator

- CEV (zile): concediu pentru evenimente remunerat, fara consumul concediului contractual; exemple deces in familie, nastere partener, casatorie si alte evenimente. Nu se stabilesc durate sau drepturi dupa tipul evenimentului.
- CFS (zile): concediu fara salar.
- Ore nem: ore nemotivate, nu zile.
- Form prof: numarul de zile de formare profesionala din luna curenta.
- CM (zile): totalul zilelor de concediu medical.
- Deleg (zile): zile de delegatie.

Utilizatorul solicita notificarea optiunilor lipsa si a clarificarilor. Verificarea codului: nomenclatorul are numai categoriile ODIHNA/MEDICAL/SPECIAL, fara atribut remunerat/neremunerat sau consum de concediu contractual. Nu exista coloana de ore nemotivate; form_prof nu are inca unitatea ZILE fixata. CEV/CFS/form_prof/deleg sunt configurabile ca surse, dar implicit nedefinite; CM este implicit numarata dupa codul CM, nu dupa toate codurile categoriei MEDICAL.

Clasificarea va necesita asociere explicita a fiecarui tip la coloana de raport, fara deducere dupa denumire. Nu se vor include automat toate absentele remunerate in CEV, intrucat CM/formare/delegatie au coloane proprii; incadrarea trebuie sa fie neambigua.

Intrebari transmise: zile calendaristice inregistrate versus zile programate la lucru (eventual reguli diferite per coloana); numai aprobate sau toate active; ore nemotivate introduse explicit versus calculate dupa o formula specificata. Numarul de zile acordat pentru fiecare eveniment nu este dedus sau limitat automat.

### Raspunsuri confirmate pentru F07-T28

Utilizatorul confirma explicit, 2026-10-31:
1. CEV/CFS/formare profesionala/CM/delegatii: se numara numai zilele in care persoana trebuia sa lucreze, nu toate zilele calendaristice inregistrate.
2. In totalurile lunare intra numai absentele aprobate. Inregistrarile inactive sau neaprobate nu se includ.
3. Orele nemotivate se introduc explicit in ore/minute; nu se deduc automat din lipsa pontajului sau dintr-o diferenta fata de norma.

Verificare tehnica: absentele au deja stare de aprobare, iar motorul are filtrul NumaiAprobate; lipseste verificarea zilei programate la numararea absentelor. Este necesara o sursa explicita de program, independenta de prezenta efectiva (absenta poate inlocui celula de pontaj).

Clarificare suplimentara transmisa: Planificare pentru ture, luni-vineri fara sarbatori pentru program normal; daca lipseste planificarea pentru ture, se semnaleaza date incomplete, fara presupunere. Aceasta propunere nu este considerata confirmata inainte de raspuns.

### Precizare utilizator: norma datorata dupa zile libere si concedii

Raspuns la intrebarea privind sursa programului: utilizatorul precizeaza ca orele totale se calculeaza din norma si numarul de zile lucratoare din luna, minus total zile libere si total concedii. Interpretarea de lucru este o norma lunara datorata, distincta de orele efectiv lucrate. Aceasta precizare nu confirma implicit propunerea anterioara Planificare/L-V pentru numararea tuturor tipurilor de absenta.

Calculul trebuie sa evite scaderea repetata a aceleiasi date: zilele deja excluse din calendarul lucrator nu se scad din nou, iar concediile aprobate suprapuse cu o zi libera nu produc o a doua reducere. Raspunsul anterior privind includerea numai a absentelor aprobate ramane valabil.

Clarificare trimisa: care dintre CO, CEV, CFS, CM, formare profesionala si delegatie reduc norma lunara si daca zilele libere compensatorii efectiv acordate se scad. Nu sunt asimilate toate rubricile de raport concediilor care reduc norma fara confirmare. Contractele lunare fixe si contractele fara norma prestabilita necesita pastrarea distinctiei fata de formula cu norma zilnica.

### Reducerile normei confirmate, 2026-10-31

Utilizatorul confirma ca norma lunara se reduce pentru CO, CEV, CFS si CM. In raspuns apare CVE, interpretat ca inversare a literelor din CEV, denumirea definita anterior. Formarea profesionala si delegatia nu fac parte din lista confirmata de reduceri. Pentru zilele compensatorii utilizatorul nu are inca o regula: nu se face reducere automata.

Se pastreaza confirmarea anterioara: numai absente aprobate, numai zile in care trebuia lucrat, fara scadere dubla a aceleiasi date. Formula zilnica poate fi implementata independent de afisarea T si reportul L FIN/L Inceput, care asteapta inca precizarile cerute.

### Clarificare utilizator 2026-10-31 - F07-T31

Programul normal L-V inseamna toate zilele lucratoare, excluzand sarbatorile legale si weekendul. Sarbatorile si weekendul necesita tratare separata pe cazurile specifice. Aceasta clarificare nu stabileste procente noi, reguli de cumulare a sporurilor sau conversia ZNL in timp liber. Regula zilelor compensatorii ramane neprecizata.

### Clarificare utilizator 2026-10-31 - F07-T33

Toate concediile/absentele, inclusiv medicalul, se inregistreaza numai in zile lucratoare L-V fara sarbatori legale. Prioritate explicita: MEDICAL peste orice alt concediu, orice concediu peste program. Se suprascriu automat inregistrarile vizate, cu jurnal. Ceruta regula pentru tura de noapte care incepe in ziua precedenta: intreaga tura sau numai orele zilei de CM.

F07-T33: utilizatorul confirma anularea intregii ture de noapte care intra intr-o zi de medical, inclusiv programul inceput in ziua precedenta. Medicalul anuleaza programele tuturor marcilor asociate explicit aceleiasi persoane; inregistrarea absentei ramane pe marca selectata.

F07-T34: utilizatorul confirma ca toate operatiunile, inclusiv concediile, trebuie sa poata fi corectate sau sterse. Pentru absente, corectia explicita cu revizie este permisa inclusiv din medical in alt tip; anularea pastreaza istoricul. Prioritatea automata nu face inregistrarea imuabila.


### TODO utilizator 2026-11-02 - F09-T14, coloane calculate

Cerute totaluri contractuale in zile/ore, total concedii configurabil, coloana automata pe fiecare tip nou de concediu, total ore garda si ore garda weekend/sarbatori configurabile, plus coloane calculate proprii din aplicatie. [Specificatie completa](../todo-rapoarte-configurabile.md).

Intrebari deschise: formula numita Total ore lucratoare este exprimata zile contractuale minus zile concediu (rezultat in zile); trebuie confirmata conversia prin norma sau schimbata denumirea. Intervalul garzii este descris 00 ziua de intrare - 24 ora de iesire: trebuie confirmata intersectia cu orele efective sau perioada extinsa, atribuirea intre luni si evitarea dublei numarari weekend/sarbatoare. Trebuie stabilita selectia tipurilor/categoriilor din totalul concediilor, separat de efectul lor asupra obligatiei. Nu se schimba formula implementata si nu se presupun reguli noi de plata.
