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
        NOT NULL,
    parent UUID DEFAULT ('00000000-0000-0000-0000-000000000000'::UUID),
    children UUID[]
);


CREATE TABLE IF NOT EXISTS note_tag (
    note_id     UUID REFERENCES notes (id) ON DELETE CASCADE,
    tag_name    TEXT NOT NULL
                 CONSTRAINT tag_length CHECK (char_length(tag_name) <= 255),
    PRIMARY KEY (note_id, tag_name)
);

CREATE TABLE IF NOT EXISTS attaches (
    id UUID PRIMARY KEY,
    path TEXT
        NOT NULL
        CONSTRAINT path_length CHECK (char_length(path) <= 255),
    note_id UUID REFERENCES notes (id) ON DELETE CASCADE
        NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    note_id         UUID        NOT NULL,
    created         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    message_info    JSON
);

CREATE TABLE IF NOT EXISTS collaborators (
    note_id UUID REFERENCES notes (id) ON DELETE CASCADE,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (note_id, user_id)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION delete_children_notes()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    AS $BODY$
    BEGIN
        DELETE FROM notes WHERE id = ANY(SELECT UNNEST(OLD.children));
        RETURN OLD;
    END;
$BODY$;

CREATE OR REPLACE TRIGGER trigger_delete_children_notes
    BEFORE DELETE
    ON notes
    FOR EACH ROW
    EXECUTE FUNCTION delete_children_notes();

CREATE OR REPLACE FUNCTION insert_message()
    RETURNS trigger
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
    INSERT INTO messages(note_id, created, message_info) VALUES (NEW.id, CURRENT_TIMESTAMP, NEW.data);
    RETURN NEW;
END;
$BODY$;

CREATE OR REPLACE TRIGGER trigger_insert_message
    AFTER UPDATE
    ON notes
    FOR EACH ROW
EXECUTE FUNCTION insert_message();

