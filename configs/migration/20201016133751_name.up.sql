-- DDL generated by Postico 1.5.6
-- Not all database features are supported. Do not use for backup.

-- Table Definition ----------------------------------------------

CREATE TABLE IF NOT EXISTS protein_event (
    user_id character varying(9) PRIMARY KEY,
    utc_time_to_drink timestamp without time zone,
    drink_time_interval_min integer NOT NULL DEFAULT 0
);

