CREATE TABLE gists (
    id VARCHAR(100) PRIMARY KEY,

    hash TEXT NOT NULL,
    result TEXT NOT NULL,
    language VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (language, hash)
);
