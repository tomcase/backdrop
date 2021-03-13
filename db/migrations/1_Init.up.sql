BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS settings (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	name VARCHAR(50) NOT NULL UNIQUE,
	value VARCHAR(50) NOT NULL DEFAULT '',
	PRIMARY KEY (id)
);

INSERT INTO settings (name, value)
VALUES ('DEVICE_ID', '');

INSERT INTO settings (name, value)
VALUES ('CLIENT_ID', 'VIz3jyOacujEuQ');

INSERT INTO settings (name, value)
VALUES ('OUTPUT_DIR', '/Users/tom/Pictures/Backgrounds');

COMMIT;

