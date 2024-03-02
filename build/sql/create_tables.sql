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