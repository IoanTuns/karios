-- F07-T27: soldul initial este separat de dreptul anual de concediu.
-- +goose Up
CREATE TABLE initializari_pontaj (
 angajat_id TEXT PRIMARY KEY REFERENCES angajati(id),
 luna TEXT NOT NULL CHECK(length(luna)=7),
 minute INTEGER NOT NULL,
 co_zile INTEGER NOT NULL CHECK(co_zile>=0),
 revizie INTEGER NOT NULL DEFAULT 1,
 modificat_la TEXT NOT NULL
) STRICT;
-- +goose Down
DROP TABLE initializari_pontaj;
