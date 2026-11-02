# D083 - Explicatii optionale pe modelul de raport

Stare: propusa. Data: 2026-11-02. Responsabil: Codex. F09-T15.

Cerere explicita: modul de calcul si semnificatia coloanelor se afiseaza numai la alegerea utilizatorului. Propunere implementata: bifa in Continut, comuna tuturor surselor, cu continut.ascunde_explicatii optional in JSON. Un model nou porneste fara explicatii; modelele deja salvate pastreaza comportamentul existent pana la debifare, evitand schimbarea tacita a livrarilor lor. Duplicarea pastreaza alegerea. Fara migrare, fara schimbarea calculului sau a datelor. PDF si Excel omit legenda/formulele/SQL cand sunt ascunse; istoricul tehnic ramane separat si nu este ascuns de aceasta optiune. Versiunile rapoartelor cresc pentru noua capacitate. Referintele tehnice noi pastreaza fisierele vechi.
