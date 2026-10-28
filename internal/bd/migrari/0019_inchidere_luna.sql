-- F12-T01: Închiderea și blocarea lunii la editare
-- +goose Up
CREATE TABLE stare_luni (
    an INTEGER NOT NULL,
    luna INTEGER NOT NULL,
    stare TEXT NOT NULL DEFAULT 'DESCHISA' CHECK(stare IN ('DESCHISA', 'BLOCATA', 'INCHISA')),
    modificat_la TEXT NOT NULL,
    modificat_de TEXT NOT NULL DEFAULT 'Utilizator',
    motiv TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (an, luna)
) STRICT;

-- +goose Down
DROP TABLE stare_luni;
