# F16-T23: cursor, zoom si importul zilelor anuale

Codex, 2026-10-31. Versiune 0.89.4, manual 1.7. Include revizia editoriala 0.89.3 din F02-T24. Utilizatorul a confirmat ca importul vizeaza numarul anual de zile, nu cereri de concediu cu intervale.

## Cauze si remedieri

- Dialog avea efect dependent de callbackul inchide si ocupat: fiecare tastare din FormularAngajat creeaza un callback nou, iar cleanup inchidea dialogul, apoi showModal il redeschidea. Focusul revenea la primul camp. Reprodus cu tastare reala Input.insertText in Prenume: test esuat inaintea corectiei.
- Efectul modal este acum legat de montare/demontare; un ref pastreaza callbackul si starea ocupat actuale. Escape respecta modificarile nesalvate si starea curenta. Aceasta este implementarea comuna folosita de celelalte formulare; cautarea showModal/close nu a identificat un alt ciclu modal independent.
- Butoanele S/M/L existau in Pontaj si Planificare (mod grila); nu au fost demonstrate ca sterse. Sunt protejate de comprimarea flex si verificate la 1024/1366px. Meniul de marire/restaurare este disponibil implicit in dialoguri, cu tinta de cel putin 32px.
- Importul genera antetele Zile_Concediu_An si Norma_Ore_Zi, dar parserul nu recunostea formele normalizate. Aliasele oficiale sunt acum recunoscute. Nu se inventeaza drepturi de concediu; valoarea vine din fisier.
- Previzualizarea importului include concediul anual, pastreaza zero distinct de lipsa si afiseaza norma in ore/minute fara rotunjire.

## Verificari

- UI Edge RO si EN, date fictive: tastare in Prenume si insertie in mijloc, focus/cursor in Nume/Marca/Functie, maximizare pastrata dupa tastare, refuz Escape pastreaza formularul modificat, Escape acceptat inchide.
- Densitate S/M/L in Pontaj si Planificare, vizibilitate la latimi 1024/1366; fiecare buton schimba clasa grilei.
- Previzualizare import cu 22 zile si 4h30: ambele afisate; aceasta proba simuleaza numai raspunsul Wails. Persistenta reala este verificata de testul Go separat.
- Test Go CSV/XLSX cu sabloanele reale generate de aplicatie: zilele ajung in baza, norma partiala este recunoscuta si actualizarea unui angajat existent la zero zile este pastrata.
- Regresie ajutor (dialog comun), 18 teste frontend, 3 probe manual si control traduceri trecute. Prima proba Planificare a cerut completarea API-urilor din simulator; eroarea nu provenea din aplicatie. Prima compilare a testului de import a folosit un nume de helper gresit, corectat la pregateste.

Nu s-a efectuat audit WCAG complet sau proba operator in WebView2. Nu se modifica baza prin migrari sau calculul salariilor. Dependente: Wails 2.17.0/npm curente, Go tranzitive urmarite separat in F14-T21. Pas urmator: build oficial si review, fara push implicit.
