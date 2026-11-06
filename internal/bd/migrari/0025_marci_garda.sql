-- +goose Up
CREATE TABLE angajati_garzi (
 angajat_id TEXT PRIMARY KEY REFERENCES angajati(id),
 principal_id TEXT NOT NULL REFERENCES angajati(id),
 CHECK(angajat_id <> principal_id)
);
CREATE INDEX angajati_garzi_principal ON angajati_garzi(principal_id);
-- +goose StatementBegin
CREATE TRIGGER angajati_garzi_fara_lant_insert BEFORE INSERT ON angajati_garzi
WHEN EXISTS(SELECT 1 FROM angajati_garzi WHERE angajat_id=NEW.principal_id OR principal_id=NEW.angajat_id)
BEGIN SELECT RAISE(ABORT,'garda_lant_invalid'); END;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER angajati_garzi_fara_lant_update BEFORE UPDATE ON angajati_garzi
WHEN EXISTS(SELECT 1 FROM angajati_garzi WHERE angajat_id<>OLD.angajat_id AND (angajat_id=NEW.principal_id OR principal_id=NEW.angajat_id))
BEGIN SELECT RAISE(ABORT,'garda_lant_invalid'); END;
-- +goose StatementEnd
