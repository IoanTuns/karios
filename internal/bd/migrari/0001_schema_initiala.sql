-- Schema inițială (F06-T02). Doar tabelele de bază, ale căror cerințe sunt clare.
-- Planificarea (F11), închiderea lunii și segmentele înghețate (F12), coloanele
-- raportului (F07, F09), conturile (F15) și procentele versionate (F04-T03, F04-T06)
-- vin în migrări proprii, odată cu taskurile lor.
--
-- Convenții:
--   * Identificatorii sunt UUID (TEXT), ca datele să poată fi sincronizate (F13).
--   * Momentele de timp sunt TEXT în UTC, format ISO 8601 („2026-10-25T05:00:00Z”).
--   * Denumirile afișate sunt JSON cu câte o variantă pe limbă (i18n.Text, D004).
--   * Coloanele de sincronizare: creat_la, modificat_la, sters_la (ștergere logică),
--     proprietar_id, id_remote, stare_sinc.

-- +goose Up

-- Setările aplicației (D008: stau în baza criptată, nu în config.toml).
CREATE TABLE setari (
    cheie        TEXT PRIMARY KEY,
    valoare      TEXT NOT NULL,
    modificat_la TEXT NOT NULL
) STRICT;

CREATE TABLE sectii (
    id            TEXT PRIMARY KEY,
    denumire      TEXT NOT NULL CHECK (json_valid(denumire)),
    activ         INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL,
    sters_la      TEXT,
    proprietar_id TEXT,
    id_remote     TEXT,
    stare_sinc    TEXT NOT NULL DEFAULT 'local'
) STRICT;

CREATE TABLE angajati (
    id            TEXT PRIMARY KEY,
    nume          TEXT NOT NULL,
    marca         TEXT NOT NULL,
    functie       TEXT,
    sectie_id     TEXT REFERENCES sectii (id),
    norma_minute  INTEGER NOT NULL DEFAULT 480 CHECK (norma_minute > 0),
    activ         INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL,
    sters_la      TEXT,
    proprietar_id TEXT,
    id_remote     TEXT,
    stare_sinc    TEXT NOT NULL DEFAULT 'local'
) STRICT;

-- Marca e unică doar printre angajații neșterși.
CREATE UNIQUE INDEX angajati_marca ON angajati (marca) WHERE sters_la IS NULL;

-- Tipurile de schimb și codurile lor scurte din grilă (7, 12, N12, G24).
CREATE TABLE tipuri_schimb (
    id             TEXT PRIMARY KEY,
    cod            TEXT NOT NULL,
    denumire       TEXT NOT NULL CHECK (json_valid(denumire)),
    ora_start      INTEGER NOT NULL CHECK (ora_start BETWEEN 0 AND 1439),     -- minute de la 00:00
    durata_minute  INTEGER NOT NULL CHECK (durata_minute BETWEEN 1 AND 2880),
    pauza_minute   INTEGER NOT NULL DEFAULT 0 CHECK (pauza_minute >= 0),
    tip            TEXT NOT NULL CHECK (tip IN ('NORMAL', 'TURA', 'NOAPTE', 'GARDA')),
    activ          INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    creat_la       TEXT NOT NULL,
    modificat_la   TEXT NOT NULL,
    sters_la       TEXT,
    proprietar_id  TEXT,
    id_remote      TEXT,
    stare_sinc     TEXT NOT NULL DEFAULT 'local'
) STRICT;

CREATE UNIQUE INDEX tipuri_schimb_cod ON tipuri_schimb (cod) WHERE sters_la IS NULL;

CREATE TABLE coduri_absenta (
    cod           TEXT PRIMARY KEY,
    denumire      TEXT NOT NULL CHECK (json_valid(denumire)),
    activ         INTEGER NOT NULL DEFAULT 1 CHECK (activ IN (0, 1)),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL
) STRICT;

CREATE TABLE sarbatori (
    data          TEXT PRIMARY KEY CHECK (data IS date(data)),               -- AAAA-LL-ZZ
    denumire      TEXT NOT NULL CHECK (json_valid(denumire)),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL
) STRICT;

-- Categoriile de ore (ZI, NOAPTE, WEEKEND, SARBATOARE…). Procentele nu stau aici:
-- trebuie versionate în timp (F04-T06) și depind de modelul de sporuri (F04-T03).
CREATE TABLE categorii_ore (
    cod           TEXT PRIMARY KEY,
    denumire      TEXT NOT NULL CHECK (json_valid(denumire)),
    creat_la      TEXT NOT NULL,
    modificat_la  TEXT NOT NULL
) STRICT;

-- Regulile de tarif, versionate în timp (D003).
CREATE TABLE seturi_reguli (
    id              TEXT PRIMARY KEY,
    denumire        TEXT NOT NULL,
    valabil_de_la   TEXT NOT NULL CHECK (valabil_de_la IS date(valabil_de_la)),
    valabil_pana_la TEXT CHECK (valabil_pana_la IS NULL OR valabil_pana_la IS date(valabil_pana_la)),
    creat_la        TEXT NOT NULL,
    modificat_la    TEXT NOT NULL
) STRICT;

-- Intervalele acoperă exact 00:00–24:00 pe fiecare tip de zi; verificarea o face
-- calcul.SetReguli.Valideaza la salvare.
CREATE TABLE reguli_interval (
    id            TEXT PRIMARY KEY,
    set_id        TEXT NOT NULL REFERENCES seturi_reguli (id) ON DELETE CASCADE,
    tip_zi        TEXT NOT NULL CHECK (tip_zi IN ('LUCRATOARE', 'SAMBATA', 'DUMINICA', 'SARBATOARE')),
    ora_start     INTEGER NOT NULL CHECK (ora_start BETWEEN 0 AND 1439),
    ora_sfarsit   INTEGER NOT NULL CHECK (ora_sfarsit BETWEEN 1 AND 1440 AND ora_sfarsit > ora_start),
    categorie_cod TEXT NOT NULL REFERENCES categorii_ore (cod)
) STRICT;

CREATE INDEX reguli_interval_set ON reguli_interval (set_id, tip_zi, ora_start);

-- Ce s-a lucrat efectiv. O tură de noapte e un singur rând (D003); absențele au
-- cod_absenta. Modelul absențelor se poate schimba după F06-T05.
CREATE TABLE pontaj (
    id              TEXT PRIMARY KEY,
    angajat_id      TEXT NOT NULL REFERENCES angajati (id),
    data            TEXT NOT NULL CHECK (data IS date(data)),                -- ziua în care începe tura
    start           TEXT,                                                   -- moment UTC
    sfarsit         TEXT,                                                   -- moment UTC
    pauza_minute    INTEGER NOT NULL DEFAULT 0 CHECK (pauza_minute >= 0),
    tip_schimb_id   TEXT REFERENCES tipuri_schimb (id),
    cod_absenta     TEXT REFERENCES coduri_absenta (cod),
    observatii      TEXT,
    creat_la        TEXT NOT NULL,
    modificat_la    TEXT NOT NULL,
    sters_la        TEXT,
    proprietar_id   TEXT,
    id_remote       TEXT,
    stare_sinc      TEXT NOT NULL DEFAULT 'local',
    CHECK ((start IS NULL) = (sfarsit IS NULL)),
    CHECK (start IS NULL OR sfarsit > start),
    CHECK (start IS NOT NULL OR cod_absenta IS NOT NULL)
) STRICT;

CREATE INDEX pontaj_angajat_data ON pontaj (angajat_id, data) WHERE sters_la IS NULL;

-- Jurnalul modificărilor (F06-T04): tabelul există de la început; scrierea în el
-- vine cu F06-T04, iar „cine” devine contul din aplicație cu F15.
CREATE TABLE jurnal_modificari (
    id            INTEGER PRIMARY KEY,
    tabel         TEXT NOT NULL,
    rand_id       TEXT NOT NULL,
    camp          TEXT,
    valoare_veche TEXT,
    valoare_noua  TEXT,
    cand          TEXT NOT NULL,
    cine          TEXT NOT NULL
) STRICT;

CREATE INDEX jurnal_modificari_rand ON jurnal_modificari (tabel, rand_id);
