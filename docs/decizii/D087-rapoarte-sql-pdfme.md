# D087 · Toate rapoartele din SQL, cu editor vizual pdfme

Stare: decisă (direcția, motorul PDF, conversia modelelor) · propusă (detaliile tehnice) · Data: 2026-11-05 · Taskuri: F09-T20, F09-T21, F09-T22, F09-T23

Decisă de: Ioan Tuns, 2026-11-05, pentru: folosirea pdfme, trecerea tuturor rapoartelor pe SQL, generarea PDF cu pdfme în aplicație și conversia automată a modelelor salvate. Restul secțiunii „Decizie” rămâne propus până la acceptare.

## Context

Editorul de rapoarte (D025, D082, D085) cere SQL scris de mână, Markdown într-o căsuță de text și o planșă SVG unde elementele se pot doar muta. Utilizatorul îl găsește greu de folosit și cere un editor vizual open-source. Mai cere ca rapoartele să se genereze din SQL, de preferat toate, ca să poată fi schimbate pe viitor fără cod nou.

Astăzi există două motoare:

- rapoartele standard (`PONTAJ_LUNAR`, `PONTAJ_DETALIAT`, `FISA_INDIVIDUALA`, `ABSENTE_ANUALE`, `GRAFIC_GARZI`), construite în Go (`Serviciu.Pregateste`) și desenate cu maroto;
- rapoartele flexibile (`SQL_LIBER`, `LISTA_GARZI`, `CALCULATE_LUNAR`), cu SQL pe un instantaneu SQLite izolat (D082) și PDF desenat cu gofpdf.

## Decizie

1. **Editorul** este Designer-ul pdfme (`@pdfme/ui`), integrat în ecranul de definire a rapoartelor. Șablonul pdfme (JSON, coordonate în mm) se salvează în definiția raportului, cu revizie și jurnal, ca până acum.
2. **Datele vin numai din SQL.** Fiecare raport, inclusiv cele standard, are o interogare pe instantaneul izolat din D082, cu aceleași limite și același authorizer. Rapoartele standard devin modele predefinite, livrate cu aplicația, pe care utilizatorul le poate duplica și modifica.
3. **Calculele rămân în Go.** Norma, orele lucrate, orele eligibile, concediile care reduc obligația și orice regulă de plată se calculează în `internal/calcul` și apar în instantaneu ca tabele gata calculate (de exemplu `lunar`, `pontaj_zile`). SQL-ul doar selectează, ordonează, grupează și afișează. Regulile de plată nu se mută în SQL și nu se pot schimba din editor. Testele de referință ale calculului rămân neschimbate.
4. **PDF-ul îl generează pdfme în aplicație** (`@pdfme/generator`), din șablon și din rezultatul SQL, deci arată ca în editor. Fonturile Go Regular/Bold/Italic/BoldItalic, deja folosite, se încarcă în pdfme pentru diacritice. Go primește PDF-ul, îl validează cu pdfcpu (deja în proiect) și adaugă anexa cu modul de calcul și istoricul tehnic (D009), care nu se pot scoate din editor.
5. **Elemente obligatorii:** subsolul tehnic D009 (versiunea raportului, revizia, aplicația, momentul, actorul) și numerotarea paginilor sunt elemente fixe, adăugate de aplicație, nu șterse din editor.
6. **Excel** se exportă din același rezultat SQL, ca acum (excelize), nu din șablonul pdfme.
7. **Conversia modelelor salvate** se face automat, o singură dată, prin cod de conversie versionat: fiecare model devine SQL + șablon pdfme, cu textele, stilurile, etichetele și imaginile lui. Conversia creează o revizie nouă, cu jurnal, și păstrează originalul în istoricul reviziilor. Motorul vechi (maroto și desenarea gofpdf) se șterge abia după ce rapoartele convertite dau același conținut pe date anonime.
8. **Ce nu are pdfme și trebuie adăugat** printr-un plugin propriu: stil pe rând și pe celulă (bold la weekend și sărbători, celulele editate din D085) și tabelul pontajului detaliat, cu coloane pe zilele lunii (28 până la 31).

## Biblioteci și licențe

Versiuni stabile la 2026-11-05, fixate exact în `package.json` și `package-lock.json`:

- `@pdfme/ui`, `@pdfme/generator`, `@pdfme/common`, `@pdfme/schemas` **6.2.8**, licență **MIT** (din registrul npm). Sursa: <https://github.com/pdfme/pdfme>.
- Dependențele de producție aduse: 208 MIT, plus `@swc/helpers` Apache-2.0, `dompurify` MPL-2.0 sau Apache-2.0, `ieee754` și `jpeg-js` BSD-3-Clause, `lucide`/`lucide-react` ISC, `pako` MIT și Zlib, `tslib` 0BSD. Toate sunt permisive și compatibile cu distribuția aplicației. Notificările de licență intră în livrare, ca la celelalte biblioteci.
- Nu se adaugă nimic în Go: pdfcpu, ncruces/go-sqlite3 și excelize sunt deja în proiect.

## Implementare (0.92.33-0.92.42)

Punctele 1-8 sunt implementate (F09-T20...F09-T25). Conversia (punctul 7) ruleaza la primul acces la rapoarte dupa autentificare, nu la pornire, pentru ca scrierile cer autentificare (D045) si jurnalul trebuie sa poarte utilizatorul. La cererea utilizatorului, si modelele Markdown/Pozitionare libera se convertesc aproximativ, iar motorul vechi s-a sters in acelasi task, fara proba pe date reale inainte. Validarea in Go a incadrarii in pagina a disparut odata cu motorul vechi.

## Variante respinse

- **Go desenează șablonul pdfme** (gofpdf): totul ar rămâne în Go, dar ar trebui reimplementată logica pdfme (tabele paginate, text dinamic), iar PDF-ul ar putea diferi de editor. Respinsă de utilizator.
- **Modele noi alături de cele vechi**, fără conversie: mai puțin risc, dar două motoare o perioadă lungă. Respinsă de utilizator.
- **GrapesJS** (HTML/CSS): are nevoie de un browser pentru PDF, respins deja în D082.
- **JasperReports, FastReport, BIRT**: Java sau .NET, prea grele pentru o aplicație portabilă.
- **Carbone, docxtemplater**: module plătite sau licențe care nu sunt OSI, iar pentru PDF le trebuie LibreOffice.
- **Calculul în SQL**: ar muta regulile de plată din `internal/calcul`, în afara testelor confirmate de contabilitate.

## Consecințe

- Utilizatorul editează vizual: trage, redimensionează, aliniază cu ghidaje, are panou de proprietăți și previzualizare identică cu PDF-ul.
- Un raport nou nu mai cere cod Go, doar SQL și șablon, cât timp datele există în instantaneu. Date noi în instantaneu cer în continuare cod Go și teste.
- PDF-ul se generează în interfață. O generare fără interfață (de exemplu din asistentul AI) trebuie să treacă prin aceeași interfață sau să rămână doar pe Excel.
- Testele cu fișier de referință ale rapoartelor se mută în două straturi: rezultatul SQL (Go) și PDF-ul pdfme (vitest). Fiecare raport convertit primește o versiune nouă (D009).
- Pachetul interfeței crește (antd, pdf.js). Se măsoară la build.
- Manualul și ghidurile rapoartelor (`ghid-rapoarte.md`, `ghid-rapoarte-flexibile.md`) se rescriu.
