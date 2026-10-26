# Configurarea orelor și totalurilor

În **Configurare pontaj** alegi repartizarea turelor și regulile folosite de grilă, centralizator, pontajul detaliat și fișa individuală. Salvarea este globală, are revizie și jurnal; rapoartele deja exportate rămân documente ale configurației folosite la generare.

## Ture peste miezul nopții

- **Toate orele în ziua începerii:** o tură vineri 22–07 apare vineri cu 9 ore și intră integral în luna începerii.
- **Împărțire pe zile calendaristice:** aceeași tură apare cu 2 ore vineri și 7 sâmbătă. Trecerea între luni include continuarea din luna precedentă și exclude orele din luna următoare.

În ambele moduri, ZNL clasifică fragmentele efective 00–24 din weekenduri și zilele libere din calendar. Ore norm clasifică luni–vineri fără acele zile. Exemplul de vineri are T=9, ZNL=7 și Ore norm=2. Adaugă zilele libere personalizate și sărbătorile anului următor în **Sărbători legale** când o tură trece peste 31 decembrie.

Duratele sunt nominale, conform programelor. Nu se ajustează automat schimbarea orei și nu se poziționează pauze. Un program cu pauză nenulă rămâne necalculabil până la stabilirea repartizării pauzei. Un program fără oră de început poate furniza numai totalul în modul zilei de început.

## Coloane finale

Implicit sunt definite Ore lucrate/T = total ore, ZNL, Ore norm și zilele active CO/CM. Pentru celelalte coloane rezultatul este „—”. Zero înseamnă o regulă definită fără înregistrări relevante.

Poți alege total ore, ore ZNL, ore norm, ore din programe selectate sau zile din tipuri de absență selectate. Selecția acceptă mai multe opțiuni (Ctrl/clic în Windows); absențele pot fi limitate la cele aprobate. Coloanele de ore acceptă surse de ore, iar coloanele de zile acceptă absențe; Form prof și L/FIN preiau unitatea sursei alese. Pentru ore, factorul explicit poate fi 0,01–10; 1 păstrează orele fără ponderare. „Ore gardă 200%” nu are automat sursă sau factor 2. Configurațiile concrete se completează după primirea regulilor de la contabilitate (F09-T10/F01).

Salvările concurente sunt refuzate; reîncarcă apoi reaplică modificările. Codurile redenumite sau eliminate din nomenclator trebuie reselecționate în reguli, fără conversii automate.

## Completarea grilei

Selectează una sau mai multe celule, alege programul după denumire, interval și regim, apoi apasă **Aplică programul**. Alternativ introdu **Ore manuale pe zi** și apasă **Aplică orele manuale**; nu se deduce un program din numărul de ore. F2 sau dublu clic focalizează câmpul manual. Orele manuale aparțin zilei selectate; sunt acceptate 0–48 ore, cu precizie de un minut și separator zecimal punct sau virgulă.

Grila și copierea afișează ore; codul programului rămâne în înregistrare. Lipirea TSV acceptă ore, nu coduri. Calculul modificărilor este o previzualizare; apasă **Salvează schița** pentru persistență. Câmpul manual gol sau Delete șterge înregistrările din zilele selectate. Orele continuate din ziua anterioară se modifică la data începerii turei, indicată în explicația celulei. Celulele cu absențe active sunt protejate.

Un cod istoric necunoscut apare „—”; înlocuiește-l explicit cu programul corect sau cu ore manuale. O absență activă anulează programul propriei zile. Orele continuate dintr-o tură anterioară rămân vizibile; compatibilitatea cu concediul nu este decisă automat.
