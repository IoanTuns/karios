# D031 · Precompletarea lunii, avertismente vizuale în grila de pontaj și copie de siguranță la cerere

Stare: propusă · Data: 2026-10-26 · Propunere: Antigravity · Legat de: F08-T07, F05-T16, I029, I030, I031

## Context

După implementarea verificărilor și previzualizării pentru cererile de concediu ([D030](D030-suprapuneri-inlocuitor-previzualizare-cereri.md)), în lista de idei din `docs/idei.md` au rămas trei îmbunătățiri operaționale cerute pentru fluxul zilnic de lucru:

1. **Precompletarea rapidă a lunii în grila de pontaj (`I029`):** completarea zilelor lucrătoare luni–vineri pentru toți angajații afișați cerea selectarea manuală a fiecărui interval între weekenduri, sărbători legale și zile cu absență.
2. **Avertismente vizuale neblocante în grila de pontaj (`I030`, în legătură cu `I004`):** utilizatorul are nevoie să observe dintr-o privire dacă o persoană are pontat lucru imediat după o gardă de 24 h (`G24`), dacă o zi depășește 24 de ore calculate sau dacă totalul lunar de ore (`T`) diferă de norma lunară (`ore_norma`), fără a inventa sau bloca vreo regulă de calcul salarial.
3. **Copie de siguranță la cerere din Setări (`I031`):** aplicația creează deja copii criptate înainte de migrări (`F05-T06`), dar utilizatorul nu avea un buton în interfață pentru a salva pe loc o copie criptată înaintea unor modificări masive în pontaj sau concedii.

## Variante

1. **Așteptarea închiderii lunii (`F12`) și a confirmării complete `F01`** pentru orice unelte de asistare în grilă și de backup manual.
2. **Adăugarea de unelte neblocante în grilă și a backupului manual criptat în Setări:**
   - Buton `Precompletează L–V` în bara de unelte a grilei de pontaj, care aplică schimbul ales doar în celulele goale din zilele de luni până vineri care nu sunt sărbători legale și nu au absențe active.
   - Evidențiere vizuală neblocantă (`celula-avertisment` și `total-diferenta`) cu detalii în `title` pentru celulele cu lucru după `G24`, zilele cu peste 24 h calculate și rândurile unde `T` diferă de `ore_norma`.
   - Secțiune `Copie de siguranță` în `Setări` care apelează `stocare.BackupManual` (`VACUUM INTO` cu cheia din seif) și salvează fișierul `pontix-manual-<data-ora>-v<versiune>.db` în folderul `karios/backup/`, fără a afecta rotația celor `MaxBackupuri` copii automate la migrare.

## Decizie propusă

Se alege **varianta 2**:

1. **Precompletarea lunii (`F08-T07`, `I029`):** în `Pontaj.tsx`, butonul `Precompletează L–V` (`pontaj.precompleteaza_lv`) parcurge rândurile vizibile (filtrate după secție și departament) și completează cu `programCompletare` toate celulele goale din zilele luni–vineri care nu sunt sărbători legale în calendarul anului și nu au absențe active. Modificările intră în schița curentă (cu previzualizare și salvare explicită).
2. **Avertismente vizuale în grila de pontaj (`F08-T07`, `I030`):**
   - Celula primește clasa `celula-avertisment` și mesaj în `title` dacă în ziua precedentă din lună este completat codul `G24` iar în ziua curentă există program/ore lucrate (`pontaj.avertisment_dupa_g24`), sau dacă totalul calculat pe zi depășește 24 de ore (`minute > 1440`, `pontaj.avertisment_peste_24h`).
   - Celula de total `T` primește clasa `total-diferenta` și mesaj în `title` (`pontaj.diferenta_norma`) când atât `T`, cât și `ore_norma` sunt definite în ore și au valori diferite.
   - Toate avertismentele sunt exclusiv vizuale și nu modifică niciun calcul din `internal/calcul`.
3. **Backup manual criptat (`F05-T16`, `I031`):** `stocare.BackupManual` deschide seiful din folderul de date, creează în `backup/` o copie criptată AES-256-XTS cu `bd.Backup` (`VACUUM INTO`) numită `pontix-manual-<UTC>-v<Aplicatie>.db` și întoarce calea completă. Funcția `pastreazaUltimele` ignoră prefixul `pontix-manual-`, astfel încât copiile manuale nu șterg și nu sunt șterse de copiile automate făcute înaintea migrărilor.

