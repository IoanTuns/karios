-- F09-T08, D026: imaginile rămân în cele trei zone ale fiecărei definiții.
-- Textele, reviziile și momentele editărilor existente sunt păstrate.
-- +goose Up
UPDATE definitii_rapoarte
SET antet = CASE WHEN json_type(antet, '$.imagini') IS NULL
                 THEN json_set(antet, '$.imagini', json('{}')) ELSE antet END,
    subsol = CASE WHEN json_type(subsol, '$.imagini') IS NULL
                  THEN json_set(subsol, '$.imagini', json('{}')) ELSE subsol END
WHERE json_type(antet, '$.imagini') IS NULL OR json_type(subsol, '$.imagini') IS NULL;
