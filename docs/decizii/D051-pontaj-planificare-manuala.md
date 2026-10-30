# D051 · Programele speciale introduse manual în Pontaj

Stare: propusă. Propusă de: Codex, 2026-10-29. Task F07-T21, cerere Ioan Tuns privind modificarea normal → S7 care nu apare în Planificare.

Salvarea din ecranul Pontaj actualizează în aceeași tranzacție poziția de planificare pentru programele speciale, după clasificarea existentă. Pentru schimbarea între secțiuni se dezactivează poziția veche și se creează alta. Revenirea la normal/gol dezactivează numai o poziție care corespunde programului anterior din celulă. O planificare independentă divergentă nu este suprascrisă: o nouă tură specială incompatibilă este refuzată.

Validarea folosește lotul final, revizii, starea lunii, absențe, categorii și suprapuneri, inclusiv lunile vecine. Jurnalul păstrează autorul real. Nu se modifică pontajele efective sau regulile de calcul. Metoda internă de salvare a schiței fără sincronizare rămâne pentru fluxurile care gestionează separat planificarea, inclusiv AI; interfața manuală folosește explicit varianta sincronizată.

Pentru date deja salvate, utilizatorul poate actualiza explicit planificarea celulei speciale selectate, fără schimbarea orelor ori rescrierea întregii luni. Se refolosesc aceleași validări și tranzacție. Nu se repară automat baza operațională la citire și nu se introduce migrare de date.
