# Review UI, UX, performanță și securitate · Pontix 0.59.0 · 28 octombrie 2026

Claude · Bază analizată: `156c1bc` (ramura `F11-T01-planificare-publicare`). Constatările sunt propuneri; review-ul nu modifică codul.

## Metodă și limite

- Citire de cod: `main.go`, `app*.go`, `internal/ai`, `internal/seif`, `internal/stocare`, `internal/bd`, `internal/portabil`, `internal/diagnoza`, `frontend/src` (App, Pontaj, Planificare, AsistentAI, Toast, style.css).
- Rulat `go vet ./...` (fără probleme) și `go test ./...` (toate pachetele trec).
- **Neverificat:** aplicația nu a fost pornită; nu s-au făcut măsurători de performanță, audit WCAG complet sau teste cu `-race`. Pentru UI se bazează pe cod și pe [analiza UI din 26.10](analiza-ui-2026-10-26.md), nu pe capturi noi.

Prioritate: **P1** = de rezolvat înainte de utilizare cu date reale; **P2** = lotul următor; **P3** = îmbunătățire.

---

## 1. UI (interfață vizuală)

### Pros

- Identitate vizuală coerentă: paletă temperată, spațiere și tipografie consecvente (confirmat și în analiza din 26.10).
- Mod întunecat complet (`[data-theme="dark"]`), bară laterală pliabilă și trei densități de grilă (S/M/L), toate persistate local.
- `prefers-reduced-motion` respectat în `style.css`.
- Coloanele nume/marcă fixe în grile, weekend și sărbători marcate distinct, tooltip-uri pe celule.
- Toasturi cu `aria-live="polite"` și `role="status"`; `ErrorBoundary` evită ecranul alb la o eroare React.

### Cons

| ID | Prio | Constatare | Propunere |
| --- | --- | --- | --- |
| UI-1 | P2 | Tema implicită e `light` fixă; `prefers-color-scheme` din Windows nu e folosit la prima pornire (`App.tsx:72`). | Valoare implicită din `matchMedia('(prefers-color-scheme: dark)')` când nu există preferință salvată. |
| UI-2 | P2 | Modul întunecat e scris ca suprascrieri punctuale pe selectori (`style.css:684+`), nu ca variabile CSS. Orice componentă nouă riscă să rămână deschisă la culoare în dark mode. | Culori ca tokenuri pe `:root`, redefinite în `[data-theme="dark"]`. |
| UI-3 | P2 | Contrastele sub 4,5:1 identificate pe 26.10 (ore pe weekend 3,82:1, pe selecție 3,93:1, „—” 3,14:1) nu apar rezolvate; în dark mode orele sunt `#6ee7b7` pe fundaluri variate, neverificate. | Verificare de contrast pe ambele teme, inclusiv pentru stările weekend/selecție/sărbătoare. |
| UI-4 | P3 | Text pentru utilizator scris direct în cod, contrar regulii §8: `aria-label="Notificări de sistem"` (`Toast.tsx:44`), titlul și filtrul dialogului de salvare Excel (`app_export.go:59-62`), mesajele `StareMediu.Mesaj` din `internal/ai/mediu.go`, textele din `DetaliiDezvoltare` și din `diagnoza.CitesteUltimeleLinii`. În engleză apar în română. | Chei noi în `ro.json` / `en.json`. |

### Review

Partea vizuală e matură pentru stadiul proiectului. Problema principală e tehnică: tema întunecată e construită prin suprascrieri și se va degrada pe măsură ce se adaugă ecrane. Trecerea la tokenuri CSS e o investiție mică acum și una mare mai târziu. Textele scrise direct în cod încalcă o regulă explicită a proiectului și ar trebui prinse automat (de exemplu un test care caută șiruri cu diacritice în `.go`/`.tsx` din afara cataloagelor).

---

## 2. UX (fluxuri și comportament)

### Pros

- Grila de pontaj e gândită pentru introducere repetitivă: navigare din tastatură (Home/End/PageUp/PageDown, cifre directe), selecție prin tragere, TSV, Undo/Redo.
- Modificările nesalvate sunt semnalate; conflictele de revizie sunt detectate la salvare.
- Asistentul AI propune acțiuni pe care utilizatorul le confirmă explicit prin carduri; răspunsurile analitice (concedii, acoperire, weekenduri) sunt calculate determinist, instant, fără LLM.
- Code-splitting cu `React.lazy` la nivel de modul.

### Cons

| ID | Prio | Constatare | Propunere |
| --- | --- | --- | --- |
| UX-1 | **P1** | Protecția la închidere nu funcționează în Wails: `beforeunload` (`Pontaj.tsx:87`, `Rapoarte.tsx:64`, `ReguliPontaj.tsx:21`) nu oprește închiderea ferestrei native, iar `main.go` nu are `OnBeforeClose`. Închiderea cu X pierde modificările nesalvate fără avertisment. | `OnBeforeClose` în `options.App`, care întreabă frontendul dacă există modificări și afișează un dialog de confirmare. |
| UX-2 | **P1** | În asistent, un mesaj scurt ca „da”, „ok”, „sigur” aplică **toate** propunerile neaplicate din ultimul mesaj care le are (`AsistentAI.tsx:206-227`), oricât de veche ar fi propunerea și fără să le mai arate. Modificările ajung direct în pontaj și planificare. | Confirmarea prin text doar pentru propunerea din ultimul mesaj al asistentului, cu rezumat și buton de confirmare; sau renunțarea la confirmarea prin text. |
| UX-3 | **P1** | Asistentul folosește mereu luna curentă (`AsistentAI.tsx:109-110`, trimisă la 265, 297), nu luna deschisă în Pontaj/Planificare. Întrebările despre altă lună primesc context greșit, iar o propunere pentru altă lună se salvează cu `an/luna` curente (vezi SEC-4). | Luna/anul din contextul de lucru al aplicației, afișate în antetul asistentului. |
| UX-4 | P2 | Ștergerea unei ture prin asistent (`sterge_tura`, `sterge_interval`) golește pontajul, dar nu dezactivează poziția din planificare: cu `codSchimb = ""` nu se găsește niciun schimb (`app_ai.go:468-477`), deci blocul de planificare nu rulează. Utilizatorul vede rezultate diferite în cele două ecrane. | La ștergere, căutarea pozițiilor existente după angajat și dată, indiferent de cod, și dezactivarea lor. |
| UX-5 | P2 | Erorile la salvarea planificării din asistent sunt ignorate (`app_ai.go:519`); răspunsul e `Succes: true` chiar dacă planificarea nu s-a salvat. | Întoarcerea erorii sau a unui succes parțial explicit. |
| UX-6 | P2 | Problemele P1 din analiza 26.10 rămân relevante: pierderea contextului (lună, filtre) la schimbarea modulului (I047), bara de acțiuni nefixă în dialoguri lungi (I048), motivele comune pentru „—” (I049). | Rămân în lista I045–I052. |
| UX-7 | P3 | Exportul Excel scrie fișierul, dar o eroare de scriere (`os.WriteFile`) nu ajunge la utilizator: funcția întoarce conținutul fără `Cale` și fără `Eroare` (`app_export.go:69-76`). | Eroare explicită cu cod i18n. |

### Review

Fluxul principal (pontaj în grilă) e bun și bine gândit pentru operatori. Riscurile UX sunt concentrate în două locuri: **închiderea ferestrei**, unde protecția există doar aparent, și **asistentul AI**, unde confirmarea prin text, luna fixă și sincronizarea incompletă cu planificarea pot produce modificări neintenționate sau greu de explicat. Asistentul scrie direct în datele de pontaj, deci ar trebui tratat cu aceeași strictețe ca formularele: confirmare explicită, context vizibil, rezultat raportat corect.

---

## 3. Performanță

### Pros

- SQLite configurat corect pentru o aplicație desktop: WAL, `synchronous=NORMAL`, `temp_store=MEMORY`, `busy_timeout=5000` (`internal/bd/bd.go:36-41`).
- Previzualizarea calculului e întârziată (debounce) cu `setTimeout` în `Pontaj.tsx:108`, nu la fiecare tastă.
- Hărțile derivate din grile sunt memorate (`useMemo`) pentru calcul, absențe, înlocuiri.
- Bundle inițial redus la ~252 KB prin `React.lazy` (0.54.0); serverul LLM pornește cu parametri potriviți CPU (fire = nuclee − 2, batch 512, un singur slot).
- Profilare locală a operațiilor (`diagnoza.Cronometreaza`) cu prag de 50 ms.

### Cons

| ID | Prio | Constatare | Propunere |
| --- | --- | --- | --- |
| PERF-1 | P2 | Grilele nu sunt virtualizate și nici componentizate pe rânduri: orice schimbare de stare (inclusiv `setSelectie` la fiecare `onMouseEnter` în timpul tragerii, `Pontaj.tsx:584`) re-randează toate celulele cu closure-uri noi. La 50×31 e acceptabil; la 200–300 de angajați (6 000–9 000 de celule) tragerea va sacada. | Rând ca `React.memo` cu props stabile; selecția derivată per rând; virtualizare pe verticală dacă secțiile depășesc ~100 de persoane. |
| PERF-2 | P3 | Calcule refăcute la fiecare randare: `new Map` pentru sărbători (`Pontaj.tsx:159`, fără `useMemo`), `new Intl.DateTimeFormat` pentru fiecare antet de zi (`Pontaj.tsx:554`), `date.sectii.find` pentru fiecare rând. | `useMemo` și un formatter creat o singură dată per limbă. |
| PERF-3 | P2 | `InterogheazaAI` citește la fiecare mesaj toate datele, absențele pe tot anul, planificarea și schița, apoi caută schimburile liniar pentru fiecare poziție (`app_ai.go:134-258`, O(poziții × schimburi)). Contextul crește cu numărul de angajați și poate depăși iar cele 8192 de tokeni (problema din 0.55.0 revine la secții mari). | Hartă `ID → schimb`; citirea absențelor doar pe lună; limită și rezumare a contextului trimis modelului, cu mesaj clar când se depășește. |
| PERF-4 | P3 | `ai.Manager.EsteActiv` și `Porneste` țin mutexul în timpul cererilor HTTP (`/health` cu timeout 800 ms) și al așteptării de până la 45 s la pornire (`internal/ai/proces.go:39-127`). În acest timp orice apel `Port()`/`EsteActiv()` din UI se blochează. `DetecteazaMediu` face și el un `/health` la fiecare apel. | Verificarea de sănătate în afara secțiunii critice; stare „pornește…” expusă UI-ului. |
| PERF-5 | P3 | `diagnoza.scrieLinieJurnal` face `MkdirAll` + `Stat` + `OpenFile` la fiecare linie; `CitesteUltimeleLinii` citește tot fișierul (până la 5 MB) în memorie. | Fișier ținut deschis și rotire la prag; citire de la coadă. |

### Review

Pentru volumul actual (o secție, zeci de angajați) performanța e adecvată, iar baza e configurată corect. Riscurile apar la scalare: grila re-randată complet și contextul AI construit din toate datele cresc liniar cu personalul. Recomandarea e o măsurătoare cu o bază generată de 300 de angajați (`genereaza-date.ps1`) înainte de optimizări, ca prioritățile să vină din cifre.

---

## 4. Securitate

### Pros

- Baza e criptată integral (AES-256-XTS prin VFS `xts`); cheia de date e aleatoare, păstrată în seif cu AES-256-GCM, cu id-ul copiei în datele autentificate (`internal/seif/seif.go:186-209`).
- Argon2id cu parametri salvați per copie, ca să poată fi măriți ulterior; scriere atomică a seifului (temp + `Sync` + `Rename`), permisiuni `0600`/`0700`.
- Cheia se dă prin `PRAGMA hexkey` la conexiune, nu în DSN; backup-urile sunt criptate cu aceeași cheie, iar backup-ul automat precede orice migrare.
- Lipsa seifului lângă o bază existentă e raportată, nu mascată prin crearea altei chei (`internal/stocare/stocare.go:44-50`).
- LLM-ul rulează doar local (`--host 127.0.0.1`), fără ieșire în rețea; niciun `dangerouslySetInnerHTML` sau `eval` în frontend.

### Cons

| ID | Prio | Constatare | Propunere |
| --- | --- | --- | --- |
| SEC-1 | **P1** | Cheia aplicației e în cod (`internal/seif/seif.go:33`, D011). Oricine are acces la repo sau la executabil (șirul e vizibil cu `strings`) poate decripta orice folder `karios/`. Criptarea protejează doar împotriva cuiva care are fișierele fără executabil. | Asumat prin D011/I006/I008; recomand ca etapa 2 (secret de utilizator, F15) să fie condiție înainte de date reale pe calculatoare partajate. |
| SEC-2 | **P1** | Serverul LLM local nu are autentificare și folosește un port fix (28491). (a) Orice proces sau alt utilizator Windows de pe aceeași mașină poate interoga modelul. (b) Mai grav: dacă alt program ascultă deja pe 28491 și răspunde 200 la `/health`, aplicația îl consideră motorul ei (`proces.go:50-53`) și îi trimite promptul cu **numele, mărcile, funcțiile, absențele și turele angajaților**. | `--api-key` aleator generat la fiecare pornire și trimis în antet; port liber ales la pornire; verificarea că procesul care răspunde e copilul lansat de aplicație. |
| SEC-3 | P2 | `PornesteAI(numeModel)` face `filepath.Join(m.folder, numeModel)` fără validare (`proces.go:62-64`): un nume cu `..` duce la orice fișier de pe disc. Frontendul e de încredere, dar e o linie de apărare ușor de adăugat. | Acceptarea doar a numelor din `ModeleDisponibile` (`filepath.Base` + verificare în listă). |
| SEC-4 | P2 | Integritatea datelor din acțiunile AI: identificarea angajatului e tolerantă cu `strings.Contains` în ambele sensuri (`app_ai.go:320-323`). „Ana” se potrivește cu „Ioana”, „Pop” cu „Popescu”; se ia primul rezultat. Datele propunerii nu sunt verificate că aparțin lunii `an/luna` în care se salvează (`app_ai.go:429-463`). | Potrivire exactă după marcă, iar după nume doar dacă e unică; refuz cu mesaj când sunt mai mulți candidați; gruparea zilelor pe luni la salvare. |
| SEC-5 | P2 | `llama-server.exe` din `karios/ai/` e lansat fără verificarea sumei de control. Cine poate scrie în folder obține execuție de cod la pornirea aplicației (pornirea automată e implicit activă). | Hash SHA-256 așteptat, fixat la build pentru binarul livrat. |
| SEC-6 | P2 | Data race: `a.aiManager` e inițializat leneș fără sincronizare (`app_ai.go:95-100`), iar `pornire` îl apelează dintr-o gorutină (`app.go:30`) în paralel cu apelurile din UI. Pot rezulta două managere și două procese `llama-server`, unul rămas orfan după închidere. | `sync.Once` sau inițializare în `NewApp`. |
| SEC-7 | P3 | Jurnalul de diagnoză e declarat „fără date personale”, dar scrie `err.Error()` și stack-uri JS complete; erorile AI conțin numele sau marca căutată (`app_ai.go:332`). Jurnalul e necriptat, în `karios/jurnale/`. | Filtrarea/anonimizarea parametrilor sau scrierea doar a codului de eroare. |
| SEC-8 | P3 | `diagnoza.RecupereazaPanic` nu e folosit nicăieri în afara pachetului: o panică într-o metodă legată Wails nu e înregistrată. | `defer diagnoza.RecupereazaPanic(...)` în metodele publice ale `App`, sau un wrapper comun. |

### Review

Fundația criptografică e bine făcută (algoritmi corecți, AEAD cu date asociate, KDF parametrizat, scriere atomică), iar limitarea principală — cheia din cod — e documentată și asumată. Riscul nou și cel mai serios vine din integrarea AI: un server local fără autentificare, pe port fix, căruia i se trimit date personale doar pe baza unui `/health` reușit, e o scurgere GDPR posibilă pe un calculator partajat. Împreună cu potrivirea tolerantă a angajaților (care poate scrie pontaj pe persoana greșită), asistentul ar trebui securizat înainte de folosirea cu date reale.

---

## Recomandare de ordine

1. **SEC-2**, **UX-1**, **UX-2** — riscuri de scurgere de date și de pierdere sau modificare neintenționată a datelor.
2. **SEC-4**, **UX-3**, **UX-4**, **SEC-6** — corectitudinea acțiunilor asistentului.
3. **SEC-3**, **SEC-5**, **UI-4**, **UX-5**, **UX-7** — apărare în adâncime și conformitate cu regulile proiectului.
4. Măsurătoare pe 300 de angajați, apoi **PERF-1** și **PERF-3** după rezultate.
5. **UI-1..3**, **PERF-2/4/5**, **SEC-7/8**.
