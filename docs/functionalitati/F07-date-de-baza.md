# F07 · Date de bază

Etapă: V1 · Stare: în lucru · Depinde de: F03, F06

## Scop

Utilizatorul poate gestiona angajații și listele de care depinde pontajul.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F07-T38 | Efecte configurabile ale absentelor si introducere din Planificare | gata | Codex | 0.92.23, migrare 0029; surse norma/program, instantaneu si audit, aprobari si reducerea obligatiei, formular comun din grila. [Rezultate](../plan-F07-T38.md), D080 propusa. Pas urmator: proba operator, fara push/merge. |
| F07-T37 | Stergere definitiva configurari inactive/nefolosite, numai administrator | gata | Codex | 0.92.15: domeniu confirmat explicit; stergere admin stricta, inactive/nefolosite, verificari referinte/revizii, audit si rollback. [Specificatie](../stergere-configurari-admin.md), [rezultate](../rezultate-F07-T37.md), [D075](../decizii/D075-stergere-configurari-admin.md). Pas urmator: build oficial si proba operator; fara push/merge. |
| F07-T36 | Modificare data absenta, anulare si reactivare vizibile | gata | Codex | Mutare atomica, stergere/anulare si reactivare, stare si tipuri editabile. [Rezultate](../rezultate-F07-T36.md); pas urmator: proba operator. |
| F07-T35 | Tipuri si regimuri editabile, eliminare General/anterior si legaturi coerente | gata | Codex | Nomenclatoare si migrare 0027, referinte stabile si comportament explicit. [Rezultate](../rezultate-F07-T35.md); pas urmator: proba operator. |
| F07-T34 | Corectare explicita si anulare concedii, inclusiv medicale | gata | Codex | Corectie cu revizie curenta, anulare cu istoric. [Rezultate](../rezultate-F07-T34.md); pas urmator: proba operator. |
| F07-T33 | Absente numai in zile lucratoare si prioritate medical > alte concedii > program | gata | Codex | Medicalul anuleaza intreaga tura si programele marcilor asociate, cu jurnal si protectie la reintroducere. [Rezultate](../rezultate-F07-T33.md); pas urmator: proba operator. |
| F07-T32 | Concedii si schimbarea tipului zilei din meniul Completare Pontaj | gata | Codex | Selectie exacta, formular comun de absente si conversie atomica a schitei de lucru. [Rezultate](../rezultate-F07-T32.md); pas urmator: proba operator. |
| F07-T31 | Program normal L-V exclude sarbatorile legale din calendar | gata | Codex | Corectie filtre AI individual/grup; weekend/sarbatoare separate, fara reguli noi de plata. [Rezultate](../rezultate-F07-T31.md); pas urmator: proba operator. |
| F07-T30 | Marcaj automat al numelor pentru marcile subordonate in toate contextele | gata | Codex | Legatura explicita principala/garda, fara modificarea numelui stocat. Liste, selectii, pontaj, planificare, absente, AI si exporturi. [Rezultate](../rezultate-F07-T30.md); pas urmator: proba operator. |
| F07-T29 | Inlocuire denumire Departament cu Compartiment in interfata, ajutor, rapoarte si manual | gata | Codex | Cerinta utilizatorului; cheile si datele existente raman compatibile. [Rezultate](../rezultate-F07-T29.md); pas urmator: proba operator. |
| F07-T28 | Clasificare explicita absente remunerate/neplatite si coloane CEV/CFS/FP/CM/Deleg/Ore nem | in lucru | Codex | Confirmate zilele programate, numai absente aprobate si introducere explicita ore/minute nemotivate. Mai necesita confirmarea sursei programului (ture/normal). Definitii si optiuni lipsa la F01. |
| F07-T27 | Norma lunara datorata, T ture si solduri L FIN/L Inceput | in lucru | Codex | Norma contractuala, T=ture, L FIN=efectiv-minus-norma, report luna urmatoare in L Inceput dupa Marca. Clarificari privind valoarea T, solduri negative si absente. Fara schimbarea formulelor pana la raspuns. |
| F07-T25 | Grupuri profesionale, pozitii, contracte configurabile, filtre si reguli conditionate | gata | Codex | Cerinte si etape in plan-F07-T25.md; pastrare date, fara formule monetare presupuse. |
| F07-T26 | Baza valorica in lei si conversia ZNL in timp liber acordat | blocat | Codex | Lipsesc baza de calcul a orei si formula de conversie/acordare. Intrebari la F01. Nu este blocata configurarea si echivalarea in ore confirmata pentru F07-T25. |
| F07-T21 | Sincronizare atomică Pontaj → Planificare pentru programele speciale | gata | Codex | Versiunea 0.88.9, salvare manuală și reparare explicită a celulelor speciale deja salvate. Testele Go complete, vet, 13 teste frontend și build Windows trecute. Integrarea în main autorizată explicit de utilizator la 30 octombrie; fără push. D051 rămâne propusă. Pas următor: probă vizuală cu operatorul și actualizarea manualului pentru noua versiune. |
| F07-T14 | Categorii de program atribuite angajaților și precompletare cu protecția absențelor | gata | Codex | 0.29.0 pe F07-T14-categorii-program-angajati, din 7880029. D036 acceptată: buton L–V vizibil, programe permise + implicit, gărzile explicite. Migrare 0016, șase categorii editabile, atribuire individuală și filtru în pontaj; previzualizare pentru rânduri vizibile/selectate cu absențe, calendar, norme și suprapuneri, aplicare în schiță/Undo/Redo/salvare explicită. TypeScript/Vite, gofmt, vet, toate testele Go și Edge trecute. Referințele rapoartelor păstrate. [Ghid](../ghid-categorii-program.md), [rezultate](../rezultate-F07-T14.md). Pas următor: probă cu operator; F04-T09 pentru anomalia preexistentă a calendarului în totaluri, I055 pentru atribuire în lot. |
| F07-T01 | Administrarea angajaților: nume, prenume, marcă, tip contract, normă parțială, concediu anual, funcție, secție, activ | gata | Codex | 2026-10-25, 0.10.0: interfață, metode Wails și serviciu Go peste sqlc; listare/căutare/filtre, creare/editare, dezactivare/reactivare, secții și norma globală. Teste Go pe baza criptată, UI în Edge cu punte simulată și build Windows trecute; [ghid](../ghid-primele-functionalitati.md). Urmează F07-T03 și prototipul F08-T01; [D013](../decizii/D013-primele-functionalitati-demo.md) rămâne propusă |
| F07-T02 | Tipuri de schimb și codurile lor scurte (`7`, `12`, `N12`, `G24`) | gata | Codex | 2026-10-25, 0.12.0: ecran de listare/creare/editare/dezactivare, denumiri RO/EN și intervale explicite. Teste Go și UI Edge, fără coduri sau ore precompletate; [D015](../decizii/D015-nomenclator-schimburi.md). Urmează grila F08 |
| F07-T03 | Administrarea sărbătorilor pentru fiecare an | gata | Codex | 2026-10-25, 0.11.0: set național 2026 (17 sărbători în 16 date), migrarea de date 0004, listare pe an, creare/editare RO/EN și completare fără suprascriere. Teste Go și UI Edge; [D014](../decizii/D014-sarbatori-2026.md). Pas următor: nomenclator schimburi și grila F08 |
| F07-T04 | Modul concedii și absențe: odihnă, medical, tipuri speciale configurabile | gata | Codex | 2026-10-25, 0.16.0: nomenclator CO/CM și tipuri speciale configurabile, înregistrare explicită pe zile/intervale, editare/anulare cu revizii și jurnal, sumar anual și protecție în pontaj. Migrarea 0008, păstrarea configurațiilor și a dreptului contractual, conflictele și turele traversând ziua verificate pe SQLite criptat. Go/vet, TypeScript/Vite, UI Edge și build Windows trecute. Urmează regulile de consum/sold/plată din F01 și I010; [D020](../decizii/D020-concedii-absente.md) rămâne propusă |
| F07-T05 | Categorii de ore și coloanele centralizatorului | de făcut | | O singură sursă pentru legătura categorie → coloană |
| F07-T06 | Schema angajaților, norma întreagă globală și sărbătorile pe ani | gata | Codex | 2026-10-25, versiunea 0.9.0, migrările 0002 și 0003; păstrarea datelor vechi, constrângeri, normă efectivă și backup înainte de migrare verificate prin teste. Build Windows verificat. Pas următor: F06-T03, F07-T01 și F07-T03; [D012](../decizii/D012-date-angajati-norme-sarbatori.md) rămâne propusă pentru reprezentarea tehnică |
| F07-T07 | Set de date demo fictive și încărcare explicită, tranzacțională | gata | Codex | 2026-10-25, 0.10.0: 12 angajați și 3 secții, JSON inclus în aplicație; confirmare în UI și refuz în backend dacă există angajați/pontaje. Testate anularea integrală la eroare și păstrarea normei globale. Datele demo nu sunt încărcate în baza operațională de agent |
| F07-T08 | Departamente multiple în fiecare secție, atribuire angajați și filtre | gata | Codex | 2026-10-25, 0.14.0: migrarea 0006, creare în Setări, atribuire filtrată după secție, filtre în Angajați/Pontaj. Verificate integritatea relației și păstrarea configurațiilor vechi. Urmează F07-T09; [D018](../decizii/D018-sectii-departamente.md) rămâne propusă |
| F07-T09 | Setările inițiale ale schimburilor din PDF | gata | Codex | 2026-10-25, 0.15.0: șase coduri inițiale, migrarea 0007, completare explicită și sfârșit nominal în UI. 7 = 07–14 și 12 = 14–02 a doua zi, conform precizărilor utilizatorului. Testate păstrarea configurațiilor/inactivelor/ștergerilor logice, redenumirea, repetabilitatea și anularea tranzacției la eroare. Go, vet, TypeScript/Vite, Edge și build Windows trecute. Urmează F07-T04/F07-T05 după specificațiile F01; [D019](../decizii/D019-schimburi-initiale-pdf.md) rămâne propusă pentru reprezentarea tehnică |
| F07-T10 | Programe de lucru și schimb special de intrare | gata | Codex | 2026-10-25, implementare 0.17.0, livrată în 0.17.1 cu F06-T07: Z7/S7/N9/IN9/Z12/N12/G24 și regim separat, formular/filtre/sugestii în pontaj. Migrarea 0009 păstrează coduri, personalizări, ștergeri logice și referințe. Go, vet, TypeScript/Vite, UI Edge și build Windows trecute. Urmează regulile de plată și rotație din F01; [D021](../decizii/D021-programe-lucru.md) rămâne propusă pentru reprezentarea tehnică |
| F07-T11 | Categorii în nomenclatorul de schimburi (`LAB` pentru precompletare și `CUSTOM`) | gata | Antigravity | 2026-10-25, 0.18.0: migrarea 0011, câmp `categorie` separat de `tip` și `regim`; setul `Z7`, `S7`, `N9`, `IN9`, `Z12`, `N12`, `G24` încadrat în `LAB`, restul în `CUSTOM`, fără suprascrierea intervalelor sau personalizărilor. Filtru în tabel, alegere în formular și la precompletare, sugestii în pontaj. Teste Go pe SQLite criptat, `go vet`, TypeScript/Vite și UI Edge trecute. Urmează regulile de plată și coloanele din F01/F07-T05; [D023](../decizii/D023-categorii-schimburi.md) rămâne propusă pentru reprezentarea tehnică |
| F07-T12 | Cereri de concediu semnate atașate, număr de înregistrare, înlocuitor și stare de aprobare | gata | Antigravity | 2026-10-25, 0.19.0: migrarea 0012, tabel `absente_documente` (BLOB criptat în SQLite pentru PDF/PNG/JPEG până la 10 MB) și câmpuri pe `absente_zile` (`nr_inregistrare`, `inlocuitor_id`, `stare_aprobare`, `aprobat_de`, `document_id`), opționale la salvare și editabile ulterior pe zi sau lot de zile selectate; evidențierea stării în pontaj. Teste Go pe SQLite criptat, `go vet`, TypeScript/Vite și UI Edge trecute. Urmează regulile de consum/sold/plată din F01 și I010; [D024](../decizii/D024-cereri-concediu-aprobare.md) rămâne propusă pentru reprezentarea tehnică |
| F07-T13 | Verificarea suprapunerilor înlocuitorului, evidențierea înlocuirii în pontaj, diferența aritmetică în sumarul de concedii și previzualizarea cererii semnate | gata | Antigravity | 2026-10-26, 0.23.0: refuz tranzacțional și în formular când înlocuitorul are absență activă nerespinsă în aceeași zi (`absente.inlocuitor_absent`) sau când titularul este deja desemnat înlocuitor în aceeași zi (`absente.titular_este_inlocuitor`), avertizare la dublă înlocuire, evidențierea înlocuirii pe rândul înlocuitorului în `Pontaj.tsx`, diferența aritmetică `Rămas` și zilele în așteptare în sumarul anual, plus previzualizare modală a cererii semnate. Teste Go, `go vet`, TypeScript/Vite și UI Edge trecute; [D030](../decizii/D030-suprapuneri-inlocuitor-previzualizare-cereri.md) rămâne propusă pentru reprezentarea tehnică |

## Schema F07-T06

| Informație | Reprezentare |
| --- | --- |
| Nume și prenume | `angajati.nume`, `angajati.prenume`; numele vechi rămâne intact, prenumele se completează manual |
| Număr de marcă | `angajati.marca`, text pentru păstrarea zerourilor inițiale, unic între angajații neșterși |
| Tip contract | `angajati.tip_contract`: `INTREAGA` implicit / `PARTIALA` |
| Norma întreagă | `setari.valoare` pentru cheia `norma_intreaga_minute`, implicit `420` minute = 7 ore/zi |
| Norma parțială | `angajati.norma_partiala_minute`, obligatorie numai pentru `PARTIALA`; minute întregi pozitive, cel mult 1440/zi |
| Norma efectivă curentă | `angajati_norme.norma_efectiva_minute`; schimbarea globală se vede imediat pentru contractele întregi |
| Concediu anual contractual | `angajati.concediu_zile_an`, întreg nenegativ sau `NULL` dacă nu a fost completat; nu reprezintă soldul disponibil |
| Sărbători pe an | `sarbatori.data` și `denumire` JSON; `an` derivat automat din dată, index `(an, data)` |

Exemple pentru viitorul cod de acces la date (nu sunt încă expuse prin interfață):

```sql
SELECT id, nume, prenume, marca, tip_contract, norma_efectiva_minute, concediu_zile_an
FROM angajati_norme WHERE sters_la IS NULL ORDER BY nume, prenume;

UPDATE setari SET valoare = '420', modificat_la = '2026-10-25T10:00:00Z'
WHERE cheie = 'norma_intreaga_minute';

SELECT data, denumire FROM sarbatori WHERE an = 2027 ORDER BY data;
```

Compatibilitate: `norma_minute_anterioara` păstrează vechea valoare exclusiv pentru verificare; este `NULL` pentru angajații noi și nu se folosește în noile calcule. Tipul `INTREAGA` al înregistrărilor migrate este valoarea implicită, nu o clasificare contractuală dedusă. Înaintea folosirii operative se verifică manual tipul, prenumele, norma parțială și concediul.

Pasul următor: F06-T03 pentru accesul la date și F07-T01/F07-T03 pentru formulare; istoricul contractual este înregistrat ca I010, înainte de F12. Nu se modifică motorul de calcul sau regulile de plată în acest task.

| F07-T24 | Categorii configurabile de schimburi si nomenclator functii asociate angajatilor/Pontajului | gata | Codex | Pastrare LAB/CUSTOM si functii existente, migrare noua, CRUD cu dezactivare si revizie, selectare functie si afisare in Pontaj. |

F07-T24: verificari finale gofmt/vet/suita Go, 18 teste frontend, UI nomenclatoare RO/EN, focus si ajutor trecute. Pas urmator: livrare cumulativa cu F09-T11, fara push sau merge.



Predare F07-T25: [rezultate si limite](../rezultate-F07-T25.md), toate verificarile automate trecute. Build oficial dupa commit, apoi proba operator si review. F07-T26 ramane blocat pentru formulele neprecizate.

F07-T27: calculul pur al normei zilnice pe calendar L-V implementat si testat (fara dubla scadere, numai CO/CEV/CFS/CM aprobate, fara libere compensatorii presupuse). Integrarea in raport si restul rubricilor raman in lucru, pana la clarificarile T/solduri si zile programate pentru ture.

Predare calcul F07-T27: gofmt/vet/suita Go completa si testele manualului trecute. Codul de baza este documentat in D063; nu s-a generat un executabil nou si nu s-a schimbat raportul. Pas urmator: clasificarea configurabila a absentelor (F07-T28), apoi conectarea normei si soldurilor dupa clarificarea T si a zilelor programate pentru ture.

Continuare F07-T27 (2026-10-31): utilizatorul cere Ore lucratoare = zile lucratoare din luna x norma zilnica, L Fin ca diferenta fata de orele efective, L Inceput dupa Marca din L Fin anterior, initializare ore si CO efectuat anterior. Responsabil Codex, in lucru pe F07-T27-solduri-initializare. Sunt cerute explicit clarificari pentru schimbarea regulii de scadere a concediilor, sensul diferentei si includerea soldului initial in L Fin. [Plan](../plan-F07-T27-solduri.md).

Predare F07-T27, cerinta solduri: implementata in 0.92.18, raport detaliat v10 si manual 1.29; [rezultate si limite](../plan-F07-T27-solduri.md), [D077](../decizii/D077-solduri-lunare-initializare.md). Formula confirmata L Fin=L Inceput+efectiv-ore lucratoare. T=ture ramane separat, neimplementat in acest lot. Urmeaza build si proba operator, fara push/merge.
