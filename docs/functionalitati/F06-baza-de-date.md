# F06 · Baza de date și migrările

Etapă: V1 · Stare: în lucru · Depinde de: F05

## Scop

Datele stau într-un singur fișier SQLite, iar un fișier dintr-o versiune mai veche se actualizează automat la pornire.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — SQLite, `modernc.org/sqlite`, `sqlc`
- [D004](../decizii/D004-traduceri.md) — denumirile ca `Text` (JSON)

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F06-T01 | Alege sistemul de migrări: `goose` sau `PRAGMA user_version`, nu ambele | gata | Claude | 2026-10-24: propus goose v3.29.0, doar înainte, cu `istoric_migrari` și protecție la downgrade; verificat în prototip ([D010](../decizii/D010-migrari-goose.md)) |
| F06-T02 | Schema inițială: tabelele din documentul de design, denumiri `Text`, coloane de sincronizare; deschiderea bazei criptate în modul WAL | gata | Claude | 2026-10-25, versiunea 0.5.0: `internal/bd` (`Deschide`, `Migreaza`) și migrarea `0001_schema_initiala.sql` cu 11 tabele. Cheia prin `PRAGMA hexkey` la fiecare conexiune; WAL, cheile străine active. Încă nefolosit de `main.go`: are nevoie de seif (F05-T10) |
| F06-T03 | Codul de acces la date generat cu `sqlc` | gata | Codex | 2026-10-25, 0.9.1: `sqlc.yaml`, `genereaza-date.ps1`, interogări și cod generat în `internal/bd/acces`. Prima tranșă: angajați, secții, setări și jurnal; generatorul citește direct migrările 0001–0003. Pas următor: consumarea în F07-T01; interogările pentru celelalte fluxuri se adaugă odată cu ele; [D013](../decizii/D013-primele-functionalitati-demo.md) |
| F06-T04 | `jurnal_modificari` de la prima schemă | în lucru | Claude | 2026-10-28: acoperire completă prin [D044](../decizii/D044-jurnal-audit-sigilat.md): închiderea lunii și regulile configurabile scrise în tranzacție și în jurnal; lanț HMAC pe fiecare rând. 0.60.0: gata pentru fluxurile existente; identificarea contului rămâne la F15. Pas următor: F06-T09 |
| F06-T05 | Absențele: în `pontaj` sau într-un tabel separat; zi cu ore lucrate și absență parțială | de făcut | | Depinde de F01-T02 |
| F06-T06 | Istoricul migrărilor (dată, versiunea aplicației) și refuzul unei baze migrate de o versiune mai nouă | de făcut | | D009 punctul 3 |
| F06-T07 | Revizii unice independente de rezoluția ceasului pentru schițe și absențe | gata | Codex | 2026-10-25, 0.17.1: migrarea 0010, contor per rând incrementat în SQL și DTO text. Conținutul/momentele existente păstrate. Test determinist cu timestamp fix și migrare, plus cele două cazuri de conflict repetate de 10 ori, trecute. Go/vet, UI Edge și build Windows trecute. [D022](../decizii/D022-revizii-concurenta.md) rămâne propusă. Pas următor: F06-T06 și conturi F15 |
| F06-T08 | Sigiliul tabelelor, verificarea la pornire, blocarea editării și confirmarea incidentului | în lucru | Claude | [D044](../decizii/D044-jurnal-audit-sigilat.md). 0.60.0: `internal/audit`, migrarea 0021, `stocare` verifică înainte de migrare și redeschide cu trigger-e, banner și dialog de confirmare. Risc asumat până la F15: oricine poate confirma un incident. Rămas: măsurarea verificării la pornire pe o bază mare și proba în aplicația pornită (neverificate) |
| F06-T09 | Detectarea înlocuirii bazei cu o copie mai veche, validă: contor autentificat în seif | de făcut | | [D044](../decizii/D044-jurnal-audit-sigilat.md) punctul 7 |
