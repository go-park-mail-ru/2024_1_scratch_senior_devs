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

CREATE TABLE IF NOT EXISTS attaches (
    id UUID PRIMARY KEY,
    path TEXT
        NOT NULL
        CONSTRAINT path_length CHECK (char_length(path) <= 255),
    note_id UUID
        NOT NULL
);

INSERT INTO users (id, description, username, password_hash, create_time)
    VALUES ('57247dd2-b768-4665-a290-8dad9506616a', '', 'mizhgun', 'f969248d621bcded4a3582a1c3b17a71eedfefa9120c36ee3bd1957438cd55b9', CURRENT_TIMESTAMP);
