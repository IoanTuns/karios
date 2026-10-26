-- name: ListaSarbatori :many
SELECT data, denumire FROM sarbatori WHERE an = ? ORDER BY data;

-- name: CitesteSarbatoare :one
SELECT data, denumire FROM sarbatori WHERE data = ?;

-- name: AdaugaSarbatoare :execrows
INSERT INTO sarbatori (data, denumire, creat_la, modificat_la) VALUES (?, ?, ?, ?)
ON CONFLICT (data) DO NOTHING;

-- name: ModificaSarbatoare :execrows
UPDATE sarbatori SET data = sqlc.arg(data_noua), denumire = sqlc.arg(denumire), modificat_la = sqlc.arg(modificat_la)
WHERE data = sqlc.arg(data_initiala);
