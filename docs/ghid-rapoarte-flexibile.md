# Rapoarte flexibile - SQL, coloane calculate si aspect

Disponibil din 0.92.25. Setari > Definire rapoarte > Adauga model de raport sau Duplica. Alegeti sursa noua, completati denumirea si verificati PDF-ul inainte de salvare. Modificarea definitiei retrage previzualizarea; este necesara o noua verificare pentru salvare.

## Sursele noi

- Lista garzilor - un rand pe garda: nume, marca, data, program, inceput/sfarsit, ore si ore eligibile.
- Totaluri lunare si coloane calculate: un rand pe marca, cu totalurile contractuale, absentele si garzile.
- Raport SQL personalizat: SELECT pe seturile de raportare afisate de butonul Campuri disponibile pentru raport.

SQL gol foloseste interogarea implicita a sursei. Ordinea si aliasurile din SELECT stabilesc coloanele tabelului. Pentru alte ordini folositi SELECT explicit, ORDER BY si GROUP BY. Filtrele din ecranul Rapoarte se aplica datelor inainte de SQL; o interogare nu poate recupera marcile excluse. Selectia unei marci este exacta, nu adauga automat marca principala ori marcile de garda.

## Exemplu - angajatul si data garzii pe fiecare rand

```sql
SELECT angajat_id, angajat, marca, data, program,
       inceput, sfarsit, ore, ore_eligibile
FROM garzi
ORDER BY angajat, data
```

Marcati angajat_id in campul Coloane ascunse in tabel daca nu doriti identificatorul pe document. Ramane disponibil calculelor si sablonului. Gardele care intersecteaza inceputul/sfarsitul lunii au orele limitate la luna, chiar daca data intrarii este in luna precedenta. garzi_zile ofera cate un fragment pentru fiecare zi intersectata.

## Seturile de date

| Tabel | Continut |
| --- | --- |
| angajati | ID, nume afisat cu marcaj de garda, nume/prenume, marca, functie, sectie/compartiment si ID-uri, contract, norma zilnica, activ, concediu anual, marca principala |
| lunar | Totaluri contractuale, concedii selectate, ore lucratoare/lucrate, ore garda/eligibile, L_INCEPUT/L_FIN si coloanele concediu_* |
| absente | Marca interna, data, cod, categorie, activ si stare de aprobare |
| tipuri_absenta | Cod, denumire, categorie, natura, efect asupra orelor, stare |
| categorii_absenta | Cod, denumire, natura, efect, stare |
| programe | ID, cod, denumire, ora de inceput in minute, durata, tip/regim, stare |
| planificare | Marca interna, data, program si sectiune, stare; include doua zile anterioare lunii pentru ture traversante |
| pontaj_schita | Valoarea salvata pe marca/data in luna |
| pontaj_efectiv | Intervalele active de inceput/sfarsit, fara cele retrase |
| garzi | Intervale active de garda, o inregistrare pe interval, orele lunii si cele eligibile |
| garzi_zile | Fragmentele garzilor pe fiecare zi a lunii, orele si indicatoarele weekend/sarbatoare |
| calendar | Data, weekend, sarbatoare si zi lucratoare pentru luna |

Schema exacta si campurile sunt afisate in editor. Nu sunt disponibile conturi, parole, chei, documente atasate sau jurnalul operational. Datele afisate sunt o vedere de raportare, nu schema bruta integrala a bazei.

Parametri SQL: :an, :luna, :angajat_id, :sectie_id, :compartiment_id, :de_la si :pana_la. Ultimul este prima zi a lunii urmatoare, limita exclusiva. Parametrii se leaga separat, fara concatenarea valorilor in SQL.

## Totaluri calculate

- total_zile_cf_contract: zile luni-vineri fara sarbatori, conform normarii contractuale zilnice. Nu se deduc concediile aici.
- total_ore_cf_contract: zilele contractuale x norma zilnica.
- total_ore_lucratoare: rezultatul comun al obligatiei lunare, cu scaderea concediilor aprobate care reduc norma. Deplasarile cu ore creditate nu se scad.
- total_ore_lucrate: rezultatul motorului comun de pontaj, inclusiv absentele aprobate cu ore creditate.
- total_zile_concediu: tipurile/categoriile bifate in editor, numai zile lucratoare active aprobate, fara dublare pe marca/data. Nicio selectie inseamna zero. Acest total informativ nu rescrie obligatia lunara.
- concediu_CO, concediu_CM si alte concediu_*: se creeaza automat in catalog pentru fiecare tip, inclusiv cele inactive pentru istoric. Underscore si cratima din cod se codifica _5f si _2d ca numele SQL sa fie unice; folositi numele afisat in schema.
- total_ore_garda: suma intervalelor active efective, iar in lipsa lor programul activ pentru aceeasi marca/data. Efectivele si planul nu se aduna de doua ori.
- total_ore_garda_weekend_sarbatori: numai intersectia cu zilele bifate (weekend, sarbatori sau ambele), 00-24, o singura data daca ziua indeplineste ambele conditii.

Contractele fara norma, inclusiv garda, au totalurile contractuale nedefinite, afisate cu liniuta, nu cu zero presupus. Se utilizeaza normele contractuale curente disponibile aplicatiei; aceasta functie nu introduce istoricul contractelor sau o perioada de angajare partiala inexistenta in modelul actual. Intervalele efective respecta schimbarea orei de vara/iarna; planificarea foloseste durata nominala.

## Coloane proprii

Adauga coloana calculata: cod unic, denumiri RO/EN, unitate Ore/Zile/Numar, precizie 0-6 zecimale, activ si formula ori SQL scalar. Coloanele se calculeaza in ordinea afisata; sagetile permit mutarea. Datele lipsa se propaga, iar zero este o optiune explicita pentru rezultat. Nu exista executie Go/JavaScript sau formule circulare.

```text
ROUND(total_ore_lucrate - total_ore_lucratoare, 2)
```

Operatii: + - * / si paranteze; ABS(x), MIN(x,y), MAX(x,y), ROUND(x,n). SQL scalar trebuie sa returneze un singur rand si o singura valoare numerica/NULL. Exemplu, daca angajat_id exista in SELECT-ul principal:

```sql
SELECT COUNT(*)
FROM absente a
JOIN calendar c ON c.data = a.data
WHERE a.angajat_id = :angajat_id
  AND a.activ = 1 AND a.stare_aprobare = 'APROBATA'
  AND c.lucratoare = 1
```

In coloana SQL, :angajat_id este marca interna a randului daca este disponibila, altfel filtrul explicit. Fara niciuna, parametrul este refuzat; includeti angajat_id si ascundeti-l din tabel. Celelalte coloane ale randului pot fi folosite drept parametri. Formulele custom nu modifica regulile oficiale ale pontajului sau platii.

## Aspectul

> Istoric: din 0.92.42 (F09-T23) aspectele Tabel automat, Markdown si Pozitionare libera de mai jos nu mai sunt desenate de aplicatie. Modelele existente se convertesc automat in editorul vizual (vezi „Conversia modelelor vechi” la final); descrierea ramane pentru intelegerea conversiei.

Tabel automat: coloanele urmeaza SELECT-ul, apoi coloanele calculate. Antetele se repeta la continuarea pe alte pagini.

Markdown: titluri #, paragrafe, liste cu -, evidentiere **text**, tabele simple cu | si substitutia {{tabel}} pentru lista rezultata. Subset documentat de Markdown, fara HTML/CSS sau scripturi executabile.

```md
# {{titlu}}
Perioada: {{perioada}}

{{tabel}}
```

Pozitionare libera: Adauga element, alege text/campuri si coordonatele X/Y, latime/inaltime in mm fata de marginea utila, font in puncte. Elementele se trag pe plansa sau se muta prin campurile numerice. {{tabel}} plaseaza lista completa; alte campuri, de exemplu {{angajat}} sau {{marca}}, iau prima inregistrare in lipsa repetarii.

Pentru o structura proprie pe fiecare rand de date bifati Repeta elementul pentru fiecare rand si configurati distanta dintre randuri. Campurile repetate primesc valorile fiecarei inregistrari si continua automat pe alte pagini. Nu combinati tabelul integral cu elemente repetate, pentru a nu duplica datele. Celulele trebuie sa incapa in dimensiunile alese; reduceti textul, cresteti dimensiunile sau schimbati pagina daca previzualizarea refuza aspectul.

Antetul/subsolul configurat si imaginile sunt pastrate; explicatiile SQL/formulelor sunt optionale prin bifa din Continut, iar istoricul tehnic este la final. Campurile ascunse raman disponibile sablonului. Excel exporta valorile si coloanele tabelului rezultat, nu coordonatele PDF. Investigatia anterioara a blocarii Excel pe sistemul utilizatorului ramane amanata, nu este declarata rezolvata de acest task.

## Limite si verificare

O interogare SELECT, maximum 80 coloane si 10000 randuri, 8 MiB rezultat, 256 KiB celula, doua secunde pe interogare; pregatirea are termen total 20 secunde. Sunt permise citirea si functiile SQL uzuale enumerate de motor, nu ATTACH, PRAGMA, scrieri sau functii de fisiere. Un calcul SQL per rand este limitat si ca numar de executii (2000 per raport).

Verifica datele si previzualizeaza PDF este obligatoriu in editor inainte de salvarea definitiei curente. Dupa salvare, modelul activ apare in Rapoarte pentru generare cu filtrele alese, PDF, tiparire si exporturile existente. Nu se tipareste automat.

### Explicatii optionale in raport

In Setari > Definire rapoarte > Editeaza > Continut, bifa Afiseaza modul de calcul si semnificatia coloanelor include explicatiile in anexa finala PDF si in Excel. Debifati pentru un raport fara explicatii, formule sau SQL. Salvati modelul; la modelele flexibile refaceti previzualizarea inainte de salvare. Modelele noi pornesc debifate; duplicarea pastreaza alegerea, iar modelele existente pastreaza afisarea pana la debifare. Istoricul tehnic ramane separat la final.

## Lista lunara de garzi: data, ziua, nume si observatii

Alegeti Lista garzilor, luna/anul dorit si filtrele pentru toate persoanele. Un rand pentru fiecare garda care incepe in intervalul selectat. Implicit acesta este luna aleasa. Cu extinderea zilelor libere, parametrii de_la/pana_la delimiteaza grupul complet de la margini; pana_la este exclusiv.

Ordinea coloanelor este Data (numarul zilei din luna), Ziua (numele romanesc al zilei saptamanii), Nume si prenume, Observatii (ora de start - ora de final, format24h). Data iesirii este adaugata in paranteze daca garda se termina in alta zi. In aceasta varianta nr_crt si marca nu sunt afisate, conform noii liste de coloane cerute. Sortarea ramane dupa data completa si ora, apoi nume/marca pentru egalitati; nu se sorteaza dupa eticheta zilei sau dupa numarul zilei ca text.

```sql
SELECT
  CAST(strftime('%d', g.data) AS INTEGER) AS data,
  CASE strftime('%w', g.data)
    WHEN '0' THEN 'Duminică'
    WHEN '1' THEN 'Luni'
    WHEN '2' THEN 'Marți'
    WHEN '3' THEN 'Miercuri'
    WHEN '4' THEN 'Joi'
    WHEN '5' THEN 'Vineri'
    WHEN '6' THEN 'Sâmbătă'
  END AS ziua,
  g.angajat AS nume_si_prenume,
  substr(g.inceput, 12, 5) || ' - ' || substr(g.sfarsit, 12, 5)
    || CASE WHEN substr(g.sfarsit, 1, 10) <> g.data
       THEN ' (' || strftime('%d.%m.%Y', g.sfarsit) || ')'
       ELSE '' END AS observatii,
  g.data AS data_calendar
FROM garzi g
WHERE g.data >= :de_la AND g.data < :pana_la
ORDER BY g.inceput, g.angajat, g.marca, g.rowid;
```

Selectati intregul camp SQL din aplicatie (Ctrl+A) si inlocuiti cu o singura copie din [fisierul SQL curat](exemple/garzi-lunar.sql), apoi previzualizati si salvati modelul. Lipirea repetata a doua interogari produce eroare SQL invalid.

### Alinierea antetului si a subsolului PDF flexibil

In Tabel automat si Markdown, cele trei texte din antet apar simultan: stanga, centrat pe pagina si dreapta, fiecare in treimea lui. Textul se imparte automat pe linii in zona aleasa; tabelul incepe dupa cea mai inalta zona. Imaginile isi pastreaza proportiile si se aliniaza in zona lor. Subsolul respecta aceleasi alinieri. Antetul si imaginile se repeta si pe anexa finala, fara suprapunerea calculului/istoricului cu antetul. Pozitionarea libera continua sa foloseasca elementele si coordonatele explicite pentru textul din corp; nu se muta automat elementele deja configurate.

Pentru titlul central al listei de garzi, introduceti textul in Antet > Centru. In Markdown folositi doar {{tabel}} daca nu doriti un al doilea titlu in corp. Etichetele {{an}}, {{luna}} si {{perioada}} se inlocuiesc in texte; pentru sursa flexibila luna este numerica, iar perioada are forma AAAA-LL.

## Formatare si editare libera a rapoartelor

In Setari > Definire rapoarte > Editeaza, deschideti Formatare texte si etichete. Alegeti campul: titlu, introducere/paragrafe Markdown, cele sase zone antet/subsol, corpul tabelului sau capetele de tabel. Puteti alege bold, italic, marime 6-40 pt si aliniere. Marimea0 si alegerea Implicit mostenesc stilul existent. Reseteaza stilul elimina particularizarea campului. Stilurile se pastreaza in model, impreuna cu revizia lui. Pentru o eticheta proprie adaugati codul exact al coloanei din rezultatul SQL (sau codul coloanei raportului standard), apoi denumirile RO/EN. Coloana si eticheta ei pot fi formatate separat. In Pozitionare libera, elementele au si bold/italic/aliniere, iar marimea ramane in campul Font al elementului. Formatarea este pe camp/celula, nu un editor Word pentru fragmente arbitrare de caractere.

Pentru un raport deja generat: Rapoarte > Genereaza si previzualizeaza > Editeaza textele si celulele raportului. Modificati titlul, antetul/subsolul, introducerea, etichetele sau celulele. Focalizati un camp pentru a-i alege stilul. Tabelul de editare are25 randuri pe pagina; modificarile se pastreaza cand navigati intre paginile editorului. Celulele accepta text pe mai multe linii.

Aplica si actualizeaza PDF valideaza toate modificarile impreuna si creeaza o copie noua pentru exportul PDF/Excel. Pana la aplicare, PDF-ul afisat este cel anterior si exportul este dezactivat. Renunta la modificari restabileste ultima copie aplicata. Copia este marcata drept personalizata; detaliile autor/moment/identificator sunt in istoricul de la final si in jurnalul aplicatiei, cu valorile/stilurile inainte si dupa. Istoricul si calculele nu pot fi sterse din editorul de celule. Maximum1000 schimbari si2MiB per aplicare,4000 caractere per camp. O copie care nu incape pe pagina este refuzata; originalul ramane disponibil.

Modificarea unei celule din raport nu actualizeaza pontajul sau planificarea si nu recalculeaza alte celule/totaluri. Pentru corectarea datelor operationale folositi meniul lor si regenerati raportul. Editarile copiei apartin acestei livrari, nu modelului; salvati PDF/Excel inainte de inchidere. Regenerarea din model produce din nou valorile sursa. Stilurile si etichetele salvate in model se reutilizeaza la generari viitoare.

## Bold pentru weekend/sarbatori si grupuri la limita lunii

In definitia raportului flexibil, zona Zile speciale si intervalul raportului:

1. Folositi SQL-ul actualizat [garzi-lunar.sql](exemple/garzi-lunar.sql), care include data completa in data_calendar.
2. Introduceti data_calendar la Coloane ascunse pentru a pastra doar Data, Ziua, Nume si prenume, Observatii pe document.
3. La Coloana cu data completa alegeti/scrieti data_calendar.
4. Bifati Bold pentru sambata si duminica si/sau Bold pentru sarbatori configurate. Se evidentiaza randul pe baza datei complete, inclusiv cand coloana este ascunsa; nu se deduce data din textul zilei sau din numarul1-31. Stilurile explicite ale unei coloane/celule pot suprascrie stilul mostenit.
5. Bifati Include grupurile de zile libere la ambele margini ale lunii pentru grupurile consecutive care ating prima/ultima zi a lunii. SQL-ul trebuie sa foloseasca :de_la si :pana_la pentru interval. Doar filtrele :an/:luna nu includ automat datele din alta luna.
6. Refaceti previzualizarea si salvati modelul.

Calendarul de sarbatori configurat este sursa regulii. Pentru noiembrie2026, intervalul complet este31octombrie-1decembrie (inceputul lunii este duminica). Pentru decembrie2026 este28noiembrie-31decembrie: grupul28-30noiembrie/1decembrie se pastreaza integral. Nu se sare peste o zi lucratoare pentru a adauga alte sarbatori indepartate. Limita tehnica este366 zile suplimentare, cu eroare explicita pentru un calendar excesiv. Textul {{de_la}} arata prima data inclusa si {{pana_la}} ultima data inclusa; parametrul SQL :pana_la este prima data exclusa. Titlul poate continua sa indice luna aleasa. Totalurile contractuale si cele din catalogul lunar raman ale lunii selectate; lista garzilor poate acoperi intervalul extins.

[garzi-lunar.sql](exemple/garzi-lunar.sql) arata zilele cu garzi inregistrate. Pentru un tabel care arata toate zilele intervalului, inclusiv randuri goale unde nu exista garda, folositi varianta [garzi-lunar-calendar.sql](exemple/garzi-lunar-calendar.sql). Ambele folosesc data_calendar pentru bold si aceiasi parametri de interval.

## Editor vizual pdfme (0.92.33, D087)

Aspectul `PDFME` pastreaza SQL-ul, calculele si filtrele de mai sus; se schimba doar desenarea paginii. Sablonul este JSON pdfme (coordonate in mm), salvat in `continut.flexibil.sablon`, cu revizia modelului.

- Tabelul `tabel` (tip `table`, nu doar citire) primeste randurile vizibile, ca JSON. `head` are cate un titlu pentru fiecare coloana vizibila; la nepotrivire generarea se opreste cu mesajul `rapoarte.pdfme_coloane`. Se accepta un singur astfel de tabel.
- Elementele doar citire de tip text folosesc variabile `{nume}`: variabilele raportului, zonele antet/subsol (`antet_stanga` ... `subsol_dreapta`, cu etichetele `{{an}}` deja inlocuite), `introducere`, `nota`, `generat_la` si coloanele SQL. Elementele din `basePdf.staticSchema` se repeta pe fiecare pagina; `{currentPage}`/`{totalPages}` dau paginarea.
- Tipuri acceptate: text, text cu variabile, tabel, imagine, linie, dreptunghi, elipsa. Pagina de baza este goala (latime/inaltime/margini), nu un PDF importat. Sablon maximum 4 MiB, 50 pagini, 300 elemente pe pagina.
- Fonturi: `Go`, `GoBold`, `GoItalic`, `GoBoldItalic`, din aplicatie; nu se descarca nimic.
- Fluxul: Go pregateste documentul si intrarile, interfata genereaza PDF-ul cu `@pdfme/generator`, `FinalizeazaRaportPdfme` il valideaza cu pdfcpu (maximum 64 MiB), adauga anexa si il pune in cache pentru salvare/tiparire. Copia editata a raportului (D085) se regenereaza la fel; stilurile pe celula si boldul pe zile se aplica din 0.92.37 (F09-T22, sectiunea de mai jos).

### Fonturi Windows in sablon (0.92.34, D088)

`fontName` poate fi, pe langa fonturile Go, numele complet al unui font Windows (tabelul `name`, ID 4), de exemplu `Arial Bold`. Aplicatia citeste antetele din `%WINDIR%\Fonts` si `%LOCALAPPDATA%\Microsoft\Windows\Fonts` o singura data pe sesiune si exclude fonturile cu licenta restrictiva sau doar bitmap. Datele se cer dupa nume (`FontSistem`), niciodata dupa cale; colectiile `.ttc` sunt extrase ca font separat. Fonturile cu fsType „fara subset” se includ intregi. La generare se incarca toate fonturile din `fontName` (elemente, pagina fixa, stilurile tabelului); cele lipsa sunt inlocuite cu `Go` in copia sablonului folosita la generare, iar sablonul salvat ramane neschimbat.

### Prefixul angajatului (0.92.35, D089)

`angajati.prefix` contine titulatura in limba raportului (gol daca lipseste; marca de garda o preia pe a marcii principale). `angajati.angajat` si numele din celelalte seturi contin numele cu prefix; `nume` si `prenume` raman fara prefix. Exemplu: `SELECT prefix, nume, prenume, marca FROM angajati WHERE prefix='Dr.' ORDER BY nume`.

### Spatiul de lucru si actualizarea tabelului (0.92.36, F09-T25)

Tabelul `tabel` din sablon pastreaza in proprietatea `pontixColoane` codurile coloanelor SQL din care a fost facut; Go ignora proprietatea. **Actualizeaza tabelul din SQL** compara aceste coduri cu coloanele vizibile curente: pastreaza titlul si latimea coloanelor ramase, adauga coloanele noi cu latimea medie, scoate coloanele disparute si normalizeaza latimile la 100%. Sabloanele mai vechi, fara coduri, sunt potrivite dupa titlul coloanei. Alegerea aspectului vizual pe un model fara sablon ruleaza SQL-ul in modul tabel si adauga `angajat_id` la coloanele ascunse daca exista si alte coloane.

### Stiluri pe rand si pe celula in tabelul pdfme (0.92.37, F09-T22)

Tabelul pdfme nu are stiluri pe rand, asa ca la generare aplicatia inlocuieste doar desenarea tabelului `tabel`; editorul, panoul de proprietati si impartirea pe pagini raman cele din pdfme. Go trimite intrarea `pontix_stiluri` (JSON): `randuri` (indice in tot corpul -> stil; aici ajung randurile aldine pentru weekend/sarbatori), `celule` (`"rand:coloana"` -> stil, din editarea copiei D085, coloana fiind indicele coloanei vizibile), `corp` si `cap` (stilurile `tabel` si `antet_tabel`). Un stil are `aldin`, `cursiv`, `font_pt`, `aliniere` (L/C/R). Ordinea de aplicare: sectiune, rand, celula. Aldinul foloseste `GoBold`/`GoItalic`/`GoBoldItalic` sau varianta „Nume Bold/Italic” a fontului Windows, daca exista. Inaltimile randurilor sunt cele calculate de pdfme pentru textul normal, ca paginarea sa nu se schimbe; un text aldin foarte lung poate depasi usor celula. Stilurile pe coloana (`coloana:cod`) nu se transmit: in modul vizual coloanele se formateaza in editor. Celulele cu imagini se deseneaza cu codul original, fara stiluri.

## Rapoartele standard ca tabele calculate (0.92.38, D090)

| Tabel | Sursa | Coloane |
| --- | --- | --- |
| pontaj_lunar | Pontaj lunar | rand, grup, angajat, marca, sectie, departament, functie, contract, norma_contract, z01-z31 |
| pontaj_detaliat | Pontaj detaliat | ca mai sus, plus total_L_INCEPUT si totalurile configurate (total_T, total_ZNL, total_CO ...) |
| fisa_individuala | Fisa individuala (cere angajat) | rand, grup, data, zi, ore, descriere, categorie |
| absente_anuale | Evidenta anuala a absentelor | rand, grup, data, angajat, marca, sectie, departament, contract, norma_contract, cod, descriere, categorie |
| grafic_garzi | Grafic garzi | rand, grup, angajat, marca, ..., z01-z31, total_garzi, total_ore |

Un tabel se calculeaza doar daca SQL-ul (sau un calcul SQL) il numeste, cu filtrele raportului. Valorile sunt textele afisate de raportul standard. `rand` este pozitia randului in raportul calculat; cu ea, aplicatia preia titlurile, culorile pe celula si pe coloana, randurile de grup/aldine si legenda doar pentru coloanele alese. Daca SQL-ul combina doua tabele standard sau nu pastreaza `rand`, raportul ramane corect, dar fara aceste metadate.

### Modelele predefinite si stilurile din raport (0.92.39, D090)

Modelele `raport-pontaj-sql`, `raport-pontaj-detaliat-sql`, `raport-fisa-sql`, `raport-absente-sql`, `raport-garzi-sql` (migrarea 0033) au SQL-ul `SELECT rand, <coloanele modelului vechi> FROM <tabel> ORDER BY rand`, cu `rand` ascuns, si sablon pdfme generat de `frontend/src/modeleSql.generator.test.ts` (`GENEREAZA_MODELE=<cale> npx vitest run src/modeleSql.generator.test.ts`). Tabelul sablonului are `pontixTitluri: true`: capul afiseaza titlurile din raport (`pontix_stiluri.titluri`), in limba lui. `pontix_stiluri` mai poarta `cap_coloane` (culoarea coloanei pe cap), fundalul celulelor din `Culori` (WEEKEND `#E7EAE7`, SARBATOARE `#FFF0C7`, ABSENTA `#DFEFE4`), randurile de grup (aldin, `#E7EFEA`) si `ascunse`, indicii coloanelor `z29`-`z31` inexistente in luna; desenarea le scoate si pastreaza proportiile celorlalte. Tabelele standard au mereu `z01`-`z31`, cele in plus fiind goale.

## Conversia modelelor vechi si stergerea motorului vechi (0.92.42, F09-T23)

Aspectele Tabel, Markdown si Pozitionare libera descrise mai sus nu mai sunt desenate de aplicatie: `ConvertesteModeleVechi` (internal/rapoarte/conversie.go) le trece pe pdfme la primul acces la rapoarte (`ListaModeleRapoarte`, `PregatesteRaport`), prin `Salveaza` (revizie noua, definitia veche in jurnal). Sabloanele se construiesc in Go (`SablonDinColoane`, pe proprietatile implicite pdfme 6.2.8 din `pdfme_baza.json`, verificat identic cu generatorul TypeScript). Modelele standard: SQL `SELECT rand, <coloanele modelului> FROM <tabel> ORDER BY rand`, `rand` ascuns. Tabel: sablon fara coloane, marcat `pontixTitluri`, care se adapteaza la generare (`SablonPentruDocument`). Markdown: titluri `#`, paragrafe, liste `-` (cu `•`), aldin `**...**` prin `textFormat: inline-markdown`, `{{tabel}}` devine tabelul; tabelele scrise cu `|` raman text. Pozitionare libera: elementele devin texte la `x+margine`, `y+margine+10` mm, cu aldin/cursiv/aliniere; `{{tabel}}` devine tabel; elementele repetate pe rand devin tabelul de date. `{{variabila}}` devine `{variabila}`. Imaginile din antet devin imagini pe prima pagina, cele din subsol imagini fixe pe fiecare pagina. Raportul pe un tabel standard foloseste variabilele raportului standard (`perioada` localizata, `luna` ca nume, persoana aleasa).
