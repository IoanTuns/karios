-- F07-T06, D012: sărbătorile se definesc și se citesc separat pentru fiecare an.
-- Anul derivă din data completă, fără o a doua valoare editabilă care să o contrazică.
-- +goose Up

ALTER TABLE sarbatori ADD COLUMN an INTEGER
    GENERATED ALWAYS AS (CAST(substr(data, 1, 4) AS INTEGER)) VIRTUAL;

CREATE INDEX sarbatori_an_data ON sarbatori (an, data);
