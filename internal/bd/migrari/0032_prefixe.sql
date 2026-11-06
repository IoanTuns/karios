-- F07-T40 (D089): prefixe (titulaturi) pentru angajati, optionale.
-- +goose Up
CREATE TABLE prefixe (
 id TEXT PRIMARY KEY, denumire TEXT NOT NULL CHECK(json_valid(denumire)),
 activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN (0,1)), modificat_la TEXT NOT NULL
);
-- Lista de pornire, editabila; forma engleza porneste identica (fara traduceri inventate).
INSERT INTO prefixe SELECT lower(hex(randomblob(16))),json_object('ro',p,'en',p),1,strftime('%Y-%m-%dT%H:%M:%fZ','now')
FROM (SELECT 'Dr.' AS p UNION ALL SELECT 'Biol.' UNION ALL SELECT 'Chim.' UNION ALL SELECT 'Farm.'
      UNION ALL SELECT 'As.' UNION ALL SELECT 'Ing.' UNION ALL SELECT 'Ec.' UNION ALL SELECT 'Psih.');
CREATE TABLE angajati_prefixe (
 angajat_id TEXT PRIMARY KEY REFERENCES angajati(id), prefix_id TEXT NOT NULL REFERENCES prefixe(id)
);
-- +goose StatementBegin
CREATE TRIGGER prefix_angajat_inexistent BEFORE INSERT ON angajati_prefixe
WHEN NOT EXISTS(SELECT 1 FROM prefixe WHERE id=NEW.prefix_id)
BEGIN SELECT RAISE(ABORT,'prefix_inexistent'); END;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER prefix_utilizat BEFORE DELETE ON prefixe
WHEN EXISTS(SELECT 1 FROM angajati_prefixe WHERE prefix_id=OLD.id)
BEGIN SELECT RAISE(ABORT,'prefix_utilizat'); END;
-- +goose StatementEnd
