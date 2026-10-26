# D015 · Nomenclatorul configurabil de schimburi

Stare: propusă · Data: 2026-10-25 · Taskuri: F07-T02

## Cerință

Ioan Tuns acceptă `N12` drept cod scurt și cere un tabel pentru configurarea lui. Folosim tabelul `tipuri_schimb` deja existent din 0001: cod, denumire traductibilă, ora de început, durată, pauză, tip și activ/inactiv.

## Propunere

- Ecran Schimburi cu listare, creare, editare și dezactivare/reactivare, fără ștergere.
- Codul se normalizează la majuscule, fără spații, maximum 16 caractere ASCII (A–Z, cifre, cratimă sau subliniere); verificarea duplicatelor ignoră majusculele. Denumirea rămâne Unicode traductibilă. Acestea sunt reguli tehnice de introducere pentru grilă, nu reguli de plată.
- Ora de început se cere explicit, între 00:00 și 23:59. Durata se introduce explicit în minute (1–2880, ca în schemă); pauza este între 0 și durata schimbului.
- Tipurile existente: NORMAL, TURA, NOAPTE, GARDA. Nu deducem ora sau tipul din cod și nu adăugăm automat N12/7/G24 cu intervale presupuse.
- Denumiri RO/EN în `i18n.Text`, validări în Go, tranzacții și jurnal cu actor `local`.
- Grila inițială va putea sugera codurile active. Interpretarea absențelor și calculul coloanelor din raport rămân în F01/F04; setarea unei pauze în nomenclator nu decide cum se scade la plată.

Fără biblioteci și fără migrare nouă.
