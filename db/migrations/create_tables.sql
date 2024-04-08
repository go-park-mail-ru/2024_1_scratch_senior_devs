CREATE TABLE IF NOT EXISTS profile (
    id              UUID        PRIMARY KEY,
    description     TEXT
                    CONSTRAINT description_length
                    CHECK (char_length(username) <= 40),
    username        TEXT        NOT NULL UNIQUE
                    CONSTRAINT name_length
                    CHECK (char_length(username) <= 20),
    password_hash   TEXT        NOT NULL
                    CONSTRAINT password_hash_length
                    CHECK (char_length(password_hash) <= 300),
    create_time     TIMESTAMP   NOT NULL,
    image_path      TEXT        NOT NULL
                    DEFAULT ('default.jpg')
                    CONSTRAINT image_path_length
                    CHECK (char_length(image_path) <= 40),
    secret          TEXT
                    CONSTRAINT secret_length
                    CHECK (char_length(secret) <= 40)
);

CREATE TABLE IF NOT EXISTS note (
    id          UUID        PRIMARY KEY,
    data        TEXT
                CONSTRAINT data_length
                CHECK (char_length(data) <= 4000),
    create_time TIMESTAMP   NOT NULL,
    update_time TIMESTAMP   NOT NULL,
    owner_id    UUID        NOT NULL
                REFERENCES profile (id)
                ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS attach (
    id      UUID    PRIMARY KEY,
    path    TEXT    NOT NULL UNIQUE
            CONSTRAINT path_length
            CHECK (char_length(path) <= 40),
    note_id UUID
            REFERENCES note (id)
            ON DELETE SET NULL
);