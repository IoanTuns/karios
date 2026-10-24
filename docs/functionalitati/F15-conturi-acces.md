# F15 · Conturi și acces

Etapă: V1 · Stare: de făcut · Depinde de: F05, F06

## Scop

Fiecare persoană lucrează în contul ei, cu parolă și întrebări de recuperare. Pe același calculator și același utilizator Windows, ultimul cont intră automat. Administratorul gestionează conturile.

## Decizii

- [D008](../decizii/D008-criptarea-datelor.md) — punctul 3: conturi, auto-login, roluri, recuperare

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F15-T01 | Modelul de conturi: nume, rol (administrator / operator), activ, hash parolă, întrebări și hash-urile răspunsurilor | de făcut | | Argon2id; răspunsuri normalizate |
| F15-T02 | Prima pornire: crearea contului administrator și a întrebărilor lui | de făcut | | Creează și seiful (F05-T10) |
| F15-T03 | Pornirea: auto-login pentru ultimul cont pe același Windows și calculator; altfel login cu parolă | de făcut | | |
| F15-T04 | „Ieșire din cont”: oprește auto-login-ul până la următorul login cu parolă | de făcut | | |
| F15-T05 | Parolă uitată: întrebările de recuperare sau resetarea de către administrator | de făcut | | Refacerea copiilor P și R din seif |
| F15-T06 | Limitarea încercărilor greșite, cu așteptare progresivă și jurnal | de făcut | | |
| F15-T07 | Ecranul de administrare a conturilor | de făcut | | Ultimul administrator nu poate fi șters |
| F15-T08 | Contul curent ca „Cine” în jurnalul de modificări | de făcut | | Folosit de F12 |
