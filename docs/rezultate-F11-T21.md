# F11-T21 - Programe multiple pe marca si zi

Responsabil Codex.2026-11-04. Ramura F11-T21-programe-multiple. Aplicatie0.92.31, migrare0031, manual1.41. Rapoarte pontaj16, detaliat17, individual16, absente13, garzi13, flexibil5. D086 propusa; cerinta functionala confirmata explicit de utilizator.

## Implementare

Pozitii de Planificare distincte, fiecare cu ID si revizie. Proiectie canonica zilnica in Pontaj (coduri separate prin +, separator deja interzis in codurile individuale), maximum1024 caractere. Migrarea copiaza valorile/reviziile/timpii existenti si pastreaza integritatea referentiala. Nu sunt biblioteci noi sau reguli de plata noi.

Validarea compara intervalele nominale: suprapuneri pe marca, inclusiv zile/luni vecine, refuzate; capetele adiacente permise. Publicarea grupeaza toate pozitiile active ale zilei, chiar la selectarea uneia sau la sectiuni diferite, verifica reviziile/tokenul, scrie ciorna o singura data si marcheaza fiecare pozitie. Un program inactiv/neeligibil blocheaza intregul set. Republicarea este idempotenta. Anularea retrage numai contributia publicata, fara a pierde celelalte coduri sau a restaura reciproc doua anulari din acelasi lot. Editarea necesita republicare.

Sincronizarea inversa pastreaza setul, refuza pierderea unui plan independent, valideaza categorii si intervale. Ciornele multiple sunt validate si fata de zilele vecine ale aceluiasi lot. Calculul expandeaza codurile in programele originale, cu atribuirea/calendarul configurat; absenta se crediteaza o singura data. PROGRAM crediteaza suma programelor zilei. Medicalul elimina intreaga tura intersectata, pe familia de marci, pastrand programul neintersectat din ziua precedenta. Restul prioritatilor existente raman aplicate.

Interfata: sumarul are Adauga program pe zi, selectare/editare/eliminare separata, totaluri pe toate pozitiile. Grila afiseaza toate codurile; o actiune rapida ambigua pentru o zi multipla este refuzata cu indrumare catre sumar. Programele normale pot fi planificate explicit; NORMAL/LUCRATOARE exclude weekend/sarbatori. Textele noi sunt in RO/EN. Rapoartele flexibile listeaza fiecare garda, graficul lunar arata toate codurile, calculele si comparativul insumeaza intervalele.

## Verificari

Teste anonime adaugate: migrare cu ciorna/revizie9 pastrata si foreign_key_check; programe adiacente si suprapuse, ture vecine, publicare partiala si republicare, anulare individuala/integrala, blocare atomica cu program inactiv, sincronizare inversa, program normal+tura, excludere weekend/sarbatoare, medical cu pastrarea intervalului neintersectat, calcul17ore din5+12 cu7ore weekend, absenta creditata o singura data, doua randuri SQL/PDF/Excel pentru aceeasi marca/zi. Trecute.

TypeScript/Vite si24 teste frontend trecute. Scenariul Edge RO/EN de adaugare/salvare/redeschidere/eliminare individuala a trecut. Simularea UI verifica interactiunea, nu calculele; acestea sunt probate in motorul Go real pe date anonime.

Prima suita completa a semnalat numere de migrari fixate in testele tehnice vechi, limita16 pentru cod simplu si vechea asteptare de pastrare a Pontajului la retragerea unui program publicat. Numerele tehnice au crescut cu1, limita16 ramane pentru codurile simple, iar asteptarea de retragere reflecta noul flux documentat. Datele/asertiunile contabile nu au fost relaxate. Referintele tehnice pentru rapoarte au fisiere NOI pe noile versiuni; cele existente raman pastrate.

Scenariul UI a necesitat corectii in harness: injectarea autentificarii demonstrative, selectarea programului nou gol dupa ordonarea intervalelor si evitarea suprascrierii citirii Planificarii cu un mock gol din ajutor. Nu s-au accesat date reale.

Dependente verificate: Go1.27.3, npm outdated gol. Module directe actuale exceptand Excelize2.11.1 amanat deja deoarece cere efp prerelease; indirectele raman urmarite in F14-T21. Nu s-au instalat programe noi. Verificarile finale si livrarea se consemneaza mai jos.

## Predare

Pas urmator: proba operator si review pentru integrare. Fara push sau merge. Buildul oficial, manualele si sumele de verificare urmeaza dupa verificarea finala.

Verificare finala: gofmt -l -w ., go vet ./..., go test -p 1 -timeout 180s ./... trecute integral. TypeScript/Vite,24 teste frontend, programe multiple Edge RO/EN, regresie UX30 RO/EN,3 teste manuale,7 teste livrare si documentatie0 erori trecute. Dialogul nativ de imprimare si baza operationala reala nu au fost accesate. Urmeaza buildul oficial din commit curat.

Livrare oficiala0.92.31 din `058c410`: `build/bin/pontix-0.92.31-20261104-094519-9043d77021f54ab3a0a551d6107cd21a/pontix.exe`. Manuale MD/HTML/PDF0.92.31, editia1.41 regenerate; ZIP langa director. Manifest complet, versiune/commit, SHA256 executabil/manuale/ZIP si CRC ZIP verificate. SHA256 executabil: `CE1CC9B46B2D4F8D5901683FA230CA7A0F494988310368BFC0C3FA13E60E8174`. Modele Wails identice semantic, ordinea sursei pastrata. Fara push/merge.
