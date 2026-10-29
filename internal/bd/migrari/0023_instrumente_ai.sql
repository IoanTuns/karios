-- F14-T25: instantanee de anulare și șabloane locale.
-- +goose Up
CREATE TABLE operatii_ai (
    id TEXT PRIMARY KEY,
    creat_la TEXT NOT NULL,
    plan TEXT NOT NULL,
    instantanee TEXT NOT NULL,
    anulat_la TEXT NOT NULL DEFAULT ''
) STRICT;
CREATE TABLE sabloane_ai (
    id TEXT PRIMARY KEY,
    nume TEXT NOT NULL CHECK(length(nume) BETWEEN 1 AND 100),
    intentie TEXT NOT NULL,
    modificat_la TEXT NOT NULL
) STRICT;
-- +goose Down
DROP TABLE sabloane_ai;
DROP TABLE operatii_ai;
