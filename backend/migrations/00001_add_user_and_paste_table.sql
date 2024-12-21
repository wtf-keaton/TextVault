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

CREATE TABLE Pastes (
    ID SERIAL PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Hash VARCHAR(64) UNIQUE NOT NULL,
    AuthorID BIGINT NOT NULL,
    FOREIGN KEY (AuthorID) REFERENCES Users(ID)
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