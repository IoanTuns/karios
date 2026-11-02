-- +goose Up
ALTER TABLE coduri_absenta ADD COLUMN efect_ore TEXT NOT NULL DEFAULT 'ANTERIOR' CHECK(efect_ore IN ('ANTERIOR','FARA_ORE','NORMA','PROGRAM'));
ALTER TABLE absente_zile ADD COLUMN minute_creditate INTEGER CHECK(minute_creditate BETWEEN 0 AND 2880);

-- +goose Down
ALTER TABLE absente_zile DROP COLUMN minute_creditate;
ALTER TABLE coduri_absenta DROP COLUMN efect_ore;
