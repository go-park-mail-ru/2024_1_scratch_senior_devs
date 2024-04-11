-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS profile (
    id              UUID                        PRIMARY KEY,
    description     TEXT
                    CONSTRAINT description_length
                    CHECK (char_length(username) <= 40),
    username        TEXT                        NOT NULL UNIQUE
                    CONSTRAINT name_length
                    CHECK (char_length(username) <= 20),
    password_hash   TEXT                        NOT NULL
                    CONSTRAINT password_hash_length
                    CHECK (char_length(password_hash) <= 300),
    create_time     TIMESTAMP WITH TIME ZONE    NOT NULL,
    image_path      TEXT                        NOT NULL
                    DEFAULT ('default.jpg')
                    CONSTRAINT image_path_length
                    CHECK (char_length(image_path) <= 40),
    secret          TEXT
                    CONSTRAINT secret_length
                    CHECK (char_length(secret) <= 40)
);

---- create above / drop below ----

DROP TABLE profile CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
