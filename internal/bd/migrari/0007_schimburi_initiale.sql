-- F07-T09, D019: șabloanele din PDF și orele confirmate de Ioan Tuns.
-- Pauza folosește implicitul tehnic existent, 0; nicio regulă de plată nouă.
-- Configurațiile preexistente, inclusiv inactive/șterse logic, nu sunt suprascrise.
-- +goose Up

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
SELECT '10000000-0000-4000-8000-000000000001', '7', '{"ro":"Program 7 ore","en":"7-hour schedule"}', 420, 420, 0, 'NORMAL', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = '7')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
SELECT '10000000-0000-4000-8000-000000000002', '12', '{"ro":"Tură 12 ore","en":"12-hour shift"}', 840, 720, 0, 'TURA', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = '12')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
SELECT '10000000-0000-4000-8000-000000000003', 'N12', '{"ro":"Tură de noapte 12 ore","en":"12-hour night shift"}', 1140, 720, 0, 'NOAPTE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'N12')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
SELECT '10000000-0000-4000-8000-000000000004', 'G24', '{"ro":"Gardă 24 ore","en":"24-hour on-call shift"}', 420, 1440, 0, 'GARDA', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'G24')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
SELECT '10000000-0000-4000-8000-000000000005', 'S1', '{"ro":"Schimb 1","en":"Shift 1"}', 360, 480, 0, 'TURA', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'S1')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
SELECT '10000000-0000-4000-8000-000000000006', 'S3', '{"ro":"Schimb 3 de noapte","en":"Night shift 3"}', 1320, 480, 0, 'NOAPTE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'S3')
ON CONFLICT (id) DO NOTHING;
