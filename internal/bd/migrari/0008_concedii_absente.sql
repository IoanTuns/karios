-- F07-T04: evidență explicită pe zile; fără reguli de plată.
-- +goose Up
ALTER TABLE coduri_absenta ADD COLUMN categorie TEXT NOT NULL DEFAULT 'SPECIAL'
    CHECK (categorie IN ('ODIHNA', 'MEDICAL', 'SPECIAL'));

INSERT INTO coduri_absenta (cod, denumire, categorie, creat_la, modificat_la)
SELECT 'CO', '{"ro":"Concediu de odihnă","en":"Annual leave"}', 'ODIHNA',
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'CO')
  AND NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'CO');
INSERT INTO coduri_absenta (cod, denumire, categorie, creat_la, modificat_la)
SELECT 'CM', '{"ro":"Concediu medical","en":"Medical leave"}', 'MEDICAL',
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'CM')
  AND NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'CM');

CREATE TABLE absente_zile (
    angajat_id TEXT NOT NULL REFERENCES angajati (id),
    data TEXT NOT NULL CHECK (data IS date(data)),
    cod TEXT NOT NULL REFERENCES coduri_absenta (cod),
    activ INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    creat_la TEXT NOT NULL,
    modificat_la TEXT NOT NULL,
    PRIMARY KEY (angajat_id, data)
) STRICT;
CREATE INDEX absente_zile_data ON absente_zile (data, angajat_id);
