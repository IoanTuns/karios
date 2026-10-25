-- name: ListaTipuriAbsenta :many
SELECT * FROM coduri_absenta ORDER BY cod;

-- name: CitesteTipAbsenta :one
SELECT * FROM coduri_absenta WHERE cod = ?;

-- name: CautaCodAbsenta :one
SELECT cod FROM coduri_absenta WHERE upper(cod) = upper(?);

-- name: SalveazaTipAbsenta :exec
INSERT INTO coduri_absenta (cod, denumire, categorie, activ, creat_la, modificat_la)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(cod) DO UPDATE SET denumire = excluded.denumire, categorie = excluded.categorie,
    activ = excluded.activ, modificat_la = excluded.modificat_la;

-- name: TipAbsentaFolosit :one
SELECT EXISTS (SELECT 1 FROM absente_zile WHERE cod = sqlc.arg(cod))
    OR EXISTS (SELECT 1 FROM pontaj WHERE cod_absenta = sqlc.arg(cod));

-- name: ListaAbsenteInterval :many
SELECT * FROM absente_zile WHERE data BETWEEN sqlc.arg(de_la) AND sqlc.arg(pana_la)
ORDER BY data, angajat_id;

-- name: CitesteZiAbsenta :one
SELECT * FROM absente_zile WHERE angajat_id = ? AND data = ?;

-- name: SalveazaZiAbsenta :exec
INSERT INTO absente_zile (angajat_id, data, cod, activ, creat_la, modificat_la)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(angajat_id, data) DO UPDATE SET cod = excluded.cod, activ = excluded.activ,
    modificat_la = excluded.modificat_la;

-- name: PontajExistentPeZi :one
SELECT EXISTS (SELECT 1 FROM pontaj WHERE angajat_id = sqlc.arg(angajat_id) AND sters_la IS NULL
    AND (data = sqlc.arg(data) OR (julianday(start) < julianday(sqlc.arg(sfarsit_utc)) AND julianday(sfarsit) > julianday(sqlc.arg(inceput_utc)))));
