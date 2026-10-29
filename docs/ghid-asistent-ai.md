# Asistentul AI: interpretare, clarificări și confirmare

F14-T22 · 0.82.0 · D048 propusă.

1. Scrie cererea, de exemplu: „Toți angajații în octombrie și noiembrie lucrează de la 7 la 14”.
2. În cardul „Am înțeles cererea astfel”, verifică persoanele, programul și datele. Poți schimba orice câmp. Dacă mai multe programe au același interval, alege codul exact.
3. Alege zilele și comportamentul pentru celulele completate. Nu există o alegere implicită dacă nu ai precizat-o în comandă. Luni–vineri reprezintă zilele săptămânii, fără a deduce reguli de sărbători sau de plată.
4. Apasă „Arată modificările”. Verifică lista completă a persoanelor, lunile, numărul celulelor noi, înlocuirile și excluderile. Tabelul are pagini de 50 de rânduri.
5. Bifează că ai verificat și confirmă numărul exact al modificărilor. Dacă datele s-au schimbat, aplicația cere o nouă previzualizare. Un plan nu se salvează parțial între luni.
6. După salvare apare numărul verificat și câte un buton pentru pontajul fiecărei luni. „Conversație nouă” șterge conversația, nu modificările deja salvate.

## Grupuri și contexte

Grupurile disponibile sunt: toți angajații activi, o secție, un departament, o categorie de program sau o listă de persoane bifate. Pentru grupuri se citesc persoanele din baza de date, nu dintr-o enumerare produsă de model. Grupurile personalizate persistente și combinarea mai multor filtre rămân extensii posibile, nu sunt implementate în acest task.

Filtrul „fără pontaj” se evaluează pentru fiecare lună: un angajat cu celule completate în octombrie poate rămâne eligibil în noiembrie. Previzualizarea explică excluderile. Absențele și lunile închise nu sunt suprascrise; programele nepermise și suprapunerile sunt excluse. Modificarea unei ture existente într-o altă secțiune se face în ecranul Planificare.

Pentru model, contextul se selectează după cerere: planificare, absențe, pontaj sau ajutor. Se trimit numai persoanele identificate explicit și înregistrările lor relevante, în limita bugetului; datele omise sunt semnalate modelului ca extras incomplet. Analizele deterministe folosesc datele complete. Nu se garantează precizia răspunsurilor libere ale modelului dintr-un extras.

O singură intenție ghidată folosește un program și un interval continuu de cel mult un an, maximum 10.000 de celule. Cererile complexe, datele relative sau selecțiile neidentificate se corectează în formular. Nu se creează programe sau reguli salariale din text. Ștergerea turelor păstrează fluxul anterior de confirmare; nu este transformată într-o planificare de formular.

## Verificare și limite

`internal/ai/intentie_test.go` conține 40 de comenzi RO/EN și probe pentru contexte. `app_ai_plan_test.go` verifică traseul cerere–clarificări–previzualizare–confirmare, toate persoanele pe două luni, tokenuri expirate și rollback la o eroare în a doua lună. Testele vechi pentru grupuri au fost adaptate pentru clarificarea explicită, păstrând verificările absențelor și celulelor ocupate; nu s-au modificat teste de calcul confirmate de contabilitate.

Proba interfeței folosește exclusiv persoane fictive și un backend Wails simulat: `npm --prefix frontend run dev -- --host 127.0.0.1 --port 5175`, apoi `http://127.0.0.1:5175/probe/plan-ai.html` sau `?limba=en`. Parcurge automat alegerea zilelor, păstrarea celulelor, previzualizarea, bifarea verificării și confirmarea; rezultatul este în atributul `data-rezultat` al paginii. Proba nu este inclusă în build-ul de producție. Aceasta verifică interfața; testele Go verifică salvarea reală în baze temporare. Edge headless: trecut RO/EN la 1024×900 și RO la 800×900, fără depășirea lățimii. Capturile și DOM-ul sunt în `build/bin/proba-ui-F14-T22/` (artefacte ignorate de Git). Proba inițială în sandbox nu a pornit browserul; reluarea cu profil separat în afara sandboxului a reușit.

Modelul local real, timpul de răspuns și memoria pe un calculator cu 8 GB nu au fost măsurate. Contextul maxim configurat al motorului nu a fost schimbat în acest task. Decizia D048 rămâne propusă pentru evaluare umană.
