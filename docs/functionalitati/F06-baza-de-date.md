# F06 · Baza de date și migrările

Etapă: V1 · Stare: de făcut · Depinde de: F05

## Scop

Datele stau într-un singur fișier SQLite, iar un fișier dintr-o versiune mai veche se actualizează automat la pornire.

## Decizii

- [D001](../decizii/D001-stack-tehnologic.md) — SQLite, `modernc.org/sqlite`, `sqlc`
- [D004](../decizii/D004-traduceri.md) — denumirile ca `Text` (JSON)

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F06-T01 | Alege sistemul de migrări: `goose` sau `PRAGMA user_version`, nu ambele | de făcut | | Necesită o decizie |
| F06-T02 | Schema inițială: tabelele din documentul de design, denumiri `Text`, coloane de sincronizare | de făcut | | |
| F06-T03 | Codul de acces la date generat cu `sqlc` | de făcut | | |
| F06-T04 | `jurnal_modificari` de la prima schemă | de făcut | | Mai ieftin acum decât adăugat peste toate scrierile mai târziu |
| F06-T05 | Absențele: în `pontaj` sau într-un tabel separat; zi cu ore lucrate și absență parțială | de făcut | | Depinde de F01-T02 |
