# Analiza UX · Pontix 0.92.59

Data: 2026-11-07 · Autor: Claude · Task: [F16-T36](functionalitati/F16-diagnoza-si-optimizari-ui.md) · Cerere: Ioan Tuns („fă o analiză ux, pro, cons, recomandări”).

## Metodă și limite

- Capturi Edge headless cu puntea simulată și date fictive (10 angajați demo), pe 7 ecrane principale, la **1024×640** (ecranul țintă din analizele anterioare) și 1440×900. Scenariu nou, reluabil: `node verificari/interfata.mjs --capturi-ux` (capturi și `masuri-ux.json` în `build/bin/verificare-ui-…`).
- Măsurători: înălțimea paginii și derularea orizontală pe fiecare ecran; numărarea euristică a textelor românești fără diacritice din `ro.json`.
- **Nu** s-au observat utilizatori reali și nu s-au folosit date reale. Performanța pe 300 de angajați e măsurată separat (F06-T13). Ecranele secundare (dialoguri, editorul pdfme, planul AI) au fost analizate în rapoartele anterioare și nu se repetă aici.

| Ecran (1024×640) | Înălțimea paginii | Derulare orizontală |
| --- | --- | --- |
| Angajați | 1067 px | nu |
| Pontaj | 1121 px | nu |
| Planificare | 1736 px | nu |
| Concedii și absențe | 974 px | nu |
| Rapoarte | 640 px | nu |
| Asistent AI | 640 px | nu |
| Setări (General) | 2466 px | nu |

## Ce funcționează bine (pro)

1. **Structură clară și stabilă.** Bara laterală are 7 intrări cu pictograme, Setări are submeniu și file (D034), fiecare ecran are breadcrumb, titlu și o frază care spune la ce folosește.
2. **Nicio derulare orizontală a paginii**, la nicio lățime măsurată; tabelele largi derulează în interiorul lor, cu antet și coloana de nume fixe.
3. **Grila Pontaj e matură:** legendă vizibilă (weekend, sărbătoare, selecție, modificat, avertisment, înlocuire), bara de inspecție a celulei, densități Compact/Normal/Aerisit, contor „celule nesalvate”, tastatură și meniu contextual, motivele valorilor „—” (F05-T20).
4. **Explicații în locul presupunerilor:** texte de ajutor lângă câmpuri, „Cum să faci” pe fiecare ecran, avertismente care spun ce lipsește și unde se corectează, nicio regulă de plată inventată.
5. **Siguranța datelor e vizibilă:** „Date păstrate local”, confirmare la părăsire cu modificări nesalvate, previzualizare înainte de salvare la absențe și planul AI, starea „doar citire” când altă instanță are datele.
6. **Rapoarte simple de folosit:** un singur formular (model, an, limbă, secție, angajat) și un buton principal.
7. **Asistentul AI e prudent:** totul local, propunerile cer confirmare, planul are stare și reluare (F14-T35).

## Probleme (contra)

Gravitate: **P1** = afectează munca zilnică, **P2** = încetinește sau derutează, **P3** = finisare.

| ID | Gravitate | Constatare | Dovadă |
| --- | --- | --- | --- |
| UX59-01 | P1 | **Antetul fiecărei pagini consumă 200–270 px din 640.** „Activitățile lunii”, zoomul (− 100% +) și „Cum să faci” stau pe rânduri separate de titlu; la Angajați, Planificare, Absențe și Asistent, „Cum să faci” are un rând propriu. | Capturi `*-1024.png`: conținutul începe la y≈250–270. |
| UX59-02 | P1 | **În Pontaj, la 1024×640 nu se vede niciun rând complet de angajat fără derulare.** Între antet și grilă stau: banda „Selecție…”, filtrul Stare, Filtre profesionale, perioada, două rânduri de butoane, densitatea, bara de inspecție și legenda. Primul angajat începe la y≈600. | `pontaj-1024.png`. |
| UX59-03 | P1 | **Totalurile noi nu se văd fără derulare orizontală.** OL e fixată la dreapta; CO, CM, CFP, CVE, DLG, FMP sunt după zilele lunii, deci în afara ecranului și la 1440 px. | `pontaj-1440.png`: se vede doar OL după ziua 16. |
| UX59-04 | P2 | **Texte fără diacritice în interfață:** cel puțin 115 din 2145 de texte românești (≈5%, numărare euristică), de ex. „Precompleteaza zile lucratoare”, „Verificarea configurarii si pasii recomandati”, „Selectie”, „Stare angajati”, „Nesalvate in plan”. Amestecul cu texte corecte arată neîngrijit. | `ro.json`; capturile Pontaj, Planificare, Setări, Absențe. |
| UX59-05 | P2 | **Două „zoomuri” diferite pe același ecran:** mărirea paginii (− 100% +) în antet și densitatea grilei („Zoom: Compact / Normal / Aerisit”). Același cuvânt pentru două lucruri. | `pontaj-1024.png`. |
| UX59-06 | P2 | **Date în format tehnic:** banda de selecție („2026-10-01 - 2026-10-01. Marci: DEMO-001”), bara de inspecție („2026-10-01”), cuvântul „mărci” în loc de persoane. Rapoartele și planul AI folosesc deja formatul „1 oct. 2026”. | `pontaj-1024.png`, `planificare-1024.png`. |
| UX59-07 | P2 | **Bara de unelte din Pontaj e aglomerată și fără ierarhie:** 9 butoane pe 3 rânduri, „Salvează schița” (acțiunea principală) în mijloc, „Închide luna” la dreapta cu aceeași greutate vizuală, „Precompletează” dezactivat fără explicație. | `pontaj-1024.png`. |
| UX59-08 | P2 | **Setări › General are 2466 px (aproape 4 ecrane)**: normă, secții, compartimente, nomenclatoare, prefixe și altele pe o singură pagină lungă, cu verificarea configurației pliată sus. | `masuri-ux.json`, `setari-1024.png`. |
| UX59-09 | P2 | **Asistentul AI are puțin loc pentru conversație la 1024×640** (≈130 px vizibili), din cauza antetului (UX59-01) și a barei motorului; în starea „Inactiv (fișiere lipsă)” nu există un pas următor evident. | `asistent-1024.png`. |
| UX59-10 | P3 | **Rândul „Persoane înregistrate” din Pontaj arată „—” în fiecare zi** când nu există date; ocupă un rând fără informație. | `pontaj-tabel-1024.png`. |
| UX59-11 | P3 | **Filtrele arată diferit de la un ecran la altul:** „Stare [Active]” singur pe un rând în Pontaj, „Stare angajați” în Planificare, filtre pe grilă la Angajați și Absențe. | Capturile ecranelor. |
| UX59-12 | P3 | **Butoane dezactivate fără motiv afișat** (Export Excel în Rapoarte, Precompletează în Pontaj, Trimite în Asistent). | `rapoarte-1024.png`, `pontaj-1024.png`, `asistent-1024.png`. |

## Recomandări, în ordine

| # | Rezolvă | Recomandare | Efort | Criteriu de acceptare |
| --- | --- | --- | --- | --- |
| R1 | UX59-01, 09 | **Antet compact comun:** titlul și descrierea pe stânga, „Activitățile lunii”, „Cum să faci” și mărirea pe un singur rând la dreapta; mărirea poate trece în subsolul barei laterale, lângă limbă și temă. | M | Pe toate ecranele, conținutul începe sub y=150 la 1024×640; nimic nu dispare. |
| R2 | UX59-02, 07 | **Pontaj pe un rând de unelte:** perioada + „Salvează schița” (principal) + „Completare” + „Mai multe ▾” (Precompletează, După categorii, Export, Asistent, Ajutor); „Închide luna” separat, ca acțiune de final de lună; banda „Selecție” mutată în bara de inspecție. | M | La 1024×640, minimum 5 rânduri de angajați vizibile fără derularea paginii; toate acțiunile rămân accesibile cu tastatura. |
| R3 | UX59-03 | **Totalurile vizibile împreună:** cele 7 coloane fixate la dreapta (în densitatea Compact, cu lățimea redusă), sau un comutator „Zile / Totaluri”. | S–M | La 1024 și 1440 px, OL…FMP se văd fără derulare orizontală. |
| R4 | UX59-04 | **Trecere pe diacritice în `ro.json`**, cu un test care semnalează textele românești fără diacritice (lista de excepții pentru coduri). | S | Testul trece; 0 texte semnalate. |
| R5 | UX59-05, 06 | **Termeni și date unitare:** „Mărire” (pagina) vs „Densitate” (grila); toate datele afișate prin formatterul comun („1 oct. 2026”); „persoane” în loc de „mărci” în textele pentru utilizator. | S | Nicio dată ISO în textele interfeței (verificare în scenariile Edge). |
| R6 | UX59-08 | **Setări › General pe secțiuni:** file sau ancore (Normă, Structură, Nomenclatoare, Prefixe), cu verificarea configurației deschisă când are probleme. | M | Fiecare secțiune încape în ≤2 ecrane la 1024×640. |
| R7 | UX59-12 | **Motivul la butoanele dezactivate** (tooltip și text accesibil): „Alege mai întâi o selecție”, „Generează întâi raportul”. | S | Fiecare buton dezactivat are motiv citit de cititorul de ecran. |
| R8 | UX59-09 | **Asistent pe toată înălțimea:** conversația ocupă spațiul rămas; în starea „fișiere lipsă”, un buton „Cum instalez motorul” către ajutor. | S | ≥300 px pentru conversație la 1024×640. |
| R9 | UX59-10, 11 | **Finisări:** rândul „Persoane înregistrate” ascuns când e gol; componenta de filtre comună pe toate ecranele. | S | — |

**Ordinea propusă:** R4 și R5 (rapide, vizibile peste tot), apoi R1 (câștig pe toate ecranele), R2 și R3 (ecranul folosit zilnic), apoi R6–R9.

## Ce nu s-a verificat

- Fluxuri complete cu un operator real și pe date reale (300 de angajați, luni închise).
- Tema întunecată pe toate ecranele (verificată anterior doar în Pontaj și Planificare).
- Scalările Windows de 125–150% în WebView2 (s-au măsurat doar zoomurile 150% și 200% din aplicație).
