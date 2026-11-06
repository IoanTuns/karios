# D059 - Nomenclatoare pentru categorii de schimburi si functii

Stare: acceptată. Data: 2026-10-31. Propusa de: Codex. Task F07-T24, cerere utilizator pentru categorii multiple si functii selectabile · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Categorii de schimburi cu ID stabil si denumiri i18n.Text, initial LAB/CUSTOM. Migrarea 0024 inlocuieste restrictia enumerata a coloanei categorie pastrand valorile, ID-urile si legaturile programelor. Trigger-ele protejeaza referintele categoriilor; nu se reconstruieste tabelul de programe si nu se modifica migrari din main. Regimul de calcul al programului este separat de categorie.

Functiile au ID stabil, denumiri RO/EN, stare activa si revizie. Asocierea angajat-functie foloseste un tabel separat cu chei straine, pastrand campul text pentru compatibilitatea rapoartelor si importurilor. Redenumirea sincronizeaza textul curent si jurnalizeaza angajatii afectati in aceeasi tranzactie. Migrarea preia valorile existente; importurile ulterioare fara asociere pastreaza textul pentru selectare/standardizare explicita.

Asocierea este curenta, la fel ca vechiul camp functie; istoricul contractual pe luni este o functionalitate separata. Dezactivarea interzice asocieri noi si le pastreaza pe cele vechi. Identitatea nu se redenumeste si nu se ofera stergere. Revizia previne suprascrierea editarilor concurente. Se folosesc numai tabele din lista interna explicita in SQL-ul dinamic. Fara biblioteci noi sau reguli de plata.

Implementarea functionalitatii cerute nu transforma aceasta propunere intr-o decizie acceptata; acceptarea arhitecturii se consemneaza de om.
