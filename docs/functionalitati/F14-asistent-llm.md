# F14 · Asistent LLM local

Etapă: Ulterior · Stare: amânat · Depinde de: F07, F08

## Scop

Utilizatorul poate da comenzi în limbaj natural („pune-l pe Ion de noapte joi”). Modelul nu calculează ore sau tarife, iar fiecare modificare se confirmă.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F14-T01 | Setul de 30–50 de comenzi reale pentru validare | amânat | | În română și engleză (D004) |
| F14-T02 | `llama-server` pe `127.0.0.1`, profil de model în `config.toml` | amânat | | |
| F14-T03 | Compararea modelelor disponibile la acel moment | amânat | | Lista din documentul de design se va fi învechit |
| F14-T04 | Interfața Go `Asistent`: JSON impus prin schemă, confirmare, jurnal | amânat | | Fără SQL generat de model |
