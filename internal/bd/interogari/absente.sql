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
SELECT z.angajat_id, z.data, z.cod, z.activ, z.creat_la, z.modificat_la, z.revizie,
    z.nr_inregistrare, COALESCE(z.inlocuitor_id, '') AS inlocuitor_id,
    z.stare_aprobare, z.aprobat_de,
    COALESCE(z.document_id, '') AS document_id,
    COALESCE(d.nume_fisier, '') AS document_nume,
    COALESCE(d.tip_mime, '') AS document_mime
FROM absente_zile z
LEFT JOIN absente_documente d ON d.id = z.document_id
WHERE z.data BETWEEN sqlc.arg(de_la) AND sqlc.arg(pana_la)
ORDER BY z.data, z.angajat_id;

-- name: CitesteZiAbsenta :one
SELECT * FROM absente_zile WHERE angajat_id = ? AND data = ?;

-- name: AdaugaDocumentAbsenta :exec
INSERT INTO absente_documente (id, nume_fisier, tip_mime, continut, creat_la)
VALUES (?, ?, ?, ?, ?);

-- name: CitesteDocumentAbsenta :one
SELECT id, nume_fisier, tip_mime, continut, creat_la
FROM absente_documente WHERE id = ?;

-- name: ExistaDocumentAbsenta :one
SELECT id, nume_fisier, tip_mime FROM absente_documente WHERE id = ?;

-- name: SalveazaZiAbsenta :one
INSERT INTO absente_zile (angajat_id, data, cod, activ, nr_inregistrare, inlocuitor_id, stare_aprobare, aprobat_de, document_id, creat_la, modificat_la)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(angajat_id, data) DO UPDATE SET cod = excluded.cod, activ = excluded.activ,
    nr_inregistrare = excluded.nr_inregistrare, inlocuitor_id = excluded.inlocuitor_id,
    stare_aprobare = excluded.stare_aprobare, aprobat_de = excluded.aprobat_de,
    document_id = excluded.document_id,
    modificat_la = excluded.modificat_la, revizie = absente_zile.revizie + 1
RETURNING revizie;

-- name: PontajExistentPeZi :one
SELECT EXISTS (SELECT 1 FROM pontaj WHERE angajat_id = sqlc.arg(angajat_id) AND sters_la IS NULL
    AND (data = sqlc.arg(data) OR (julianday(start) < julianday(sqlc.arg(sfarsit_utc)) AND julianday(sfarsit) > julianday(sqlc.arg(inceput_utc)))));
