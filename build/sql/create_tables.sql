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
        CONSTRAINT image_path_length CHECK (char_length(image_path) <= 255)
);

CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY,
    data JSON,
    create_time TIMESTAMP
        NOT NULL,
    update_time TIMESTAMP,
    owner_id UUID REFERENCES users (id)
        NOT NULL
);


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION add_draft_note()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
DECLARE
    name_text json;
BEGIN
    IF NEW.id IS NOT NULL AND NEW.username IS NOT NULL THEN
        name_text := format('{
            "title": "You-note <3",
            "content": "Привет, %s!"
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


INSERT INTO users (id, description, username, password_hash, create_time)
    VALUES ('3d67ba58-9023-42b5-a059-d626d7587f1e', 'У меня много заметок!', 'mizhgun', 'f969248d621bcded4a3582a1c3b17a71eedfefa9120c36ee3bd1957438cd55b9', CURRENT_TIMESTAMP);
INSERT INTO notes (id, data, create_time, update_time, owner_id)
    VALUES ('f732e6ae-07fe-4f20-aca1-37e8f115d082', '{
            "title": "Моя первая заметка",
            "content": "Очень важный текст, который определенно стоит запомнить"
        }', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '3d67ba58-9023-42b5-a059-d626d7587f1e'),
            ('81e769de-4cd3-499b-90a7-98b4aca6f9d2', '{
            "title": "Моя вторая заметка",
            "content": "Текст поважнее"
        }', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '3d67ba58-9023-42b5-a059-d626d7587f1e'),
            ('6e246734-e712-4182-ad59-c5cd8debff8b', '{
            "title": "Моя третья заметка",
            "content": "Не такой важный текст, но тоже хочется не забыть"
        }', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '3d67ba58-9023-42b5-a059-d626d7587f1e');

INSERT INTO users (id, description, username, password_hash, create_time)
    VALUES ('57247dd2-b768-4665-a290-8dad9506616a', 'Заметок не делаю, всё держу в голове!', 'alladan', 'f969248d621bcded4a3582a1c3b17a71eedfefa9120c36ee3bd1957438cd55b9', CURRENT_TIMESTAMP);
DELETE FROM notes WHERE owner_id = '57247dd2-b768-4665-a290-8dad9506616a'
