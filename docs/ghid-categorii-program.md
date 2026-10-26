# Categorii de program și precompletare · 0.29.0

## Configurarea categoriilor

Deschide **Setări → Categorii de program**. Sunt disponibile șase categorii inițiale editabile:

| Categorie | Programe inițiale | Implicit L–V |
| --- | --- | --- |
| Dimineață | Z7 | Z7 |
| După-amiază | S7 | S7 |
| Gărzi | IN9, Z12, N12, G24 | fără implicit |
| Dimineață și gărzi | Z7, IN9, Z12, N12, G24 | Z7 |
| După-amiază și gărzi | S7, IN9, Z12, N12, G24 | S7 |
| Program mixt | Z7, S7, N9, IN9, Z12, N12, G24 | fără implicit |

Setul folosește programele existente, cu intervalele lor păstrate. Dacă unul a fost eliminat sau dezactivat, nu este recreat prin această funcție. Poți adăuga categorii proprii, redenumi sau dezactiva categoriile, bifa programele permise și alege un implicit dintre ele. **Fără implicit** păstrează alegerea manuală; nu generează rotații. După un conflict de salvare, închide formularul și folosește **Reîncarcă categoriile** înainte de reaplicare.

Categoriile de program ale angajaților sunt distincte de categoriile **LAB / CUSTOM** ale schimburilor. Gărzile se aleg explicit; combinarea schimburilor într-o categorie nu distribuie automat nopțile, weekendurile sau gărzile pe persoane.

## Atribuirea angajaților

În **Angajați → Editează**, alege **Categoria de program**, apoi salvează. Categoria apare și sub contract în listă. Angajații existenți rămân inițial fără categorie: nu se deduce programul lor din nume, funcție sau pontajul vechi.

Precompletarea folosește categoria și norma curente, inclusiv când alegi o lună trecută. Datele deja înregistrate nu se convertesc. Pentru norma parțială, definește un schimb cu durata individuală potrivită și folosește-l în categoria angajatului; un implicit de 7 ore nu completează automat o normă de 4 ore. Același principiu se aplică dacă norma întreagă globală diferă de durata implicitului.

O categorie dezactivată rămâne vizibilă la angajații deja atribuiți, dar nu poate fi atribuită altora și nu mai produce propuneri automate.

## Completarea pontajului

Butonul **Precompletează L–V** este din nou în bara vizibilă. Pentru un program unic, alege întâi programul în **Completare**, apoi apasă butonul. Sunt completate numai celulele goale din rândurile vizibile, fără weekenduri, sărbători sau absențe active. Angajații cu categorie sunt omiși dacă programul ales nu este permis sau durata diferă de norma lor. **Precompletează selecția L–V** rămâne în Completare pentru rândurile selectate.

Pentru programe individuale:

1. Opțional, în **Filtre**, alege secția, departamentul și **Categoria de program**. Poți vedea separat persoanele **Fără categorie de program**.
2. Apasă **După categorii**. Previzualizarea folosește programul implicit al fiecărui angajat din rândurile vizibile. Bifează **Numai angajații din selecția curentă** dacă vrei doar persoanele din dreptunghiul selectat; sunt vizate toate zilele eligibile ale lunii pentru acele persoane.
3. Verifică numărul propunerilor și explicațiile zilelor omise. Tabelul afișează primele 100 de propuneri; aplicarea include întregul număr indicat.
4. Apasă **Aplică în schiță**, apoi verifică grila. **Anulează pasul / Refă pasul** operează asupra întregului lot.
5. Apasă **Salvează schița** pentru păstrarea în bază. Anularea previzualizării nu schimbă datele; aplicarea în schiță nu salvează automat.

## Protecțiile completării pe categorii

- Nu înlocuiește celule existente, modificări locale sau zile cu pontaj efectiv.
- Ocolește weekendurile, sărbătorile și zilele libere personalizate din calendar.
- Respectă toate zilele cu absență activă: concediu, medical și orice tip special definit, indiferent de starea aprobării. Absențele anulate nu blochează ziua.
- Verifică suprapunerile cu turele salvate sau din schița curentă, inclusiv din luna precedentă. O tură care ar continua într-o zi cu absență, într-un weekend sau într-o zi liberă nu este propusă.
- Lasă pentru verificare manuală ziua următoare unei gărzi de cel puțin 24 h și intervalele neclare. Aceasta este o măsură conservatoare de precompletare, nu o regulă salarială sau o stabilire a repausului legal.
- Omite programele inactive, fără oră de început, cu pauză nerepartizată sau cu durată diferită de norma individuală. Gărzile și intrările speciale nu sunt deduse automat.

Aplicarea și salvarea unui program cunoscut din afara categoriei active sunt refuzate. Înregistrările istorice rămân păstrate; orele manuale rămân o introducere explicită, fără deducerea unui schimb din numărul de ore. Verificările de absențe și reviziile concurente existente rămân active la salvare.
