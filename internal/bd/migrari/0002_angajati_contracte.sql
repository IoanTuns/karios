-- F07-T06, D012: datele contractuale și norma întreagă globală.
-- Numele vechi nu se desparte automat; norma veche se păstrează pentru verificare.
-- +goose Up

ALTER TABLE angajati ADD COLUMN norma_minute_anterioara INTEGER;
UPDATE angajati SET norma_minute_anterioara = norma_minute;
ALTER TABLE angajati DROP COLUMN norma_minute;
ALTER TABLE angajati ADD COLUMN prenume TEXT NOT NULL DEFAULT '';
ALTER TABLE angajati ADD COLUMN tip_contract TEXT NOT NULL DEFAULT 'INTREAGA'
    CHECK (tip_contract IN ('INTREAGA', 'PARTIALA'));
ALTER TABLE angajati ADD COLUMN norma_partiala_minute INTEGER
    CHECK ((tip_contract = 'INTREAGA' AND norma_partiala_minute IS NULL)
        OR (tip_contract = 'PARTIALA' AND norma_partiala_minute IS NOT NULL
            AND norma_partiala_minute BETWEEN 1 AND 1440));
ALTER TABLE angajati ADD COLUMN concediu_zile_an INTEGER CHECK (concediu_zile_an >= 0);

-- Valoarea se păstrează canonic în minute întregi (1–1440). Restul setărilor
-- rămân text liber. Cheia normei nu poate fi eliminată sau redenumită.
-- +goose StatementBegin
CREATE TRIGGER setari_norma_inserare BEFORE INSERT ON setari
WHEN NEW.cheie = 'norma_intreaga_minute'
    AND (CAST(CAST(NEW.valoare AS INTEGER) AS TEXT) <> NEW.valoare
        OR CAST(NEW.valoare AS INTEGER) NOT BETWEEN 1 AND 1440)
BEGIN
    SELECT RAISE(ABORT, 'bd.norma_intreaga_invalida');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER setari_norma_modificare BEFORE UPDATE ON setari
WHEN (OLD.cheie = 'norma_intreaga_minute' AND NEW.cheie <> OLD.cheie)
    OR (NEW.cheie = 'norma_intreaga_minute'
        AND (CAST(CAST(NEW.valoare AS INTEGER) AS TEXT) <> NEW.valoare
            OR CAST(NEW.valoare AS INTEGER) NOT BETWEEN 1 AND 1440))
BEGIN
    SELECT RAISE(ABORT, 'bd.norma_intreaga_invalida');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER setari_norma_stergere BEFORE DELETE ON setari
WHEN OLD.cheie = 'norma_intreaga_minute'
BEGIN
    SELECT RAISE(ABORT, 'bd.norma_intreaga_invalida');
END;
-- +goose StatementEnd

INSERT INTO setari (cheie, valoare, modificat_la)
VALUES ('norma_intreaga_minute', '420', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
ON CONFLICT (cheie) DO NOTHING;

-- Verifică și o eventuală valoare preexistentă, fără a o suprascrie.
UPDATE setari SET valoare = valoare WHERE cheie = 'norma_intreaga_minute';

-- Nu copiem norma globală în angajați: contractele întregi urmează setarea.
-- norma_minute_anterioara este doar informație de migrare, nu normă curentă.
CREATE VIEW angajati_norme AS
SELECT a.*,
    CASE a.tip_contract
        WHEN 'INTREAGA' THEN CAST((SELECT valoare FROM setari
            WHERE cheie = 'norma_intreaga_minute') AS INTEGER)
        WHEN 'PARTIALA' THEN a.norma_partiala_minute
    END AS norma_efectiva_minute
FROM angajati AS a;
