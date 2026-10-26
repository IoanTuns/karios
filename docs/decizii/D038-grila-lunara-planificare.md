# D038 · Grilă lunară tabelară de planificare și publicare unificată

Stare: acceptată · Data: 2026-10-26 · Propusă de: Antigravity · Decisă de: Ioan Tuns · Task: F11-T03

## Context

Utilizatorul a solicitat ca modulul de planificare să permită, pentru fiecare angajat, alegerea lunii și bifarea zilei, formându-se direct o pereche `(zi, tip program)`. Aceasta include toate schimburile speciale gestionate în planificare:
- După-amiază (`S7` - 15:00–22:00)
- Intrare weekend / zi liberă (`IN9` - 22:00–07:00, regim INTRARE_ZI_LIBERA)
- Gărzi (`G24` - 24 ore)
- Ture de weekend / sărbători (`Z12` 07:00–19:00, `N12` 19:00–07:00)

Fluxul anterior din F11-T01 se baza pe un dialog modal de adăugare și o listă verticală plată de înregistrări, ceea ce făcea anevoioasă vizualizarea și completarea lunară pe angajați.

## Decizie

1. **Grilă lunară tabelară:**
   - Ecranul **Planificare** adoptă ca vedere principală o matrice lunară (angajații pe rânduri, zilele lunii 1–28/29/30/31 pe coloane), similară ca ergonomie cu grila de pontaj, dar dedicată planificării turelor speciale.
   - Weekendurile și sărbătorile legale sunt evidențiate vizual.
   - Absențele active ale angajatului (concedii de odihnă, medicale, speciale) sunt afișate distinct pe fundalul celulei (`celula-absenta`) pentru a evita planificarea în zile indisponibile.

2. **Atribuire directă pereche zi – program:**
   - O paletă rapidă cu cipurile schimburilor speciale (`S7`, `IN9`, `G24`, `Z12`, `N12`) și acțiunea de ștergere permite aplicarea directă dintr-un clic pe celula sau celulele selectate.
   - Selecția suportă clic simplu, `Shift+Click` și navigare din tastatură.
   - Se validează ca programul aplicat să fie permis de categoria de program a angajatului (dacă angajatul are categorie activă).
   - Stivă Undo/Redo (`Ctrl+Z` / `Ctrl+Y`) pentru anularea/refacerea pașilor locali înainte de salvare.

3. **Publicare unificată:**
   - Butonul **Previzualizează și publică** suportă publicarea unificată a întregului plan lunar (sau a selecției curente), fără a obliga utilizatorul să filtreze pe o singură secțiune dacă dorește publicarea simultană a gărzilor, turelor de weekend și a după-amiezilor.
   - Se păstrează tranzacționalitatea, verificarea suprapunerilor, protecția absențelor și a pontajelor efective introduse între timp (conform D037).

4. **Vedere istorică / detaliată:**
   - Se păstrează un comutator între vizualizarea tabelară (Grilă) și vizualizarea detaliată (Listă) pentru auditul reviziilor și al stărilor de publicare.

