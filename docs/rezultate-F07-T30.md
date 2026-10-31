# F07-T30 - Nume marcate pentru ierarhie

Responsabil: Codex. Ramura F07-T30-marcare-ierarhie. Aplicatie 0.92.6, manual 1.17. [D066](decizii/D066-nume-ierarhic.md) propusa.

## Modificare

Marcile subordonate de garda sunt afisate Nume Prenume (Garda), folosind traducerea RO/EN. Datele personale si marcile nu se modifica. Omonimii si contractele de garda fara persoana principala nu sunt marcati automat.

Formatterul comun frontend acopera lista Angajati, selectii, Pontaj, Planificare/lista/Gantt, absente, inlocuitori, Gantt absente, Jurnal, previzualizarea importului pentru marcile existente si rapoarte/AI. Formularele pastreaza numele editabil curat. Backendul aplica marcajul in toate sursele PDF, Excel din document, ambele foi Excel clasice, planuri/sumare/exporturi AI noi si numele din validarile planificarii. Numele subordonat din AI este protejat de amprenta legaturii explicite. Instantaneele vechi si fisierele deja exportate nu sunt modificate; afisarile UI care au ID folosesc identitatea curenta, cu fallback la instantaneu.

Rapoarte pontaj/detaliat/fisa v8, absente/garzi v6; export Excel clasic 1.1.0. Referintele vechi sunt pastrate. Fara modificarea calculelor sau a schemei BD.

## Verificari

Edge cu punte simulata RO/EN: creare marca de garda, marcaj in numele listei, Pontaj, Planificare, selectii Absente/Rapoarte, plus fluxurile existente de functii/filtre/reguli. Go: formatter bilingv si omonim, toate sursele de raport, ambele foi Excel, marcaj AI si eliminarea asocierii/invalideaza token vechi. TypeScript/Vite, 18 teste frontend, 3 teste manual, 7 teste livrare si control documentatie trecute. gofmt, go vet si suita Go completa trecute. Testele noi au fost ajustate pentru normalizarea existenta a spatiilor si pentru pozitia numelui in antetul fisei, apoi au trecut. Build-ul oficial se consemneaza dupa generare.

Dependente verificate: Go/npm actuale, 62 actualizari tranzitive Go urmarite in F14-T21; fara biblioteci noi. Nu au fost deschise date reale. Proba interactiva WebView2 si imprimarea fizica raman la operator. Pas urmator: review si proba operator; fara push/merge in main.
