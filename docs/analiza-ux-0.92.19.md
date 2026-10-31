# Analiza UX Pontix 0.92.19

Data: 2026-10-31. Responsabil: Codex. Task: F16-T30.
Baza: cod c98109b, livrare documentata ee706dc, inspectia componentelor curente si rezultatele F16-T29. Analiza euristica; nu este un studiu cu operatori si nu estimeaza timpi sau rate de eroare masurate. Nu s-au rulat probe UI noi in acest task. Nu s-au accesat date reale. Recomandarile sunt propuneri, fara implementare sau schimbarea regulilor de calcul. Blocajul Excel ramane amanat explicit de utilizator; nu este reinvestigat aici.

## Concluzie

Aplicatia ofera instrumentele necesare unui operator experimentat, dar solicita inca prea multe cunostinte despre modelul intern. Prioritatea UX recomandata este trecerea de la meniuri si configurari la sarcini concrete: pregatesc luna, completez, verific, public, inchid. Problemele de derulare si contrast remediate in F16-T29 nu sunt raportate din nou drept defecte deschise.

## Avantaje

- Planificarea afiseaza separat modificarile nesalvate, datele salvate si publicarea in Pontaj; previzualizarea arata blocajele si inregistrarile efective afectate.
- Grilele pastreaza reperele de identificare la derulare. Zoomul general si densitatea grilei permit adaptarea spatiului de lucru.
- Filtrele profesionale, cautarea si selectia pe interval sustin operatiuni repetate pentru mai multe persoane.
- Formularul angajatului este grupat; soldurile initiale sunt separate de modificarile curente. Numele afisat diferentiaza Garda fara schimbarea numelui salvat.
- Absentele explica prioritatea si marcile afectate; protejarea lunilor inchise si istoricul reduc riscul modificarilor necontrolate.
- Ajutorul contextual, verificarea configurarii si AI cu previzualizare ofera sprijin fara a cere memorarea tuturor regulilor.
- Actiunile administrative rare sunt separate; stergerea este limitata si verificata in backend.

## Probleme si propuneri prioritizate

P1 = risc de greseala, blocare sau rezultat neinteles. P2 = eficienta si descoperire. P3 = confort si personalizare. Severitatea este evaluare euristica, nu incident masurat.

| ID | Prioritate | Constatare / limita | Imbunatatire propusa | Criteriu de acceptare |
| --- | --- | --- | --- | --- |
| UX30-01 | P1 | Verificarea ghidata numara lipsuri, dar toate legaturile de verificare deschid generic Angajati. Utilizatorul trebuie sa regaseasca persoanele. Sursa: ConfigurareGhidata.tsx. | Deschidere cu filtrul exact, lista persoanelor vizate si revenire la lista de verificari; separare informatii de probleme blocante. | Un clic pe lipsa categoriei afiseaza numai persoanele respective; Garda independenta nu este prezentata ca eroare. |
| UX30-02 | P1 | Relatia contract - categorie - program - tip - regim ramane dificil de inteles, chiar cu explicatiile de excludere. | Previzualizare a eligibilitatii direct in formularul angajatului si al categoriei: programe permise, excluse si motiv; legaturi spre configurarea exacta. | Operatorul poate afla de ce TN9/G12 nu este disponibil fara sa incerce pe rand meniurile. Nu se schimba eligibilitatea automat. |
| UX30-03 | P1 | In editorul regulilor din Planificare, unele raspunsuri backend ajung direct la setEroare(r.eroare.cod), fara traducere. Constatare statica, nu eroare reprodusa vizual. | Traducere prin cod si parametri, mesaj orientat spre actiune; codul tehnic separat pentru diagnostic. | Nicio eroare a editorului nu afiseaza exclusiv o cheie interna in RO/EN. |
| UX30-04 | P1 | La schimbarea filtrelor profesionale in Rapoarte se invalideaza rezultatul, dar nu si selectie.angajat. Un angajat selectat poate disparea din optiuni in timp ce ID-ul ramane in cerere. Sursa: Rapoarte.tsx; risc de verificat, comportamentul final backend nu este afirmat. | Revalidarea persoanei selectate la fiecare filtru; resetare explicita si mesaj daca nu mai corespunde. Curatarea mesajului de succes vechi. | Selecteaza persoana, aplica un filtru incompatibil: formularul nu poate trimite ID ascuns si explica resetarea. |
| UX30-05 | P1 | Publicarea si absentele au rezumate, dar operatorul trebuie sa reconstruiasca efectul final din contoare si reguli. | Previzualizare comparativa inainte/dupa pentru operatiile cu suprascriere, pe marca si data; distinctie intre valori retrase si cele create, cu sursa si motiv. | Un caz program -> CO -> CM arata exact rezultatul, toate marcile asociate si tura de noapte intreaga, fara a modifica regulile aprobate. |
| UX30-06 | P2 | Fluxul lunar este distribuit intre Planificare, Pontaj, Absente si Rapoarte. Starea publicarii este mai clara, dar nu exista aici dovada unei singure liste executabile care coordoneaza toate etapele. | Panou lunar cu legaturi: configurari lipsa, plan nepublicat, modificari nesalvate, verificari, raport si inchidere. Refolosirea verificarilor existente, fara calcule paralele. | Operatorul vede ce ramane de facut pentru luna si ajunge direct la inregistrarile relevante; avertismentele nu blocheaza fara regula existenta. |
| UX30-07 | P2 | Chipurile filtrelor profesionale sunt texte fara eliminare individuala si fara numele dimensiunii. Rezultate goale pot cere cautare in mai multe grupuri de filtre. Sursa: FiltreProfesionale.tsx. | Chipuri de forma Functie: Medic, cu eliminare individuala; explicatie si resetare contextualizata la rezultate goale. | Se poate elimina un singur filtru cu tastatura, pastrand luna si celelalte criterii. |
| UX30-08 | P2 | Completarea, editarea directa, clicul dreapta si AI ofera mai multe cai pentru acelasi rezultat. Operatorul nou poate confunda celula activa cu intervalul selectat. | Rezumat persistent al selectiei langa actiune: persoane, zile si marci; comenzi principale Program / Absenta / Golire, cu actiuni rare secundare. | Inainte de aplicare se intelege domeniul operatiei; filtrarea nu extinde selectia pe ascuns. |
| UX30-09 | P2 | Protectiile de abandon folosesc frecvent window.confirm: previn pierderea, dar ofera actiuni generice si nu permit salvare direct din dialog. | Dialog comun cu numarul modificarilor si actiunile Salveaza si continua / Renunta la modificari / Ramai. Pentru publicare se pastreaza etapa separata. | Inchiderea, schimbarea lunii si navigarea au acelasi comportament; esecul salvarii nu continua navigarea. |
| UX30-10 | P2 | S/M/L reprezinta densitatea, separat de zoomul general. Titlurile ajuta, dar literele singure nu explica diferenta unui utilizator nou. | Etichete Compact / Normal / Aerisit sau indiciu vizibil, cu procentul zoomului general separat si resetare. | Operatorul poate identifica independent ambele setari fara hover; accesibil pe tastatura. |
| UX30-11 | P2 | AccesTabele ofera acelasi aria-label generic tuturor regiunilor. Testele de contrast si geometrie nu echivaleaza cu verificarea unui cititor de ecran. | Denumiri specifice regiunilor; audit manual tastatura/cititor de ecran pentru grile, dialoguri si modificari dinamice. | Sunt anuntate numele tabelului si reperele celulei; nu exista capcane de focus. Probe Windows DPI 125/150% separate de zoomul aplicatiei. |
| UX30-12 | P3 | Preferintele si filtrele pot fi numeroase pentru aceleasi sarcini recurente; nu s-a confirmat lipsa tuturor mecanismelor de memorare existente. | Evaluarea si, unde lipseste, salvarea vederilor personale pentru fiecare ecran; revenire clara la vederea implicita. | O vedere salveaza numai criterii si prezentare, nu selectii periculoase sau date de pontaj; utilizatorii nu isi suprascriu preferintele. |

## Ordine recomandata

1. Lot A: UX30-03 si UX30-04, defecte concrete de tratare a starii/erorilor; apoi UX30-01 si UX30-02 pentru configurare.
2. Lot B: UX30-05, UX30-08 si UX30-09 pentru claritatea efectelor si recuperarea din greseli.
3. Lot C: UX30-06, UX30-07 si UX30-10 pentru eficienta lunara.
4. Lot D: UX30-11 si evaluarea UX30-12. Accesibilitatea ramane criteriu in toate loturile, nu se amana remedierea unui blocaj demonstrat.

## Validare recomandata cu operatori

Sesiuni pe date fictive, fara explicatii oferite inainte de sarcina: configurare angajat in ture; explicarea unui program lipsa; publicare peste pontaj existent; introducere CM peste CO cu marca Garda asociata; corectare/anulare absenta; initializare L Inceput si CO executat; raport filtrat si inchidere de luna. Se noteaza reusita fara ajutor, pasii gresiti, timpul si increderea in rezultat, fara praguri inventate prezentate ca masuratori actuale.

Prima runda trebuie sa includa un operator nou si unul experimentat, ecran mic si scalare Windows, teme RO/EN dupa utilizare. La export se verifica separat salvarea fisierului si deschiderea in programul extern; investigatia blocajului Excel se reia numai cand utilizatorul o solicita.

## Limite si predare

Analiza nu certifica formule contabile, conformitate legala sau accesibilitate completa. Nu se propun procente sau reguli de plata noi. Nu s-au rulat teste Go/frontend deoarece nu s-a schimbat codul. Dependentele nu se modifica intr-un task exclusiv de documentare. Pas urmator: alegerea lotului de implementare; fara push/merge, fara publicare externa.
