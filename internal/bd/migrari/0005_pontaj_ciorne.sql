-- F08-T01, D016: schița vizuală, separată de intervalele de pontaj efectiv.
-- Revizia fiecărei celule este modificat_la; verificarea concurenței se face
-- în tranzacția serviciului înainte de actualizare.
-- +goose Up

CREATE TABLE pontaj_ciorne (
    angajat_id   TEXT NOT NULL REFERENCES angajati (id),
    data         TEXT NOT NULL CHECK (data IS date(data)),
    valoare      TEXT NOT NULL DEFAULT '' CHECK (length(valoare) <= 16
        AND instr(valoare, char(0)) = 0 AND instr(valoare, char(9)) = 0 AND instr(valoare, char(10)) = 0 AND instr(valoare, char(13)) = 0),
    modificat_la TEXT NOT NULL,
    PRIMARY KEY (angajat_id, data)
) STRICT;

CREATE INDEX pontaj_ciorne_data ON pontaj_ciorne (data, angajat_id);
