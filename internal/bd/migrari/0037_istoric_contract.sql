-- F07-T43 (I010, D101): schimbarile normei zilnice si ale dreptului anual la concediu, cu data
-- de la care sunt valabile. Valorile din fisa angajatului raman cele de dinaintea primei schimbari.
-- Un camp NULL inseamna ca valoarea respectiva nu se schimba la acea data.
-- +goose Up
CREATE TABLE istoric_contract (
    angajat_id       TEXT NOT NULL REFERENCES angajati (id),
    de_la            TEXT NOT NULL CHECK (de_la GLOB '[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]' AND date(de_la) = de_la),
    norma_minute     INTEGER CHECK (norma_minute IS NULL OR norma_minute BETWEEN 1 AND 1440),
    concediu_zile_an INTEGER CHECK (concediu_zile_an IS NULL OR concediu_zile_an BETWEEN 0 AND 366),
    creat_la         TEXT NOT NULL,
    CHECK (norma_minute IS NOT NULL OR concediu_zile_an IS NOT NULL),
    PRIMARY KEY (angajat_id, de_la)
);

-- +goose Down
DROP TABLE istoric_contract;
