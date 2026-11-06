# Integrarea în `main` · 2026-11-06

Pregătită de Claude. Merge-ul și push-ul le face un om (AGENTS.md §2–3).

## Pe scurt

- `main` = `origin/main` = `e4b6a8f` (0.89.3, 2026-10-30). Nimeni nu a mai scris pe `main`.
- Vârful lanțului: ramura `F05-T11-doar-citire`, commit `b8b7603`, versiunea **0.92.45**.
- `main` este strămoș direct al vârfului, deci integrarea este un **fast-forward**: fără commit de merge și fără conflicte.
- 97 de commituri, din circa 45 de taskuri. Fiecare ramură de task din lanț pornește din cea anterioară, deci vârful le conține pe toate.

## Ce intră

| Zonă | Taskuri | Pe scurt |
| --- | --- | --- |
| Date de bază (F07) | T24–T40 | nomenclatoare, mărci de gardă, structura contractelor, norma lunară și soldurile, compartimente, zile lucrătoare, absențe (din pontaj, prioritate, efecte, categorii), ștergerea configurărilor de către administrator, prefixul angajatului |
| Rapoarte (F09) | T11–T26 | legendă, editor flexibil SQL, explicații opționale, gărzi, aliniere, personalizare, editorul vizual pdfme, fonturi Windows, stiluri pe rând și celulă, rapoartele standard din SQL, conversia modelelor și ștergerea motorului vechi, antet repetat |
| Planificare (F11) | T20–T21 | prioritate planificare/pontaj, programe multiple pe aceeași marcă și zi |
| Interfață (F16) | T23–T33 | focus și zoom, fereastră adaptivă, ore 24h, filtre de stare, audituri UI/UX și remediile lor, anexa rapoartelor, analiza UI pdfme |
| Livrare, stocare | F02-T24, F05-T11, F06-T06, F14-T21 | build 0.89.3, a doua instanță doar în citire, istoricul migrărilor și refuzul bazei mai noi, verificarea dependențelor |

Detaliile sunt în [versiuni.md](versiuni.md) (0.89.4 … 0.92.45) și în [jurnal.md](jurnal.md).

## Baza de date

Intră 11 migrări noi, toate adăugate; nicio migrare existentă în `main` nu este modificată:

0024 nomenclatoare · 0025 mărci de gardă · 0026 structura contractelor · 0027 clasificări schimburi · 0028 inițializare solduri · 0029 efecte absențe · 0030 categorii absențe · 0031 programe multiple · 0032 prefixe · 0033 modele standard SQL · 0034 istoric migrări.

Migrările nu se pot desface. Aplicația face singură o copie a bazei înaintea lor (D002). Din 0.92.44, o versiune mai veche refuză baza migrată (protecție la downgrade).

## Ce se întâmplă la prima pornire a versiunii integrate

1. Copia automată a bazei, apoi cele 11 migrări.
2. La prima deschidere a rapoartelor, după autentificare, modelele de raport vechi se convertesc la editorul vizual (F09-T23). Versiunea veche a fiecărui model rămâne în Jurnal modificări. Modelele Markdown și Poziționare liberă se convertesc aproximativ și trebuie verificate.
3. Fișierul `karios/pontix.lock` ține evidența instanței care scrie; a doua copie pornește doar în citire (F05-T11).

## Verificări pe `b8b7603` (2026-11-06)

- `gofmt -l .`: curat. `go vet ./...`: fără probleme. `go test ./...`: toate trec (57 s).
- Interfață: TypeScript curat, vitest 45 trecute (1 generator sărit intenționat), build reușit.
- Documentație: 0 erori. Manuale: 3 teste OK. Livrare: 7 teste OK.
- Scenarii Edge (punte simulată): 8 din 9 trec (rapoarte pdfme, rapoarte flexibile, ux30, focus, nomenclatoare, programe multiple, absențe din pontaj, ștergere admin).
- `--ajutor` **pică**: „Navigarea AI ocolește confirmările” (0 confirmări în loc de 2). Pică identic pe `07dc462`, deci exista înainte de lucrul din 2026-11-05/06. Consemnat ca [I117](idei.md); poate fi un defect real (navigarea din asistentul AI fără confirmarea modificărilor nesalvate) și trebuie investigat separat.

## Neverificat

- Prima pornire pe datele reale (migrările 0024–0034 și conversia rapoartelor).
- Funcțiile de după 0.92.36 în fereastra aplicației: rapoartele standard pe SQL, conversia, antetul repetat, două executabile deschise simultan, mesajul pentru baza mai nouă. Utilizatorul a confirmat funcționarea pentru 0.92.36.
- Semnarea executabilului și proba pe un calculator curat (F05-T08, F05-T07).

## Pași pentru integrare (de făcut de un om)

1. **Copie a datelor reale:** din aplicația curentă, Creează backup, sau o copie a întregului folder `karios/` cu aplicația închisă.
2. **Merge**, în checkout-ul principal (`D:\karios\karios`, pe `main`, curat):

   ```powershell
   git merge --ff-only F05-T11-doar-citire
   ```

   `--ff-only` garantează că `main` doar avansează la vârful ramurii (codul de la `b8b7603`, plus commiturile de documentație ale acestei note), fără commit de merge. Dacă refuză, înseamnă că `main` s-a schimbat între timp; atunci nu forța.
3. **Verificare în checkout-ul principal** (opțională, dar recomandată):

   ```powershell
   gofmt -l .; go vet ./...; go test ./...
   cd frontend; npm ci; npm test; npm run build; cd ..
   ```

4. **Build și probă pe date reale:** `.\build.ps1`, apoi pornirea pe copia datelor de la pasul 1.
5. **Push** (cere aprobare): `git push origin main`.
6. **Curățenie**, opțional, după push:
   - Ramurile de task din lanț pot fi șterse cu `git branch -d <ramura>`. `-d` refuză orice ramură neintegrată.
   - Worktree-ul `D:\karios\karios-claude` (ramura `F16-T15-raporteaza-problema`, 2026-10-30) **nu se integrează**: funcția ei este deja în lanț, iar un merge ar da conflicte în jurnal, versiuni și cataloage. Poate fi închis după ce confirmi că nu are lucru nesalvat.

## Riscuri

- Volumul: 0.89.3 → 0.92.45 aduce schimbări mari de date (contracte, absențe, solduri) și de rapoarte. Copia de la pasul 1 este plasa de siguranță.
- Conversia rapoartelor scrie în modelele salvate. Originalele rămân în jurnal, dar restaurarea lor e manuală; copia de la pasul 1 e mai simplă.
- Fișierul `docs/exemple/garzi-lunar.sql` are în `D:\karios\karios-codex` o modificare nesalvată a utilizatorului; nu face parte din lanț și nu este atinsă.
