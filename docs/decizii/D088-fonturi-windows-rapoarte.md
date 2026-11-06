# D088 · Fonturile instalate în Windows în rapoartele pdfme

Stare: acceptată · Data: 2026-11-05 · Taskuri: F09-T24 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Regulă la acceptare: pentru rapoartele oficiale se folosesc fonturile Windows standard (Arial, Calibri, Times New Roman) sau fontul inclus în aplicație, pentru că unele fonturi comerciale nu permit încorporarea în PDF, iar aplicația nu verifică licența fiecărui font.

Ioan Tuns cere fonturile din Windows în editorul de rapoarte (I114), după D087.

## Context

pdfme are nevoie de datele complete ale fontului, atât în editor, cât și la generare. Nu poate citi fișiere locale: acceptă adrese web doar de pe gazde publice. Un calculator Windows are sute de fonturi, adică sute de MB. Unele fonturi interzic includerea în documente, prin câmpul `fsType` din tabelul `OS/2`. Altele sunt colecții (`.ttc`), pe care pdfme nu le poate deschide direct.

## Decizie

1. Go citește doar antetele fonturilor din `%WINDIR%\Fonts` și `%LOCALAPPDATA%\Microsoft\Windows\Fonts` (`.ttf`, `.otf`, `.ttc`, `.otc`): numele complet din tabelul `name` și `fsType`. Lista se păstrează în memorie până la închiderea aplicației.
2. Se exclud fonturile cu licență restrictivă (`fsType` bitul 1) și cele doar bitmap (bitul 9). Fonturile cu „fără subset” (bitul 8) se includ întregi în PDF.
3. Interfața primește doar numele. Un font se încarcă numai când este ales în editor sau când apare într-un șablon. Cererea se face după nume, iar Go caută numele în lista proprie; interfața nu trimite căi de fișiere.
4. Dintr-o colecție `.ttc`, Go extrage fontul cerut ca fișier TrueType separat (copiază tabelele lui), ca pdfme să îl poată folosi.
5. Fonturile Go rămân rezerva. Dacă șablonul folosește un font care lipsește pe calculatorul curent, raportul se generează cu fontul Go și utilizatorul vede un avertisment cu fonturile lipsă.
6. Fonturile Windows nu intră în repo și nici în pachetul livrat. Ajung doar în PDF-urile generate, ca subset, cum permite licența fiecăruia.
7. Nu se adaugă bibliotecă nouă: tabelele `name`, `OS/2` și antetul colecției se citesc direct, după specificația OpenType.

## Variante respinse

- **Toate fonturile încărcate la pornirea editorului:** sute de MB în memorie și o pornire lentă.
- **Server HTTP local pentru fonturi:** pdfme refuză gazdele locale, iar un port deschis nu are rost aici.
- **Copierea fonturilor în aplicație:** licențele Microsoft nu permit redistribuirea.

## Consecințe

Un șablon cu fonturi Windows arată la fel pe orice calculator Windows care are fontul, iar PDF-ul generat arată la fel peste tot, pentru că fontul e inclus. Pe un calculator fără fontul respectiv, aspectul diferă, dar utilizatorul este avertizat.
