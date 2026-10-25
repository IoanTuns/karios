# D017 · Build fără ștergerea datelor aplicației

Stare: propusă · Data: 2026-10-25 · Taskuri: F05-T15

## Problemă observată

La verificarea F08-T01, scriptul existent a rulat `wails build -clean`. Wails curăță întregul `build/bin`, unde aplicația portabilă poate avea datele în `karios/`. Încercarea a eșuat pe `pontix.db`, blocat de instanța deschisă. Curățarea poate elimina alte fișiere înainte să întâlnească fișierul blocat; nu pretindem că tentativa a fost atomică sau că toate fișierele au rămas intacte.

Verificăm doar metadatele și existența fișierelor, fără afișarea datelor angajaților sau a cheilor. Nu oprim procesul utilizatorului și nu ștergem ori recreăm datele.

Verificarea ulterioară a confirmat lipsa `seif` și a folderului `backup`; baza criptată a rămas pe disc. Instanța fusese închisă de utilizator. Ioan Tuns a confirmat că datele erau exclusiv demo și a autorizat regenerarea. Datele demo au fost regenerate și redeschise cu succes într-un folder separat `build/bin/pontix-0.13.0/karios`: 12 angajați, 3 secții și 16 date de sărbătoare. Baza veche a fost păstrată în locul inițial, fără încercare de înlocuire cu o cheie nouă.

## Propunere

- Eliminăm opțiunea `-clean`; build-ul nu trebuie să curețe folderul portabil.
- `build.ps1` primește `-NumeExecutabil`, implicit `pontix.exe`, limitat la un nume `.exe` fără cale. Exemplu: `./build.ps1 -NumeExecutabil pontix-0.13.0.exe`.
- Dacă ieșirea există, verificăm accesul exclusiv înainte de build; un executabil blocat produce un mesaj clar, fără încercare de ștergere. Pentru o versiune deschisă construim sub nume alternativ.
- Verificarea proprietăților Windows citește exact fișierul de ieșire ales.
- Testul final construiește un executabil alternativ și verifică existența fișierelor de date înainte și după; nu pornește noul executabil peste instanța deschisă.

O curățare manuală a datelor rămâne o acțiune separată, cu aprobarea utilizatorului conform AGENTS.md.
