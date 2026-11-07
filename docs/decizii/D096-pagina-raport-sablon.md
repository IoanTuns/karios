# D096 · Setări de pagină dau dimensiunea șablonului vizual

Stare: propusă · Data: 2026-11-07 · Taskuri: F09-T32

## Context

În modul editor vizual (pdfme), dimensiunea paginii venea din șablon (`basePdf.width/height`), fixată la crearea lui. Formatul și orientarea din **Setări de pagină** schimbau doar anexa. Ioan Tuns a raportat (2026-11-07) un raport setat Portret care ieșea peisaj: șablonul fusese creat pe A4 peisaj (implicitul unui model nou).

## Decizie

- Formatul și orientarea din Setări de pagină sunt pagina raportului. La schimbarea lor, șablonul vizual se potrivește: noua lățime și înălțime; elementele se întind proporțional pe lățime (x și lățimea), iar cele din jumătatea de jos păstrează distanța față de marginea de jos (subsolul rămâne jos). Fonturile și marginile rămân.
- La deschiderea unui model al cărui șablon nu se potrivește cu Setări de pagină, potrivirea se face automat, cu un mesaj; modelul trebuie salvat.

## Variante respinse

- **Orientarea doar din editorul vizual** (fără câmpuri în formular): câmpurile există și sunt așteptate de utilizatori.
- **Potrivirea la generare, în Go**: PDF-ul ar fi diferit de ce arată editorul.
- **Scalarea fonturilor**: textul ar deveni prea mic sau prea mare; lățimile coloanelor se potrivesc oricum procentual.

## Consecințe

- Modelele salvate cu șablonul nepotrivit se generează ca înainte până sunt deschise și salvate.
- După o schimbare de orientare, elementele așezate liber pot cere reglaje în editor.
