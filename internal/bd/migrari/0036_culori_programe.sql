-- F07-T42 (D097): culoarea fiecarui program (tura, garda), afisata in grilele Planificare
-- si Pontaj si in rapoartele PDF. NULL = culoarea implicita dupa tipul programului.
-- +goose Up
ALTER TABLE tipuri_schimb ADD COLUMN culoare TEXT
    CHECK (culoare IS NULL OR culoare GLOB '#[0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F]');

-- +goose Down
ALTER TABLE tipuri_schimb DROP COLUMN culoare;
