# F07-T35 - Tipuri si regimuri configurabile

Responsabil: Codex. Ramura F07-T35-tipuri-regimuri-configurabile. Aplicatie 0.92.12, manual 1.23. D072 propusa; regula General/anterior in Normal solicitata explicit de utilizator.

Migrare noua 0027: nomenclatoare tipuri_program/regimuri_program, referinte stabile clasificari_schimburi, conversie GENERAL in LUCRATOARE afisat Normal. Tipul schimbului, intervalele, durata, starea si identitatea programului se pastreaza. Compatibilitatea cu scrieri vechi GENERAL este normalizata prin triggere. Migrarile existente nu au fost modificate.

Editorii reutilizeaza salvarile cu revizii si jurnal ale nomenclatoarelor. Denumiri RO/EN, activare/dezactivare, comportament explicit la creare. Un comportament deja utilizat nu se schimba implicit; se creeaza o definitie noua, apoi se asociaza explicit in program. O definitie inactiva ramane pe programele existente, dar nu se atribuie la programe noi. Selectiile din formular si filtrele sunt dinamice.

Legaturi: programul salveaza ID-urile definitiilor si comportamentele validate in aceeasi tranzactie. Categoriile de program pastreaza ID-urile acelorasi programe; contractele folosesc categoriile existente, iar pontajul, planificarea, precompletarea, AI si calculele primesc comportamentele motorului existent. Redenumirea nu modifica calculul. Regimul personalizat se afiseaza si in selectorul de completare Pontaj. Nu se inventeaza sporuri, procente sau compensari.

Teste: migrare din schema 26 cu program Noapte personalizat/inactiv, conversie regim, pastrarea intervalului si verificare FK; nomenclatoare personalizate, autoritatea comportamentului salvat fata de payload, pastrarea N12 in categoria si planificarea persoanei, numar utilizari, redenumire, conflict de revizie, blocarea reinterpretarii comportamentului utilizat, dezactivare si refuz asignari noi. Testele de migrare au fost actualizate numai pentru numarul tehnic 27 si conversia General ceruta; referintele contabile nu s-au schimbat. O inlocuire prea larga a numarului de migrari a afectat temporar o asteptare tehnica pentru doua marci; a fost corectata, iar testele datebaza si bd au trecut.

TypeScript/Vite si Edge nomenclatoare RO/EN trecute: adaugare tip/regim si asociere program. 21 teste frontend, 3 teste manual, 7 teste livrare, control documentatie si Edge focus/zoom trecute. gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute integral dupa actualizarea contoarelor tehnice. Go stabil 1.27.3, npm actual; actualizarile tranzitive Go sunt urmarite separat in F14-T21. Fara biblioteca noua.

Pas urmator: build oficial si proba operator in WebView2. Datele reale nu au fost accesate. Niciun push sau merge in main.

Prima suita Go completa a identificat inca doua contoare de migrari in testele de actualizare a rapoartelor; au fost crescute pentru migrarea 0027. Verificarile de continut ale rapoartelor au ramas intacte. Spatiile finale regenerate de Wails au fost normalizate, fara modificarea logicii bindingurilor.
