-- F07-T35: clasificari editabile cu semantica explicita.
-- +goose Up
CREATE TABLE tipuri_program (
 id TEXT PRIMARY KEY,
 denumire TEXT NOT NULL,
 comportament TEXT NOT NULL DEFAULT 'NORMAL' CHECK(comportament IN ('NORMAL','TURA','NOAPTE','GARDA')),
 activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN(0,1)),
 modificat_la TEXT NOT NULL
) STRICT;
CREATE TABLE regimuri_program (
 id TEXT PRIMARY KEY,
 denumire TEXT NOT NULL,
 comportament TEXT NOT NULL DEFAULT 'LUCRATOARE' CHECK(comportament IN ('LUCRATOARE','NELUCRATOARE','INTRARE_ZI_LIBERA')),
 activ INTEGER NOT NULL DEFAULT 1 CHECK(activ IN(0,1)),
 modificat_la TEXT NOT NULL
) STRICT;
INSERT INTO tipuri_program VALUES
 ('NORMAL','{"ro":"Normal","en":"Normal"}','NORMAL',1,'initial'),
 ('TURA','{"ro":"Tura","en":"Shift"}','TURA',1,'initial'),
 ('NOAPTE','{"ro":"Noapte","en":"Night"}','NOAPTE',1,'initial'),
 ('GARDA','{"ro":"Garda","en":"On-call"}','GARDA',1,'initial');
INSERT INTO regimuri_program VALUES
 ('LUCRATOARE','{"ro":"Normal","en":"Normal"}','LUCRATOARE',1,'initial'),
 ('NELUCRATOARE','{"ro":"Weekend / sarbatoare","en":"Weekend / holiday"}','NELUCRATOARE',1,'initial'),
 ('INTRARE_ZI_LIBERA','{"ro":"Intrare speciala in zi libera","en":"Special start on a day off"}','INTRARE_ZI_LIBERA',1,'initial');
UPDATE tipuri_schimb SET regim='LUCRATOARE' WHERE regim='GENERAL';
CREATE TABLE clasificari_schimburi (
 schimb_id TEXT PRIMARY KEY REFERENCES tipuri_schimb(id) ON DELETE CASCADE,
 tip_id TEXT NOT NULL REFERENCES tipuri_program(id),
 regim_id TEXT NOT NULL REFERENCES regimuri_program(id)
) STRICT;
INSERT INTO clasificari_schimburi SELECT id,tip,regim FROM tipuri_schimb;
CREATE INDEX clasificari_tip ON clasificari_schimburi(tip_id);
CREATE INDEX clasificari_regim ON clasificari_schimburi(regim_id);
-- +goose StatementBegin
CREATE TRIGGER regim_general_insert AFTER INSERT ON tipuri_schimb WHEN NEW.regim='GENERAL'
BEGIN
 UPDATE tipuri_schimb SET regim='LUCRATOARE' WHERE id=NEW.id;
END;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER regim_general_update AFTER UPDATE OF regim ON tipuri_schimb WHEN NEW.regim='GENERAL'
BEGIN
 UPDATE tipuri_schimb SET regim='LUCRATOARE' WHERE id=NEW.id;
END;
-- +goose StatementEnd
-- +goose Down
DROP TRIGGER regim_general_update;
DROP TRIGGER regim_general_insert;
DROP TABLE clasificari_schimburi;
DROP TABLE regimuri_program;
DROP TABLE tipuri_program;
