-- F07-T10, D021: programele reale confirmate de utilizator.
-- +goose Up
ALTER TABLE tipuri_schimb ADD COLUMN regim TEXT NOT NULL DEFAULT 'GENERAL'
    CHECK (regim IN ('GENERAL', 'LUCRATOARE', 'NELUCRATOARE', 'INTRARE_ZI_LIBERA'));

-- Numai metadatele șabloanelor originale nemodificate; intervalele rămân.
UPDATE tipuri_schimb SET regim = 'NELUCRATOARE'
WHERE sters_la IS NULL AND activ = 1 AND pauza_minute = 0 AND creat_la = modificat_la
    AND ((id = '10000000-0000-4000-8000-000000000003' AND cod = 'N12' AND ora_start = 1140 AND durata_minute = 720 AND tip = 'NOAPTE'
        AND denumire = '{"ro":"Tură de noapte 12 ore","en":"12-hour night shift"}')
      OR (id = '10000000-0000-4000-8000-000000000004' AND cod = 'G24' AND ora_start = 420 AND durata_minute = 1440 AND tip = 'GARDA'
        AND denumire = '{"ro":"Gardă 24 ore","en":"24-hour on-call shift"}'));

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, regim, activ, creat_la, modificat_la)
SELECT '20000000-0000-4000-8000-000000000001', 'Z7', '{"ro":"Zi 07–14","en":"Day 07–14"}', 420, 420, 0, 'NORMAL', 'LUCRATOARE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ','now'), strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'Z7')
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'Z7')
ON CONFLICT(id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, regim, activ, creat_la, modificat_la)
SELECT '20000000-0000-4000-8000-000000000002', 'S7', '{"ro":"Seară 15–22","en":"Evening 15–22"}', 900, 420, 0, 'TURA', 'LUCRATOARE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ','now'), strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'S7')
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'S7')
ON CONFLICT(id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, regim, activ, creat_la, modificat_la)
SELECT '20000000-0000-4000-8000-000000000003', 'N9', '{"ro":"Noapte 22–07","en":"Night 22–07"}', 1320, 540, 0, 'NOAPTE', 'LUCRATOARE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ','now'), strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'N9')
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'N9')
ON CONFLICT(id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, regim, activ, creat_la, modificat_la)
SELECT '20000000-0000-4000-8000-000000000004', 'IN9', '{"ro":"Intrare în weekend/sărbătoare 22–07","en":"Weekend/holiday entry 22–07"}', 1320, 540, 0, 'NOAPTE', 'INTRARE_ZI_LIBERA', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ','now'), strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'IN9')
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'IN9')
ON CONFLICT(id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, regim, activ, creat_la, modificat_la)
SELECT '20000000-0000-4000-8000-000000000005', 'Z12', '{"ro":"Zi liberă 07–19","en":"Non-working day 07–19"}', 420, 720, 0, 'TURA', 'NELUCRATOARE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ','now'), strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'Z12')
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'Z12')
ON CONFLICT(id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, regim, activ, creat_la, modificat_la)
SELECT '20000000-0000-4000-8000-000000000006', 'N12', '{"ro":"Tură de noapte 12 ore","en":"12-hour night shift"}', 1140, 720, 0, 'NOAPTE', 'NELUCRATOARE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ','now'), strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'N12')
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'N12')
ON CONFLICT(id) DO NOTHING;

INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, regim, activ, creat_la, modificat_la)
SELECT '20000000-0000-4000-8000-000000000007', 'G24', '{"ro":"Gardă 24 ore","en":"24-hour on-call shift"}', 420, 1440, 0, 'GARDA', 'NELUCRATOARE', 1,
    strftime('%Y-%m-%dT%H:%M:%SZ','now'), strftime('%Y-%m-%dT%H:%M:%SZ','now')
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = 'G24')
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = 'G24')
ON CONFLICT(id) DO NOTHING;
