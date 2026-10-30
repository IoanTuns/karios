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
| F15-T01 | Modelul de conturi: nume, rol (administrator / operator), activ, hash parolă, întrebări și hash-urile răspunsurilor | gata | Claude | 0.68.0: tabelul `conturi`, secretele în seif (P/R), fără hash în bază (D045) |
| F15-T02 | Prima pornire: crearea contului administrator și a întrebărilor lui | gata | Claude | 0.68.0: ecran obligatoriu când nu există conturi; scrierea seifului în tranzacție |
| F15-T03 | Pornirea: auto-login pentru ultimul cont pe același Windows și calculator; altfel login cu parolă | gata | Claude | 0.68.0: `sesiune_locala` cu SID și calculator; neverificat pe Windows real |
| F15-T04 | „Ieșire din cont”: oprește auto-login-ul până la următorul login cu parolă | gata | Claude | 0.68.0: „Ieșire din cont” în bara laterală, dezactivată cu modificări nesalvate în pontaj |
| F15-T05 | Parolă uitată: întrebările de recuperare sau resetarea de către administrator | gata | Claude | 0.68.0: 3 întrebări din listă, răspunsuri normalizate; resetare de administrator |
| F15-T06 | Limitarea încercărilor greșite, cu așteptare progresivă și jurnal | gata | Claude | 0.68.0: 1/5/30/120 min de la a 5-a încercare greșită; jurnal pentru fiecare încercare |
| F15-T07 | Ecranul de administrare a conturilor | gata | Claude | 0.68.0: Setări → Conturi (doar administrator); ultimul administrator protejat |
| F15-T08 | Contul curent ca „Cine” în jurnalul de modificări | gata | Claude | 0.68.0: `contextLucru` poartă numele contului; folosit de F12 |
| F15-T09 | Ecranul de creare a contului: întrebările în lățimea cardului, parolele citite corect și la completarea automată, buton „Arată” | gata | Claude | 0.81.1: `FormData` la trimitere, `onInput`, „Arată”, CSS; neverificat vizual |
| F15-T10 | Seiful fără cheie în cod: copia W (DPAPI pe utilizatorul Windows) și copia K (cheia de recuperare) | în lucru | Claude | [D054](../decizii/D054-chei-fara-cod-recuperare.md) |
| F15-T11 | Trecerea datelor existente: cheie de date nouă, recriptarea bazei și a backup-urilor, P0/R0 refăcute la autentificare, ștergerea copiei A | în lucru | Claude | D054 punctul 4 |
| F15-T12 | Ecranul „Deblochează datele” (parolă sau cheie de recuperare) când copia W lipsește | în lucru | Claude | D054 punctul 2 |
| F15-T13 | Cheia de recuperare obligatorie: afișată administratorului, confirmată prin retastare | în lucru | Claude | D054 punctul 3 |
