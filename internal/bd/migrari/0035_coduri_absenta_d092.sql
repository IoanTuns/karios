-- F08-T12 (D092): redenumirea codurilor de absenta CEV->CVE, CFS->CFP, DEPL->DLG, FP->FMP.
-- Se redenumeste doar cand codul vechi exista si cel nou nu: daca exista ambele, nu se
-- amesteca nimic. Tipul isi pastreaza denumirea, categoria, efectul si starea; zilele
-- inregistrate trec pe codul nou. Regulile si modelele de raport care pomenesc codurile
-- vechi sunt traduse la citire (calcul.CoduriRedenumite).
-- +goose Up
CREATE TEMP TABLE redenumiri_d092 (vechi TEXT PRIMARY KEY, nou TEXT NOT NULL);
INSERT INTO redenumiri_d092
SELECT v.vechi, v.nou FROM (
    SELECT 'CEV' AS vechi, 'CVE' AS nou UNION ALL SELECT 'CFS', 'CFP'
    UNION ALL SELECT 'DEPL', 'DLG' UNION ALL SELECT 'FP', 'FMP'
) v
WHERE EXISTS (SELECT 1 FROM coduri_absenta WHERE cod = v.vechi)
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE cod = v.nou);
CREATE TEMP TABLE coduri_d092 AS SELECT * FROM coduri_absenta WHERE cod IN (SELECT vechi FROM redenumiri_d092);
UPDATE coduri_d092 SET cod = (SELECT nou FROM redenumiri_d092 WHERE vechi = coduri_d092.cod);
INSERT INTO coduri_absenta SELECT * FROM coduri_d092;
UPDATE absente_zile SET cod = (SELECT nou FROM redenumiri_d092 WHERE vechi = absente_zile.cod)
WHERE cod IN (SELECT vechi FROM redenumiri_d092);
UPDATE pontaj SET cod_absenta = (SELECT nou FROM redenumiri_d092 WHERE vechi = pontaj.cod_absenta)
WHERE cod_absenta IN (SELECT vechi FROM redenumiri_d092);
DELETE FROM coduri_absenta WHERE cod IN (SELECT vechi FROM redenumiri_d092);
DROP TABLE coduri_d092;
DROP TABLE redenumiri_d092;

-- +goose Down
CREATE TEMP TABLE redenumiri_d092 (vechi TEXT PRIMARY KEY, nou TEXT NOT NULL);
INSERT INTO redenumiri_d092
SELECT v.vechi, v.nou FROM (
    SELECT 'CVE' AS vechi, 'CEV' AS nou UNION ALL SELECT 'CFP', 'CFS'
    UNION ALL SELECT 'DLG', 'DEPL' UNION ALL SELECT 'FMP', 'FP'
) v
WHERE EXISTS (SELECT 1 FROM coduri_absenta WHERE cod = v.vechi)
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE cod = v.nou);
CREATE TEMP TABLE coduri_d092 AS SELECT * FROM coduri_absenta WHERE cod IN (SELECT vechi FROM redenumiri_d092);
UPDATE coduri_d092 SET cod = (SELECT nou FROM redenumiri_d092 WHERE vechi = coduri_d092.cod);
INSERT INTO coduri_absenta SELECT * FROM coduri_d092;
UPDATE absente_zile SET cod = (SELECT nou FROM redenumiri_d092 WHERE vechi = absente_zile.cod)
WHERE cod IN (SELECT vechi FROM redenumiri_d092);
UPDATE pontaj SET cod_absenta = (SELECT nou FROM redenumiri_d092 WHERE vechi = pontaj.cod_absenta)
WHERE cod_absenta IN (SELECT vechi FROM redenumiri_d092);
DELETE FROM coduri_absenta WHERE cod IN (SELECT vechi FROM redenumiri_d092);
DROP TABLE coduri_d092;
DROP TABLE redenumiri_d092;
