-- F04-T10 (I131, D105): tabelele de tarif din D003 (categorii_ore, seturi_reguli,
-- reguli_interval) nu sunt folosite: coloanele vin din Configurare pontaj (D092), iar sporurile
-- sunt scoase prin D100. Se sterg, cu tot cu datele (daca exista).
-- +goose Up
DROP INDEX IF EXISTS reguli_interval_set;
DROP TABLE IF EXISTS reguli_interval;
DROP TABLE IF EXISTS seturi_reguli;
DROP TABLE IF EXISTS categorii_ore;

-- +goose Down
CREATE TABLE categorii_ore (
    cod           TEXT PRIMARY KEY,
    denumire      TEXT NOT NULL CHECK (json_valid(denumire)),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL
) STRICT;

-- Regulile de tarif, versionate în timp (D003).
CREATE TABLE seturi_reguli (
    id              TEXT PRIMARY KEY,
    denumire        TEXT NOT NULL,
    valabil_de_la   TEXT NOT NULL CHECK (valabil_de_la IS date(valabil_de_la)),
    valabil_pana_la TEXT CHECK (valabil_pana_la IS NULL OR valabil_pana_la IS date(valabil_pana_la)),
    creat_la        TEXT NOT NULL,
    modificat_la    TEXT NOT NULL
) STRICT;

-- Intervalele acoperă exact 00:00–24:00 pe fiecare tip de zi (verificarea din cod a fost ștearsă odată cu tabelele).
CREATE TABLE reguli_interval (
    id            TEXT PRIMARY KEY,
    set_id        TEXT NOT NULL REFERENCES seturi_reguli (id) ON DELETE CASCADE,
    tip_zi        TEXT NOT NULL CHECK (tip_zi IN ('LUCRATOARE', 'SAMBATA', 'DUMINICA', 'SARBATOARE')),
    ora_start     INTEGER NOT NULL CHECK (ora_start BETWEEN 0 AND 1439),
    ora_sfarsit   INTEGER NOT NULL CHECK (ora_sfarsit BETWEEN 1 AND 1440 AND ora_sfarsit > ora_start),
    categorie_cod TEXT NOT NULL REFERENCES categorii_ore (cod)
) STRICT;

CREATE INDEX reguli_interval_set ON reguli_interval (set_id, tip_zi, ora_start);
