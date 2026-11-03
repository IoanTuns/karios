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
4. Alege secția, compartimentul și categoria de program corespunzătoare.
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

### Secții, compartimente și categorii

În **Setări**, pregătește secțiile și compartimentele. Atribuie angajatului categoria de program potrivită. Categoria stabilește programele permise și, dacă este definit, programul implicit folosit la completare. O categorie nu este o regulă salarială.

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

În editorul persoanelor, folosește **Arată doar selectați** pentru a verifica lista aleasă. În previzualizare, **Doar excepții** izolează celulele ocupate, blocate sau identice; dezactivează filtrul pentru a revedea toate rândurile. Rezumatul arată distinct persoanele efectiv afectate de modificări; numărul nu este neapărat egal cu numărul persoanelor selectate. Verifică aceste diferențe înainte de confirmare.
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

Acoperire: meniurile principale și subsecțiunile din Setări, fluxurile de date și instrumentele AI din versiunea {{VERSIUNE}}. Revizia editorială 1.8 include optimizările UI din Pontaj: meniul contextual clic dreapta, clipboardul avansat, mini-harta de acoperire pe zile și panoul drawer pentru Asistentul AI. Include ajutorul contextual pe toate ecranele, 21 de articole comune manualului și AI, sincronizarea Pontaj → Planificare și formatarea cheii de recuperare. Referința exactă a livrării este commitul de pe copertă. Actualizarea versiunii nu înlocuiește revizuirea instrucțiunilor: buildul cere confirmarea versiunii revizuite în sursele manualului.

Manualul nu substituie regulile de salarizare și aprobările organizației. Nu certifică validarea interactivă a fiecărui ecran pe calculatorul utilizatorului. În cazul unei diferențe între instrucțiuni și aplicație, notează versiunea și cere verificarea fluxului înainte de aplicări în lot.

Manualele Markdown, HTML și PDF sunt regenerate pentru fiecare build și salvate lângă executabil, cu versiunea în nume. HTML-ul se deschide în browser fără internet și are căutare locală. PDF-ul este pentru citire și tipărire. Versiunea, ediția, data buildului și commitul identifică documentele; manifestul build.json păstrează și sumele de verificare.

Ajutorul contextual și accesul la articole din AI sunt disponibile din această versiune. Reformularea opțională cu modelul și tururile vizuale automate rămân extensii viitoare. Instrucțiunile nu execută operații și nu înlocuiesc validările sau confirmările.

## 20. Ghiduri interactive comune aplicației și asistentului

Articolele de mai jos sunt generate din catalogul comun RO/EN. Aceleași instrucțiuni sunt afișate în panoul «Cum să faci» și în răspunsurile de utilizare ale asistentului. Identificatorii stabili permit citarea și deschiderea articolului.

{{ARTICOLE_AJUTOR}}

### Clarificari de operare pentru 0.88.15

Indicatorul zilnic numara persoanele cu inregistrari din lista filtrata; nu confirma necesarul de personal sau deficitul. Panoul AI cere confirmare inainte de schimbarea perioadei ori deschiderea concediilor daca exista modificari nesalvate. Ghidurile functioneaza fara model; configurarea motorului este pliabila. Pe ecrane inguste panoul AI se suprapune peste grila si se inchide din antet. Shift+F10 deschide actiunile celulei; Tab parcurge butoanele, Escape revine la celula. Cautarea in ghiduri afiseaza un rezultat curent; o citare veche ofera deschiderea ghidului versiunii actuale.

### Pachetul de livrare 0.89.1

Buildul include o arhiva ZIP cu executabilul, manualele, manifestul, LICENSE, instructiunile ACTUALIZARE.md, NOTE-VERSIUNE.md si CHECKLIST-ACCEPTANTA.md. SHA256SUMS.txt verifica fisierele incluse; fisierul .zip.sha256 de langa arhiva verifica descarcarea. Pachetul standard nu include baza, seiful, backupuri sau modele AI. Extrageti arhiva inainte de pornire si urmati instructiunile de actualizare. Checklistul este de completat de operator, nu o lista de teste deja executate.

### Raporteaza o problema sau trimite o parere

Din meniul lateral alege Raporteaza o problema. Selecteaza tipul, descrie situatia si adauga pasii de reproducere. Jurnalul de diagnoza este optional. Genereaza raportul si verifica textul inainte de a-l trimite; elimina orice date personale ramase. Salveaza fisierul pentru trimitere pe canalul organizatiei sau deschide formularul GitHub precompletat. Publicarea pe GitHub este confirmata de tine in browser si necesita contul tau; aplicatia nu trimite automat raportul. Pentru descrieri lungi, pastreaza fisierul complet, deoarece formularul GitHub poate contine text scurtat.

### Mesaje de eroare si date indisponibile

Panoul de eroare urmeaza limba interfetei. O data absenta sau invalida nu trebuie interpretata ca o data calendaristica valida; verifica inregistrarea sursa in Jurnal, Restaurare sau Planificare. Daca apare o eroare, pastreaza contextul si foloseste raportarea problemei, fara date personale.

### Introducere, import si vizualizare (0.89.4)

Dialogurile pastreaza focusul si cursorul la tastare. Butonul din antet permite marire/restaurare; S/M/L schimba densitatea grilelor Pontaj si Planificare (mod grila). Importul angajatilor recunoaste Zile_Concediu_An ca numar anual de zile si Norma_Ore_Zi din sablonul oficial CSV/XLSX. Acest camp nu creeaza cereri de concediu cu date calendaristice: acestea se inregistreaza in Concedii si absente.

### Categorii de schimburi si functii predefinite (0.90.0)

In Setari -> Schimburi, zona Categorii de schimburi permite adaugarea si editarea denumirilor RO/EN si a starii active. Apoi editeaza sau adauga programul si alege categoria lui. O categorie poate contine mai multe programe. LAB si CUSTOM sunt pastrate la actualizare. Precompletarea cu sabloane ramane pentru LAB; o categorie noua nu primeste automat programe.

In Setari (pagina generala), zona Functii permite predefinirea listei. La adaugarea sau editarea angajatului, alege Functie din selector; functia aleasa apare si sub nume in Pontaj. Functiile existente sunt preluate automat la migrare. O valoare text importata care nu are inca o asociere in nomenclator este pastrata la angajat; pentru standardizare, defineste functia si selecteaz-o explicit.

Dezactivarea nu sterge datele sau asocierile curente, dar impiedica noi asocieri. Redenumirea functiei actualizeaza eticheta curenta a angajatilor asociati, cu jurnal; asocierea este curenta si nu reprezinta istoricul contractual pe luni. Un conflict de revizie cere reincarcarea sectiunii inainte de o noua editare. Verifica inainte modificarile nesalvate: formularul listei este local sectiunii.

## Marca de garda si raportul detaliat

1. In Angajati, pastrati inregistrarea principala cu marca de program normal.
2. Adaugati sau editati inregistrarea cu marca distincta de garda. La Angajat principal selectati inregistrarea de program normal, identificata si prin marca. Nu se face asociere automata dupa nume.
3. In Planificare introduceti garzile pe marca de garda, in sectiunea GARZI.
4. Generati Pontaj detaliat cu totaluri. Marca principala si marca de garda apar pe randuri separate, grupate. Selectarea unei persoane include si marcile asociate, cu respectarea filtrelor de sectie/compartiment.

Planificarea GARZI activa completeaza numai intrarile lipsa din Pontaj. O inregistrare Pontaj existenta (inclusiv goala) si absentele au prioritate; publicarea nu dubleaza orele. Raportul nu modifica Pontajul. Orele si trecerea intre zile folosesc configuratia de calcul existenta; nu se deduc sporuri sau procente. Legatura marcilor este curenta, nu istoric contractual.

## Fereastra la pornire

Dimensiunea initiala se adapteaza zonei de lucru a monitorului principal si scalarii Windows. Pe ecrane mici se reduce si limita minima, astfel incat fereastra sa poata incapea. Puteti redimensiona sau maximiza manual. Adaptarea se recalculeaza la fiecare pornire.

## Grupuri profesionale, pozitii si contracte (0.92.0)

In Setari configurati Grupuri profesionale (ex. Medici, Asistenti), cu ordinea numerica de afisare; numerele mici apar primele. In Functii atribuiti grupul si bifati Conducere/sef pentru functiile de conducere. In Angajati selectati functia. Sectia si compartimentul delimiteaza echipa; in rapoartele filtrate sefii acelei echipe apar primii, urmati alfabetic de personal. La egalitate de nume se foloseste marca. Marcile normale si de garda ale aceleiasi persoane raman impreuna. In lipsa unei functii proprii, marca de garda foloseste functia principala pentru filtrare/incadrare, prin legatura explicita.

Tipuri de contract: norma globala, individuala, zilnica fixa, lunara sau fara norma prestabilita. Norma fixa se introduce in minute (420 = 7 ore, 300 = 5 ore). Garda este un tip predefinit fara norma. Selectati contractul fiecarei marci. Asocierea unei marci de garda in formular selecteaza contractul GARDA daca acesta este activ. Migrarea actualizeaza numai marcile deja legate explicit ca garzi; nu identifica persoane dupa nume. Contractele si functiile pot fi dezactivate fara stergerea istoricului de modificari.

Categoria de programe a contractului, daca este stabilita, are prioritate fata de categoria angajatului si limiteaza programele permise la salvarea Pontajului si Planificarii. Precompletarea dupa categorii nu completeaza norme lunare sau contracte fara norma cu ore zilnice presupuse. In Definire rapoarte puteti adauga coloanele Tip de contract si Norma contractuala. Coloana existenta Ore norm isi pastreaza semnificatia de ore lucrate in zile normale, nu devine norma contractuala.

Angajati, Pontaj, Planificare si Rapoarte au filtre pentru grup, functie, pozitie conducere/executie, contract si garda. Garda inseamna incadrarea marcii/contractului, nu existenta unei ture in luna aleasa. Butonul Reseteaza filtrele elimina doar aceste filtre profesionale. Filtrele de sectie/compartiment si perioada raman independente.

## Reguli conditionate si plata totala 200%

In Configurare pontaj, regulile generale raman baza. Adaugati o regula conditionata si alegeti profilul: grup/functie, contract, pozitie si conditie de garda. Lasarea unui filtru gol inseamna toate valorile. Regula este initial inactiva; activati-o dupa verificarea selectiei si salvati configuratia.

Exemplul confirmat pentru medici: alegeti grupul corespunzator, Garda, contractul relevant, coloana Ore garda, sursa Ore din programele selectate in weekend/sarbatori, codurile eligibile (ex. G24) si factorul total 2. O ora eligibila produce doua ore echivalente; factorul reprezinta plata TOTALA 200%, nu inca 200% adaugat bazei. Orele efectiv lucrate din T nu se dubleaza. Aceasta coloana este o echivalare in ore, nu o suma in lei; baza valorica nu este definita inca.

Pentru asistenti, selectati grupul lor si coloana ZNL cu sursa Ore ZNL si factorul 1. Evidenta arata orele de weekend/sarbatoare; nu scade aceste ore ca zile libere deja efectuate si nu deduce o conversie in zile. Conversia, acumularea dreptului si acordarea timpului liber necesita o regula confirmata separat.

Clasificarea weekend/sarbatoare se face pe fragmente calendaristice 00-24, inclusiv la trecerea peste miezul noptii, fara dubla numarare cand sarbatoarea cade in weekend. Se respecta optiunea existenta de atribuire intre zile/luni. Doua reguli conditionate active care se aplica aceleiasi coloane si aceleiasi persoane sunt refuzate la calcul; nu se alege arbitrar una dintre ele.

Regulile si incadrarile sunt curente. Modificarea lor poate schimba regenerarea rapoartelor vechi; pastrati exporturile validate. Nu este implementat un istoric contractual cu perioade de valabilitate.

In Rapoarte generati mai intai previzualizarea, apoi exportati PDF sau Excel. Excel-ul foloseste exact documentul previzualizat, cu aceleasi filtre, grupuri si randuri. Exportul clasic din grila Pontaj ramane registrul lunar complet.

Din 0.92.1, salvarea unei marci de garda care foloseste functia principala pastreaza mostenirea. O functie principala dezactivata nu impiedica editarea celorlalte date ale marcii de garda. Atribuirea unei functii proprii diferite ramane explicita.

## Precizari pentru norma datorata (implementare in curs)

CO, CEV, CFS si CM aprobate sunt reducerile confirmate ale normei. Formarea profesionala si delegatia nu sunt in aceasta lista. Zilele libere compensatorii nu sunt scazute automat, deoarece regula lor nu este stabilita.

Motorul de baza pentru contracte cu norma zilnica si calendar luni-vineri exclude zilele libere calendaristice si concediile aprobate fara scadere dubla. Norma zilnica absenta ramane nedefinita, distincta de o norma calculata de zero ore. Acest calcul nu este inca legat de coloanele raportului: semnificatia finala a T si a soldurilor L FIN/L Inceput, respectiv identificarea zilelor programate pentru ture sunt in clarificare. Raportul existent nu trebuie considerat modificat doar pe baza acestei precizari.

## Legenda coloanelor in rapoarte (0.92.3)

Raportul detaliat nu mai afiseaza coloana Ore norma. Eliminarea este numai de prezentare; nu modifica sursele motorului sau grila Pontaj.

Subsolul PDF explica fiecare coloana vizibila, pe rand separat: semnificatie, sursa de calcul activa, coduri selectate, factor si filtrul de aprobare, unde sunt aplicabile. Zilele lunii au o explicatie comuna. Regulile conditionate sunt prezentate cu profilurile carora li se aplica. Legenda se repeta pe fiecare pagina; exportul Excel din Rapoarte o include sub tabel.

O regula neconfigurata produce liniuta, nu zero. Legenda descrie configuratia executata, inclusiv limitarile existente: norma datorata, marcajul T, noua separare a orelor de garda si soldurile lunare nu sunt inca integrate in acest raport. Explicatiile nu confirma automat o regula de plata.

## Zoom pentru intreaga interfata

In antetul fiecarui ecran folositi minus pentru micsorare, plus pentru marire si butonul procentual pentru revenire la 100%. Sunt disponibile niveluri intre 75% si 200%; alegerea este pastrata local. Scurtaturi: Ctrl+- (micsorare), Ctrl++ (marire), Ctrl+0 (resetare), inclusiv in dialoguri. Zoomul global este independent de densitatea S/M/L din Pontaj si Planificare. PDF-urile si Excel-urile exportate isi pastreaza dimensiunile.

## Marcarea automata a marcilor subordonate

O marca de garda asociata explicit unei persoane principale este afisata ca Nume Prenume (Garda), cu marcaj tradus in limba aleasa. Marcajul apare in liste, selectii, Pontaj, Planificare, absente/Gantt, AI si rapoarte PDF/Excel. Numele si prenumele editabile nu se modifica; nu adaugati manual sufixul. Marca principala si persoanele cu acelasi nume, fara legatura, nu primesc marcajul. Un contract de garda independent nu este considerat automat marca subordonata.

La schimbarea sau eliminarea asocierii, afisarea curenta se actualizeaza dupa salvare. Fisierele deja exportate si instantaneele AI salvate pastreaza continutul istoric; o noua generare foloseste incadrarea curenta. In prezent, singurul subtip ierarhic disponibil este marca de garda.

## Program normal si calendarul zilelor lucratoare

Programul normal luni-vineri se aplica numai zilelor lucratoare: sambata, duminica si sarbatorile legale din calendar sunt excluse. Filtrul AI Zile lucratoare si precompletarea folosesc aceasta regula. Verificati calendarul de sarbatori pentru fiecare an vizat, inclusiv cand intervalul trece in anul urmator.

Weekend inseamna strict sambata-duminica; o sarbatoare de marti nu este inclusa automat intr-o cerere doar pentru weekend. Optiunea Toate zilele sau o data explicita permite planificarea distincta a unei garzi/ture in acea zi. Munca efectiva in weekend sau sarbatoare nu este stearsa: este clasificata separat conform regulilor configurate, pe fragmente calendaristice 00-24. Suprapunerea sarbatoare-weekend se numara o singura data.

Pentru contractul normal, norma exclude aceste zile. Pentru garzi si ture se folosesc incadrarea si regulile lor specifice; programul normal nu se aplica automat. Interpretarea 200%, ZNL si timpul liber compensatoriu sunt tratate prin regulile confirmate separat; aceasta corectie de calendar nu defineste noi procente sau conversii in zile. Regula compensatoriilor ramane neprecizata.

Schimbarea calendarului impune regenerarea previzualizarii AI. Exporturile deja salvate raman neschimbate; regenerarea rapoartelor foloseste calendarul curent.

## Concediu sau alt tip de absenta direct din Pontaj

1. Selectati celula sau intervalul de zile si persoane in grila Pontaj.
2. Deschideti Completare si alegeti Concediu / alt tip de absenta.
3. Alegeti unul dintre tipurile configurate (de exemplu CO sau CM), completati aprobarea, numarul de inregistrare, inlocuitorul si documentul, dupa caz.
4. Verificati selectia din formular, apoi Salvati. Se folosesc numai zilele lucratoare din selectie; weekendurile si sarbatorile sunt excluse.

Operatia inlocuieste programul din schita si salveaza imediat absenta. Nu mai este necesara salvarea schitei pentru aceste zile. Codul concediului este afisat in celula; aceeasi inregistrare apare in Concedii si absente. Puteti redeschide meniul pentru a schimba tipul unei absente existente. Aprobarea nu este acordata automat: pentru o absenta noua starea initiala este In asteptare.

Modificarile nesalvate din afara selectiei se pastreaza. Istoricul local Anuleaza/Refa se reinitializeaza dupa conversia salvata; corectiile ulterioare se fac prin formularul de absente. Programul planificat suprapus este anulat conform prioritatii concediilor; medicalul poate afecta si tura precedenta sau marcile asociate. Lunile inchise si pontajul efectiv inregistrat sunt protejate. O eroare de validare sau un conflict de revizie anuleaza intregul lot si pastreaza schita initiala.

Pentru revenire de la absenta la program de lucru, anulati absenta din evidenta Concedii si absente si apoi aplicati programul dorit in Pontaj; simpla aplicare a unui program nu suprascrie o absenta activa.


## Prioritatea concediilor si a medicalului

Orice concediu, inclusiv medical, se inregistreaza numai luni-vineri, fara sarbatorile din calendarul aplicatiei. Verificati calendarul pentru anul folosit. Formularul exclude automat weekendurile si sarbatorile; serviciul refuza si cererile directe care contin aceste zile.

Concediul inlocuieste programul din schita si planificarea acelei zile. Medicalul are prioritate fata de celelalte concedii si programe. Daca o tura inceputa in ziua precedenta intra in ziua medicalului, se anuleaza intreaga tura. Sunt anulate si programele marcilor asociate explicit aceleiasi persoane, inclusiv garzile. Absenta medicala nu se dubleaza pe aceste marci. O cerere noua nu suprascrie automat medicalul. Pentru corectarea explicita folositi Editeaza sau formularul din Pontaj pe absenta existenta; revizia inregistrarii este verificata.

Salvarea este atomica si jurnalizata. O luna inchisa, inclusiv luna precedenta afectata de o tura de noapte, blocheaza operatia intreaga. Pontajele efective cu ore de inceput/sfarsit nu se sterg automat; un conflict trebuie verificat de operator. Programele anulate nu se restaureaza automat la anularea concediului. Inregistrarile istorice din weekend/sarbatori nu sunt migrate automat; verificati-le separat.


## Corectarea si anularea concediilor, inclusiv medicale

In Concedii si absente, folositi Editeaza pe rand pentru schimbarea tipului si a detaliilor. Pentru mai multe zile, selectati randurile si alegeti editarea selectiei. Aceste actiuni permit inclusiv corectarea unui medical in alt tip de concediu.

Pentru eliminare, folositi Anuleaza ziua sau Anuleaza zilele selectate si confirmati. Absenta devine inactiva si nu mai blocheaza programul; istoricul ramane disponibil. Anularea medicalului nu reface automat programele pe care le-a anulat. Pentru revenire la lucru aplicati din nou programul corect.

O corectie facuta intre timp de alt operator produce un conflict: reincarcati datele inainte de a repeta. Lunile inchise trebuie redeschise prin fluxul autorizat. Prioritatea medicalului protejeaza fata de suprascrieri automate, fara sa impiedice corectarea explicita a unei inregistrari curente.


## Formatul orelor

Aplicatia afiseaza orele in format de 24 de ore: 00:00, 07:00, 19:00, 23:59. Nu se foloseste AM/PM, indiferent de limba interfetei sau a sistemului. Regula se aplica si jurnalului, auditului, restaurarii, blocarii conturilor si ultimei modificari din Pontaj. Datele raman localizate; fusul orar si duratele/totalurile de ore nu se modifica. Rapoartele folosesc deja formatul 24h.

La definirea unui schimb, introduceti ora de inceput ca HH:MM, de exemplu 07:00 sau 19:00; sunt acceptate ore de la 00:00 la 23:59.


## Tipuri de schimb si regimuri configurabile

In Setari > Schimburi sunt disponibile listele Tipuri de schimb si Regimuri de lucru. Puteti adauga definitii, edita denumirile in romana/engleza si dezactiva intrarile. La adaugare selectati comportamentul in aplicatie. Tipurile folosesc Normal, Tura, Noapte sau Garda; regimurile folosesc Normal, Weekend/sarbatoare sau Intrare speciala in zi libera.

Regimul General/anterior este convertit automat la Normal la actualizarea bazei. Tipul schimbului ramane distinct: un schimb Noapte sau Garda nu devine Normal prin aceasta conversie. Intervalele, duratele, identitatile si legaturile existente se pastreaza.

In Adauga/Editeaza program selectati tipul si regimul din listele configurate. Acelasi program este folosit in categoriile de program, contractele asociate acelor categorii, pontaj, planificare, precompletare si calcul. Denumirea personalizata nu determina singura comportamentul; se foloseste comportamentul selectat explicit. Noapte/Garda si regimurile speciale raman gestionate conform motorului existent. Configuratia nu stabileste automat procente noi de spor.

Puteti redenumi o definitie deja folosita. Pentru schimbarea comportamentului creati o definitie noua si alegeti-o explicit in program; comportamentul unei definitii deja folosite este protejat pentru a evita reclasificarea retroactiva implicita. Dezactivarea ascunde definitia din alegerile noi; programele existente isi pastreaza asocierea si pot fi editate. Listele se anuleaza prin dezactivare, pastrand referintele si istoricul.


## Schimbarea datei si reactivarea unei absente

In Concedii si absente, pe randul zilei alegeti Editeaza. Puteti schimba data; cele doua campuri de data raman sincronizate pentru editarea unei singure zile. La salvare vechea zi se anuleaza si noua zi se inregistreaza atomic. Daca noua zi este ocupata, nelucratoare, intr-o luna inchisa sau apare un conflict, operatia intreaga este refuzata. Pastrati Activ bifat pentru mutare.

Sterge ziua elimina ziua din evidenta activa prin anulare cu istoric. Bifati Arata si zilele anulate pentru a o vedea din nou si folositi Reactiveaza. In formular puteti debifa Activ pentru dezactivare sau bifa pentru reactivare. Reaparitia programelor anulate anterior de concediu nu este automata. In lista Tipuri de absenta sunt disponibile butoane Dezactiveaza/Reactiveaza, precum si starea Activ din editor. Dezactivarea unui tip nu anuleaza zilele existente.


## Filtrarea dupa starea activa

Listele cu stare au filtru Active / Inactive / Toate, implicit Active. Este disponibil pentru nomenclatoare (functii, grupuri, contracte, categorii de schimb, tipuri si regimuri), programe de lucru, categorii de program, zile si tipuri de absenta, conturi, modele de raport si reguli conditionate. Lista Angajati porneste pe activi.

In Pontaj filtrul controleaza angajatii afisati; alegeti Inactive sau Toate pentru istoricul persoanelor dezactivate. In Planificare exista un filtru pentru angajati si, in modul Lista, unul separat pentru pozitiile active/retrase. Filtrarea nu sterge datele si nu modifica totalurile salvate sau formulele rapoartelor. Un program, tip sau model dezactivat se regaseste prin filtrul Inactive; reactivati-l din editorul sau. Regula noua, creata initial inactiva, este afisata prin trecerea explicita la Toate.

In evidenta absentelor, Inactive inlocuieste vechiul filtru Arata si zilele anulate; Toate afiseaza ambele stari. Pentru reactivare folositi butonul de pe rand. Ecranele fara stare activ/inactiv (de exemplu jurnal, backupuri si calendarul sarbatorilor) isi pastreaza filtrele proprii; nu se introduc stari artificiale pentru aceste date. Optiunile din formularele de creare raman active, iar la editare este pastrata si valoarea inactiva deja asociata, pentru a nu pierde referinta.


## Stergere definitiva a configurarilor nefolosite - administrator

Autentificati-va ca administrator. In lista configurarii alegeti Stare > Inactive (sau Toate), apoi Sterge definitiv pe randul dorit si confirmati numele. Inregistrarile active trebuie mai intai dezactivate. Operatorii nu au buton, iar endpointul refuza si apelurile directe fara administrator autentificat.

Disponibil pentru programe, categorii de program, tipuri de absenta, functii, grupuri, tipuri de contract, categorii de schimb, tipuri/regimuri de lucru si modele de raport. Definitii de baza necesare compatibilitatii raman protejate. O configurare folosita de pontaj, planificare, concedii, angajati, reguli, sabloane sau operatii AI nu poate fi stearsa; mesajul indica sursa utilizarii. Modificati configurarea care o refera sau pastrati-o inactiva daca exista istoric.

Stergerea unei categorii neatribuite elimina numai lista ei proprie de programe permise; programele nu sunt sterse. Stergerea unui program nefolosit elimina numai metadatele sale de clasificare. Nu se elimina asocieri din alte configurari pentru a forta stergerea. Pontajele, concediile, planificarile, angajatii si conturile nu se sterg prin aceasta actiune. Jurnalul pastreaza snapshotul configurarii eliminate. Stergerea este definitiva; nu exista buton de restaurare pentru configurarea stearsa.

### Prioritatea Planificarii, concediilor si medicalului

La publicarea Planificarii, programele si orele manuale existente in schita Pontaj sunt inlocuite implicit; nu mai trebuie bifata o permisiune de inlocuire. Previzualizarea arata valoarea veche si cea noua. Daca datele se modifica intre previzualizare si confirmare, repetati previzualizarea.

Concediile si zilele libere au prioritate fata de program si anuleaza planificarea aferenta. Medicalul are prioritate maxima si anuleaza intreaga tura de noapte intersectata pe toate marcile asociate. Absentele se inregistreaza numai in zile lucratoare, fara weekend si sarbatori legale. Lunile inchise raman protejate.

#### Suprascrierea pontajelor efective

La publicarea Planificarii sunt inlocuite si pontajele efective ale aceleiasi marci si date. Previzualizarea indica numarul lor. Inregistrarile vechi sunt retrase din calcul si pastrate in audit; nu se sterg fizic. Concediile retrag pontajele efective intersectate, iar medicalul aplica regula tuturor marcilor asociate si intregii ture intersectate. O luna inchisa sau blocata impiedica operatia, inclusiv pentru o tura inceputa in luna precedenta. Anularea concediului nu restaureaza automat programul retras: programul dorit trebuie reaplicat explicit.

### Solduri lunare si initializarea pontajului

In Angajati > Editeaza completati Luna de initializare, L Inceput initial (ore cu semn) si CO efectuat anterior (zile). Soldul orelor este valabil la inceputul lunii alese. Zilele CO initiale reprezinta consumul din acelasi an, anterior acelei luni; nu le mai introduceti ca absente datate. Dreptul anual contractual ramane separat. Pentru corectie editati valorile; pentru anularea valorilor initiale introduceti zero. Modificarile care ar afecta luni inchise/blocate sunt refuzate.

Ore lucratoare = (zile luni-vineri fara sarbatori legale - zile CO/CEV/CFS/CM aprobate) x norma zilnica. CEV/CFS folosesc codurile atribuite in regulile de pontaj; celelalte tipuri speciale nu se scad automat. Formarea profesionala si delegatiile nu reduc norma prin simpla incadrare la Special.

Raportul detaliat afiseaza L Inceput imediat dupa Marca. L Fin = L Inceput + total ore efectiv lucrate - Ore lucratoare. L Inceput al lunii urmatoare este L Fin al lunii precedente, inclusiv la trecerea dintre ani. Sunt acceptate solduri negative. Daca exista ore neinterpretabile, soldul final ramane nedefinit pana la corectarea lor. Marcile de garda raman fara norma prestabilita si fara sold de norma calculat. Formulele sunt explicate si in subsolul raportului.

Fara initializare explicita, calculul incepe cu sold zero in prima luna cu date a marcii. Soldurile se recalculeaza din datele existente; pentru un istoric anterior folosirii aplicatiei, introduceti soldul la luna de referinta. Norma contractuala folosita este cea configurata curent; nu se presupun modificari contractuale istorice neinregistrate.

### Claritate si navigare in interfata (0.92.19)

Marcajul Garda apare in nume atat pentru contractele cu atribut de garda, cat si pentru marcile asociate unui angajat principal. Numele salvat ramane nemodificat. In lista angajatilor sunt distincte asocierea cu marca principala si contractul de garda independent. Un contract independent nu creeaza automat o legatura intre persoane.

Formularul angajatului este grupat in Identitate, Incadrare, Contract/program si Stare. Soldurile initiale se deschid separat; rezumatul arata valorile existente. Filtrele profesionale sunt intr-un panou retractabil; filtrele active raman afisate deasupra lui. Cautarea, perioada si starea activa/inactiva raman disponibile.

Tabelele au repere fixe pentru identificare si antete la derulare. In Angajati raman vizibile numele si marca; grilele isi pastreaza identificarea si totalul fix. Focalizati zona tabelului pentru derulare cu tastatura; Alt+Home revine la inceput, Alt+End merge la capatul din dreapta. La zoom mare reperele grilelor se compacteaza pentru a pastra accesul la date.

In sumarul Planificarii, De ce lipsesc programe explica excluderea prin stare, tip/regim sau categorie si ofera acces spre configurare/ajutor. Programul de noapte trebuie clasificat ca Noapte, nu doar denumit astfel. Etapele Nesalvat, Salvat in plan si Publicat in Pontaj sunt afisate textual; previzualizarea arata si numarul pontajelor efective retrase. Modificarile nesalvate sunt semnalate la navigare si inchidere.

Formularul de absente afiseaza zilele selectate, marcile afectate si prioritatea program/concediu/medical. La anulare sau respingere nu se retrag programe si nu se restaureaza automat orele retrase anterior. Weekendurile si sarbatorile legale sunt excluse din absente.

Setari contine o verificare ghidata a configurarilor, cu legaturi catre angajati, schimburi, categorii si regulile de pontaj. Aceasta semnaleaza puncte de verificat, nu declara automat gresite contractele independente sau lipsa unei categorii. Stergerea definitiva este o actiune rara pentru administratori, disponibila la configurari inactive in Alte actiuni. Referintele blocante ofera, unde este posibil, acces spre sursa folosirii.

Salvarea Excel verifica registrul si continutul scris temporar inainte de a inlocui destinatia; erorile nu mai sunt mascate printr-o descarcare alternativa. Investigarea blocajului de deschidere raportat de utilizator este amanata la cererea lui; nu se declara remediata cauza acelui blocaj.

## Activitatile lunii si configurarea asistata (0.92.20)

Butonul Activitatile lunii din antet deschide situatia datelor salvate pentru anul si luna alese: configurari lipsa, pozitii nepublicate, zile de absenta in asteptare si acces la verificarea Pontajului, raport si inchiderea lunii. Contoarele nu includ modificarile nesalvate si nu reprezinta o validare contabila. Legaturile pastreaza perioada aleasa; inchiderea se face in continuare din Pontaj, cu motiv si verificarile existente.

In Setari, Verificarea configurarii deschide Angajati cu criteriul selectat, de exemplu lipsa categoriei. Poti elimina criteriul sau reveni la verificari. Garda independenta este informatie de revizuit, nu eroare. In formularul angajatului si al categoriei, deschide Programe disponibile in Planificare pentru lista programelor permise si motivele de excludere. Lista explica eligibilitatea pentru Planificare, nu publica si nu completeaza automat programul normal.

## Revizuirea modificarilor si navigarea (0.92.20)

La adaugarea sau editarea concediilor/absentelor, completeaza formularul si apasa Previzualizeaza modificarile. Simularea foloseste aceeasi operatie ca salvarea, intr-o tranzactie anulata: arata schita, programul planificat, absenta si pontajul efectiv afectate. Pentru medical se includ marcile asociate si intreaga tura intersectata, dupa regulile existente. Confirmarea revalideaza datele; daca intre timp se schimba datele sau apare un blocaj, poate fi necesara o noua previzualizare. Modificarea formularului invalideaza comparatia. Apasa Confirma si salveaza pentru aplicare; inchiderea previzualizarii nu salveaza. Anularea unei absente nu restaureaza automat programele retrase.

Grilele arata marcile si intervalul selectat langa comenzi. La schimbarea lunii sau navigarea dintr-o grila cu modificari, poti alege Ramai aici, Renunta la modificari sau Salveaza si continua. O salvare esuata nu continua navigarea. Salvarea Planificarii ramane separata de publicarea in Pontaj. Protectia nativa la inchiderea ferestrei Windows ramane cea existenta. Alte formulare isi pastreaza propriile validari si confirmari de abandon.

Filtrele profesionale active au eticheta si valoare si se pot elimina individual din butonul lor; celelalte criterii si perioada raman neschimbate. In Rapoarte, daca persoana selectata nu mai corespunde filtrelor, selectia se goleste cu mesaj explicativ. Compact/Normal/Aerisit controleaza densitatea grilei; procentul din antet controleaza zoomul intregii aplicatii. Regiunile de tabel au denumiri specifice pentru accesul prin tastatura; Alt+Home/End deruleaza la capetele tabelului cand regiunea are focus. Aceste facilitati nu reprezinta certificare de accesibilitate cu cititor de ecran.

## Anexa finala a rapoartelor PDF (0.92.21)

Modul de calcul, semnificatia coloanelor si istoricul tehnic de generare/revizii apar dupa tabel, intr-o anexa care incepe pe pagina noua. Nu se mai repeta in subsolul fiecarei pagini de pontaj. Numerotarea continua pana la final. Daca explicatiile depasesc o pagina, anexa continua fara trunchiere; istoricul de generare este la sfarsit. Istoricul nominal al operatiunilor ramane in Jurnal modificari.


## Tipuri de absenta si introducere din Planificare

In **Concedii si absente > Categorii de absenta > Adauga categorie de absenta**, definiti codul, denumirea si natura. La **Efect asupra orelor** alegeti:

- **Fara ore lucrate; reduce obligatia lunara**: ziua aprobata nu produce ore lucrate si scade o norma zilnica din Ore lucratoare.
- **Ore lucrate din norma zilnica**: ziua aprobata produce norma zilnica din contract, fara reducerea obligatiei lunare (de exemplu deplasare).
- **Ore lucrate din programul planificat**: ziua aprobata produce durata programului salvat anterior pe acea zi, fara reducerea obligatiei lunare.

Ultimele doua optiuni sunt disponibile pentru natura Alta absenta / deplasare. Pentru odihna si medical nu se acorda ore lucrate. Apoi folositi **Adauga tip de absenta**, completati codul si denumirea si selectati categoria definita. Tipul preia automat regula categoriei; aceasta este afisata in formular.

Categoriile pot fi adaugate, redenumite, dezactivate si reactivate. Lista are filtrul Active/Inactive/Toate, implicit Active. O categorie inactiva nu se poate atribui unor tipuri noi; asocierile existente si istoricul se pastreaza. Denumirea este folosita in selector, filtre, sumar si raportul de absente. Natura (odihna/medical/alta) stabileste evidenta concediului contractual si prioritatea medicalului, independent de numele categoriei. Natura si regula unei categorii folosite in istoric nu se modifica retroactiv; pentru alta regula creati o categorie noua. Categoriile initiale pastreaza regulile existente; migrarea separa automat categoriile cu efecte diferite deja configurate. Orele creditate nu definesc plata sau sporul.

In **Planificare > Grila lunara**, selectati angajatul/angajatii si zilele, apoi **Concediu / absenta pe selectie**. Salvati sau anulati mai intai modificarile nesalvate ale planului. Formularul comun permite alegerea tipului, aprobarii si documentelor, apoi previzualizarea efectelor si confirmarea salvarii. Weekendurile si sarbatorile legale sunt excluse. Efectul medicalului ramane prioritar, inclusiv pe marcile asociate si turele de noapte intersectate.

Programul necesar creditarii trebuie salvat inaintea absentei. Daca lipseste programul ori norma, salvarea este refuzata, fara modificari partiale. Durata creditata se pastreaza chiar daca programul este anulat de absenta; schimbarea aprobarii sau reactivarea nu recalculeaza durata din setarile curente. Zilele in asteptare/respinse nu produc ore creditate. Anularea elimina efectul din totaluri si nu restaureaza programul vechi automat.


## Rapoarte flexibile: SQL, Markdown si coloane proprii

Din **Setari > Definire rapoarte**, adaugati sau duplicati un model. La **Datele raportului** alegeti **Lista garzilor - un rand pe garda**, **Totaluri lunare si coloane calculate** sau **Raport SQL personalizat**. Folositi **Campuri disponibile pentru raport** pentru schema si completati un singur SELECT, ori lasati SQL gol pentru modelul implicit.

Exemplu pentru angajat si data garzii pe fiecare rand:

```sql
SELECT angajat, marca, data, program, inceput, sfarsit, ore
FROM garzi
ORDER BY angajat, data
```

Datele sunt limitate la perioada si filtrele raportului. Parametrii :an, :luna, :angajat_id, :sectie_id, :compartiment_id, :de_la si :pana_la sunt disponibili; :pana_la este prima zi a lunii urmatoare, exclusiva. SQL foloseste vederi de raportare locale, fara conturi, parole sau documente atasate. Nu se pot executa scrieri, PRAGMA ori ATTACH.

**Adauga coloana calculata** defineste codul, denumirea, unitatea, precizia si o formula sau SELECT scalar. Formule acceptate: + - * /, paranteze, ABS, MIN, MAX si ROUND. Exemplu: ROUND(total_ore_lucrate - total_ore_lucratoare, 2). Coloanele se evalueaza in ordinea afisata. O valoare lipsa ramane nedefinita daca nu alegeti explicit rezultatul zero. Pentru SQL per angajat includeti angajat_id in SELECT; il puteti trece la **Coloane ascunse in tabel**.

Catalogul lunar contine zile/ore conform contractului, ore lucratoare/lucrate, total concedii selectabil pe tipuri si categorii, total ore garda si ore eligibile weekend/sarbatori. Numai concediile aprobate din zile lucratoare intra in totaluri. Tipurile de concediu nou definite produc automat campuri concediu_* in catalog. Totalul selectat nu schimba obligatia lunara: deplasarile creditate nu se scad. Garda fara norma nu primeste obligatii de ore inventate.

Alegeti **Tabel automat**, **Markdown** sau **Pozitionare libera**. Markdown suporta titluri, paragrafe, liste, evidentiere si tabele simple; {{tabel}} insereaza lista rezultata. Nu este HTML executabil. Editorul liber permite tragerea campurilor si coordonate exacte in mm. Campurile obisnuite folosesc prima inregistrare; bifati **Repeta elementul pentru fiecare rand de date** si configurati pasul randurilor pentru o structura proprie pe toate inregistrarile. {{tabel}} si campurile repetate sunt alternative, pentru a evita dublarea.

Apasati **Verifica datele si previzualizeaza PDF**, verificati rezultatul, apoi **Salveaza**. Orice modificare necesita reverificare. Modelul activ apare in **Rapoarte**; datele se citesc din nou la generare. Formulele apar optional in anexa finala, conform bifei din Continut; istoricul tehnic ramane la final. Exportul Excel pastreaza tabelul, nu coordonatele vizuale.

Ghidul tehnic din repo, docs/ghid-rapoarte-flexibile.md, contine campurile, exemplele si limitele. Maximum 80 coloane/10000 randuri, cu limite de timp si volum; un aspect care nu incape este refuzat, fara decupare tacita.

### Explicatii optionale in raport

In Setari > Definire rapoarte > Editeaza > Continut, bifa Afiseaza modul de calcul si semnificatia coloanelor include explicatiile in anexa finala PDF si in Excel. Debifati pentru un raport fara explicatii, formule sau SQL. Salvati modelul; la modelele flexibile refaceti previzualizarea inainte de salvare. Modelele noi pornesc debifate; duplicarea pastreaza alegerea, iar modelele existente pastreaza afisarea pana la debifare. Istoricul tehnic ramane separat la final.
