# Fișa de acceptanță - Pontix {{VERSIUNE}}

Commit {{COMMIT}}. Data: __________. Cine face proba: __________________.
PC 1 (curat): Windows ______, 64 biți, utilizator fără drepturi de administrator: da / nu.
PC 2 (pentru stick): Windows ______. Stick USB: ______ GB.

**Reguli:**
- Folosiți **numai date fictive** și foldere separate de datele reale.
- La fiecare probă marcați **T** (trecut), **E** (eșuat) sau **N** (neexecutat) și notați ce ați văzut.
- Nu bifați nimic pe baza testelor automate ale programatorului. Fișa înseamnă ce ați văzut voi.

## Ce pregătiți înainte

- Arhiva acestei versiuni (`pontix-{{VERSIUNE}}-….zip`).
- Arhiva unei versiuni vechi, de exemplu 0.92.92, cu un folder `karios` care are **date fictive**: câțiva angajați, o lună de pontaj, un concediu, o zi liberă adăugată (punte), un program cu pauză.
- Un PC pe care Pontix n-a rulat niciodată (PC 1) și un al doilea PC (PC 2).
- Un stick USB gol.

## A. Prima pornire pe un PC curat (PC 1)

| # | Proba (ce faceți) | Rezultat așteptat | T/E/N | Observații |
| --- | --- | --- | --- | --- |
| A1 | Extrageți toată arhiva într-un folder al vostru, de exemplu `Documente\Pontix`. Nu porniți din ZIP. | Folderul conține `pontix.exe`, manualele (MD, HTML, PDF), `ACTUALIZARE.md`, această fișă. | | |
| A2 | Porniți `pontix.exe` ca utilizator obișnuit, fără „Executare ca administrator”. | Aplicația pornește. Windows poate afișa un avertisment SmartScreen, pentru că executabilul nu e încă semnat: notați textul exact. | | |
| A3 | Dacă aplicația nu pornește și cere WebView2, notați mesajul. | Pe Windows 11, WebView2 există deja. Pe Windows 10 poate lipsi; notați dacă a trebuit instalat. | | |
| A4 | Creați primul administrator. Notați cheia de recuperare pe hârtie. | Contul se creează. Cheia de recuperare se afișează și trebuie confirmată. | | |
| A5 | Închideți aplicația, apoi o reporniți. | Se deschide direct sau cere autentificarea; datele se păstrează. | | |
| A6 | Verificați că lângă `pontix.exe` au apărut folderul `karios` (și, eventual, `ai`). | Datele stau doar acolo, nu în altă parte a calculatorului. | | |

## B. Rularea de pe un stick USB (PC 1 → PC 2)

| # | Proba | Rezultat așteptat | T/E/N | Observații |
| --- | --- | --- | --- | --- |
| B1 | Închideți aplicația. Copiați **tot** folderul (cu `karios`) pe stick. | Copierea reușește. | | |
| B2 | Porniți `pontix.exe` de pe stick, pe PC 1. | Pornește; datele de la A sunt acolo. | | |
| B3 | Introduceți stickul în PC 2 și porniți aplicația de pe stick. | Pe un alt utilizator Windows, datele sunt blocate: aplicația cere parola contului sau cheia de recuperare. Cu parola, se deschid. | | |
| B4 | Lăsați aplicația deschisă de pe stick pe PC 2 și porniți-o a doua oară. | A doua fereastră se deschide **doar în citire** (nu poate salva). | | |
| B5 | Închideți aplicația și scoateți stickul în siguranță. Porniți din nou de pe stick. | Datele sunt intacte. | | |

## C. Actualizarea de la o versiune veche cu date (PC 1)

Treceți prin migrările noi ale bazei de date (0037–0042).

| # | Proba | Rezultat așteptat | T/E/N | Observații |
| --- | --- | --- | --- | --- |
| C1 | În versiunea veche, faceți un backup din Setări. Închideți aplicația și păstrați o copie a întregului folder vechi. | Backupul există; copia folderului e făcută. | | |
| C2 | Extrageți versiunea nouă într-un folder separat. Copiați acolo **tot** folderul `karios` vechi. Porniți. | Pornește; în `karios\backup` apare automat o copie făcută înainte de actualizare. | | |
| C3 | Verificați angajații, o lună de pontaj, concediile, planificarea. | Toate sunt ca în versiunea veche. | | |
| C4 | Setări › Sărbători legale: puntea adăugată. | Apare cu **Tip = Zi liberă adăugată**; sărbătorile legale apar ca **Sărbătoare legală**. | | |
| C5 | Setări › Schimburi: programul cu pauză. | Are **Cum se scade pauza = Nestabilit**; în Pontaj, orele lui apar „—”, cu motivul. | | |
| C6 | Rapoarte: generați Pontaj detaliat pe luna veche. | Raportul se generează; cifrele sunt cele de dinainte, iar unde coloana nu se aplică celula e goală. | | |
| C7 | Porniți din nou versiunea veche pe **copia** păstrată la C1. | Versiunea veche merge în continuare pe copia ei. (Pe folderul deja actualizat nu se mai întoarce.) | | |

## D. Funcțiile noi, pe scurt

| # | Proba | Rezultat așteptat | T/E/N | Observații |
| --- | --- | --- | --- | --- |
| D1 | Alegeți modul pauzei pentru programul cu pauză (de exemplu „Din total”). Salvați. | În Pontaj, orele lui apar fără pauză (de exemplu 8 în loc de 8,5). | | |
| D2 | Pontaj: o zi de 10 ore pentru un angajat cu normă de 8. | Coloana **OS** arată 2. Pe un rând de gardă, OS e gol. | | |
| D3 | Configurare pontaj: „Zilele libere adăugate în calendar, în ZNL” = „Nu intră”. Salvați. | Orele lucrate în punte nu mai intră în ZNL. | | |
| D4 | Pontaj › Mai multe › **Verifică luna**. | Apare lista constatărilor (de exemplu programele cu pauză nestabilită); nu se schimbă nimic. | | |
| D5 | Pontaj › Mai multe › **Export pentru salarizare**, cu „CSV generic” și cu „Excel generic”. | Fișierele se deschid în Excel; diacriticele sunt corecte; un rând pe marcă; orele și zilele corespund Pontajului. | | |
| D6 | Rapoarte: toate secțiile, bifa **O pagină nouă pentru fiecare secție**. | Fiecare secție începe pe o pagină nouă în PDF. | | |
| D7 | Fișa unui angajat › **Istoric normă și concediu**: o schimbare de normă de la mijlocul lunii. | Ore lucrătoare se calculează cu norma veche până la dată și cu cea nouă după. | | |
| D8 | Setări › General › Limbă: **Scrie fișierul de lucru** pentru `en`, schimbați un text în fișier, **Verifică**, **Importă**. | Textul schimbat apare în interfață în engleză; **Șterge traducerile importate** îl readuce. | | |
| D9 | Configurare pontaj: „Turele din nopțile cu schimbarea orei” = „Orele efective”. O tură 19–07 în noaptea de 24/25 octombrie. | Tura are 13 ore (cu „Durata din program”, 12). | | |

## E. Siguranța datelor

| # | Proba | Rezultat așteptat | T/E/N | Observații |
| --- | --- | --- | --- | --- |
| E1 | Faceți un backup, schimbați o celulă, restaurați backupul, reporniți. | Celula revine la valoarea din backup. | | |
| E2 | Închideți o lună; încercați să modificați o celulă; redeschideți luna, cu motiv. | Luna închisă nu se poate modifica; redeschiderea cere motiv și apare în jurnal. | | |
| E3 | Deblocați datele cu **cheia de recuperare** (pe PC 2 sau alt utilizator Windows). | Datele se deschid. | | |
| E4 | Meniu › Raportează o problemă, cu date fictive. | Raportul se poate citi înainte de trimitere și nu conține parole sau chei; nu se trimite nimic automat. | | |

## Concluzie

Probleme care blochează livrarea: ________________________________________________

Alte observații: _______________________________________________________________

Rezultat: ☐ Acceptat pentru folosire  ☐ Acceptat cu rezerve  ☐ Neacceptat

Semnătura: __________________

Fișa nu autorizează publicarea și nici folosirea datelor pentru salarizare fără confirmarea responsabilului.
