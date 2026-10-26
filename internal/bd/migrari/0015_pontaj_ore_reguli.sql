-- F09-T09, D028: reguli explicite și al doilea raport lunar; codurile din schițe rămân intacte.
-- +goose Up
CREATE TABLE configuratie_pontaj (
    id INTEGER PRIMARY KEY CHECK(id=1),
    configuratie TEXT NOT NULL CHECK(json_valid(configuratie)),
    revizie INTEGER NOT NULL DEFAULT 1 CHECK(revizie>0),
    modificat_la TEXT NOT NULL
) STRICT;
INSERT INTO configuratie_pontaj(id,configuratie,modificat_la) VALUES(1,
'{"atribuire":"ZI_INCEPUT","reguli":[{"coloana":"ore_lucrate","sursa":"TOTAL_ORE","coduri":[],"factor_sutimi":100},{"coloana":"T","sursa":"TOTAL_ORE","coduri":[],"factor_sutimi":100},{"coloana":"ZNL","sursa":"ORE_ZNL","coduri":[],"factor_sutimi":100},{"coloana":"CO","sursa":"ZILE_ABSENTA","coduri":["CO"],"factor_sutimi":100},{"coloana":"CEV","sursa":"NEDEFINITA","coduri":[],"factor_sutimi":100},{"coloana":"CFS","sursa":"NEDEFINITA","coduri":[],"factor_sutimi":100},{"coloana":"ore_norma","sursa":"ORE_NORMA","coduri":[],"factor_sutimi":100},{"coloana":"form_prof","sursa":"NEDEFINITA","coduri":[],"factor_sutimi":100},{"coloana":"CM","sursa":"ZILE_ABSENTA","coduri":["CM"],"factor_sutimi":100},{"coloana":"L_FIN","sursa":"NEDEFINITA","coduri":[],"factor_sutimi":100},{"coloana":"deleg","sursa":"NEDEFINITA","coduri":[],"factor_sutimi":100},{"coloana":"ore_garda","sursa":"NEDEFINITA","coduri":[],"factor_sutimi":100}]}',strftime('%Y-%m-%dT%H:%M:%SZ','now'));

INSERT INTO definitii_rapoarte(id,denumire,antet,continut,subsol,pagina,creat_la,modificat_la)
VALUES('raport-pontaj-detaliat','{"ro":"Pontaj detaliat","en":"Detailed timesheet"}',
'{"stanga":{"ro":"","en":""},"centru":{"ro":"PONTAJ DETALIAT · {{perioada}}","en":"DETAILED TIMESHEET · {{perioada}}"},"dreapta":{"ro":"","en":""},"imagini":{}}',
'{"sursa":"PONTAJ_DETALIAT","coloane":["angajat","marca","zile","totaluri"],"text":{"ro":"","en":""}}',
'{"stanga":{"ro":"Întocmit: __________________","en":"Prepared by: __________________"},"centru":{"ro":"","en":""},"dreapta":{"ro":"Verificat: __________________","en":"Checked by: __________________"},"imagini":{}}',
'{"format":"A3","orientare":"PEISAJ","margine_mm":10,"font_pt":6}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now'));

UPDATE definitii_rapoarte
SET continut=json_set(continut,'$.coloane',
    json((SELECT json_group_array(CASE WHEN value='cod' THEN 'ore' ELSE value END) FROM json_each(continut,'$.coloane')))),
    revizie=revizie+1,modificat_la=strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE json_extract(continut,'$.sursa')='FISA_INDIVIDUALA'
  AND EXISTS(SELECT 1 FROM json_each(continut,'$.coloane') WHERE value='cod');
