# D082 - Rapoarte SQL, formule si aspect configurabil

Stare: propusa. Data: 2026-11-02. Taskuri F09-T13/F09-T14. Responsabil Codex.

## Cerinta si confirmari

Utilizatorul autorizeaza implementarea intregului TODO: campuri SQL, coloane calculate definite din aplicatie, coloane pe tip nou de concediu, selectie concedii, garzi eligibile si sabloane Markdown ori pozitionare libera. Confirma: ore lucratoare = ore contractuale minus concediile care reduc obligatia; deplasarile creditate nu se scad. Garzile eligibile sunt intersectiile efective cu zilele 00-24, fara dublare weekend/sarbatoare. Concediile pot fi selectate pe tip si categorie, deduplicate pe marca/data.

## Propunere tehnica implementata

Definitia existenta JSON primeste continut.flexibil optional; fara migrare si fara rescrierea modelelor existente. Trei surse noi: SQL_LIBER, CALCULATE_LUNAR si LISTA_GARZI. Reviziile, duplicarea si jurnalul definitiilor raman comune. Coloanele custom apartin modelului si se pot reutiliza duplicandu-l.

SQL se executa numai pe un instantaneu SQLite in memorie, populat din campuri de raportare alese explicit, in aceeasi tranzactie de citire a bazei operationale. Nu sunt copiate conturi, parole, chei, documente BLOB sau jurnalul. Filtrele se aplica inainte de interogarea utilizatorului. Authorizer-ul permite SELECT/READ si functii de citire enumerate, nu scrieri, PRAGMA, ATTACH, extensii sau functii de fisiere. O singura instructiune, parametri legati separat, maximum 80 coloane, 10000 randuri, 8 MiB rezultat, 256 KiB celula, 2 secunde per interogare si 20 secunde pregatirea datelor. Aceste limite nu sunt un buget global de heap al aplicatiei.

Formulele numerice folosesc un arbore de expresii validat din biblioteca standard Go, nu executie Go/JavaScript sau evaluare generala. Operatii + - * / si ABS/MIN/MAX/ROUND, evaluare in ordinea coloanelor; fara cicluri ori campuri inexistente. Datele lipsa se propaga, cu optiune explicita zero pentru rezultat. Unitati si precizie configurabile; formulele utilizatorului nu schimba motorul normei contractuale existent. SQL scalar cere un singur rezultat numeric/NULL. Pentru calcul per angajat, angajat_id trebuie inclus in rezultat sau ales ca filtru; poate fi ascuns din tabel.

Totalurile reutilizeaza calculul comun al pontajului. Garzile active efective au prioritate fata de program pentru aceeasi marca/data; intrarile retrase nu participa, iar planul activ completeaza lipsa efectivelor. Planificarea foloseste durata nominala; efectivele folosesc intervalele in fusul Europe/Bucharest si durata reala, inclusiv DST. Fragmentele se limiteaza la luna si calendar, fara dublare. garzi_zile expune fragmente pe zi pentru reguli SQL proprii. Marcile asociate raman distincte.

Markdown suporta titluri, paragrafe, liste, evidentiere **text**, tabele simple si {{tabel}}; este un subset documentat, nu un motor CommonMark complet sau HTML executabil. Variabilele pot fi campuri rezultate si campuri globale. Editorul vizual foloseste mm, coordonate numerice si tragere; elementele pot fi repetate pe fiecare rand de date, cu pas configurabil si paginare. Tabelul integral si campurile repetate sunt variante alternative, fara dublarea implicita a datelor. Coordonatele/textul care nu incap sunt refuzate. Istoricul si formulele raman in anexa finala.

PDF-ul flexibil are versiunea 1 separata. Rapoartele vechi isi pastreaza motorul si versiunile; catalogul surselor are referinta tehnica noua modele-v6, fara modificarea modele-v5 sau a cazurilor de calcul confirmate.

## Biblioteci si licente

Nu se adauga un pachet sau program nou. gofpdf v1.4.3, deja inclus de maroto si in go.mod/go.sum, este utilizat direct pentru coordonate si Unicode; licenta MIT verificata in LICENSE al modulului local si notificarea existenta docs/licente/rapoarte/gofpdf-v1.4.3.txt. Clasificarea devine directa in go.mod, versiunea ramane fixata. SQLite ncruces v0.36.0 si fonturile Go Regular/Bold existente sunt reutilizate. Nu se instaleaza un browser ori un motor HTML/PDF extern.

## Alternative

SQL pe conexiunea operationala respins: nu este necesar accesul la secrete si scrieri. Expresii arbitrare Go/JS respinse: calculul trebuie limitat si testabil. HTML/CSS executabil nu este introdus; alternativa de pozitionare ceruta este editorul vizual si PDF nativ. O copie de PDF vechi cu eticheta noua este exclusa; buildul regenereaza manualele.
