CREATE TABLE counters (
    id SERIAL PRIMARY KEY,
    value INTEGER DEFAULT 0
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE syncs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    last_sync REAL NOT NULL,
    is_synced BOOLEAN NOT NULL,
    status TEXT NOT NULL
);
