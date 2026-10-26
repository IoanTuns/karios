-- F07-T14: categorii reutilizabile de program, distincte de LAB/CUSTOM.
-- +goose Up
CREATE TABLE categorii_program (
    id TEXT PRIMARY KEY,
    denumire TEXT NOT NULL CHECK(json_valid(denumire)),
    schimb_implicit_id TEXT REFERENCES tipuri_schimb(id),
    activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN (0,1)),
    revizie INTEGER NOT NULL DEFAULT 1 CHECK(revizie > 0),
    creat_la TEXT NOT NULL,
    modificat_la TEXT NOT NULL
) STRICT;
CREATE TABLE categorii_program_schimburi (
    categorie_id TEXT NOT NULL REFERENCES categorii_program(id),
    schimb_id TEXT NOT NULL REFERENCES tipuri_schimb(id),
    PRIMARY KEY(categorie_id, schimb_id)
) STRICT;
ALTER TABLE angajati ADD COLUMN categorie_program_id TEXT REFERENCES categorii_program(id);
CREATE INDEX angajati_categorie_program ON angajati(categorie_program_id) WHERE sters_la IS NULL;

INSERT INTO categorii_program(id,denumire,creat_la,modificat_la) VALUES
('program-dimineata','{"ro":"Dimineață","en":"Morning"}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now')),
('program-dupa-amiaza','{"ro":"După-amiază","en":"Afternoon"}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now')),
('program-garzi','{"ro":"Gărzi","en":"On-call shifts"}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now')),
('program-dimineata-garzi','{"ro":"Dimineață și gărzi","en":"Morning and on-call shifts"}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now')),
('program-dupa-amiaza-garzi','{"ro":"După-amiază și gărzi","en":"Afternoon and on-call shifts"}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now')),
('program-mixt','{"ro":"Program mixt","en":"Mixed schedule"}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now'));

INSERT INTO categorii_program_schimburi(categorie_id,schimb_id)
SELECT c.id,s.id FROM categorii_program c JOIN tipuri_schimb s ON s.sters_la IS NULL
WHERE (c.id IN ('program-dimineata','program-dimineata-garzi','program-mixt') AND s.cod='Z7')
   OR (c.id IN ('program-dupa-amiaza','program-dupa-amiaza-garzi','program-mixt') AND s.cod='S7')
   OR (c.id='program-mixt' AND s.cod='N9')
   OR (c.id IN ('program-garzi','program-dimineata-garzi','program-dupa-amiaza-garzi','program-mixt') AND s.cod IN ('IN9','Z12','N12','G24'));
UPDATE categorii_program SET schimb_implicit_id=(SELECT s.id FROM tipuri_schimb s
    WHERE s.cod=CASE WHEN categorii_program.id IN ('program-dimineata','program-dimineata-garzi') THEN 'Z7' ELSE 'S7' END
      AND s.sters_la IS NULL AND s.activ=1)
WHERE id IN ('program-dimineata','program-dimineata-garzi','program-dupa-amiaza','program-dupa-amiaza-garzi');
