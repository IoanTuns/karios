-- F09-T07, D024: definiții independente cu cele trei componente configurabile.
-- 0011 este rezervată categoriilor de schimburi din taskul paralel F07-T11.
-- +goose Up
CREATE TABLE definitii_rapoarte (
    id TEXT PRIMARY KEY,
    denumire TEXT NOT NULL CHECK (json_valid(denumire)),
    antet TEXT NOT NULL CHECK (json_valid(antet)),
    continut TEXT NOT NULL CHECK (json_valid(continut)),
    subsol TEXT NOT NULL CHECK (json_valid(subsol)),
    pagina TEXT NOT NULL CHECK (json_valid(pagina)),
    activ INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0,1)),
    revizie INTEGER NOT NULL DEFAULT 1 CHECK (revizie > 0),
    creat_la TEXT NOT NULL,
    modificat_la TEXT NOT NULL
) STRICT;

INSERT INTO definitii_rapoarte (id,denumire,antet,continut,subsol,pagina,creat_la,modificat_la)
VALUES ('raport-pontaj','{"ro":"Pontaj lunar","en":"Monthly timesheet"}','{"stanga":{"ro":"","en":""},"centru":{"ro":"PONTAJ · {{perioada}}","en":"TIMESHEET · {{perioada}}"},"dreapta":{"ro":"","en":""}}','{"sursa":"PONTAJ_LUNAR","coloane":["angajat","marca","zile"],"text":{"ro":"","en":""}}','{"stanga":{"ro":"Întocmit: __________________","en":"Prepared by: __________________"},"centru":{"ro":"","en":""},"dreapta":{"ro":"Verificat: __________________","en":"Checked by: __________________"}}','{"format":"A3","orientare":"PEISAJ","margine_mm":10,"font_pt":8}',
    strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now'));

INSERT INTO definitii_rapoarte (id,denumire,antet,continut,subsol,pagina,creat_la,modificat_la)
VALUES ('raport-fisa','{"ro":"Fișă individuală","en":"Individual timesheet"}','{"stanga":{"ro":"","en":""},"centru":{"ro":"FIȘĂ INDIVIDUALĂ · {{angajat}}\n{{perioada}}","en":"INDIVIDUAL TIMESHEET · {{angajat}}\n{{perioada}}"},"dreapta":{"ro":"","en":""}}','{"sursa":"FISA_INDIVIDUALA","coloane":["data","zi","cod","descriere"],"text":{"ro":"","en":""}}','{"stanga":{"ro":"Întocmit: __________________","en":"Prepared by: __________________"},"centru":{"ro":"","en":""},"dreapta":{"ro":"Verificat: __________________","en":"Checked by: __________________"}}','{"format":"A4","orientare":"PORTRET","margine_mm":10,"font_pt":8}',
    strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now'));

INSERT INTO definitii_rapoarte (id,denumire,antet,continut,subsol,pagina,creat_la,modificat_la)
VALUES ('raport-absente','{"ro":"Evidența absențelor","en":"Absence register"}','{"stanga":{"ro":"","en":""},"centru":{"ro":"EVIDENȚA ABSENȚELOR · {{an}}","en":"ABSENCE REGISTER · {{an}}"},"dreapta":{"ro":"","en":""}}','{"sursa":"ABSENTE_ANUALE","coloane":["data","angajat","marca","cod","descriere","categorie"],"text":{"ro":"","en":""}}','{"stanga":{"ro":"Întocmit: __________________","en":"Prepared by: __________________"},"centru":{"ro":"","en":""},"dreapta":{"ro":"Verificat: __________________","en":"Checked by: __________________"}}','{"format":"A4","orientare":"PEISAJ","margine_mm":10,"font_pt":8}',
    strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now'));
