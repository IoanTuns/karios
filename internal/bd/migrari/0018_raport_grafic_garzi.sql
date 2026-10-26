-- F11, D040: Raport dedicat pentru graficul lunar de gărzi și schimburi.
-- +goose Up
INSERT INTO definitii_rapoarte(id,denumire,antet,continut,subsol,pagina,creat_la,modificat_la)
VALUES('raport-garzi','{"ro":"Grafic gărzi și schimburi","en":"On-call and shift schedule"}',
'{"stanga":{"ro":"Unitate: {{sectie}}\n{{departament}}","en":"Unit: {{sectie}}\n{{departament}}"},"centru":{"ro":"GRAFIC LUNAR DE GĂRZI ȘI SCHIMBURI · {{perioada}}","en":"MONTHLY ON-CALL AND SHIFT SCHEDULE · {{perioada}}"},"dreapta":{"ro":"APROBAT,\nManager","en":"APPROVED,\nManager"},"imagini":{}}',
'{"sursa":"GRAFIC_GARZI","coloane":["angajat","marca","functie","zile","total_garzi","total_ore"],"text":{"ro":"","en":""}}',
'{"stanga":{"ro":"Întocmit: __________________\nCoordonator gărzi","en":"Prepared by: __________________\nOn-call coordinator"},"centru":{"ro":"Avizat: __________________\nȘef secție","en":"Endorsed: __________________\nDepartment head"},"dreapta":{"ro":"Aprobat: __________________\nDirector medical","en":"Approved: __________________\nMedical director"},"imagini":{}}',
'{"format":"A3","orientare":"PEISAJ","margine_mm":10,"font_pt":7}',strftime('%Y-%m-%dT%H:%M:%SZ','now'),strftime('%Y-%m-%dT%H:%M:%SZ','now'));
