# F14 · Asistent LLM local

Etapă: Ulterior · Stare: amânat · Depinde de: F07, F08

## Scop

Utilizatorul poate da comenzi în limbaj natural („pune-l pe Ion de noapte joi”). Conform deciziei ferme [D042](../decizii/D042-asistent-llm-exclusiv-local.md), se folosește **exclusiv un model LLM local** (rulat pe `127.0.0.1`, fără servicii cloud sau transmitere de date în exterior). Modelul nu calculează ore sau tarife, nu generează SQL direct, iar fiecare modificare se confirmă explicit de către utilizator.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F14-T01 | Setul de 30–50 de comenzi reale pentru validare | amânat | | În română și engleză (D004) |
| F14-T02 | `llama-server` pe `127.0.0.1`, supraveghetor portabil și autoverificare mediu | finalizat | Antigravity | Verificare portabilitate, căi karios/ai/, pornire/oprire proces copil fără consolă |
| F14-T03 | Compararea modelelor disponibile la acel moment | amânat | | Qwen 2.5 3B/7B GGUF Q4_K_M recomandat |
| F14-T04 | Interfața Go `Asistent`: client local, JSON impus prin schemă, confirmare explicită | finalizat | Antigravity | Fără SQL direct; fiecare acțiune cere confirmare în UI (v0.36.0) |
