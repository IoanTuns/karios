# Definirea, generarea și tipărirea rapoartelor

Versiunea 0.21.0 este pregătită pe ramura izolată `F09-T08-imagini-paginare`, cu categoriile F07-T11 și cererile de concediu F07-T12. Include imagini și rapoarte v2. Actualizarea păstrează modelele definite în proba 0.20.0 și completează migrarea cererilor care lipsea din acea probă; backupul automat rămâne înaintea aplicării migrărilor.

## Definire rapoarte

Deschide meniul **Definire rapoarte**. Sunt pregătite trei modele: Pontaj lunar, Fișă individuală și Evidența absențelor. **Adaugă model de raport** creează un model nou; **Duplică** copiază setările într-o definiție independentă.

Fiecare raport are trei componente proprii:

1. **Antet** — texte și imagini stânga, centru și dreapta, de exemplu sigla unității, titlul și perioada.
2. **Conținut** — sursa datelor, textul introductiv repetat înaintea tabelului pe pagini, coloanele și ordinea lor. Bifează coloanele dorite și folosește săgețile pentru reordonare.
3. **Subsol** — texte și imagini stânga, centru și dreapta pentru explicații, semnături sau aprobări.

În fiecare zonă de antet/subsol poți încărca o imagine **PNG sau JPEG**, împreună cu textul. Importul păstrează conținutul în modelul din baza criptată; fișierul original poate fi mutat ulterior fără să afecteze raportul. Imaginile sunt comune limbilor modelului. Alege alt fișier pentru înlocuire sau **Elimină imaginea** pentru eliminarea din model, apoi salvează.

**Înălțime maximă (mm)** acceptă 5–35 mm, implicit 15 mm. Imaginea își păstrează proporțiile și încape în lățimea zonei, fără decupare. Fiecare imagine poate avea cel mult 2 MiB, 4096 px pe axă și 4 milioane de pixeli. PNG permite transparență. SVG și URL-urile externe nu sunt acceptate.

Completează variantele română/engleză folosind selectorul limbii textelor. Câmpurile afișate, precum `{{perioada}}`, `{{sectie}}`, `{{angajat}}` sau `{{generat_la}}`, se înlocuiesc la generare. Sunt acceptate numai câmpurile listate. Textele sunt texte simple; HTML-ul, SQL-ul sau codul nu se execută.

Pagina se configurează separat pentru fiecare raport: A4/A3, portret/peisaj, margini între 5 și 25 mm și font de 6–14 pt. Pentru grila cu 31 de zile, punctul de pornire este A3 peisaj. Antetul și subsolul se repetă pe pagini, împreună cu numărul paginii.

Fiecare raport poate avea **una sau mai multe pagini**. Conținutul continuă automat când pagina se umple; nu trebuie creat câte un model pentru fiecare pagină. Antetul, subsolul, imaginile și capul de tabel se repetă pe toate paginile; numerotarea este „Pagina X din Y”. Previzualizarea afișează numărul real de pagini. Dimensiunile imaginilor și textele lungi reduc spațiul disponibil conținutului.

Salvează modelul înainte de generare. Un model inactiv rămâne în listă și poate fi reactivat prin editare, dar nu apare la generare. La modificări concurente, salvarea refuză versiunea veche și păstrează formularul pentru a putea copia textele înainte de reîncărcare.

## Rapoarte

1. Deschide **Rapoarte**, alege un model activ, anul/luna și limba documentului.
2. Opțional filtrează secția, departamentul și angajatul. Fișa individuală cere un angajat. Evidența absențelor folosește întreg anul.
3. Apasă **Generează și previzualizează**. PDF-ul reflectă modelul și datele citite împreună la acel moment; schimbarea filtrelor retrage previzualizarea veche.
4. **Salvează PDF** deschide dialogul de salvare, unde alegi locul și numele. Anularea dialogului nu salvează fișierul.
5. **Deschide pentru tipărire** deschide exact PDF-ul previzualizat în aplicația asociată din Windows. Alege acolo **Tipărire / Ctrl+P**, imprimanta, paginile și numărul de exemplare. Nu se trimite automat o comandă către imprimanta implicită.

Copiile deschise pentru tipărire se păstrează în `karios/rapoarte/` de lângă executabil. Salvarea în alt folder este o acțiune explicită din dialog. Rapoartele generate sunt documente lizibile și se păstrează în locația aleasă de utilizator.

## Ce conțin primele rapoarte

- **Pontaj lunar:** centralizatorul existent afișează orele, cu weekendurile și sărbătorile marcate. Codurile programelor se păstrează în date pentru calcule. Angajații inactivi apar dacă au înregistrări sau ore continuate în perioada cerută.
- **Pontaj detaliat:** raport lunar separat, implicit A3 peisaj, cu zilele lunii și cele 12 coloane finale: Ore lucrate, T, ZNL, CO, CEV, CFS, Ore norm, Form prof, CM, L/FIN, Deleg și Ore gardă 200%. Grupul totalurilor este obligatoriu la sfârșit; antetul, subsolul și imaginile se configurează independent.
- **Fișă individuală:** zilele lunii, orele și descrierile programelor. Un cod istoric neinterpretat produce „—”, fără a inventa ore. Absențele au zero ore din propria înregistrare; pot exista ore continuate dintr-o tură începută anterior.
- **Evidența absențelor:** zilele active înregistrate explicit în anul ales, cu tipul și categoria lor. Zilele anulate nu sunt incluse.

Documentele sunt marcate ca **documente de lucru**, deoarece grila curentă păstrează schițe. Nu deduc salarii, sporuri, solduri de concediu sau totaluri de plată. Aceste reguli și centralizatorul salarial final așteaptă F01.

Subsolul tehnic D009 rămâne prezent pe lângă subsolul configurabil: numele/versiunea raportului, revizia definiției, versiunea aplicației, starea regulilor, momentul generării și actorul local până la implementarea conturilor.

## Verificări și licențe

Modelele se testează pe SQLite criptat; sunt verificate independența definițiilor, reviziile, refuzul salvărilor concurente și jurnalul atomic. Referințele curente sunt centralizator v3, fișă individuală v3, pontaj detaliat v2 și absențe v2; referințele anterioare sunt păstrate. PDF-urile sunt validate cu pdfcpu, iar imaginile PNG/JPEG sunt extrase și verificate pe fiecare pagină. Fonturile cu diacritice sunt incluse.

Pentru PDF-uri fictive de probă și verificarea UI:

```powershell
go test ./internal/rapoarte -run 'TestReferinteRapoarteSiPDF|TestPaginarePDFCuAntetSiSubsol|TestPDFImaginiPeFiecarePagina' -director-demo <cale-absoluta-catre-build/bin/rapoarte-demo>
cd frontend
npm run build
cd ..
node verificari/interfata.mjs
```

Referințele sunt inițializate numai în acest task prin `-actualizeaza-rapoarte`; ulterior o schimbare cere verificarea și creșterea versiunii raportului conform D009. `verificari/randare-pdf.ps1` poate reda o pagină în PNG cu motorul PDF inclus în Windows, fără instalări.

Interacțiunile UI se verifică în Edge cu punte Wails simulată; backendul PDF se testează separat, pe date fictive reale. Testele nu trimit documente către imprimantă. Bibliotecile și notificările de licență sunt în [D025](decizii/D025-definitii-rapoarte.md) și [licențe](licente/rapoarte/README.md).
