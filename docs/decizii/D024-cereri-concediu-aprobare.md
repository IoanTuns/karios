# D024 · Cereri de concediu semnate atașate, număr de înregistrare, înlocuitor și aprobare

Stare: propusă · Data: 2026-10-25 · Task: F07-T12

Ioan Tuns a cerut ca înregistrările de concedii să poată avea atașată cererea de concediu semnată, un număr de înregistrare, persoana care înlocuiește angajatul aflat în concediu și starea de aprobare. La clarificare a ales:
- păstrarea fișierelor atașate (PDF/JPG/PNG) direct în baza de date SQLite criptată (`BLOB`), incluse automat în backup și protejate cu aceeași cheie ([D008](D008-criptarea-datelor.md));
- câmpuri opționale la salvare, care pot fi adăugate la înregistrare sau completate ulterior pe zi/interval;
- stare de aprobare (`IN_ASTEPTARE`, `APROBATA`, `RESPINSA`) și câmp opțional `aprobat_de`, cu afișarea zilelor active în grila de pontaj și evidențierea stării lor.

## Reprezentare propusă

- Migrarea `0012_cereri_concediu.sql`:
  - creează tabelul `absente_documente (id, nume_fisier, tip_mime, continut, creat_la) STRICT` pentru stocarea fișierelor semnate (`application/pdf`, `image/png`, `image/jpeg`, maximum 10 MB);
  - extinde `absente_zile` cu:
    - `nr_inregistrare TEXT NOT NULL DEFAULT ''`,
    - `inlocuitor_id TEXT REFERENCES angajati (id) CHECK (inlocuitor_id IS NULL OR inlocuitor_id <> angajat_id)`,
    - `stare_aprobare TEXT NOT NULL DEFAULT 'IN_ASTEPTARE' CHECK (stare_aprobare IN ('IN_ASTEPTARE', 'APROBATA', 'RESPINSA'))`,
    - `aprobat_de TEXT NOT NULL DEFAULT ''`,
    - `document_id TEXT REFERENCES absente_documente (id)`.
- Când se înregistrează sau se editează un lot de zile (interval sau selecție multiplă) cu un document atașat nou, fișierul este salvat o singură dată în `absente_documente` în aceeași tranzacție, iar toate zilele lotului primesc același `document_id`.
- Citirea listei de absențe și a schiței de pontaj încarcă doar metadatele documentului (`document_id`, `document_nume`, `document_mime`), fără a citi `BLOB`-ul la fiecare afișare; conținutul documentului se citește la cerere prin `CitesteDocumentAbsenta(id)` pentru vizualizare/descărcare.
- Înlocuitorul (`inlocuitor_id`) este opțional, se alege dintre ceilalți angajați și nu poate fi aceeași persoană cu angajatul în concediu.
- Înregistrările existente din `absente_zile` își păstrează datele și primesc implicit `stare_aprobare = 'IN_ASTEPTARE'`, fără număr de înregistrare, înlocuitor sau document atașat până la completarea lor.
- În grila de pontaj (`Pontaj.tsx`), celulele de absență activă rămân protejate și afișează în titlu/stil starea aprobării, numărul de înregistrare și înlocuitorul (dacă sunt completate).

Reprezentarea tehnică rămâne propusă; cerințele funcționale sunt explicite de la Ioan Tuns.
