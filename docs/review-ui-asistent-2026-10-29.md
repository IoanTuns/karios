# Review UI · Asistent Karios 0.82.1 · 29 octombrie 2026

F14-T23 · Codex · Bază: `f093186`, funcționalitatea F14-T22. Concluziile și prioritățile sunt propuneri, legate de D049. Nu s-au modificat interfața, calculele sau datele utilizatorului.

Fluxul oferă control util asupra persoanelor și modificărilor. Principalele dificultăți sunt stabilitatea poziției la editare, spațiul foarte mic al zonei de lucru și delimitarea propunerii de rezultatul efectiv. Recomand remedierea primului lot înaintea integrării F14-T22 în main.

## Metodă și limite

- Inspectate `AsistentAI.tsx`, `PlanAsistentAI.tsx`, `asistentSesiune.ts`, CSS și cataloagele RO/EN. Verificate contractul planului și tokenul, fără audit complet al backendului.
- Revăzute capturile existente F14-T22 ale formularului izolat. Acestea verificau alegerea opțiunilor, confirmarea și încadrarea în lățime; nu acopereau integrarea în conversație.
- Probe noi în Edge headless: componenta reală `AsistentAI`, sesiunea reală, componenta reală de plan, CSS de producție și un înveliș cu clasele layoutului aplicației. Puntea Wails, datele, portul și rezultatele sunt simulate. Nu este executabilul Wails și nu este o probă cu modelul local real.
- Scenarii: editare RO/EN, previzualizare cu 129 de celule, rezultat, eroare de token expirat și selecție nominală cu 200 de persoane fictive. Cererea inițială este aceeași ca în captura utilizatorului.
- Capturi solicitate la 1024×640 și 1440×900. Browserul a raportat `innerWidth/innerHeight` de 998×547, respectiv 1414×807 în aceste probe; dimensiunile DOM sunt raportate ca atare, fără a le confunda cu dimensiunea PNG-ului. Nu s-a certificat comportamentul în WebView2 la toate scalările Windows.
- Instrumentarea a numărat apelurile `scrollIntoView`, fără a schimba comportamentul metodei. Animația smooth face ca pozițiile verticale să varieze în timp; concluziile se bazează pe apelurile reproduse, dimensiunile elementelor, structura codului și capturi, nu pe un prag exact de deplasare.
- Nu s-au făcut sesiuni cu operatori, audit complet de accesibilitate, măsurători pe hardware cu 8 GB sau benchmarkuri ale modelului. Lista de 200 de persoane este o probă de ergonomie, nu o măsurare de performanță.

Artefacte locale, ignorate de Git: `build/bin/audit-ui-F14-T23/`. Conțin capturi, DOM, JSON cu măsurători, script și copia fixture-ului temporar. Prima execuție a întâlnit fișierul DOM încă blocat de redirecționarea browserului; reluarea după închiderea fluxului a reușit.

## Pro

| Ce funcționează bine | Beneficiu |
| --- | --- |
| Persoanele, programul și intervalul pot fi corectate înainte de scriere | Utilizatorul poate repara interpretarea fără să formuleze din nou întreaga comandă |
| Zilele și tratarea celulelor ocupate au alegeri explicite | Dispar două presupuneri importante din fluxul de planificare |
| Sunt disponibile grupuri și selecții nominale, plus lista completă din previzualizare | „Toți” poate fi verificat concret, nu doar acceptat din textul modelului |
| Sunt afișate valori vechi, valori propuse și motive de excludere | Există informația necesară unei confirmări în cunoștință de cauză |
| Confirmarea include numărul modificărilor; datele schimbate invalidează planul | Confirmarea este legată de un plan verificat, nu de o afirmație liberă a modelului |
| Există mesaje de eroare, rezultat verificat și salturi către fiecare lună | Se poate continua verificarea în pontaj după salvare |
| Etichete asociate controalelor, legendă pentru grupuri, `aria-pressed`, antete de tabel și focus CSS vizibil | Bază bună pentru utilizarea cu tastatura; accesibilitatea completă nu este certificată |
| Catalog RO/EN și rearanjarea câmpurilor în funcție de lățime | Formularul izolat se adaptează și textele noi sunt traductibile |

## Cons și îmbunătățiri

P1 = primul lot recomandat înainte de integrare. P2 = următorul lot. S/M reprezintă efort relativ, nu termen de livrare. Toate rândurile sunt propuneri neimplementate.

| ID / prioritate | Constatare și dovadă | Îmbunătățire | Criteriu de acceptare | Efort |
| --- | --- | --- | --- | --- |
| I072 · P1 | Editarea formularului modifică obiectul conversației, iar efectul din `AsistentAI.tsx:124` derulează la capăt pentru orice schimbare a conversației. Proba RO/EN a înregistrat un al doilea apel după clic pe opțiunea zilelor, fără mesaj nou. Este o problemă de integrare pe care proba izolată nu o acoperea | Derulare automată numai la mesaj nou și numai dacă utilizatorul urmărea capătul conversației; păstrarea poziției și focusului la editarea intenției | Schimbarea selecției, datei sau opțiunii zilelor într-un card mai vechi nu produce apel automat la capăt; tastatura și poziția rămân la controlul editat | S |
| I073 · P1 | Formularul are aproximativ 731 px înălțime într-o zonă de conversație de 228 px în proba mică. Se cumulează derularea paginii, conversației și tabelului; tabelul poate avea 420 px. Cardul de operații moștenește limita de 80% a bulei, plus mai multe chenare și paddinguri. După „Arată modificările”, previzualizarea este inserată mai jos fără focus sau salt dedicat; confirmarea se află după tabel | Plan activ pe lățimea zonei utile, rezumat compact și etape clare Interpretare → Verificare → Rezultat. Acțiune principală persistentă în zona planului, fără acoperirea rândurilor. La generarea planului, aducerea rezumatului verificării în vedere. Comenzile motorului și sugestiile devin secundare după începutul lucrului | La un viewport CSS de minimum 1024×640 sunt vizibile contextul activ și acțiunea principală; nu trebuie derulată întreaga pagină pentru confirmare. Editarea și verificarea folosesc o singură zonă principală de derulare, iar tabelul rămâne consultabil integral | M |
| I074 · P1 | După aplicare, aceeași tabelă păstrează coloana „Propus” și stările „Completare” / „Înlocuire”. Rândurile păstrate ori excluse continuă să afișeze Z7 ca valoare propusă. Captura rezultatului conține un rând „6 → Z7 · Păstrată”, sub titlul „Rezultat verificat”. Aceasta este ambiguitate de prezentare; nu dovedește o scriere greșită în bază | Delimitare explicită între previzualizare și rezultat: „Se va salva” înainte, „Valoare finală” și „Salvat / Păstrat / Neaplicat” după. Sumar final cu modificările efectuate și excluderile; formularul vechi poate fi pliat într-un rezumat | Un rând păstrat cu valoarea 6 afișează 6 ca rezultat final; o absență nu este prezentată drept Z7 aplicat; numărul de rânduri salvate corespunde rezultatului confirmării | M |
| I075 · P2 | Tabelul listează toate celulele în pagini de 50, fără căutare sau filtru de stare. În proba cu 129 de rânduri, două excepții sunt amestecate cu 127 de completări. O singură frază cu șase contoare este greu de scanat; lista persoanelor nu distinge selectații de cei care chiar primesc modificări | Sumar separat pe modificări, păstrate și excluse; filtru „Doar excepții”, căutare după persoană și sumar pe lună. Afișarea numărului selectat și a celui afectat | Excluderile se găsesc fără parcurgerea tuturor paginilor; fiecare lună are număr de modificări și excluderi; totalurile rămân identice cu planul, indiferent de filtru | S/M |
| I076 · P2 | Selecția nominală afișează 200 de checkboxuri într-o zonă de maximum 240 px, fără căutare. Proba confirmă 200 de controale și absența unui câmp de căutare. Nu există „doar selectați” sau comenzi pentru selecția vizibilă | Căutare după nume/marcă, vedere „Selectați”, comenzi explicite „Selectează rezultatele” și „Deselectează rezultatele”, cu numărul persoanelor clar vizibil | Un angajat din lista de 200 se găsește direct prin nume sau marcă; filtrarea nu pierde selecțiile și comenzile de grup afectează numai rezultatele indicate | S/M |
| I077 · P2 | `complet` verifică prezența câmpurilor, nu ordinea datelor; eroarea este globală. La token expirat, proba confirmă dispariția previzualizării și focus pe BODY. Mesajul explică motivul, dar utilizatorul trebuie să regăsească singur „Arată modificările” | Validare lângă câmp, identificarea clarificărilor încă necesare și focus pe primul câmp invalid. Pentru plan expirat: mesaj cu „Actualizează previzualizarea”, păstrarea contextului și evidențierea diferențelor după recalculare | Intervalul inversat identifică datele relevante; eroarea are acțiune accesibilă cu tastatura; noul plan nu se confirmă automat și cere verificare din nou | M |
| I078 · P2 | Intențiile și rezultatele sunt în sesiune, dar planul neconfirmat și pagina sunt `useState` în componentă. Remontarea pierde previzualizarea. Mai multe intenții rămân în conversație fără o stare vizibilă „înlocuită” sau „anulată”. Aceste constatări provin din cod, nu dintr-o probă de navigare în executabil | Stări vizibile pentru plan: în editare, pregătit, expirat, aplicat, anulat. Comenzi „Revizuiește” / „Anulează”. La revenire, restabilirea zonei de lucru ori un mesaj explicit că este necesară o nouă verificare | Este clar care plan este activ; un card vechi nu pare mai nou decât revizuirea sa; revenirea nu prezintă un token vechi drept confirmabil și explică orice regenerare necesară | M |
| I079 · P2 | Cererea inițială rămâne citată chiar după editarea grupului, fără o etichetă distinctă pentru textul inițial. Sumarul folosește luni ISO și multe numere într-o frază; tabelele au antete de 11 px și celule de 12 px. În browserul de probă, data nativă apare în format american deși interfața este RO; acesta poate depinde de mediul browserului | Rezumat curent lizibil: „10 angajați · 1 octombrie–30 noiembrie 2026 · Z7 · luni–vineri · păstrează existente”. Separarea cererii originale de interpretarea actuală. Mesaje și stări lizibile, inclusiv avertisment textual la înlocuire; rezumatul datei nu depinde exclusiv de formatul controlului nativ | După schimbarea grupului, rezumatul actual reflectă noua selecție. Perioada este neambiguă în RO/EN. Numerele și excepțiile se disting fără citirea întregului paragraf și fără a depinde numai de culoare | S/M |
| I080 · P2 | Bara motorului arată starea și portul, dar nu identifică modelul activ sau limita contextului; utilizatorul a trebuit să întrebe separat cum află modelul. Această lipsă îngreunează înțelegerea erorilor de context și compararea modelelor | Detalii pliabile „Despre asistent”: modelul efectiv folosit, limita configurată și, dacă motorul o raportează, utilizarea contextului ultimei cereri. Portul și comenzile tehnice rămân în zona secundară | Modelul afișat corespunde procesului activ, nu primului fișier găsit în folder; nu se afișează consumuri de tokenuri sau memorie inventate. Valorile indisponibile sunt indicate explicit | M |

Surse principale: `frontend/src/AsistentAI.tsx:123`, `frontend/src/AsistentAI.tsx:355`, `frontend/src/PlanAsistentAI.tsx:18`, `frontend/src/PlanAsistentAI.tsx:32`, `frontend/src/PlanAsistentAI.tsx:64`, `frontend/src/PlanAsistentAI.tsx:105`, `frontend/src/PlanAsistentAI.tsx:124`, `frontend/src/style.css:501`, `frontend/src/style.css:514`, `frontend/src/style.css:515`, `frontend/src/style.css:819` și cheile `ai.plan_*` din cataloage. Numerele de linie se referă la baza auditului.

## Ordine recomandată

1. I072: stabilizarea editării. Include un test cu formularul în componenta reală de conversație și cel puțin două mesaje/carduri, nu numai formularul izolat.
2. I073–I074: spațiu pentru verificare, acțiune persistentă și rezultat fără ambiguități. Include scenarii cu excluderi și întoarcere la interpretare.
3. I075–I077: găsirea excepțiilor, selecții nominale și recuperarea din erori.
4. I078–I080: ciclul complet al planului, rezumat, lizibilitate și identificarea modelului.

Se păstrează selecția deterministă, clarificările necesare, protecția absențelor, verificarea tokenului și confirmarea explicită. Reducerea numărului de ecrane sau chenare nu trebuie să elimine aceste verificări.

## Reproducere și probe locale

- Editare: creează un plan, derulează la câmpurile de sus și schimbă zilele sau persoanele. Efectul global este apelat chiar dacă nu se adaugă un mesaj.
- Verificare: generează un plan cu peste 50 de celule pe o fereastră mică. Urmărește unde apar noul sumar, excluderile și confirmarea.
- Rezultat: păstrează o celulă completată și exclude o absență. După aplicare, verifică sensul coloanei „Propus” pentru acele rânduri.
- Eroare: modifică datele după previzualizare, apoi confirmă. Planul trebuie refuzat; verifică mesajul, focusul și posibilitatea de reluare.
- Persoane: deschide selecția nominală cu o echipă mare și caută o marcă situată la finalul listei.

Capturi locale: [previzualizare](../build/bin/audit-ui-F14-T23/previzualizare-ro.png), [rezultat](../build/bin/audit-ui-F14-T23/rezultat-ro.png), [eroare](../build/bin/audit-ui-F14-T23/eroare-ro.png), [persoane](../build/bin/audit-ui-F14-T23/persoane-ro.png), [editare RO](../build/bin/audit-ui-F14-T23/editare-ro.png), [editare EN](../build/bin/audit-ui-F14-T23/editare-en.png).

Review-ul nu implementează remedierile și nu validează fidelitatea răspunsurilor modelului real. Pas următor: evaluarea umană a D049 și alocarea primului lot de implementare.
