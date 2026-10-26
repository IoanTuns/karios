# D022 · Revizii independente de ceas

Stare: propusă · Data: 2026-10-25 · Task: F06-T07

La verificarea F07-T10, două teste existente au arătat că două salvări apropiate pot primi aceeași valoare `modificat_la`. Aceasta era folosită și drept revizie; un client vechi putea astfel suprascrie modificarea intermediară. Repetarea suitei a trecut, confirmând caracterul intermitent, nu corectitudinea mecanismului.

## Corecție propusă

- Migrarea 0010 adaugă `revizie INTEGER NOT NULL DEFAULT 1 CHECK (revizie > 0)` în `pontaj_ciorne` și `absente_zile`, fără a schimba conținutul sau momentele salvate.
- Inserarea începe cu 1. Orice modificare efectivă incrementează revizia în aceeași comandă SQL și o întoarce prin RETURNING. Citirea și validarea sunt în tranzacție; loturile păstrează anularea integrală la conflict.
- DTO-ul continuă să expună revizia ca text, pentru comparații exacte în JavaScript și compatibilitate de structură. Se verifică revizia rândului, nu ordinea momentelor de timp. Anularea unei absențe păstrează rândul și contorul; reînregistrarea continuă contorul.
- `modificat_la` rămâne timestamp de audit, independent de revizie. Un client încărcat înainte de migrare trebuie să reîncarce datele; tokenul său vechi este refuzat.
- Nu introducem biblioteci și nu modificăm testele confirmate de contabilitate. Testul nou fixează timestampul prin trigger în baza temporară, pentru a verifica determinist două modificări cu același moment afișat.

Corecția este separată de programele de lucru și se livrează împreună cu ele în 0.17.1.
