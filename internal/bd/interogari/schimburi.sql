-- name: ListaSchimburi :many
SELECT id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ
FROM tipuri_schimb WHERE sters_la IS NULL ORDER BY cod;

-- name: CitesteSchimb :one
SELECT id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ
FROM tipuri_schimb WHERE id = ? AND sters_la IS NULL;

-- name: CitesteCodSchimb :one
SELECT id FROM tipuri_schimb WHERE upper(cod) = upper(?) AND sters_la IS NULL;

-- name: AdaugaSchimb :exec
INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ModificaSchimb :execrows
UPDATE tipuri_schimb SET cod = ?, denumire = ?, ora_start = ?, durata_minute = ?, pauza_minute = ?, tip = ?, activ = ?, modificat_la = ?
WHERE id = ? AND sters_la IS NULL;

-- name: SchimbaActivSchimb :execrows
UPDATE tipuri_schimb SET activ = ?, modificat_la = ? WHERE id = ? AND sters_la IS NULL;

-- name: AdaugaSchimbInitial :execrows
INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, activ, creat_la, modificat_la)
SELECT sqlc.arg(id), sqlc.arg(cod), sqlc.arg(denumire), sqlc.arg(ora_start), sqlc.arg(durata_minute), 0, sqlc.arg(tip), 1, sqlc.arg(creat_la), sqlc.arg(modificat_la)
WHERE NOT EXISTS (SELECT 1 FROM tipuri_schimb WHERE upper(cod) = upper(sqlc.arg(cod)))
  AND NOT EXISTS (SELECT 1 FROM coduri_absenta WHERE upper(cod) = upper(sqlc.arg(cod)))
ON CONFLICT (id) DO NOTHING;
