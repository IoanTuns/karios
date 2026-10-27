-- F11-T08, D041: Motor de reguli configurabile (Rule Engine)
-- +goose Up
CREATE TABLE reguli_configurabile (
    id TEXT PRIMARY KEY,
    cod TEXT NOT NULL UNIQUE,
    denumire TEXT NOT NULL,
    descriere TEXT NOT NULL DEFAULT '',
    tip_regula TEXT NOT NULL CHECK(tip_regula IN ('REPAUS_MINIM', 'PLAFON_LUNAR', 'SUCCESIUNE_INTERZISA')),
    severitate TEXT NOT NULL DEFAULT 'AVERTISMENT' CHECK(severitate IN ('INFO', 'AVERTISMENT', 'BLOCANT')),
    coduri_schimb TEXT NOT NULL,
    valoare_numerica INTEGER NOT NULL DEFAULT 0,
    activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN (0,1)),
    creat_la TEXT NOT NULL,
    modificat_la TEXT NOT NULL
) STRICT;

INSERT INTO reguli_configurabile (id, cod, denumire, descriere, tip_regula, severitate, coduri_schimb, valoare_numerica, activ, creat_la, modificat_la)
VALUES
('reg-repaus-g24', 'REPAUS_G24', 'Repaus minim obligatoriu după gardă', 'După o gardă de 24 de ore (G24), angajatul trebuie să beneficieze de minim 24h repaus liber.', 'REPAUS_MINIM', 'AVERTISMENT', 'G24', 24, 1, '2026-10-27T00:00:00Z', '2026-10-27T00:00:00Z'),
('reg-plafon-garzi', 'PLAFON_GARZI', 'Plafon maxim gărzi pe lună', 'Un angajat nu poate fi planificat la mai mult de 6 gărzi de 24h pe parcursul aceleiași luni.', 'PLAFON_LUNAR', 'AVERTISMENT', 'G24', 6, 1, '2026-10-27T00:00:00Z', '2026-10-27T00:00:00Z'),
('reg-succesiune-g24-z7', 'SUCCESIUNE_G24_Z7', 'Interdicție tură de zi imediat după gardă', 'Interzisă planificarea schimbului Z7 în ziua imediat următoare unei gărzi G24.', 'SUCCESIUNE_INTERZISA', 'AVERTISMENT', 'G24,Z7', 0, 1, '2026-10-27T00:00:00Z', '2026-10-27T00:00:00Z');

-- +goose Down
DROP TABLE reguli_configurabile;
