# Analiza UI · Rapoarte cu editor vizual pdfme (0.92.34)

Data: 2026-11-05 · Responsabil: Claude · Task: F16-T33 · Baza: commit 0e46691 (0.92.34)

## Metoda si limite

Am parcurs fluxul complet in Edge headless, cu interfata construita si puntea Wails simulata (`node verificari/interfata.mjs --rapoarte-pdfme`, scenariu nou `verificari/rapoarte-pdfme-scenarii.mjs`). Pasii: lista modelelor, model nou SQL, verificare, aspect „Editor vizual (pdfme)”, adaugarea unui font Windows, generarea PDF-ului. Fonturile sunt reale (Go din modulul x/image, Arial din Windows); datele sunt fictive. Capturi la 1024×640 si 1440×960 cu dialogul marit, plus masuratori de pozitie si derulare. Rezultat: PDF generat (34 KB), fara erori JavaScript, fara derulare orizontala a paginii.

Limite:

- Puntea este simulata; backendul real e acoperit de testele Go.
- Edge headless nu afiseaza PDF intr-un iframe, deci previzualizarea PDF nu a fost vazuta, doar generarea ei.
- Selectarea unui element pe pagina nu a putut fi simulata, deci panoul de proprietati al elementului nu a fost capturat.
- Nu s-a facut proba in WebView2, cu operatori, cu tastatura sau cu cititor de ecran.

Constatarile marcate „din cod” nu au fost reproduse vizual.

## Concluzie

Editorul vizual rezolva problema de baza: pagina se construieste prin tragere, iar PDF-ul arata ca in editor. Problema principala nu este pdfme, ci locul in care l-am pus: la sfarsitul unui formular lung, intr-un dialog cu derulare. In modul vizual, formularul mai arata si optiuni care nu se aplica. Inainte de rafinari, editorul are nevoie de un spatiu de lucru propriu.

## Avantaje

- Un sablon complet (antet, titlu, tabel, numarul paginii) se obtine dintr-un clic, din coloanele SQL.
- Tabelul continua singur pe pagini, cu capul repetat. Diacriticele ies corect cu fonturile Go si Windows.
- La 1440×960, cu dialogul marit, zona de lucru are 1311×720 px si pagina A4 la 100% se citeste bine (captura 10).
- Butonul „Verifica datele si previzualizeaza PDF” sta imediat sub editor.
- Nepotrivirea dintre coloanele tabelului si SQL are un mesaj clar, cu cifre si actiunea de urmat.
- Interfata nu are derulare orizontala si nu a dat erori JavaScript la nicio rezolutie.

## Probleme si propuneri

P1 = blocare sau rezultat neinteles; P2 = efort si descoperire; P3 = confort.

| ID | Prioritate | Constatare | Propunere | Criteriu de acceptare |
| --- | --- | --- | --- | --- |
| UIR-01 | P1 | Editorul incepe la circa 2.950 px de la inceputul dialogului, dupa 34 de campuri si butoane (antet, continut, SQL, concedii, calcule). La 1024×640 zona editorului (520 px) e mai inalta decat corpul vizibil al dialogului (circa 440 px), deci exista doua derulari una in alta. Campul de fonturi si selectorul de aspect nu se vad impreuna cu pagina (capturi 04, 08). | In modul vizual, editorul se deschide intr-un spatiu de lucru separat, pe tot ecranul: file „Date (SQL)” / „Aspect” / „Pagina”, sau un buton „Deschide editorul vizual”. Dialogul se mareste automat, iar editorul ocupa inaltimea disponibila. | La 1024×640 pagina, bara de elemente si panoul de proprietati sunt vizibile fara derularea dialogului; o singura bara de derulare. |
| UIR-02 | P1 | In modul vizual raman vizibile sectiunile Antet si Subsol, Formatare texte si etichete, bold pentru weekend si sarbatori. Antetul si subsolul se folosesc doar ca variabile `{antet_stanga}`; stilurile si boldul nu au efect in pdfme (F09-T22). Utilizatorul le poate completa fara rezultat. | Ascunderea sau marcarea explicita „nu se aplica in editorul vizual”. Antetul si subsolul se editeaza direct pe pagina; textele din formular raman doar pentru modelele vechi. | In modul vizual nu exista niciun control vizibil care sa nu schimbe PDF-ul, sau fiecare astfel de control are mentiunea respectiva. |
| UIR-03 | P1 | Din cod (RaportFlexibil.tsx, schimbarea aspectului): daca aspectul vizual se alege inainte de „Verifica raportul”, sablonul se face fara coloane, iar prima verificare se opreste cu nepotrivire „0 fata de N”. | La alegerea aspectului vizual, aplicatia ruleaza SQL-ul ca sa afle coloanele (ca butonul „Sablon nou din coloanele SQL”) sau cere intai verificarea. | Model nou → SQL → aspect vizual → Verifica: primul PDF iese fara eroare. |
| UIR-04 | P1 | Capul tabelului arata codurile SQL (`angajat_id`, `total_ore_lucrate`); in aplicatia reala Go inlocuieste doar `_` cu spatiu. Identificatorul intern apare ca prima coloana. | Capul tabelului se ia din etichetele coloanelor (Formatare, denumirile coloanelor calculate), cu fallback lizibil. `angajat_id` se ascunde implicit. | Sablonul nou din `SELECT angajat_id, angajat, …` are capete in limba utilizatorului si fara identificatorul intern. |
| UIR-05 | P1 | Daca SQL-ul se schimba dupa crearea sablonului, tabelul nu se actualizeaza; problema apare abia la verificare. Singura iesire este un sablon nou, care pierde tot aspectul. | Detectarea coloanelor adaugate sau eliminate si actiunea „Actualizeaza tabelul”, care pastreaza latimile, stilurile si restul paginii. | Adaugarea unei coloane in SQL nu cere refacerea sablonului; elementele existente raman neschimbate. |
| UIR-06 | P2 | Interfata editorului e in engleza („Field List”, „Bulk update field names”, panoul de proprietati), iar elementele au nume tehnice (`antet_stanga_txt`, `titlu_txt`). | Etichetele pdfme se traduc din cataloagele i18n (optiunea `labels`). Elementele din sablonul de pornire primesc nume lizibile („Antet stanga”, „Titlu”, „Tabel date”). | Editorul nu afiseaza texte in engleza in RO; lista elementelor e inteleasa fara documentatie. |
| UIR-07 | P2 | Sintaxa variabilelor difera: formularul arata `{{titlu}}`, pdfme cere `{titlu}`. Coloanele SQL disponibile nu sunt listate langa editor. | In modul vizual, o lista de campuri cu sintaxa `{…}` (variabilele raportului si coloanele SQL), cu copiere sau inserare prin clic. | Utilizatorul poate pune `{marca}` intr-un text fara sa consulte manualul. |
| UIR-08 | P2 | Fontul Windows se adauga dintr-un camp de deasupra editorului, departe de panoul de proprietati unde se alege fontul; confirmarea e doar o notificare temporara (captura 05). | Alegerea direct langa editor, vizibila impreuna cu panoul; lista fonturilor deja adaugate sau folosite in sablon; pentru bold sau italic, varianta sugerata automat. | Dupa adaugare, fontul se poate alege pentru elementul selectat fara derulare. |
| UIR-09 | P2 | La 1024×640 pagina A4 are circa 440 px latime; textul de 9 pt e greu de citit, iar „Field List” ocupa aproape jumatate din latime (capturi 04, 08). | Sub 1280 px panoul lateral se restrange implicit, iar pagina se deschide cu „potriveste latimea”. | La 1024 px pagina foloseste cel putin 70% din latimea editorului. |
| UIR-10 | P2 | Previzualizarea PDF apare intr-un iframe sub editor si sub tabelul de date, deci inca o derulare (captura 07). Afisarea in WebView2 nu este verificata. | Previzualizare intr-un panou alaturat sau intr-un dialog, cu numarul de pagini si comutare rapida editor/PDF; proba in aplicatia construita. | Dupa „Verifica”, PDF-ul se vede fara derulare; confirmat in WebView2. |
| UIR-11 | P3 | Bara de elemente are doar pictograme (text cu variabile, text, tabel, imagine, linie, dreptunghi, elipsa); diferenta dintre cele doua tipuri de text nu e evidenta. | Etichete sau indicii in romana; eventual doar tipurile utile in rapoarte. | Fiecare pictograma are un nume accesibil in RO. |
| UIR-12 | P3 | Lucrul cu mouse-ul domina; tastatura si cititorul de ecran nu au fost verificate pentru editor. Selectarea elementelor nu a putut fi automatizata. | Audit manual al tastaturii (selectie, mutare cu sagetile, panou de proprietati) si al numelor accesibile. | Un element se poate selecta, muta si redimensiona fara mouse, sau limita este documentata. |

## Ordine recomandata

1. UIR-01 si UIR-02 impreuna: spatiul de lucru separat pentru aspect rezolva derularea si ascunde controalele fara efect.
2. UIR-03, UIR-04, UIR-05: legatura dintre SQL si tabel (coloane, etichete, actualizare).
3. UIR-06, UIR-07, UIR-08: limba, variabile, fonturi.
4. UIR-09 si UIR-10, apoi P3, dupa proba in WebView2.

Toate sunt propuneri; nu schimba calculele sau regulile de plata.

## Dovezi

Capturile si `masuratori-pdfme.json` sunt in `build/bin/verificare-ui-2026-11-05T17-57-45-245Z/` (local, ignorat de git): 01 lista modelelor, 02 model nou, 03 inceputul dialogului, 04 editorul, 05 font adaugat, 07 previzualizare, 08–09 la 1024×640, 10 la 1440×960 marit.
