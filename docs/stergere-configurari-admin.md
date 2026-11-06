# F07-T37 - Stergere definitiva pentru curatarea configurarilor

Responsabil: Codex. Stare: domeniu confirmat explicit de utilizator (raspuns: da, 1 si 2). Implementare si teste in lucru; nicio stergere din datele reale.

## Cerinta

Administratorul trebuie sa poata elimina configurarile inactive si irelevante rezultate din configurarea initiala. Operatorul nu primeste acest drept.

## Domeniul confirmat

| Lista | Stergere propusa | Utilizari care o blocheaza |
| --- | --- | --- |
| Programe de lucru | Program inactiv, fara utilizari | Pontaj, planificare, categorii, reguli, sabloane sau operatii AI |
| Categorii de program | Categorie inactiva, fara persoane/contracte asociate | Angajati, contracte, reguli |
| Tipuri de absenta | Tip inactiv, nefolosit | Orice zi de absenta, inclusiv anulata, pontaj sau regula |
| Functii, grupuri, contracte | Definitie inactiva, nefolosita | Angajati, functii subordonate, categorii, reguli |
| Categorii de schimb, tipuri si regimuri | Definitie inactiva, nefolosita | Programe care o refera |
| Modele de raport | Model personalizat inactiv | Referinte configurate; modelele structurale necesare raman protejate |

Pontajele, zilele de concediu, planificarile, angajatii, conturile si jurnalul NU intra in aceasta curatare. Definitii de baza necesare compatibilitatii nu se elimina. Refuzul trebuie sa explice utilizarea care blocheaza stergerea.

## Efecte asupra asocierilor

Un program are metadate proprii in clasificari_schimburi, create automat. Eliminarea programului ar elimina si aceste metadate proprii prin relatia existenta ON DELETE CASCADE. Nu se elimina programe din categorii pentru a forta stergerea: asocierea cu o categorie blocheaza operatia.

O categorie de program neatribuita poate contine o lista proprie de programe permise. Stergerea categoriei elimina lista ei proprie, pastrand programele. Orice referinta externa blocheaza operatia.

## Autorizare si concurenta

Verificare stricta in endpoint: administrator autentificat, inclusiv cand baza nu are inca utilizatori. Butonul ascuns pentru operator nu este suficient. Confirmarea indica numele inregistrarii si caracterul definitiv. Se verifica revizia si starea in aceeasi tranzactie cu verificarile de utilizare si jurnalizarea stergerii. La orice eroare, nicio modificare partiala.

## Verificari necesare inainte de livrare

- Fara sesiune si cu operator: refuz pe backend, nu doar in interfata.
- Administrator: stergere a unei definitii inactive/nefolosite si pastrarea jurnalului.
- Definitie activa, revizie veche, ID invalid sau inexistent: refuz.
- Utilizare prin FK, cod in schita, reguli JSON, sabloane/operatii AI: refuz cu motiv.
- Asociere proprie eliminata numai in limitele confirmate; nicio eliminare din date operationale.
- Tranzactie si rollback integral la eroare.
- Interfata RO/EN: confirmare, renuntare, eroare si actualizarea listei; operator fara buton.
- gofmt, go vet, suita Go, TypeScript/frontend, verificari documentatie si build versionat cu manual.

## Starea aprobarii automate

Prima incercare de implementare a fost respinsa inaintea executarii comenzii. Motivul: stergeri definitive din mai multe tabele si eliminarea unor asocieri, fara teste complete si fara confirmarea exacta a efectelor ireversibile. Nici ramura, nici fisierele acelei incercari nu au fost create. Dupa confirmarea explicita a utilizatorului, implementarea si testele au fost autorizate. Nicio baza cu date reale nu a fost modificata.

Confirmare primita: toate configurarile enumerate, cu pastrarea datelor operationale si blocarea celor utilizate. O categorie neatribuita trebuie sa poata fi eliminata impreuna cu lista ei proprie de programe permise (nu cu programele); editorul nu permite salvarea unei categorii goale. Orice referinta externa a categoriei blocheaza stergerea.
