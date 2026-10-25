-- name: ListaCiorneLuna :many
SELECT angajat_id, data, valoare, modificat_la, revizie FROM pontaj_ciorne
WHERE data >= sqlc.arg(data_de_la) AND data <= sqlc.arg(data_pana_la)
ORDER BY angajat_id, data;

-- name: CitesteCiorna :one
SELECT valoare, modificat_la, revizie FROM pontaj_ciorne WHERE angajat_id = ? AND data = ?;

-- name: SalveazaCiorna :one
INSERT INTO pontaj_ciorne (angajat_id, data, valoare, modificat_la) VALUES (?, ?, ?, ?)
ON CONFLICT (angajat_id, data) DO UPDATE SET valoare = excluded.valoare, modificat_la = excluded.modificat_la, revizie = pontaj_ciorne.revizie + 1
RETURNING revizie;
