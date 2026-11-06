-- F07-T24: categorii configurabile si functii predefinite.
-- +goose Up
CREATE TABLE categorii_schimburi (
 id TEXT PRIMARY KEY, denumire TEXT NOT NULL CHECK(json_valid(denumire)),
 activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN (0,1)), modificat_la TEXT NOT NULL
);
INSERT INTO categorii_schimburi VALUES ('LAB','{"ro":"Laborator","en":"Laboratory"}',1,strftime('%Y-%m-%dT%H:%M:%fZ','now'));
INSERT INTO categorii_schimburi VALUES ('CUSTOM','{"ro":"Personalizate","en":"Custom"}',1,strftime('%Y-%m-%dT%H:%M:%fZ','now'));
ALTER TABLE tipuri_schimb ADD COLUMN categorie_noua TEXT NOT NULL DEFAULT 'CUSTOM';
UPDATE tipuri_schimb SET categorie_noua=categorie;
ALTER TABLE tipuri_schimb DROP COLUMN categorie;
ALTER TABLE tipuri_schimb RENAME COLUMN categorie_noua TO categorie;
-- +goose StatementBegin
CREATE TRIGGER categorie_schimb_insert BEFORE INSERT ON tipuri_schimb
WHEN NOT EXISTS(SELECT 1 FROM categorii_schimburi WHERE id=NEW.categorie)
BEGIN SELECT RAISE(ABORT,'categorie_inexistenta'); END;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER categorie_schimb_update BEFORE UPDATE OF categorie ON tipuri_schimb
WHEN NOT EXISTS(SELECT 1 FROM categorii_schimburi WHERE id=NEW.categorie)
BEGIN SELECT RAISE(ABORT,'categorie_inexistenta'); END;
-- +goose StatementEnd
CREATE TABLE functii (
 id TEXT PRIMARY KEY, denumire TEXT NOT NULL CHECK(json_valid(denumire)),
 activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN (0,1)), modificat_la TEXT NOT NULL
);
INSERT INTO functii SELECT lower(hex(randomblob(16))),json_object('ro',nume,'en',nume),1,strftime('%Y-%m-%dT%H:%M:%fZ','now')
FROM (SELECT DISTINCT trim(functie) AS nume FROM angajati WHERE functie IS NOT NULL AND trim(functie)<>'');
CREATE TABLE angajati_functii (
 angajat_id TEXT PRIMARY KEY REFERENCES angajati(id), functie_id TEXT NOT NULL REFERENCES functii(id)
);
INSERT INTO angajati_functii SELECT a.id,f.id FROM angajati a JOIN functii f ON json_extract(f.denumire,'$.ro')=trim(a.functie);
-- Un import care schimba textul functiei nu pastreaza o asociere veche incompatibila.
-- +goose StatementBegin
CREATE TRIGGER functie_angajat_text AFTER UPDATE OF functie ON angajati
WHEN NEW.functie IS NOT OLD.functie
BEGIN DELETE FROM angajati_functii WHERE angajat_id=NEW.id; END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER categorie_schimb_delete BEFORE DELETE ON categorii_schimburi
WHEN EXISTS(SELECT 1 FROM tipuri_schimb WHERE categorie=OLD.id)
BEGIN SELECT RAISE(ABORT,'categorie_utilizata'); END;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER categorie_schimb_identitate BEFORE UPDATE OF id ON categorii_schimburi
WHEN NEW.id<>OLD.id AND EXISTS(SELECT 1 FROM tipuri_schimb WHERE categorie=OLD.id)
BEGIN SELECT RAISE(ABORT,'categorie_utilizata'); END;
-- +goose StatementEnd
