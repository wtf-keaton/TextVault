-- +goose Up
CREATE TABLE Users (
    ID SERIAL PRIMARY KEY,
    Username VARCHAR(50) NOT NULL,
    Email VARCHAR(100) UNIQUE NOT NULL,
    PasswordHash VARCHAR(255) NOT NULL,
    IsAdmin BOOLEAN DEFAULT FALSE,
    IsBanned BOOLEAN DEFAULT FALSE,
    
    UNIQUE (Username, Email)
);

CREATE INDEX idx_user_username ON Users (Username);
CREATE INDEX idx_user_email ON Users (Email);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Pastes (
    ID UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Language VARCHAR(50) NOT NULL,
    AuthorID SERIAL -- Can be NULL
);

CREATE INDEX idx_paste_title ON Pastes (Title);
CREATE INDEX idx_paste_author_id ON Pastes (AuthorID);

-- +goose Down
DROP INDEX IF EXISTS idx_paste_author_id;
DROP INDEX IF EXISTS idx_paste_title;
DROP TABLE IF EXISTS Pastes;

DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_user_username;
DROP TABLE IF EXISTS Users;