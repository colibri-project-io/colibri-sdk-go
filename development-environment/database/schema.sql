CREATE TABLE IF NOT EXISTS profiles (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    birthday TIMESTAMP NOT NULL,
    profile_id INTEGER NOT NULL,
    CONSTRAINT fk_users_profiles FOREIGN KEY (profile_id) REFERENCES profiles (id)
);

CREATE TABLE IF NOT EXISTS contacts (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    email       TEXT UNIQUE
);