# Versiuni

Istoricul versiunilor aplicației, cea mai nouă prima ([D009](decizii/D009-versionare.md)). Fiecare commit care schimbă codul crește versiunea din `internal/versiune/versiune.go` și adaugă aici o secțiune `## X.Y.Z · AAAA-LL-ZZ`.

Pe lângă versiunea aplicației, fiecare raport și fiecare schimbare a bazei de date au versiunea lor; ele se trec în secțiunea versiunii aplicației care le aduce.

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
