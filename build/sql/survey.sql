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
    number          INT,
    survey_id UUID REFERENCES surveys (id) ON DELETE CASCADE
                    NOT NULL
);

CREATE TABLE results (
    id              UUID PRIMARY KEY,
    question_id     UUID
                    NOT NULL,
    voice           INT
                    NOT NULL
);
