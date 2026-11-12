-- F04-T05 (D107): cum se scade pauza unui program. NESTABILIT pastreaza comportamentul vechi
-- (programul cu pauza ramane necalculat) pana cand utilizatorul alege modul; pauza_ora este ora
-- de inceput a pauzei (minute de la miezul noptii) pentru modul LA_ORA.
-- +goose Up
ALTER TABLE tipuri_schimb ADD COLUMN mod_pauza TEXT NOT NULL DEFAULT 'NESTABILIT'
    CHECK (mod_pauza IN ('NESTABILIT', 'NU_SE_SCADE', 'DIN_TOTAL', 'LA_ORA', 'LA_MIJLOC'));
ALTER TABLE tipuri_schimb ADD COLUMN pauza_ora INTEGER CHECK (pauza_ora IS NULL OR pauza_ora BETWEEN 0 AND 1439);

-- +goose Down
ALTER TABLE tipuri_schimb DROP COLUMN pauza_ora;
ALTER TABLE tipuri_schimb DROP COLUMN mod_pauza;
