# D056 · Doar motoare AI cunoscute: suma llama-server fixată în aplicație

Stare: propusă · Data: 2026-10-30 · Taskuri: F14-T31

## Context

Review-ul din 2026-10-28 (SEC-5): aplicația pornește `karios/ai/llama-server.exe` fără nicio verificare, iar pornirea automată e activă implicit. Cine poate scrie în folder, de exemplu pe un stick cu aplicația portabilă, obține rularea propriului program la pornirea Pontix.

## Decizie

1. **Lista sumelor permise e în cod** (`internal/ai/server_verificat.go`): SHA-256 al fiecărui `llama-server` livrat. Înainte de lansare, aplicația calculează suma binarului; dacă nu e în listă, nu-l pornește și afișează `ai.server_necunoscut`.
2. **Binarul compilat nativ** (`build.ps1 -CuCompilareAI`) se compilează înaintea aplicației, iar suma lui intră în executabil prin `-ldflags -X pontix/internal/ai.HashServerCompilat=…`.
3. **Actualizarea llama.cpp** adaugă suma noului binar în listă, în același commit.
4. **Modelul `.gguf` nu se verifică:** e un fișier de date, nu cod executabil; utilizatorul poate alege alt model. Riscul lui (răspunsuri greșite) e acoperit de confirmarea fiecărei propuneri.

## Variante respinse

- **Suma din `build.json`, de lângă executabil:** cine poate înlocui binarul poate rescrie și manifestul, deci verificarea nu protejează nimic.
- **Prima pornire acceptă orice binar și îi reține suma (TOFU):** un binar înlocuit înainte de prima pornire ar fi acceptat; în plus, ar cere un ecran de confirmare pentru fiecare actualizare.
- **Semnătură digitală a binarului:** cere un certificat de semnare de cod și o infrastructură pe care proiectul nu o are.

## Consecințe

- Un `llama-server.exe` copiat din altă sursă (altă versiune llama.cpp) nu mai pornește până nu i se adaugă suma în cod.
- Rămâne intervalul mic dintre verificare și lansare; cine îl poate exploata poate deja rula programe ca utilizatorul curent.
