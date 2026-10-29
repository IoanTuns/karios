# D045 · Conturi: întrebări, parolă, roluri și primul administrator

Stare: acceptată · Data: 2026-10-29 · Decisă de: Ioan Tuns (alegerile din conversația cu Claude, 2026-10-29) · Taskuri: F15-T01 … F15-T08

Completează [D008](D008-criptarea-datelor.md) punctul 3, care a decis conturile, auto-login-ul pe același Windows și calculator, rolurile, recuperarea și blocarea progresivă.

## Context

Pentru implementarea F15 lipseau patru detalii: cum se aleg întrebările de recuperare, ce parolă se acceptă, ce nu poate face un operator și cum apare primul cont într-o bază care are deja date.

## Decizie

1. **Întrebările de recuperare** se aleg dintr-o listă fixă de 10, câte 3 diferite pe cont. Răspunsurile se normalizează (litere mici, fără diacritice, spații strânse) și nu se păstrează nicăieri în clar. Administratorul poate și reseta parola unui cont.
2. **Parola:** cel puțin 8 caractere, cu cel puțin o literă și o cifră.
3. **Roluri:** operatorul poate face tot, mai puțin gestiunea conturilor și confirmarea incidentelor de audit, rezervate administratorului.
4. **Primul administrator:** la prima pornire fără conturi, inclusiv după actualizarea unei baze cu date, aplicația cere obligatoriu crearea contului de administrator, înainte de orice altceva. Datele existente rămân.

### Implementare (propunerea lui Claude, în cadrul D008)

- **Secretul nu se dublează în bază:** parola deschide copia `P:<id cont>` din seif, iar răspunsurile copia `R:<id cont>` (Argon2id, cum prevede D008). Verificarea înseamnă deschiderea copiei și compararea cheii obținute cu cheia de date. Trecerea la etapa 2 (I006) rămâne ștergerea copiei A.
- **Baza păstrează** contul (nume, nume afișat, rol, activ, întrebările alese, încercările greșite, blocarea) în tabelul `conturi`, sigilat ca restul (D044).
- **Auto-login:** baza reține ultimul cont, SID-ul utilizatorului Windows și numele calculatorului; dacă se potrivesc toate trei, contul intră fără parolă. „Ieșire din cont” oprește auto-login-ul până la următoarea autentificare cu parolă.
- **Blocare:** de la a 5-a încercare greșită consecutivă, contul se blochează 1 min, apoi 5 min, 30 min, iar de la a 8-a 120 min. Fiecare încercare, reușită sau nu, intră în jurnal.
- **Autorul din jurnal** devine numele contului (`ana` sau `ana/asistent_ai`).
- **Ultimul administrator activ** nu poate fi dezactivat sau retrogradat.
- Până la autentificare, orice scriere în bază e refuzată de aceleași trigger-e ca la incidentul de audit, cu excepția operațiilor de autentificare.

## Variante respinse

- **Întrebări scrise de utilizator** — întrebările slabe se ghicesc ușor (Ioan Tuns).
- **Doar resetare de către administrator** — singurul administrator nu și-ar mai putea recupera contul.
- **Parolă de minim 10 sau 12 caractere** — mai incomod pentru personal (Ioan Tuns).
- **Hash-ul parolei și în bază** — ar fi un al doilea secret de ținut sincronizat cu seiful.

## Consecințe

- La prima pornire după actualizare, nimeni nu lucrează până nu se creează administratorul.
- Dacă seiful pierde copia unui cont, contul nu se mai poate autentifica; administratorul îi resetează parola.
- Cât timp există copia A (etapa 1), cine are `.exe`-ul poate citi baza fără cont; conturile protejează aplicația, nu fișierul.
