-- name: ListaCiorneLuna :many
SELECT angajat_id, data, valoare, modificat_la FROM pontaj_ciorne
WHERE data >= sqlc.arg(data_de_la) AND data <= sqlc.arg(data_pana_la)
ORDER BY angajat_id, data;

-- name: CitesteCiorna :one
SELECT valoare, modificat_la FROM pontaj_ciorne WHERE angajat_id = ? AND data = ?;

-- name: SalveazaCiorna :exec
INSERT INTO pontaj_ciorne (angajat_id, data, valoare, modificat_la) VALUES (?, ?, ?, ?)
ON CONFLICT (angajat_id, data) DO UPDATE SET valoare = excluded.valoare, modificat_la = excluded.modificat_la;
