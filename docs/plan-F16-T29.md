# F16-T29 - Remedieri audit UI28

Responsabil: Codex. Stare: in lucru. Utilizatorul a autorizat toate recomandarile din auditul 0.92.18.

Ordine: export Excel si tabele derulante; identitate Garda/asocieri/programe excluse; salvare/publicare si efecte absente; gruparea formularelor, filtre, solduri; stari textuale, actiuni rare, lizibilitate, ajutor de configurare. Fiecare UI28-01 ... UI28-14 va avea rezultat si verificare explicite aici. Fara schimbari de formule de plata, asociere automata dupa nume sau atingerea datelor reale.

Verificare Excel: aplicatia Microsoft Excel este disponibila prin COM pe acest sistem. Se vor genera doar registre anonime si se va verifica deschiderea lor ascunsa, read-only, fara instalari. Separat se verifica salvarea completa si redeschiderea fisierului final.

Plan de verificari: teste Go relevante si suita completa, TypeScript/frontend, Edge RO/EN pe dialoguri, tabele la ambele capete ale derularii, zoom si teme; manual versionat si build. Nu se declara validare de accesibilitate asistiva fara o proba dedicata.

## Rezultate pe fiecare punct al auditului

| Punct | Rezultat | Verificare |
| --- | --- | --- |
| UI28-01 | Marcajul Garda unificat dupa contract sau asociere, inclusiv rapoarte, export si AI; numele salvat ramane intact | Test frontend si Go pentru contract independent, teste de ierarhie; corectata dublarea marcajului in randurile AI |
| UI28-02 | Contractul independent si marca principala asociata sunt explicate separat; medicalul afiseaza marcile efectiv asociate | Formulare/absente RO/EN; nu se asociaza dupa nume |
| UI28-03 | Etape textuale Nesalvat/Salvat in plan/Publicat in Pontaj; numar efective retrase in previzualizare; protectie navigare si inchidere pentru plan/sumar | TypeScript/Vite, suita Go si scenarii UI |
| UI28-04 | Fiecare program exclus are motiv: inactiv, clasificare, categorie inactiva sau nepermisa; acces configurare/ajutor | Test de eligibilitate si nomenclatoare RO/EN |
| UI28-05 | Formular grupat: identitate, incadrare, contract/program, solduri retractabile, stare; validarea deschide grupul care contine eroarea | Focus/cursor, formular maximizat si campuri de initializare RO/EN |
| UI28-06 | Filtre profesionale retractabile, selectiile active vizibile; perioada/cautarea/starea raman accesibile | Nomenclatoare si grile RO/EN |
| UI28-07 | Rezumat al zilelor si marcilor afectate de absenta, prioritate medical si tura intreaga; anularea/respingerile sunt explicate separat | Scenarii conversie absente RO/EN; regulile de calcul nu se schimba |
| UI28-08 | Solduri initiale separate, rezumat cu luna/ore/CO, exemplu de ore negative si norma curenta explicita | Probe formular RO/EN; calculul ramane cel confirmat anterior |
| UI28-09 | Stari textuale si indicii suplimentare; culorile nu reprezinta singura explicatie | Contrast-ecrane: 0 probleme pe ambele teme; contrast-grila: toate starile minimum 4.5:1 |
| UI28-10 | Stergere numai admin/inactive in Alte actiuni, confirmare pastrata, legatura spre sectiunea utilizarii blocante cand este disponibila | Scenarii admin/operator RO/EN; autorizarea backend este pastrata |
| UI28-11 | Indicii si etichete mai lizibile, repere compacte la zoom mare, butoane de dialog pastrate | Edge 1024/1366, zoom global 100/150/200%, densitati S/M/L; DPI Windows nu a fost modificat |
| UI28-12 | Verificare ghidata in Setari si legaturi intre configurari, explicatii pentru tip/regim/categorie si ajutor contextual | Compilare, control traduceri si scenarii nomenclatoare |
| UI28-13 | Antete, nume si marca fixate; totaluri/actiuni la capatul opus, acces la tabel prin tastatura, Alt+Home/End | Probe geometrice fara suprapunere pe Angajati/Pontaj/Planificare, ambele teme si capete ale derularii |
| UI28-14 | Investigatia blocajului de deschidere AMANATA explicit de utilizator. Se pastreaza salvarea verificata si raportarea erorilor, deja implementate | Inaintea amanarii, doua registre anonime cu 300 persoane s-au deschis in Excel COM in 379/85 ms. Blocajul utilizatorului nu a fost reprodus si nu este declarat rezolvat |

## Verificari si limite

Gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute. 24 teste frontend si TypeScript/Vite trecute; Edge focus/scroll/zoom si campuri RO/EN, nomenclatoare RO/EN, absente din Pontaj RO/EN si admin/operator RO/EN trecute. Contrastul verificat pe ambele teme. 3 teste manual, 7 livrare si control documentatie trecute. Nu s-a efectuat un audit cu cititor de ecran sau pe configuratia reala DPI a utilizatorului.

Primele probe au identificat selectoare de test legate de vechea pozitie a campurilor; au fost inlocuite cu nume stabile. Prima suita Go a identificat marcaj Garda dublat in randurile AI si un test pentru vechea regula de asociere exclusiva; codul a fost corectat, iar testul tehnic a fost aliniat cerintei de contract independent. Testele de referinta contabila nu au fost schimbate. Referinte noi de raport pentru versiunile incrementate, cu pastrarea celor vechi.

Aplicatie 0.92.19, manual 1.30. Rapoarte: pontaj 10, detaliat 11, individual 10, absente 7, garzi 7. Export complet Excel 1.2.0. Fara migrari sau dependente noi. Dependente verificate: Go stabil 1.27.3, npm fara actualizari, 62 actualizari tranzitive Go raman urmarite separat in F14-T21. Nicio baza reala accesata.

Pas urmator: build oficial si verificarea integritatii pachetului, apoi proba operator; Excel ramane amanat. Fara push/merge.

Livrare verificata: `build/bin/pontix-0.92.19-20261031-210831-9f92eaecbda146c1a8101c64c3510452`, executabil `pontix.exe`, aplicatie 0.92.19, commit c98109b, manual 1.30 in MD/HTML/PDF. Manifest complet; SHA256 executabil/manuale si ZIP/CRC verificate. Excel ramane amanat la cererea utilizatorului. Pas urmator: proba operator; fara push/merge.
