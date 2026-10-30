# Manualele fiecărei livrări

F02-T18. Sursa editabilă este `docs/manual-utilizare-ro.md`, cu marcaje pentru versiune, ediție, data buildului și commit. Documentele generate stau lângă executabil, în directorul nou creat de `build.ps1`, nu peste sursă sau manualele unei versiuni precedente.

## La fiecare versiune

1. Actualizează instrucțiunile pentru comportamentul versiunii noi.
2. Trece versiunea aplicației în `versiune_revizuita` din `docs/manual/manual.json`. Buildul refuză o valoare diferită de versiunea aplicației.
3. Crește `editie` când schimbi conținutul editorial; aceasta este distinctă de versiunea aplicației. Data și commitul sunt luate automat din build.
4. Rulează `python verificari/manuale.py`, apoi `./build.ps1`.

Fiecare build generează `manual-utilizare-ro-X.Y.Z.md`, `.html` și `.pdf`. Manifestul `build.json` conține versiunea, commitul, ediția, data și SHA256 pentru fiecare document. `stare: complet` se scrie numai după verificarea executabilului și a tuturor manualelor. La eroare rămâne `stare: esuat`; acel director nu se distribuie ca livrare completă.

Generatorul folosește biblioteca standard Python și Microsoft Edge existent, în mod headless, pentru PDF. Nu descarcă programe sau pachete. Sunt necesare comanda `python` și Edge în locația Windows standard. Pentru rularea directă se poate specifica `--browser`. Profilul browserului este izolat în temporarele buildului; nu se folosește profilul personal. Python și Edge sunt unelte de construire, nu sunt necesare pentru rularea aplicației. HTML-ul este autonom și se deschide fără internet; PDF-ul se citește cu un cititor PDF compatibil.

## Previzualizare editorială

```powershell
python docs/manual/genereaza.py --versiune 0.88.13 --commit proba --data 2026-10-30 --destinatie build/bin/manual-proba-nou --fara-pdf
```

Destinația trebuie să nu conțină deja manuale cu aceleași nume. `--fara-pdf` este numai pentru previzualizare; buildul oficial generează obligatoriu toate cele trei formate. Convențiile Markdown suportate sunt cele folosite în manual: titluri, paragrafe, liste, tabele, citate, text aldin și cod inline.

Deschide HTML-ul cu `?verificare=1` pentru proba locală de căutare și cuprins; succesul este indicat în atributul `data-verificare="trecut"`. Proba verifică filtrarea, lipsa rezultatelor, resetarea și țintele cuprinsului. Nu execută operații în aplicație.

## Limitări

Regenerarea asigură identitatea documentelor livrate; nu poate revizui automat corectitudinea instrucțiunilor. Revizuirea editorială rămâne obligatorie la fiecare versiune. Din F02-T21, articolele interne sunt comune panoului de ajutor, răspunsurilor AI și capitolului 20 al manualului. Se editează metadatele în internal/ajutor/catalog.json și textele în cataloagele i18n RO/EN. Manualul nu inventează reguli de salarizare și nu conține date operaționale.
