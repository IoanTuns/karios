# F16-T26 - Ore in format 24h

Responsabil: Codex. Ramura F16-T26-ore-24h. Aplicatie 0.92.11, manual 1.22. D071 propusa.

Formatterul comun foloseste h23 (00-23), fara AM/PM indiferent de limba. Orele din jurnal, restaurare, audit, mesaje de blocare cont si ultima modificare Pontaj folosesc aceeasi regula. Afisarile directe toLocaleTimeString/toLocaleString au fost eliminate. Datele si fusul orar raman cele anterioare; duratele si totalurile nu sunt transformate in ore de ceas. Rapoartele Go folosesc deja 15:04, fara schimbare de continut sau versiune de raport.

La definirea schimbului, controlul nativ time putea impune 12h dupa Windows. A fost inlocuit cu HH:MM validat, intre 00:00 si 23:59; validarea backend existenta ramane. Traducere pentru format RO/EN.

21 teste frontend trecute, inclusiv 3 cazuri pentru ro/en/en-US: miezul noptii 00, dimineata 07, seara 19, 23 si lipsa dayPeriod, inclusiv optiune hour12 suprascrisa. TypeScript/Vite trecut. Edge nomenclatoare RO/EN: accepta 00:00/07:00/19:00/23:59, refuza 24:00/7:00/19:60/7 PM; programul se salveaza prin fluxul existent. Edge focus/cursor/zoom trecut. 3 teste manual, 7 teste livrare si control documentatie trecute.

gofmt, go vet ./... si go test -p 1 -timeout 120s ./... trecute, inclusiv reverificarea finala dupa campul HH:MM si traduceri. Dependente: Go stabil 1.27.3 si npm actuale; 62 actualizari tranzitive Go urmarite separat prin F14-T21. Fara migrare sau biblioteca noua. Fara schimbarea testelor contabile de referinta.

Pas urmator: build oficial si proba operator in WebView2. Datele reale nu au fost accesate. Fara push/merge in main.
