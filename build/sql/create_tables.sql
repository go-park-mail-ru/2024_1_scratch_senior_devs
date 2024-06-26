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
    id              UUID        PRIMARY KEY,
    data            JSON,
    create_time     TIMESTAMP   NOT NULL,
    update_time     TIMESTAMP   NOT NULL,
    owner_id        UUID        NOT NULL
                    REFERENCES users (id),
    parent          UUID        DEFAULT ('00000000-0000-0000-0000-000000000000'::UUID),
    children        UUID[],
    tags            TEXT[],
    collaborators   UUID[],
    icon            TEXT
                    CONSTRAINT icon_length CHECK (char_length(icon) <= 255),
    header          TEXT
                    CONSTRAINT header_length CHECK (char_length(header) <= 255),
    is_public       BOOLEAN     NOT NULL
                    DEFAULT false
);

CREATE TABLE IF NOT EXISTS favorites (
    note_id UUID REFERENCES notes (id),
    user_id UUID REFERENCES users (id),
    PRIMARY KEY(note_id, user_id)
);

CREATE TABLE IF NOT EXISTS all_tags (
    tag_name    TEXT
                CONSTRAINT tag_name_length CHECK (char_length(tag_name) <= 255),
    user_id     UUID
                REFERENCES users (id),
    PRIMARY KEY (tag_name, user_id)
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

CREATE OR REPLACE FUNCTION update_tags()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    AS $BODY$
    BEGIN
        UPDATE notes SET tags = array_replace(tags, OLD.tag_name, NEW.tag_name) WHERE owner_id=OLD.user_id;
        RETURN OLD;
    END;
$BODY$;

CREATE OR REPLACE TRIGGER trigger_update_tags
    AFTER UPDATE
    ON all_tags
    FOR EACH ROW
    EXECUTE FUNCTION  update_tags();
