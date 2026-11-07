# D091 · Folderul AI separat de folderul de date

Stare: acceptată (cerința) · propusă (detaliile tehnice) · Data: 2026-11-07 · Taskuri: F14-T36

Decisă de: Ioan Tuns, 2026-11-07, pentru: folderul `ai` scos din `karios/` și pus la același nivel cu el, ca baza de date și LLM-ul să fie separate.

## Context

Din D042 și F14, motorul local (`llama-server.exe`) și modelele `.gguf` stau în `karios/ai/`, adică în același folder cu baza criptată, seiful, copiile și jurnalele (D002). Modelele au câțiva GB. Folderul `karios/` se copiază la mutarea sau actualizarea aplicației și se tratează ca date sensibile. Amestecat cu ele, LLM-ul îngreunează copierea, copiile de rezervă făcute manual și înțelegerea a ce conține folderul de date.

## Decizie

1. **Așezarea:** lângă `pontix.exe` stau două foldere: `karios/` pentru datele operaționale (bază, seif, blocare, copii, WebView2, jurnale) și `ai/` pentru motorul și modelele LLM, plus fișierul de dezactivare a pornirii automate. Aplicația nu scrie nimic în afara acestor două foldere. Aceasta modifică regula din [D002](D002-portabilitate.md) („un singur folder de date”) doar pentru componenta AI.
2. **Mutarea automată:** la pornire, dacă există `karios/ai/` și nu există încă `ai/`, folderul se mută printr-o redenumire pe același disc (fără copierea modelelor). Dacă redenumirea nu reușește, de exemplu pentru că motorul rulează din vechiul folder într-o altă instanță, aplicația folosește în continuare `karios/ai/` și reîncearcă la următoarea pornire. Dacă există ambele foldere, nu se amestecă nimic: se folosește `ai/`, iar `karios/ai/` rămâne neatins.
3. **Livrarea:** `build.ps1` pune `llama-server.exe` verificat în `ai/` (modelul tot nu se include automat).
4. **Verificarea motorului** (D056) rămâne neschimbată; se aplică noului loc.

## Variante respinse

- **Copierea în loc de redenumire:** dublează temporar câțiva GB și poate eșua la jumătate; redenumirea e atomică pe același disc.
- **Folder `karios-ai/`:** utilizatorul a cerut `ai` la același nivel cu `karios`.
- **Mutarea și a WebView2 sau a jurnalelor:** nu au fost cerute; rămân date ale aplicației în `karios/`.

## Consecințe

- La mutarea aplicației pe alt calculator, se copiază `karios/` pentru date și, dacă se folosește asistentul, `ai/`. Manualul spune asta.
- O copie a folderului `karios/` nu mai include modelele.
- Cine poate scrie în `ai/` poate înlocui motorul; protecția rămâne verificarea din D056.
