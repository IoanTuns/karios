-- F03-T11 (D004, D106): textele interfetei corectate sau traduse de utilizator, importate din
-- folderul karios/traduceri. Se pastreaza doar textele diferite de cataloagele incluse in .exe.
-- +goose Up
CREATE TABLE traduceri (
    limba        TEXT NOT NULL CHECK (length(limba) BETWEEN 2 AND 5),
    cheie        TEXT NOT NULL CHECK (length(cheie) BETWEEN 1 AND 200),
    text         TEXT NOT NULL CHECK (length(text) BETWEEN 1 AND 4000),
    modificat_la TEXT NOT NULL,
    PRIMARY KEY (limba, cheie)
) STRICT;

-- +goose Down
DROP TABLE traduceri;
