# Versiuni

Istoricul versiunilor aplicației, cea mai nouă prima ([D009](decizii/D009-versionare.md)). Fiecare commit care schimbă codul crește versiunea din `internal/versiune/versiune.go` și adaugă aici o secțiune `## X.Y.Z · AAAA-LL-ZZ`.

Pe lângă versiunea aplicației, fiecare raport și fiecare schimbare a bazei de date au versiunea lor; ele se trec în secțiunea versiunii aplicației care le aduce.

## 0.19.0 · 2026-10-25

Cereri de concediu semnate atașate, număr de înregistrare, persoană înlocuitoare și stare de aprobare pentru concedii și absențe (F07-T12, D024).

- Zilele de concediu și absență includ opțional numărul de înregistrare al cererii, persoana care înlocuiește angajatul absent (validată ca fiind diferită de angajatul în concediu), starea aprobării (`IN_ASTEPTARE`, `APROBATA`, `RESPINSA`), numele celui care a aprobat și un document semnat atașat (PDF/PNG/JPG până la 10 MB), stocat în tabela `absente_documente` din baza SQLite criptată și partajat la salvarea unui interval.
- Ecranul Concedii și absențe afișează și filtrează după starea aprobării, permite atașarea și descărcarea cererii semnate și oferă completarea ulterioară a cererii, înlocuitorului sau aprobării pentru o zi ori pentru un lot de zile selectate.
- Grila de pontaj afișează în explicația celulei și prin stil vizual starea aprobării, numărul de înregistrare și persoana înlocuitoare pentru zilele de absență active.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0012_cereri_concediu`.**

## 0.18.0 · 2026-10-25

Categorii în nomenclatorul de schimburi: `LAB` pentru precompletare și `CUSTOM` pentru restul (F07-T11, D023).

- Câmpul `categorie` (`LAB` / `CUSTOM`) este separat de `tip` și `regim` în `tipuri_schimb`.
- Codurile `G24`, `IN9`, `N12`, `N9`, `Z7`, `Z12`, `S7` intră în categoria de precompletare `LAB`; celelalte coduri existente (`7`, `12`, `S1`, `S3`) și schimburile create manual intră în `CUSTOM`, cu păstrarea intervalelor și personalizărilor.
- Ecranul Schimburi oferă coloană și filtru pentru categorie, alegerea categoriei la precompletare și în formularul de adăugare/editare; sugestiile din pontaj includ categoria.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0011_categorii_schimburi`.**

## 0.17.2 · 2026-10-25

Actualizare de dependențe (F07-T11, D007).

- `react-i18next` actualizat 17.1.8 → 17.1.9 (MIT).

Rapoarte: nicio schimbare. Baza de date: nicio schimbare.

## 0.17.1 · 2026-10-25

Corecție a detectării conflictelor la salvări rapide (F06-T07, D022).

- Revizie numerică per rând, incrementată în SQL, separată de timestampul de audit în schițe și absențe.
- Două modificări în același tic al ceasului au revizii distincte; o salvare pe baza reviziei vechi refuză întregul lot. Conținutul și momentele înregistrărilor existente se păstrează la migrare.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0010_revizii_concurenta`.** Include programele și IN9 din 0.17.0.

## 0.17.0 · 2026-10-25

Programele reale și intrarea specială în weekend/sărbătoare (F07-T10, D021).

- Z7 07–14, S7 15–22, N9 22–07; IN9 22–07 distinct, început în seara precedentă zilei libere. Pentru weekenduri/sărbători: Z12 07–19, N12 19–07 și G24 07–07.
- Regim separat de tipul activității, vizibil și editabil în nomenclator, filtrabil în listă și afișat în sugestiile din pontaj. GENERAL păstrează configurațiile vechi neclasificate.
- Completarea adaugă numai lipsurile. Codurile istorice și personalizările se păstrează; numai metadatele șabloanelor originale nemodificate N12/G24 sunt încadrate automat în regimul pentru zile libere.
- Fără alocare automată de ture, reguli noi de plată sau interpretarea golului 14–15 ca pauză scăzută.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0009_programe_lucru`.**

## 0.16.0 · 2026-10-25

Modul Concedii și absențe (F07-T04, D020).

- Tipuri de odihnă, medical și zile speciale; CO/CM inițiale și tipuri noi cu cod/denumiri RO/EN, editare și dezactivare. Codurile nu se confundă cu schimburile, iar categoria unui tip folosit se păstrează.
- Înregistrare explicită pe zile sau intervale cu selecție individuală, weekenduri/sărbători vizibile și selectare opțională luni–vineri fără sărbători. Editare, anulare individuală/în lot cu istoric și revizii pentru conflicte.
- Sumar anual pe angajat și categorie, cu dreptul contractual curent afișat separat. Pontajul arată absențele și le protejează împotriva suprascrierii; înregistrarea refuză zilele cu pontaj existent, inclusiv ture care traversează ziua.
- Evidența nu produce solduri legale, reportări, indemnizații ori calcule de plată. Baza standard Go de fusuri orare este inclusă pentru executabilul portabil.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0008_concedii_absente`.**

## 0.15.0 · 2026-10-25

Schimburi inițiale din PDF și intervalele precizate de Ioan Tuns (F07-T09, D019).

- 7: 07–14; 12: 14–02 a doua zi; N12: 19–07; G24: 07–07; S1: 06–14; S3: 22–06.
- Migrarea adaugă numai codurile lipsă și păstrează configurațiile existente, inclusiv cele inactive ori șterse logic. Pauza păstrează implicitul tehnic existent de 0 minute, editabil; nu se introduc reguli de plată.
- Ecranul Schimburi arată sfârșitul nominal și trecerea în ziua următoare, plus completare explicită a lipsurilor din setul inițial.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0007_schimburi_initiale`.**

## 0.14.0 · 2026-10-25

Departamente în cadrul secțiilor (F07-T08, D018).

- Setări: creare de departamente sub o secție și afișare ierarhică.
- Angajați: departament opțional, ales numai dintre cele ale secției; schimbarea secției golește selecția. Filtre pe secție și departament în listă și grila lunară.
- Migrarea păstrează secțiile și angajații existenți, fără atribuiri automate; validarea relației se face în serviciu și prin constrângeri/trigger în bază. Schimburile existente se păstrează.

Rapoarte: nicio schimbare. **Baza de date: migrarea `0006_departamente`.**

## 0.13.0 · 2026-10-25

Prima grilă lunară după structura modelului primit (F08-T01, D016).

- An/lună, filtru pe secție, zilele reale ale lunii, weekenduri și sărbători, antet și identificarea angajatului fixate la derulare.
- Introducere din tastatură, selecție cu Shift, copiere/lipire TSV și completare pe intervale; sugestii din schimburile active.
- Salvarea explicită și tranzacțională a schiței text, cu revizii per celulă, refuzul suprascrierii la conflict și jurnal. Avertizare pentru modificări nesalvate.
- Coloanele de total din model sunt prezente, necalculate. Codurile necunoscute nu se transformă în ore sau plăți.

Rapoarte: niciun PDF generat încă. **Baza de date: migrarea `0005_pontaj_ciorne`.**

## 0.12.1 · 2026-10-25

Corecția build-ului pentru păstrarea datelor portabile (F05-T15, D017).

- Eliminată curățarea `build/bin`, care poate conține folderul de date `karios/`.
- Parametru `-NumeExecutabil` pentru o ieșire alternativă; căile invalide și fișierele blocate sunt refuzate înainte de compilare.
- Proprietățile Windows sunt verificate pentru fișierul de ieșire ales. Test PowerShell de regresie pentru refuzuri.

Rapoarte: nicio schimbare. Baza de date: nicio migrare nouă.

## 0.12.0 · 2026-10-25

Nomenclatorul configurabil de schimburi (F07-T02, D015).

- Ecran Schimburi: cod, denumiri RO/EN, oră de început, durată și pauză în minute, tip și activ/inactiv.
- Creare/editare și dezactivare/reactivare, validări și jurnal tranzacțional. Coduri normalizate la majuscule, fără duplicate, inclusiv între schimburile inactive.
- Nu se introduc automat ore sau reguli de plată pentru N12 ori alte coduri.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.11.0 · 2026-10-25

Sărbători legale România 2026 și administrarea calendarului pe ani (F07-T03, D014).

- Migrarea de date include lista națională: 17 sărbători nominale în 16 date distincte; 1 iunie combină Ziua Copilului și a doua zi de Rusalii. Datele mobile sunt explicit ortodoxe. Înregistrările preexistente nu sunt suprascrise.
- Ecran Sărbători: selectare an, listare, creare/editare cu denumiri RO/EN, validarea datei reale și completarea repetabilă a datelor lipsă din 2026.
- Operațiile din interfață sunt tranzacționale și jurnalizate; nu se introduc reguli de plată.

Rapoarte: niciunul încă. **Baza de date: migrarea `0004_sarbatori_romania_2026`.**

## 0.10.0 · 2026-10-25

Primele fluxuri complete în interfața desktop (F07-T01, F07-T07).

- Ecran Angajați: căutare și filtre, adăugare, editare, dezactivare/reactivare, normă efectivă și concediu anual contractual. Formular cu validare și protecție la abandonarea modificărilor.
- Setări: norma întreagă globală și adăugarea secțiilor. Norma parțială rămâne individuală.
- Date demo fictive: 12 angajați și 3 secții; încărcare explicită, tranzacțională, numai în lipsa angajaților și pontajelor. Fără modificarea normei globale.
- Serviciu Go cu operații tranzacționale și jurnal pentru aceste fluxuri; erori structurate și texte RO/EN din cataloagele comune.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.9.1 · 2026-10-25

Acces la date generat cu sqlc v1.31.1 (F06-T03), pentru angajați, secții, setări și jurnal. Configurația folosește direct migrările existente; `genereaza-date.ps1` fixează versiunea generatorului fără instalare globală. Nu adaugă o dependență în executabil.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.9.0 · 2026-10-25

Schema datelor de bază cerute pentru angajați și sărbători (F07-T06, D012).

- Angajați: nume și prenume separate, marcă, contract `INTREAGA` implicit sau `PARTIALA`, norma parțială individuală în minute și zile anuale de concediu conform contractului (`NULL` până la completare).
- Norma întreagă globală: `setari.norma_intreaga_minute`, implicit `420` (7 ore/zi). Vederea `angajati_norme` expune `norma_efectiva_minute`: globală pentru norma întreagă, individuală pentru norma parțială.
- Datele existente se păstrează: numele vechi nu se desparte automat, norma veche rămâne în `norma_minute_anterioara`; noile informații contractuale necesită verificare/completare manuală. Norma veche nu intră în noul calcul al normei curente.
- Sărbători definite pe dată completă, cu `an` derivat și index pentru filtrare anuală; nu se introduc automat sărbători legale.
- Constrângeri pentru tipul contractului, norma parțială, concediu nenegativ și norma globală pozitivă; cheia normei globale nu poate fi eliminată sau redenumită.

Această versiune pregătește baza de date; formularele, accesul la date prin sqlc și istoricul contractual rămân în taskuri separate.

Rapoarte: niciunul încă. **Baza de date: migrările `0002_angajati_contracte` și `0003_sarbatori_an`.**

## 0.8.1 · 2026-10-25

Actualizarea dependențelor pentru F07-T06 (D007): SQLite WASM `v6.3.35304` → `v6.3.35401` în aplicație și prototip; `react-i18next` `17.1.7` → `17.1.8`. Versiuni fixate exact în manifestele și fișierele de blocare existente.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.8.0 · 2026-10-25

Backup automat înaintea actualizărilor bazei de date (F05-T06, D002 punctul 3).

- Înainte de orice migrare pe o bază care are deja date, aplicația face o copie criptată în `karios/backup/pontix-AAAA-LL-ZZ-HHMMSS-vX.Y.Z.db` (`bd.Backup`, cu `VACUUM INTO`).
- Dacă backup-ul eșuează, migrarea nu pornește și aplicația arată eroarea `stocare.backup_esuat`.
- Se păstrează cele mai noi **7** backup-uri automate; alte fișiere din `backup/` nu se ating.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.7.0 · 2026-10-25

Aplicația își deschide datele criptate la pornire (F05-T12, D011).

- Cheia aplicației stă în cod (decizia D011): datele din `karios/` nu se pot citi fără executabil.
- Pachetul `internal/stocare`: la prima pornire creează `karios/seif` și `karios/pontix.db`; la următoarele le deschide și aplică migrările lipsă.
- Dacă baza există, dar seiful lipsește, aplicația nu creează un seif nou (ar face baza ilizibilă), ci afișează eroarea `seif.deteriorat`.
- Baza se închide la ieșirea din aplicație.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă (prima pornire aplică `0001_schema_initiala`).

## 0.6.0 · 2026-10-25

Seiful cheii de date (F05-T10, D008 punctul 2).

- Pachetul `internal/seif`: generează cheia de date și o păstrează în `karios/seif`, în copii criptate cu AES-256-GCM: copia „A” cu cheia aplicației și copii deschise cu o parolă sau cu răspunsuri (Argon2id).
- Cheia aplicației se injectează la build (`seif.CheieAplicatie`); fără ea, seiful refuză să funcționeze.
- Textele noi: `seif.fara_cheie_aplicatie`, `seif.deteriorat`, `seif.secret_gresit`.
- Nu e încă folosit de aplicație; integrarea vine cu F05-T12.

Rapoarte: niciunul încă. Baza de date: nicio migrare nouă.

## 0.5.0 · 2026-10-25

Baza de date criptată și schema inițială (F06-T02).

- Pachetul `internal/bd`: `Deschide` (SQLite criptat AES-256-XTS, modul WAL, cheile străine active, cheia verificată la deschidere) și `Migreaza` (goose).
- Dependențe noi: `ncruces/go-sqlite3` v0.35.6 și `pressly/goose/v3` v3.29.0 (D008, D010).
- Textele noi: `bd.nu_se_deschide`, `bd.cheie_invalida`.
- Pachetul nu e încă folosit de aplicație; legătura vine cu seiful cheii (F05-T10).

Rapoarte: niciunul încă. **Baza de date: migrarea `0001_schema_initiala`** — tabelele `setari`, `sectii`, `angajati`, `tipuri_schimb`, `coduri_absenta`, `sarbatori`, `categorii_ore`, `seturi_reguli`, `reguli_interval`, `pontaj`, `jurnal_modificari`.

## 0.4.0 · 2026-10-24

Mesaj clar când folderul de date nu permite scrierea (F05-T05, D002 punctul 5).

- La pornire, aplicația verifică faptul că poate scrie în `karios/`. Dacă nu poate, afișează o fereastră Windows cu un mesaj tradus, în limba Windows, care spune unde să fie mutată aplicația.
- `i18n.Eroare` e acum tipul comun de eroare traductibilă; `calcul.Eroare` e un alias.
- `i18n.LimbaSistem`: limba interfeței Windows, dacă are catalog.
- Textele noi: `aplicatie.eroare_pornire`, `portabil.fara_scriere`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.3 · 2026-10-24

Scriptul de build (F05-T13).

- `build.ps1` construiește `pontix.exe` cu versiunea în proprietățile fișierului și commit-ul curent (cu sufixul `-modificat` dacă există fișiere necomise).
- `wails.json` are secțiunea `info` (nume, versiune, copyright); testul `TestWailsJSON` verifică faptul că versiunea e aceeași ca în `versiune.go`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.2 · 2026-10-24

Prototipul migrărilor (F06-T01). Aplicația nu se schimbă.

- `prototipuri/sqlite-criptat/migrari_test.go`: goose v3.29.0 pe baza criptată; rezultatele sunt în D010.
- Hook-ul `pre-commit` protejează doar migrările existente în `main`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.1 · 2026-10-24

Prototipul bazei de date criptate (F05-T09). Aplicația nu se schimbă.

- `prototipuri/sqlite-criptat/`, modul Go separat: `ncruces/go-sqlite3` v0.35.6 cu VFS-urile `adiantum` și `xts`. Toate verificările din D008 trec; rezultatele sunt în D008.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.3.0 · 2026-10-24

Datele lângă executabil (F05-T03, D002).

- Pachetul `internal/portabil` creează la pornire folderul `karios/` lângă `.exe`, cu `webview/` și `backup/`.
- Cache-ul WebView2 stă în `karios/webview/` (`WebviewUserDataPath`). Verificat: aplicația nu mai creează nimic în `%APPDATA%` sau `%LOCALAPPDATA%`.

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.2.0 · 2026-10-24

Scheletul aplicației desktop (F05-T02).

- Aplicația Wails v2.17.0 cu interfața React 19.3.1 + TypeScript 7.0.2, construită cu Vite 8.3.5. `pontix.exe` are 11,5 MB.
- Interfața ia textele din cataloagele Go prin i18next 26.5.10 / react-i18next 17.1.7 (F03-T06, parțial); funcția nouă `i18n.Catalog`.
- Fereastra arată numele, descrierea și versiunea aplicației; commit-ul se injectează la build.
- Textele noi: `aplicatie.titlu`, `aplicatie.subtitlu`, `aplicatie.versiune`.
- Dependențele Go actualizate la ultimele versiuni stabile (D007); pachetele npm fixate exact (`save-exact`).

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.

## 0.1.0 · 2026-10-24

Prima versiune numerotată. Cuprinde tot codul de până acum.

- Motorul de calcul: `Segmenteaza`, `Valideaza`, calendarul cu sărbători (F04-T01, F04-T02).
- Traducerile: cataloagele RO și EN, folderul de traduceri, `i18n.Text`, erori traductibile (F03-T01 … T04).
- Go 1.27.3 (F02-T09).
- Versionarea: pachetul `internal/versiune`, verificarea la commit (F02-T10).

Rapoarte: niciunul încă. Baza de date: nicio migrare încă.
