# D104 · Rapoartele cu o pagină nouă pentru fiecare secție

Stare: propusă · Data: 2026-11-12 · Taskuri: F09-T02 · Cerută de: Ioan Tuns, 2026-11-12 („Poți să le implementezi pe toate și să construiești o versiune nouă?”)

Implementare: implementată din 0.92.97; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

## Context

F09-T02 cerea pentru centralizatorul lunar „o pagină pe secție”. Celelalte cerințe (A4/A3 peisaj, antetul tabelului repetat) existau deja. Rapoartele se desenează cu pdfme în interfață (D087), ca un singur tabel care curge pe pagini. pdfme nu are un salt de pagină forțat în mijlocul tabelului.

## Decizie

- **Opțiunea.** În ecranul Rapoarte apare bifa **O pagină nouă pentru fiecare secție**, doar când:
  - raportul e standard: nu e SQL personalizat și nu e fișa individuală;
  - sunt alese toate secțiile, toate departamentele și toți angajații;
  - angajații din raport sunt în cel puțin două secții.
- **Generarea.** Interfața generează raportul pe rând, pentru fiecare secție, exact ca la alegerea manuală a secției.
  - Secțiile vin în ordinea din Setări.
  - Angajații fără secție vin la final, ca o parte separată, cu numele secției „Fără secție”. Valoarea tehnică a cererii este `SectieID = "-"`.
- **Unirea.** Go unește PDF-urile cu pdfcpu (`api.MergeRaw`), care e deja în proiect. Fiecare secție începe pe pagină nouă. Rezultatul se previzualizează, se salvează și se tipărește ca orice raport.
- **Ce rămâne la fel.** Fiecare parte are antetul, subsolul și numerotarea ei: „Pagina 1 din 2” se reia la fiecare secție. Conținutul raportului nu se schimbă, deci versiunile rapoartelor rămân aceleași.
- **Exportul Excel** al rezultatului unit conține toți angajații, pe o singură foaie, ca raportul fără secție.
- **Editarea copiei.** Rezultatul unit nu se poate edita ca o copie a raportului: editarea ar regenera un singur PDF, fără saltul de pagină.

## Variante respinse

- **Salt de pagină în tabelul pdfme.** pdfme nu îl oferă. Un tabel separat pe secție în același șablon ar cere schimbarea generatorului de șabloane pentru toate rapoartele.
- **Un fișier PDF pe secție.** Operatorul ar trebui să salveze și să tipărească mai multe fișiere.
- **Numerotare continuă pe tot documentul.** Ar trebui rescrisă în PDF-ul unit, iar fiecare secție tipărită separat ar avea numere fără sens.

## Consecințe

- Pentru N secții, generarea durează de aproximativ N ori mai mult decât un raport cu toate secțiile.
- Un raport SQL care își filtrează singur datele după `sectie_id` ar primi „-” pentru partea fără secție. De aceea opțiunea nu apare la rapoartele SQL.
