-- F11-T01: planul este separat de schița publicată în pontaj.
-- +goose Up
CREATE TABLE planificare (
    id TEXT PRIMARY KEY,
    angajat_id TEXT NOT NULL REFERENCES angajati(id),
    data TEXT NOT NULL CHECK(data IS date(data)),
    schimb_id TEXT NOT NULL REFERENCES tipuri_schimb(id),
    sectiune TEXT NOT NULL CHECK(sectiune IN ('GARZI','DUPA_AMIAZA')),
    activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN (0,1)),
    revizie INTEGER NOT NULL DEFAULT 1 CHECK(revizie>0),
    publicat_plan_revizie INTEGER NOT NULL DEFAULT 0,
    publicat_pontaj_revizie INTEGER NOT NULL DEFAULT 0,
    publicat_valoare TEXT NOT NULL DEFAULT '',
    publicat_la TEXT NOT NULL DEFAULT '',
    creat_la TEXT NOT NULL,
    modificat_la TEXT NOT NULL
) STRICT;
CREATE UNIQUE INDEX planificare_angajat_data ON planificare(angajat_id,data) WHERE activ=1;
CREATE INDEX planificare_perioada ON planificare(data,sectiune);
