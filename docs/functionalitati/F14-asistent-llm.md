# F14 · Asistent LLM local

Etapă: Ulterior · Stare: în lucru · Depinde de: F07, F08

## Scop

Utilizatorul poate da comenzi în limbaj natural („pune-l pe Ion de noapte joi”). Conform deciziei ferme [D042](../decizii/D042-asistent-llm-exclusiv-local.md), se folosește **exclusiv un model LLM local** (rulat pe `127.0.0.1`, fără servicii cloud sau transmitere de date în exterior). Modelul nu calculează ore sau tarife, nu generează SQL direct, iar fiecare modificare se confirmă explicit de către utilizator.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F14-T01 | Setul de 30–50 de comenzi reale pentru validare | amânat | | În română și engleză (D004) |
| F14-T02 | `llama-server` pe `127.0.0.1`, supraveghetor portabil și autoverificare mediu | finalizat | Antigravity | Verificare portabilitate, căi karios/ai/, pornire/oprire proces copil fără consolă (v0.36.0) |
| F14-T03 | Compararea modelelor disponibile la acel moment | amânat | | Qwen 2.5 3B/7B GGUF Q4_K_M recomandat |
| F14-T04 | Interfața Go `Asistent`: client local, JSON impus prin schemă, confirmare explicită | finalizat | Antigravity | Fără SQL direct; fiecare acțiune cere confirmare în UI (v0.36.0) |
| F14-T05 | Pachet portabil AI Ready și pornire automată supravegheată | finalizat | Antigravity | Script creaza-pachet-ai.ps1, karios/ai/ pre-populat, opțiune de pornire automată la lansare (v0.37.0) |
| F14-T06 | Compilare nativă llama.cpp pe Windows cu toolchain portabil și integrare deploy | finalizat | Antigravity | GCC 16 + CMake portabil în tools/, scripts/compileaza-llama.ps1, build.ps1 -CuCompilareAI (v0.38.0) |
| F14-T07 | Interogarea ca sarcină în fundal: conversația ținută în Go, eveniment Wails la final, buton „Oprește”, notificare pe alt ecran | de făcut | | Ioan Tuns, 2026-10-29: azi răspunsul se pierde la schimbarea ecranului. Propunerile nu se aplică niciodată singure |
| F14-T08 | Propuneri de grup pentru „toți angajații / secția X / categoria Y”, recunoscute din cod, cu lista persoanelor și o singură confirmare | de făcut | | Captura 2026-10-29: „toți angajații” s-a aplicat doar unui angajat |
| F14-T09 | Răspuns din cod pentru listări („listează / câți / cine”); fără propuneri repetate sau deja aplicate | de făcut | | Captura 2026-10-29: „Listează toți angajații” a reprodus o propunere aplicată |
| F14-T10 | Luna din contextul ecranului (UX-3) și confirmare doar prin buton (UX-2) | de făcut | | [review](../review-ui-ux-performanta-securitate-2026-10-28.md) |
| F14-T11 | Modificările aplicate de asistent apar în jurnal ca „asistent AI, confirmat de utilizator” | de făcut | | Trasabilitate, D044 |
| F14-T12 | Serverul local cu cheie API aleatoare la fiecare pornire și verificarea procesului propriu (SEC-2) | de făcut | | [review](../review-ui-ux-performanta-securitate-2026-10-28.md) |
| F14-T13 | Identificarea angajatului doar exact (marcă) sau după nume unic; refuz la mai mulți candidați (SEC-4) | de făcut | | [review](../review-ui-ux-performanta-securitate-2026-10-28.md) |
| F14-T14 | Inițializarea sincronizată a managerului AI, fără procese `llama-server` duble (SEC-6) | de făcut | | [review](../review-ui-ux-performanta-securitate-2026-10-28.md) |
