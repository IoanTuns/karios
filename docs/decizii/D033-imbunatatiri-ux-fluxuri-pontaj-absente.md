# D033 · Îmbunătățiri UX în fluxurile de lucru: căutare și navigare lunară în pontaj, selecție prin tragere, precompletare pe rândurile selectate, grupare cereri de concediu și salt contextual

Stare: propusă · Data: 2026-10-26 · Propunere: Antigravity · Legat de: F08-T09, I038, I039, I040, I041, I042, I043

## Context

În urma analizei UX a fluxurilor de lucru din `Pontix 0.25.0`, au fost identificate șase fricțiuni în operarea lunară:

1. **Căutarea angajatului și schimbarea lunii în grila de pontaj (`I038`):** la 50+ angajați lipsea un câmp de căutare rapidă după nume/marcă în `Pontaj`, iar trecerea la luna precedentă sau următoare cerea deschiderea selectorului de lună și schimbarea separată a anului la granița Decembrie/Ianuarie.
2. **Selecția dreptunghiulară cu mouse-ul și încadrarea la derulare din tastatură (`I039`):** selecția multiplă cu mouse-ul funcționa doar cu `Shift+Click`, nu și prin tragere (`Click & Drag`), iar la navigarea cu săgețile pe ecrane mici celula focalizată putea ieși din zona vizibilă fără `scrollIntoView`.
3. **Precompletarea L–V doar pe rândurile selectate (`I040`):** butonul `Precompletează L–V` aplica schimbul pe toate rândurile vizibile, fără o opțiune de a precompleta L–V doar pentru angajații din selecția curentă.
4. **Căutarea în evidența absențelor și selectarea întregii cereri de concediu (`I041`):** în `Concedii și absențe`, selectarea tuturor zilelor unei cereri cu mai multe zile cerea bifarea fiecărui rând individual, și lipsea căutarea text după numărul de înregistrare, cod sau numele înlocuitorului.
5. **Schimbarea de context între `Pontaj` și `Concedii și absențe` (`I042`):** când utilizatorul inspecta o celulă în `Pontaj` și dorea să verifice sau să înregistreze o absență pentru acel angajat, trebuia să comute manual ecranul și să refiltreze anul și angajatul.
6. **Filtrarea rapidă a rândurilor cu avertismente sau modificări în `Pontaj` (`I043`):** înainte de salvarea sau tipărirea lunii, utilizatorul nu putea filtra grila exclusiv pe rândurile cu avertismente (`după G24`, `>24h`, `T ≠ Ore norm`), modificări nesalvate sau înlocuiri.

## Variante

1. **Păstrarea fluxurilor existente** fără scurtături între ecrane sau filtre suplimentare în grilă.
2. **Implementarea celor șase îmbunătățiri UX în ordinea priorității (`I038` → `I043`):**
   - Câmp de căutare după nume/marcă (tolerant la diacritice) în `Pontaj.tsx` și butoane `◀` / `▶` pentru luna precedentă/următoare cu actualizarea automată a anului.
   - Selecție prin tragere cu mouse-ul (`onMouseEnter` când butonul stâng este apăsat) și `scrollIntoView({block: 'nearest', inline: 'nearest'})` la navigarea din tastatură.
   - Buton `Precompletează selecția L–V` (`precompleteazaSelectieLV`) care completează zilele lucrătoare luni–vineri goale doar pentru rândurile cuprinse în selecția curentă.
   - Câmp de căutare text în `Absente.tsx` (după angajat, nr. înregistrare, înlocuitor, cod, aprobator, document) și buton `Selectează cererea` pe rândurile care împart același `document_id` sau `nr_inregistrare`.
   - Legătură contextuală `Deschide în Concedii și absențe` în bara `.inspectie-celula` din `Pontaj.tsx`, care comută în `Absente.tsx` prefiltrat pe anul și angajatul celulei active.
   - Selector de filtrare a rândurilor în `Pontaj.tsx` (`Toți angajații`, `Doar cu avertismente`, `Doar cu modificări nesalvate`, `Doar înlocuitori`).

## Decizie propusă

Se alege **varianta 2** (taskul `F08-T09`). Toate îmbunătățirile operează exclusiv în interfața React, păstrează intacte regulile de calcul și schema bazei de date și sunt acoperite de verificarea automată `verificari/interfata.mjs`.

