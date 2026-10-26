-- F07-T08, D018: fiecare departament aparține unei secții.
-- Angajații existenți rămân în secțiile lor, fără departament atribuit automat.
-- +goose Up

CREATE TABLE departamente (
    id            TEXT PRIMARY KEY,
    sectie_id     TEXT NOT NULL REFERENCES sectii (id),
    denumire      TEXT NOT NULL CHECK (json_valid(denumire)),
    activ         INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL,
    sters_la      TEXT,
    proprietar_id TEXT,
    id_remote     TEXT,
    stare_sinc    TEXT NOT NULL DEFAULT 'local'
) STRICT;

CREATE INDEX departamente_sectie ON departamente (sectie_id) WHERE sters_la IS NULL;
ALTER TABLE angajati ADD COLUMN departament_id TEXT REFERENCES departamente (id);
CREATE INDEX angajati_departament ON angajati (departament_id) WHERE sters_la IS NULL;

-- +goose StatementBegin
CREATE TRIGGER angajati_departament_inserare BEFORE INSERT ON angajati
WHEN NEW.departament_id IS NOT NULL AND NOT EXISTS (
    SELECT 1 FROM departamente WHERE id = NEW.departament_id AND sectie_id = NEW.sectie_id
)
BEGIN
    SELECT RAISE(ABORT, 'date.departament_invalid');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER angajati_departament_modificare BEFORE UPDATE OF departament_id, sectie_id ON angajati
WHEN NEW.departament_id IS NOT NULL AND NOT EXISTS (
    SELECT 1 FROM departamente WHERE id = NEW.departament_id AND sectie_id = NEW.sectie_id
)
BEGIN
    SELECT RAISE(ABORT, 'date.departament_invalid');
END;
-- +goose StatementEnd

-- Nu mutăm un departament sub altă secție lăsând angajații cu o relație invalidă.
-- +goose StatementBegin
CREATE TRIGGER departamente_sectie_modificare BEFORE UPDATE OF sectie_id ON departamente
WHEN NEW.sectie_id <> OLD.sectie_id AND EXISTS (
    SELECT 1 FROM angajati WHERE departament_id = OLD.id AND sectie_id IS NOT NEW.sectie_id
)
BEGIN
    SELECT RAISE(ABORT, 'date.departament_invalid');
END;
-- +goose StatementEnd
