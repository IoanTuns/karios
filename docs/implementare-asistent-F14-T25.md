# F14-T25 · Instrumente pentru asistent

Responsabil: Codex. Stare: în lucru. Cerere: Ioan Tuns, implementarea tuturor celor opt propuneri. Ramură separată, fără push sau merge.

## Criterii

1. Anularea ultimei operații de planificare AI, disponibilă și după repornire: previzualizare, confirmare, verificarea reviziilor și lunilor, rollback integral la conflict. Se restaurează numai celulele și pozițiile schimbate de operație; jurnalul rămâne. Nu se anulează operațiile istorice fără instantanee.
2. Selecția explicată pentru fiecare angajat, inclus/exclus, cu criteriu de grup și excepții pe zile (absență, categorie, pontaj existent).
3. Clarificări pentru câmpuri lipsă și confirmarea anului preluat din context; date și ordine validate înainte de previzualizare.
4. Sumar pe lună și persoană, distinct selectat/afectat, păstrat/blocat/identic; căutare și filtrare a detaliilor.
5. Șabloane persistente, denumite de utilizator, încărcate în editor cu perioadă și parametri editabili; niciun șablon nu aplică direct modificări.
6. Verificarea pontajului pentru perioada și zilele alese: goluri, coduri invalide, absențe, categorie și suprapuneri conform verificărilor existente. Fără reguli salariale noi și fără corectare automată.
7. Modelul efectiv pornit și contextul configurat, fără afișarea cheii API sau a consumului inventat; resetarea conversației explicată.
8. Export local CSV al propunerii reconstruite de backend și verificate prin token, cu sumar, selecție, valori/stări; rezultatul aplicat exportat din instantaneul operației. Dialog de salvare, fără transmitere externă, celule protejate împotriva formulelor.

## Verificări

Teste backend pentru anulare, conflicte, persistență, audit și atomicitate, selecție/sumar, șabloane și verificări. Regresii de interpretare și probe UI RO/EN cu date fictive. Verificări Go, frontend, dependențe și build oficial. Limitări de mediu declarate în predare.

## Implementare

Versiunea 0.84.0. Migrarea 0023 adaugă `operatii_ai` și `sabloane_ai`. Planul salvat și instantaneele sunt scrise în aceeași tranzacție cu aplicarea. Anularea compară reviziile și instantaneele, reverifică restricțiile curente, reconstruiește tokenul și restaurează atomic, păstrând auditul și reviziile crescătoare. Celula inițial absentă devine goală; poziția nou creată devine inactivă. Operațiile istorice fără instantanee nu sunt anulabile. Nu există redo.

Explicațiile și rezumatele sunt construite în backend și incluse în token. Șabloanele nu conțin cererea liberă sau tokenuri; încărcarea confirmă din nou perioada. Verificarea folosește intervalele reale, nu protecțiile de zi folosite pentru a interzice scrierea peste pontaj efectiv. Regulile configurabile existente se evaluează pe schița lunii complete pentru persoanele selectate; nicio regulă salarială nu este adăugată. Exportul este local și verificat din sursa backend.

Interfața are „Plan nou / verificare” și „Instrumente”, fără dependență de motorul LLM pentru funcțiile deterministe. [Ghidul](ghid-asistent-ai.md) descrie toate cele opt funcții și limitele.

## Probe UI reproductibile

Pornește Vite cu `npm --prefix frontend run dev -- --host 127.0.0.1 --port 5177`, apoi rulează `frontend/probe/ruleaza-instrumente.ps1`. Nu instalează nimic, folosește Edge existent în mod headless, profile temporare separate și date fictive. Probe: clarificări, verificare, model, șabloane și selecție nominală, sumar/explicații/filtre/export, anulare, conflict și token de anulare expirat. Matrice RO/EN, ferestre 1024×733 și 800×733. DOM și capturi: `build/bin/proba-ui-F14-T25/`, ignorate de Git. Probele vechi de editare/rezultat rămân în `frontend/probe/asistent-ui.html`.

Toate cele 24 de combinații finale au trecut: 12 pentru funcțiile noi, inclusiv expirarea anulării, ștergerea șablonului, afișarea planificării separate și traducerea regulilor configurate; 12 regresii pentru editare/focus, istoric, rezultat, eroare, previzualizare și motor absent. Scriptul a necesitat citire partajată și așteptarea DOM-ului complet, deoarece procesul browserului poate păstra deschis fișierul redirecționat. Au fost inspectate vizual rezultatul EN pe fereastră mică și conflictul RO.

Testele de migrare schimbă numai numerele așteptate pentru schema nouă. Testele AI existente confirmă acum explicit anul; nu sunt schimbate teste de calcul sau cazuri confirmate de contabilitate. O probă nouă de suprapunere a necesitat un interval efectiv intersectat real; apoi a detectat că protecția zilei de scriere producea un fals pozitiv la ture adiacente. Verificarea consultativă a fost corectată să ignore această protecție sintetică, păstrând intervalele efective și protecția originală la scriere.

Go stabil verificat: 1.27.3. `npm outdated` nu raportează actualizări. Dependențele runtime sunt curente; actualizările tranzitive ale uneltelor rămân în taskul separat F14-T21. Nu există biblioteci noi sau licențe schimbate. `gofmt -l -w .`, `go vet ./...`, `go test ./...`, compilarea TypeScript/Vite și `git diff --check` au trecut. Rulările concurente Go de durată au fost reduse, una a fost oprită, alta s-a încheiat cu succes; regresiile noi au trecut și izolat cu timeout. Reluarea finală folosește `go test -p 1 -timeout 120s ./...` și dependențele locale, fără schimbarea versiunilor. Build-ul se consemnează la predare.

Nu s-au rulat probe cu modelul real, operatori, WebView2 interactiv sau măsurători de RAM pe 8 GB. Nu se pretinde că verificarea pontajului certifică legislația ori regulile salariale. Exportul CSV este testat ca date și token; dialogul nativ de salvare necesită probă umană.
