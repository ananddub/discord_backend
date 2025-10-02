-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    profile_pic VARCHAR(255),
    bio VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    color_code VARCHAR(255) DEFAULT '#FFFFFF',
    background_color VARCHAR(255) DEFAULT 'white',
    background_pic VARCHAR(255) DEFAULT NULL,
    status VARCHAR(50) DEFAULT 'offline' NOT NULL
);

-- +goose Down
DROP TABLE users;
