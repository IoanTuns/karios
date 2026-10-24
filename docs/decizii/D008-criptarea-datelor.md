# D008 · Criptarea datelor operaționale

Stare: propusă (varianta de cheie decisă de Ioan Tuns) · Data: 2026-10-24 · Taskuri: F05-T09, F05-T10, F05-T12

## Context

Cerință (Ioan Tuns, 2026-10-24): toate datele operaționale stau lângă `.exe`, în folderul `karios/`, și **nicio informație stocată de aplicație nu poate fi citită decât prin aplicație**. Forma de protecție trebuie să fie reversibilă de către aplicație.

O codare (base64, XOR, compresie) nu îndeplinește cerința: se decodează fără nicio cheie. E nevoie de criptare autentificată, iar întrebarea reală e de unde ia aplicația cheia.

## Decizie

### 1. Ce se criptează

| Fișier | Conținut | Protecție |
| --- | --- | --- |
| `karios/pontix.db` (+ `-wal`, `-shm`) | toate datele, inclusiv setările | criptată pe pagini (punctul 4) |
| `karios/backup/*.db` | copiile de siguranță | criptate cu aceeași cheie |
| `karios/seif` | cheia de date împachetată, ultimul utilizator, întrebările de recuperare | criptat; răspunsurile doar ca hash |
| `karios/pontix.lock` | cine are aplicația deschisă | criptat |
| `karios/webview/` | cache-ul WebView2 | interfața nu păstrează date în browser (fără `localStorage`); folderul se golește la închidere |

- **Setările trec din `config.toml` în baza criptată.** Lângă `.exe` nu rămâne niciun fișier lizibil cu informații.
- **Traducerile corectate** (D004) se importă din aplicație și se păstrează în bază. Fișierele JSON aduse de utilizator sunt doar intrare pentru import.
- **Excepție voită:** exporturile (PDF, Excel) sunt create la cererea utilizatorului, acolo unde alege el, și sunt lizibile.

### 2. Cheia — etapa 1: cheie ascunsă în `.exe` (decizie Ioan Tuns, 2026-10-24)

Se folosesc două niveluri de chei. Așa, nivelul de protecție se poate schimba mai târziu fără recriptarea bazei:

- **Cheia de date:** aleatoare (256 biți), generată la prima pornire, diferită pentru fiecare folder `karios/`. Cu ea se criptează baza și backup-urile.
- **Cheia aplicației:** ascunsă în `.exe`. Cu ea se împachetează cheia de date în `karios/seif` (AES-256-GCM).

**Cheia aplicației nu intră în repo** (AGENTS.md, secțiunea 4). Se injectează la build (`-ldflags`) dintr-un fișier secret, ignorat de git, păstrat de un om. **Aceeași cheie trebuie folosită la toate versiunile.** Dacă se pierde, versiunile noi nu mai pot deschide datele vechi (F05-T12).

### 3. Accesul: validarea ultimului utilizator, apoi întrebări de recuperare

1. `karios/seif` reține **ultimul utilizator**: utilizatorul Windows (SID) și calculatorul.
2. **Același utilizator, același calculator** → aplicația se deschide direct.
3. **Utilizator sau calculator diferit** (folder mutat, alt cont) → aplicația cere **răspunsurile la întrebările de recuperare**. Dacă sunt corecte, noul utilizator devine „ultimul utilizator”.
4. **Întrebările se stabilesc la prima pornire:** 3 întrebări alese sau scrise de utilizator, cu răspunsuri normalizate (fără diferențe de majuscule, diacritice sau spații). Răspunsurile se păstrează doar ca hash Argon2id, niciodată în clar.
5. **După 5 încercări greșite**, așteptarea crește progresiv (1 min, 5 min, 30 min…), iar fiecare încercare se notează în jurnalul de modificări.

**Împachetare dublă de la început (decizie Ioan Tuns, 2026-10-24):** cheia de date se împachetează *și* cu o cheie derivată din răspunsurile la întrebări (Argon2id peste răspunsurile normalizate). Seiful conține deci două copii împachetate ale aceleiași chei de date:

| Copie | Deschisă cu | Folosită când |
| --- | --- | --- |
| A | cheia din `.exe` | ultimul utilizator e același (etapa 1) |
| B | cheia derivată din răspunsuri | recuperare; în etapa 2, singura copie |

Trecerea la etapa 2 înseamnă doar ștergerea copiei A, fără recriptarea bazei. Schimbarea întrebărilor sau a răspunsurilor reface doar copia B.

### 4. Cum se criptează baza

- **Propunere:** driverul SQLite `github.com/ncruces/go-sqlite3`, tot Go pur, fără CGO, cu VFS de criptare pe pagini (Adiantum sau AES-XTS). Toată baza, inclusiv jurnalul WAL, e criptată pe disc, iar interogările SQL funcționează normal.
- **Ar înlocui** `modernc.org/sqlite` din D001, care nu are criptare.
- **Respinse:** SQLCipher (cere CGO); criptarea câmp cu câmp în aplicație (fără căutare și sortare în SQL, ușor de uitat un câmp).
- **De confirmat prin prototip (F05-T09):** modul WAL, `VACUUM INTO` cu baza criptată, `sqlc`, performanța pe o lună de pontaj, ultima versiune stabilă (D007).

## Ce protejează și ce nu (etapa 1)

| Situație | Protejat? |
| --- | --- |
| Cineva deschide `karios/pontix.db` cu un program SQLite sau într-un editor | **da**, vede doar octeți criptați |
| Folderul `karios/` copiat sau furat, fără `.exe` | **da** |
| Folderul mutat pe alt calculator, deschis cu aplicația, fără răspunsurile la întrebări | **da**, cât timp nu se ocolește aplicația |
| Cineva cu `.exe` care analizează programul și scoate cheia | **nu**: o poate face un specialist; programele Go se analizează relativ ușor |
| Răspunsuri la întrebări ușor de ghicit (data nașterii, numele orașului) | **slab**; de aceea întrebările sunt alese de utilizator |

Pentru o protecție reală împotriva unui atacator hotărât e nevoie de etapa 2 (parolă sau răspunsuri care derivă cheia). Ideea e înregistrată ca [I006](../idei.md).

## Variante respinse

- **Codare (base64, XOR)** — nu protejează nimic.
- **Windows DPAPI ca singură cheie** — folderul mutat pe alt calculator nu s-ar mai putea deschide; contrazice D002.
- **Parolă la fiecare pornire, din prima etapă** — amânată la cererea lui Ioan Tuns (I006).

## Întrebări deschise

1. **Utilizatori:** rămâne utilizatorul Windows singura identitate, sau se fac conturi în aplicație? Conturile ar completa „Cine” în jurnalul de modificări (F12) și s-ar lega de autentificarea viitoare (F13).
2. **Blocarea automată** după inactivitate: da sau nu? Cu etapa 1, ar însemna reverificarea utilizatorului sau a întrebărilor.
3. **Cine păstrează cheia aplicației** și unde (F05-T12).

## Consecințe

- Datele nu pot fi citite cu alte programe; accesul trece prin aplicație.
- Apare un ecran pentru stabilirea întrebărilor la prima pornire și unul pentru recuperare, cu texte traduse (D004). Limba lor se ia din Windows, fiindcă setările sunt încă criptate în acel moment.
- Cheia aplicației devine un secret critic al proiectului (F05-T12).
- După acceptare se schimbă: D001 (driverul SQLite), D004 (traducerile din folder devin import), F05-T04 (setările în bază).
