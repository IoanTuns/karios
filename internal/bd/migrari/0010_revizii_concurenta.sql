-- F06-T07, D022: revizia nu depinde de rezoluția ceasului.
-- +goose Up
ALTER TABLE pontaj_ciorne ADD COLUMN revizie INTEGER NOT NULL DEFAULT 1 CHECK (revizie > 0);
ALTER TABLE absente_zile ADD COLUMN revizie INTEGER NOT NULL DEFAULT 1 CHECK (revizie > 0);
