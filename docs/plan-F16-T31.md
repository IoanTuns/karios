# F16-T31: implementare UX30

Responsabil: Codex. Stare: gata. Cerere: toate cele sase recomandari din analiza UX 0.92.19.

1. Erori traduse si selectie rapoarte revalidata.
2. Configurare cu liste filtrate si eligibilitate explicata.
3. Comparatie inainte/dupa pentru suprascrieri.
4. Domeniul selectiei si salvare/abandon coerente.
5. Panou lunar cu legaturi si stari reale.
6. Filtre eliminabile si accesibilitate tastatura, denumiri de regiuni.

Se refolosesc API-urile si regulile existente. Nu se inventeaza calcule, reguli de plata sau restaurare automata. Fara dependente noi; Excel amanat. Urmeaza inspectie, implementare, verificari si livrare cu manual.

## Implementare

| Recomandare | Rezultat |
| --- | --- |
| 1 | Erori ale editorului de reguli traduse cu parametri; selectia persoanei din Rapoarte se revalideaza si nu se trimite un ID ascuns dupa schimbarea filtrelor. |
| 2 | Verificarea configurarii deschide lista cu criteriul exact, eliminabil si cu revenire; formular angajat/categorie cu programe permise/excluse si motive. |
| 3 | Previzualizare obligatorie la adaugare/editare/reactivare absenta si conversie Pontaj; aceeasi operatie in tranzactie anulata, inclusiv documente si revizii. Comparatie pe sursa, marca si zi, inclusiv medical pe marci asociate si tura intreaga. Publicarea Planificarii pastreaza comparatia existenta si adauga marca. |
| 4 | Rezumat vizibil al marcilor si intervalului; dialog comun Ramai/Renunta/Salveaza si continua la navigarea/schimbarea perioadei grilelor si sumarului Planificarii. Salvarea esuata pastreaza datele si nu navigheaza. Publicarea este separata. |
| 5 | Panou Activitatile lunii cu contoare reale ale datelor salvate, stare inchisa/deschisa, verificari si legaturi care pastreaza anul/luna. Nu se declara validare contabila automata. |
| 6 | Filtre cu nume, valoare si eliminare individuala; mesaj fara rezultate; denumiri specifice de regiuni, focus vizibil, ID-uri unice de dialog si Escape numai pentru dialogul activ. Compact/Normal/Aerisit explicite. |

D079 ramane propusa. Reguli de calcul, rapoarte si schema BD neschimbate. Aplicatie 0.92.20, manual 1.31; fara dependente noi. Go stabil 1.27.3 si npm fara actualizari; actualizarile tranzitive Go existente raman urmarite in F14-T21, nu sunt incluse in acest lot UX.

## Probe si limite

Teste backend dedicate pentru rollbackul simularii: pontaj efectiv, schita, planificare, absenta, document atasat, revizii si jurnal; medical pe marca asociata, conversie Pontaj, zi nelucratoare si luna inchisa. Comparatia nu rezerva datele: salvarea finala aplica verificarile existente pe starea curenta.

Probe Edge pe date fictive: UX30 RO/EN (filtrare configurare, eligibilitate, invalidare selectie raport, salvare esuata si reusita, abandon fara salvare, panou lunar si perioada, Escape in dialog suprapus); absente din Pontaj RO/EN cu previzualizare si confirmare; focus/cursor/zoom; contrast ambele teme. Simularea UI este separata de backendul real, verificat in Go.

Protectia nativa la inchiderea Windows ramane cea existenta; dialogul comun nou acopera navigarea interna din grile. Alte formulare isi pastreaza validarile si confirmarea proprie de abandon. Nu s-a efectuat certificare cu cititor de ecran sau schimbarea DPI real Windows. Nu s-au accesat baze reale. Excel ramane amanat explicit. Pas urmator: verificari finale, build oficial si proba operator. Fara push/merge.

Verificari finale: gofmt -l -w ., go vet ./... si go test -p 1 -timeout 120s ./... trecute; TypeScript/Vite, 24 teste frontend, probe UX30 RO/EN inclusiv abandon si Escape suprapus, nomenclatoare RO, absente RO/EN, focus/zoom si contrast ambele teme trecute. 3 teste manuale, 7 livrare, documentatie 0 erori. Primele probe au identificat un selector de editare UI inadecvat si API simulat de backup lipsa; fixturele au fost corectate, nu regulile aplicatiei. Testul nou de medical tine cont ca CO anterior retrage deja pontajul efectiv de pe marca principala.
