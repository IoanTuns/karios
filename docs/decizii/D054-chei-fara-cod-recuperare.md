# D054 · Cheile datelor fără secret în cod, cu cheie de recuperare

Stare: acceptată · Data: 2026-10-30 · Decisă de: Ioan Tuns · Înlocuiește: [D011](D011-cheia-aplicatiei-in-cod.md) · Modifică: [D008](D008-criptarea-datelor.md) punctul 2 · Taskuri: F15-T10 … F15-T13

## Context

Repo-ul este public și rămâne public (Ioan Tuns, 2026-10-30). Cheia aplicației din `internal/seif` (D011) e deci publică: oricine poate deschide copia A a oricărui seif și citi baza, fără parolă. Cerința: cheia trebuie scoasă din cod, **iar baza trebuie să rămână recuperabilă dacă se întâmplă ceva cu aplicația**.

## Decizie (Ioan Tuns: „da la ambele”)

1. **Nicio cheie secretă în cod.** Securitatea depinde doar de secrete care stau la utilizatori. Algoritmul e public, deci orice versiune a aplicației, chiar compilată din nou din repo, deschide datele cu un secret valid.
2. **Copiile cheii de date din seif:**

   | Copie | Se deschide cu | Folosită când |
   | --- | --- | --- |
   | `W:<SID>` | protecția datelor din Windows (DPAPI), pentru utilizatorul Windows curent pe acest calculator | pornirea obișnuită: fără parolă la deschiderea datelor; auto-login-ul D008 rămâne |
   | `P:<cont>` / `R:<cont>` | parola, respectiv răspunsurile unui cont (Argon2id) | deblocarea pe alt calculator sau alt utilizator Windows; recuperarea parolei |
   | `K` | **cheia de recuperare**, generată de aplicație (25 de caractere în 5 grupe) | ultima soluție, când niciun cont nu mai funcționează |

   La deblocarea cu parolă sau cu cheia de recuperare, aplicația adaugă copia `W` pentru utilizatorul Windows curent.
3. **Cheia de recuperare e obligatorie.** Se afișează administratorului după autentificare, care o tastează înapoi pentru confirmare. Până la confirmare, aplicația nu permite lucrul: administratorul vede cheia, ceilalți un mesaj să-l cheme. Cât timp nu e confirmată, se generează alta la fiecare pornire.
4. **Trecerea datelor existente** (o singură dată, la prima pornire a versiunii noi):
   - seiful se deschide cu cheia veche, publică, păstrată în cod doar pentru această trecere;
   - se generează o **cheie de date nouă**; baza și backup-urile se recriptează cu ea, iar cheia veche nu mai deschide nimic din datele actuale;
   - copiile `P`/`R` vechi (legate de cheia veche) se păstrează temporar ca `P0`/`R0` și se refac la prima autentificare a fiecărui cont;
   - copia A se șterge; trecerea se notează în jurnal.
   - Fișierele noi (bază și seif) se scriu alături și se înlocuiesc la final; o trecere întreruptă se termină la pornirea următoare.

## Ce rămâne în afara protecției

- Datele copiate în altă parte **înainte** de trecere rămân citibile cu cheia publică veche.
- Dacă se pierd simultan toate parolele și răspunsurile, cheia de recuperare și profilul Windows, datele nu se mai pot citi. De aceea cheia de recuperare se confirmă obligatoriu și se recomandă cel puțin doi administratori.
- Un cont care nu s-a autentificat după trecere (are doar `P0`) nu poate debloca datele pe alt calculator până nu se autentifică o dată pe calculatorul obișnuit.

## Variante respinse

- **Cheia injectată la build dintr-un fișier secret** — scoasă din cod, dar extrasă ușor din `.exe`; executabilul circulă.
- **Doar parolă (fără W)** — sigur, dar fără auto-login (era motivul amânării I006).
- **Fără cheie de recuperare** — pierderea tuturor parolelor ar pierde datele.
