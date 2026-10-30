# Analiză UX pentru release pilot · Pontix 0.88.14

F16-T20 · Codex · 2026-10-30 · Bază: 9aafb62, integrare main f15e623. Propuneri pentru review, fără schimbări de cod.

## Concluzie

Aplicația susține bine operarea repetitivă de către un utilizator instruit. Înainte de release, prioritatea este încrederea în date și în efectul acțiunilor: protecție consecventă la navigare, indicatori explicați și diferențierea funcțiilor disponibile cu/fără model AI. Nu recomand o redesenare generală.

## Metodă și limite

Audit euristic al App, Pontaj, AsistentAI, Ajutor, Formulare, CSS și manualului. Inspectată captura reală din proba recentă Edge 1024×640, build/bin/verificare-ui-2026-10-30T17-14-05-006Z/ajutor-drawer-ro.png; sunt reutilizate rezultatele RO/light și EN/dark ale F02-T21, nu prezentate drept probe noi. Capturile istorice trimise de utilizator oferă context, nu dovadă pentru starea actuală. Problemele deduse din cod sunt marcate; nu s-a reprodus pierderea datelor în această analiză. Fără studiu cu operatori, audit WCAG complet, probă pe Windows curat sau benchmark nou.

## Ce funcționează bine

- Grila arată orele și păstrează identitatea persoanei și totalul la vedere; tastatura, selecțiile, clipboardul și anularea reduc munca repetitivă.
- Separarea Pontaj / Planificare / Absențe corespunde unor activități distincte; bara celulei explică programul și modificarea.
- Contextul de lucru și protecția navigării principale sunt implementate; nu repropunem funcții deja existente.
- Ajutorul contextual este local, cu pași, rezultat, probleme frecvente și versiune. Citările AI deschid ghidul și din drawer.
- Temele și contrastul au verificări dedicate. Dialogul de ajutor are focus vizibil, Escape și revenirea focusului verificate în F02-T21.
- Manualele au versiunea livrării; restaurarea are pași și confirmare explicite.

## Probleme prioritare

P1 = de rezolvat sau verificat înainte de pilot cu date reale; P2 = îmbunătățire de operare. Prioritatea exprimă impactul, nu o vulnerabilitate sau o pierdere de date deja demonstrată.

| ID | Prioritate / dovadă | Impact | Propunere și criteriu de acceptare |
| --- | --- | --- | --- |
| UX01 | P1, cod: Pontaj.tsx, callbackurile drawerului onDeschidePontaj/onDeschideConcediu schimbă perioada/ecranul fără poateAbandona; navigarea principală din App folosește confirmarea | Comportament diferit pentru aceeași schimbare de context; risc pentru schița nesalvată | Un singur mecanism de verificare, cu Salvează / Renunță / Rămâi când este posibil. Probă: modifică o celulă, cere altă lună din AI, refuză; luna și valoarea rămân. Același test pentru salt la concedii. Se verifică separat anul cererii: callbackul actual folosește anul grilei. |
| UX02 | P1, cod: acoperireZile colorează după zero persoane și sub 25% din randuri; randuri este lista filtrată, fără necesar de personal configurat | Culoarea poate sugera deficit sau acoperire suficientă fără criteriu operațional; filtrarea schimbă pragul | Până există necesar confirmat, afișează neutru „Persoane cu înregistrări în selecția afișată”, explică baza și evită „optim/deficit”. Test: filtrele și lipsa datelor nu produc verdict despre necesarul echipei. Nu se inventează praguri. |
| UX03 | P1 pentru claritatea pilotului, cod: card-lipsa-fisiere-ai apare în conversație deși ghidurile deterministe nu cer model | Operatorul poate crede că întregul asistent este indisponibil; configurația tehnică domină sarcina | Mesaj distinct: „Ghidurile sunt disponibile. Comenzile care folosesc modelul necesită configurare.” Detaliile binar/model într-o zonă pliabilă. Test fără model: operatorul găsește și folosește un ghid. |
| UX04 | P2, cod/CSS: drawer de 420px, max-width 45vw, pe același rând cu grila; AsistentAI complet în interior | La rezoluții mici, conversația și grila concurează pentru spațiu; multe comenzi înaintea conținutului | Mod compact: conversație și perioada active, opțiuni avansate pliate; pe ecrane înguste panou suprapus sau comutare. Măsurare la 1024×640 și 1366×768: formularul și închiderea accesibile, rânduri utile vizibile, fără derulare orizontală a formularului. |
| UX05 | P2, captură și cod: căutarea filtrează lista de ghiduri, dar articolul ales rămâne afișat chiar dacă nu mai corespunde; lista și corpul au derulare separată | „Fără rezultate” în stânga poate coexista cu un articol vechi în dreapta; acțiunile sunt jos | Selectează primul rezultat sau explică păstrarea articolului; la zero rezultate oferă resetare. Bară de acțiuni persistentă fără să acopere textul. Test căutare fără rezultate și zoom 200%. |
| UX06 | P2, cod: sursa afișează implicit ID tehnic și versiune; eroarea pentru articol lipsă/versiune diferită folosește același mesaj | Utilizatorul nu știe dacă sursa lipsește sau este doar veche | Titlu lizibil drept sursă principală, ID în detalii; mesaje distincte și acțiune „Deschide ghidul versiunii curente”. Test ID invalid și versiune veche. |
| UX07 | P2, cod: meniul contextual are Escape și butoane, dar nu are un mecanism explicit de deschidere prin tastatură/focus pe prima acțiune | Funcțiile descoperite prin clic dreapta sunt mai greu accesibile fără mouse | Shift+F10 sau alternativă vizibilă, focus inițial și revenire la celulă. Test integral fără mouse; nu se elimină comenzile existente din toolbar. |
| UX08 | P1 ca validare de release, manual: actualizarea cere copierea întregului folder karios; recuperarea implică backup, seif și cheie | Utilizatorul poate confunda backupul bazei cu o copie portabilă completă | Checklist scurt: copie completă, aplicație închisă, executabil separat, folder karios complet, verificare după pornire. În Setări: locația datelor și ultima copie reușită, dacă nu există deja echivalent. Acceptare prin probă pe copie fictivă și Windows separat; existența documentației nu dovedește recuperarea testată. |

## Plan recomandat

1. Lot înainte de pilot: UX01, UX02, UX03 și proba de actualizare/recuperare UX08. Fiecare implementare primește task separat, versiune și teste potrivite.
2. Lot de ergonomie: UX04–UX07, pe baza observațiilor operatorilor. Păstrăm comenzile rapide pentru utilizatorii avansați.
3. Probă cu 3–5 operatori: prima pornire, completarea unei săptămâni, schimbare S7 și verificarea Planificării, concediu pe interval, raport, backup/restaurare pe copie, întrebare de ajutor fără model. Date exclusiv fictive. Înregistrăm finalizare fără ajutor, erori, solicitări de sprijin și înțelegerea stării salvat/nesalvat; nu declarăm un scor de utilizabilitate fără măsurare.

## Predare

Analiza nu modifică regulile de calcul, pragurile sau interfața și nu acceptă decizii în numele utilizatorului. Pas următor: alegerea și alocarea taskurilor pentru lotul de pilot. Actualizările de dependențe nu sunt necesare pentru acest document; nu s-au rulat din nou teste de cod.
