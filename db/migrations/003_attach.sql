-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS attach (
    id          UUID    PRIMARY KEY,
    file_path   TEXT    NOT NULL UNIQUE
                CONSTRAINT path_length
                CHECK (char_length(path) <= 40),
    note_id     UUID
                REFERENCES note (id)
                ON DELETE SET NULL
);

---- create above / drop below ----

DROP TABLE attach CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
