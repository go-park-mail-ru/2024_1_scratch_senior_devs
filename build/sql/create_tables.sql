CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    description TEXT
        CONSTRAINT description_length CHECK (char_length(username) <= 255),
    username TEXT
        NOT NULL
        UNIQUE
        CONSTRAINT name_length CHECK (char_length(username) <= 255),
    password_hash TEXT
        NOT NULL
        CONSTRAINT password_hash_length CHECK (char_length(password_hash) <= 511),
    create_time TIMESTAMP
        NOT NULL,
    image_path TEXT DEFAULT ('default.jpg')
        NOT NULL
        CONSTRAINT image_path_length CHECK (char_length(image_path) <= 255),
    secret TEXT
        CONSTRAINT secret_length CHECK (char_length(secret) <= 255)
);

CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY,
    data JSON,
    create_time TIMESTAMP
        NOT NULL,
    update_time TIMESTAMP
        NOT NULL,
    owner_id UUID REFERENCES users (id)
        NOT NULL
);

CREATE TABLE IF NOT EXISTS attaches (
    id UUID PRIMARY KEY,
    path TEXT
        NOT NULL
        CONSTRAINT path_length CHECK (char_length(path) <= 255),
    note_id UUID REFERENCES notes (id)
        NOT NULL
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION add_draft_note()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    AS $BODY$
DECLARE
    name_text json;
BEGIN
    IF NEW.id IS NOT NULL AND NEW.username IS NOT NULL THEN
        name_text := format('{
			"title": "You-note ❤️",
			"content": [
				{
					"id": "1",
					"type": "div",
					"content": [
						{
							"id": "2",
							"content": "Привет, %s!"
						}
					]
				}
			]
		}', NEW.username);
        INSERT INTO notes (id, data, create_time, update_time, owner_id)
        VALUES (uuid_generate_v4(), name_text, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NEW.id);
    END IF;

    RETURN NEW;
END;
$BODY$;

CREATE OR REPLACE TRIGGER add_note_on_new_user
    AFTER INSERT
    ON users
    FOR EACH ROW
    EXECUTE FUNCTION add_draft_note();
