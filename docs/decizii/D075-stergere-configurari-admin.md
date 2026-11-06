# D075 - Stergere configurari inactive

Stare: acceptată. Data: 2026-10-31. Task: F07-T37 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

## Context

Utilizatorul a confirmat explicit domeniul: configurari inactive nefolosite, numai administrator. Datele operationale si jurnalul se pastreaza.

## Propunere tehnica

Endpoint unic cu lista de tabele permise, administrator autentificat si verificarea rolului/starii actuale in baza. Revizia, referintele, eliminarea si auditul se verifica in tranzactie. Referintele istorice, inclusiv inactive, blocheaza stergerea. Sunt verificate si referintele prin cod din schite si cele JSON din setari/sabloane/operatii AI. Definitii structurale sunt protejate.

Metadatele proprii ale programului si lista proprie a categoriei se elimina numai impreuna cu definitia nefolosita; programele din lista raman. Snapshotul complet ramane in jurnal.

## Variante respinse

- Stergere in cascada a datelor operationale: contrara domeniului confirmat.
- Ascunderea butonului fara autorizare backend: insuficienta.
- Stergere a inregistrarilor active: necesita dezactivare explicita inainte.

## Consecinte

Curatarea este conservatoare: o referinta istorica poate cere pastrarea unei configurari inactive. Nu se adauga migrari sau biblioteci. [Specificatie](../stergere-configurari-admin.md).
