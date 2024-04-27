CREATE TABLE surveys (
    id UUID         PRIMARY KEY,
    created_at      TIMESTAMP
                    NOT NULL
);

CREATE TABLE questions (
    id              UUID    PRIMARY KEY,
    title           TEXT
                    CONSTRAINT title_length CHECK (char_length(title) <= 255),
--     min_mark        INT,
--     skip            INT,
    question_type   TEXT
                    CONSTRAINT question_type_length CHECK (char_length(question_type) <= 255),
    number          INT
                    NOT NULL,
    survey_id UUID REFERENCES surveys (id) ON DELETE CASCADE
                    NOT NULL
);

CREATE TABLE results (
    id              UUID PRIMARY KEY,
    question_id     UUID REFERENCES questions (id) ON DELETE CASCADE
                    NOT NULL,
    voice           INT
                    NOT NULL
);

--- INSERT INTO surveys (id, created_at) VALUES ('855cb60a-9ed0-485b-bd3c-6ac29f429ebf', CURRENT_TIMESTAMP);
--- INSERT INTO questions ( id, title, question_type, number, survey_id ) VALUES ('611f9ded-1340-4b3d-9202-ce7aec7fe418','first questiojn', 'CSAT', 1, '855cb60a-9ed0-485b-bd3c-6ac29f429ebf');