-- +goose Up
CREATE TABLE categorii_absente (
 cod TEXT PRIMARY KEY,
 denumire TEXT NOT NULL CHECK(json_valid(denumire)),
 natura TEXT NOT NULL CHECK(natura IN ('ODIHNA','MEDICAL','SPECIAL')),
 efect_ore TEXT NOT NULL CHECK(efect_ore IN ('ANTERIOR','FARA_ORE','NORMA','PROGRAM')),
 activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN(0,1)),
 revizie INTEGER NOT NULL DEFAULT 1,
 CHECK(natura='SPECIAL' OR efect_ore IN ('ANTERIOR','FARA_ORE'))
);
INSERT INTO categorii_absente(cod,denumire,natura,efect_ore) VALUES
 ('ODIHNA','{"ro":"Concediu de odihna","en":"Annual leave"}','ODIHNA','ANTERIOR'),
 ('MEDICAL','{"ro":"Concediu medical","en":"Medical leave"}','MEDICAL','ANTERIOR'),
 ('SPECIAL','{"ro":"Alte absente (reguli existente)","en":"Other absences (existing rules)"}','SPECIAL','ANTERIOR');
INSERT INTO categorii_absente(cod,denumire,natura,efect_ore)
 SELECT DISTINCT categorie||'_'||efect_ore,
 json_object('ro',CASE efect_ore WHEN 'NORMA' THEN 'Deplasari - norma zilnica' WHEN 'PROGRAM' THEN 'Deplasari - program planificat' ELSE 'Concedii cu reducerea obligatiei' END,
 'en',CASE efect_ore WHEN 'NORMA' THEN 'Travel - daily contract hours' WHEN 'PROGRAM' THEN 'Travel - planned hours' ELSE 'Leave reducing obligation' END),categorie,efect_ore
 FROM coduri_absenta WHERE efect_ore<>'ANTERIOR';
ALTER TABLE coduri_absenta ADD COLUMN categorie_id TEXT REFERENCES categorii_absente(cod);
UPDATE coduri_absenta SET categorie_id=CASE WHEN efect_ore='ANTERIOR' THEN categorie ELSE categorie||'_'||efect_ore END;

-- +goose Down
ALTER TABLE coduri_absenta DROP COLUMN categorie_id;
DROP TABLE categorii_absente;
