# D066 - Marcarea numelor in ierarhie

Stare: propusa. Codex, 2026-10-31. Task F07-T30.

Utilizatorul solicita identificarea automata a inregistrarilor ierarhice in toate contextele. Marcajul de afisare propus este Nume Prenume (Garda), tradus prin catalagele RO/EN. Este aplicat numai unei marci cu legatura explicita angajat_principal_id; nu se deduce din nume, functie sau contract. Un contract de garda independent nu este automat subordonat. In prezent singurul subtip ierarhic disponibil este marca de garda; tipuri noi necesita model explicit.

Formatter comun pentru frontend si formatter comun Go pentru rapoarte/exporturi. Numele si prenumele din BD si formularele de editare raman curate. Cheile si marcile nu se modifica. AI salveaza marcajul in noile instantanee si include legaturile in amprenta previzualizarii, astfel incat o schimbare de ierarhie impune regenerarea. Istoricul si fisierele exportate nu sunt rescrise. Fara migrare sau biblioteca noua.
