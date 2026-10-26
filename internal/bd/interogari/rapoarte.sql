-- name: ListaDefinitiiRapoarte :many
SELECT * FROM definitii_rapoarte ORDER BY id;

-- name: CitesteDefinitieRaport :one
SELECT * FROM definitii_rapoarte WHERE id = ?;

-- name: AdaugaDefinitieRaport :exec
INSERT INTO definitii_rapoarte (id,denumire,antet,continut,subsol,pagina,activ,creat_la,modificat_la)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ModificaDefinitieRaport :one
UPDATE definitii_rapoarte SET denumire = ?, antet = ?, continut = ?, subsol = ?, pagina = ?, activ = ?,
    revizie = revizie + 1, modificat_la = ? WHERE id = ? AND revizie = ?
RETURNING revizie;

