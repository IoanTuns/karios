# D043 · Sistem de diagnoză locală, captare a erorilor și profilare a interogărilor

Stare: acceptată · Decisă de: Ioan Tuns · Data: 2026-10-28 · Taskuri: F15-T01

## Context

Pontix funcționează exclusiv local și offline (D042), fără servicii externe de telemetrie, monitorizare sau raportare automată a erorilor (cum ar fi Sentry sau Datadog). 

În cazul unei excepții JavaScript neprinse pe frontend, a unei panici în Go sau a unei interogări lente SQLite / calcul intens, sistemul nu dispunea până acum de un mecanism integrat de captare, stocare și analiză persistentă. Diagnosticul devenea anevoios, iar utilizatorii nu aveau cum să furnizeze date tehnice precise echipei de asistență.

## Decizie

1. **Captare erori pe frontend (React & WebView2):**
   - Interceptare globală prin `window.addEventListener('error')` și `window.addEventListener('unhandledrejection')`.
   - Înfășurarea componentelor majore într-un `ErrorBoundary` React pentru a preveni înghețarea întregii interfețe la erori locale de randare.
   - Trimiterea automată a excepțiilor frontend către backend prin metoda Wails `InregistreazaEroareFrontend(sursa, mesaj, stack)`.

2. **Jurnalizare persistentă și rotativă în backend (`internal/diagnoza`):**
   - Jurnalele se salvează exclusiv pe disc local în directorul portabil de date: `karios/jurnale/pontix-AAAA-LL.log`.
   - Dimensiunea fișierului de log este limitată la maximum 5 MB, cu rotire automată și păstrarea a maximum 3 arhive lunare pentru a preveni umplerea spațiului pe disc.
   - **Confidențialitate strictă (GDPR / Regula 4 AGENTS.md):** jurnalele înregistrează doar timestamp, severitate, modul, fișier, linie, mesaj de eroare tehnic și stivă de execuție (`stack trace`), **fără** date cu caracter personal (fără nume, CNP-uri, adrese sau salarii).

3. **Măsurarea timpului de execuție pentru interogări (Query Profiling):**
   - Monitorizarea duratei de execuție (`time.Since(start)`) pentru interogările majore SQLite și operațiile de calcul grele.
   - Definirea unui prag de avertisment (*Slow Query*, ex: > 50ms pentru SQL, > 500ms pentru operații agregate).
   - Menținerea în memorie a ultimelor N metrici de latență (minim, mediu, maxim) expuse prin `RaportDiagnoza()`.

4. **Interfață și suport tehnic:**
   - Integrarea unui panou de diagnoză în fereastra modală *Detalii dezvoltare, integritate și cerințe sistem* din Setări.
   - Butoane rapide pentru utilizator: `📁 Deschide folderul de jurnale` și `📋 Copiază logul de diagnoză`.

## Variante respinse

- **Telemetrie sau logare în cloud** — Respinsă: încalcă cerința fundamentală de izolare totală offline (D042) și confidențialitatea medicală GDPR.
- **Salvarea erorilor exclusiv în consolă (`console.error` / `log.Println`)** — Respinsă: pe Windows, aplicația compilată cu `-H windowsgui` nu are consolă atașată, iar mesajele se pierd la închiderea ferestrei.

## Consecințe

- Orice eroare tehnică poate fi depistată imediat prin inspectarea fișierului local din `karios/jurnale/`.
- Utilizatorii pot exporta facil un raport tehnic complet în caz de asistență.
- Interogările lente sunt semnalizate prompt pentru optimizare.

