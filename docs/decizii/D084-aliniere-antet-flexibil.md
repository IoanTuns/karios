# D084 - Zonele antetului si subsolului flexibil

Stare: înlocuită de D087. Data: 2026-11-03. Responsabil Codex. F09-T18 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Zonele antetului se aranjează acum vizual în editorul pdfme ([D087](D087-rapoarte-sql-pdfme.md)); motorul care le folosea a fost șters în F09-T23.

Captura utilizatorului confirma un defect: editorul ofera trei zone, dar PDF-ul flexibil le stivuieste la stanga. Propunere: TABEL si Markdown folosesc trei coloane simultane egale, cu aliniere L/C/R; inaltimea antetului este maximul textelor dupa impachetare, nu suma lor. Imaginile se ancoreaza la stanga/centru/dreapta in zona proprie, pastrand proportiile. Subsolul foloseste aceleasi alinieri si rezerva maxima de spatiu. Antetul configurat se repeta si pe anexa; continutul anexei porneste dupa antet, fara suprapunere. Coordonatele continutului LIBER raman explicite, fara translatarea elementelor existente; antetul textual al corpului LIBER nu este introdus peste elementele pozitionate manual. Fara sintaxa Markdown/HTML noua si fara biblioteca noua.
