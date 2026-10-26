# D032 · Îmbunătățiri de ergonomie în interfață: bară de inspecție a celulei, Undo/Redo în pontaj, coloană T fixă, filtre de secțiune și selector de limbă

Stare: propusă · Data: 2026-10-26 · Propunere: Antigravity · Legat de: F08-T08, I032, I033, I034, I035, I036, I037

## Context

În urma analizei ergonomiei interfeței `Pontix 0.24.0`, au fost identificate șase direcții concrete de îmbunătățire a experienței de lucru zilnice:

1. **Vizibilitatea detaliilor celulei din grila de pontaj (`I032`):** codul intern al schimbului (`Z7`, `N12`, `G24`), detaliile absenței (stare aprobare, număr de înregistrare, înlocuitor), colegii înlocuiți și avertismentele erau vizibile doar în atributul `title` la survol cu mouse-ul, fiind inaccesibile la navigarea din tastatură; în plus, aplicarea unui schimb cerea două acțiuni (alegere din listă + clic pe `Aplică programul`).
2. **Anularea și refacerea pașilor în schița de pontaj (`I033`):** o lipire TSV sau o precompletare pe mai multe celule nu putea fi anulată rapid cu `Ctrl+Z` / `Ctrl+Y`.
3. **Comutarea limbii interfeței (`I034`):** deși cataloagele RO și EN sunt încărcate la pornire în `react-i18next`, lipsea un selector vizibil în bara laterală pentru schimbarea limbii pe loc.
4. **Vizibilitatea totalului `T` la derularea orizontală (`I035`):** în grila cu 31 de zile și 12 coloane de total, coloana `T` ieșea din ecran la derularea primelor zile ale lunii.
5. **Ergonomia secțiunilor din `Concedii și absențe` și a selecțiilor multiple din `Configurare pontaj` (`I036`):** afișarea simultană a celor trei tabele din `Concedii și absențe` cerea derulare verticală lungă, iar `<select multiple>` din `Configurare pontaj` cerea `Ctrl+Click`.
6. **Diferențierea vizuală în meniul lateral (`I037`):** mai multe ecrane împărțeau aceeași pictogramă SVG în bara laterală.

## Variante

1. **Păstrarea interfeței actuale** până la finalizarea tuturor funcționalităților de calcul (`F01`, `F04`) și de conturi (`F15`).
2. **Implementarea incrementală a îmbunătățirilor de ergonomie fără schimbarea regulilor de calcul sau a schemei bazei de date:**
   - Bară de inspecție a celulei active (`.inspectie-celula`) și paletă de butoane rapide pentru schimburile active în `Pontaj.tsx`.
   - Stivă `Undo / Redo` (`Ctrl+Z` / `Ctrl+Y` și butoane în bara de unelte) pentru modificările din sesiunea curentă a schiței.
   - Înghețarea coloanei `T` în dreapta (`position: sticky; right: 0`) în tabelul grilei de pontaj.
   - Comutator de secțiuni în `Absente.tsx` (`Toate secțiunile`, `Zile înregistrate`, `Sumar pe angajat`, `Tipuri de absență`) și butoane-cip de comutare rapidă lângă `<select multiple>` în `ReguliPontaj.tsx`.
   - Selector de limbă `RO / EN` în subsolul barei laterale și pictograme SVG distincte pentru toate cele 9 ecrane din `App.tsx`.

## Decizie propusă

Se alege **varianta 2** (taskul `F08-T08`):

1. **Bară de inspecție și paletă rapidă de schimburi (`I032`):** `Pontaj.tsx` afișează sub bara de unelte panoul `.inspectie-celula` cu detaliile celulei aflate la capătul selecției curente (angajat, dată, cod schimb / ore manuale, ore calculate, absență, înlocuire și avertismente), actualizat instantaneu la navigarea din tastatură sau cu mouse-ul, împreună cu butoane rapide pentru aplicarea dintr-un clic a oricărui schimb activ pe selecția curentă.
2. **Istoric `Undo / Redo` în schiță (`I033`):** `aplicaCelule` păstrează stiva de stări anterioare (`istoricInapoi`, `istoricInainte`, până la 50 de pași), accesibilă prin `Ctrl+Z` / `Ctrl+Y` în grilă și prin butoanele `Anulează pasul` / `Refă pasul`, resetată la salvarea sau reîncărcarea lunii.
3. **Selector de limbă în bara laterală (`I034`):** în subsolul barei laterale din `App.tsx` se adaugă un `<select name="limba_interfata">` care schimbă pe loc limba interfeței (`ro` / `en`) prin `i18n.changeLanguage`. Persistența în baza de date rămâne în `F03-T05` / `F05-T04`.
4. **Coloană `T` fixă la derulare (`I035`):** celulele `<th>` și `<td>` ale coloanei `T` din `.grila-pontaj` primesc clasa `fix-total-t` (`position: sticky; right: 0`), rămânând vizibile în dreapta pe tot parcursul derulării orizontale.
5. **Filtre de secțiune în `Absente.tsx` și comutatoare rapide în `ReguliPontaj.tsx` (`I036`):** `Absente.tsx` permite comutarea între vizualizarea completă (implicită) și fiecare dintre cele trei secțiuni (`zile`, `sumar`, `tipuri`), iar `ReguliPontaj.tsx` adaugă butoane-cip pentru bifarea/debifarea cu un singur clic a codurilor din regulă, păstrând și `<select multiple>`.
6. **Pictograme SVG distincte (`I037`):** `Simbol` din `App.tsx` oferă câte o pictogramă proprie pentru fiecare dintre cele 9 ecrane.

