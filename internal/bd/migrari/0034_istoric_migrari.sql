-- F06-T06 (D009 punctul 3): ce migrari s-au aplicat, cand si din ce versiune a aplicatiei.
-- goose_db_version tine migrarea si momentul; aici se adauga versiunea aplicatiei.
-- +goose Up
CREATE TABLE istoric_migrari (
 versiune INTEGER NOT NULL,
 aplicatie TEXT NOT NULL,
 aplicata_la TEXT NOT NULL
);
