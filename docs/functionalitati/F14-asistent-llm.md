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
| F14-T07 | Interogarea ca sarcină în fundal: conversația ținută în Go, eveniment Wails la final, buton „Oprește”, notificare pe alt ecran | gata | Claude | 0.61.0: sarcină Go cu eveniment `ai:sarcina_terminata`, oprire, sesiune în `asistentSesiune.ts`, semnal în meniu. Conversația rămâne în frontend (sessionStorage), nu în Go: se pierde la închiderea aplicației, ca înainte. Neverificat cu modelul real |
| F14-T08 | Propuneri de grup pentru „toți angajații / secția X / categoria Y”, recunoscute din cod, cu lista persoanelor și o singură confirmare | gata | Claude | 0.62.0: `internal/ai/grup.go`, `aplicaPropunereGrup`; o zi anume sau un interval cu mai multe schimburi rămân la model. Eroarea salvării planificării e încă ignorată (UX-5 din review) |
| F14-T09 | Răspuns din cod pentru listări („listează / câți / cine”); fără propuneri repetate sau deja aplicate | gata | Claude | 0.63.0: `internal/ai/listari.go`, `faraRepetari` în `asistentSesiune.ts`. O propunere identică cu una încă neconfirmată se omite: se confirmă cardul de mai sus |
| F14-T10 | Luna din contextul ecranului (UX-3) și confirmare doar prin buton (UX-2) | gata | Claude | 0.64.0: luna din pontaj în antetul asistentului, `peLuni` la aplicare, fără confirmare prin text. Planificarea nu raportează încă luna înapoi la `App` (se folosește doar contextul pontajului) |
| F14-T11 | Modificările aplicate de asistent apar în jurnal ca „asistent AI, confirmat de utilizator” | gata | Claude | 0.67.0: `internal/datebaza/actor.go`; `salveazaPropunere` marchează contextul. Evenimentele de audit (inițializare, migrare, incident) rămân `local` |
| F14-T12 | Serverul local cu cheie API aleatoare la fiecare pornire și verificarea procesului propriu (SEC-2) | gata | Claude | 0.72.0: port liber + `--api-key` + proces propriu urmărit; include SEC-3 (numele modelului). Neverificat cu binarul real |
| F14-T13 | Identificarea angajatului doar exact (marcă) sau după nume unic; refuz la mai mulți candidați (SEC-4) | gata | Claude | 0.73.0: `identificaAngajat` în `app_ai.go`. Afișarea din cardul propunerii folosește încă potrivirea tolerantă din interfață (doar afișare; aplicarea e exactă) |
| F14-T14 | Inițializarea sincronizată a managerului AI, fără procese `llama-server` duble (SEC-6) | gata | Claude | 0.74.0: `sync.Once`; detectorul de curse (`-race`) neverificat, fără gcc |
| F14-T15 | Recunoașterea comenzilor cu greșeli de tastare și afișarea interpretării | gata | Claude | 0.77.0: `corectura.go`, `esteToti`, `lunaDin`; interpretarea pe card. „fără pontaj” e F14-T16 |
| F14-T16 | Filtrul „fără pontaj” în comenzile de grup și completarea doar a celulelor goale | gata | Claude | 0.78.0: filtrul pe luna vizată; `doar_goale` implicit adevărat pentru propunerile de grup |
| F14-T17 | Istoricul trimis modelului fără text de imitat | gata | Claude | 0.79.0: `actiuni` în cerere, secțiune în promptul de sistem; neverificat cu modelul real |
| F14-T18 | Detectarea răspunsurilor care pretind o planificare fără nicio acțiune | gata | Claude | 0.80.0: listă de expresii de confirmare; răspunsurile analitice din cod nu trec prin verificare |
| F14-T19 | Propunerile identice cu unele deja aplicate se arată, cu „deja aplicată, aplică din nou?” | gata | Claude | 0.81.0: în `asistentSesiune.ts`; fără test automat (frontend fără teste unitare) |
