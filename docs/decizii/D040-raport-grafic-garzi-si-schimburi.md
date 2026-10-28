# D040 · Raport dedicat pentru graficul lunar de gărzi și schimburi

Stare: acceptată · Data: 2026-10-26 · Propusă de: Antigravity · Decisă de: Ioan Tuns · Task: F11-T07

## Context

În spitale, clinici și unități cu activitate continuă din România, „Graficul lunar de gărzi și schimburi” este un document oficial distinct de pontajul general salarial:
1. Este întocmit de coordonatorul de gărzi / medicul șef de secție înainte de începerea lunii și aprobat de directorul medical și manager.
2. Centralizează distribuția personalului pe ture și gărzi (`G24`, `Z12`, `N12`, `IN9`, `S7`), cu evidențierea weekendurilor și a sărbătorilor legale.
3. Prezintă totalurile lunare de gărzi și ore pentru fiecare angajat, precum și o linie de sinteză zilnică a acoperirii (`TOTAL / ZI`).
4. Utilizatorul a cerut un raport distinct pentru gărzi și schimburi, separat de pontajul general.

## Decizie

1. **Sursă dedicată de raport:**
   - Adăugarea sursei `GRAFIC_GARZI` în motorul de rapoarte (`internal/rapoarte`), având coloanele `angajat`, `marca`, `functie`, `zile`, `total_garzi` și `total_ore`.
   - Raportul are versiunea proprie `internal/rapoarte/garzi/versiune.go` (`const Versiune = 1`).
   - În baza de date se adaugă definiția inițială `raport-garzi` prin migrarea **0018**, configurată implicit în format A3 Peisaj, cu antet oficial și bloc de semnături în subsol (Întocmit, Avizat, Aprobat).

2. **Generare și totaluri:**
   - Datele sunt citite direct din tabela `planificare` pentru luna, secția și departamentul selectate.
   - Fiecare celulă a lunii afișează codul de gardă/schimb alocat sau `—` dacă persoana este liberă în acea zi.
   - Se calculează automat numărul de gărzi/ture efectuate în lună și totalul de ore planificate.
   - Rândul final de sinteză `TOTAL / ZI` indică numărul de persoane de serviciu în fiecare zi a lunii, permițând verificarea imediată a acoperirii weekendurilor și a sărbătorilor.

3. **Accesibilitate directă:**
   - Pe lângă prezența în secțiunea generală „Rapoarte”, ecranul „Planificare” include butonul direct `Grafic gărzi (PDF)` în bara de acțiuni.
   - Clicul generează instant documentul PDF cu previzualizare modală (opțiune de maximizare), buton de salvare fișier și buton de trimitere directă către imprimantă.
