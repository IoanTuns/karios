-- F11-T21: pozitii distincte pe zi; proiectie zilnica a codurilor in Pontaj.
-- +goose Up
DROP INDEX planificare_angajat_data;
CREATE INDEX planificare_angajat_data ON planificare(angajat_id,data) WHERE activ=1;
ALTER TABLE pontaj_ciorne RENAME TO pontaj_ciorne_vechi;
CREATE TABLE pontaj_ciorne (
 angajat_id TEXT NOT NULL REFERENCES angajati(id),
 data TEXT NOT NULL CHECK(data IS date(data)),
 valoare TEXT NOT NULL DEFAULT '' CHECK(length(valoare)<=1024 AND instr(valoare,char(0))=0 AND instr(valoare,char(9))=0 AND instr(valoare,char(10))=0 AND instr(valoare,char(13))=0),
 modificat_la TEXT NOT NULL,
 revizie INTEGER NOT NULL DEFAULT 1 CHECK(revizie>0),
 PRIMARY KEY(angajat_id,data)
) STRICT;
INSERT INTO pontaj_ciorne SELECT angajat_id,data,valoare,modificat_la,revizie FROM pontaj_ciorne_vechi;
DROP TABLE pontaj_ciorne_vechi;
CREATE INDEX pontaj_ciorne_data ON pontaj_ciorne(data,angajat_id);
