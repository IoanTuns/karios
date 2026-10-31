# F11-T20 - Prioritatea surselor de pontaj

Responsabil: Codex. In lucru.

Cerinta confirmata: Planificarea scrie peste valorile existente din Pontaj. Concediile/zilele libere suprascriu programul in Pontaj si Planificare. Medicalul are prioritate maxima, pe toate marcile asociate, inclusiv anularea integrala a turei de noapte intersectate. Numai zile lucratoare pentru absente; lunile inchise si reviziile concurente raman protejate.

Implementare: publicarea planificarii inlocuieste implicit valorile existente, fara optiunea de a omite celulele ocupate. Se pastreaza previzualizarea/tokenul pentru concurenta, jurnalul si refuzul peste absente. Salvarea planului ramane etapa de pregatire, publicarea este scrierea in Pontaj. Verificare explicita pentru pontajul efectiv, nu numai schite.

Teste: inlocuire program/ore, concurenta, absente prioritare, medical pe marci si noapte, rollback, suita Go/UI si manual versionat. Pas urmator: inventar complet surse si ajustarea publicarii.

## Predare intermediara

Inlocuirea implicita a schitelor este implementata. Pontajele efective raman protejate; intrebarea de confirmare a retragerii lor din calcul, cu pastrarea istoricului, este in asteptare dupa respingerea automata. Nicio operatiune respinsa nu a fost executata. [D076](decizii/D076-prioritate-publicare-pontaj.md) propusa. Dependente verificate: Go stabil 1.27.3, npm fara actualizari; cele 62 actualizari tranzitive Go sunt urmarite separat in F14-T21. Nicio biblioteca instalata.

Verificari pentru partea independenta: gofmt, go vet ./..., go test -p 1 -timeout 120s ./..., 22 teste frontend, TypeScript/Vite, Edge focus/zoom, 3 teste manual, 7 teste livrare si control documentatie trecute. Niciun build oficial nou livrat: taskul ramane in lucru pana la raspunsul despre pontajele efective. Pas urmator: aplicarea raspunsului, testele de integrare suplimentare si build versionat. Nu s-a facut push/merge.
