-- name: ListaPlanificare :many
SELECT * FROM planificare WHERE data>=sqlc.arg(data_de_la) AND data<=sqlc.arg(data_pana_la)
ORDER BY data,angajat_id,id;

-- name: CitestePozitiePlanificare :one
SELECT * FROM planificare WHERE id=?;

-- name: AdaugaPozitiePlanificare :exec
INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,activ,creat_la,modificat_la)
VALUES (?,?,?,?,?,?,?,?);

-- name: ModificaPozitiePlanificare :execrows
UPDATE planificare SET schimb_id=?,activ=?,modificat_la=?,revizie=revizie+1
WHERE id=? AND revizie=?;

-- name: MarcheazaPlanificarePublicata :execrows
UPDATE planificare SET publicat_plan_revizie=?,publicat_pontaj_revizie=?,publicat_valoare=?,publicat_la=?
WHERE id=? AND revizie=? AND activ=1;
