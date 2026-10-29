-- F15, D008 punctul 3, D045: conturile aplicației și sesiunea locală pentru auto-login.
-- Parolele și răspunsurile nu stau aici: deschid copiile P:<id> și R:<id> din seif.
-- +goose Up
CREATE TABLE conturi (
    id               TEXT PRIMARY KEY,
    nume             TEXT NOT NULL UNIQUE COLLATE NOCASE CHECK (length(trim(nume)) BETWEEN 2 AND 40),
    nume_afisat      TEXT NOT NULL DEFAULT '',
    rol              TEXT NOT NULL CHECK (rol IN ('ADMINISTRATOR', 'OPERATOR')),
    activ            INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    intrebari        TEXT NOT NULL DEFAULT '[]' CHECK (json_valid(intrebari)),
    incercari_esuate INTEGER NOT NULL DEFAULT 0 CHECK (incercari_esuate >= 0),
    blocat_pana_la   TEXT NOT NULL DEFAULT '',
    creat_la         TEXT NOT NULL,
    modificat_la     TEXT NOT NULL
) STRICT;

CREATE TABLE sesiune_locala (
    id           INTEGER PRIMARY KEY CHECK (id = 1),
    cont_id      TEXT REFERENCES conturi (id),
    sid          TEXT NOT NULL DEFAULT '',
    calculator   TEXT NOT NULL DEFAULT '',
    autologin    INTEGER NOT NULL DEFAULT 0 CHECK (autologin IN (0, 1)),
    modificat_la TEXT NOT NULL
) STRICT;

-- +goose Down
DROP TABLE sesiune_locala;
DROP TABLE conturi;
