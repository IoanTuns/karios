# D047 · Contextul asistentului AI

Stare: propusă · Data: 2026-10-29 · Propusă de: Codex

Task: F14-T20. Eroarea raportată arată 9658 tokenuri pentru un context de 8192.

Propunere implementată pentru verificare: context local de 32768 tokenuri, în loc de 8192. La HTTP 400 cu tipul exact `exceed_context_size_error`, clientul elimină cel mai vechi schimb conversațional până la următorul mesaj al utilizatorului și reîncearcă. Instrucțiunile, datele aplicației și ultima cerere rămân întregi. Alte erori nu se reîncearcă. Dacă nu mai există schimburi vechi, eroarea ajunge la utilizator.

Alternativa de a mări doar contextul amână aceeași eroare pentru conversațiile lungi. Nu estimăm tokenurile după caractere și nu trunchiem datele pontajului. Contextul mai mare consumă memorie suplimentară; funcționarea cu modelul real trebuie verificată. Eliminarea istoricului poate pierde referințe la schimburi vechi.

Nu se adaugă biblioteci. Acceptarea deciziei rămâne la om.
