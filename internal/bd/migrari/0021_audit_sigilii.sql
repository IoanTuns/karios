-- F06-T08, F12-T04, D044: sigiliile tabelelor, lanțul jurnalului și amprenta lunii.
-- Prima sigilare o face aplicația, imediat după această migrare.
-- +goose Up
CREATE TABLE sigilii (
    tabel   TEXT PRIMARY KEY,
    valoare BLOB NOT NULL
) STRICT;

ALTER TABLE jurnal_modificari ADD COLUMN mac BLOB;
ALTER TABLE stare_luni ADD COLUMN amprenta TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE stare_luni DROP COLUMN amprenta;
ALTER TABLE jurnal_modificari DROP COLUMN mac;
DROP TABLE sigilii;
