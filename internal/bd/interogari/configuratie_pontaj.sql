-- name: CitesteConfiguratiePontaj :one
SELECT configuratie,revizie FROM configuratie_pontaj WHERE id=1;

-- name: SalveazaConfiguratiePontaj :one
UPDATE configuratie_pontaj SET configuratie=?,revizie=revizie+1,modificat_la=?
WHERE id=1 AND revizie=? RETURNING revizie;
