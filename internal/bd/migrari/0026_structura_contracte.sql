-- +goose Up
CREATE TABLE grupuri_profesionale(id TEXT PRIMARY KEY,denumire TEXT NOT NULL CHECK(json_valid(denumire)),activ INTEGER NOT NULL CHECK(activ IN(0,1)),modificat_la TEXT NOT NULL,ordine INTEGER NOT NULL DEFAULT 0);
ALTER TABLE functii ADD COLUMN grup_id TEXT REFERENCES grupuri_profesionale(id);
ALTER TABLE functii ADD COLUMN ordine INTEGER NOT NULL DEFAULT 0;
ALTER TABLE functii ADD COLUMN conducere INTEGER NOT NULL DEFAULT 0 CHECK(conducere IN(0,1));
CREATE TABLE tipuri_contract(id TEXT PRIMARY KEY,denumire TEXT NOT NULL CHECK(json_valid(denumire)),activ INTEGER NOT NULL CHECK(activ IN(0,1)),modificat_la TEXT NOT NULL,normare TEXT NOT NULL DEFAULT 'FARA' CHECK(normare IN('GLOBALA','INDIVIDUALA','ZILNICA','LUNARA','FARA')),norma_minute INTEGER CHECK(norma_minute>0),categorie_program_id TEXT REFERENCES categorii_program(id),garda INTEGER NOT NULL DEFAULT 0 CHECK(garda IN(0,1)),CHECK((normare IN('ZILNICA','LUNARA') AND norma_minute IS NOT NULL) OR (normare IN('GLOBALA','INDIVIDUALA','FARA') AND norma_minute IS NULL)));
INSERT INTO tipuri_contract VALUES('INTREAGA','{"ro":"Norma intreaga","en":"Full time"}',1,'initial','GLOBALA',NULL,NULL,0),('PARTIALA','{"ro":"Norma partiala","en":"Part time"}',1,'initial','INDIVIDUALA',NULL,NULL,0),('GARDA','{"ro":"Garda","en":"On-call"}',1,'initial','FARA',NULL,NULL,1);
CREATE TABLE angajati_contracte(angajat_id TEXT PRIMARY KEY REFERENCES angajati(id),contract_id TEXT NOT NULL REFERENCES tipuri_contract(id));
INSERT INTO angajati_contracte SELECT id,tip_contract FROM angajati;
UPDATE angajati_contracte SET contract_id='GARDA' WHERE angajat_id IN(SELECT angajat_id FROM angajati_garzi);

-- +goose StatementBegin
CREATE TRIGGER contract_istoric_actualizat AFTER UPDATE OF tip_contract ON angajati WHEN NEW.tip_contract<>OLD.tip_contract
BEGIN UPDATE angajati_contracte SET contract_id=NEW.tip_contract WHERE angajat_id=NEW.id AND contract_id IN('INTREAGA','PARTIALA'); END;
-- +goose StatementEnd
