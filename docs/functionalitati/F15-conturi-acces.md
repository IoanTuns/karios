# F15 · Conturi și acces

Etapă: V1 · Stare: în lucru · Depinde de: F05, F06

## Scop

Fiecare persoană lucrează în contul ei, cu parolă și întrebări de recuperare. Pe același calculator și același utilizator Windows, ultimul cont intră automat. Administratorul gestionează conturile.

## Decizii

- [D008](../decizii/D008-criptarea-datelor.md) — punctul 3: conturi, auto-login, roluri, recuperare
- [D045](../decizii/D045-conturi-detalii.md) — întrebări din listă, parola, rolurile, primul administrator

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F15-T01 | Modelul de conturi: nume, rol (administrator / operator), activ, hash parolă, întrebări și hash-urile răspunsurilor | în lucru | Claude | Argon2id; răspunsuri normalizate |
| F15-T02 | Prima pornire: crearea contului administrator și a întrebărilor lui | în lucru | Claude | Creează și seiful (F05-T10) |
| F15-T03 | Pornirea: auto-login pentru ultimul cont pe același Windows și calculator; altfel login cu parolă | în lucru | Claude |  |
| F15-T04 | „Ieșire din cont”: oprește auto-login-ul până la următorul login cu parolă | în lucru | Claude |  |
| F15-T05 | Parolă uitată: întrebările de recuperare sau resetarea de către administrator | în lucru | Claude | Refacerea copiilor P și R din seif |
| F15-T06 | Limitarea încercărilor greșite, cu așteptare progresivă și jurnal | în lucru | Claude |  |
| F15-T07 | Ecranul de administrare a conturilor | în lucru | Claude | Ultimul administrator nu poate fi șters |
| F15-T08 | Contul curent ca „Cine” în jurnalul de modificări | în lucru | Claude | Folosit de F12 |
