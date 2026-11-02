# TODO - Rapoarte configurabile si coloane calculate

Data: 2026-11-02. Cerinte utilizator, consemnate de Codex. Taskuri: F09-T13, F09-T14. Stare: de facut; nu sunt implementate prin acest document. Nu se schimba formulele sau rapoartele existente.

## F09-T13 - Date si structura libera a raportului

- Alegerea campurilor din baza de date pentru fiecare raport, inclusiv o interogare SQL definita de utilizator.
- Definirea libera a structurii, ordinii si amplasarii campurilor, independent de sursele fixe actuale.
- Integrare Markdown pentru titluri, texte, tabele si variabile. Markdown simplu nu ofera coordonate arbitrare; HTML/CSS sau editorul vizual sunt variante propuse pentru pozitionare precisa, nu decizii aprobate.
- Surse si sabloane separate: un set de date poate fi afisat in mai multe forme.
- Caz cerut explicit: lista cu cate un rand pe garda, cu angajatul si data garzii. Coloane suplimentare posibile: marca distincta de garda, program, inceput, sfarsit si durata.
- Propuneri tehnice de evaluat: SQL numai SELECT, parametri pentru perioada si filtre, campuri rezultate disponibile ca variabile, reutilizarea calculelor aplicatiei, previzualizare si versionarea sabloanelor. Nu se introduce acum o biblioteca sau un motor de executie.

## F09-T14 - Catalogul coloanelor calculate solicitate

| Coloana | Definitie ceruta / observatii |
| --- | --- |
| Total zile cf contract | Totalul zilelor conform contractului pentru zilele lucratoare ale perioadei. Calendarul normal exclude weekendul si sarbatorile legale, conform regulii deja confirmate. Absentele se evidentiaza separat; regula perioadelor partiale de contract trebuie precizata la implementare. |
| Total ore cf contract | Totalul orelor contractuale pentru zilele lucratoare ale perioadei. Pentru norma zilnica fixa, sursa este norma contractuala si zilele aferente; se respecta contractele fara norma prestabilita, inclusiv garda, fara a inventa obligatii. |
| Total zile concediu | Utilizatorul trebuie sa poata defini din aplicatie ce concedii/absente intra in acest total. Trebuie distins totalul informativ de totalul zilelor care reduc obligatia lunara, ca deplasarile creditate sa nu reduca automat norma. |
| Total zile pentru fiecare tip de concediu | La definirea unui tip nou se genereaza automat o coloana calculata corespunzatoare. Asocierea foloseste ID/cod stabil, nu numele afisat. Tipurile inactive pastreaza istoricul; coloanele disponibile si vizibilitatea lor trebuie administrate. |
| Total ore lucratoare | Cerere exacta: Total zile cf contract - Total zile concediu. Rezultatul acestei scaderi are unitatea zile, nu ore; formula si denumirea necesita confirmare inainte de implementare. |
| Total ore garda | Suma orelor de garda pentru perioada, pastrand distinctia marcilor de program normal/garda si evitand dublarea aceleiasi inregistrari. |
| Total ore garda weekend/sarbatori | Configurabil din aplicatie. Utilizatorul precizeaza intervalul: 00 ziua de intrare - 24 ora de iesire. Se consemneaza formularea fara a presupune ca toate orele dintre capete sunt ore efectiv lucrate. Necesita confirmarea regulii exacte de intersectare cu garda si calendarul. |
| Coloane specifice noi | Utilizatorul poate defini direct din aplicatie coloane calculate proprii, cu denumire, unitate, surse si formula, si le poate include/pozitiona in rapoarte. Sintaxa si operatiile disponibile urmeaza sa fie stabilite. |

## Reguli existente de pastrat pana la alte confirmari

- Numai absentele aprobate intra in totalurile lunare; zilele de absenta sunt zile lucratoare, fara weekend si sarbatori legale.
- Efectul categoriei absentei este deja configurabil: fara ore lucrate si cu reducerea obligatiei, sau creditare din norma/program, fara reducerea obligatiei. Totalul concediilor ales de utilizator nu trebuie sa anuleze implicit aceasta distinctie.
- Soldul confirmat ramane L Fin = L Inceput + ore lucrate - ore lucratoare.
- Garda nu are norma prestabilita, T sau ZNL; orele eligibile de weekend/sarbatori au fost discutate separat pe intervale calendaristice 00-24. Cererea noua de configurabilitate nu introduce procente ori reguli monetare noi.

## Clarificari necesare inaintea implementarii

1. Total ore lucratoare: confirmarea unitatii si formulei. Propunere de verificat, nu regula aprobata: (zile contractuale - zile de concediu care reduc obligatia) x norma zilnica; pentru norme variabile, suma orelor contractuale pe zile ramase. Alternativ, rezultatul zile minus zile se numeste Total zile lucratoare.
2. Total zile concediu: selectia se face pe tipuri, categorii sau ambele? Se pastreaza evitarea numararii duble a unei date cand tipurile sunt reunite.
3. Garda weekend/sarbatori: se insumeaza exclusiv intersectia orelor efective ale garzii cu zilele eligibile 00:00-24:00 sau un interval extins integral de la 00:00 in ziua intrarii pana la 24:00 in ziua iesirii? Cum se atribuie intre luni si cum se trateaza ziua simultan weekend si sarbatoare, fara dublare?
4. Coloane noi: se definesc pe angajat/luna, pe fiecare inregistrare sau ambele; cum se afiseaza datele lipsa, rotunjirile si unitatile? Definitiile trebuie validate si versionate, cu teste pentru formulele de calcul.

Pas urmator: clarificarea formulelor si proiectarea editorului intr-un task de implementare alocat explicit. Acest TODO nu autorizeaza modificarea regulilor existente si nu reprezinta o livrare de functionalitate.
