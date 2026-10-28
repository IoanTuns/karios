# D036 · Categorii de program și precompletare pe angajat

Stare: acceptată · Data: 2026-10-26 · Propusă de: Codex · Decisă de: Ioan Tuns · Task: F07-T14

Utilizatorul cere explicit categorii de program pentru angajații care lucrează numai dimineața, numai după-amiaza, gărzi sau combinații, pentru reducerea muncii de completare, cu respectarea concediilor, medicalelor și altor zile libere înregistrate. Semnalează și un default dispărut după compactarea UI; clarificarea exactă este cerută separat.

Modelul propus: nomenclator de categorii de program, cu denumiri RO/EN, lista programelor selectate și opțional un program implicit pentru zile lucrătoare. Fiecare angajat poate avea o categorie; combinațiile sunt categorii cu mai multe programe. Categoriile sunt distincte de LAB/CUSTOM din nomenclatorul schimburilor. Angajații existenți rămân fără atribuire automată. Configurarea păstrează revizii și jurnal.

Precompletarea pe categorii este o acțiune explicită și revizuibilă, completează numai zilele goale ale angajaților activi selectați/vizibili, exclude weekendurile și zilele libere din calendar și respectă toate absențele active, indiferent de tip. Nu se alterează înregistrările existente. O categorie fără implicit (de exemplu numai gărzi), un program inactiv/incomplet sau o durată nepotrivită normei individuale cer alegere manuală; nu se inventează un program. Se evită suprapunerea turelor și trecerea peste o zi cu absență.

Programul implicit nu înseamnă rotație. Combinațiile nu distribuie automat gărzile fără un model explicit primit de la utilizator. Prima livrare poate funcționa fără reguli de rotație; acestea rămân în F11 dacă sunt cerute ulterior. Calculele și raportarea orelor continuă să folosească programele salvate.

Confirmare explicită: defaultul lipsă este „Butonul de precompletare L–V”, care revine în bara vizibilă. Modelul confirmat este „Programe permise + implicit; gărzile se aleg explicit”. Categorii inițiale editabile: Dimineață, După-amiază, Gărzi, Dimineață și gărzi, După-amiază și gărzi, Program mixt; angajații existenți nu primesc categorie dedusă.
