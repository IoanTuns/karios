-- F07-T12, D024: cereri de concediu semnate atașate, număr de înregistrare, înlocuitor și stare de aprobare.
-- +goose Up
CREATE TABLE absente_documente (
    id          TEXT PRIMARY KEY,
    nume_fisier TEXT NOT NULL CHECK (length(trim(nume_fisier)) BETWEEN 1 AND 255),
    tip_mime    TEXT NOT NULL CHECK (tip_mime IN ('application/pdf', 'image/png', 'image/jpeg')),
    continut    BLOB NOT NULL CHECK (length(continut) BETWEEN 1 AND 10485760),
    creat_la    TEXT NOT NULL
) STRICT;

ALTER TABLE absente_zile ADD COLUMN nr_inregistrare TEXT NOT NULL DEFAULT ''
    CHECK (length(nr_inregistrare) <= 64);
ALTER TABLE absente_zile ADD COLUMN inlocuitor_id TEXT REFERENCES angajati (id);
ALTER TABLE absente_zile ADD COLUMN stare_aprobare TEXT NOT NULL DEFAULT 'IN_ASTEPTARE'
    CHECK (stare_aprobare IN ('IN_ASTEPTARE', 'APROBATA', 'RESPINSA'));
ALTER TABLE absente_zile ADD COLUMN aprobat_de TEXT NOT NULL DEFAULT ''
    CHECK (length(aprobat_de) <= 120);
ALTER TABLE absente_zile ADD COLUMN document_id TEXT REFERENCES absente_documente (id);

-- +goose StatementBegin
CREATE TRIGGER absente_zile_inlocuitor_inserare BEFORE INSERT ON absente_zile
WHEN NEW.inlocuitor_id IS NOT NULL AND NEW.inlocuitor_id = NEW.angajat_id
BEGIN
    SELECT RAISE(ABORT, 'absente.inlocuitor_invalid');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER absente_zile_inlocuitor_modificare BEFORE UPDATE OF angajat_id, inlocuitor_id ON absente_zile
WHEN NEW.inlocuitor_id IS NOT NULL AND NEW.inlocuitor_id = NEW.angajat_id
BEGIN
    SELECT RAISE(ABORT, 'absente.inlocuitor_invalid');
END;
-- +goose StatementEnd
