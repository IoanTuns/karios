# D057 · Manual comun pentru documente, ajutor contextual și AI

Stare: propusă. Propusă de: Codex, 2026-10-30. Task F02-T17. Cerere Ioan Tuns: manual în română și evaluarea unui manual interactiv în AI și a ghidurilor în fiecare meniu.

## Propunere

Manualul are o sursă editabilă în repo, versiune de document și versiune de aplicație acoperită. HTML-ul local cu căutare și PDF-ul sunt livrabile documentare. Nu necesită rețea sau un model AI. Funcțiile nelivrate nu se prezintă ca existente.

Pentru integrarea ulterioară, sursa se separă în articole cu ID stabil, ecran, sarcină, pași, rezultat, erori, versiune minimă și traduceri. Același articol alimentează panoul „Cum să faci”, căutarea globală și răspunsul din AI. Se evită copierea unor instrucțiuni divergente în prompturi, cataloage și PDF.

Întrebările de utilizare merg întâi la o căutare locală deterministă în articole. AI-ul poate reformula un extras relevant, cu titlul și linkul articolului, fără a primi întregul manual la fiecare mesaj. Conținutul documentației nu acordă permisiuni și nu ocolește validările sau confirmările operațiilor. Lipsa unui articol produce un răspuns explicit, nu o funcție inventată.

Panoul contextual prezintă pași și legături către ecrane. Tururile care evidențiază butoane sunt opționale și trebuie legate de identificatori stabili ai interfeței. Niciun tur nu salvează sau publică date în locul utilizatorului.

## Alternative

- Întregul manual în promptul sistemului: respins ca propunere din cauza consumului de context și a latenței pe calculatoare modeste.
- Texte independente pentru fiecare meniu și fiecare răspuns AI: nerecomandat, deoarece actualizările se pot contrazice.
- Numai PDF: păstrat pentru distribuire și tipărire, insuficient pentru ajutor contextual.

Integrarea în aplicație nu este implementată de F02-T17 și nu schimbă executabilul recompilat din main. Ea necesită task de cod separat, versiune nouă, texte RO/EN și verificări de navigare și rutare AI. Decizia rămâne propusă.
