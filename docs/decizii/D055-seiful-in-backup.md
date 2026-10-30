# D055 · O copie a seifului lângă backup-uri

Stare: acceptată · Data: 2026-10-30 · Decisă de: Ioan Tuns („începe cu propunerile tale în ordine”) · Modifică: [D054](D054-chei-fara-cod-recuperare.md) · Taskuri: F05-T27

## Context

După D054, cheia de date stă doar în fișierul `karios/seif`, în copii criptate (W, P/R, K). Backup-urile din `karios/backup/` conțin doar baza. Dacă `seif` se pierde sau se strică, nici backup-urile, nici cheia de recuperare de pe hârtie nu mai deschid datele: copia K e chiar în fișierul pierdut.

## Decizie

1. **Oglinda seifului:** la fiecare salvare a seifului (parolă nouă, cheie de recuperare, copie W, contorul jurnalului), același conținut se scrie și în `karios/backup/seif`. La fiecare pornire deblocată, oglinda se rescrie, ca să existe și pentru folderele vechi.
2. **Refacerea la pornire:** dacă `karios/seif` lipsește sau e deteriorat, iar `karios/backup/seif` e valid, seiful se reface din oglindă. Seiful stricat se păstrează ca `seif.deteriorat`. Refacerea intră în jurnalul de audit (`seif_refacut`).
3. **Siguranța:** oglinda conține doar copiile criptate ale cheii de date. Fără parola unui cont, cheia de recuperare sau utilizatorul Windows potrivit, nu deschide nimic. Copierea folderului `backup/` pe un stick sau pe alt calculator ia cu ea tot ce trebuie pentru recuperare: baza și seiful.
4. Oglinda nu se rotește și nu are istoric: cheia de date nu se mai schimbă după D054, deci aceeași oglindă deschide toate backup-urile.

## Variante respinse

- **Seiful în fiecare fișier de backup** (arhivă bază + seif): schimbă formatul backup-ului, restaurarea și verificarea; câștigul e același cu o singură oglindă, pentru că cheia de date e aceeași.
- **Oglinda scrisă doar la backup-ul zilnic:** după schimbarea cheii de recuperare, oglinda ar mai păstra o zi copia K veche, care ar deschide în continuare datele.

## Consecințe

- Ștergerea sau coruperea lui `seif` nu mai înseamnă pierderea datelor, dacă `backup/` e întreg.
- Pentru o pierdere a discului rămâne nevoia de a copia periodic `backup/` în altă parte (de spus în ghid).
- Cine are acces la `backup/seif` are aceleași șanse ca la `seif`: forța brută pe parole (Argon2id) sau cheia de recuperare.
