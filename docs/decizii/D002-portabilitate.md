# D002 · Reguli pentru portabilitate

Stare: acceptată · Data: 2026-10-24 · Decisă de: Ioan Tuns · Taskuri: F05-T03 … F05-T07, F05-T11

## Context

Aplicația trebuie să fie un singur folder care se poate copia sau muta: fără instalare, fără drepturi de administrator. Revizuirea din 2026-10-24 a găsit locuri unde, cu setările implicite, regula nu se respectă sau datele se pot pierde.

**Cerință (Ioan Tuns, 2026-10-24):** toate datele operaționale, în afara executabilului, stau lângă `.exe`, într-un folder cu numele repo-ului, `karios/`. Nimic din ce stochează aplicația nu poate fi citit decât prin aplicație. Criptarea e tratată în [D008](D008-criptarea-datelor.md).

Structura folderului:

```
pontix.exe
karios/
  pontix.db            criptată (+ pontix.db-wal, pontix.db-shm cât timp rulează)
  seif                 cheia de date împachetată (D008)
  pontix.lock          cine are aplicația deschisă
  backup/              copiile făcute cu VACUUM INTO, criptate
  webview/             cache-ul WebView2, fără date
```

## Decizie

### 1. Cache-ul WebView2 în `karios/webview/`

- **Problema:** Wails pornește WebView2 cu folderul de date implicit în profilul utilizatorului (`%APPDATA%` / `%LOCALAPPDATA%`). Acolo ajung cache, cookie-uri și `localStorage`, care rămân pe calculator după mutarea folderului.
- **Regula:** opțiunea Wails pentru folderul de date WebView2 (`WebviewUserDataPath`) se setează la `karios/webview/`. Interfața nu păstrează date în browser (D008).
- **Verificare:** după prima pornire pe un PC curat, în `%APPDATA%` și `%LOCALAPPDATA%` nu apare niciun folder al aplicației.

### 2. Fusul orar inclus în `.exe`

- **Problema:** Windows nu are o bază de fusuri orare pe care Go s-o poată folosi; pe un PC fără Go aplicația n-ar porni.
- **Regula:** `import _ "time/tzdata"` (aproximativ 450 KB în plus). **Aplicat** în `internal/calcul/fus.go`.

### 3. Backup doar cu `VACUUM INTO`

- **Problema:** în modul WAL, modificările recente stau în `pontix.db-wal`. Copierea doar a fișierului `.db` cu aplicația deschisă le pierde; copierea `.db` + `.wal` nu e atomică.
- **Regula:** backup-ul se face numai din aplicație, cu `VACUUM INTO 'karios/backup/pontix-AAAA-LL-ZZ-HHMM.db'`, criptat cu aceeași cheie (D008).
- **Backup automat** înainte de orice migrare a bazei (la fiecare versiune nouă, D007) și la închiderea lunii (F12). Se păstrează ultimele N copii, cu N configurabil.
- **Verificare:** test care scrie date, face backup fără checkpoint și deschide copia.

### 4. Discuri din rețea — amânat

**Decizie (Ioan Tuns, 2026-10-24):** se poate adăuga ulterior, ca opțiune atașată de utilizator. În prima versiune, `karios/` stă local, lângă `.exe`. Ideea e înregistrată ca [I005](../idei.md).

Ce trebuie știut când se va face: WAL nu funcționează pe discuri din rețea, iar pe SMB doi utilizatori simultan pot corupe baza în orice mod. Același risc îl au folderele sincronizate (OneDrive, Dropbox).

### 5. Folder fără drept de scriere

- **Problema:** dacă aplicația e pusă în `Program Files` sau pe un mediu protejat la scriere, SQLite eșuează cu o eroare greu de înțeles.
- **Regula:** la pornire, aplicația încearcă să creeze un fișier temporar în `karios/`. Dacă nu reușește, afișează un mesaj tradus (D004) care spune unde să fie mutată aplicația.
- **Verificare:** pornire din `Program Files` fără drepturi de administrator.

### 6. O singură instanță care scrie — aprobat

**Decizie (Ioan Tuns, 2026-10-24):** aprobat.

- **Problema:** două copii pornite pe același folder ar scrie simultan în aceeași bază.
- **Regula:** `karios/pontix.lock` e ținut deschis cât rulează aplicația. A doua instanță arată cine o are deschisă (calculator, utilizator, de la ce oră) și **pornește doar în citire**.
- **Verificare:** pornirea simultană a două copii; a doua nu poate salva nimic.

## Variante respinse

- **Date în `%APPDATA%`** — contrazice cerința de portabilitate.
- **Backup prin copierea fișierului** — nesigur în modul WAL.
- **A doua instanță refuzată complet** — a fost preferată pornirea doar în citire (punctul 6).

## Consecințe

- Mutarea aplicației înseamnă mutarea a două lucruri: `pontix.exe` și folderul `karios/`.
- Aplicația e pentru un singur utilizator odată. Lucrul simultan, dacă va fi nevoie, se face prin sincronizare cu un server (F13), nu printr-o bază comună.
- **Excepție acceptată:** tokenurile de autentificare (F13) stau în Windows Credential Manager și nu se mută cu folderul.
