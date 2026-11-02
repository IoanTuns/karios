# F07-T38 - Efectele absentelor si accesul din Planificare

Responsabil: Codex. Stare: gata. Ramura: F07-T38-efecte-absente-planificare.

Cerinta utilizatorului din 2026-11-02: tipuri configurabile fara ore lucrate sau cu ore considerate lucrate (deplasari), introducere si din Planificare.

Confirmari: tipurile fara ore lucrate reduc si obligatia lunara; pentru ore creditate sunt necesare ambele optiuni, norma zilnica din contract si durata programului planificat. Sursa se alege explicit pe tip. Lipsa programului/normei nu produce ore presupuse. Numai absentele aprobate participa la totaluri. Se pastreaza prioritatea medical > alte absente > program, calendarul lucrator, previzualizarea si auditul.

Etape: configurare persistenta compatibila; calcul comun si teste anonime; formular comun in Planificare; traduceri/manual/versiune; verificari si build.

Nu se adauga reguli monetare. Excel ramane amanat conform cerintei anterioare. Pas urmator: implementare si verificare.


## Rezultat 0.92.23

Migrare 0029, efect explicit pe tip si minute creditate persistente pe zi. Formularul tipului si lista afiseaza efectul; previzualizarea arata minutele creditate. Serverul ignora duratele propuse de client si le calculeaza din sursa aleasa. Tipurile folosite nu pot schimba categoria/efectul. Varianta ANTERIOR pastreaza regulile istorice existente, fara backfill presupus. [D080](decizii/D080-efecte-absente.md) ramane propusa pentru reprezentarea tehnica; regulile functionale sunt confirmate explicit de utilizator.

Planificare: selectie exacta de persoane/zile, eliminare weekend/sarbatori, acces la formularul comun, previzualizare si salvare. Modificarile nesalvate ale planului trebuie salvate sau anulate intai; nu se abandoneaza implicit. Programul necesar orelor trebuie sa existe salvat. Confirmarea salveaza atomic absenta si suprascrierile prin serviciul existent, apoi reincarca planificarea.

Teste: surse norma (420 min) si program (720 min), fara ore, obligatia lunara neredusa pentru creditare si redusa pentru concediu, aprobari/anulare, protectia efectului folosit, eroare la lipsa programului si rollback integral al previzualizarii. Calculul pur testeaza aprobata/in asteptare/respinsa si date insuficiente. Edge RO/EN: formular din Planificare, previzualizare, confirmare, creare tip cu efect NORMA. Date exclusiv anonime.

Dependentele frontend au commit separat 0a9c788 (0.92.22); npm test: 9 fisiere, 24 teste trecute. TypeScript/Vite, manuale (3), livrare (7), documentatie (0 erori) trecute. Go: gofmt, vet si suita completa. Testele tehnice de migrari au primit numai numarul noii migrari; referintele rapoartelor au versiuni noi, cele vechi raman intacte. Referinta particulara a garzilor era numita fix v12; mutata in v13 pentru noua generatie.

Incercari corectate: sqlc nu este instalat in PATH; accesoriile SQL mici sunt scrise separat, fara modificarea fisierelor generate. Comanda unittest cu modul nu gasea importul documentatie; executia directa verificari/test_livrare.py a trecut. Primele probe UI au selectat o duminica sau au asteptat prea putin incarcarea/selectia grilei; scenariul asteapta explicit ziua selectata. Nu au fost accesate date operationale.

Limite: instantaneul duratei este pastrat la editarea aprobarii/reactivare; anularea nu restaureaza programul. Pentru o sursa/durata diferita deja folosita se creeaza un tip distinct. Numai programe cu durata cunoscuta si fara pauza intra in creditare, conform calculului nominal existent. Aceasta functie nu defineste remunerarea ori rubricile CEV/CFS/FP/Deleg, care pastreaza maparile explicite existente. Excel ramane amanat.

Pas urmator: verificarea de catre operator a tipurilor configurate si a unui caz real, fara publicare sau merge automat. Livrarea se consemneaza in jurnal dupa build.
