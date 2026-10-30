# Pontix — Manual de utilizare

Ediția {{EDITIE}} · Aplicația {{VERSIUNE}} · Build: {{DATA_BUILD}} · Commit: {{COMMIT}}.

Pentru operatorii de pontaj și administratorii aplicației. Exemplele sunt fictive. Manualul descrie această versiune; denumirile și pașii pot diferi în versiuni mai vechi.

Versiunile HTML și PDF includ aceleași articole folosite de ajutorul din aplicație. Butonul «Cum să faci» din antet deschide ghidurile secțiunii; căutarea găsește articole din toate secțiunile. În AI, întrebările de utilizare primesc pași și surse fără motorul local pornit.

## 1. Început rapid: ordinea lucrului

Pentru o lună nouă, urmează această ordine:

1. Verifică anul, luna, lista angajaților și categoriile de program.
2. Verifică zilele libere în calendarul anului.
3. Înregistrează concediile și absențele cunoscute.
4. Completează programul obișnuit în Pontaj, apoi salvează schița.
5. Stabilește turele speciale în Planificare și verifică publicarea lor în Pontaj.
6. Compară programul planificat cu pontajul și rezolvă diferențele.
7. Generează și verifică rapoartele, apoi exportă documentele necesare.
8. Creează un backup și închide luna când datele sunt verificate.

**Înaintea oricărei completări în lot, verifică persoanele și perioada.** Un filtru poate ascunde persoane; un număr de ore afișat nu identifică întotdeauna programul din spatele lui.

## 2. Pornire, autentificare și recuperare

### Prima pornire

1. Deschide `pontix.exe` dintr-un folder în care ai drept de scriere.
2. Dacă aplicația nu are conturi, creează primul administrator.
3. Completează numele de utilizator, parola și cele trei întrebări diferite de recuperare, cu răspunsurile lor.
4. Parola trebuie să aibă cel puțin 8 caractere, o literă și o cifră. Repetă parola în câmpul de confirmare.
5. După autentificare, notează cheia de recuperare afișată și retasteaz-o pentru confirmare. Lipirea cheii este dezactivată în această etapă.

Cheia de recuperare are 25 de caractere, în cinci grupuri de câte cinci. La confirmare și deblocare, cratimele se adaugă automat, iar literele devin majuscule: tastează caracterele cheii. Păstreaz-o într-un loc sigur, separat de calculator. Nu este parola contului. Până la confirmarea cheii, administratorul vede ecranul de confirmare, iar ceilalți utilizatori sunt îndrumați să apeleze la administrator.

### Pornirea obișnuită

Aplicația poate redeschide datele folosind protecția contului Windows și poate relua sesiunea locală, conform stării salvate. De aceea nu se cere întotdeauna parola la fiecare pornire. Din **Contul meu** poți schimba parola și întrebările de recuperare sau poți ieși din cont.

### Alt calculator sau alt utilizator Windows

La mesajul **Deblochează datele**, folosește numele și parola unui cont valid ori cheia de recuperare. După deblocare, aplicația pregătește accesul pentru utilizatorul Windows curent. Dacă ai schimbat calculatorul după o actualizare majoră a protecției datelor, autentifică mai întâi conturile pe calculatorul obișnuit.

Dacă ai uitat parola, folosește recuperarea prin întrebările contului. Administratorul poate administra conturile și poate genera o cheie de recuperare nouă; după înlocuire, cheia veche nu mai funcționează. Pierderea tuturor căilor de acces poate face datele nerecuperabile.

## 3. Orientarea în aplicație

Meniul principal conține **Angajați**, **Pontaj**, **Planificare**, **Concedii și absențe**, **Rapoarte**, **Asistent AI** și **Setări**.

În Setări găsești și zilele libere, schimburile, categoriile de program, configurarea pontajului, definițiile rapoartelor, jurnalul modificărilor și, pentru administrator, conturile. În unele ecrane denumirile sunt grupate în file sau subsecțiuni.

- **Anul și luna** stabilesc perioada ecranului curent. Nu presupune că toate ecranele au rămas pe aceeași perioadă.
- **Căutarea și filtrele** restrâng persoanele afișate. Dacă lipsește un angajat, golește mai întâi căutarea și filtrele.
- **Tema** poate urma inițial Windows. După ce alegi manual tema luminoasă sau întunecată, alegerea este păstrată.
- **Zoomul grilei** schimbă densitatea afișării, fără să schimbe datele.
- **Salvează** confirmă modificările formularului sau ale schiței. Schimbarea unui câmp nu înseamnă întotdeauna salvare.
- **Anulează pasul / Refă pasul** privesc modificările din sesiunea ecranului. Nu înlocuiesc restaurarea unui backup sau anularea unei operații AI salvate.

La avertizarea despre modificări nesalvate, salvează sau confirmă abandonarea lor. Revenirea ulterioară nu recuperează automat ce ai abandonat.

## 4. Angajați: adăugare, modificare și import

### Adaugă un angajat

1. Deschide **Angajați** și apasă butonul pentru adăugare.
2. Introdu numele, prenumele și marca. Folosește o marcă unică.
3. Alege tipul contractului. Pentru normă parțială completează orele zilnice; pentru normă întreagă verifică norma globală din Setări.
4. Alege secția, departamentul și categoria de program corespunzătoare.
5. Completează celelalte informații necesare, inclusiv dreptul anual de concediu dacă este cunoscut.
6. Salvează și verifică în listă persoana și marca.

**Dezactivarea** păstrează istoricul. Un angajat inactiv poate apărea în perioade pentru care are deja înregistrări. Pentru o informație greșită folosește editarea; nu crea o a doua persoană cu altă marcă pentru același angajat.

### Import din Excel sau CSV

1. Deschide **Importă din Excel / CSV** și descarcă șablonul din aplicație.
2. Completează șablonul, păstrând antetele așteptate. Nu folosi automat un fișier cu alte coloane.
3. Selectează fișierul și așteaptă analiza.
4. Verifică erorile pe rânduri, numărul de persoane noi și actualizările propuse.
5. Verifică opțiunea de actualizare a înregistrărilor existente înainte de confirmare.
6. Execută importul și verifică rezultatul în listă.

Datele demonstrative sunt pentru testare. Încarcă-le doar într-o copie separată destinată probelor, nu în evidența operațională.

## 5. Pregătirea programelor și calendarului

### Secții, departamente și categorii

În **Setări**, pregătește secțiile și departamentele. Atribuie angajatului categoria de program potrivită. Categoria stabilește programele permise și, dacă este definit, programul implicit folosit la completare. O categorie nu este o regulă salarială.

Pentru schimbarea categoriei mai multor angajați, folosește selecția și atribuirea în lot din lista angajaților, verificând persoanele bifate înainte de aplicare.

### Schimburi

În **Setări → Schimburi**, verifică programul după cod, interval, durată, regim și starea activă. Denumirile și intervalele sunt configurabile; tabelul următor reprezintă setul inițial, nu înlocuiește configurația propriei baze.

| Cod | Interval inițial | Durată nominală |
| --- | --- | --- |
| Z7 | 07:00–14:00 | 7 ore |
| S7 | 15:00–22:00 | 7 ore |
| N9 | 22:00–07:00 a doua zi | 9 ore |
| IN9 | 22:00–07:00 a doua zi, intrare specială în zi liberă | 9 ore |
| Z12 | 07:00–19:00 | 12 ore |
| N12 | 19:00–07:00 a doua zi | 12 ore |
| G24 | 07:00–07:00 a doua zi | 24 ore |

IN9 se introduce în ziua în care începe tura: de exemplu, vineri seara pentru intrarea în sâmbătă. Regimul zilei libere nu se deduce doar din numărul de ore.

Un program dezactivat rămâne în istoric. Nu redenumi și nu modifica intervalele unui program folosit fără să verifici efectul asupra datelor și rapoartelor. Un program cu pauză nerepartizată poate rămâne necalculabil.

### Zile libere

Deschide **Setări → Zile libere legale**, alege anul și verifică lista. Poți adăuga, modifica sau importa datele folosind șablonul oferit. Verifică și anul următor pentru turele care trec peste 31 decembrie. Butonul de completare pentru 2026 completează lista aferentă acelui an; nu presupune existența listelor pentru orice an.

## 6. Concedii și absențe: înregistrarea corectă

1. Deschide **Concedii și absențe** și alege anul.
2. Apasă **Înregistrează zile** și selectează angajatul, tipul și intervalul.
3. Verifică datele bifate. Weekendurile și sărbătorile pot fi selectate; folosește explicit selecția luni–vineri fără sărbători numai dacă se potrivește cererii și programului persoanei.
4. Completează numărul de înregistrare, aprobarea și înlocuitorul, după caz. Cererile deschise din AI cer număr de înregistrare înainte de salvare.
5. Poți atașa cererea semnată în PDF, PNG sau JPG, de cel mult 10 MB.
6. Verifică numărul zilelor și salvează. Recitește înregistrările și starea aprobării.

**În așteptare**, **Aprobată** și **Respinsă** au semnificații diferite. Salvarea unei cereri nu înseamnă automat aprobarea ei. Un răspuns conversațional al AI-ului nu înregistrează concediul.

Dacă în aceeași zi există deja program sau o absență incompatibilă, aplicația poate refuza salvarea. Verifică pontajul, tura începută în ziua precedentă și înregistrarea existentă, apoi corectează explicit datele. Nu completa ore peste o celulă de absență protejată.

### Modifică sau retrage o cerere

Folosește **Editează** pentru o zi sau selecția zilelor cererii și **Completează cererea / aprobarea** pentru metadate comune. **Selectează cererea** ajută la bifarea zilelor care aparțin aceluiași document sau număr de înregistrare. Verifică persoana și datele înainte de **Anulează ziua** ori **Anulează zilele selectate**.

Documentul atașat poate fi previzualizat sau descărcat. Sumarul anual prezintă înregistrările și dreptul contractual curent; diferența afișată nu rezolvă automat reportările sau indemnizațiile.

## 7. Pontaj: completează și salvează luna

### Completează un program

1. Deschide **Pontaj** și alege anul și luna.
2. Golește filtrele dacă vrei să lucrezi pe întreaga echipă sau restrânge explicit persoanele dorite.
3. Selectează celula ori dreptunghiul de zile și persoane.
4. Deschide **Completare**, alege programul și apasă **Aplică programul**.
5. Verifică în inspecția de deasupra grilei codul, intervalul, data și persoana.
6. Apasă **Salvează schița** și verifică mesajul de succes.

**Z7 și S7 au ambele 7 ore.** Grila poate afișa „7” în ambele cazuri. Codul din inspecția celulei arată dacă este program de zi sau de seară. O schimbare Z7 → S7 poate fi corect salvată chiar dacă numărul din celulă nu se schimbă.

### Completează ore manuale

Alege selecția, deschide Completare și introdu orele manuale pe zi. Sunt acceptate valori până la 48 de ore, cu precizie de un minut, folosind punct sau virgulă ca separator zecimal. Apasă **Aplică orele manuale**, verifică și salvează. Orele manuale nu identifică un program și nu permit aceleași verificări de interval ca un schimb cu oră de început.

F2 sau dublu clic deschide completarea manuală. Câmpul manual gol ori Delete golește înregistrarea din selecție. Verifică selecția înainte de această acțiune.

### Completează în lot

**Precompletează L–V** completează celulele eligibile conform programului ales. **După categorii** folosește categoriile atribuite persoanelor și arată o previzualizare. Verifică persoanele, intervalul și motivele zilelor omise. Aplicarea intră în schiță; apasă apoi **Salvează schița**.

Celulele cu absențe sunt protejate. Copierea și lipirea din clipboard sau Excel (Ctrl+C / Ctrl+V) acceptă atât ore zecimale, cât și coduri de schimb cunoscute (ex: G24, 12, 8). Clic dreapta pe o celulă deschide un meniu contextual rapid cu schimburi frecvente, ștergere, copiere/lipire și salt la concedii. În antetul grilei, rândul de acoperire (mini-hartă heatmap) arată personalul prezent pe fiecare zi. Butonul **Asistent AI** din bara de unelte deschide un panou lateral retractabil (drawer) pentru interogări și planificări directe pe grilă.

### Scurtături de lucru

| Acțiune | Tastatură / interacțiune |
| --- | --- |
| Deplasare | Săgeți, Enter, Tab |
| Selectarea unui dreptunghi | Shift + săgeți, Shift + clic sau tragere cu mouse-ul |
| Copiere / Lipire selecție | Ctrl+C / Ctrl+V |
| Meniu contextual rapid | Clic dreapta pe celulă sau pe selecție |
| Anularea/refacerea unui pas nesalvat | Ctrl+Z / Ctrl+Y |
| Copiere/lipire de ore | Ctrl+C / Ctrl+V |
| Completare manuală | F2 sau dublu clic |

## 8. Planificare: ture, publicare și comparație

Planificarea păstrează programul intenționat. Pontajul păstrează schița de lucru folosită la ore și rapoarte. Nu considera automat cele două ecrane identice.

### Planifică o tură specială

1. Deschide **Planificare** și verifică anul, luna și filtrele.
2. În grila lunară, alege programul special, de exemplu S7, și celulele vizate.
3. Verifică persoanele, datele și programele. Folosește anularea/refacerea pașilor dacă este necesar.
4. Salvează planificarea prin comenzile ecranului și verifică starea pozițiilor.
5. Deschide **Previzualizează și publică** pentru a transfera planul în schița de pontaj.
6. Verifică valorile dinainte/după și celulele excluse. Dacă dorești înlocuirea unor valori existente, alege explicit această opțiune.
7. Confirmă publicarea și verifică aceeași persoană și aceeași zi în Pontaj.

O tură nu poate suprascrie o absență sau o înregistrare incompatibilă. Programele trebuie permise de categoria angajatului. Dacă datele s-au schimbat după previzualizare, reia verificarea.

### Interpretează stările

- **Nepublicată:** poziția este în plan, fără o publicare curentă corespunzătoare.
- **Publicată:** versiunea publicată corespunde schiței verificate.
- **Modificată:** planul ori schița nu mai corespund ultimei publicări.
- **Anulată:** poziția a fost dezactivată, cu istoric păstrat.

În comparația planificat–realizat, verifică separat turele identice, cele schimbate, înregistrările neplanificate, turele neefectuate și absențele. Folosește vederile de acoperire, echitate sau regulile disponibile ca sprijin pentru verificare, nu ca aprobare automată a planului.

### Sincronizare la salvarea din Pontaj

Un program special introdus manual în Pontaj, de exemplu S7, actualizează și Planificarea când apeși **Salvează schița**. Sincronizarea are loc la salvare, nu în timpul tastării. S7 apare în secțiunea **După-amiază**. Revenirea la un program normal sau gol retrage numai poziția care corespunde programului anterior.

Pentru un program special salvat într-o versiune mai veche, selectează celula și apasă **Actualizează planificarea** în bara de detalii. Salvează întâi celelalte modificări din grilă. Repetarea nu creează dubluri și nu schimbă orele.

O planificare divergentă, modificată separat, nu este suprascrisă. Absențele, suprapunerile, categoriile, lunile închise și reviziile sunt verificate; la eroare întregul lot este anulat. Retragerea unei poziții direct din Planificare nu șterge automat schița publicată anterior.

## 9. Ore, totaluri și ture peste miezul nopții

În **Setări → Configurare pontaj**, configurația stabilește repartizarea orelor și sursele coloanelor. Modificarea se poate reflecta în perioadele recalculate și în rapoartele generate ulterior. Salvează configurația numai după verificarea regulilor folosite de organizație.

| Afișaj | Cum se citește |
| --- | --- |
| T / Ore lucrate | Totalul orelor conform configurației |
| ZNL | Ore atribuite intervalelor de weekend sau zile libere conform calendarului și regulilor |
| Ore norm | Categoria configurată de ore; nu este automat norma contractuală lunară |
| 0 | Regulă definită, fără valoare relevantă în datele analizate |
| — | Rezultat nedefinit sau necalculabil; verifică regula și programul |

Pentru o tură 22:00–07:00, configurația poate afișa toate cele 9 ore în ziua începerii sau le poate împărți 2 ore și 7 ore pe zile calendaristice. Orele dintr-o celulă pot proveni din tura începută în ziua precedentă; modifică înregistrarea de la data începerii turei.

Duratele sunt nominale. Nu presupune ajustarea automată pentru schimbarea orei ori pentru poziționarea pauzelor. O coloană numită „Ore gardă 200%” nu aplică automat un factor de plată. Procentele, indemnizațiile și interpretările salariale trebuie să provină din configurația validată de contabilitate.

## 10. Rapoarte, Excel și tipărire

### Generează un PDF

1. Salvează mai întâi datele ecranului în care ai lucrat.
2. Deschide **Rapoarte** și alege modelul, perioada și limba documentului.
3. Alege filtrele și angajatul, dacă raportul îl cere.
4. Apasă **Generează și previzualizează**.
5. Verifică antetul, perioada, persoanele, paginile și totalurile.
6. Alege **Salvează PDF** pentru fișier sau **Deschide pentru tipărire** pentru aplicația de PDF din Windows.

Tipărirea se confirmă în aplicația PDF, alegând imprimanta și paginile. Schimbarea filtrelor impune o nouă generare. Un PDF deja exportat nu se actualizează când modifici ulterior datele.

Din **Pontaj → Export Excel (.xlsx)** generezi registrul Excel al perioadei. Verifică fișierul exportat înainte de utilizare. Exportul de documente nu înlocuiește backupul bazei și nu dovedește închiderea lunii.

### Personalizează un model

În **Setări → Definire rapoarte**, duplică un model când vrei o variantă independentă. Configurează antetul, coloanele, subsolul, pagina și textele în limbile necesare. Poți folosi imaginile acceptate de editor. Salvează, apoi generează din nou pentru a verifica paginarea. Folosește doar câmpurile de înlocuire afișate în editor.

## 11. Asistent AI: comenzi și confirmare

Asistentul rulează local. Unele cereri sunt recunoscute direct de aplicație, iar altele folosesc modelul. Motorul trebuie să fie disponibil pentru răspunsurile care necesită inferență; planurile și verificările deterministe pot avea flux separat.

### Pregătește o comandă clară

Precizează persoana sau grupul, perioada cu an, codul programului, zilele vizate și ce se întâmplă cu valorile existente. Folosește marca pentru a evita confuzia între nume asemănătoare.

Exemplu fictiv de planificare:

> Planifică toți angajații pe Z7 în noiembrie 2026, luni–vineri, păstrează înregistrările existente.

1. Verifică **Luna de lucru** înainte de trimitere.
2. În **Plan activ → Interpretare**, corectează grupul, persoanele, programul și perioada.
3. Confirmă anul dacă a fost preluat din context sau din șablon.
4. Alege zilele și păstrarea ori înlocuirea celulelor existente.
5. Apasă **Arată modificările** și verifică rezumatele și tabelul.
6. Bifează verificarea și confirmă aplicarea.
7. Verifică rezultatul și deschide pontajul lunii indicate.

**Salvat**, **Păstrat**, **Neaplicat** și **Deja identic** descriu rezultate diferite. Cererea „toți angajații” nu garantează modificarea tuturor celulelor: pot exista absențe, programe nepermise, celule păstrate sau luni închise. Verifică explicațiile de selecție și motivele pe zile.

### Solicită un concediu

Exemplu fictiv:

> Ana Exemplu merge în concediu de odihnă între 26 și 31 octombrie 2026.

Cererea recunoscută afișează **Deschide formularul de concediu**. Verifică angajatul, tipul și datele în formular, completează numărul de înregistrare și salvează. Starea inițială este în așteptare. Dacă lipsesc date, completează clarificările. Pentru o persoană ambiguă, folosește marca sau numele complet.

Întrebările despre cine lipsește folosesc perioada disponibilă; schimbă luna de lucru când întrebi despre altă lună. Un răspuns liber precum „nu ai permisiunea” nu trebuie tratat singur ca diagnostic al contului: verifică formularul și eroarea concretă a aplicației. Dacă AI-ul nu recunoaște cererea, folosește direct **Concedii și absențe**.

### Motor și conversație

În **Opțiuni motor** verifică starea, modelul efectiv activ și contextul configurat. Limita configurată este 32768 de tokenuri; aceasta nu este o măsurătoare a memoriei sau a contextului consumat. Prima cerere poate dura mai mult din cauza pornirii și încălzirii motorului. Nu trimite repetat aceeași comandă doar fiindcă răspunsul întârzie; urmărește starea și folosește oprirea dacă este necesar.

**Conversație nouă** golește istoricul conversației, nu pontajul. Repetă perioada și persoanele necesare în noua cerere. Motorul acceptă executabilul verificat de aplicație; un fișier `llama-server.exe` înlocuit cu altul poate fi refuzat. Folosește pachetul AI compatibil, fără a instala arbitrar alt binar peste el.

## 12. Instrumentele AI: șabloane, verificare, export și anulare

### Șabloane

În **Șabloane de planificare**, salvează parametrii sub un nume ușor de recunoscut. Încărcarea unui șablon deschide editorul, cere verificarea perioadei și nu aplică automat modificări. Ștergerea șablonului nu șterge pontajul. Șabloanele sunt disponibile și în **Instrumente**.

### Verificarea pontajului

Deschide **Plan nou / verificare**, alege persoanele, perioada și zilele și folosește **Verifică fără modificări**. Citește semnalele pentru goluri, coduri necunoscute, programe nepermise, absențe sau suprapuneri. Orele manuale nu pot fi verificate complet ca intervale. Regulile configurate se evaluează conform implementării și nu reprezintă o certificare salarială sau legală. Funcția nu corectează automat datele.

### Exportul unui plan

**Exportă planul CSV** salvează local rezumatul și toate rândurile planului, nu doar pagina filtrată. O propunere este marcată ca neaplicată. Rezultatul unei operații salvate este istoric; verifică starea exportului dacă operația a fost anulată ulterior.

### Anularea ultimei operații AI

1. Deschide **Instrumente → Verifică ultima operație**.
2. Citește valorile de restaurat și eventualele conflicte.
3. Bifează confirmarea și apasă **Anulează operația AI** numai după verificare.
4. Verifică rezultatul în pontaj.

Sunt disponibile operațiile pentru care aplicația a păstrat instantanee, nu orice acțiune istorică sau mesaj din conversație. O editare ulterioară ori o lună închisă poate bloca anularea integrală. Nu forța restaurarea prin copierea manuală a unor fișiere vechi peste baza deschisă. Jurnalul operației rămâne păstrat.

## 13. Închiderea și redeschiderea lunii

1. Salvează schița și verifică persoanele, absențele, planificarea și rapoartele.
2. Creează un backup.
3. În Pontaj, apasă **Închide luna (salarii)**.
4. Completează motivul și confirmă.
5. Verifică starea lunii și documentele rezultate din datele finalizate.

Închiderea blochează modificările vizate de starea lunii. Nu calculează și nu plătește salarii doar prin apăsarea butonului. Pentru o corectură ulterioară, folosește redeschiderea explicită, cu motiv, modifică datele, regenerează documentele și închide din nou conform procedurii organizației.

## 14. Backup, restaurare și actualizarea aplicației

### Creează o copie de siguranță

În **Setări**, folosește **Creează backup** și verifică mesajul cu locația. Există și copii automate, inclusiv la prima pornire din zi și înaintea unor operații de actualizare. Nu te baza exclusiv pe copiile de pe același disc: păstrează periodic o copie separată a folderului de date și cheia de recuperare într-un loc sigur.

### Restaurează din aplicație

1. Ca administrator, deschide zona de restaurare din Setări.
2. Alege copia după dată și verifică rezultatul verificării ei.
3. Citește avertizarea: datele curente vor fi înlocuite cu cele din copia aleasă.
4. Confirmă programarea restaurării.
5. Închide aplicația și pornește-o din nou; restaurarea se aplică la următoarea pornire.
6. Verifică persoanele, perioada și datele. Aplicația păstrează o copie de siguranță a datelor de dinaintea restaurării.

O copie modificată din afara aplicației poate fi refuzată. Dacă te răzgândești înainte de repornire, folosește anularea restaurării programate.

### Mută sau actualizează aplicația

Înainte de actualizare, creează un backup din Setări, verifică locația lui și păstrează separat cheia de recuperare. Copia bazei nu înlocuiește copia completă a folderului de date.

Aplicația portabilă păstrează datele în folderul `karios` de lângă executabil. Un executabil pus într-un director nou nu mută automat baza veche.

1. Închide aplicația și păstrează o copie completă a folderului vechi.
2. Pune executabilul nou într-un director separat.
3. Pentru continuarea datelor existente, copiază întregul folder `karios` lângă executabilul nou, nu doar baza de date.
4. Pornește aplicația și parcurge eventuala deblocare și confirmarea cheii de recuperare.
5. Verifică datele înainte de a renunța la copia veche.

Folderul include baza, seiful și subdirectoare precum `backup`, `ai`, `rapoarte`, `jurnale` și `webview`. Oglinda seifului din `backup/seif` ajută la refacerea seifului pierdut sau deteriorat, dar nu înlocuiește parolele ori cheia de recuperare. Nu o trata drept parolă și nu o șterge ca pe un fișier temporar.

## 15. Jurnal, conturi și probleme de integritate

În **Setări → Jurnalul modificărilor**, filtrează după tabel, angajat, lună sau text. Citește valorile vechi/noi, momentul și autorul. Modificările confirmate prin AI sunt identificate în jurnal. Exportul CSV este pentru consultare și nu înlocuiește backupul.

Administratorul gestionează **Conturi**, rolurile, activarea și resetarea parolelor. Rolul de administrator al conturilor nu înseamnă că AI-ul poate decide singur orice modificare de date. Folosește confirmările formularelor.

Dacă apare un incident de audit, scrierea poate fi blocată. Verifică dacă baza a fost copiată, înlocuită cu una veche sau modificată în afara aplicației. Păstrează fișierele și cere verificarea persoanei responsabile; confirmă incidentul numai după stabilirea cauzei și completarea motivului cerut. Restaurarea din aplicație este preferabilă suprascrierii manuale a bazei.

Pentru suport, notează versiunea aplicației, ecranul, anul/luna, pașii și eroarea exactă. Jurnalele de diagnoză sunt în `karios/jurnale`; nu atașa parole, chei, cereri semnate sau baza cu date personale la un raport public.

## 16. Probleme frecvente și pași de rezolvare

| Situație | Ce verifici / ce faci |
| --- | --- |
| Am schimbat Z7 în S7, dar în celulă văd tot 7 | Verifică programul din inspecție; ambele au 7 ore. Salvează schița. |
| S7 apare în Pontaj, nu și în Planificare | Salvează schița. Pentru o celulă salvată anterior, selecteaz-o și apasă Actualizează planificarea. Verifică aceeași lună și persoană și rezolvă eventualul conflict semnalat. |
| Nu apare un angajat | Golește filtrele și căutarea; verifică starea activă și perioada. |
| Celula nu poate fi editată | Verifică absența activă și starea lunii. Modifică absența din modulul ei. |
| Salvarea raportează conflict | Datele s-au schimbat între timp. Păstrează separat ce ai de reaplicat, reîncarcă și verifică noua stare. |
| Previzualizarea AI a expirat | Reia previzualizarea și verifică diferențele înainte de o nouă confirmare. |
| AI-ul a răspuns cu o listă în loc să înregistreze concediul | Folosește o cerere explicită, cu nume/marcă și date; verifică formularul de concediu. Dacă nu apare, folosește direct modulul. |
| Eroare de context AI | Pornește o conversație nouă și trimite o cerere concentrată, cu perioada și persoanele. Verifică versiunea aplicației și motorul compatibil. |
| Totalul este „—” | Verifică programul, pauza și sursa coloanei în configurare. Nu înlocui automat „—” cu zero. |
| Nu pot restaura un backup | Verifică rolul, copia și mesajul de integritate; nu modifica fișierele pentru a ocoli eroarea. |
| După mutare, aplicația pare goală | Verifică dacă întregul folder `karios` a fost copiat lângă executabilul pornit. Închide aplicația înainte de transfer. |
| Motorul AI este necunoscut | Folosește motorul livrat și verificat pentru aplicație; nu schimba arbitrar executabilul. |

## 17. Verificare înainte de predarea lunii

- Anul și luna sunt corecte în fiecare ecran și raport.
- Lista persoanelor este completă; filtrele nu ascund înregistrări relevante.
- Concediile au datele, aprobările și documentele corecte.
- Programele sunt verificate după cod, nu numai după ore.
- Diferențele planificat–realizat și semnalele de verificare au fost analizate.
- Schițele și formularele sunt salvate; rapoartele au fost regenerate după ultimele corecții.
- Totalurile nedefinite au fost clarificate cu persoana responsabilă de reguli.
- Backupul este verificat, iar accesul la cheia de recuperare este asigurat.
- Închiderea lunii este consemnată conform procedurii organizației.

## 18. Glosar

| Termen | Înțeles în aplicație |
| --- | --- |
| Marcă | Identificatorul angajatului, util inclusiv în comenzi AI |
| Schimb / program | Cod cu durată și, după caz, oră de început și regim |
| Schiță | Înregistrările de lucru salvate din grila Pontaj |
| Planificare | Turele intenționate, păstrate separat |
| Publicare | Transferul verificat al planificării în schița de pontaj |
| Previzualizare | Imaginea modificărilor propuse înainte de confirmare |
| Revizie / conflict | Protecție împotriva rescrierii unei modificări făcute între timp |
| Seif | Fișierul care păstrează copiile protejate ale cheii datelor |
| Cheie de recuperare | Secret de recuperare distinct de parola contului |
| Backup | Copie de siguranță, folosită pentru refacerea unei stări anterioare |

## 19. Despre acest manual

Acoperire: meniurile principale și subsecțiunile din Setări, fluxurile de date și instrumentele AI din versiunea {{VERSIUNE}}. Revizia editorială 1.5 include optimizările UI din Pontaj: meniul contextual clic dreapta, clipboardul avansat, mini-harta de acoperire pe zile și panoul drawer pentru Asistentul AI. Include ajutorul contextual pe toate ecranele, 21 de articole comune manualului și AI, sincronizarea Pontaj → Planificare și formatarea cheii de recuperare. Referința exactă a livrării este commitul de pe copertă. Actualizarea versiunii nu înlocuiește revizuirea instrucțiunilor: buildul cere confirmarea versiunii revizuite în sursele manualului.

Manualul nu substituie regulile de salarizare și aprobările organizației. Nu certifică validarea interactivă a fiecărui ecran pe calculatorul utilizatorului. În cazul unei diferențe între instrucțiuni și aplicație, notează versiunea și cere verificarea fluxului înainte de aplicări în lot.

Manualele Markdown, HTML și PDF sunt regenerate pentru fiecare build și salvate lângă executabil, cu versiunea în nume. HTML-ul se deschide în browser fără internet și are căutare locală. PDF-ul este pentru citire și tipărire. Versiunea, ediția, data buildului și commitul identifică documentele; manifestul build.json păstrează și sumele de verificare.

Ajutorul contextual și accesul la articole din AI sunt disponibile din această versiune. Reformularea opțională cu modelul și tururile vizuale automate rămân extensii viitoare. Instrucțiunile nu execută operații și nu înlocuiesc validările sau confirmările.

## 20. Ghiduri interactive comune aplicației și asistentului

Articolele de mai jos sunt generate din catalogul comun RO/EN. Aceleași instrucțiuni sunt afișate în panoul «Cum să faci» și în răspunsurile de utilizare ale asistentului. Identificatorii stabili permit citarea și deschiderea articolului.

{{ARTICOLE_AJUTOR}}

### Clarificari de operare pentru 0.88.15

Indicatorul zilnic numara persoanele cu inregistrari din lista filtrata; nu confirma necesarul de personal sau deficitul. Panoul AI cere confirmare inainte de schimbarea perioadei ori deschiderea concediilor daca exista modificari nesalvate. Ghidurile functioneaza fara model; configurarea motorului este pliabila. Pe ecrane inguste panoul AI se suprapune peste grila si se inchide din antet. Shift+F10 deschide actiunile celulei; Tab parcurge butoanele, Escape revine la celula. Cautarea in ghiduri afiseaza un rezultat curent; o citare veche ofera deschiderea ghidului versiunii actuale.

### Pachetul de livrare 0.88.16

Buildul include o arhiva ZIP cu executabilul, manualele, manifestul, LICENSE, instructiunile ACTUALIZARE.md, NOTE-VERSIUNE.md si CHECKLIST-ACCEPTANTA.md. SHA256SUMS.txt verifica fisierele incluse; fisierul .zip.sha256 de langa arhiva verifica descarcarea. Pachetul standard nu include baza, seiful, backupuri sau modele AI. Extrageti arhiva inainte de pornire si urmati instructiunile de actualizare. Checklistul este de completat de operator, nu o lista de teste deja executate.
