# F07-T33 - Prioritatea concediilor

Responsabil: Codex. Ramura F07-T33-prioritate-absente. Aplicatie 0.92.9, manual 1.20. D069 propusa; regulile functionale confirmate de utilizator sunt documentate in F01.

Concedii numai L-V fara sarbatori din calendar. Medical > alte concedii > program. Inregistrarea modifica atomic schita, planificarea, absenta si jurnalul. Medicalul anuleaza intreaga tura suprapusa, inclusiv ziua precedenta si programele tuturor marcilor explicit asociate persoanei; nu dubleaza absentele. Salvarea ulterioara a programului este protejata in Pontaj si Planificare. Lunile inchise si pontajul efectiv raman protejate.

Teste noi: weekend/sarbatoare pentru CO si CM; CO inlocuit de medical fara revizie initiala; medical protejat contra altui concediu; doua marci asociate si persoana separata; anulare integrala tura; reintroducere program refuzata; miezul noptii fara suprapunere; luna precedenta blocata cu rollback. Testele tehnice anterioare au fost adaptate la suprascrierea ceruta si la zile lucratoare. Nu s-au schimbat teste contabile de referinta, procente sau reguli de plata.

TypeScript/Vite, 18 teste frontend, 3 teste manual, 7 teste livrare, control documentatie si Edge RO conversie Pontaj trecute. gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute. Edge RO/EN trecute. Build oficial complet din 730716d: `build/bin/pontix-0.92.9-20261031-151335-d2bdb97efb0f4e5cac1726f6dfbd5a53/pontix.exe`. Manual 1.20 MD/HTML/PDF alaturi. Identitate, hashuri executabil/manuale/arhiva si integritate ZIP verificate. Go stabil 1.27.3 si npm fara actualizari; actualizarile tranzitive Go raman urmarite prin F14-T21. Fara dependente sau migrari noi.

Limitari: calendarul configurat trebuie verificat; absentele istorice nelucratoare nu sunt migrate, programele nu sunt restaurate automat la anulare. Pontajele efective necesita verificare separata. Datele reale nu au fost accesate. Pas urmator: proba operator si review. Fara push/merge.
