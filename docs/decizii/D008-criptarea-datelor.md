# D008 · Criptarea datelor și accesul

Stare: acceptată · Data: 2026-10-24 · Decisă de: Ioan Tuns · Taskuri: F05-T09, F05-T10, F05-T12, F15

> Notă tehnică (2026-10-25, F06-T02): AES-256-XTS folosește două chei AES-256, deci cheia de date are 64 de octeți (512 biți), nu 256 de biți. Nivelul de protecție rămâne AES-256.

## Context

Cerință (Ioan Tuns, 2026-10-24): toate datele operaționale stau lângă `.exe`, în folderul `karios/`, și **nicio informație stocată de aplicație nu poate fi citită decât prin aplicație**. Forma de protecție trebuie să fie reversibilă de către aplicație.

O codare (base64, XOR, compresie) nu îndeplinește cerința: se decodează fără nicio cheie. E nevoie de criptare autentificată, iar întrebarea reală e de unde ia aplicația cheia.

## Decizie

### 1. Ce se criptează

| Fișier | Conținut | Protecție |
| --- | --- | --- |
| `karios/pontix.db` (+ `-wal`, `-shm`) | toate datele, inclusiv setările și conturile | criptată pe pagini (punctul 4) |
| `karios/backup/*.db` | copiile de siguranță | criptate cu aceeași cheie |
| `karios/seif` | copiile împachetate ale cheii de date, ultimul cont | criptat; parolele și răspunsurile doar ca hash |
| `karios/pontix.lock` | cine are aplicația deschisă | criptat |
| `karios/webview/` | cache-ul WebView2 | interfața nu păstrează date în browser (fără `localStorage`); folderul se golește la închidere |

- **Setările trec din `config.toml` în baza criptată.** Lângă `.exe` nu rămâne niciun fișier lizibil cu informații.
- **Traducerile corectate** (D004) se importă din aplicație și se păstrează în bază. Fișierele JSON aduse de utilizator sunt doar intrare pentru import.
- **Excepție voită:** exporturile (PDF, Excel) sunt create la cererea utilizatorului, acolo unde alege el, și sunt lizibile.

### 2. Cheia — etapa 1: cheie ascunsă în `.exe` (decizie Ioan Tuns)

- **Cheia de date:** aleatoare (256 biți), generată la prima pornire, diferită pentru fiecare folder `karios/`. Cu ea se criptează baza și backup-urile.
- **Seiful păstrează mai multe copii împachetate (AES-256-GCM) ale aceleiași chei de date:**

| Copie | Deschisă cu | Folosită când |
| --- | --- | --- |
| A | cheia aplicației, ascunsă în `.exe` | pornire normală, auto-login (etapa 1) |
| P (câte una pe cont) | cheia derivată din parola contului (Argon2id) | login cu parolă; în etapa 2, fără copia A |
| R (câte una pe cont) | cheia derivată din răspunsurile la întrebările de recuperare ale contului | parolă uitată |

- **Împachetarea multiplă e decisă de la început** (Ioan Tuns). Trecerea la etapa 2 înseamnă ștergerea copiei A, fără recriptarea bazei ([I006](../idei.md)). Schimbarea unei parole sau a unor răspunsuri reface doar copia acelui cont.
- **Cheia aplicației nu intră în repo** (AGENTS.md, secțiunea 4). Se injectează la build (`-ldflags`) dintr-un fișier secret, ignorat de git, păstrat de un om. **Aceeași cheie trebuie folosită la toate versiunile.** Dacă se pierde, versiunile noi nu mai pot deschide datele vechi (F05-T12).

### 3. Accesul: conturi în aplicație, cu auto-login pentru ultimul cont (decizii Ioan Tuns)

- **Conturi în aplicație.** Fiecare persoană are cont propriu, cu parolă și 3 întrebări de recuperare. Parolele și răspunsurile se păstrează doar ca hash Argon2id; răspunsurile se normalizează (majuscule, diacritice, spații).
- **Roluri.** Contul creat la prima pornire e **administrator**: creează și dezactivează conturi și poate reseta parola unui utilizator. Celelalte conturi sunt **operatori**. Ultimul administrator nu poate fi șters sau retrogradat.
- **Pornirea:**
  1. Seiful reține ultimul cont, împreună cu utilizatorul Windows (SID) și calculatorul pe care s-a autentificat.
  2. **Același utilizator Windows, același calculator** → **auto-login** în ultimul cont.
  3. **Altfel** (alt cont Windows, alt PC, folder mutat) → ecran de login, cu numele ultimului cont precompletat; se cere parola.
  4. **„Ieșire din cont”** oprește auto-login-ul până la următorul login cu parolă. Așa, cine lasă calculatorul altcuiva poate închide accesul.
- **Parolă uitată:** răspunsurile la întrebările contului sau resetarea de către un administrator.
- **Încercări greșite:** după 5, așteptarea crește progresiv (1 min, 5 min, 30 min…); fiecare încercare se notează în jurnalul de modificări.
- **„Cine”** din jurnalul de modificări (F12) este contul din aplicație.
- **Fără blocare automată la inactivitate în V1** (Ioan Tuns). Se poate adăuga ulterior.

**Risc asumat (auto-login):** oricine folosește aceeași sesiune Windows pe același calculator intră ca ultimul cont. În acest caz, „Cine” din jurnal arată contul, nu neapărat persoana. Măsura e „Ieșire din cont”.

### 4. Cum se criptează baza

- **Propunere:** driverul SQLite `github.com/ncruces/go-sqlite3` (v0.35.6, licență MIT), tot Go pur, fără CGO, cu **VFS-ul `xts` (AES-256-XTS)**. Toată baza, inclusiv jurnalul WAL, e criptată pe disc, iar interogările SQL funcționează normal, prin interfața standard `database/sql`, pe care o folosește și `sqlc`.
- **Ar înlocui** `modernc.org/sqlite` din D001, care nu are criptare.
- **Respinse:** SQLCipher (cere CGO); criptarea câmp cu câmp în aplicație (fără căutare și sortare în SQL, ușor de uitat un câmp).

**Rezultatele prototipului (F05-T09, 2026-10-24),** `prototipuri/sqlite-criptat/`:

| Verificare | `adiantum` | `xts` |
| --- | --- | --- |
| Modul WAL pe Windows | da | da |
| Nimic lizibil pe disc (baza și jurnalul WAL) | da | da |
| `VACUUM INTO` dă un backup criptat, care se redeschide cu cheia | da | da |
| Cheie greșită sau deschidere fără criptare → refuz | da | da |
| 4650 de rânduri (50 angajați × 31 zile × 3 segmente): scriere / totaluri | 19–34 ms / 7–11 ms | 22–51 ms / 7–16 ms |

Un test de control confirmă că verificarea găsește textul într-o bază necriptată, deci rezultatele „nimic lizibil” sunt reale.

**De ce `xts`:** folosește doar AES, un algoritm standard aprobat NIST/FIPS. Procesoarele actuale au instrucțiuni AES, deci diferența de viteză față de `adiantum` e neglijabilă la volumul nostru.

**Limite găsite la prototip (valabile pentru ambele variante):**
- **Criptarea nu detectează modificările.** Cine are fișierul poate strica sau înlocui pagini fără să știe ce conțin, iar aplicația nu-și dă seama. **Măsură:** fiecare backup primește o semnătură (HMAC-SHA256 cu o cheie derivată din cheia de date), verificată înainte de restaurare.
- **Criptarea e deterministă.** Comparând două copii ale aceluiași fișier, se vede care pagini s-au schimbat, nu și ce conțin. Pentru amenințările din etapa 1 e acceptabil.

## Ce protejează și ce nu (etapa 1)

| Situație | Protejat? |
| --- | --- |
| Cineva deschide `karios/pontix.db` cu un program SQLite sau într-un editor | **da**, vede doar octeți criptați |
| Folderul `karios/` copiat sau furat, fără `.exe` | **da** |
| Folderul mutat pe alt calculator și deschis cu aplicația, fără parola unui cont | **da**, cât timp nu se ocolește aplicația |
| Altcineva în aceeași sesiune Windows, pe același calculator | **nu**: auto-login (risc asumat, vezi punctul 3) |
| Cineva cu `.exe` care analizează programul și scoate cheia | **nu**: o poate face un specialist; programele Go se analizează relativ ușor |

Protecția reală împotriva unui atacator hotărât vine în etapa 2 ([I006](../idei.md)): fără copia A și fără auto-login, cheia de date se deschide doar cu parola unui cont.

## Variante respinse

- **Codare (base64, XOR)** — nu protejează nimic.
- **Windows DPAPI ca singură cheie** — folderul mutat pe alt calculator nu s-ar mai putea deschide; contrazice D002.
- **Utilizatorul Windows ca identitate** — înlocuit de conturile din aplicație (Ioan Tuns).
- **Parolă la fiecare pornire** — înlocuită de auto-login pentru ultimul cont (Ioan Tuns).
- **Toate conturile egale** — înlocuit de rolurile administrator și operator (Ioan Tuns).

## Întrebări deschise

1. **Cine păstrează cheia aplicației** și unde (F05-T12).

## Consecințe

- Datele nu pot fi citite cu alte programe; accesul trece prin aplicație și printr-un cont.
- Ecrane noi, cu texte traduse (D004): crearea administratorului la prima pornire, login, recuperare, gestionarea conturilor. Limba ecranelor dinaintea deblocării se ia din Windows.
- Conturile devin o funcționalitate de V1 ([F15](../functionalitati/F15-conturi-acces.md)); jurnalul de modificări (F12) depinde de ele.
- Cheia aplicației devine un secret critic al proiectului (F05-T12).
- După acceptare se schimbă: D001 (driverul SQLite), D004 (traducerile din folder devin import), F05-T04 (setările în bază).
