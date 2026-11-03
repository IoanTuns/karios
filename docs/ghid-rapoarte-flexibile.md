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
