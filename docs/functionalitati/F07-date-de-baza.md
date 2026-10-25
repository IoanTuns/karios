# F07 · Date de bază

Etapă: V1 · Stare: în lucru · Depinde de: F03, F06

## Scop

Utilizatorul poate gestiona angajații și listele de care depinde pontajul.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F07-T01 | Administrarea angajaților: nume, prenume, marcă, tip contract, normă parțială, concediu anual, funcție, secție, activ | gata | Codex | 2026-10-25, 0.10.0: interfață, metode Wails și serviciu Go peste sqlc; listare/căutare/filtre, creare/editare, dezactivare/reactivare, secții și norma globală. Teste Go pe baza criptată, UI în Edge cu punte simulată și build Windows trecute; [ghid](../ghid-primele-functionalitati.md). Urmează F07-T03 și prototipul F08-T01; [D013](../decizii/D013-primele-functionalitati-demo.md) rămâne propusă |
| F07-T02 | Tipuri de schimb și codurile lor scurte (`7`, `12`, `N12`, `G24`) | gata | Codex | 2026-10-25, 0.12.0: ecran de listare/creare/editare/dezactivare, denumiri RO/EN și intervale explicite. Teste Go și UI Edge, fără coduri sau ore precompletate; [D015](../decizii/D015-nomenclator-schimburi.md). Urmează grila F08 |
| F07-T03 | Administrarea sărbătorilor pentru fiecare an | gata | Codex | 2026-10-25, 0.11.0: set național 2026 (17 sărbători în 16 date), migrarea de date 0004, listare pe an, creare/editare RO/EN și completare fără suprascriere. Teste Go și UI Edge; [D014](../decizii/D014-sarbatori-2026.md). Pas următor: nomenclator schimburi și grila F08 |
| F07-T04 | Coduri de absență | de făcut | | Lista finală din F01-T02 |
| F07-T05 | Categorii de ore și coloanele centralizatorului | de făcut | | O singură sursă pentru legătura categorie → coloană |
| F07-T06 | Schema angajaților, norma întreagă globală și sărbătorile pe ani | gata | Codex | 2026-10-25, versiunea 0.9.0, migrările 0002 și 0003; păstrarea datelor vechi, constrângeri, normă efectivă și backup înainte de migrare verificate prin teste. Build Windows verificat. Pas următor: F06-T03, F07-T01 și F07-T03; [D012](../decizii/D012-date-angajati-norme-sarbatori.md) rămâne propusă pentru reprezentarea tehnică |
| F07-T07 | Set de date demo fictive și încărcare explicită, tranzacțională | gata | Codex | 2026-10-25, 0.10.0: 12 angajați și 3 secții, JSON inclus în aplicație; confirmare în UI și refuz în backend dacă există angajați/pontaje. Testate anularea integrală la eroare și păstrarea normei globale. Datele demo nu sunt încărcate în baza operațională de agent |
| F07-T08 | Departamente multiple în fiecare secție, atribuire angajați și filtre | gata | Codex | 2026-10-25, 0.14.0: migrarea 0006, creare în Setări, atribuire filtrată după secție, filtre în Angajați/Pontaj. Verificate integritatea relației și păstrarea configurațiilor vechi. Urmează F07-T09; [D018](../decizii/D018-sectii-departamente.md) rămâne propusă |
| F07-T09 | Setările inițiale ale schimburilor din PDF | de făcut | | Ioan Tuns a precizat sursa: PDF-ul, nu valori inventate. Se extrag exact valorile disponibile; schimburile deja configurate se păstrează |

## Schema F07-T06

| Informație | Reprezentare |
| --- | --- |
| Nume și prenume | `angajati.nume`, `angajati.prenume`; numele vechi rămâne intact, prenumele se completează manual |
| Număr de marcă | `angajati.marca`, text pentru păstrarea zerourilor inițiale, unic între angajații neșterși |
| Tip contract | `angajati.tip_contract`: `INTREAGA` implicit / `PARTIALA` |
| Norma întreagă | `setari.valoare` pentru cheia `norma_intreaga_minute`, implicit `420` minute = 7 ore/zi |
| Norma parțială | `angajati.norma_partiala_minute`, obligatorie numai pentru `PARTIALA`; minute întregi pozitive, cel mult 1440/zi |
| Norma efectivă curentă | `angajati_norme.norma_efectiva_minute`; schimbarea globală se vede imediat pentru contractele întregi |
| Concediu anual contractual | `angajati.concediu_zile_an`, întreg nenegativ sau `NULL` dacă nu a fost completat; nu reprezintă soldul disponibil |
| Sărbători pe an | `sarbatori.data` și `denumire` JSON; `an` derivat automat din dată, index `(an, data)` |

Exemple pentru viitorul cod de acces la date (nu sunt încă expuse prin interfață):

```sql
SELECT id, nume, prenume, marca, tip_contract, norma_efectiva_minute, concediu_zile_an
FROM angajati_norme WHERE sters_la IS NULL ORDER BY nume, prenume;

UPDATE setari SET valoare = '420', modificat_la = '2026-10-25T10:00:00Z'
WHERE cheie = 'norma_intreaga_minute';

SELECT data, denumire FROM sarbatori WHERE an = 2027 ORDER BY data;
```

Compatibilitate: `norma_minute_anterioara` păstrează vechea valoare exclusiv pentru verificare; este `NULL` pentru angajații noi și nu se folosește în noile calcule. Tipul `INTREAGA` al înregistrărilor migrate este valoarea implicită, nu o clasificare contractuală dedusă. Înaintea folosirii operative se verifică manual tipul, prenumele, norma parțială și concediul.

Pasul următor: F06-T03 pentru accesul la date și F07-T01/F07-T03 pentru formulare; istoricul contractual este înregistrat ca I010, înainte de F12. Nu se modifică motorul de calcul sau regulile de plată în acest task.
