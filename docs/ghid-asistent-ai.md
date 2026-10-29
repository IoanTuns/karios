# Asistentul AI: interpretare, clarificări și confirmare

F14-T22, F14-T24, F14-T25 · 0.84.0 · D048, D049 și D050 propuse.

1. Scrie cererea, de exemplu: „Toți angajații în octombrie și noiembrie lucrează de la 7 la 14”.
2. Deschide „Plan activ”. În etapa „Interpretare”, verifică persoanele, programul și datele. Poți schimba orice câmp. Dacă mai multe programe au același interval, alege codul exact. Istoricul rămâne în „Conversație”, iar un plan mai vechi se deschide din istoric sau din selectorul de planuri.
3. Confirmă anul când a fost preluat din context sau dintr-un șablon. Alege zilele și comportamentul pentru celulele completate. Nu există o alegere implicită dacă nu ai precizat-o în comandă. Luni–vineri reprezintă zilele săptămânii, fără a deduce reguli de sărbători sau de plată.
4. Apasă „Arată modificările”. În etapa „Verificare”, verifică lista completă a persoanelor, lunile, numărul celulelor noi, înlocuirile și excluderile. Tabelul are pagini de 50 de rânduri. Contextul și acțiunile rămân la vedere, iar conținutul central se derulează. „Revizuiește interpretarea” revine la câmpuri și invalidează vechea verificare.
5. Bifează că ai verificat și confirmă numărul exact al modificărilor. Dacă datele s-au schimbat, aplicația cere o nouă previzualizare. Un plan nu se salvează parțial între luni.
6. După salvare, etapa „Rezultat” arată valorile finale și stările „Salvat”, „Păstrat”, „Neaplicat” sau „Deja identic”. O celulă păstrată cu valoarea 6 rămâne afișată ca 6; o absență exclusă nu apare ca program aplicat. Numărul verificat și butoanele către pontajul fiecărei luni rămân în subsol. „Conversație nouă” șterge conversația, nu modificările deja salvate.

Comenzile de pornire/oprire și pornirea automată sunt în „Opțiuni motor”. În timpul editării nu se mută poziția sau focusul la capătul conversației. În istoric, un mesaj nou derulează automat numai dacă utilizatorul urmărea capătul. Schimbarea planului selectat cere din nou previzualizarea unui plan neconfirmat; intenția editată și rezultatele salvate rămân în sesiune. Persistența previzualizărilor între planuri este urmărită prin I078, separat de primul lot P1.

## Grupuri și contexte

Grupurile disponibile sunt: toți angajații activi, o secție, un departament, o categorie de program sau o listă de persoane bifate. Pentru grupuri se citesc persoanele din baza de date, nu dintr-o enumerare produsă de model. Grupurile personalizate persistente și combinarea mai multor filtre rămân extensii posibile, nu sunt implementate în acest task.

Filtrul „fără pontaj” se evaluează pentru fiecare lună: un angajat cu celule completate în octombrie poate rămâne eligibil în noiembrie. Previzualizarea explică excluderile. Absențele și lunile închise nu sunt suprascrise; programele nepermise și suprapunerile sunt excluse. Modificarea unei ture existente într-o altă secțiune se face în ecranul Planificare.

Pentru model, contextul se selectează după cerere: planificare, absențe, pontaj sau ajutor. Se trimit numai persoanele identificate explicit și înregistrările lor relevante, în limita bugetului; datele omise sunt semnalate modelului ca extras incomplet. Analizele deterministe folosesc datele complete. Nu se garantează precizia răspunsurilor libere ale modelului dintr-un extras.

O singură intenție ghidată folosește un program și un interval continuu de cel mult un an, maximum 10.000 de celule. Cererile complexe, datele relative sau selecțiile neidentificate se corectează în formular. Nu se creează programe sau reguli salariale din text. Ștergerea turelor păstrează fluxul anterior de confirmare; nu este transformată într-o planificare de formular.

## Instrumente disponibile din 0.84.0

„Plan nou / verificare” deschide editorul și funcționează fără motor LLM. În selecția nominală poți căuta după nume sau marcă și selecta/deselecta rezultatele filtrate. Angajații inactivi nu pot fi selectați. Perioada inversată sau mai mare de un an nu poate fi previzualizată.

În etapa de verificare, „Rezumat pe luni” și „Rezumat pe persoane” disting persoanele selectate de cele afectate efectiv. „De ce sunt incluși sau excluși angajații?” prezintă criteriul de grup, inclusiv angajații inactivi sau din alte grupuri. Selecția nu garantează schimbarea tuturor zilelor: motivele de păstrare/excludere sunt în tabelul zilnic. Tabelul se filtrează după persoană, lună și stare, cu paginare recalculată după filtrare.

„Șabloane de planificare” salvează parametrii editorului sub un nume ales de utilizator. Sunt locale, în baza criptată, maximum 100; textul liber al cererii nu se păstrează în șablon. Încărcarea cere verificarea perioadei și nu aplică modificări. Secțiunile/categoriile și mărcile sunt rezolvate din datele curente la previzualizare; dacă nu mai există, selecția trebuie corectată. Ștergerea șablonului cere confirmare și nu schimbă pontajul. Șabloanele sunt disponibile și în „Instrumente”.

„Verifică pontajul existent” cere doar persoanele, perioada, anul confirmat și zilele vizate, fără a cere programul propunerii. Semnalează goluri fără absență, coduri necunoscute, programe inactive/nepermise, diferențe față de planificare, absențe și suprapuneri. Orele manuale sunt marcate ca intervale neverificabile. Regulile configurabile deja existente se evaluează pe schița lunii complete pentru persoanele selectate, inclusiv când verifici numai câteva zile; plafonul lunar nu este calculat doar pe intervalul vizibil. Sunt afișate codul, denumirea și severitatea regulii încălcate. Nu se introduc reguli de plată, nu se certifică respectarea legislației și nu se corectează automat date.

„Exportă planul CSV” salvează local într-un fișier ales prin dialog. Exportă toate rândurile, nu numai pagina/filtrul vizibil, cu selecția explicată și rezumatele. Backendul reconstruiește și verifică tokenul propunerii; rezultatul aplicat se citește din operația salvată. CSV are separator `;`, UTF-8 cu BOM și neutralizarea formulelor. Valorile finale apar separat de cele propuse. Exportul unui rezultat este istoric; dacă operația a fost anulată, acest lucru este indicat explicit. Nu se trimit date în exterior.

În „Instrumente” → „Verifică ultima operație”, consultă valorile de restaurat, bifează verificarea și apasă „Anulează operația AI”. Numai planificările aplicate de fluxul ghidat din 0.84.0 au instantanee pentru anulare. Operațiile din versiunile vechi și ștergerile din fluxul vechi nu sunt reconstruite din conversație. O editare ulterioară, inclusiv salvarea aceleiași valori cu revizie nouă, schimbarea planificării/publicării/pontajului efectiv ori închiderea lunii blochează întreaga anulare. Restricțiile curente pentru valorile restaurate sunt reverificate. Dacă apare o eroare în a doua lună, prima lună se restaurează tranzacțional la starea de dinaintea tentativei de anulare. Jurnalul rămâne, reviziile cresc, pozițiile nou create se dezactivează, iar o celulă inițial absentă rămâne goală cu revizie. Operațiile și șabloanele rămân disponibile după repornire; conversația nu este necesară. Anularea unei operații nu garantează anularea uneia mai vechi care atinge aceleași celule: reviziile schimbate impun refuzul, nu rescrierea editărilor.

„Opțiuni motor” arată numele modelului din procesul propriu și contextul configurat, 32768 de tokenuri. Nu confundă primul model de pe disc cu modelul activ. Consumul curent de tokenuri și RAM nu este măsurat. „Conversație nouă” elimină istoricul trimis modelului, păstrând pontajul, operațiile și șabloanele.

## Probe automate

`internal/ai/intentie_test.go` conține 40 de comenzi RO/EN și probe pentru contexte. `app_ai_plan_test.go` verifică traseul cerere–clarificări–previzualizare–confirmare, toate persoanele pe două luni, tokenuri expirate și rollback la o eroare în a doua lună. Testele vechi pentru grupuri au fost adaptate pentru clarificarea explicită, păstrând verificările absențelor și celulelor ocupate; nu s-au modificat teste de calcul confirmate de contabilitate.

Proba interfeței folosește exclusiv persoane fictive și un backend Wails simulat: `npm --prefix frontend run dev -- --host 127.0.0.1 --port 5175`, apoi `http://127.0.0.1:5175/probe/plan-ai.html` sau `?limba=en`. Parcurge automat alegerea zilelor, păstrarea celulelor, previzualizarea, bifarea verificării și confirmarea; rezultatul este în atributul `data-rezultat` al paginii. Proba nu este inclusă în build-ul de producție. Aceasta verifică interfața; testele Go verifică salvarea reală în baze temporare. Edge headless: trecut RO/EN la 1024×900 și RO la 800×900, fără depășirea lățimii. Capturile și DOM-ul sunt în `build/bin/proba-ui-F14-T22/` (artefacte ignorate de Git). Proba inițială în sandbox nu a pornit browserul; reluarea cu profil separat în afara sandboxului a reușit.

Modelul local real, timpul de răspuns și memoria pe un calculator cu 8 GB nu au fost măsurate. Contextul maxim configurat al motorului nu a fost schimbat în acest task. Decizia D048 rămâne propusă pentru evaluare umană.

F14-T24 adaugă proba ecranului complet la `/probe/asistent-ui.html`, cu parametrii `limba=ro|en` și `caz=editare|previzualizare|istoric|rezultat|eroare|fara-motor`. Ea verifică un plan mai vechi în conversație, poziția și focusul, confirmarea vizibilă, tokenul expirat și valorile finale. Sunt folosite numai date fictive și Wails simulat. [Rezultate și limite](rezultate-ui-F14-T24.md).
