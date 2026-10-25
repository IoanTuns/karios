# F07 · Date de bază

Etapă: V1 · Stare: în lucru · Depinde de: F03, F06

## Scop

Utilizatorul poate gestiona angajații și listele de care depinde pontajul.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F07-T01 | Administrarea angajaților: nume, prenume, marcă, tip contract, normă parțială, concediu anual, funcție, secție, activ | de făcut | | Schema: F07-T06; interfața și accesul la date rămân de implementat |
| F07-T02 | Tipuri de schimb și codurile lor scurte (`7`, `12`, `N12`, `G24`) | de făcut | | De stabilit orele exacte pentru fiecare cod |
| F07-T03 | Administrarea sărbătorilor pentru fiecare an | de făcut | | Denumiri `Text`; schema și filtrarea pe an: F07-T06. Precompletarea nu înlocuiește editarea manuală |
| F07-T04 | Coduri de absență | de făcut | | Lista finală din F01-T02 |
| F07-T05 | Categorii de ore și coloanele centralizatorului | de făcut | | O singură sursă pentru legătura categorie → coloană |
| F07-T06 | Schema angajaților, norma întreagă globală și sărbătorile pe ani | în lucru | Codex | Alocat de Ioan Tuns în conversația din 2026-10-25. Ramură `F07-T06-schema-angajati`; migrare nouă, teste de păstrare a datelor; [D012](../decizii/D012-date-angajati-norme-sarbatori.md) |
