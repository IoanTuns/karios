-- name: ListaAngajati :many
SELECT id, nume, prenume, marca, functie, sectie_id, departament_id, tip_contract,
    norma_partiala_minute, concediu_zile_an, activ
FROM angajati WHERE sters_la IS NULL
ORDER BY nume COLLATE NOCASE, prenume COLLATE NOCASE, marca;

-- name: CitesteAngajat :one
SELECT id, nume, prenume, marca, functie, sectie_id, departament_id, tip_contract,
    norma_partiala_minute, concediu_zile_an, activ
FROM angajati WHERE id = ? AND sters_la IS NULL;

-- name: AdaugaAngajat :exec
INSERT INTO angajati (id, nume, prenume, marca, functie, sectie_id, departament_id, tip_contract,
    norma_partiala_minute, concediu_zile_an, activ, creat_la, modificat_la)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ModificaAngajat :execrows
UPDATE angajati SET nume = ?, prenume = ?, marca = ?, functie = ?, sectie_id = ?, departament_id = ?,
    tip_contract = ?, norma_partiala_minute = ?, concediu_zile_an = ?, activ = ?, modificat_la = ?
WHERE id = ? AND sters_la IS NULL;

-- name: SchimbaActivAngajat :execrows
UPDATE angajati SET activ = ?, modificat_la = ? WHERE id = ? AND sters_la IS NULL;

-- name: ListaSectii :many
SELECT id, denumire FROM sectii WHERE sters_la IS NULL AND activ = 1 ORDER BY denumire;

-- name: AdaugaSectie :exec
INSERT INTO sectii (id, denumire, creat_la, modificat_la) VALUES (?, ?, ?, ?);

-- name: CitesteSectieActiva :one
SELECT id FROM sectii WHERE id = ? AND sters_la IS NULL AND activ = 1;

-- name: ListaDepartamente :many
SELECT id, sectie_id, denumire FROM departamente WHERE sters_la IS NULL AND activ = 1
ORDER BY sectie_id, denumire;

-- name: CitesteDepartamentActiv :one
SELECT id, sectie_id FROM departamente WHERE id = ? AND sters_la IS NULL AND activ = 1;

-- name: AdaugaDepartament :exec
INSERT INTO departamente (id, sectie_id, denumire, creat_la, modificat_la) VALUES (?, ?, ?, ?, ?);

-- name: CitesteSetare :one
SELECT valoare FROM setari WHERE cheie = ?;

-- name: ActualizeazaSetare :execrows
UPDATE setari SET valoare = ?, modificat_la = ? WHERE cheie = ?;

-- name: NumaraAngajati :one
SELECT count(*) FROM angajati;

-- name: NumaraPontaje :one
SELECT count(*) FROM pontaj;

-- name: AdaugaJurnal :exec
INSERT INTO jurnal_modificari (tabel, rand_id, valoare_veche, valoare_noua, cand, cine)
VALUES (?, ?, ?, ?, ?, ?);
