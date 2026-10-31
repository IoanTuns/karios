# F07-T27 - Ore lucratoare, sold lunar si initializare

Responsabil: Codex. Stare: in lucru. Ramura F07-T27-solduri-initializare.

## Cerinta curenta

- Redenumire Ore lucrate in Ore lucratoare, cu calcul din zilele lucratoare ale lunii si norma zilnica.
- L Fin: diferenta dintre ore lucratoare si ore efectiv lucrate; sensul exact si participarea L Inceput sunt in curs de confirmare deoarece difera de cerinta anterioara.
- L Inceput imediat dupa Marca, preluat din L Fin al lunii precedente.
- Initializare pe angajat a soldului de ore si a concediului de odihna deja efectuat anterior folosirii aplicatiei.

## Clarificari trimise

1. Norma bruta lunara fara scaderea CO/CEV/CFS/CM, sau regula anterioara care scade concediile?
2. L Fin = obligatie minus efectiv, ori sold initial plus aceasta diferenta, ori sensul anterior efectiv minus obligatie cu sold initial?

Nu se schimba formulele sau testele de referinta pana la raspuns. Marcile de garda raman fara norma prestabilita, conform regulii confirmate anterior.

## Integrare identificata

Calcul comun in internal/datebaza/pontaj_calcul.go si internal/calcul/pontaj_ore.go, raport detaliat in internal/rapoarte/document.go. Norma datorata cu concedii este deja implementata ca functie pura separata, dar nu este integrata complet. L_FIN este actualmente o rubrica configurabila, nu un sold lunar complet implementat. Formularul angajatului contine dreptul anual de CO, fara consum initial.

## Propunere pentru initializare

O data de referinta la inceput de luna si un sold de ore semnat per marca. Soldurile ulterioare se calculeaza cronologic, inclusiv trecerea decembrie/ianuarie, fara a dubla soldul initial. CO deja efectuat se refera la anul datei de referinta si numai la perioada anterioara acesteia; nu reprezinta reducerea dreptului contractual. CO introdus ulterior se adauga o singura data la consumul anual. Se pastreaza revizie si audit la editare. Este necesara avertizarea impotriva introducerii acelorasi zile ca sold initial si ca absente calendaristice.

Nicio migrare sau schimbare de cod aplicata inca. Urmeaza confirmarea regulilor, apoi migrare noua, interfata de initializare, calcul comun, raport/legenda si teste lunare/trecere de an/concurenta. Versiunile aplicatiei, raportului si manualului se vor actualiza impreuna cu implementarea.

Confirmare explicita primita: se scad CO/CEV/CFS/CM aprobate; L Fin = L Inceput + ore lucrate - ore lucratoare. Inceputa implementarea cu initializare pe marca si luna de referinta, CO initial aferent anului respectiv.

## Rezultat 0.92.18

Cerinta curenta implementata: calcul comun al orelor lucratoare cu scaderea CO/CEV/CFS/CM aprobate, sold lunar in minute, L Inceput dupa Marca in raport detaliat v10, initializare in formularul angajatului si consum CO initial in sumarul anual. Dreptul anual nu este schimbat. Migrare 0028, revizie si audit; editarea initializarii este refuzata daca afecteaza luni inchise. CO initial nu se dubleaza cu absente anterioare datate. Valorile zero corecteaza/anuleaza soldul initial fara stergere fizica.

CEV/CFS folosesc codurile explicite din regulile de pontaj; Special singur nu scade norma. Norma curenta este folosita si pentru lunile recalculate; nu exista un istoric contractual de norma, iar aplicatia nu presupune unul. Fara initializare, prima luna cu date porneste de la zero. Marcile de garda nu primesc norma sau sold calculat. Marcajul T=ture din discutia anterioara ramane in afara acestei schimbari.

Verificari trecute: gofmt, go vet ./..., go test -p 1 -timeout 120s ./..., 22 teste frontend, TypeScript/Vite, Edge focus/zoom si noile campuri RO/EN, Edge absente din Pontaj RO/EN, 3 teste manual, 7 livrare si control documentatie. Test suplimentar de raport cu L Inceput=1.5 dupa Marca, cu vet pe pachet, trecut. Referintele tehnice noi v10 sunt separate; referintele vechi si cazurile contabile protejate nu au fost modificate. Contoarele de migrari din testele tehnice au crescut cu unu pentru migrarea noua.

Testele noi acopera norma partiala, garda fara norma, concedii neaprobate, report intre ani, semnul soldurilor, minute, revizii concurente, CO initial si dubla numarare, blocarea lunilor inchise. La primele rulari au fost identificate numai actualizarile necesare ale contoarelor tehnice si referintelor raportului schimbat. Dependente verificate: Go stabil 1.27.3, npm fara actualizari; 62 actualizari tranzitive Go sunt urmarite separat in F14-T21. Fara instalari sau biblioteci noi.

Manual editie 1.29, raport detaliat v10. Pas urmator: build oficial, integritate pachet si proba operator WebView2. Nicio baza reala accesata, fara push/merge.
