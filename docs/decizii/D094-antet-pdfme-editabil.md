# D094 · Antetul rapoartelor pdfme se editează pe pagină

Stare: propusă · Data: 2026-11-07 · Taskuri: F09-T29

Implementare: implementată din 0.92.72; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

## Context

În șabloanele pdfme făcute de aplicație, antetul (trei zone și titlul) și subsolul (două zone și numărul paginii) stau în `basePdf.staticSchema`, ca să se repete pe fiecare pagină (F09-T26). Designer-ul pdfme le arată în zona roșie a marginilor, dar nu le poate selecta. Ioan Tuns a cerut (2026-11-07) să le poată edita: de exemplu să șteargă titlul sau să pună data curentă lângă el.

## Decizie

- În editorul vizual, butonul **Editează antetul** mută temporar elementele fixe pe prima pagină, marcate `pontixStatic`, cu marginile reduse la marginea paginii. Acolo se mută, se modifică și se șterg ca orice element; elementele adăugate în acest mod devin și ele fixe.
- **Gata cu antetul** (și, automat, verificarea sau închiderea editorului) pune la loc în `staticSchema` tot ce nu era în corpul paginii la intrare, ca text numai citire (cu `{variabile}`). Marginea de sus și cea de jos se recalculează din elementele fixe (+2 mm), iar tabelul coboară sub antet dacă s-ar suprapune.
- Variabilă nouă `{data_curenta}`: data generării, `ZZ.LL.AAAA`. `{generat_la}` rămâne cu ora.

## Variante respinse

- **Câmpuri separate în formular pentru titlu și dată** — ar fi acoperit doar cazul cerut; pe pagină se poate schimba orice element (poziție, font, conținut).
- **Antetul mutat definitiv în corpul paginii** — pdfme îl desenează atunci doar pe prima pagină; antetul repetat pe fiecare pagină se pierde.
- **Numele `{data}`** — se ciocnește cu o coloană SQL numită `data` (coloanele SQL sunt și ele variabile).

## Consecințe

- Modelul salvat are mereu antetul în `staticSchema`; modul de editare există doar în editor.
- Rapoartele existente nu se schimbă până nu li se editează antetul; nicio versiune de raport nouă.
