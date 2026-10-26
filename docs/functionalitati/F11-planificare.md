# F11 · Planificare și planificat vs. realizat

Etapă: V1.1 · Stare: în lucru · Depinde de: F08

## Scop

Graficul planificat se face separat de ce s-a lucrat efectiv, iar diferențele se văd într-un raport.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F11-T01 | Planificare separată pentru gărzi și după-amiază, cu publicare în pontaj | gata | Codex, Antigravity | Cerere explicită a utilizatorului după F07-T14. Ramură F11-T01-planificare-publicare, I056/D037. Planul se salvează independent în tabela planificare (migrarea 0017); publicare explicită cu previzualizare, calcul de ore și protecția absențelor, programelor permise și înregistrărilor existente. Traduceri complete RO/EN, teste trecute. |
| F11-T03 | Grilă lunară tabelară de planificare (pereche zi – program) și publicare unificată | gata | Antigravity | Cerință confirmată prin I057 și D038. Grilă lunară cu angajați pe rânduri și zilele 1–31 pe coloane, atribuire directă pentru S7, IN9, G24, Z12, N12, evidențiere weekend/sărbători/absențe, Undo/Redo, validare categorii și previzualizare/publicare unificată. Versiunea 0.31.0. |
| F11-T04 | Sumar editabil pe angajat pentru gărzi și ture | gata | Antigravity | Cerință confirmată prin I058 și D039. Dialog modal accesibil din grilă și vedere dedicată cu totaluri ore/ture (G24, Z12, N12, IN9, S7) și adăugare/editare/eliminare rapidă pe zile. Versiunea 0.32.0. |
| F11-T05 | Raport de acoperire weekend și zile nelucrătoare (G24, Z12, N12, IN9) | gata | Antigravity | Cerință confirmată prin I059 și D039. Verificare automată acoperire pe sâmbete, duminici și sărbători legale, badge-uri alocare/lipsă, carduri sintetice și filtru doar deficit. Versiunea 0.32.0. |
| F11-T06 | Verificare echitate sărbători (Crăciun/Revelion/Paști) și marcaje ture consecutive | gata | Antigravity | Cerință confirmată prin I060 și D039. Corelare rotație Crăciun/Revelion/Paști între ani și marcaje vizuale neblocante (celula-atentie, tooltip, inspecție poziție) pentru ture consecutive. Versiunea 0.32.0. |
| F11-T02 | Raportul planificat vs. realizat | de făcut | | |
