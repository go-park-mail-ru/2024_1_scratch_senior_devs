-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS note (
    id          UUID                        PRIMARY KEY,
    note_data   TEXT
                CONSTRAINT data_length
                CHECK (char_length(data) <= 4000),
    create_time TIMESTAMP WITH TIME ZONE    NOT NULL,
    update_time TIMESTAMP WITH TIME ZONE    NOT NULL,
    owner_id    UUID                        NOT NULL
                REFERENCES profile (id)
                ON DELETE CASCADE
);

---- create above / drop below ----

DROP TABLE note CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
