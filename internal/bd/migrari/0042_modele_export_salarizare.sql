-- F10-T03 (D108): modelele exportului pentru salarizare (CSV sau Excel, coloane configurabile).
-- Doua modele generice la instalare; modelele pentru programe anume se adauga dupa documentatia lor.
-- +goose Up
CREATE TABLE modele_export_salarizare (
    id            TEXT PRIMARY KEY,
    denumire      TEXT NOT NULL CHECK (length(denumire) BETWEEN 1 AND 100),
    configuratie  TEXT NOT NULL CHECK (json_valid(configuratie)),
    activ         INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    revizie       INTEGER NOT NULL DEFAULT 1 CHECK (revizie > 0),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL
) STRICT;

INSERT INTO modele_export_salarizare (id, denumire, configuratie, creat_la, modificat_la) VALUES
    ('export-csv-generic', 'CSV generic', '{"format": "CSV", "separator": ";", "bom": true, "separator_zecimal": ",", "zecimale": 2, "forma_ore": "ZECIMAL", "antet": true, "coloane": [{"sursa": "marca", "titlu": "Marca"}, {"sursa": "nume_complet", "titlu": "Nume"}, {"sursa": "total:OL", "titlu": "Ore lucrate"}, {"sursa": "total:OS", "titlu": "Ore suplimentare"}, {"sursa": "total:ZNL", "titlu": "Ore ZNL"}, {"sursa": "total:CO", "titlu": "Zile CO"}, {"sursa": "total:CM", "titlu": "Zile CM"}]}', '2026-11-12T00:00:00Z', '2026-11-12T00:00:00Z'),
    ('export-xlsx-generic', 'Excel generic', '{"format": "XLSX", "separator": ";", "bom": true, "separator_zecimal": ",", "zecimale": 2, "forma_ore": "ZECIMAL", "antet": true, "coloane": [{"sursa": "marca", "titlu": "Marca"}, {"sursa": "nume_complet", "titlu": "Nume"}, {"sursa": "total:OL", "titlu": "Ore lucrate"}, {"sursa": "total:OS", "titlu": "Ore suplimentare"}, {"sursa": "total:ZNL", "titlu": "Ore ZNL"}, {"sursa": "total:CO", "titlu": "Zile CO"}, {"sursa": "total:CM", "titlu": "Zile CM"}]}', '2026-11-12T00:00:00Z', '2026-11-12T00:00:00Z');

-- +goose Down
DROP TABLE modele_export_salarizare;
