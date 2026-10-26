# D042 · Asistent LLM exclusiv local (fără servicii cloud sau API externe)

Stare: acceptată · Data: 2026-10-26 · Propusă de: Antigravity · Decisă de: Ioan Tuns · Task: F14

## Context

În cadrul discuției privind integrarea capacităților de asistență AI în Pontix (comenzi în limbaj natural pentru planificare, inspecție pontaj, interogări orare), utilizatorul a stabilit o cerință fermă și obligatorie:
> *„Se folosește doar model local! Notează asta”*

Aplicația gestionează date cu caracter personal (nume, prenume, mărci, concedii medicale, program de lucru, gărzi în unități sanitare). Orice apel către un API cloud extern ar presupune transmiterea datelor prin rețea, ridicând riscuri de securitate și probleme de conformitate GDPR.

## Decizie

1. **Exclusivitate locală (Offline):**
   - Nicio funcționalitate AI nu va transmite date către servicii externe sau furnizori cloud (exclus OpenAI, Anthropic, Google Cloud API etc.).
   - Modelul va rula **strict local** pe calculatorul utilizatorului (pe `127.0.0.1`), prin intermediul unui server de inferență integrat sau local (ex: `llama-server` / `llama.cpp` sau `Ollama`).

2. **Arhitectură și cerințe tehnice:**
   - Modele compatibile: modele compacte cuantizate la 4 biți (format GGUF Q4_K_M, 3B–7B parametri, ex: `Qwen 2.5 Instruct`, `Llama 3.2 Instruct`), capabile să ruleze eficient pe CPU sau GPU fără conexiune la internet.
   - Punct de integrare: endpoint local standard OpenAI-compatible (`http://127.0.0.1:<port>/v1/chat/completions`).

3. **Integritate și control:**
   - Modelul emite strict apeluri de funcții (Function Calling / JSON tipizat cu schemă impusă din Go).
   - Fără generare directă de cod SQL.
   - Orice modificare a planului sau pontajului propusă de modelul local necesită confirmarea explicită a utilizatorului înainte de salvare.
