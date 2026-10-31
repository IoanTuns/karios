# F11-T20 - Prioritatea surselor de pontaj

Responsabil: Codex. In lucru.

Cerinta confirmata: Planificarea scrie peste valorile existente din Pontaj. Concediile/zilele libere suprascriu programul in Pontaj si Planificare. Medicalul are prioritate maxima, pe toate marcile asociate, inclusiv anularea integrala a turei de noapte intersectate. Numai zile lucratoare pentru absente; lunile inchise si reviziile concurente raman protejate.

Implementare: publicarea planificarii inlocuieste implicit valorile existente, fara optiunea de a omite celulele ocupate. Se pastreaza previzualizarea/tokenul pentru concurenta, jurnalul si refuzul peste absente. Salvarea planului ramane etapa de pregatire, publicarea este scrierea in Pontaj. Verificare explicita pentru pontajul efectiv, nu numai schite.

Teste: inlocuire program/ore, concurenta, absente prioritare, medical pe marci si noapte, rollback, suita Go/UI si manual versionat. Pas urmator: inventar complet surse si ajustarea publicarii.

## Predare intermediara

Inlocuirea implicita a schitelor este implementata. Pontajele efective raman protejate; intrebarea de confirmare a retragerii lor din calcul, cu pastrarea istoricului, este in asteptare dupa respingerea automata. Nicio operatiune respinsa nu a fost executata. [D076](decizii/D076-prioritate-publicare-pontaj.md) propusa. Dependente verificate: Go stabil 1.27.3, npm fara actualizari; cele 62 actualizari tranzitive Go sunt urmarite separat in F14-T21. Nicio biblioteca instalata.

Verificari pentru partea independenta: gofmt, go vet ./..., go test -p 1 -timeout 120s ./..., 22 teste frontend, TypeScript/Vite, Edge focus/zoom, 3 teste manual, 7 teste livrare si control documentatie trecute. Niciun build oficial nou livrat: taskul ramane in lucru pana la raspunsul despre pontajele efective. Pas urmator: aplicarea raspunsului, testele de integrare suplimentare si build versionat. Nu s-a facut push/merge.

## Continuare dupa confirmare

Confirmare explicita primita: Da, si pontajele efective. Implementata retragerea logica auditata, fara stergerea fizica. Limitele de asteptare descrise anterior reprezinta starea sesiunii anterioare, acum depasita. Publicarea inlocuieste efectivele aceleiasi date; absentele retrag efectivele intersectate, iar medicalul pe toate marcile explicite asociate. Protejate lunile inchise, inclusiv luna de inceput/sfarsit a turei. Previzualizarea indica numarul efectivelor retrase si detecteaza schimbarea lor in token.

Testele vechi care cereau refuz la existenta unui efectiv au fost adaptate la cerinta explicit confirmata; nu sunt teste de referinta contabila. Prima rulare a evidentiat o extindere accidentala a ignorarii efectivelor la editarea manuala/AI; corectia limiteaza exceptia la publicare. Urmeaza suita completa, livrarea si verificarea pachetului.

## Verificari finale 0.92.17

Gofmt, go vet ./..., go test -p 1 -timeout 120s ./... trecute; apoi teste suplimentare si vet pe pachetul datebaza pentru miezul noptii si rollback audit. 22 teste frontend, TypeScript/Vite, Edge absente din Pontaj RO/EN, 3 teste manual, 7 livrare si control documentatie trecute. Testele acopera auditul integral, tokenul modificat de schimbarea orelor efective, medical pe toate marcile, luna precedenta blocata, interval adiacent fara suprapunere, rollback la esec audit. Afisarea numarului de efective retrase a fost verificata prin compilare TypeScript; nu s-a executat o proba end-to-end WebView2 pe date reale.

Dependente reverificate: Go 1.27.3 stabil, npm fara actualizari, 62 tranzitive Go raman in F14-T21. Manual editie 1.28. Fara migrari sau biblioteci noi. Pas urmator: build oficial si integritate pachet, apoi proba operator. Fara push/merge si fara acces la date reale.
