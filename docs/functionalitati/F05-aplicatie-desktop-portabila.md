# F05 · Aplicația desktop portabilă

Etapă: V1 · Stare: în lucru · Depinde de: F02

## Scop

Un singur `.exe` care pornește pe orice Windows, fără instalare și fără drepturi de administrator, și ține totul în folderul lui.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — Wails + React
- [D002](../decizii/D002-portabilitate.md) — regulile pentru portabilitate (propusă)

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F05-T19 | Ergonomia pontajului și salvarea formularelor lungi, primul lot după audit | gata | Codex | 0.28.0 pe F05-T19-ergonomie-pontaj, din 7350a70; I045–I048, D035 acceptată explicit. Grupare Filtre/Completare/Ajutor, legendă și contrast, eliminarea avertismentului T/Ore norm, context de sesiune cu selecție/derulare și confirmarea abandonării, dialoguri cu corp derulant și salvare persistentă la reguli. Edge 1024×640: 7 rânduri complete, fără derularea paginii; scenarii de conflict/renunțare/RO-EN și verificările Go trecute. Fără schimbarea calculelor sau datelor operaționale. [Rezultate](../rezultate-ui-F05-T19.md), [lista UI](../lista-lucru-ui.md). Pas următor: F05-T20 și probă cu operator. |
| F05-T20 | Motive pentru valori necalculate și validarea câmpurilor de configurare | de făcut | | I049; [lista UI](../lista-lucru-ui.md). |
| F05-T21 | Ergonomia listei de absențe și a cererilor | de făcut | | I050; [lista UI](../lista-lucru-ui.md). |
| F05-T22 | Inserarea și validarea variabilelor în editorul de rapoarte | de făcut | | I051; [lista UI](../lista-lucru-ui.md). |
| F05-T23 | Simplificarea controalelor și navigării din Setări | de făcut | | I052, respectând D034; [lista UI](../lista-lucru-ui.md). |
| F05-T24 | Avertizare la închiderea ferestrei cu modificări nesalvate (`OnBeforeClose`; `beforeunload` nu oprește Wails) (UX-1) | gata | Claude | 0.69.0: `OnBeforeClose` + `SeteazaModificariNesalvate`; dialogul nativ neverificat în aplicația pornită |
| F05-T25 | Restaurarea unui backup din aplicație, cu verificarea sigiliilor copiei și notare în jurnal | gata | Claude | 0.70.0: verificare pe copie temporară, aplicare la pornire, copie de siguranță, jurnal; doar administrator |
| F05-T26 | Backup automat zilnic la prima pornire din zi, cu păstrarea ultimelor 7 | gata | Claude | 0.71.0: la prima pornire din zi, rotația comună de 7; fără copie pentru bază nouă sau cu incident |
| F05-T18 | Analiza UI/UX a versiunii 0.27.0: avantaje, limite și îmbunătățiri prioritizate | gata | Codex | Ramură F05-T18-audit-ui din 2f55173. Revizuire de cod și Edge pe date fictive la 1440 și 1024 px, grilă 50×31, capturi și măsurători separate. Confirmate: grilă la y=629 pe 1024×640, resetarea perioadei/căutării, texte depășite și comparație ambiguă T/Ore norm, contrast redus și salvare departe de editare. Fluxurile UI, gofmt/vet și testele Go au trecut. Numai documentație; UI, calculele, versiunea și bazele sunt nemodificate. Pas următor: alegerea lotului din I045–I052/D035 pentru implementare și probă cu un operator. [audit](../analiza-ui-2026-10-26.md), [D035](../decizii/D035-prioritati-ui-dupa-audit.md) |
| F05-T01 | Stabilește ultima versiune stabilă de Wails și verifică ce implică pentru proiect | gata | Claude | 2026-10-24: **Wails v2.17.0** (v3 e încă beta, exclus de D007; vezi I007). Node.js **24.22.1 LTS**. Wails v2 are opțiunea `WebviewUserDataPath` cerută de D002 |
| F05-T02 | Scheletul Wails + React + TypeScript | gata | Claude | 2026-10-24, versiunea 0.2.0: Wails 2.17.0, React 19.3.1, TypeScript 7.0.2, Vite 8.3.5; `pontix.exe` 11,5 MB. Nepornit încă: întâi F05-T03 (altfel WebView2 scrie în `%APPDATA%`) |
| F05-T03 | Datele în `karios/` lângă `.exe`; cache-ul WebView2 în `karios/webview/` | gata | Claude | 2026-10-24, versiunea 0.3.0: `internal/portabil`, `WebviewUserDataPath`. Verificat: după 8 s de rulare, niciun folder nou în `%APPDATA%` sau `%LOCALAPPDATA%`; cache-ul în `karios/webview/` |
| F05-T04 | Setările aplicației | de făcut | | Include F03-T05. Dacă D008 e acceptată: în baza criptată, nu în `config.toml` |
| F05-T05 | Verificarea dreptului de scriere în `karios/`, cu mesaj tradus | gata | Claude | 2026-10-24, versiunea 0.4.0: fereastră Windows de eroare în limba Windows. Verificat pe `.exe`: folder fără drept de scriere (ACL) și fișier `karios` care blochează folderul; în ambele cazuri apare mesajul și nu se creează nimic. Modul WAL a trecut la F06-T02 |
| F05-T06 | Backup cu `VACUUM INTO`, automat înainte de migrări; cel mult 7 copii păstrate | gata | Claude | 2026-10-25, versiunea 0.8.0: backup criptat cu `VACUUM INTO` în `karios/backup/` înaintea oricărei migrări pe o bază cu date; dacă backup-ul eșuează, migrarea nu pornește. Se păstrează cele mai noi 7 (decizie Ioan Tuns). Verificat prin teste (10 actualizări simulate); pe `.exe` nu, fiindcă încă nu există o a doua migrare. Backup-ul la închiderea lunii: F12-T01 |
| F05-T07 | Probă pe un PC Windows curat și de pe un stick USB | de făcut | | |
| F05-T08 | Semnarea digitală a `.exe` | de făcut | | De ales furnizorul și costul; fără ea apar avertismente SmartScreen |
| F05-T09 | Prototip SQLite criptat cu `ncruces/go-sqlite3`: WAL, `VACUUM INTO`, `sqlc`, performanță | gata | Claude | 2026-10-24, versiunea 0.3.1: toate verificările trec pe `adiantum` și `xts`; propus `xts`. Limite: fără detectarea modificărilor (backup-uri semnate), criptare deterministă. Rezultate în D008 |
| F05-T10 | Seiful: cheia de date cu copiile A (cheia din `.exe`), P și R (pe cont) | gata | Claude | 2026-10-25, versiunea 0.6.0: `internal/seif`. Copia A (cheia aplicației, AES-256-GCM) și copii deschise cu un secret (Argon2id 64 MiB), id-ul copiei în datele autentificate, scriere atomică. 8 teste, inclusiv copie mutată, fișier modificat, etapa 2. Neintegrat încă în aplicație: depinde de F05-T12 |
| F05-T11 | Fișierul de blocare: a doua instanță pornește doar în citire | de făcut | | D002 punctul 6 (aprobat) |
| F05-T12 | Gestionarea cheii aplicației: generare, păstrare în afara repo-ului, injectare la build; integrarea seifului și a bazei în pornirea aplicației | gata | Claude | 2026-10-25, versiunea 0.7.0: cheia în cod ([D011](../decizii/D011-cheia-aplicatiei-in-cod.md)); `internal/stocare` deschide seiful și baza la pornire. Verificat pe `.exe`: prima pornire, a doua pornire după oprire forțată, seif lipsă (nu se creează unul nou) |
| F05-T13 | Script de build: versiunea din `versiune.go` în proprietățile `.exe` și commit-ul prin `-ldflags` | gata | Claude | 2026-10-24, versiunea 0.3.3: `build.ps1`; versiunea în proprietățile `.exe` (verificată în Explorer: 0.3.3, Pontix, copyright) și commit-ul prin `-ldflags`. Scriptul verifică numele și versiunea prin Windows Shell, aceeași sursă ca Proprietăți din Explorer |
| F05-T14 | Semnătura backup-urilor (HMAC-SHA256), verificată înainte de restaurare | de făcut | | D008 punctul 4: criptarea nu detectează modificările |
| F05-T15 | Build fără curățarea folderului de date și cu nume alternativ pentru executabil | gata | Codex | 2026-10-25, 0.12.1: eliminat `-clean`, parametru `-NumeExecutabil`, validarea căii și refuzul fișierului blocat înainte de build. Test PowerShell și build real alternativ reușite. După incident, utilizatorul a autorizat regenerarea datelor exclusiv demo, în folder separat; baza veche păstrată. [D017](../decizii/D017-build-fara-stergerea-datelor.md) |
| F05-T16 | Copie de siguranță criptată la cerere din ecranul Setări | gata | Antigravity | 2026-10-26, 0.24.0: `stocare.BackupManual` (`VACUUM INTO` criptat cu cheia din seif în `karios/backup/pontix-manual-<UTC>-v<Aplicatie>.db`), fără a afecta rotația celor 7 copii automate la migrare; buton în `Setări` cu afișarea căii salvate. Pas următor: F05-T14 (semnătura HMAC-SHA256 a backup-urilor). [D031](../decizii/D031-precompletare-avertismente-pontaj-backup-manual.md) rămâne propusă |
| F05-T17 | Gruparea secțiunilor de configurare (General, Zile libere legale, Schimburi, Configurare pontaj, Definire rapoarte) în Setări și temă vizuală albastru-ardezie pentru Setări | gata | Antigravity | 2026-10-26, 0.27.0: meniu lateral cu 5 intrări principale (`Angajați`, `Pontaj`, `Concedii și absențe`, `Rapoarte`, `Setări`), sub-meniu sub `Setări` și taburi în partea de sus a ecranului Setări pentru cele 5 secțiuni, redenumirea `ui.sarbatori` în `Zile libere legale` și temă albastru-ardezie (`#2b4c7e` / `#eef3f9`) pentru zona Setări. I044, [D034](../decizii/D034-grupare-meniu-setari-tema-vizuala.md) |
