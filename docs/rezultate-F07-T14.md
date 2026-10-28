# F07-T14 · Categorii de program și precompletare

26 octombrie 2026 · Codex · Versiune 0.29.0 · Ramură `F07-T14-categorii-program-angajati`.

Utilizatorul a confirmat că lipsea accesul vizibil la **Precompletează L–V**, și a ales modelul **programe permise + implicit, gărzile alese explicit**. D036 este acceptată pe baza acestor răspunsuri.

## Implementare

- Migrare 0016 cu categorii, legături la ID-urile stabile ale schimburilor și atribuire opțională pe angajat; șase categorii inițiale editabile, fără reclasificarea persoanelor existente.
- Administrare în Setări, denumiri RO/EN, implicit opțional, dezactivare, revizii, jurnal și reîncărcare după conflict. Formularul angajatului și lista afișează categoria; Pontajul are filtru dedicat.
- Butonul L–V este vizibil; noul buton După categorii cere propuneri din backend pe un instantaneu tranzacțional al datelor curente și al modificărilor locale.
- Planificator pur în `internal/calcul`, cu calendar, toate absențele active, pontaje efective, norme individuale, ture precedente și suprapuneri. Intervalele neclare sau garda fără alegere explicită sunt omise cu motiv.
- Previzualizare revizuibilă, limitare la angajații din selecție, aplicare în schiță, Undo/Redo pentru lot și salvare explicită. Programele cunoscute din afara categoriei active sunt refuzate în UI și tranzacțional la salvare. Înregistrările istorice neschimbate rămân valide; orele manuale nu sunt transformate în programe deduse.

## Verificări

Testele noi acoperă dimineață/după-amiază/gărzi, norme parțiale, programe fără implicit/inactive/incomplete, CO/CM/tip special, zi liberă personalizată, weekend, trecerea de lună, suprapuneri și intervale adiacente fără suprapunere. Sunt verificate lipsa scrierilor la previzualizare, conflictul categoriilor, păstrarea atribuirilor inactive, refuzul atribuirilor inactive noi și rollbackul lotului dacă un program nu este permis.

Proba UI Edge trece scenariile noi: șase categorii inițiale, creare Mix de test, conflict și reîncărcare, atribuire angajat, filtrare, 60 propuneri pentru întregul lot și 39 pentru doi angajați selectați, cu programe Z7/S7 distincte, protecția absențelor și zilei libere, aplicare fără salvare automată, Undo/Redo, salvare și refuzul programului nepermis. Regresiile existente pentru angajați, absențe, rapoarte, contexte și RO/EN trec. La 1024×640 rămân 7 rânduri complete, cu grila la y=270 și 320 px înălțime.

Capturile finale sunt locale în `build/bin/verificare-ui-2026-10-26T13-38-07-577Z/`: [categorii](../build/bin/verificare-ui-2026-10-26T13-38-07-577Z/categorii-program.png), [previzualizare](../build/bin/verificare-ui-2026-10-26T13-38-07-577Z/precompletare-categorii.png), [grilă compactă](../build/bin/verificare-ui-2026-10-26T13-38-07-577Z/pontaj-compact-1024.png). Puntea Wails este simulată în proba UI; backendul este verificat separat cu Go pe SQLite criptat. Nu se pretinde probă interactivă a ferestrei desktop sau tipărire fizică.

Testele tehnice ale migrărilor își cresc numărătoarea cu migrarea nouă. Proba istorică 0006 verifică organizarea prin SQL pe schema sa istorică, iar proba de ștergere a Z7/S7 își elimină întâi referințele noi din fixture. Nu sunt schimbate rezultate contabile sau fișiere de referință ale rapoartelor.

Verificarea finală `npm run build`, `gofmt -l -w .`, `go vet ./...`, `go test ./...` și `node verificari/interfata.mjs` a trecut. Modulele declarate Go și npm nu au actualizări stabile la verificare; toolchain-ul Go este 1.27.3. Nu s-au adăugat biblioteci.

## Limită descoperită și predare

I054/F04-T09 urmărește o problemă existentă în `pontaj_calcul.go`: citirea sărbătorilor cu an `nil` produce un calendar gol pentru totalurile ZNL/Ore norm. Noua precompletare citește explicit anii necesari și ocolește corect zilele libere. Corecția adaptorului vechi și testele de concordanță grilă/PDF sunt un task separat, conform regulii de scop; nu se declară rezolvate prin această livrare.

Nu s-au deschis baze operaționale și nu s-au copiat date reale. O probă cu operatorul rămâne utilă pentru alegerea categoriilor efective. Atribuirea în lot a categoriilor este idee separată I055; în această versiune atribuirea este individuală, iar precompletarea lunii este pe loturi.

Build final: [Pontix 0.29.0](../build/bin/pontix-0.29.0-20261026-154815-e387bc164751414fa41836b299a0e727/pontix.exe), din commitul `d6c0a05`, construit cu `build.ps1` în director nou. Produsul și versiunea Windows sunt verificate, manifestul este complet, SHA256 `E26B376A7436F17034141C2E1ACDA493ED8E98976FC9FC443222E72A8DD0883E`. Sufixul `-modificat` al commitului din manifest reflectă numai documentul PR modificat în paralel și exclus din task; codul aplicației este comis. Ghidurile sunt copiate lângă executabil. Datele altor versiuni nu sunt transferate automat.
