-- name: ListaCategoriiProgram :many
SELECT * FROM categorii_program ORDER BY denumire;

-- name: CitesteCategorieProgram :one
SELECT * FROM categorii_program WHERE id = ?;

-- name: ListaSchimburiCategorieProgram :many
SELECT schimb_id FROM categorii_program_schimburi WHERE categorie_id = ? ORDER BY schimb_id;

-- name: AdaugaCategorieProgram :exec
INSERT INTO categorii_program(id,denumire,schimb_implicit_id,activ,creat_la,modificat_la)
VALUES (?,?,?,?,?,?);

-- name: ModificaCategorieProgram :execrows
UPDATE categorii_program SET denumire=?,schimb_implicit_id=?,activ=?,modificat_la=?,revizie=revizie+1
WHERE id=? AND revizie=?;

-- name: GolesteSchimburiCategorieProgram :exec
DELETE FROM categorii_program_schimburi WHERE categorie_id=?;

-- name: AdaugaSchimbCategorieProgram :exec
INSERT INTO categorii_program_schimburi(categorie_id,schimb_id) VALUES (?,?);

-- name: ListaPontajePrecompletare :many
SELECT angajat_id,data,start,sfarsit FROM pontaj
WHERE sters_la IS NULL AND data >= sqlc.arg(data_de_la) AND data <= sqlc.arg(data_pana_la);
