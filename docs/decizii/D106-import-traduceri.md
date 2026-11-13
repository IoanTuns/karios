# D106 · Importul traducerilor din folder în baza criptată

Stare: acceptată · Data: 2026-11-12 · Taskuri: F03-T11 · Cerută de: Ioan Tuns, 2026-11-12 („Mergi pentru F03-T11”) · Decisă de: Ioan Tuns, 2026-11-13 („accept toate”)

## Context

D004 spune că un folder de lângă `.exe` poate corecta textele sau adăuga limbi. Nota ei ulterioară cere ca traducerile corectate să ajungă în baza criptată (D008). Funcția `i18n.AdaugaFolder` exista, dar nu era apelată nicăieri. Ioan Tuns a cerut implementarea fără să aleagă formatul și regulile, așa că le propun aici.

## Decizie

- **Folderul.** Fișierele se pun în `karios/traduceri/`, lângă celelalte date ale aplicației. Aplicația scrie doar în `karios/` și `ai/` (D091).
- **Formatul.** Un fișier pe limbă, `<cod>.json` (de exemplu `en.json`, `hu.json`), plat, cheie → text, la fel ca cataloagele incluse (D004).
- **Fișierul de lucru.** În **Setări › General › Limbă**, butonul **Scrie fișierul de lucru** pune în folder toate textele de acum ale unei limbi. Cele lipsă apar în limba implicită. Traducătorul îl editează. Un fișier existent cu același nume se păstrează ca `.vechi`.
- **Verificarea.** Butonul **Verifică folderul** arată, pe limbă: texte noi, modificate, revenite la textul inclus, neschimbate și neimportate, fiecare cu motivul. Nu scrie nimic.
- **Ce se refuză:**
  - o cheie care nu există în aplicație;
  - un text cu alți parametri `{{…}}` decât originalul;
  - un text mai lung de 4.000 de caractere.

  Un fișier care nu se poate citi oprește verificarea și e numit în mesaj. Textele goale sunt ignorate.
- **Importul.** Îl poate face doar administratorul. Se salvează în tabelul `traduceri` doar textele diferite de cele incluse în `.exe`. Un text adus înapoi la cel inclus șterge corectura.
  - Totul se face într-o singură tranzacție.
  - În jurnal intră un rând pe limbă, cu numărul de texte, nu un rând pe text.
  - Tabelul e sub sigiliul de audit (D044).
- **Efectul.** Corecturile au prioritate față de cataloagele incluse și se aplică imediat, în interfață și în mesajele din Go. Se încarcă la deschiderea bazei. După deblocare, interfața își reîncarcă textele.
- **O limbă nouă** (de exemplu `hu.json`) apare în lista limbilor interfeței, cu codul ei ca nume. Textele netraduse apar în română. Rapoartele rămân în română și engleză.
- **Ștergerea.** Butonul **Șterge traducerile importate** șterge corecturile unei limbi; textele revin la cele incluse.

## Variante respinse

- **Citirea directă a folderului la fiecare pornire,** fără import. Ar fi contrazis D004/D008, pentru că textele n-ar fi ajuns în baza criptată. Un fișier stricat ar fi putut bloca pornirea.
- **Excel pentru traducători.** E un al doilea format de întreținut. JSON-ul plat e chiar formatul cataloagelor și se deschide în orice editor.
- **Importul de către orice utilizator.** Textele se schimbă pentru toți, deci decide administratorul.

## Consecințe

- Migrarea 0039 adaugă tabelul `traduceri`.
- Selectorul de limbă din bara laterală arată în continuare doar română și engleză. O limbă nouă se alege din fila Limbă.
