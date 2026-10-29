# F14-T24 · Primul lot UI al asistentului · 0.83.0

Codex · 2026-10-29 · I072–I074 din [review](review-ui-asistent-2026-10-29.md). Implementare autorizată de Ioan Tuns după review; D049 rămâne propusă, fără schimbarea stării în numele omului.

## Schimbări

- Istoricul și planul activ au zone separate. Planul folosește lățimea disponibilă și spațiul vertical rămas în fereastră, recalculat la redimensionare. Planurile mai vechi se pot deschide din istoric sau selector.
- Mesajul păstrează același ID la editarea intenției. Efectul de urmărire detectează numai mesaje noi și nu derulează în timpul editării. Derularea se limitează la istoric și se produce numai dacă utilizatorul urmărea capătul.
- Interpretarea, verificarea și rezultatul sunt etape distincte. Trecerea între etape aduce focusul pe titlul etapei, fără a derula pagina; editarea unui câmp nu mută focusul.
- Antetul arată selecția, programul, perioada, zilele și comportamentul pentru celulele existente. Acțiunile sunt într-un subsol separat de conținutul derulabil. Comenzile motorului sunt pliate, iar sugestiile se afișează la începutul conversației.
- Rezultatul afișează valori finale: pentru rândurile salvate folosește valoarea confirmată; pentru rândurile păstrate, identice sau excluse păstrează valoarea anterioară. Stările nu mai descriu rezultatul drept simplă propunere. Salvarea și tokenul din backend nu au fost modificate.
- La revenirea în interpretare se invalidează verificarea și bifa. Nu se aplică automat niciun plan, inclusiv după eroare.

## Verificări

| Probă | Rezultat |
| --- | --- |
| Editare într-un plan mai vechi, între mesaje din istoric | Trecut: aceeași selecție de plan, focus pe același control și poziție neschimbată; fără apel nou de urmărire a istoricului |
| Mesaj nou în timp ce utilizatorul citește începutul istoricului | Trecut: citirea nu este întreruptă |
| Mesaj nou când utilizatorul urmărește capătul | Trecut: un singur apel de derulare în istoric |
| Verificarea unui plan cu 129 de celule | Trecut: contextul și confirmarea rămân vizibile și după derularea tabelului; confirmarea rămâne dezactivată fără bifare |
| Rezultat cu o celulă păstrată, o absență, o celulă identică și o înlocuire | Trecut: valoare finală 6 / — / Z7 / Z7; stări Păstrat / Neaplicat / Deja identic / Salvat. Numărul aplicat este 126 și se deschide luna corectă |
| Token expirat | Trecut: planul nu se aplică; se revine în interpretare cu eroarea și acțiunea de previzualizare vizibile; vechiul token nu mai are confirmare disponibilă |
| Stare inițială fără motor instalat | Trecut: instrucțiunile sunt în conținutul derulabil, introducerea comenzii rămâne accesibilă și nu există un plan gol selectabil |
| Ecran complet RO/EN și lățime îngustă | Trecut: fără derulare a paginii, cu acțiuni vizibile. În proba de 800×640, viewportul raportat este 774×547 |
| Formular izolat, proba existentă RO/EN | Trecut după adaptarea selectorilor la etape și subsol |
| TypeScript și build frontend | `npm --prefix frontend run build` trecut |
| Verificări Go înainte de commit | `gofmt -l -w .`, `go vet ./...`, `go test ./...` trecute |
| Dependențe | Go stabil raportat: 1.27.3; npm outdated fără actualizări; modulele directe și versiunile runtime fixate folosite de aplicație nu raportează actualizări noi. Intrările tranzitive din graful de unelte rămân în F14-T21. Nicio bibliotecă nouă |

Au trecut nouă combinații ale probei complete și două ale formularului izolat. Exemple de dimensiuni CSS raportate de browser:

- viewport 1024×640: plan de 800 px lățime, zona centrală de verificare aproximativ 172 px înălțime, confirmare la y=559–593, în fereastră;
- viewport 998×547: plan de 774 px lățime, confirmare la y=466–500;
- viewport 774×547: plan de 550 px lățime, confirmare la y=466–500;
- viewport 1414×807: plan de 1144 px lățime, zona centrală de verificare aproximativ 321 px înălțime.

PNG-urile au dimensiunile cerute browserului; acestea nu sunt confundate cu `innerWidth` și `innerHeight`. La înălțimi mici se văd puține rânduri simultan, dar contextul și acțiunile rămân fixe, iar conținutul este integral consultabil prin derulare.

## Probe și limite

Fixture-urile sunt în `frontend/probe/asistent-ui.html`, `frontend/probe/asistent-ui.jsx` și proba existentă `frontend/probe/plan-ai.jsx`. Se rulează cu Vite, de exemplu `npm --prefix frontend run dev -- --host 127.0.0.1 --port 5177`, apoi `/probe/asistent-ui.html?caz=rezultat&limba=ro`. Rezultatul este în atributul `data-rezultat`; la succes este `trecut`. Ecranele de probă nu sunt incluse în build-ul de producție.

Capturi, DOM și scriptul local de lansare sunt în `build/bin/proba-ui-F14-T24/`, ignorat de Git: [editare](../build/bin/proba-ui-F14-T24/editare-ro-1024x640.png), [verificare](../build/bin/proba-ui-F14-T24/previzualizare-ro-1024x640.png), [rezultat la viewport 1024×640](../build/bin/proba-ui-F14-T24/rezultat-ro-1050x733.png), [eroare](../build/bin/proba-ui-F14-T24/eroare-ro-1024x640.png), [rezultat EN](../build/bin/proba-ui-F14-T24/rezultat-en-1024x640.png).

Primele încercări ale scriptului local au întâlnit o eroare de sintaxă PowerShell și o încheiere întârziată a browserului/DOM incomplet. Scriptul a fost corectat pentru profil nou per probă și citirea DOM-ului complet, apoi cazurile au fost reluate și rezultatele inspectate.

Probele UI folosesc componentele reale, CSS-ul real și o punte Wails simulată cu persoane fictive. Nu sunt o certificare a executabilului WebView2 la toate scalările Windows, o probă cu operatori sau o măsurare de RAM pe 8 GB. Testele backend existente trec, dar fidelitatea modelului real nu a fost reevaluată aici.

I075–I080 rămân în backlog. În particular, previzualizarea unui plan neconfirmat se reconstruiește după selectarea altui plan; nu este păstrat un token ascuns pentru confirmare ulterioară.
