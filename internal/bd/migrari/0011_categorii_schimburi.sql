-- F07-T11, D023: categorii în tabelul de schimburi (LAB și CUSTOM).
-- +goose Up
ALTER TABLE tipuri_schimb ADD COLUMN categorie TEXT NOT NULL DEFAULT 'CUSTOM'
    CHECK (categorie IN ('LAB', 'CUSTOM'));

-- Codurile Z7, S7, N9, IN9, Z12, N12, G24 intră în categoria de precompletare LAB;
-- restul (inclusiv 7, 12, S1, S3 și schimburile definite manual) rămân în CUSTOM.
-- Intervalele, starea activ, regimul și momentele de audit se păstrează.
UPDATE tipuri_schimb SET categorie = 'LAB'
WHERE sters_la IS NULL AND upper(cod) IN ('Z7', 'S7', 'N9', 'IN9', 'Z12', 'N12', 'G24');
