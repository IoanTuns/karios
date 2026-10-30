# F02-T21 · Verificare înainte de PR

Data: 2026-10-30. Responsabil: Codex. Verificare solicitată de utilizator după actualizarea main.

Referințe: origin/main f15e623 (PR #2, include 632938e), ajutor F02-T21-ajutor-contextual 9b60d87, bază comună 12470c3. Ramurile au câte două commituri proprii. Main local și codul au rămas neatinse; verificarea folosește git merge-tree --write-tree, fără merge în index.

## Rezultat: integrarea nu este pregătită pentru PR

### P1 · Patru conflicte de merge

- docs/jurnal.md: intrările noi trebuie păstrate din ambele ramuri.
- docs/manual-utilizare-ro.md: se reunesc descrierile optimizărilor Pontaj și ale ajutorului. Nu se elimină instrucțiunile pentru clipboard, meniul contextual, heatmap sau drawer și nici capitolul 20 cu articole comune.
- docs/versiuni.md: ramurile au publicat conținut diferit sub aceeași versiune 0.88.13.
- frontend/src/style.css: ambele ramuri adaugă stiluri la sfârșit. Se păstrează atât stilurile drawerului/meniului/heatmapului, cât și cele ale panoului de ajutor, fără alegerea integrală a unei singure variante.

### P1 · Coliziune de versiune care nu apare ca un conflict Git

internal/versiune/versiune.go, wails.json și docs/manual/manual.json au automat același rezultat 0.88.13 în ambele ramuri. Un merge simplu ar ascunde faptul că 0.88.13 din main nu include ajutorul. Integrarea trebuie să devină 0.88.14 (dacă între timp nu apare altă versiune), cu istoric separat și manual revizuit pentru toate funcțiile. Ediția manualului se crește la 1.3; se regenerează documentele și executabilul împreună.

### P2 · Citările de ajutor nu funcționează în noul drawer AI

origin/main adaugă în frontend/src/Pontaj.tsx o instanță AsistentAI fără onAjutor. Componenta mea frontend/src/AsistentAI.tsx afișează butoanele de citare, dar handlerul onAjutor?.(...) nu face nimic în acea instanță. Nu există eroare TypeScript: callbackul este opțional. În ecranul principal AI funcționează, fiind transmis de App.tsx.

Înainte de PR se transmite callbackul (inclusiv versiunea citată) din App prin Pontaj la AsistentAI sau se oferă un mecanism comun echivalent. Se verifică și protecția modificărilor nesalvate la navigarea din ghid.

## Ce se păstrează fără suprascriere

- Blobul frontend/src/Pontaj.tsx din rezultatul simulat este identic byte cu byte cu origin/main: optimizările noi nu sunt suprascrise de ramura mea.
- JSON-urile RO și EN se reunesc valid. Pentru fiecare limbă se păstrează toate cele 11 chei noi/modificate din main și toate cele 142 din ajutor, fără valori pierdute.
- Nu sunt schimbări concurente în App.tsx sau AsistentAI.tsx din main; problema drawerului este funcțională, nu un conflict de text.

## Verificări și pașii necesari

Au fost inspectate diferențele și arborele de merge simulat; JSON-urile și păstrarea Pontajului au fost verificate programatic. Nu s-a compilat un rezultat cu marcaje de conflict și nu s-au declarat drept valide testele vechi pentru noua combinație.

Pași următori: rezolvarea celor patru conflicte păstrând ambele funcționalități; versiunea nouă și revizuirea manualului; conectarea citărilor în drawer; Go/vet, frontend, probe UI pentru drawer și pentru funcțiile noi din Pontaj; build și manuale regenerate din același commit. Nu s-au făcut merge, push sau PR în această verificare.

## Remediere pentru PR

Integrarea pe F02-T21-ajutor-pr pastreaza ambele blocuri CSS si intrarile de jurnal, separa istoricul 0.88.13/0.88.14 si revizuieste manualul la editia 1.3. Callbackul App -> Pontaj -> AsistentAI permite deschiderea ghidurilor din drawer. Cheia inexistenta nav.asistent_ai este inlocuita cu ui.asistent_ai. Scenariile UI RO/EN verifica explicit citarea din drawer.

Verificari: frontend build si 13 teste Vitest, 3 teste manual si scenarii UI RO/light + EN/dark 1024x640 trecute. go vet trecut; suita Go completa trecuta. Dependente: Go 1.27.3 si npm curente; actualizarile Go tranzitive raman urmarite separat in F14-T21.
